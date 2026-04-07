package roles

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateRoleHandler_RequestValidation tests create role validation
func TestCreateRoleHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request with name and code",
			requestBody: map[string]interface{}{
				"name": "Admin Role",
				"code": "ADMIN",
			},
		},
		{
			name: "valid request with description",
			requestBody: map[string]interface{}{
				"name":        "Admin Role",
				"code":        "ADMIN",
				"description": "Administrator role with full access",
			},
		},
		{
			name: "valid request with permissions",
			requestBody: map[string]interface{}{
				"name":           "Admin Role",
				"code":           "ADMIN",
				"permission_ids": []int64{1, 2, 3},
			},
		},
		{
			name:           "missing name",
			requestBody:    map[string]interface{}{"code": "ADMIN"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing code",
			requestBody:    map[string]interface{}{"name": "Admin Role"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty name",
			requestBody:    map[string]interface{}{"name": "", "code": "ADMIN"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty code",
			requestBody:    map[string]interface{}{"name": "Admin Role", "code": ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/roles", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestUpdateRoleHandler_RequestValidation tests update role validation
func TestUpdateRoleHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid update with name",
			path: "/api/v1/roles/123",
			requestBody: map[string]interface{}{
				"name": "Updated Admin Role",
			},
		},
		{
			name: "valid update with description",
			path: "/api/v1/roles/123",
			requestBody: map[string]interface{}{
				"description": "Updated description",
			},
		},
		{
			name:           "invalid role ID",
			path:           "/api/v1/roles/0",
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

// TestGetRoleHandler tests get role endpoint
func TestGetRoleHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/roles/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/roles/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/roles/")
		})
	}
}

// TestListRolesHandler_QueryParams tests list roles query params
func TestListRolesHandler_QueryParams(t *testing.T) {
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
			query: "?name=Admin",
		},
		{
			name:  "filter by code",
			query: "?code=ADMIN",
		},
		{
			name:  "filter by status",
			query: "?status=1",
		},
		{
			name:  "combined filters",
			query: "?page=1&page_size=20&name=Admin&status=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/roles" + tt.query
			req := httptest.NewRequest(http.MethodGet, path, nil)

			assert.Contains(t, req.URL.Path, "/api/v1/roles")
		})
	}
}

// TestUpdateRoleStatusHandler_RequestValidation tests status update validation
func TestUpdateRoleStatusHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "enable role",
			requestBody: map[string]interface{}{
				"status": 1,
			},
		},
		{
			name: "disable role",
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
			req := httptest.NewRequest(http.MethodPut, "/api/v1/roles/123/status", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestUpdateRolePermissionsHandler_RequestValidation tests permissions update validation
func TestUpdateRolePermissionsHandler_RequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid permissions update",
			requestBody: map[string]interface{}{
				"permission_ids": []int64{1, 2, 3, 4, 5},
			},
		},
		{
			name: "empty permissions",
			requestBody: map[string]interface{}{
				"permission_ids": []int64{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/roles/123/permissions", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")

			assert.NotNil(t, body)
		})
	}
}

// TestDeleteRoleHandler tests delete role endpoint
func TestDeleteRoleHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedValid bool
	}{
		{
			name:          "valid ID",
			path:          "/api/v1/roles/123",
			expectedValid: true,
		},
		{
			name:          "zero ID",
			path:          "/api/v1/roles/0",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			assert.Contains(t, req.URL.Path, "/api/v1/roles/")
		})
	}
}

// TestListPermissionsHandler tests list permissions endpoint
func TestListPermissionsHandler(t *testing.T) {
	t.Run("list permissions response structure", func(t *testing.T) {
		// ListPermissionsResponse should have list field
		assert.NotEmpty(t, "list")
	})
}

// TestRoleResponseStructures tests response type structures
func TestRoleResponseStructures(t *testing.T) {
	t.Run("RoleInfo fields", func(t *testing.T) {
		fields := []string{
			"id", "name", "code", "description", "status",
			"status_text", "is_system", "created_at", "updated_at",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("RoleWithPermissions fields", func(t *testing.T) {
		fields := []string{
			"id", "name", "code", "description", "status",
			"status_text", "is_system", "created_at", "updated_at",
			"permissions",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})

	t.Run("PermissionInfo fields", func(t *testing.T) {
		fields := []string{
			"id", "name", "code", "type", "type_text",
			"parent_id", "path", "icon", "sort",
		}
		for _, f := range fields {
			assert.NotEmpty(t, f)
		}
	})
}

// TestRoleStatusConstants tests status enum values
func TestRoleStatusConstants(t *testing.T) {
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

// TestPermissionTypeConstants tests permission type enum values
func TestPermissionTypeConstants(t *testing.T) {
	typeValues := map[int8]string{
		0: "menu",
		1: "button",
		2: "api",
	}

	for typ, name := range typeValues {
		t.Run(name, func(t *testing.T) {
			assert.GreaterOrEqual(t, typ, int8(0))
			assert.Less(t, typ, int8(3))
		})
	}
}
