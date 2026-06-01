#!/usr/bin/env python3
"""Deploy Reddit-style feed UI + forum sort API to domestic VPS (port 8091)."""
import os
import sys
import tarfile
from pathlib import Path

import paramiko

HOST, USER = "122.51.233.225", "liyongquan"
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = f"/home/{USER}/projects/ai-forum"
ROOT = Path(__file__).resolve().parents[1]
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
PORT = "8091"

INCLUDE = [
    "frontend",
    "services/forum",
    "docker-compose.cn-mirror.yml",
]


def pack() -> Path:
    out = ROOT / "reddit-feed-deploy.tgz"

    def _filter(tarinfo):
        if "node_modules" in tarinfo.name.replace("\\", "/"):
            return None
        return tarinfo

    with tarfile.open(out, "w:gz") as tar:
        for name in INCLUDE:
            path = ROOT / name
            if path.is_dir():
                tar.add(path, arcname=name, filter=_filter)
            elif path.is_file():
                tar.add(path, arcname=name)
    print(f"Packed {out} ({out.stat().st_size // 1024} KB)")
    return out


def _print(text: str) -> None:
    sys.stdout.buffer.write(text.encode("utf-8", "replace"))
    sys.stdout.buffer.write(b"\n")


def run(c, cmd, t=3600):
    _print(f"$ {cmd[:120]}")
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        _print(out[-12000:])
    return out


def main():
    if not PW:
        raise SystemExit("Set DEPLOY_PASSWORD environment variable")
    tgz = pack()
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=PW, timeout=60)
    sftp = c.open_sftp()
    sftp.put(str(tgz), "/tmp/reddit-feed-deploy.tgz")
    sftp.close()
    run(c, f"cd {INSTALL} && tar -xzf /tmp/reddit-feed-deploy.tgz")
    run(
        c,
        f"sed -i 's/\"8888:80\"/\"{PORT}:80\"/' {INSTALL}/docker-compose.server.yml 2>/dev/null; "
        f"grep -A2 'nginx:' {INSTALL}/docker-compose.server.yml | head -5; "
        f"grep ports {INSTALL}/docker-compose.server.yml || true",
        30,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} build frontend forum-service 2>&1 | tail -50",
        3600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} up -d frontend forum-service nginx",
        300,
    )
    run(c, f"sleep 8 && curl -sS -o /dev/null -w 'HTTP8091=%{{http_code}}\\n' http://127.0.0.1:{PORT}/")
    run(
        c,
        f"curl -sS http://127.0.0.1:{PORT}/ 2>/dev/null | head -c 500 | tr '\\n' ' '",
        30,
    )
    c.close()
    _print(f"\nLive: http://{HOST}:{PORT}/community")


if __name__ == "__main__":
    main()
