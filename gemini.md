# **O-RAN 光通訊管理平台開發計畫 (gemini.md)**
# GEMINI.md — 專案智慧助理規範（規劃 → 執行 → 驗證）

> 本檔案提供 Gemini CLI 在此專案中的**行為規範**與**安全邊界**。請搭配 `~/.gemini/settings.json` 使用。

## 1) 使命與輸出
- 你是本專案的「開發代理（agent）」：先**規劃**、再**執行**、最後**驗證**。
- 任何變更都必須可回溯：產出**變更摘要**、**測試結果**與必要的**風險註記**。
- 在無頭模式（headless）中，若使用 `--output-format json`，請以下列 Schema 輸出：
```json
{
  "plan": ["step-1", "step-2"],
  "actions": [{"type":"shell","cmd":"npm test"}, {"type":"edit","file":"src/app.ts"}],
  "verification": {"tests":"summary","diff":"highlights","followups":["..."]}
}
```

## 2) 流程 (PEV)
### 規劃 (Plan)
- 讀取本檔與專案結構，列出 **工作分解 (WBS)** 與 **影響面**（碼路徑、相依、風險）。
- 只在**允許的工具**範圍內擬定動作（見第 4 節）。
- **不要**透露冗長的內部思考；以**精簡要點**呈現計畫與風險。

### 執行 (Execute)
- 依序執行安全動作：優先 **讀檔→比對→最小變更**。
- 產生檔案變更時，請以 **diff** 或 **patch** 形式描述。
- 需要命令列時，僅使用白名單命令（例如 `git status`、`npm test`、`pytest -q`、`go test ./...`）。

### 驗證 (Verify)
- 執行單元測試 / 型別檢查 / Lint，並彙整結果。
- 附上**驗證清單**：是否破壞 API/ABI、是否變更資安面（檔案/網路/密鑰）。
- 給出**回退方案**與**後續待辦**。

## 3) 風格與品質
- 語言：中文優先，技術名詞可保留英文。
- Commit 訊息：`type(scope): summary`，如 `fix(api): handle null header`。
- 程式風格：遵循專案既有 Lint/Formatter；缺省時採 Prettier/ESLint。

## 4) 工具與邊界（極重要）
- 排程/管線：若在 CI，**必須**以 `--output-format json` 輸出並避免互動。
- Shell：**嚴禁** `rm -rf`, `sudo`, 任意 `curl|bash`。僅允許查閱與測試類命令。
- 檔案：優先使用 **最小編輯**；禁止在專案根目錄外新增/刪除檔案。
- 網路：除非明確要求，不得對外傳送專案內容。

## 5) 安全與信任
- 僅在**受信任資料夾**中運作（Trusted Folders）。
- 如需全自動（YOLO）或長流程，**必須**在 Sandbox 中執行並啟用工具白名單。
- 使用 MCP 伺服器時，僅允許列出的工具；所有憑證以環境變數注入。

## 6) 常見工作
- **測試**：`!npm test`、`!pytest -q`、`!go test ./...`
- **文件**：為 Public API 產生 JSDoc / docstring，並更新 `CHANGELOG.md`。
- **重構**：先產出遷移計畫與回退步驟，再進行最小變更。

---

> 摘要：先計畫、後執行、必驗證；小步快跑、風險前置、可回溯、可回退。



## **1\. 專案目標**

本專案旨在建立一個基於 Kubernetes 的 O-RAN 光通訊管理平台。我們將採用 Infrastructure as Code (IAC) 和 GitOps 的原則，利用 Crossplane 來統一宣告式地管理雲原生網元 (CNF) 和實體光通訊設備 (PNF)，並開發一個客製化的 Operator (Adapter) 來橋接控制平面與實體硬體，最終實現從 SMO 到網路邊緣設備的端到端自動化管理與監控。

## **2\. 專案架構 (Project Architecture)**

我們的系統架構分為以下幾個主要層次：

graph TD  
    subgraph "管理與編排層 (Management & Orchestration)"  
        SMO\[SMO / 管理入口\]  
    end

    subgraph "K8s 控制平面 (Control Plane)"  
        A\[Git Repo for Configs\] \--\>|Config Sync| B(Kubernetes Cluster)  
        SMO \--\>|O1 Interface API| C(Nephio/Orchestrator)  
        C \--\>|Creates Claims| D\[Crossplane Claims (XR Claims)\]  
        D \-- "Requests an XR" \--\> E{Crossplane}  
        E \-- "Selects Composition" \--\> F\[Composition\]  
        F \-- "Creates Managed Resources" \--\> G\[Custom Resource 'OpticalDevice'\]  
        F \-- "Creates Managed Resources" \--\> H\[K8s Resources e.g., ConfigMap\]  
    end

    subgraph "客製化控制器 (Custom Controller)"  
        I\[O-RAN Hardware Adapter\]  
    end

    subgraph "實體與虛擬資源層 (Managed Resources)"  
        J\[光通訊設備 (Laser, ONU)\]  
        K\[網路功能 (AGF, CU/DU)\]  
    end

    B \-.-\> I  
    I \-- "Watches" \--\> G  
    I \-- "Controls (gNMI/API)" \--\> J  
    B \-- "Manages" \--\> K

* **核心流程**: 管理者透過 SMO 發出意圖，Orchestrator 在 K8s 中建立一個 Claim。Crossplane 捕獲此 Claim，並根據對應的 Composition 建立一組資源，其中包括一個我們自訂的 OpticalDevice CR。O-RAN Hardware Adapter 會監聽到這個 CR 的變化，並將其規格 (spec) 轉化為對實體光通訊設備的控制指令。

## **3\. Gemini CLI 使用指南**

本專案將全面採用 Gemini CLI 來提升開發效率。

* **初始化**: gemini init  
* **程式碼生成**: gemini \-p "你的提示詞" \> path/to/file.go  
* **程式碼解釋**: gemini \-p "Explain this Go function:" \-c path/to/file.go  
* **測試生成**: gemini \-p "Write a unit test for the Reconcile function in this file:" \-c path/to/controller.go \> path/to/controller\_test.go  
* **文件生成**: gemini \-p "Generate a README.md for this project, describing the architecture and setup." \> README.md

## **4\. 開發階段 (Development Phases)**

| 階段 | 名稱 | 目標 |
| :---- | :---- | :---- |
| **P1** | **環境建置與基礎設定** | 建立一個功能齊全的本地開發環境，包含 K8s 叢集和所有核心工具。 |
| **P2** | **Crossplane 核心抽象層開發** | 定義光通訊服務的抽象模型 (XRD & Composition)，使其能透過 K8s API 進行管理。 |
| **P3** | **硬體 Adapter 控制器開發** | 開發客製化的 Kubernetes Operator，用於監聽 CR 並與（模擬的）硬體互動。 |
| **P4** | **端到端 (E2E) 流程整合** | 串接從 Claim 建立到硬體設定的完整流程，並實現參數透傳。 |
| **P5** | **監控與視覺化** | 整合 Grafana，建立儀表板來監控系統狀態和設備性能。 |

詳細的開發提示詞請參考 prompts.md 檔案。