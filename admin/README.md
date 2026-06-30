# TopMentor 运营后台

Vue3 + Pinia + Element Plus PC 管理端。

## 功能

- **财务大盘**：充值、未消耗课时、已完成订单、学霸结算与提现
- **学霸审核**：待审核列表、通过/驳回（证件 URL 脱敏展示）
- **课件管理**：绘本课件 CRUD、上下架

## 启动

```bash
# 确保 backend 已启动且已执行 migrate（含 007_admin_seed.sql）
cd admin
npm install
npm run dev
# → http://localhost:5174
```

默认账号（见 `.env`）：

- 用户名：`admin`
- 密码：`admin123`

## API 前缀

所有请求代理至 `http://localhost:8080/api/v1/admin/*`

参考：`docs/01-PRD.md` §2.3、`docs/03-SDS.md` 运营后台 API
