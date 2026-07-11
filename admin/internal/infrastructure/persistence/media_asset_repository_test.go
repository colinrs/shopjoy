package persistence_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/colinrs/shopjoy/admin/internal/domain/media"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	// Use Regexp matcher: GORM emits full column-list INSERT, which exact-match
	// would require hard-coding GORM's internal column order. Regexp keeps the
	// test focused on intent (we INSERT into media_assets) and resilient to
	// GORM version changes.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return gdb, mock
}

func TestMediaAssetRepository_Insert_OK(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewMediaAssetRepository(gdb)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `media_assets`").
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Insert(context.Background(), &media.Asset{PublicID: "p1", URL: "u", TenantID: 7})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestMediaAssetRepository_Insert_Duplicate(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewMediaAssetRepository(gdb)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `media_assets`").
		WillReturnError(errors.New("Error 1062: Duplicate entry 'cloudinary|p1' for key 'uk_provider_public'"))
	mock.ExpectRollback()

	err := repo.Insert(context.Background(), &media.Asset{PublicID: "p1", URL: "u", TenantID: 7})
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}
