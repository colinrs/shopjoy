package product

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrDuplicateCategory = errors.New("duplicate category")
	ErrInvalidCategory   = errors.New("invalid category")
	ErrHasChildren       = errors.New("category has children")
	ErrHasProducts       = errors.New("category has products")
)

type CategoryStatus int

const (
	CategoryStatusDisabled CategoryStatus = iota
	CategoryStatusEnabled
)

type Category struct {
	ID       int64
	TenantID shared.TenantID
	ParentID int64
	Name     string
	Code     string
	Level    int
	Sort     int
	Icon     string
	Image    string
	Status   CategoryStatus
	Audit    shared.AuditInfo
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) Enable() {
	c.Status = CategoryStatusEnabled
}

func (c *Category) Disable() {
	c.Status = CategoryStatusDisabled
}

func (c *Category) IsActive() bool {
	return c.Status == CategoryStatusEnabled
}

func (c *Category) IsRoot() bool {
	return c.ParentID == 0
}

type Brand struct {
	ID          int64
	TenantID    shared.TenantID
	Name        string
	Logo        string
	Description string
	Website     string
	Sort        int
	Status      shared.Status
	Audit       shared.AuditInfo
}

func (b *Brand) TableName() string {
	return "brands"
}

type Attribute struct {
	ID         int64
	TenantID   shared.TenantID
	Name       string
	Code       string
	InputType  AttributeInputType
	Options    []string
	IsRequired bool
	Status     shared.Status
	Audit      shared.AuditInfo
}

type AttributeInputType int

const (
	InputTypeText AttributeInputType = iota
	InputTypeSelect
	InputTypeMultiSelect
	InputTypeNumber
	InputTypeBoolean
)

type SKU struct {
	ID         int64
	ProductID  int64
	Code       string
	Price      shared.Money
	Stock      int
	Attributes map[string]string
	Status     shared.Status
	Audit      shared.AuditInfo
}

func (s *SKU) TableName() string {
	return "skus"
}

func (s *SKU) IsAvailable() bool {
	return s.Status == shared.StatusEnabled && s.Stock > 0
}

func (s *SKU) DeductStock(quantity int) error {
	if s.Stock < quantity {
		return ErrInsufficientStock
	}
	s.Stock -= quantity
	return nil
}

func (s *SKU) RestoreStock(quantity int) {
	s.Stock += quantity
}

type CategoryRepository interface {
	Create(ctx context.Context, db *gorm.DB, category *Category) error
	Update(ctx context.Context, db *gorm.DB, category *Category) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Category, error)
	FindByParentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, parentID int64) ([]*Category, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Category, error)
	FindTree(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Category, error)
}

type BrandRepository interface {
	Create(ctx context.Context, db *gorm.DB, brand *Brand) error
	Update(ctx context.Context, db *gorm.DB, brand *Brand) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Brand, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query BrandQuery) ([]*Brand, int64, error)
}

type BrandQuery struct {
	shared.PageQuery
	Name   string
	Status shared.Status
}
