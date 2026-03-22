package promotion

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

// CreatePromotionRequest 创建促销请求
type CreatePromotionRequest struct {
	Name        string
	Description string
	Type        pkgpromotion.Type
	Priority    int
	Currency    string
	Scope       pkgpromotion.PromotionScope
	StartAt     time.Time
	EndAt       time.Time
	Rules       []CreatePromotionRuleRequest
}

// CreatePromotionRuleRequest 创建促销规则请求
type CreatePromotionRuleRequest struct {
	ConditionType  pkgpromotion.ConditionType
	ConditionValue int64
	ActionType     pkgpromotion.ActionType
	ActionValue    int64
	MaxDiscount    int64
}

// UpdatePromotionRequest 更新促销请求
type UpdatePromotionRequest struct {
	ID          int64
	Name        string
	Description string
	StartAt     time.Time
	EndAt       time.Time
}

// PromotionResponse 促销响应
type PromotionResponse struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        int                      `json:"type"`
	Status      int                      `json:"status"`
	StartAt     string                   `json:"start_at"`
	EndAt       string                   `json:"end_at"`
	Rules       []*PromotionRuleResponse `json:"rules"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// PromotionRuleResponse 促销规则响应
type PromotionRuleResponse struct {
	ID             int64 `json:"id"`
	PromotionID    int64 `json:"promotion_id"`
	ConditionType  int   `json:"condition_type"`
	ConditionValue int64 `json:"condition_value"`
	ActionType     int   `json:"action_type"`
	ActionValue    int64 `json:"action_value"`
	MaxDiscount    int64 `json:"max_discount"`
}

// PromotionListResponse 促销列表响应
type PromotionListResponse struct {
	List     []*PromotionResponse `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// QueryPromotionRequest 查询促销请求
type QueryPromotionRequest struct {
	Name     string
	Status   pkgpromotion.Status
	Type     pkgpromotion.Type
	Page     int
	PageSize int
}

// PromotionApp 促销应用服务接口
type PromotionApp interface {
	CreatePromotion(ctx context.Context, tenantID shared.TenantID, req CreatePromotionRequest) (*PromotionResponse, error)
	UpdatePromotion(ctx context.Context, tenantID shared.TenantID, req UpdatePromotionRequest) (*PromotionResponse, error)
	GetPromotion(ctx context.Context, tenantID shared.TenantID, id int64) (*PromotionResponse, error)
	ListPromotions(ctx context.Context, tenantID shared.TenantID, req QueryPromotionRequest) (*PromotionListResponse, error)
	DeletePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error
	ActivatePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error
	DeactivatePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error
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

func (a *promotionApp) CreatePromotion(ctx context.Context, tenantID shared.TenantID, req CreatePromotionRequest) (*PromotionResponse, error) {
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
			TenantID:    tenantID,
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Status:      pkgpromotion.StatusPending,
			Priority:    req.Priority,
			StartAt:     req.StartAt.UTC(),
			EndAt:       req.EndAt.UTC(),
			Scope:       req.Scope,
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

func (a *promotionApp) UpdatePromotion(ctx context.Context, tenantID shared.TenantID, req UpdatePromotionRequest) (*PromotionResponse, error) {
	p, err := a.promotionRepo.FindByID(ctx, a.db, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Only allow update if promotion is not active
	if p.Status == pkgpromotion.StatusActive {
		return nil, ErrPromotionIsActive
	}

	p.Name = req.Name
	p.Description = req.Description
	p.StartAt = req.StartAt.UTC()
	p.EndAt = req.EndAt.UTC()
	p.Audit.Update(0)

	if err := a.promotionRepo.Update(ctx, a.db, p); err != nil {
		return nil, err
	}

	return toPromotionResponse(p), nil
}

func (a *promotionApp) GetPromotion(ctx context.Context, tenantID shared.TenantID, id int64) (*PromotionResponse, error) {
	p, err := a.promotionRepo.FindByID(ctx, a.db, tenantID, id)
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

func (a *promotionApp) ListPromotions(ctx context.Context, tenantID shared.TenantID, req QueryPromotionRequest) (*PromotionListResponse, error) {
	query := pkgpromotion.Query{
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

	promotions, total, err := a.promotionRepo.FindList(ctx, a.db, tenantID, query)
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

func (a *promotionApp) DeletePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return err
	}

	// Only allow delete if promotion is not active
	if p.Status == pkgpromotion.StatusActive {
		return ErrPromotionIsActive
	}

	return a.db.Transaction(func(tx *gorm.DB) error {
		// Delete promotion (soft delete)
		if err := a.promotionRepo.Delete(ctx, tx, tenantID, id); err != nil {
			return err
		}
		return nil
	})
}

func (a *promotionApp) ActivatePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, tenantID, id)
	if err != nil {
		return err
	}

	// Validate time range
	now := time.Now().UTC()
	if now.After(p.EndAt) {
		return ErrPromotionAlreadyEnded
	}

	// Update status
	p.Status = pkgpromotion.StatusActive
	p.Audit.Update(0)

	return a.promotionRepo.Update(ctx, a.db, p)
}

func (a *promotionApp) DeactivatePromotion(ctx context.Context, tenantID shared.TenantID, id int64) error {
	p, err := a.promotionRepo.FindByID(ctx, a.db, tenantID, id)
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
		Rules:       rules,
		CreatedAt:   p.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}