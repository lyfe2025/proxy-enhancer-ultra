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
# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# 数据库配置
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=proxy_enhancer
DB_SSL_MODE=disable

# JWT配置
JWT_SECRET=your-super-secret-key-here

# 前端配置
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_TITLE=智能反向代理平台
VITE_DEV_PORT=3000
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
