package points

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== Status Constants ====================

// EarnRuleStatus 积分获取规则状态
type EarnRuleStatus int

const (
	EarnRuleStatusDraft EarnRuleStatus = iota
	EarnRuleStatusActive
	EarnRuleStatusInactive
)

func (s EarnRuleStatus) String() string {
	switch s {
	case EarnRuleStatusDraft:
		return "draft"
	case EarnRuleStatusActive:
		return "active"
	case EarnRuleStatusInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

func (s EarnRuleStatus) IsValid() bool {
	return s >= EarnRuleStatusDraft && s <= EarnRuleStatusInactive
}

// RedeemRuleStatus 积分兑换规则状态
type RedeemRuleStatus int

const (
	RedeemRuleStatusInactive RedeemRuleStatus = iota
	RedeemRuleStatusActive
)

func (s RedeemRuleStatus) String() string {
	switch s {
	case RedeemRuleStatusInactive:
		return "inactive"
	case RedeemRuleStatusActive:
		return "active"
	default:
		return "unknown"
	}
}

func (s RedeemRuleStatus) IsValid() bool {
	return s >= RedeemRuleStatusInactive && s <= RedeemRuleStatusActive
}

// RedemptionStatus 兑换记录状态
type RedemptionStatus int

const (
	RedemptionStatusPending RedemptionStatus = iota
	RedemptionStatusCompleted
	RedemptionStatusCancelled
)

func (s RedemptionStatus) String() string {
	switch s {
	case RedemptionStatusPending:
		return "pending"
	case RedemptionStatusCompleted:
		return "completed"
	case RedemptionStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

func (s RedemptionStatus) IsValid() bool {
	return s >= RedemptionStatusPending && s <= RedemptionStatusCancelled
}

// ==================== Scenario Constants ====================

// EarnScenario 积分获取场景
type EarnScenario string

const (
	EarnScenarioOrderPayment  EarnScenario = "ORDER_PAYMENT"
	EarnScenarioSignIn        EarnScenario = "SIGN_IN"
	EarnScenarioProductReview EarnScenario = "PRODUCT_REVIEW"
	EarnScenarioFirstOrder    EarnScenario = "FIRST_ORDER"
)

func (s EarnScenario) IsValid() bool {
	switch s {
	case EarnScenarioOrderPayment, EarnScenarioSignIn, EarnScenarioProductReview, EarnScenarioFirstOrder:
		return true
	default:
		return false
	}
}

// ==================== Calculation Type Constants ====================

// CalculationType 积分计算类型
type CalculationType string

const (
	CalculationTypeFixed  CalculationType = "FIXED"
	CalculationTypeRatio  CalculationType = "RATIO"
	CalculationTypeTiered CalculationType = "TIERED"
)

func (c CalculationType) IsValid() bool {
	switch c {
	case CalculationTypeFixed, CalculationTypeRatio, CalculationTypeTiered:
		return true
	default:
		return false
	}
}

// ==================== Condition Type Constants ====================

// ConditionType 条件类型
type ConditionType string

const (
	ConditionTypeNone             ConditionType = "NONE"
	ConditionTypeNewUser          ConditionType = "NEW_USER"
	ConditionTypeFirstOrder       ConditionType = "FIRST_ORDER"
	ConditionTypeSpecificProducts ConditionType = "SPECIFIC_PRODUCTS"
	ConditionTypeMinAmount        ConditionType = "MIN_AMOUNT"
)

func (c ConditionType) IsValid() bool {
	switch c {
	case ConditionTypeNone, ConditionTypeNewUser, ConditionTypeFirstOrder, ConditionTypeSpecificProducts, ConditionTypeMinAmount:
		return true
	default:
		return false
	}
}

// ==================== Transaction Type Constants ====================

// TransactionType 积分交易类型
type TransactionType string

const (
	TransactionTypeEarn     TransactionType = "EARN"
	TransactionTypeRedeem   TransactionType = "REDEEM"
	TransactionTypeAdjust   TransactionType = "ADJUST"
	TransactionTypeExpire   TransactionType = "EXPIRE"
	TransactionTypeFreeze   TransactionType = "FREEZE"
	TransactionTypeUnfreeze TransactionType = "UNFREEZE"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeEarn, TransactionTypeRedeem, TransactionTypeAdjust, TransactionTypeExpire, TransactionTypeFreeze, TransactionTypeUnfreeze:
		return true
	default:
		return false
	}
}

// ==================== Adjustment Type Constants ====================

// AdjustmentType 调整类型
type AdjustmentType string

const (
	AdjustmentTypeAdd    AdjustmentType = "ADD"
	AdjustmentTypeDeduct AdjustmentType = "DEDUCT"
)

func (a AdjustmentType) IsValid() bool {
	switch a {
	case AdjustmentTypeAdd, AdjustmentTypeDeduct:
		return true
	default:
		return false
	}
}

// ==================== TierConfig 阶梯配置 ====================

// TierConfig 阶梯配置
type TierConfig struct {
	Threshold *int64         `json:"threshold"` // null for last tier (unlimited)
	Ratio     decimal.Decimal `json:"ratio"`     // Points per currency unit
}

// TierConfigs 阶梯配置列表
type TierConfigs []TierConfig

// Scan implements sql.Scanner interface
func (t *TierConfigs) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into TierConfigs", value)
	}
	return json.Unmarshal(bytes, t)
}

// Value implements driver.Valuer interface
func (t TierConfigs) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// ==================== EarnRule 积分获取规则 ====================

// EarnRule 积分获取规则
type EarnRule struct {
	ID               int64           `gorm:"column:id;primaryKey"`
	TenantID         shared.TenantID `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id"`
	Name             string          `gorm:"column:name;type:varchar(255);not null"`
	Description      string          `gorm:"column:description;type:text"`
	Scenario         EarnScenario    `gorm:"column:scenario;type:varchar(50);not null"`
	CalculationType  CalculationType `gorm:"column:calculation_type;type:varchar(20);not null"`
	FixedPoints      int64           `gorm:"column:fixed_points;type:bigint;default:0"`
	Ratio            decimal.Decimal `gorm:"column:ratio;type:decimal(10,4)"`
	Tiers            TierConfigs     `gorm:"column:tiers;type:json"`
	ConditionType    ConditionType   `gorm:"column:condition_type;type:varchar(50);not null;default:'NONE'"`
	ConditionValue   string          `gorm:"column:condition_value;type:text"`
	ExpirationMonths int             `gorm:"column:expiration_months;type:int;not null;default:12"`
	Status           EarnRuleStatus  `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	Priority         int             `gorm:"column:priority;type:int;not null;default:0"`
	StartAt          *int64          `gorm:"column:start_at"`
	EndAt            *int64          `gorm:"column:end_at"`
	DeletedAt        *int64          `gorm:"column:deleted_at;index"`
	Audit            shared.AuditInfo `gorm:"embedded"`
}

func (e *EarnRule) TableName() string {
	return "earn_rules"
}

// IsActive 检查规则是否处于激活状态
func (e *EarnRule) IsActive() bool {
	if e.Status != EarnRuleStatusActive {
		return false
	}

	now := time.Now().UTC().UnixMilli()

	// 检查开始时间
	if e.StartAt != nil && now < *e.StartAt {
		return false
	}

	// 检查结束时间
	if e.EndAt != nil && now > *e.EndAt {
		return false
	}

	return true
}

// CalculatePoints 计算积分数量
func (e *EarnRule) CalculatePoints(orderAmount decimal.Decimal) (int64, error) {
	if !e.IsActive() {
		return 0, code.ErrPromotionNotActive
	}

	switch e.CalculationType {
	case CalculationTypeFixed:
		return e.FixedPoints, nil

	case CalculationTypeRatio:
		if e.Ratio.IsZero() {
			return 0, nil
		}
		points := orderAmount.Mul(e.Ratio)
		return points.IntPart(), nil

	case CalculationTypeTiered:
		return e.calculateTieredPoints(orderAmount), nil

	default:
		return 0, code.ErrParam
	}
}

// calculateTieredPoints 计算阶梯积分
func (e *EarnRule) calculateTieredPoints(orderAmount decimal.Decimal) int64 {
	if len(e.Tiers) == 0 {
		return 0
	}

	amountInt := orderAmount.IntPart()

	// 找到匹配的阶梯
	for i, tier := range e.Tiers {
		// 最后一个阶梯或无阈值限制
		if tier.Threshold == nil {
			return amountInt * tier.Ratio.IntPart()
		}

		// 当前阶梯匹配
		if amountInt < *tier.Threshold {
			// 如果是第一个阶梯，不匹配
			if i == 0 {
				return 0
			}
			// 使用前一个阶梯的比例
			return amountInt * e.Tiers[i-1].Ratio.IntPart()
		}
	}

	// 超过所有阶梯，使用最后一个阶梯
	lastTier := e.Tiers[len(e.Tiers)-1]
	return amountInt * lastTier.Ratio.IntPart()
}

// Activate 激活规则
func (e *EarnRule) Activate(updatedBy int64) error {
	if e.Status == EarnRuleStatusActive {
		return code.ErrParam
	}
	e.Status = EarnRuleStatusActive
	e.Audit.Update(updatedBy)
	return nil
}

// Deactivate 停用规则
func (e *EarnRule) Deactivate(updatedBy int64) error {
	if e.Status == EarnRuleStatusInactive {
		return code.ErrParam
	}
	e.Status = EarnRuleStatusInactive
	e.Audit.Update(updatedBy)
	return nil
}

// ==================== RedeemRule 积分兑换规则 ====================

// RedeemRule 积分兑换规则
type RedeemRule struct {
	ID             int64            `gorm:"column:id;primaryKey"`
	TenantID       shared.TenantID  `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_id"`
	Name           string           `gorm:"column:name;type:varchar(255);not null"`
	Description    string           `gorm:"column:description;type:text"`
	CouponID       int64            `gorm:"column:coupon_id;type:bigint;not null"`
	PointsRequired int64            `gorm:"column:points_required;type:bigint;not null"`
	TotalStock     int64            `gorm:"column:total_stock;type:bigint;not null;default:0"`
	UsedStock      int64            `gorm:"column:used_stock;type:bigint;not null;default:0"`
	PerUserLimit   int              `gorm:"column:per_user_limit;type:int;not null;default:1"`
	Status         RedeemRuleStatus `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	StartAt        *int64          `gorm:"column:start_at"`
	EndAt          *int64          `gorm:"column:end_at"`
	DeletedAt      *int64           `gorm:"column:deleted_at;index"`
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (r *RedeemRule) TableName() string {
	return "redeem_rules"
}

// IsActive 检查规则是否处于激活状态
func (r *RedeemRule) IsActive() bool {
	if r.Status != RedeemRuleStatusActive {
		return false
	}

	now := time.Now().UTC().UnixMilli()

	if r.StartAt != nil && now < *r.StartAt {
		return false
	}

	if r.EndAt != nil && now > *r.EndAt {
		return false
	}

	return true
}

// HasStock 检查是否有库存
func (r *RedeemRule) HasStock() bool {
	return r.UsedStock < r.TotalStock
}

// RemainingStock 获取剩余库存
func (r *RedeemRule) RemainingStock() int64 {
	remaining := r.TotalStock - r.UsedStock
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Activate 激活规则
func (r *RedeemRule) Activate(updatedBy int64) error {
	if r.Status == RedeemRuleStatusActive {
		return code.ErrParam
	}
	r.Status = RedeemRuleStatusActive
	r.Audit.Update(updatedBy)
	return nil
}

// Deactivate 停用规则
func (r *RedeemRule) Deactivate(updatedBy int64) error {
	if r.Status == RedeemRuleStatusInactive {
		return code.ErrParam
	}
	r.Status = RedeemRuleStatusInactive
	r.Audit.Update(updatedBy)
	return nil
}

// ==================== PointsAccount 积分账户 ====================

// PointsAccount 积分账户
type PointsAccount struct {
	ID            int64            `gorm:"column:id;primaryKey"`
	TenantID      shared.TenantID  `gorm:"column:tenant_id;type:bigint;not null;uniqueIndex:uniq_tenant_user"`
	UserID        int64            `gorm:"column:user_id;type:bigint;not null;uniqueIndex:uniq_tenant_user"`
	Balance       int64            `gorm:"column:balance;type:bigint;not null;default:0"`
	FrozenBalance int64            `gorm:"column:frozen_balance;type:bigint;not null;default:0"`
	TotalEarned   int64            `gorm:"column:total_earned;type:bigint;not null;default:0"`
	TotalRedeemed int64            `gorm:"column:total_redeemed;type:bigint;not null;default:0"`
	TotalExpired  int64            `gorm:"column:total_expired;type:bigint;not null;default:0"`
	DeletedAt     *int64           `gorm:"column:deleted_at;index"`
	Audit         shared.AuditInfo `gorm:"embedded"`
}

func (a *PointsAccount) TableName() string {
	return "points_accounts"
}

// AvailableBalance 获取可用余额
func (a *PointsAccount) AvailableBalance() int64 {
	return a.Balance - a.FrozenBalance
}

// AddPoints 增加积分
func (a *PointsAccount) AddPoints(points int64, updatedBy int64) error {
	if points <= 0 {
		return code.ErrParam
	}
	a.Balance += points
	a.TotalEarned += points
	a.Audit.Update(updatedBy)
	return nil
}

// DeductPoints 扣减积分
func (a *PointsAccount) DeductPoints(points int64, updatedBy int64) error {
	if points <= 0 {
		return code.ErrParam
	}
	if a.AvailableBalance() < points {
		return code.ErrSharedInsufficientAmount
	}
	a.Balance -= points
	a.TotalRedeemed += points
	a.Audit.Update(updatedBy)
	return nil
}

// Freeze 冻结积分
func (a *PointsAccount) Freeze(points int64, updatedBy int64) error {
	if points <= 0 {
		return code.ErrParam
	}
	available := a.AvailableBalance()
	if available < points {
		return code.ErrSharedInsufficientAmount
	}
	a.FrozenBalance += points
	a.Audit.Update(updatedBy)
	return nil
}

// Unfreeze 解冻积分
func (a *PointsAccount) Unfreeze(points int64, updatedBy int64) error {
	if points <= 0 {
		return code.ErrParam
	}
	if a.FrozenBalance < points {
		return code.ErrSharedInsufficientAmount
	}
	a.FrozenBalance -= points
	a.Audit.Update(updatedBy)
	return nil
}

// Expire 过期积分
func (a *PointsAccount) Expire(points int64, updatedBy int64) error {
	if points <= 0 {
		return code.ErrParam
	}
	if a.Balance < points {
		return code.ErrSharedInsufficientAmount
	}
	a.Balance -= points
	a.TotalExpired += points
	a.Audit.Update(updatedBy)
	return nil
}

// ==================== PointsTransaction 积分交易记录 ====================

// PointsTransaction 积分交易记录
type PointsTransaction struct {
	ID            int64            `gorm:"column:id;primaryKey"`
	TenantID      shared.TenantID  `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_user"`
	UserID        int64            `gorm:"column:user_id;type:bigint;not null;index:idx_tenant_user"`
	AccountID     int64            `gorm:"column:account_id;type:bigint;not null;index:idx_account_id"`
	Points        int64            `gorm:"column:points;type:bigint;not null"` // Positive = earn, negative = deduct
	BalanceAfter  int64            `gorm:"column:balance_after;type:bigint;not null"`
	Type          TransactionType  `gorm:"column:type;type:varchar(20);not null;index:idx_type"`
	ReferenceType string           `gorm:"column:reference_type;type:varchar(50)"`
	ReferenceID   string           `gorm:"column:reference_id;type:varchar(100);index:idx_reference"`
	Description   string           `gorm:"column:description;type:text"`
	ExpiresAt     *int64          `gorm:"column:expires_at;index:idx_expires_at"`
	DeletedAt     *int64           `gorm:"column:deleted_at;index"`
	Audit         shared.AuditInfo `gorm:"embedded"`
}

func (t *PointsTransaction) TableName() string {
	return "points_transactions"
}

// IsEarn 是否为获取交易
func (t *PointsTransaction) IsEarn() bool {
	return t.Points > 0
}

// IsDeduction 是否为扣减交易
func (t *PointsTransaction) IsDeduction() bool {
	return t.Points < 0
}

// ==================== PointsRedemption 积分兑换记录 ====================

// PointsRedemption 积分兑换记录
type PointsRedemption struct {
	ID           int64            `gorm:"column:id;primaryKey"`
	TenantID     shared.TenantID  `gorm:"column:tenant_id;type:bigint;not null;index:idx_tenant_user"`
	UserID       int64            `gorm:"column:user_id;type:bigint;not null;index:idx_tenant_user"`
	RedeemRuleID int64            `gorm:"column:redeem_rule_id;type:bigint;not null;index:idx_redeem_rule"`
	CouponID     int64            `gorm:"column:coupon_id;type:bigint;not null"`
	UserCouponID int64            `gorm:"column:user_coupon_id;type:bigint"`
	PointsUsed   int64            `gorm:"column:points_used;type:bigint;not null"`
	Status       RedemptionStatus `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status"`
	CompletedAt  *int64          `gorm:"column:completed_at"`
	DeletedAt    *int64           `gorm:"column:deleted_at;index"`
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (r *PointsRedemption) TableName() string {
	return "points_redemptions"
}

// Complete 完成兑换
func (r *PointsRedemption) Complete(userCouponID int64, updatedBy int64) error {
	if r.Status != RedemptionStatusPending {
		return code.ErrParam
	}
	r.Status = RedemptionStatusCompleted
	r.UserCouponID = userCouponID
	now := time.Now().UTC().UnixMilli()
	r.CompletedAt = &now
	r.Audit.Update(updatedBy)
	return nil
}

// Cancel 取消兑换
func (r *PointsRedemption) Cancel(updatedBy int64) error {
	if r.Status != RedemptionStatusPending {
		return code.ErrParam
	}
	r.Status = RedemptionStatusCancelled
	r.Audit.Update(updatedBy)
	return nil
}

// ==================== Repository Interfaces ====================

// EarnRuleQuery 积分获取规则查询条件
type EarnRuleQuery struct {
	shared.PageQuery
	Name            string
	Status          EarnRuleStatus
	Scenario        EarnScenario
	CalculationType CalculationType
}

// EarnRuleRepository 积分获取规则仓储接口
type EarnRuleRepository interface {
	Create(ctx context.Context, db *gorm.DB, rule *EarnRule) error
	Update(ctx context.Context, db *gorm.DB, rule *EarnRule) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*EarnRule, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query EarnRuleQuery) ([]*EarnRule, int64, error)
	FindByScenario(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, scenario EarnScenario) ([]*EarnRule, error)
	UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status EarnRuleStatus) error
	GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*EarnRuleStats, error)
}

// EarnRuleStats 积分获取规则统计
type EarnRuleStats struct {
	Total  int64
	Active int64
}

// RedeemRuleQuery 积分兑换规则查询条件
type RedeemRuleQuery struct {
	shared.PageQuery
	Name   string
	Status RedeemRuleStatus
}

// RedeemRuleRepository 积分兑换规则仓储接口
type RedeemRuleRepository interface {
	Create(ctx context.Context, db *gorm.DB, rule *RedeemRule) error
	Update(ctx context.Context, db *gorm.DB, rule *RedeemRule) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*RedeemRule, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query RedeemRuleQuery) ([]*RedeemRule, int64, error)
	UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status RedeemRuleStatus) error
	IncrementUsedStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, quantity int64) error
	GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*RedeemRuleStats, error)
}

// RedeemRuleStats 积分兑换规则统计
type RedeemRuleStats struct {
	Total         int64
	Active        int64
	TotalRedeemed int64
}

// PointsAccountQuery 积分账户查询条件
type PointsAccountQuery struct {
	shared.PageQuery
	UserID int64
	Email  string
}

// PointsAccountRepository 积分账户仓储接口
type PointsAccountRepository interface {
	Create(ctx context.Context, db *gorm.DB, account *PointsAccount) error
	Update(ctx context.Context, db *gorm.DB, account *PointsAccount) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*PointsAccount, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*PointsAccount, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query PointsAccountQuery) ([]*PointsAccount, int64, error)
	GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*PointsAccountStats, error)
}

// PointsAccountStats 积分账户统计
type PointsAccountStats struct {
	Total        int64
	TotalBalance int64
	Active       int64
}

// PointsTransactionQuery 积分交易查询条件
type PointsTransactionQuery struct {
	shared.PageQuery
	UserID    int64
	AccountID int64
	Type      TransactionType
	StartTime *time.Time
	EndTime   *time.Time
}

// PointsTransactionRepository 积分交易仓储接口
type PointsTransactionRepository interface {
	Create(ctx context.Context, db *gorm.DB, transaction *PointsTransaction) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*PointsTransaction, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query PointsTransactionQuery) ([]*PointsTransaction, int64, error)
	GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query PointsTransactionQuery) (*PointsTransactionStats, error)
}

// PointsTransactionStats 积分交易统计
type PointsTransactionStats struct {
	TotalEarned   int64
	TotalRedeemed int64
}

// PointsRedemptionQuery 积分兑换记录查询条件
type PointsRedemptionQuery struct {
	shared.PageQuery
	UserID    int64
	Status    RedemptionStatus
	StartTime *time.Time
	EndTime   *time.Time
}

// PointsRedemptionRepository 积分兑换记录仓储接口
type PointsRedemptionRepository interface {
	Create(ctx context.Context, db *gorm.DB, redemption *PointsRedemption) error
	Update(ctx context.Context, db *gorm.DB, redemption *PointsRedemption) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*PointsRedemption, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query PointsRedemptionQuery) ([]*PointsRedemption, int64, error)
	CountByUserAndRule(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID, ruleID int64) (int64, error)
}