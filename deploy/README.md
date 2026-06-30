# TopMentor Linux 生产部署（域名 HTTPS）

主域名：`toppeertalk.com`（HTTPS 由 Caddy + Let's Encrypt **自动申请并续期**）。

## 架构

```
Internet :443/:80
    └── Caddy（自动 HTTPS）
            ├── toppeertalk.com      → web/dist（H5）
            │     /api/*、/uploads/* → api:8080
            └── admin.toppeertalk.com → admin/dist
                  /api/*             → api:8080
```

## 前置条件

1. Linux 服务器 + Docker + Docker Compose v2
2. **DNS 已解析**到服务器公网 IP：

| 类型 | 主机记录 | 值 |
|------|----------|-----|
| A | `@` | 服务器 IP |
| A | `www` | 服务器 IP |
| A | `admin` | 服务器 IP |

3. 安全组 / 防火墙开放 **80**、**443**

## 一键部署（域名 HTTPS）

```bash
git clone <your-repo> TopMentor
cd TopMentor

cp deploy/env.production.example .env
# 或直接使用已写好的配置：cp deploy/env.server .env
nano .env
# 使用 env.server 时密码已生成，可按需修改

chmod +x deploy/deploy.sh
./deploy/deploy.sh
```

## 从 IP 测试环境切换过来

若之前跑过 `./deploy/deploy-staging.sh`：

```bash
chmod +x deploy/switch-to-domain.sh
./deploy/switch-to-domain.sh
```

或手动：

```bash
docker compose -f docker-compose.prod.yml -f docker-compose.staging.yml down

# .env 确保为：
# DOMAIN=toppeertalk.com
# SITE_URL=https://toppeertalk.com
# CORS_ORIGINS=https://toppeertalk.com,https://www.toppeertalk.com,https://admin.toppeertalk.com

./deploy/deploy.sh
```

## 验证

```bash
curl -I https://toppeertalk.com
curl -I https://admin.toppeertalk.com
docker compose -f docker-compose.prod.yml exec api wget -qO- http://127.0.0.1:8080/api/v1/health
```

浏览器：

- H5：`https://toppeertalk.com/welcome`
- 后台：`https://admin.toppeertalk.com`

登录：手机号 + 验证码 `123456`（mock）；后台 `admin` / `.env` 中 `ADMIN_PASSWORD`。

## 修改域名

编辑 `.env` 中 `DOMAIN`、`ADMIN_DOMAIN`、`ACME_EMAIL`、`SITE_URL`、`CORS_ORIGINS`，然后：

```bash
docker compose -f docker-compose.prod.yml up -d --force-recreate caddy
```

## 更新发布

```bash
git pull
./deploy/deploy.sh
```

## 常用运维

```bash
docker compose -f docker-compose.prod.yml logs -f
docker compose -f docker-compose.prod.yml logs -f caddy   # 证书问题
docker compose -f docker-compose.prod.yml down
```

## IP 测试模式（可选）

域名未就绪时见 `deploy/deploy-staging.sh`（HTTP + IP，端口 80/8081）。

## 本地开发

本地用根目录 `docker-compose.yml`（仅 Postgres + Redis），不要在本机占用 80/443。
