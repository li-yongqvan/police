#!/usr/bin/env python3
import os, sys, time, paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
M = "docker.m.daocloud.io/library"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

def run(cmd, t=3600):
    print("$", cmd[:120])
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        print(out[-12000:])
    return o.channel.recv_exit_status(), out

patch = f"""set -e
cd {INSTALL}
find services frontend -name Dockerfile 2>/dev/null | while read f; do
  sed -i 's|^FROM golang:|FROM {M}/golang:|g' "$f"
  sed -i 's|^FROM alpine:|FROM {M}/alpine:|g' "$f"
  sed -i 's|^FROM node:|FROM {M}/node:|g' "$f"
  sed -i 's|^FROM nginx:|FROM {M}/nginx:|g' "$f"
done
sed -i 's/\\r$//' .env scripts/*.sh 2>/dev/null || true
grep ^FROM services/user/Dockerfile services/forum/Dockerfile frontend/Dockerfile 2>/dev/null || true
"""
run(patch, 60)

run(f"docker pull {M}/golang:1.22-alpine && docker pull {M}/alpine:3.19 && docker pull {M}/node:20-alpine", 600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} build --pull 2>&1", 3600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d 2>&1", 600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 40 && docker compose {cf} restart nginx", 120)
run(f"cd {INSTALL} && bash scripts/seed-content.sh 2>&1 || true", 180)
_, t = run(f"curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/; curl -fsS http://127.0.0.1:8888/forum-api/api/v1/stats/community 2>/dev/null | head -c 100; echo", 60)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps", 60)
c.close()
print(f"\nhttp://{HOST}:8888/")
