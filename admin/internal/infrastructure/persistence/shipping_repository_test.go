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
	// Args: tenantID, marketID, is_default=true, is_active=true, LIMIT 1 (GORM's First()).
	const tenantID = int64(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, true, true, 1).
		WillReturnRows(rows)

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, tenantID, marketID)
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
	const tenantID = int64(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	// Second query: fallback to market_id=0 → global default row.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, int64(0), true, true, 1).
		WillReturnRows(emptyRows.AddRow(
			fallbackID, int64(1), int64(0), "CNY", "Global-Default",
			true, true, "standard", int64(0),
			time.Now(), time.Now(), sql.NullTime{},
		))

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, tenantID, marketID)
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
	const tenantID = int64(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, int64(0), true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, tenantID, marketID)
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
	const tenantID = int64(1)
	// Count query: tenant_id, market_id.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(tenantID, marketID).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).
		AddRow(int64(1), tenantID, marketID, "HKD", "tmpl-a",
			false, true, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{}).
		AddRow(int64(2), tenantID, marketID, "HKD", "tmpl-b",
			true, true, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{})

	// Find: tenant_id, market_id, LIMIT (GORM emits LIMIT but no OFFSET for our offset=0).
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, tenantID, marketID, "", nil, 1, 10)
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
	const tenantID = int64(1)
	// tenant_id is always present; market_id omitted (marketID=0 = all); name LIKE included.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(tenantID, "%name%").
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, "%name%", 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, tenantID, 0, "name", nil, 1, 10)
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

	const (
		marketID = int64(3)
		tenantID = int64(1)
	)
	isActive := true

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	// Count args: tenant_id, market_id, name LIKE, is_active.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, "%express%", isActive).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(int64(11), tenantID, marketID, "USD", "Express",
		false, isActive, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{})

	// Find args: tenant_id, market_id, name LIKE, is_active, LIMIT.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, "%express%", isActive, 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, tenantID, marketID, "express", &isActive, 1, 10)
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
	const tenantID = int64(1)
	mock.ExpectBegin()
	// UPDATE args (GORM ordering): is_default=false, updated_at=now,
	// tenant_id (WHERE), is_default=true (WHERE), market_id=marketID (WHERE).
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `shipping_templates` SET `is_default`=?,`updated_at`=? WHERE tenant_id = ? AND is_default = ? AND market_id = ?")).
		WithArgs(false, sqlmock.AnyArg(), tenantID, true, marketID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	if err := repo.UnsetAllDefaultByMarket(context.Background(), gdb, tenantID, marketID); err != nil {
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
	// marketID=0: WHERE clause has tenant_id + is_default=true (no market_id).
	const tenantID = int64(1)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `shipping_templates` SET `is_default`=?,`updated_at`=? WHERE tenant_id = ? AND is_default = ?")).
		WithArgs(false, sqlmock.AnyArg(), tenantID, true).
		WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	if err := repo.UnsetAllDefaultByMarket(context.Background(), gdb, tenantID, 0); err != nil {
		t.Fatalf("UnsetAllDefaultByMarket(0): %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindDefaultByMarket_TenantIsolation pins the C3 fix:
// FindDefaultByMarket must scope the lookup to the supplied tenantID.
// Seed two templates in the same market but different tenants; each call
// must return only its own tenant's row, never the other tenant's.
//
// This is a SQL-shape assertion: the mocked query receives (tenantID, marketID,
// is_default, is_active, LIMIT) — proving tenant_id appears in WHERE.
func TestShippingRepo_FindDefaultByMarket_TenantIsolation(t *testing.T) {
	const (
		marketID    = int64(7)
		tenant1ID   = int64(100)
		tenant2ID   = int64(200)
		tenant1Row  = int64(101)
		tenant2Row  = int64(202)
	)
	row := func(id, tenant int64) *sqlmock.Rows {
		return sqlmock.NewRows([]string{
			"id", "tenant_id", "market_id", "currency", "name",
			"is_default", "is_active", "carrier_code", "warehouse_id",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			id, tenant, marketID, "CNY", "T",
			true, true, "standard", int64(0),
			time.Now(), time.Now(), sql.NullTime{},
		)
	}

	// Tenant 1 → tenant1Row
	gdb1, mock1 := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()
	mock1.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenant1ID, marketID, true, true, 1).
		WillReturnRows(row(tenant1Row, tenant1ID))
	got1, err := repo.FindDefaultByMarket(context.Background(), gdb1, tenant1ID, marketID)
	if err != nil {
		t.Fatalf("tenant1 FindDefaultByMarket: %v", err)
	}
	if got1 == nil || got1.ID != tenant1Row || got1.TenantID != tenant1ID {
		t.Fatalf("tenant1: expected id=%d tenant=%d, got %+v", tenant1Row, tenant1ID, got1)
	}
	if err := mock1.ExpectationsWereMet(); err != nil {
		t.Fatalf("tenant1 expectations: %v", err)
	}

	// Tenant 2 → tenant2Row (proves the call does NOT silently return tenant 1's row).
	gdb2, mock2 := newMockDB(t)
	mock2.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenant2ID, marketID, true, true, 1).
		WillReturnRows(row(tenant2Row, tenant2ID))
	got2, err := repo.FindDefaultByMarket(context.Background(), gdb2, tenant2ID, marketID)
	if err != nil {
		t.Fatalf("tenant2 FindDefaultByMarket: %v", err)
	}
	if got2 == nil || got2.ID != tenant2Row || got2.TenantID != tenant2ID {
		t.Fatalf("tenant2: expected id=%d tenant=%d, got %+v", tenant2Row, tenant2ID, got2)
	}
	if err := mock2.ExpectationsWereMet(); err != nil {
		t.Fatalf("tenant2 expectations: %v", err)
	}
}

// TestShippingRepo_FindDefaultByMarket_TenantFallbackIsolation pins the
// fallback path under tenant scope: when tenant A has no market=7 default
// but has a market=0 fallback, the call must return that fallback — and
// must NOT return tenant B's fallback (which would leak cross-tenant data).
func TestShippingRepo_FindDefaultByMarket_TenantFallbackIsolation(t *testing.T) {
	const (
		marketID  = int64(7)
		tenantAID = int64(300)
		fallbackA = int64(301)
	)
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	// Step 1: market-scoped lookup for tenant A → no row.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantAID, marketID, true, true, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	// Step 2: tenant-scoped fallback to market_id=0 → tenant A's global default.
	fallbackRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		fallbackA, tenantAID, int64(0), "CNY", "Global",
		true, true, "standard", int64(0),
		time.Now(), time.Now(), sql.NullTime{},
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantAID, int64(0), true, true, 1).
		WillReturnRows(fallbackRows)

	got, err := repo.FindDefaultByMarket(context.Background(), gdb, tenantAID, marketID)
	if err != nil {
		t.Fatalf("FindDefaultByMarket tenant fallback: %v", err)
	}
	if got == nil || got.ID != fallbackA || got.TenantID != tenantAID {
		t.Fatalf("expected tenant A fallback id=%d tenant=%d, got %+v", fallbackA, tenantAID, got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_FindListByMarket_TenantIsolation pins that the listing
// query is scoped to tenantID — cross-tenant rows must not leak into the
// result set. Asserts WHERE receives (tenantID, marketID, ...) args.
func TestShippingRepo_FindListByMarket_TenantIsolation(t *testing.T) {
	const (
		tenantID = int64(400)
		marketID = int64(7)
	)
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `shipping_templates`")).
		WithArgs(tenantID, marketID).
		WillReturnRows(countRows)

	listRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "market_id", "currency", "name",
		"is_default", "is_active", "carrier_code", "warehouse_id",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(int64(501), tenantID, marketID, "CNY", "Mine",
		true, true, "standard", int64(0), time.Now(), time.Now(), sql.NullTime{})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `shipping_templates`")).
		WithArgs(tenantID, marketID, 10).
		WillReturnRows(listRows)

	got, total, err := repo.FindListByMarket(context.Background(), gdb, tenantID, marketID, "", nil, 1, 10)
	if err != nil {
		t.Fatalf("FindListByMarket: %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Fatalf("expected 1 row, got total=%d len=%d", total, len(got))
	}
	if got[0].TenantID != tenantID {
		t.Fatalf("expected tenant=%d, got %d", tenantID, got[0].TenantID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// TestShippingRepo_UnsetAllDefaultByMarket_TenantIsolation pins the critical
// safety property: UnsetAllDefaultByMarket must NOT touch another tenant's
// default when promoting tenant A's template. Asserts WHERE receives
// (false, ts, tenantID, true, marketID) — GORM keeps the .Where() order.
func TestShippingRepo_UnsetAllDefaultByMarket_TenantIsolation(t *testing.T) {
	const (
		tenantID = int64(500)
		marketID = int64(7)
	)
	gdb, mock := newMockDB(t)
	repo := persistence.NewShippingTemplateRepository()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `shipping_templates` SET `is_default`=?,`updated_at`=? WHERE tenant_id = ? AND is_default = ? AND market_id = ?")).
		WithArgs(false, sqlmock.AnyArg(), tenantID, true, marketID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.UnsetAllDefaultByMarket(context.Background(), gdb, tenantID, marketID); err != nil {
		t.Fatalf("UnsetAllDefaultByMarket: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

// silence unused-import warning for shipping pkg when the only reference is
// in the rows-set helper above (kept for type clarity).
var _ = shipping.ShippingTemplate{}
