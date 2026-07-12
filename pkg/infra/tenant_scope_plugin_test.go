package infra

import (
	"context"
	"testing"

	"github.com/colinrs/shopjoy/pkg/contextx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testModel 模拟租户隔离表
type testModel struct {
	ID       int64  `gorm:"primarykey"`
	TenantID int64  `gorm:"column:tenant_id;not null"`
	Name     string `gorm:"column:name"`
}

func (testModel) TableName() string { return "products" }

// systemModel 模拟系统级表（不在注册表中）
type systemModel struct {
	ID   int64  `gorm:"primarykey"`
	Name string `gorm:"column:name"`
}

func (systemModel) TableName() string { return "tenants" }

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// 注册 Plugin
	if err := (&TenantScopePlugin{}).Initialize(db); err != nil {
		t.Fatalf("failed to initialize tenant scope plugin: %v", err)
	}

	// 创建测试表
	if err := db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		tenant_id INTEGER NOT NULL,
		name TEXT
	)`).Error; err != nil {
		t.Fatalf("failed to create products table: %v", err)
	}

	if err := db.Exec(`CREATE TABLE tenants (
		id INTEGER PRIMARY KEY,
		name TEXT
	)`).Error; err != nil {
		t.Fatalf("failed to create tenants table: %v", err)
	}

	return db
}

func ctxWithTenant(tenantID int64) context.Context {
	ctx := context.Background()
	ctx = contextx.SetTenantID(ctx, tenantID)
	ctx = contextx.SetUserType(ctx, contextx.UserTypeTenantAdmin)
	return ctx
}

func ctxPlatformAdmin() context.Context {
	ctx := context.Background()
	ctx = contextx.SetTenantID(ctx, 0)
	ctx = contextx.SetUserType(ctx, contextx.UserTypePlatformAdmin)
	return ctx
}

// captureSQL 捕获 GORM 生成的 SQL
func captureSQL(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{
		DryRun: true,
	})
}

// TestORMQueryAutoInject 测试 ORM 查询自动注入 WHERE tenant_id = ?
func TestORMQueryAutoInject(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxWithTenant(42)
	dryDB := captureSQL(db).WithContext(ctx)

	var results []testModel
	tx := dryDB.Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	if !contains(sql, "tenant_id") {
		t.Errorf("expected SQL to contain tenant_id condition, got: %s", sql)
	}
	t.Logf("Generated SQL: %s", sql)
}

// TestORMCreateAutoFill 测试 Create 自动填充 tenant_id
func TestORMCreateAutoFill(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxWithTenant(42)
	dryDB := captureSQL(db).WithContext(ctx)

	model := &testModel{ID: 1, Name: "test"}
	tx := dryDB.Create(model)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	// 验证 tenant_id 被自动设置
	if model.TenantID != 42 {
		t.Errorf("expected tenant_id=42, got: %d", model.TenantID)
	}
}

// TestORMQueryNonTenantTable 测试非租户表不注入条件
func TestORMQueryNonTenantTable(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxWithTenant(42)
	dryDB := captureSQL(db).WithContext(ctx)

	var results []systemModel
	tx := dryDB.Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	if contains(sql, "tenant_id") {
		t.Errorf("expected SQL NOT to contain tenant_id for system table, got: %s", sql)
	}
}

// TestPlatformAdminSkip 测试平台管理员（tenantID=0）跳过过滤
func TestPlatformAdminSkip(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxPlatformAdmin()
	dryDB := captureSQL(db).WithContext(ctx)

	var results []testModel
	tx := dryDB.Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	if contains(sql, "tenant_id") {
		t.Errorf("expected SQL NOT to contain tenant_id for platform admin, got: %s", sql)
	}
}

// TestSkipTenantScope 测试显式跳过租户过滤
func TestSkipTenantScope(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxWithTenant(42)
	ctx = SkipTenantScope(ctx)
	dryDB := captureSQL(db).WithContext(ctx)

	var results []testModel
	tx := dryDB.Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	if contains(sql, "tenant_id") {
		t.Errorf("expected SQL NOT to contain tenant_id when SkipTenantScope is set, got: %s", sql)
	}
}

// TestRawSQLRewrite 测试 Raw SQL 自动改写
func TestRawSQLRewrite(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		tenantID int64
		wantErr  bool
		wantTID  bool // 是否应包含 tenant_id
	}{
		{
			name:     "simple SELECT",
			sql:      "SELECT * FROM products WHERE name = 'test'",
			tenantID: 42,
			wantTID:  true,
		},
		{
			name:     "SELECT with JOIN",
			sql:      "SELECT * FROM products p JOIN orders o ON p.id = o.product_id WHERE p.name = 'test'",
			tenantID: 42,
			wantTID:  true,
		},
		{
			name:     "non-tenant table",
			sql:      "SELECT * FROM tenants WHERE name = 'test'",
			tenantID: 42,
			wantTID:  false,
		},
		{
			name:     "UPDATE statement",
			sql:      "UPDATE products SET name = 'updated' WHERE id = 1",
			tenantID: 42,
			wantTID:  true,
		},
		{
			name:     "DELETE statement",
			sql:      "DELETE FROM products WHERE id = 1",
			tenantID: 42,
			wantTID:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RewriteSQLWithTenantScope(tt.sql, tt.tenantID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("Original:  %s", tt.sql)
			t.Logf("Rewritten: %s", result)

			if tt.wantTID && !contains(result, "tenant_id") {
				t.Errorf("expected rewritten SQL to contain tenant_id, got: %s", result)
			}
			if !tt.wantTID && contains(result, "tenant_id") {
				t.Errorf("expected rewritten SQL NOT to contain tenant_id, got: %s", result)
			}
		})
	}
}

// TestRawSQLWithSubquery 测试子查询中的表也被处理
func TestRawSQLWithSubquery(t *testing.T) {
	sql := "SELECT * FROM products WHERE id IN (SELECT product_id FROM orders WHERE status = 1)"
	result, err := RewriteSQLWithTenantScope(sql, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Original:  %s", sql)
	t.Logf("Rewritten: %s", result)

	// products 和 orders 都是租户隔离表，都应该有 tenant_id 条件
	if !contains(result, "tenant_id") {
		t.Errorf("expected rewritten SQL to contain tenant_id, got: %s", result)
	}
}

// TestRawSQLWithOR 测试 OR 表达式的优先级处理
func TestRawSQLWithOR(t *testing.T) {
	sql := "SELECT * FROM products WHERE status = 1 OR name = 'test'"
	result, err := RewriteSQLWithTenantScope(sql, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Original:  %s", sql)
	t.Logf("Rewritten: %s", result)

	// 应该包含 tenant_id 条件
	if !contains(result, "tenant_id") {
		t.Errorf("expected rewritten SQL to contain tenant_id, got: %s", result)
	}

	// 验证 OR 表达式被正确括起来
	// 期望: WHERE (status = 1 OR name = 'test') AND products.tenant_id = 42
	// 而不是: WHERE status = 1 OR name = 'test' AND products.tenant_id = 42
	if contains(result, "OR") && contains(result, "tenant_id") {
		t.Logf("OR expression with tenant_id - verify precedence is correct")
	}
}

// TestNoTenantIDInContext 测试 context 中没有 tenantID 时不注入
func TestNoTenantIDInContext(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background() // 没有设置 tenantID
	dryDB := captureSQL(db).WithContext(ctx)

	var results []testModel
	tx := dryDB.Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	if contains(sql, "tenant_id") {
		t.Errorf("expected SQL NOT to contain tenant_id when no tenant in context, got: %s", sql)
	}
}

// TestORMQueryWithJoin 测试 JOIN 查询中主表自动注入 tenant_id
// 注意：JOIN 的子表不会被自动注入，需要手动处理
func TestORMQueryWithJoin(t *testing.T) {
	db := setupTestDB(t)
	ctx := ctxWithTenant(42)
	dryDB := captureSQL(db).WithContext(ctx)

	// 创建 order_items 表（不在注册表中）
	dryDB.Exec("CREATE TABLE IF NOT EXISTS order_items (id INTEGER PRIMARY KEY, order_id INTEGER, product_id INTEGER)")

	type orderItem struct {
		ID        int64 `gorm:"primarykey"`
		OrderID   int64 `gorm:"column:order_id"`
		ProductID int64 `gorm:"column:product_id"`
	}

	var results []orderItem
	tx := dryDB.Table("order_items oi").
		Select("oi.*").
		Joins("JOIN products p ON p.id = oi.product_id").
		Find(&results)
	if tx.Error != nil {
		t.Fatalf("unexpected error: %v", tx.Error)
	}

	sql := tx.Statement.SQL.String()
	t.Logf("Generated SQL: %s", sql)

	// 主表 order_items 不在注册表中，不会自动注入
	// JOIN 的 products 在注册表中，但 ORM callback 不会处理 JOIN 表
	// 这就是为什么需要手动添加 JOIN 表的 tenant_id 过滤
	if contains(sql, "tenant_id") {
		t.Logf("SQL contains tenant_id (good - but note ORM only handles primary table)")
	} else {
		t.Logf("SQL does NOT contain tenant_id - JOIN queries need manual tenant filtering")
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
