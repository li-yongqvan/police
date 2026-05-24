# Phase 1: 基础设施与项目骨架 - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-05-21
**Phase:** 01-基础设施与项目骨架
**Areas discussed:** 存储方案, JWT策略, 数据库Migration, 服务鉴权

---

## 存储方案

| Option | Description | Selected |
|--------|-------------|----------|
| 本地磁盘 | MVP直接用本地磁盘目录，预留S3接口 | ✓ |
| MinIO/S3 | MVP就接入MinIO，一步到位支持后续扩容 | |
| 数据库存储 | 只存base64到数据库字段 | |

**User's choice:** 本地磁盘(推荐)
**Notes:** 附件路径按日期分目录

## JWT策略

| Option | Description | Selected |
|--------|-------------|----------|
| 短过期+刷新Token | JWT过期30分钟+Redis存储refresh token，被ban后即时失效 | ✓ |
| 长过期单次Token | JWT过期24小时，无刷新机制 | |
| 无状态JWT | 无状态JWT，被ban靠黑名单表 | |

**User's choice:** 短过期+刷新Token(推荐)

## 数据库Migration

| Option | Description | Selected |
|--------|-------------|----------|
| golang-migrate | 使用golang-migrate管理schema版本 | ✓ |
| 手动SQL脚本 | 手动执行SQL脚本 | |

**User's choice:** golang-migrate(推荐)

## 服务鉴权

| Option | Description | Selected |
|--------|-------------|----------|
| 内网信任 | 服务间通过内网REST直接调用，不额外鉴权 | ✓ |
| 内部密钥 | 服务间调用携带内部API密钥 | |

**User's choice:** 内网信任(推荐)

---

## Deferred Ideas

None.
