"""Kill stuck deploy, re-upload, run deploy-domestic with logging."""
import importlib.util
import os
import sys
import time
from pathlib import Path

import paramiko

ROOT = Path(__file__).resolve().parents[1]
HOST = "122.51.233.225"
USER = "liyongquan"
PASSWORD = os.environ.get("DEPLOY_PASSWORD", "Liyongquan@123")
INSTALL = f"/home/{USER}/projects/ai-forum"
LOG = f"{INSTALL}/deploy.log"

spec = importlib.util.spec_from_file_location("rd", ROOT / "scripts" / "remote-deploy-domestic.py")
rd = importlib.util.module_from_spec(spec)
spec.loader.exec_module(rd)
pack_project = rd.pack_project


def run(client, cmd, timeout=7200):
    print(f"$ {cmd[:120]}...")
    _, stdout, stderr = client.exec_command(cmd, get_pty=True, timeout=timeout)
    out = stdout.read().decode("utf-8", errors="replace")
    code = stdout.channel.recv_exit_status()
    if out:
        print(out[-8000:])
    err = stderr.read().decode("utf-8", errors="replace")
    if err:
        print(err[-2000:], file=sys.stderr)
    return code


def main():
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(HOST, username=USER, password=PASSWORD, timeout=30)

    run(client, "pkill -f deploy-domestic.sh || true; pkill -f 'curl.*docker-compose' || true")
    run(client, f"rm -f {INSTALL}/.docker/cli-plugins/docker-compose 2>/dev/null; mkdir -p ~/.docker/cli-plugins")

    tgz = pack_project()
    sftp = client.open_sftp()
    sftp.put(str(tgz), "/tmp/ai-forum-deploy.tgz")
    sftp.close()

    run(client, f"mkdir -p {INSTALL} && tar -xzf /tmp/ai-forum-deploy.tgz -C {INSTALL}")
    run(client, f"chmod +x {INSTALL}/scripts/*.sh")

    cmd = f"cd {INSTALL} && nohup bash scripts/deploy-domestic.sh > deploy.log 2>&1 & echo $!"
    _, stdout, _ = client.exec_command(cmd)
    pid = stdout.read().decode().strip()
    print(f"Started deploy pid {pid}, tailing {LOG}...")

    for _ in range(360):
        time.sleep(10)
        _, stdout, _ = client.exec_command(f"tail -15 {LOG} 2>/dev/null; pgrep -f deploy-domestic.sh >/dev/null && echo RUNNING || echo DONE")
        chunk = stdout.read().decode()
        lines = chunk.strip().splitlines()
        if lines:
            print("\n".join(lines[-18:]))
        if "DONE" in chunk and "RUNNING" not in chunk.split("DONE")[0][-20:]:
            break
        if "=== Deploy complete ===" in chunk:
            break

    code, out, _ = client.exec_command(f"grep -q 'Deploy complete' {LOG} && echo OK || echo FAIL; curl -sI http://127.0.0.1:8888/ | head -3")
    result = out.read().decode()
    print(result)
    client.close()
    return 0 if "OK" in result else 1


if __name__ == "__main__":
    raise SystemExit(main())
