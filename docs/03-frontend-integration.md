# Step 03: 前端整合 + 首次聊天成功

## 做了什麼
1. `web/index.html` — 瀏覽器聊天介面（登入 + 聊天室）
2. `cmd/server/main.go` — 用 `StaticFS` 提供靜態檔案，避免路由衝突
3. 修正 `ws.go` — 補上 `go h.writePump(conn, client)`

## 踩過的坑

### 路由衝突：StaticFile("/") vs GET("/ws")
Gin 的 `StaticFile("/", ...)` 會註冊 `/*filepath` 萬用路由，跟 `/ws` 衝突導致 panic。
解法：用 `StaticFS("/web", ...)` 掛在子路徑，根路徑用 `Redirect` 導向 `/web/index.html`。

### 沒有 writePump = 訊息送不出去
每個 WebSocket 連線需要兩個 goroutine：
- `readPump`：讀 WebSocket → 塞進 Hub.Broadcast
- `writePump`：讀 Client.Send → 寫回 WebSocket

少了 writePump，Hub 廣播的訊息堆在 Client.Send channel 裡沒人讀，瀏覽器什麼都收不到。

### 前端分割訊息依賴格式
前端用 `text.indexOf(': ')` 分割名字和訊息，所以後端格式必須是 `"名字: 訊息"`（冒號後面要有空格）。
`fmt.Sprintf("%s: %s", name, message)` — 注意 `: ` 不是 `:`。

## 新概念

### Gin 靜態檔案服務
```go
r.StaticFS("/web", http.Dir("./web"))   // /web/index.html 對應 ./web/index.html
r.StaticFile("/favicon.ico", "./web/favicon.ico")  // 單一檔案
```

### WebSocket 連線流程
```
瀏覽器 JS                    Go Server
   │                            │
   ├── new WebSocket("/ws") ──→ upgrader.Upgrade()
   │                            ├── 建立 Client
   │                            ├── hub.Register <- client
   │                            ├── go readPump()
   │                            └── go writePump()
   │                            │
   ├── ws.send("hello") ─────→ readPump: ReadMessage()
   │                            ├── hub.Broadcast <- "名字: hello"
   │                            │
   │                            Hub.Run(): Broadcast → 所有 Client.Send
   │                            │
   ←── ws.onmessage ─────────── writePump: WriteMessage()
```
