package adminuser

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

const (
	accessTokenExpire  = 2 * time.Hour
	refreshTokenExpire = 7 * 24 * time.Hour
)

type service struct {
	repo      adminuser.Repository
	roleRepo  role.Repository
	db        *gorm.DB
	jwtSecret []byte
}

// NewService 创建管理员应用服务
func NewService(repo adminuser.Repository, roleRepo role.Repository, db *gorm.DB, jwtSecret string) Service {
	return &service{
		repo:      repo,
		roleRepo:  roleRepo,
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *service) RegisterTenantAdmin(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.repo.Exists(ctx, s.db, req.Email, req.Mobile)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, code.ErrAdminDuplicateUser
	}

	// 创建管理员账号
	user := &adminuser.AdminUser{
		Email:    req.Email,
		Mobile:   req.Mobile,
		RealName: req.RealName,
		Type:     adminuser.TypeTenantAdmin,
		Status:   adminuser.StatusActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 事务：创建管理员 + 创建店铺
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(ctx, tx, user); err != nil {
			return err
		}
		// TODO: 创建店铺，将 tenant_id 关联到用户
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 生成 token 并返回
	return s.generateLoginResponse(user)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var user *adminuser.AdminUser
	var err error

	// 尝试通过邮箱、用户名、手机号查找
	user, err = s.repo.FindByEmail(ctx, s.db, req.Account)
	if err != nil {
		user, err = s.repo.FindByUsername(ctx, s.db, req.Account)
		if err != nil {
			user, err = s.repo.FindByMobile(ctx, s.db, req.Account)
			if err != nil {
				return nil, code.ErrAdminUserNotFound
			}
		}
	}

	// 检查状态
	if !user.CanLogin() {
		return nil, code.ErrAdminAccountDisabled
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, code.ErrAdminWrongPassword
	}

	// 更新最后登录时间
	user.UpdateLastLogin(req.IP)
	if err := s.repo.Update(ctx, s.db, user); err != nil {
		return nil, err
	}

	return s.generateLoginResponse(user)
}

func (s *service) GetByID(ctx context.Context, id int64) (*AdminUserResponse, error) {
	user, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(user), nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*AdminUserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, s.db, email)
	if err != nil {
		return nil, err
	}
	return s.toResponse(user), nil
}

func (s *service) UpdateProfile(ctx context.Context, req UpdateProfileRequest) (*AdminUserResponse, error) {
	user, err := s.repo.FindByID(ctx, s.db, req.UserID)
	if err != nil {
		return nil, err
	}

	// 权限检查：只能修改自己的资料，或者平台超管可以修改所有人的
	// TODO: 获取当前操作者信息并检查权限

	user.UpdateProfile(req.RealName, req.Avatar, req.Mobile)
	if req.Email != "" && req.Email != user.Email {
		user.Email = req.Email
	}

	if err := s.repo.Update(ctx, s.db, user); err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

func (s *service) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return code.ErrAdminPasswordMismatch
	}

	user, err := s.repo.FindByID(ctx, s.db, req.UserID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return code.ErrAdminWrongPassword
	}

	// 设置新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	return s.repo.Update(ctx, s.db, user)
}

func (s *service) List(ctx context.Context, operatorID int64, req ListRequest) (*ListResponse, error) {
	// 获取操作者信息
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return nil, err
	}

	// 权限检查：非平台超管只能查看自己租户的管理员
	query := adminuser.Query{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}

	if req.Type > 0 {
		query.Type = adminuser.Type(req.Type)
	}
	if req.Status > 0 {
		query.Status = adminuser.Status(req.Status)
	}

	if operator.IsSuperAdmin() {
		// 平台超管可以查看所有，或指定租户
		query.TenantID = req.TenantID
	} else {
		// 商家管理员只能查看自己租户的
		query.TenantID = operator.TenantID
	}

	list, total, err := s.repo.FindList(ctx, s.db, query)
	if err != nil {
		return nil, err
	}

	respList := make([]*AdminUserResponse, len(list))
	for i, u := range list {
		respList[i] = s.toResponse(u)
	}

	return &ListResponse{
		List:     respList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *service) Disable(ctx context.Context, operatorID, targetID int64) error {
	if operatorID == targetID {
		return code.ErrAdminCannotDeleteSelf
	}

	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return err
	}

	target, err := s.repo.FindByID(ctx, s.db, targetID)
	if err != nil {
		return err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return code.ErrAdminPermissionDenied
	}

	if err := target.Disable(); err != nil {
		return err
	}

	return s.repo.Update(ctx, s.db, target)
}

func (s *service) Enable(ctx context.Context, operatorID, targetID int64) error {
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return err
	}

	target, err := s.repo.FindByID(ctx, s.db, targetID)
	if err != nil {
		return err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return code.ErrAdminPermissionDenied
	}

	if err := target.Enable(); err != nil {
		return err
	}

	return s.repo.Update(ctx, s.db, target)
}

// Create 创建管理员
func (s *service) Create(ctx context.Context, operatorID int64, req CreateAdminUserRequest) (*AdminUserResponse, error) {
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return nil, err
	}

	// 确定目标租户ID
	targetTenantID := req.TenantID
	if !operator.IsSuperAdmin() {
		// 非平台超管只能在自己租户创建用户
		targetTenantID = operator.TenantID
	}

	// BR-010: 用户名在租户内唯一
	if req.Username != "" {
		exists, err := s.repo.ExistsByUsername(ctx, s.db, targetTenantID, req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, code.ErrAdminUsernameExists
		}
	}

	// BR-011: 邮箱在租户内唯一
	if req.Email != "" {
		exists, err := s.repo.ExistsByEmail(ctx, s.db, targetTenantID, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, code.ErrAdminEmailExists
		}
	}

	// BR-012: 每个租户只能有一个主账号 (Type=2)
	if adminuser.Type(req.Type) == adminuser.TypeTenantAdmin && targetTenantID > 0 {
		count, err := s.repo.CountMainAccount(ctx, s.db, targetTenantID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, code.ErrAdminMainAccountExists
		}
	}

	// 创建用户
	user := &adminuser.AdminUser{
		TenantID: targetTenantID,
		Username: req.Username,
		Email:    req.Email,
		Mobile:   req.Mobile,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Type:     adminuser.Type(req.Type),
		Status:   adminuser.StatusActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 事务：创建用户 + 分配角色
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(ctx, tx, user); err != nil {
			return err
		}

		// 分配角色
		if len(req.RoleIDs) > 0 {
			if err := s.roleRepo.AssignToUser(ctx, tx, shared.TenantID(targetTenantID), user.ID, req.RoleIDs); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

// Update 更新管理员
func (s *service) Update(ctx context.Context, operatorID int64, req UpdateAdminUserRequest) (*AdminUserResponse, error) {
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return nil, err
	}

	target, err := s.repo.FindByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return nil, code.ErrAdminPermissionDenied
	}

	// 更新资料
	target.UpdateProfile(req.RealName, req.Avatar, req.Mobile)

	// 更新邮箱（如果提供且不同）
	if req.Email != "" && req.Email != target.Email {
		// BR-011: 邮箱在租户内唯一
		exists, err := s.repo.ExistsByEmail(ctx, s.db, target.TenantID, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, code.ErrAdminEmailExists
		}
		target.Email = req.Email
	}

	if err := s.repo.Update(ctx, s.db, target); err != nil {
		return nil, err
	}

	return s.toResponse(target), nil
}

// Delete 删除管理员
func (s *service) Delete(ctx context.Context, operatorID, targetID int64) error {
	// BR-013: 不能删除自己
	if operatorID == targetID {
		return code.ErrAdminCannotDeleteSelf
	}

	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return err
	}

	target, err := s.repo.FindByID(ctx, s.db, targetID)
	if err != nil {
		return err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return code.ErrAdminPermissionDenied
	}

	return s.repo.Delete(ctx, s.db, targetID)
}

// AssignRoles 分配角色
func (s *service) AssignRoles(ctx context.Context, operatorID int64, req AssignRolesRequest) error {
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return err
	}

	target, err := s.repo.FindByID(ctx, s.db, req.AdminUserID)
	if err != nil {
		return err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return code.ErrAdminPermissionDenied
	}

	return s.roleRepo.AssignToUser(ctx, s.db, shared.TenantID(target.TenantID), req.AdminUserID, req.RoleIDs)
}

// ResetPassword 重置密码
func (s *service) ResetPassword(ctx context.Context, operatorID, targetID int64) (*ResetPasswordResponse, error) {
	operator, err := s.repo.FindByID(ctx, s.db, operatorID)
	if err != nil {
		return nil, err
	}

	target, err := s.repo.FindByID(ctx, s.db, targetID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if !operator.CanManageTenant(target.TenantID) {
		return nil, code.ErrAdminPermissionDenied
	}

	// 生成随机临时密码（12位）
	tempPassword, err := generateRandomPassword(12)
	if err != nil {
		return nil, err
	}

	// 设置新密码
	if err := target.SetPassword(tempPassword); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, s.db, target); err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{
		TemporaryPassword: tempPassword,
	}, nil
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// 生成登录响应
func (s *service) generateLoginResponse(user *adminuser.AdminUser) (*LoginResponse, error) {
	accessToken, err := s.generateToken(user.ID, user.TenantID, user.Type, false)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.TenantID, user.Type, true)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenExpire.Seconds()),
		User:         *s.toResponse(user),
	}, nil
}

// 生成 JWT token
func (s *service) generateToken(userID, tenantID int64, userType adminuser.Type, isRefresh bool) (string, error) {
	expire := accessTokenExpire
	if isRefresh {
		expire = refreshTokenExpire
	}

	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"type":      userType,
		"exp":       time.Now().Add(expire).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// 转换为响应对象
func (s *service) toResponse(user *adminuser.AdminUser) *AdminUserResponse {
	resp := &AdminUserResponse{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		RealName:  user.RealName,
		Avatar:    user.Avatar,
		Type:      int(user.Type),
		TypeText:  user.Type.String(),
		Status:    int(user.Status),
		CreatedAt: time.Unix(user.CreatedAt, 0).UTC().Format(time.RFC3339),
	}

	if user.LastLoginAt != nil {
		resp.LastLoginAt = time.Unix(*user.LastLoginAt, 0).UTC().Format(time.RFC3339)
	}

	return resp
}
