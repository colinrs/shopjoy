package promotion

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	pkgcoupon "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// CreateCouponRequest 创建优惠券请求
type CreateCouponRequest struct {
	Name         string
	Code         string
	Description  string
	Type         pkgcoupon.CouponType
	Value        decimal.Decimal
	MinAmount    decimal.Decimal
	MaxDiscount  decimal.Decimal
	TotalCount   int
	PerUserLimit int
	StartAt      time.Time
	EndAt        time.Time
}

// UpdateCouponRequest 更新优惠券请求
type UpdateCouponRequest struct {
	ID           int64
	Name         string
	Description  string
	MinAmount    decimal.Decimal
	MaxDiscount  decimal.Decimal
	TotalCount   int
	PerUserLimit int
	StartAt      time.Time
	EndAt        time.Time
}

// CouponResponse 优惠券响应
type CouponResponse struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	Code         string          `json:"code"`
	Description  string          `json:"description"`
	Type         int             `json:"type"`
	Value        decimal.Decimal `json:"value"`
	MinAmount    decimal.Decimal `json:"min_amount"`
	MaxDiscount  decimal.Decimal `json:"max_discount"`
	TotalCount   int             `json:"total_count"`
	UsedCount    int             `json:"used_count"`
	PerUserLimit int             `json:"per_user_limit"`
	Status       int             `json:"status"`
	StartAt      string          `json:"start_at"`
	EndAt        string          `json:"end_at"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// CouponListResponse 优惠券列表响应
type CouponListResponse struct {
	List     []*CouponResponse `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// QueryCouponRequest 查询优惠券请求
type QueryCouponRequest struct {
	Name     string
	Status   pkgcoupon.CouponStatus
	Type     pkgcoupon.CouponType
	Page     int
	PageSize int
}

// GenerateCouponCodesRequest 生成优惠券码请求
type GenerateCouponCodesRequest struct {
	CouponID int64
	Count    int
}

// GenerateCouponCodesResponse 生成优惠券码响应
type GenerateCouponCodesResponse struct {
	Codes []string `json:"codes"`
}

// IssueCouponToUserRequest 向用户发放优惠券请求
type IssueCouponToUserRequest struct {
	CouponID int64
	UserID   int64
}

// IssueCouponToUserResponse 向用户发放优惠券响应
type IssueCouponToUserResponse struct {
	UserCouponID int64  `json:"user_coupon_id"`
	CouponCode   string `json:"coupon_code"`
	ExpireAt     string `json:"expire_at"`
}

// UserCouponResponse 用户优惠券响应
type UserCouponResponse struct {
	ID         int64  `json:"id"`
	CouponID   int64  `json:"coupon_id"`
	CouponName string `json:"coupon_name"`
	CouponCode string `json:"coupon_code"`
	Status     int    `json:"status"`
	UsedAt     string `json:"used_at,omitempty"`
	OrderID    int64  `json:"order_id,omitempty"`
	ReceivedAt string `json:"received_at"`
	ExpireAt   string `json:"expire_at"`
}

// UserCouponListResponse 用户优惠券列表响应
type UserCouponListResponse struct {
	List     []*UserCouponResponse `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}

// CouponApp 优惠券应用服务接口
type CouponApp interface {
	CreateCoupon(ctx context.Context, tenantID shared.TenantID, req CreateCouponRequest) (*CouponResponse, error)
	UpdateCoupon(ctx context.Context, tenantID shared.TenantID, req UpdateCouponRequest) (*CouponResponse, error)
	GetCoupon(ctx context.Context, tenantID shared.TenantID, id int64) (*CouponResponse, error)
	ListCoupons(ctx context.Context, tenantID shared.TenantID, req QueryCouponRequest) (*CouponListResponse, error)
	DeleteCoupon(ctx context.Context, tenantID shared.TenantID, id int64) error
	GenerateCouponCodes(ctx context.Context, tenantID shared.TenantID, req GenerateCouponCodesRequest) (*GenerateCouponCodesResponse, error)
	IssueCouponToUser(ctx context.Context, tenantID shared.TenantID, req IssueCouponToUserRequest) (*IssueCouponToUserResponse, error)
	ListUserCoupons(ctx context.Context, tenantID shared.TenantID, userID int64, status pkgcoupon.UserCouponStatus, page, pageSize int) (*UserCouponListResponse, error)
}

type couponApp struct {
	db             *gorm.DB
	couponRepo     pkgcoupon.CouponRepository
	userCouponRepo pkgcoupon.UserCouponRepository
	idGen          snowflake.Snowflake
}

// NewCouponApp 创建优惠券应用服务
func NewCouponApp(
	db *gorm.DB,
	couponRepo pkgcoupon.CouponRepository,
	userCouponRepo pkgcoupon.UserCouponRepository,
	idGen snowflake.Snowflake,
) CouponApp {
	return &couponApp{
		db:             db,
		couponRepo:     couponRepo,
		userCouponRepo: userCouponRepo,
		idGen:          idGen,
	}
}

func (a *couponApp) CreateCoupon(ctx context.Context, tenantID shared.TenantID, req CreateCouponRequest) (*CouponResponse, error) {
	// Input validation
	if req.Name == "" {
		return nil, code.ErrCouponNameRequired
	}
	if !req.Type.IsValid() {
		return nil, code.ErrCouponTypeInvalid
	}
	if req.Value.IsZero() || req.Value.IsNegative() {
		return nil, code.ErrCouponValueRequired
	}
	if req.StartAt.IsZero() || req.EndAt.IsZero() {
		return nil, code.ErrCouponTimeRequired
	}
	if req.StartAt.After(req.EndAt) {
		return nil, code.ErrCouponInvalidTimeRange
	}

	id, err := a.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	c := &pkgcoupon.Coupon{
		ID:           id,
		TenantID:     tenantID,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Type:         req.Type,
		Value:        req.Value,
		MinAmount:    req.MinAmount,
		MaxDiscount:  req.MaxDiscount,
		TotalCount:   req.TotalCount,
		UsedCount:    0,
		PerUserLimit: req.PerUserLimit,
		Status:       pkgcoupon.CouponStatusInactive,
		StartAt:      req.StartAt.UTC(),
		EndAt:        req.EndAt.UTC(),
		Audit:        shared.NewAuditInfo(0),
	}

	// Generate code if not provided
	if c.Code == "" {
		c.Code = fmt.Sprintf("CPN%s%d", time.Now().Format("20060102"), c.ID)
	}

	if err := a.couponRepo.Create(ctx, a.db, c); err != nil {
		return nil, err
	}

	return toCouponResponse(c), nil
}

func (a *couponApp) UpdateCoupon(ctx context.Context, tenantID shared.TenantID, req UpdateCouponRequest) (*CouponResponse, error) {
	c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Only allow update if coupon is not active
	if c.Status == pkgcoupon.CouponStatusActive {
		return nil, code.ErrCouponCannotDelete
	}

	c.Name = req.Name
	c.Description = req.Description
	c.MinAmount = req.MinAmount
	c.MaxDiscount = req.MaxDiscount
	c.TotalCount = req.TotalCount
	c.PerUserLimit = req.PerUserLimit
	c.StartAt = req.StartAt.UTC()
	c.EndAt = req.EndAt.UTC()
	c.Audit.Update(0)

	if err := a.couponRepo.Update(ctx, a.db, c); err != nil {
		return nil, err
	}

	return toCouponResponse(c), nil
}

func (a *couponApp) GetCoupon(ctx context.Context, tenantID shared.TenantID, id int64) (*CouponResponse, error) {
	c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toCouponResponse(c), nil
}

func (a *couponApp) ListCoupons(ctx context.Context, tenantID shared.TenantID, req QueryCouponRequest) (*CouponListResponse, error) {
	query := pkgcoupon.CouponQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
	}
	if req.Status != 0 {
		query.Status = &req.Status
	}
	if req.Type != 0 {
		query.Type = &req.Type
	}
	query.PageQuery.Validate()

	coupons, total, err := a.couponRepo.FindList(ctx, a.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	resp := &CouponListResponse{
		List:     make([]*CouponResponse, len(coupons)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	for i, c := range coupons {
		resp.List[i] = toCouponResponse(c)
	}

	return resp, nil
}

func (a *couponApp) DeleteCoupon(ctx context.Context, tenantID shared.TenantID, id int64) error {
	c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return err
	}

	// Only allow delete if coupon is not active
	if c.Status == pkgcoupon.CouponStatusActive {
		return code.ErrCouponCannotDelete
	}

	return a.db.Transaction(func(tx *gorm.DB) error {
		// Soft delete coupon
		if err := a.couponRepo.Delete(ctx, tx, tenantID, id); err != nil {
			return err
		}
		return nil
	})
}

func (a *couponApp) GenerateCouponCodes(ctx context.Context, tenantID shared.TenantID, req GenerateCouponCodesRequest) (*GenerateCouponCodesResponse, error) {
	c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, req.CouponID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, req.Count)
	for i := 0; i < req.Count; i++ {
		code := fmt.Sprintf("%s-%s-%04d", c.Code, time.Now().Format("20060102"), i+1)
		codes = append(codes, code)
	}

	return &GenerateCouponCodesResponse{
		Codes: codes,
	}, nil
}

func (a *couponApp) IssueCouponToUser(ctx context.Context, tenantID shared.TenantID, req IssueCouponToUserRequest) (*IssueCouponToUserResponse, error) {
	// Find coupon
	c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, req.CouponID)
	if err != nil {
		return nil, err
	}

	// Check if coupon is active
	if !c.IsActive() {
		return nil, code.ErrCouponNotActive
	}

	// Check usage limit
	if c.TotalCount > 0 && c.UsedCount >= c.TotalCount {
		return nil, code.ErrCouponUserLimitReached
	}

	var userCouponID int64

	err = a.db.Transaction(func(tx *gorm.DB) error {
		// Generate user coupon ID
		id, err := a.idGen.NextID(ctx)
		if err != nil {
			return err
		}
		userCouponID = id

		// Create user coupon
		userCoupon := &pkgcoupon.UserCoupon{
			ID:         id,
			TenantID:   tenantID,
			UserID:     req.UserID,
			CouponID:   c.ID,
			Status:     pkgcoupon.UserCouponStatusUnused,
			ReceivedAt: time.Now().UTC(),
			ExpireAt:   c.EndAt,
		}

		if err := a.userCouponRepo.Create(ctx, tx, userCoupon); err != nil {
			return err
		}

		// Atomically increment coupon usage count (avoid race condition)
		if err := a.couponRepo.IncrementUsage(ctx, tx, tenantID, c.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &IssueCouponToUserResponse{
		UserCouponID: userCouponID,
		CouponCode:   c.Code,
		ExpireAt:     c.EndAt.Format(time.RFC3339),
	}, nil
}

func (a *couponApp) ListUserCoupons(ctx context.Context, tenantID shared.TenantID, userID int64, status pkgcoupon.UserCouponStatus, page, pageSize int) (*UserCouponListResponse, error) {
	var statusPtr *pkgcoupon.UserCouponStatus
	if status != 0 {
		statusPtr = &status
	}

	userCoupons, err := a.userCouponRepo.FindByUserID(ctx, a.db, tenantID, userID, statusPtr)
	if err != nil {
		return nil, err
	}

	// Paginate results
	total := int64(len(userCoupons))
	start := (page - 1) * pageSize
	if start > int(total) {
		start = int(total)
	}
	end := start + pageSize
	if end > int(total) {
		end = int(total)
	}

	pageUserCoupons := userCoupons[start:end]

	resp := &UserCouponListResponse{
		List:     make([]*UserCouponResponse, 0, len(pageUserCoupons)),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	// Load coupon details
	for _, uc := range pageUserCoupons {
		c, err := a.couponRepo.FindByID(ctx, a.db, tenantID, uc.CouponID)
		if err != nil {
			continue
		}

		userCouponResp := &UserCouponResponse{
			ID:         uc.ID,
			CouponID:   uc.CouponID,
			CouponName: c.Name,
			CouponCode: c.Code,
			Status:     int(uc.Status),
			ReceivedAt: uc.ReceivedAt.Format(time.RFC3339),
			ExpireAt:   uc.ExpireAt.Format(time.RFC3339),
		}

		if uc.UsedAt != nil {
			userCouponResp.UsedAt = uc.UsedAt.Format(time.RFC3339)
		}
		if uc.OrderID != 0 {
			userCouponResp.OrderID = uc.OrderID
		}

		resp.List = append(resp.List, userCouponResp)
	}

	return resp, nil
}

// toCouponResponse 转换为响应DTO
func toCouponResponse(c *pkgcoupon.Coupon) *CouponResponse {
	return &CouponResponse{
		ID:           c.ID,
		Name:         c.Name,
		Code:         c.Code,
		Description:  c.Description,
		Type:         int(c.Type),
		Value:        c.Value,
		MinAmount:    c.MinAmount,
		MaxDiscount:  c.MaxDiscount,
		TotalCount:   c.TotalCount,
		UsedCount:    c.UsedCount,
		PerUserLimit: c.PerUserLimit,
		Status:       int(c.Status),
		StartAt:      c.StartAt.Format(time.RFC3339),
		EndAt:        c.EndAt.Format(time.RFC3339),
		CreatedAt:    c.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    c.Audit.UpdatedAt.Format(time.RFC3339),
	}
}
