#!/usr/bin/env bash
# One-shot campus pilot bootstrap on production host (/opt/ai-forum).
set -eu
cd /opt/ai-forum

COMPOSE_FILES="-f docker-compose.yml"
[ -f docker-compose.server.yml ] && COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"

echo "==> Build & restart application stack"
docker compose $COMPOSE_FILES build user-service forum-service admin-service frontend
docker compose $COMPOSE_FILES up -d user-service forum-service admin-service frontend nginx
sleep 8
docker compose $COMPOSE_FILES restart nginx
sleep 5

echo "==> Seed pilot boards & sample posts (idempotent)"
if [ -f scripts/seed-content.sh ]; then
  bash scripts/seed-content.sh || echo "WARN: seed-content had issues (tables may already exist)"
fi

echo "==> Smoke check (local gateway)"
mkdir -p logs /opt/backups/ai-forum
cat > .env.smoke <<'EOF'
BASE_URL=http://127.0.0.1:8888
SMOKE_USER=demo_student
SMOKE_PASS=demo123456
EOF
if [ -f scripts/smoke-vps-remote.sh ]; then
  BASE_URL=http://127.0.0.1:8888 SMOKE_USER=demo_student SMOKE_PASS=demo123456 bash scripts/smoke-vps-remote.sh
elif [ -f scripts/smoke-test.sh ]; then
  BASE_URL=http://127.0.0.1:8888 SMOKE_USER=demo_student SMOKE_PASS=demo123456 bash scripts/smoke-test.sh
fi

echo "==> Install cron (backup daily 3:00, smoke every 15 min)"
if [ -f scripts/install-pilot-cron.sh ]; then
  bash scripts/install-pilot-cron.sh
fi

echo ""
echo "=============================================="
echo " 校内试运行已就绪"
echo " 学生入口: https://api.shgenren.dpdns.org"
echo " 管理端:   同上，使用 demo_admin / 中台账号登录"
echo " 下一步:"
echo "   1) 中台批量生成邀请码并发放"
echo "   2) 管理端处理待审核与举报"
echo "   3) 见 docs/pilot-launch-guide.md"
echo "=============================================="
