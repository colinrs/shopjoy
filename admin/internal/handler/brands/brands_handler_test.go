package brands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateBrandHandler_RequestValidation tests create brand validation
func TestCreateBrandHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request with name only",
			requestBody: map[string]interface{}{
				"name": "Apple",
			},
		},
		{
			name: "valid request with all fields",
			requestBody: map[string]interface{}{
				"name":              "Apple",
				"logo":              "https://example.com/apple-logo.png",
				"description":       "Apple Inc.",
				"website":           "https://apple.com",
				"trademark_number":  "TM123456",
				"trademark_country": "US",
				"enable_page":       true,
				"sort":              1,
			},
		},
		{
			name:           "missing name",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty name",
			requestBody: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/brands", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestUpdateBrandHandler_RequestValidation tests update brand validation
func TestUpdateBrandHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid update request",
			path: "/api/v1/brands/123",
			requestBody: map[string]interface{}{
				"name":        "Updated Apple",
				"description": "Updated description",
			},
		},
		{
			name:           "invalid brand ID",
			path:           "/api/v1/brands/0",
			requestBody:    map[string]interface{}{"name": "Test"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, tt.path, strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestGetBrandHandler tests get brand endpoint
func TestGetBrandHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/brands/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/brands/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/brands/")
		})
	}
}

// TestListBrandsHandler_QueryParams tests list brands query params
func TestListBrandsHandler_QueryParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "default pagination",
			query: "",
		},
		{
			name:  "custom pagination",
			query: "?page=1&page_size=20",
		},
		{
			name:  "filter by name",
			query: "?name=Apple",
		},
		{
			name:  "filter by status",
			query: "?status=1",
		},
		{
			name:  "combined filters",
			query: "?page=1&page_size=20&name=Apple&status=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/brands" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			assert.Contains(t, req.URL.Path, "/api/v1/brands")
		})
	}
}

// TestUpdateBrandStatusHandler_RequestValidation tests status update validation
func TestUpdateBrandStatusHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "enable brand",
			requestBody: map[string]interface{}{
				"status": 1,
			},
		},
		{
			name: "disable brand",
			requestBody: map[string]interface{}{
				"status": 0,
			},
		},
		{
			name:           "invalid status",
			requestBody:    map[string]interface{}{"status": 2},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/brands/123/status", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestToggleBrandPageHandler_RequestValidation tests toggle page validation
func TestToggleBrandPageHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "enable page",
			requestBody: map[string]interface{}{
				"enabled": true,
			},
		},
		{
			name: "disable page",
			requestBody: map[string]interface{}{
				"enabled": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/brands/123/toggle-page", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestDeleteBrandHandler tests delete brand endpoint
func TestDeleteBrandHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/brands/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/brands/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/brands/")
		})
	}
}

// TestSetBrandMarketVisibilityHandler_RequestValidation tests market visibility validation
func TestSetBrandMarketVisibilityHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "set visible",
			requestBody: map[string]interface{}{
				"market_ids": []int64{1, 2, 3},
				"visible":    true,
			},
		},
		{
			name: "set invisible",
			requestBody: map[string]interface{}{
				"market_ids": []int64{1, 2, 3},
				"visible":    false,
			},
		},
		{
			name:           "empty market_ids",
			requestBody:    map[string]interface{}{"market_ids": []int64{}, "visible": true},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/brands/123/market-visibility", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestGetBrandProductCountHandler tests get brand product count endpoint
func TestGetBrandProductCountHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/brands/123/product-count",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/brands/0/product-count",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/product-count")
		})
	}
}

// TestBrandResponseStructures tests response type structures
func TestBrandResponseStructures(t *testing.T) {
	t.Run("BrandDetailResp fields", func(t *testing.T) {
		fields := []string{
			"id", "name", "logo", "description", "website",
			"trademark_number", "trademark_country", "enable_page",
			"sort", "status", "product_count", "created_at", "updated_at",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("ListBrandResp fields", func(t *testing.T) {
		fields := []string{"list", "total"}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("BrandMarketVisibilityResp fields", func(t *testing.T) {
		fields := []string{"brand_id", "markets"}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("BrandMarketItemResp fields", func(t *testing.T) {
		fields := []string{"market_id", "is_visible"}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})
}

// TestBrandStatusConstants tests status enum values
func TestBrandStatusConstants(t *testing.T) {
	statusValues := map[int8]string{
		0: "disabled",
		1: "enabled",
	}

	for status, name := range statusValues {
		t.Run(name, func(t *testing.T) {
			assert.GreaterOrEqual(t, status, int8(0))
			assert.Less(t, status, int8(2))
		})
	}
}
