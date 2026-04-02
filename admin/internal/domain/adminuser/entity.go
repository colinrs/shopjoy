package adminuser

import (
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"golang.org/x/crypto/bcrypt"
)

// Type 管理员类型
type Type int

const (
	TypeSuperAdmin  Type = 1 // 平台超级管理员
	TypeTenantAdmin Type = 2 // 商家管理员
	TypeTenantSub   Type = 3 // 商家子账号
)

func (t Type) String() string {
	switch t {
	case TypeSuperAdmin:
		return "super_admin"
	case TypeTenantAdmin:
		return "tenant_admin"
	case TypeTenantSub:
		return "tenant_sub"
	default:
		return "unknown"
	}
}

// Status 用户状态
type Status int

const (
	StatusActive   Status = 1 // 正常
	StatusDisabled Status = 2 // 禁用
	StatusDeleted  Status = 3 // 已删除
)

// AdminUser 管理后台用户实体
type AdminUser struct {
	application.Model
	TenantID    int64      `gorm:"column:tenant_id;not null;default:0;index"`
	Username    string     `gorm:"column:username;uniqueIndex;size:64"`
	Email       string     `gorm:"column:email;uniqueIndex;not null;size:128"`
	Mobile      string     `gorm:"column:mobile;uniqueIndex;size:20"`
	Password    string     `gorm:"column:password;not null;size:255"`
	RealName    string     `gorm:"column:real_name;size:32"`
	Avatar      string     `gorm:"column:avatar;size:255"`
	Type        Type       `gorm:"column:type;not null;default:1;index"`
	Status      Status     `gorm:"column:status;not null;default:1;index"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	LastLoginIP string     `gorm:"column:last_login_ip;size:45"`
}

func (u *AdminUser) TableName() string {
	return "admin_users"
}

// SetPassword 设置密码（bcrypt加密）
func (u *AdminUser) SetPassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return code.ErrAdminPasswordTooWeak
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// CheckPassword 验证密码
func (u *AdminUser) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

// CanLogin 是否可以登录
func (u *AdminUser) CanLogin() bool {
	return u.Status == StatusActive
}

// UpdateLastLogin 更新最后登录时间
func (u *AdminUser) UpdateLastLogin(ip string) {
	now := time.Now().UTC()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
}

// IsSuperAdmin 是否平台超管
func (u *AdminUser) IsSuperAdmin() bool {
	return u.Type == TypeSuperAdmin && u.TenantID == 0
}

// IsTenantAdmin 是否商家管理员
func (u *AdminUser) IsTenantAdmin() bool {
	return u.Type == TypeTenantAdmin && u.TenantID > 0
}

// IsTenantSub 是否商家子账号
func (u *AdminUser) IsTenantSub() bool {
	return u.Type == TypeTenantSub && u.TenantID > 0
}

// CanManageTenant 是否可以管理指定租户
func (u *AdminUser) CanManageTenant(tenantID int64) bool {
	if u.IsSuperAdmin() {
		return true
	}
	return u.TenantID == tenantID
}

// Disable 禁用账号
func (u *AdminUser) Disable() error {
	if u.Status == StatusDeleted {
		return code.ErrAdminAlreadyDeleted
	}
	u.Status = StatusDisabled
	return nil
}

// Enable 启用账号
func (u *AdminUser) Enable() error {
	if u.Status == StatusDeleted {
		return code.ErrAdminAlreadyDeleted
	}
	u.Status = StatusActive
	return nil
}

// UpdateProfile 更新基础信息
func (u *AdminUser) UpdateProfile(realName, avatar, mobile string) {
	if realName != "" {
		u.RealName = realName
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	if mobile != "" {
		u.Mobile = mobile
	}
}
