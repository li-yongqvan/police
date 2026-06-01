#!/usr/bin/env python3
import os, secrets, sys, time, paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
jwt, pg = secrets.token_hex(32), secrets.token_hex(16)

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER, password=PW, timeout=60)

def run(cmd, t=3600):
    print("$", cmd[:110])
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    print(out[-15000:])
    return o.channel.recv_exit_status(), out

run(f"cd {INSTALL} && cp .env.production.example .env", 30)
run(
    f"""cd {INSTALL} && sed -i 's/^JWT_SECRET=.*/JWT_SECRET={jwt}/' .env && """
    f"""sed -i 's/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD={pg}/' .env && """
    f"""sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD={pg}/' .env && """
    f"""sed -i 's/^DB_HOST=.*/DB_HOST=postgres/' .env && """
    f"""sed -i 's/^REDIS_HOST=.*/REDIS_HOST=redis/' .env && """
    f"""grep -E '^(JWT_SECRET|POSTGRES_PASSWORD|DB_HOST|REDIS_HOST)=' .env""",
    30,
)

cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} up -d --build 2>&1", 3600)

run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && sleep 35 && docker compose {cf} restart nginx", 120)
run(f"cd {INSTALL} && bash scripts/seed-content.sh 2>&1 || true", 180)
run(f"curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/; curl -fsS http://127.0.0.1:8888/forum-api/api/v1/stats/community 2>/dev/null | head -c 120; echo", 60)
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} ps", 60)

c.close()
print(f"\nOpen: http://{HOST}:8888/  (demo_student / demo123456)")
