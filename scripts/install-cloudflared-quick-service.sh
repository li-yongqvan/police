#!/usr/bin/env bash
# systemd unit for Quick Tunnel (no domain). Survives reboot; URL may change each start.
set -eu
SCRIPT="/opt/ai-forum/scripts/run-cloudflare-quick-tunnel.sh"
UNIT="/etc/systemd/system/cloudflared-quick.service"

sudo tee "$UNIT" >/dev/null <<EOF
[Unit]
Description=Cloudflare Quick Tunnel for AI Forum
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/bin/bash ${SCRIPT}
ExecStop=/bin/bash -c 'kill \$(cat /var/run/cloudflared-quick.pid 2>/dev/null) 2>/dev/null || true'
TimeoutStartSec=120

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable cloudflared-quick.service
echo "Enabled cloudflared-quick.service — URL in /var/log/cloudflared-quick.log"
echo "After reboot: sudo systemctl start cloudflared-quick && grep trycloudflare /var/log/cloudflared-quick.log"
