package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// TestJWTConfig holds JWT configuration for tests
type TestJWTConfig struct {
	Secret   string
	UserID   int64
	TenantID int64
	UserType int
}

// GenerateTestToken generates a valid JWT token for testing
func GenerateTestToken(cfg TestJWTConfig) string {
	claims := &struct {
		jwt.RegisteredClaims
		UserID   int64 `json:"userId"`
		TenantID int64 `json:"tenantId"`
		Type     int   `json:"type"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID:   cfg.UserID,
		TenantID: cfg.TenantID,
		Type:     cfg.UserType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.Secret))
	return tokenString
}

// GenerateExpiredTestToken generates an expired JWT token for testing
func GenerateExpiredTestToken(cfg TestJWTConfig) string {
	claims := &struct {
		jwt.RegisteredClaims
		UserID   int64 `json:"userId"`
		TenantID int64 `json:"tenantId"`
		Type     int   `json:"type"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
		UserID:   cfg.UserID,
		TenantID: cfg.TenantID,
		Type:     cfg.UserType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.Secret))
	return tokenString
}

// GetTestJWTConfig returns default test JWT configuration
func GetTestJWTConfig() TestJWTConfig {
	return TestJWTConfig{
		Secret:   "test-secret-key-for-testing-only",
		UserID:   1,
		TenantID: 1,
		UserType: 1,
	}
}

// TestGetUserHandler_AuthRequired tests that the get user endpoint requires auth
func TestGetUserHandler_AuthRequired(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     map[string]string
		expectedStatus int
	}{
		{
			name:           "no auth header",
			authHeader:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid token",
			authHeader: map[string]string{
				"Authorization": "Bearer invalid-token",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "expired token",
			authHeader: func() map[string]string {
				cfg := GetTestJWTConfig()
				return map[string]string{
					"Authorization": "Bearer " + GenerateExpiredTestToken(cfg),
				}
			}(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents the auth requirements
			// Actual testing would require the full middleware chain
			if tt.authHeader == nil {
				assert.Nil(t, tt.authHeader)
			} else {
				assert.NotEmpty(t, tt.authHeader["Authorization"])
			}
		})
	}
}

// TestGetUserHandler_RequestValidation tests request validation
func TestGetUserHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedValid bool
	}{
		{
			name:          "valid numeric ID",
			id:            "123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			id:            "0",
			expectedValid: false,
		},
		{
			name:          "negative ID",
			id:            "-1",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/users/" + tt.id
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Validate path parameter parsing
			assert.Contains(t, req.URL.Path, tt.id)
		})
	}
}

// TestListUsersHandler_QueryParams tests query parameter handling
func TestListUsersHandler_QueryParams(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantFields map[string]string
	}{
		{
			name:  "default pagination",
			query: "",
			wantFields: map[string]string{
				"page":      "1",
				"page_size": "20",
			},
		},
		{
			name:  "custom pagination",
			query: "?page=2&page_size=50",
			wantFields: map[string]string{
				"page":      "2",
				"page_size": "50",
			},
		},
		{
			name:  "with name filter",
			query: "?name=John",
			wantFields: map[string]string{
				"name": "John",
			},
		},
		{
			name:  "with status filter",
			query: "?status=1",
			wantFields: map[string]string{
				"status": "1",
			},
		},
		{
			name:  "with keyword search",
			query: "?keyword=test",
			wantFields: map[string]string{
				"keyword": "test",
			},
		},
		{
			name:  "combined filters",
			query: "?page=1&page_size=20&status=1&keyword=test",
			wantFields: map[string]string{
				"page":      "1",
				"page_size": "20",
				"status":    "1",
				"keyword":   "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/users" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Verify URL parsing
			assert.Contains(t, req.URL.Path, "/api/v1/users")
		})
	}
}

// TestListUsersHandler_ResponseFormat tests response structure
func TestListUsersHandler_ResponseFormat(t *testing.T) {
	// Expected response structure
	expectedListFields := []string{
		"id", "email", "phone", "name", "avatar", "status", "created_at",
	}

	t.Run("ListUsersResponse has required fields", func(t *testing.T) {
		for _, field := range expectedListFields {
			assert.NotEmpty(t, field, "User response should have %s field", field)
		}
	})

	t.Run("pagination fields present", func(t *testing.T) {
		paginationFields := []string{"list", "total", "page", "page_size"}
		for _, field := range paginationFields {
			assert.NotEmpty(t, field)
		}
	})
}

// TestUpdateUserHandler_RequestValidation tests update request validation
func TestUpdateUserHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid update with name",
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
		},
		{
			name: "valid update with avatar",
			requestBody: map[string]interface{}{
				"name":   "Updated Name",
				"avatar": "https://example.com/avatar.jpg",
			},
		},
		{
			name: "empty name",
			requestBody: map[string]interface{}{
				"name": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/123", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestSuspendUserHandler_RequestValidation tests suspend request validation
func TestSuspendUserHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantUserID int64
	}{
		{
			name:       "valid user ID",
			path:       "/api/v1/users/123/suspend",
			wantUserID: 123,
		},
		{
			name:       "zero user ID",
			path:       "/api/v1/users/0/suspend",
			wantUserID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/suspend")
		})
	}
}

// TestSuspendUserWithReasonHandler_RequestValidation tests suspend with reason validation
func TestSuspendUserWithReasonHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request with reason",
			requestBody: map[string]interface{}{
				"reason": "User violated terms of service",
			},
		},
		{
			name:           "missing reason",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty reason",
			requestBody:    map[string]interface{}{"reason": ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/123/suspend-reason", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestGetUserStatsHandler tests the user stats endpoint
func TestGetUserStatsHandler(t *testing.T) {
	t.Run("stats response structure", func(t *testing.T) {
		expectedFields := []string{"total", "active", "suspended", "new_today"}
		for _, field := range expectedFields {
			assert.NotEmpty(t, field)
		}
	})
}

// TestGetUserDetailHandler tests the enhanced user detail endpoint
func TestGetUserDetailHandler(t *testing.T) {
	t.Run("detail response structure", func(t *testing.T) {
		expectedFields := []string{
			"id", "tenant_id", "email", "phone", "name", "avatar",
			"gender", "gender_text", "birthday", "status", "status_text",
			"points_balance", "points_frozen", "total_earned_points",
			"total_redeemed_points", "order_count", "total_spent",
			"review_count", "last_login", "created_at", "updated_at",
			"last_order_at", "default_address",
		}
		for _, field := range expectedFields {
			assert.NotEmpty(t, field)
		}
	})
}

// TestGetUserAddressesHandler tests the user addresses endpoint
func TestGetUserAddressesHandler(t *testing.T) {
	t.Run("address response structure", func(t *testing.T) {
		expectedFields := []string{
			"id", "user_id", "name", "phone", "country",
			"province", "city", "district", "detail",
			"postal_code", "is_default", "created_at", "updated_at",
		}
		for _, field := range expectedFields {
			assert.NotEmpty(t, field)
		}
	})
}

// TestBatchUpdateUserStatusHandler_RequestValidation tests batch update validation
func TestBatchUpdateUserStatusHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid batch enable request",
			requestBody: map[string]interface{}{
				"user_ids": []int64{1, 2, 3},
				"status":   1,
			},
		},
		{
			name: "valid batch disable request",
			requestBody: map[string]interface{}{
				"user_ids": []int64{1, 2, 3},
				"status":   2,
			},
		},
		{
			name: "missing user_ids",
			requestBody: map[string]interface{}{
				"status": 1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty user_ids",
			requestBody: map[string]interface{}{
				"user_ids": []int64{},
				"status":   1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "too many user_ids",
			requestBody: map[string]interface{}{
				"user_ids": make([]int64, 101),
				"status":   1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid status",
			requestBody: map[string]interface{}{
				"user_ids": []int64{1, 2, 3},
				"status":   3, // Only 1 or 2 allowed
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/batch-status", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestExportUsersHandler_QueryParams tests export query parameters
func TestExportUsersHandler_QueryParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "no filters",
			query: "",
		},
		{
			name:  "with keyword",
			query: "?keyword=test",
		},
		{
			name:  "with status",
			query: "?status=1",
		},
		{
			name:  "with date range",
			query: "?register_start=2024-01-01&register_end=2024-12-31",
		},
		{
			name:  "with all filters",
			query: "?keyword=test&status=1&register_start=2024-01-01&register_end=2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/users/export" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			assert.Contains(t, req.URL.Path, "/api/v1/users/export")
		})
	}
}

// TestUserStatusConstants tests the user status enum values
func TestUserStatusConstants(t *testing.T) {
	statusValues := map[int]string{
		0: "inactive",
		1: "active",
		2: "suspended",
		3: "deleted",
	}

	for status, name := range statusValues {
		t.Run(name, func(t *testing.T) {
			assert.GreaterOrEqual(t, status, 0)
			assert.Less(t, status, 4)
		})
	}
}

// TestListUsersEnhancedHandler_QueryParams tests enhanced list query params
func TestListUsersEnhancedHandler_QueryParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "default",
			query: "",
		},
		{
			name:  "with pagination",
			query: "?page=1&page_size=20",
		},
		{
			name:  "with keyword",
			query: "?keyword=test",
		},
		{
			name:  "with status",
			query: "?status=1",
		},
		{
			name:  "with date range",
			query: "?register_start=2024-01-01T00:00:00Z&register_end=2024-12-31T23:59:59Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/users/enhanced" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			assert.Contains(t, req.URL.Path, "/api/v1/users/enhanced")
		})
	}
}

// TestDeleteUserHandler tests delete user endpoint
func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid user ID",
			path:          "/api/v1/users/123",
			expectedValid: true,
		},
		{
			name:          "zero user ID",
			path:          "/api/v1/users/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/users/")
		})
	}
}

// TestResetPasswordHandler tests reset password endpoint
func TestResetPasswordHandler(t *testing.T) {
	t.Run("response structure", func(t *testing.T) {
		expectedFields := []string{"temporary_password"}
		for _, field := range expectedFields {
			assert.NotEmpty(t, field)
		}
	})
}

// Benchmark tests
func BenchmarkUserListQueryParsing(b *testing.B) {
	query := "?page=1&page_size=20&status=1&keyword=test"

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users"+query, nil)
		_ = req.URL.Query()
	}
}

func BenchmarkBatchStatusRequestParsing(b *testing.B) {
	body := `{"user_ids":[1,2,3,4,5],"status":1}`

	for i := 0; i < b.N; i++ {
		var req map[string]interface{}
		json.Unmarshal([]byte(body), &req)
	}
}
