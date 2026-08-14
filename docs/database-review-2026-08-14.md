# 数据库审查报告（2026-08-14）

> 范围：生产 `ai_forum` 库的 `schema_auth` / `schema_admin` + 迁移链一致性（默认范围，不含 Discourse PG 与历史 `schema_forum` 数据，仅记录其残留）。
> 方式：SSH 到 `122.51.233.225` → `docker exec` 进入 `ai-forum-postgres-1` 以容器内 `ai_forum` 角色执行 psql。
> 原则：全程只读，未执行任何 DDL/DML；未将任何凭据写入本文档。

---

## 1. 结论摘要

| # | 级别 | 问题 | 一句话影响 |
|---|------|------|-----------|
| P1-1 | 高 | `schema_auth_migrations` 处于 `version=13, dirty=true`，且 `user_follows` 表在生产不存在 | user-service 每次启动只警告不迁移，后续任何新迁移都无法自动应用；关注/粉丝 API 一调用即 500 |
| P1-2 | 高 | 运行中的 user-service 镜像内迁移文件与仓库不一致（多出一个残缺的 `009_qq_number_login_accounts.up.sql`，无 down） | 镜像重构建后行为不可预测，迁移链不可复现 |
| P1-3 | 高 | 生产库没有任何数据库备份（无 pg_dump 产物、无定时任务） | 任何误操作/磁盘故障都可能导致数据丢失 |
| P1-4 | 高 | 两个服务共用 `ai_forum` 角色连库，且该角色是**超级用户**（rolsuper=t），三 schema 无权限隔离 | 任一口子被突破即等于整库被接管 |
| P2-1 | 中 | 迁移编号 010/011 缺口无法从仓库溯源（009 在 `72798d0` 被改名 012，010/011 从未进过仓库，但生产已到 13） | 迁移链不完整，只能靠 dirty 状态推断历史 |
| P2-2 | 中 | `public.schema_forum_migrations` 残留（`version=13, dirty=true`），forum 时代死表 | 混淆迁移状态，新环境初始化语义不清 |
| P2-3 | 中 | 外键缺失：`user_roles.user_id`（bigint，与 `users.id` integer 类型不一致）、`invite_codes.created_by/used_by`、`operation_logs.operator_id` 均无 FK | 跨表数据一致性只靠应用层保证 |
| P2-4 | 中 | 冗余索引 3 组（username、code、stat_date 各有"唯一+非唯一"双索引） | 每次写入多维护一套索引 |
| P2-5 | 中 | `pg_hba` 对 loopback 全 `trust`、`ssl=off`、`listen_addresses=*`（5432 未对外发布） | 当前暴露面小（容器网内），但任何同网容器/主机进程可免密连库 |
| P2-6 | 中 | 真实姓名 + QQ 号明文存储（`member_qq_records` 53 行、`oauth_identities.provider_user_id` 14 行） | 个人隐私数据无加密/脱敏/归档策略 |
| P3-1 | 低 | `statistics_daily` 空表（0 行）且当前 Go 代码已无任何写入 | 死表，迁移 003 遗留 |
| P3-2 | 低 | `system_config` 保留 11 个 forum 时代键（board_*/post_*/comment_*/upload_*），社区已迁 Discourse | 配置语义过期，易误导运营 |
| P3-3 | 低 | 生产仍活跃 21 个演示账号（demo*/admin0*/plat0*，全部 active） | 试运行后未下线，弱口令风险敞口 |
| P3-4 | 低 | 磁盘 40G 已用 78%（可用 8.6G）；Redis 默认 RDB 无 AOF | 备份前需注意空间；Redis 重启有会话回滚风险（可接受） |

---

## 2. 环境快照

- 容器：`ai-forum-postgres-1`（postgres:16，healthy）、`ai-forum-redis-1`（redis:7-alpine）、`ai-forum-user-service-1`、`ai-forum-admin-service-1`、前端/Nginx。Discourse `app` 容器在默认 bridge 网，与 PG 所在 `ai-forum_ai-forum-net` **网络隔离**，互不可达。
- 库：`ai_forum`，owner=`ai_forum`（唯一可登录角色，superuser）。
- 数据量：全库 9 MB；users 76 行（level 2×14 / level 5×62）、invite_codes 1、oauth_identities 14、member_qq_records 53、roles 2、user_roles 14、sensitive_words 10、statistics_daily 0。
- 扩展：`pgcrypto`（public，1.3）、`plpgsql`。
- 序列正常：`users.id` max=182 与序列一致。
- Redis：16 个键 = `refresh:<userId>`（TTL 约 12–30 天，滑动会话）+ `rl:login:ip:<ip>`（限流，短 TTL）。
- 备份：`~/backups/`、`~/deploy-backups/` 内均为**代码**备份（2026-07-14 / 07-20），**无任何数据库备份**；crontab 为空。

## 3. 迁移状态明细

| 迁移表（均在 public） | version | dirty | 说明 |
|----------------------|---------|-------|------|
| `schema_auth_migrations` | 13 | **true** | user-service 启动日志：`Migration warning: failed to run migrations: Dirty database version 13. Fix and force version.`（2026-08-11 起每次启动都出现） |
| `schema_admin_migrations` | 7 | false | 干净，与仓库 001–007 一致 |
| `schema_forum_migrations` | 13 | true | forum-service 时代死表，schema_forum 本身已不存在 |

**仓库迁移文件（services/user/migrations/）**：001–008、012、013。009 在提交 `72798d0`（2026-07-16 "fix: avoid qq login migration version conflict"）被改名 012；**010/011 从未出现在仓库历史**，但生产迁移版本已到 13，说明 010/011 曾在生产被应用过（内容不可考，疑为当时部署目录下的临时迁移，后被删除）。

**运行中镜像的迁移文件 ≠ 仓库**：镜像内多了 `009_qq_number_login_accounts.up.sql`（且**没有对应 down**），其内容相当于仓库 012 的 INSERT 后半段。即镜像构建自一次"改名不彻底"的中间状态。

**已确认在生产生效的迁移**：
- 004/007（演示账号数据，21 个账号在库）✅
- 005（users 表有 department/squad/grade/profile_completed 四列）✅
- 008（member_qq_records 53 行）✅
- 012（51 个 QQ 号用户名、密码统一 bcrypt('123456')）✅
- 013（user_follows 建表）❌ **表不存在**
- admin 007（post/comment_requires_level 已置 '0'，2026-07-20）✅

**影响面**：`user_follows` 表缺失但代码仍引用——`services/user/internal/service/user_service.go:519-652` 的 Follow/Unfollow/IsFollowing/GetFollowing/GetFollowers/FollowCounts 全部直查该表；`frontend/src/api.js:202-245` 已封装对应 API（当前无 .vue 视图接入，但 API 一经调用即 500）。

## 4. 结构 / 权限 / 安全明细

- 表结构：所有表都有主键；**全库仅 2 个外键**（`oauth_identities.user_id → users.id`、`user_roles.role_id → roles.id`）；CHECK 约束只有 NOT NULL，无取值域约束（如 level/status 无枚举）。
- 权限：只有 `ai_forum` 一个可登录角色，且 superuser；user-service 与 admin-service 均用 `DB_USER=ai_forum / DB_NAME=ai_forum`，仅靠应用层 `DB_SCHEMA` 区分。
- `pg_hba.conf`：`local all all trust`、`host 127.0.0.1/::1 trust`、`local/host replication trust`、`host all all all scram-sha-256`；`ssl=off`、`listen_addresses=*`。5432 端口未发布到宿主，风险面当前限于容器网络与宿主机进程。
- 密码存储：76 个用户全部 bcrypt（`$2a$10$` 前缀）✅；未发现明文密码列。
- OAuth：`raw_profile`（jsonb）仅含 name/qq_nickname/qq_number 三个键，**不含 token** ✅。
- PII：`member_qq_records`（真实姓名+QQ号，53 行）与 `oauth_identities.provider_user_id`（QQ 号，14 行）明文保存，无脱敏/归档/删除策略。
- 邀请码：明文 `varchar`，现仅 `DEMO2026` 一条且已 used。
- Redis 键面干净：仅 refresh 会话与登录限流，无异常键。

## 5. 修复计划（草案，待确认后执行）

> 按 AGENTS.md 规则：确认前不动生产库；执行前必须先备份。建议按 4 个 Ticket 串行推进。

**DB-FIX-01 迁移链修复 + 镜像对齐（P1-1/P1-2/P2-1）**
1. `pg_dump` 全库备份并验证。
2. `UPDATE schema_auth_migrations SET version=12, dirty=false`（即 golang-migrate 的 force 12），重启 user-service，由其自动应用 013 建表，恢复 dirty=false。
3. 验证 `schema_auth.user_follows` 存在、启动日志无警告、follow API 正常。
4. 用仓库当前迁移目录重建 user-service 镜像（自动移除残缺 009.up），再发布。

**DB-FIX-02 备份机制（P1-3）**
1. 建立每日 `pg_dump` cron（如 02:30，输出到 `~/backups/db/`，保留 30 天，磁盘占用约 <10 MB/天）。
2. 做一次恢复演练（在本地 docker PG 上 restore 验证可恢复性）。

**DB-FIX-03 降权与网络加固（P1-4/P2-5）**
1. 新建 `ai_forum_user`、`ai_forum_admin` 两个非超级用户，分别授予各自 schema 的 CRUD + sequence 权限（user-service 不跨 schema，admin-service 经内部 API 访问 user，均已核实）。
2. 修改服务器 `.env`/compose 环境切换服务账号；`ai_forum` 保留为 DBA 账号。
3. `pg_hba`：删除 loopback trust，统一 scram；评估 ssl（内网 + 不发布端口前提下可暂缓）。

**DB-FIX-04 清理与一致性（P2-2~P2-6/P3-1~P3-3）**
1. 备份后 `DROP TABLE public.schema_forum_migrations`。
2. 删除 3 组冗余索引。
3. `user_roles.user_id` bigint→integer 对齐并补 FK；评估 `invite_codes.created_by/used_by`（存在系统值 0）与 `operation_logs.operator_id` 补 FK。
4. 评估 `statistics_daily` 与 forum 时代 system_config 键的去留（Discourse 侧统计已独立）。
5. 与用户确认后下线/清理 21 个演示账号或改强密码。
6. 评估 `member_qq_records` 脱敏/归档策略（试运行结束后）。

## 6. 待确认事项

1. 上述 4 个 Ticket 的拆分与顺序是否认可（尤其 DB-FIX-03 需要改服务器 .env，涉及服务重启窗口）。
2. 迁移 010/011 的历史内容是否需要继续溯源（可查旧镜像层，但收益有限，建议接受现状并记录）。
3. 演示账号：下线 / 保留 / 改密码？
4. 修复在什么时间窗执行（当前论坛有真实流量，登录/SSO 依赖 user-service）。