# TopMentor

聚焦 6-14 岁少儿的 1 对 1 音视频双语陪伴平台（清华等名校学霸 × 中产家庭）。

## 文档

| 文档 | 说明 |
| :--- | :--- |
| [docs/01-PRD.md](docs/01-PRD.md) | 产品需求与业务规则 |
| [docs/02-HLD.md](docs/02-HLD.md) | 架构与技术选型 |
| [docs/03-SDS.md](docs/03-SDS.md) | 数据库、API、目录结构 |
| [docs/06-Classroom.md](docs/06-Classroom.md) | M4 课室（Agora + 心跳） |

## 运营后台

```bash
cd admin
npm install
npm run dev
# → http://localhost:5174  账号 admin / admin123
```

## 快速开始（M1）

### 前置

- Go 1.18+
- Docker Desktop

### 启动

```bash
# 1. 复制环境变量
cp .env.example .env

# 2. 启动数据库
docker compose up -d

# 3. 数据库迁移
cd backend && go run ./cmd/migrate

# 4. 启动 API
go run ./cmd/server

# 5. 健康检查
curl http://localhost:8080/api/v1/health
```

### 测试

```bash
cd backend && go test ./...
```

## 仓库结构

```
TopMentor/
├── backend/    # Go Gin API
├── web/        # H5 移动 Web（Vue3 + Vant）
├── admin/      # Vue3 运营后台（Element Plus）
└── docs/       # 设计文档
```

## H5 前端启动

```bash
cd web
npm install
npm run dev
# → http://localhost:5173
```

登录：手机号 + 验证码 `123456`（开发 mock）

## 后端 API

```bash
# 发送验证码
curl -X POST http://localhost:8080/api/v1/auth/sms/send \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000"}'

# 短信登录
curl -X POST http://localhost:8080/api/v1/auth/sms/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","code":"123456"}'

# 学霸广场
curl "http://localhost:8080/api/v1/mentors?school=清华"
```

## Linux 生产部署（域名 HTTPS）

```bash
cp deploy/env.production.example .env   # 修改密码与密钥
chmod +x deploy/deploy.sh
./deploy/deploy.sh
```

- H5：`https://toppeertalk.com/welcome`
- 后台：`https://admin.toppeertalk.com`

从 IP 测试切域名：`./deploy/switch-to-domain.sh`。详见 [deploy/README.md](deploy/README.md)。
