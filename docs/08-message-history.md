# Step 08: 聊天記錄持久化

## 做了什麼
1. `internal/model/message.go` — Message struct（sender_id, receiver_id, content）
2. `internal/repository/message.go` — SaveMessage / GetMessages
3. `internal/handler/message.go` — GET /api/messages?peer_id=1
4. `ws.go` — 送私訊時同時存 DB
5. `app.html` — 開聊天時從 API 拉歷史訊息
6. `main.go` — AutoMigrate Message + 註冊路由

## 新概念

### GORM 複合條件查詢
```go
DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
    userID, peerID, peerID, userID)
```
查兩人之間的所有對話（A→B 和 B→A），按時間排序。
等於 SQL：
```sql
SELECT * FROM messages
WHERE (sender_id = 1 AND receiver_id = 2) OR (sender_id = 2 AND receiver_id = 1)
ORDER BY created_at ASC LIMIT 100
```

### 不需要 service 層的判斷標準
handler 和 repository 之間如果沒有業務邏輯（加密、驗證、轉換），直接呼叫 repository 就好。
需要時才抽 service，不提前設計空殼。

### 前端只在第一次開聊天時拉歷史
```javascript
if (!conversations[user.name]) {
    // 第一次開 → 從 API 拉
}
// 之後的新訊息由 WebSocket 即時推送，不重複拉
```
