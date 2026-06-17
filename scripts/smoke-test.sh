#!/usr/bin/env bash
# Smoke test via gateway (production) or direct services (dev).
set -euo pipefail

BASE_URL="${BASE_URL:-}"
USER_URL="${USER_URL:-http://127.0.0.1:8001}"
FORUM_URL="${FORUM_URL:-http://127.0.0.1:8002}"
ADMIN_URL="${ADMIN_URL:-http://127.0.0.1:8003}"
SMOKE_USER="${SMOKE_USER:-demo_student}"
SMOKE_PASS="${SMOKE_PASS:-demo123456}"

if [ -n "$BASE_URL" ]; then
  BASE_URL="${BASE_URL%/}"
  echo "==> Gateway smoke: $BASE_URL"
  curl -fsS "${BASE_URL}/" >/dev/null && echo "frontend OK"
  curl -fsS "${BASE_URL}/forum-api/api/v1/stats/community" >/dev/null && echo "forum stats OK"
  LOGIN_JSON=$(printf '%s' "{\"username\":\"${SMOKE_USER}\",\"password\":\"${SMOKE_PASS}\"}")
  TOKEN=$(curl -fsS -X POST "${BASE_URL}/user-api/api/v1/login" \
    -H 'Content-Type: application/json' \
    --data-binary "$LOGIN_JSON" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
  [ -n "${TOKEN}" ] || { echo "login failed"; exit 1; }
  echo "login OK"
  curl -fsS -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/user-api/api/v1/users/me" >/dev/null
  echo "users/me OK"
  curl -fsS "${BASE_URL}/forum-api/api/v1/boards" >/dev/null
  echo "boards OK"
  curl -fsS "${BASE_URL}/forum-api/api/v1/posts?limit=5" >/dev/null
  echo "posts OK"
  echo "Gateway smoke test passed."
  exit 0
fi

curl -fsS "${USER_URL}/health" >/dev/null && echo "user-service OK"
curl -fsS "${FORUM_URL}/health" >/dev/null && echo "forum-service OK"
curl -fsS "${ADMIN_URL}/health" >/dev/null && echo "admin-service OK"

LOGIN_JSON=$(printf '%s' "{\"username\":\"${SMOKE_USER}\",\"password\":\"${SMOKE_PASS}\"}")
TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/login" \
  -H 'Content-Type: application/json' \
  --data-binary "$LOGIN_JSON" 2>/dev/null | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p' || true)

if [ -z "${TOKEN}" ]; then
  TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/demo-login" \
    -H 'Content-Type: application/json' \
    -d '{"role":"student"}' | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
  [ -z "${TOKEN}" ] && TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/demo-login" \
    -H 'Content-Type: application/json' \
    -d '{"role":"student"}' | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
fi

[ -n "${TOKEN}" ] || { echo "login failed"; exit 1; }
echo "login OK"

curl -fsS -H "Authorization: Bearer ${TOKEN}" "${USER_URL}/api/v1/users/me" >/dev/null
echo "users/me OK"

curl -fsS "${FORUM_URL}/api/v1/boards" >/dev/null
echo "boards OK"

curl -fsS "${FORUM_URL}/api/v1/posts?limit=5" >/dev/null
echo "posts OK"

echo "Smoke test passed."
