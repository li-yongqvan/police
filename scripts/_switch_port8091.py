#!/usr/bin/env python3
import os, paramiko
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = "/home/liyongquan/projects/ai-forum"
PORT = "8091"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=PW, timeout=60)

def run(cmd, t=300):
    print("$", cmd[:95])
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    print(o.read().decode("utf-8", "replace")[-3000:])

run(f"sed -i 's/\"8888:80\"/\"{PORT}:80\"/' {INSTALL}/docker-compose.server.yml && grep ports {INSTALL}/docker-compose.server.yml", 30)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d nginx", 120)
run(f"curl -sS -o /dev/null -w 'local{PORT}=%{{http_code}}\\n' http://127.0.0.1:{PORT}/", 30)
c.close()
print(f"Try http://122.51.233.225:{PORT}/")
