package persistence_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
)

// TestRoleRepository_AssignToUser_HardDelete verifies that AssignToUser issues
// a *hard* DELETE on user_roles, not a soft delete via UPDATE ... SET deleted_at.
//
// Why this matters: the user_roles PRIMARY KEY is (user_id, role_id), which
// excludes deleted_at. If the repository only sets deleted_at (GORM default for
// models embedding application.Model) instead of physically removing rows, a
// subsequent INSERT of the same (user_id, role_id) collides with the primary
// key and MySQL raises Error 1062 "Duplicate entry ... for key
// 'user_roles.PRIMARY'". This test pins the contract: re-assigning roles must
// physically replace the prior rows.
func TestRoleRepository_AssignToUser_HardDelete(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewRoleRepository()

	const userID = int64(8)
	roleIDs := []int64{1, 7, 5, 6, 4}

	mock.ExpectBegin()
	// Hard DELETE — primary contract under test.
	mock.ExpectExec("DELETE FROM `user_roles`").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, int64(len(roleIDs))))
	// Allow any number of subsequent INSERTs (GORM emits one per role_id).
	for i := range len(roleIDs) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_roles`")).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}
	mock.ExpectCommit()

	if err := repo.AssignToUser(context.Background(), gdb, userID, roleIDs); err != nil {
		t.Fatalf("AssignToUser: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestRoleRepository_AssignToRole_HardDelete guards the same contract for
// the role_permissions junction table, which has the same composite-PK layout.
func TestRoleRepository_AssignToRole_HardDelete(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewPermissionRepository()

	const roleID = int64(2)
	permIDs := []int64{10, 11, 12}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `role_permissions`").
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewResult(0, int64(len(permIDs))))
	for i := range len(permIDs) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `role_permissions`")).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}
	mock.ExpectCommit()

	if err := repo.AssignToRole(context.Background(), gdb, roleID, permIDs); err != nil {
		t.Fatalf("AssignToRole: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
