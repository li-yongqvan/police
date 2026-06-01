#!/usr/bin/env python3
import os, paramiko
PW = os.environ.get("DEPLOY_PASSWORD", "")
INSTALL = "/home/liyongquan/projects/ai-forum"
cf = "-f docker-compose.yml -f docker-compose.server.yml -f docker-compose.cn-mirror.yml"

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect("122.51.233.225", username="liyongquan", password=PW, timeout=60)

def run(cmd, t=120):
    print("$", cmd[:100])
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    print(out[-4000:])
    return out

# Restart docker nginx first
run(f"cd {INSTALL} && export PATH=$HOME/.docker/cli-plugins:$PATH && docker compose {cf} restart nginx && sleep 3 && curl -sS -o /dev/null -w 'local8888=%{{http_code}}\\n' http://127.0.0.1:8888/", 60)

conf = """server {
    listen 80;
    server_name 122.51.233.225;
    client_max_body_size 32m;
    location / {
        proxy_pass http://127.0.0.1:8888;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 30s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
"""
run(f"cat > /tmp/ai-forum-proxy.conf << 'EOF'\n{conf}\nEOF", 30)

# Try sudo with same password (common on student VPS)
sudo_pw = PW
run(f"echo '{sudo_pw}' | sudo -S cp /tmp/ai-forum-proxy.conf /etc/nginx/conf.d/ai-forum.conf 2>&1", 30)
run(f"echo '{sudo_pw}' | sudo -S nginx -t 2>&1", 30)
run(f"echo '{sudo_pw}' | sudo -S systemctl reload nginx 2>&1", 30)
run("curl -sS -o /dev/null -w 'local80=%{http_code}\n' http://127.0.0.1:80/ -H 'Host: 122.51.233.225'", 30)
run("curl -sS http://127.0.0.1:80/ -H 'Host: 122.51.233.225' 2>&1 | head -8", 30)

c.close()
print("\nTry: http://122.51.233.225/  and  http://122.51.233.225:8888/")
