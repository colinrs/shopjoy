package user

import (
	"context"
	crand "crypto/rand"
	"time"

	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	db       *gorm.DB
	userRepo domain.Repository
	idGen    snowflake.Snowflake
}

func NewService(db *gorm.DB, userRepo domain.Repository, idGen snowflake.Snowflake) Service {
	return &ServiceImpl{
		db:       db,
		userRepo: userRepo,
		idGen:    idGen,
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

	u.Audit.Update(0)
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
