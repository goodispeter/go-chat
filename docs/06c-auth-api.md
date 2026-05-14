# Step 06c: 註冊/登入 API（三層架構）

## 安裝的套件
```bash
go get golang.org/x/crypto/bcrypt
```

## 做了什麼
1. `internal/repository/user.go` — DB 操作
2. `internal/service/auth.go` — 業務邏輯（bcrypt 加密/比對）
3. `internal/handler/auth.go` — HTTP handler

## 新概念

### 三層架構（等於 Java Spring 的 Controller → Service → Repository）
```
handler/auth.go    ← 只管 HTTP request/response
    ↓ 呼叫
service/auth.go    ← 業務邏輯（加密、驗證）
    ↓ 呼叫
repository/user.go ← 只管 DB 操作
```

每層的職責：
| 層 | 知道什麼 | 不知道什麼 |
|----|---------|----------|
| handler | HTTP、gin.Context | 密碼怎麼加密、SQL |
| service | bcrypt、業務規則 | HTTP 狀態碼 |
| repository | GORM、SQL | 為什麼要查 |

### bcrypt 自動加鹽
```go
// 加密：每次產生不同的 hash（自動加鹽）
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 比對：從 hash 裡取出 salt 再比
bcrypt.CompareHashAndPassword(hash, []byte(inputPassword))
```
不用自己管理 salt，重啟 server 也能比對，因為 salt 存在 hash 字串裡。

### Gin 的 ShouldBindJSON
```go
var req RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```
等於 Spring 的 `@RequestBody` + `@Valid`。`binding:"required"` 等於 `@NotNull`。
