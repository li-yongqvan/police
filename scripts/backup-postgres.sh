#!/usr/bin/env bash
# Daily PostgreSQL backup for AI Forum (run via cron on production host).
set -eu

# Production project root; fall back to repo layout when run from a checkout.
PROJECT_ROOT="${AI_FORUM_PROJECT_ROOT:-/home/liyongquan/projects/ai-forum}"
if [ ! -d "$PROJECT_ROOT" ]; then
  PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
fi
cd "$PROJECT_ROOT"

# cron has a minimal PATH; make sure docker compose plugin is reachable.
export PATH="$HOME/.docker/cli-plugins:/usr/local/bin:/usr/bin:/bin:$PATH"

BACKUP_DIR="${BACKUP_DIR:-/home/liyongquan/backups/db}"
RETAIN_DAYS="${RETAIN_DAYS:-30}"
COMPOSE_FILES="-f docker-compose.yml"
[ -f docker-compose.server.yml ] && COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"

mkdir -p "$BACKUP_DIR"
STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="$BACKUP_DIR/ai_forum_${STAMP}.dump.gz"

echo "==> Dumping database to $OUT"
docker compose $COMPOSE_FILES exec -T postgres \
  pg_dump -U "${POSTGRES_USER:-ai_forum}" "${POSTGRES_DB:-ai_forum}" -Fc | gzip -9 >"$OUT"

find "$BACKUP_DIR" -name 'ai_forum_*.dump.gz' -mtime +"$RETAIN_DAYS" -delete 2>/dev/null || true
echo "==> Backup done ($(du -h "$OUT" | cut -f1))"