#!/usr/bin/env bash
# Append AI Forum cron lines safely (idempotent).
set -eu
MARK_BEGIN="# BEGIN ai-forum pilot"
MARK_END="# END ai-forum pilot"
CRON=/etc/crontab

if grep -qF "$MARK_BEGIN" "$CRON" 2>/dev/null; then
  sed -i "/$MARK_BEGIN/,/$MARK_END/d" "$CRON"
fi

sed -i '/ai-forum\/scripts/d' "$CRON"
sed -i '/1panel-v2\.1\.8/d' "$CRON"
sed -i '/AGENT_DEPLOY/d' "$CRON"
sed -i '/New project 6/d' "$CRON"
sed -i '/trycloudflare/d' "$CRON"

cat >>"$CRON" <<EOF
$MARK_BEGIN
0 3 * * * root /opt/ai-forum/scripts/backup-postgres.sh >> /var/log/ai-forum-backup.log 2>&1
*/15 * * * * root /opt/ai-forum/scripts/cron-smoke.sh
$MARK_END
EOF
echo "Cron installed:"
grep -A3 "$MARK_BEGIN" "$CRON"
