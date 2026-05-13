# Step 02: WebSocket 核心架構

## 做了什麼
1. `internal/chat/hub.go` — Hub（聊天室中樞）管理所有連線和廣播
2. `internal/handler/ws.go` — WebSocket handler，處理連線升級和讀寫
3. `cmd/server/main.go` — 接上 Hub 和 WsHandler

## 架構圖
```
瀏覽器 A ──WebSocket──→ readPump goroutine ──→ Hub.Broadcast channel
瀏覽器 B ──WebSocket──→ readPump goroutine ──→ Hub.Broadcast channel
                                                      ↓
                                                Hub.Run() 主迴圈
                                                (select 監聽三個 channel)
                                                      ↓
                                            廣播到每個 Client.Send channel
                                                      ↓
                        writePump goroutine ←── Client A 的 Send channel
                        writePump goroutine ←── Client B 的 Send channel
                              ↓                        ↓
                        寫回 WebSocket A          寫回 WebSocket B
```

## 新概念

### select — channel 專用的 switch
```go
select {
case client := <-h.Register:    // 有人加入
case client := <-h.Unregister:  // 有人離開
case message := <-h.Broadcast:  // 有人說話
}
```
同時監聽多個 channel，誰有資料就處理誰。Java 沒有對應語法，最接近的是 NIO Selector。

### WebSocket 升級
HTTP 是 request-response，WebSocket 是雙向長連線。
`upgrader.Upgrade(w, r, nil)` 把普通 HTTP 連線升級成 WebSocket。
等於 Java 的 `@ServerEndpoint` 或 Spring 的 `WebSocketHandler`。

### 每個連線兩個 goroutine
- **readPump**：從 WebSocket 讀訊息 → 塞進 Hub.Broadcast
- **writePump**：從 Client.Send 讀訊息 → 寫回 WebSocket

這是 Go WebSocket server 的標準 pattern。Java 通常用一個 thread + callback（onMessage / onClose）。

### range channel
```go
for message := range client.Send {
    // Send channel 被 close 之前，會持續讀取
    // close 之後自動跳出迴圈
}
```

### 內層 select 防阻塞
```go
select {
case client.Send <- message:  // 送成功
default:                       // 送不進去 → client 卡住了，踢掉
    close(client.Send)
    delete(h.Clients, client)
}
```

## 目前專案結構
```
go-chat/
├── cmd/server/main.go          ← 進入點，組裝 Hub + Handler
├── internal/
│   ├── chat/
│   │   └── hub.go              ← Hub：管理連線 + 廣播
│   └── handler/
│       └── ws.go               ← WebSocket handler
└── docs/
```
