package review

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainReview "github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

type Service interface {
	ListReviews(ctx context.Context, tenantID shared.TenantID, req ListReviewsRequest) (*ListReviewsResponse, error)
	GetReview(ctx context.Context, tenantID shared.TenantID, id int64) (*ReviewDetailDTO, error)
	ApproveReview(ctx context.Context, tenantID shared.TenantID, id int64) error
	HideReview(ctx context.Context, tenantID shared.TenantID, id int64, reason string) error
	ShowReview(ctx context.Context, tenantID shared.TenantID, id int64) error
	DeleteReview(ctx context.Context, tenantID shared.TenantID, id int64) error
	ToggleFeatured(ctx context.Context, tenantID shared.TenantID, id int64, featured bool) error
	CreateReply(ctx context.Context, tenantID shared.TenantID, adminID int64, adminName string, reviewID int64, req CreateReplyRequest) (*ReplyDTO, error)
	UpdateReply(ctx context.Context, tenantID shared.TenantID, reviewID int64, req UpdateReplyRequest) (*ReplyDTO, error)
	DeleteReply(ctx context.Context, tenantID shared.TenantID, reviewID int64) error
	BatchApprove(ctx context.Context, tenantID shared.TenantID, ids []int64) (*BatchOperationResult, error)
	BatchHide(ctx context.Context, tenantID shared.TenantID, ids []int64, reason string) (*BatchOperationResult, error)
	GetStats(ctx context.Context, tenantID shared.TenantID) (*domainReview.OverallStats, error)
	GetProductStats(ctx context.Context, tenantID shared.TenantID, productID int64) (*domainReview.ReviewStats, error)
}

type service struct {
	db         *gorm.DB
	reviewRepo domainReview.Repository
	replyRepo  domainReview.ReplyRepository
	idGen      snowflake.Snowflake
}

func NewService(db *gorm.DB, reviewRepo domainReview.Repository, replyRepo domainReview.ReplyRepository, idGen snowflake.Snowflake) Service {
	return &service{
		db:         db,
		reviewRepo: reviewRepo,
		replyRepo:  replyRepo,
		idGen:      idGen,
	}
}

func (s *service) ListReviews(ctx context.Context, tenantID shared.TenantID, req ListReviewsRequest) (*ListReviewsResponse, error) {
	query := domainReview.Query{
		TenantID:  tenantID,
		ProductID: req.ProductID,
		UserID:    req.UserID,
		HasImage:  req.HasImage,
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	if req.Status != "" {
		status := domainReview.ParseStatus(req.Status)
		query.Status = status
	}
	if req.RatingMin != nil {
		query.RatingMin = req.RatingMin
	}
	if req.RatingMax != nil {
		query.RatingMax = req.RatingMax
	}
	if req.StartTime != nil {
		query.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		query.EndTime = req.EndTime
	}

	reviews, total, err := s.reviewRepo.FindList(ctx, s.db, query)
	if err != nil {
		return nil, err
	}

	// Get reply status for all reviews
	reviewIDs := make([]int64, len(reviews))
	for i, r := range reviews {
		reviewIDs[i] = int64(r.ID)
	}

	replies, err := s.replyRepo.FindByReviewIDs(ctx, s.db, reviewIDs)
	if err != nil {
		return nil, err
	}

	replyMap := make(map[int64]bool)
	for _, reply := range replies {
		replyMap[reply.ReviewID] = true
	}

	list := make([]*ReviewListItem, len(reviews))
	for i, r := range reviews {
		list[i] = &ReviewListItem{
			ID:            int64(r.ID),
			OrderID:       r.OrderID,
			ProductID:     r.ProductID,
			SKUCode:       r.SKUCode,
			UserName:      r.DisplayName(),
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
			HasReply:      replyMap[int64(r.ID)],
			CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		}
	}

	return &ListReviewsResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *service) GetReview(ctx context.Context, tenantID shared.TenantID, id int64) (*ReviewDetailDTO, error) {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Load reply
	reply, err := s.replyRepo.FindByReviewID(ctx, s.db, id)
	if err != nil && !errors.Is(err, code.ErrReplyNotFound) {
		return nil, err
	}
	rev.Reply = reply

	return FromDomainReview(rev), nil
}

func (s *service) ApproveReview(ctx context.Context, tenantID shared.TenantID, id int64) error {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rev.Approve(); err != nil {
		return err
	}

	return s.reviewRepo.Update(ctx, s.db, rev)
}

func (s *service) HideReview(ctx context.Context, tenantID shared.TenantID, id int64, reason string) error {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rev.Hide(reason); err != nil {
		return err
	}

	return s.reviewRepo.Update(ctx, s.db, rev)
}

func (s *service) ShowReview(ctx context.Context, tenantID shared.TenantID, id int64) error {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rev.Show(); err != nil {
		return err
	}

	return s.reviewRepo.Update(ctx, s.db, rev)
}

func (s *service) DeleteReview(ctx context.Context, tenantID shared.TenantID, id int64) error {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rev.SoftDelete(); err != nil {
		return err
	}

	return s.reviewRepo.Delete(ctx, s.db, tenantID, int64(rev.ID))
}

func (s *service) ToggleFeatured(ctx context.Context, tenantID shared.TenantID, id int64, featured bool) error {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rev.SetFeatured(featured); err != nil {
		return err
	}

	return s.reviewRepo.Update(ctx, s.db, rev)
}

func (s *service) CreateReply(ctx context.Context, tenantID shared.TenantID, adminID int64, adminName string, reviewID int64, req CreateReplyRequest) (*ReplyDTO, error) {
	rev, err := s.reviewRepo.FindByID(ctx, s.db, tenantID, reviewID)
	if err != nil {
		return nil, err
	}

	if !rev.CanReply() {
		return nil, code.ErrReviewCannotReplyHidden
	}

	// Check if already has reply
	_, err = s.replyRepo.FindByReviewID(ctx, s.db, reviewID)
	if err == nil {
		return nil, code.ErrReviewAlreadyReplied
	}
	if !errors.Is(err, code.ErrReplyNotFound) {
		return nil, err
	}

	reply, err := domainReview.NewReviewReply(reviewID, tenantID.Int64(), adminID, adminName, req.Content)
	if err != nil {
		return nil, err
	}

	if err := s.replyRepo.Create(ctx, s.db, reply); err != nil {
		return nil, err
	}

	return &ReplyDTO{
		ID:        int64(reply.ID),
		Content:   reply.Content,
		AdminName: reply.AdminName,
		CreatedAt: reply.CreatedAt.Format(time.RFC3339),
		UpdatedAt: reply.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) UpdateReply(ctx context.Context, tenantID shared.TenantID, reviewID int64, req UpdateReplyRequest) (*ReplyDTO, error) {
	reply, err := s.replyRepo.FindByReviewID(ctx, s.db, reviewID)
	if err != nil {
		return nil, err
	}

	if err := reply.Update(req.Content); err != nil {
		return nil, err
	}

	if err := s.replyRepo.Update(ctx, s.db, reply); err != nil {
		return nil, err
	}

	return &ReplyDTO{
		ID:        int64(reply.ID),
		Content:   reply.Content,
		AdminName: reply.AdminName,
		CreatedAt: reply.CreatedAt.Format(time.RFC3339),
		UpdatedAt: reply.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) DeleteReply(ctx context.Context, tenantID shared.TenantID, reviewID int64) error {
	return s.replyRepo.Delete(ctx, s.db, reviewID)
}

func (s *service) BatchApprove(ctx context.Context, tenantID shared.TenantID, ids []int64) (*BatchOperationResult, error) {
	if len(ids) == 0 {
		return nil, code.ErrReviewBatchEmpty
	}
	if len(ids) > 100 {
		return nil, code.ErrReviewBatchLimitExceeded
	}

	count, err := s.reviewRepo.BatchUpdateStatus(ctx, s.db, tenantID, ids, domainReview.StatusApproved, "")
	if err != nil {
		return nil, err
	}

	return &BatchOperationResult{
		SuccessCount: int(count),
		FailedCount:  len(ids) - int(count),
		Errors:       []string{},
	}, nil
}

func (s *service) BatchHide(ctx context.Context, tenantID shared.TenantID, ids []int64, reason string) (*BatchOperationResult, error) {
	if len(ids) == 0 {
		return nil, code.ErrReviewBatchEmpty
	}
	if len(ids) > 100 {
		return nil, code.ErrReviewBatchLimitExceeded
	}

	count, err := s.reviewRepo.BatchUpdateStatus(ctx, s.db, tenantID, ids, domainReview.StatusHidden, reason)
	if err != nil {
		return nil, err
	}

	return &BatchOperationResult{
		SuccessCount: int(count),
		FailedCount:  len(ids) - int(count),
		Errors:       []string{},
	}, nil
}

type reviewStatsRow struct {
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
	FeaturedCount    int64
	RepliedCount     int64
}

func (s *service) GetStats(ctx context.Context, tenantID shared.TenantID) (*domainReview.OverallStats, error) {
	var row reviewStatsRow

	mainWhere := "deleted_at IS NULL"
	subWhere := "r2.deleted_at IS NULL"
	var args []interface{}
	if tenantID != 0 {
		mainWhere += " AND tenant_id = ?"
		subWhere += " AND r2.tenant_id = ?"
		args = append(args, tenantID.Int64(), tenantID.Int64())
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(*)                                                        AS total_reviews,
			SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END)                    AS pending_reviews,
			SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END)                    AS approved_reviews,
			SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END)                    AS hidden_reviews,
			COALESCE(AVG(overall_rating), 0)                                AS average_rating,
			COALESCE(AVG(quality_rating), 0)                                AS quality_avg_rating,
			COALESCE(AVG(value_rating), 0)                                  AS value_avg_rating,
			SUM(CASE WHEN overall_rating >= 4.5 THEN 1 ELSE 0 END)         AS five_star_count,
			SUM(CASE WHEN overall_rating >= 3.5 AND overall_rating < 4.5 THEN 1 ELSE 0 END) AS four_star_count,
			SUM(CASE WHEN overall_rating >= 2.5 AND overall_rating < 3.5 THEN 1 ELSE 0 END) AS three_star_count,
			SUM(CASE WHEN overall_rating >= 1.5 AND overall_rating < 2.5 THEN 1 ELSE 0 END) AS two_star_count,
			SUM(CASE WHEN overall_rating < 1.5 THEN 1 ELSE 0 END)          AS one_star_count,
			SUM(CASE WHEN images IS NOT NULL AND images != '[]' AND images != '' THEN 1 ELSE 0 END) AS with_image_count,
			SUM(CASE WHEN is_featured = true THEN 1 ELSE 0 END)            AS featured_count,
			(SELECT COUNT(DISTINCT rr.review_id) FROM review_replies rr
			 INNER JOIN reviews r2 ON rr.review_id = r2.id
			 WHERE %s)                                                      AS replied_count
		FROM reviews
		WHERE %s
	`, subWhere, mainWhere)

	err := s.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	var replyRate float64
	if row.TotalReviews > 0 {
		replyRate = float64(row.RepliedCount) / float64(row.TotalReviews) * 100
	}

	return &domainReview.OverallStats{
		TotalReviews:     row.TotalReviews,
		PendingReviews:   row.PendingReviews,
		ApprovedReviews:  row.ApprovedReviews,
		HiddenReviews:    row.HiddenReviews,
		AverageRating:    row.AverageRating,
		QualityAvgRating: row.QualityAvgRating,
		ValueAvgRating:   row.ValueAvgRating,
		FiveStarCount:    row.FiveStarCount,
		FourStarCount:    row.FourStarCount,
		ThreeStarCount:   row.ThreeStarCount,
		TwoStarCount:     row.TwoStarCount,
		OneStarCount:     row.OneStarCount,
		WithImageCount:   row.WithImageCount,
		ReplyRate:        replyRate,
		FeaturedCount:    row.FeaturedCount,
	}, nil
}

func (s *service) GetProductStats(ctx context.Context, tenantID shared.TenantID, productID int64) (*domainReview.ReviewStats, error) {
	// Implement product stats calculation
	return &domainReview.ReviewStats{}, nil
}
