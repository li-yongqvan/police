#!/usr/bin/env python3
import os, paramiko
HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

cmds = [
    "pkill -f 'docker compose' 2>/dev/null; pkill -f redeploy-cn 2>/dev/null; true",
    f"ls -la {INSTALL}/scripts/redeploy-cn.sh {INSTALL}/docker-compose.cn-mirror.yml {INSTALL}/.env 2>&1",
    f"head -5 {INSTALL}/.env 2>/dev/null",
    "docker images | head -15",
    "docker ps -a | head -15",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker pull docker.m.daocloud.io/library/nginx:alpine 2>&1 | tail -5",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose -f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml pull nginx postgres redis 2>&1 | tail -20",
]
for cmd in cmds:
    print("===", cmd[:90])
    _, o, e = c.exec_command(cmd, timeout=300)
    print(o.read().decode() + e.read().decode())
c.close()
