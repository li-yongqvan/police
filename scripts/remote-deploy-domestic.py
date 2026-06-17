#!/usr/bin/env python3
"""Upload and deploy AI forum to domestic server (no sudo)."""
from __future__ import annotations

import os
import sys
import tarfile
import tempfile
from pathlib import Path

try:
    import paramiko
except ImportError:
    print("Run: python -m pip install paramiko", file=sys.stderr)
    sys.exit(1)

ROOT = Path(__file__).resolve().parents[1]
HOST = os.environ.get("DEPLOY_HOST", "122.51.233.225")
USER = os.environ.get("DEPLOY_USER", "liyongquan")
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL_DIR = os.environ.get("INSTALL_DIR", "~/projects/ai-forum")
REMOTE_TGZ = "/tmp/ai-forum-deploy.tgz"

EXCLUDE_PARTS = {
    "node_modules",
    ".git",
    "reference",
    "frontend-standalone",
    ".cursor",
    "pilot-deploy",
}


def pack_project() -> Path:
    tgz = ROOT / "pilot-deploy.tgz"
    with tempfile.TemporaryDirectory() as tmp:
        staging = Path(tmp) / "ai-forum"
        staging.mkdir()
        for item in ROOT.rglob("*"):
            if not item.is_file():
                continue
            rel = item.relative_to(ROOT)
            if any(p in EXCLUDE_PARTS for p in rel.parts):
                continue
            if rel.suffix in {".tgz", ".tar.gz"}:
                continue
            dest = staging / rel
            dest.parent.mkdir(parents=True, exist_ok=True)
            dest.write_bytes(item.read_bytes())
        with tarfile.open(tgz, "w:gz") as tar:
            tar.add(staging, arcname=".")
    print(f"Packed {tgz} ({tgz.stat().st_size // 1024} KB)")
    return tgz


def main() -> int:
    if not PASSWORD:
        print("Set DEPLOY_PASSWORD", file=sys.stderr)
        return 1

    tgz = pack_project()
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    print(f"Connecting {USER}@{HOST}...")
    client.connect(HOST, username=USER, password=PASSWORD, timeout=30)

    sftp = client.open_sftp()
    print(f"Uploading {REMOTE_TGZ}...")
    sftp.put(str(tgz), REMOTE_TGZ)
    sftp.close()

    remote_install = INSTALL_DIR.replace("~", f"/home/{USER}")
    steps = f"""
set -euo pipefail
mkdir -p {remote_install}
tar -xzf {REMOTE_TGZ} -C {remote_install}
cd {remote_install}
chmod +x scripts/*.sh 2>/dev/null || true
bash scripts/deploy-domestic.sh
"""
    print("Deploying (build may take 5-15 min)...")
    _, stdout, stderr = client.exec_command(steps, get_pty=True)
    out = stdout.read().decode("utf-8", errors="replace")
    err = stderr.read().decode("utf-8", errors="replace")
    code = stdout.channel.recv_exit_status()
    print(out)
    if err:
        print(err, file=sys.stderr)
    client.close()

    if code != 0:
        print(f"Failed: exit {code}", file=sys.stderr)
        return code

    print(f"\nTry: http://{HOST}:8888/")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
