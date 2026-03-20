package storefront

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type Shop struct {
	ID           int64
	TenantID     shared.TenantID
	Name         string
	Description  string
	Logo         string
	Banner       string
	ContactPhone string
	ContactEmail string
	Address      string
	SocialLinks  map[string]string
	SEO          SEOConfig
	Status       shared.Status
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (s *Shop) TableName() string {
	return "shops"
}

type SEOConfig struct {
	Title       string
	Description string
	Keywords    []string
}

type Theme struct {
	ID          int64
	TenantID    shared.TenantID
	Name        string
	Code        string
	Description string
	Thumbnail   string
	Config      ThemeConfig
	IsActive    bool
	IsCustom    bool
}

func (t *Theme) TableName() string {
	return "themes"
}

type ThemeConfig struct {
	Colors     map[string]string
	Fonts      map[string]string
	Layout     string
	CustomCSS  string
	Components map[string]interface{}
}

type Page struct {
	ID       int64
	TenantID shared.TenantID
	Name     string
	Slug     string
	Type     PageType
	Content  string
	SEO      SEOConfig
	Status   shared.Status
	Sort     int
	Audit    shared.AuditInfo `gorm:"embedded"`
}

type PageType int

const (
	PageTypeHome PageType = iota
	PageTypeProduct
	PageTypeCollection
	PageTypeCustom
)

func (p *Page) TableName() string {
	return "pages"
}

type Navigation struct {
	ID       int64
	TenantID shared.TenantID
	Name     string
	Position string
	Items    []NavItem
	Status   shared.Status
}

func (n *Navigation) TableName() string {
	return "navigations"
}

type NavItem struct {
	ID       int64
	NavID    int64
	ParentID int64
	Name     string
	Link     string
	Type     string
	TargetID int64
	Sort     int
}

func (ni *NavItem) TableName() string {
	return "nav_items"
}

type ShopRepository interface {
	FindByTenantID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Shop, error)
	Save(ctx context.Context, db *gorm.DB, shop *Shop) error
}

type ThemeRepository interface {
	Create(ctx context.Context, db *gorm.DB, theme *Theme) error
	Update(ctx context.Context, db *gorm.DB, theme *Theme) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Theme, error)
	FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Theme, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Theme, error)
	SetActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
}

type PageRepository interface {
	Create(ctx context.Context, db *gorm.DB, page *Page) error
	Update(ctx context.Context, db *gorm.DB, page *Page) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Page, error)
	FindBySlug(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, slug string) (*Page, error)
	FindByType(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType PageType) (*Page, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Page, error)
}

type NavigationRepository interface {
	Save(ctx context.Context, db *gorm.DB, nav *Navigation) error
	FindByPosition(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, position string) (*Navigation, error)
}