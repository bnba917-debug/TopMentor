#!/usr/bin/env bash
# 从 IP 测试环境切换到域名 HTTPS 部署
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> 停止测试环境（若曾运行 deploy-staging.sh）..."
docker compose -f docker-compose.prod.yml -f docker-compose.staging.yml down 2>/dev/null || true

if [[ ! -f .env ]]; then
  echo "错误: 未找到 .env，请先 cp deploy/env.production.example .env"
  exit 1
fi

# shellcheck disable=SC1091
source <(grep -E '^(DOMAIN|ADMIN_DOMAIN|SITE_URL)=' .env | sed 's/^/export /')

DOMAIN="${DOMAIN:-toppeertalk.com}"
ADMIN_DOMAIN="${ADMIN_DOMAIN:-admin.toppeertalk.com}"

echo "==> 域名配置: H5=https://${DOMAIN}  后台=https://${ADMIN_DOMAIN}"
echo "    请确认 DNS 已解析到本机，且安全组已开放 80/443"
echo ""

bash "$ROOT/deploy/deploy.sh"

echo ""
echo "验证命令:"
echo "  curl -I https://${DOMAIN}"
echo "  curl -I https://${ADMIN_DOMAIN}"
