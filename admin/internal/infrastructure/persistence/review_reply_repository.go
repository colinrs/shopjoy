package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type replyRepo struct{}

func NewReplyRepository() review.ReplyRepository {
	return &replyRepo{}
}

type replyModel struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	ReviewID  int64  `gorm:"column:review_id;not null;uniqueIndex"`
	TenantID  int64  `gorm:"column:tenant_id;not null"`
	AdminID   int64  `gorm:"column:admin_id;not null"`
	AdminName string `gorm:"column:admin_name;size:100;not null"`
	Content   string `gorm:"column:content;type:text;not null"`
	CreatedAt int64  `gorm:"column:created_at;not null"`
	UpdatedAt int64  `gorm:"column:updated_at;not null"`
}

func (replyModel) TableName() string {
	return "review_replies"
}

func (m *replyModel) toEntity() *review.ReviewReply {
	entity := &review.ReviewReply{
		ReviewID:  m.ReviewID,
		TenantID:  m.TenantID,
		AdminID:   m.AdminID,
		AdminName: m.AdminName,
		Content:   m.Content,
	}
	if m.CreatedAt > 0 {
		entity.CreatedAt = time.Unix(m.CreatedAt, 0)
	}
	if m.UpdatedAt > 0 {
		entity.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	}
	return entity
}

func fromReplyEntity(r *review.ReviewReply) *replyModel {
	return &replyModel{
		ReviewID:  r.ReviewID,
		TenantID:  r.TenantID,
		AdminID:   r.AdminID,
		AdminName: r.AdminName,
		Content:   r.Content,
	}
}

func (r *replyRepo) Create(ctx context.Context, db *gorm.DB, reply *review.ReviewReply) error {
	model := fromReplyEntity(reply)
	return db.WithContext(ctx).Create(model).Error
}

func (r *replyRepo) Update(ctx context.Context, db *gorm.DB, reply *review.ReviewReply) error {
	model := fromReplyEntity(reply)
	return db.WithContext(ctx).
		Model(&replyModel{}).
		Where("review_id = ?", reply.ReviewID).
		Updates(map[string]interface{}{
			"content":    model.Content,
			"updated_at": model.UpdatedAt,
		}).Error
}

func (r *replyRepo) Delete(ctx context.Context, db *gorm.DB, reviewID int64) error {
	return db.WithContext(ctx).
		Where("review_id = ?", reviewID).
		Delete(&replyModel{}).Error
}

func (r *replyRepo) FindByReviewID(ctx context.Context, db *gorm.DB, reviewID int64) (*review.ReviewReply, error) {
	var model replyModel
	err := db.WithContext(ctx).
		Where("review_id = ?", reviewID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrReplyNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *replyRepo) FindByReviewIDs(ctx context.Context, db *gorm.DB, reviewIDs []int64) ([]*review.ReviewReply, error) {
	if len(reviewIDs) == 0 {
		return []*review.ReviewReply{}, nil
	}

	var models []replyModel
	err := db.WithContext(ctx).
		Where("review_id IN ?", reviewIDs).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	replies := make([]*review.ReviewReply, len(models))
	for i, m := range models {
		replies[i] = m.toEntity()
	}
	return replies, nil
}