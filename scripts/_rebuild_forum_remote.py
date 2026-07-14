#!/usr/bin/env python3
"""Rebuild forum-service on domestic VPS with latest liked-state API."""
import os
import sys
from pathlib import Path

import paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
ROOT = Path(__file__).resolve().parents[1]
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

FILES = [
    "services/forum/Dockerfile",
    "services/forum/cmd/main.go",
    "services/forum/internal/middleware/auth.go",
    "services/forum/internal/handler/notification_handler.go",
    "services/forum/internal/handler/post_handler.go",
    "services/forum/internal/service/extras_service.go",
    "services/forum/internal/service/forum_service.go",
    "services/forum/internal/model/post.go",
]


def run(c, cmd, t=3600):
    print(f"$ {cmd[:100]}")
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    sys.stdout.buffer.write(out[-10000:].encode("utf-8", "replace"))
    sys.stdout.buffer.write(b"\n")


def main():
    if not PW:
        raise SystemExit("Set DEPLOY_PASSWORD")
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=PW, timeout=60)
    sftp = c.open_sftp()
    for rel in FILES:
        sftp.put(str(ROOT / rel), f"{INSTALL}/{rel}")
    sftp.close()
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} build forum-service --no-cache 2>&1 | tail -30",
        3600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} up -d --force-recreate forum-service",
        300,
    )
    run(
        c,
        'sleep 4 && curl -sS "http://127.0.0.1:8091/forum-api/api/v1/posts?sort=hot&limit=1" | head -c 280',
        30,
    )
    c.close()
    print(f"\nForum API live on http://{HOST}:8091/")


if __name__ == "__main__":
    main()
