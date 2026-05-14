# Step 06e: JWT Middleware

## 做了什麼
建立 `internal/middleware/jwt.go`，驗證 request 的 JWT token，保護需要登入的路由。

## 新概念

### Gin Middleware 格式
```go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 驗證邏輯
        c.Next()  // 放行
    }
}
```
回傳 `gin.HandlerFunc`，等於 Java 的 `Filter`。

### c.Abort() vs c.Next()
- `c.Abort()` — 中斷，後面的 handler 不執行（等於 Java Filter 不呼叫 `chain.doFilter()`）
- `c.Next()` — 放行，繼續執行下一個 handler（等於 `chain.doFilter()`）

### c.Set / c.Get — Request 內傳值
```go
// middleware 裡存
c.Set("username", claims["username"])

// handler 裡取
username := c.GetString("username")
```
等於 Java 的 `request.setAttribute()` / `request.getAttribute()`。
在同一個 request 的生命週期內傳遞資料。

### Authorization Header 格式
```
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9...
```
`Bearer` 是 OAuth2 標準前綴，用 `strings.TrimPrefix` 去掉後拿到純 token。

### 驗證流程
```
Request → JWTAuth middleware
           ├── 沒有 header → 401 + Abort
           ├── token 無效  → 401 + Abort
           └── token 有效  → c.Set(user info) → c.Next() → handler
```
