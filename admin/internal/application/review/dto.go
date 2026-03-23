package review

import (
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
)

// ListReviewsRequest list reviews request
type ListReviewsRequest struct {
	ProductID  int64
	Status     string
	RatingMin  *int
	RatingMax  *int
	HasImage   bool
	Keyword    string
	StartTime  *time.Time
	EndTime    *time.Time
	Page       int
	PageSize   int
}

// ReviewListItem review list item DTO
type ReviewListItem struct {
	ID            int64
	OrderID       int64
	ProductID     int64
	ProductName   string
	SKUCode       string
	UserName      string
	IsAnonymous   bool
	IsVerified    bool
	QualityRating int
	ValueRating   int
	OverallRating string
	Content       string
	Images        []string
	Status        string
	IsFeatured    bool
	HelpfulCount  int
	HasReply      bool
	CreatedAt     string
}

// ListReviewsResponse list reviews response
type ListReviewsResponse struct {
	List     []*ReviewListItem
	Total    int64
	Page     int
	PageSize int
}

// ReviewDetailDTO review detail DTO
type ReviewDetailDTO struct {
	ID            int64
	TenantID      int64
	OrderID       int64
	ProductID     int64
	ProductName   string
	SKUCode       string
	UserID        int64
	UserName      string
	IsAnonymous   bool
	IsVerified    bool
	QualityRating int
	ValueRating   int
	OverallRating string
	Content       string
	Images        []string
	Status        string
	IsFeatured    bool
	HelpfulCount  int
	CreatedAt     string
	UpdatedAt     string
	Reply         *ReplyDTO
}

// ReplyDTO reply DTO
type ReplyDTO struct {
	ID        int64
	Content   string
	AdminName string
	CreatedAt string
	UpdatedAt string
}

// CreateReplyRequest create reply request
type CreateReplyRequest struct {
	Content string
}

// UpdateReplyRequest update reply request
type UpdateReplyRequest struct {
	Content string
}

// BatchOperationResult batch operation result
type BatchOperationResult struct {
	SuccessCount int
	FailedCount  int
	Errors       []string
}

// FromDomainReview converts domain review to DTO
func FromDomainReview(r *review.Review) *ReviewDetailDTO {
	dto := &ReviewDetailDTO{
		ID:            r.ID,
		TenantID:      r.TenantID.Int64(),
		OrderID:       r.OrderID,
		ProductID:     r.ProductID,
		SKUCode:       r.SKUCode,
		UserID:        r.UserID,
		UserName:      r.UserName,
		IsAnonymous:   r.IsAnonymous,
		IsVerified:    r.IsVerified,
		QualityRating: r.QualityRating,
		ValueRating:   r.ValueRating,
		OverallRating: formatRating(r.OverallRating),
		Content:       r.Content,
		Images:        r.Images,
		Status:        r.Status.String(),
		IsFeatured:    r.IsFeatured,
		HelpfulCount:  r.HelpfulCount,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     r.UpdatedAt.Format(time.RFC3339),
	}

	if r.Reply != nil {
		dto.Reply = &ReplyDTO{
			ID:        r.Reply.ID,
			Content:   r.Reply.Content,
			AdminName: r.Reply.AdminName,
			CreatedAt: r.Reply.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.Reply.UpdatedAt.Format(time.RFC3339),
		}
	}

	return dto
}

func formatRating(r float64) string {
	return fmt.Sprintf("%.2f", r)
}