import paramiko
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password="Liyongquan@123", timeout=30)
cmds = [
    "docker ps -a",
    "ls -la ~/projects/ai-forum 2>&1 | head -8",
    "curl -sI http://127.0.0.1:8888/ 2>&1 | head -8",
    "tail -30 ~/projects/ai-forum/.deploy.log 2>/dev/null || echo no log",
]
for cmd in cmds:
    _, o, _ = c.exec_command(cmd)
    print("===", cmd, "===")
    print(o.read().decode())
c.close()
