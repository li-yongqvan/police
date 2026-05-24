# Phase 1: 基础设施与项目骨架 - Context

**Gathered:** 2026-05-21
**Status:** Ready for planning

<domain>
## Phase Boundary

搭建完整的开发基础设施，包括Docker Compose编排、PostgreSQL数据库schema（3个独立schema）、Nginx网关、以及三个Go微服务（user-service、forum-service、admin-service）的项目骨架，确保服务间可以通过REST互相通信。此阶段不实现业务逻辑，只搭建可运行的骨架。

</domain>

<decisions>
## Implementation Decisions

### 存储方案
- **D-01:** 附件/图片存储MVP阶段使用本地文件系统（磁盘目录），预留S3/MinIO接口
- **D-02:** 本地存储路径约定为 `/data/uploads/`，按日期分子目录

### JWT鉴权策略
- **D-03:** JWT access token过期时间30分钟
- **D-04:** 使用Redis存储refresh token，支持token刷新
- **D-05:** 用户被封禁后，通过Redis删除其refresh token实现即时失效

### 数据库Migration
- **D-06:** 使用golang-migrate管理PostgreSQL schema版本
- **D-07:** 每个服务的migration文件放在各自服务的 `migrations/` 目录下
- **D-08:** 启动时自动运行migration（up命令）

### 服务间通信
- **D-09:** MVP阶段服务间调用（/internal/v1/前缀）不额外鉴权，信任Docker Compose内部网络隔离
- **D-10:** 后期迁移到gRPC/消息队列时再添加服务间认证

### Claude's Discretion
- 项目目录结构的具体子目录命名由Claude按Go项目最佳实践决定
- Go模块版本选择由Claude按最新稳定版决定

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture
- `AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md` — 微服务架构设计、端口分配、API路由、目录结构、Docker Compose编排方案
- `.planning/ROADMAP.md` — Phase 1目标、成功标准、技术说明
- `.planning/REQUIREMENTS.md` — INFRA-01/02/03、FE-04需求定义
- `.planning/PROJECT.md` — 项目约束、技术栈、关键决策

### No external specs
No external specs beyond the above — decisions fully captured in this CONTEXT.md.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- 无现有代码 — 绿野项目

### Established Patterns
- 无现有模式 — 首次搭建

### Integration Points
- Nginx将作为唯一入口，转发到三个后端服务
- PostgreSQL将创建三个独立schema：schema_auth、schema_forum、schema_admin
- Redis将用于JWT refresh token存储和会话管理

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches for Go + Gin project structure.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-infrastructure*
*Context gathered: 2026-05-21*
