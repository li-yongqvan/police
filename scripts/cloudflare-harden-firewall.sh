#!/usr/bin/env bash
# Optional: block direct public access to forum port after Tunnel works.
# Keeps SSH (22). Adjust if you use a non-default SSH port.
set -eu

SSH_PORT="${SSH_PORT:-22}"
FORUM_PORT="${FORUM_PORT:-8888}"

if ! command -v ufw >/dev/null 2>&1; then
  echo "ufw not installed. Skip or: apt install ufw" >&2
  exit 1
fi

echo "This will DENY inbound ${FORUM_PORT}/tcp (forum gateway)."
echo "Users must use Cloudflare Tunnel URL only."
if [ "${FORCE:-}" != "1" ]; then
  read -r -p "Continue? [y/N] " ans
  [ "$ans" = "y" ] || [ "$ans" = "Y" ] || exit 0
fi

sudo ufw allow "${SSH_PORT}/tcp"
sudo ufw deny "${FORUM_PORT}/tcp"
sudo ufw --force enable
sudo ufw status

echo "Verify Tunnel URL still works from your phone before closing this session."
