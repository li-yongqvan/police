#!/usr/bin/env python3
import os
import sys
import paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"


def run(c, cmd, t=3600):
    print(f"$ {cmd[:110]}")
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    sys.stdout.buffer.write(out[-10000:].encode("utf-8", "replace"))
    sys.stdout.buffer.write(b"\n")
    return out


def main():
    if not PW:
        raise SystemExit("Set DEPLOY_PASSWORD")
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=PW, timeout=60)
    run(
        c,
        "docker pull docker.m.daocloud.io/library/node:20-alpine && "
        "docker pull docker.m.daocloud.io/library/nginx:alpine",
        600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} build frontend --no-cache 2>&1 | tail -80",
        3600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} up -d --force-recreate frontend nginx",
        300,
    )
    run(c, "sleep 6 && curl -sS -o /dev/null -w 'HTTP8091=%{http_code}\\n' http://127.0.0.1:8091/")
    run(c, "curl -sS http://127.0.0.1:8091/ | grep -oE 'index-[^\"]+\\.js' | head -3")
    run(c, "curl -sS http://127.0.0.1:8091/ | grep -oE 'CommunityHome-[^\"]+\\.js' | head -1")
    c.close()
    print(f"\nOpen: http://{HOST}:8091/community (Ctrl+Shift+R)")


if __name__ == "__main__":
    main()
