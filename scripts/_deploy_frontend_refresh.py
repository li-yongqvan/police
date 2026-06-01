#!/usr/bin/env python3
"""Upload frontend-related changes and rebuild on domestic VPS."""
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

INCLUDE = [
    "frontend",
    "mobile-web",
    "docker-compose.cn-mirror.yml",
]


def pack() -> Path:
    out = ROOT / "frontend-refresh.tgz"

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
    _print(f"$ {cmd[:100]}")
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        _print(out[-8000:])
    return out


def main():
    if not PW:
        raise SystemExit("Set DEPLOY_PASSWORD")
    tgz = pack()
    c = paramiko.SSHClient()
    c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    c.connect(HOST, username=USER, password=PW, timeout=60)
    sftp = c.open_sftp()
    sftp.put(str(tgz), "/tmp/frontend-refresh.tgz")
    sftp.close()
    run(c, f"cd {INSTALL} && tar -xzf /tmp/frontend-refresh.tgz")
    run(
        c,
        f"docker pull docker.m.daocloud.io/library/node:20-alpine && "
        f"docker pull docker.m.daocloud.io/library/nginx:alpine",
        600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} build frontend --no-cache 2>&1 | tail -40",
        3600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} up -d frontend nginx",
        300,
    )
    run(c, "curl -sS -o /dev/null -w 'HTTP8091=%{http_code}\\n' http://127.0.0.1:8091/")
    c.close()
    _print(f"\nLive: http://{HOST}:8091/")


if __name__ == "__main__":
    main()
