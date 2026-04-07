package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// APIEndpoint documents an API endpoint for testing coverage tracking
type APIEndpoint struct {
	Method       string
	Path         string
	AuthRequired bool
	Group        string
	Description  string
}

// AllAdminAPIEndpoints lists all admin API endpoints
var AllAdminAPIEndpoints = []APIEndpoint{
	// Auth (no auth required)
	{"POST", "/api/v1/auth/login", false, "auth", "管理员登录"},
	{"POST", "/api/v1/auth/register", false, "auth", "注册租户管理员"},

	// Users
	{"GET", "/api/v1/users/:id", true, "users", "获取用户详情"},
	{"GET", "/api/v1/users", true, "users", "获取用户列表"},
	{"PUT", "/api/v1/users/:id", true, "users", "更新用户信息"},
	{"POST", "/api/v1/users/:id/suspend", true, "users", "禁用用户"},
	{"POST", "/api/v1/users/:id/activate", true, "users", "启用用户"},
	{"DELETE", "/api/v1/users/:id", true, "users", "删除用户"},
	{"POST", "/api/v1/users/:id/reset-password", true, "users", "重置用户密码"},
	{"GET", "/api/v1/users/stats", true, "users", "获取用户统计"},
	{"GET", "/api/v1/users/:id/detail", true, "users", "获取用户详情（含聚合数据）"},
	{"GET", "/api/v1/users/:id/addresses", true, "users", "获取用户地址列表"},
	{"POST", "/api/v1/users/:id/suspend-reason", true, "users", "冻结用户（带原因）"},
	{"GET", "/api/v1/users/stats/enhanced", true, "users", "获取用户统计（增强版）"},
	{"GET", "/api/v1/users/export", true, "users", "导出用户列表"},
	{"GET", "/api/v1/users/enhanced", true, "users", "获取增强版用户列表"},
	{"POST", "/api/v1/users/batch-status", true, "users", "批量更新用户状态"},

	// Admin Users
	{"GET", "/api/v1/admin-users", true, "admin_users", "获取管理员列表"},
	{"POST", "/api/v1/admin-users", true, "admin_users", "创建管理员"},
	{"GET", "/api/v1/admin-users/:id", true, "admin_users", "获取管理员详情"},
	{"PUT", "/api/v1/admin-users/:id", true, "admin_users", "更新管理员"},
	{"DELETE", "/api/v1/admin-users/:id", true, "admin_users", "删除管理员"},
	{"POST", "/api/v1/admin-users/:id/reset-password", true, "admin_users", "重置管理员密码"},
	{"PUT", "/api/v1/admin-users/:id/roles", true, "admin_users", "分配角色"},
	{"POST", "/api/v1/admin-users/:id/disable", true, "admin_users", "禁用管理员"},
	{"POST", "/api/v1/admin-users/:id/enable", true, "admin_users", "启用管理员"},
	{"PUT", "/api/v1/admin-users/profile", true, "admin_users", "更新管理员资料"},
	{"PUT", "/api/v1/admin-users/password", true, "admin_users", "修改管理员密码"},

	// Products
	{"POST", "/api/v1/products", true, "products", "创建商品"},
	{"PUT", "/api/v1/products/:id", true, "products", "更新商品"},
	{"GET", "/api/v1/products/:id", true, "products", "获取商品详情"},
	{"GET", "/api/v1/products", true, "products", "获取商品列表"},
	{"GET", "/api/v1/products/export", true, "products", "导出商品"},
	{"POST", "/api/v1/products/:id/on-sale", true, "products", "上架商品"},
	{"POST", "/api/v1/products/:id/off-sale", true, "products", "下架商品"},
	{"PUT", "/api/v1/products/:id/stock", true, "products", "更新库存"},
	{"DELETE", "/api/v1/products/:id", true, "products", "删除商品"},
	{"POST", "/api/v1/products/batch-update", true, "products", "批量更新商品"},

	// SKUs
	{"POST", "/api/v1/skus", true, "products", "创建SKU变体"},
	{"PUT", "/api/v1/skus/:id", true, "products", "更新SKU变体"},
	{"GET", "/api/v1/skus/:id", true, "products", "获取SKU详情"},
	{"GET", "/api/v1/products/:product_id/skus", true, "products", "获取商品SKU列表"},
	{"DELETE", "/api/v1/skus/:id", true, "products", "删除SKU变体"},

	// Product Localizations
	{"POST", "/api/v1/product-localizations", true, "products", "创建商品本地化"},
	{"PUT", "/api/v1/product-localizations/:id", true, "products", "更新商品本地化"},
	{"GET", "/api/v1/product-localizations/:id", true, "products", "获取商品本地化详情"},
	{"GET", "/api/v1/products/:product_id/localizations", true, "products", "获取商品本地化列表"},
	{"DELETE", "/api/v1/product-localizations/:id", true, "products", "删除商品本地化"},

	// Categories
	{"POST", "/api/v1/categories", true, "categories", "创建分类"},
	{"PUT", "/api/v1/categories/:id", true, "categories", "更新分类"},
	{"GET", "/api/v1/categories/:id", true, "categories", "获取分类详情"},
	{"GET", "/api/v1/categories", true, "categories", "获取分类列表"},
	{"GET", "/api/v1/categories/tree", true, "categories", "获取分类树"},
	{"PUT", "/api/v1/categories/:id/status", true, "categories", "更新分类状态"},
	{"DELETE", "/api/v1/categories/:id", true, "categories", "删除分类"},
	{"PUT", "/api/v1/categories/sort", true, "categories", "更新分类排序"},
	{"PUT", "/api/v1/categories/:id/move", true, "categories", "移动分类"},
	{"GET", "/api/v1/categories/:id/product-count", true, "categories", "获取分类下商品数量"},
	{"PUT", "/api/v1/categories/:id/market-visibility", true, "categories", "设置分类市场可见性"},
	{"GET", "/api/v1/categories/:id/market-visibility", true, "categories", "获取分类市场可见性"},

	// Brands
	{"POST", "/api/v1/brands", true, "brands", "创建品牌"},
	{"PUT", "/api/v1/brands/:id", true, "brands", "更新品牌"},
	{"GET", "/api/v1/brands/:id", true, "brands", "获取品牌详情"},
	{"GET", "/api/v1/brands", true, "brands", "获取品牌列表"},
	{"PUT", "/api/v1/brands/:id/status", true, "brands", "更新品牌状态"},
	{"DELETE", "/api/v1/brands/:id", true, "brands", "删除品牌"},
	{"PUT", "/api/v1/brands/:id/toggle-page", true, "brands", "切换品牌专区"},
	{"GET", "/api/v1/brands/:id/product-count", true, "brands", "获取品牌下商品数量"},
	{"PUT", "/api/v1/brands/:id/market-visibility", true, "brands", "设置品牌市场可见性"},
	{"GET", "/api/v1/brands/:id/market-visibility", true, "brands", "获取品牌市场可见性"},

	// Roles
	{"GET", "/api/v1/roles", true, "roles", "获取角色列表"},
	{"POST", "/api/v1/roles", true, "roles", "创建角色"},
	{"GET", "/api/v1/roles/:id", true, "roles", "获取角色详情"},
	{"PUT", "/api/v1/roles/:id", true, "roles", "更新角色"},
	{"DELETE", "/api/v1/roles/:id", true, "roles", "删除角色"},
	{"PUT", "/api/v1/roles/:id/status", true, "roles", "更新角色状态"},
	{"PUT", "/api/v1/roles/:id/permissions", true, "roles", "更新角色权限"},
	{"GET", "/api/v1/permissions", true, "roles", "获取所有权限列表"},

	// Dashboard
	{"GET", "/api/v1/dashboard/overview", true, "dashboard", "获取仪表盘概览数据"},
	{"GET", "/api/v1/dashboard/sales-trend", true, "dashboard", "获取销售趋势"},
	{"GET", "/api/v1/dashboard/order-status", true, "dashboard", "获取订单状态分布"},
	{"GET", "/api/v1/dashboard/top-products", true, "dashboard", "获取热销商品TOP"},
	{"GET", "/api/v1/dashboard/pending-orders", true, "dashboard", "获取待处理订单"},
	{"GET", "/api/v1/dashboard/activities", true, "dashboard", "获取最近活动"},
	{"GET", "/api/v1/dashboard", true, "dashboard", "获取仪表盘所有数据"},

	// Fulfillment - Shipments
	{"POST", "/api/v1/shipments", true, "shipments", "创建发货单"},
	{"POST", "/api/v1/shipments/batch", true, "shipments", "批量创建发货单"},
	{"POST", "/api/v1/shipments/batch-tracking", true, "shipments", "批量更新追踪号"},
	{"GET", "/api/v1/shipments", true, "shipments", "发货单列表"},
	{"GET", "/api/v1/shipments/:id", true, "shipments", "发货单详情"},
	{"PUT", "/api/v1/shipments/:id", true, "shipments", "更新发货单"},
	{"PUT", "/api/v1/shipments/:id/status", true, "shipments", "更新发货单状态"},
	{"GET", "/api/v1/orders/:id/shipments", true, "shipments", "获取订单的发货单列表"},
	{"GET", "/api/v1/carriers", true, "shipments", "物流公司列表"},
	{"GET", "/api/v1/shipments/export", true, "shipments", "导出发货单"},
	{"PUT", "/api/v1/shipments/:id/cancel", true, "shipments", "取消发货单"},

	// Fulfillment - Refunds
	{"GET", "/api/v1/refunds", true, "refunds", "退款列表"},
	{"GET", "/api/v1/refunds/:id", true, "refunds", "退款详情"},
	{"PUT", "/api/v1/refunds/:id/approve", true, "refunds", "批准退款"},
	{"PUT", "/api/v1/refunds/:id/reject", true, "refunds", "拒绝退款"},
	{"GET", "/api/v1/refund-reasons", true, "refunds", "退款原因列表"},
	{"GET", "/api/v1/refunds/statistics", true, "refunds", "退款统计"},
	{"GET", "/api/v1/refunds/export", true, "refunds", "导出退款"},

	// Fulfillment - Orders
	{"GET", "/api/v1/orders", true, "fulfillment_orders", "订单列表（含履约状态筛选）"},
	{"GET", "/api/v1/orders/:id", true, "fulfillment_orders", "订单详情（含履约信息）"},
	{"PUT", "/api/v1/orders/:id/ship", true, "fulfillment_orders", "订单发货（创建发货单）"},
	{"GET", "/api/v1/orders/fulfillment-summary", true, "fulfillment_orders", "履约摘要统计"},
	{"PUT", "/api/v1/orders/:id/remark", true, "fulfillment_orders", "更新订单备注"},
	{"PUT", "/api/v1/orders/:id/adjust-price", true, "fulfillment_orders", "订单改价"},
	{"GET", "/api/v1/orders/export", true, "fulfillment_orders", "导出订单"},
	{"PUT", "/api/v1/orders/:id/cancel", true, "fulfillment_orders", "取消订单"},
	{"POST", "/api/v1/orders/batch-cancel", true, "fulfillment_orders", "批量取消订单"},
	{"POST", "/api/v1/orders/:id/remind-payment", true, "fulfillment_orders", "发送支付提醒"},

	// Fulfillment Statistics
	{"GET", "/api/v1/fulfillment/statistics", true, "fulfillment_statistics", "综合履约统计"},
	{"GET", "/api/v1/fulfillment/statistics/export", true, "fulfillment_statistics", "导出履约统计"},

	// Payments
	{"GET", "/api/v1/payments/transactions", true, "payments", "交易列表"},
	{"GET", "/api/v1/payments/transactions/:id", true, "payments", "交易详情"},
	{"GET", "/api/v1/payments/order/:id", true, "payments", "订单支付信息"},
	{"POST", "/api/v1/payments/:id/refund", true, "payments", "发起退款"},
	{"GET", "/api/v1/payments/stats", true, "payments", "支付统计"},

	// Regions
	{"GET", "/api/v1/regions", true, "regions", "地区列表"},

	// Warehouses
	{"POST", "/api/v1/warehouses", true, "warehouses", "创建仓库"},
	{"GET", "/api/v1/warehouses/:id", true, "warehouses", "获取仓库详情"},
	{"PUT", "/api/v1/warehouses/:id", true, "warehouses", "更新仓库"},
	{"DELETE", "/api/v1/warehouses/:id", true, "warehouses", "删除仓库"},
	{"PUT", "/api/v1/warehouses/:id/status", true, "warehouses", "更新仓库状态"},
	{"PUT", "/api/v1/warehouses/:id/default", true, "warehouses", "设置默认仓库"},
	{"GET", "/api/v1/warehouses", true, "warehouses", "仓库列表"},

	// Shipping Templates
	{"POST", "/api/v1/shipping-templates", true, "shipping_templates", "创建运费模板"},
	{"GET", "/api/v1/shipping-templates/:id", true, "shipping_templates", "获取运费模板详情"},
	{"PUT", "/api/v1/shipping-templates/:id", true, "shipping_templates", "更新运费模板"},
	{"DELETE", "/api/v1/shipping-templates/:id", true, "shipping_templates", "删除运费模板"},
	{"GET", "/api/v1/shipping-templates", true, "shipping_templates", "运费模板列表"},
	{"PUT", "/api/v1/shipping-templates/:id/default", true, "shipping_templates", "设置默认模板"},

	// Shipping Mappings
	{"POST", "/api/v1/shipping-mappings", true, "shipping_mappings", "创建配送映射"},
	{"GET", "/api/v1/shipping-mappings", true, "shipping_mappings", "配送映射列表"},
	{"DELETE", "/api/v1/shipping-mappings/:id", true, "shipping_mappings", "删除配送映射"},

	// Shipping Calculator
	{"POST", "/api/v1/shipping/calculate", true, "shipping_calculator", "计算运费"},

	// Product Markets
	{"GET", "/api/v1/product-markets", true, "product_markets", "商品市场列表"},
	{"POST", "/api/v1/product-markets/push", true, "product_markets", "推送商品到市场"},
	{"PUT", "/api/v1/product-markets/:id", true, "product_markets", "更新商品市场"},
	{"DELETE", "/api/v1/product-markets/:id", true, "product_markets", "从市场移除"},

	// Shop Settings
	{"GET", "/api/v1/shop/settings", true, "shop", "获取店铺设置"},
	{"PUT", "/api/v1/shop/settings", true, "shop", "更新店铺设置"},
	{"GET", "/api/v1/shop/shipping-settings", true, "shop", "获取配送设置"},
	{"PUT", "/api/v1/shop/shipping-settings", true, "shop", "更新配送设置"},
	{"GET", "/api/v1/shop/payment-settings", true, "shop", "获取支付设置"},
	{"PUT", "/api/v1/shop/payment-settings", true, "shop", "更新支付设置"},
	{"GET", "/api/v1/shop/business-hours", true, "shop", "获取营业时间"},
	{"PUT", "/api/v1/shop/business-hours", true, "shop", "更新营业时间"},
	{"GET", "/api/v1/shop/notification-settings", true, "shop", "获取通知设置"},
	{"PUT", "/api/v1/shop/notification-settings", true, "shop", "更新通知设置"},

	// Storefront
	{"GET", "/api/v1/storefront/themes", true, "storefront", "主题列表"},
	{"GET", "/api/v1/storefront/pages", true, "storefront", "页面列表"},

	// Reviews
	{"GET", "/api/v1/reviews", true, "reviews", "评价列表"},
	{"GET", "/api/v1/reviews/:id", true, "reviews", "评价详情"},
	{"PUT", "/api/v1/reviews/:id/reply", true, "reviews", "回复评价"},

	// Promotions
	{"GET", "/api/v1/promotions", true, "promotions", "促销活动列表"},
	{"POST", "/api/v1/promotions", true, "promotions", "创建促销活动"},
	{"GET", "/api/v1/promotions/:id", true, "promotions", "促销活动详情"},
	{"PUT", "/api/v1/promotions/:id", true, "promotions", "更新促销活动"},
	{"DELETE", "/api/v1/promotions/:id", true, "promotions", "删除促销活动"},

	// Coupons
	{"GET", "/api/v1/coupons", true, "coupons", "优惠券列表"},
	{"POST", "/api/v1/coupons", true, "coupons", "创建优惠券"},
	{"GET", "/api/v1/coupons/:id", true, "coupons", "优惠券详情"},
	{"PUT", "/api/v1/coupons/:id", true, "coupons", "更新优惠券"},
	{"DELETE", "/api/v1/coupons/:id", true, "coupons", "删除优惠券"},
	{"POST", "/api/v1/coupons/:id/generate-codes", true, "coupons", "生成优惠券码"},

	// Points
	{"GET", "/api/v1/points/records", true, "points", "积分记录列表"},
	{"GET", "/api/v1/points/balance/:user_id", true, "points", "用户积分余额"},

	// Uploads
	{"POST", "/api/v1/uploads", true, "uploads", "上传文件"},
	{"DELETE", "/api/v1/uploads/:id", true, "uploads", "删除上传文件"},
}

// TestAllEndpointsHaveDocumentation verifies all endpoints are documented
func TestAllEndpointsHaveDocumentation(t *testing.T) {
	endpointGroups := make(map[string][]string)
	for _, ep := range AllAdminAPIEndpoints {
		endpointGroups[ep.Group] = append(endpointGroups[ep.Group], ep.Method+" "+ep.Path)
	}

	// Verify we have endpoints in expected groups
	expectedGroups := []string{
		"auth", "users", "admin_users", "products", "categories",
		"brands", "roles", "dashboard", "shipments", "refunds",
		"fulfillment_orders", "fulfillment_statistics", "payments",
		"regions", "warehouses", "shipping_templates", "shipping_mappings",
		"shipping_calculator", "product_markets", "shop", "storefront",
		"reviews", "promotions", "coupons", "points", "uploads",
	}

	t.Run("endpoint groups coverage", func(t *testing.T) {
		for _, group := range expectedGroups {
			_, exists := endpointGroups[group]
			assert.True(t, exists, "Missing endpoint group: %s", group)
		}
	})

	t.Run("auth endpoints have correct auth requirement", func(t *testing.T) {
		for _, ep := range AllAdminAPIEndpoints {
			if ep.Group == "auth" {
				assert.False(t, ep.AuthRequired, "Auth endpoint %s should NOT require auth", ep.Path)
			} else {
				assert.True(t, ep.AuthRequired, "Non-auth endpoint %s should require auth", ep.Path)
			}
		}
	})

	t.Run("total endpoint count", func(t *testing.T) {
		// This provides a baseline for test coverage
		assert.Greater(t, len(AllAdminAPIEndpoints), 100, "Should have over 100 endpoints documented")
	})
}

// TestEndpointPathValidation validates endpoint path formats
func TestEndpointPathValidation(t *testing.T) {
	for _, ep := range AllAdminAPIEndpoints {
		t.Run(ep.Method+" "+ep.Path, func(t *testing.T) {
			// Paths should start with /api/v1
			assert.True(t, len(ep.Path) >= 7, "Path too short")
			assert.True(t, ep.Path[:7] == "/api/v1" || ep.Path[:8] == "/api/v1/", "Path should start with /api/v1")

			// Valid HTTP methods
			validMethods := map[string]bool{
				"GET":    true,
				"POST":   true,
				"PUT":    true,
				"DELETE": true,
			}
			assert.True(t, validMethods[ep.Method], "Invalid HTTP method: %s", ep.Method)

			// Description should not be empty
			assert.NotEmpty(t, ep.Description, "Endpoint should have a description")
		})
	}
}

// TestAPIEndpointGroups provides a summary of endpoints per group
func TestAPIEndpointGroups(t *testing.T) {
	endpointCounts := make(map[string]int)
	for _, ep := range AllAdminAPIEndpoints {
		endpointCounts[ep.Group]++
	}

	t.Log("Endpoint counts by group:")
	for group, count := range endpointCounts {
		t.Logf("  %s: %d endpoints", group, count)
	}

	// Total count
	total := 0
	for _, count := range endpointCounts {
		total += count
	}
	t.Logf("Total: %d endpoints", total)

	assert.Greater(t, total, 0, "Should have documented endpoints")
}
