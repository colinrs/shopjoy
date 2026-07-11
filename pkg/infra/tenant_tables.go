package infra

// tenantScopedTables 记录所有需要租户隔离的表。
// 新增表时必须在此注册，配合 linter 检查遗漏。
// 对照 sql/ 目录下各领域 schema.sql 中有 tenant_id 列的表。
var tenantScopedTables = map[string]bool{
	// User Domain
	"users":               true,
	"admin_users":         true,
	"roles":               true,
	"user_addresses":      true,
	"user_operation_logs": true,

	// Product Domain
	"categories":              true,
	"brands":                  true,
	"products":                true,
	"skus":                    true,
	"product_markets":         true,
	"product_localizations":   true,
	"category_markets":        true,
	"brand_markets":           true,
	"warehouses":              true,
	"warehouse_inventories":   true,
	"inventory_logs":          true,
	"markets":                 true,

	// Order Domain
	"orders":     true,
	"carts":      true,
	"cart_items": true,

	// Payment Domain
	"payments":            true,
	"order_payments":      true,
	"payment_transactions": true,
	"payment_refunds":     true,
	"webhook_events":      true,

	// Fulfillment Domain
	"shipments":                  true,
	"shipment_items":             true,
	"refunds":                    true,
	"shipping_templates":         true,
	"shipping_zones":             true,
	"shipping_template_mappings": true,

	// Promotion Domain
	"promotions":       true,
	"promotion_rules":  true,
	"promotion_usage":  true,
	"coupons":          true,
	"user_coupons":     true,

	// Points Domain
	"earn_rules":          true,
	"redeem_rules":        true,
	"points_accounts":     true,
	"points_transactions": true,
	"points_redemptions":  true,

	// Review Domain
	"reviews":       true,
	"review_replies": true,
	"review_stats":  true,

	// Storefront Domain
	"shops":            true,
	"themes":           true,
	"pages":            true,
	"navigations":      true,
	"decorations":      true,
	"page_versions":    true,
	"seo_configs":      true,
	"theme_audit_logs": true,

	// Shop Domain
	"shop_settings": true,
}

// IsTenantScopedTable 检查表是否需要租户隔离
func IsTenantScopedTable(tableName string) bool {
	return tenantScopedTables[tableName]
}
