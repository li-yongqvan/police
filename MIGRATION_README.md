# AI智联论坛 迁移包说明

## 包含内容

本迁移包 (`ai-forum-migration.tar.gz`) 包含AI智联论坛项目的完整源码和配置文件，已移除所有Cloudflare相关配置，适合国内云服务器部署。

### 主要文件

- `campus-forum/` - 论坛项目主目录
  - `frontend/` - Vue 3 前端源码（不含node_modules）
  - `services/` - Go后端服务源码
  - `shared/` - 共享数据和种子文件
  - `docs/` - 项目文档

- 配置文件
  - `docker-compose.yml` - Docker编排配置
  - `.env.example` - 环境变量模板
  - `nginx/` - Nginx网关配置
  - `migrations/` - 数据库迁移脚本

- 文档
  - `README.md` - 项目说明
  - `DEPLOY_CN.md` - 国内部署指南
  - `AGENT_DEPLOY.md` - 自动化部署文档
  - `AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md` - 详细需求文档

### 排除内容

以下文件已从迁移包中排除，需要在部署时重新生成：

- `node_modules/` - 前端依赖（97MB，通过npm install安装）
- `admin-service` - 编译后的中台服务二进制文件
- `forum-service` - 编译后的论坛服务二进制文件
- `user-service` - 编译后的用户服务二进制文件
- `.cloudflared/` - Cloudflare配置（已删除）
- `*.log` - 日志文件

## 快速开始

### 1. 解压迁移包

```bash
tar -xzf ai-forum-migration.tar.gz
cd campus-forum
```

### 2. 配置环境

```bash
cp .env.example .env
# 编辑 .env 文件设置数据库密码、JWT密钥等
```

### 3. 安装依赖

```bash
# 安装前端依赖（使用国内镜像）
cd frontend
npm install --registry=https://registry.npmmirror.com
cd ..

# 安装后端依赖
cd services/user && go mod tidy && cd ../..
cd services/forum && go mod tidy && cd ../..
cd services/admin && go mod tidy && cd ../..
```

### 4. 启动服务

```bash
# 使用Docker Compose启动
docker-compose up -d

# 或手动启动开发模式
./start-services.sh
```

## 国内部署优化

详细的国内部署优化指南请参考 `DEPLOY_CN.md`，主要包括：

1. **镜像源配置** - 使用国内npm和Go模块镜像
2. **数据库优化** - PostgreSQL和Redis性能调优
3. **安全配置** - 防火墙、HTTPS、数据备份
4. **监控维护** - 服务状态检查、性能监控

## 技术栈

- **后端**: Go 1.22 + Gin 框架
- **前端**: Vue 3 + Vite
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **网关**: Nginx
- **部署**: Docker Compose

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 80 | API网关 + 前端入口 |
| user-service | 8001 | 用户认证、授权、资料管理 |
| forum-service | 8002 | 帖子、评论、板块、点赞、收藏 |
| admin-service | 8003 | 系统配置、内容审核、数据统计 |
| PostgreSQL | 5432 | 共享数据库（三schema隔离） |
| Redis | 6379 | 缓存与会话管理 |
| 前端（开发） | 4173 | Vue 3 开发服务器 |

## 注意事项

1. **已移除Cloudflare** - 本迁移包不包含任何Cloudflare相关配置
2. **依赖需要重新安装** - node_modules和编译产物需要重新生成
3. **数据库需要初始化** - 首次部署需要运行数据库迁移脚本
4. **环境变量需要配置** - 复制.env.example为.env并填写实际配置

## 项目规模

- 迁移包大小: 241KB
- 文件数量: 279个
- 完整项目大小（含依赖）: 约190MB

## 支持与文档

- 详细需求文档: `AI智联论坛-需求文档+开发边界+MVP范围+微服务架构设计.md`
- 部署指南: `DEPLOY_CN.md`
- 自动化部署: `AGENT_DEPLOY.md`

---

**迁移包生成时间**: 2026-05-30
**适用环境**: 国内云服务器（Ubuntu 20.04+ / CentOS 7+）
**推荐配置**: 4核8G内存，50GB SSD