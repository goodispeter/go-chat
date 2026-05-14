# Step 06d: JWT Token

## 安裝的套件
```bash
go get github.com/golang-jwt/jwt/v5
```

## 做了什麼
1. `internal/service/jwt.go` — GenerateToken / ParseToken
2. `internal/service/auth.go` — Login 改為回傳 JWT token
3. `internal/handler/auth.go` — Login 回傳 `{"token": "..."}` 給前端
4. `cmd/server/main.go` — 註冊 `/api/register` 和 `/api/login` 路由

## 新概念

### JWT（JSON Web Token）
一個加密過的字串，裡面放使用者資訊。登入成功後 server 發 token 給 client，之後每個 request 帶上 token 就不用重新登入。

```
eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.xxxxx
        header              payload            signature
```

等於 Java 的 `Jwts.builder().setClaims(...).signWith(key).compact()`。

### HS256 vs ES256
- **HS256**（我們用的）= 對稱加密，一個 secret key 簽名和驗證
- **ES256** = 非對稱加密，私鑰簽名、公鑰驗證

HS256 適合單一 server。微服務架構用 ES256（各服務只需要公鑰就能驗證）。

### Claims
token 裡面放的資料：
```go
claims := jwt.MapClaims{
    "user_id":  userID,
    "username": username,
    "exp":      time.Now().Add(24 * time.Hour).Unix(), // 過期時間
}
```

### Gin 路由群組
```go
api := r.Group("/api")          // 建立 /api 前綴群組
api.POST("/register", handler)  // 實際路徑：/api/register
api.POST("/login", handler)     // 實際路徑：/api/login
```
等於 Spring 的 `@RequestMapping("/api")`。

## 測試指令
```bash
go run ./cmd/server/

# 註冊
curl -X POST http://localhost:8080/api/register -d '{"username":"peter","password":"123456"}'

# 登入（會回傳 token）
curl -X POST http://localhost:8080/api/login -d '{"username":"peter","password":"123456"}'
```
