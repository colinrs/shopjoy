package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type shopRepo struct{}

func NewShopRepository() storefront.ShopRepository {
	return &shopRepo{}
}

type shopModel struct {
	ID              int64  `gorm:"column:id;primaryKey"`
	TenantID        int64  `gorm:"column:tenant_id;not null;uniqueIndex"`
	Name            string `gorm:"column:name;type:varchar(255);not null"`
	Description     string `gorm:"column:description;type:text"`
	Logo            string `gorm:"column:logo;type:varchar(500)"`
	Banner          string `gorm:"column:banner;type:varchar(500)"`
	ContactPhone    string `gorm:"column:contact_phone;type:varchar(20)"`
	ContactEmail    string `gorm:"column:contact_email;type:varchar(255)"`
	Address         string `gorm:"column:address;type:text"`
	SocialLinks     string `gorm:"column:social_links;type:json"`
	SEOTitle        string `gorm:"column:seo_title;type:varchar(255)"`
	SEODescription  string `gorm:"column:seo_description;type:text"`
	SEOKeywords     string `gorm:"column:seo_keywords;type:varchar(500)"`
	Status          int8   `gorm:"column:status;not null;default:1;index"`
	CurrentThemeID  *int64 `gorm:"column:current_theme_id"`
	ThemeConfig     string `gorm:"column:theme_config;type:text"`
	CreatedAt       int64  `gorm:"column:created_at;not null"`
	UpdatedAt       int64  `gorm:"column:updated_at;not null"`
}

func (shopModel) TableName() string {
	return "shops"
}

func (m *shopModel) toEntity() *storefront.Shop {
	var socialLinks map[string]string
	if m.SocialLinks != "" {
		json.Unmarshal([]byte(m.SocialLinks), &socialLinks)
	}

	var themeConfig *storefront.ThemeConfig
	if m.ThemeConfig != "" {
		var config storefront.ThemeConfig
		if err := json.Unmarshal([]byte(m.ThemeConfig), &config); err == nil {
			themeConfig = &config
		}
	}

	return &storefront.Shop{
		Model: application.Model{ID: m.ID, CreatedAt: time.Unix(m.CreatedAt, 0).UTC(), UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC()},
		TenantID:       shared.TenantID(m.TenantID),
		Name:           m.Name,
		Description:    m.Description,
		Logo:           m.Logo,
		Banner:         m.Banner,
		ContactPhone:   m.ContactPhone,
		ContactEmail:   m.ContactEmail,
		Address:        m.Address,
		SocialLinks:    socialLinks,
		SEO:            storefront.SEOConfig{
			Title:       m.SEOTitle,
			Description: m.SEODescription,
		},
		Status:         shared.Status(m.Status),
		CurrentThemeID: m.CurrentThemeID,
		ThemeConfig:    themeConfig,
	}
}

func fromShopEntity(s *storefront.Shop) *shopModel {
	socialLinks, _ := json.Marshal(s.SocialLinks)
	var themeConfig string
	if s.ThemeConfig != nil {
		data, _ := json.Marshal(s.ThemeConfig)
		themeConfig = string(data)
	}

	return &shopModel{
		ID:             s.Model.ID,
		TenantID:       s.TenantID.Int64(),
		Name:           s.Name,
		Description:    s.Description,
		Logo:           s.Logo,
		Banner:         s.Banner,
		ContactPhone:   s.ContactPhone,
		ContactEmail:   s.ContactEmail,
		Address:        s.Address,
		SocialLinks:    string(socialLinks),
		SEOTitle:       s.SEO.Title,
		SEODescription: s.SEO.Description,
		SEOKeywords:    keywordsToString(s.SEO.Keywords),
		Status:         int8(s.Status),
		CurrentThemeID: s.CurrentThemeID,
		ThemeConfig:    themeConfig,
		CreatedAt:      s.Model.CreatedAt.Unix(),
		UpdatedAt:      s.Model.UpdatedAt.Unix(),
	}
}

func keywordsToString(keywords []string) string {
	if len(keywords) == 0 {
		return ""
	}
	data, _ := json.Marshal(keywords)
	return string(data)
}

func (r *shopRepo) FindByTenantID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*storefront.Shop, error) {
	var model shopModel
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *shopRepo) Save(ctx context.Context, db *gorm.DB, shop *storefront.Shop) error {
	model := fromShopEntity(shop)

	// Check if exists
	var existing shopModel
	err := db.WithContext(ctx).
		Where("tenant_id = ?", model.TenantID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	shop.Model.ID = existing.ID
	model.ID = existing.ID
	return db.WithContext(ctx).Model(&shopModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"name":             model.Name,
			"description":      model.Description,
			"logo":             model.Logo,
			"banner":           model.Banner,
			"contact_phone":    model.ContactPhone,
			"contact_email":    model.ContactEmail,
			"address":          model.Address,
			"social_links":     model.SocialLinks,
			"seo_title":        model.SEOTitle,
			"seo_description":  model.SEODescription,
			"seo_keywords":     model.SEOKeywords,
			"status":           model.Status,
			"current_theme_id": model.CurrentThemeID,
			"theme_config":     model.ThemeConfig,
			"updated_at":       model.UpdatedAt,
		}).Error
}