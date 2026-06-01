#!/usr/bin/env python3
import os, sys, paramiko
HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

def run(cmd, t=3600):
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    raw = o.read() + e.read()
    text = raw.decode("utf-8", "replace").encode("ascii", "replace").decode("ascii")
    print(text[-5000:])
    return text

print("=== BUILD ===")
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} build 2>&1", 3600)
print("=== UP ===")
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d 2>&1", 600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 40 && docker compose {cf} restart nginx", 120)
run(f"cd {INSTALL} && bash scripts/seed-content.sh 2>&1", 180)
t = run(f"curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}' http://127.0.0.1:8888/; echo; cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps", 60)
c.close()
sys.exit(0 if "HTTP8888=200" in t else 1)
