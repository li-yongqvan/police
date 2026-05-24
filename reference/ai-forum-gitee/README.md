# AI智联论坛（AI Forum）

一个面向高校学院的人工智能技术交流社区论坛，支持用户发帖、评论、资源分享、活动报名，配备完整的中台管控系统（内容审核、权限管理、数据统计）。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22 + Gin 框架 |
| 前端 | Vue 3 + Vite |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 网关 | Nginx |
| 部署 | Docker Compose |

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd ai-forum
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置 JWT_SECRET 等参数
```

### 3. 一键启动

```bash
docker-compose up -d
```

### 4. 验证服务

```bash
# 检查所有服务健康状态
curl http://localhost:8001/health   # user-service
curl http://localhost:8002/health   # forum-service
curl http://localhost:8003/health   # admin-service
curl http://localhost/              # 前端
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 80 | API 网关 + 前端入口 |
| user-service | 8001 | 用户认证、授权、资料管理 |
| forum-service | 8002 | 帖子、评论、板块、点赞、收藏 |
| admin-service | 8003 | 系统配置、内容审核、数据统计 |
| PostgreSQL | 5432 | 共享数据库（三schema隔离） |
| Redis | 6379 | 缓存与会话管理 |
| 前端（开发） | 3000 | Vue 3 开发服务器 |

## 目录结构

```
ai-forum/
├── docker-compose.yml          # Docker 编排配置
├── .env.example                # 环境变量模板
├── nginx/
│   └── nginx.conf              # Nginx 反向代理配置
├── migrations/
│   └── init/                   # 全局数据库初始化脚本
│       ├── 000_init_schemas.up.sql
│       └── 000_init_schemas.down.sql
├── services/
│   ├── user/                   # 用户服务 (:8001)
│   │   ├── cmd/main.go         # 入口
│   │   ├── internal/           # 业务逻辑
│   │   │   ├── handler/        # HTTP 处理器
│   │   │   ├── service/        # 业务服务层
│   │   │   ├── model/          # 数据模型
│   │   │   └── middleware/     # 中间件（JWT、CORS）
│   │   ├── pkg/                # 公共包
│   │   │   ├── database/       # PostgreSQL 连接
│   │   │   ├── redis/          # Redis 连接
│   │   │   └── jwt/            # JWT 工具
│   │   ├── migrations/         # 服务专属迁移
│   │   └── Dockerfile
│   ├── forum/                  # 论坛服务 (:8002)
│   │   ├── internal/client/    # 服务间通信（admin-service）
│   │   └── ...
│   └── admin/                  # 中台服务 (:8003)
│       ├── internal/client/    # 服务间通信（user/forum）
│       └── ...
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── views/              # 页面组件
│   │   ├── router/             # 路由配置
│   │   ├── api/                # API 封装
│   │   ├── App.vue
│   │   └── main.js
│   ├── Dockerfile
│   └── package.json
└── README.md
```

## API 路由

### 用户服务 (`/api/v1/` → user-service)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/register` | 注册（邀请码） |
| POST | `/api/v1/login` | 登录 |
| GET | `/api/v1/users/:id` | 获取个人资料 |
| PUT | `/api/v1/users/:id` | 更新个人资料 |
| POST | `/api/v1/users/:id/avatar` | 上传头像 |

### 论坛服务 (`/api/v1/posts`, `/api/v1/boards` → forum-service)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/boards` | 获取板块列表 |
| GET | `/api/v1/posts` | 获取帖子列表 |
| POST | `/api/v1/posts` | 发帖 |
| GET | `/api/v1/posts/:id` | 帖子详情 |
| PUT | `/api/v1/posts/:id` | 编辑帖子 |
| DELETE | `/api/v1/posts/:id` | 删除帖子 |
| POST | `/api/v1/posts/:id/comments` | 评论 |
| POST | `/api/v1/posts/:id/like` | 点赞 |
| POST | `/api/v1/posts/:id/collect` | 收藏 |
| POST | `/api/v1/attachments/upload` | 附件上传 |

### 中台服务 (`/api/v1/admin/` → admin-service)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/config` | 获取系统配置 |
| PUT | `/api/v1/admin/config` | 更新系统配置 |
| GET | `/api/v1/admin/audit/pending` | 待审核列表 |
| POST | `/api/v1/admin/audit/:id/approve` | 审核通过 |
| POST | `/api/v1/admin/audit/:id/reject` | 审核驳回 |
| POST | `/api/v1/admin/users/:id/ban` | 封禁用户 |
| GET | `/api/v1/admin/stats/overview` | 数据概览 |
| GET | `/api/v1/admin/stats/daily` | 每日统计 |
| POST | `/api/v1/admin/sensitive-words` | 添加敏感词 |

## 数据库

MVP 阶段共享 PostgreSQL 实例，三服务使用独立 schema 隔离：

| 服务 | Schema |
|------|--------|
| user-service | `schema_auth` |
| forum-service | `schema_forum` |
| admin-service | `schema_admin` |

## 开发

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

### 后端开发

```bash
cd services/user
go mod tidy
go run cmd/main.go
```

## License

Internal project - AI Forum Team
