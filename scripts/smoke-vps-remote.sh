#!/usr/bin/env bash
set -eu
BASE="${BASE_URL:-http://127.0.0.1:8888}"
BASE="${BASE%/}"
SMOKE_USER="${SMOKE_USER:-demo_student}"
SMOKE_PASS="${SMOKE_PASS:-demo123456}"
LOG=/tmp/ai-forum-smoke-result.txt
exec >"$LOG" 2>&1

echo "=== VPS smoke $(date -Is) ==="
echo "BASE=$BASE"

curl -fsS -o /dev/null "$BASE/" && echo "frontend OK"
curl -fsS "$BASE/forum-api/api/v1/stats/community" && echo ""

export SMOKE_USER SMOKE_PASS BASE_URL
python3 <<'PY'
import json, os, urllib.request
base = os.environ.get("BASE_URL", "http://127.0.0.1:8888").rstrip("/")
user = os.environ.get("SMOKE_USER", "demo_student")
pwd = os.environ.get("SMOKE_PASS", "demo123456")
body = json.dumps({"username": user, "password": pwd}).encode()
req = urllib.request.Request(
    base + "/user-api/api/v1/login",
    data=body,
    headers={"Content-Type": "application/json"},
    method="POST",
)
with urllib.request.urlopen(req, timeout=30) as resp:
    data = json.loads(resp.read().decode())
token = data.get("access_token") or data.get("token")
if not token:
    raise SystemExit("login failed: " + str(data))
print("login OK")
for path, label in [
    ("/user-api/api/v1/users/me", "users/me"),
    ("/forum-api/api/v1/boards", "boards"),
    ("/forum-api/api/v1/posts?limit=5", "posts"),
]:
    r = urllib.request.Request(base + path, headers={"Authorization": "Bearer " + token})
    with urllib.request.urlopen(r, timeout=30) as resp:
        if resp.status != 200:
            raise SystemExit(label + " status " + str(resp.status))
    print(label + " OK")
print("=== VPS smoke PASSED ===")
PY

cat "$LOG"
