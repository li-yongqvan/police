#!/usr/bin/env bash
# Named Tunnel via dashboard token (no browser on server).
# 1) Zero Trust → Networks → Tunnels → Create → Docker → copy token
# 2) export CLOUDFLARE_TUNNEL_TOKEN='eyJh...'
# 3) bash scripts/setup-named-tunnel-token.sh
set -eu
cd "$(dirname "$0")/.."

if [ -z "${CLOUDFLARE_TUNNEL_TOKEN:-}" ]; then
  echo "Set CLOUDFLARE_TUNNEL_TOKEN from Cloudflare Tunnel install command." >&2
  exit 1
fi

# Stop quick tunnel if running
if [ -f /var/run/cloudflared-quick.pid ]; then
  kill "$(cat /var/run/cloudflared-quick.pid)" 2>/dev/null || true
  rm -f /var/run/cloudflared-quick.pid
fi
systemctl stop cloudflared-quick 2>/dev/null || true

mkdir -p /etc/cloudflared
echo "$CLOUDFLARE_TUNNEL_TOKEN" > /etc/cloudflared/tunnel.token
chmod 600 /etc/cloudflared/tunnel.token

UNIT=/etc/systemd/system/cloudflared-named.service
cat >"$UNIT" <<'EOF'
[Unit]
Description=Cloudflare Named Tunnel (AI Forum)
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=5
EnvironmentFile=-/etc/cloudflared/env
ExecStart=/usr/bin/cloudflared tunnel --no-autoupdate run --token ${CLOUDFLARE_TUNNEL_TOKEN}

[Install]
WantedBy=multi-user.target
EOF

# systemd does not expand env from file in ExecStart easily — use wrapper
WRAPPER=/usr/local/bin/cloudflared-named-run.sh
cat >"$WRAPPER" <<'WRAP'
#!/bin/bash
set -a
source /etc/cloudflared/tunnel.token.env 2>/dev/null || true
set +a
exec /usr/bin/cloudflared tunnel --no-autoupdate run --token "$CLOUDFLARE_TUNNEL_TOKEN"
WRAP
chmod +x "$WRAPPER"

echo "CLOUDFLARE_TUNNEL_TOKEN=$(cat /etc/cloudflared/tunnel.token)" > /etc/cloudflared/tunnel.token.env
chmod 600 /etc/cloudflared/tunnel.token.env

sed -i 's|ExecStart=.*|ExecStart=/usr/local/bin/cloudflared-named-run.sh|' "$UNIT"

systemctl daemon-reload
systemctl enable cloudflared-named.service
systemctl restart cloudflared-named.service
sleep 3
systemctl status cloudflared-named.service --no-pager | head -15
echo "Configure Public Hostname api.shgenren.dpdns.org → http://127.0.0.1:8888 in Zero Trust dashboard."
