#!/usr/bin/env bash
# Fail if production .env still uses placeholder secrets.
set -eu
ENV_FILE="${1:-.env}"
[ -f "$ENV_FILE" ] || { echo "Missing $ENV_FILE"; exit 1; }

# shellcheck disable=SC1090
set -a && source "$ENV_FILE" && set +a

fail=0
check_not_default() {
  local name=$1 val=$2
  if [ -z "$val" ] || [[ "$val" == *replace* ]] || [[ "$val" == *changeme* ]] || [ "$val" = "ai_forum_secret" ]; then
    echo "FAIL: $name must be set to a strong unique value"
    fail=1
  else
    echo "OK: $name"
  fi
}

check_not_default JWT_SECRET "$JWT_SECRET"
check_not_default POSTGRES_PASSWORD "$POSTGRES_PASSWORD"
check_not_default DB_PASSWORD "${DB_PASSWORD:-$POSTGRES_PASSWORD}"

if [ "$fail" -ne 0 ]; then
  echo "Run: openssl rand -hex 32"
  exit 1
fi
echo "Production secrets check passed."
