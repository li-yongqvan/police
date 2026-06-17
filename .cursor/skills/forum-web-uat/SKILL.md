---
name: forum-web-uat
description: >-
  Designs and runs human full-feature UAT for forum/Web B/S products using
  distilled methods from James Bach, Cem Kaner, Michael Bolton, Rex Black,
  Jeff Offutt, and James Whittaker. Use when creating exploratory test plans,
  risk-driven walkthroughs, Web input security matrices, release gates, or
  feedback templates for AI-forum-style community platforms.
---

# Forum Web UAT（六大方法论融合）

## When to use

- 编写或更新**人工全功能体验 / UAT** 文档
- 组织内测、灰度、答辩演示前的验收
- 把测试者反馈整理成 AI 可执行的修改单
- 本项目默认产出：`docs/human-full-experience-uat.md`

## Expert distillate (read [masters-reference.md](masters-reference.md) for detail)

| Master | Use in plan as |
|--------|----------------|
| **Bach** | Risk map + session charters; explore beyond scripts |
| **Kaner** | Context table → scope in/out |
| **Bolton** | Tag cases with `[P][J][R][E]` |
| **Black** | Day plan, defect grades, release gate |
| **Offutt** | Web input matrix per form/API field |
| **Whittaker** | Attacker checklist (XSS, authz, upload, rate limit) |

## Workflow: generate human UAT for this repo

### Step 1 — Kaner context (5 min)

Fill or confirm:

| 维度 | 本项目默认 |
|------|------------|
| 阶段 | MVP 校内试运行 |
| 用户规模 | 5～50 |
| 主路径 | 邀请码注册 → 登录 → 浏览 → 发帖 → 评论 → 审核 |
| 不测/弱测 | 万级压测、全浏览器矩阵 |

### Step 2 — Bach risk map

Order testing time: 🔴 权限/安全/数据 → 🟠 主流程 → 🟡 体验 → 🟢 文案样式.

### Step 3 — Bolton environment matrix

At minimum test: **本地 8091** and note if **Docker/生产** differs (demo-login, HTTPS, 网关前缀).

### Step 4 — Black lifecycle

Embed in doc:

1. Day0: health + `scripts/smoke-test.sh`
2. Day1–2: student UAT sessions (charters)
3. Day3: admin UAT
4. Day4: Offutt matrix + Whittaker attacks
5. Day5: defect triage + release gate checklist

### Step 5 — Build the human doc sections

Required sections in output markdown:

1. 测试者须知（账号、环境、预计时长）
2. 情境说明（Kaner）
3. 风险与 Session 章程（Bach）— 3～5 条 charter，每条 60–90 min
4. 四要素走查表（Bolton）— 学生 / 协会 admin / platform_admin
5. Web 输入安全矩阵（Offutt）— 表格式
6. 攻击式清单（Whittaker）— 可勾选
7. **体验反馈单模板**（给 AI / 开发）— 必填字段见下
8. 上线准入（Black）— 勾选项

### Step 6 — Map to this codebase

Routes and roles from `frontend/src/router.js` and README demo accounts:

- Public: `/`, `/register`, `/oauth/qq`
- Community: `/community/*` (home, boards, posts, profile, messages, my/*, circle, rank, about)
- Admin: `/admin/*` (platform_admin only: invites, sensitive, roles)

Scripts: `scripts/smoke-test.sh`, `full-experience-test.sh`, `security-smoke.sh`.

Seed: recommend `scripts/seed/002_pilot_gx_content.sql` before UAT.

## Session charter template (Bach)

```markdown
### Charter <id>: <title>
- **时长**: 60–90 min
- **角色**: demo_student | demo_admin | demo_platform_admin | 新注册
- **使命**: <一句话要发现什么>
- **风险焦点**: [R] ...
- **环境**: [E] 本地 | Docker | 云
- **完成定义**: 走完 <路径列表> 并至少提交 1 条反馈单（无问题则填「通过」）
```

## Feedback template for testers → AI (mandatory in UAT doc)

```markdown
## 反馈 #<n>
- **Charter/场景**:
- **类型**: BUG | UI | 文案 | 安全 | 建议
- **严重度**: P0 | P1 | P2 | P3
- **[P] 功能**:
- **[E] 环境**:
- **复现步骤**: 1. … 2. …
- **预期**:
- **实际**:
- **证据**: 截图 / Network URL+状态码（勿贴 Token）
- **希望怎么改**:
```

## Defect grading (Black)

| 级 | 定义 |
|----|------|
| P0 | 核心不可用或可 exploited 安全 |
| P1 | 主流程受损 |
| P2 | 次要功能 |
| P3 | 体验/样式 |

Release block: 任何 open P0/P1；smoke 全 PASS；UAT 主路径 100%.

## Offutt input matrix (minimal columns)

For each input: **位置** | **合法样例** | **边界** | **非法/空** | **安全 payload** | **预期**

Forum inputs: 注册（邀请码、用户名、密码）、登录、发帖标题/正文、评论、资料昵称/简介、搜索、举报原因、管理端配置、附件文件名.

## Whittaker minimum attacks

- Stored XSS: title, body, comment, bio
- Vertical authz: student → admin API
- Horizontal: tamper `user_id` / `post_id` in URL or API
- Upload: `.exe`, oversized file
- Rate limit: 5+ failed logins
- Info leak: stack trace in API body

## Agent rules when applying this skill

1. For foolproof click-by-click guides, prefer `docs/step-by-step-uat-guide.md`; for methodology UAT, use `docs/human-full-experience-uat.md`.
2. Do not duplicate `test-plan.md` unless extending automation sections.
3. Keep printable checklist ≤ 2 pages optional; full detail in step-by-step guide.
4. Every manual case should have **expected result** and Bolton tags where useful.
5. Do not claim tests passed without runtime evidence if user asked for verification.
6. Link feedback template prominently so testers know how to report to AI.

## Additional resources

- Expert profiles: [masters-reference.md](masters-reference.md)
- Step-by-step human guide: [docs/step-by-step-uat-guide.md](../../docs/step-by-step-uat-guide.md)
- Project UAT output: [docs/human-full-experience-uat.md](../../docs/human-full-experience-uat.md)
- API automation: [docs/test-plan.md](../../docs/test-plan.md) §8
