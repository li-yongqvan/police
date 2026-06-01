#!/usr/bin/env python3
"""Upload ai-forum-migration.tar.gz to domestic VPS and deploy."""
from __future__ import annotations

import os
import sys
import time
from pathlib import Path

import paramiko

ROOT = Path(__file__).resolve().parents[1]
HOST = os.environ.get("DEPLOY_HOST", "122.51.233.225")
USER = os.environ.get("DEPLOY_USER", "liyongquan")
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = os.environ.get("INSTALL_DIR", f"/home/{USER}/projects/ai-forum")
TGZ_LOCAL = ROOT / "ai-forum-migration.tar.gz"
REMOTE_TGZ = "/tmp/ai-forum-migration.tar.gz"

# Missing from migration tar but required on domestic host (port 8888 + deploy)
EXTRA_FILES = [
    "docker-compose.server.yml",
    "docker-compose.prod.yml",
    ".env.production.example",
    "scripts/deploy-domestic.sh",
    "scripts/init-db.sh",
    "scripts/seed-content.sh",
    "infra/domestic-nginx/ai-forum.conf.example",
]


def run(client, cmd: str, timeout: int = 7200) -> int:
    print(f"$ {cmd[:140]}...")
    _, stdout, stderr = client.exec_command(cmd, get_pty=True, timeout=timeout)
    out = stdout.read().decode("utf-8", errors="replace")
    err = stderr.read().decode("utf-8", errors="replace")
    code = stdout.channel.recv_exit_status()
    if out:
        print(out[-12000:])
    if err.strip():
        print(err[-3000:], file=sys.stderr)
    return code


def main() -> int:
    if not PASSWORD:
        print("Set DEPLOY_PASSWORD", file=sys.stderr)
        return 1
    if not TGZ_LOCAL.is_file():
        print(f"Missing {TGZ_LOCAL}", file=sys.stderr)
        return 1

    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    print(f"Connecting {USER}@{HOST}...")
    client.connect(HOST, username=USER, password=PASSWORD, timeout=60)

    sftp = client.open_sftp()
    print(f"Uploading {TGZ_LOCAL.name}...")
    sftp.put(str(TGZ_LOCAL), REMOTE_TGZ)
    for rel in EXTRA_FILES:
        local = ROOT / rel
        if not local.is_file():
            print(f"skip missing local {rel}")
            continue
        remote = f"{INSTALL}/{rel}".replace("\\", "/")
        remote_dir = os.path.dirname(remote)
        try:
            sftp.stat(remote_dir)
        except OSError:
            # mkdir -p via ssh later
            pass
        print(f"  + {rel}")
    sftp.close()

    run(client, f"mkdir -p {INSTALL}")
    run(client, f"tar -xzf {REMOTE_TGZ} -C {INSTALL}")
  # upload extras after extract
    sftp = client.open_sftp()

    def ensure_dir(path: str) -> None:
        parts = path.strip("/").split("/")
        cur = ""
        for p in parts:
            cur += "/" + p
            try:
                sftp.stat(cur)
            except OSError:
                try:
                    sftp.mkdir(cur)
                except OSError:
                    pass

    for rel in EXTRA_FILES:
        local = ROOT / rel
        if not local.is_file():
            continue
        remote = f"{INSTALL}/{rel}".replace("\\", "/")
        ensure_dir(os.path.dirname(remote))
        sftp.put(str(local), remote)
    sftp.close()

    run(client, f"chmod +x {INSTALL}/scripts/*.sh 2>/dev/null || true")

    run(client, "pkill -f deploy-domestic.sh 2>/dev/null || true", timeout=30)

    deploy = f"cd {INSTALL} && nohup bash scripts/deploy-domestic.sh > deploy.log 2>&1 & echo $!"
    _, stdout, _ = client.exec_command(deploy)
    pid = stdout.read().decode().strip()
    print(f"Deploy started pid={pid}, tailing deploy.log...")

    for i in range(120):
        time.sleep(15)
        _, stdout, _ = client.exec_command(
            f"tail -20 {INSTALL}/deploy.log 2>/dev/null; "
            f"pgrep -f deploy-domestic.sh >/dev/null && echo STATUS=RUNNING || echo STATUS=DONE; "
            f"curl -fsS -o /dev/null -w 'HTTP8888=%{{http_code}}\\n' http://127.0.0.1:8888/ 2>/dev/null || echo HTTP8888=fail; "
            f"curl -fsS -o /dev/null -w 'HTTP80=%{{http_code}}\\n' http://127.0.0.1:80/ 2>/dev/null || echo HTTP80=fail"
        )
        chunk = stdout.read().decode()
        print(chunk[-2500:])
        if "STATUS=DONE" in chunk:
            if "HTTP8888=200" in chunk or "Deploy complete" in chunk:
                break
            if i > 10:
                break

    print("\n=== Access ===")
    print(f"  http://{HOST}/")
    print(f"  http://{HOST}:8888/")
    print("  demo_student / demo123456")
    client.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
