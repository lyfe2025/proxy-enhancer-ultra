# 数据库注释文档

本文档描述了项目中所有数据库表和字段的中文注释说明。

## 概述

为了提高代码的可读性和维护性，我们为所有数据库表和字段添加了详细的中文注释。这些注释不仅在Go代码中可见，同时也被应用到了PostgreSQL数据库中。

## 表结构说明

### 1. 用户管理相关表

#### users (用户表)
存储系统用户信息

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| username | VARCHAR(50) | 用户名，唯一标识符 |
| email | VARCHAR(255) | 邮箱地址，用于登录和通知 |
| password_hash | VARCHAR(255) | 密码哈希值，不返回给前端 |
| last_login_at | TIMESTAMP | 最后登录时间，可为空 |
| is_active | BOOLEAN | 账户状态，默认激活 |

#### roles (角色表)
存储系统角色信息

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | INTEGER | 角色ID，自增主键 |
| name | VARCHAR(50) | 角色名称，如 admin、user |
| description | TEXT | 角色功能描述 |
| created_at | TIMESTAMP | 创建时间 |

#### permissions (权限表)
存储系统权限信息

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | INTEGER | 权限ID，自增主键 |
| name | VARCHAR(100) | 权限名称，如 user:create |
| resource | VARCHAR(50) | 权限作用的资源，如 user、proxy |
| action | VARCHAR(50) | 权限操作类型，如 create、read、update、delete |
| created_at | TIMESTAMP | 创建时间 |

### 2. 关联表

#### user_roles (用户角色关联表)
多对多关系表

| 字段名 | 类型 | 注释 |
|--------|------|------|
| user_id | UUID | 用户ID，外键 |
| role_id | INTEGER | 角色ID，外键 |
| created_at | TIMESTAMP | 关联建立时间 |

#### role_permissions (角色权限关联表)
多对多关系表

| 字段名 | 类型 | 注释 |
|--------|------|------|
| role_id | INTEGER | 角色ID，外键 |
| permission_id | INTEGER | 权限ID，外键 |
| created_at | TIMESTAMP | 关联建立时间 |

### 3. 代理配置相关表

#### proxy_configs (代理配置表)
存储代理服务器配置信息

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| name | VARCHAR(100) | 配置名称，便于识别 |
| target_url | VARCHAR(500) | 代理目标地址 |
| proxy_domain | VARCHAR(255) | 代理服务域名 |
| is_active | BOOLEAN | 配置启用状态 |
| settings | JSONB | 配置参数，JSON格式存储 |

#### domains (域名配置表)
存储代理域名配置

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| proxy_config_id | UUID | 所属代理配置ID |
| domain | VARCHAR(255) | 域名地址 |
| is_primary | BOOLEAN | 是否为主要域名 |
| ssl_enabled | BOOLEAN | SSL/TLS 配置状态 |

### 4. 业务规则相关表

#### rules (业务规则表)
存储业务规则配置

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| name | VARCHAR(100) | 规则名称 |
| description | TEXT | 规则功能描述 |
| trigger_conditions | JSONB | 规则触发条件，JSON格式 |
| actions | JSONB | 规则执行动作，JSON格式 |
| is_active | BOOLEAN | 规则启用状态 |
| priority | INTEGER | 规则执行优先级，数字越大优先级越高 |

### 5. 弹窗和数据收集相关表

#### popups (弹窗配置表)
存储弹窗配置信息

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| title | VARCHAR(200) | 弹窗显示标题 |
| content | TEXT | 弹窗显示内容 |
| style_config | JSONB | 弹窗样式配置，JSON格式 |
| form_config | JSONB | 表单字段配置，JSON格式 |
| is_active | BOOLEAN | 弹窗启用状态 |

#### submissions (数据提交表)
存储用户通过弹窗提交的数据

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| popup_id | UUID | 所属弹窗ID |
| form_data | JSONB | 用户提交的表单数据，JSON格式 |
| user_ip | INET | 用户IP地址 |
| user_agent | TEXT | 用户浏览器信息 |
| referrer_url | TEXT | 提交时的来源页面 |

### 6. 监控和日志相关表

#### proxy_logs (代理访问日志表)
存储代理服务访问日志

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| proxy_config_id | UUID | 所属代理配置ID |
| method | VARCHAR(10) | HTTP请求方法，如GET、POST |
| url | TEXT | 请求的完整URL |
| user_agent | TEXT | 用户浏览器信息 |
| user_ip | INET | 用户IP地址 |
| status_code | INTEGER | HTTP响应状态码 |
| response_time | BIGINT | 响应耗时，单位毫秒 |
| error_message | TEXT | 错误信息，失败时记录 |

#### system_metrics (系统监控指标表)
存储系统监控指标数据

| 字段名 | 类型 | 注释 |
|--------|------|------|
| id | UUID | 主键ID，UUID类型 |
| created_at | TIMESTAMP | 创建时间，自动设置 |
| updated_at | TIMESTAMP | 更新时间，自动维护 |
| deleted_at | TIMESTAMP | 软删除时间，支持逻辑删除 |
| metric_name | VARCHAR(100) | 指标名称，如CPU使用率、内存使用量 |
| metric_value | DOUBLE PRECISION | 指标的具体数值 |
| tags | JSONB | 指标标签，用于分类和过滤，JSON格式 |
| timestamp | TIMESTAMP | 指标记录的时间戳 |

## 使用方法

### 更新注释

如果需要更新数据库注释，可以使用以下方法：

1. **修改Go模型文件**：
   ```go
   // 在 internal/models/models.go 中修改字段注释
   Username string `json:"username" gorm:"uniqueIndex;size:50;not null" comment:"用户名，唯一标识符"`
   ```

2. **更新SQL脚本**：
   修改 `internal/database/add_column_comments.sql` 文件中对应的注释。

3. **执行更新脚本**：
   ```bash
   ./scripts/update-db-comments.sh
   ```

### 查看数据库注释

在PostgreSQL中查看表和列注释：

```sql
-- 查看表注释
SELECT schemaname, tablename, obj_description(oid, 'pg_class') as table_comment 
FROM pg_tables JOIN pg_class ON pg_tables.tablename = pg_class.relname 
WHERE schemaname = 'public' ORDER BY tablename;

-- 查看列注释
SELECT 
    column_name,
    data_type,
    col_description((SELECT oid FROM pg_class WHERE relname = 'users'), ordinal_position) as comment
FROM information_schema.columns 
WHERE table_name = 'users' 
ORDER BY ordinal_position;
```

## 注释规范

### 命名规范
- 注释使用简洁明了的中文
- 包含字段的用途和重要约束
- 对于外键字段，明确说明关联关系
- 对于JSON字段，说明存储的数据格式

### 示例
```go
// 好的注释示例
Username string `json:"username" comment:"用户名，唯一标识符"`
Email    string `json:"email" comment:"邮箱地址，用于登录和通知"`
Settings string `json:"settings" comment:"配置参数，JSON格式存储"`

// 避免的注释示例
Username string `json:"username" comment:"用户名"`           // 太简单
Email    string `json:"email" comment:"email address"`     // 使用英文
```

## 维护建议

1. **保持同步**：Go代码中的注释和数据库中的注释应该保持一致
2. **及时更新**：当字段含义或用途发生变化时，及时更新注释
3. **代码审查**：在代码审查时检查注释的准确性和完整性
4. **文档更新**：重要的表结构变化应该同时更新此文档

通过完善的注释系统，新加入项目的开发者能够快速理解数据库结构，提高开发效率和代码质量。
