# Missing APIs Design Document

## Overview

This document outlines the design for missing backend APIs in the ShopJoy admin system. Each API follows the existing project patterns (go-zero framework) and RESTful conventions.

---

## 1. Role Management APIs

### Context
The role entity and repository already exist (`admin/internal/domain/role/entity.go`), but there's no corresponding API definition file. The database schema includes `roles`, `permissions`, `user_roles`, and `role_permissions` tables.

### API Definition (`admin/desc/role.api`)

```go
syntax = "v1"

info (
    title:   "Role Management API"
    desc:    "角色管理相关接口"
    version: "v1"
)

type (
    // ===================== Role Types =====================
    RoleInfo {
        ID          int64  `json:"id"`
        Name        string `json:"name"`
        Code        string `json:"code"`
        Description string `json:"description"`
        Status      int8   `json:"status"`
        StatusText  string `json:"status_text"`
        IsSystem    bool   `json:"is_system"`
        CreatedAt   string `json:"created_at"`
        UpdatedAt   string `json:"updated_at"`
    }

    PermissionInfo {
        ID       int64  `json:"id"`
        Name     string `json:"name"`
        Code     string `json:"code"`
        Type     int8   `json:"type"` // 0=menu, 1=button, 2=api
        TypeText string `json:"type_text"`
        ParentID int64  `json:"parent_id"`
        Path     string `json:"path"`
        Icon     string `json:"icon"`
        Sort     int    `json:"sort"`
    }

    RoleWithPermissions {
        RoleInfo
        Permissions []*PermissionInfo `json:"permissions"`
    }

    // ===================== Request/Response =====================
    ListRolesRequest {
        Page     int    `form:"page,default=1"`
        PageSize int    `form:"page_size,default=20"`
        Name     string `form:"name,optional"`
        Code     string `form:"code,optional"`
        Status   int8   `form:"status,optional"`
    }

    ListRolesResponse {
        List     []*RoleInfo `json:"list"`
        Total    int64       `json:"total"`
        Page     int         `json:"page"`
        PageSize int         `json:"page_size"`
    }

    CreateRoleRequest {
        Name        string   `json:"name" validate:"required"`
        Code        string   `json:"code" validate:"required"`
        Description string   `json:"description,optional"`
        PermissionIDs []int64 `json:"permission_ids,optional"`
    }

    CreateRoleResponse {
        ID int64 `json:"id"`
    }

    UpdateRoleRequest {
        ID          int64  `path:"id"`
        Name        string `json:"name,optional"`
        Description string `json:"description,optional"`
    }

    RoleIDRequest {
        ID int64 `path:"id"`
    }

    UpdateRolePermissionsRequest {
        ID            int64   `path:"id"`
        PermissionIDs []int64 `json:"permission_ids"`
    }

    UpdateRoleStatusRequest {
        ID     int64 `path:"id"`
        Status int8  `json:"status"` // 0=disabled, 1=enabled
    }

    ListPermissionsResponse {
        List []*PermissionInfo `json:"list"`
    }
)

@server (
    group:      roles
    middleware: AuthMiddleware
)
service admin-api {
    @doc "获取角色列表"
    @handler ListRolesHandler
    get /api/v1/roles (ListRolesRequest) returns (ListRolesResponse)

    @doc "创建角色"
    @handler CreateRoleHandler
    post /api/v1/roles (CreateRoleRequest) returns (CreateRoleResponse)

    @doc "获取角色详情"
    @handler GetRoleHandler
    get /api/v1/roles/:id (RoleIDRequest) returns (RoleWithPermissions)

    @doc "更新角色"
    @handler UpdateRoleHandler
    put /api/v1/roles/:id (UpdateRoleRequest) returns (RoleInfo)

    @doc "删除角色"
    @handler DeleteRoleHandler
    delete /api/v1/roles/:id (RoleIDRequest)

    @doc "更新角色状态"
    @handler UpdateRoleStatusHandler
    put /api/v1/roles/:id/status (UpdateRoleStatusRequest) returns (RoleInfo)

    @doc "更新角色权限"
    @handler UpdateRolePermissionsHandler
    put /api/v1/roles/:id/permissions (UpdateRolePermissionsRequest)

    @doc "获取所有权限列表"
    @handler ListPermissionsHandler
    get /api/v1/permissions (ListPermissionsResponse) returns (ListPermissionsResponse)
}
```

### Authentication Requirements
- All endpoints require `AuthMiddleware`
- Role creation/update restricted to admin type 1 (platform super admin) or 2 (tenant admin)
- System roles (`is_system = true`) cannot be deleted or have code changed

### Error Responses
| Error Code | HTTP Status | Message |
|------------|-------------|---------|
| 100001 | 404 | role not found |
| 100002 | 409 | duplicate role |
| 100003 | 400 | invalid role |
| 100004 | 400 | cannot modify system role |

### Business Logic Considerations
1. **Role hierarchy**: System roles cannot be modified by tenant admins
2. **Permission scoping**: Tenant admins can only assign permissions they possess
3. **Audit trail**: All role changes should be logged with operator information
4. **Default permissions**: New roles start with no permissions

---

## 2. Order Management APIs (Additional Endpoints)

### Context
The `fulfillment.api` already has comprehensive order management. Missing endpoints for canceling orders and sending payment reminders.

### Additional Types to Add (`admin/desc/fulfillment.api`)

```go
// ===================== Order Cancellation Types =====================
CancelOrderReq {
    ID     int64  `path:"id"`
    Reason string `json:"reason" validate:"required,max=200"`
}

CancelOrderResp {
    OrderID    int64  `json:"order_id"`
    OrderNo    string `json:"order_no"`
    Status     string `json:"status"`
    CancelledAt string `json:"cancelled_at"`
}

// ===================== Payment Reminder Types =====================
RemindPaymentReq {
    ID int64 `path:"id"`
}

RemindPaymentResp {
    OrderID    int64  `json:"order_id"`
    OrderNo    string `json:"order_no"`
    RemindedAt string `json:"reminded_at"`
    Message    string `json:"message"`
}
```

### Additional Endpoints

```go
@server (
    group:      fulfillment_orders
    middleware: AuthMiddleware
)
service admin-api {
    // ... existing endpoints ...

    @doc "取消订单"
    @handler CancelOrderHandler
    put /api/v1/orders/:id/cancel (CancelOrderReq) returns (CancelOrderResp)

    @doc "发送支付提醒"
    @handler RemindPaymentHandler
    post /api/v1/orders/:id/remind-payment (RemindPaymentReq) returns (RemindPaymentResp)
}
```

### Error Responses
| Error Code | HTTP Status | Message |
|------------|-------------|---------|
| 40001 | 404 | order not found |
| 40002 | 400 | invalid order status (cannot cancel completed/shipped orders) |
| 40014 | 400 | cancel reason required |
| 40015 | 400 | payment reminder already sent recently |
| 40016 | 400 | order already paid |

### Business Logic Considerations
1. **Cancel restrictions**: Cannot cancel orders that are already shipped, delivered, or completed
2. **Refund handling**: If order was paid, cancellation should trigger automatic refund
3. **Inventory restoration**: Upon cancellation, restore inventory for all order items
4. **Payment reminder throttling**: Limit to 1 reminder per hour per order
5. **Notification**: Send email/SMS notification to customer upon cancellation

---

## 3. Dashboard Statistics API

### Context
The frontend dashboard (`shop-admin/src/views/dashboard/index.vue`) currently uses mock data. Need a unified backend API to aggregate real statistics.

### API Definition (`admin/desc/dashboard.api`)

```go
syntax = "v1"

info (
    title:   "Dashboard Statistics API"
    desc:    "仪表盘统计数据接口"
    version: "v1"
)

type (
    // ===================== Overview Stats =====================
    DashboardOverviewRequest {
        // No parameters - returns current tenant's overview
    }

    DashboardOverviewResponse {
        // Today's metrics
        TodayOrders    int64  `json:"today_orders"`
        TodaySales     string `json:"today_sales"`     // Decimal as string
        TodayGrowth    string `json:"today_growth"`    // Percentage string "12.5%"
        YesterdaySales string `json:"yesterday_sales"` // For comparison

        // Cumulative metrics
        TotalProducts  int64  `json:"total_products"`
        TotalUsers     int64  `json:"total_users"`
        NewUsersToday  int64  `json:"new_users_today"`

        // Currency
        Currency       string `json:"currency"`
    }

    // ===================== Sales Trend =====================
    SalesTrendRequest {
        Period string `form:"period,default=week"` // week, month, year
    }

    SalesTrendData {
        Date   string `json:"date"`   // "2024-03-18"
        Sales  string `json:"sales"`  // Decimal as string
        Orders int64  `json:"orders"`
    }

    SalesTrendResponse {
        Period string             `json:"period"`
        Data   []*SalesTrendData  `json:"data"`
        Currency string           `json:"currency"`
    }

    // ===================== Order Status Distribution =====================
    OrderStatusDistributionRequest {}

    OrderStatusItem {
        Status    string `json:"status"`
        StatusText string `json:"status_text"`
        Count     int64  `json:"count"`
        Percentage string `json:"percentage"` // "35.5%"
        Color     string `json:"color"`
    }

    OrderStatusDistributionResponse {
        List []*OrderStatusItem `json:"list"`
        Total int64             `json:"total"`
    }

    // ===================== Top Products =====================
    TopProductsRequest {
        Limit int `form:"limit,default=5"`
        Period string `form:"period,default=week"` // week, month, all
    }

    TopProductItem {
        ProductID   int64  `json:"product_id"`
        ProductName string `json:"product_name"`
        Image       string `json:"image"`
        Sales       int64  `json:"sales"`
        Revenue     string `json:"revenue"`
    }

    TopProductsResponse {
        List     []*TopProductItem `json:"list"`
        Currency string            `json:"currency"`
    }

    // ===================== Pending Orders =====================
    PendingOrdersRequest {
        Limit int `form:"limit,default=5"`
    }

    PendingOrderItem {
        OrderID   int64  `json:"order_id"`
        OrderNo   string `json:"order_no"`
        PayAmount string `json:"pay_amount"`
        Status    string `json:"status"`
        StatusText string `json:"status_text"`
        CreatedAt string `json:"created_at"`
        UserName  string `json:"user_name,optional"`
    }

    PendingOrdersResponse {
        List []*PendingOrderItem `json:"list"`
        Total int64              `json:"total"`
    }

    // ===================== Recent Activities =====================
    RecentActivitiesRequest {
        Limit int `form:"limit,default=10"`
    }

    ActivityItem {
        ID        int64  `json:"id"`
        Type      string `json:"type"` // order_created, payment_received, product_low_stock, etc.
        Content   string `json:"content"`
        Time      string `json:"time"`
        Operator  string `json:"operator,optional"`
    }

    RecentActivitiesResponse {
        List []*ActivityItem `json:"list"`
    }

    // ===================== Unified Dashboard API =====================
    GetDashboardRequest {}

    GetDashboardResponse {
        Overview    *DashboardOverviewResponse    `json:"overview"`
        StatusDistribution *OrderStatusDistributionResponse `json:"status_distribution"`
        PendingOrders []*PendingOrderItem         `json:"pending_orders"`
        TopProducts []*TopProductItem             `json:"top_products"`
        RecentActivities []*ActivityItem          `json:"recent_activities"`
    }
)

@server (
    group:      dashboard
    middleware: AuthMiddleware
)
service admin-api {
    @doc "获取仪表盘概览数据"
    @handler GetDashboardOverviewHandler
    get /api/v1/dashboard/overview (DashboardOverviewRequest) returns (DashboardOverviewResponse)

    @doc "获取销售趋势"
    @handler GetSalesTrendHandler
    get /api/v1/dashboard/sales-trend (SalesTrendRequest) returns (SalesTrendResponse)

    @doc "获取订单状态分布"
    @handler GetOrderStatusDistributionHandler
    get /api/v1/dashboard/order-status (OrderStatusDistributionRequest) returns (OrderStatusDistributionResponse)

    @doc "获取热销商品TOP"
    @handler GetTopProductsHandler
    get /api/v1/dashboard/top-products (TopProductsRequest) returns (TopProductsResponse)

    @doc "获取待处理订单"
    @handler GetPendingOrdersHandler
    get /api/v1/dashboard/pending-orders (PendingOrdersRequest) returns (PendingOrdersResponse)

    @doc "获取最近活动"
    @handler GetRecentActivitiesHandler
    get /api/v1/dashboard/activities (RecentActivitiesRequest) returns (RecentActivitiesResponse)

    @doc "获取仪表盘所有数据（聚合接口）"
    @handler GetDashboardHandler
    get /api/v1/dashboard (GetDashboardRequest) returns (GetDashboardResponse)
}
```

### Performance Considerations
1. **Caching**: Cache overview stats for 1 minute, use Redis with tenant-based keys
2. **Async aggregation**: The unified `/api/v1/dashboard` endpoint should fetch data in parallel
3. **Materialized views**: Consider pre-computed daily/weekly aggregations for large datasets

### Business Logic Considerations
1. **Tenant isolation**: All statistics scoped to current tenant
2. **Currency**: Respect market currency settings
3. **Timezone**: All times in UTC, frontend converts for display

---

## 4. Shop/Tenant Settings APIs

### Context
The `tenants` table exists with tenant configuration. Need APIs for shop administrators to manage their shop settings.

### API Definition (`admin/desc/shop.api`)

```go
syntax = "v1"

info (
    title:   "Shop Settings API"
    desc:    "店铺设置相关接口"
    version: "v1"
)

type (
    // ===================== Shop Settings Types =====================
    ShopSettings {
        // Basic Info
        ID           int64  `json:"id"`
        Name         string `json:"name"`
        Code         string `json:"code"`
        Logo         string `json:"logo"`
        Description  string `json:"description,optional"`

        // Contact Info
        ContactName  string `json:"contact_name"`
        ContactPhone string `json:"contact_phone"`
        ContactEmail string `json:"contact_email"`
        Address      string `json:"address,optional"`

        // Domain Settings
        Domain       string `json:"domain"`
        CustomDomain string `json:"custom_domain,optional"`

        // Branding
        PrimaryColor   string `json:"primary_color,optional"`
        SecondaryColor string `json:"secondary_color,optional"`
        Favicon        string `json:"favicon,optional"`

        // Business Settings
        DefaultCurrency string `json:"default_currency"`
        DefaultLanguage string `json:"default_language"`
        Timezone        string `json:"timezone"`

        // Status
        Status    int8   `json:"status"`
        StatusText string `json:"status_text"`
        Plan      int8   `json:"plan"`
        PlanText  string `json:"plan_text"`
        ExpireAt  string `json:"expire_at,optional"`

        // Timestamps
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
    }

    UpdateShopSettingsRequest {
        // Basic Info
        Name         string `json:"name,optional" validate:"omitempty,min=2,max=100"`
        Logo         string `json:"logo,optional"`
        Description  string `json:"description,optional" validate:"omitempty,max=500"`

        // Contact Info
        ContactName  string `json:"contact_name,optional"`
        ContactPhone string `json:"contact_phone,optional"`
        ContactEmail string `json:"contact_email,optional" validate:"omitempty,email"`
        Address      string `json:"address,optional"`

        // Domain Settings
        CustomDomain string `json:"custom_domain,optional"`

        // Branding
        PrimaryColor   string `json:"primary_color,optional" validate:"omitempty,hexcolor"`
        SecondaryColor string `json:"secondary_color,optional" validate:"omitempty,hexcolor"`
        Favicon        string `json:"favicon,optional"`

        // Business Settings
        DefaultLanguage string `json:"default_language,optional"`
        Timezone        string `json:"timezone,optional"`
    }

    // ===================== Business Hours =====================
    BusinessHours {
        DayOfWeek int8   `json:"day_of_week"` // 0=Sunday, 1=Monday...
        OpenTime  string `json:"open_time"`   // "09:00"
        CloseTime string `json:"close_time"`  // "18:00"
        IsClosed  bool   `json:"is_closed"`
    }

    UpdateBusinessHoursRequest {
        Hours []*BusinessHours `json:"hours"`
    }

    // ===================== Notification Settings =====================
    NotificationSettings {
        OrderCreated      bool `json:"order_created"`
        OrderPaid         bool `json:"order_paid"`
        OrderShipped      bool `json:"order_shipped"`
        OrderCancelled    bool `json:"order_cancelled"`
        LowStockAlert     bool `json:"low_stock_alert"`
        LowStockThreshold int  `json:"low_stock_threshold"`
        RefundRequested   bool `json:"refund_requested"`
        NewReview         bool `json:"new_review"`
    }

    UpdateNotificationSettingsRequest {
        NotificationSettings
    }

    // ===================== Payment Settings =====================
    PaymentSettings {
        StripeEnabled   bool   `json:"stripe_enabled"`
        StripePublicKey string `json:"stripe_public_key,optional"`
        // Add other payment methods as needed
    }

    UpdatePaymentSettingsRequest {
        StripeEnabled   bool   `json:"stripe_enabled"`
        StripeSecretKey string `json:"stripe_secret_key,optional"` // Encrypted
    }

    // ===================== Shipping Settings =====================
    ShippingSettings {
        FreeShippingThreshold string `json:"free_shipping_threshold"`
        DefaultShippingFee    string `json:"default_shipping_fee"`
        Currency              string `json:"currency"`
    }

    UpdateShippingSettingsRequest {
        FreeShippingThreshold string `json:"free_shipping_threshold,optional"`
        DefaultShippingFee    string `json:"default_shipping_fee,optional"`
    }
)

@server (
    group:      shop
    middleware: AuthMiddleware
)
service admin-api {
    @doc "获取店铺设置"
    @handler GetShopSettingsHandler
    get /api/v1/shop returns (ShopSettings)

    @doc "更新店铺设置"
    @handler UpdateShopSettingsHandler
    put /api/v1/shop (UpdateShopSettingsRequest) returns (ShopSettings)

    @doc "获取营业时间设置"
    @handler GetBusinessHoursHandler
    get /api/v1/shop/business-hours returns ([]BusinessHours)

    @doc "更新营业时间设置"
    @handler UpdateBusinessHoursHandler
    put /api/v1/shop/business-hours (UpdateBusinessHoursRequest)

    @doc "获取通知设置"
    @handler GetNotificationSettingsHandler
    get /api/v1/shop/notifications returns (NotificationSettings)

    @doc "更新通知设置"
    @handler UpdateNotificationSettingsHandler
    put /api/v1/shop/notifications (UpdateNotificationSettingsRequest) returns (NotificationSettings)

    @doc "获取支付设置"
    @handler GetPaymentSettingsHandler
    get /api/v1/shop/payment returns (PaymentSettings)

    @doc "更新支付设置"
    @handler UpdatePaymentSettingsHandler
    put /api/v1/shop/payment (UpdatePaymentSettingsRequest) returns (PaymentSettings)

    @doc "获取运费设置"
    @handler GetShippingSettingsHandler
    get /api/v1/shop/shipping returns (ShippingSettings)

    @doc "更新运费设置"
    @handler UpdateShippingSettingsHandler
    put /api/v1/shop/shipping (UpdateShippingSettingsRequest) returns (ShippingSettings)
}
```

### Error Responses
| Error Code | HTTP Status | Message |
|------------|-------------|---------|
| 90001 | 404 | tenant not found |
| 90003 | 400 | invalid domain |
| 90004 | 403 | tenant is inactive |
| 90009 | 400 | invalid custom domain format |
| 90010 | 400 | custom domain already in use |

### Business Logic Considerations
1. **Multi-tenant isolation**: Each tenant can only modify their own settings
2. **Plan restrictions**: Some features may be restricted based on plan level
3. **Domain verification**: Custom domain changes may require DNS verification
4. **Sensitive data**: Payment secret keys should be encrypted at rest
5. **Audit logging**: All setting changes should be logged

---

## 5. Warehouse Management APIs (Verification)

### Context
Warehouse management APIs already exist in `admin/desc/inventory.api`. The following endpoints are available:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/warehouses` | POST | Create warehouse |
| `/api/v1/warehouses/:id` | PUT | Update warehouse |
| `/api/v1/warehouses/:id` | GET | Get warehouse detail |
| `/api/v1/warehouses` | GET | List warehouses |
| `/api/v1/warehouses/:id/status` | PUT | Update warehouse status |
| `/api/v1/warehouses/:id/default` | PUT | Set default warehouse |
| `/api/v1/warehouses/:id` | DELETE | Delete warehouse |

### Suggested Additional Functionality

```go
// ===================== Warehouse Statistics =====================
GetWarehouseStatsReq {
    ID int64 `path:"id"`
}

WarehouseStatsResp {
    WarehouseID      int64  `json:"warehouse_id"`
    WarehouseName    string `json:"warehouse_name"`
    TotalSKUs        int64  `json:"total_skus"`
    TotalStock       int64  `json:"total_stock"`
    LowStockSKUs     int64  `json:"low_stock_skus"`
    OutOfStockSKUs   int64  `json:"out_of_stock_skus"`
    TotalValue       string `json:"total_value"` // Total inventory value
    Currency         string `json:"currency"`
    LastStockUpdate  string `json:"last_stock_update"`
}

// ===================== Warehouse Transfer =====================
TransferStockReq {
    FromWarehouseID int64  `json:"from_warehouse_id"`
    ToWarehouseID   int64  `json:"to_warehouse_id"`
    SKUCode         string `json:"sku_code"`
    Quantity        int    `json:"quantity"`
    Remark          string `json:"remark,optional"`
}

TransferStockResp {
    TransferID int64  `json:"transfer_id"`
    FromWarehouse string `json:"from_warehouse"`
    ToWarehouse   string `json:"to_warehouse"`
    SKUCode       string `json:"sku_code"`
    Quantity      int    `json:"quantity"`
    Status        string `json:"status"`
    CreatedAt     string `json:"created_at"`
}

// Add to inventory.api
@server (
    group:      warehouses
    middleware: AuthMiddleware
)
service admin-api {
    // ... existing endpoints ...

    @doc "获取仓库统计"
    @handler GetWarehouseStatsHandler
    get /api/v1/warehouses/:id/stats (GetWarehouseStatsReq) returns (WarehouseStatsResp)

    @doc "库存调拨"
    @handler TransferStockHandler
    post /api/v1/inventory/transfer (TransferStockReq) returns (TransferStockResp)
}
```

---

## 6. New Error Codes to Add

Add these error codes to `pkg/code/code.go`:

```go
// ==================== Role Module Additional Errors (100xxx) ====================
ErrRoleCannotModifySystem = &Err{HTTPCode: http.StatusBadRequest, Code: 100004, Msg: "cannot modify system role"}
ErrRoleInUse              = &Err{HTTPCode: http.StatusBadRequest, Code: 100005, Msg: "role is in use by users"}

// ==================== Order Module Additional Errors (40xxx) ====================
ErrOrderCannotCancel      = &Err{HTTPCode: http.StatusBadRequest, Code: 40014, Msg: "order cannot be cancelled in current status"}
ErrOrderCancelReasonRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 40015, Msg: "cancel reason is required"}
ErrPaymentReminderSent    = &Err{HTTPCode: http.StatusTooManyRequests, Code: 40016, Msg: "payment reminder already sent recently"}
ErrOrderAlreadyPaid       = &Err{HTTPCode: http.StatusBadRequest, Code: 40017, Msg: "order already paid, cannot send reminder"}

// ==================== Shop Module Additional Errors (110xxx) ====================
ErrShopInvalidDomain      = &Err{HTTPCode: http.StatusBadRequest, Code: 110003, Msg: "invalid domain format"}
ErrShopDomainInUse        = &Err{HTTPCode: http.StatusConflict, Code: 110004, Msg: "custom domain already in use"}
ErrShopPlanRestricted     = &Err{HTTPCode: http.StatusForbidden, Code: 110005, Msg: "feature not available in current plan"}

// ==================== Dashboard Module (220xxx) ====================
ErrDashboardDataUnavailable = &Err{HTTPCode: http.StatusServiceUnavailable, Code: 220001, Msg: "dashboard data temporarily unavailable"}

// ==================== Inventory Module Additional Errors (170xxx) ====================
ErrInventoryTransferFailed     = &Err{HTTPCode: http.StatusBadRequest, Code: 170006, Msg: "stock transfer failed"}
ErrInsufficientStockForTransfer = &Err{HTTPCode: http.StatusBadRequest, Code: 170007, Msg: "insufficient stock for transfer"}
ErrSameWarehouseTransfer       = &Err{HTTPCode: http.StatusBadRequest, Code: 170008, Msg: "cannot transfer to same warehouse"}
```

---

## 7. Database Schema Additions

### Shop Settings Table

```sql
CREATE TABLE IF NOT EXISTS `shop_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `primary_color` VARCHAR(7) DEFAULT '#6366F1' COMMENT '主题色',
    `secondary_color` VARCHAR(7) DEFAULT '#818CF8' COMMENT '次要颜色',
    `favicon` VARCHAR(255) DEFAULT '' COMMENT 'Favicon URL',
    `default_language` VARCHAR(10) DEFAULT 'zh-CN' COMMENT '默认语言',
    `timezone` VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '时区',
    `free_shipping_threshold` DECIMAL(19,4) DEFAULT 0 COMMENT '免运费门槛',
    `default_shipping_fee` DECIMAL(19,4) DEFAULT 0 COMMENT '默认运费',
    `low_stock_threshold` INT DEFAULT 10 COMMENT '低库存预警阈值',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='店铺设置表';

CREATE TABLE IF NOT EXISTS `business_hours` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `day_of_week` TINYINT NOT NULL COMMENT '星期几 0=周日',
    `open_time` TIME COMMENT '营业开始时间',
    `close_time` TIME COMMENT '营业结束时间',
    `is_closed` TINYINT NOT NULL DEFAULT 0 COMMENT '是否休息',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_day` (`tenant_id`, `day_of_week`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='营业时间表';

CREATE TABLE IF NOT EXISTS `notification_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_created` TINYINT NOT NULL DEFAULT 1,
    `order_paid` TINYINT NOT NULL DEFAULT 1,
    `order_shipped` TINYINT NOT NULL DEFAULT 1,
    `order_cancelled` TINYINT NOT NULL DEFAULT 1,
    `low_stock_alert` TINYINT NOT NULL DEFAULT 1,
    `low_stock_threshold` INT NOT NULL DEFAULT 10,
    `refund_requested` TINYINT NOT NULL DEFAULT 1,
    `new_review` TINYINT NOT NULL DEFAULT 1,
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知设置表';

-- Stock Transfer Log
CREATE TABLE IF NOT EXISTS `stock_transfers` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `transfer_no` VARCHAR(32) NOT NULL COMMENT '调拨单号',
    `from_warehouse_id` BIGINT NOT NULL COMMENT '源仓库ID',
    `to_warehouse_id` BIGINT NOT NULL COMMENT '目标仓库ID',
    `sku_code` VARCHAR(64) NOT NULL COMMENT 'SKU编码',
    `quantity` INT NOT NULL COMMENT '调拨数量',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0=待处理, 1=已完成, 2=已取消',
    `remark` VARCHAR(500) DEFAULT '' COMMENT '备注',
    `operator_id` BIGINT NOT NULL COMMENT '操作人ID',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_transfer_no` (`transfer_no`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_sku_code` (`sku_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存调拨记录表';
```

---

## 8. Implementation Priority

| Priority | Module | Effort | Impact |
|----------|--------|--------|--------|
| P0 | Dashboard Statistics | Medium | High - Replaces mock data with real metrics |
| P0 | Role Management | Medium | High - Security/access control foundation |
| P1 | Order Cancel/Remind | Low | Medium - Essential order operations |
| P1 | Shop Settings | Medium | Medium - Shop branding configuration |
| P2 | Warehouse Statistics | Low | Low - Nice to have for inventory overview |
| P2 | Stock Transfer | Medium | Medium - Multi-warehouse operations |

---

## 9. API Versioning Strategy

All APIs follow the existing versioning pattern:
- URI versioning: `/api/v1/...`
- No breaking changes within v1
- New fields can be added to responses
- Deprecated fields marked but not removed
- Future v2 would introduce breaking changes with migration guide