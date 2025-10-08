# **架構設計文件 (ARCHITECTURE.md)**

## **1\. 核心控制迴圈 (Core Control Loop)**

本專案的核心是 O-RAN Hardware Adapter。它遵循 Kubernetes Operator 模式，其主要職責是**持續地將 Kubernetes API 中 OpticalDevice 資源的期望狀態 (Spec)，同步到外部實體設備的實際狀態上**。

### **Reconcile 流程圖**

graph TD  
    A(Start Reconcile for OpticalDevice 'X') \--\> B{Get OpticalDevice 'X' from K8s};  
    B \--\> C{Is 'X' being deleted?};  
    C \-- Yes \--\> D\[Perform Cleanup on Hardware\];  
    D \--\> E{Cleanup Successful?};  
    E \-- Yes \--\> F\[Remove Finalizer from 'X'\];  
    E \-- No \--\> G\[Requeue with Error\];  
    F \--\> H(End);  
    C \-- No \--\> I\[Ensure Finalizer exists on 'X'\];  
    I \--\> J\[Read Spec from 'X'\];  
    J \--\> K\[Generate Hardware Config Payload\];  
    K \--\> L{Call Hardware Simulator API};  
    L \-- Success \--\> M\[Read back actual state from Simulator\];  
    M \--\> N\[Update Status of 'X' with actual state & 'Ready' Condition\];  
    N \--\> H;  
    L \-- Failure \--\> O\[Update Status with 'Error' Condition\];  
    O \--\> G;

### **關鍵設計：**

* **Finalizer (最終處理器)**: 我們將為每一個 OpticalDevice 資源添加一個 finalizer (例如 hardware.ran.example.com/finalizer)。當使用者刪除 OpticalDevice 時，K8s 不會立即移除它，而是設置一個 deletionTimestamp。我們的控制器會捕捉到這個狀態，執行對應的硬體資源釋放或關閉指令。操作成功後，再移除 finalizer，K8s 才會真正刪除該資源。這確保了不會產生孤兒硬體資源。  
* **冪等性 (Idempotency)**: 控制器的每一次 Reconcile 都應該是冪等的。無論執行一次還是十次，只要期望狀態不變，硬體的最終狀態都應該相同。我們會透過比較 Spec 與 Status 中的 observedGeneration 或設定雜湊值來避免不必要的硬體 API 呼叫。  
* **狀態回報 (Status Reporting)**: OpticalDevice 的 Status 欄位至關重要。它不僅僅是一個 Ready 狀態，還應該包含從硬體回讀的實際狀態（observedState）和詳細的 Conditions，以便使用者和上層系統（如 SMO）能清晰地了解設備的當前狀況。

## **2\. 元件互動與 API 介面**

* **Adapter \<-\> Hardware Simulator**:  
  * **介面**: RESTful HTTP API  
  * **Endpoint**: POST /configure  
  * **Request Body**: {"hostname": "...", "port": ..., "bandwidth": "...", "laserPower": "..."}  
  * **Response Body**: {"status": "configured", "observedBandwidth": "...", "lastUpdated": "..."}  
  * **Endpoint**: DELETE /configure  
  * **Request Body**: {"hostname": "..."}  
  * **Response Body**: {"status": "deconfigured"}

這個明確的介面定義使得 Adapter 和模擬器可以獨立開發與測試。