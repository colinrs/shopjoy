package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type themeRepo struct{}

func NewThemeRepository() storefront.ThemeRepository {
	return &themeRepo{}
}

type themeModel struct {
	ID            int64  `gorm:"column:id;primaryKey"`
	TenantID      int64  `gorm:"column:tenant_id;not null;index"`
	Name          string `gorm:"column:name;type:varchar(100);not null"`
	Code          string `gorm:"column:code;type:varchar(100);not null;index"`
	Description   string `gorm:"column:description;type:text"`
	Thumbnail     string `gorm:"column:thumbnail;type:varchar(500)"`
	PreviewImage  string `gorm:"column:preview_image;type:varchar(500)"`
	Config        string `gorm:"column:config;type:json"`
	ConfigSchema  string `gorm:"column:config_schema;type:text"`
	DefaultConfig string `gorm:"column:default_config;type:text"`
	IsActive      int    `gorm:"column:is_active;not null;default:0;index"`
	IsCustom      int    `gorm:"column:is_custom;not null;default:0"`
	IsPreset      int    `gorm:"column:is_preset;not null;default:1"`
	CreatedAt     int64  `gorm:"column:created_at;not null"`
	UpdatedAt     int64  `gorm:"column:updated_at;not null"`
}

func (themeModel) TableName() string {
	return "themes"
}

func (m *themeModel) toEntity() *storefront.Theme {
	var config storefront.ThemeConfig
	if m.Config != "" {
		json.Unmarshal([]byte(m.Config), &config)
	}

	var configSchema storefront.ThemeConfigSchema
	if m.ConfigSchema != "" {
		json.Unmarshal([]byte(m.ConfigSchema), &configSchema)
	}

	var defaultConfig storefront.ThemeConfig
	if m.DefaultConfig != "" {
		json.Unmarshal([]byte(m.DefaultConfig), &defaultConfig)
	}

	return &storefront.Theme{
		ID:            m.ID,
		TenantID:      shared.TenantID(m.TenantID),
		Name:          m.Name,
		Code:          m.Code,
		Description:   m.Description,
		Thumbnail:     m.Thumbnail,
		PreviewImage:  m.PreviewImage,
		Config:        config,
		ConfigSchema:  configSchema,
		DefaultConfig: defaultConfig,
		IsActive:      m.IsActive == 1,
		IsCustom:      m.IsCustom == 1,
		IsPreset:      m.IsPreset == 1,
	}
}

func fromThemeEntity(t *storefront.Theme) *themeModel {
	now := time.Now().Unix()
	config, _ := json.Marshal(t.Config)
	configSchema, _ := json.Marshal(t.ConfigSchema)
	defaultConfig, _ := json.Marshal(t.DefaultConfig)

	isActive := 0
	if t.IsActive {
		isActive = 1
	}
	isCustom := 0
	if t.IsCustom {
		isCustom = 1
	}
	isPreset := 0
	if t.IsPreset {
		isPreset = 1
	}

	return &themeModel{
		ID:            t.ID,
		TenantID:      t.TenantID.Int64(),
		Name:          t.Name,
		Code:          t.Code,
		Description:   t.Description,
		Thumbnail:     t.Thumbnail,
		PreviewImage:  t.PreviewImage,
		Config:        string(config),
		ConfigSchema:  string(configSchema),
		DefaultConfig: string(defaultConfig),
		IsActive:      isActive,
		IsCustom:      isCustom,
		IsPreset:      isPreset,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (r *themeRepo) Create(ctx context.Context, db *gorm.DB, theme *storefront.Theme) error {
	model := fromThemeEntity(theme)
	return db.WithContext(ctx).Create(model).Error
}

func (r *themeRepo) Update(ctx context.Context, db *gorm.DB, theme *storefront.Theme) error {
	model := fromThemeEntity(theme)
	config, _ := json.Marshal(theme.Config)
	configSchema, _ := json.Marshal(theme.ConfigSchema)
	defaultConfig, _ := json.Marshal(theme.DefaultConfig)

	// Only allow updating custom themes (not preset themes)
	return db.WithContext(ctx).Model(&themeModel{}).
		Where("id = ? AND tenant_id = ? AND is_preset = 0", theme.ID, theme.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":           model.Name,
			"description":    model.Description,
			"thumbnail":      model.Thumbnail,
			"preview_image":  model.PreviewImage,
			"config":         string(config),
			"config_schema":  string(configSchema),
			"default_config": string(defaultConfig),
			"is_active":      model.IsActive,
			"is_custom":      model.IsCustom,
			"updated_at":     model.UpdatedAt,
		}).Error
}

func (r *themeRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*storefront.Theme, error) {
	var model themeModel
	err := db.WithContext(ctx).
		Where("id = ? AND (tenant_id = ? OR is_preset = 1)", id, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *themeRepo) FindByCode(ctx context.Context, db *gorm.DB, code string) (*storefront.Theme, error) {
	var model themeModel
	err := db.WithContext(ctx).
		Where("code = ? AND is_preset = 1", code).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *themeRepo) FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*storefront.Theme, error) {
	var model themeModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = 1", tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *themeRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*storefront.Theme, error) {
	var models []themeModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? OR is_preset = 1", tenantID.Int64()).
		Order("is_preset ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	themes := make([]*storefront.Theme, len(models))
	for i, m := range models {
		themes[i] = m.toEntity()
	}
	return themes, nil
}

func (r *themeRepo) FindPresets(ctx context.Context, db *gorm.DB) ([]*storefront.Theme, error) {
	var models []themeModel
	err := db.WithContext(ctx).
		Where("is_preset = 1").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	themes := make([]*storefront.Theme, len(models))
	for i, m := range models {
		themes[i] = m.toEntity()
	}
	return themes, nil
}

func (r *themeRepo) SetActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deactivate all themes for this tenant
		if err := tx.Model(&themeModel{}).
			Where("tenant_id = ? AND is_custom = 1", tenantID.Int64()).
			Update("is_active", 0).Error; err != nil {
			return err
		}

		// Activate the specified theme
		return tx.Model(&themeModel{}).
			Where("id = ? AND (tenant_id = ? OR is_preset = 1)", id, tenantID.Int64()).
			Update("is_active", 1).Error
	})
}