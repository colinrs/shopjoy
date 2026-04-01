// Package review 评价领域层
package review

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// Status 评价状态
type Status int

const (
	StatusPending  Status = iota // 待审核
	StatusApproved               // 已批准
	StatusHidden                 // 已隐藏
	StatusDeleted                // 已删除
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusApproved:
		return "approved"
	case StatusHidden:
		return "hidden"
	case StatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// ParseStatus parses status string to Status
func ParseStatus(s string) *Status {
	switch s {
	case "pending":
		status := StatusPending
		return &status
	case "approved":
		status := StatusApproved
		return &status
	case "hidden":
		status := StatusHidden
		return &status
	case "deleted":
		status := StatusDeleted
		return &status
	default:
		return nil
	}
}

// IsValid validates status
func (s Status) IsValid() bool {
	return s >= StatusPending && s <= StatusDeleted
}

// CanApprove checks if review can be approved
func (s Status) CanApprove() bool {
	return s == StatusPending
}

// CanHide checks if review can be hidden
func (s Status) CanHide() bool {
	return s == StatusPending || s == StatusApproved
}

// CanShow checks if hidden review can be shown
func (s Status) CanShow() bool {
	return s == StatusHidden
}

// Review 评价实体
type Review struct {
	application.Model
	TenantID      shared.TenantID // 租户ID
	OrderID       int64           // 订单ID
	ProductID     int64           // 商品ID
	SKUCode       string          // SKU代码
	UserID        int64           // 用户ID
	UserName      string          // 用户名（用于显示）
	QualityRating int             // 质量评分 (1-5)
	ValueRating   int             // 性价比评分 (1-5)
	OverallRating float64         // 综合评分
	Content       string          // 评价内容
	Images        []string        // 图片URL列表
	Status        Status          // 状态
	IsAnonymous   bool            // 是否匿名
	IsVerified    bool            // 是否已验证购买
	IsFeatured    bool            // 是否精选
	HelpfulCount  int             // 有帮助数
	Reply         *ReviewReply    // 商家回复
}

// TableName returns table name
func (Review) TableName() string {
	return "reviews"
}

// NewReview creates a new review
func NewReview(tenantID shared.TenantID, orderID, productID int64, skuCode string,
	userID int64, userName string, qualityRating, valueRating int, content string,
	images []string, isAnonymous bool) (*Review, error) {
	if qualityRating < 1 || qualityRating > 5 || valueRating < 1 || valueRating > 5 {
		return nil, code.ErrReviewInvalidRating
	}
	if len(content) > 1000 {
		return nil, code.ErrReviewContentTooLong
	}

	overallRating := float64(qualityRating+valueRating) / 2.0

	// Auto-approve high ratings (>= 4)
	status := StatusPending
	if overallRating >= 4.0 {
		status = StatusApproved
	}

	return &Review{
		TenantID:      tenantID,
		OrderID:       orderID,
		ProductID:     productID,
		SKUCode:       skuCode,
		UserID:        userID,
		UserName:      userName,
		QualityRating: qualityRating,
		ValueRating:   valueRating,
		OverallRating: overallRating,
		Content:       content,
		Images:        images,
		Status:        status,
		IsAnonymous:   isAnonymous,
		IsVerified:    true, // All reviews are verified purchases
		IsFeatured:    false,
		HelpfulCount:  0,
	}, nil
}

// Approve approves a pending review
func (r *Review) Approve() error {
	if !r.Status.CanApprove() {
		return code.ErrReviewCannotApprove
	}
	r.Status = StatusApproved
	return nil
}

// Hide hides a review
func (r *Review) Hide(reason string) error {
	if !r.Status.CanHide() {
		return code.ErrReviewCannotHide
	}
	r.Status = StatusHidden
	return nil
}

// Show shows a hidden review
func (r *Review) Show() error {
	if !r.Status.CanShow() {
		return code.ErrReviewCannotShow
	}
	r.Status = StatusApproved
	return nil
}

// SoftDelete soft deletes a review
func (r *Review) SoftDelete() error {
	if r.Status == StatusDeleted {
		return code.ErrReviewAlreadyDeleted
	}
	r.Status = StatusDeleted
	return nil
}

// SetFeatured sets featured status
func (r *Review) SetFeatured(featured bool) error {
	if featured && r.Status != StatusApproved {
		return code.ErrReviewCannotFeature
	}
	r.IsFeatured = featured
	return nil
}

// CanReply checks if review can be replied
func (r *Review) CanReply() bool {
	return r.Status != StatusHidden && r.Status != StatusDeleted
}

// DisplayName returns user display name (anonymous-aware)
func (r *Review) DisplayName() string {
	if r.IsAnonymous {
		return "Anonymous"
	}
	return r.UserName
}

// ReviewReply 商家回复实体
type ReviewReply struct {
	application.Model
	ReviewID  int64  // 评价ID
	TenantID  int64  // 租户ID
	AdminID   int64  // 管理员ID
	AdminName string // 管理员名称
	Content   string // 回复内容
}

// TableName returns table name
func (ReviewReply) TableName() string {
	return "review_replies"
}

// NewReviewReply creates a new review reply
func NewReviewReply(reviewID, tenantID, adminID int64, adminName, content string) (*ReviewReply, error) {
	if content == "" {
		return nil, code.ErrReviewReplyEmpty
	}
	if len(content) > 500 {
		return nil, code.ErrReplyContentTooLong
	}

	return &ReviewReply{
		ReviewID:  reviewID,
		TenantID:  tenantID,
		AdminID:   adminID,
		AdminName: adminName,
		Content:   content,
	}, nil
}

// Update updates reply content
func (r *ReviewReply) Update(content string) error {
	if content == "" {
		return code.ErrReviewReplyEmpty
	}
	if len(content) > 500 {
		return code.ErrReplyContentTooLong
	}
	r.Content = content
	return nil
}

// ReviewStats 商品评价统计
type ReviewStats struct {
	application.Model
	TenantID         int64   // 租户ID
	ProductID        int64   // 商品ID
	TotalReviews     int     // 总评价数
	AverageRating    float64 // 平均评分
	QualityAvgRating float64 // 质量平均分
	ValueAvgRating   float64 // 性价比平均分
	Rating1Count     int     // 1星数量
	Rating2Count     int     // 2星数量
	Rating3Count     int     // 3星数量
	Rating4Count     int     // 4星数量
	Rating5Count     int     // 5星数量
	WithImageCount   int     // 有图片数量
}

// TableName returns table name
func (ReviewStats) TableName() string {
	return "review_stats"
}

// StringArray for JSON array storage
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, s)
}

// Query 查询条件
type Query struct {
	TenantID  shared.TenantID
	ProductID int64
	Status    *Status
	RatingMin *int
	RatingMax *int
	HasImage  bool
	Keyword   string
	StartTime *time.Time
	EndTime   *time.Time
	Page      int
	PageSize  int
}

// Validate validates query
func (q *Query) Validate() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
}

// Offset returns offset
func (q Query) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// Limit returns limit
func (q Query) Limit() int {
	return q.PageSize
}

// Repository 评价仓储接口
type Repository interface {
	Create(ctx context.Context, db *gorm.DB, review *Review) error
	Update(ctx context.Context, db *gorm.DB, review *Review) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Review, error)
	FindByIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64) ([]*Review, error)
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Review, int64, error)
	FindByProductID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64) ([]*Review, error)
	BatchUpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64, status Status, reason string) (int64, error)
}

// ReplyRepository 回复仓储接口
type ReplyRepository interface {
	Create(ctx context.Context, db *gorm.DB, reply *ReviewReply) error
	Update(ctx context.Context, db *gorm.DB, reply *ReviewReply) error
	Delete(ctx context.Context, db *gorm.DB, reviewID int64) error
	FindByReviewID(ctx context.Context, db *gorm.DB, reviewID int64) (*ReviewReply, error)
	FindByReviewIDs(ctx context.Context, db *gorm.DB, reviewIDs []int64) ([]*ReviewReply, error)
}

// StatsRepository 统计仓储接口
type StatsRepository interface {
	GetByProduct(ctx context.Context, db *gorm.DB, tenantID, productID int64) (*ReviewStats, error)
	GetOverallStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*OverallStats, error)
	UpdateProductStats(ctx context.Context, db *gorm.DB, tenantID, productID int64) error
}

// OverallStats 总体统计
type OverallStats struct {
	TotalReviews     int64
	PendingReviews   int64
	ApprovedReviews  int64
	HiddenReviews    int64
	AverageRating    float64
	QualityAvgRating float64
	ValueAvgRating   float64
	FiveStarCount    int64
	FourStarCount    int64
	ThreeStarCount   int64
	TwoStarCount     int64
	OneStarCount     int64
	WithImageCount   int64
	ReplyRate        float64
	FeaturedCount    int64
}
