#!/usr/bin/env bash
# Cron-friendly smoke check. Set BASE_URL/DISCOURSE_URL/SMOKE_USER/SMOKE_PASS in $ROOT/.env.smoke
set -eu
ROOT="${AI_FORUM_PROJECT_ROOT:-/home/liyongquan/projects/ai-forum}"
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