#!/usr/bin/env python3
import os, sys, time, paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def run(c, cmd, t=600):
    print("$", cmd[:100])
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = (o.read().decode("utf-8", "replace") + e.read().decode("utf-8", "replace"))
    if out.strip():
        print(out[-8000:])
    return out


def main():
    if not PASSWORD:
        sys.exit("Set DEPLOY_PASSWORD")
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=PASSWORD, timeout=60)
    run(c, "pkill -f deploy-domestic.sh 2>/dev/null; pkill -f redeploy-cn.sh 2>/dev/null; true", 20)

    sftp = c.open_sftp()
    sftp.put(os.path.join(ROOT, "docker-compose.cn-mirror.yml"), f"{INSTALL}/docker-compose.cn-mirror.yml")
    sh = f"""#!/usr/bin/env bash
set -euo pipefail
cd {INSTALL}
export PATH="$HOME/.docker/cli-plugins:$PATH"
CF="-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
docker compose $CF pull nginx postgres redis || true
docker compose $CF up -d --build
sleep 40
docker compose $CF restart nginx || true
sleep 10
[ -f scripts/seed-content.sh ] && bash scripts/seed-content.sh || true
curl -fsS -o /dev/null -w "HTTP8888=%{{http_code}}\\n" http://127.0.0.1:8888/ || true
docker compose $CF ps
echo REDEPLOY_DONE
"""
    path = f"{INSTALL}/scripts/redeploy-cn.sh"
    with sftp.file(path, "w") as f:
        f.write(sh)
    sftp.close()
    run(c, f"chmod +x {path} && nohup bash {path} > {INSTALL}/redeploy.log 2>&1 &", 30)

    for i in range(80):
        time.sleep(30)
        t = run(c, f"tail -25 {INSTALL}/redeploy.log; pgrep -f redeploy-cn && echo RUNNING || echo IDLE", 120)
        if "REDEPLOY_DONE" in t:
            break
        if "IDLE" in t and i > 3:
            break
        if "HTTP8888=200" in t:
            break

    c.close()
    print(f"\nhttp://{HOST}:8888/  http://{HOST}/")


if __name__ == "__main__":
    main()
