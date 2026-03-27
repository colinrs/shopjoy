package review

import (
	"context"
	"errors"
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

func (s *service) GetStats(ctx context.Context, tenantID shared.TenantID) (*domainReview.OverallStats, error) {
	// Implement stats calculation from database
	// This is a simplified implementation
	return &domainReview.OverallStats{}, nil
}

func (s *service) GetProductStats(ctx context.Context, tenantID shared.TenantID, productID int64) (*domainReview.ReviewStats, error) {
	// Implement product stats calculation
	return &domainReview.ReviewStats{}, nil
}