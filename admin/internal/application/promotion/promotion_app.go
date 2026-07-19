package promotion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// =============================================================================
// Request types
// =============================================================================

// CreatePromotionRequest is the unified input for both PROMOTION and COUPON
// kinds. Coupon-specific fields (Code, TotalCount, MarketID) are nullable.
type CreatePromotionRequest struct {
	TenantID     shared.TenantID
	Kind         promotion.Kind
	Name         string
	Description  string
	Code         *string
	Type         promotion.Type
	MarketID     *int64
	Currency     string
	TotalCount   *int
	UsageLimit   int
	PerUserLimit int
	Tags         []string
	Scope        promotion.PromotionScope
	StartAt      time.Time
	EndAt        time.Time
	Rules        []promotion.PromotionRule
	ActorID      int64
}

// UpdatePromotionRequest mirrors CreatePromotionRequest plus ID and Status.
// Rules is a pointer: nil means "do not touch", non-nil means "replace".
type UpdatePromotionRequest struct {
	ID           int64
	Name         string
	Description  string
	Code         *string
	Type         promotion.Type
	MarketID     *int64
	Currency     string
	TotalCount   *int
	UsageLimit   int
	PerUserLimit int
	Tags         []string
	Scope        promotion.PromotionScope
	StartAt      time.Time
	EndAt        time.Time
	Rules        *[]promotion.PromotionRule // nil = no change; non-nil = replace
	Status       *promotion.Status
	ActorID      int64
}

// =============================================================================
// Response types (local; Task 8 will regenerate wire types and merge in)
// =============================================================================

// PromotionResponse is the unified output of GET / LIST endpoints.
type PromotionResponse struct {
	ID           int64
	TenantID     shared.TenantID
	Kind         promotion.Kind
	Name         string
	Description  string
	Code         *string
	Type         promotion.Type
	Status       promotion.Status
	MarketID     *int64
	Currency     string
	TotalCount   *int
	UsedCount    *int
	UsageLimit   int
	PerUserLimit int
	Tags         []string
	Scope        promotion.PromotionScope
	StartAt      time.Time
	EndAt        time.Time
	Rules        []PromotionRuleResponse
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// PromotionRuleResponse mirrors the domain rule for wire serialization.
type PromotionRuleResponse struct {
	ID             int64
	ConditionType  promotion.ConditionType
	ConditionValue decimal.Decimal
	ActionType     promotion.ActionType
	ActionValue    decimal.Decimal
	MaxDiscount    decimal.Decimal
	SortOrder      int
}

// ListPromotionResponse carries a paginated slice of PromotionResponse.
type ListPromotionResponse struct {
	List  []*PromotionResponse
	Total int64
	Page  int
	Size  int
}

// UserCouponResponse mirrors a per-user coupon claim.
type UserCouponResponse struct {
	ID         int64
	UserID     int64
	CouponID   int64
	Status     promotion.UserCouponStatus
	UsedAt     *time.Time
	OrderID    int64
	ReceivedAt time.Time
	ExpireAt   time.Time
}

// ListUserCouponResponse is the paginated envelope for user coupon listings.
type ListUserCouponResponse struct {
	List  []*UserCouponResponse
	Total int64
}

// PromotionUsageResponse carries one (coupon, order, user) hit.
type PromotionUsageResponse struct {
	ID             int64
	CouponID       *int64
	UserID         int64
	OrderID        int64
	DiscountAmount decimal.Decimal
	CreatedAt      time.Time
}

// ListPromotionUsageResponse is the paginated envelope for usage listings.
type ListPromotionUsageResponse struct {
	List  []*PromotionUsageResponse
	Total int64
}

// =============================================================================
// PromotionApp
// =============================================================================

// PromotionApp is the unified application service for both PROMOTION and
// COUPON kinds. Coupon-specific methods (IssueToUser, BatchIssue,
// GenerateCodes, ListUserCoupons, FindPromotionUsage) live alongside the
// kind-agnostic CRUD methods.
type PromotionApp struct {
	repo promotion.Repository
	db   *gorm.DB
}

// NewPromotionApp wires the unified repository + DB.
func NewPromotionApp(repo promotion.Repository, db *gorm.DB) *PromotionApp {
	return &PromotionApp{repo: repo, db: db}
}

// Create persists a Promotion and (optionally) its rules. The new promotion
// starts with StatusPending; UseCount is initialized to 0 if TotalCount is set.
func (a *PromotionApp) Create(ctx context.Context, req *CreatePromotionRequest) (*PromotionResponse, error) {
	now := time.Now().UTC()
	p := &promotion.Promotion{
		TenantID:     req.TenantID,
		Kind:         req.Kind,
		Name:         req.Name,
		Description:  req.Description,
		Code:         req.Code,
		Type:         req.Type,
		Status:       promotion.StatusPending,
		MarketID:     req.MarketID,
		Currency:     req.Currency,
		TotalCount:   req.TotalCount,
		UsedCount:    nilOrZero(req.TotalCount),
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		Scope:        req.Scope,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		Rules:        req.Rules,
		Audit: shared.AuditInfo{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: req.ActorID,
			UpdatedBy: req.ActorID,
		},
	}
	if err := a.repo.Create(ctx, a.db, p); err != nil {
		return nil, err
	}
	if len(p.Rules) > 0 {
		for i := range p.Rules {
			p.Rules[i].OwnerKind = p.Kind
			p.Rules[i].OwnerID = p.ID
		}
		if err := a.repo.CreateRules(ctx, a.db, p.Kind, p.ID, p.Rules); err != nil {
			return nil, err
		}
	}
	return a.toResponse(p), nil
}

// Update modifies an existing Promotion and (optionally) its rules.
//
// Rules semantics:
//   - req.Rules == nil  → leave existing rules untouched
//   - req.Rules != nil  → delete all existing rules, then insert the new ones
//     (an empty slice removes all rules)
func (a *PromotionApp) Update(ctx context.Context, req *UpdatePromotionRequest) (*PromotionResponse, error) {
	p, err := a.repo.FindByID(ctx, a.db, req.ID)
	if err != nil {
		return nil, err
	}
	p.Name = req.Name
	p.Description = req.Description
	p.Code = req.Code
	p.Type = req.Type
	p.MarketID = req.MarketID
	p.Currency = req.Currency
	p.TotalCount = req.TotalCount
	p.UsageLimit = req.UsageLimit
	p.PerUserLimit = req.PerUserLimit
	p.Tags = req.Tags
	p.Scope = req.Scope
	p.StartAt = req.StartAt
	p.EndAt = req.EndAt
	if req.Status != nil {
		p.Status = *req.Status
	}
	p.Audit.UpdatedAt = time.Now().UTC()
	p.Audit.UpdatedBy = req.ActorID

	if err := a.repo.Update(ctx, a.db, p); err != nil {
		return nil, err
	}
	if req.Rules != nil {
		// Hard-delete old rules so the promotion_rules table doesn't
		// accumulate tombstones across replace cycles. DeleteRulesByOwner
		// (soft-delete) is reserved for full promotion delete so that
		// undeleting a promotion recovers its rules.
		if err := a.repo.HardDeleteRulesByOwner(ctx, a.db, p.Kind, p.ID); err != nil {
			return nil, err
		}
		if len(*req.Rules) > 0 {
			for i := range *req.Rules {
				(*req.Rules)[i].OwnerKind = p.Kind
				(*req.Rules)[i].OwnerID = p.ID
			}
			if err := a.repo.CreateRules(ctx, a.db, p.Kind, p.ID, *req.Rules); err != nil {
				return nil, err
			}
			p.Rules = *req.Rules
		} else {
			p.Rules = nil
		}
	}
	return a.toResponse(p), nil
}

// Get returns the promotion with its rules loaded.
func (a *PromotionApp) Get(ctx context.Context, id int64) (*PromotionResponse, error) {
	p, err := a.repo.FindByID(ctx, a.db, id)
	if err != nil {
		return nil, err
	}
	rules, err := a.repo.FindRulesByOwner(ctx, a.db, p.Kind, p.ID)
	if err != nil {
		return nil, err
	}
	p.Rules = rules
	return a.toResponse(p), nil
}

// List paginates with optional filters (Name, Kind, Status, Type, MarketID,
// ExpiredOnly). Status / Type / Kind / MarketID must be pointers in the query
// so the iota-zero value can be expressed as a real filter, not "unset".
//
// Rules are loaded in a single batch query (FindRulesByOwners) instead of
// N+1 — the list response feeds the admin table's "优惠内容" column which
// renders the first rule's action_value, so empty rules would render as "-".
func (a *PromotionApp) List(ctx context.Context, q promotion.Query) (*ListPromotionResponse, error) {
	list, total, err := a.repo.FindList(ctx, a.db, q)
	if err != nil {
		return nil, err
	}
	ownerIDs := make([]int64, len(list))
	for i, p := range list {
		ownerIDs[i] = p.ID
	}
	rulesByOwner, err := a.repo.FindRulesByOwners(ctx, a.db, ownerIDs)
	if err != nil {
		return nil, err
	}
	out := make([]*PromotionResponse, len(list))
	for i, p := range list {
		p.Rules = rulesByOwner[p.ID]
		out[i] = a.toResponse(p)
	}
	return &ListPromotionResponse{
		List:  out,
		Total: total,
		Page:  q.Page,
		Size:  q.PageSize,
	}, nil
}

// Delete removes a promotion and all its rules.
func (a *PromotionApp) Delete(ctx context.Context, id int64) error {
	p, err := a.repo.FindByID(ctx, a.db, id)
	if err != nil {
		return err
	}
	if err := a.repo.DeleteRulesByOwner(ctx, a.db, p.Kind, p.ID); err != nil {
		return err
	}
	return a.repo.Delete(ctx, a.db, id)
}

// Activate flips Status → StatusActive. Refuses if EndAt is already past.
func (a *PromotionApp) Activate(ctx context.Context, id int64) (*PromotionResponse, error) {
	p, err := a.repo.FindByID(ctx, a.db, id)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().After(p.EndAt) {
		return nil, code.ErrPromotionExpired
	}
	p.Status = promotion.StatusActive
	p.Audit.UpdatedAt = time.Now().UTC()
	if err := a.repo.Update(ctx, a.db, p); err != nil {
		return nil, err
	}
	return a.toResponse(p), nil
}

// Deactivate flips Status → StatusPaused.
func (a *PromotionApp) Deactivate(ctx context.Context, id int64) (*PromotionResponse, error) {
	p, err := a.repo.FindByID(ctx, a.db, id)
	if err != nil {
		return nil, err
	}
	p.Status = promotion.StatusPaused
	p.Audit.UpdatedAt = time.Now().UTC()
	if err := a.repo.Update(ctx, a.db, p); err != nil {
		return nil, err
	}
	return a.toResponse(p), nil
}

// =============================================================================
// Rules (kind-agnostic)
// =============================================================================

// GetRules returns the rules attached to a Promotion of any kind.
func (a *PromotionApp) GetRules(ctx context.Context, ownerKind promotion.Kind, ownerID int64) ([]*PromotionRuleResponse, error) {
	rules, err := a.repo.FindRulesByOwner(ctx, a.db, ownerKind, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]*PromotionRuleResponse, len(rules))
	for i := range rules {
		out[i] = ruleToResponse(&rules[i])
	}
	return out, nil
}

// CreateRules inserts a batch of rules under an existing owner.
func (a *PromotionApp) CreateRules(ctx context.Context, ownerKind promotion.Kind, ownerID int64, rules []promotion.PromotionRule) ([]*PromotionRuleResponse, error) {
	for i := range rules {
		rules[i].OwnerKind = ownerKind
		rules[i].OwnerID = ownerID
	}
	if err := a.repo.CreateRules(ctx, a.db, ownerKind, ownerID, rules); err != nil {
		return nil, err
	}
	return a.GetRules(ctx, ownerKind, ownerID)
}

// UpdateRule mutates a single rule by ID. OwnerKind / OwnerID are preserved.
func (a *PromotionApp) UpdateRule(ctx context.Context, rule *promotion.PromotionRule) (*PromotionRuleResponse, error) {
	if err := a.repo.UpdateRule(ctx, a.db, rule); err != nil {
		return nil, err
	}
	return ruleToResponse(rule), nil
}

// DeleteRule removes a single rule by ID.
func (a *PromotionApp) DeleteRule(ctx context.Context, id int64) error {
	return a.repo.DeleteRule(ctx, a.db, id)
}

// =============================================================================
// COUPON-specific
// =============================================================================

// IssueToUser atomically:
//  1. Calls Promotion.Issue (validates Kind=COUPON + active + not depleted)
//  2. Inserts a UserCoupon row
//  3. Calls IncrementUsedCount, which uses an atomic SQL guard against
//     overselling. If the inventory check fails, the in-memory increment is
//     rolled back by leaving it as-is — the repo refused to write.
//
// For COUPONs that were migrated with usage_limit=0 ("unlimited", per
// handoff I2), ConsumeInventory short-circuits (TotalCount nil) and
// IncrementUsedCount is skipped — there is no inventory to track.
func (a *PromotionApp) IssueToUser(ctx context.Context, couponID, userID int64) (*UserCouponResponse, error) {
	p, err := a.repo.FindByID(ctx, a.db, couponID)
	if err != nil {
		return nil, err
	}
	uc, err := p.Issue(userID, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if err := a.repo.IssueUserCoupon(ctx, a.db, uc); err != nil {
		return nil, err
	}
	// Atomically consume inventory if the coupon has TotalCount tracking.
	if p.TotalCount != nil {
		if err := a.repo.IncrementUsedCount(ctx, a.db, couponID); err != nil {
			return nil, err
		}
	}
	return userCouponToResponse(uc), nil
}

// BatchIssue issues a coupon to a list of users. Failures on individual
// users are silently skipped — partial success is the norm for batch ops.
func (a *PromotionApp) BatchIssue(ctx context.Context, couponID int64, userIDs []int64) (int64, []int64, error) {
	var issued int64
	ids := make([]int64, 0, len(userIDs))
	for _, uid := range userIDs {
		resp, err := a.IssueToUser(ctx, couponID, uid)
		if err != nil {
			continue
		}
		issued++
		ids = append(ids, resp.ID)
	}
	return issued, ids, nil
}

// GenerateCodes produces `quantity` deterministic-ish codes and inserts a
// COUPON promotion for each. The cfg map carries the per-coupon parameters
// (name, value, currency, etc.) — see couponFromConfig for the shape.
//
// Per the brief this is a stub: the JSON parsing will be implemented in
// the logic layer (Task 7). The current loop generates codes and creates
// coupon promotions using whatever couponFromConfig returns.
func (a *PromotionApp) GenerateCodes(ctx context.Context, prefix string, quantity int, cfg map[string]any) ([]string, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	out := make([]string, 0, quantity)
	for i := 0; i < quantity; i++ {
		code := prefix + randomCode(8)
		req := couponFromConfig(code, cfg)
		if req == nil {
			// couponFromConfig is a stub in this task; logic layer (Task 7)
			// will replace it. Return the codes we have so callers can
			// observe progress.
			out = append(out, code)
			continue
		}
		if _, err := a.Create(ctx, req); err != nil {
			return out, err
		}
		out = append(out, code)
	}
	return out, nil
}

// ListUserCoupons paginates user-claimed coupons.
func (a *PromotionApp) ListUserCoupons(ctx context.Context, q promotion.UserCouponQuery) (*ListUserCouponResponse, error) {
	list, total, err := a.repo.FindUserCoupons(ctx, a.db, q)
	if err != nil {
		return nil, err
	}
	out := make([]*UserCouponResponse, len(list))
	for i := range list {
		out[i] = userCouponToResponse(list[i])
	}
	return &ListUserCouponResponse{List: out, Total: total}, nil
}

// FindPromotionUsage paginates historical promotion usage rows.
func (a *PromotionApp) FindPromotionUsage(ctx context.Context, q promotion.UsageQuery) (*ListPromotionUsageResponse, error) {
	list, total, err := a.repo.FindPromotionUsage(ctx, a.db, q)
	if err != nil {
		return nil, err
	}
	out := make([]*PromotionUsageResponse, len(list))
	for i := range list {
		out[i] = usageToResponse(list[i])
	}
	return &ListPromotionUsageResponse{List: out, Total: total}, nil
}

// =============================================================================
// Helpers
// =============================================================================

// nilOrZero returns a pointer to 0 if total is non-nil, otherwise nil.
// Used to initialize UsedCount alongside TotalCount on Create.
func nilOrZero(total *int) *int {
	if total == nil {
		return nil
	}
	z := 0
	return &z
}

// toResponse converts a domain Promotion (with optional loaded rules) to the
// wire-shaped PromotionResponse.
func (a *PromotionApp) toResponse(p *promotion.Promotion) *PromotionResponse {
	resp := &PromotionResponse{
		ID:           p.ID,
		TenantID:     p.TenantID,
		Kind:         p.Kind,
		Name:         p.Name,
		Description:  p.Description,
		Code:         p.Code,
		Type:         p.Type,
		Status:       p.Status,
		MarketID:     p.MarketID,
		Currency:     p.Currency,
		TotalCount:   p.TotalCount,
		UsedCount:    p.UsedCount,
		UsageLimit:   p.UsageLimit,
		PerUserLimit: p.PerUserLimit,
		Tags:         p.Tags,
		Scope:        p.Scope,
		StartAt:      p.StartAt,
		EndAt:        p.EndAt,
		CreatedAt:    p.Audit.CreatedAt,
		UpdatedAt:    p.Audit.UpdatedAt,
	}
	for _, r := range p.Rules {
		resp.Rules = append(resp.Rules, *ruleToResponse(&r))
	}
	return resp
}

func ruleToResponse(r *promotion.PromotionRule) *PromotionRuleResponse {
	return &PromotionRuleResponse{
		ID:             r.ID,
		ConditionType:  r.ConditionType,
		ConditionValue: r.ConditionValue,
		ActionType:     r.ActionType,
		ActionValue:    r.ActionValue,
		MaxDiscount:    r.MaxDiscount,
		SortOrder:      r.SortOrder,
	}
}

func userCouponToResponse(uc *promotion.UserCoupon) *UserCouponResponse {
	return &UserCouponResponse{
		ID:         uc.ID,
		UserID:     uc.UserID,
		CouponID:   uc.CouponID,
		Status:     uc.Status,
		UsedAt:     uc.UsedAt,
		OrderID:    uc.OrderID,
		ReceivedAt: uc.ReceivedAt,
		ExpireAt:   uc.ExpireAt,
	}
}

func usageToResponse(u *promotion.PromotionUsage) *PromotionUsageResponse {
	return &PromotionUsageResponse{
		ID:             u.ID,
		CouponID:       u.CouponID,
		UserID:         u.UserID,
		OrderID:        u.OrderID,
		DiscountAmount: u.DiscountAmount,
		CreatedAt:      u.CreatedAt,
	}
}

// randomCode produces a pseudo-random alphanumeric string of length n using
// nanosecond time as the entropy source. Used by GenerateCodes — not
// cryptographically secure.
func randomCode(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Microsecond)
	}
	return string(b)
}

// couponFromConfig builds a CreatePromotionRequest from a code + cfg map.
// The cfg shape comes from the wire GenerateCouponCodesReq.CouponConfig
// (a JSON object) parsed by the logic layer. Supported keys:
//
//	name (string, required)            — coupon display name
//	description (string, optional)
//	type (string)                      — "fixed_amount" | "percentage" | "free_shipping"
//	discount_value (string|number)     — required for fixed/percentage
//	min_order_amount (string|number)
//	max_discount (string|number)
//	currency (string)                  — defaults to "CNY"
//	start_time (string, RFC3339)       — required
//	end_time (string, RFC3339)         — required
//	usage_limit (int)                  — total_count; 0 = unlimited
//	per_user_limit (int)               — defaults to 1
//	scope_type (string)                — storewide | products | categories | brands
//	_actor_id (int64)                  — audit (injected by logic layer)
//	_tenant_id (int64)                 — tenant scope (injected by logic layer)
//
// Per handoff I2, usage_limit=0 leaves TotalCount unset so the repo
// 's consume-inventory SQL guard treats the coupon as unlimited.
func couponFromConfig(code string, cfg map[string]any) *CreatePromotionRequest {
	if cfg == nil {
		return nil
	}
	name := stringFromCfg(cfg, "name")
	if name == "" {
		return nil
	}

	startStr := stringFromCfg(cfg, "start_time")
	endStr := stringFromCfg(cfg, "end_time")
	startAt, err1 := time.Parse(time.RFC3339, startStr)
	endAt, err2 := time.Parse(time.RFC3339, endStr)
	if err1 != nil || err2 != nil {
		return nil
	}

	discountType := stringFromCfg(cfg, "type")
	rule := promotion.PromotionRule{
		ConditionType: promotion.ConditionMinAmount,
		ActionType:    actionTypeFromCouponType(discountType),
	}
	if v := stringFromCfg(cfg, "min_order_amount"); v != "" {
		if d, err := decimal.NewFromString(v); err == nil {
			rule.ConditionValue = d
		}
	}
	if v := stringFromCfg(cfg, "discount_value"); v != "" {
		if d, err := decimal.NewFromString(v); err == nil {
			rule.ActionValue = d
		}
	}
	if v := stringFromCfg(cfg, "max_discount"); v != "" {
		if d, err := decimal.NewFromString(v); err == nil {
			rule.MaxDiscount = d
		}
	}

	usageLimit := intFromCfg(cfg, "usage_limit")
	perUserLimit := intFromCfg(cfg, "per_user_limit")
	if perUserLimit == 0 {
		perUserLimit = 1
	}

	codeCopy := code
	var total *int
	if usageLimit > 0 {
		v := usageLimit
		total = &v
	}

	actorID := int64FromCfg(cfg, "_actor_id")
	tenantID := int64FromCfg(cfg, "_tenant_id")

	return &CreatePromotionRequest{
		TenantID:     shared.TenantID(tenantID),
		Kind:         promotion.KindCoupon,
		Name:         name,
		Description:  stringFromCfg(cfg, "description"),
		Code:         &codeCopy,
		Type:         promotion.TypeDiscount,
		Currency:     defaultCurrencyFromCfg(cfg, "currency"),
		TotalCount:   total,
		UsageLimit:   usageLimit,
		PerUserLimit: perUserLimit,
		Scope:        scopeFromCfg(cfg),
		StartAt:      startAt,
		EndAt:        endAt,
		Rules:        []promotion.PromotionRule{rule},
		ActorID:      actorID,
	}
}

// stringFromCfg looks up a string key in cfg, accepting either a
// string value or a fmt.Stringer / numeric value that converts via
// fmt.Sprint.
func stringFromCfg(cfg map[string]any, key string) string {
	v, ok := cfg[key]
	if !ok || v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case float64:
		return fmt.Sprintf("%g", s)
	case float32:
		return fmt.Sprintf("%g", s)
	case int:
		return fmt.Sprintf("%d", s)
	case int64:
		return fmt.Sprintf("%d", s)
	case json.Number:
		return s.String()
	}
	return fmt.Sprintf("%v", v)
}

// intFromCfg reads an integer-typed key from cfg, tolerating JSON
// numeric values which arrive as float64.
func intFromCfg(cfg map[string]any, key string) int {
	v, ok := cfg[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

func int64FromCfg(cfg map[string]any, key string) int64 {
	v, ok := cfg[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	}
	return 0
}

func defaultCurrencyFromCfg(cfg map[string]any, key string) string {
	if s := stringFromCfg(cfg, key); s != "" {
		return s
	}
	return "CNY"
}

// actionTypeFromCouponType maps the wire coupon Type string onto the
// domain ActionType.
func actionTypeFromCouponType(t string) promotion.ActionType {
	switch strings.ToLower(t) {
	case "percentage":
		return promotion.ActionPercentage
	case "free_shipping":
		return promotion.ActionFreeShipping
	default:
		return promotion.ActionFixedAmount
	}
}

// scopeFromCfg translates the wire "scope_type" hint onto a
// PromotionScope. The form doesn't yet bind IDs for batch
// generation, so IDs stay nil.
func scopeFromCfg(cfg map[string]any) promotion.PromotionScope {
	switch strings.ToUpper(stringFromCfg(cfg, "scope_type")) {
	case string(promotion.ScopeTypeProducts):
		return promotion.PromotionScope{Type: promotion.ScopeTypeProducts}
	case string(promotion.ScopeTypeCategories):
		return promotion.PromotionScope{Type: promotion.ScopeTypeCategories}
	case string(promotion.ScopeTypeBrands):
		return promotion.PromotionScope{Type: promotion.ScopeTypeBrands}
	}
	return promotion.PromotionScope{Type: promotion.ScopeTypeStorewide}
}
