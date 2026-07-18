package promotion

import (
	"context"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// sanitizeTags enforces a max length of 64 chars per entry and
// drops empty entries. Returns nil for empty result so persistence
// can write SQL NULL consistently.
func sanitizeTags(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if len(t) > 64 {
			t = t[:64]
		}
		out = append(out, t)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// CreatePromotionRequest 创建促销请求
type CreatePromotionRequest struct {
	Name        string
	Description string
	Type        pkgpromotion.Type
	Priority    int
	Currency    string
	Scope       pkgpromotion.PromotionScope
	StartAt      time.Time
	EndAt        time.Time
	Rules        []CreatePromotionRuleRequest
	UsageLimit   int
	PerUserLimit int
	Tags         []string
}

// CreatePromotionRuleRequest 创建促销规则请求
type CreatePromotionRuleRequest struct {
	ConditionType  pkgpromotion.ConditionType
	ConditionValue decimal.Decimal
	ActionType     pkgpromotion.ActionType
	ActionValue    decimal.Decimal
	MaxDiscount    decimal.Decimal
}

// UpdatePromotionRequest 更新促销请求
//
// The previous version of this struct only carried Name / Description /
// StartAt / EndAt, so every other field on UpdatePromotionReq (Type,
// DiscountType, DiscountValue, MinOrderAmount, MaxDiscount,
// UsageLimit, PerUserLimit, ProductIDs, CategoryIDs, MarketIDs, Tags)
// was silently dropped between the wire and the persistence layer.
//
// Type and Scope are now honored. Discount mechanics still flow
// through the dedicated promotion-rules endpoints, not through this
// update, mirroring the Create-side split. UsageLimit, PerUserLimit
// and Tags have no DB columns today and remain no-ops until a
// schema migration lands; see [UpdatePromotion] for the storage map.
type UpdatePromotionRequest struct {
	ID          int64
	Name        string
	Description string
	Type        pkgpromotion.Type
	Scope       pkgpromotion.PromotionScope
	StartAt      time.Time
	EndAt        time.Time
	UsageLimit   int
	PerUserLimit int
	Tags         []string
}

// PromotionResponse 促销响应
type PromotionResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Status      int    `json:"status"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	// ScopeType mirrors the stored Scope.Type so the wire can carry
	// the same flag the form posts without forcing the frontend to
	// re-derive it from product_ids / category_ids arrays.
	ScopeType string                   `json:"scope_type"`
	Rules        []*PromotionRuleResponse `json:"rules"`
	UsageLimit   int                      `json:"usage_limit"`
	PerUserLimit int                      `json:"per_user_limit"`
	Tags         []string                 `json:"tags"`
	CreatedAt    string                   `json:"created_at"`
	UpdatedAt    string                   `json:"updated_at"`
}

// PromotionRuleResponse 促销规则响应
type PromotionRuleResponse struct {
	ID             int64           `json:"id"`
	PromotionID    int64           `json:"promotion_id"`
	ConditionType  int             `json:"condition_type"`
	ConditionValue decimal.Decimal `json:"condition_value"`
	ActionType     int             `json:"action_type"`
	ActionValue    decimal.Decimal `json:"action_value"`
	MaxDiscount    decimal.Decimal `json:"max_discount"`
}

// PromotionListResponse 促销列表响应
type PromotionListResponse struct {
	List     []*PromotionResponse `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// QueryPromotionRequest 查询促销请求
//
// Status and Type are pointers so that "no filter" (nil) can be
// distinguished from a filter that happens to equal the iota-zero
// value (e.g. StatusPending = 0, TypeDiscount = 0). The old
// `req.Status != 0` sentinel collided with those zero values and
// silently dropped the filter for the default status/type.
//
// ExpiredOnly is a separate flag because the wire status "expired"
// is computed from EndAt at response time and is never a stored
// enum value — it cannot be expressed as StatusExpired.
type QueryPromotionRequest struct {
	Name        string
	Status      *pkgpromotion.Status
	Type        *pkgpromotion.Type
	ExpiredOnly bool
	Page        int
	PageSize    int
}

// PromotionApp 促销应用服务接口
type PromotionApp interface {
	CreatePromotion(ctx context.Context, req CreatePromotionRequest) (*PromotionResponse, error)
	UpdatePromotion(ctx context.Context, req UpdatePromotionRequest) (*PromotionResponse, error)
	GetPromotion(ctx context.Context, id int64) (*PromotionResponse, error)
	ListPromotions(ctx context.Context, req QueryPromotionRequest) (*PromotionListResponse, error)
	DeletePromotion(ctx context.Context, id int64) error
	ActivatePromotion(ctx context.Context, id int64) error
	DeactivatePromotion(ctx context.Context, id int64) error
}

type promotionApp struct {
	db            *gorm.DB
	promotionRepo pkgpromotion.Repository
	idGen         snowflake.Snowflake
}

// NewPromotionApp 创建促销应用服务
func NewPromotionApp(db *gorm.DB, promotionRepo pkgpromotion.Repository, idGen snowflake.Snowflake) PromotionApp {
	return &promotionApp{
		db:            db,
		promotionRepo: promotionRepo,
		idGen:         idGen,
	}
}

func (a *promotionApp) CreatePromotion(ctx context.Context, req CreatePromotionRequest) (*PromotionResponse, error) {
	req.Tags = sanitizeTags(req.Tags)

	// Input validation
	if req.Name == "" {
		return nil, code.ErrPromotionNameRequired
	}
	if req.Currency == "" {
		return nil, code.ErrPromotionCurrencyRequired
	}
	if req.StartAt.IsZero() || req.EndAt.IsZero() {
		return nil, code.ErrPromotionTimeRequired
	}
	if req.StartAt.After(req.EndAt) {
		return nil, code.ErrPromotionInvalidTimeRange
	}
	if !req.Type.IsValid() {
		return nil, code.ErrPromotionTypeInvalid
	}
	// TODO(Task 5): replace code.ErrPromotionInvalid with code.ErrPromotionUsageLimitInvalid / code.ErrPromotionPerUserLimitInvalid
	if req.UsageLimit < 0 {
		return nil, code.ErrPromotionInvalid
	}
	if req.PerUserLimit < 0 {
		return nil, code.ErrPromotionInvalid
	}
	if !req.Scope.Type.IsValid() {
		return nil, code.ErrPromotionScopeInvalid
	}

	var result *pkgpromotion.Promotion

	err := a.db.Transaction(func(tx *gorm.DB) error {
		id, err := a.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		p := &pkgpromotion.Promotion{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Status:      pkgpromotion.StatusPending,
			Priority:    req.Priority,
			StartAt:     req.StartAt.UTC(),
			EndAt:       req.EndAt.UTC(),
			Scope:       req.Scope,
			UsageLimit:   req.UsageLimit,
			PerUserLimit: req.PerUserLimit,
			Tags:         req.Tags,
			Currency:    req.Currency,
			Audit:       shared.NewAuditInfo(0),
			Rules:       make([]pkgpromotion.PromotionRule, 0, len(req.Rules)),
		}

		// Create rules
		for _, ruleReq := range req.Rules {
			ruleID, err := a.idGen.NextID(ctx)
			if err != nil {
				return err
			}

			rule := pkgpromotion.PromotionRule{
				ID:             ruleID,
				PromotionID:    id,
				ConditionType:  ruleReq.ConditionType,
				ConditionValue: ruleReq.ConditionValue,
				ActionType:     ruleReq.ActionType,
				ActionValue:    ruleReq.ActionValue,
				MaxDiscount:    ruleReq.MaxDiscount,
			}
			p.Rules = append(p.Rules, rule)
		}

		if err := a.promotionRepo.Create(ctx, tx, p); err != nil {
			return err
		}

		result = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toPromotionResponse(result), nil
}

func (a *promotionApp) UpdatePromotion(ctx context.Context, req UpdatePromotionRequest) (*PromotionResponse, error) {
	req.Tags = sanitizeTags(req.Tags)

	p, err := a.promotionRepo.FindByID(ctx, a.db, req.ID)
	if err != nil {
		return nil, err
	}

	// Only allow update if promotion is not active
	if p.Status == pkgpromotion.StatusActive {
		return nil, code.ErrPromotionCannotDelete
	}

	// TODO(Task 5): replace code.ErrPromotionInvalid with code.ErrPromotionUsageLimitInvalid / code.ErrPromotionPerUserLimitInvalid
	if req.UsageLimit < 0 {
		return nil, code.ErrPromotionInvalid
	}
	if req.PerUserLimit < 0 {
		return nil, code.ErrPromotionInvalid
	}

	// Type change is allowed for non-active promotions. Changing the
	// promotion classification mid-flight is what makes the Update
	// path different from a status toggle; restart rules separately
	// via the promotion-rules endpoints if needed.
	if !req.Type.IsValid() {
		return nil, code.ErrPromotionTypeInvalid
	}
	p.Name = req.Name
	p.Description = req.Description
	p.Type = req.Type
	p.Scope = req.Scope
	p.UsageLimit = req.UsageLimit
	p.PerUserLimit = req.PerUserLimit
	p.Tags = req.Tags
	p.StartAt = req.StartAt.UTC()
	p.EndAt = req.EndAt.UTC()
	p.Audit.Update(0)

	if err := a.promotionRepo.Update(ctx, a.db, p); err != nil {
		return nil, err
	}

	return toPromotionResponse(p), nil
}

func (a *promotionApp) GetPromotion(ctx context.Context, id int64) (*PromotionResponse, error) {
	p, err := a.promotionRepo.FindByID(ctx, a.db, id)
	if err != nil {
		return nil, err
	}
	// Load promotion rules
	rules, err := a.promotionRepo.FindRulesByPromotionID(ctx, a.db, p.ID)
	if err != nil {
		return nil, err
	}
	p.Rules = rules
	return toPromotionResponse(p), nil
}

func (a *promotionApp) ListPromotions(ctx context.Context, req QueryPromotionRequest) (*PromotionListResponse, error) {
	query := pkgpromotion.Query{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name:        req.Name,
		ExpiredOnly: req.ExpiredOnly,
	}
	// Status / Type are already pointers in the request; assign directly.
	// A nil pointer means "no filter", which the storage layer honors.
	// Previously this code used `if req.Status != 0` as the sentinel,
	// which silently dropped filters for the iota-zero values
	// (StatusPending = 0, TypeDiscount = 0).
	query.Status = req.Status
	query.Type = req.Type
	query.PageQuery.Validate()

	promotions, total, err := a.promotionRepo.FindList(ctx, a.db, query)
	if err != nil {
		return nil, err
	}

	resp := &PromotionListResponse{
		List:     make([]*PromotionResponse, len(promotions)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	for i, p := range promotions {
		resp.List[i] = toPromotionResponse(p)
	}

	return resp, nil
}

func (a *promotionApp) DeletePromotion(ctx context.Context, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, id)
	if err != nil {
		return err
	}

	// Only allow delete if promotion is not active
	if p.Status == pkgpromotion.StatusActive {
		return code.ErrPromotionCannotDelete
	}

	return a.db.Transaction(func(tx *gorm.DB) error {
		// Delete promotion (soft delete)
		if err := a.promotionRepo.Delete(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (a *promotionApp) ActivatePromotion(ctx context.Context, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, id)
	if err != nil {
		return err
	}

	// Validate time range
	now := time.Now().UTC()
	if now.After(p.EndAt) {
		return code.ErrPromotionExpired
	}

	// Update status
	p.Status = pkgpromotion.StatusActive
	p.Audit.Update(0)

	return a.promotionRepo.Update(ctx, a.db, p)
}

func (a *promotionApp) DeactivatePromotion(ctx context.Context, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, id)
	if err != nil {
		return err
	}

	// Update status
	p.Status = pkgpromotion.StatusPaused
	p.Audit.Update(0)

	return a.promotionRepo.Update(ctx, a.db, p)
}

// toPromotionResponse 转换为响应DTO
func toPromotionResponse(p *pkgpromotion.Promotion) *PromotionResponse {
	rules := make([]*PromotionRuleResponse, 0, len(p.Rules))
	for _, r := range p.Rules {
		rules = append(rules, &PromotionRuleResponse{
			ID:             r.ID,
			PromotionID:    r.PromotionID,
			ConditionType:  int(r.ConditionType),
			ConditionValue: r.ConditionValue,
			ActionType:     int(r.ActionType),
			ActionValue:    r.ActionValue,
			MaxDiscount:    r.MaxDiscount,
		})
	}

	return &PromotionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Type:        int(p.Type),
		Status:      int(p.Status),
		StartAt:     p.StartAt.Format(time.RFC3339),
		EndAt:       p.EndAt.Format(time.RFC3339),
		ScopeType:    string(p.Scope.Type),
		Rules:        rules,
		UsageLimit:   p.UsageLimit,
		PerUserLimit: p.PerUserLimit,
		Tags:         p.Tags,
		CreatedAt:    p.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.Audit.UpdatedAt.Format(time.RFC3339),
	}
}
