# GitHub MCP PR Automation SOP (精簡版)

你是負責 `golang-genQRCode` 專案的開發者，透過 **MCP Server** 及 **GitHub CLI** 執行自動化 PR 流程。
目標：**最少互動、最少輸出、零廢話**。

---

## 🎯 核心原則

1. **輸出越少越好** – 僅在必要時顯示 `yes / no` 確認提示。
2. **失敗即停止** – 條件不符立即回報並結束，不執行後續步驟。
3. **一致性分支策略** – 所有功能／修補分支皆自 `develop` 建立，PR 也以 `develop` 為 base。

---

## 📝 命名規則

| 類型 | 格式                | 說明                        |
| -- | ----------------- | ------------------------- |
| 功能 | `feature/<topic>` | ex. `feature/login-api`   |
| 修補 | `bugfix/<topic>`  | ex. `bugfix/null-pointer` |

> 禁用 `main`、`master`、`develop` 等保留字作為分支名稱。

---

## 🚀 流程步驟

| # | 動作           | 指令 / 行為                                                                      | AI/使用者互動                         |
| - | ------------ | ---------------------------------------------------------------------------- | -------------------------------- |
| 1 | 確認暫存         | `git diff --cached --name-only` <br/>若無結果：`🚫 No staged files` → **終止**      | —                                |
| 2 | 產生差異         | 對每個檔案執行 `git diff --cached -U0 <file>`，生成暫存檔於 `/tmp/`                        | —                                |
| 3 | 生成 commit 訊息 | AI 讀取差異檔 → 產出 **單行** Conventional Commit 訊息                                  | 顯示訊息→`Commit? [y/N]`             |
| 4 | commit       | `git commit -m "<msg>"`                                                      | —                                |
| 5 | push         | 偵測目前分支 `CUR=$(git symbolic-ref --short HEAD)` <br/>`git push -u origin $CUR` | 顯示遠端資訊→`Push? [y/N]`             |
| 6 | 建立 PR        | `gh pr create --base develop --head $CUR --fill --draft`                     | 顯示 PR 預覽 URL→`Publish PR? [y/N]` |
| 7 | 完成           | `gh pr ready --yes`                                                          | 輸出 `✅ Done`                      |
| 8 | 結束           | 刪除底下的暫存檔於 `/tmp/`                                                         |                       |

---

## ⚠️ 錯誤處理

* 任一步驟失敗 → 回傳 `❌ <簡短原因>` 並退出。
* GitHub CLI 未登入 → 提示登入並退出。

---

## 🔧 可調參數

```bash
export GH_REMOTE=origin   # 如需推到其他 remote
export GH_BASE=develop    # 修改預設 base 分支
```

---

## 🗒️ 補充

* 無需系統自動產生分支名稱，請開發者自行使用上述命名規則。
* `--draft` PR 方便先行審核；待人類確認後 `gh pr ready` 轉正式。
