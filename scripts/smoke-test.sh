#!/usr/bin/env bash
# Smoke test via gateway (production) or direct services (dev).
# Current architecture: Vue SPA (:8888) + Discourse (:8080) + user-service + admin-service.
set -euo pipefail

BASE_URL="${BASE_URL:-}"
DISCOURSE_URL="${DISCOURSE_URL:-http://122.51.233.225:8080}"
USER_URL="${USER_URL:-http://127.0.0.1:8001}"
ADMIN_URL="${ADMIN_URL:-http://127.0.0.1:8003}"
SMOKE_USER="${SMOKE_USER:-demo_student}"
SMOKE_PASS="${SMOKE_PASS:-}"

if [ -n "$BASE_URL" ]; then
  BASE_URL="${BASE_URL%/}"
  echo "==> Gateway smoke: $BASE_URL"
  curl -fsS "${BASE_URL}/" >/dev/null && echo "frontend OK"
  curl -fsS "${DISCOURSE_URL%/}/latest.json" >/dev/null && echo "discourse OK"
  if [ -n "$SMOKE_PASS" ]; then
    LOGIN_JSON=$(printf '%s' "{\"username\":\"${SMOKE_USER}\",\"password\":\"${SMOKE_PASS}\"}")
    TOKEN=$(curl -fsS -X POST "${BASE_URL}/api/v1/login" \
      -H 'Content-Type: application/json' \
      --data-binary "$LOGIN_JSON" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
    [ -n "${TOKEN}" ] || { echo "login failed"; exit 1; }
    echo "login OK"
    curl -fsS -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/user-api/api/v1/users/me" >/dev/null
    echo "users/me OK"
  else
    echo "login skipped (SMOKE_PASS empty)"
  fi
  echo "Gateway smoke test passed."
  exit 0
fi

curl -fsS "${USER_URL}/health" >/dev/null && echo "user-service OK"
curl -fsS "${ADMIN_URL}/health" >/dev/null && echo "admin-service OK"
curl -fsS "${DISCOURSE_URL%/}/latest.json" >/dev/null && echo "discourse OK"

if [ -n "$SMOKE_PASS" ]; then
  LOGIN_JSON=$(printf '%s' "{\"username\":\"${SMOKE_USER}\",\"password\":\"${SMOKE_PASS}\"}")
  TOKEN=$(curl -fsS -X POST "${USER_URL}/api/v1/login" \
    -H 'Content-Type: application/json' \
    --data-binary "$LOGIN_JSON" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
  [ -n "${TOKEN}" ] || { echo "login failed"; exit 1; }
  echo "login OK"
  curl -fsS -H "Authorization: Bearer ${TOKEN}" "${USER_URL}/api/v1/users/me" >/dev/null
  echo "users/me OK"
else
  echo "login skipped (SMOKE_PASS empty)"
fi

echo "Smoke test passed."