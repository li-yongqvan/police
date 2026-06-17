# AI 智联论坛 MVP

一个面向学院级人工智能协会展示的最小可运行 MVP。项目由 `Vue 3 + Vite` 前端和 `Go + Gin` 三个最小服务组成，重点覆盖：

- 演示登录 / 角色切换
- 三大论坛板块
- 发帖、评论、点赞、附件展示
- 最小中台：审核、封禁、配置、统计
- 共享 JSON 种子数据

## 目录结构

```text
frontend/           Vue 3 + Vite 前端
services/user/      用户服务
services/forum/     论坛服务
services/admin/     中台服务
shared/mock-data/   共享种子数据
docs/               研究与展示说明
```

## 本地启动

先分别启动三个 Go 服务：

```powershell
cd services/user
go run ./cmd/server
```

```powershell
cd services/forum
go run ./cmd/server
```

```powershell
cd services/admin
go run ./cmd/server
```

再启动前端：

```powershell
cd frontend
npm install
npm run dev -- --host 127.0.0.1 --port 4173
```

打开 `http://127.0.0.1:4173` 即可进入演示。

## 默认端口

- `user-service`: `8001`
- `forum-service`: `8002`
- `admin-service`: `8003`
- `frontend`: `4173`

## 展示建议动线

1. 用学生账号登录，浏览三大板块和帖子流
2. 打开技术问答详情页，展示评论与资源附件
3. 发布一篇新帖，演示公开发布或进入审核队列
4. 切到管理员或中台管理员，进入中台
5. 展示待审核内容、板块开关、发帖限制和数据概览
