package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type promotionRepo struct{}

func NewPromotionRepository() promotion.Repository {
	return &promotionRepo{}
}

// decodeJSONInt64Slice parses a MySQL JSON column value into an int64 slice.
// Returns nil on missing / malformed / empty payloads.
func decodeJSONInt64Slice(s string) []int64 {
	if s == "" {
		return nil
	}
	var out []int64
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}

// encodeJSONInt64Slice renders an int64 slice as a MySQL JSON column value.
func encodeJSONInt64Slice(in []int64) string {
	if len(in) == 0 {
		return "null"
	}
	b, err := json.Marshal(in)
	if err != nil {
		return "null"
	}
	return string(b)
}

// decodeJSONStringSlice parses a MySQL JSON column value into a string slice.
func decodeJSONStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}

// encodeJSONStringSlice renders a string slice as a MySQL JSON column value.
func encodeJSONStringSlice(in []string) string {
	if len(in) == 0 {
		return "null"
	}
	b, err := json.Marshal(in)
	if err != nil {
		return "null"
	}
	return string(b)
}

// promotionModel represents the database model for the merged promotions table.
// Both system promotions and claim-based coupons live here, distinguished by
// the `kind` column. Coupon-specific fields (Code, MarketID, TotalCount,
// UsedCount) are nullable per the live schema.
type promotionModel struct {
	application.Model
	TenantID     int64            `gorm:"column:tenant_id;not null;index"`
	Kind         string           `gorm:"column:kind;type:enum('PROMOTION','COUPON');not null;default:PROMOTION"`
	Name         string           `gorm:"column:name;size:255;not null"`
	Code         *string          `gorm:"column:code;size:100"`
	Description  string           `gorm:"column:description;type:text"`
	Type         int              `gorm:"column:type;not null;default:0"`
	Status       int              `gorm:"column:status;not null;default:0;index"`
	Priority     int              `gorm:"column:priority;not null;default:0"`
	MarketID     *int64           `gorm:"column:market_id;index"`
	Currency     string           `gorm:"column:currency;size:10;not null;default:CNY"`
	TotalCount   *int             `gorm:"column:total_count"`
	UsedCount    *int             `gorm:"column:used_count"`
	UsageLimit   int              `gorm:"column:usage_limit;not null;default:0"`
	PerUserLimit int              `gorm:"column:per_user_limit;not null;default:1"`
	Tags         string           `gorm:"column:tags;type:json"` // JSON-encoded []string
	ScopeType    string           `gorm:"column:scope_type;size:32;not null;default:STOREWIDE"`
	ScopeIDs     string           `gorm:"column:scope_ids;type:json"`   // JSON array of int64
	ExcludeIDs   string           `gorm:"column:exclude_ids;type:json"` // JSON array of int64
	StartAt      time.Time        `gorm:"column:start_at;not null"`
	EndAt        time.Time        `gorm:"column:end_at;not null"`
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (*promotionModel) TableName() string { return "promotions" }

// toEntity converts the persistence model into the domain entity.
// Nullable JSON columns are decoded into Go slices/pointers; nil on missing.
func (m *promotionModel) toEntity() *promotion.Promotion {
	p := &promotion.Promotion{
		ID:           m.ID,
		TenantID:     shared.TenantID(m.TenantID),
		Kind:         promotion.Kind(m.Kind),
		Name:         m.Name,
		Description:  m.Description,
		Code:         m.Code,
		Type:         promotion.Type(m.Type),
		Status:       promotion.Status(m.Status),
		Priority:     m.Priority,
		MarketID:     m.MarketID,
		Currency:     m.Currency,
		TotalCount:   m.TotalCount,
		UsedCount:    m.UsedCount,
		UsageLimit:   m.UsageLimit,
		PerUserLimit: m.PerUserLimit,
		Scope: promotion.PromotionScope{
			Type:       promotion.ScopeType(m.ScopeType),
			IDs:        decodeJSONInt64Slice(m.ScopeIDs),
			ExcludeIDs: decodeJSONInt64Slice(m.ExcludeIDs),
		},
		StartAt: m.StartAt.UTC(),
		EndAt:   m.EndAt.UTC(),
		Tags:    decodeJSONStringSlice(m.Tags),
		Audit: shared.AuditInfo{
			CreatedAt: m.Audit.CreatedAt.UTC(),
			UpdatedAt: m.Audit.UpdatedAt.UTC(),
			CreatedBy: m.Audit.CreatedBy,
			UpdatedBy: m.Audit.UpdatedBy,
		},
		DeletedAt: nil,
	}
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		p.DeletedAt = &t
	}
	return p
}

// fromPromotionEntity converts a domain entity into a persistence model.
// Code/MarketID/TotalCount/UsedCount pass through as pointers.
func fromPromotionEntity(p *promotion.Promotion) *promotionModel {
	m := &promotionModel{
		TenantID:     p.TenantID.Int64(),
		Kind:         string(p.Kind),
		Name:         p.Name,
		Description:  p.Description,
		Code:         p.Code,
		Type:         int(p.Type),
		Status:       int(p.Status),
		Priority:     p.Priority,
		MarketID:     p.MarketID,
		Currency:     p.Currency,
		TotalCount:   p.TotalCount,
		UsedCount:    p.UsedCount,
		UsageLimit:   p.UsageLimit,
		PerUserLimit: p.PerUserLimit,
		ScopeType:    string(p.Scope.Type),
		ScopeIDs:     encodeJSONInt64Slice(p.Scope.IDs),
		ExcludeIDs:   encodeJSONInt64Slice(p.Scope.ExcludeIDs),
		StartAt:      p.StartAt.UTC(),
		EndAt:        p.EndAt.UTC(),
		Tags:         encodeJSONStringSlice(p.Tags),
		Audit: shared.AuditInfo{
			CreatedAt: p.Audit.CreatedAt.UTC(),
			UpdatedAt: p.Audit.UpdatedAt.UTC(),
			CreatedBy: p.Audit.CreatedBy,
			UpdatedBy: p.Audit.UpdatedBy,
		},
	}
	m.ID = p.ID
	return m
}

// promotionRuleModel represents the database model for PromotionRule.
// Rules attach to either a PROMOTION or a COUPON via (owner_kind, owner_id).
// The legacy `promotion_id` column is nullable: populated only when the rule
// belongs to a system promotion (owner_kind = PROMOTION). Coupon rules have
// promotion_id = NULL.
type promotionRuleModel struct {
	application.Model
	OwnerKind           string          `gorm:"column:owner_kind;type:enum('PROMOTION','COUPON');not null;default:PROMOTION"`
	OwnerID             int64           `gorm:"column:owner_id;not null"`
	PromotionID         *int64          `gorm:"column:promotion_id"`
	ConditionType       int             `gorm:"column:condition_type;not null;default:0"`
	ConditionValue      decimal.Decimal `gorm:"column:condition_value;type:decimal(19,4);not null;default:0"`
	ActionType          int             `gorm:"column:action_type;not null;default:0"`
	ActionValue         decimal.Decimal `gorm:"column:action_value;type:decimal(19,4);not null;default:0"`
	MaxDiscountAmount   decimal.Decimal `gorm:"column:max_discount_amount;type:decimal(19,4);not null;default:0"`
	MaxDiscountCurrency string          `gorm:"column:max_discount_currency;size:10;default:CNY"`
	Currency            string          `gorm:"column:currency;size:10;not null;default:CNY"`
	SortOrder           int             `gorm:"column:sort_order;not null;default:0"`
	// NOTE: promotion_rules does NOT have created_by/updated_by columns.
	// Do NOT embed shared.AuditInfo here — it would inject those columns into
	// INSERT statements and fail with "Unknown column 'created_by'".
}

func (*promotionRuleModel) TableName() string { return "promotion_rules" }

// toEntity converts a rule model into the domain PromotionRule entity.
// When the rule belongs to a coupon, OwnerID is the promotion.id of the
// COUPON row, and the legacy PromotionID stays nil.
func (m *promotionRuleModel) toEntity() *promotion.PromotionRule {
	ownerID := m.OwnerID
	if m.OwnerKind == string(promotion.KindPromotion) && m.PromotionID != nil {
		ownerID = *m.PromotionID
	}
	return &promotion.PromotionRule{
		ID:             m.ID,
		OwnerKind:      promotion.Kind(m.OwnerKind),
		OwnerID:        ownerID,
		ConditionType:  promotion.ConditionType(m.ConditionType),
		ConditionValue: m.ConditionValue,
		ActionType:     promotion.ActionType(m.ActionType),
		ActionValue:    m.ActionValue,
		MaxDiscount:    m.MaxDiscountAmount,
		Currency:       m.Currency,
		SortOrder:      m.SortOrder,
		CreatedAt:      m.CreatedAt.UTC(),
		UpdatedAt:      m.UpdatedAt.UTC(),
	}
}

// fromPromotionRuleEntity converts a domain PromotionRule into a model row.
// promotion_id is populated ONLY when owner_kind == PROMOTION.
func fromPromotionRuleEntity(r *promotion.PromotionRule) *promotionRuleModel {
	m := &promotionRuleModel{
		OwnerKind:           string(r.OwnerKind),
		OwnerID:             r.OwnerID,
		ConditionType:       int(r.ConditionType),
		ConditionValue:      r.ConditionValue,
		ActionType:          int(r.ActionType),
		ActionValue:         r.ActionValue,
		MaxDiscountAmount:   r.MaxDiscount,
		MaxDiscountCurrency: r.Currency,
		Currency:            r.Currency,
		SortOrder:           r.SortOrder,
		// NOTE: promotion_rules table has no created_by/updated_by columns.
		// application.Model (embedded above) handles created_at/updated_at.
	}
	if r.OwnerKind == promotion.KindPromotion {
		pid := r.OwnerID
		m.PromotionID = &pid
	}
	m.ID = r.ID
	return m
}

// userCouponModel represents the database model for user_coupons.
// coupon_id numerically equals promotions.id for COUPON rows.
type userCouponModel struct {
	application.Model
	TenantID   int64            `gorm:"column:tenant_id;not null;index"`
	UserID     int64            `gorm:"column:user_id;not null;index"`
	CouponID   int64            `gorm:"column:coupon_id;not null;index"`
	Status     int              `gorm:"column:status;not null;index"`
	UsedAt     *time.Time       `gorm:"column:used_at"`
	OrderID    int64            `gorm:"column:order_id"`
	ReceivedAt time.Time        `gorm:"column:received_at;not null"`
	ExpireAt   time.Time        `gorm:"column:expire_at;not null;index"`
	Audit      shared.AuditInfo `gorm:"embedded"`
}

func (*userCouponModel) TableName() string { return "user_coupons" }

func (m *userCouponModel) toEntity() *promotion.UserCoupon {
	return &promotion.UserCoupon{
		ID:         m.ID,
		TenantID:   shared.TenantID(m.TenantID),
		UserID:     m.UserID,
		CouponID:   m.CouponID,
		Status:     promotion.UserCouponStatus(m.Status),
		UsedAt:     m.UsedAt,
		OrderID:    m.OrderID,
		ReceivedAt: m.ReceivedAt,
		ExpireAt:   m.ExpireAt,
		CreatedAt:  m.Audit.CreatedAt,
		UpdatedAt:  m.Audit.UpdatedAt,
	}
}

// promotionUsageModel represents the database model for promotion_usage.
type promotionUsageModel struct {
	application.Model
	TenantID       int64            `gorm:"column:tenant_id;not null;index"`
	PromotionID    int64            `gorm:"column:promotion_id;not null;index"`
	RuleID         *int64           `gorm:"column:rule_id;index"`
	OrderID        int64            `gorm:"column:order_id;not null;index"`
	UserID         int64            `gorm:"column:user_id;not null;index"`
	DiscountAmount decimal.Decimal  `gorm:"column:discount_amount;type:decimal(19,4);not null"`
	Currency       string           `gorm:"column:currency;size:10;not null"`
	OriginalAmount decimal.Decimal  `gorm:"column:original_amount;type:decimal(19,4);not null"`
	FinalAmount    decimal.Decimal  `gorm:"column:final_amount;type:decimal(19,4);not null"`
	CouponID       *int64           `gorm:"column:coupon_id;index"`
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (*promotionUsageModel) TableName() string { return "promotion_usage" }

func (m *promotionUsageModel) toEntity() *promotion.PromotionUsage {
	return &promotion.PromotionUsage{
		ID:             m.ID,
		TenantID:       shared.TenantID(m.TenantID),
		PromotionID:    m.PromotionID,
		RuleID:         m.RuleID,
		OrderID:        m.OrderID,
		UserID:         m.UserID,
		DiscountAmount: m.DiscountAmount,
		Currency:       m.Currency,
		OriginalAmount: m.OriginalAmount,
		FinalAmount:    m.FinalAmount,
		CouponID:       m.CouponID,
		CreatedAt:      m.Audit.CreatedAt,
	}
}

// Create inserts a new promotion.
func (r *promotionRepo) Create(ctx context.Context, db *gorm.DB, p *promotion.Promotion) error {
	model := fromPromotionEntity(p)
	return db.WithContext(ctx).Create(model).Error
}

// Update updates an existing promotion. The legacy `code` column is only
// relevant for COUPON rows; the update is idempotent for PROMOTION rows.
func (r *promotionRepo) Update(ctx context.Context, db *gorm.DB, p *promotion.Promotion) error {
	updates := map[string]any{
		"name":           p.Name,
		"description":    p.Description,
		"type":           int(p.Type),
		"status":         int(p.Status),
		"priority":       p.Priority,
		"start_at":       p.StartAt.UTC(),
		"end_at":         p.EndAt.UTC(),
		"scope_type":     string(p.Scope.Type),
		"scope_ids":      encodeJSONInt64Slice(p.Scope.IDs),
		"exclude_ids":    encodeJSONInt64Slice(p.Scope.ExcludeIDs),
		"usage_limit":    p.UsageLimit,
		"per_user_limit": p.PerUserLimit,
		"tags":           encodeJSONStringSlice(p.Tags),
		"currency":       p.Currency,
		"updated_at":     time.Now().UTC(),
		"updated_by":     p.Audit.UpdatedBy,
	}
	if p.MarketID != nil {
		updates["market_id"] = *p.MarketID
	} else {
		updates["market_id"] = nil
	}
	if p.TotalCount != nil {
		updates["total_count"] = *p.TotalCount
	}
	if p.UsedCount != nil {
		updates["used_count"] = *p.UsedCount
	}
	if p.Code != nil {
		updates["code"] = *p.Code
	}

	return db.WithContext(ctx).
		Model(&promotionModel{}).
		Where("id = ? AND deleted_at IS NULL", p.ID).
		Updates(updates).Error
}

// Delete soft-deletes a promotion by id.
func (r *promotionRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	now := time.Now().UTC()
	result := db.WithContext(ctx).
		Model(&promotionModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrPromotionNotFound
	}
	return nil
}

// FindByID returns a promotion (any kind) by id.
func (r *promotionRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*promotion.Promotion, error) {
	var model promotionModel
	err := db.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPromotionNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByCode looks up a COUPON by its (unique) code. PROMOTIONs do not have
// codes and will never match this query.
func (r *promotionRepo) FindByCode(ctx context.Context, db *gorm.DB, codeStr string) (*promotion.Promotion, error) {
	var model promotionModel
	err := db.WithContext(ctx).
		Where("kind = ? AND code = ? AND deleted_at IS NULL", promotion.KindCoupon, codeStr).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPromotionNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindList returns promotions with pagination and filters. All enum filters
// use pointer types so "no filter" can be distinguished from the iota-zero
// value (e.g. PROMOTION, ACTIVE, DISCOUNT).
func (r *promotionRepo) FindList(ctx context.Context, db *gorm.DB, query promotion.Query) ([]*promotion.Promotion, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&promotionModel{}).Where("deleted_at IS NULL")

	if query.TenantID > 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", query.TenantID.Int64())
	}
	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Kind != nil && query.Kind.IsValid() {
		dbQuery = dbQuery.Where("kind = ?", *query.Kind)
	}
	if query.Status != nil && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}
	if query.Type != nil && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", *query.Type)
	}
	if query.MarketID != nil {
		dbQuery = dbQuery.Where("market_id = ?", *query.MarketID)
	}
	if query.ExpiredOnly {
		// "expired" is derived from EndAt, not a stored enum value.
		dbQuery = dbQuery.Where("end_at <= ?", time.Now().UTC())
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

// CreateRules inserts multiple rules for a single owner (promotion or coupon).
// The owner_kind / owner_id columns are stamped from the parameters; the
// legacy `promotion_id` is derived automatically by fromPromotionRuleEntity.
func (r *promotionRepo) CreateRules(ctx context.Context, db *gorm.DB, ownerKind promotion.Kind, ownerID int64, rules []promotion.PromotionRule) error {
	if len(rules) == 0 {
		return nil
	}
	models := make([]*promotionRuleModel, len(rules))
	for i, rule := range rules {
		rule.OwnerKind = ownerKind
		rule.OwnerID = ownerID
		models[i] = fromPromotionRuleEntity(&rule)
	}
	return db.WithContext(ctx).Create(&models).Error
}

// FindRulesByOwner returns all rules attached to a given owner.
func (r *promotionRepo) FindRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind promotion.Kind, ownerID int64) ([]promotion.PromotionRule, error) {
	var models []promotionRuleModel
	err := db.WithContext(ctx).
		Where("owner_kind = ? AND owner_id = ?", ownerKind, ownerID).
		Order("sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	rules := make([]promotion.PromotionRule, len(models))
	for i, m := range models {
		rules[i] = *m.toEntity()
	}
	return rules, nil
}

// FindRulesByOwners fetches every rule whose owner_id appears in the
// given set and groups the result by owner_id. Owner_kind is intentionally
// not part of the filter: owner_id is unique within the promotions table,
// so matching on id alone is sufficient. One round-trip replaces N.
func (r *promotionRepo) FindRulesByOwners(ctx context.Context, db *gorm.DB, ownerIDs []int64) (map[int64][]promotion.PromotionRule, error) {
	out := make(map[int64][]promotion.PromotionRule, len(ownerIDs))
	if len(ownerIDs) == 0 {
		return out, nil
	}
	var models []promotionRuleModel
	err := db.WithContext(ctx).
		Where("owner_id IN ?", ownerIDs).
		Order("sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	for _, m := range models {
		ownerID := m.OwnerID
		if m.OwnerKind == string(promotion.KindPromotion) && m.PromotionID != nil {
			ownerID = *m.PromotionID
		}
		rule := m.toEntity()
		rule.OwnerID = ownerID
		out[ownerID] = append(out[ownerID], *rule)
	}
	return out, nil
}

// UpdateRule updates an existing rule row.
func (r *promotionRepo) UpdateRule(ctx context.Context, db *gorm.DB, rule *promotion.PromotionRule) error {
	model := fromPromotionRuleEntity(rule)
	now := time.Now().UTC()
	return db.WithContext(ctx).
		Model(&promotionRuleModel{}).
		Where("id = ?", rule.ID).
		Updates(map[string]any{
			"condition_type":        model.ConditionType,
			"condition_value":       model.ConditionValue,
			"action_type":           model.ActionType,
			"action_value":          model.ActionValue,
			"max_discount_amount":   model.MaxDiscountAmount,
			"max_discount_currency": model.MaxDiscountCurrency,
			"currency":              model.Currency,
			"sort_order":            model.SortOrder,
			"updated_at":            now,
		}).Error
}

// DeleteRule deletes a single rule by id.
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

// DeleteRulesByOwner removes every rule belonging to a given owner.
func (r *promotionRepo) DeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind promotion.Kind, ownerID int64) error {
	return db.WithContext(ctx).
		Where("owner_kind = ? AND owner_id = ?", ownerKind, ownerID).
		Delete(&promotionRuleModel{}).Error
}

// HardDeleteRulesByOwner permanently removes rules for an owner, bypassing
// GORM's soft-delete (deleted_at). Used by PromotionApp.Update when replacing
// the ruleset so old rows don't accumulate as tombstones across replace
// cycles. Compare with DeleteRulesByOwner (soft-delete) which is kept for
// full promotion delete so an undeleted promotion recovers its rules.
func (r *promotionRepo) HardDeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind promotion.Kind, ownerID int64) error {
	return db.WithContext(ctx).
		Unscoped().
		Where("owner_kind = ? AND owner_id = ?", ownerKind, ownerID).
		Delete(&promotionRuleModel{}).Error
}

// FindActiveCoupons returns claim-eligible coupons for a market (or any market
// when marketID is nil). Eligibility = status active, time window open, and
// inventory not depleted (total_count IS NULL OR used_count < total_count).
func (r *promotionRepo) FindActiveCoupons(ctx context.Context, db *gorm.DB, marketID *int64) ([]*promotion.Promotion, error) {
	now := time.Now().UTC()
	q := db.WithContext(ctx).Model(&promotionModel{}).
		Where("deleted_at IS NULL").
		Where("kind = ?", promotion.KindCoupon).
		Where("status = ?", promotion.StatusActive).
		Where("start_at <= ? AND end_at >= ?", now, now).
		Where("(total_count IS NULL OR used_count < total_count)")
	if marketID != nil {
		q = q.Where("market_id IS NULL OR market_id = ?", *marketID)
	}
	var models []promotionModel
	if err := q.Order("priority DESC, created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*promotion.Promotion, len(models))
	for i, m := range models {
		out[i] = m.toEntity()
	}
	return out, nil
}

// IncrementUsedCount atomically consumes one unit of coupon inventory. The
// inventory check is enforced at the database level (used_count < total_count)
// so concurrent callers cannot oversell.
func (r *promotionRepo) IncrementUsedCount(ctx context.Context, db *gorm.DB, couponID int64) error {
	res := db.WithContext(ctx).Exec(`
		UPDATE promotions
		SET used_count = used_count + 1
		WHERE id = ? AND kind = 'COUPON'
		  AND (total_count IS NULL OR used_count < total_count)
	`, couponID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return code.ErrCouponUsedUp
	}
	return nil
}

// IssueUserCoupon writes a new user_coupons row pointing at a COUPON
// (coupon_id numerically equals promotions.id for COUPON rows, per the
// merge handoff in the design spec).
func (r *promotionRepo) IssueUserCoupon(ctx context.Context, db *gorm.DB, uc *promotion.UserCoupon) error {
	return db.WithContext(ctx).Create(&userCouponModel{
		TenantID:   uc.TenantID.Int64(),
		UserID:     uc.UserID,
		CouponID:   uc.CouponID,
		Status:     int(uc.Status),
		UsedAt:     uc.UsedAt,
		OrderID:    uc.OrderID,
		ReceivedAt: uc.ReceivedAt,
		ExpireAt:   uc.ExpireAt,
	}).Error
}

// FindUserCoupons returns paginated user_coupons rows. When CouponID is set,
// the join ensures we only return claims against rows still tagged COUPON —
// this guards against the case where a row is later converted.
func (r *promotionRepo) FindUserCoupons(ctx context.Context, db *gorm.DB, query promotion.UserCouponQuery) ([]*promotion.UserCoupon, int64, error) {
	page, size := normalizePage(query.Page, query.Size)

	dbQuery := db.WithContext(ctx).Model(&userCouponModel{})
	if query.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *query.UserID)
	}
	if query.CouponID != nil {
		// Join to promotions to filter by kind='COUPON' so we never return
		// user_coupons pointing at a non-COUPON row (defensive against bad IDs).
		dbQuery = dbQuery.
			Joins("JOIN promotions p ON p.id = user_coupons.coupon_id").
			Where("p.kind = ?", promotion.KindCoupon).
			Where("user_coupons.coupon_id = ?", *query.CouponID)
	}
	if query.Status != nil {
		dbQuery = dbQuery.Where("user_coupons.status = ?", *query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []userCouponModel
	err := dbQuery.Order("user_coupons.received_at DESC").
		Offset(page * size).
		Limit(size).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]*promotion.UserCoupon, len(models))
	for i, m := range models {
		out[i] = m.toEntity()
	}
	return out, total, nil
}

// FindPromotionUsage returns paginated promotion_usage rows.
func (r *promotionRepo) FindPromotionUsage(ctx context.Context, db *gorm.DB, query promotion.UsageQuery) ([]*promotion.PromotionUsage, int64, error) {
	page, size := normalizePage(query.Page, query.Size)

	dbQuery := db.WithContext(ctx).Model(&promotionUsageModel{})
	if query.CouponID != nil {
		dbQuery = dbQuery.Where("coupon_id = ?", *query.CouponID)
	}
	if query.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *query.UserID)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []promotionUsageModel
	err := dbQuery.Order("created_at DESC").
		Offset(page * size).
		Limit(size).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]*promotion.PromotionUsage, len(models))
	for i, m := range models {
		out[i] = m.toEntity()
	}
	return out, total, nil
}

// normalizePage applies the same defaults as shared.PageQuery.Validate so that
// query types without Validate (UserCouponQuery, UsageQuery) still get safe
// pagination values.
func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
