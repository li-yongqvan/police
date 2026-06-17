#!/usr/bin/env python3
"""Deploy top-nav fix: RankView, CircleView, campus-circle migration."""
import os
import sys
import tarfile
from pathlib import Path

import paramiko

HOST = os.environ.get("DEPLOY_HOST", "122.51.233.225")
USER = os.environ.get("DEPLOY_USER", "liyongquan")
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = os.environ.get("INSTALL_DIR", f"/home/{USER}/projects/ai-forum")
ROOT = Path(__file__).resolve().parents[1]
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"
PORT = os.environ.get("DEPLOY_PORT", "8091")

INCLUDE = [
    "frontend",
    "mobile-web",
    "services/forum",
    "docker-compose.cn-mirror.yml",
]


def pack() -> Path:
    out = ROOT / "nav-fix-deploy.tgz"

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
    _, o, _ = c.exec_command(cmd, get_pty=True, timeout=t)
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
    sftp.put(str(tgz), "/tmp/nav-fix-deploy.tgz")
    sftp.close()
    run(c, f"cd {INSTALL} && tar -xzf /tmp/nav-fix-deploy.tgz")
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} build frontend forum-service 2>&1 | tail -55",
        3600,
    )
    run(
        c,
        f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && "
        f"docker compose {cf} up -d --force-recreate frontend forum-service nginx",
        300,
    )
    for p in (PORT, "8888"):
        run(c, f"curl -sS -o /dev/null -w 'HTTP{p}=%{{http_code}}\\n' http://127.0.0.1:{p}/ 2>/dev/null || true")
    run(
        c,
        f"grep -n 'community/circle' {INSTALL}/frontend/src/composables/useGxNav.js | head -3 || echo 'WARN: nav fix not on disk'",
        30,
    )
    run(
        c,
        f"curl -sS http://127.0.0.1:{PORT}/forum-api/api/v1/boards 2>/dev/null | head -c 500",
        30,
    )
    c.close()
    _print(f"\nDeployed. Try: http://{HOST}:{PORT}/community/circle")
    _print(f"              http://{HOST}:{PORT}/community/rank")


if __name__ == "__main__":
    main()
