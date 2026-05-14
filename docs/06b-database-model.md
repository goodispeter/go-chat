# Step 06b: PostgreSQL 連線 + User Model

## 安裝的套件
```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

## 做了什麼
1. `internal/database/db.go` — GORM 連線 PostgreSQL
2. `internal/model/user.go` — User struct

## 新概念

### GORM 等於 Java 的 Hibernate
```go
// GORM                              // Java JPA
gorm.Open(postgres.Open(dsn))       // EntityManagerFactory
DB.Create(&user)                    // repository.save(user)
DB.Where("username = ?", x).First() // findByUsername(x)
DB.AutoMigrate(&User{})             // hibernate.ddl-auto=update
```

### Struct Tag 控制 GORM 行為
```go
type User struct {
    ID       uint   `gorm:"primaryKey"`           // 主鍵
    Username string `gorm:"uniqueIndex;not null"`  // 唯一索引 + 不可空
    Password string `json:"-"`                     // API 回應時隱藏
}
```
- `gorm:"uniqueIndex"` 等於 `@Column(unique = true)`
- `json:"-"` 等於 `@JsonIgnore`

### AutoMigrate
`DB.AutoMigrate(&User{})` 會自動建 table、加欄位。
不會刪欄位（安全考量），等於 Hibernate 的 `ddl-auto=update`。
