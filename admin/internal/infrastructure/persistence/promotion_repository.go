package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type promotionRepo struct{}

func NewPromotionRepository() promotion.Repository {
	return &promotionRepo{}
}

// promotionModel represents the database model for Promotion
type promotionModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID    int64  `gorm:"column:tenant_id;not null;index"`
	Name        string `gorm:"column:name;size:200;not null;index"`
	Description string `gorm:"column:description;type:text"`
	Type        int    `gorm:"column:type;not null;index"`
	Status      int    `gorm:"column:status;not null;index"`
	Priority    int    `gorm:"column:priority;not null;default:0"`
	StartAt     int64  `gorm:"column:start_at;not null;index"`
	EndAt       int64  `gorm:"column:end_at;not null;index"`
	ScopeType   string `gorm:"column:scope_type;size:32;not null"`
	ScopeIDs    string `gorm:"column:scope_ids;type:json"`       // JSON array of int64
	ExcludeIDs  string `gorm:"column:exclude_ids;type:json"`     // JSON array of int64
	Currency    string `gorm:"column:currency;size:10;not null"` // ISO 4217
	CreatedBy   int64  `gorm:"column:created_by;not null"`
	UpdatedBy   int64  `gorm:"column:updated_by;not null"`
	DeletedAt   *int64 `gorm:"column:deleted_at;index"`
	CreatedAt   int64  `gorm:"column:created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at"`
}

func (promotionModel) TableName() string {
	return "promotions"
}

func (m *promotionModel) toEntity() *promotion.Promotion {
	// Parse JSON arrays for scope
	var scopeIDs, excludeIDs []int64
	if m.ScopeIDs != "" {
		json.Unmarshal([]byte(m.ScopeIDs), &scopeIDs)
	}
	if m.ExcludeIDs != "" {
		json.Unmarshal([]byte(m.ExcludeIDs), &excludeIDs)
	}

	var deletedAt *time.Time
	if m.DeletedAt != nil {
		t := time.Unix(*m.DeletedAt, 0)
		deletedAt = &t
	}

	return &promotion.Promotion{
		ID:          m.ID,
		TenantID:    shared.TenantID(m.TenantID),
		Name:        m.Name,
		Description: m.Description,
		Type:        promotion.Type(m.Type),
		Status:      promotion.Status(m.Status),
		Priority:    m.Priority,
		StartAt:     time.Unix(m.StartAt, 0),
		EndAt:       time.Unix(m.EndAt, 0),
		Scope: promotion.PromotionScope{
			Type:       promotion.ScopeType(m.ScopeType),
			IDs:        scopeIDs,
			ExcludeIDs: excludeIDs,
		},
		Currency: m.Currency,
		Audit: shared.AuditInfo{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
		DeletedAt: deletedAt,
	}
}

func fromPromotionEntity(p *promotion.Promotion) *promotionModel {
	// Serialize JSON arrays for scope
	scopeIDsJSON, _ := json.Marshal(p.Scope.IDs)
	excludeIDsJSON, _ := json.Marshal(p.Scope.ExcludeIDs)

	var deletedAt *int64
	if p.DeletedAt != nil {
		ts := p.DeletedAt.Unix()
		deletedAt = &ts
	}

	return &promotionModel{
		ID:          p.ID,
		TenantID:    p.TenantID.Int64(),
		Name:        p.Name,
		Description: p.Description,
		Type:        int(p.Type),
		Status:      int(p.Status),
		Priority:    p.Priority,
		StartAt:     p.StartAt.Unix(),
		EndAt:       p.EndAt.Unix(),
		ScopeType:   string(p.Scope.Type),
		ScopeIDs:    string(scopeIDsJSON),
		ExcludeIDs:  string(excludeIDsJSON),
		Currency:    p.Currency,
		CreatedBy:   p.Audit.CreatedBy,
		UpdatedBy:   p.Audit.UpdatedBy,
		DeletedAt:   deletedAt,
		CreatedAt:   p.Audit.CreatedAt,
		UpdatedAt:   p.Audit.UpdatedAt,
	}
}

// promotionRuleModel represents the database model for PromotionRule
type promotionRuleModel struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	PromotionID    int64  `gorm:"column:promotion_id;not null;index"`
	ConditionType  int    `gorm:"column:condition_type;not null"`
	ConditionValue int64  `gorm:"column:condition_value;not null"`
	ActionType     int    `gorm:"column:action_type;not null"`
	ActionValue    int64  `gorm:"column:action_value;not null"`
	MaxDiscount    int64  `gorm:"column:max_discount;not null;default:0"`
	Currency       string `gorm:"column:currency;size:10;not null"`
	SortOrder      int    `gorm:"column:sort_order;not null;default:0"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
}

func (promotionRuleModel) TableName() string {
	return "promotion_rules"
}

func (m *promotionRuleModel) toEntity() promotion.PromotionRule {
	return promotion.PromotionRule{
		ID:             m.ID,
		PromotionID:    m.PromotionID,
		ConditionType:  promotion.ConditionType(m.ConditionType),
		ConditionValue: m.ConditionValue,
		ActionType:     promotion.ActionType(m.ActionType),
		ActionValue:    m.ActionValue,
		MaxDiscount:    m.MaxDiscount,
		Currency:       m.Currency,
		SortOrder:      m.SortOrder,
		CreatedAt:      time.Unix(m.CreatedAt, 0),
		UpdatedAt:      time.Unix(m.UpdatedAt, 0),
	}
}

func fromPromotionRuleEntity(r *promotion.PromotionRule) *promotionRuleModel {
	return &promotionRuleModel{
		ID:             r.ID,
		PromotionID:    r.PromotionID,
		ConditionType:  int(r.ConditionType),
		ConditionValue: r.ConditionValue,
		ActionType:     int(r.ActionType),
		ActionValue:    r.ActionValue,
		MaxDiscount:    r.MaxDiscount,
		Currency:       r.Currency,
		SortOrder:      r.SortOrder,
		CreatedAt:      r.CreatedAt.Unix(),
		UpdatedAt:      r.UpdatedAt.Unix(),
	}
}

// Create inserts a new promotion
func (r *promotionRepo) Create(ctx context.Context, db *gorm.DB, p *promotion.Promotion) error {
	model := fromPromotionEntity(p)
	return db.WithContext(ctx).Create(model).Error
}

// Update updates an existing promotion
func (r *promotionRepo) Update(ctx context.Context, db *gorm.DB, p *promotion.Promotion) error {
	model := fromPromotionEntity(p)
	return db.WithContext(ctx).
		Model(&promotionModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", p.ID, p.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":          model.Name,
			"description":   model.Description,
			"type":          model.Type,
			"status":        model.Status,
			"priority":      model.Priority,
			"start_at":      model.StartAt,
			"end_at":        model.EndAt,
			"scope_type":    model.ScopeType,
			"scope_ids":     model.ScopeIDs,
			"exclude_ids":   model.ExcludeIDs,
			"currency":      model.Currency,
			"updated_by":    model.UpdatedBy,
			"updated_at":    model.UpdatedAt,
		}).Error
}

// Delete soft deletes a promotion
func (r *promotionRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&promotionModel{}).Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can delete all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().Unix()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrPromotionNotFound
	}
	return nil
}

// FindByID finds a promotion by ID
func (r *promotionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*promotion.Promotion, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model promotionModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPromotionNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindActive finds all active promotions for a tenant
func (r *promotionRepo) FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*promotion.Promotion, error) {
	query := db.WithContext(ctx).Model(&promotionModel{}).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	now := time.Now().Unix()
	var models []promotionModel
	err := query.
		Where("status = ?", promotion.StatusActive).
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("priority DESC, created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	promotions := make([]*promotion.Promotion, len(models))
	for i, m := range models {
		promotions[i] = m.toEntity()
	}
	return promotions, nil
}

// FindActiveByCurrency finds all active promotions for a tenant filtered by currency
func (r *promotionRepo) FindActiveByCurrency(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, currency string) ([]*promotion.Promotion, error) {
	query := db.WithContext(ctx).Model(&promotionModel{}).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	now := time.Now().Unix()
	var models []promotionModel
	err := query.
		Where("status = ?", promotion.StatusActive).
		Where("currency = ?", currency).
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("priority DESC, created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	promotions := make([]*promotion.Promotion, len(models))
	for i, m := range models {
		promotions[i] = m.toEntity()
	}
	return promotions, nil
}

// FindList finds promotions with pagination and filters
func (r *promotionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query promotion.Query) ([]*promotion.Promotion, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&promotionModel{}).Where("deleted_at IS NULL")

	// Tenant filter: Platform admin (TenantID == 0) can access all tenant data
	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Status != nil && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}
	if query.Type != nil && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", *query.Type)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []promotionModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	promotions := make([]*promotion.Promotion, len(models))
	for i, m := range models {
		promotions[i] = m.toEntity()
	}
	return promotions, total, nil
}

// CreateRules inserts multiple promotion rules
func (r *promotionRepo) CreateRules(ctx context.Context, db *gorm.DB, rules []promotion.PromotionRule) error {
	if len(rules) == 0 {
		return nil
	}

	models := make([]*promotionRuleModel, len(rules))
	for i, rule := range rules {
		models[i] = fromPromotionRuleEntity(&rule)
	}
	return db.WithContext(ctx).Create(&models).Error
}

// FindRulesByPromotionID finds all rules for a promotion
func (r *promotionRepo) FindRulesByPromotionID(ctx context.Context, db *gorm.DB, promotionID int64) ([]promotion.PromotionRule, error) {
	var models []promotionRuleModel
	err := db.WithContext(ctx).
		Where("promotion_id = ?", promotionID).
		Order("sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	rules := make([]promotion.PromotionRule, len(models))
	for i, m := range models {
		rules[i] = m.toEntity()
	}
	return rules, nil
}

// FindRulesByPromotionIDs finds all rules for multiple promotions in a single query
func (r *promotionRepo) FindRulesByPromotionIDs(ctx context.Context, db *gorm.DB, promotionIDs []int64) (map[int64][]promotion.PromotionRule, error) {
	if len(promotionIDs) == 0 {
		return make(map[int64][]promotion.PromotionRule), nil
	}

	var models []promotionRuleModel
	err := db.WithContext(ctx).
		Where("promotion_id IN ?", promotionIDs).
		Order("promotion_id ASC, sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	rulesByPromotion := make(map[int64][]promotion.PromotionRule)
	for _, m := range models {
		rulesByPromotion[m.PromotionID] = append(rulesByPromotion[m.PromotionID], m.toEntity())
	}
	return rulesByPromotion, nil
}

// UpdateRule updates a promotion rule
func (r *promotionRepo) UpdateRule(ctx context.Context, db *gorm.DB, rule *promotion.PromotionRule) error {
	model := fromPromotionRuleEntity(rule)
	return db.WithContext(ctx).
		Model(&promotionRuleModel{}).
		Where("id = ?", rule.ID).
		Updates(map[string]interface{}{
			"condition_type":  model.ConditionType,
			"condition_value": model.ConditionValue,
			"action_type":     model.ActionType,
			"action_value":    model.ActionValue,
			"max_discount":    model.MaxDiscount,
			"currency":        model.Currency,
			"sort_order":      model.SortOrder,
			"updated_at":      model.UpdatedAt,
		}).Error
}

// DeleteRule deletes a promotion rule
func (r *promotionRepo) DeleteRule(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.WithContext(ctx).
		Model(&promotionRuleModel{}).
		Where("id = ?", id).
		Delete(&promotionRuleModel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrPromotionRuleNotFound
	}
	return nil
}