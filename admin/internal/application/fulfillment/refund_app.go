package fulfillment

import (
	"context"
	"encoding/json"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RefundDetailResponse represents the refund detail response DTO for API
type RefundDetailResponse struct {
	ID             int64    `json:"id"`
	RefundNo       string   `json:"refund_no"`
	OrderID        string   `json:"order_id"`
	UserID         int64    `json:"user_id"`
	Type           int8     `json:"type"`
	TypeText       string   `json:"type_text"`
	Status         int8     `json:"status"`
	StatusText     string   `json:"status_text"`
	ReasonType     string   `json:"reason_type"`
	Reason         string   `json:"reason"`
	Description    string   `json:"description"`
	Images         []string `json:"images"`
	Amount         string   `json:"amount"`
	Currency       string   `json:"currency"`
	RejectReason   string   `json:"reject_reason"`
	ApprovedAt     string   `json:"approved_at"`
	ApprovedBy     int64    `json:"approved_by"`
	CompletedAt    string   `json:"completed_at"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// RefundListResponse represents the refund list response
type RefundListResponse struct {
	List     []*RefundDetailResponse `json:"list"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

// QueryRefundRequest represents the refund query request
type QueryRefundRequest struct {
	Page       int
	PageSize   int
	RefundNo   string
	OrderID    string
	UserID     int64
	Status     fulfillment.RefundStatus
	ReasonType string
	StartTime  time.Time
	EndTime    time.Time
}

// RefundStatisticsResponse represents the refund statistics response
type RefundStatisticsResponse struct {
	TotalRefunds    int64                    `json:"total_refunds"`
	TotalAmount     string                   `json:"total_amount"`
	Currency        string                   `json:"currency"`
	RefundRate      string                   `json:"refund_rate"`
	PendingCount    int64                    `json:"pending_count"`
	ApprovedCount   int64                    `json:"approved_count"`
	RejectedCount   int64                    `json:"rejected_count"`
	CompletedCount  int64                    `json:"completed_count"`
	ReasonBreakdown []RefundReasonStatsResp  `json:"reason_breakdown"`
	DailyTrend      []RefundDailyStatsResp   `json:"daily_trend"`
	TopProducts     []RefundProductStatsResp `json:"top_products"`
}

// RefundReasonStatsResp represents refund reason statistics
type RefundReasonStatsResp struct {
	ReasonType string `json:"reason_type"`
	ReasonName string `json:"reason_name"`
	Count      int64  `json:"count"`
	Percentage string `json:"percentage"`
}

// RefundDailyStatsResp represents daily refund statistics
type RefundDailyStatsResp struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
	Amount string `json:"amount"`
}

// RefundProductStatsResp represents product refund statistics
type RefundProductStatsResp struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	RefundCount int64  `json:"refund_count"`
	RefundRate  string `json:"refund_rate"`
}

// RefundReasonListResponse represents the refund reason list response
type RefundReasonListResponse struct {
	List []*RefundReasonItemResponse `json:"list"`
}

// RefundReasonItemResponse represents a single refund reason
type RefundReasonItemResponse struct {
	ID       int64 `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
	IsActive bool   `json:"is_active"`
}

// RefundApp is the refund application service interface
type RefundApp interface {
	GetRefund(ctx context.Context, tenantID shared.TenantID, id int64) (*RefundDetailResponse, error)
	ListRefunds(ctx context.Context, tenantID shared.TenantID, req QueryRefundRequest) (*RefundListResponse, error)
	ApproveRefund(ctx context.Context, tenantID shared.TenantID, id int64, approvedBy int64) (*RefundDetailResponse, error)
	RejectRefund(ctx context.Context, tenantID shared.TenantID, id int64, rejectReason string, updatedBy int64) (*RefundDetailResponse, error)
	GetRefundStatistics(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time) (*RefundStatisticsResponse, error)
}

type refundApp struct {
	db               *gorm.DB
	refundRepo       fulfillment.RefundRepository
	refundReasonRepo fulfillment.RefundReasonRepository
	idGen            snowflake.Snowflake
}

// NewRefundApp creates a new refund application service
func NewRefundApp(
	db *gorm.DB,
	refundRepo fulfillment.RefundRepository,
	refundReasonRepo fulfillment.RefundReasonRepository,
	idGen snowflake.Snowflake,
) RefundApp {
	return &refundApp{
		db:               db,
		refundRepo:       refundRepo,
		refundReasonRepo: refundReasonRepo,
		idGen:            idGen,
	}
}

func (a *refundApp) GetRefund(ctx context.Context, tenantID shared.TenantID, id int64) (*RefundDetailResponse, error) {
	refund, err := a.refundRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toRefundDetailResponse(refund), nil
}

func (a *refundApp) ListRefunds(ctx context.Context, tenantID shared.TenantID, req QueryRefundRequest) (*RefundListResponse, error) {
	query := fulfillment.RefundQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		RefundNo:   req.RefundNo,
		OrderID:    req.OrderID,
		UserID:     req.UserID,
		ReasonType: req.ReasonType,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}
	if req.Status.IsValid() {
		query.Status = req.Status
	}
	query.PageQuery.Validate()

	refunds, total, err := a.refundRepo.FindList(ctx, a.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	resp := &RefundListResponse{
		List:     make([]*RefundDetailResponse, len(refunds)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	for i, r := range refunds {
		resp.List[i] = toRefundDetailResponse(r)
	}

	return resp, nil
}

func (a *refundApp) ApproveRefund(ctx context.Context, tenantID shared.TenantID, id int64, approvedBy int64) (*RefundDetailResponse, error) {
	refund, err := a.refundRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Validate status - only pending refunds can be approved
	if refund.Status != fulfillment.RefundStatusPending {
		return nil, code.ErrRefundInvalidStatus
	}

	// Approve the refund
	if err := refund.Approve(approvedBy); err != nil {
		return nil, err
	}

	// Update in database
	if err := a.refundRepo.Update(ctx, a.db, refund); err != nil {
		return nil, err
	}

	return toRefundDetailResponse(refund), nil
}

func (a *refundApp) RejectRefund(ctx context.Context, tenantID shared.TenantID, id int64, rejectReason string, updatedBy int64) (*RefundDetailResponse, error) {
	// Validate reject reason is provided
	if rejectReason == "" {
		return nil, code.ErrRefundRejectReasonRequired
	}

	refund, err := a.refundRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Validate status - only pending refunds can be rejected
	if refund.Status != fulfillment.RefundStatusPending {
		return nil, code.ErrRefundInvalidStatus
	}

	// Reject the refund
	if err := refund.Reject(rejectReason, updatedBy); err != nil {
		return nil, err
	}

	// Update in database
	if err := a.refundRepo.Update(ctx, a.db, refund); err != nil {
		return nil, err
	}

	return toRefundDetailResponse(refund), nil
}

func (a *refundApp) GetRefundStatistics(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time) (*RefundStatisticsResponse, error) {
	// Get counts by status
	pendingCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusPending)
	approvedCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusApproved)
	rejectedCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusRejected)
	completedCount, _ := a.refundRepo.CountByStatus(ctx, a.db, tenantID, fulfillment.RefundStatusCompleted)

	totalRefunds := pendingCount + approvedCount + rejectedCount + completedCount

	// Get refunds for the time period for amount calculation
	query := fulfillment.RefundQuery{
		StartTime: startTime,
		EndTime:   endTime,
	}
	query.Validate()

	refunds, _, err := a.refundRepo.FindList(ctx, a.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	var totalAmount decimal.Decimal
	reasonMap := make(map[string]int64)
	for _, r := range refunds {
		totalAmount = totalAmount.Add(r.Amount)
		reasonMap[r.ReasonType]++
	}

	// Get refund reasons for breakdown
	reasons, _ := a.refundReasonRepo.FindActive(ctx, a.db)
	reasonBreakdown := make([]RefundReasonStatsResp, 0, len(reasons))
	for _, reason := range reasons {
		count := reasonMap[reason.Code]
		var percentage string
		if totalRefunds > 0 {
			percentage = formatPercentage(float64(count) / float64(totalRefunds) * 100)
		} else {
			percentage = "0.00%"
		}
		reasonBreakdown = append(reasonBreakdown, RefundReasonStatsResp{
			ReasonType: reason.Code,
			ReasonName: reason.Name,
			Count:      count,
			Percentage: percentage,
		})
	}

	// Calculate refund rate (simplified - in real implementation, would compare with total orders)
	refundRate := "0.00%"
	if totalRefunds > 0 {
		refundRate = formatPercentage(float64(completedCount) / float64(totalRefunds) * 100)
	}

	return &RefundStatisticsResponse{
		TotalRefunds:    totalRefunds,
		TotalAmount:     formatAmount(totalAmount),
		Currency:        "CNY",
		RefundRate:      refundRate,
		PendingCount:    pendingCount,
		ApprovedCount:   approvedCount,
		RejectedCount:   rejectedCount,
		CompletedCount:  completedCount,
		ReasonBreakdown: reasonBreakdown,
		DailyTrend:      []RefundDailyStatsResp{},
		TopProducts:     []RefundProductStatsResp{},
	}, nil
}

// toRefundDetailResponse converts a refund entity to detail response DTO
func toRefundDetailResponse(r *fulfillment.Refund) *RefundDetailResponse {
	// Parse images from JSON
	var images []string
	if r.Images != "" {
		json.Unmarshal([]byte(r.Images), &images)
	}

	return &RefundDetailResponse{
		ID:           r.ID,
		RefundNo:     r.RefundNo,
		OrderID:      r.OrderID,
		UserID:       r.UserID,
		Type:         int8(r.Type),
		TypeText:     getRefundTypeText(r.Type),
		Status:       int8(r.Status),
		StatusText:   getRefundStatusText(r.Status),
		ReasonType:   r.ReasonType,
		Reason:       r.Reason,
		Description:  r.Description,
		Images:       images,
		Amount:       formatAmount(r.Amount),
		Currency:     r.Currency,
		RejectReason: r.RejectReason,
		ApprovedAt:   formatTimeToString(r.ApprovedAt),
		ApprovedBy:   r.ApprovedBy,
		CompletedAt:  formatTimeToString(r.CompletedAt),
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    r.UpdatedAt.Format(time.RFC3339),
	}
}

// getRefundTypeText returns the text representation of refund type
func getRefundTypeText(t fulfillment.RefundType) string {
	switch t {
	case fulfillment.RefundTypeFull:
		return "Full Refund"
	case fulfillment.RefundTypePartial:
		return "Partial Refund"
	default:
		return "Unknown"
	}
}

// getRefundStatusText returns the text representation of refund status
func getRefundStatusText(s fulfillment.RefundStatus) string {
	switch s {
	case fulfillment.RefundStatusPending:
		return "Pending"
	case fulfillment.RefundStatusApproved:
		return "Approved"
	case fulfillment.RefundStatusRejected:
		return "Rejected"
	case fulfillment.RefundStatusCompleted:
		return "Completed"
	case fulfillment.RefundStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// formatAmount formats amount (in yuan) to string
func formatAmount(amount decimal.Decimal) string {
	return shared.NewMoney(amount, "CNY").String()
}

// formatAmountFromInt64 formats amount (in yuan) to string
func formatAmountFromInt64(amount int64) string {
	return shared.NewMoney(decimal.NewFromInt(amount), "CNY").String()
}

// formatPercentage formats a percentage value
func formatPercentage(value float64) string {
	return shared.NewMoney(decimal.NewFromFloat(value*100), "").String()
}

// formatTimeToString formats a time.Time pointer to RFC3339 string
func formatTimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}