package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== Earn Rule Repository ====================

type earnRuleRepo struct{}

func NewEarnRuleRepository() points.EarnRuleRepository {
	return &earnRuleRepo{}
}

type earnRuleModel struct {
	application.Model
	TenantID         int64              `gorm:"column:tenant_id;not null;index:idx_tenant_id"`
	Name             string             `gorm:"column:name;type:varchar(255);not null"`
	Description      string             `gorm:"column:description;type:text"`
	Scenario         string             `gorm:"column:scenario;type:varchar(50);not null"`
	CalculationType  string             `gorm:"column:calculation_type;type:varchar(20);not null"`
	FixedPoints      int64              `gorm:"column:fixed_points;type:bigint;default:0"`
	Ratio            decimal.Decimal    `gorm:"column:ratio;type:decimal(10,4)"`
	Tiers            points.TierConfigs `gorm:"column:tiers;type:json"`
	ConditionType    string             `gorm:"column:condition_type;type:varchar(50);not null;default:'NONE'"`
	ConditionValue   string             `gorm:"column:condition_value;type:text"`
	ExpirationMonths int                `gorm:"column:expiration_months;type:int;not null;default:12"`
	Status           int                `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	Priority         int                `gorm:"column:priority;type:int;not null;default:0"`
	StartAt          *time.Time         `gorm:"column:start_at"`
	EndAt            *time.Time         `gorm:"column:end_at"`
}

func (earnRuleModel) TableName() string {
	return "earn_rules"
}

func (m *earnRuleModel) toEntity() *points.EarnRule {
	return &points.EarnRule{
		Model:            application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt},
		TenantID:         shared.TenantID(m.TenantID),
		Name:             m.Name,
		Description:      m.Description,
		Scenario:         points.EarnScenario(m.Scenario),
		CalculationType:  points.CalculationType(m.CalculationType),
		FixedPoints:      m.FixedPoints,
		Ratio:            m.Ratio,
		Tiers:            m.Tiers,
		ConditionType:    points.ConditionType(m.ConditionType),
		ConditionValue:   m.ConditionValue,
		ExpirationMonths: m.ExpirationMonths,
		Status:           points.EarnRuleStatus(m.Status),
		Priority:         m.Priority,
		StartAt:          m.StartAt,
		EndAt:            m.EndAt,
	}
}

func fromEarnRuleEntity(e *points.EarnRule) *earnRuleModel {
	return &earnRuleModel{
		Model:            application.Model{ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, DeletedAt: e.DeletedAt},
		TenantID:         e.TenantID.Int64(),
		Name:             e.Name,
		Description:      e.Description,
		Scenario:         string(e.Scenario),
		CalculationType:  string(e.CalculationType),
		FixedPoints:      e.FixedPoints,
		Ratio:            e.Ratio,
		Tiers:            e.Tiers,
		ConditionType:    string(e.ConditionType),
		ConditionValue:   e.ConditionValue,
		ExpirationMonths: e.ExpirationMonths,
		Status:           int(e.Status),
		Priority:         e.Priority,
		StartAt:          e.StartAt,
		EndAt:            e.EndAt,
	}
}

func (r *earnRuleRepo) Create(ctx context.Context, db *gorm.DB, rule *points.EarnRule) error {
	model := fromEarnRuleEntity(rule)
	return db.WithContext(ctx).Create(model).Error
}

func (r *earnRuleRepo) Update(ctx context.Context, db *gorm.DB, rule *points.EarnRule) error {
	model := fromEarnRuleEntity(rule)
	return db.WithContext(ctx).
		Model(&earnRuleModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rule.ID, rule.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":              model.Name,
			"description":       model.Description,
			"scenario":          model.Scenario,
			"calculation_type":  model.CalculationType,
			"fixed_points":      model.FixedPoints,
			"ratio":             model.Ratio,
			"tiers":             model.Tiers,
			"condition_type":    model.ConditionType,
			"condition_value":   model.ConditionValue,
			"expiration_months": model.ExpirationMonths,
			"status":            model.Status,
			"priority":          model.Priority,
			"start_at":          model.StartAt,
			"end_at":            model.EndAt,
			"updated_at":        model.UpdatedAt,
		}).Error
}

func (r *earnRuleRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&earnRuleModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrNotFound
	}
	return nil
}

func (r *earnRuleRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.EarnRule, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model earnRuleModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *earnRuleRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.EarnRuleQuery) ([]*points.EarnRule, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&earnRuleModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.Scenario != "" && query.Scenario.IsValid() {
		dbQuery = dbQuery.Where("scenario = ?", query.Scenario)
	}
	if query.CalculationType != "" && query.CalculationType.IsValid() {
		dbQuery = dbQuery.Where("calculation_type = ?", query.CalculationType)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []earnRuleModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	rules := make([]*points.EarnRule, len(models))
	for i, m := range models {
		rules[i] = m.toEntity()
	}
	return rules, total, nil
}

func (r *earnRuleRepo) FindByScenario(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, scenario points.EarnScenario) ([]*points.EarnRule, error) {
	query := db.WithContext(ctx).Model(&earnRuleModel{}).Where("deleted_at IS NULL AND scenario = ?", scenario)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var models []earnRuleModel
	err := query.Order("priority DESC, created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	rules := make([]*points.EarnRule, len(models))
	for i, m := range models {
		rules[i] = m.toEntity()
	}
	return rules, nil
}

func (r *earnRuleRepo) UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status points.EarnRuleStatus) error {
	query := db.WithContext(ctx).Model(&earnRuleModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	return query.Updates(map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}).Error
}

func (r *earnRuleRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.EarnRuleStats, error) {
	query := db.WithContext(ctx).Model(&earnRuleModel{}).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var total, active int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	activeQuery := db.WithContext(ctx).Model(&earnRuleModel{}).Where("deleted_at IS NULL AND status = ?", points.EarnRuleStatusActive)
	if tenantID != 0 {
		activeQuery = activeQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := activeQuery.Count(&active).Error; err != nil {
		return nil, err
	}

	return &points.EarnRuleStats{
		Total:  total,
		Active: active,
	}, nil
}

// ==================== Redeem Rule Repository ====================

type redeemRuleRepo struct{}

func NewRedeemRuleRepository() points.RedeemRuleRepository {
	return &redeemRuleRepo{}
}

type redeemRuleModel struct {
	application.Model
	TenantID       int64      `gorm:"column:tenant_id;not null;index:idx_tenant_id"`
	Name           string     `gorm:"column:name;type:varchar(255);not null"`
	Description    string     `gorm:"column:description;type:text"`
	CouponID       int64      `gorm:"column:coupon_id;type:bigint;not null"`
	PointsRequired int64      `gorm:"column:points_required;type:bigint;not null"`
	TotalStock     int64      `gorm:"column:total_stock;type:bigint;not null;default:0"`
	UsedStock      int64      `gorm:"column:used_stock;type:bigint;not null;default:0"`
	PerUserLimit   int        `gorm:"column:per_user_limit;type:int;not null;default:1"`
	Status         int        `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	StartAt        *time.Time `gorm:"column:start_at"`
	EndAt          *time.Time `gorm:"column:end_at"`
}

func (redeemRuleModel) TableName() string {
	return "redeem_rules"
}

func (m *redeemRuleModel) toEntity() *points.RedeemRule {
	return &points.RedeemRule{
		Model:          application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt},
		TenantID:       shared.TenantID(m.TenantID),
		Name:           m.Name,
		Description:    m.Description,
		CouponID:       m.CouponID,
		PointsRequired: m.PointsRequired,
		TotalStock:     m.TotalStock,
		UsedStock:      m.UsedStock,
		PerUserLimit:   m.PerUserLimit,
		Status:         points.RedeemRuleStatus(m.Status),
		StartAt:        m.StartAt,
		EndAt:          m.EndAt,
	}
}

func fromRedeemRuleEntity(r *points.RedeemRule) *redeemRuleModel {
	return &redeemRuleModel{
		Model:          application.Model{ID: r.ID, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt, DeletedAt: r.DeletedAt},
		TenantID:       r.TenantID.Int64(),
		Name:           r.Name,
		Description:    r.Description,
		CouponID:       r.CouponID,
		PointsRequired: r.PointsRequired,
		TotalStock:     r.TotalStock,
		UsedStock:      r.UsedStock,
		PerUserLimit:   r.PerUserLimit,
		Status:         int(r.Status),
		StartAt:        r.StartAt,
		EndAt:          r.EndAt,
	}
}

func (r *redeemRuleRepo) Create(ctx context.Context, db *gorm.DB, rule *points.RedeemRule) error {
	model := fromRedeemRuleEntity(rule)
	return db.WithContext(ctx).Create(model).Error
}

func (r *redeemRuleRepo) Update(ctx context.Context, db *gorm.DB, rule *points.RedeemRule) error {
	model := fromRedeemRuleEntity(rule)
	return db.WithContext(ctx).
		Model(&redeemRuleModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rule.ID, rule.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":            model.Name,
			"description":     model.Description,
			"coupon_id":       model.CouponID,
			"points_required": model.PointsRequired,
			"total_stock":     model.TotalStock,
			"used_stock":      model.UsedStock,
			"per_user_limit":  model.PerUserLimit,
			"status":          model.Status,
			"start_at":        model.StartAt,
			"end_at":          model.EndAt,
			"updated_at":      model.UpdatedAt,
		}).Error
}

func (r *redeemRuleRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrNotFound
	}
	return nil
}

func (r *redeemRuleRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.RedeemRule, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model redeemRuleModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *redeemRuleRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.RedeemRuleQuery) ([]*points.RedeemRule, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []redeemRuleModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	rules := make([]*points.RedeemRule, len(models))
	for i, m := range models {
		rules[i] = m.toEntity()
	}
	return rules, total, nil
}

func (r *redeemRuleRepo) UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status points.RedeemRuleStatus) error {
	query := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	return query.Updates(map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}).Error
}

func (r *redeemRuleRepo) IncrementUsedStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, quantity int64) error {
	query := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	return query.UpdateColumn("used_stock", gorm.Expr("used_stock + ?", quantity)).Error
}

func (r *redeemRuleRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.RedeemRuleStats, error) {
	query := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var total, active int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	activeQuery := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("deleted_at IS NULL AND status = ?", points.RedeemRuleStatusActive)
	if tenantID != 0 {
		activeQuery = activeQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := activeQuery.Count(&active).Error; err != nil {
		return nil, err
	}

	var totalRedeemed int64
	sumQuery := db.WithContext(ctx).Model(&redeemRuleModel{}).Where("deleted_at IS NULL")
	if tenantID != 0 {
		sumQuery = sumQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	sumQuery.Select("COALESCE(SUM(used_stock), 0)").Scan(&totalRedeemed)

	return &points.RedeemRuleStats{
		Total:         total,
		Active:        active,
		TotalRedeemed: totalRedeemed,
	}, nil
}

// ==================== Points Account Repository ====================

type pointsAccountRepo struct{}

func NewPointsAccountRepository() points.PointsAccountRepository {
	return &pointsAccountRepo{}
}

type pointsAccountModel struct {
	application.Model
	TenantID      int64 `gorm:"column:tenant_id;not null;uniqueIndex:uniq_tenant_user"`
	UserID        int64 `gorm:"column:user_id;not null;uniqueIndex:uniq_tenant_user"`
	Balance       int64 `gorm:"column:balance;type:bigint;not null;default:0"`
	FrozenBalance int64 `gorm:"column:frozen_balance;type:bigint;not null;default:0"`
	TotalEarned   int64 `gorm:"column:total_earned;type:bigint;not null;default:0"`
	TotalRedeemed int64 `gorm:"column:total_redeemed;type:bigint;not null;default:0"`
	TotalExpired  int64 `gorm:"column:total_expired;type:bigint;not null;default:0"`
}

func (pointsAccountModel) TableName() string {
	return "points_accounts"
}

func (m *pointsAccountModel) toEntity() *points.PointsAccount {
	return &points.PointsAccount{
		Model:         application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt},
		TenantID:      shared.TenantID(m.TenantID),
		UserID:        m.UserID,
		Balance:       m.Balance,
		FrozenBalance: m.FrozenBalance,
		TotalEarned:   m.TotalEarned,
		TotalRedeemed: m.TotalRedeemed,
		TotalExpired:  m.TotalExpired,
	}
}

func fromPointsAccountEntity(a *points.PointsAccount) *pointsAccountModel {
	return &pointsAccountModel{
		Model:         application.Model{ID: a.ID, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt, DeletedAt: a.DeletedAt},
		TenantID:      a.TenantID.Int64(),
		UserID:        a.UserID,
		Balance:       a.Balance,
		FrozenBalance: a.FrozenBalance,
		TotalEarned:   a.TotalEarned,
		TotalRedeemed: a.TotalRedeemed,
		TotalExpired:  a.TotalExpired,
	}
}

func (r *pointsAccountRepo) Create(ctx context.Context, db *gorm.DB, account *points.PointsAccount) error {
	model := fromPointsAccountEntity(account)
	return db.WithContext(ctx).Create(model).Error
}

func (r *pointsAccountRepo) Update(ctx context.Context, db *gorm.DB, account *points.PointsAccount) error {
	model := fromPointsAccountEntity(account)
	return db.WithContext(ctx).
		Model(&pointsAccountModel{}).
		Where("id = ? AND tenant_id = ?", account.ID, account.TenantID.Int64()).
		Updates(map[string]interface{}{
			"balance":        model.Balance,
			"frozen_balance": model.FrozenBalance,
			"total_earned":   model.TotalEarned,
			"total_redeemed": model.TotalRedeemed,
			"total_expired":  model.TotalExpired,
			"updated_at":     model.UpdatedAt,
		}).Error
}

func (r *pointsAccountRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsAccount, error) {
	query := db.WithContext(ctx)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model pointsAccountModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pointsAccountRepo) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*points.PointsAccount, error) {
	query := db.WithContext(ctx).Where("user_id = ?", userID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model pointsAccountModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pointsAccountRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsAccountQuery) ([]*points.PointsAccount, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&pointsAccountModel{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.UserID > 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}

	// Note: Email filtering would require a join with users table, simplified for now

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []pointsAccountModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	accounts := make([]*points.PointsAccount, len(models))
	for i, m := range models {
		accounts[i] = m.toEntity()
	}
	return accounts, total, nil
}

func (r *pointsAccountRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.PointsAccountStats, error) {
	query := db.WithContext(ctx).Model(&pointsAccountModel{})
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var total, totalBalance, active int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	sumQuery := db.WithContext(ctx).Model(&pointsAccountModel{})
	if tenantID != 0 {
		sumQuery = sumQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	sumQuery.Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance)

	activeQuery := db.WithContext(ctx).Model(&pointsAccountModel{}).Where("balance > 0")
	if tenantID != 0 {
		activeQuery = activeQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := activeQuery.Count(&active).Error; err != nil {
		return nil, err
	}

	return &points.PointsAccountStats{
		Total:        total,
		TotalBalance: totalBalance,
		Active:       active,
	}, nil
}

// ==================== Points Transaction Repository ====================

type pointsTransactionRepo struct{}

func NewPointsTransactionRepository() points.PointsTransactionRepository {
	return &pointsTransactionRepo{}
}

type pointsTransactionModel struct {
	application.Model
	TenantID      int64      `gorm:"column:tenant_id;not null;index:idx_tenant_user"`
	UserID        int64      `gorm:"column:user_id;not null;index:idx_tenant_user"`
	AccountID     int64      `gorm:"column:account_id;not null;index:idx_account_id"`
	Points        int64      `gorm:"column:points;type:bigint;not null"`
	BalanceAfter  int64      `gorm:"column:balance_after;type:bigint;not null"`
	Type          string     `gorm:"column:type;type:varchar(20);not null;index:idx_type"`
	ReferenceType string     `gorm:"column:reference_type;type:varchar(50)"`
	ReferenceID   string     `gorm:"column:reference_id;type:varchar(100);index:idx_reference"`
	Description   string     `gorm:"column:description;type:text"`
	ExpiresAt     *time.Time `gorm:"column:expires_at"`
}

func (pointsTransactionModel) TableName() string {
	return "points_transactions"
}

func (m *pointsTransactionModel) toEntity() *points.PointsTransaction {
	return &points.PointsTransaction{
		Model:         application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt},
		TenantID:      shared.TenantID(m.TenantID),
		UserID:        m.UserID,
		AccountID:     m.AccountID,
		Points:        m.Points,
		BalanceAfter:  m.BalanceAfter,
		Type:          points.TransactionType(m.Type),
		ReferenceType: m.ReferenceType,
		ReferenceID:   m.ReferenceID,
		Description:   m.Description,
		ExpiresAt:     m.ExpiresAt,
	}
}

func fromPointsTransactionEntity(t *points.PointsTransaction) *pointsTransactionModel {
	return &pointsTransactionModel{
		Model:         application.Model{ID: t.ID, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt, DeletedAt: t.DeletedAt},
		TenantID:      t.TenantID.Int64(),
		UserID:        t.UserID,
		AccountID:     t.AccountID,
		Points:        t.Points,
		BalanceAfter:  t.BalanceAfter,
		Type:          string(t.Type),
		ReferenceType: t.ReferenceType,
		ReferenceID:   t.ReferenceID,
		Description:   t.Description,
		ExpiresAt:     t.ExpiresAt,
	}
}

func (r *pointsTransactionRepo) Create(ctx context.Context, db *gorm.DB, transaction *points.PointsTransaction) error {
	model := fromPointsTransactionEntity(transaction)
	return db.WithContext(ctx).Create(model).Error
}

func (r *pointsTransactionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsTransaction, error) {
	query := db.WithContext(ctx)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model pointsTransactionModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pointsTransactionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsTransactionQuery) ([]*points.PointsTransaction, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&pointsTransactionModel{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.UserID > 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.AccountID > 0 {
		dbQuery = dbQuery.Where("account_id = ?", query.AccountID)
	}
	if query.Type != "" && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}
	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", *query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []pointsTransactionModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*points.PointsTransaction, len(models))
	for i, m := range models {
		transactions[i] = m.toEntity()
	}
	return transactions, total, nil
}

func (r *pointsTransactionRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsTransactionQuery) (*points.PointsTransactionStats, error) {
	dbQuery := db.WithContext(ctx).Model(&pointsTransactionModel{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", *query.EndTime)
	}

	var totalEarned, totalRedeemed int64

	earnQuery := db.WithContext(ctx).Model(&pointsTransactionModel{}).Where("type = ?", points.TransactionTypeEarn)
	if tenantID != 0 {
		earnQuery = earnQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if query.StartTime != nil {
		earnQuery = earnQuery.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		earnQuery = earnQuery.Where("created_at < ?", *query.EndTime)
	}
	earnQuery.Select("COALESCE(SUM(points), 0)").Scan(&totalEarned)

	redeemQuery := db.WithContext(ctx).Model(&pointsTransactionModel{}).Where("type = ?", points.TransactionTypeRedeem)
	if tenantID != 0 {
		redeemQuery = redeemQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if query.StartTime != nil {
		redeemQuery = redeemQuery.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		redeemQuery = redeemQuery.Where("created_at < ?", *query.EndTime)
	}
	redeemQuery.Select("COALESCE(SUM(ABS(points)), 0)").Scan(&totalRedeemed)

	return &points.PointsTransactionStats{
		TotalEarned:   totalEarned,
		TotalRedeemed: totalRedeemed,
	}, nil
}

// ==================== Points Redemption Repository ====================

type pointsRedemptionRepo struct{}

func NewPointsRedemptionRepository() points.PointsRedemptionRepository {
	return &pointsRedemptionRepo{}
}

type pointsRedemptionModel struct {
	application.Model
	TenantID     int64      `gorm:"column:tenant_id;not null;index:idx_tenant_user"`
	UserID       int64      `gorm:"column:user_id;not null;index:idx_tenant_user"`
	RedeemRuleID int64      `gorm:"column:redeem_rule_id;not null;index:idx_redeem_rule"`
	CouponID     int64      `gorm:"column:coupon_id;type:bigint;not null"`
	UserCouponID int64      `gorm:"column:user_coupon_id;type:bigint"`
	PointsUsed   int64      `gorm:"column:points_used;type:bigint;not null"`
	Status       int        `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	CompletedAt  *time.Time `gorm:"column:completed_at"`
}

func (pointsRedemptionModel) TableName() string {
	return "points_redemptions"
}

func (m *pointsRedemptionModel) toEntity() *points.PointsRedemption {
	return &points.PointsRedemption{
		Model:        application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt},
		TenantID:     shared.TenantID(m.TenantID),
		UserID:       m.UserID,
		RedeemRuleID: m.RedeemRuleID,
		CouponID:     m.CouponID,
		UserCouponID: m.UserCouponID,
		PointsUsed:   m.PointsUsed,
		Status:       points.RedemptionStatus(m.Status),
		CompletedAt:  m.CompletedAt,
	}
}

func fromPointsRedemptionEntity(r *points.PointsRedemption) *pointsRedemptionModel {
	return &pointsRedemptionModel{
		Model:        application.Model{ID: r.ID, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt, DeletedAt: r.DeletedAt},
		TenantID:     r.TenantID.Int64(),
		UserID:       r.UserID,
		RedeemRuleID: r.RedeemRuleID,
		CouponID:     r.CouponID,
		UserCouponID: r.UserCouponID,
		PointsUsed:   r.PointsUsed,
		Status:       int(r.Status),
		CompletedAt:  r.CompletedAt,
	}
}

func (r *pointsRedemptionRepo) Create(ctx context.Context, db *gorm.DB, redemption *points.PointsRedemption) error {
	model := fromPointsRedemptionEntity(redemption)
	return db.WithContext(ctx).Create(model).Error
}

func (r *pointsRedemptionRepo) Update(ctx context.Context, db *gorm.DB, redemption *points.PointsRedemption) error {
	model := fromPointsRedemptionEntity(redemption)
	return db.WithContext(ctx).
		Model(&pointsRedemptionModel{}).
		Where("id = ? AND tenant_id = ?", redemption.ID, redemption.TenantID.Int64()).
		Updates(map[string]interface{}{
			"user_coupon_id": model.UserCouponID,
			"status":         model.Status,
			"completed_at":   model.CompletedAt,
			"updated_at":     model.UpdatedAt,
		}).Error
}

func (r *pointsRedemptionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsRedemption, error) {
	query := db.WithContext(ctx)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model pointsRedemptionModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pointsRedemptionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsRedemptionQuery) ([]*points.PointsRedemption, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&pointsRedemptionModel{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.UserID > 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.Status >= points.RedemptionStatusPending && query.Status <= points.RedemptionStatusCancelled {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", *query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []pointsRedemptionModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	redemptions := make([]*points.PointsRedemption, len(models))
	for i, m := range models {
		redemptions[i] = m.toEntity()
	}
	return redemptions, total, nil
}

func (r *pointsRedemptionRepo) CountByUserAndRule(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID, ruleID int64) (int64, error) {
	query := db.WithContext(ctx).Model(&pointsRedemptionModel{}).
		Where("user_id = ? AND redeem_rule_id = ?", userID, ruleID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
