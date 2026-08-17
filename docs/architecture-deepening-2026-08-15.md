# 架构深化：授权知识收敛（候选 1）· 决策与实施记录

> 日期：2026-08-15 · 分支：codex/discourse-rebuild · 流程：improve-codebase-architecture（grilling / domain-modeling / codebase-design）
> 关联：docs/architecture-ownership-scan-2026-08-14.md（候选扫描）、docs/database-review-2026-08-14.md（DB-FIX-03 衔接）、CONTEXT.md（Authorization Language）、docs/adr/0003-authorization-ownership-and-session-token-contract.md

## 1. 目标

- 消除扫描报告中的无主共享物 O-2「角色与等级授权语义」：为授权知识建立唯一所有者，使模块变深（接口小、实现多）。
- 用户目标：项目模块化、深模块化；本候选同时是 DB-FIX-03（数据库账号降权）的前置——user-service 不再跨 schema 后才能真正降权。

## 2. 决策树共识（grilling，逐题经用户确认）

| # | 问题 | 结论 |
|---|------|------|
| Q1 | 深模块归属 | 安在 user-service（BC-Identity） |
| Q2 | 角色权威 | admin 新增内部端点返回权威角色名（优先级/回退判定留在 admin）；user 登录/SSO 调用一次 + Redis 60s 消费侧缓存；admin 不可用→上次缓存→否则 student+日志；删 admin role:* 死缓存 |
| Q3 | level 语义 | 降为纯展示属性：移出 JWT claim；删 RequireLevel 与 auth 中间件 level 上下文/DB 回查；资料与统计保留展示；值域锁定 0-5 写入契约 |
| Q4 | roles.permissions | 删除列/模型/种子；角色列表只返回 name/description；权限 = 角色名 |
| Q5 | RequireLevel | 随 Q3-A 消解（删除） |
| Q6 | token 契约 | user 唯一签发；admin 删 GenerateToken/Claims，只留 map 版 ValidateToken；claim 清单写进 api-contract.md + ADR；不抽共享 Go 包 |
| Q7 | 测试 | 按接口即测试面双侧补齐：admin 解析/端点表驱动测试；user 登录流程 mock + 降级阶梯 + token 契约测试 |
| Q8 | 部署 | admin → user 两步灰度；迁移 008 与 admin 新代码同批；完成后衔接 DB-FIX-03 降权 |

## 3. 领域模型变更（domain-modeling）

- CONTEXT.md 新增 Authorization Language：Role / Role Authority / Authoritative Role Name / Role Resolution / Level / Session Token Contract（提交 dd7cbcd）。
- 新增 ADR-0003《Authorization Ownership and Session Token Contract》：授权唯一所有者、BC-Identity 唯一签发、拒绝共享 Go JWT 包的理由、与 DB-FIX-03 的衔接。

## 4. 接口设计（codebase-design · design-it-twice，4 个并行子代理）

- Kuhn（最小接口）：2 个入口，user 业务端口无 error，fresh-first。
- Arendt（最优默认调用方）：与 Kuhn 同构，强调 ResolveAppRole 门面零改动 8 个调用点。
- Kepler（最大灵活性）：端点带 resolved_at，端口返回 RoleResolution{Source,ResolvedAt}+error，RoleSource 组合 + FailClosed 旋钮，cache-first。
- Aristotle（端口与适配器）：user 独立包 internal/role，Resolver(Name,error)+Cache 端口，GenerateToken 拒收域外角色，replace-don't-layer。

选定混合方案：

- 端点：GET /internal/v1/users/:id/role → 200 {"user_id": id, "role": "student|admin|platform_admin"}；400 非法 id；500 存储失败；永不 404（未知用户→student）。
- admin 侧：纯函数 ResolveRoleName（优先级/回退出 SQL 进 Go）+ 一条 join 查询；不建单一实现的 RoleSource 端口。
- user 侧：独立包 internal/role；业务端口 Resolver{ Resolve(ctx,id) string }（无 error，调用方零失败分支）；Fetcher + Cache 两个内部端口。
- 阶梯顺序：cache-first 读穿（60s 内不重复请求 admin）；失败→上次缓存→否则 student+日志；不做负缓存。
- token：GenerateToken 拒收域外角色；admin 保留 map-only ValidateToken。
- 测试：replace-don't-layer，新接缝测试上位后删除被取代的旧 handler 测试。

## 5. 实施计划书（Ticket 拆分，逐票 commit + push）

- ARCH-DEEP-1 · admin Role Authority 落地：role_authority.go（ResolveRoleName + ResolveAuthoritativeRole）、role_handler.GetUserRole、internal 路由、删除 role:* Redis 死缓存；表驱动 + httptest 测试。
- ARCH-DEEP-2 · admin 删除 permissions：迁移 008_drop_role_permissions、model 与 ListRoles 清理（前端零引用，已核实）。
- ARCH-DEEP-3 · user internal/role 消费模块：端口/装饰器/HTTP/Redis 适配器；删除 resolveJWTReole；GenerateToken/Claims 去 level；中间件删 level 回查与 RequireLevel；cmd/main.go 装配；阶梯/登录/token 契约测试；删除 TestLoginResponseJSONIncludesResolvedRole。
- ARCH-DEEP-4 · admin 纯消费 + 契约文档：删 admin GenerateToken/Claims；middleware 消费测试；api-contract.md 补 level 行、令牌契约与内部端点契约。

## 6. 部署顺序与风险（Q8-A）

1. 备份生产库 → 2. 发布 admin（Ticket 1+2 同批，含迁移 008）→ 3. 容器内 curl 验证新端点 → 4. 发布 user（Ticket 3+4）→ 5. 冒烟（admin 登录、降级演练）→ 6. 衔接 DB-FIX-03 降权（独立 Ticket，需确认服务器 .env 改动与重启窗口）。

风险与边界：角色变更 ≤60s 生效（与锁定 TTL 一致）；admin 全挂且缓存冷启动→登录降级 student（锁定行为）；本次不做负缓存/单飞、共享 Go 包、system_config 论坛时代键清理（O-6）、env 拆分（O-4）、admin 遗留 role:<id> 死键线上清理；前端零改动。

## 7. 状态

- 已提交并推送：CONTEXT.md + ADR-0003（dd7cbcd）、本文档（c57acf6）。
- 已完成：ARCH-DEEP-1（dcbdb4e）、ARCH-DEEP-2（c7c7918）、ARCH-DEEP-3（2815139）已提交并推送；ARCH-DEEP-4 随本票提交。
- 待办：DB-FIX-03 降权需用户确认服务器窗口；生产迁移 008 按"admin 先发、备份在前"顺序部署；跨服务联调冒烟未做。
- 待执行：Ticket ARCH-DEEP-1~4（待用户批准后逐票执行）；DB-FIX-03 衔接待服务器窗口确认。
