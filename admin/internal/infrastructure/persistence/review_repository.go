package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type reviewRepo struct{}

func NewReviewRepository() review.Repository {
	return &reviewRepo{}
}

type reviewModel struct {
	ID            int64      `gorm:"column:id;primaryKey"`
	TenantID      int64      `gorm:"column:tenant_id;not null;index:idx_tenant_product"`
	OrderID       int64      `gorm:"column:order_id;not null;index"`
	ProductID     int64      `gorm:"column:product_id;not null;index:idx_tenant_product;index:idx_product_status"`
	SKUCode       string     `gorm:"column:sku_code;size:64;not null"`
	UserID        int64      `gorm:"column:user_id;not null;index:idx_tenant_user"`
	UserName      string     `gorm:"column:user_name;size:100;not null"`
	QualityRating int        `gorm:"column:quality_rating;not null"`
	ValueRating   int        `gorm:"column:value_rating;not null"`
	OverallRating float64    `gorm:"column:overall_rating;type:decimal(3,2);not null"`
	Content       string     `gorm:"column:content;type:text;not null"`
	Images        string     `gorm:"column:images;type:json"`
	Status        int        `gorm:"column:status;not null;default:0;index:idx_status;index:idx_product_status"`
	IsAnonymous   bool       `gorm:"column:is_anonymous;not null;default:false"`
	IsVerified    bool       `gorm:"column:is_verified;not null;default:false"`
	IsFeatured    bool       `gorm:"column:is_featured;not null;default:false"`
	HelpfulCount  int        `gorm:"column:helpful_count;not null;default:0"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null"`
}

func (reviewModel) TableName() string {
	return "reviews"
}

func (m *reviewModel) toEntity() *review.Review {
	var images []string
	if m.Images != "" {
		_ = json.Unmarshal([]byte(m.Images), &images)
	}

	return &review.Review{
		Model:         application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt},
		TenantID:      shared.TenantID(m.TenantID),
		OrderID:       m.OrderID,
		ProductID:     m.ProductID,
		SKUCode:       m.SKUCode,
		UserID:        m.UserID,
		UserName:      m.UserName,
		QualityRating: m.QualityRating,
		ValueRating:   m.ValueRating,
		OverallRating: m.OverallRating,
		Content:       m.Content,
		Images:        images,
		Status:        review.Status(m.Status),
		IsAnonymous:   m.IsAnonymous,
		IsVerified:    m.IsVerified,
		IsFeatured:    m.IsFeatured,
		HelpfulCount:  m.HelpfulCount,
	}
}

func fromReviewEntity(r *review.Review) *reviewModel {
	imagesJSON, _ := json.Marshal(r.Images)
	return &reviewModel{
		TenantID:      r.TenantID.Int64(),
		OrderID:       r.OrderID,
		ProductID:     r.ProductID,
		SKUCode:       r.SKUCode,
		UserID:        r.UserID,
		UserName:      r.UserName,
		QualityRating: r.QualityRating,
		ValueRating:   r.ValueRating,
		OverallRating: r.OverallRating,
		Content:       r.Content,
		Images:        string(imagesJSON),
		Status:        int(r.Status),
		IsAnonymous:   r.IsAnonymous,
		IsVerified:    r.IsVerified,
		IsFeatured:    r.IsFeatured,
		HelpfulCount:  r.HelpfulCount,
	}
}

func (r *reviewRepo) Create(ctx context.Context, db *gorm.DB, rev *review.Review) error {
	model := fromReviewEntity(rev)
	return db.WithContext(ctx).Create(model).Error
}

func (r *reviewRepo) Update(ctx context.Context, db *gorm.DB, rev *review.Review) error {
	model := fromReviewEntity(rev)
	return db.WithContext(ctx).
		Model(&reviewModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rev.ID, rev.TenantID.Int64()).
		Updates(map[string]interface{}{
			"status":        model.Status,
			"is_featured":   model.IsFeatured,
			"helpful_count": model.HelpfulCount,
			"updated_at":    model.UpdatedAt,
		}).Error
}

func (r *reviewRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&reviewModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrReviewNotFound
	}
	return nil
}

func (r *reviewRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*review.Review, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model reviewModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrReviewNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *reviewRepo) FindByIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64) ([]*review.Review, error) {
	if len(ids) == 0 {
		return []*review.Review{}, nil
	}

	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var models []reviewModel
	err := query.Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, err
	}

	reviews := make([]*review.Review, len(models))
	for i, m := range models {
		reviews[i] = m.toEntity()
	}
	return reviews, nil
}

func (r *reviewRepo) FindList(ctx context.Context, db *gorm.DB, query review.Query) ([]*review.Review, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&reviewModel{}).Where("deleted_at IS NULL")

	if query.TenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", query.TenantID.Int64())
	}

	if query.ProductID > 0 {
		dbQuery = dbQuery.Where("product_id = ?", query.ProductID)
	}
	if query.Status != nil && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}
	if query.RatingMin != nil {
		dbQuery = dbQuery.Where("overall_rating >= ?", float64(*query.RatingMin))
	}
	if query.RatingMax != nil {
		dbQuery = dbQuery.Where("overall_rating <= ?", float64(*query.RatingMax))
	}
	if query.HasImage {
		dbQuery = dbQuery.Where("images IS NOT NULL AND images != '[]' AND images != ''")
	}
	if query.Keyword != "" {
		dbQuery = dbQuery.Where("content LIKE ?", fmt.Sprintf("%%%s%%", query.Keyword))
	}
	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime.Unix())
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime.Unix())
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []reviewModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	reviews := make([]*review.Review, len(models))
	for i, m := range models {
		reviews[i] = m.toEntity()
	}
	return reviews, total, nil
}

func (r *reviewRepo) BatchUpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64, status review.Status, reason string) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	query := db.WithContext(ctx).Model(&reviewModel{}).
		Where("id IN ?", ids).
		Where("deleted_at IS NULL")

	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	// Filter by valid status transitions
	switch status {
	case review.StatusApproved:
		query = query.Where("status = ?", review.StatusPending)
	case review.StatusHidden:
		query = query.Where("status IN ?", []int{int(review.StatusPending), int(review.StatusApproved)})
	}

	now := time.Now().UTC()
	result := query.Updates(map[string]interface{}{
		"status":     int(status),
		"updated_at": now,
	})
	return result.RowsAffected, result.Error
}

func (r *reviewRepo) FindByProductID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64) ([]*review.Review, error) {
	dbQuery := db.WithContext(ctx).Model(&reviewModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	dbQuery = dbQuery.Where("product_id = ?", productID)

	var models []reviewModel
	err := dbQuery.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	reviews := make([]*review.Review, len(models))
	for i, m := range models {
		reviews[i] = m.toEntity()
	}
	return reviews, nil
}
