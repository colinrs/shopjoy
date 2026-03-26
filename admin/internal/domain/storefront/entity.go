package storefront

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type Shop struct {
	ID              int64
	TenantID        shared.TenantID
	Name            string
	Description     string
	Logo            string
	Banner          string
	ContactPhone    string
	ContactEmail    string
	Address         string
	SocialLinks     map[string]string
	SEO             SEOConfig
	Status          shared.Status
	CurrentThemeID  *int64
	ThemeConfig     *ThemeConfig
	DeletedAt       *int64
	Audit           shared.AuditInfo `gorm:"embedded"`
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
	ID            int64
	TenantID      shared.TenantID
	Name          string
	Code          string
	Description   string
	Thumbnail     string
	PreviewImage  string
	Config        ThemeConfig
	ConfigSchema  ThemeConfigSchema
	DefaultConfig ThemeConfig
	IsActive      bool
	IsCustom      bool
	IsPreset      bool
	DeletedAt     *int64
}

func (t *Theme) TableName() string {
	return "themes"
}

type ThemeConfig struct {
	Colors     map[string]string      `json:"colors"`
	Fonts      map[string]string      `json:"fonts"`
	Layout     string                 `json:"layout"`
	CustomCSS  string                 `json:"custom_css"`
	Components map[string]interface{} `json:"components"`
}

type ThemeConfigSchema struct {
	Fields []ThemeConfigField `json:"fields"`
}

type ThemeConfigField struct {
	Key     string      `json:"key"`
	Label   string      `json:"label"`
	Type    string      `json:"type"` // color, select, text
	Options []SelectOpt `json:"options,omitempty"`
	Default string      `json:"default"`
}

type SelectOpt struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Page struct {
	ID          int64
	TenantID    shared.TenantID
	Name        string
	Slug        string
	Type        PageType
	Content     string
	SEO         SEOConfig
	Status      shared.Status
	Sort        int
	IsPublished bool
	PublishedAt *int64
	Version     int
	Audit       shared.AuditInfo `gorm:"embedded"`
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

// Decoration represents a decoration block on a page
type Decoration struct {
	ID          int64
	TenantID    shared.TenantID
	PageID      int64
	BlockType   string      // banner, product_grid, rich_text, image_carousel, etc.
	BlockConfig BlockConfig // JSON configuration
	SortOrder   int
	IsActive    bool
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   *int64
}

func (d *Decoration) TableName() string {
	return "decorations"
}

// BlockConfig represents the configuration of a decoration block
type BlockConfig map[string]interface{}

// PageVersion represents a snapshot of page decoration at a point in time
type PageVersion struct {
	ID        int64
	TenantID  shared.TenantID
	PageID    int64
	Version   int
	Blocks    []BlockSnapshot // JSON snapshot of blocks
	CreatedBy int64
	CreatedAt int64
	DeletedAt *int64
}

func (v *PageVersion) TableName() string {
	return "page_versions"
}

// BlockSnapshot represents a snapshot of a decoration block for versioning
type BlockSnapshot struct {
	BlockType   string                 `json:"block_type"`
	BlockConfig map[string]interface{} `json:"block_config"`
	SortOrder   int                    `json:"sort_order"`
}

// SEOConfigEntity represents SEO configuration for a page type
type SEOConfigEntity struct {
	ID          int64
	TenantID    shared.TenantID
	PageType    string // global, home, category, product, custom
	PageID      *int64 // NULL for global/page type defaults
	Title       string
	Description string
	Keywords    string
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   *int64
}

func (s *SEOConfigEntity) TableName() string {
	return "seo_configs"
}

type Navigation struct {
	ID        int64
	TenantID  shared.TenantID
	Name      string
	Position  string
	Items     []NavItem
	Status    shared.Status
	DeletedAt *int64
}

func (n *Navigation) TableName() string {
	return "navigations"
}

type NavItem struct {
	ID        int64
	NavID     int64
	ParentID  int64
	Name      string
	Link      string
	Type      string
	TargetID  int64
	Sort      int
	DeletedAt *int64
}

func (ni *NavItem) TableName() string {
	return "nav_items"
}

// BlockOrder represents the sort order for reordering blocks
type BlockOrder struct {
	ID        int64
	SortOrder int
}

// ThemeAuditLog represents an audit log entry for theme changes
type ThemeAuditLog struct {
	ID         int64
	TenantID   shared.TenantID
	Action     string // switch_theme, update_config
	ThemeID    int64
	ThemeName  string
	ThemeCode  string
	OldConfig  string // JSON string
	NewConfig  string // JSON string
	UserID     int64
	UserName   string
	IPAddress  string
	UserAgent  string
	CreatedAt  int64
	DeletedAt  *int64
}

func (l *ThemeAuditLog) TableName() string {
	return "theme_audit_logs"
}

// Audit action constants
const (
	AuditActionSwitchTheme  = "switch_theme"
	AuditActionUpdateConfig = "update_config"
)

// Repository interfaces

type ShopRepository interface {
	FindByTenantID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Shop, error)
	Save(ctx context.Context, db *gorm.DB, shop *Shop) error
}

type ThemeRepository interface {
	Create(ctx context.Context, db *gorm.DB, theme *Theme) error
	Update(ctx context.Context, db *gorm.DB, theme *Theme) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Theme, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Theme, error)
	FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Theme, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Theme, error)
	FindPresets(ctx context.Context, db *gorm.DB) ([]*Theme, error)
	SetActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
}

type PageRepository interface {
	Create(ctx context.Context, db *gorm.DB, page *Page) error
	Update(ctx context.Context, db *gorm.DB, page *Page) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Page, error)
	FindBySlug(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, slug string) (*Page, error)
	FindByType(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType PageType) (*Page, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*Page, int64, error)
	CountAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error)
}

type DecorationRepository interface {
	Create(ctx context.Context, db *gorm.DB, d *Decoration) error
	Update(ctx context.Context, db *gorm.DB, d *Decoration) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Decoration, error)
	FindByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) ([]*Decoration, error)
	Reorder(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orders []BlockOrder) error
	DeleteByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) error
}

type PageVersionRepository interface {
	Create(ctx context.Context, db *gorm.DB, v *PageVersion) error
	FindByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, page, pageSize int) ([]*PageVersion, int64, error)
	FindByVersion(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, version int) (*PageVersion, error)
	DeleteOldest(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, keepCount int) error
	CountByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) (int64, error)
}

type SEOConfigRepository interface {
	Save(ctx context.Context, db *gorm.DB, config *SEOConfigEntity) error
	FindByPageType(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType string, pageID *int64) (*SEOConfigEntity, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*SEOConfigEntity, int64, error)
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType string, pageID *int64) error
	CountAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error)
}

type ThemeAuditLogRepository interface {
	Create(ctx context.Context, db *gorm.DB, log *ThemeAuditLog) error
	FindByTenantID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*ThemeAuditLog, int64, error)
}

type NavigationRepository interface {
	Save(ctx context.Context, db *gorm.DB, nav *Navigation) error
	FindByPosition(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, position string) (*Navigation, error)
}
