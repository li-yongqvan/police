#!/usr/bin/env python3
import os, paramiko
PW = os.environ.get("DEPLOY_PASSWORD", "")
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=PW, timeout=60)
cmds = [
    "ss -tlnp | grep -E ':8888|:80 ' || netstat -tlnp 2>/dev/null | grep -E ':8888|:80 '",
    "curl -sS -o /dev/null -w 'local8888=%{http_code}\n' http://127.0.0.1:8888/",
    "curl -sS -o /dev/null -w 'pub8888=%{http_code}\n' --connect-timeout 5 http://122.51.233.225:8888/ || echo pub_fail",
    "curl -sS -o /dev/null -w 'local80=%{http_code}\n' http://127.0.0.1:80/ 2>&1 | tail -1",
    "docker port ai-forum-nginx-1 2>/dev/null || true",
    "ps aux | grep -E 'nginx|1panel' | grep -v grep | head -10",
]
for cmd in cmds:
    print("===", cmd)
    _, o, e = c.exec_command(cmd, timeout=30)
    print(o.read().decode() + e.read().decode())
c.close()
