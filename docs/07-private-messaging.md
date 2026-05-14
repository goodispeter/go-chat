# Step 07: 私訊功能

## 做了什麼
1. `hub.go` — ClientsByUserID map（O(1) 查找）、SendToUser、broadcastUserList（id:name 格式）
2. `ws.go` — readPump 改用 JSON 解析（type/to/text）、Client 帶 UserID
3. `app.html` — 全新私訊介面、JSON 發送、解析 id:name 名單

## 後端架構

### WebSocket 訊息協定
前端 → 後端（JSON）：
```json
{"type": "pm", "to": 1, "text": "你好"}
```

後端 → 前端（字串）：
```
[System] Users: 1:Alice,2:Bob,3:Carol
[PM:Alice:14:30:05] 你好
```

### Hub 的兩個 map
```go
Clients       map[*Client]bool     // 遍歷所有連線（廣播用）
ClientsByUserID map[uint]*Client   // O(1) 查找特定使用者（私訊用）
```
Register 時兩個都加，Unregister 時兩個都刪。

### JWT claims 的 float64 問題
```go
UserID: uint(c.GetFloat64("user_id"))
```
JWT MapClaims 把所有數字解析成 float64（JSON 規格），取出時要轉型。

## 前端架構

### activeChat 從 string 改為 object
```javascript
activeChat = { id: 1, name: "Alice" }  // 之前只是 "Alice"
```
送訊息用 id，顯示用 name。

### conversations 用 name 當 key
```javascript
conversations = {
    "Alice": [{ from: "Alice", text: "hi", time: "14:30", mine: false }],
    "Bob": [...]
}
```
用 name 是因為收到 PM 時只知道 sender name，要能找到對應的對話。
