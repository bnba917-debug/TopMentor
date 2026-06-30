# 02-HLD.md - 架构蓝图 (Architecture Blueprint)

## 1. 技术栈选型（已确认）

| 分层 | 技术选型 | 核心理由 |
| :--- | :--- | :--- |
| **用户端 (H5)** | **Vue3 + Vite + TS + Vant** | 移动 Web，浏览器直接访问；家长/学霸路由分入口 |
| **运营后台** | Vue3 + Pinia + Element Plus | 轻量 SPA，适合 CRUD 与审核流 |
| **后端** | **Go 1.18+ + Gin** | 高并发、低延迟，适合 RTC 心跳与时间槽锁 |
| **数据库** | PostgreSQL 15 | 强事务，保障课时与钱包账目 |
| **缓存** | Redis 7 | 时间槽分布式锁、房间状态、会话 |
| **音视频** | **声网 Agora Web SDK** | 浏览器 WebRTC 1v1 + 白板 |
| **对象存储** | 阿里云 OSS | 自荐视频、课件 PDF/图片 |
| **支付** | 微信 H5 支付 / 支付宝 H5 | 课时包充值；企业付款提现 |
| **登录** | 短信验证码 + JWT | 不依赖小程序；微信内可扩展 OAuth |

> Node.js 方案已弃用，统一采用 Go 后端以降低 RTC 场景下的 GC 抖动风险。

---

## 2. 核心架构拓扑

```
                    ┌─────────────────────────────┐
                    │  H5 移动 Web (Vue3 + Vant)   │
                    │  / 家长端    /mentor 学霸端  │
                    │  浏览器 / 微信内链接均可打开  │
                    └────┬───────────────┬────────┘
                         │ HTTPS         │ WebRTC
                         ▼               ▼
┌────────────────────────────────┐   ┌─────────────────────────────┐
│     Go Gin API (/api/v1)       │   │  第三方云基础设施            │
├────────────────────────────────┤   ├─────────────────────────────┤
│ auth      短信登录 / JWT       │   │ Agora Web 浏览器音视频+白板  │
│ user      家长档案 / 课时      │   │ 微信支付   充值 / 企业付款    │
│ mentor    学霸 / 时间槽        │   │ 阿里云 OSS 视频 / 课件       │
│ booking   预约 / Redis 锁      │   └─────────────────────────────┘
│ room      进房 Token / 心跳    │
│ wallet    扣费 / 提现          │   ┌─────────────────────────────┐
│ admin     审核 / 课件 / 财务   │   │ Vue3 Admin (PC 浏览器)       │
└────┬───────────────────┬───────┘   └─────────────────────────────┘
     │                   │
     ▼                   ▼
┌─────────────┐   ┌─────────────┐
│ PostgreSQL  │   │ Redis 7     │
│ 强一致资产   │   │ 锁 / 缓存   │
└─────────────┘   └─────────────┘
```

---

## 3. 后端分层与模块边界

```
handler/   → 解析 HTTP、binding 校验、调用 service、写 response
service/   → 业务规则、事务编排、调用 repository + redis + 第三方 SDK
repository/→ 纯 SQL / GORM 数据访问，不含业务判断
model/     → 表实体、请求/响应 DTO
pkg/       → 可复用工具（response、redis、jwt、wx、agora）
```

**依赖方向：** handler → service → repository → DB，禁止反向依赖。

---

## 4. 核心机制设计

### 4.1 1v1 扣费与心跳 (Billing Heartbeat)

1. **上课前 15 分钟：** `available_lessons -= 1`，`locked_lessons += 1`（事务 + FOR UPDATE）
2. **双端进房：** Agora 回调或 join API 将订单置为 `ACTIVE`
3. **心跳：** 客户端每 30s POST `/rooms/:orderId/heartbeat`
4. **异常挂断：** 连续 3 次缺失心跳 → 调 Agora REST 关房 → 按 `actual_minutes / 45` 结算课时

### 4.2 动态时间槽锁 (Time-Slot Lock)

```
SET mentor:slot:{mentor_id}:{slot_date}:{start_time} {user_id} NX EX 10
```

- 抢到锁 → 进入 PostgreSQL 事务扣课时并写订单
- 10s 内未完成 → 锁自动释放，其他用户可重试

### 4.3 防跳单内容过滤

- IM/留言管道接入敏感词 + 正则（11 位数字、微信/QQ 关键词）
- 命中后替换为 `***` 并写审计日志

### 4.4 取消预约规则

- 距开课 **> 2 小时：** 全额退还 locked/available 课时，无违约金
- 距开课 **≤ 2 小时：** 扣除 1 课时作为违约金，另一方可选补偿策略（M3 实现）

---

## 5. 部署与环境

| 环境 | 说明 |
| :--- | :--- |
| local | docker-compose 启动 PG + Redis，backend 本机 `go run` |
| staging | 单机 Docker 镜像 + 托管 PG/Redis |
| production | K8s 或云 VM 多副本 + 读写分离 PG（后期） |

---

## 6. 可观测性（M1 预留）

- 结构化日志：`slog` JSON 格式，含 `request_id`
- 健康检查：`GET /api/v1/health` 探测 DB ping + Redis ping
- 指标：后续接入 Prometheus `/metrics`
