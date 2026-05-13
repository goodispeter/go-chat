# Go Chat — 即時聊天室

## 專案目標
用 Go 做一個 WebSocket 即時聊天室，透過實作學習 goroutine、channel、併發管理。

## 技術選型
| 組件 | 選擇 | 原因 |
|------|------|------|
| HTTP 路由 | Gin | 路由簡潔、支援中間件、社群最大 |
| WebSocket | gorilla/websocket | Go 生態最成熟的 WebSocket 套件 |
| 前端 | 內嵌 HTML（embed） | 不需要額外前端專案 |
| 儲存 | 記憶體 | 聊天室不需要持久化 |

## 功能清單
- [ ] WebSocket 連線建立
- [ ] 訊息廣播（一人發話全員收到）
- [ ] 加入/離開通知
- [ ] 線上人數顯示
- [ ] 瀏覽器前端介面

## 目前進度
Step 01: Gin HTTP server 啟動成功（port 8080）
