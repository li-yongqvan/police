#!/usr/bin/env python3
"""Point host nginx :80 -> forum gateway on :8091."""
import os
import paramiko

PW = os.environ.get("DEPLOY_PASSWORD", "Liyongquan@123")
HOST = "122.51.233.225"

conf = """server {
    listen 80;
    listen [::]:80;
    server_name 122.51.233.225 _;

    client_max_body_size 32m;

    location / {
        proxy_pass http://127.0.0.1:8091;
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

c = paramiko.SSHClient()
c.set_missing_host_key_policy(paramiko.AutoAddPolicy())
c.connect(HOST, username=USER if (USER := "liyongquan") else "liyongquan", password=PW, timeout=60)


def run(cmd, t=60):
    print(f"$ {cmd[:110]}")
    _, o, e = c.exec_command(cmd, get_pty=True, timeout=t)
    out = o.read().decode("utf-8", "replace")
    if out.strip():
        print(out[-3000:])
    return out


run(f"cat > /tmp/ai-forum-proxy.conf << 'EOF'\n{conf}\nEOF")
run(f"echo '{PW}' | sudo -S cp /tmp/ai-forum-proxy.conf /etc/nginx/conf.d/ai-forum.conf 2>&1")
run(f"echo '{PW}' | sudo -S rm -f /etc/nginx/sites-enabled/default 2>&1 || true")
run(f"echo '{PW}' | sudo -S nginx -t 2>&1")
run(f"echo '{PW}' | sudo -S systemctl reload nginx 2>&1")
run("curl -sS http://127.0.0.1:80/ | grep -oE 'index-[^\"]+\\.js' | head -1")
run("curl -sS -o /dev/null -w 'port80=%{http_code}\\n' http://127.0.0.1:80/")
c.close()
print(f"\nUse: http://{HOST}/community/circle")
