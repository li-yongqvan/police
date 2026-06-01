#!/usr/bin/env python3
"""Replace broken migration tree with full pilot-deploy.tgz and deploy."""
import os, secrets, sys, time, paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TGZ = os.path.join(ROOT, "pilot-deploy.tgz")
M = "docker.m.daocloud.io/library"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
jwt, pg = secrets.token_hex(32), secrets.token_hex(16)

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

def run(cmd, t=3600):
    print("$", cmd[:110])
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        print(out[-6000:])
    return out

print("Upload pilot-deploy.tgz...")
sftp = c.open_sftp()
sftp.put(TGZ, "/tmp/pilot-deploy.tgz")
sftp.close()

run(f"mkdir -p {INSTALL} && tar -xzf /tmp/pilot-deploy.tgz -C {INSTALL}")

patch = f"""cd {INSTALL}
cp .env.production.example .env
sed -i 's/^JWT_SECRET=.*/JWT_SECRET={jwt}/' .env
sed -i 's/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD={pg}/' .env
sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD={pg}/' .env
sed -i 's/^DB_HOST=.*/DB_HOST=postgres/' .env
sed -i 's/^REDIS_HOST=.*/REDIS_HOST=redis/' .env
sed -i 's/\\r$//' .env scripts/*.sh 2>/dev/null || true
find services frontend -name Dockerfile | while read f; do
  sed -i 's|^FROM golang:|FROM {M}/golang:|g' "$f"
  sed -i 's|^FROM alpine:|FROM {M}/alpine:|g' "$f"
  sed -i 's|^FROM node:|FROM {M}/node:|g' "$f"
  sed -i 's|^FROM nginx:|FROM {M}/nginx:|g' "$f"
  grep -q GOPROXY "$f" || sed -i '/^WORKDIR \\/app/a ENV GOPROXY=https://goproxy.cn,direct' "$f"
done
chmod +x scripts/*.sh
"""
run(patch, 60)

run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} build 2>&1 | tail -25", 3600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d 2>&1", 300)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 45 && docker compose {cf} restart nginx", 120)
run(f"cd {INSTALL} && bash scripts/seed-content.sh 2>&1 | tail -15", 180)
out = run(f"curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/; cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps", 60)
c.close()
print(f"\nhttp://{HOST}:8888/")
print("OK" if "HTTP8888=200" in out else "Check logs")
