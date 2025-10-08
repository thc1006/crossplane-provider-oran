# **O-RAN 光通訊設備管理與協調專案 \- 會議逐字稿 (修正與整合版)**

## **專案目標與核心挑戰**

我們的核心目標是建立一個自動化的 O-RAN 光通訊設備管理與協調架構。目前遇到的主要挑戰是，如何將傳統的光通訊元件（例如雷射、ONU 等 PNF, Physical Network Function）有效地納入到現代化的、以 Kubernetes (K8s) 為中心的雲原生管理平台中。

過去，許多設定是透過手動上傳設定檔或直接下 command line 指令完成的。現在的目標是從上層的 SMO (Service Management and Orchestration)，透過 O1 管理介面，將整個流程自動化串接起來，實現真正的 Infrastructure as Code (IAC)。

## **架構設計與關鍵技術**

我們將沿用並擴展現有的架構，其基本流程如下：

1. **SMO 層**: 作為最高層級的管理者，發出意圖（Intent）。  
2. **Orchestrator (以 Nephio 為例)**: 接收來自 SMO 的意圖，並將其轉化為 K8s 世界可以理解的部署管理指令。  
3. **Crossplane**: 這是我們架構的核心。我們將使用 Crossplane 來作為 IAC 的自動化平台。  
   * **Composite Resource (XR)**: 我們會定義高階的 XR 來抽象化光通訊服務，例如 "光纖網路服務"。  
   * **Composition**: XR 會透過 Composition 來組合多個底層的 Managed Resources。  
   * **XRD (Composite Resource Definition)**: 用於定義 XR 的 Schema。  
   * **Claim**: 開發者或管理者可以透過提交一個簡單的 Claim 來請求一個 XR 實例。  
4. **Config Sync**: 用於同步 Git Repo 中的 YAML 設定檔到 K8s 叢集中，確保 GitOps 流程。  
5. **Custom Operator/Adapter**: 這是需要我們重點開發的部分。Crossplane 負責管理 K8s 內的資源，但要控制 K8s 外部的實體硬體（如光通訊設備），我們需要一個 Operator/Adapter。這個 Adapter 會監聽 K8s 中的 Custom Resource (CR)，並將其狀態轉化為對硬體設備的具體控制指令（例如透過 Netconf, gNMI 或傳統 API）。我們將其命名為 O-RAN Hardware Adapter。  
6. **被管理的設備**:  
   * **雲原生網元 (CNF)**: 如 CU-CP, CU-UP 等，這些本身就可以在 K8s 中進行部署與管理。  
   * **實體網元 (PNF)**: 如 AGF (Access Gateway Function)、DPU、DPDK 網卡，以及我們的重點——光通訊元件。這些需要透過 Adapter 來進行橋接控制。

## **開發與呈現**

* **系統雛形對接**: 今年的工作項目重點是完成這個系統雛形，並實現端到端的對接。  
* **視覺化**: 為了 Demo 的效果，我們會使用 **Grafana** 來接上相關數據，將資源管理、頻寬控制（例如 TRTCM 參數設定）等狀態以儀表板的形式美觀地呈現出來。  
* **圖表呈現**: 在簡報和文件中，Local Controller 的圖示應該清晰地表達其軟硬體關係，例如 Controller 的方框疊加在代表實體主機的電腦圖示上，以節省空間並清晰傳達概念。Adapter 應該被包含在 Controller 的邏輯框圖內。

## **結論**

專案方向明確：**利用 Crossplane 作為核心抽象層，開發一個客製化的 Adapter 來橋接雲原生控制平面與實體光通訊設備，並透過 SMO 實現端到端的自動化管理與部署，最終搭配 Grafana 進行視覺化呈現。**