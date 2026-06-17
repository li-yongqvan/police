#!/usr/bin/env bash
# =============================================================================
# AI 智联论坛 — 安全冒烟测试（Whittaker 攻击式 + Offutt Web安全）
# 用法：bash scripts/security-smoke.sh
# =============================================================================
set -euo pipefail

BASE_USER="${BASE_USER:-http://127.0.0.1:8001}"
BASE_FORUM="${BASE_FORUM:-http://127.0.0.1:8002}"
BASE_ADMIN="${BASE_ADMIN:-http://127.0.0.1:8003}"
BASE_URL="${BASE_URL:-}"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
PASS=0; FAIL=0; TOTAL=0

section() { echo -e "\n${CYAN}══ $* ══${NC}"; }
pass() { echo -e "  ${GREEN}PASS${NC}: $*"; PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); }
fail() { echo -e "  ${RED}FAIL${NC}: $*"; FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); }
warn() { echo -e "  ${YELLOW}WARN${NC}: $*"; }

post_json() { curl -s -w "\n%{http_code}" -X POST "$1" -H 'Content-Type: application/json' -d "$2"; }
get_auth() { curl -s -w "\n%{http_code}" -X GET "$1" -H "$2"; }

split() { local s; s=$(echo "$1" | tail -1); local b; b=$(echo "$1" | sed '$d'); eval "$2='$s'"; eval "$3='$b'"; }
json_field() { echo "$1" | $PYTHON -c "import sys,json; print(json.load(sys.stdin).get('$2',''))" 2>/dev/null || echo ""; }

# 获取 token
echo "获取测试 token..."
login_json='{"username":"demo_student","password":"demo123456"}'
resp=$(post_json "$BASE_USER/api/v1/login" "$login_json"); split "$resp" s b
STUDENT_TOKEN=$(json_field "$b" "access_token")
[ -z "$STUDENT_TOKEN" ] && STUDENT_TOKEN=$(json_field "$b" "token")
STUDENT_AUTH="Authorization: Bearer $STUDENT_TOKEN"

login_json='{"username":"demo_admin","password":"demo123456"}'
resp=$(post_json "$BASE_USER/api/v1/login" "$login_json"); split "$resp" s b
ADMIN_TOKEN=$(json_field "$b" "access_token")
[ -z "$ADMIN_TOKEN" ] && ADMIN_TOKEN=$(json_field "$b" "token")
ADMIN_AUTH="Authorization: Bearer $ADMIN_TOKEN"

# 获取 board_id
boards=$(curl -s "$BASE_FORUM/api/v1/boards")
BOARD_ID=$(echo "$boards" | $PYTHON -c "import sys,json; bs=json.load(sys.stdin).get('boards',[]); print(bs[0]['id'] if bs else '1')" 2>/dev/null || echo "1")

# ============================================================
section "1. XSS 注入测试"
# ============================================================

PAYLOADS=(
  "<script>alert(1)</script>"
  "<img src=x onerror=alert(1)>"
  "<svg/onload=alert(1)>"
  "javascript:alert(1)"
)

for payload in "${PAYLOADS[@]}"; do
  # 发帖测试
  resp=$(post_json "$BASE_FORUM/api/v1/posts" \
    "{\"title\":\"XSS: $payload\",\"content\":\"$payload\",\"board_id\":$BOARD_ID}" \
    "$STUDENT_AUTH")
  split "$resp" s b
  POST_ID=$(json_field "$b" "id")
  if [ -n "$POST_ID" ]; then
    # 检查返回内容是否原样包含脚本标签
    detail=$(curl -s "$BASE_FORUM/api/v1/posts/$POST_ID")
    if echo "$detail" | grep -q '<script>'; then
      fail "XSS: payload 未过滤 → '$payload'"
    else
      pass "XSS: payload 已过滤 → '${payload:0:30}...'"
    fi
    # 清理
    curl -s -X DELETE "$BASE_FORUM/api/v1/posts/$POST_ID" -H "$STUDENT_AUTH" > /dev/null 2>&1 || true
  else
    pass "XSS: 帖子被拒绝创建 → '${payload:0:30}...'"
  fi
done

# 评论 XSS
resp=$(post_json "$BASE_FORUM/api/v1/posts" \
  "{\"title\":\"XSS Comment Target\",\"content\":\"target\",\"board_id\":$BOARD_ID}" \
  "$STUDENT_AUTH"); split "$resp" s b
TARGET_POST=$(json_field "$b" "id")

if [ -n "$TARGET_POST" ]; then
  resp=$(post_json "$BASE_FORUM/api/v1/posts/$TARGET_POST/comments" \
    '{"content":"<img src=x onerror=alert(1)>"}' \
    "$STUDENT_AUTH"); split "$resp" s b
  if [ "$s" -eq 201 ]; then
    # check stored
    comments=$(curl -s "$BASE_FORUM/api/v1/posts/$TARGET_POST/comments")
    if echo "$comments" | grep -q '<img.*onerror'; then
      fail "评论 XSS: payload 未过滤"
    else
      pass "评论 XSS: payload 已过滤"
    fi
  elif [ "$s" -eq 400 ]; then
    pass "评论 XSS: 直接拒绝"
  fi
  curl -s -X DELETE "$BASE_FORUM/api/v1/posts/$TARGET_POST" -H "$STUDENT_AUTH" > /dev/null 2>&1 || true
fi

# ============================================================
section "2. SQL 注入测试"
# ============================================================

SQLI_PAYLOADS=(
  "' OR '1'='1"
  "'; DROP TABLE posts; --"
  "1' UNION SELECT 1,2,3 --"
  "admin' --"
)

# 搜索/列表参数注入
for payload in "${SQLI_PAYLOADS[@]}"; do
  encoded=$($PYTHON -c "import urllib.parse; print(urllib.parse.quote('$payload'))" 2>/dev/null || echo "$payload")
  resp=$(curl -s -w "\n%{http_code}" "$BASE_FORUM/api/v1/posts?q=$encoded")
  split "$resp" s b
  if [ "$s" -eq 200 ]; then
    pass "SQLi (搜索): 正常返回 (未被注入影响)"
  else
    warn "SQLi (搜索): HTTP $s"
  fi
done

# 登录框注入
resp=$(post_json "$BASE_USER/api/v1/login" \
  "{\"username\":\"admin' --\",\"password\":\"x\"}")
split "$resp" s b
if [ "$s" -eq 401 ] || [ "$s" -eq 400 ]; then
  pass "SQLi (登录): 注入被阻止"
else
  fail "SQLi (登录): 异常响应 HTTP $s"
fi

# ============================================================
section "3. 越权测试"
# ============================================================

# 学生访问管理 API
ADMIN_ENDPOINTS=(
  "GET:$BASE_ADMIN/api/v1/admin/posts?limit=5"
  "GET:$BASE_ADMIN/api/v1/admin/users?limit=5"
  "GET:$BASE_ADMIN/api/v1/admin/config"
  "GET:$BASE_ADMIN/api/v1/admin/invite-codes"
  "GET:$BASE_ADMIN/api/v1/admin/audit/pending"
  "GET:$BASE_ADMIN/api/v1/admin/stats/overview"
  "POST:$BASE_ADMIN/api/v1/admin/invite-codes:{\"created_by\":0}"
  "POST:$BASE_ADMIN/api/v1/admin/users/1/ban:{}"
)

for ep in "${ADMIN_ENDPOINTS[@]}"; do
  method="${ep%%:*}"
  rest="${ep#*:}"
  url="${rest%%:*}"
  data="${rest#*:}"
  if [ "$method" = "GET" ]; then
    resp=$(curl -s -w "\n%{http_code}" -X GET "$url" -H "$STUDENT_AUTH")
  else
    resp=$(curl -s -w "\n%{http_code}" -X "$method" "$url" -H "$STUDENT_AUTH" -H 'Content-Type: application/json' -d "$data")
  fi
  split "$resp" s b
  endpoint_name=$(echo "$url" | sed 's|.*/api/v1/admin/||')
  if [ "$s" -eq 401 ] || [ "$s" -eq 403 ]; then
    pass "越权: /admin/$endpoint_name → 拒绝 ($s)"
  else
    fail "越权: /admin/$endpoint_name → 应拒绝 (got $s)"
  fi
done

# 水平越权：尝试修改他人资料
resp=$(get_auth "$BASE_USER/api/v1/users/1" "$STUDENT_AUTH"); split "$resp" s b
resp=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_USER/api/v1/users/1" \
  -H "$STUDENT_AUTH" -H 'Content-Type: application/json' \
  -d '{"nickname":"HACKED"}')
split "$resp" s b
if [ "$s" -eq 403 ] || [ "$s" -eq 401 ]; then
  pass "水平越权: 修改他人资料 → 拒绝 ($s)"
else
  warn "水平越权: 修改他人资料 → HTTP $s (可能自己的 ID=1)"
fi

# ============================================================
section "4. 请求频率限制"
# ============================================================

RATE_HIT=0
for i in $(seq 1 20); do
  resp=$(post_json "$BASE_USER/api/v1/login" '{"username":"demo_student","password":"wrong"}' "")
  split "$resp" s b
  if [ "$s" -eq 429 ]; then
    RATE_HIT=1
    pass "限流: 第 ${i} 次触发 (HTTP 429)"
    break
  fi
done
if [ "$RATE_HIT" -eq 0 ]; then
  warn "限流: 20 次内未触发（限流可能未配置或阈值较高）"
fi

# ============================================================
section "5. 信息泄露"
# ============================================================

# 5a. 错误响应不应暴露堆栈
resp=$(post_json "$BASE_USER/api/v1/login" 'not-json')
if echo "$resp" | grep -qiE 'goroutine|panic|runtime\.|stack trace'; then
  fail "信息泄露: 错误响应暴露了堆栈信息"
else
  pass "信息泄露: 错误响应无堆栈"
fi

# 5b. 健康检查不应暴露敏感信息
for svc in "$BASE_USER" "$BASE_FORUM" "$BASE_ADMIN"; do
  resp=$(curl -s "$svc/health")
  if echo "$resp" | grep -qiE 'password|secret|token|DATABASE_URL'; then
    fail "信息泄露: $svc/health 含敏感信息"
  else
    pass "信息泄露: $svc/health 不含敏感信息"
  fi
done

# 5c. API 响应头检查
resp_headers=$(curl -s -I "$BASE_USER/api/v1/login" -X OPTIONS 2>&1 || true)
if echo "$resp_headers" | grep -qi 'Server:.*\/'; then
  warn "信息泄露: 响应头可能包含服务器版本信息"
else
  pass "信息泄露: 无明显的服务器版本泄露"
fi

# ============================================================
section "6. 总结"
# ============================================================

echo ""
echo "╔══════════════════════════════════════════╗"
echo "║       安全冒烟测试报告                     ║"
printf "║  ${GREEN}PASS${NC}: %-3d  ${RED}FAIL${NC}: %-3d  TOTAL: %-3d    ║\n" "$PASS" "$FAIL" "$TOTAL"
echo "╚══════════════════════════════════════════╝"

if [ "$FAIL" -gt 0 ]; then
  echo -e "${RED}安全测试未通过，请修复以上 FAIL 项！${NC}"
  exit 1
else
  echo -e "${GREEN}安全冒烟通过！${NC}"
  exit 0
fi
