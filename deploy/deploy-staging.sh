#!/usr/bin/env bash
# TopMentor 服务器测试部署（HTTP + IP，无需域名证书）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export COMPOSE_FILE="docker-compose.prod.yml:docker-compose.staging.yml"
export NODE_IMAGE="${NODE_IMAGE:-node:20-alpine}"

if [[ ! -f .env ]]; then
  echo "错误: 未找到 .env"
  echo "请先执行: cp deploy/env.production.example .env"
  exit 1
fi

# 测试模式建议 SITE_URL 用 IP
if grep -q '^SITE_URL=https://toppeertalk.com' .env 2>/dev/null; then
  echo "提示: 测试阶段可将 .env 中 SITE_URL 改为 http://你的服务器IP"
fi

bash "$ROOT/deploy/deploy.sh"

SERVER_IP="${SERVER_IP:-$(curl -s --max-time 3 ifconfig.me 2>/dev/null || curl -s --max-time 3 ip.sb 2>/dev/null || echo '你的服务器IP')}"

echo ""
echo "=== 测试访问地址（HTTP）==="
echo "  H5:    http://${SERVER_IP}/welcome"
echo "  后台:  http://${SERVER_IP}:8081/"
echo ""
echo "登录验证码仍为 mock 模式下的 123456（见 .env SMS_MOCK_CODE）"
