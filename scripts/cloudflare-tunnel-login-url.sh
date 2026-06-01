#!/usr/bin/env bash
# Print Cloudflare login URL for headless server (open on your PC browser).
set -eu
LOG=/opt/ai-forum/CF_LOGIN_URL.txt
: >"$LOG"
echo "Starting cloudflared tunnel login — URL will be saved to $LOG"
timeout 120 cloudflared tunnel login 2>&1 | tee -a "$LOG" || true
grep -oE 'https://[^ ]+' "$LOG" | grep -E 'dash.cloudflare.com|cloudflare.com' | head -3
