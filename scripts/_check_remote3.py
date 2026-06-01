import paramiko
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password="Liyongquan@123", timeout=30)
cmds = [
    "ps aux | grep -E 'docker|compose|curl|go build' | grep -v grep",
    "ls -la ~/.docker/cli-plugins/ 2>&1",
    "docker compose version 2>&1",
    "cd ~/projects/ai-forum && tail -5 scripts/deploy-domestic.sh",
]
for cmd in cmds:
    _, o, _ = c.exec_command(cmd)
    print("===", cmd, "===")
    print(o.read().decode())
c.close()
