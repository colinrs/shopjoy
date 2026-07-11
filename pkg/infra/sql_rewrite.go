package infra

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/format"
	"github.com/pingcap/tidb/pkg/parser/opcode"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
)

// RewriteSQLWithTenantScope 使用 SQL AST 解析为注册表中的表注入 tenant_id 条件
// 覆盖 SELECT、UPDATE、DELETE 语句中的 FROM/JOIN/子查询
func RewriteSQLWithTenantScope(sql string, tenantID int64) (string, error) {
	p := parser.New()
	stmts, _, err := p.ParseSQL(sql)
	if err != nil {
		return "", fmt.Errorf("parse SQL: %w", err)
	}

	for _, stmt := range stmts {
		rewriter := &tenantScopeRewriter{tenantID: tenantID}
		stmt.Accept(rewriter)
	}

	var buf strings.Builder
	restoreFlags := format.RestoreStringSingleQuotes | format.RestoreNameBackQuotes
	for i, stmt := range stmts {
		if i > 0 {
			buf.WriteString(";")
		}
		restoreCtx := format.NewRestoreCtx(restoreFlags, &buf)
		if err := stmt.Restore(restoreCtx); err != nil {
			return "", fmt.Errorf("restore SQL: %w", err)
		}
	}
	return buf.String(), nil
}

// tenantScopeRewriter 实现 ast.Visitor 接口，遍历 SQL AST 注入 tenant_id 条件
type tenantScopeRewriter struct {
	tenantID int64
}

// Enter 实现 ast.Visitor 接口 — 在进入节点时注入条件
func (r *tenantScopeRewriter) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		r.rewriteSelect(node)
	case *ast.DeleteStmt:
		r.rewriteDelete(node)
	case *ast.UpdateStmt:
		r.rewriteUpdate(node)
	case *ast.InsertStmt:
		r.rewriteInsert(node)
	}
	return in, false // false = 继续遍历子节点
}

// Leave 实现 ast.Visitor 接口
func (r *tenantScopeRewriter) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// rewriteSelect 为 SELECT 语句注入 WHERE tenant_id = ?
func (r *tenantScopeRewriter) rewriteSelect(sel *ast.SelectStmt) {
	if sel.From == nil || sel.From.TableRefs == nil {
		return
	}

	tableAliases := r.collectTenantScopedTables(sel.From.TableRefs)
	if len(tableAliases) == 0 {
		return
	}

	conditions := r.buildConditions(tableAliases)
	if len(conditions) == 0 {
		return
	}

	r.injectWhere(sel.Where, conditions, func(newWhere ast.ExprNode) {
		sel.Where = newWhere
	})
}

// rewriteDelete 为 DELETE 语句注入 WHERE tenant_id = ?
func (r *tenantScopeRewriter) rewriteDelete(del *ast.DeleteStmt) {
	if del.TableRefs == nil || del.TableRefs.TableRefs == nil {
		return
	}

	tableAliases := r.collectTenantScopedTables(del.TableRefs.TableRefs)
	if len(tableAliases) == 0 {
		return
	}

	conditions := r.buildConditions(tableAliases)
	if len(conditions) == 0 {
		return
	}

	r.injectWhere(del.Where, conditions, func(newWhere ast.ExprNode) {
		del.Where = newWhere
	})
}

// rewriteUpdate 为 UPDATE 语句注入 WHERE tenant_id = ?
func (r *tenantScopeRewriter) rewriteUpdate(upd *ast.UpdateStmt) {
	if upd.TableRefs == nil || upd.TableRefs.TableRefs == nil {
		return
	}

	tableAliases := r.collectTenantScopedTables(upd.TableRefs.TableRefs)
	if len(tableAliases) == 0 {
		return
	}

	conditions := r.buildConditions(tableAliases)
	if len(conditions) == 0 {
		return
	}

	r.injectWhere(upd.Where, conditions, func(newWhere ast.ExprNode) {
		upd.Where = newWhere
	})
}

// rewriteInsert 为 INSERT 语句自动添加 tenant_id 列和值
func (r *tenantScopeRewriter) rewriteInsert(ins *ast.InsertStmt) {
	if ins.Table == nil || ins.Table.TableRefs == nil {
		return
	}

	tableName := r.extractTableNameFromRefs(ins.Table)
	if !IsTenantScopedTable(tableName) {
		return
	}

	// 检查是否已经有 tenant_id 列
	if ins.Columns != nil {
		for _, col := range ins.Columns {
			if strings.EqualFold(col.Name.O, "tenant_id") {
				return // 已有，不重复添加
			}
		}
	}

	// 添加 tenant_id 列
	if ins.Columns == nil {
		ins.Columns = []*ast.ColumnName{}
	}
	ins.Columns = append(ins.Columns, &ast.ColumnName{
		Name: ast.NewCIStr("tenant_id"),
	})

	// 添加对应的值
	if ins.Lists != nil && len(ins.Lists) > 0 {
		tenantValue := ast.NewValueExpr(r.tenantID, "", "")
		for i := range ins.Lists {
			ins.Lists[i] = append(ins.Lists[i], tenantValue)
		}
	}
}

// collectTenantScopedTables 从 JOIN 树中收集所有租户隔离表及其别名
func (r *tenantScopeRewriter) collectTenantScopedTables(tableRefs *ast.Join) map[string]string {
	result := make(map[string]string)
	r.walkTableRefs(tableRefs, result)
	return result
}

// walkTableRefs 递归遍历 JOIN 树
func (r *tenantScopeRewriter) walkTableRefs(join *ast.Join, result map[string]string) {
	if join == nil {
		return
	}
	if join.Left != nil {
		r.walkTableSource(join.Left, result)
	}
	if join.Right != nil {
		r.walkTableSource(join.Right, result)
	}
}

// walkTableSource 处理单个表源（表名或子查询）
func (r *tenantScopeRewriter) walkTableSource(node ast.ResultSetNode, result map[string]string) {
	if node == nil {
		return
	}

	switch src := node.(type) {
	case *ast.TableSource:
		if src.Source == nil {
			return
		}
		// 子查询不在此层处理，但子查询内部的表会被 Enter 递归访问
		if _, ok := src.Source.(*ast.SelectStmt); ok {
			return
		}
		tableName := r.extractTableNameFromSource(src)
		if tableName == "" || !IsTenantScopedTable(tableName) {
			return
		}
		alias := src.AsName.O
		if alias == "" {
			alias = tableName
		}
		result[alias] = tableName
	case *ast.Join:
		r.walkTableRefs(src, result)
	}
}

// extractTableNameFromRefs 从 TableRefsClause 提取表名
func (r *tenantScopeRewriter) extractTableNameFromRefs(refs *ast.TableRefsClause) string {
	if refs == nil || refs.TableRefs == nil {
		return ""
	}
	return r.extractTableNameFromJoin(refs.TableRefs)
}

// extractTableNameFromJoin 从 JOIN 树中提取第一个表名
func (r *tenantScopeRewriter) extractTableNameFromJoin(join *ast.Join) string {
	if join == nil {
		return ""
	}
	if join.Left != nil {
		if ts, ok := join.Left.(*ast.TableSource); ok {
			return r.extractTableNameFromSource(ts)
		}
		if j, ok := join.Left.(*ast.Join); ok {
			return r.extractTableNameFromJoin(j)
		}
	}
	if join.Right != nil {
		if ts, ok := join.Right.(*ast.TableSource); ok {
			return r.extractTableNameFromSource(ts)
		}
	}
	return ""
}

// extractTableNameFromSource 从 TableSource 提取表名
func (r *tenantScopeRewriter) extractTableNameFromSource(src *ast.TableSource) string {
	if src == nil || src.Source == nil {
		return ""
	}
	if tableName, ok := src.Source.(*ast.TableName); ok {
		return tableName.Name.O
	}
	return ""
}

// buildConditions 为一组表别名构建 tenant_id 条件
func (r *tenantScopeRewriter) buildConditions(tableAliases map[string]string) []ast.ExprNode {
	var conditions []ast.ExprNode

	for alias := range tableAliases {
		colName := &ast.ColumnNameExpr{
			Name: &ast.ColumnName{
				Table: ast.NewCIStr(alias),
				Name:  ast.NewCIStr("tenant_id"),
			},
		}
		val := ast.NewValueExpr(r.tenantID, "", "")
		condition := &ast.BinaryOperationExpr{
			Op: opcode.EQ,
			L:  colName,
			R:  val,
		}
		conditions = append(conditions, condition)
	}

	return conditions
}

// injectWhere 将条件注入到现有 WHERE 中（用 AND 连接）
func (r *tenantScopeRewriter) injectWhere(existing ast.ExprNode, conditions []ast.ExprNode, setter func(ast.ExprNode)) {
	if len(conditions) == 0 {
		return
	}

	var newConds ast.ExprNode = conditions[0]
	for i := 1; i < len(conditions); i++ {
		newConds = &ast.BinaryOperationExpr{
			Op: opcode.LogicAnd,
			L:  newConds,
			R:  conditions[i],
		}
	}

	if existing != nil {
		newConds = &ast.BinaryOperationExpr{
			Op: opcode.LogicAnd,
			L:  existing,
			R:  newConds,
		}
	}

	setter(newConds)
}
