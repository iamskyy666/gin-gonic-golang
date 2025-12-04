A **small, production-oriented CRUD REST API** in Go using **Gin** and GORM. We’ll show a clean folder structure, explain why things are placed where they are, and give working code for a `User` resource (Create / Read / Update / Delete). We’ll also include best practices: config via env, logging, graceful shutdown, middleware, repository/service/handler separation, Dockerfile, and tests scaffold.

We’ll use:

* Gin for HTTP routing
* GORM (with PostgreSQL driver) for DB access
* logrus for structured logs (replaceable)
* `.env` for configuration
* A repository → service → handler separation for testability
* Graceful shutdown and context propagation

---

# 1 — Project layout (recommended)

```
project/
├─ cmd/
│  └─ server/
│     └─ main.go
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ db/
│  │  └─ postgres.go
│  ├─ logger/
│  │  └─ logger.go
│  ├─ migrations/
│  │  └─ migrate.go
│  ├─ models/
│  │  └─ user.go
│  ├─ repository/
│  │  └─ user_repo.go
│  ├─ service/
│  │  └─ user_service.go
│  ├─ handler/
│  │  └─ user_handler.go
│  ├─ router/
│  │  └─ router.go
│  └─ middleware/
│     ├─ recovery.go
│     ├─ logger_middleware.go
│     └─ cors.go
├─ Dockerfile
├─ docker-compose.yml
├─ .env.example
├─ go.mod
├─ go.sum
└─ Makefile
```

---

# 2 — Configuration (.env example)

`.env.example`

```
ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:password@db:5432/mydb?sslmode=disable
LOG_LEVEL=info
```

Load with `github.com/joho/godotenv` or use os.Getenv directly in production (we’ll show a small loader).

---

# 3 — `go.mod`

```go
module github.com/ourorg/ourapi

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.4
    gorm.io/driver/postgres v1.4.6
    github.com/sirupsen/logrus v1.10.0
    github.com/joho/godotenv v1.5.1
)
```

(Adjust versions as needed.)

---

# 4 — Key files & code

Below are the minimal, ready-to-run files. We keep things short but complete so we can run it locally.

### `cmd/server/main.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ourorg/ourapi/internal/config"
	"github.com/ourorg/ourapi/internal/db"
	"github.com/ourorg/ourapi/internal/logger"
	"github.com/ourorg/ourapi/internal/migrations"
	"github.com/ourorg/ourapi/internal/router"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.New(cfg)
	log.WithField("env", cfg.Env).Info("starting service")

	// Init DB
	gormDB, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to db")
	}
	// Run migrations
	if err := migrations.Run(gormDB, log); err != nil {
		log.WithError(err).Fatal("migrations failed")
	}

	// Build router
	r := router.SetupRouter(gormDB, log, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Infof("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("listen error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("server forced to shutdown")
	}
	log.Info("server exiting")
}
```

---

### `internal/config/config.go`

```go
package config

import (
	"errors"
	"os"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	LogLevel    string
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:         getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
```

---

### `internal/logger/logger.go`

```go
package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ourorg/ourapi/internal/config"
)

func New(cfg *config.Config) *logrus.Entry {
	level, _ := logrus.ParseLevel(strings.ToLower(cfg.LogLevel))
	log := logrus.New()
	log.SetLevel(level)
	// For production set JSONFormatter
	if cfg.Env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log.WithField("service", "ourapi")
}
```

---

### `internal/db/postgres.go`

```go
package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ourorg/ourapi/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dialector := postgres.Open(cfg.DatabaseURL)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
```

---

### `internal/migrations/migrate.go`

```go
package migrations

import (
	"github.com/ourorg/ourapi/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, log *logrus.Entry) error {
	// AutoMigrate is simple for examples. In real proj, use proper migration tool.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	log.Info("migrations applied")
	return nil
}
```

---

### `internal/models/user.go`

```go
package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

---

### `internal/repository/user_repo.go`

```go
package repository

import (
	"errors"

	"github.com/ourorg/ourapi/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	List(offset, limit int) ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) List(offset, limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
```

---

### `internal/service/user_service.go`

```go
package service

import (
	"errors"

	"github.com/ourorg/ourapi/internal/models"
	"github.com/ourorg/ourapi/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	CreateUser(u *models.User) error
	GetUser(id uint) (*models.User, error)
	ListUsers(offset, limit int) ([]models.User, error)
	UpdateUser(u *models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(u *models.User) error {
	// business rules (e.g., validate email etc.) can go here
	return s.repo.Create(u)
}

func (s *userService) GetUser(id uint) (*models.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (s *userService) ListUsers(offset, limit int) ([]models.User, error) {
	return s.repo.List(offset, limit)
}

func (s *userService) UpdateUser(u *models.User) error {
	existing, err := s.repo.GetByID(u.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	return s.repo.Update(u)
}

func (s *userService) DeleteUser(id uint) error {
	// we could verify existence first if desired
	return s.repo.Delete(id)
}
```

---

### `internal/handler/user_handler.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ourorg/ourapi/internal/models"
	"github.com/ourorg/ourapi/internal/service"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	usrSvc service.UserService
	log    *logrus.Entry
}

func NewUserHandler(svc service.UserService, log *logrus.Entry) *UserHandler {
	return &UserHandler{usrSvc: svc, log: log}
}

func (h *UserHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/", h.create)
	rg.GET("/", h.list)
	rg.GET("/:id", h.get)
	rg.PUT("/:id", h.update)
	rg.DELETE("/:id", h.delete)
}

type createUserRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
}

func (h *UserHandler) create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("invalid create request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := h.usrSvc.CreateUser(u); err != nil {
		h.log.WithError(err).Error("create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) list(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	users, err := h.usrSvc.ListUsers(offset, limit)
	if err != nil {
		h.log.WithError(err).Error("list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	u, err := h.usrSvc.GetUser(uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		h.log.WithError(err).Error("get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get failed"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &models.User{
		ID:    uint(id),
		Name:  req.Name,
		Email: req.Email,
	}
	if err := h.usrSvc.UpdateUser(u); err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		h.log.WithError(err).Error("update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usrSvc.DeleteUser(uint(id)); err != nil {
		h.log.WithError(err).Error("delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}
```

---

### `internal/router/router.go`

```go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ourorg/ourapi/internal/handler"
	"github.com/ourorg/ourapi/internal/logger"
	"github.com/ourorg/ourapi/internal/repository"
	"github.com/ourorg/ourapi/internal/service"
	"github.com/ourorg/ourapi/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, log *logrus.Entry, cfg *config.Config) *gin.Engine {
	// Choose Gin mode by env
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// middlewares (custom and built-in)
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(log))
	r.Use(CORSMiddleware())

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// wire dependencies
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userH := handler.NewUserHandler(userSvc, log)

	api := r.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")
	userH.Register(users)

	return r
}
```

---

### `internal/middleware/logger_middleware.go`

```go
package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(log *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		log.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"latency": latency.String(),
			"ip":      c.ClientIP(),
		}).Info("request completed")
	}
}
```

(We placed middleware under `router` in the example for brevity; in the layout we can keep in `internal/middleware`.)

---

# 5 — Dockerfile & docker-compose

`Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ourapi ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /ourapi /ourapi
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/ourapi"]
```

`docker-compose.yml` (dev)

```yaml
version: "3.8"
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  api:
    build: .
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://postgres:password@db:5432/mydb?sslmode=disable
      ENV: development
      PORT: 8080
    ports:
      - "8080:8080"

volumes:
  pgdata:
```

---

# 6 — Best practices & explanations (why we did this)

1. **Folder separation** (internal vs cmd): `internal` holds the application code (not importable by other modules). `cmd/server` is the app entrypoint.
2. **Repository → Service → Handler layers**:

   * Repository: DB-specific queries.
   * Service: Business logic and validation.
   * Handler: HTTP layer; bind/json/requests/responses.
     This improves testability and separation-of-concerns.
3. **Use dependency injection**: pass `db`, `logger` and repos/services into handlers so we can mock them in tests.
4. **Graceful shutdown**: we listen to SIGINT/SIGTERM and shut down with a context timeout.
5. **Structured logging**: we used logrus and JSON in production. Structured logs are easier to ingest.
6. **Validation**: use Gin binding tags (`binding:"required,email"`) to validate request payloads.
7. **Pagination & limits**: always offer `limit` and `offset` to list endpoints to avoid returning unbounded results.
8. **Migrations**: for demo we used `AutoMigrate`, but in real prod use a migration tool (e.g., `golang-migrate/migrate`) so migrations are tracked and reversible.
9. **Error handling**: map internal errors to proper HTTP codes; don’t leak internals.
10. **Secrets/config**: keep in env (or secret stores in production), do not commit `.env` with secrets to git.
11. **Database connection pooling**: we used GORM; configure `sql.DB` pool settings if necessary (MaxIdleConns, MaxOpenConns).
12. **Health & readiness**: provide `/health` (and possibly `/ready`) endpoints for orchestration.
13. **CORS, rate-limiting, auth**: add middleware for CORS, rate limits, and authentication (JWT/OAuth) as needed.
14. **Observability**: add metrics (Prometheus), tracing (OpenTelemetry), and request/response logs in production.
15. **Testing**: mock repository interfaces in unit tests; add integration tests against a test DB (dockerized).

---

# 7 — Quick tests & examples

Example curl commands:

* Create:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"Alice","email":"alice@example.com"}' http://localhost:8080/api/v1/users
```

* List:

```bash
curl http://localhost:8080/api/v1/users
```

* Get:

```bash
curl http://localhost:8080/api/v1/users/1
```

* Update:

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Alicia","email":"alicia@example.com"}' http://localhost:8080/api/v1/users/1
```

* Delete:

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

# 8 — Testing strategy (quick)

* Unit tests for service layer with mocked repository (table-driven tests).
* Handler tests using `httptest` and a mock service.
* Integration tests that run against a disposable Docker Postgres; use `TestMain` to set up/tear down.

---

# 9 — Next steps / enhancements we recommend

* Add JWT auth middleware + role-based access control.
* Introduce request/response DTOs and transform models to avoid leaking DB models.
* Add OpenAPI (Swagger) docs — `swaggo/swag`.
* Add Prometheus metrics endpoints and distributed tracing.
* Replace `AutoMigrate` with structured migrations using `migrate`.
* Add CI pipelines (lint, vet, unit tests, build).
* Add image build & multi-arch if required.

---

Let’s add JWT-based authentication to the project we sketched. We’ll:

* Add configuration for JWT secrets and TTL.
* Store password hashes on `User` and add registration + login endpoints.
* Create an `AuthService` that issues JWTs.
* Add an `AuthMiddleware` that protects routes and sets the authenticated user ID in Gin context.
* Show minimal repository/service/handler changes and where to wire everything.

We’ll use `github.com/golang-jwt/jwt/v5` and `golang.org/x/crypto/bcrypt`. Keep secrets in env and follow best practices (rotate secrets, use refresh tokens for long sessions, store tokens securely on clients, use HTTPS).

Below are the essential code changes — drop them into the same folders from the earlier layout.

---

## 1) config: add JWT fields

`internal/config/config.go` — extend `Config`:

```go
type Config struct {
    Env         string
    Port        string
    DatabaseURL string
    LogLevel    string
    JWTSecret   string
    JWTTTL      time.Duration // e.g., in minutes
}

func Load() (*Config, error) {
    ttlMinutes, _ := strconv.Atoi(getEnv("JWT_TTL_MINUTES", "1440")) // default 24h
    cfg := &Config{
        Env:         getEnv("ENV", "development"),
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        JWTSecret:   getEnv("JWT_SECRET", ""),
        JWTTTL:      time.Duration(ttlMinutes) * time.Minute,
    }
    if cfg.DatabaseURL == "" {
        return nil, errors.New("DATABASE_URL is required")
    }
    if cfg.JWTSecret == "" {
        return nil, errors.New("JWT_SECRET is required")
    }
    return cfg, nil
}
```

Add `import` for `time` and `strconv`.

Add `.env.example`:

```
JWT_SECRET=replace_with_secure_random
JWT_TTL_MINUTES=1440
```

---

## 2) models: add password hash

`internal/models/user.go` — add PasswordHash (omit from JSON responses):

```go
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `gorm:"type:varchar(100);not null" json:"name"`
    Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

Note: `json:"-"` prevents the password hash from being serialized.

---

## 3) repository: get by email

`internal/repository/user_repo.go` — add method:

```go
func (r *userRepo) GetByEmail(email string) (*models.User, error) {
    var u models.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}
```

And add it to the `UserRepository` interface:

```go
GetByEmail(email string) (*models.User, error)
```

---

## 4) service: password handling + auth service

### a) make user creation hash password

`internal/service/user_service.go` — update `CreateUser` signature to accept raw password or create a helper. Example: create `CreateUserWithPassword`.

```go
import "golang.org/x/crypto/bcrypt"

// Add to interface
CreateUser(u *models.User, rawPassword string) error

// Implementation
func (s *userService) CreateUser(u *models.User, rawPassword string) error {
    if len(rawPassword) < 6 {
        return errors.New("password too short")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hash)
    return s.repo.Create(u)
}
```

Also keep existing `UpdateUser` etc.

### b) new auth service to issue tokens

`internal/service/auth_service.go`

```go
package service

import (
    "time"
    "errors"

    "github.com/golang-jwt/jwt/v5"
    "github.com/ourorg/ourapi/internal/config"
    "github.com/ourorg/ourapi/internal/models"
    "github.com/ourorg/ourapi/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Authenticate(email, password string) (*models.User, error) // returns user if success
    GenerateToken(user *models.User) (string, error)
    ParseToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type authService struct {
    repo repository.UserRepository
    cfg  *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
    return &authService{repo: repo, cfg: cfg}
}

func (a *authService) Authenticate(email, password string) (*models.User, error) {
    u, err := a.repo.GetByEmail(email)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, errors.New("invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    return u, nil
}

func (a *authService) GenerateToken(user *models.User) (string, error) {
    now := time.Now()
    exp := now.Add(a.cfg.JWTTTL)
    claims := jwt.RegisteredClaims{
        Subject:   fmt.Sprint(user.ID),
        IssuedAt:  jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(exp),
        // optionally add Audience, Issuer, ID ...
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(a.cfg.JWTSecret))
}

func (a *authService) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
        // ensure signing method
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(a.cfg.JWTSecret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}
```

Add imports `fmt` and other packages. (We used `fmt.Sprint` for subject.)

---

## 5) handler: auth endpoints (register + login)

`internal/handler/auth_handler.go`

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/ourorg/ourapi/internal/models"
    "github.com/ourorg/ourapi/internal/service"
    "github.com/sirupsen/logrus"
)

type AuthHandler struct {
    userSvc service.UserService
    authSvc service.AuthService
    log     *logrus.Entry
}

func NewAuthHandler(us service.UserService, as service.AuthService, log *logrus.Entry) *AuthHandler {
    return &AuthHandler{userSvc: us, authSvc: as, log: log}
}

type registerRequest struct {
    Name     string `json:"name" binding:"required,min=2"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.POST("/register", h.register)
    rg.POST("/login", h.login)
}

func (h *AuthHandler) register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    u := &models.User{
        Name:  req.Name,
        Email: req.Email,
    }
    if err := h.userSvc.CreateUser(u, req.Password); err != nil {
        h.log.WithError(err).Warn("register failed")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Do not return password or sensitive info
    c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
}

func (h *AuthHandler) login(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    u, err := h.authSvc.Authenticate(req.Email, req.Password)
    if err != nil {
        h.log.WithError(err).Warn("login failed")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }
    token, err := h.authSvc.GenerateToken(u)
    if err != nil {
        h.log.WithError(err).Error("token generation failed")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}
```

---

## 6) middleware: JWT auth middleware

`internal/middleware/auth_middleware.go`

```go
package middleware

import (
    "errors"
    "net/http"
    "strings"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/ourorg/ourapi/internal/service"
    "github.com/sirupsen/logrus"
)

func AuthMiddleware(authSvc service.AuthService, log *logrus.Entry) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }
        tokenStr := parts[1]
        claims, err := authSvc.ParseToken(tokenStr)
        if err != nil {
            log.WithError(err).Warn("token parse failed")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        // claims.Subject should be the user ID
        if claims.Subject == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            return
        }
        uid64, err := strconv.ParseUint(claims.Subject, 10, 32)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject in token"})
            return
        }
        // set user id in context
        c.Set("userID", uint(uid64))
        c.Next()
    }
}
```

We used `service.AuthService` to parse token so secret handling is centralized.

---

## 7) router: wire auth routes and protect endpoints

`internal/router/router.go` — update wiring to create AuthService and AuthHandler, and protect the user routes:

```go
// wire dependencies
userRepo := repository.NewUserRepository(db)
userSvc := service.NewUserService(userRepo)
authSvc := service.NewAuthService(userRepo, cfg) // auth service needs repo + cfg
userH := handler.NewUserHandler(userSvc, log)
authH := handler.NewAuthHandler(userSvc, authSvc, log)

// public auth routes
api := r.Group("/api")
v1 := api.Group("/v1")
authGroup := v1.Group("/auth")
authH.RegisterRoutes(authGroup)

// protected users routes
users := v1.Group("/users")
users.Use(middleware.AuthMiddleware(authSvc, log))
userH.Register(users)
```

Add imports for `middleware` and `handler` paths accordingly.

---

## 8) Using the authenticated user in handlers

Inside handlers where we need current user ID:

```go
uid, exists := c.Get("userID")
if exists {
    userID := uid.(uint) // assert
    // use userID
}
```

For safer code, assert with type switch.

---

## 9) Example requests

Register:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
# => returns {"token":"<JWT_TOKEN>"}
```

Call a protected endpoint:

```bash
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## 10) Best-practice notes & next steps

* **Keep JWT secret secure**: set via secrets manager (HashiCorp/ environ / K8s secret), not in code or committed files.
* **Token expiry & refresh tokens**: we used a short-lived access token pattern; for long sessions implement rotating refresh tokens stored server-side or as HTTP-only secure cookies.
* **Revoke tokens**: JWTs are stateless; to revoke create token blacklist or use short TTL + refresh token revocation mechanism.
* **Claims**: include minimal claims (sub, exp, iat). Avoid sensitive data in token payload.
* **HTTPS only**: always use HTTPS in production.
* **Rate-limit login**: prevent brute-force.
* **Password policy**: enforce and consider 2FA.
* **Use `scrypt`/`argon2` for password hashing** if desired; bcrypt is acceptable.
* **Testing**: mock `AuthService` in handler tests and test middleware behaviour.

---

