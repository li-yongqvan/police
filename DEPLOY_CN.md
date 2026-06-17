# AI智联论坛 国内部署指南

## 项目概述

AI智联论坛是一个面向高校学院的人工智能技术交流社区论坛，已移除所有Cloudflare相关配置，适合国内云服务器部署。

## 技术栈

- **后端**: Go 1.22 + Gin 框架
- **前端**: Vue 3 + Vite
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **网关**: Nginx
- **部署**: Docker Compose

## 快速部署步骤

### 1. 上传项目文件

将 `ai-forum-migration.tar.gz` 上传到国内云服务器，然后解压：

```bash
tar -xzf ai-forum-migration.tar.gz
cd campus-forum
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置以下参数：
# - JWT_SECRET: 设置一个安全的密钥
# - POSTGRES_PASSWORD: 设置数据库密码
# - 其他数据库连接参数
```

### 3. 安装依赖

#### 前端依赖
```bash
cd frontend
npm install
# 或使用国内镜像
npm install --registry=https://registry.npmmirror.com
cd ..
```

#### 后端依赖（如果需要重新编译）
```bash
cd services/user
go mod tidy
cd ../forum
go mod tidy
cd ../admin
go mod tidy
cd ../..
```

### 4. 启动服务

#### 方式一：使用Docker Compose（推荐）
```bash
# 确保已安装 Docker 和 Docker Compose
docker-compose up -d
```

#### 方式二：手动启动（开发模式）
```bash
# 启动后端服务
cd services/user && go run ./cmd/server &
cd ../forum && go run ./cmd/server &
cd ../admin && go run ./cmd/server &

# 启动前端
cd frontend
npm run dev -- --host 0.0.0.0 --port 4173
```

### 5. 访问服务

- **前端**: http://your-server-ip:4173 (开发模式) 或 http://your-server-ip (Docker模式)
- **用户服务**: http://your-server-ip:8001
- **论坛服务**: http://your-server-ip:8002
- **中台服务**: http://your-server-ip:8003

## 国内部署优化建议

### 1. 使用国内镜像源

#### npm镜像
```bash
# 设置npm使用淘宝镜像
npm config set registry https://registry.npmmirror.com

# 或在项目中创建 .npmrc 文件
echo "registry=https://registry.npmmirror.com" > frontend/.npmrc
```

#### Go模块代理
```bash
# 设置Go模块代理
export GOPROXY=https://goproxy.cn,direct
```

### 2. 数据库优化

#### PostgreSQL配置优化
编辑 `docker-compose.yml` 中的 postgres 服务，添加性能优化参数：

```yaml
postgres:
  image: postgres:16
  command: >
    postgres
    -c shared_buffers=256MB
    -c effective_cache_size=768MB
    -c work_mem=4MB
    -c maintenance_work_mem=64MB
    -c max_connections=200
  # ... 其他配置
```

### 3. Redis配置优化

```yaml
redis:
  image: redis:7-alpine
  command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
  # ... 其他配置
```

### 4. Nginx优化

编辑 `nginx/nginx.conf`，添加以下优化配置：

```nginx
# 启用gzip压缩
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

# 静态文件缓存
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 生产环境配置

### 1. 安全配置

#### 修改默认密码
```bash
# 编辑 .env 文件
JWT_SECRET=your-super-secret-key-here
POSTGRES_PASSWORD=strong-database-password
```

#### 防火墙配置
```bash
# 只开放必要端口
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 22/tcp    # SSH
# 关闭其他端口
sudo ufw enable
```

### 2. HTTPS配置（推荐）

使用Let's Encrypt免费SSL证书：

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加：0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. 数据备份

#### PostgreSQL备份脚本
```bash
#!/bin/bash
# backup-db.sh
DATE=$(date +%Y%m%d_%H%M%S)
docker exec postgres pg_dump -U ai_forum ai_forum > backup_$DATE.sql
gzip backup_$DATE.sql

# 保留最近30天备份
find . -name "backup_*.sql.gz" -mtime +30 -delete
```

#### 定时备份
```bash
# 添加到crontab
0 2 * * * /path/to/backup-db.sh
```

## 监控和维护

### 1. 服务状态检查
```bash
# 检查所有服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f user-service
docker-compose logs -f forum-service
docker-compose logs -f admin-service
```

### 2. 性能监控
```bash
# 监控资源使用
docker stats

# 监控数据库
docker exec postgres psql -U ai_forum -c "SELECT * FROM pg_stat_activity;"
```

## 常见问题解决

### 1. 端口冲突
```bash
# 检查端口占用
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :5432

# 修改docker-compose.yml中的端口映射
```

### 2. 数据库连接失败
```bash
# 检查PostgreSQL状态
docker-compose logs postgres

# 重置数据库
docker-compose down -v
docker-compose up -d
```

### 3. 前端构建失败
```bash
# 清除缓存重新安装
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

## 项目结构说明

```
campus-forum/
├── frontend/           # Vue 3 前端源码
├── services/           # Go后端服务
│   ├── user/          # 用户服务
│   ├── forum/         # 论坛服务
│   └── admin/         # 中台服务
├── shared/            # 共享数据
├── docs/              # 文档
└── start-services.sh  # 启动脚本
```

## 联系支持

如有部署问题，请检查：
1. 服务器系统要求：Ubuntu 20.04+ / CentOS 7+
2. 最低配置：2核4G内存
3. 推荐配置：4核8G内存
4. 确保Docker和Docker Compose已正确安装

---

**注意**: 本项目已移除所有Cloudflare相关配置，可直接在国内云服务器部署使用。