# M1 基础骨架 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: subagent-driven-development or inline execution with TDD.

**Goal:** 搭建 TopMentor monorepo，Go 后端可连接 PostgreSQL/Redis，迁移就绪，健康检查 API 可测。

**Architecture:** Gin 分层 API；SQL 文件迁移；docker-compose 本地依赖。

**Tech Stack:** Go 1.22, Gin, pgx, go-redis, testify

---

### Task 1: 基础设施

**Files:**
- Create: `docker-compose.yml`, `.env.example`, `README.md`
- Create: `backend/migrations/001_init.sql`

- [x] docker-compose 定义 postgres:15 + redis:7
- [x] `.env.example` 对齐 SDS §4

### Task 2: 后端骨架

**Files:**
- Create: `backend/go.mod`, `backend/cmd/server/main.go`, `backend/cmd/migrate/main.go`
- Create: `backend/internal/config/config.go`
- Create: `backend/pkg/response/response.go`

### Task 3: 健康检查（TDD）

**Files:**
- Create: `backend/internal/handler/health_test.go`
- Create: `backend/internal/handler/health.go`
- Create: `backend/internal/service/health.go`

**验证命令:**
```bash
docker compose up -d
cd backend && go run ./cmd/migrate
cd backend && go test ./...
cd backend && go run ./cmd/server
curl http://localhost:8080/api/v1/health
```
