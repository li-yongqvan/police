#!/usr/bin/env python3
"""Continue domestic deploy on 122.51.233.225."""
import os
import sys
import time

import paramiko

HOST = os.environ.get("DEPLOY_HOST", "122.51.233.225")
USER = os.environ.get("DEPLOY_USER", "liyongquan")
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"


def run(client, cmd: str, timeout: int = 7200) -> tuple[int, str]:
    print(f"$ {cmd[:120]}")
    _, stdout, stderr = client.exec_command(cmd, get_pty=True, timeout=timeout)
    out = stdout.read().decode("utf-8", errors="replace")
    err = stderr.read().decode("utf-8", errors="replace")
    code = stdout.channel.recv_exit_status()
    text = (out + err).strip()
    if text:
        print(text[-10000:])
    return code, text


def main() -> int:
    if not PASSWORD:
        print("Set DEPLOY_PASSWORD", file=sys.stderr)
        return 1

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(HOST, username=USER, password=PASSWORD, timeout=60)

    print("==> Status check")
    run(
        client,
        f"ls -la {INSTALL}/docker-compose.yml {INSTALL}/scripts/deploy-domestic.sh 2>&1; "
        f"tail -25 {INSTALL}/deploy.log 2>/dev/null || echo 'no deploy.log'; "
        "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null | head -15",
        timeout=60,
    )

    _, text = run(client, "pgrep -f deploy-domestic.sh >/dev/null && echo DEPLOY_RUNNING || echo DEPLOY_IDLE", 30)
    if "DEPLOY_RUNNING" not in text:
        print("==> Starting deploy-domestic.sh")
        run(
            client,
            f"cd {INSTALL} && nohup bash scripts/deploy-domestic.sh >> deploy.log 2>&1 & sleep 2; pgrep -af deploy-domestic || true",
            30,
        )
    else:
        print("==> Deploy already running")

    for i in range(100):
        time.sleep(20)
        _, chunk = run(
            client,
            f"tail -18 {INSTALL}/deploy.log 2>/dev/null; "
            "pgrep -f deploy-domestic.sh >/dev/null && echo STATUS=RUNNING || echo STATUS=DONE; "
            "curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}} ' http://127.0.0.1:8888/ 2>/dev/null || echo HTTP8888=fail; "
            "curl -fsS -o /dev/null -w 'HTTP80=%{{http_code}}\\n' http://127.0.0.1:80/ 2>/dev/null || echo HTTP80=fail; "
            "curl -fsS http://127.0.0.1:8888/forum-api/api/v1/stats/community 2>/dev/null | head -c 80 || true",
            90,
        )
        if "STATUS=DONE" in chunk and ("HTTP8888=200" in chunk or "Deploy complete" in chunk or i >= 6):
            break
        if "STATUS=DONE" in chunk and "failed" in chunk.lower() and i >= 3:
            break

    run(
        client,
        f"cd {INSTALL} && (docker compose -f docker-compose.yml -f docker-compose.server.yml ps 2>/dev/null || true)",
        60,
    )

    client.close()
    print(f"\n=== Access ===\n  http://{HOST}:8888/\n  http://{HOST}/\n  demo_student / demo123456")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
