package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type carrierRepo struct{}

func NewCarrierRepository() fulfillment.CarrierRepository {
	return &carrierRepo{}
}

// carrierModel represents the database model for Carrier
type carrierModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:false"`
	Code        string    `gorm:"column:code;size:20;not null;uniqueIndex"`
	Name        string    `gorm:"column:name;size:100;not null"`
	TrackingURL string    `gorm:"column:tracking_url;size:255;not null;default:''"`
	IsActive    int       `gorm:"column:is_active;not null;default:1"`
	Sort        int       `gorm:"column:sort;not null;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
}

func (carrierModel) TableName() string {
	return "carriers"
}

func (m *carrierModel) toEntity() *fulfillment.Carrier {
	return &fulfillment.Carrier{
		Model:       application.Model{ID: m.ID, CreatedAt: m.CreatedAt.UTC()},
		Code:        m.Code,
		Name:        m.Name,
		TrackingURL: m.TrackingURL,
		IsActive:    m.IsActive == 1,
		Sort:        m.Sort,
	}
}

func fromCarrierEntity(c *fulfillment.Carrier) *carrierModel {
	isActive := 0
	if c.IsActive {
		isActive = 1
	}
	return &carrierModel{
		ID:          c.Model.ID,
		Code:        c.Code,
		Name:        c.Name,
		TrackingURL: c.TrackingURL,
		IsActive:    isActive,
		Sort:        c.Sort,
		CreatedAt:   c.Model.CreatedAt,
	}
}

// Create inserts a new carrier
func (r *carrierRepo) Create(ctx context.Context, db *gorm.DB, carrier *fulfillment.Carrier) error {
	model := fromCarrierEntity(carrier)
	return db.WithContext(ctx).Create(model).Error
}

// FindByID finds a carrier by ID
func (r *carrierRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*fulfillment.Carrier, error) {
	var model carrierModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrCarrierNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByCode finds a carrier by code
func (r *carrierRepo) FindByCode(ctx context.Context, db *gorm.DB, codeStr string) (*fulfillment.Carrier, error) {
	var model carrierModel
	err := db.WithContext(ctx).Where("code = ?", codeStr).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrCarrierNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindAll finds all carriers
func (r *carrierRepo) FindAll(ctx context.Context, db *gorm.DB) ([]*fulfillment.Carrier, error) {
	var models []carrierModel
	err := db.WithContext(ctx).Order("sort ASC, id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	carriers := make([]*fulfillment.Carrier, len(models))
	for i, m := range models {
		carriers[i] = m.toEntity()
	}
	return carriers, nil
}

// FindActive finds all active carriers
func (r *carrierRepo) FindActive(ctx context.Context, db *gorm.DB) ([]*fulfillment.Carrier, error) {
	var models []carrierModel
	err := db.WithContext(ctx).
		Where("is_active = 1").
		Order("sort ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	carriers := make([]*fulfillment.Carrier, len(models))
	for i, m := range models {
		carriers[i] = m.toEntity()
	}
	return carriers, nil
}

// Update updates a carrier
func (r *carrierRepo) Update(ctx context.Context, db *gorm.DB, carrier *fulfillment.Carrier) error {
	model := fromCarrierEntity(carrier)
	return db.WithContext(ctx).
		Model(&carrierModel{}).
		Where("id = ?", carrier.Model.ID).
		Updates(map[string]any{
			"name":         model.Name,
			"tracking_url": model.TrackingURL,
			"is_active":    model.IsActive,
			"sort":         model.Sort,
		}).Error
}

// FindByCodes finds carriers by multiple codes
// Returns a map keyed by carrier code
func (r *carrierRepo) FindByCodes(ctx context.Context, db *gorm.DB, codes []string) (map[string]*fulfillment.Carrier, error) {
	if len(codes) == 0 {
		return make(map[string]*fulfillment.Carrier), nil
	}

	var models []carrierModel
	err := db.WithContext(ctx).Where("code IN ?", codes).Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*fulfillment.Carrier)
	for _, m := range models {
		result[m.Code] = m.toEntity()
	}
	return result, nil
}
