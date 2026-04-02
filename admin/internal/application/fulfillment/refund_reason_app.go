package fulfillment

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"gorm.io/gorm"
)

// RefundReasonApp is the refund reason application service interface
type RefundReasonApp interface {
	ListRefundReasons(ctx context.Context) (*RefundReasonListResponse, error)
}

type refundReasonApp struct {
	db               *gorm.DB
	refundReasonRepo fulfillment.RefundReasonRepository
}

// NewRefundReasonApp creates a new refund reason application service
func NewRefundReasonApp(db *gorm.DB, refundReasonRepo fulfillment.RefundReasonRepository) RefundReasonApp {
	return &refundReasonApp{
		db:               db,
		refundReasonRepo: refundReasonRepo,
	}
}

func (a *refundReasonApp) ListRefundReasons(ctx context.Context) (*RefundReasonListResponse, error) {
	reasons, err := a.refundReasonRepo.FindActive(ctx, a.db)
	if err != nil {
		return nil, err
	}

	list := make([]*RefundReasonItemResponse, len(reasons))
	for i, r := range reasons {
		list[i] = &RefundReasonItemResponse{
			ID:       r.ID,
			Code:     r.Code,
			Name:     r.Name,
			Sort:     r.Sort,
			IsActive: r.IsActive,
		}
	}

	return &RefundReasonListResponse{
		List: list,
	}, nil
}
