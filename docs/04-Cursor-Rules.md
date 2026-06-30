# 04-Cursor-Rules.md - AI Coding Rules

当你（Cursor）在本项目生成、重构或修复代码时，**必须严格遵守以下规范**。

## 1. 编程范式与质量控制

1. **类型安全**
   - 前端全面使用 TypeScript，禁止 `any`；API 请求/响应需定义 `interface` 或 `type`
   - 后端 Go 结构体使用 `binding` 标签校验入参；禁止字符串拼接 SQL
2. **错误处理**
   - 禁止未捕获 panic 直达客户端
   - 统一响应格式见 `03-SDS.md` §3.1：`{ "code": 0, "msg": "ok", "data": {} }`
3. **分层原则**
   - `handler → service → repository`，禁止 handler 直接访问 DB

## 2. 前端 H5 规范

1. 视频连线、白板、学霸卡片封装为独立 Vue 组件，单文件不超过 500 行
2. 移动优先：使用 `rem`/`vw` + flex；UI 组件库推荐 **Vant 4**
3. 须兼容 iOS Safari、Android Chrome、微信内置浏览器
4. 调用摄像头/麦克风前检测 HTTPS；课室页集成 **Agora Web SDK**

## 3. 数据库与事务红线

1. 变更 `users.available_lessons`、`mentors.balance` **必须**在事务内且 `SELECT ... FOR UPDATE`
2. 高频字段 `openid`、`mentor_id`、`slot_date` 必须有 B-Tree 索引（见 `03-SDS.md`）

## 4. 执行工作流

收到具体任务时：

1. 阅读 `01-PRD.md` 确认业务逻辑与验收标准
2. 阅读 `03-SDS.md` 确认表结构、API 路径、错误码
3. 阅读 `02-HLD.md` 确认模块边界与机制
4. 按本文件规范编码，完成后输出 **CheckList**（改动文件、环境变量、验证命令）

## 5. 测试要求

- 后端：每个 service 公共方法至少 1 个单元测试；HTTP handler 用 `httptest`
- 前端：关键 composable / 工具函数需 Vitest 覆盖
- 提交前运行：`cd backend && go test ./...`

## 6. 目录速查

| 路径 | 职责 |
| :--- | :--- |
| `backend/internal/handler/` | 路由与入参 |
| `backend/internal/service/` | 业务逻辑 |
| `backend/internal/repository/` | 数据访问 |
| `backend/migrations/` | SQL 迁移 |
| `web/` | H5 移动 Web（家长 + 学霸） |
| `admin/` | 运营后台 |
