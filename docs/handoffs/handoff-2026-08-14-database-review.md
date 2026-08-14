# Handoff · 2026-08-14 · 数据库审查交接

> 新会话建议先读本文件，然后说："接着做"。

## 1. 现在做了什么

- 本对话已收尾：MUI-05/06/07 Discourse 移动端打磨全部完成并上线（本次对话 MUI-07 收口）。
- 用户下一步准备对数据库进行一次审查，在新对话中进行；本文件即为此交接。
- 当前分支 codex/discourse-rebuild，HEAD 29b650c（已推送），工作区干净。
- 注意：AGENTS.md 的项目介绍仍描述 forum-service 与三 schema，但**与仓库实际不符**，审查时以本文件和仓库实际为准。

## 2. 已经完成了什么

### 本对话完成项（MUI，与 DB 间接相关）
- 972ddad /top 日期筛选控件字号修复；ae26f0 MUI-07 溢出修复 + 主题同步工具加固；29b650c MUI-07 剩余表面审计（20 场景全绿）。
- Discourse 主题同步工具链：discourse-themes/tools/（backup/sync/restore），部署记录在 docs/discourse-theme-deploy-runbook.md。

### 数据库事实盘点（2026-08-14 现查，供审查直接用）
- 现役主库：PostgreSQL 16（docker compose 服务名 postgres），库名 i_forum，用户 i_forum（模板见 .env.example；本地 .env 有真实值，勿写入任何文档）。
- 现役 schema 只有两个：schema_auth（user-service）、schema_admin（admin-service）；由 migrations/init/000_init_schemas.up.sql 在容器首次初始化时经 docker-entrypoint-initdb.d 创建。
- 现役 Go 服务只有 services/user 与 services/admin；**services/forum 已不在仓库**。旧 forum-service 代码在 eference/ai-forum-gitee/services/forum/（含迁移 001_boards_posts / 002_comments_interactions / 003_attachments），仅历史参考，勿当现役。
- 各服务迁移（golang-migrate 编号，up/down 成对）：
  - services/user/migrations/：001_users_table、002_invite_codes_table、003_operation_logs、004_demo_accounts、005_profile_fields、006_oauth_identities、007_demo_batch_accounts、008_member_qq_records(+seed)、**012_qq_number_login_accounts、013_user_follows**（009–011 编号缺失，需对照 schema_migrations 实际应用版本）。
  - services/admin/migrations/：001_sensitive_words_seed、002_admin_tables、003_statistics、004_demo_user_roles、005_demo_batch_roles、006_demo_batch_roles_fixup、007_allow_level_zero_post_comment（名字带论坛语义，需确认归属 schema_admin 的合理性）。
- 社区内容已迁移到独立 Discourse 实例（直连 http://122.51.233.225:8080，前端入口 :8888），其 PostgreSQL 在 Discourse 容器栈自带数据卷中，**与 ai_forum 的 PG 是两套库**，审查时勿混淆。
- 相关参考文档：shared/api-contract.md、docs/fullstack-integration-audit-plan.md、docs/项目介绍与使用指南-领域驱动设计视角.md、DEPLOY.md、docs/demo-accounts.md。

## 3. 卡在哪里

- 本次交接本身无阻塞；但审查开始前有三个事实需要先确认（建议列入审查第一步）：
  1. schema_auth 与 schema_admin 的 schema_migrations 已应用版本 vs 迁移文件编号（重点：user 的 009–011 跳号）。
  2. 生产 PG 中是否还残留旧 schema_forum（老 forum-service 时代数据），是否属于审查/清理范围。
  3. 审查范围是否包含 Discourse 自己的 PostgreSQL（如需，走 Discourse 容器 pp 的 psql，与本项目 PG 分开）。

## 4. 下一步做什么

1. 新对话首句："读 docs/handoffs/handoff-2026-08-14-database-review.md，然后接着做"。
2. 与用户确认审查范围（建议默认：现役 i_forum 的 schema_auth/schema_admin + 迁移链一致性；可选扩展：Discourse PG、历史 forum schema 残留）。
3. 连接方式（二选一）：
   - 服务器：ssh liyongquan@122.51.233.225（密钥直连，有 Docker 权限）→ docker compose ps 确认 PG 容器 → docker exec -it <postgres容器> psql -U ai_forum -d ai_forum。
   - 本地：docker compose up -d postgres 后同样 psql（注意本地库是全新初始化，不一定有生产数据）。
4. 审查清单建议：schema/表清单 vs 迁移文件；schema_migrations 版本一致性（user 跳号、admin 007）；表结构（约束/索引/外键/默认值）；权限与角色（三个 app 用户是否共用 ai_forum）；敏感数据存储（密码哈希、OAuth/QQ 记录、邀请码）；Redis 键与 PG 一致性；备份与恢复策略（pgdata volume、服务器端备份现状）。
5. 产出：审查报告 + 问题清单，放 docs/；如要修复，按 AGENTS.md 规则先出计划书、用户确认后再动，且生产库变更前必须备份。

## 5. 哪些坑不要再踩

- 本机 PowerShell 5.1：写文件必须 UTF-8 无 BOM + LF（[System.IO.File]::WriteAllText(path, text, (New-Object System.Text.UTF8Encoding(False)))）；Get-Content 默认 GBK 读无 BOM UTF-8 会乱码（用 -Encoding UTF8 或 ReadAllText）；远程 SSH 命令只用单引号、嵌套用 ''；ConvertTo-Json 序列化长字符串会损坏（长内容走 raw 文件传输）。
- 本会话 pply_patch 工具不可用（WindowsApps 访问被拒），改文件用 PowerShell 脚本或按行改写。
- git 纪律：提交前设 $env:DEV_RECORD_SKIP='1'；每 Ticket 一次 commit，格式 <type>: <description> (#Ticket)；提交后立即 push codex/discourse-rebuild；严禁直接动 master；work/ 已被 gitignore。
- AGENTS.md 规则：不确定先问；用户质疑先判断合理性再生成计划书；调 skill 先说明原因；批量删文件/上传隐私/窃取账号等操作立即停止。
- 数据库特有：不要把 eference/ai-forum-gitee 当现役代码；Discourse PG 与 ai_forum PG 两套勿混淆；schema_forum 不在现役 compose 初始化里；生产库任何 DDL/DML 前先备份。
## 6. 后续（2026-08-14 审查已完成）
- 审查报告：`docs/database-review-2026-08-14.md`（P1/P2/P3 问题清单 + DB-FIX-01~04 修复 Ticket 草案）。
- 下一步：等待用户确认修复计划；确认后按 DB-FIX-01 起串行执行，生产库变更前必须先备份。
- 审查全程只读，未对生产库执行任何 DDL/DML。