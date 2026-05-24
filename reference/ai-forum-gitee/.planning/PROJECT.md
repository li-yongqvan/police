# AI智联论坛（AI Forum）

## What This Is

一个面向高校学院的人工智能技术交流社区论坛，支持用户发帖、评论、资源分享，配备完整的中台管控系统（内容审核、权限管理、数据统计）。MVP已实现邀请码注册登录、三大板块发帖评论、敏感词审核、中台管理、数据可视化等完整闭环。

## Core Value

让AI学习者能快速找到技术同伴、交流心得、获取资源，同时通过中台实现论坛全流程可控可管可分析。

## Requirements

### Validated

- ✓ 用户通过邀请码注册/登录，注册后为Level 0游客，中台审核升级 — v1.0
- ✓ 用户等级体系：Level 0 游客 → Level 1 初级会员 → Level 2 中级会员 → Level 3 高级会员 → Level 4+，等级影响发帖/评论等操作权限 — v1.0
- ✓ 用户可查看和编辑个人资料（头像、昵称、签名、等级展示） — v1.0
- ✓ 用户可在三大核心板块（AI学习交流区、协会公告&活动区、技术问答求助区）发帖 — v1.0
- ✓ 支持图文发帖、回复评论、点赞、收藏 — v1.0
- ✓ 支持网盘链接/文档图片上传分享 — v1.0
- ✓ 管理员可设置精华帖、置顶帖 — v1.0
- ✓ 简易管理员后台：帖子审核、删帖、用户管理 — v1.0
- ✓ 中台基础控制：发帖限制、板块开关、简单权限分配 — v1.0
- ✓ 中台基础监管：敏感词过滤、帖子人工审核、封禁账号/删帖 — v1.0
- ✓ 中台数据统计：注册用户数、日发帖量、板块活跃度、基础报表 — v1.0
- ✓ 移动端自适应浏览发帖 — v1.0
- ✓ Nginx API网关反向代理三个Go微服务 — v1.0

### Active

（下一里程碑的新需求将在此添加）

### Out of Scope

- 私信聊天、好友关注 — MVP阶段砍掉，第二阶段再做
- 细分复杂技术子板块 — MVP仅三大板块
- 赛事报名系统、线上表单 — 第二阶段
- 风控高级过滤、复杂数据统计（用户画像、流失预警）— 第三阶段
- 校外互通、社会用户专区 — 第四阶段
- 代码高亮、论文排版、数据集专区 — 第二阶段
- 中台高级功能（应急控制、自定义敏感词库、风险预警、数据导出）— 第二阶段
- 企业招聘、直播、线下活动系统 — 远期规划
- 客户端APP — 仅网页端+H5移动端
- AI模型训练、算力部署 — 仅做交流讨论

## Context

- **项目名称**：AI智联论坛（学院级→社会级通用人工智能交流社区）
- **定位**：本校学院内部AI技术交流、学习打卡、赛事组队、资源共享、师生答疑
- **状态**：MVP v1.0已完成
- **架构**：3个Go Gin微服务 + Vue 3前端 + PostgreSQL + Redis + Nginx，Docker Compose编排
- **代码**：约 30000 行，Go 后端 + Vue 前端
- **Git**：尚未初始化正式仓库

## Constraints

- **Tech stack**: Go + Gin（后端）、Vue 3 + Vite（前端）、PostgreSQL + Redis（数据层）、Nginx（网关）
- **Deployment**: Docker Compose（MVP阶段），共享PostgreSQL实例+独立schema隔离
- **Server cost**: 初期轻量化低成本部署
- **Timeline**: MVP快速上线，避免复杂技术开发导致延期
- **Content compliance**: 依托中台监管功能降低合规风险

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| 微服务架构（3服务） | 轻量高性能，适合高并发，预留渐进式扩展 | ✓ Good |
| REST通信（MVP阶段） | 简单直接，预留gRPC/消息队列接口 | ✓ Good |
| 共享PostgreSQL+分schema | 降低MVP部署成本，后期可无缝迁移独立数据库 | ✓ Good |
| 纯邀请码注册 | 控制用户来源质量，中台集中管控邀请码发放 | ✓ Good |
| 用户等级体系(MVP) | Level 0~5+，注册为Level 0，中台手动审核升级，等级影响权限 | ✓ Good |
| 敏感词审核+发帖自动校验 | 发帖时调用admin-service做敏感词校验，命中则进入待审核 | ✓ Good |
| 内网服务间API | /internal/v1/路由仅用于服务间通信，不对外暴露 | ✓ Good |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-05-22 after v1.0 MVP milestone*
