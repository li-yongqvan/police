#!/usr/bin/env bash
# E2E Integration Test Script for AI Forum
# Tests all critical API flows end-to-end
# Usage: bash e2e-test.sh
# Expects services running: user-service (:8001), forum-service (:8002), admin-service (:8003)

set -euo pipefail

BASE_USER="http://localhost:8001"
BASE_FORUM="http://localhost:8002"
BASE_ADMIN="http://localhost:8003"

PASS=0
FAIL=0
TOTAL=0

assert_status() {
  local label="$1" expected="$2" actual="$3"
  TOTAL=$((TOTAL + 1))
  if [ "$actual" -eq "$expected" ]; then
    echo "  PASS: $label (HTTP $actual)"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: $label (expected HTTP $expected, got HTTP $actual)"
    FAIL=$((FAIL + 1))
  fi
}

assert_json() {
  local label="$1" expected="$2" actual="$3"
  TOTAL=$((TOTAL + 1))
  if [ "$actual" = "$expected" ]; then
    echo "  PASS: $label"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: $label (expected '$expected', got '$actual')"
    FAIL=$((FAIL + 1))
  fi
}

echo "========================================="
echo " AI Forum E2E Integration Tests"
echo "========================================="
echo ""

# =============================================
# 1. Health Checks
# =============================================
echo "--- 1. Health Checks ---"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_USER/health")
assert_status "user-service health" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_FORUM/health")
assert_status "forum-service health" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_ADMIN/health")
assert_status "admin-service health" 200 "$STATUS"

echo ""

# =============================================
# 2. Invite Code Generation (admin-service -> user-service)
# =============================================
echo "--- 2. Invite Code Generation ---"

ADMIN_TOKEN="Bearer test_admin_token"

# Generate single invite code via user-service internal API
RESPONSE=$(curl -s -X POST "$BASE_USER/internal/v1/invite-codes" \
  -H "Content-Type: application/json" \
  -d '{"created_by": 0}')
CODE=$(echo "$RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['code'])" 2>/dev/null || echo "")
assert_json "invite code generated" "" "$([ -n "$CODE" ] && echo "$CODE" || echo "")"

if [ -z "$CODE" ]; then
  echo "  SKIP: cannot continue without invite code"
  echo ""
  echo "========================================="
  echo " Results: $PASS passed, $FAIL failed out of $TOTAL tests"
  echo "========================================="
  exit 1
fi

echo ""

# =============================================
# 3. User Registration and Login
# =============================================
echo "--- 3. Registration & Login ---"

USERNAME="e2e_user_$(date +%s)"
PASSWORD="TestPass123!"

# Register
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_USER/api/v1/register" \
  -H "Content-Type: application/json" \
  -d "{\"invite_code\":\"$CODE\",\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
STATUS=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
assert_status "register with invite code" 201 "$STATUS"

# Extract user ID and token
USER_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('user_id','') or json.load(sys.stdin).get('id',''))" 2>/dev/null || echo "")
ACCESS_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || echo "")

if [ -z "$ACCESS_TOKEN" ]; then
  echo "  WARN: No access_token in response, trying login..."
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_USER/api/v1/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
  STATUS=$(echo "$RESPONSE" | tail -1)
  BODY=$(echo "$RESPONSE" | sed '$d')
  assert_status "login after registration" 200 "$STATUS"
  ACCESS_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || echo "")
  USER_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('user_id',''))" 2>/dev/null || echo "")
fi

AUTH="Authorization: Bearer $ACCESS_TOKEN"

# Get profile
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_USER/api/v1/users/$USER_ID" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "get user profile" 200 "$STATUS"

echo ""

# =============================================
# 4. Board Listing (forum-service)
# =============================================
echo "--- 4. Board Listing ---"

RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_FORUM/api/v1/boards")
STATUS=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
assert_status "list boards" 200 "$STATUS"

BOARD_ID=$(echo "$BODY" | python3 -c "import sys,json; boards=json.load(sys.stdin).get('boards',[]); print(boards[0]['id'] if boards else '1')" 2>/dev/null || echo "1")

echo ""

# =============================================
# 5. Create Post (forum-service)
# =============================================
echo "--- 5. Post Creation ---"

POST_TITLE="E2E Test Post $(date +%s)"
POST_CONTENT="This is an E2E integration test post content."

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_FORUM/api/v1/posts" \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d "{\"title\":\"$POST_TITLE\",\"content\":\"$POST_CONTENT\",\"board_id\":$BOARD_ID}")
STATUS=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')
assert_status "create post" 201 "$STATUS"

POST_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))" 2>/dev/null || echo "")

# List posts
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_FORUM/api/v1/posts")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list posts" 200 "$STATUS"

# Get post detail
if [ -n "$POST_ID" ]; then
  RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_FORUM/api/v1/posts/$POST_ID")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "get post detail" 200 "$STATUS"
fi

echo ""

# =============================================
# 6. Comment on Post
# =============================================
echo "--- 6. Comments ---"

if [ -n "$POST_ID" ]; then
  COMMENT_CONTENT="E2E test comment content"
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_FORUM/api/v1/posts/$POST_ID/comments" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d "{\"content\":\"$COMMENT_CONTENT\"}")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "create comment" 201 "$STATUS"

  # List comments
  RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_FORUM/api/v1/posts/$POST_ID/comments")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "list comments" 200 "$STATUS"
fi

echo ""

# =============================================
# 7. Like and Collect post
# =============================================
echo "--- 7. Interactions ---"

if [ -n "$POST_ID" ]; then
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_FORUM/api/v1/posts/$POST_ID/like" -H "$AUTH")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "like post" 200 "$STATUS"

  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_FORUM/api/v1/posts/$POST_ID/collect" -H "$AUTH")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "collect post" 200 "$STATUS"
fi

echo ""

# =============================================
# 8. Admin Operations
# =============================================
echo "--- 8. Admin Operations ---"

# List pending audits
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/audit/pending?page=1&limit=10" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list pending audits" 200 "$STATUS"

# List users
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/users?page=1&limit=10" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list users (admin)" 200 "$STATUS"

# List posts (admin)
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/posts?page=1&limit=10" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list posts (admin)" 200 "$STATUS"

# List boards (admin)
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/boards" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list boards (admin)" 200 "$STATUS"

# List invite codes (admin)
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/invite-codes?page=1&limit=10" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list invite codes" 200 "$STATUS"

# List sensitive words
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/sensitive-words" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list sensitive words" 200 "$STATUS"

# List roles
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/roles" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "list roles" 200 "$STATUS"

echo ""

# =============================================
# 9. Stats Endpoints
# =============================================
echo "--- 9. Statistics ---"

# Admin stats overview
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/stats/overview" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "stats overview" 200 "$STATUS"

# Admin daily stats
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/stats/daily?days=7" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "daily stats" 200 "$STATUS"

echo ""

# =============================================
# 10. Config Endpoints
# =============================================
echo "--- 10. Config ---"

RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_ADMIN/api/v1/admin/config" -H "$AUTH")
STATUS=$(echo "$RESPONSE" | tail -1)
assert_status "get config" 200 "$STATUS"

echo ""

# =============================================
# 11. Post Deletion (own post)
# =============================================
echo "--- 11. Post Deletion ---"

if [ -n "$POST_ID" ]; then
  RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_FORUM/api/v1/posts/$POST_ID" -H "$AUTH")
  STATUS=$(echo "$RESPONSE" | tail -1)
  assert_status "delete own post" 200 "$STATUS"
fi

echo ""

# =============================================
# Summary
# =============================================
echo "========================================="
echo " Results: $PASS passed, $FAIL failed out of $TOTAL tests"
echo "========================================="

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
exit 0
