#!/usr/bin/env bash
# Named Cloudflare Tunnel (stable). Requires Cloudflare account.
# Optional: TUNNEL_HOSTNAME=forum.example.com TUNNEL_ZONE=example.com
set -eu
cd "$(dirname "$0")/.."

TUNNEL_NAME="${TUNNEL_NAME:-ai-forum}"
ORIGIN_URL="${ORIGIN_URL:-http://127.0.0.1:8888}"
CF_DIR="infra/cloudflared"
CONFIG="${CF_DIR}/config.yml"
CREDS="${CF_DIR}/credentials.json"
EXAMPLE="${CF_DIR}/config.yml.example"

install_cloudflared() {
  if command -v cloudflared >/dev/null 2>&1; then return 0; fi
  ARCH="$(uname -m)"
  case "$ARCH" in
    x86_64) PKG=cloudflared-linux-amd64.deb ;;
    aarch64) PKG=cloudflared-linux-arm64.deb ;;
    *) echo "Unsupported arch: $ARCH" >&2; exit 1 ;;
  esac
  curl -fsSL "https://github.com/cloudflare/cloudflared/releases/latest/download/${PKG}" -o "/tmp/${PKG}"
  sudo dpkg -i "/tmp/${PKG}" || sudo apt-get install -f -y
}

install_cloudflared
mkdir -p "$CF_DIR"

if [ ! -f "$CREDS" ]; then
  echo "==> Log in to Cloudflare (browser URL will appear)"
  cloudflared tunnel login

  echo "==> Create tunnel: $TUNNEL_NAME"
  cloudflared tunnel create "$TUNNEL_NAME" || true

  CREDS_SRC="$(find "$HOME/.cloudflared" -maxdepth 1 -type f -name '*.json' 2>/dev/null | head -1)"
  if [ ! -f "$CREDS_SRC" ]; then
    CREDS_SRC="$HOME/.cloudflared/${TUNNEL_NAME}.json"
  fi
  if [ ! -f "$CREDS_SRC" ]; then
    echo "credentials not found under ~/.cloudflared — copy tunnel JSON to $CREDS" >&2
    exit 1
  fi
  cp "$CREDS_SRC" "$CREDS"
  chmod 600 "$CREDS"
fi

TUNNEL_ID="$(cloudflared tunnel list 2>/dev/null | awk -v n="$TUNNEL_NAME" '$0 ~ n {print $1; exit}')"
if [ -z "$TUNNEL_ID" ]; then
  TUNNEL_ID="$(python3 -c "import json; print(json.load(open('$CREDS'))['TunnelID'])" 2>/dev/null || true)"
fi

if [ -n "${TUNNEL_HOSTNAME:-}" ] && [ -n "${TUNNEL_ZONE:-}" ]; then
  echo "==> DNS route: $TUNNEL_HOSTNAME"
  cloudflared tunnel route dns "$TUNNEL_NAME" "$TUNNEL_HOSTNAME" || \
    cloudflared tunnel route dns --overwrite-dns "$TUNNEL_ID" "$TUNNEL_HOSTNAME"
fi

if [ ! -f "$CONFIG" ]; then
  if [ -n "${TUNNEL_HOSTNAME:-}" ]; then
    sed "s/forum.example.com/${TUNNEL_HOSTNAME}/g" "$EXAMPLE" >"$CONFIG"
  else
    cat >"$CONFIG" <<EOF
tunnel: ${TUNNEL_NAME}
credentials-file: /etc/cloudflared/credentials.json
ingress:
  - service: ${ORIGIN_URL}
  - service: http_status:404
EOF
    echo "Wrote catch-all ingress (no hostname). Add TUNNEL_HOSTNAME in Zero Trust dashboard or re-run with TUNNEL_HOSTNAME=..."
  fi
fi

COMPOSE_FILES="-f docker-compose.yml"
[ -f docker-compose.server.yml ] && COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.server.yml"
[ -f docker-compose.host.yml ] && COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.host.yml"
COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.cloudflare.yml"

echo "==> Start named tunnel container"
docker compose $COMPOSE_FILES up -d cloudflared

echo "Done. Configure public hostname in Cloudflare Zero Trust if not using TUNNEL_HOSTNAME."
echo "Docs: docs/cloudflare-tunnel-setup.md"
