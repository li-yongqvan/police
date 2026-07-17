# ADR-0001: 建立开发记录制度

- Status: Accepted
- Date: 2026-07-03
- Scope: project workflow

## Context

AI 智联论坛已经进入试运行和持续优化阶段，改动横跨前端、后端、数据库迁移、部署脚本和云服务器状态。仅依赖 Git commit message 难以解释每次改动背后的用户目标、验证过程和发布风险。

## Decision

建立 `docs/dev-records/` 作为项目级开发记录目录：

- 用 `timeline.md` 做时间线索引。
- 用 `entries/` 保存阶段性或提交级记录。
- 用 `adr/` 保存长期有效的技术和流程决策。
- 用 `scripts/dev-record.ps1` 支持本地 post-commit 自动生成草稿记录。

## Consequences

开发记录需要人工维护，不应只依赖自动生成内容。自动记录只负责捕捉提交、统计和工作区状态；真正有价值的“为什么”和“怎么验收”需要在提交后补充。

该制度会增加少量文档成本，但能降低后续排查部署漂移、接口契约变化和多 agent 协作覆盖的成本。
