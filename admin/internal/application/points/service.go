package points

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== DTO Types ====================

// EarnRuleDTO represents earn rule response
type EarnRuleDTO struct {
	ID               int64              `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Scenario         string             `json:"scenario"`
	CalculationType  string             `json:"calculation_type"`
	FixedPoints      int64              `json:"fixed_points"`
	Ratio            decimal.Decimal    `json:"ratio"`
	Tiers            points.TierConfigs `json:"tiers,omitempty"`
	ConditionType    string             `json:"condition_type"`
	ConditionValue   string             `json:"condition_value"`
	ExpirationMonths int                `json:"expiration_months"`
	Status           string             `json:"status"`
	Priority         int                `json:"priority"`
	StartAt          *time.Time         `json:"start_at,omitempty"`
	EndAt            *time.Time         `json:"end_at,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// RedeemRuleDTO represents redeem rule response
type RedeemRuleDTO struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CouponID       int64      `json:"coupon_id"`
	PointsRequired int64      `json:"points_required"`
	TotalStock     int64      `json:"total_stock"`
	UsedStock      int64      `json:"used_stock"`
	PerUserLimit   int        `json:"per_user_limit"`
	Status         string     `json:"status"`
	StartAt        *time.Time `json:"start_at,omitempty"`
	EndAt          *time.Time `json:"end_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// PointsAccountDTO represents points account response
type PointsAccountDTO struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Balance       int64     `json:"balance"`
	FrozenBalance int64     `json:"frozen_balance"`
	TotalEarned   int64     `json:"total_earned"`
	TotalRedeemed int64     `json:"total_redeemed"`
	TotalExpired  int64     `json:"total_expired"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PointsTransactionDTO represents transaction response
type PointsTransactionDTO struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	AccountID     int64      `json:"account_id"`
	Points        int64      `json:"points"`
	BalanceAfter  int64      `json:"balance_after"`
	Type          string     `json:"type"`
	ReferenceType string     `json:"reference_type"`
	ReferenceID   string     `json:"reference_id"`
	Description   string     `json:"description"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// PointsRedemptionDTO represents redemption response
type PointsRedemptionDTO struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	RedeemRuleID int64      `json:"redeem_rule_id"`
	CouponID     int64      `json:"coupon_id"`
	UserCouponID int64      `json:"user_coupon_id,omitempty"`
	PointsUsed   int64      `json:"points_used"`
	Status       string     `json:"status"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ==================== Request Types ====================

// CreateEarnRuleRequest create earn rule request
type CreateEarnRuleRequest struct {
	Name             string
	Description      string
	Scenario         points.EarnScenario
	CalculationType  points.CalculationType
	FixedPoints      int64
	Ratio            decimal.Decimal
	Tiers            points.TierConfigs
	ConditionType    points.ConditionType
	ConditionValue   string
	ExpirationMonths int
	Priority         int
	StartAt          *time.Time
	EndAt            *time.Time
}

// UpdateEarnRuleRequest update earn rule request
type UpdateEarnRuleRequest struct {
	ID               int64
	Name             string
	Description      string
	Scenario         points.EarnScenario
	CalculationType  points.CalculationType
	FixedPoints      int64
	Ratio            decimal.Decimal
	Tiers            points.TierConfigs
	ConditionType    points.ConditionType
	ConditionValue   string
	ExpirationMonths int
	Priority         int
	StartAt          *time.Time
	EndAt            *time.Time
}

// CreateRedeemRuleRequest create redeem rule request
type CreateRedeemRuleRequest struct {
	Name           string
	Description    string
	CouponID       int64
	PointsRequired int64
	TotalStock     int64
	PerUserLimit   int
	StartAt        *time.Time
	EndAt          *time.Time
}

// UpdateRedeemRuleRequest update redeem rule request
type UpdateRedeemRuleRequest struct {
	ID             int64
	Name           string
	Description    string
	CouponID       int64
	PointsRequired int64
	TotalStock     int64
	PerUserLimit   int
	StartAt        *time.Time
	EndAt          *time.Time
}

// AdjustPointsRequest adjust points request
type AdjustPointsRequest struct {
	AccountID      int64
	AdjustmentType points.AdjustmentType
	Points         int64
	Reason         string
	OperatorID     int64
}

// ==================== Service Interface ====================

// Service points application service interface
type Service interface {
	// Earn Rules
	CreateEarnRule(ctx context.Context, tenantID shared.TenantID, req CreateEarnRuleRequest, operatorID int64) (*EarnRuleDTO, error)
	UpdateEarnRule(ctx context.Context, tenantID shared.TenantID, req UpdateEarnRuleRequest, operatorID int64) (*EarnRuleDTO, error)
	DeleteEarnRule(ctx context.Context, tenantID shared.TenantID, id int64) error
	GetEarnRule(ctx context.Context, tenantID shared.TenantID, id int64) (*EarnRuleDTO, error)
	ListEarnRules(ctx context.Context, tenantID shared.TenantID, query points.EarnRuleQuery) ([]*EarnRuleDTO, int64, *points.EarnRuleStats, error)
	ActivateEarnRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error
	DeactivateEarnRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error

	// Redeem Rules
	CreateRedeemRule(ctx context.Context, tenantID shared.TenantID, req CreateRedeemRuleRequest, operatorID int64) (*RedeemRuleDTO, error)
	UpdateRedeemRule(ctx context.Context, tenantID shared.TenantID, req UpdateRedeemRuleRequest, operatorID int64) (*RedeemRuleDTO, error)
	DeleteRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64) error
	GetRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64) (*RedeemRuleDTO, error)
	ListRedeemRules(ctx context.Context, tenantID shared.TenantID, query points.RedeemRuleQuery) ([]*RedeemRuleDTO, int64, *points.RedeemRuleStats, error)
	ActivateRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error
	DeactivateRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error

	// Accounts
	GetAccount(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsAccountDTO, error)
	GetAccountByUser(ctx context.Context, tenantID shared.TenantID, userID int64) (*PointsAccountDTO, error)
	ListAccounts(ctx context.Context, tenantID shared.TenantID, query points.PointsAccountQuery) ([]*PointsAccountDTO, int64, *points.PointsAccountStats, error)
	AdjustPoints(ctx context.Context, tenantID shared.TenantID, req AdjustPointsRequest) (*PointsTransactionDTO, error)
	GetAccountTransactions(ctx context.Context, tenantID shared.TenantID, accountID int64, query points.PointsTransactionQuery) ([]*PointsTransactionDTO, int64, error)

	// Transactions
	GetTransaction(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsTransactionDTO, error)
	ListTransactions(ctx context.Context, tenantID shared.TenantID, query points.PointsTransactionQuery) ([]*PointsTransactionDTO, int64, *points.PointsTransactionStats, error)

	// Redemptions
	GetRedemption(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsRedemptionDTO, error)
	ListRedemptions(ctx context.Context, tenantID shared.TenantID, query points.PointsRedemptionQuery) ([]*PointsRedemptionDTO, int64, error)

	// Statistics
	GetStats(ctx context.Context, tenantID shared.TenantID, startTime, endTime *time.Time) (*PointsStatsDTO, error)
	GetTrend(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time, granularity string) ([]TrendDataPoint, error)
	GetTopUsers(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time, limit int) ([]TopUserDTO, error)
	GetExpiringPoints(ctx context.Context, tenantID shared.TenantID, days int) ([]ExpiringPointsDTO, int64, error)
}

// PointsStatsDTO represents points statistics
type PointsStatsDTO struct {
	TotalIssued        int64   `json:"total_issued"`
	TotalRedeemed      int64   `json:"total_redeemed"`
	TotalExpired       int64   `json:"total_expired"`
	OutstandingBalance int64   `json:"outstanding_balance"`
	RedemptionRate     string  `json:"redemption_rate"`
	ActiveUsers        int64   `json:"active_users"`
	PeriodStart        string  `json:"period_start"`
	PeriodEnd          string  `json:"period_end"`
}

// TrendDataPoint represents trend data point
type TrendDataPoint struct {
	Date     string `json:"date"`
	Earned   int64  `json:"earned"`
	Redeemed int64  `json:"redeemed"`
	Expired  int64  `json:"expired"`
}

// TopUserDTO represents top user data
type TopUserDTO struct {
	UserID       int64 `json:"user_id"`
	PointsEarned int64 `json:"points_earned"`
}

// ExpiringPointsDTO represents expiring points data
type ExpiringPointsDTO struct {
	Date      string `json:"date"`
	Points    int64  `json:"points"`
	UserCount int64  `json:"user_count"`
}

// ==================== Service Implementation ====================

type service struct {
	db              *gorm.DB
	earnRuleRepo    points.EarnRuleRepository
	redeemRuleRepo  points.RedeemRuleRepository
	accountRepo     points.PointsAccountRepository
	transactionRepo points.PointsTransactionRepository
	redemptionRepo  points.PointsRedemptionRepository
	idGen           snowflake.Snowflake
}

// NewService creates a new points service
func NewService(
	db *gorm.DB,
	earnRuleRepo points.EarnRuleRepository,
	redeemRuleRepo points.RedeemRuleRepository,
	accountRepo points.PointsAccountRepository,
	transactionRepo points.PointsTransactionRepository,
	redemptionRepo points.PointsRedemptionRepository,
	idGen snowflake.Snowflake,
) Service {
	return &service{
		db:              db,
		earnRuleRepo:    earnRuleRepo,
		redeemRuleRepo:  redeemRuleRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		redemptionRepo:  redemptionRepo,
		idGen:           idGen,
	}
}

// ==================== Earn Rules ====================

func (s *service) CreateEarnRule(ctx context.Context, tenantID shared.TenantID, req CreateEarnRuleRequest, operatorID int64) (*EarnRuleDTO, error) {
	if req.Name == "" {
		return nil, code.ErrParam
	}
	if !req.Scenario.IsValid() {
		return nil, code.ErrParam
	}
	if !req.CalculationType.IsValid() {
		return nil, code.ErrParam
	}

	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	rule := &points.EarnRule{
		ID:               id,
		TenantID:         tenantID,
		Name:             req.Name,
		Description:      req.Description,
		Scenario:         req.Scenario,
		CalculationType:  req.CalculationType,
		FixedPoints:      req.FixedPoints,
		Ratio:            req.Ratio,
		Tiers:            req.Tiers,
		ConditionType:    req.ConditionType,
		ConditionValue:   req.ConditionValue,
		ExpirationMonths: req.ExpirationMonths,
		Status:           points.EarnRuleStatusDraft,
		Priority:         req.Priority,
		StartAt:          timeToInt64(req.StartAt),
		EndAt:            timeToInt64(req.EndAt),
		Audit:            shared.NewAuditInfo(operatorID),
	}

	if err := s.earnRuleRepo.Create(ctx, s.db, rule); err != nil {
		return nil, err
	}

	return toEarnRuleDTO(rule), nil
}

func (s *service) UpdateEarnRule(ctx context.Context, tenantID shared.TenantID, req UpdateEarnRuleRequest, operatorID int64) (*EarnRuleDTO, error) {
	rule, err := s.earnRuleRepo.FindByID(ctx, s.db, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if rule.Status == points.EarnRuleStatusActive {
		return nil, code.ErrParam
	}

	rule.Name = req.Name
	rule.Description = req.Description
	rule.Scenario = req.Scenario
	rule.CalculationType = req.CalculationType
	rule.FixedPoints = req.FixedPoints
	rule.Ratio = req.Ratio
	rule.Tiers = req.Tiers
	rule.ConditionType = req.ConditionType
	rule.ConditionValue = req.ConditionValue
	rule.ExpirationMonths = req.ExpirationMonths
	rule.Priority = req.Priority
	rule.StartAt = timeToInt64(req.StartAt)
	rule.EndAt = timeToInt64(req.EndAt)
	rule.Audit.Update(operatorID)

	if err := s.earnRuleRepo.Update(ctx, s.db, rule); err != nil {
		return nil, err
	}

	return toEarnRuleDTO(rule), nil
}

func (s *service) DeleteEarnRule(ctx context.Context, tenantID shared.TenantID, id int64) error {
	rule, err := s.earnRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if rule.Status == points.EarnRuleStatusActive {
		return nil // Cannot delete active rule
	}

	return s.earnRuleRepo.Delete(ctx, s.db, tenantID, id)
}

func (s *service) GetEarnRule(ctx context.Context, tenantID shared.TenantID, id int64) (*EarnRuleDTO, error) {
	rule, err := s.earnRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toEarnRuleDTO(rule), nil
}

func (s *service) ListEarnRules(ctx context.Context, tenantID shared.TenantID, query points.EarnRuleQuery) ([]*EarnRuleDTO, int64, *points.EarnRuleStats, error) {
	rules, total, err := s.earnRuleRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, nil, err
	}

	stats, err := s.earnRuleRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, 0, nil, err
	}

	dtos := make([]*EarnRuleDTO, len(rules))
	for i, r := range rules {
		dtos[i] = toEarnRuleDTO(r)
	}

	return dtos, total, stats, nil
}

func (s *service) ActivateEarnRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error {
	rule, err := s.earnRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rule.Activate(operatorID); err != nil {
		return err
	}

	return s.earnRuleRepo.Update(ctx, s.db, rule)
}

func (s *service) DeactivateEarnRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error {
	rule, err := s.earnRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rule.Deactivate(operatorID); err != nil {
		return err
	}

	return s.earnRuleRepo.Update(ctx, s.db, rule)
}

// ==================== Redeem Rules ====================

func (s *service) CreateRedeemRule(ctx context.Context, tenantID shared.TenantID, req CreateRedeemRuleRequest, operatorID int64) (*RedeemRuleDTO, error) {
	if req.Name == "" {
		return nil, code.ErrParam
	}
	if req.CouponID == 0 {
		return nil, code.ErrParam
	}
	if req.PointsRequired <= 0 {
		return nil, code.ErrParam
	}

	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	rule := &points.RedeemRule{
		ID:             id,
		TenantID:       tenantID,
		Name:           req.Name,
		Description:    req.Description,
		CouponID:       req.CouponID,
		PointsRequired: req.PointsRequired,
		TotalStock:     req.TotalStock,
		UsedStock:      0,
		PerUserLimit:   req.PerUserLimit,
		Status:         points.RedeemRuleStatusInactive,
		StartAt:        timeToInt64(req.StartAt),
		EndAt:          timeToInt64(req.EndAt),
		Audit:          shared.NewAuditInfo(operatorID),
	}

	if err := s.redeemRuleRepo.Create(ctx, s.db, rule); err != nil {
		return nil, err
	}

	return toRedeemRuleDTO(rule), nil
}

func (s *service) UpdateRedeemRule(ctx context.Context, tenantID shared.TenantID, req UpdateRedeemRuleRequest, operatorID int64) (*RedeemRuleDTO, error) {
	rule, err := s.redeemRuleRepo.FindByID(ctx, s.db, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if rule.Status == points.RedeemRuleStatusActive {
		return nil, code.ErrParam
	}

	rule.Name = req.Name
	rule.Description = req.Description
	rule.CouponID = req.CouponID
	rule.PointsRequired = req.PointsRequired
	rule.TotalStock = req.TotalStock
	rule.PerUserLimit = req.PerUserLimit
	rule.StartAt = timeToInt64(req.StartAt)
	rule.EndAt = timeToInt64(req.EndAt)
	rule.Audit.Update(operatorID)

	if err := s.redeemRuleRepo.Update(ctx, s.db, rule); err != nil {
		return nil, err
	}

	return toRedeemRuleDTO(rule), nil
}

func (s *service) DeleteRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64) error {
	rule, err := s.redeemRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if rule.Status == points.RedeemRuleStatusActive {
		return nil // Cannot delete active rule
	}

	return s.redeemRuleRepo.Delete(ctx, s.db, tenantID, id)
}

func (s *service) GetRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64) (*RedeemRuleDTO, error) {
	rule, err := s.redeemRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toRedeemRuleDTO(rule), nil
}

func (s *service) ListRedeemRules(ctx context.Context, tenantID shared.TenantID, query points.RedeemRuleQuery) ([]*RedeemRuleDTO, int64, *points.RedeemRuleStats, error) {
	rules, total, err := s.redeemRuleRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, nil, err
	}

	stats, err := s.redeemRuleRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, 0, nil, err
	}

	dtos := make([]*RedeemRuleDTO, len(rules))
	for i, r := range rules {
		dtos[i] = toRedeemRuleDTO(r)
	}

	return dtos, total, stats, nil
}

func (s *service) ActivateRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error {
	rule, err := s.redeemRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rule.Activate(operatorID); err != nil {
		return err
	}

	return s.redeemRuleRepo.Update(ctx, s.db, rule)
}

func (s *service) DeactivateRedeemRule(ctx context.Context, tenantID shared.TenantID, id int64, operatorID int64) error {
	rule, err := s.redeemRuleRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := rule.Deactivate(operatorID); err != nil {
		return err
	}

	return s.redeemRuleRepo.Update(ctx, s.db, rule)
}

// ==================== Accounts ====================

func (s *service) GetAccount(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsAccountDTO, error) {
	account, err := s.accountRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toPointsAccountDTO(account), nil
}

func (s *service) GetAccountByUser(ctx context.Context, tenantID shared.TenantID, userID int64) (*PointsAccountDTO, error) {
	account, err := s.accountRepo.FindByUserID(ctx, s.db, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return toPointsAccountDTO(account), nil
}

func (s *service) ListAccounts(ctx context.Context, tenantID shared.TenantID, query points.PointsAccountQuery) ([]*PointsAccountDTO, int64, *points.PointsAccountStats, error) {
	accounts, total, err := s.accountRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, nil, err
	}

	stats, err := s.accountRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, 0, nil, err
	}

	dtos := make([]*PointsAccountDTO, len(accounts))
	for i, a := range accounts {
		dtos[i] = toPointsAccountDTO(a)
	}

	return dtos, total, stats, nil
}

func (s *service) AdjustPoints(ctx context.Context, tenantID shared.TenantID, req AdjustPointsRequest) (*PointsTransactionDTO, error) {
	if req.Points <= 0 {
		return nil, code.ErrParam
	}
	if !req.AdjustmentType.IsValid() {
		return nil, code.ErrParam
	}

	var transaction *points.PointsTransaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		account, err := s.accountRepo.FindByID(ctx, tx, tenantID, req.AccountID)
		if err != nil {
			return err
		}

		var pointsAmount int64
		switch req.AdjustmentType {
		case points.AdjustmentTypeAdd:
			if err := account.AddPoints(req.Points, req.OperatorID); err != nil {
				return err
			}
			pointsAmount = req.Points
		case points.AdjustmentTypeDeduct:
			if err := account.DeductPoints(req.Points, req.OperatorID); err != nil {
				return err
			}
			pointsAmount = -req.Points
		}

		if err := s.accountRepo.Update(ctx, tx, account); err != nil {
			return err
		}

		// Create transaction record
		id, err := s.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		transaction = &points.PointsTransaction{
			ID:           id,
			TenantID:     tenantID,
			UserID:       account.UserID,
			AccountID:    account.ID,
			Points:       pointsAmount,
			BalanceAfter: account.Balance,
			Type:         points.TransactionTypeAdjust,
			Description:  req.Reason,
			Audit:        shared.NewAuditInfo(req.OperatorID),
		}

		return s.transactionRepo.Create(ctx, tx, transaction)
	})

	if err != nil {
		return nil, err
	}

	return toPointsTransactionDTO(transaction), nil
}

func (s *service) GetAccountTransactions(ctx context.Context, tenantID shared.TenantID, accountID int64, query points.PointsTransactionQuery) ([]*PointsTransactionDTO, int64, error) {
	query.AccountID = accountID
	transactions, total, err := s.transactionRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*PointsTransactionDTO, len(transactions))
	for i, t := range transactions {
		dtos[i] = toPointsTransactionDTO(t)
	}

	return dtos, total, nil
}

// ==================== Transactions ====================

func (s *service) GetTransaction(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsTransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toPointsTransactionDTO(transaction), nil
}

func (s *service) ListTransactions(ctx context.Context, tenantID shared.TenantID, query points.PointsTransactionQuery) ([]*PointsTransactionDTO, int64, *points.PointsTransactionStats, error) {
	transactions, total, err := s.transactionRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, nil, err
	}

	stats, err := s.transactionRepo.GetStats(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, nil, err
	}

	dtos := make([]*PointsTransactionDTO, len(transactions))
	for i, t := range transactions {
		dtos[i] = toPointsTransactionDTO(t)
	}

	return dtos, total, stats, nil
}

// ==================== Redemptions ====================

func (s *service) GetRedemption(ctx context.Context, tenantID shared.TenantID, id int64) (*PointsRedemptionDTO, error) {
	redemption, err := s.redemptionRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toPointsRedemptionDTO(redemption), nil
}

func (s *service) ListRedemptions(ctx context.Context, tenantID shared.TenantID, query points.PointsRedemptionQuery) ([]*PointsRedemptionDTO, int64, error) {
	redemptions, total, err := s.redemptionRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*PointsRedemptionDTO, len(redemptions))
	for i, r := range redemptions {
		dtos[i] = toPointsRedemptionDTO(r)
	}

	return dtos, total, nil
}

// ==================== Statistics ====================

func (s *service) GetStats(ctx context.Context, tenantID shared.TenantID, startTime, endTime *time.Time) (*PointsStatsDTO, error) {
	// Get account stats
	accountStats, err := s.accountRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	// Get transaction stats
	txQuery := points.PointsTransactionQuery{
		StartTime: startTime,
		EndTime:   endTime,
	}
	txStats, err := s.transactionRepo.GetStats(ctx, s.db, tenantID, txQuery)
	if err != nil {
		return nil, err
	}

	// Calculate redemption rate
	var redemptionRate string
	if txStats.TotalEarned > 0 {
		rate := decimal.NewFromInt(txStats.TotalRedeemed).Div(decimal.NewFromInt(txStats.TotalEarned)).Mul(decimal.NewFromInt(100))
		redemptionRate = rate.StringFixed(2) + "%"
	} else {
		redemptionRate = "0%"
	}

	// Note: We would need to query expired transactions specifically, simplified for now

	return &PointsStatsDTO{
		TotalIssued:        txStats.TotalEarned,
		TotalRedeemed:      txStats.TotalRedeemed,
		TotalExpired:       0,
		OutstandingBalance: accountStats.TotalBalance,
		RedemptionRate:     redemptionRate,
		ActiveUsers:        accountStats.Active,
		PeriodStart:        formatTime(startTime),
		PeriodEnd:          formatTime(endTime),
	}, nil
}

func (s *service) GetTrend(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time, granularity string) ([]TrendDataPoint, error) {
	// Simplified implementation - in production, this would query aggregated data
	var result []TrendDataPoint

	// Generate data points based on granularity
	switch granularity {
	case "daily":
		for t := startTime; t.Before(endTime); t = t.AddDate(0, 0, 1) {
			// Would query actual data for each day
			result = append(result, TrendDataPoint{
				Date:     t.Format("2006-01-02"),
				Earned:   0,
				Redeemed: 0,
				Expired:  0,
			})
		}
	case "weekly":
		for t := startTime; t.Before(endTime); t = t.AddDate(0, 0, 7) {
			result = append(result, TrendDataPoint{
				Date:     t.Format("2006-01-02"),
				Earned:   0,
				Redeemed: 0,
				Expired:  0,
			})
		}
	}

	return result, nil
}

func (s *service) GetTopUsers(ctx context.Context, tenantID shared.TenantID, startTime, endTime time.Time, limit int) ([]TopUserDTO, error) {
	// Simplified implementation - would query aggregated data with ranking
	return []TopUserDTO{}, nil
}

func (s *service) GetExpiringPoints(ctx context.Context, tenantID shared.TenantID, days int) ([]ExpiringPointsDTO, int64, error) {
	// Simplified implementation - would query transactions with expiration dates
	return []ExpiringPointsDTO{}, 0, nil
}

// ==================== Helper Functions ====================

// timeToInt64 converts *time.Time to *int64 (Unix seconds)
func timeToInt64(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	ts := t.Unix()
	return &ts
}

// int64ToTime converts *int64 (Unix seconds) to *time.Time
func int64ToTime(ts *int64) *time.Time {
	if ts == nil {
		return nil
	}
	t := time.Unix(*ts, 0)
	return &t
}

func toEarnRuleDTO(r *points.EarnRule) *EarnRuleDTO {
	return &EarnRuleDTO{
		ID:               r.ID,
		Name:             r.Name,
		Description:      r.Description,
		Scenario:         string(r.Scenario),
		CalculationType:  string(r.CalculationType),
		FixedPoints:      r.FixedPoints,
		Ratio:            r.Ratio,
		Tiers:            r.Tiers,
		ConditionType:    string(r.ConditionType),
		ConditionValue:   r.ConditionValue,
		ExpirationMonths: r.ExpirationMonths,
		Status:           r.Status.String(),
		Priority:         r.Priority,
		StartAt:          int64ToTime(r.StartAt),
		EndAt:            int64ToTime(r.EndAt),
		CreatedAt:        time.Unix(r.Audit.CreatedAt, 0),
		UpdatedAt:        time.Unix(r.Audit.UpdatedAt, 0),
	}
}

func toRedeemRuleDTO(r *points.RedeemRule) *RedeemRuleDTO {
	return &RedeemRuleDTO{
		ID:             r.ID,
		Name:           r.Name,
		Description:    r.Description,
		CouponID:       r.CouponID,
		PointsRequired: r.PointsRequired,
		TotalStock:     r.TotalStock,
		UsedStock:      r.UsedStock,
		PerUserLimit:   r.PerUserLimit,
		Status:         r.Status.String(),
		StartAt:        int64ToTime(r.StartAt),
		EndAt:          int64ToTime(r.EndAt),
		CreatedAt:      time.Unix(r.Audit.CreatedAt, 0),
		UpdatedAt:      time.Unix(r.Audit.UpdatedAt, 0),
	}
}

func toPointsAccountDTO(a *points.PointsAccount) *PointsAccountDTO {
	return &PointsAccountDTO{
		ID:            a.ID,
		UserID:        a.UserID,
		Balance:       a.Balance,
		FrozenBalance: a.FrozenBalance,
		TotalEarned:   a.TotalEarned,
		TotalRedeemed: a.TotalRedeemed,
		TotalExpired:  a.TotalExpired,
		CreatedAt:     time.Unix(a.Audit.CreatedAt, 0),
		UpdatedAt:     time.Unix(a.Audit.UpdatedAt, 0),
	}
}

func toPointsTransactionDTO(t *points.PointsTransaction) *PointsTransactionDTO {
	return &PointsTransactionDTO{
		ID:            t.ID,
		UserID:        t.UserID,
		AccountID:     t.AccountID,
		Points:        t.Points,
		BalanceAfter:  t.BalanceAfter,
		Type:          string(t.Type),
		ReferenceType: t.ReferenceType,
		ReferenceID:   t.ReferenceID,
		Description:   t.Description,
		ExpiresAt:     int64ToTime(t.ExpiresAt),
		CreatedAt:     time.Unix(t.Audit.CreatedAt, 0),
	}
}

func toPointsRedemptionDTO(r *points.PointsRedemption) *PointsRedemptionDTO {
	return &PointsRedemptionDTO{
		ID:           r.ID,
		UserID:       r.UserID,
		RedeemRuleID: r.RedeemRuleID,
		CouponID:     r.CouponID,
		UserCouponID: r.UserCouponID,
		PointsUsed:   r.PointsUsed,
		Status:       r.Status.String(),
		CompletedAt:  int64ToTime(r.CompletedAt),
		CreatedAt:    time.Unix(r.Audit.CreatedAt, 0),
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}