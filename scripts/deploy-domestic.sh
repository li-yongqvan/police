#!/usr/bin/env bash
# Domestic deploy without sudo / Cloudflare. Binds gateway on :8888 (host :80 often taken).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
INSTALL_DIR="${INSTALL_DIR:-$HOME/projects/ai-forum}"

compose() {
  if docker compose version >/dev/null 2>&1; then
    docker compose "$@"
  elif command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    echo "Docker Compose not found" >&2
    exit 1
  fi
}

install_compose_plugin() {
  if docker compose version >/dev/null 2>&1; then
    return 0
  fi
  echo "==> Installing Docker Compose plugin to ~/.docker/cli-plugins"
  mkdir -p "$HOME/.docker/cli-plugins"
  ARCH="$(uname -m)"
  case "$ARCH" in
    x86_64) BIN=docker-compose-linux-x86_64 ;;
    aarch64|arm64) BIN=docker-compose-linux-aarch64 ;;
    *) echo "Unsupported arch: $ARCH" >&2; exit 1 ;;
  esac
  VER=v2.32.4
  DEST="$HOME/.docker/cli-plugins/docker-compose"
  URLS=(
    "https://mirror.ghproxy.com/https://github.com/docker/compose/releases/download/${VER}/${BIN}"
    "https://ghfast.top/https://github.com/docker/compose/releases/download/${VER}/${BIN}"
    "https://github.com/docker/compose/releases/download/${VER}/${BIN}"
  )
  for url in "${URLS[@]}"; do
    echo "Trying $url"
    if curl -fsSL --connect-timeout 20 --max-time 300 "$url" -o "$DEST"; then
      chmod +x "$DEST"
      if docker compose version >/dev/null 2>&1; then
        echo "Compose OK: $(docker compose version)"
        return 0
      fi
    fi
  done
  echo "Failed to install Docker Compose plugin" >&2
  exit 1
}

echo "==> Domestic deploy at $INSTALL_DIR"

install_compose_plugin

if [ "$ROOT" != "$INSTALL_DIR" ]; then
  mkdir -p "$INSTALL_DIR"
  rsync -a --delete \
    --exclude node_modules --exclude .git --exclude pilot-deploy.tgz \
    "$ROOT/" "$INSTALL_DIR/"
  cd "$INSTALL_DIR"
fi

if [ ! -f .env ]; then
  cp .env.production.example .env
  JWT="$(openssl rand -hex 32)"
  PG="$(openssl rand -hex 16)"
  sed -i "s/^JWT_SECRET=.*/JWT_SECRET=$JWT/" .env
  sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$PG/" .env
  sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=$PG/" .env
  echo "Generated .env secrets"
fi

chmod +x scripts/*.sh 2>/dev/null || true

COMPOSE_FILES="-f docker-compose.yml -f docker-compose.server.yml"

echo "==> Build and start (http://<host>:8888)"
compose $COMPOSE_FILES up -d --build

echo "Waiting for migrations..."
sleep 25
compose $COMPOSE_FILES restart nginx || true
sleep 5

bash scripts/seed-content.sh || true

curl -fsS http://127.0.0.1:8888/ >/dev/null && echo "frontend OK"
curl -fsS http://127.0.0.1:8888/forum-api/api/v1/stats/community | head -c 120 && echo

PUBLIC_IP="$(curl -fsS --max-time 3 ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')"
echo ""
echo "=== Deploy complete ==="
echo "  Direct:  http://${PUBLIC_IP}:8888/"
echo "  Local:   http://127.0.0.1:8888/"
echo "  Admin:   demo_platform_admin / demo123456"
echo ""
echo "Port 80 is used by host nginx. Ask ops to proxy / -> 127.0.0.1:8888"
echo "  See: infra/domestic-nginx/ai-forum.conf.example"
