# 05-Payment.md - M3 支付方案

## 1. 设计目标

个人项目无小程序、短期无商户号时，用 **Mock 充值**跑通业务；办个体户后只改 `.env` 即可切换真实支付，**无需改业务代码**。

## 2. 支付通道

| channel | 场景 | 配置项 |
| :--- | :--- | :--- |
| `mock` | 开发 / 内测，立即到账 | `PAYMENT_MODE=mock` |
| `wechat_jsapi` | 微信内置浏览器 | `WX_APP_ID` + `WX_MCH_ID` + `WX_MCH_API_KEY` |
| `wechat_h5` | 手机 Safari/Chrome | 同上 |
| `alipay_h5` | 支付宝浏览器 | `ALIPAY_APP_ID` + 密钥 |

## 3. 环境变量

```bash
PAYMENT_MODE=mock          # mock | live
SITE_URL=https://m.example.com

# live 模式 — 微信
WX_APP_ID=
WX_MCH_ID=
WX_MCH_API_KEY=
WX_NOTIFY_URL=https://your-domain.com/api/v1/payment/notify/wechat

# live 模式 — 支付宝
ALIPAY_APP_ID=
ALIPAY_PRIVATE_KEY=
ALIPAY_PUBLIC_KEY=
ALIPAY_NOTIFY_URL=
```

**切换步骤：** `PAYMENT_MODE=live` → 填商户密钥 → 重启 backend。

> `PAYMENT_MODE=mock` 时，任意 channel 请求都会走 Mock 立即到账（便于联调）。

## 4. API

| Method | Path | 鉴权 | 说明 |
| :--- | :--- | :--- | :--- |
| GET | `/packages` | 否 | 课时包列表 |
| GET | `/payment/channels` | 否 | 当前可用通道与模式 |
| POST | `/recharge` | 是 | 发起充值 |
| GET | `/recharge/:id` | 是 | 充值订单状态 |
| GET | `/users/lessons` | 是 | 课时余额 |
| POST | `/payment/notify/wechat` | 否 | 微信异步回调（待接入） |

### POST /recharge

```json
{
  "package_id": 1,
  "channel": "mock"
}
```

**Mock 成功响应：**

```json
{
  "code": 0,
  "data": {
    "order_id": "Rabc...",
    "status": "PAID",
    "mock_paid": true,
    "available_lessons": 1,
    "lesson_count": 1,
    "amount_cents": 9900,
    "channel": "mock"
  }
}
```

**Live 待支付响应：**

```json
{
  "code": 0,
  "data": {
    "order_id": "Rabc...",
    "status": "PENDING",
    "pay_url": "https://...",
    "jsapi_params": { "appId": "...", "paySign": "..." }
  }
}
```

## 5. 代码结构

```
backend/pkg/payment/
├── provider.go      # 工厂 + 接口
├── mock.go          # Mock 立即到账
├── wechat_jsapi.go  # 微信 JSAPI（占位，填密钥后实现 TODO）
├── wechat_h5.go     # 微信 H5（占位）
└── alipay_h5.go     # 支付宝 H5（占位）
```

## 6. 事务保证

Mock 与回调完成支付均走 `CompleteAndCreditLessons`：

1. `recharge_orders` 行锁，防重复支付
2. `users.available_lessons` 行锁后增加课时
3. 同一事务提交

## 7. 错误码

| code | 含义 |
| :--- | :--- |
| 40003 | 支付通道未配置 |
| 50101 | 回调接口待接入 |

---

## 8. 个体户进件 Checklist（TopMentor 适用）

个人项目要接 **live 支付** 前，按此清单自检（可打印给家人一起看）。

### 8.1 执照是否「能用」

| 检查项 | 要求 | 自查 |
| :--- | :--- | :--- |
| 主体类型 | 个体工商户或公司 | ☐ |
| 执照状态 | 在营、未异常 | ☐ |
| **经营范围** | 含以下至少一类：教育咨询、信息咨询、文化艺术交流、互联网信息服务、软件/技术服务 | ☐ |
| 仅「家电销售/维修」 | **不适合** TopMentor，需变更范围或新办 | ☐ |

> 经营范围原文可在营业执照或「国家企业信用信息公示系统」查询。

### 8.2 微信支付（建议优先）

| 步骤 | 说明 |
| :--- | :--- |
| 1 | 注册 [微信支付商户平台](https://pay.weixin.qq.com/) |
| 2 | 提交个体户执照、法人身份证、银行账户 |
| 3 | 选择类目：教育培训 / 教育咨询 / 在线课程 等（以审核页为准） |
| 4 | （微信内 JSAPI）认证**服务号**并与商户号绑定 |
| 5 | 配置支付授权目录、JS 接口安全域名（H5 域名） |
| 6 | 拿到 `WX_MCH_ID`、`WX_MCH_API_KEY`，AppID 填 `WX_APP_ID` |
| 7 | 设置回调 URL → 本项目 `WX_NOTIFY_URL=/api/v1/payment/notify/wechat` |
| 8 | `.env` 设 `PAYMENT_MODE=live`，重启 backend |

### 8.3 支付宝（可选补充）

| 步骤 | 说明 |
| :--- | :--- |
| 1 | [支付宝开放平台](https://open.alipay.com/) 创建应用 |
| 2 | 签约「手机网站支付」 |
| 3 | 配置 `ALIPAY_APP_ID`、应用私钥、支付宝公钥 |
| 4 | 回调 → `ALIPAY_NOTIFY_URL` |

### 8.4 用「家人执照」前必须说清的事

| 问题 | 说明 |
| :--- | :--- |
| 钱进谁账户 | 进该个体户对公/法人账户 |
| 谁报税 | 该执照主体 |
| 封号风险 | 经营范围与 TopMentor 不符可能被拒或冻结 |
| 建议 | 长期运营最好**自己办**含咨询/教育范围的个体户 |

### 8.5 与本项目的配置对应

```bash
PAYMENT_MODE=live
SITE_URL=https://m.你的域名.com
WX_APP_ID=服务号AppID
WX_MCH_ID=商户号
WX_MCH_API_KEY=APIv2或v3密钥
WX_NOTIFY_URL=https://你的域名.com/api/v1/payment/notify/wechat
CORS_ORIGINS=https://m.你的域名.com
```

H5 在微信内打开 → 前端 `channel` 选 `wechat_jsapi`；外部浏览器 → `wechat_h5` 或 `alipay_h5`。

### 8.6 进件完成后的代码工作

当前 `pkg/payment/wechat_*.go` 为占位，有商户号后需补：

1. 统一下单 API（金额、订单号、回调 URL）
2. JSAPI 签名参数 / H5 `mweb_url`
3. 回调验签 + 调用 `CompleteAndCreditLessons`

**Mock 充值与预约扣课时逻辑无需改动。**
