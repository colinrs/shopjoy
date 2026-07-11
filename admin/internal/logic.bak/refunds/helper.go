package refunds

import (
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

// mapRefundStatus converts int8 status to RefundStatus
func mapRefundStatus(status int8) fulfillment.RefundStatus {
	return fulfillment.RefundStatus(status)
}

// convertRefundToDetailResp converts RefundDetailResponse to RefundDetailResp
func convertRefundToDetailResp(r *appfulfillment.RefundDetailResponse) *types.RefundDetailResp {
	return &types.RefundDetailResp{
		ID:           r.ID,
		RefundNo:     r.RefundNo,
		OrderID:      r.OrderID,
		UserID:       r.UserID,
		Type:         fulfillment.RefundType(r.Type).String(),
		TypeText:     r.TypeText,
		Status:       fulfillment.RefundStatus(r.Status).String(),
		StatusText:   r.StatusText,
		ReasonType:   r.ReasonType,
		Reason:       r.Reason,
		Description:  r.Description,
		Images:       r.Images,
		Amount:       r.Amount,
		Currency:     r.Currency,
		RejectReason: r.RejectReason,
		ApprovedAt:   r.ApprovedAt,
		ApprovedBy:   r.ApprovedBy,
		CompletedAt:  r.CompletedAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

// convertRefundReasonToResp converts RefundReasonItemResponse to RefundReasonResp
func convertRefundReasonToResp(r *appfulfillment.RefundReasonItemResponse) *types.RefundReasonResp {
	return &types.RefundReasonResp{
		ID:       r.ID,
		Code:     r.Code,
		Name:     r.Name,
		Sort:     r.Sort,
		IsActive: r.IsActive,
	}
}

// parseTime parses RFC3339 time string
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, timeStr)
}
