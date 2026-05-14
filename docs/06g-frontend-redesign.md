# Step 06g: 前端重構 — 登入/註冊 + 聊天室

## 做了什麼
1. `web/login.html` — 全新登入/註冊頁面（JWT 認證流程）
2. `web/chat.html` — 重新設計的聊天室介面
3. `internal/middleware/jwt.go` — 支援 query parameter 傳 token（WebSocket 需要）
4. `cmd/server/main.go` — 接上 DB 連線、AutoMigrate、根路徑導向登入頁

## 設計決策

### 風格
- 深色主題，靈感來自 Discord/Slack
- 字體：Outfit（UI）+ JetBrains Mono（品牌/時間戳）
- 配色：紫色系主色（#6c5ce7）+ 暗色背景（#0a0a0f）
- 噪點紋理（noise texture）增加質感
- 背景漸變光暈（ambient gradient）增加層次

### 登入頁功能
- Sign In / Sign Up 頁籤切換
- 表單驗證（前端 + 後端）
- 密碼確認欄位
- Loading spinner 按鈕
- 註冊成功後自動登入
- 有 token 時自動跳轉聊天室

### 聊天室功能
- 從 JWT 取得 username（不用手動輸入）
- 使用者頭像（取名字首字母 + 基於名字的顏色）
- 系統訊息用圓角標籤樣式
- 使用者訊息有 avatar + 名字 + 時間戳
- 斷線偵測 + 重新連線按鈕
- Logout 按鈕（清除 token）

## 新概念

### WebSocket 無法帶 Authorization Header
瀏覽器的 `new WebSocket()` API 不支援自定義 header。
解法：用 query parameter 傳 token：`ws://host/ws?token=xxx`
middleware 改成先查 header，沒有再查 query parameter。

### localStorage 存 JWT
```javascript
localStorage.setItem('token', data.token)  // 登入後存
localStorage.getItem('token')               // 連 WebSocket 時取
localStorage.removeItem('token')            // 登出時刪
```
等於瀏覽器版的 key-value store，關閉分頁後資料還在。

### JWT 前端解碼（不驗簽）
```javascript
const payload = JSON.parse(atob(token.split('.')[1]));
```
JWT 的 payload 是 base64 編碼，前端可以解碼拿 username 顯示。
注意：這只是「看」不是「驗」，驗證在 server 端做。
