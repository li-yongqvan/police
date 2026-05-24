#!/usr/bin/env bash
set -euo pipefail

USER_URL="${USER_URL:-http://127.0.0.1:8001}"
FORUM_URL="${FORUM_URL:-http://127.0.0.1:8002}"
ADMIN_URL="${ADMIN_URL:-http://127.0.0.1:8003}"

curl -fsS "${USER_URL}/health" >/dev/null && echo "user-service OK"
curl -fsS "${FORUM_URL}/health" >/dev/null && echo "forum-service OK"
curl -fsS "${ADMIN_URL}/health" >/dev/null && echo "admin-service OK"

TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/demo-login" \
  -H "Content-Type: application/json" \
  -d '{"role":"student"}' | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')

if [ -z "${TOKEN}" ]; then
  TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/demo-login" \
    -H "Content-Type: application/json" \
    -d '{"role":"student"}' | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
fi

[ -n "${TOKEN}" ] || { echo "demo-login failed"; exit 1; }
echo "demo-login OK"

curl -fsS -H "Authorization: Bearer ${TOKEN}" "${USER_URL}/api/v1/users/me" >/dev/null
echo "users/me OK"

curl -fsS "${FORUM_URL}/api/v1/boards" >/dev/null
echo "boards OK"

curl -fsS "${FORUM_URL}/api/v1/posts?limit=5" >/dev/null
echo "posts OK"

echo "Smoke test passed."
