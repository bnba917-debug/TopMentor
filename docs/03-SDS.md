# 03-SDS.md - 开发蓝图 (Engineering Blueprint)

## 1. 仓库目录结构 (Monorepo Layout)

```
TopMentor/
├── docs/                          # 产品/架构/开发文档
├── backend/                       # Go (Gin) API 服务
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/                # 环境变量与配置加载
│   │   ├── handler/               # HTTP 路由与入参校验
│   │   ├── service/               # 业务逻辑与事务
│   │   ├── repository/            # 数据库访问
│   │   └── model/                 # 领域实体与 DTO
│   ├── migrations/                # SQL 迁移脚本（按序号递增）
│   ├── pkg/
│   │   ├── response/              # 统一 API 响应封装
│   │   └── redis/                 # Redis 客户端封装
│   └── go.mod
├── web/                           # H5 移动 Web（Vue3 + Vite + TS + Vant）
│   ├── src/pages/                 # 家长端页面
│   ├── src/pages-mentor/          # 学霸端页面（或 routes/mentor）
│   └── vite.config.ts
├── admin/                         # Vue3 + Pinia + Element Plus 运营后台
├── docker-compose.yml             # 本地 PostgreSQL 15 + Redis 7
├── .env.example                   # 环境变量模板
└── README.md
```

---

## 2. 数据库物理模型设计 (Schema)

AI 请严格按照以下 PostgreSQL 建表语句生成 Model/Entity 与迁移文件。

### 2.1 核心用户与导师

```sql
-- 1. 用户表 (家长与学员档案)
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    openid          VARCHAR(64) UNIQUE NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    child_name      VARCHAR(50),
    child_age       INT DEFAULT 6 CHECK (child_age BETWEEN 6 AND 14),
    english_level   VARCHAR(20) DEFAULT 'beginner', -- beginner | intermediate | advanced
    available_lessons INT DEFAULT 0 CHECK (available_lessons >= 0),
    locked_lessons  INT DEFAULT 0 CHECK (locked_lessons >= 0),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_openid ON users(openid);

-- 2. 导师表 (学霸档案)
CREATE TABLE mentors (
    id              SERIAL PRIMARY KEY,
    openid          VARCHAR(64) UNIQUE NOT NULL,
    real_name       VARCHAR(50) NOT NULL,
    school_name     VARCHAR(100) NOT NULL,
    major           VARCHAR(100) NOT NULL,
    gender          VARCHAR(10) DEFAULT 'unknown', -- male | female | unknown
    english_score   VARCHAR(100),
    intro_video_url VARCHAR(512),
    tags            TEXT[] DEFAULT '{}',
    is_verified     SMALLINT DEFAULT 0, -- 0-待审 1-通过 2-驳回
    balance         NUMERIC(10, 2) DEFAULT 0.00 CHECK (balance >= 0),
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_mentors_openid ON mentors(openid);
CREATE INDEX idx_mentors_verified ON mentors(is_verified);

-- 3. 导师入驻审核材料 (敏感信息仅存后端，前端脱敏展示)
CREATE TABLE mentor_applications (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    id_card_url     VARCHAR(512) NOT NULL,
    student_card_url VARCHAR(512) NOT NULL,
    english_proof_url VARCHAR(512),
    reject_reason   TEXT,
    reviewed_by     INT,
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 2.2 排课与订单

```sql
-- 4. 导师时间槽
CREATE TABLE mentor_slots (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    slot_date       DATE NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    status          SMALLINT DEFAULT 0, -- 0-空闲 1-已预约 2-学霸关闭
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (mentor_id, slot_date, start_time)
);
CREATE INDEX idx_mentor_slots_query ON mentor_slots(mentor_id, slot_date, status);

-- 5. 课程订单
CREATE TABLE course_orders (
    id                  VARCHAR(64) PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(id),
    mentor_id           INT NOT NULL REFERENCES mentors(id),
    slot_id             INT NOT NULL REFERENCES mentor_slots(id),
    status              VARCHAR(20) DEFAULT 'PENDING',
    -- PENDING | RESERVED | ACTIVE | COMPLETED | CANCELLED
    agora_channel_name  VARCHAR(128),
    actual_minutes      INT DEFAULT 0,
    mentor_feedback     TEXT,
    feedback_submitted_at TIMESTAMPTZ,
    cancelled_by        VARCHAR(10), -- user | mentor | system
    cancel_reason       TEXT,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_course_orders_user ON course_orders(user_id, status);
CREATE INDEX idx_course_orders_mentor ON course_orders(mentor_id, status);

-- 6. 成长报告 (课后 24h 内填报)
CREATE TABLE growth_reports (
    id              SERIAL PRIMARY KEY,
    order_id        VARCHAR(64) NOT NULL UNIQUE REFERENCES course_orders(id),
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    user_id         INT NOT NULL REFERENCES users(id),
    speaking_score  SMALLINT CHECK (speaking_score BETWEEN 1 AND 5),
    confidence_score SMALLINT CHECK (confidence_score BETWEEN 1 AND 5),
    vocabulary_score SMALLINT CHECK (vocabulary_score BETWEEN 1 AND 5),
    comment         TEXT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 2.3 资产与支付

```sql
-- 7. 课时包商品
CREATE TABLE lesson_packages (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,
    lesson_count    INT NOT NULL CHECK (lesson_count > 0),
    price_cents     INT NOT NULL CHECK (price_cents > 0),
    is_trial        BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 8. 充值订单 (微信支付)
CREATE TABLE recharge_orders (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id),
    package_id      INT NOT NULL REFERENCES lesson_packages(id),
    amount_cents    INT NOT NULL,
    wx_transaction_id VARCHAR(64),
    status          VARCHAR(20) DEFAULT 'PENDING', -- PENDING | PAID | REFUNDED
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_recharge_orders_user ON recharge_orders(user_id);

-- 9. 钱包流水 (导师余额变动审计)
CREATE TABLE wallet_transactions (
    id              SERIAL PRIMARY KEY,
    mentor_id       INT NOT NULL REFERENCES mentors(id),
    order_id        VARCHAR(64) REFERENCES course_orders(id),
    amount          NUMERIC(10, 2) NOT NULL,
    type            VARCHAR(20) NOT NULL, -- EARN | WITHDRAW | ADJUST
    balance_after   NUMERIC(10, 2) NOT NULL,
    remark          TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_tx_mentor ON wallet_transactions(mentor_id, created_at DESC);

-- 10. 互动课件 (H5 白板绘本)
CREATE TABLE courseware (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(100) NOT NULL,
    cover_url       VARCHAR(512),
    content_url     VARCHAR(512) NOT NULL,
    sort_order      INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 2.4 种子数据

```sql
INSERT INTO lesson_packages (name, lesson_count, price_cents, is_trial) VALUES
  ('体验课 1 节', 1, 9900, TRUE),
  ('标准包 10 节', 10, 89900, FALSE),
  ('进阶包 30 节', 30, 249900, FALSE);
```

---

## 3. 统一 API 规范

### 3.1 基础约定

| 项 | 约定 |
| :--- | :--- |
| Base URL | `/api/v1` |
| 认证 | Header `Authorization: Bearer <jwt>` |
| 成功响应 | `{ "code": 0, "msg": "ok", "data": { ... } }` |
| 错误响应 | `{ "code": <业务码>, "msg": "<可读说明>" }` |
| 分页 | Query: `page=1&page_size=20`，响应 `data.list` + `data.total` |

### 3.2 业务错误码

| code | 含义 |
| :--- | :--- |
| 0 | 成功 |
| 40001 | 参数校验失败 |
| 40101 | 未登录或 Token 失效 |
| 40301 | 无权限 |
| 40401 | 资源不存在 |
| 40002 | 验证码错误或已过期 |
| 40003 | 支付通道未配置 |
| 40901 | 时间槽已被占用 |
| 40902 | 课时余额不足 |
| 50001 | 系统内部错误 |

### 3.3 MVP 阶段 API 清单

#### Auth

| Method | Path | 说明 |
| :--- | :--- | :--- |
| POST | `/auth/sms/send` | 发送短信验证码 |
| POST | `/auth/sms/login` | 手机号 + 验证码登录，返回 JWT |
| POST | `/auth/wx-login` | （可选）微信内网页 OAuth，返回 JWT |
| PUT | `/users/profile` | 更新孩子档案 |

#### 家长端

| Method | Path | 说明 |
| :--- | :--- | :--- |
| GET | `/mentors` | 学霸广场列表（筛选：school, gender, tag） |
| GET | `/mentors/:id` | 学霸详情含 30s 自荐视频 |
| GET | `/mentors/:id/slots` | 某学霸可预约时间周历 |
| POST | `/bookings` | 预约时间槽（Redis 锁 + 扣课时） |
| GET | `/users/lessons` | 课时账户余额 |
| GET | `/packages` | 课时包列表 |
| GET | `/payment/channels` | 可用支付通道 |
| POST | `/recharge` | 发起充值（body: package_id, channel） |
| GET | `/recharge/:id` | 充值订单详情 |
| GET | `/users/lessons` | 课时账户余额 |
| POST | `/rooms/:orderId/join` | 获取 Agora Token 进房 |
| POST | `/rooms/:orderId/heartbeat` | 课中心跳（30s） |
| GET | `/reports/:orderId` | 查看成长报告 |

#### 导师端

| Method | Path | 说明 |
| :--- | :--- | :--- |
| POST | `/mentor/apply` | 提交入驻材料 |
| GET | `/mentor/slots` | 我的时间矩阵 |
| PUT | `/mentor/slots` | 批量设置/释放时间槽 |
| POST | `/mentor/reports` | 提交成长报告 |
| GET | `/mentor/wallet` | 钱包余额与流水 |
| POST | `/mentor/withdraw` | 提现至微信零钱 |

#### 运营后台

| Method | Path | 说明 |
| :--- | :--- | :--- |
| GET | `/admin/mentors/pending` | 待审核学霸列表 |
| PUT | `/admin/mentors/:id/review` | 通过/驳回 |
| CRUD | `/admin/courseware` | 课件管理 |
| GET | `/admin/finance/summary` | 财务大盘 |

#### 系统

| Method | Path | 说明 |
| :--- | :--- | :--- |
| GET | `/health` | 健康检查（DB + Redis） |

---

## 4. 环境变量 (.env)

```bash
# Server
APP_ENV=development
SERVER_PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=topmentor
DB_PASSWORD=topmentor_dev
DB_NAME=topmentor
DB_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=change-me-in-production
JWT_EXPIRE_HOURS=168

# 微信
WX_APP_ID=
WX_APP_SECRET=
WX_MCH_ID=
WX_MCH_API_KEY=

# 声网 Agora
AGORA_APP_ID=
AGORA_APP_CERTIFICATE=

# 阿里云 OSS
OSS_ENDPOINT=
OSS_BUCKET=
OSS_ACCESS_KEY=
OSS_SECRET_KEY=
```

---

## 5. 本地开发启动顺序

1. `docker compose up -d` 启动 PostgreSQL 与 Redis
2. `cd backend && go run ./cmd/migrate` 执行数据库迁移
3. `cd backend && go run ./cmd/server` 启动 API（默认 `:8080`）
4. `curl http://localhost:8080/api/v1/health` 验证服务

---

## 6. MVP 分期交付范围

| 阶段 | 交付物 | 对应 PRD |
| :--- | :--- | :--- |
| **M1 基础骨架** | Monorepo、DB 迁移、健康检查、统一响应 | — |
| **M2 用户与导师** | 微信登录、档案、学霸广场列表 | F-U-01, F-U-02, F-M-01 |
| **M3 预约与资产** | 时间槽、Redis 锁、课时扣减 | F-U-03, F-U-04, F-M-02 |
| **M4 课室** | Agora 进房、心跳、白板课件 | F-U-05 |
| **M5 闭环** | 成长报告、提现、运营后台 | F-U-06, F-M-04~05, F-A-* |

当前迭代目标：**M1 基础骨架**（本文档 §5 所列启动流程可跑通）。
