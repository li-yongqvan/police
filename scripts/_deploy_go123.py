#!/usr/bin/env python3
import os, paramiko
HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
M = "docker.m.daocloud.io/library"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

def run(cmd, t=3600):
    print("$", cmd[:100])
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    print(out[-10000:])
    return out

run(f"""cd {INSTALL} && sed -i 's/golang:1.22-alpine/golang:1.23-alpine/g' services/*/Dockerfile && grep FROM services/user/Dockerfile""")
run(f"docker pull {M}/golang:1.23-alpine", 600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && DOCKER_BUILDKIT=1 docker compose {cf} build 2>&1", 3600)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d 2>&1", 300)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 45 && docker compose {cf} ps && curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/", 120)
c.close()
print(f"http://{HOST}:8888/")
