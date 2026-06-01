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
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps -a",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} logs nginx --tail 40",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} logs user-service --tail 15",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} logs forum-service --tail 15",
    "curl -v http://127.0.0.1:8888/ 2>&1 | tail -25",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} exec -T nginx wget -qO- http://frontend/ 2>&1 | head -5",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} exec -T nginx wget -qO- http://user-service:8001/health 2>&1",
]
for cmd in cmds:
    print("\n===", cmd[:100], "===")
    _, o, e = c.exec_command(cmd, timeout=90)
    print(o.read().decode("utf-8", "replace") + e.read().decode("utf-8", "replace"))
c.close()
