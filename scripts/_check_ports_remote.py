#!/usr/bin/env python3
import os, paramiko
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = "/home/liyongquan/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=PW, timeout=60)

cmds = [
    f"grep -A2 'ports:' {INSTALL}/docker-compose.server.yml",
    "docker port ai-forum-nginx-1 2>/dev/null || true",
    "curl -sS -o /dev/null -w 'local8888=%{http_code}\n' http://127.0.0.1:8888/ 2>&1 || echo local8888=fail",
    "curl -sS -o /dev/null -w 'local8091=%{http_code}\n' http://127.0.0.1:8091/ 2>&1 || echo local8091=fail",
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps --format 'table {{{{.Names}}}}\t{{{{.Status}}}}\t{{{{.Ports}}}}'",
]
for cmd in cmds:
    print("===", cmd[:90])
    _, o, e = c.exec_command(cmd, timeout=30)
    print(o.read().decode() + e.read().decode())
c.close()
