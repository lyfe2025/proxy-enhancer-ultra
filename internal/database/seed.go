package database

import (
	"fmt"
	"proxy-enhancer-ultra/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Seed 创建默认数据
func (d *Database) Seed() error {
	d.logger.Info("Starting database seeding...")

	// 创建默认角色
	if err := d.createDefaultRoles(); err != nil {
		return fmt.Errorf("failed to create default roles: %w", err)
	}

	// 创建默认管理员用户
	if err := d.createDefaultAdmin(); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	d.logger.Info("Database seeding completed successfully")
	return nil
}

// createDefaultRoles 创建默认角色
func (d *Database) createDefaultRoles() error {
	// 创建admin角色
	var adminCount int64
	d.DB.Model(&models.Role{}).Where("name = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		adminRole := &models.Role{
			Name:        "admin",
			Description: "系统管理员角色",
		}
		if err := d.DB.Create(adminRole).Error; err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}
		d.logger.Info("Admin role created successfully")
	} else {
		d.logger.Info("Admin role already exists, skipping creation")
	}

	// 创建user角色
	var userCount int64
	d.DB.Model(&models.Role{}).Where("name = ?", "user").Count(&userCount)
	if userCount == 0 {
		userRole := &models.Role{
			Name:        "user",
			Description: "普通用户角色",
		}
		if err := d.DB.Create(userRole).Error; err != nil {
			return fmt.Errorf("failed to create user role: %w", err)
		}
		d.logger.Info("User role created successfully")
	} else {
		d.logger.Info("User role already exists, skipping creation")
	}

	return nil
}

// createDefaultAdmin 创建默认管理员用户
func (d *Database) createDefaultAdmin() error {
	// 检查是否已存在管理员用户
	var count int64
	d.DB.Model(&models.User{}).Where("username = ? OR email = ?", "admin", "admin@example.com").Count(&count)
	if count > 0 {
		d.logger.Info("Admin user already exists, skipping creation")
		// 检查是否已分配admin角色
		var admin models.User
		d.DB.Where("username = ?", "admin").First(&admin)
		return d.assignAdminRole(admin.ID)
	}

	// 创建默认管理员
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := &models.User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	if err := d.DB.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// 为admin用户分配admin角色
	if err := d.assignAdminRole(admin.ID); err != nil {
		return fmt.Errorf("failed to assign admin role: %w", err)
	}

	d.logger.Info("Default admin user created successfully (username: admin, password: admin123456)")
	return nil
}

// assignAdminRole 为用户分配admin角色
func (d *Database) assignAdminRole(userID uuid.UUID) error {
	// 获取admin角色
	var adminRole models.Role
	if err := d.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to find admin role: %w", err)
	}

	// 检查是否已经分配了admin角色
	var existingUserRole models.UserRole
	err := d.DB.Where("user_id = ? AND role_id = ?", userID, adminRole.ID).First(&existingUserRole).Error
	if err == nil {
		d.logger.Info("Admin role already assigned to user, skipping")
		return nil
	}

	// 创建用户角色关联
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: adminRole.ID,
	}

	if err := d.DB.Create(userRole).Error; err != nil {
		return fmt.Errorf("failed to assign admin role to user: %w", err)
	}

	d.logger.Info("Admin role assigned to user successfully")
	return nil
}
