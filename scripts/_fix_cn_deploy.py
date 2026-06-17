#!/usr/bin/env python3
"""Fix .env, Docker mirror, redeploy on domestic VPS."""
import os
import secrets
import sys
import time

import paramiko

HOST = os.environ.get("DEPLOY_HOST", "122.51.233.225")
USER = os.environ.get("DEPLOY_USER", "liyongquan")
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def run(client, cmd: str, timeout: int = 7200) -> str:
    print(f"$ {cmd[:130]}")
    _, stdout, stderr = client.exec_command(cmd, get_pty=True, timeout=timeout)
    out = stdout.read().decode("utf-8", errors="replace")
    err = stderr.read().decode("utf-8", errors="replace")
    code = stdout.channel.recv_exit_status()
    text = out + err
    if text.strip():
        print(text[-12000:])
    if code != 0:
        print(f"[exit {code}]", file=sys.stderr)
    return text


def main() -> int:
    if not PASSWORD:
        print("Set DEPLOY_PASSWORD", file=sys.stderr)
        return 1

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(HOST, username=USER, password=PASSWORD, timeout=60)

    run(client, "pkill -f deploy-domestic.sh 2>/dev/null || true", 30)

    # Upload mirror compose
    sftp = client.open_sftp()
    sftp.put(
        os.path.join(ROOT, "docker-compose.cn-mirror.yml"),
        f"{INSTALL}/docker-compose.cn-mirror.yml",
    )
    sftp.close()

    jwt = secrets.token_hex(32)
    pg = secrets.token_hex(16)
    env_script = f"""set -e
cd {INSTALL}
if [ ! -f .env ]; then
  cp .env.production.example .env 2>/dev/null || cp .env.example .env
fi
grep -q '^JWT_SECRET=replace' .env && sed -i 's/^JWT_SECRET=.*/JWT_SECRET={jwt}/' .env || true
grep -q '^POSTGRES_PASSWORD=replace' .env && sed -i 's/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD={pg}/' .env || true
grep -q '^DB_PASSWORD=replace' .env && sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD={pg}/' .env || true
grep -q '^POSTGRES_PASSWORD=$' .env 2>/dev/null && sed -i 's/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD={pg}/' .env || true
echo "env ok"
"""
    run(client, env_script, 60)

    mirror_json = """{
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://docker.1panel.live",
    "https://hub.rat.dev"
  ]
}
"""
    run(
        client,
        f"""echo '{mirror_json}' | sudo tee /etc/docker/daemon.json >/dev/null 2>&1 && sudo systemctl restart docker && sleep 5 && echo mirror_ok || echo mirror_skip_no_sudo""",
        120,
    )

    compose = (
        "docker compose -f docker-compose.yml -f docker-compose.server.yml "
        "-f docker-compose.cn-mirror.yml"
    )
    redeploy_sh = f"""#!/usr/bin/env bash
set -euo pipefail
cd {INSTALL}
export PATH="$HOME/.docker/cli-plugins:$PATH"
CF="-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
echo "==> pull base images"
docker compose $CF pull nginx postgres redis || true
echo "==> build & up"
docker compose $CF up -d --build
echo "==> wait"
sleep 35
docker compose $CF restart nginx || true
sleep 10
if [ -f scripts/init-db.sh ]; then bash scripts/init-db.sh || true; fi
if [ -f scripts/seed-content.sh ]; then bash scripts/seed-content.sh || true; fi
curl -fsS -o /dev/null -w "HTTP8888=%{{http_code}}\\n" http://127.0.0.1:8888/ || true
echo DONE
"""
    remote_sh = f"{INSTALL}/scripts/redeploy-cn.sh"
    sftp = client.open_sftp()
    with sftp.file(remote_sh, "w") as f:
        f.write(redeploy_sh)
    sftp.close()
    run(client, f"chmod +x {remote_sh} && nohup bash {remote_sh} > {INSTALL}/redeploy.log 2>&1 & echo started", 30)

    for i in range(90):
        time.sleep(25)
        text = run(
            client,
            f"tail -20 {INSTALL}/redeploy.log 2>/dev/null; "
            "pgrep -f 'docker compose' >/dev/null && echo BUILD=running || echo BUILD=idle; "
            "curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/ 2>/dev/null || echo HTTP8888=fail",
            90,
        )
        if "HTTP8888=200" in text:
            break
        if "BUILD=idle" in text and ("error" in text.lower() or "Error" in text) and i > 4:
            break
        if "BUILD=idle" in text and "HTTP8888=200" in text:
            break
        if "BUILD=idle" in text and i > 15:
            break

    run(client, f"cd {INSTALL} && docker compose -f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml ps", 60)

    client.close()
    print(f"\n=== Try in browser ===\n  http://{HOST}:8888/\n  http://{HOST}/\n  demo_student / demo123456")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
