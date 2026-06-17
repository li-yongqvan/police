#!/usr/bin/env bash
# Ephemeral Quick Tunnel (no domain). URL changes when process restarts.
set -eu
ORIGIN_URL="${ORIGIN_URL:-http://127.0.0.1:8888}"
LOG="${CLOUDFLARED_QUICK_LOG:-/var/log/cloudflared-quick.log}"
PIDFILE="${CLOUDFLARED_QUICK_PID:-/var/run/cloudflared-quick.pid}"

if ! command -v cloudflared >/dev/null 2>&1; then
  echo "Installing cloudflared..."
  ARCH="$(uname -m)"
  case "$ARCH" in
    x86_64) PKG=cloudflared-linux-amd64.deb ;;
    aarch64) PKG=cloudflared-linux-arm64.deb ;;
    *) echo "Unsupported arch: $ARCH" >&2; exit 1 ;;
  esac
  curl -fsSL "https://github.com/cloudflare/cloudflared/releases/latest/download/${PKG}" -o "/tmp/${PKG}"
  sudo dpkg -i "/tmp/${PKG}" || sudo apt-get install -f -y
fi

if [ -f "$PIDFILE" ] && kill -0 "$(cat "$PIDFILE")" 2>/dev/null; then
  echo "Quick tunnel already running (pid $(cat "$PIDFILE"))"
else
  : >"$LOG"
  nohup cloudflared tunnel --url "$ORIGIN_URL" >>"$LOG" 2>&1 &
  echo $! | sudo tee "$PIDFILE" >/dev/null
  echo "Started cloudflared quick tunnel (pid $(cat "$PIDFILE"))"
fi

echo "Waiting for public URL..."
for _ in $(seq 1 30); do
  URL="$(grep -oE 'https://[a-zA-Z0-9-]+\.trycloudflare\.com' "$LOG" | head -1 || true)"
  if [ -n "$URL" ]; then
    echo ""
    echo "==> Mobile / mainland access URL:"
    echo "$URL"
    echo ""
    echo "Log: $LOG"
    exit 0
  fi
  sleep 1
done

echo "URL not found yet. Check: tail -f $LOG" >&2
exit 1
