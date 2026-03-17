# ShopJoy 业务功能开发计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans and superpowers:subagent-driven-development to implement this plan task-by-task.

**Goal:** 实现完整的电商SaaS业务功能，包括用户认证、商品管理、订单流程、支付集成

**Architecture:** 严格遵循DDD分层架构，按限界上下文逐个实现。每个上下文包含：Domain层(实体/值对象/仓储接口) → Application层(服务/DTO) → Infrastructure层(仓储实现) → Interface层(API/Handler)

**Tech Stack:** go-zero, GORM, PostgreSQL, Redis, JWT, bcrypt

---

## Phase 1: Identity & Access - 用户注册/登录/认证

### Task 1.1: 创建用户仓储实现

**Files:**
- Create: `admin/internal/infrastructure/persistence/user_repository.go`

**Step 1: 实现用户仓储接口**

```go
package persistence

import (
    "context"
    "errors"
    
    "github.com/colinrs/shopjoy/admin/internal/domain/user"
    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() user.Repository {
    return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, db *gorm.DB, u *user.User) error {
    return db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) Update(ctx context.Context, db *gorm.DB, u *user.User) error {
    return db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*user.User, error) {
    var u user.User
    err := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).First(&u).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, user.ErrUserNotFound
    }
    return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email string) (*user.User, error) {
    var u user.User
    err := db.WithContext(ctx).Where("email = ? AND tenant_id = ?", email, tenantID.Int64()).First(&u).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, user.ErrUserNotFound
    }
    return &u, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, phone string) (*user.User, error) {
    var u user.User
    err := db.WithContext(ctx).Where("phone = ? AND tenant_id = ?", phone, tenantID.Int64()).First(&u).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, user.ErrUserNotFound
    }
    return &u, err
}

func (r *UserRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query user.Query) ([]*user.User, int64, error) {
    var users []*user.User
    var total int64
    
    dbQuery := db.WithContext(ctx).Model(&user.User{}).Where("tenant_id = ?", tenantID.Int64())
    
    if query.Name != "" {
        dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
    }
    if query.Email != "" {
        dbQuery = dbQuery.Where("email = ?", query.Email)
    }
    if query.Phone != "" {
        dbQuery = dbQuery.Where("phone = ?", query.Phone)
    }
    
    if err := dbQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := dbQuery.Offset(query.Offset()).Limit(query.Limit()).Find(&users).Error
    return users, total, err
}

func (r *UserRepository) Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email, phone string) (bool, error) {
    var count int64
    err := db.WithContext(ctx).Model(&user.User{}).
        Where("tenant_id = ? AND (email = ? OR phone = ?)", tenantID.Int64(), email, phone).
        Count(&count).Error
    return count > 0, err
}
```

**Step 2: 验证编译通过**

Run: `go build ./admin/...`
Expected: 编译成功

**Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/user_repository.go
git commit -m "feat: implement user repository"
```

---

### Task 1.2: 创建用户应用服务实现

**Files:**
- Create: `admin/internal/application/user/service_impl.go`

**Step 1: 实现用户应用服务**

```go
package user

import (
    "context"
    "time"
    
    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "github.com/colinrs/shopjoy/pkg/snowflake"
    domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
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
        return nil, domain.ErrDuplicateUser
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
            u.Birthday = &t
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
            u.Birthday = &t
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
        return errors.New("passwords do not match")
    }
    
    u, err := s.userRepo.FindByID(ctx, s.db, req.TenantID, req.UserID)
    if err != nil {
        return err
    }
    
    if !u.CheckPassword(req.OldPassword) {
        return domain.ErrWrongPassword
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

func toUserResponse(u *domain.User) *UserResponse {
    resp := &UserResponse{
        ID:       u.ID,
        TenantID: u.TenantID.Int64(),
        Email:    u.Email,
        Phone:    u.Phone,
        Name:     u.Name,
        Avatar:   u.Avatar,
        Gender:   int(u.Gender),
        Status:   int(u.Status),
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
```

**Step 2: 添加errors包导入**

Edit: 在文件头部添加 `"errors"` 导入

**Step 3: 验证编译通过**

Run: `go build ./admin/...`
Expected: 编译成功

**Step 4: Commit**

```bash
git add admin/internal/application/user/service_impl.go
git commit -m "feat: implement user application service"
```

---

### Task 1.3: 创建JWT认证工具

**Files:**
- Create: `pkg/auth/jwt.go`

**Step 1: 实现JWT工具**

```go
package auth

import (
    "errors"
    "time"
    
    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "github.com/golang-jwt/jwt/v4"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
    UserID   int64             `json:"user_id"`
    TenantID shared.TenantID   `json:"tenant_id"`
    Email    string            `json:"email"`
    jwt.RegisteredClaims
}

type JWTManager struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
}

func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
    return &JWTManager{
        secretKey:     []byte(secretKey),
        accessExpiry:  accessExpiry,
        refreshExpiry: refreshExpiry,
    }
}

func (j *JWTManager) GenerateTokenPair(userID int64, tenantID shared.TenantID, email string) (accessToken, refreshToken string, err error) {
    accessToken, err = j.generateAccessToken(userID, tenantID, email)
    if err != nil {
        return "", "", err
    }
    
    refreshToken, err = j.generateRefreshToken(userID, tenantID)
    if err != nil {
        return "", "", err
    }
    
    return accessToken, refreshToken, nil
}

func (j *JWTManager) generateAccessToken(userID int64, tenantID shared.TenantID, email string) (string, error) {
    claims := Claims{
        UserID:   userID,
        TenantID: tenantID,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTManager) generateRefreshToken(userID int64, tenantID shared.TenantID) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   string(rune(userID)),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })
    
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrTokenExpired
        }
        return nil, ErrInvalidToken
    }
    
    if !token.Valid {
        return nil, ErrInvalidToken
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, ErrInvalidToken
    }
    
    return claims, nil
}
```

**Step 2: 验证编译通过**

Run: `go build ./pkg/...`
Expected: 编译成功

**Step 3: Commit**

```bash
git add pkg/auth/jwt.go
git commit -m "feat: add JWT authentication manager"
```

---

### Task 1.4: 创建API定义 - 用户认证接口

**Files:**
- Create: `admin/desc/user.api`

**Step 1: 定义用户相关API**

```go
syntax = "v1"

info (
    title:   "User API"
    desc:    "用户管理相关接口"
    version: "v1"
)

type (
    RegisterRequest {
        Email    string `json:"email" validate:"required,email"`
        Phone    string `json:"phone"`
        Password string `json:"password" validate:"required,min=6"`
        Name     string `json:"name" validate:"required"`
    }
    
    RegisterResponse {
        ID        int64  `json:"id"`
        Email     string `json:"email"`
        Name      string `json:"name"`
        CreatedAt string `json:"created_at"`
    }
    
    LoginRequest {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required"`
    }
    
    LoginResponse {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int64  `json:"expires_in"`
        User         UserInfo `json:"user"`
    }
    
    UserInfo {
        ID       int64  `json:"id"`
        Email    string `json:"email"`
        Name     string `json:"name"`
        Avatar   string `json:"avatar"`
    }
    
    GetUserRequest {
        ID int64 `path:"id"`
    }
    
    GetUserResponse {
        ID        int64  `json:"id"`
        Email     string `json:"email"`
        Phone     string `json:"phone"`
        Name      string `json:"name"`
        Avatar    string `json:"avatar"`
        Status    int    `json:"status"`
        CreatedAt string `json:"created_at"`
    }
    
    ListUsersRequest {
        Page     int    `form:"page,default=1"`
        PageSize int    `form:"page_size,default=20"`
        Name     string `form:"name,optional"`
        Email    string `form:"email,optional"`
    }
    
    ListUsersResponse {
        List     []*GetUserResponse `json:"list"`
        Total    int64              `json:"total"`
        Page     int                `json:"page"`
        PageSize int                `json:"page_size"`
    }
    
    UpdateUserRequest {
        ID     int64  `path:"id"`
        Name   string `json:"name"`
        Avatar string `json:"avatar,optional"`
    }
    
    ChangePasswordRequest {
        OldPassword     string `json:"old_password"`
        NewPassword     string `json:"new_password" validate:"min=6"`
        ConfirmPassword string `json:"confirm_password"`
    }
)

service admin-api {
    @handler RegisterHandler
    post /users/register (RegisterRequest) returns (RegisterResponse)
    
    @handler LoginHandler
    post /users/login (LoginRequest) returns (LoginResponse)
    
    @handler GetUserHandler
    get /users/:id (GetUserRequest) returns (GetUserResponse)
    
    @handler ListUsersHandler
    get /users (ListUsersRequest) returns (ListUsersResponse)
    
    @handler UpdateUserHandler
    put /users/:id (UpdateUserRequest) returns (GetUserResponse)
    
    @handler ChangePasswordHandler
    put /users/password (ChangePasswordRequest) returns ()
}
```

**Step 2: 更新admin.api引入user.api**

Edit: `admin/desc/admin.api`

Add import at top:
```go
import "user.api"
```

**Step 3: 生成API代码**

Run: `cd admin && make api`
Expected: 代码生成成功

**Step 4: Commit**

```bash
git add admin/desc/user.api admin/desc/admin.api
git commit -m "feat: add user API definitions"
```

---

### Task 1.5: 实现登录/注册Handler

**Files:**
- Modify: `admin/internal/handler/login_handler.go` (create)
- Modify: `admin/internal/handler/register_handler.go` (create)

**Step 1: 实现注册Handler**

```go
package handler

import (
    "net/http"
    
    "github.com/colinrs/shopjoy/pkg/response"
    "github.com/colinrs/shopjoy/admin/internal/application/user"
    "github.com/colinrs/shopjoy/admin/internal/svc"
    "github.com/colinrs/shopjoy/admin/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.RegisterRequest
        if err := httpx.Parse(r, &req); err != nil {
            response.ParamError(w, err)
            return
        }
        
        // Get tenant from context
        tenantID, err := getTenantID(r)
        if err != nil {
            response.ParamError(w, err)
            return
        }
        
        appReq := user.CreateUserRequest{
            TenantID: tenantID,
            Email:    req.Email,
            Phone:    req.Phone,
            Password: req.Password,
            Name:     req.Name,
        }
        
        resp, err := svcCtx.UserService.Register(r.Context(), appReq)
        if err != nil {
            response.Error(w, err)
            return
        }
        
        response.Success(w, types.RegisterResponse{
            ID:        resp.ID,
            Email:     resp.Email,
            Name:      resp.Name,
            CreatedAt: resp.CreatedAt,
        })
    }
}
```

**Step 2: 实现登录Handler**

```go
package handler

import (
    "net/http"
    
    "github.com/colinrs/shopjoy/pkg/auth"
    "github.com/colinrs/shopjoy/pkg/response"
    "github.com/colinrs/shopjoy/admin/internal/svc"
    "github.com/colinrs/shopjoy/admin/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.LoginRequest
        if err := httpx.Parse(r, &req); err != nil {
            response.ParamError(w, err)
            return
        }
        
        tenantID, err := getTenantID(r)
        if err != nil {
            response.ParamError(w, err)
            return
        }
        
        // Find user and verify password
        userResp, err := svcCtx.UserService.GetByEmail(r.Context(), tenantID, req.Email)
        if err != nil {
            response.Error(w, err)
            return
        }
        
        // Generate tokens
        accessToken, refreshToken, err := svcCtx.JWTManager.GenerateTokenPair(
            userResp.ID,
            tenantID,
            userResp.Email,
        )
        if err != nil {
            response.Error(w, err)
            return
        }
        
        response.Success(w, types.LoginResponse{
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
            ExpiresIn:    3600,
            User: types.UserInfo{
                ID:     userResp.ID,
                Email:  userResp.Email,
                Name:   userResp.Name,
                Avatar: userResp.Avatar,
            },
        })
    }
}
```

**Step 3: 添加辅助函数到routes.go**

Add to `admin/internal/handler/routes.go`:

```go
func getTenantID(r *http.Request) (shared.TenantID, error) {
    tenantIDStr := r.Header.Get("X-Tenant-ID")
    if tenantIDStr == "" {
        return 0, errors.New("tenant id required")
    }
    var id int64
    if _, err := fmt.Sscanf(tenantIDStr, "%d", &id); err != nil {
        return 0, err
    }
    return shared.TenantID(id), nil
}
```

**Step 4: 更新ServiceContext**

Edit: `admin/internal/svc/service_context.go`

Add fields:
```go
UserService user.Service
JWTManager  *auth.JWTManager
```

**Step 5: Commit**

```bash
git add admin/internal/handler/login_handler.go admin/internal/handler/register_handler.go
git add admin/internal/svc/service_context.go
git commit -m "feat: implement login and register handlers"
```

---

## 后续Phase概述

### Phase 2: RBAC权限管理
- 角色仓储实现
- 权限定义和中间件
- API权限控制

### Phase 3: 租户管理
- 租户仓储实现
- 租户创建和配置
- 多租户数据隔离

### Phase 4+: 商品、订单、支付等业务功能

---

**执行方式：**
1. 逐个Task执行，每个Task完成后验证并提交
2. 使用superpowers:subagent-driven-development进行并行开发
3. 每个Phase完成后进行集成测试
