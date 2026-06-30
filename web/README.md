# TopMentor H5 移动 Web

家长端 H5：**手机浏览器或微信内打开链接即可**，无需小程序审核。

## 技术栈

- Vue 3 + Vite + TypeScript
- Vant 4（移动 UI）
- Pinia + Vue Router
- Axios（API 代理至后端）

## 快速开始

```bash
# 1. 确保后端与 Docker 已启动
cd ../backend && go run ./cmd/server

# 2. 安装并启动 H5
cd web
npm install
npm run dev
```

浏览器访问：**http://localhost:5173**

## 页面

| 路径 | 说明 |
| :--- | :--- |
| `/login` | 手机号 + 短信验证码登录 |
| `/mentors` | 学霸广场（无需登录可浏览） |
| `/mentors/:id` | 学霸详情 |
| `/profile` | 孩子档案（需登录） |

## 开发说明

- 短信验证码开发模式固定 **123456**（`SMS_MOCK_MODE=true`）
- API 通过 Vite 代理转发至 `http://localhost:8080`
- 生产环境须部署 **HTTPS**（后续 WebRTC 课室需要）
