package storefront

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// DTO types for application layer

type ThemeDTO struct {
	ID           int64               `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	PreviewImage string              `json:"preview_image"`
	Thumbnail    string              `json:"thumbnail"`
	IsPreset     bool                `json:"is_preset"`
	IsCurrent    bool                `json:"is_current"`
	Config       *ThemeConfigDTO     `json:"config,omitempty"`
	ConfigSchema *ThemeConfigSchemaDTO `json:"config_schema,omitempty"`
}

type ThemeConfigDTO struct {
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	FontFamily     string `json:"font_family"`
	ButtonStyle    string `json:"button_style"`
}

type ThemeConfigSchemaDTO struct {
	Fields []ThemeConfigFieldDTO `json:"fields"`
}

type ThemeConfigFieldDTO struct {
	Key     string        `json:"key"`
	Label   string        `json:"label"`
	Type    string        `json:"type"`
	Options []SelectOptDTO `json:"options,omitempty"`
	Default string        `json:"default"`
}

type SelectOptDTO struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type CurrentThemeDTO struct {
	Theme  *ThemeDTO     `json:"theme"`
	Config ThemeConfigDTO `json:"config"`
}

type PageDTO struct {
	ID          int64            `json:"id"`
	PageType    string           `json:"page_type"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	IsPublished bool             `json:"is_published"`
	Version     int              `json:"version"`
}

type PageDetailDTO struct {
	Page        *PageDTO        `json:"page"`
	Decorations []*DecorationDTO `json:"decorations"`
}

type DecorationDTO struct {
	ID          int64                  `json:"id"`
	BlockType   string                 `json:"block_type"`
	BlockConfig map[string]any         `json:"block_config"`
	SortOrder   int                    `json:"sort_order"`
}

type BlockOrderDTO struct {
	ID        int64 `json:"id"`
	SortOrder int   `json:"sort_order"`
}

type VersionDTO struct {
	ID        int64  `json:"id"`
	Version   int    `json:"version"`
	CreatedBy int64  `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
}

type VersionDetailDTO struct {
	Version *VersionDTO     `json:"version"`
	Blocks  []*DecorationDTO `json:"blocks"`
}

type SEOConfigDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

type PageSEOConfigDTO struct {
	PageType string        `json:"page_type"`
	PageID   *int64        `json:"page_id,omitempty"`
	Config   SEOConfigDTO  `json:"config"`
}

// Service interfaces

type ThemeService interface {
	ListThemes(ctx context.Context, tenantID shared.TenantID) ([]*ThemeDTO, error)
	GetCurrentTheme(ctx context.Context, tenantID shared.TenantID) (*CurrentThemeDTO, error)
	SwitchTheme(ctx context.Context, tenantID shared.TenantID, themeID int64) error
	UpdateThemeConfig(ctx context.Context, tenantID shared.TenantID, config ThemeConfigDTO) error
}

type PageService interface {
	ListPages(ctx context.Context, tenantID shared.TenantID) ([]*PageDTO, error)
	GetPage(ctx context.Context, tenantID shared.TenantID, pageID int64) (*PageDetailDTO, error)
	GetPageBySlug(ctx context.Context, tenantID shared.TenantID, slug string) (*PageDetailDTO, error)
	SaveDraft(ctx context.Context, tenantID shared.TenantID, pageID int64, blocks []*DecorationDTO, userID int64) error
	PublishPage(ctx context.Context, tenantID shared.TenantID, pageID int64, userID int64) error
	UnpublishPage(ctx context.Context, tenantID shared.TenantID, pageID int64) error
}

type DecorationService interface {
	GetDecorations(ctx context.Context, tenantID shared.TenantID, pageID int64) ([]*DecorationDTO, error)
	AddDecoration(ctx context.Context, tenantID shared.TenantID, pageID int64, blockType string, blockConfig map[string]any, sortOrder int) (*DecorationDTO, error)
	UpdateDecoration(ctx context.Context, tenantID shared.TenantID, decorationID int64, blockConfig map[string]any) error
	DeleteDecoration(ctx context.Context, tenantID shared.TenantID, decorationID int64) error
	ReorderBlocks(ctx context.Context, tenantID shared.TenantID, pageID int64, orders []BlockOrderDTO) error
}

type VersionService interface {
	ListVersions(ctx context.Context, tenantID shared.TenantID, pageID int64, limit int) ([]*VersionDTO, error)
	GetVersion(ctx context.Context, tenantID shared.TenantID, pageID int64, version int) (*VersionDetailDTO, error)
	RestoreVersion(ctx context.Context, tenantID shared.TenantID, pageID int64, version int, userID int64) error
}

type SEOService interface {
	GetGlobalSEO(ctx context.Context, tenantID shared.TenantID) (*SEOConfigDTO, error)
	UpdateGlobalSEO(ctx context.Context, tenantID shared.TenantID, config SEOConfigDTO) error
	GetPageSEO(ctx context.Context, tenantID shared.TenantID, pageType string, pageID *int64) (*PageSEOConfigDTO, error)
	UpdatePageSEO(ctx context.Context, tenantID shared.TenantID, pageType string, pageID *int64, config SEOConfigDTO) error
	ListPageSEO(ctx context.Context, tenantID shared.TenantID) ([]*PageSEOConfigDTO, error)
}

// Service implementation

type themeService struct {
	db        *gorm.DB
	themeRepo storefront.ThemeRepository
	shopRepo  storefront.ShopRepository
}

func NewThemeService(db *gorm.DB, themeRepo storefront.ThemeRepository, shopRepo storefront.ShopRepository) ThemeService {
	return &themeService{
		db:        db,
		themeRepo: themeRepo,
		shopRepo:  shopRepo,
	}
}

func (s *themeService) ListThemes(ctx context.Context, tenantID shared.TenantID) ([]*ThemeDTO, error) {
	themes, err := s.themeRepo.FindAll(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	shop, _ := s.shopRepo.FindByTenantID(ctx, s.db, tenantID)
	var currentThemeID int64
	if shop != nil && shop.CurrentThemeID != nil {
		currentThemeID = *shop.CurrentThemeID
	}

	dtos := make([]*ThemeDTO, len(themes))
	for i, t := range themes {
		dtos[i] = &ThemeDTO{
			ID:           t.ID,
			Code:         t.Code,
			Name:         t.Name,
			Description:  t.Description,
			PreviewImage: t.PreviewImage,
			Thumbnail:    t.Thumbnail,
			IsPreset:     t.IsPreset,
			IsCurrent:    t.ID == currentThemeID,
		}
	}
	return dtos, nil
}

func (s *themeService) GetCurrentTheme(ctx context.Context, tenantID shared.TenantID) (*CurrentThemeDTO, error) {
	shop, err := s.shopRepo.FindByTenantID(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, nil
	}

	var themeID int64 = 1001 // Default to Classic theme
	if shop.CurrentThemeID != nil {
		themeID = *shop.CurrentThemeID
	}

	theme, err := s.themeRepo.FindByID(ctx, s.db, tenantID, themeID)
	if err != nil {
		return nil, err
	}
	if theme == nil {
		return nil, nil
	}

	config := themeConfigToDTO(shop.ThemeConfig)
	if config == nil {
		config = themeConfigToDTO(&theme.DefaultConfig)
	}

	return &CurrentThemeDTO{
		Theme: &ThemeDTO{
			ID:           theme.ID,
			Code:         theme.Code,
			Name:         theme.Name,
			Description:  theme.Description,
			PreviewImage: theme.PreviewImage,
			Thumbnail:    theme.Thumbnail,
			IsPreset:     theme.IsPreset,
			IsCurrent:    true,
		},
		Config: *config,
	}, nil
}

func (s *themeService) SwitchTheme(ctx context.Context, tenantID shared.TenantID, themeID int64) error {
	theme, err := s.themeRepo.FindByID(ctx, s.db, tenantID, themeID)
	if err != nil {
		return err
	}
	if theme == nil {
		return code.ErrThemeNotFound
	}

	shop, err := s.shopRepo.FindByTenantID(ctx, s.db, tenantID)
	if err != nil {
		return err
	}
	if shop == nil {
		return code.ErrShopNotFound
	}

	shop.CurrentThemeID = &themeID
	return s.shopRepo.Save(ctx, s.db, shop)
}

func (s *themeService) UpdateThemeConfig(ctx context.Context, tenantID shared.TenantID, config ThemeConfigDTO) error {
	shop, err := s.shopRepo.FindByTenantID(ctx, s.db, tenantID)
	if err != nil {
		return err
	}
	if shop == nil {
		return code.ErrShopNotFound
	}

	themeConfig := dtoToThemeConfig(config)
	shop.ThemeConfig = &themeConfig
	return s.shopRepo.Save(ctx, s.db, shop)
}

func themeConfigToDTO(config *storefront.ThemeConfig) *ThemeConfigDTO {
	if config == nil {
		return nil
	}

	primaryColor := ""
	secondaryColor := ""
	fontFamily := ""
	buttonStyle := ""

	if colors, ok := config.Colors["primary"]; ok {
		primaryColor = colors
	}
	if colors, ok := config.Colors["secondary"]; ok {
		secondaryColor = colors
	}
	if fonts, ok := config.Fonts["heading"]; ok {
		fontFamily = fonts
	}
	if comp, ok := config.Components["button_style"]; ok {
		if bs, ok := comp.(string); ok {
			buttonStyle = bs
		}
	}

	return &ThemeConfigDTO{
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
		FontFamily:     fontFamily,
		ButtonStyle:    buttonStyle,
	}
}

func dtoToThemeConfig(dto ThemeConfigDTO) storefront.ThemeConfig {
	return storefront.ThemeConfig{
		Colors: map[string]string{
			"primary":   dto.PrimaryColor,
			"secondary": dto.SecondaryColor,
		},
		Fonts: map[string]string{
			"heading": dto.FontFamily,
			"body":    dto.FontFamily,
		},
		Components: map[string]any{
			"button_style": dto.ButtonStyle,
		},
	}
}