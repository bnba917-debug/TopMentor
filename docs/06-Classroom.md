# 06-Classroom.md - M4 课室方案

## 1. 能力范围

| 功能 | 状态 |
| :--- | :--- |
| Agora RTC Token 签发 | ✅ mock / live |
| 进房 `POST /rooms/:orderId/join` | ✅ |
| 心跳 `POST /rooms/:orderId/heartbeat?role=user` | ✅ 30s |
| 下课 `POST /rooms/:orderId/complete` | ✅ |
| 我的预约 `GET /users/orders` | ✅ |
| H5 课室页 `/classroom/:orderId` | ✅ |
| 白板课件 | M4+ 待接 |
| 断线强关 / 按比例退费 | M4+ 待接 |

## 2. 环境变量

```bash
AGORA_APP_ID=              # 声网控制台 App ID
AGORA_APP_CERTIFICATE=     # 证书（启用 App 证书后）
AGORA_MOCK_MODE=true       # 未配置 App ID 时默认 true
LESSON_DURATION_MINUTES=45
```

**切换 live：** 填写 `AGORA_APP_ID` + `AGORA_APP_CERTIFICATE`，设 `AGORA_MOCK_MODE=false`，重启 backend。H5 生产环境须 **HTTPS**。

## 3. API

### 进房

```http
POST /api/v1/rooms/{orderId}/join
Authorization: Bearer <token>
Content-Type: application/json

{"role":"user"}   // user | mentor
```

响应含 `app_id`, `channel`, `token`, `uid`, `mock_mode`, `end_at`（倒计时结束时间）。

### 心跳

```http
POST /api/v1/rooms/{orderId}/heartbeat?role=user
```

### 下课

```http
POST /api/v1/rooms/{orderId}/complete
```

## 4. 联调路径

1. 登录 → 购买课时 → 预约学霸  
2. 「我的」→「我的预约」→「进入课室」  
3. Mock 模式：显示模拟画面 + 倒计时 + 自动心跳  
4. Live 模式：浏览器请求麦克风/摄像头权限后 Agora 进房  

**双端联调（mock）：** 两个浏览器分别登录，同一订单一个 `role=user`、一个 `role=mentor` 进房（mock 下 mentor 不校验身份）。

## 5. 代码结构

```
backend/pkg/agora/       # Token 签发
backend/pkg/roomsession/ # Redis 心跳 / 计时
backend/internal/service/room.go
web/src/views/ClassroomView.vue
web/src/composables/useClassroom.ts
```
