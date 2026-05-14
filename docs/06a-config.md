# Step 06a: Config 設定管理

## 安裝的套件
```bash
go get github.com/joho/godotenv
```

## 做了什麼
1. 建立 `internal/config/config.go` 統一管理設定
2. 用 godotenv 載入 `.env` 檔案
3. `.gitignore` 排除 `.env` 機敏資訊

## 新概念

### 環境變數管理
Go 沒有 Spring 的 `application.yml`，社群主流做法是 `.env` 檔 + `os.Getenv`：
- `.env` — 開發環境設定
- `.env.production` — 正式環境設定
- `APP_ENV=production go run ./cmd/server/` — 用環境變數切換

### godotenv
`godotenv.Load(".env")` 把 `.env` 裡的 key=value 載入到環境變數。
等於 Java Spring Boot 的 `application-dev.properties`。

### 不進 git 的機敏資訊
密碼、JWT secret、API key 放 `.env`，`.gitignore` 排除。
團隊用 `.env.example` 放範本（不含真實值），新人 copy 一份改值。
