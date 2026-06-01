#!/usr/bin/env bash
# Daily PostgreSQL backup for AI Forum (run via cron on production host).
set -eu
cd /opt/ai-forum 2>/dev/null || cd "$(dirname "$0")/.."

BACKUP_DIR="${BACKUP_DIR:-/opt/backups/ai-forum}"
RETAIN_DAYS="${RETAIN_DAYS:-7}"
COMPOSE_FILES="-f docker-compose.yml"
[ -f docker-compose.server.yml ] && COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"

mkdir -p "$BACKUP_DIR"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="$BACKUP_DIR/ai_forum_${STAMP}.sql.gz"

echo "==> Dumping database to $OUT"
docker compose $COMPOSE_FILES exec -T postgres \
  pg_dump -U "${POSTGRES_USER:-ai_forum}" "${POSTGRES_DB:-ai_forum}" | gzip -9 >"$OUT"

find "$BACKUP_DIR" -name 'ai_forum_*.sql.gz' -mtime +"$RETAIN_DAYS" -delete 2>/dev/null || true
echo "==> Backup done ($(du -h "$OUT" | cut -f1))"
