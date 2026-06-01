#!/usr/bin/env python3
import os, paramiko
HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

cmds = [
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps -a 2>&1",
    f"tail -30 {INSTALL}/deploy.log 2>/dev/null || true",
    f"tail -30 {INSTALL}/redeploy.log 2>/dev/null || true",
    "curl -fsS -o /dev/null -w 'HTTP8888=%{http_code}\n' http://127.0.0.1:8888/ 2>&1 || echo HTTP8888=fail",
    "curl -fsS -o /dev/null -w 'HTTP80=%{http_code}\n' http://127.0.0.1:80/ 2>&1 || echo HTTP80=fail",
    f"grep ^FROM {INSTALL}/services/user/Dockerfile 2>/dev/null || true",
]
for cmd in cmds:
    print("===", cmd[:90])
    _, o, e = c.exec_command(cmd, timeout=60)
    print(o.read().decode() + e.read().decode())
c.close()
