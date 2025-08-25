package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel 基础模型 - 包含所有表的公共字段
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()" comment:"主键ID，UUID类型"` // 主键ID，自动生成UUID
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`                                 // 创建时间，自动设置
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime" comment:"更新时间"`                                 // 更新时间，自动维护
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" comment:"软删除时间"`                                                  // 软删除时间，支持逻辑删除
}

// User 用户模型 - 存储系统用户信息
type User struct {
	BaseModel
	Username     string     `json:"username" gorm:"uniqueIndex;size:50;not null" comment:"用户名，唯一"` // 用户名，唯一标识符
	Email        string     `json:"email" gorm:"uniqueIndex;size:255;not null" comment:"邮箱地址，唯一"`  // 邮箱地址，用于登录和通知
	PasswordHash string     `json:"-" gorm:"size:255;not null" comment:"密码哈希值"`                    // 密码哈希值，不返回给前端
	LastLoginAt  *time.Time `json:"last_login_at" comment:"最后登录时间"`                                // 最后登录时间，可为空
	IsActive     bool       `json:"is_active" gorm:"default:true" comment:"是否激活"`                  // 账户状态，默认激活
	UserRoles    []UserRole `json:"user_roles" gorm:"foreignKey:UserID" comment:"用户角色关联"`          // 用户角色关联列表
}

// Role 角色模型 - 存储系统角色信息
type Role struct {
	ID              int              `json:"id" gorm:"primaryKey;autoIncrement" comment:"角色ID，自增主键"`     // 角色ID，自增主键
	Name            string           `json:"name" gorm:"uniqueIndex;size:50;not null" comment:"角色名称，唯一"` // 角色名称，如 admin、user
	Description     string           `json:"description" gorm:"type:text" comment:"角色描述"`                // 角色功能描述
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`            // 创建时间
	UserRoles       []UserRole       `json:"user_roles" gorm:"foreignKey:RoleID" comment:"用户角色关联"`       // 用户角色关联列表
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:RoleID" comment:"角色权限关联"` // 角色权限关联列表
}

// Permission 权限模型 - 存储系统权限信息
type Permission struct {
	ID              int              `json:"id" gorm:"primaryKey;autoIncrement" comment:"权限ID，自增主键"`           // 权限ID，自增主键
	Name            string           `json:"name" gorm:"uniqueIndex;size:100;not null" comment:"权限名称，唯一"`      // 权限名称，如 user:create
	Resource        string           `json:"resource" gorm:"size:50;not null" comment:"资源名称"`                  // 权限作用的资源，如 user、proxy
	Action          string           `json:"action" gorm:"size:50;not null" comment:"操作类型"`                    // 权限操作类型，如 create、read、update、delete
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime" comment:"创建时间"`                  // 创建时间
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:PermissionID" comment:"角色权限关联"` // 角色权限关联列表
}

// UserRole 用户角色关联模型 - 多对多关系表
type UserRole struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey" comment:"用户ID"` // 用户ID，外键
	RoleID    int       `json:"role_id" gorm:"primaryKey" comment:"角色ID"`           // 角色ID，外键
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" comment:"关联创建时间"`  // 关联建立时间
	User      User      `json:"user" gorm:"foreignKey:UserID" comment:"关联的用户"`      // 关联的用户对象
	Role      Role      `json:"role" gorm:"foreignKey:RoleID" comment:"关联的角色"`      // 关联的角色对象
}

// RolePermission 角色权限关联模型 - 多对多关系表
type RolePermission struct {
	RoleID       int        `json:"role_id" gorm:"primaryKey" comment:"角色ID"`                  // 角色ID，外键
	PermissionID int        `json:"permission_id" gorm:"primaryKey" comment:"权限ID"`            // 权限ID，外键
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime" comment:"关联创建时间"`         // 关联建立时间
	Role         Role       `json:"role" gorm:"foreignKey:RoleID" comment:"关联的角色"`             // 关联的角色对象
	Permission   Permission `json:"permission" gorm:"foreignKey:PermissionID" comment:"关联的权限"` // 关联的权限对象
}

// ProxyConfig 代理配置模型 - 存储代理服务器配置信息
type ProxyConfig struct {
	BaseModel
	Name        string   `json:"name" gorm:"size:100;not null" comment:"配置名称"`                  // 配置名称，便于识别
	TargetURL   string   `json:"target_url" gorm:"size:500;not null" comment:"目标URL"`           // 代理目标地址
	ProxyDomain string   `json:"proxy_domain" gorm:"size:255;not null" comment:"代理域名"`          // 代理服务域名
	IsActive    bool     `json:"is_active" gorm:"default:false" comment:"是否启用"`                 // 配置启用状态
	Settings    string   `json:"settings" gorm:"type:jsonb;default:'{}'" comment:"配置设置，JSON格式"` // 配置参数，JSON格式存储
	Domains     []Domain `json:"domains" gorm:"foreignKey:ProxyConfigID" comment:"关联的域名列表"`     // 关联的域名配置列表
}

// Domain 域名模型 - 存储代理域名配置
type Domain struct {
	BaseModel
	ProxyConfigID uuid.UUID   `json:"proxy_config_id" gorm:"type:uuid;not null" comment:"代理配置ID"`     // 所属代理配置ID
	Domain        string      `json:"domain" gorm:"size:255;not null;index" comment:"域名"`             // 域名地址
	IsPrimary     bool        `json:"is_primary" gorm:"default:false" comment:"是否为主域名"`               // 是否为主要域名
	SSLEnabled    bool        `json:"ssl_enabled" gorm:"default:false" comment:"是否启用SSL"`             // SSL/TLS 配置状态
	ProxyConfig   ProxyConfig `json:"proxy_config" gorm:"foreignKey:ProxyConfigID" comment:"关联的代理配置"` // 关联的代理配置对象
}

// Rule 规则模型 - 存储业务规则配置
type Rule struct {
	BaseModel
	Name              string `json:"name" gorm:"size:100;not null" comment:"规则名称"`                        // 规则名称
	Description       string `json:"description" gorm:"type:text" comment:"规则描述"`                         // 规则功能描述
	TriggerConditions string `json:"trigger_conditions" gorm:"type:jsonb;not null" comment:"触发条件，JSON格式"` // 规则触发条件，JSON格式
	Actions           string `json:"actions" gorm:"type:jsonb;not null" comment:"执行动作，JSON格式"`            // 规则执行动作，JSON格式
	IsActive          bool   `json:"is_active" gorm:"default:true" comment:"是否启用"`                        // 规则启用状态
	Priority          int    `json:"priority" gorm:"default:0;index" comment:"优先级，数字越大优先级越高"`             // 规则执行优先级
}

// Popup 弹窗模型 - 存储弹窗配置信息
type Popup struct {
	BaseModel
	Title       string       `json:"title" gorm:"size:200;not null" comment:"弹窗标题"`                     // 弹窗显示标题
	Content     string       `json:"content" gorm:"type:text" comment:"弹窗内容"`                           // 弹窗显示内容
	StyleConfig string       `json:"style_config" gorm:"type:jsonb;default:'{}'" comment:"样式配置，JSON格式"` // 弹窗样式配置
	FormConfig  string       `json:"form_config" gorm:"type:jsonb;default:'{}'" comment:"表单配置，JSON格式"`  // 表单字段配置
	IsActive    bool         `json:"is_active" gorm:"default:true" comment:"是否启用"`                      // 弹窗启用状态
	Submissions []Submission `json:"submissions" gorm:"foreignKey:PopupID" comment:"用户提交的数据列表"`         // 用户提交的数据列表
}

// Submission 数据提交模型 - 存储用户通过弹窗提交的数据
type Submission struct {
	BaseModel
	PopupID     uuid.UUID `json:"popup_id" gorm:"type:uuid;not null;index" comment:"弹窗ID"`    // 所属弹窗ID
	FormData    string    `json:"form_data" gorm:"type:jsonb;not null" comment:"表单数据，JSON格式"` // 用户提交的表单数据
	UserIP      string    `json:"user_ip" gorm:"type:inet" comment:"用户IP地址"`                  // 用户IP地址
	UserAgent   string    `json:"user_agent" gorm:"type:text" comment:"用户浏览器信息"`              // 用户浏览器信息
	ReferrerURL string    `json:"referrer_url" gorm:"type:text" comment:"来源页面URL"`            // 提交时的来源页面
	Popup       Popup     `json:"popup" gorm:"foreignKey:PopupID" comment:"关联的弹窗配置"`          // 关联的弹窗配置对象
}

// ProxyLog 代理日志模型 - 存储代理服务访问日志
type ProxyLog struct {
	BaseModel
	ProxyConfigID uuid.UUID `json:"proxy_config_id" gorm:"type:uuid;index" comment:"代理配置ID"` // 所属代理配置ID
	Method        string    `json:"method" gorm:"size:10;not null" comment:"HTTP请求方法"`       // HTTP请求方法，如GET、POST
	URL           string    `json:"url" gorm:"type:text;not null" comment:"请求URL"`           // 请求的完整URL
	UserAgent     string    `json:"user_agent" gorm:"type:text" comment:"用户浏览器信息"`           // 用户浏览器信息
	UserIP        string    `json:"user_ip" gorm:"type:inet" comment:"用户IP地址"`               // 用户IP地址
	StatusCode    int       `json:"status_code" comment:"HTTP响应状态码"`                         // HTTP响应状态码
	ResponseTime  int64     `json:"response_time" comment:"响应时间，毫秒"`                         // 响应耗时，单位毫秒
	ErrorMessage  string    `json:"error_message" gorm:"type:text" comment:"错误信息"`           // 错误信息，失败时记录
}

// SystemMetric 系统指标模型 - 存储系统监控指标数据
type SystemMetric struct {
	BaseModel
	MetricName  string    `json:"metric_name" gorm:"size:100;not null;index" comment:"指标名称"` // 指标名称，如CPU使用率、内存使用量
	MetricValue float64   `json:"metric_value" comment:"指标数值"`                               // 指标的具体数值
	Tags        string    `json:"tags" gorm:"type:jsonb;default:'{}'" comment:"指标标签，JSON格式"` // 指标标签，用于分类和过滤
	Timestamp   time.Time `json:"timestamp" gorm:"index" comment:"指标记录时间"`                   // 指标记录的时间戳
}

// TableName 设置表名 - 定义各模型对应的数据库表名

func (User) TableName() string {
	return "users" // 用户表
}

func (Role) TableName() string {
	return "roles" // 角色表
}

func (Permission) TableName() string {
	return "permissions" // 权限表
}

func (UserRole) TableName() string {
	return "user_roles" // 用户角色关联表
}

func (RolePermission) TableName() string {
	return "role_permissions" // 角色权限关联表
}

func (ProxyConfig) TableName() string {
	return "proxy_configs" // 代理配置表
}

func (Domain) TableName() string {
	return "domains" // 域名配置表
}

func (Rule) TableName() string {
	return "rules" // 业务规则表
}

func (Popup) TableName() string {
	return "popups" // 弹窗配置表
}

func (Submission) TableName() string {
	return "submissions" // 数据提交表
}

func (ProxyLog) TableName() string {
	return "proxy_logs" // 代理访问日志表
}

func (SystemMetric) TableName() string {
	return "system_metrics" // 系统监控指标表
}
