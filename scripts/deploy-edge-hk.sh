#!/usr/bin/env bash
# Deploy edge nginx on Hong Kong (or CN) VPS — run ON THE EDGE SERVER
set -eu

ORIGIN_HOST="${ORIGIN_HOST:-107.172.138.10}"
EDGE_DOMAIN="${EDGE_DOMAIN:-}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONF_SRC="${REPO_ROOT}/infra/edge-nginx/nginx.conf"
CACHE_DIR="/var/cache/nginx/ai-forum"

if [ -z "$EDGE_DOMAIN" ]; then
  echo "Set EDGE_DOMAIN=forum.example.com" >&2
  exit 1
fi

if [ ! -f "$CONF_SRC" ]; then
  echo "Missing $CONF_SRC" >&2
  exit 1
fi

sudo apt-get update -qq
sudo apt-get install -y nginx certbot python3-certbot-nginx
sudo mkdir -p "$CACHE_DIR"

TMP="$(mktemp)"
sed "s/ORIGIN_HOST/${ORIGIN_HOST}/g; s/YOUR_DOMAIN/${EDGE_DOMAIN}/g" "$CONF_SRC" > "$TMP"
sudo cp "$TMP" "/etc/nginx/conf.d/ai-forum-edge.conf"
rm -f "$TMP"

sudo nginx -t
sudo systemctl enable nginx
sudo systemctl reload nginx

echo "==> Request TLS certificate"
sudo certbot --nginx -d "$EDGE_DOMAIN" --non-interactive --agree-tos -m "admin@${EDGE_DOMAIN}" || true

echo "Done. Users should open: https://${EDGE_DOMAIN}/"
echo "Origin: http://${ORIGIN_HOST}:8888 (lock down firewall to edge IP only)"
