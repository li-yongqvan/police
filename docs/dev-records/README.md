# 开发记录制度

本目录用于沉淀项目的重要开发节点、技术决策和上线记录。它不是替代 Git 历史，而是给人看的上下文索引：为什么改、改了什么、如何验证、后续风险是什么。

## 目录结构

- `timeline.md`：开发时间线索引。
- `entries/`：每次重要提交或阶段性工作的记录。
- `adr/`：Architecture Decision Record，记录长期有效的技术/流程决策。
- `templates/`：手写记录模板。

## 自动记录

本地 `.git/hooks/post-commit` 可调用：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/dev-record.ps1 -PostCommit -Category chore
```

脚本会生成一条 `docs/dev-records/entries/*.md`，并追加到 `timeline.md`。记录生成后仍需要人工补充“Context / Outcome / Follow-ups”，避免只留下机器流水账。

## 建议流程

1. 功能或修复提交后，由 hook 自动生成记录。
2. 当天收尾时补充用户目标、验证方式和遗留风险。
3. 涉及架构、部署、账号体系、权限模型等长期约定时，在 `adr/` 新增 ADR。
4. 合并或发布前检查 `timeline.md` 链接是否可用。
