package infra

import (
	"strings"

	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const tenantScopePluginName = "tenant_scope"

// TenantScopePlugin 自动为租户隔离表注入 WHERE tenant_id = ? 条件
// 同时为 Create 操作自动填充 tenant_id 列
type TenantScopePlugin struct{}

func (p *TenantScopePlugin) Name() string {
	return tenantScopePluginName
}

func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
	// Query/Delete/Update: 注入 WHERE tenant_id = ?
	if err := db.Callback().Query().Before("gorm:query").
		Register(tenantScopePluginName+":before_query", p.onQuery); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").
		Register(tenantScopePluginName+":before_delete", p.onQuery); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").
		Register(tenantScopePluginName+":before_update", p.onQuery); err != nil {
		return err
	}

	// Create: 自动填充 tenant_id
	if err := db.Callback().Create().Before("gorm:create").
		Register(tenantScopePluginName+":before_create", p.onCreate); err != nil {
		return err
	}

	// Raw: SQL 改写
	if err := db.Callback().Raw().Before("gorm:raw").
		Register(tenantScopePluginName+":before_raw", p.onRaw); err != nil {
		return err
	}

	return nil
}

// onQuery 为 Query/Delete/Update 注入 WHERE tenant_id = ?
func (p *TenantScopePlugin) onQuery(db *gorm.DB) {
	tenantID, ok := p.getTenantID(db)
	if !ok {
		return
	}

	tableName := p.getTableName(db)
	if !IsTenantScopedTable(tableName) {
		return
	}

	db.Statement.AddClause(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Name: "tenant_id"}, Value: tenantID},
		},
	})
}

// onCreate 为 Create 自动填充 tenant_id 列
func (p *TenantScopePlugin) onCreate(db *gorm.DB) {
	tenantID, ok := p.getTenantID(db)
	if !ok {
		return
	}

	tableName := p.getTableName(db)
	if !IsTenantScopedTable(tableName) {
		return
	}

	db.Statement.SetColumn("tenant_id", tenantID)
}

// onRaw 处理 db.Raw() — 通过 SQL AST 改写
func (p *TenantScopePlugin) onRaw(db *gorm.DB) {
	tenantID, ok := p.getTenantID(db)
	if !ok {
		return
	}

	sqlStr := db.Statement.SQL.String()
	if sqlStr == "" {
		return
	}

	rewritten, err := RewriteSQLWithTenantScope(sqlStr, tenantID)
	if err != nil {
		// 解析失败，记录日志，不改写（安全降级）
		logx.Errorf("tenant_scope: failed to rewrite SQL: %v, original: %s", err, sqlStr)
		return
	}

	db.Statement.SQL.Reset()
	db.Statement.SQL.WriteString(rewritten)
}

// getTenantID 从 context 获取租户 ID，返回 (tenantID, 是否需要过滤)
func (p *TenantScopePlugin) getTenantID(db *gorm.DB) (int64, bool) {
	ctx := db.Statement.Context
	if ctx == nil {
		return 0, false
	}

	// 检查是否显式跳过
	if IsSkipTenantScope(ctx) {
		return 0, false
	}

	tenantID, ok := contextx.GetTenantID(ctx)
	if !ok || tenantID == 0 {
		return 0, false
	}

	return tenantID, true
}

// getTableName 获取当前操作的表名
func (p *TenantScopePlugin) getTableName(db *gorm.DB) string {
	// 优先使用 Statement.Table
	if db.Statement.Table != "" {
		return db.Statement.Table
	}
	// 其次从 Schema 获取
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	// 从 SQL 中提取表名（兜底）
	sql := db.Statement.SQL.String()
	return extractTableNameFromSQL(sql)
}

// extractTableNameFromSQL 从 SQL 语句中提取第一个表名
// 用于 Statement.Table 和 Schema 都为空时的兜底
func extractTableNameFromSQL(sql string) string {
	upper := strings.ToUpper(strings.TrimSpace(sql))

	// SELECT ... FROM table
	if idx := strings.Index(upper, "FROM"); idx >= 0 {
		rest := strings.TrimSpace(sql[idx+4:])
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			name := strings.Trim(parts[0], "`\"'")
			// 跳过子查询开头的 (
			if name == "(" {
				return ""
			}
			// 去掉别名（如 "users u" -> "users"）
			return name
		}
	}

	// INSERT INTO table
	if idx := strings.Index(upper, "INTO"); idx >= 0 {
		rest := strings.TrimSpace(sql[idx+4:])
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			return strings.Trim(parts[0], "`\"'")
		}
	}

	// UPDATE table
	if strings.HasPrefix(upper, "UPDATE") {
		rest := strings.TrimSpace(sql[6:])
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			return strings.Trim(parts[0], "`\"'")
		}
	}

	// DELETE FROM table
	if strings.HasPrefix(upper, "DELETE") {
		if idx := strings.Index(upper, "FROM"); idx >= 0 {
			rest := strings.TrimSpace(sql[idx+4:])
			parts := strings.Fields(rest)
			if len(parts) > 0 {
				return strings.Trim(parts[0], "`\"'")
			}
		}
	}

	return ""
}

// Close 实现 gorm.Plugin 接口（空实现）
func (p *TenantScopePlugin) Close() error {
	return nil
}

// TranslateError 实现 gorm.Plugin 接口（透传错误）
func (p *TenantScopePlugin) TranslateError(err error) error {
	return err
}
