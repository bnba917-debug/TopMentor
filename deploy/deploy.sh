#!/usr/bin/env bash
# TopMentor Linux 域名 HTTPS 部署
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.prod.yml}"
NODE_IMAGE="${NODE_IMAGE:-node:20-alpine}"

if [[ ! -f .env ]]; then
  echo "错误: 未找到 .env"
  echo "请先执行: cp deploy/env.production.example .env"
  echo "然后编辑 .env 中的 DOMAIN、DB_PASSWORD、JWT_SECRET、ADMIN_PASSWORD"
  exit 1
fi

# 读取域名（供输出提示）
DOMAIN="$(grep -E '^DOMAIN=' .env 2>/dev/null | cut -d= -f2- | tr -d '\r' || echo toppeertalk.com)"
ADMIN_DOMAIN="$(grep -E '^ADMIN_DOMAIN=' .env 2>/dev/null | cut -d= -f2- | tr -d '\r' || echo admin.toppeertalk.com)"
DOMAIN="${DOMAIN:-toppeertalk.com}"
ADMIN_DOMAIN="${ADMIN_DOMAIN:-admin.toppeertalk.com}"

if ! command -v docker >/dev/null 2>&1; then
  echo "错误: 未安装 Docker"
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "错误: 需要 Docker Compose v2（docker compose）"
  exit 1
fi

# 若之前跑过 staging，先停掉避免 80 端口冲突
if docker ps --format '{{.Names}}' 2>/dev/null | grep -q topmentor-caddy; then
  echo "==> 检测到已有 Caddy 容器，将重新以域名模式启动..."
fi

echo "==> 构建 H5 前端 (web)..."
docker run --rm \
  -v "$ROOT/web:/app" \
  -w /app \
  "$NODE_IMAGE" \
  sh -c "npm ci && npm run build"

echo "==> 构建运营后台 (admin)..."
docker run --rm \
  -v "$ROOT/admin:/app" \
  -w /app \
  "$NODE_IMAGE" \
  sh -c "npm ci && npm run build"

echo "==> 启动 Docker 服务（域名 HTTPS）..."
docker compose -f "$COMPOSE_FILE" up -d --build

echo "==> 等待 API 健康检查..."
for i in $(seq 1 30); do
  if docker compose -f "$COMPOSE_FILE" exec -T api wget -q -O - http://127.0.0.1:8080/api/v1/health >/dev/null 2>&1; then
    echo "API 已就绪"
    break
  fi
  if [[ "$i" -eq 30 ]]; then
    echo "警告: API 健康检查超时，请查看日志: docker compose -f $COMPOSE_FILE logs api"
    exit 1
  fi
  sleep 2
done

echo ""
echo "部署完成。"
echo "  H5:    https://${DOMAIN}/welcome"
echo "  后台:  https://${ADMIN_DOMAIN}/"
echo ""
echo "首次访问 HTTPS 时 Caddy 会自动向 Let's Encrypt 申请证书。"
echo "若证书失败: docker compose -f $COMPOSE_FILE logs caddy"
echo "  并确认 DNS 已解析、80/443 已开放。"
echo ""
echo "常用命令:"
echo "  查看日志: docker compose -f $COMPOSE_FILE logs -f"
echo "  重启 API: docker compose -f $COMPOSE_FILE restart api"
echo "  停止服务: docker compose -f $COMPOSE_FILE down"
