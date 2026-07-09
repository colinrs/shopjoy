package user

import (
	"testing"
	"time"

	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// TestToUserResponse_CreatedAt 验证序列化器读取的是 application.Model.CreatedAt
// (由 GORM 命名约定自动填充)，而不是 shared.AuditInfo.CreatedAt (始终为零值)。
//
// 失败原因 (修复前)：
//   - domain.User 同时嵌入 application.Model.CreatedAt 和 shared.AuditInfo.CreatedAt
//   - toUserResponse 读取 u.Audit.CreatedAt，但 GORM 读库时只会回填 u.Model.CreatedAt
//   - 结果：API 返回 "0001-01-01T00:00:00Z"
//
// 修复后：toUserResponse 必须从 u.Model.CreatedAt 读取。
func TestToUserResponse_CreatedAt(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	u := &domain.User{
		Model: application.Model{
			ID:        999,
			CreatedAt: now,
			UpdatedAt: now,
		},
		TenantID: 1,
		Email:    "test@example.com",
		Name:     "tester",
		Status:   domain.StatusActive,
		// Audit 故意留空——模拟 GORM 从 DB 读取时的状态：
		//   - Model.CreatedAt 会被 GORM 填充（DB 写入/读取都走它）
		//   - Audit.CreatedAt 始终为零值
		Audit: shared.AuditInfo{},
	}

	resp := toUserResponse(u)

	want := now.Format(time.RFC3339)
	if resp.CreatedAt == "" {
		t.Fatalf("CreatedAt is empty, want %q", want)
	}
	if resp.CreatedAt != want {
		t.Errorf("CreatedAt = %q, want %q (expected to read from application.Model.CreatedAt)", resp.CreatedAt, want)
	}
}

func TestToExtendedUserResponse_CreatedAt(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	u := &domain.User{
		Model:    application.Model{ID: 999, CreatedAt: now, UpdatedAt: now},
		TenantID: 1,
		Email:    "test@example.com",
		Name:     "tester",
		Status:   domain.StatusActive,
		Audit:    shared.AuditInfo{},
	}

	resp := toExtendedUserResponse(u)
	want := now.Format(time.RFC3339)
	if resp.CreatedAt != want {
		t.Errorf("CreatedAt = %q, want %q", resp.CreatedAt, want)
	}
}

func TestToUserDetailResponse_CreatedAt(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	u := &domain.User{
		Model:    application.Model{ID: 999, CreatedAt: now, UpdatedAt: now},
		TenantID: 1,
		Email:    "test@example.com",
		Name:     "tester",
		Status:   domain.StatusActive,
		Audit:    shared.AuditInfo{},
	}

	resp := toUserDetailResponse(u)
	want := now.Format(time.RFC3339)
	if resp.CreatedAt != want {
		t.Errorf("CreatedAt = %q, want %q", resp.CreatedAt, want)
	}
	if resp.UpdatedAt != want {
		t.Errorf("UpdatedAt = %q, want %q", resp.UpdatedAt, want)
	}
}