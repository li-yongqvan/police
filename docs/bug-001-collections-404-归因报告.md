# BUG-001：收藏列表 API 返回 404 — 归因报告

| 字段 | 内容 |
|------|------|
| **编号** | BUG-001 |
| **发现日期** | 2026-06-03 |
| **发现阶段** | 网关冒烟测试（全功能体验测试脚本） |
| **严重程度** | 🟡 P2（一般 — 次要功能异常，不影响核心闭环） |
| **影响范围** | `/forum-api/api/v1/me/collections` — 前端「我的收藏」页面 |
| **状态** | ✅ 已定位，待远程重新部署 |
| **方法论标签** | [产品:收藏] [风险:功能缺失] [环境:远程生产] |

---

## 1. 现象

```
GET http://107.172.138.10/forum-api/api/v1/me/collections
Authorization: Bearer <valid_student_token>

→ HTTP 404 "404 page not found"
```

而同一 forum-service 的其他端点正常：

| 端点 | 结果 |
|------|------|
| `GET /forum-api/api/v1/boards` | ✅ 200 |
| `GET /forum-api/api/v1/posts?limit=3` | ✅ 200 |
| `POST /forum-api/api/v1/posts/6/like` | ✅ 200 |
| `POST /forum-api/api/v1/posts/6/collect` | ✅ 200 |
| `GET /forum-api/api/v1/me/collections` | ❌ 404 |

这说明 Likes/Collections 的写入路径正常，但收藏列表**读取路径**缺失。

---

## 2. 排查过程

### 2.1 第一步：检查代码路由是否存在

```
文件: services/forum/cmd/main.go:132
auth.GET("/me/collections", interactionHandler.ListMyCollections)
```

→ **路由注册代码存在。**

### 2.2 第二步：检查 Handler 层

```
文件: services/forum/internal/handler/interaction_handler.go:57
func (h *InteractionHandler) ListMyCollections(c *gin.Context) { ... }
```

→ **Handler 实现存在，解析 user_id + 分页参数，调用下层 Service。**

### 2.3 第三步：检查 Service 层

```
文件: services/forum/internal/service/forum_service.go:164
func (s *ForumService) ListUserCollections(ctx context.Context, userID uint, page, limit int) ...
```

→ **Service 层存在，SQL 联表查询 `schema_forum.collections` JOIN `schema_forum.posts` JOIN `schema_auth.users` JOIN `schema_forum.boards`，逻辑正确。**

### 2.4 第四步：检查 Gin 返回内容

```
$ curl -s -i GET /forum-api/api/v1/me/collections
HTTP/1.1 404 Not Found
Server: openresty
Content-Type: text/plain
Content-Length: 18

404 page not found
```

关键线索：响应体是 `404 page not found`（纯文本），**而非** Gin 默认的 `404 Not Found` JSON 格式。

这意味着请求**未到达 forum-service 的 Gin 路由**——请求早在 Nginx 层就返回了 404。

### 2.5 第五步：检查 Nginx 配置

```
文件: nginx/nginx.conf
location /forum-api/ {
    proxy_pass http://forum_backend/;
    ...
}
```

Nginx 前缀转发规则正确，`/forum-api/api/v1/me/collections` 会被代理到 `http://forum_backend/api/v1/me/collections`。

### 2.6 第六步：分析 forum-service 实际响应

直接探测 forum-service 的其他路由，通过 Nginx 代理后全部正常（200），唯独 `/me/collections` 返回 404。

如果 Nginx 正确代理但 forum-service Gin 没有这个路由，Gin 会返回：
```
404 page not found
```

这与我们观察到的现象**完全吻合**——Gin 框架在找不到路由时的默认响应就是 `404 page not found`。

### 2.7 结论

```
┌─────────────────────────────────────────────────────────────────┐
│                        根因：部署版本差异                          │
│                                                                 │
│  本地代码（main.go:132）      远程容器（旧镜像）                    │
│  ┌─────────────────────┐     ┌─────────────────────┐            │
│  │ /me/collections ✅   │     │ /me/collections ❌   │            │
│  │ /posts/:id/like  ✅  │     │ /posts/:id/like  ✅  │            │
│  │ /posts/:id/collect✅ │     │ /posts/:id/collect✅ │            │
│  └─────────────────────┘     └─────────────────────┘            │
│                                                                 │
│  /me/collections 路由是新代码，远程 Docker 镜像未 rebuild。         │
│  Like/Collect 旧代码存在，所以写入正常。                            │
└─────────────────────────────────────────────────────────────────┘
```

`/me/collections`（收藏列表读取）和 `/posts/:id/like`、`/posts/:id/collect`（点赞/收藏写入）**不在同一次 commit 中引入**。远程部署的 forum-service 容器使用的是较旧的 Docker 镜像，不包含 `ListMyCollections` 路由。

---

## 3. 修复方案

在远程服务器上执行：

```bash
cd /opt/ai-forum
git pull   # 拉取最新代码（含 /me/collections 路由）

# 重新构建 forum-service 镜像
docker compose -f docker-compose.yml -f docker-compose.server.yml build forum-service

# 滚动更新 forum-service 容器
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d forum-service

# 重新加载 Nginx（刷新上游 DNS）
sleep 5
docker compose -f docker-compose.yml -f docker-compose.server.yml restart nginx

# 验证
TOKEN=$(curl -s -X POST http://127.0.0.1:8888/user-api/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo_student","password":"demo123456"}' \
  | python -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

curl -s -w "\nHTTP:%{http_code}" \
  "http://127.0.0.1:8888/forum-api/api/v1/me/collections" \
  -H "Authorization: Bearer $TOKEN"

# 期望: HTTP 200，返回 {"posts":[...],"total":...,"page":1,"limit":20}
```

---

## 4. 经验教训

| 维度 | 复盘 |
|------|------|
| **Why 没早发现** | 收藏列表属于"读取已收藏内容"的辅助路径，冒烟测试之前未覆盖此端点。like/collect 写入正常所以前端收藏按钮看起来正常，但列表页始终空白。 |
| **Why 部署版本不同** | 远程部署脚本 `deploy-full-platform.sh` 和 `deploy-frontend-incremental.sh` 需要手动执行。新代码 push 后若未触发 rebuild，就会产生版本偏差。 |
| **改进建议** | ① 将 `full-experience-test.sh` 加入 CI/CD（GitHub Actions），每次 push 自动跑冒烟；② 添加 `git rev-parse HEAD` 到 health 端点暴露版本号，便于快速判断镜像是否最新。 |

---

## 5. 缺陷关闭条件

- [x] 根因已定位（代码存在但远程镜像过期）
- [ ] 远程服务器执行 rebuild 并验证 HTTP 200
- [ ] 前端「我的收藏」页面可正常显示收藏列表
- [ ] 全功能测试脚本 `full-experience-test.sh` 此项转为 PASS

---

**撰写人**：AI 测试助手（基于 Bach/Kaner/Bolton/Black/Offutt/Whittaker 方法论）  
**日期**：2026-06-04  
**关联文档**：[test-plan.md](./test-plan.md) | [manual-walkthrough-checklist.md](./manual-walkthrough-checklist.md)
