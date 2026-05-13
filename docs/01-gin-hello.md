# Step 01: Gin Hello World

## 做了什麼
1. 建立 `cmd/server/main.go` — 程式進入點
2. 用 Gin 啟動 HTTP server，`GET /` 回傳文字
3. 瀏覽器打開 `http://localhost:8080` 確認成功

## 新概念

### Gin 基本用法
```go
r := gin.Default()          // 建立 Gin 引擎（帶 logger + recovery middleware）
r.GET("/", handler)         // 註冊路由
r.Run(":8080")              // 啟動 server
```

### gin.Default() vs gin.New()
- `gin.Default()` = `gin.New()` + Logger + Recovery
- Logger：自動印出每個 request 的 method、path、狀態碼、耗時
- Recovery：handler panic 時自動回 500，不會整個 server 掛掉

### gin.Context
Gin 把 `http.ResponseWriter` + `*http.Request` 合併成 `*gin.Context`：
```go
// 標準庫
func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}

// Gin
func handler(c *gin.Context) {
    c.String(200, "hello")
}
```

### debug 模式的警告
全部是開發階段的提示，不是錯誤：
- `Running in "debug" mode` — 正式部署時用 `gin.SetMode(gin.ReleaseMode)`
- `You trusted all proxies` — 本地開發不影響，部署時設定信任的 proxy

## 跑法
```bash
go run ./cmd/server/
# 瀏覽器打開 http://localhost:8080
```
