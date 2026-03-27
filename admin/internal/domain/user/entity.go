package user

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Status int

const (
	StatusInactive Status = iota
	StatusActive
	StatusSuspended
	StatusDeleted
)

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
	GenderOther
)

type User struct {
	ID        int64
	TenantID  shared.TenantID
	Email     string
	Phone     string
	Password  string
	Name      string
	Avatar    string
	Gender    Gender
	Birthday  *time.Time
	Status    Status
	LastLogin *time.Time
	DeletedAt gorm.DeletedAt
	Audit     shared.AuditInfo `gorm:"embedded"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) SetPassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return code.ErrUserPasswordTooWeak
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

func (u *User) CanLogin() bool {
	return u.Status == StatusActive
}

func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLogin = &now
}

func (u *User) Suspend() error {
	if u.DeletedAt.Valid {
		return code.ErrUserAlreadyDeleted
	}
	u.Status = StatusSuspended
	return nil
}

func (u *User) Activate() error {
	if u.DeletedAt.Valid {
		return code.ErrUserAlreadyDeleted
	}
	u.Status = StatusActive
	return nil
}

func (u *User) SoftDelete() error {
	if u.DeletedAt.Valid {
		return code.ErrUserAlreadyDeleted
	}
	u.DeletedAt = gorm.DeletedAt{Time: time.Now().UTC()}
	return nil
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, user *User) error
	Update(ctx context.Context, db *gorm.DB, user *User) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*User, error)
	FindByEmail(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email string) (*User, error)
	FindByPhone(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, phone string) (*User, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*User, int64, error)
	Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email, phone string) (bool, error)
	GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Stats, error)
}

type Query struct {
	shared.PageQuery
	Name   string
	Email  string
	Phone  string
	Status Status
}

type Stats struct {
	Total     int64
	Active    int64
	Suspended int64
	NewToday  int64
}
