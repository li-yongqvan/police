# 内测前收尾计划书 · Pre-Pilot Closure

Date: 2026-08-17
状态：已获用户确认（决策点 1~4 全部选 A），执行中
分支：codex/discourse-rebuild
关联：docs/database-review-2026-08-14.md · docs/handoffs/handoff-2026-08-14-test-results.md · docs/plans/2026-08-14-server-memory-relief-plan.md · docs/pilot-acceptance-checklist.md

## 1. 背景与结论

- 项目门禁（docs/todo-list.md 灰度节奏）：**内测 5～10 人（教师+骨干）的触发条件 = 所有 P0/P1 关闭**。
- P0（Discourse staff 全站 500）已于 2026-08-17 修复并复测通过；剩余未关闭项集中在：数据库三项 P1（DB-FIX-01/02/03 未执行）、P1-2 审核闭环定案、生产部署基线核对。
- 重要发现：仓库代码在 2026-07-30/31（Ticket #5/#6/#7）已移除 admin-service 对 forum-service 的调用、清理概览统计接口、隐藏未实现菜单（侧边栏现仅「管理首页/用户管理/邀请码/前往 Discourse 管理」）。但 08-14 生产测试仍复现 `127.0.0.1:8002 connection refused` → 说明**生产部署的 admin-service / 前端镜像落后于仓库**。故内测前必须先核对部署基线，避免"代码已修、线上未更新"。
- 目标：完成本计划书后，按 docs/pilot-acceptance-checklist.md 组织 5～10 人内测。

## 2. 现状盘点

| 项 | 状态 | 说明 |
|----|------|------|
| P0 Discourse staff 500 | ✅ 已关闭（08-17） | 修复+复测通过；24h 内存观察未完成 |
| P1-1 forum-service 死调用 | ⚠️ 代码已清理，生产待核对 | 仓库 07-30 已移除，生产镜像疑似旧版 |
| P1-2 内容审核闭环 | ⚠️ 待定案 | 前端已隐藏菜单（走「下沉 Discourse」路线），后端 ListPendingAudit 仍为 TODO |
| DB P1-1 迁移链 dirty + user_follows 缺失 | ❌ 未执行（DB-FIX-01） | 内测会产生真实数据，必须修复 |
| DB P1-2 镜像迁移文件不一致 | ❌ 未执行（DB-FIX-01 第 4 步） | 重建镜像时自动消除 |
| DB P1-3 无数据库备份 | ❌ 未执行（DB-FIX-02） | 数据安全底线，必须在内测前完成 |
| DB P1-4 共用超级用户 | ❌ 未执行（DB-FIX-03） | 建议内测前做降权（或内测后尽快补） |
| P2/P3 体验项 | ⏸️ 不阻塞内测 | 见第 7 节清单，内测后迭代 |

## 3. Ticket 拆解（串行执行，每 Ticket 一 commit）

### Ticket 0 · 生产部署基线核对（前置，先于一切改动）

状态：✅ 已完成（2026-08-17）

目标：确认线上 admin-service / 前端是否已是仓库当前版本（含 07-30 Ticket #5/#6/#7 清理）。

动作：
1. SSH 到 122.51.233.225，核对部署版本：
   - `docker ps` 看 admin-service/frontend 镜像名与创建时间；
   - `docker inspect` 镜像标签/构建时间，与仓库 HEAD 时间（2026-08-17）对比。
2. 功能验证：demo_platform_admin 登录 :8888 → /admin → 确认「管理首页」不再请求 stats/overview/boards/audit/pending（浏览器 Network 面板）。
3. 若镜像为旧版 → 用仓库最新代码重建 admin-service 与前端镜像并发布（见 Ticket 4 部署流程）。

产出/验证：管理首页无 8002 报错；无 stats 请求。

核对结果（SSH 122.51.233.225 实测）：
- admin-service 镜像构建于 2026-07-22（早于 Ticket #5/#6/#7 的 07-30 清理），容器二进制含 forum-service×19、stats/overview×2、audit/pending×1、8002×2 → P1-1 线上确认为真，需重建部署。
- 前端镜像构建于 2026-08-12，AdminOverview JS 中 stats/overview=0、audit/pending=0、8002=0 → 前端清理已生效，无需改动。
- 服务器项目目录 /home/liyongquan/projects/ai-forum，git 分支 codex/server-deploy-mui01-20260803-172000（HEAD 8b915a9），落后于仓库；部署构建上下文为 ./services/admin、./services/user（docker-compose.server-full.yml）。
遗留风险：旧镜像若同时存在其他未同步改动，以 git 历史为准重新部署。
是否影响发布：核对本身无影响；若需重部署，论坛中断约 1～3 分钟。

### Ticket 1 · DB-FIX-01 迁移链修复 + user_follows 建表

状态：✅ 已完成（2026-08-17）

目标：生产库迁移链恢复干净，013（user_follows）落表，镜像迁移文件与仓库一致。

前置：执行任何 DDL/DML 前必须 pg_dump 全库备份并验证备份文件。

动作（生产库，全程先在测试库演练）：
1. `pg_dump` 全库备份到 `~/backups/db/`，验证文件完整（\q 后 `pg_restore --list` 或 gzip -t）。
2. `UPDATE public.schema_auth_migrations SET version=12, dirty=false`（等价 golang-migrate force 12）。
3. 重启 user-service，由其自动应用 013 建表，确认 `schema_auth.user_follows` 存在、启动日志无 Migration warning。
4. 验证 follow API：`GET /api/v1/users/me/following` 等不再 500（frontend/src/api.js 已封装，202-245 行）。
5. 用仓库当前 services/user/migrations/ 重建 user-service 镜像（自动移除残缺 `009_qq_number_login_accounts.up.sql`），发布。

验证方式：psql 查 user_follows 表结构；user-service 启动日志无 dirty 警告；follow API 200。

执行记录（SSH 122.51.233.225 实测）：
1. 备份 ~/backups/db/ai_forum_pre_pilot_20260817_154934.dump.gz，gzip -t 通过，pg_restore -l 可列出全表。
2. 首次 force 12 后重启遇阻：013_user_follows.up/down.sql 与 008 down 均带 UTF-8 BOM，psql 报 syntax error near 0xEFBBBF，迁移失败回滚为 13|t。
3. 已修复：本地仓库 3 文件去 BOM + LF，scp 同步服务器；删除服务器残缺 009_qq_number_login_accounts.up.sql；重建 ai-forum-user-service 镜像。
4. 再次 force 12 → up -d user-service → 日志 Database migrations applied successfully，迁移链 13|f。
5. user_follows 表完整（PK/唯一约束/CHECK 防自关注/FK 级联删除/两索引）；/user-api/api/v1/users/me/following 登录后返回 200 count=0，不再 500。
6. 遗留提示：镜像内 008/013 相关 BOM 修复已随新镜像生效；本地仓库 3 个迁移文件已同步修复待 commit（随本 Ticket 提交）。
风险/回滚：生产库变更前已有备份；若迁移失败，恢复备份并保持 dirty 状态不动。
是否影响发布：user-service 重启约 10～30 秒，登录短暂中断，建议低峰执行。

### Ticket 2 · DB-FIX-02 备份机制 + 恢复演练

状态：✅ 已完成（2026-08-17）

目标：建立每日自动备份，并证明可恢复。

动作：
1. 服务器 crontab 增加每日 02:30 `pg_dump`（输出 `~/backups/db/`，保留 30 天，预估 <10MB/天；先确认磁盘空间，当前 40G 已用 78%）。
2. 执行一次备份 → 在本地 Docker PG 上 restore → 验证可读。
3. 将备份/恢复命令沉淀为脚本（scripts/ 或 /opt/ai-forum/scripts/），写回 docs（DEPLOY.md 或 pilot-runbook.md）。

验证方式：crontab -l 可见；备份文件存在且 gzip -t 通过；本地 restore 后表行数一致。

执行记录（SSH 122.51.233.225 实测）：
1. 改进 scripts/backup-postgres.sh：修正项目根为 /home/liyongquan/projects/ai-forum（原 /opt/ai-forum 不存在）、pg_dump 改 -Fc custom 格式、默认 BACKUP_DIR=~/backups/db、RETAIN_DAYS=30、cron PATH 补 docker compose 插件。
2. 用户级 crontab 已装：每日 02:30 执行并写 backup.log（原 crontab 为空，历史 install-pilot-cron.sh 指向的 /opt/ai-forum 路径失效）。
3. 手动执行备份：ai_forum_20260817-163638.dump.gz（16K），gzip -t 通过，pg_restore -l 列出 30 表。
4. 恢复演练：pg_restore 到临时库 ai_forum_restore_test，退出码 0，users=76、roles=2、迁移 13|f 与生产一致，随后清理临时库。
风险/回滚：纯新增机制，无回滚风险。
是否影响发布：无在线影响。

### Ticket 3 · P1-2 审核闭环定案（需用户拍板，见第 5 节决策点 1）

状态：✅ 已完成（2026-08-17，决策点 1 选 A：治理下沉 Discourse）

默认方案（与仓库现状一致）：**治理能力正式下沉 Discourse**。
- Vue 端：维持现有「管理首页」说明与「前往 Discourse 管理」外链，不新增审核页面。
- admin-service：`ListPendingAudit` 标注废弃并移除其 TODO 路由（当前已无消费方），或保留返回空列表。
- 清理孤儿数据：确认 demo02 对话题 #19 回复产生的待处理审核记录是否残留 schema_admin；若残留，在备份后由 DBA 手工清理（或随 DB-FIX-04 一并处理）。

验证方式：Vue 管理端无指向 404 的菜单；admin-service 无 pending 相关死路由；sensitive_words 接口仍可用。

执行记录：
1. 生产 ai_forum 库核查：schema_admin 仅 6 张表（operation_logs/roles/sensitive_words/statistics_daily/system_config/user_roles），schema_auth 无审核表；operation_logs 均为 ban/unban/level 治理操作 → 无孤儿审核记录，无需 DBA 清理。
2. 本地代码：删除 model/audit_record.go 与 admin_service.go 中 ListPendingAudit（TODO 无消费方），go build ./... 通过（BUILD OK）。
3. 前端维持现状（侧边栏已仅剩 管理首页/用户管理/邀请码/前往 Discourse 管理，概览页无 stats 调用，无需改动）。
风险/回滚：删除路由前确认前端无调用（已核实前端仅 moderationMode 配置读写，无 pending 列表调用）。
是否影响发布：admin-service 重建发布，短暂中断。

### Ticket 4 · 发布与验收准备

状态：✅ 主体已完成（2026-08-17）；剩余种子内容/邀请码/值班人待用户确认（见第 9 节）

动作：
1. 若 Ticket 0 判定生产为旧镜像，或本计划书产生 admin/user 镜像变更 → 统一重建并发布（docker compose build + up -d，或按 DEPLOY.md 流程），提交后 push。
2. 内测准备项核对（docs/pilot-acceptance-checklist.md「准备」）：
   - 邀请码批量生成 50～100 个（demo_platform_admin → 邀请码），导出给辅导员；
   - 演示账号处置：见决策点 3（至少 demo-login 白名单生效、演示账号改强口令或临时下线，DB P3-3 提到 21 个演示账号全 active）；
   - 种子内容：每板块 ≥3 条帖子；
   - 值班人确定；HTTPS 域名（Cloudflare Tunnel / 正式域名）确认手机 4G/Wi-Fi 可访问；
   - `BASE_URL=... bash scripts/smoke-test.sh` 通过。
3. 更新 docs/todo-list.md（当前停留在 2026-06-04，P0 BUG-001 为 forum-service 时代遗留，应归档），并勾选验收表。

验证方式：按 pilot-acceptance-checklist.md 逐项；手机端窄屏抽查登录→论坛→管理路径。

执行记录（SSH 122.51.233.225 实测）：
1. **代码整体对齐**：发现服务器项目目录 services/admin、services/user 与仓库严重不一致（服务器部署分支保留 forum 时代代码：forum_client.go/stats_service.go/audit_handler.go 等）。用 git archive HEAD 打包 services/admin+user（96 个 tracked 文件，.env 不受影响）覆盖服务器，备份先行（/tmp/services-backup-*.tgz）。
2. **重建发布**：user-service + admin-service 镜像重建并部署，双双 healthy。新 admin 二进制 forum-service=0、stats/overview=0（P1-1 关闭）；残留 8002×1 来自依赖库非业务代码。
3. **接口回归**：demo_platform_admin 登录后 /admin/users total=76、/admin/invite-codes count=2、users/me role=platform_admin 全部 200。
4. **演示账号改强口令**（决策 3 A）：21 个 demo*/admin0*/plat0* 统一改随机 16 位强口令（bcrypt cost 12，UPDATE 21）；验证新口令登录 OK、旧密码 123456 401；口令已单独交付用户，未入库。
5. **安全项**：POST /api/v1/demo-login 公网 403（白名单生效）。
6. **域名**：https://api.shgenren.dpdns.org 服务器 curl 200、本机 30s 内 200（首次握手慢，可达）。
7. **冒烟**：smoke-test.sh 适配当前架构（原脚本仍调 forum-api 死接口）；cron-smoke.sh 路径修正；网关模式冒烟 PASS（frontend/discourse/login/users-me）；crontab 已配 */15 冒烟 + 每日 02:30 备份。
8. **遗留**：Discourse 7 板块中 6 个 topic_count=0（仅「常规」3 条），不满足验收清单「种子帖子每板块 ≥3 条」，待用户确认是否补种子内容；邀请码批量生成与值班人待确认。
风险/回滚：发布按 DEPLOY.md 回滚预案。
是否影响发布：是，本 Ticket 即发布动作。

## 4. 执行顺序与依赖

```
Ticket 0（部署基线核对，只读）
  └─→ Ticket 1（DB-FIX-01，需用户确认 + 备份）
        └─→ Ticket 2（DB-FIX-02 备份机制）
              └─→ Ticket 3（P1-2 定案，依赖决策点 1）
                    └─→ Ticket 4（发布 + 验收准备，依赖决策点 3）
                          └─→ 24h 观察后：内测 5～10 人
```

- Ticket 0 与决策点可并行；DB 类改动严格串行，一个 commit 对应一个 Ticket。
- 内测放人前还需满足：P0 修复后 24h 内存观察无反弹（今天 08-17 起算，观察至 08-18）。

## 5. 决策点（2026-08-17 已拍板：全部选 A）

1. **P1-2 审核闭环** ✅ A：治理下沉 Discourse，Vue 端不实现审核页面；ListPendingAudit 废弃；清理孤儿审核记录。
2. **DB-FIX-03/04** ✅ A：内测前执行降权与清理（DB-FIX-03 降权 + pg_hba 加固；DB-FIX-04 清理与一致性）。
3. **演示账号处置** ✅ A：内测前改强口令并保留，demo-login 白名单确认生效。
4. **Discourse 减压** ✅ A：观察 24h 后执行 Phase C（UNICORN_WORKERS 3→2）+ Phase D（Swap 扩容）。

## 6. 不阻塞内测项（P2/P3，内测后迭代）

- 中文搜索无结果 / invalid param（P2-4，需 Discourse 分词设置 + 重建索引）
- reactions/like/voting/solved 控件缺失（P2-5，插件启用与分类设置核查）
- Vue 端无暗色模式（P2-7，可先标注「仅亮色」）
- 用户管理页无搜索（P2-6，加用户名/学号搜索）
- 邀请码作废后统计不刷新 + 状态文案混排（P2-2/P3-2）
- 论坛字体 CORS（P2-3，主题字体地址改相对路径）
- #create-topic 图标按钮无 aria-label（P3-1）
- todo-list.md 归档（BUG-001 为 forum-service 时代遗留）

## 7. 交付约定

- 每 Ticket 完成后立即 commit（格式 `<type>: <description> (#<ticket>)`）+ push origin codex/discourse-rebuild；不直接动 master。
- 所有验证通过后默认同步 GitHub 并部署云服务器（除非用户明确只做本地）。
- 本计划书任何改动生产库的动作，执行前必须先备份并向用户报备。

## 9. 待用户确认（内测放人前）
1. 种子内容：是否由我批量补充（Discourse 6 个空板块 × ≥3 条）？还是值班管理员手动填充？
2. 邀请码：内测前是否现在批量生成（建议 50～100 个）？
3. 值班人与内测启动日期确认。
4. 演示账号新口令已单独交付，请存至密码管理器。

## 8. 风险汇总

| 风险 | 等级 | 缓解 |
|------|------|------|
| DB 迁移修复失败导致服务不可用 | 高 | 先演练、先备份、可回滚到备份 |
| 生产镜像与仓库不一致未被发现 | 高 | Ticket 0 前置核对 |
| 内测产生真实数据后无备份 | 高 | Ticket 2 内测前完成 |
| 演示账号弱口令被滥用 | 中 | 决策点 3 处置 |
| Discourse 内存再爬升 | 中 | 24h 观察 + Phase C/D 待命 |