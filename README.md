# Go Chat

用 Go 從零打造的即時聊天應用。從一個 Hello World 開始，逐步實作 WebSocket 聊天室、JWT 認證、PostgreSQL 持久化、私訊功能。

這是一個 **學習專案**，每一步都有完整的開發紀錄和學習筆記。

## 功能

- 註冊 / 登入（JWT 認證）
- 在線用戶列表（即時更新）
- 1 對 1 私訊（WebSocket 即時推送）
- 聊天記錄持久化（PostgreSQL）
- 響應式前端（桌面 / 手機）

## 技術棧

| 層級 | 技術 |
|------|------|
| HTTP 框架 | [Gin](https://github.com/gin-gonic/gin) |
| WebSocket | [gorilla/websocket](https://github.com/gorilla/websocket) |
| ORM | [GORM](https://gorm.io) + PostgreSQL |
| 認證 | [golang-jwt](https://github.com/golang-jwt/jwt) + bcrypt |
| 設定管理 | [godotenv](https://github.com/joho/godotenv) |
| 前端 | 純 HTML / CSS / JS（內嵌） |

## Quick Start

### 前置條件

- Go 1.21+
- PostgreSQL

### 安裝

```bash
git clone <repo-url>
cd go-chat
go mod tidy
```

### 設定

建立 `.env` 檔案：

```env
JWT_SECRET=your-dev-secret-key
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=
DB_NAME=go_chat
DB_PORT=5432
```

### 建立資料庫

```bash
psql -U postgres -c "CREATE DATABASE go_chat;"
```

### 啟動

```bash
# 開發環境
go run ./cmd/server/

# 正式環境
APP_ENV=production go run ./cmd/server/

# 或 build 後執行
go build -o go-chat.exe ./cmd/server/
./go-chat.exe
```

瀏覽器打開 `http://localhost:8080`

## 專案結構

```
go-chat/
├── cmd/server/
│   └── main.go                 ← 進入點：組裝 dependency、註冊路由
├── internal/
│   ├── chat/
│   │   └── hub.go              ← WebSocket Hub：連線管理、廣播、私訊路由
│   ├── config/
│   │   └── config.go           ← 環境變數 + .env 載入
│   ├── database/
│   │   └── db.go               ← GORM PostgreSQL 連線
│   ├── handler/
│   │   ├── auth.go             ← 註冊 / 登入 HTTP handler
│   │   ├── message.go          ← 歷史訊息 HTTP handler
│   │   └── ws.go               ← WebSocket handler（升級連線 + readPump/writePump）
│   ├── middleware/
│   │   └── jwt.go              ← JWT 驗證 middleware
│   ├── model/
│   │   ├── user.go             ← User GORM model
│   │   └── message.go          ← Message GORM model
│   ├── repository/
│   │   ├── user.go             ← User DB 操作
│   │   └── message.go          ← Message DB 操作
│   └── service/
│       ├── auth.go             ← 註冊 / 登入業務邏輯（bcrypt）
│       └── jwt.go              ← JWT 產生 / 解析
├── web/
│   ├── index.html              ← 前端路由（有 token → 聊天，沒有 → 登入）
│   ├── login.html              ← 登入 / 註冊頁面
│   ├── chat.html               ← 舊版公開聊天室（已棄用）
│   └── app.html                ← 主介面（在線列表 + 私訊）
├── docs/                       ← 開發日誌（每步一個 md）
│   ├── 00-project-init.md
│   ├── 01-gin-hello.md
│   ├── ...
│   └── journey/                ← 學習總結 HTML
│       ├── index.html          ← 總覽首頁
│       ├── 01-gin-websocket.html
│       ├── 02-auth-jwt.html
│       ├── 03-private-messaging.html
│       └── cheatsheet.html     ← Java → Go 對照表
├── .env                        ← 開發環境設定（不進 git）
├── .env.production             ← 正式環境設定（不進 git）
└── go.mod
```

## 實作步驟

每一步都有對應的開發日誌，記錄做了什麼、學到什麼、踩過的坑。

### Phase 1：WebSocket 聊天室

| Step | 主題 | 日誌 |
|------|------|------|
| 00 | 專案初始化（Gin + gorilla/websocket） | [docs/00-project-init.md](docs/00-project-init.md) |
| 01 | Gin HTTP server 啟動 | [docs/01-gin-hello.md](docs/01-gin-hello.md) |
| 02 | WebSocket 核心架構（Hub + readPump/writePump） | [docs/02-websocket-core.md](docs/02-websocket-core.md) |
| 03 | 前端整合 + 首次聊天成功 | [docs/03-frontend-integration.md](docs/03-frontend-integration.md) |
| 04 | 線上人數顯示 + sendAll 重構 | [docs/04-online-count.md](docs/04-online-count.md) |
| 05 | 訊息時間戳 | [docs/05-timestamp.md](docs/05-timestamp.md) |

### Phase 2：認證與資料庫

| Step | 主題 | 日誌 |
|------|------|------|
| 06a | Config 設定管理（godotenv） | [docs/06a-config.md](docs/06a-config.md) |
| 06b | PostgreSQL 連線 + User model | [docs/06b-database-model.md](docs/06b-database-model.md) |
| 06c | 註冊 / 登入 API（三層架構） | [docs/06c-auth-api.md](docs/06c-auth-api.md) |
| 06d | JWT token 產生 / 解析 | [docs/06d-jwt.md](docs/06d-jwt.md) |
| 06e | JWT middleware | [docs/06e-jwt-middleware.md](docs/06e-jwt-middleware.md) |
| 06f | JWT 保護 WebSocket 路由 | [docs/06f-protected-routes.md](docs/06f-protected-routes.md) |
| 06g | 前端重構（登入頁 + 聊天室 redesign） | [docs/06g-frontend-redesign.md](docs/06g-frontend-redesign.md) |

### Phase 3：私訊與持久化

| Step | 主題 | 日誌 |
|------|------|------|
| 07 | 私訊功能（WebSocket JSON + UserID） | [docs/07-private-messaging.md](docs/07-private-messaging.md) |
| 08 | 聊天記錄持久化（Message model + API） | [docs/08-message-history.md](docs/08-message-history.md) |

## 學習總結

完整的學習歷程和 Java → Go 對照表：

**[打開學習總結 →](docs/journey/index.html)**
