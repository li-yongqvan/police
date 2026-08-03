# Handoff · Discourse 重建 — Ticket 1 执行完成
> 日期：2026-07-29 | 服务器：122.51.233.225 | 分支：master | 状态：✅ Ticket 1 完成

---

## 1. 任务概述

在 Discourse 重建之前，将云服务器生产环境所有未提交的修改和未跟踪文件纳入 git，打 tag 锁定基线，并记录当前运行的 Docker 容器状态。

对应 spec：`.scratch/discourse-forum-rebuild/spec.md` → Ticket 1: Production Baseline Lock

---

## 2. 执行过程

| Step | 操作 | 结果 |
|---|---|---|
| 1 | `git add -A` 暂存全部文件 | 212 files staged，无 .env/日志/tarball 被误暂存 |
| 2 | `git commit -m "chore: production baseline before Discourse rebuild"` | commit `a05a575` |
| 3 | `git tag -a pre-discourse-baseline` | annotated tag 已创建 |
| 4 | `docker compose images` + `docker compose ps` → 快照文件 | `docs/deploy/baseline-snapshot-2026-07-29.md` |
| 5 | 验证 git status / log / tag / snapshot 文件 | 全部通过，追加提交 snapshot 文件为 `cf7f860` |

---

## 3. Git 提交

| Commit | 描述 |
|---|---|
| `a05a575` | chore: production baseline before Discourse rebuild |
| `cf7f860` | docs: add production baseline snapshot 2026-07-29 |

Tag: `pre-discourse-baseline` → `a05a575`

---

## 4. 服务器当前状态

- **IP**：122.51.233.225
- **项目路径**：/home/liyongquan/projects/ai-forum
- **分支**：master
- **HEAD**：cf7f860
- **工作区**：干净（git status --short 为空）

### 运行中的容器

| 容器 | 端口映射 | 运行时间 |
|---|---|---|
| admin-service | 0.0.0.0:8003->8003 | 47h+ (healthy) |
| forum-service | 0.0.0.0:8002->8002 | 47h+ (healthy) |
| frontend | 80/tcp (internal) | 39h+ |
| nginx | 0.0.0.0:8091->80, 0.0.0.0:8888->80 | 39h+ |
| postgres | 5432/tcp (internal) | 2d+ (healthy) |
| redis | 6379/tcp (internal) | 2d+ (healthy) |
| user-service | 0.0.0.0:8001->8001 | 47h+ (healthy) |

详细信息见：`docs/deploy/baseline-snapshot-2026-07-29.md`

---

## 5. 未执行事项

- **未 push**：所有 commits 和 tag 仅在服务器本地，未推送到 GitHub/Gitee（按 Constraints 要求）
- 未修改任何运行中的容器
- 未修改 .env 或凭证文件
- 未切换分支

---

## 6. 下一步

**Ticket 2: Discourse Proof-of-Concept Deployment**

编译完成的 goal.md 位于：`.scratch/discourse-forum-rebuild/goal.md`

新执行线程启动指令：
> 读取 `.scratch/discourse-forum-rebuild/goal.md`，按 Execution order 逐步执行。

### 关键约束提醒
- Discourse 使用独立端口（如 8080），不影响现有服务
- Discourse 使用自己的 PostgreSQL，不碰现有 schema
- 不 push 到远程仓库
- 服务器资源有限，需先检查内存和磁盘
