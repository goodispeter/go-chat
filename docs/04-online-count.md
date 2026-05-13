# Step 04: 線上人數顯示

## 做了什麼
1. Hub Register/Unregister 時推送 `[System] Online: X` 給所有人
2. 前端 JS 攔截 Online 訊息，更新 header 而不是顯示在聊天區
3. 抽出 `sendAll` method 消除重複程式碼
4. join/leave 訊息從 ws.go 移到 hub.go 統一處理

## 踩過的坑

### deadlock：select 迴圈裡自己送 Broadcast 給自己
Hub.Run() 是單一 goroutine 的 select 迴圈。如果在 Register case 裡往 Broadcast channel 送訊息，但 Broadcast 的消費者也是同一個 select 迴圈——送的人等讀，讀的人在忙，死結。

解法：Register/Unregister case 裡不經過 Broadcast channel，直接用 sendAll 發送給所有 client。

### 訊息職責歸屬
join/leave/online 訊息全部由 Hub 負責（因為只有 Hub 知道連線狀態）。ws.go 的 handler 只負責：
- 升級 WebSocket
- 建立 Client
- 送 Register/Unregister
- readPump / writePump

## 學到的概念

### 抽 method 消除重複
重複的 for + select pattern 抽成 sendAll：
```go
func (h *Hub) sendAll(message []byte) {
    for client := range h.Clients {
        select {
        case client.Send <- message:
        default:
            close(client.Send)
            delete(h.Clients, client)
        }
    }
}
```

### fmt.Appendf 取代 []byte(fmt.Sprintf(...))
```go
// 之前：先建 string 再轉 []byte
[]byte(fmt.Sprintf("[System] Online: %d", count))

// 之後：直接產生 []byte，少一次轉換
fmt.Appendf(nil, "[System] Online: %d", count)
```
