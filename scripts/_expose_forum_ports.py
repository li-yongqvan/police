#!/usr/bin/env python3
"""Expose forum on 8888 (docs default) alongside 8091."""
import os
import paramiko

PW = os.environ.get("DEPLOY_PASSWORD", "Liyongquan@123")
INSTALL = "/home/liyongquan/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

patch = """# Ubuntu cloud host — gateway on 8091 and 8888 (host :80 needs root reverse proxy)
services:
  nginx:
    ports:
      - "8091:80"
      - "8888:80"

  postgres:
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
"""

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=PW, timeout=60)

def run(cmd, t=120):
    print(f"$ {cmd[:100]}")
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    print(o.read().decode("utf-8", "replace")[-2500:])

run(f"cat > {INSTALL}/docker-compose.server.yml << 'EOF'\n{patch}\nEOF")
run(
    f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
    f"docker compose {cf} up -d nginx",
)
for p in ("8091", "8888"):
    run(f"curl -sS http://127.0.0.1:{p}/ | grep -oE 'index-[^\"]+\\.js' | head -1")
    run(f"curl -sS -o /dev/null -w 'HTTP{p}=%{{http_code}}\\n' http://127.0.0.1:{p}/")
c.close()
print("\nForum URLs:")
print("  http://122.51.233.225:8091/community")
print("  http://122.51.233.225:8888/community")
print("Port 80 is still system nginx — use :8091 or :8888, or set 1Panel reverse proxy to 8091")
