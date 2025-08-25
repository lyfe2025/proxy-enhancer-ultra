# 🎯 配置文件说明 - 完全环境变量化

## ✅ 问题已解决：无硬编码值！

现在所有配置文件都使用环境变量，**完全无硬编码值**！

## 🚀 使用方法：一个 .env 文件搞定一切

**只需要在项目根目录创建 `.env` 文件就可以了！**

## 📁 配置文件说明（现在都是环境变量驱动）

### Go 后端配置
- `config.yaml` - 配置文件（所有值都是环境变量占位符）
- `config.yaml.example` - 配置模板（参考用）

### Vue 前端配置  
- `web/vite.config.ts` - 构建工具配置（使用环境变量）
- `web/package.json` - 依赖配置（已配置好）

## 🚀 快速开始（3步搞定）

### 第1步：创建 .env 文件
在项目根目录创建 `.env` 文件：

```bash
# 环境配置
ENV=development

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s

# TLS配置（可选）
ENABLE_TLS=false
TLS_CERT_FILE=
TLS_KEY_FILE=

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=proxy_enhancer_ultra
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
DB_CONN_MAX_IDLE_TIME=5m
DB_LOG_LEVEL=warn

# JWT配置
JWT_SECRET=your_very_secure_jwt_secret_key_here
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h
JWT_ISSUER=proxy-enhancer-ultra

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
LOG_FILE_PATH=logs/app.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=28
LOG_COMPRESS=true

# 代理配置
PROXY_TIMEOUT=30s
PROXY_MAX_IDLE_CONNS=100
PROXY_IDLE_CONN_TIMEOUT=90s
PROXY_USER_AGENT=ProxyEnhancerUltra/1.0
PROXY_FOLLOW_REDIRECTS=true
PROXY_MAX_REDIRECTS=10
PROXY_BUFFER_SIZE=32768

# 监控配置
MONITORING_ENABLED=true
MONITORING_COLLECT_INTERVAL=1m
MONITORING_RETENTION_DAYS=30
MONITORING_CLEANUP_INTERVAL=24h
MONITORING_METRICS_PATH=/metrics
MONITORING_HEALTH_CHECK_PATH=/health

# CORS配置
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
CORS_EXPOSED_HEADERS=
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400

# 限流配置
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
RATE_LIMIT_WINDOW_SIZE=1m
RATE_LIMIT_CLEANUP_INTERVAL=5m
# 前端配置
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_TITLE=智能反向代理平台
VITE_DEV_PORT=3000

# Redis 配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5

# 数据库时区
DB_TIMEZONE=Asia/Shanghai
```

### 第2步：启动后端
```bash
go run cmd/server/main.go
```

### 第3步：启动前端
```bash
cd web
pnpm dev
```

## 💡 现在的优势

1. **完全灵活** - 所有配置都通过环境变量控制
2. **无硬编码** - 配置文件中的值都是环境变量占位符
3. **环境隔离** - 不同环境可以有不同的 .env 文件
4. **安全** - 敏感信息不会被提交到代码库

## 🔧 支持的环境变量

### 服务器配置
- `SERVER_HOST` - 监听地址
- `SERVER_PORT` - 监听端口
- `SERVER_READ_TIMEOUT` - 读取超时
- `SERVER_WRITE_TIMEOUT` - 写入超时
- `SERVER_IDLE_TIMEOUT` - 空闲超时

### 数据库配置
- `DB_HOST` - 数据库主机
- `DB_PORT` - 数据库端口
- `DB_USER` - 数据库用户名
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名称
- `DB_SSL_MODE` - SSL模式
- `DB_TIMEZONE` - 时区

### JWT配置
- `JWT_SECRET` - JWT密钥
- `JWT_EXPIRES_IN` - 过期时间
- `JWT_REFRESH_EXPIRES_IN` - 刷新过期时间

### 前端配置
- `VITE_API_BASE_URL` - API基础地址
- `VITE_APP_TITLE` - 应用标题
- `VITE_DEV_PORT` - 开发服务器端口

## ⚠️ 注意事项

- `.env` 文件已添加到 .gitignore，不会被提交
- 生产环境请务必修改默认密码和密钥
- 所有配置项都有合理的默认值

## 🎉 总结

**现在你拥有：**
1. **完全灵活的配置系统** - 无硬编码值
2. **统一的配置方式** - 只需要 .env 文件
3. **环境隔离** - 开发/测试/生产环境分离
4. **安全配置** - 敏感信息不会被提交

**你只需要：**
1. 创建 `.env` 文件
2. 设置必要的环境变量
3. 启动项目

**所有配置文件现在都是环境变量驱动，完全无硬编码！**
