# Step 00: 專案初始化

## 做了什麼
1. `go mod init go-chat` — 建立 Go 模組
2. `go get github.com/gin-gonic/gin` — 安裝 HTTP 框架
3. `go get github.com/gorilla/websocket` — 安裝 WebSocket 套件

## 新概念

### Gin 是什麼
Go 最熱門的 HTTP 框架，等於 Java 的 Spring MVC。解決標準庫 `net/http` 的痛點：
- 路由參數：`/users/:id` 直接拿到 id（標準庫要自己 TrimPrefix 解析）
- 路由群組：`v1 := r.Group("/api/v1")`
- 中間件：`r.Use(gin.Logger())`
- JSON 綁定：`c.ShouldBindJSON(&body)` 一行搞定

### gorilla/websocket 是什麼
Go 生態最成熟的 WebSocket 實作。WebSocket 是一種協定，讓 server 和 client 保持長連線、雙向即時通訊。
- HTTP 是「client 問、server 答」（request-response）
- WebSocket 是「兩邊隨時都能發訊息」（bidirectional）

聊天室需要 WebSocket，因為有人發訊息時 server 要「主動推」給所有人，HTTP 做不到。

### go.sum 是什麼
`go get` 下載套件後會產生 `go.sum`，記錄每個依賴的 hash 值。等於 Java 的 `gradle.lock` 或 npm 的 `package-lock.json`。確保每次 build 用的是同一版本。

## 專案結構（目前）
```
go-chat/
├── go.mod       ← 模組定義 + 依賴列表
├── go.sum       ← 依賴版本鎖定
├── CLAUDE.md    ← 專案簡介
└── docs/
    └── 00-project-init.md  ← 你在這裡
```
