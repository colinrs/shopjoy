package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testJWTSecret secret for testing
const testJWTSecret = "test-secret-key-for-testing-only"

// mockAdminUserRepo is a mock implementation of adminuser.Repository
type mockAdminUserRepo struct {
	findByIDFunc func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error)
}

func (m *mockAdminUserRepo) Create(ctx context.Context, db *gorm.DB, user *adminuser.AdminUser) error {
	return nil
}

func (m *mockAdminUserRepo) Update(ctx context.Context, db *gorm.DB, user *adminuser.AdminUser) error {
	return nil
}

func (m *mockAdminUserRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return nil
}

func (m *mockAdminUserRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, db, id)
	}
	return nil, nil
}

func (m *mockAdminUserRepo) FindByEmail(ctx context.Context, db *gorm.DB, email string) (*adminuser.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminUserRepo) FindByUsername(ctx context.Context, db *gorm.DB, username string) (*adminuser.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminUserRepo) FindByMobile(ctx context.Context, db *gorm.DB, mobile string) (*adminuser.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminUserRepo) FindList(ctx context.Context, db *gorm.DB, query adminuser.Query) ([]*adminuser.AdminUser, int64, error) {
	return nil, 0, nil
}

func (m *mockAdminUserRepo) Exists(ctx context.Context, db *gorm.DB, email, mobile string) (bool, error) {
	return false, nil
}

func (m *mockAdminUserRepo) UpdatePassword(ctx context.Context, db *gorm.DB, id int64, hashedPassword string) error {
	return nil
}

func (m *mockAdminUserRepo) ExistsByUsername(ctx context.Context, db *gorm.DB, tenantID int64, username string) (bool, error) {
	return false, nil
}

func (m *mockAdminUserRepo) ExistsByEmail(ctx context.Context, db *gorm.DB, tenantID int64, email string) (bool, error) {
	return false, nil
}

func (m *mockAdminUserRepo) CountMainAccount(ctx context.Context, db *gorm.DB, tenantID int64) (int64, error) {
	return 0, nil
}

func (m *mockAdminUserRepo) CountByRoleID(ctx context.Context, db *gorm.DB, roleID int64) (int64, error) {
	return 0, nil
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

// generateTestToken creates a valid JWT token for testing
func generateTestToken(secret string, userID, tenantID int64, userType int, expiresAt time.Time) (string, error) {
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Type:     userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateExpiredToken creates an expired JWT token for testing
func generateExpiredToken(secret string, userID, tenantID int64, userType int) (string, error) {
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Type:     userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name           string
		authHeader     string
		setupMockRepo  func() *mockAdminUserRepo
		useNilDB       bool
		expectedUserID int64
		expectedTenant int64
		expectedType   int
		expectNext     bool
		expectUsername bool
	}{
		{
			name: "valid token with super admin type",
			authHeader: func() string {
				token, _ := generateTestToken(testJWTSecret, 100, 0, 1, time.Now().Add(time.Hour))
				return "Bearer " + token
			}(),
			setupMockRepo: func() *mockAdminUserRepo {
				return &mockAdminUserRepo{
					findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
						return &adminuser.AdminUser{Username: "admin"}, nil
					},
				}
			},
			useNilDB:       false,
			expectedUserID: 100,
			expectedTenant: 0,
			expectedType:   1,
			expectNext:     true,
			expectUsername: true,
		},
		{
			name: "valid token with tenant admin type",
			authHeader: func() string {
				token, _ := generateTestToken(testJWTSecret, 200, 500, 2, time.Now().Add(time.Hour))
				return "Bearer " + token
			}(),
			setupMockRepo: func() *mockAdminUserRepo {
				return &mockAdminUserRepo{
					findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
						return &adminuser.AdminUser{Username: "tenant_admin"}, nil
					},
				}
			},
			useNilDB:       false,
			expectedUserID: 200,
			expectedTenant: 500,
			expectedType:   2,
			expectNext:     true,
			expectUsername: true,
		},
		{
			name: "valid token with tenant sub type",
			authHeader: func() string {
				token, _ := generateTestToken(testJWTSecret, 300, 500, 3, time.Now().Add(time.Hour))
				return "Bearer " + token
			}(),
			setupMockRepo: func() *mockAdminUserRepo {
				return &mockAdminUserRepo{
					findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
						return &adminuser.AdminUser{Username: "sub_user"}, nil
					},
				}
			},
			useNilDB:       false,
			expectedUserID: 300,
			expectedTenant: 500,
			expectedType:   3,
			expectNext:     true,
			expectUsername: true,
		},
		{
			name: "valid token without repo lookup (nil db)",
			authHeader: func() string {
				token, _ := generateTestToken(testJWTSecret, 100, 0, 1, time.Now().Add(time.Hour))
				return "Bearer " + token
			}(),
			setupMockRepo: func() *mockAdminUserRepo {
				return &mockAdminUserRepo{
					findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
						return &adminuser.AdminUser{Username: "should_not_be_set"}, nil
					},
				}
			},
			useNilDB:       true, // Pass nil for db so repo lookup is skipped
			expectedUserID: 100,
			expectedTenant: 0,
			expectedType:   1,
			expectNext:     true,
			expectUsername: false, // db is nil, so repo lookup is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedContext context.Context
			nextCalled := false

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				capturedContext = r.Context()

				// Verify context values
				userID, ok := contextx.GetUserID(capturedContext)
				if !ok || userID != tt.expectedUserID {
					t.Errorf("GetUserID() = %v, ok=%v, want %v, ok=true", userID, ok, tt.expectedUserID)
				}

				tenantID, ok := contextx.GetTenantID(capturedContext)
				if !ok || tenantID != tt.expectedTenant {
					t.Errorf("GetTenantID() = %v, ok=%v, want %v, ok=true", tenantID, ok, tt.expectedTenant)
				}

				userType := contextx.GetCurrentUserType(capturedContext)
				if userType != tt.expectedType {
					t.Errorf("GetCurrentUserType() = %v, want %v", userType, tt.expectedType)
				}
			})

			var testDB *gorm.DB
			if !tt.useNilDB {
				testDB = db
			}
			middleware := NewAuthMiddleware(testJWTSecret, testDB, tt.setupMockRepo())

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)

			if !nextCalled {
				t.Error("next handler was not called")
			}

			// Verify username was (or wasn't) set from repo depending on test case
			username := contextx.GetCurrentUserName(capturedContext)
			if tt.expectUsername && username == "" {
				t.Error("GetCurrentUserName() = empty, want non-empty")
			}
			if !tt.expectUsername && username != "" {
				t.Errorf("GetCurrentUserName() = %v, want empty (repo lookup skipped)", username)
			}
		})
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	token, err := generateExpiredToken(testJWTSecret, 100, 0, 1)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	nextCalled := false

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	middleware(nextHandler).ServeHTTP(rr, req)

	if nextCalled {
		t.Error("next handler should not be called for expired token")
	}

	// The error handler returns 400 status and plain text error message
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}

	// Verify the response body contains the expected error message
	body := rr.Body.String()
	if !strings.Contains(body, "Token") && !strings.Contains(body, "过期") {
		t.Errorf("response body = %v, want contains 'Token' or '过期'", body)
	}
}

func TestAuthMiddleware_MalformedToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantNext   bool
	}{
		{
			name:       "invalid base64 token",
			authHeader: "Bearer invalid-token-that-is-not-valid-base64!!!",
			wantNext:   false,
		},
		{
			name:       "token with wrong signature",
			authHeader: func() string {
				// Generate token with different secret
				token, _ := generateTestToken("wrong-secret", 100, 0, 1, time.Now().Add(time.Hour))
				return "Bearer " + token
			}(),
			wantNext: false,
		},
		{
			name:       "empty token after Bearer",
			authHeader: "Bearer ",
			wantNext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
			})

			middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)

			if nextCalled != tt.wantNext {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNext)
			}
		})
	}
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	nextCalled := false

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header set

	rr := httptest.NewRecorder()
	middleware(nextHandler).ServeHTTP(rr, req)

	if nextCalled {
		t.Error("next handler should not be called when Authorization header is missing")
	}

	// The error handler returns 400 status and plain text error message
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}

	// Verify the response body contains the expected error message
	body := rr.Body.String()
	if !strings.Contains(body, "授权") && !strings.Contains(body, "登录") {
		t.Errorf("response body = %v, want contains '授权' or '登录'", body)
	}
}

func TestAuthMiddleware_InvalidBearerFormat(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantNext   bool
	}{
		{
			name:       "missing Bearer prefix",
			authHeader: "some-token-here",
			wantNext:   false,
		},
		{
			name:       "Basic auth instead of Bearer",
			authHeader: "Basic dXNlcjpwYXNz",
			wantNext:   false,
		},
		{
			name:       "empty Bearer prefix",
			authHeader: "Bearer",
			wantNext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
			})

			middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)

			if nextCalled != tt.wantNext {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNext)
			}
		})
	}
}

func TestAuthMiddleware_UserInfoInContext(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name      string
		userID    int64
		tenantID  int64
		userType  int
		mockRepo  *mockAdminUserRepo
		checkFunc func(t *testing.T, ctx context.Context)
	}{
		{
			name:     "user info correctly extracted for super admin",
			userID:   999,
			tenantID: 0,
			userType: 1,
			mockRepo: &mockAdminUserRepo{
				findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
					if id != 999 {
						t.Errorf("FindByID called with id=%v, want 999", id)
					}
					return &adminuser.AdminUser{Username: "superadmin"}, nil
				},
			},
			checkFunc: func(t *testing.T, ctx context.Context) {
				userID, ok := contextx.GetUserID(ctx)
				if !ok || userID != 999 {
					t.Errorf("GetUserID() = %v, ok=%v, want 999, true", userID, ok)
				}

				tenantID, ok := contextx.GetTenantID(ctx)
				if !ok || tenantID != 0 {
					t.Errorf("GetTenantID() = %v, ok=%v, want 0, true", tenantID, ok)
				}

				userType := contextx.GetCurrentUserType(ctx)
				if userType != 1 {
					t.Errorf("GetCurrentUserType() = %v, want 1", userType)
				}

				username := contextx.GetCurrentUserName(ctx)
				if username != "superadmin" {
					t.Errorf("GetCurrentUserName() = %v, want superadmin", username)
				}
			},
		},
		{
			name:     "user info correctly extracted for tenant admin",
			userID:   12345,
			tenantID: 67890,
			userType: 2,
			mockRepo: &mockAdminUserRepo{
				findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
					return &adminuser.AdminUser{Username: "tenantadmin"}, nil
				},
			},
			checkFunc: func(t *testing.T, ctx context.Context) {
				userID, ok := contextx.GetUserID(ctx)
				if !ok || userID != 12345 {
					t.Errorf("GetUserID() = %v, ok=%v, want 12345, true", userID, ok)
				}

				tenantID, ok := contextx.GetTenantID(ctx)
				if !ok || tenantID != 67890 {
					t.Errorf("GetTenantID() = %v, ok=%v, want 67890, true", tenantID, ok)
				}

				userType := contextx.GetCurrentUserType(ctx)
				if userType != 2 {
					t.Errorf("GetCurrentUserType() = %v, want 2", userType)
				}
			},
		},
		{
			name:     "username not set when repo returns error",
			userID:   100,
			tenantID: 500,
			userType: 2,
			mockRepo: &mockAdminUserRepo{
				findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
					return nil, errors.New("database error")
				},
			},
			checkFunc: func(t *testing.T, ctx context.Context) {
				// User info should still be set even if username fetch fails
				userID, ok := contextx.GetUserID(ctx)
				if !ok || userID != 100 {
					t.Errorf("GetUserID() = %v, ok=%v, want 100, true", userID, ok)
				}

				// Username should be empty because repo returned error
				username := contextx.GetCurrentUserName(ctx)
				if username != "" {
					t.Errorf("GetCurrentUserName() = %v, want empty string", username)
				}
			},
		},
		{
			name:     "username not set when repo returns nil user",
			userID:   100,
			tenantID: 500,
			userType: 2,
			mockRepo: &mockAdminUserRepo{
				findByIDFunc: func(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
					return nil, nil // User not found
				},
			},
			checkFunc: func(t *testing.T, ctx context.Context) {
				username := contextx.GetCurrentUserName(ctx)
				if username != "" {
					t.Errorf("GetCurrentUserName() = %v, want empty string", username)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedContext context.Context

			token, _ := generateTestToken(testJWTSecret, tt.userID, tt.tenantID, tt.userType, time.Now().Add(time.Hour))

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedContext = r.Context()
				if tt.checkFunc != nil {
					tt.checkFunc(t, capturedContext)
				}
			})

			middleware := NewAuthMiddleware(testJWTSecret, db, tt.mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)
		})
	}
}

func TestAuthMiddleware_TenantIDInContext(t *testing.T) {
	tests := []struct {
		name            string
		tokenTenantID   int64
		expectedTenant  int64
	}{
		{
			name:           "platform admin with tenant_id 0",
			tokenTenantID:  0,
			expectedTenant: 0,
		},
		{
			name:           "tenant admin with non-zero tenant_id",
			tokenTenantID:  12345,
			expectedTenant: 12345,
		},
		{
			name:           "large tenant_id value",
			tokenTenantID:  9007199254740991, // Max safe integer
			expectedTenant: 9007199254740991,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedContext context.Context

			token, _ := generateTestToken(testJWTSecret, 100, tt.tokenTenantID, 1, time.Now().Add(time.Hour))

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedContext = r.Context()
			})

			middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)

			tenantID, ok := contextx.GetTenantID(capturedContext)
			if !ok {
				t.Error("GetTenantID() ok = false, want true")
			}
			if tenantID != tt.expectedTenant {
				t.Errorf("GetTenantID() = %v, want %v", tenantID, tt.expectedTenant)
			}
		})
	}
}

func TestAuthMiddleware_NilRepoAndDB(t *testing.T) {
	// Test that middleware works correctly when repo and db are nil
	token, _ := generateTestToken(testJWTSecret, 100, 500, 2, time.Now().Add(time.Hour))

	var capturedContext context.Context

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
	})

	// Pass nil for both db and repo - should still work for token validation
	middleware := NewAuthMiddleware(testJWTSecret, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	middleware(nextHandler).ServeHTTP(rr, req)

	// Should still extract user info from token
	userID, ok := contextx.GetUserID(capturedContext)
	if !ok || userID != 100 {
		t.Errorf("GetUserID() = %v, ok=%v, want 100, true", userID, ok)
	}

	tenantID, ok := contextx.GetTenantID(capturedContext)
	if !ok || tenantID != 500 {
		t.Errorf("GetTenantID() = %v, ok=%v, want 500, true", tenantID, ok)
	}

	userType := contextx.GetCurrentUserType(capturedContext)
	if userType != 2 {
		t.Errorf("GetCurrentUserType() = %v, want 2", userType)
	}

	// Username should be empty because repo is nil
	username := contextx.GetCurrentUserName(capturedContext)
	if username != "" {
		t.Errorf("GetCurrentUserName() = %v, want empty string (repo is nil)", username)
	}
}

func TestClaims_JWTClaims(t *testing.T) {
	// Test that our Claims struct properly maps to JWT claims
	claims := &Claims{
		UserID:   123,
		TenantID: 456,
		Type:     2,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "test-issuer",
			Subject:   "test-subject",
		},
	}

	// Verify fields are set correctly
	if claims.UserID != 123 {
		t.Errorf("Claims.UserID = %v, want 123", claims.UserID)
	}
	if claims.TenantID != 456 {
		t.Errorf("Claims.TenantID = %v, want 456", claims.TenantID)
	}
	if claims.Type != 2 {
		t.Errorf("Claims.Type = %v, want 2", claims.Type)
	}
	if claims.Issuer != "test-issuer" {
		t.Errorf("Claims.Issuer = %v, want test-issuer", claims.Issuer)
	}
}
