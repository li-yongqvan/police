#!/usr/bin/env python3
import os
import paramiko

HOST = "122.51.233.225"
USER = "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "Liyongquan@123")

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

cmds = [
    "ss -tlnp | grep -E ':80|:8091|:8888' || true",
    "docker ps --format '{{.Names}} {{.Ports}}' | grep -E 'nginx|frontend' || true",
    "curl -sS http://127.0.0.1:8091/forum-api/api/v1/boards | python3 -c \"import sys,json; b=json.load(sys.stdin); print([x['slug'] for x in b])\"",
    "curl -sS http://127.0.0.1:8091/ | grep -oE 'index-[^\"]+\\.js' | head -1",
    "test -f /opt/ai-forum/frontend/src/composables/useGxNav.js && grep circle /opt/ai-forum/frontend/src/composables/useGxNav.js | head -1 || echo 'no /opt/ai-forum nav'",
    "curl -sS -o /dev/null -w 'port80=%{http_code}\\n' http://127.0.0.1:80/ 2>/dev/null || echo port80=fail",
]
for cmd in cmds:
    print(f"\n$ {cmd[:90]}")
    _, o, _ = c.exec_command(cmd, timeout=30)
    print(o.read().decode("utf-8", "replace").strip())
c.close()
