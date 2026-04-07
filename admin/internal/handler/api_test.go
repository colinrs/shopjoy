package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	claims := &Claims{
		UserID:   cfg.UserID,
		TenantID: cfg.TenantID,
		Type:     cfg.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.Secret))
	return tokenString
}

// GenerateExpiredTestToken generates an expired JWT token for testing
func GenerateExpiredTestToken(cfg TestJWTConfig) string {
	claims := &Claims{
		UserID:   cfg.UserID,
		TenantID: cfg.TenantID,
		Type:     cfg.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.Secret))
	return tokenString
}

// ContextWithAuth creates a context with auth values set
func ContextWithAuth(ctx context.Context, userID, tenantID int64, userType int) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID, userID)
	ctx = context.WithValue(ctx, contextKeyTenantID, tenantID)
	ctx = context.WithValue(ctx, contextKeyUserType, userType)
	return ctx
}

// context keys for testing
var (
	contextKeyUserID   = "test_user_id"
	contextKeyTenantID = "test_tenant_id"
	contextKeyUserType = "test_user_type"
)

// API Test Request/Response helpers

// APIResponse represents a standard API response
type APIResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data,omitempty"`
}

// TestRequest holds test request data
type TestRequest struct {
	Method string
	Path   string
	Body   any
	Header map[string]string
}

// ExecuteRequest executes an HTTP request and returns the response
func ExecuteRequest(handler http.Handler, req TestRequest) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = json.Marshal(req.Body)
	}

	httpReq := httptest.NewRequest(req.Method, req.Path, bytes.NewReader(bodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httpReq)
	return rr
}

// ExecuteRequestWithQuery executes an HTTP request with query parameters
func ExecuteRequestWithQuery(handler http.Handler, method, path string, query map[string]string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	httpReq := httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	// Add query parameters
	q := httpReq.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httpReq)
	return rr
}

// AssertResponseOK asserts that the response is successful (2xx)
func AssertResponseOK(t *testing.T, rr *httptest.ResponseRecorder) {
	assert.True(t, rr.Code >= 200 && rr.Code < 300, "Expected 2xx status, got %d. Body: %s", rr.Code, rr.Body.String())
}

// AssertResponseError asserts that the response is an error (4xx/5xx)
func AssertResponseError(t *testing.T, rr *httptest.ResponseRecorder) {
	assert.True(t, rr.Code >= 400, "Expected error status, got %d", rr.Code)
}

// AssertResponseStatus asserts the exact response status code
func AssertResponseStatus(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, rr.Code, "Status mismatch. Body: %s", rr.Body.String())
}

// AssertJSONField asserts a field in the JSON response
func AssertJSONField(t *testing.T, body []byte, field string, expectedValue any) {
	var resp map[string]any
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		// Try to get directly from response
		val, exists := resp[field]
		assert.True(t, exists, "Field %s not found in response", field)
		assert.Equal(t, expectedValue, val)
		return
	}

	val, exists := data[field]
	assert.True(t, exists, "Field %s not found in response data", field)
	assert.Equal(t, expectedValue, val)
}

// ParseAPIResponse parses an API response
func ParseAPIResponse(t *testing.T, body []byte) APIResponse {
	var resp APIResponse
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	return resp
}

// ============================================================================
// Table-driven test helpers
// ============================================================================

// TestCase represents a test case for table-driven tests
type TestCase struct {
	Name           string
	Request        TestRequest
	ExpectedStatus int
	Setup          func(*testing.T)
	CheckResponse  func(*testing.T, *httptest.ResponseRecorder)
}

// RunTableTests runs a table of test cases
func RunTableTests(t *testing.T, handler http.HandlerFunc, testCases []TestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Setup != nil {
				tc.Setup(t)
			}

			rr := ExecuteRequest(handler, tc.Request)

			if tc.ExpectedStatus > 0 {
				AssertResponseStatus(t, rr, tc.ExpectedStatus)
			}

			if tc.CheckResponse != nil {
				tc.CheckResponse(t, rr)
			}
		})
	}
}

// ValidationTestCase represents a validation error test case
type ValidationTestCase struct {
	Name        string
	Request     TestRequest
	ExpectedMsg string
	Setup       func(*testing.T)
}

// RunValidationTests runs validation test cases
func RunValidationTests(t *testing.T, handler http.HandlerFunc, testCases []ValidationTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Setup != nil {
				tc.Setup(t)
			}

			rr := ExecuteRequest(handler, tc.Request)
			AssertResponseError(t, rr)

			body := rr.Body.Bytes()
			var resp map[string]any
			err := json.Unmarshal(body, &resp)
			assert.NoError(t, err)

			// Check for error message
			if tc.ExpectedMsg != "" {
				msg, ok := resp["msg"].(string)
				assert.True(t, ok, "Expected msg field in error response")
				assert.Contains(t, msg, tc.ExpectedMsg, "Error message should contain expected text")
			}
		})
	}
}

// ============================================================================
// Common test data factories
// ============================================================================

// ValidLoginRequest creates a valid login request
func ValidLoginRequest() map[string]any {
	return map[string]any{
		"account":  "admin@test.com",
		"password": "password123",
	}
}

// ValidCategoryRequest creates a valid category create request
func ValidCategoryRequest() map[string]any {
	return map[string]any{
		"name": "Test Category",
		"code": "TEST-CAT-001",
		"sort": 1,
	}
}

// ValidProductRequest creates a valid product create request
func ValidProductRequest() map[string]interface{} {
	return map[string]any{
		"name":        "Test Product",
		"description": "A test product description",
		"price":       "99.99",
		"currency":    "CNY",
		"category_id": 1,
		"sku":         "TEST-SKU-001",
	}
}

// ValidBrandRequest creates a valid brand create request
func ValidBrandRequest() map[string]interface{} {
	return map[string]any{
		"name":        "Test Brand",
		"description": "A test brand",
		"logo":        "https://example.com/logo.png",
	}
}

// ValidRoleRequest creates a valid role create request
func ValidRoleRequest() map[string]interface{} {
	return map[string]any{
		"name":        "Test Role",
		"code":        "TEST-ROLE-001",
		"description": "A test role",
	}
}

// ============================================================================
// Auth middleware helpers (for testing protected endpoints)
// ============================================================================

// Claims JWT claims struct (matches middleware)
type Claims struct {
	jwt.RegisteredClaims
	UserID   int64 `json:"userId"`
	TenantID int64 `json:"tenantId"`
	Type     int   `json:"type"`
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

// AuthHeader creates an authorization header with a valid token
func AuthHeader(cfg TestJWTConfig) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + GenerateTestToken(cfg),
	}
}

// AuthHeaderWithExpiredToken creates an authorization header with an expired token
func AuthHeaderWithExpiredToken(cfg TestJWTConfig) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + GenerateExpiredTestToken(cfg),
	}
}

// ============================================================================
// Error response checking
// ============================================================================

// AssertErrCode asserts the error code in the response
func AssertErrCode(t *testing.T, body []byte, expectedCode int) {
	var resp map[string]any
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)

	code, ok := resp["code"].(float64)
	assert.True(t, ok, "Expected numeric code field in error response")
	assert.Equal(t, float64(expectedCode), code, "Error code mismatch")
}

// UnwrapResponse unwraps the data field from an API response
func UnwrapResponse(t *testing.T, body []byte) map[string]any {
	var resp struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data,omitempty"`
	}
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code, "Expected success response")

	if resp.Data == nil {
		return nil
	}

	var data map[string]any
	err = json.Unmarshal(resp.Data, &data)
	assert.NoError(t, err)
	return data
}

// UnwrapResponseSlice unwraps a slice data field from an API response
func UnwrapResponseSlice(t *testing.T, body []byte) []any {
	var resp struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data,omitempty"`
	}
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code, "Expected success response")

	if resp.Data == nil {
		return nil
	}

	var data []any
	err = json.Unmarshal(resp.Data, &data)
	assert.NoError(t, err)
	return data
}
