package user

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"time"

	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	db          *gorm.DB
	userRepo    domain.Repository
	addressRepo domain.AddressRepository
	idGen       snowflake.Snowflake
}

func NewService(db *gorm.DB, userRepo domain.Repository, addressRepo domain.AddressRepository, idGen snowflake.Snowflake) Service {
	return &ServiceImpl{
		db:          db,
		userRepo:    userRepo,
		addressRepo: addressRepo,
		idGen:       idGen,
	}
}

func (s *ServiceImpl) Register(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	exists, err := s.userRepo.Exists(ctx, s.db, req.TenantID, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, code.ErrUserDuplicateUser
	}

	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ID:       id,
		TenantID: req.TenantID,
		Email:    req.Email,
		Phone:    req.Phone,
		Name:     req.Name,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Status:   domain.StatusActive,
		Audit:    shared.NewAuditInfo(0),
	}

	if req.Birthday != "" {
		if t, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			ut := shared.NewUnixTime(t)
			u.Birthday = &ut
		}
	}

	if err := u.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, s.db, u); err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (s *ServiceImpl) Update(ctx context.Context, req UpdateUserRequest) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, s.db, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	u.Name = req.Name
	u.Avatar = req.Avatar
	u.Gender = req.Gender

	if req.Birthday != "" {
		if t, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			ut := shared.NewUnixTime(t)
			u.Birthday = &ut
		}
	}

	u.Audit.Update(0)

	if err := s.userRepo.Update(ctx, s.db, u); err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (s *ServiceImpl) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return code.ErrUserPasswordMismatch
	}

	u, err := s.userRepo.FindByID(ctx, s.db, req.TenantID, req.UserID)
	if err != nil {
		return err
	}

	if !u.CheckPassword(req.OldPassword) {
		return code.ErrUserWrongPassword
	}

	if err := u.SetPassword(req.NewPassword); err != nil {
		return err
	}

	u.Audit.Update(0)
	return s.userRepo.Update(ctx, s.db, u)
}

func (s *ServiceImpl) GetByID(ctx context.Context, tenantID shared.TenantID, id int64) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

func (s *ServiceImpl) GetByEmail(ctx context.Context, tenantID shared.TenantID, email string) (*UserResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, s.db, tenantID, email)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

func (s *ServiceImpl) List(ctx context.Context, tenantID shared.TenantID, req QueryRequest) (*UserListResponse, error) {
	query := domain.Query{
		PageQuery: req.PageQuery,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    req.Status,
	}
	query.PageQuery.Validate()

	users, total, err := s.userRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	resp := &UserListResponse{
		List:     make([]*UserResponse, len(users)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	for i, u := range users {
		resp.List[i] = toUserResponse(u)
	}

	return resp, nil
}

func (s *ServiceImpl) Suspend(ctx context.Context, tenantID shared.TenantID, id int64) error {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := u.Suspend(); err != nil {
		return err
	}

	u.Audit.Update(0)
	return s.userRepo.Update(ctx, s.db, u)
}

func (s *ServiceImpl) Activate(ctx context.Context, tenantID shared.TenantID, id int64) error {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := u.Activate(); err != nil {
		return err
	}

	u.Audit.Update(0)
	return s.userRepo.Update(ctx, s.db, u)
}

func (s *ServiceImpl) Delete(ctx context.Context, tenantID shared.TenantID, id int64) error {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return err
	}

	if err := u.SoftDelete(); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, s.db, u)
}

func (s *ServiceImpl) ResetPassword(ctx context.Context, tenantID shared.TenantID, id int64) (string, error) {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return "", err
	}

	// Generate temporary password
	tempPassword := generateTempPassword()
	if err := u.SetPassword(tempPassword); err != nil {
		return "", err
	}

	u.Audit.Update(0)
	if err := s.userRepo.Update(ctx, s.db, u); err != nil {
		return "", err
	}

	return tempPassword, nil
}

func (s *ServiceImpl) GetStats(ctx context.Context, tenantID shared.TenantID) (*UserStatsResponse, error) {
	stats, err := s.userRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	return &UserStatsResponse{
		Total:     stats.Total,
		Active:    stats.Active,
		Suspended: stats.Suspended,
		NewToday:  stats.NewToday,
	}, nil
}

func generateTempPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLen = byte(len(charset))
	b := make([]byte, 12)

	// Use crypto/rand for secure random generation - no fallback to weak random
	randomBytes := make([]byte, 12)
	if _, err := crand.Read(randomBytes); err != nil {
		// crypto/rand failure is critical - log and panic
		// This should never happen on a properly configured system
		panic("failed to generate secure temporary password: " + err.Error())
	}

	for i, rb := range randomBytes {
		b[i] = charset[rb%charsetLen]
	}
	return string(b)
}

func toUserResponse(u *domain.User) *UserResponse {
	resp := &UserResponse{
		ID:        u.ID,
		TenantID:  u.TenantID.Int64(),
		Email:     u.Email,
		Phone:     u.Phone,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Gender:    int(u.Gender),
		Status:    int(u.Status),
		CreatedAt: u.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if u.Birthday != nil {
		resp.Birthday = u.Birthday.Format("2006-01-02")
	}

	if u.LastLogin != nil {
		resp.LastLogin = u.LastLogin.Format("2006-01-02 15:04:05")
	}

	return resp
}

// GetDetail returns detailed user information with points and order statistics
func (s *ServiceImpl) GetDetail(ctx context.Context, tenantID shared.TenantID, id int64) (*UserDetailResponse, error) {
	u, err := s.userRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	resp := toUserDetailResponse(u)

	// Query points account data
	type PointsAccountData struct {
		Balance      int64 `gorm:"column:balance"`
		Frozen       int64 `gorm:"column:frozen_balance"`
		TotalEarned  int64 `gorm:"column:total_earned"`
		TotalRedeemed int64 `gorm:"column:total_redeemed"`
	}
	var pointsAccount PointsAccountData
	s.db.WithContext(ctx).
		Table("points_accounts").
		Select("balance, frozen_balance, total_earned, total_redeemed").
		Where("user_id = ? AND tenant_id = ?", id, tenantID.Int64()).
		First(&pointsAccount)
	resp.PointsBalance = pointsAccount.Balance
	resp.PointsFrozen = pointsAccount.Frozen
	resp.TotalEarned = pointsAccount.TotalEarned
	resp.TotalRedeemed = pointsAccount.TotalRedeemed

	// Query order statistics
	type OrderStatsData struct {
		Count      int64 `gorm:"column:count"`
		TotalSpent int64 `gorm:"column:total_spent"`
	}
	var orderStats OrderStatsData
	s.db.WithContext(ctx).
		Table("orders").
		Select("COUNT(*) as count, COALESCE(SUM(pay_amount), 0) as total_spent").
		Where("user_id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID.Int64()).
		Scan(&orderStats)
	resp.OrderCount = orderStats.Count
	// Convert cents to yuan string with 2 decimal places
	resp.TotalSpent = formatAmountFromCents(orderStats.TotalSpent)

	// Query review count
	var reviewCount int64
	s.db.WithContext(ctx).
		Table("reviews").
		Where("user_id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID.Int64()).
		Count(&reviewCount)
	resp.ReviewCount = reviewCount

	return resp, nil
}

// ExtendedList returns a list of users with extended information
func (s *ServiceImpl) ExtendedList(ctx context.Context, tenantID shared.TenantID, req EnhancedQueryRequest) (*ExtendedListResponse, error) {
	query := domain.Query{
		PageQuery: req.PageQuery,
		Name:      req.Keyword,
		Email:     req.Keyword,
		Phone:     req.Keyword,
		Status:    req.Status,
	}
	query.PageQuery.Validate()

	users, total, err := s.userRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	resp := &ExtendedListResponse{
		List:     make([]*ExtendedUserResponse, len(users)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Collect user IDs for batch queries
	userIDs := make([]int64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	// Batch query points accounts
	pointsMap := make(map[int64]struct {
		Balance int64
	})
	if len(userIDs) > 0 {
		type PointsData struct {
			UserID  int64 `gorm:"column:user_id"`
			Balance int64 `gorm:"column:balance"`
		}
		var pointsData []PointsData
		s.db.WithContext(ctx).
			Table("points_accounts").
			Select("user_id, balance").
			Where("user_id IN ? AND tenant_id = ?", userIDs, tenantID.Int64()).
			Find(&pointsData)
		for _, pd := range pointsData {
			pointsMap[pd.UserID] = struct{ Balance int64 }{Balance: pd.Balance}
		}
	}

	// Batch query order statistics
	orderStatsMap := make(map[int64]struct {
		Count      int64
		TotalSpent int64
	})
	if len(userIDs) > 0 {
		type OrderData struct {
			UserID     int64 `gorm:"column:user_id"`
			Count      int64 `gorm:"column:count"`
			TotalSpent int64 `gorm:"column:total_spent"`
		}
		var orderData []OrderData
		s.db.WithContext(ctx).
			Table("orders").
			Select("user_id, COUNT(*) as count, COALESCE(SUM(pay_amount), 0) as total_spent").
			Where("user_id IN ? AND tenant_id = ? AND deleted_at IS NULL", userIDs, tenantID.Int64()).
			Group("user_id").
			Find(&orderData)
		for _, od := range orderData {
			orderStatsMap[od.UserID] = struct {
				Count      int64
				TotalSpent int64
			}{Count: od.Count, TotalSpent: od.TotalSpent}
		}
	}

	for i, u := range users {
		resp.List[i] = toExtendedUserResponse(u)
		// Populate points balance
		if pd, ok := pointsMap[u.ID]; ok {
			resp.List[i].PointsBalance = pd.Balance
		}
		// Populate order statistics
		if od, ok := orderStatsMap[u.ID]; ok {
			resp.List[i].OrderCount = od.Count
			resp.List[i].TotalSpent = formatAmountFromCents(od.TotalSpent)
		}
	}

	return resp, nil
}

// GetAddresses returns all addresses for a user
func (s *ServiceImpl) GetAddresses(ctx context.Context, tenantID shared.TenantID, userID int64) (*AddressListResponse, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, s.db, tenantID, userID)
	if err != nil {
		return nil, err
	}

	resp := &AddressListResponse{
		List:  make([]*AddressResponse, len(addresses)),
		Total: int64(len(addresses)),
	}

	for i, addr := range addresses {
		resp.List[i] = toAddressResponse(addr)
	}

	return resp, nil
}

// SuspendWithReason suspends a user with a specific reason
func (s *ServiceImpl) SuspendWithReason(ctx context.Context, req SuspendUserRequest) error {
	return s.Suspend(ctx, req.TenantID, req.UserID)
}

// ActivateUser activates a user
func (s *ServiceImpl) ActivateUser(ctx context.Context, tenantID shared.TenantID, userID int64) error {
	return s.Activate(ctx, tenantID, userID)
}

// DeleteUser deletes a user
func (s *ServiceImpl) DeleteUser(ctx context.Context, tenantID shared.TenantID, userID int64) error {
	return s.Delete(ctx, tenantID, userID)
}

// GetUserStats returns user statistics
func (s *ServiceImpl) GetUserStats(ctx context.Context, tenantID shared.TenantID) (*UserStats, error) {
	stats, err := s.userRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	return &UserStats{
		TotalUsers:     stats.Total,
		ActiveUsers:    stats.Active,
		SuspendedUsers: stats.Suspended,
		NewUsersToday:  stats.NewToday,
	}, nil
}

// ExportUsers exports users to a file (placeholder implementation)
func (s *ServiceImpl) ExportUsers(ctx context.Context, tenantID shared.TenantID, req EnhancedQueryRequest) ([]byte, error) {
	// Get all users matching the criteria
	query := domain.Query{
		PageQuery: shared.PageQuery{Page: 1, PageSize: 10000}, // Large page size for export
		Name:      req.Keyword,
		Email:     req.Keyword,
		Phone:     req.Keyword,
		Status:    req.Status,
	}

	users, _, err := s.userRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	// Generate CSV content
	var csvContent string
	csvContent = "ID,Email,Phone,Name,Status,CreatedAt\n"
	for _, u := range users {
		statusText := "active"
		if u.Status == domain.StatusSuspended {
			statusText = "suspended"
		} else if u.Status == domain.StatusInactive {
			statusText = "inactive"
		}
		csvContent += fmt.Sprintf("%d,%s,%s,%s,%s,%s\n",
			u.ID, u.Email, u.Phone, u.Name, statusText, u.Audit.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return []byte(csvContent), nil
}

func toUserDetailResponse(u *domain.User) *UserDetailResponse {
	resp := &UserDetailResponse{
		ID:        u.ID,
		TenantID:  u.TenantID.Int64(),
		Email:     u.Email,
		Phone:     u.Phone,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Gender:    int(u.Gender),
		GenderText: getGenderText(u.Gender),
		Status:    int(u.Status),
		StatusText: getStatusText(u.Status),
		CreatedAt: u.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.Audit.UpdatedAt.Format(time.RFC3339),
	}

	if u.Birthday != nil {
		resp.Birthday = u.Birthday.Format("2006-01-02")
	}

	if u.LastLogin != nil {
		resp.LastLogin = u.LastLogin.Format(time.RFC3339)
	}

	return resp
}

func toExtendedUserResponse(u *domain.User) *ExtendedUserResponse {
	resp := &ExtendedUserResponse{
		ID:         u.ID,
		TenantID:   u.TenantID.Int64(),
		Email:      u.Email,
		Phone:      u.Phone,
		Name:       u.Name,
		Avatar:     u.Avatar,
		Status:     int(u.Status),
		StatusText: getStatusText(u.Status),
		CreatedAt:  u.Audit.CreatedAt.Format(time.RFC3339),
	}

	if u.LastLogin != nil {
		resp.LastLogin = u.LastLogin.Format(time.RFC3339)
	}

	return resp
}

func toAddressResponse(addr *domain.UserAddress) *AddressResponse {
	return &AddressResponse{
		ID:         addr.ID,
		Name:       addr.Name,
		Phone:      addr.Phone,
		Country:    addr.Country,
		Province:   addr.Province,
		City:       addr.City,
		District:   addr.District,
		Address:    addr.Address,
		PostalCode: addr.PostalCode,
		IsDefault:  addr.IsDefault,
		CreatedAt:  addr.CreatedAt.Format(time.RFC3339),
	}
}

func getGenderText(gender domain.Gender) string {
	switch gender {
	case domain.GenderMale:
		return "男"
	case domain.GenderFemale:
		return "女"
	case domain.GenderOther:
		return "其他"
	default:
		return "未知"
	}
}

func getStatusText(status domain.Status) string {
	switch status {
	case domain.StatusActive:
		return "正常"
	case domain.StatusSuspended:
		return "已冻结"
	case domain.StatusInactive:
		return "未激活"
	default:
		return "未知"
	}
}

// formatAmountFromCents converts amount in cents to yuan string with 2 decimal places
func formatAmountFromCents(cents int64) string {
	yuan := cents / 100
	remainCents := cents % 100
	if remainCents < 0 {
		remainCents = -remainCents
	}
	if cents < 0 && yuan == 0 {
		return fmt.Sprintf("-%d.%02d", yuan, remainCents)
	}
	return fmt.Sprintf("%d.%02d", yuan, remainCents)
}
