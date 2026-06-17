import paramiko
c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password="Liyongquan@123", timeout=30)
_, o, _ = c.exec_command("ps aux | grep deploy-domestic | grep -v grep; ps aux | grep 'docker compose' | grep -v grep | head -5")
print(o.read().decode())
c.close()
