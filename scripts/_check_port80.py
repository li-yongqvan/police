#!/usr/bin/env python3
import os
import paramiko

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=os.environ.get("DEPLOY_PASSWORD", "Liyongquan@123"), timeout=60)
cmds = [
    "curl -sS http://127.0.0.1:80/ | head -12",
    "curl -sS http://127.0.0.1:8090/ | head -12",
    "curl -sS http://127.0.0.1:8091/ | head -8",
    "cat /home/liyongquan/projects/ai-forum/docker-compose.server.yml",
    "docker ps --format '{{.Names}} {{.Ports}}'",
]
for cmd in cmds:
    print(f"\n$ {cmd[:100]}")
    _, o, _ = c.exec_command(cmd, timeout=30)
    print(o.read().decode("utf-8", "replace")[:2000])
c.close()
