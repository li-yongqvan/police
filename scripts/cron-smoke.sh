#!/usr/bin/env bash
# Cron-friendly smoke check. Set BASE_URL in /opt/ai-forum/.env.smoke
set -eu
ROOT="/opt/ai-forum"
LOG="${ROOT}/logs/smoke-cron.log"
mkdir -p "${ROOT}/logs"
ENV_FILE="${ROOT}/.env.smoke"

if [ -f "$ENV_FILE" ]; then
  # shellcheck disable=SC1090
  set -a && source "$ENV_FILE" && set +a
fi

{
  echo "=== $(date -Is) ==="
  if bash "${ROOT}/scripts/smoke-test.sh"; then
    echo "PASS"
  else
    echo "FAIL"
  fi
} >>"$LOG" 2>&1
