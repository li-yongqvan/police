#!/usr/bin/env bash
# Diagnose Cloudflare Tunnel + local gateway for https://api.shgenren.dpdns.org (522/timeout).
set -eu
cd "$(dirname "$0")/.."

PUBLIC_HOST="${TUNNEL_HOSTNAME:-api.shgenren.dpdns.org}"
ORIGIN="${ORIGIN_URL:-http://127.0.0.1:8888}"
ORIGIN="${ORIGIN%/}"
FAIL=0

say() { printf '%s\n' "$*"; }
ok() { say "  [OK] $*"; }
bad() { say "  [FAIL] $*"; FAIL=1; }
warn() { say "  [WARN] $*"; }

say "=== Tunnel health $(date -Is) ==="
say "Public: https://${PUBLIC_HOST}"
say "Origin: ${ORIGIN}"
say ""

say "1) Local gateway"
if curl -fsS -o /dev/null --max-time 5 "${ORIGIN}/"; then
  ok "GET ${ORIGIN}/"
else
  bad "GET ${ORIGIN}/ — start stack: docker compose -f docker-compose.yml -f docker-compose.server.yml up -d && restart nginx"
fi
if curl -fsS --max-time 5 "${ORIGIN}/forum-api/api/v1/stats/community" >/dev/null 2>&1; then
  ok "forum-api stats"
else
  bad "forum-api — docker compose ... restart nginx"
fi
say ""

say "2) cloudflared process"
if systemctl is-active cloudflared-named >/dev/null 2>&1; then
  ok "systemd cloudflared-named is active"
  systemctl status cloudflared-named --no-pager 2>/dev/null | sed -n '1,8p' | sed 's/^/    /' || true
elif docker ps --format '{{.Names}}' 2>/dev/null | grep -q cloudflared; then
  ok "docker cloudflared container running"
  docker ps --filter name=cloudflared --format '    {{.Names}} {{.Status}}' 2>/dev/null || true
elif pgrep -x cloudflared >/dev/null 2>&1; then
  ok "cloudflared process (non-systemd)"
else
  bad "no cloudflared — run: bash scripts/restore-named-tunnel.sh  (needs CLOUDFLARE_TUNNEL_TOKEN)"
  say "    or: docker compose ... -f docker-compose.cloudflare.yml up -d cloudflared"
fi
say ""

say "3) Public URL (from this host)"
HTTP_CODE="$(curl -sS -o /dev/null -w '%{http_code}' --max-time 20 "https://${PUBLIC_HOST}/" 2>/dev/null || echo 000)"
if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "304" ]; then
  ok "https://${PUBLIC_HOST}/ → HTTP ${HTTP_CODE}"
elif [ "$HTTP_CODE" = "000" ]; then
  bad "https://${PUBLIC_HOST}/ timed out (Cloudflare 522) — tunnel not connected to origin"
  say "    Zero Trust → Networks → Tunnels → your tunnel → Public Hostname:"
  say "      ${PUBLIC_HOST} → http://127.0.0.1:8888"
  say "    Tunnel status must show Healthy / Connected."
else
  warn "https://${PUBLIC_HOST}/ → HTTP ${HTTP_CODE} (check CF dashboard / SSL mode Full)"
fi
say ""

say "4) Quick tunnel fallback (temporary URL for students)"
if [ -f /var/log/cloudflared-quick.log ]; then
  URL="$(grep -oE 'https://[a-zA-Z0-9-]+\.trycloudflare\.com' /var/log/cloudflared-quick.log 2>/dev/null | tail -1 || true)"
  if [ -n "$URL" ]; then
    warn "Quick tunnel URL (expires on restart): $URL"
  fi
fi
if [ -f PUBLIC_TUNNEL_URL.txt ]; then
  warn "PUBLIC_TUNNEL_URL.txt: $(cat PUBLIC_TUNNEL_URL.txt)"
fi
say ""

if [ "$FAIL" -eq 0 ]; then
  say "=== All local checks passed ==="
else
  say "=== Fix failures above, then re-run: bash scripts/check-tunnel-health.sh ==="
  exit 1
fi
