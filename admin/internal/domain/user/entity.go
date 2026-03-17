package user

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPhone    = errors.New("invalid phone format")
	ErrPasswordTooWeak = errors.New("password too weak")
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateUser   = errors.New("duplicate user")
	ErrWrongPassword   = errors.New("wrong password")
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
	Audit     shared.AuditInfo
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) SetPassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return ErrPasswordTooWeak
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
	if u.Status == StatusDeleted {
		return errors.New("user already deleted")
	}
	u.Status = StatusSuspended
	return nil
}

func (u *User) Activate() error {
	if u.Status == StatusDeleted {
		return errors.New("user already deleted")
	}
	u.Status = StatusActive
	return nil
}

func (u *User) SoftDelete() error {
	u.Status = StatusDeleted
	return nil
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, user *User) error
	Update(ctx context.Context, db *gorm.DB, user *User) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*User, error)
	FindByEmail(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email string) (*User, error)
	FindByPhone(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, phone string) (*User, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*User, int64, error)
	Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email, phone string) (bool, error)
}

type Query struct {
	shared.PageQuery
	Name   string
	Email  string
	Phone  string
	Status Status
}
