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
    print("$", cmd[:100])
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        print(out[-8000:])
    return o.channel.recv_exit_status(), out

# Ensure GOPROXY for go mod download inside build
run(f"""cd {INSTALL} && for f in services/*/Dockerfile; do
  grep -q GOPROXY "$f" || sed -i '/^WORKDIR \\/app/a ENV GOPROXY=https://goproxy.cn,direct' "$f"
done""", 30)

run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} build --no-cache 2>&1 | tail -40", 3600)
code, _ = run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d 2>&1", 600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 40 && docker compose {cf} ps", 120)
_, t = run("curl -fsS -o /dev/null -w 'HTTP8888=%{http_code}\n' http://127.0.0.1:8888/ 2>&1", 30)
c.close()
ok = "HTTP8888=200" in t
print("RESULT:", "OK" if ok else "FAILED")
sys.exit(0 if ok else 1)
