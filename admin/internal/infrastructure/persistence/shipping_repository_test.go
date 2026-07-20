package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/pkg/code"
)

// TestShippingRepo_FindDefaultByMarket_MarketHit pins the contract that when a
// market-scoped default template exists, FindDefaultByMarket must return it
// (i.e., it must NOT silently fall through to the market_id=0 row).
func TestShippingRepo_FindDefaultByMarket_MarketHit(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const (
		marketID     = int64(7)
		templateID   = int64(101)
		templateName = "HK-Standard"
	)

	rows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		templateID, int64(1), marketID, "HKD", templateName,
		true, true, "standard", int64(0),
		time.Now(), time.Now(), sql.NullTime{},
	)

	// Single SELECT — market-scoped row found, no fallback needed.
	// Args: marketID, is_default=true, is_active=true, LIMIT 1 (GORM's First()).
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(marketID, true, true, 1).
		WillReturnRows(rows)

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, marketID)
	if err != nil {
		t.Fatalf("FindDefaultByMarket: %v", err)
	}
	if got == nil {
		t.Fatal("expected template, got nil")
	}
	if got.ID != templateID {
		t.Errorf("expected id %d, got %d", templateID, got.ID)
	}
	if got.Name != templateName {
		t.Errorf("expected name %q, got %q", templateName, got.Name)
	}
	if got.MarketID != marketID {
		t.Errorf("expected market_id %d, got %d", marketID, got.MarketID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindDefaultByMarket_FallbackToZero pins the fallback
// contract: when no market-scoped default exists, FindDefaultByMarket MUST
// fall through to the market_id=0 (全市场通用) default. This is the
// cross-market safety net that prevents the storefront from having NO
// fallback template when a market-specific one is missing.
func TestShippingRepo_FindDefaultByMarket_FallbackToZero(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const (
		marketID   = int64(7)
		fallbackID = int64(42)
	)

	emptyRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	})

	// First query: market-scoped → no row → ErrRecordNotFound.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(marketID, true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	// Second query: fallback to market_id=0 → global default row.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(int64(0), true, true, 1).
		WillReturnRows(emptyRows.AddRow(
			fallbackID, int64(1), int64(0), "CNY", "Global-Default",
			true, true, "standard", int64(0),
			time.Now(), time.Now(), sql.NullTime{},
		))

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, marketID)
	if err != nil {
		t.Fatalf("FindDefaultByMarket: %v", err)
	}
	if got == nil {
		t.Fatal("expected fallback template, got nil")
	}
	if got.ID != fallbackID {
		t.Errorf("expected fallback id %d, got %d", fallbackID, got.ID)
	}
	if got.MarketID != 0 {
		t.Errorf("expected fallback market_id 0, got %d", got.MarketID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindDefaultByMarket_NotFound pins the "no default anywhere"
// contract: when neither market-scoped nor global default exists, the method
// must return code.ErrShippingTemplateNotFound (not raw gorm.ErrRecordNotFound
// or nil — callers rely on errors.Is(err, code.Err...)).
func TestShippingRepo_FindDefaultByMarket_NotFound(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const marketID = int64(7)

	// Both queries miss.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(marketID, true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(int64(0), true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, marketID)
	if !errors.Is(err, code.ErrShippingTemplateNotFound) {
		t.Fatalf("expected ErrShippingTemplateNotFound, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil template on not-found, got %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindListByMarket_Filtered pins the per-market filter:
// when marketID > 0 the WHERE clause must include market_id = ?.
func TestShippingRepo_FindListByMarket_Filtered(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const marketID = int64(7)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	// Count query also receives LIMIT/OFFSET args from GORM's Count+Find chain.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(marketID).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).
		AddRow(int64(1), int64(1), marketID, "HKD", "tmpl-a",
			false, true, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{}).
		AddRow(int64(2), int64(1), marketID, "HKD", "tmpl-b",
			true, true, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{})

	// Find: market_id, LIMIT (GORM emits LIMIT but no OFFSET for our offset=0).
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(marketID, 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, marketID, "", nil, 1, 10)
	if err != nil {
		t.Fatalf("FindListByMarket: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(got))
	}
	for _, tmpl := range got {
		if tmpl.MarketID != marketID {
			t.Errorf("expected market_id %d, got %d", marketID, tmpl.MarketID)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindListByMarket_AllMarkets pins the "marketID=0 means all"
// contract: when marketID = 0, the WHERE clause must NOT include a market_id
// filter — the listing must be unfiltered by market.
func TestShippingRepo_FindListByMarket_AllMarkets(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	// No WithArgs — count query must have NO market_id filter arg.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, 0, "name", nil, 1, 10)
	if err != nil {
		t.Fatalf("FindListByMarket: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 rows, got %d", len(got))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindListByMarket_NameAndActiveFilter pins the optional
// filter combination: name LIKE + isActive, both must be applied via WHERE
// args in addition to (or instead of) the market filter.
func TestShippingRepo_FindListByMarket_NameAndActiveFilter(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const marketID = int64(3)
	isActive := true

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	// Count args: market_id, name LIKE, is_active.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(marketID, "%express%", isActive).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(int64(11), int64(1), marketID, "USD", "Express",
		false, isActive, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{})

	// Find args: market_id, name LIKE, is_active, LIMIT.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(marketID, "%express%", isActive, 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, marketID, "express", &isActive, 1, 10)
	if err != nil {
		t.Fatalf("FindListByMarket: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(got) != 1 || got[0].Name != "Express" {
		t.Errorf("unexpected rows: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_UnsetAllDefaultByMarket pins the per-market scope of the
// new method: the WHERE clause must combine is_default=true AND market_id=?.
// This prevents accidentally unsetting defaults across other markets when
// promoting a market-specific template to default.
func TestShippingRepo_UnsetAllDefaultByMarket(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	const marketID = int64(7)

	// GORM auto-wraps single UPDATE in a transaction.
	mock.ExpectBegin()
	// UPDATE args (GORM ordering): is_default=false, updated_at=now,
	// is_default=true (WHERE), market_id=marketID (WHERE).
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `shipping_templates` SET `is_default`=?,`updated_at`=? WHERE is_default = ? AND market_id = ?")).
		WithArgs(false, sqlmock.AnyArg(), true, marketID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	if err := repo.UnsetAllDefaultByMarket(context.Background(), gdb, marketID); err != nil {
		t.Fatalf("UnsetAllDefaultByMarket: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_UnsetAllDefaultByMarket_AllMarkets pins that
// marketID=0 (全市场通用) means "unset ALL defaults across every market".
// This is the operator override used when there is genuinely no per-market
// default and the global one needs to be replaced.
func TestShippingRepo_UnsetAllDefaultByMarket_AllMarkets(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	mock.ExpectBegin()
	// marketID=0: WHERE clause has only is_default=true (no market_id).
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `shipping_templates` SET `is_default`=?,`updated_at`=? WHERE is_default = ?")).
		WithArgs(false, sqlmock.AnyArg(), true).
		WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	if err := repo.UnsetAllDefaultByMarket(context.Background(), gdb, 0); err != nil {
		t.Fatalf("UnsetAllDefaultByMarket(0): %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// silence unused-import warning for shipping pkg when the only reference is
// in the rows-set helper above (kept for type clarity).
var _ = shipping.ShippingTemplate{}
