#!/bin/bash
set -eu
echo "=== Containers ==="
docker ps --format '{{.Names}} {{.Status}}' | grep ai-forum

echo "=== Homepage ==="
curl -sS -H 'Host: 107.172.138.10' -o /dev/null -w 'HTTP: %{http_code}\n' http://127.0.0.1:80/
curl -sS -H 'Host: 107.172.138.10' http://127.0.0.1:80/ | grep -o 'index-[^"]*\.js' | head -1

echo "=== Frontend chunks ==="
docker exec ai-forum-frontend-1 ls /usr/share/nginx/html/assets/ | grep -E 'Register|AdminPosts|index-CUR' || true

echo "=== Login ==="
curl -sS -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_student","password":"demo123456"}' | grep -q access_token && echo OK || echo FAIL

echo "=== Admin APIs ==="
TOKEN=$(curl -sS -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_admin","password":"demo123456"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

curl -sS -o /dev/null -w "admin/posts: %{http_code}\n" -H "Authorization: Bearer $TOKEN" \
  http://127.0.0.1:8888/admin-api/api/v1/admin/posts

curl -sS -o /dev/null -w "admin/unban: %{http_code}\n" -X POST \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{}' \
  http://127.0.0.1:8888/admin-api/api/v1/admin/users/1/unban

curl -sS -o /dev/null -w "admin/boards: %{http_code}\n" -H "Authorization: Bearer $TOKEN" \
  http://127.0.0.1:8888/admin-api/api/v1/admin/boards

echo "=== Done ==="
