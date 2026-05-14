# Step 06f: JWT 保護 WebSocket 路由

## 做了什麼
1. `cmd/server/main.go` — 用 Group + Use 把 `/ws` 路由加上 JWT 驗證
2. `internal/handler/ws.go` — 從 gin.Context 取 username（不再用 query parameter）

## 新概念

### Gin 路由群組控制 middleware
```go
// 公開路由 — 不經過 JWTAuth
api := r.Group("/api")
api.POST("/register", handler.Register)
api.POST("/login", handler.Login)

// 受保護路由 — 經過 JWTAuth
protected := r.Group("/")
protected.Use(middleware.JWTAuth())
protected.GET("/ws", ws.HandleWebSocket)
```
哪個 Group `.Use()` 了 middleware，底下的路由就會走。
等於 Java Spring 的 `antMatchers("/ws").authenticated()`。

### middleware 傳值給 handler
```go
// middleware 裡存
c.Set("username", claims["username"])

// handler 裡取
username := c.GetString("username")
```
WebSocket handler 不再需要 `c.Query("name")`，使用者身份從 JWT 取得。

## 目前的 request 流程
```
POST /api/register → handler.Register（公開）
POST /api/login    → handler.Login（公開，回傳 token）
GET  /ws           → JWTAuth middleware → ws.HandleWebSocket（需要 token）
```
