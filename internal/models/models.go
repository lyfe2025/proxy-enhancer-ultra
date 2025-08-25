package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username     string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string     `json:"-" gorm:"size:255;not null"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	UserRoles    []UserRole `json:"user_roles" gorm:"foreignKey:UserID"`
}

// Role 角色模型
type Role struct {
	ID              int              `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string           `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description     string           `json:"description" gorm:"type:text"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UserRoles       []UserRole       `json:"user_roles" gorm:"foreignKey:RoleID"`
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:RoleID"`
}

// Permission 权限模型
type Permission struct {
	ID              int              `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string           `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Resource        string           `json:"resource" gorm:"size:50;not null"`
	Action          string           `json:"action" gorm:"size:50;not null"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:PermissionID"`
}

// UserRole 用户角色关联模型
type UserRole struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	RoleID    int       `json:"role_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	RoleID       int        `json:"role_id" gorm:"primaryKey"`
	PermissionID int        `json:"permission_id" gorm:"primaryKey"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Role         Role       `json:"role" gorm:"foreignKey:RoleID"`
	Permission   Permission `json:"permission" gorm:"foreignKey:PermissionID"`
}

// ProxyConfig 代理配置模型
type ProxyConfig struct {
	BaseModel
	Name        string   `json:"name" gorm:"size:100;not null"`
	TargetURL   string   `json:"target_url" gorm:"size:500;not null"`
	ProxyDomain string   `json:"proxy_domain" gorm:"size:255;not null"`
	IsActive    bool     `json:"is_active" gorm:"default:false"`
	Settings    string   `json:"settings" gorm:"type:jsonb;default:'{}'"`
	Domains     []Domain `json:"domains" gorm:"foreignKey:ProxyConfigID"`
}

// Domain 域名模型
type Domain struct {
	BaseModel
	ProxyConfigID uuid.UUID `json:"proxy_config_id" gorm:"type:uuid;not null"`
	Domain        string    `json:"domain" gorm:"size:255;not null;index"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false"`
	SSLEnabled    bool      `json:"ssl_enabled" gorm:"default:false"`
	ProxyConfig   ProxyConfig `json:"proxy_config" gorm:"foreignKey:ProxyConfigID"`
}

// Rule 规则模型
type Rule struct {
	BaseModel
	Name               string `json:"name" gorm:"size:100;not null"`
	Description        string `json:"description" gorm:"type:text"`
	TriggerConditions  string `json:"trigger_conditions" gorm:"type:jsonb;not null"`
	Actions            string `json:"actions" gorm:"type:jsonb;not null"`
	IsActive           bool   `json:"is_active" gorm:"default:true"`
	Priority           int    `json:"priority" gorm:"default:0;index"`
}

// Popup 弹窗模型
type Popup struct {
	BaseModel
	Title       string       `json:"title" gorm:"size:200;not null"`
	Content     string       `json:"content" gorm:"type:text"`
	StyleConfig string       `json:"style_config" gorm:"type:jsonb;default:'{}'"`
	FormConfig  string       `json:"form_config" gorm:"type:jsonb;default:'{}'"`
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	Submissions []Submission `json:"submissions" gorm:"foreignKey:PopupID"`
}

// Submission 数据提交模型
type Submission struct {
	BaseModel
	PopupID     uuid.UUID `json:"popup_id" gorm:"type:uuid;not null;index"`
	FormData    string    `json:"form_data" gorm:"type:jsonb;not null"`
	UserIP      string    `json:"user_ip" gorm:"type:inet"`
	UserAgent   string    `json:"user_agent" gorm:"type:text"`
	ReferrerURL string    `json:"referrer_url" gorm:"type:text"`
	Popup       Popup     `json:"popup" gorm:"foreignKey:PopupID"`
}

// ProxyLog 代理日志模型
type ProxyLog struct {
	BaseModel
	ProxyConfigID uuid.UUID `json:"proxy_config_id" gorm:"type:uuid;index"`
	Method        string    `json:"method" gorm:"size:10;not null"`
	URL           string    `json:"url" gorm:"type:text;not null"`
	UserAgent     string    `json:"user_agent" gorm:"type:text"`
	UserIP        string    `json:"user_ip" gorm:"type:inet"`
	StatusCode    int       `json:"status_code"`
	ResponseTime  int64     `json:"response_time"` // 毫秒
	ErrorMessage  string    `json:"error_message" gorm:"type:text"`
}

// SystemMetric 系统指标模型
type SystemMetric struct {
	BaseModel
	MetricName  string  `json:"metric_name" gorm:"size:100;not null;index"`
	MetricValue float64 `json:"metric_value"`
	Tags        string  `json:"tags" gorm:"type:jsonb;default:'{}'"`
	Timestamp   time.Time `json:"timestamp" gorm:"index"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (ProxyConfig) TableName() string {
	return "proxy_configs"
}

func (Domain) TableName() string {
	return "domains"
}

func (Rule) TableName() string {
	return "rules"
}

func (Popup) TableName() string {
	return "popups"
}

func (Submission) TableName() string {
	return "submissions"
}

func (ProxyLog) TableName() string {
	return "proxy_logs"
}

func (SystemMetric) TableName() string {
	return "system_metrics"
}