-- 为数据库所有表字段添加中文注释
-- 执行此脚本来为现有数据库表添加列注释

-- 用户表 (users) 注释
COMMENT ON COLUMN users.id IS '主键ID，UUID类型';
COMMENT ON COLUMN users.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN users.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN users.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN users.username IS '用户名，唯一标识符';
COMMENT ON COLUMN users.email IS '邮箱地址，用于登录和通知';
COMMENT ON COLUMN users.password_hash IS '密码哈希值，不返回给前端';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间，可为空';
COMMENT ON COLUMN users.is_active IS '账户状态，默认激活';

-- 角色表 (roles) 注释
COMMENT ON COLUMN roles.id IS '角色ID，自增主键';
COMMENT ON COLUMN roles.name IS '角色名称，如 admin、user';
COMMENT ON COLUMN roles.description IS '角色功能描述';
COMMENT ON COLUMN roles.created_at IS '创建时间';

-- 权限表 (permissions) 注释
COMMENT ON COLUMN permissions.id IS '权限ID，自增主键';
COMMENT ON COLUMN permissions.name IS '权限名称，如 user:create';
COMMENT ON COLUMN permissions.resource IS '权限作用的资源，如 user、proxy';
COMMENT ON COLUMN permissions.action IS '权限操作类型，如 create、read、update、delete';
COMMENT ON COLUMN permissions.created_at IS '创建时间';

-- 用户角色关联表 (user_roles) 注释
COMMENT ON COLUMN user_roles.user_id IS '用户ID，外键';
COMMENT ON COLUMN user_roles.role_id IS '角色ID，外键';
COMMENT ON COLUMN user_roles.created_at IS '关联建立时间';

-- 角色权限关联表 (role_permissions) 注释
COMMENT ON COLUMN role_permissions.role_id IS '角色ID，外键';
COMMENT ON COLUMN role_permissions.permission_id IS '权限ID，外键';
COMMENT ON COLUMN role_permissions.created_at IS '关联建立时间';

-- 代理配置表 (proxy_configs) 注释
COMMENT ON COLUMN proxy_configs.id IS '主键ID，UUID类型';
COMMENT ON COLUMN proxy_configs.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN proxy_configs.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN proxy_configs.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN proxy_configs.name IS '配置名称，便于识别';
COMMENT ON COLUMN proxy_configs.target_url IS '代理目标地址';
COMMENT ON COLUMN proxy_configs.proxy_domain IS '代理服务域名';
COMMENT ON COLUMN proxy_configs.is_active IS '配置启用状态';
COMMENT ON COLUMN proxy_configs.settings IS '配置参数，JSON格式存储';

-- 域名配置表 (domains) 注释
COMMENT ON COLUMN domains.id IS '主键ID，UUID类型';
COMMENT ON COLUMN domains.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN domains.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN domains.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN domains.proxy_config_id IS '所属代理配置ID';
COMMENT ON COLUMN domains.domain IS '域名地址';
COMMENT ON COLUMN domains.is_primary IS '是否为主要域名';
COMMENT ON COLUMN domains.ssl_enabled IS 'SSL/TLS 配置状态';

-- 业务规则表 (rules) 注释
COMMENT ON COLUMN rules.id IS '主键ID，UUID类型';
COMMENT ON COLUMN rules.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN rules.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN rules.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN rules.name IS '规则名称';
COMMENT ON COLUMN rules.description IS '规则功能描述';
COMMENT ON COLUMN rules.trigger_conditions IS '规则触发条件，JSON格式';
COMMENT ON COLUMN rules.actions IS '规则执行动作，JSON格式';
COMMENT ON COLUMN rules.is_active IS '规则启用状态';
COMMENT ON COLUMN rules.priority IS '规则执行优先级，数字越大优先级越高';

-- 弹窗配置表 (popups) 注释
COMMENT ON COLUMN popups.id IS '主键ID，UUID类型';
COMMENT ON COLUMN popups.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN popups.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN popups.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN popups.title IS '弹窗显示标题';
COMMENT ON COLUMN popups.content IS '弹窗显示内容';
COMMENT ON COLUMN popups.style_config IS '弹窗样式配置，JSON格式';
COMMENT ON COLUMN popups.form_config IS '表单字段配置，JSON格式';
COMMENT ON COLUMN popups.is_active IS '弹窗启用状态';

-- 数据提交表 (submissions) 注释
COMMENT ON COLUMN submissions.id IS '主键ID，UUID类型';
COMMENT ON COLUMN submissions.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN submissions.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN submissions.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN submissions.popup_id IS '所属弹窗ID';
COMMENT ON COLUMN submissions.form_data IS '用户提交的表单数据，JSON格式';
COMMENT ON COLUMN submissions.user_ip IS '用户IP地址';
COMMENT ON COLUMN submissions.user_agent IS '用户浏览器信息';
COMMENT ON COLUMN submissions.referrer_url IS '提交时的来源页面';

-- 代理访问日志表 (proxy_logs) 注释
COMMENT ON COLUMN proxy_logs.id IS '主键ID，UUID类型';
COMMENT ON COLUMN proxy_logs.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN proxy_logs.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN proxy_logs.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN proxy_logs.proxy_config_id IS '所属代理配置ID';
COMMENT ON COLUMN proxy_logs.method IS 'HTTP请求方法，如GET、POST';
COMMENT ON COLUMN proxy_logs.url IS '请求的完整URL';
COMMENT ON COLUMN proxy_logs.user_agent IS '用户浏览器信息';
COMMENT ON COLUMN proxy_logs.user_ip IS '用户IP地址';
COMMENT ON COLUMN proxy_logs.status_code IS 'HTTP响应状态码';
COMMENT ON COLUMN proxy_logs.response_time IS '响应耗时，单位毫秒';
COMMENT ON COLUMN proxy_logs.error_message IS '错误信息，失败时记录';

-- 系统监控指标表 (system_metrics) 注释
COMMENT ON COLUMN system_metrics.id IS '主键ID，UUID类型';
COMMENT ON COLUMN system_metrics.created_at IS '创建时间，自动设置';
COMMENT ON COLUMN system_metrics.updated_at IS '更新时间，自动维护';
COMMENT ON COLUMN system_metrics.deleted_at IS '软删除时间，支持逻辑删除';
COMMENT ON COLUMN system_metrics.metric_name IS '指标名称，如CPU使用率、内存使用量';
COMMENT ON COLUMN system_metrics.metric_value IS '指标的具体数值';
COMMENT ON COLUMN system_metrics.tags IS '指标标签，用于分类和过滤，JSON格式';
COMMENT ON COLUMN system_metrics.timestamp IS '指标记录的时间戳';

-- 表注释
COMMENT ON TABLE users IS '用户表 - 存储系统用户信息';
COMMENT ON TABLE roles IS '角色表 - 存储系统角色信息';
COMMENT ON TABLE permissions IS '权限表 - 存储系统权限信息';
COMMENT ON TABLE user_roles IS '用户角色关联表 - 多对多关系表';
COMMENT ON TABLE role_permissions IS '角色权限关联表 - 多对多关系表';
COMMENT ON TABLE proxy_configs IS '代理配置表 - 存储代理服务器配置信息';
COMMENT ON TABLE domains IS '域名配置表 - 存储代理域名配置';
COMMENT ON TABLE rules IS '业务规则表 - 存储业务规则配置';
COMMENT ON TABLE popups IS '弹窗配置表 - 存储弹窗配置信息';
COMMENT ON TABLE submissions IS '数据提交表 - 存储用户通过弹窗提交的数据';
COMMENT ON TABLE proxy_logs IS '代理访问日志表 - 存储代理服务访问日志';
COMMENT ON TABLE system_metrics IS '系统监控指标表 - 存储系统监控指标数据';
