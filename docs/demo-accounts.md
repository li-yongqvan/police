# 演示账号一览

> 密码规则：**批量账号密码与用户名相同**；原有三个账号仍为 `demo123456`。

## 原有账号（保留）

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 学生 | `demo_student` | `demo123456` |
| 协会管理员 | `demo_admin` | `demo123456` |
| 中台管理员 | `demo_platform_admin` | `demo123456` |

## 并行测试账号（各 6 个）

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 学生 | `demo01` … `demo06` | 与用户名相同（如 `demo01` / `demo01`） |
| 协会管理员 | `admin01` … `admin06` | 与用户名相同 |
| 中台管理员 | `plat01` … `plat06` | 与用户名相同 |

## 迁移

- 用户：`services/user/migrations/007_demo_batch_accounts.up.sql`
- 角色：`services/admin/migrations/005_demo_batch_roles.up.sql`

重启 **user-service** 后再启动 **admin-service**（`start-dev.ps1` 已加 8 秒间隔）。若管理员号登录后进不了 `/admin`，执行 `006_demo_batch_roles_fixup` 迁移或重新跑 `005_demo_batch_roles.up.sql`。

## 本机验证（2026-06-04）

- 迁移已应用：用户 21 个（含原 3 个 + 新 18 个）
- `admin01`–`06`、`plat01`–`06` 已绑定 `user_roles`
- 登录 API：`POST http://127.0.0.1:8001/api/v1/login`，body `{"username":"demo01","password":"demo01"}`
