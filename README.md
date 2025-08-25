# Proxy Enhancer Ultra

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-green.svg)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> 🚀 智能反向代理平台 - 无侵入式网站功能增强解决方案

## 📖 项目简介

Proxy Enhancer Ultra 是一个基于 Go 语言开发的智能反向代理平台，能够完全代理目标网站到自有域名，并在代理过程中注入自定义的交互功能。系统主要解决企业需要在不修改原网站代码的情况下，增加自定义交互功能（如弹窗、表单收集）的需求，适用于营销活动、用户调研、数据收集等场景。

### ✨ 核心特性

- 🔄 **完全代理**：无缝代理目标网站，保持原有样式和功能
- 🎯 **智能注入**：根据配置规则动态注入弹窗和表单
- 📊 **数据收集**：实时收集用户提交数据，支持导出和分析
- 🎨 **可视化设计**：拖拽式弹窗和表单设计器
- 🔐 **权限管理**：多角色权限控制，支持系统管理员和运营人员
- 📱 **响应式设计**：支持桌面端和移动端访问
- 🚀 **高性能**：基于 Go 的高并发处理能力
- 🗄️ **灵活存储**：支持 Supabase 和本地 PostgreSQL

## 🏗️ 技术架构

### 后端技术栈

- **核心语言**：Go 1.21+
- **Web框架**：标准库 net/http + gorilla/mux
- **反向代理**：net/http/httputil.ReverseProxy
- **数据库**：Supabase PostgreSQL + GORM
- **认证**：JWT (golang-jwt/jwt)
- **配置管理**：Viper
- **日志**：Logrus

### 前端技术栈

- **框架**：Vue 3 + TypeScript
- **UI组件库**：Element Plus
- **构建工具**：Vite
- **状态管理**：Pinia
- **HTTP客户端**：Axios
- **路由**：Vue Router 4

### 系统架构

```
┌─────────────────┐    ┌─────────────────────────────────┐
│   用户浏览器    │    │         Go统一服务               │
│                │───▶│                                 │
│ 访问代理域名   │    │  ┌─────────────┐ ┌─────────────┐ │
└─────────────────┘    │  │  代理服务   │ │  管理后台   │ │
                       │  │             │ │             │ │
                       │  │ 反向代理    │ │  Web API    │ │
                       │  │ 内容注入    │ │  JWT认证    │ │
                       │  └─────────────┘ └─────────────┘ │
                       └─────────────────────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Supabase       │
                       │  PostgreSQL     │
                       │  统一数据存储   │
                       └─────────────────┘
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 12+ (或 Supabase 账户)
- Git

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/your-username/proxy-enhancer-ultra.git
cd proxy-enhancer-ultra
```

2. **后端设置**
```bash
# 安装 Go 依赖
go mod download

# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置文件，设置数据库连接信息
vim config.yaml
```

3. **前端设置**
```bash
cd web
npm install

# 构建前端
npm run build
```

4. **启动服务**
```bash
# 启动后端服务
go run cmd/server/main.go

# 或使用二进制文件
./bin/server
```

### 环境变量配置

创建 `.env` 文件：

```bash
# 数据库配置
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=proxy_platform
DB_SSL_MODE=disable

# JWT配置
JWT_SECRET=your-secret-key

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

## 📋 功能模块

### 1. 代理管理
- 代理配置管理
- 域名绑定设置
- 目标网站配置
- 状态监控和日志

### 2. 规则配置
- 触发条件设置
- 弹窗设计器
- 表单构建器
- 展示规则配置

### 3. 数据收集
- 实时数据查看
- 数据筛选和搜索
- 统计分析图表
- 数据导出功能

### 4. 系统管理
- 用户权限管理
- 系统配置设置
- 性能监控
- 安全审计

## 🎨 用户界面

系统采用现代化的暗色主题设计：

- **主色调**：深黑色(#0a0a0a)背景，亮绿色(#00ff88)强调
- **设计风格**：现代化卡片式布局，支持微妙阴影和边框发光效果
- **响应式**：桌面优先设计，支持移动端自适应
- **交互体验**：流畅动画反馈，绿色发光悬停效果

## 📁 项目结构

```
proxy-enhancer-ultra/
├── cmd/                    # 主程序入口
├── internal/               # 内部包
│   ├── auth/              # 认证模块
│   ├── config/            # 配置管理
│   ├── database/          # 数据库层
│   ├── handlers/          # HTTP处理器
│   ├── middleware/        # 中间件
│   ├── models/            # 数据模型
│   ├── proxy/             # 代理服务
│   └── services/          # 业务服务
├── pkg/                   # 公共包
├── web/                   # 前端应用
├── configs/               # 配置文件
└── docs/                  # 文档
```

## 🔧 开发指南

### 代码规范

- 遵循 Go 官方代码规范
- 文件大小控制在 300 行以内
- 遵循单一职责原则
- 完整的错误处理和日志记录

### 测试

```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test -tags=integration ./...

# 运行端到端测试
go test -tags=e2e ./...
```

### 构建

```bash
# 构建后端
go build -o bin/server cmd/server/main.go

# 构建前端
cd web && npm run build

# 构建 Docker 镜像
docker build -t proxy-enhancer-ultra .
```

## 🚀 部署

### Docker 部署

```bash
# 构建镜像
docker build -t proxy-enhancer-ultra .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  proxy-enhancer-ultra
```

### 生产环境部署

1. **配置反向代理**（Nginx/Apache）
2. **设置 SSL 证书**
3. **配置数据库连接池**
4. **设置日志轮转**
5. **配置监控和告警**

## 📊 性能特性

- **高并发**：基于 Go 的 goroutine 并发模型
- **连接池**：数据库连接池优化
- **缓存策略**：规则缓存和响应缓存
- **负载均衡**：支持多实例部署

## 🔒 安全特性

- **JWT 认证**：安全的用户认证机制
- **CORS 配置**：跨域请求安全控制
- **速率限制**：防止恶意请求攻击
- **输入验证**：严格的数据验证和清理
- **SQL 注入防护**：使用参数化查询

## 🤝 贡献指南

我们欢迎所有形式的贡献！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 贡献规范

- 遵循现有的代码风格
- 添加适当的测试
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系我们

- **项目主页**：https://github.com/your-username/proxy-enhancer-ultra
- **问题反馈**：https://github.com/your-username/proxy-enhancer-ultra/issues
- **功能建议**：https://github.com/your-username/proxy-enhancer-ultra/discussions

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者和用户！

---

**⭐ 如果这个项目对您有帮助，请给我们一个星标！**
