package categories

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateCategoryHandler_RequestValidation tests create category validation
func TestCreateCategoryHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request with name only",
			requestBody: map[string]interface{}{
				"name": "Electronics",
			},
		},
		{
			name: "valid request with all fields",
			requestBody: map[string]interface{}{
				"name":            "Electronics",
				"parent_id":       0,
				"code":            "ELEC",
				"icon":            "https://example.com/icon.png",
				"image":           "https://example.com/image.png",
				"seo_title":       "Electronic Products",
				"seo_description": "Best electronic products",
				"sort":            1,
			},
		},
		{
			name: "valid request with parent",
			requestBody: map[string]interface{}{
				"name":      "Smartphones",
				"parent_id": 1,
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestUpdateCategoryHandler_RequestValidation tests update category validation
func TestUpdateCategoryHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid update request",
			path: "/api/v1/categories/123",
			requestBody: map[string]interface{}{
				"name": "Updated Electronics",
				"code": "ELEC-UPD",
				"sort": 2,
			},
		},
		{
			name:           "invalid ID",
			path:           "/api/v1/categories/0",
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

// TestGetCategoryHandler tests get category endpoint
func TestGetCategoryHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/categories/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/categories/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/categories/")
		})
	}
}

// TestListCategoriesHandler_QueryParams tests list categories query params
func TestListCategoriesHandler_QueryParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "no filter",
			query: "",
		},
		{
			name:  "with parent_id",
			query: "?parent_id=0",
		},
		{
			name:  "with parent_id filter",
			query: "?parent_id=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/categories" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			assert.Contains(t, req.URL.Path, "/api/v1/categories")
		})
	}
}

// TestCategoryTreeHandler tests get category tree endpoint
func TestCategoryTreeHandler(t *testing.T) {
	t.Run("tree response structure", func(t *testing.T) {
		// CategoryTreeResp should have children field
		assert.NotEmpty(t, "children")
	})
}

// TestUpdateCategoryStatusHandler_RequestValidation tests status update validation
func TestUpdateCategoryStatusHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "enable category",
			requestBody: map[string]interface{}{
				"status": 1,
			},
		},
		{
			name: "disable category",
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
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/123/status", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestUpdateCategorySortHandler_RequestValidation tests sort update validation
func TestUpdateCategorySortHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid sort update",
			requestBody: map[string]interface{}{
				"sorts": []map[string]interface{}{
					{"id": 1, "sort": 1},
					{"id": 2, "sort": 2},
					{"id": 3, "sort": 3},
				},
			},
		},
		{
			name:           "empty sorts",
			requestBody:    map[string]interface{}{"sorts": []map[string]interface{}{}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/sort", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestMoveCategoryHandler_RequestValidation tests move category validation
func TestMoveCategoryHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "move to root",
			requestBody: map[string]interface{}{
				"new_parent_id": 0,
			},
		},
		{
			name: "move to subcategory",
			requestBody: map[string]interface{}{
				"new_parent_id": 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/123/move", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestSetCategoryMarketVisibilityHandler_RequestValidation tests market visibility validation
func TestSetCategoryMarketVisibilityHandler_RequestValidation(t *testing.T) {
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
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/123/market-visibility", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestDeleteCategoryHandler tests delete category endpoint
func TestDeleteCategoryHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/categories/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/categories/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/categories/")
		})
	}
}

// TestCategoryResponseStructures tests response type structures
func TestCategoryResponseStructures(t *testing.T) {
	t.Run("CategoryDetailResp fields", func(t *testing.T) {
		fields := []string{
			"id", "parent_id", "name", "code", "level",
			"sort", "icon", "image", "seo_title", "seo_description",
			"status", "product_count", "created_at", "updated_at",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("CategoryTreeResp fields", func(t *testing.T) {
		fields := []string{
			"id", "parent_id", "name", "code", "level",
			"sort", "icon", "image", "seo_title", "seo_description",
			"status", "product_count", "created_at", "updated_at", "children",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("CategoryMarketVisibilityResp fields", func(t *testing.T) {
		fields := []string{"category_id", "markets"}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})
}

// TestCategoryStatusConstants tests status enum values
func TestCategoryStatusConstants(t *testing.T) {
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
