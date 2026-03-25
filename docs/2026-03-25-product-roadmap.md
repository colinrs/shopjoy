# ShopJoy Admin System - Product Roadmap

## Executive Summary

This document outlines the feature gaps identified in the ShopJoy e-commerce admin system and provides a prioritized implementation roadmap. The analysis covers backend APIs, frontend pages, and API integration opportunities.

---

## Gap Analysis Summary

### Priority Matrix

| Gap Category | Business Impact | Implementation Effort | Priority |
|-------------|-----------------|---------------------|----------|
| Dashboard Real Statistics | HIGH | MEDIUM | P0 |
| Order Operations (Cancel/Remind) | HIGH | LOW | P0 |
| Role Management API | MEDIUM | LOW | P1 |
| Warehouse Management Page | HIGH | MEDIUM | P1 |
| Category/Brand Advanced Features | MEDIUM | LOW | P2 |
| Shop/Tenant Settings | MEDIUM | MEDIUM | P2 |
| Review Product Stats Integration | LOW | LOW | P3 |

---

## Release Planning

### Release v1.1 - Core Operations Enhancement (Sprint 1-2)

**Theme:** Complete essential order operations and real-time business intelligence

#### Feature 1: Dashboard Real Statistics API

**Business Impact:**
- Enables data-driven decision making
- Provides real-time visibility into business health
- Reduces manual reporting overhead by 80%

**Backend APIs Required:**
```
GET /api/v1/dashboard/stats          - Overall statistics
GET /api/v1/dashboard/sales-trend    - Sales trend data (7d/30d/90d)
GET /api/v1/dashboard/order-status   - Order status distribution
GET /api/v1/dashboard/hot-products   - Top selling products
GET /api/v1/dashboard/activities     - Recent activities
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| DASH-001 | As an admin, I want to see today's real order count so I can track business activity | Given I am on the dashboard, When the page loads, Then I see the actual count of orders created today from the database |
| DASH-002 | As an admin, I want to see today's real sales revenue so I can track daily performance | Given I am on the dashboard, When the page loads, Then I see the sum of all paid order amounts for today |
| DASH-003 | As an admin, I want to see sales trend charts so I can identify patterns | Given I select a time range (week/month/year), When I view the sales chart, Then I see accurate historical data matching the database |
| DASH-004 | As an admin, I want to see pending orders count so I can prioritize work | Given I am on the dashboard, When the page loads, Then I see real count of orders requiring action |
| DASH-005 | As an admin, I want to see top 5 hot products so I can optimize inventory | Given I am on the dashboard, When I view the hot products section, Then I see products ranked by actual sales volume |

**Technical Implementation:**

Backend:
- Create `dashboard` module in admin service
- Implement aggregation queries with caching (Redis, 5-minute TTL)
- Define error codes in range 200xxx (Shared module)

Frontend:
- Replace mock data in `/views/dashboard/index.vue`
- Add loading states and error handling
- Implement auto-refresh (configurable interval)

**Files to Create/Modify:**
```
admin/desc/dashboard.api              # New API definition
admin/internal/logic/dashboard/       # Business logic
shop-admin/src/api/dashboard.ts       # Frontend API
shop-admin/src/views/dashboard/index.vue  # Update existing
```

---

#### Feature 2: Order Cancel API

**Business Impact:**
- Completes order lifecycle management
- Reduces customer support tickets by 30%
- Enables automated cancellation workflows

**Backend API Required:**
```
PUT /api/v1/orders/:id/cancel
Request: { reason: string (required, 5-200 chars) }
Response: { order_id: number, status: string, cancelled_at: string }
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| CANCEL-001 | As an admin, I want to cancel pending orders so I can handle customer requests | Given an order with status "pending", When I cancel with a reason, Then the order status becomes "cancelled" |
| CANCEL-002 | As an admin, I want to cancel paid orders before shipment for customer requests | Given a paid order not yet shipped, When I cancel with a reason, Then the order is cancelled and refund is initiated |
| CANCEL-003 | As an admin, I want to provide a cancellation reason for audit trail | Given I cancel an order, When I enter a reason, Then the reason is saved and visible in order history |
| CANCEL-004 | As an admin, I want to see cancelled orders in the order list | Given an order was cancelled, When I view orders, Then I can see the "cancelled" status and reason |

**Business Rules:**
- Can cancel orders with status: `pending`, `paid`
- Cannot cancel orders with status: `shipped`, `completed`
- Cancellation triggers refund for paid orders
- Reason is mandatory (5-200 characters)

**Error Codes:**
```go
ErrOrderCannotCancel     = &Err{HTTPCode: 400, Code: 40010, Msg: "order cannot be cancelled"}
ErrOrderCancelReasonRequired = &Err{HTTPCode: 400, Code: 40011, Msg: "cancellation reason required"}
```

---

#### Feature 3: Payment Reminder API

**Business Impact:**
- Improves order conversion rate by 15-20%
- Reduces abandoned cart rate
- Automates customer communication

**Backend API Required:**
```
POST /api/v1/orders/:id/remind-payment
Request: {} (empty body)
Response: { success: boolean, message: string }
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| REMIND-001 | As an admin, I want to send payment reminders for pending orders | Given an unpaid order, When I click "Remind", Then a notification is sent to the customer |
| REMIND-002 | As an admin, I want to see if reminder was sent recently | Given an order has pending status, When I view order details, Then I see when the last reminder was sent |
| REMIND-003 | As an admin, I want to limit reminder frequency to prevent spam | Given a reminder was sent within 24 hours, When I try to send another, Then the system shows an error message |

**Business Rules:**
- Only applicable for orders with status `pending`
- Maximum 3 reminders per order
- Minimum 24 hours between reminders
- Integration with notification service (SMS/Email)

**Error Codes:**
```go
ErrOrderNotPending       = &Err{HTTPCode: 400, Code: 40012, Msg: "order is not pending payment"}
ErrReminderLimitExceeded = &Err{HTTPCode: 429, Code: 40013, Msg: "reminder limit exceeded"}
ErrReminderTooFrequent   = &Err{HTTPCode: 429, Code: 40014, Msg: "please wait before sending another reminder"}
```

---

### Release v1.2 - Admin Operations Enhancement (Sprint 3-4)

**Theme:** Complete administrative user management and inventory operations

#### Feature 4: Role Management API

**Business Impact:**
- Enables proper access control
- Supports team scaling
- Meets compliance requirements

**Backend API Required:**
```
GET /api/v1/roles
Response: {
  list: [{
    id: number,
    name: string,
    code: string,
    description: string,
    permissions: string[],
    user_count: number,
    status: number,
    created_at: string
  }],
  total: number
}
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| ROLE-001 | As an admin, I want to see available roles when assigning roles to users | Given I am editing an admin user, When I open the role assignment dialog, Then I see a list of all active roles |
| ROLE-002 | As an admin, I want to see how many users have each role | Given I view the roles list, When I see a role, Then I can see the count of users assigned to it |
| ROLE-003 | As an admin, I want to search roles by name | Given I am in the role assignment dialog, When I type in the search box, Then roles are filtered by name |

**Technical Implementation:**

1. Create `role.api` definition file
2. Implement role logic in `admin/internal/logic/roles/`
3. Update frontend `RoleAssignDialog.vue` to use the API

**Files to Create/Modify:**
```
admin/desc/role.api                              # New API definition
admin/internal/logic/roles/list_roles_logic.go   # Business logic
shop-admin/src/api/role.ts                       # Frontend API
shop-admin/src/views/admin-users/components/RoleAssignDialog.vue
```

---

#### Feature 5: Warehouse Management Page

**Business Impact:**
- Enables multi-warehouse inventory tracking
- Supports geographic distribution strategy
- Improves fulfillment efficiency by 25%

**Backend APIs Already Available:**
```
POST   /api/v1/warehouses           - Create warehouse
PUT    /api/v1/warehouses/:id       - Update warehouse
GET    /api/v1/warehouses/:id       - Get warehouse detail
GET    /api/v1/warehouses           - List warehouses
PUT    /api/v1/warehouses/:id/status - Update status
PUT    /api/v1/warehouses/:id/default - Set as default
DELETE /api/v1/warehouses/:id       - Delete warehouse
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| WH-001 | As an admin, I want to create warehouses to manage inventory locations | Given I am on the warehouse page, When I click "Create", Then I can add a new warehouse with code, name, country, and address |
| WH-002 | As an admin, I want to edit warehouse information | Given a warehouse exists, When I click "Edit", Then I can modify its name, country, and address |
| WH-003 | As an admin, I want to set a default warehouse | Given multiple warehouses exist, When I set one as default, Then it is marked and used for default inventory operations |
| WH-004 | As an admin, I want to enable/disable warehouses | Given a warehouse exists, When I toggle its status, Then inventory operations respect the status |
| WH-005 | As an admin, I want to see warehouse inventory summary | Given I view a warehouse, When I look at the detail, Then I see total SKUs and stock quantities |
| WH-006 | As an admin, I want to delete unused warehouses | Given a warehouse has no inventory, When I delete it, Then it is removed from the system |
| WH-007 | As an admin, I want to see which warehouse is default | Given I view the warehouse list, When I look at the table, Then the default warehouse is clearly marked |

**Frontend Implementation:**

Create new page structure:
```
shop-admin/src/views/inventory/warehouses/
  index.vue                 # List view
  [id].vue                  # Detail/Edit view
  components/
    WarehouseForm.vue       # Create/Edit form
    WarehouseCard.vue       # Warehouse info card
```

Update router configuration to add the new route.

**UI/UX Design:**
- Table view with columns: Code, Name, Country, Status, Default, Actions
- Status toggle switch
- Default warehouse badge
- Create/Edit modal form
- Confirmation dialog for delete

---

### Release v1.3 - Content Management Enhancement (Sprint 5-6)

**Theme:** Advanced category and brand management with multi-market support

#### Feature 6: Category Advanced Features Integration

**Business Impact:**
- Improves catalog organization efficiency
- Supports multi-market strategy
- Enhances SEO optimization

**Backend APIs Already Available:**
```
PUT /api/v1/categories/:id/status            - Update status
PUT /api/v1/categories/sort                  - Batch sort
PUT /api/v1/categories/:id/move              - Move to new parent
GET /api/v1/categories/:id/product-count     - Get product count
PUT /api/v1/categories/:id/market-visibility - Set market visibility
GET /api/v1/categories/:id/market-visibility - Get market visibility
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| CAT-001 | As an admin, I want to enable/disable categories | Given a category, When I toggle its status, Then it is hidden/shown in the storefront |
| CAT-002 | As an admin, I want to reorder categories via drag-and-drop | Given categories at the same level, When I drag to reorder, Then the sort order is updated |
| CAT-003 | As an admin, I want to move a category to a different parent | Given a category, When I select "Move" and choose a new parent, Then the category is relocated |
| CAT-004 | As an admin, I want to see how many products are in each category | Given I view categories, When I look at the list, Then product count is displayed |
| CAT-005 | As an admin, I want to control category visibility per market | Given a category and multiple markets, When I configure visibility, Then the category shows only in selected markets |

**Frontend Implementation:**

Enhance existing `/views/categories/index.vue`:
- Add status toggle column
- Add drag-and-drop for sorting
- Add "Move" action in dropdown
- Add product count column
- Add market visibility configuration dialog

---

#### Feature 7: Brand Advanced Features Integration

**Business Impact:**
- Enables brand page management
- Supports multi-market brand strategy
- Improves brand discovery

**Backend APIs Already Available:**
```
PUT /api/v1/brands/:id/status            - Update status
PUT /api/v1/brands/:id/toggle-page       - Enable/disable brand page
GET /api/v1/brands/:id/product-count     - Get product count
PUT /api/v1/brands/:id/market-visibility - Set market visibility
GET /api/v1/brands/:id/market-visibility - Get market visibility
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| BRAND-001 | As an admin, I want to enable/disable brands | Given a brand, When I toggle its status, Then it is hidden/shown in the storefront |
| BRAND-002 | As an admin, I want to enable brand pages for featured brands | Given a brand, When I enable the brand page, Then a dedicated page is generated |
| BRAND-003 | As an admin, I want to see product count per brand | Given I view brands, When I look at the list, Then product count is displayed |
| BRAND-004 | As an admin, I want to control brand visibility per market | Given a brand and multiple markets, When I configure visibility, Then the brand shows only in selected markets |

**Frontend Implementation:**

Enhance existing `/views/brands/index.vue`:
- Add status toggle column
- Add "Brand Page" toggle switch
- Add product count column
- Add market visibility configuration dialog

---

### Release v1.4 - Platform Configuration (Sprint 7-8)

**Theme:** Multi-tenant and shop configuration management

#### Feature 8: Shop/Tenant Settings APIs

**Business Impact:**
- Enables self-service configuration
- Reduces IT support tickets by 40%
- Supports white-label customization

**Backend APIs Required:**
```
GET  /api/v1/shop/settings           - Get shop settings
PUT  /api/v1/shop/settings           - Update shop settings
GET  /api/v1/shop/payment-methods    - Get available payment methods
PUT  /api/v1/shop/payment-methods    - Configure payment methods
GET  /api/v1/shop/shipping-methods   - Get shipping configuration
PUT  /api/v1/shop/shipping-methods   - Configure shipping
GET  /api/v1/shop/notifications      - Get notification settings
PUT  /api/v1/shop/notifications      - Configure notifications
```

**User Stories:**

| ID | User Story | Acceptance Criteria |
|----|-----------|-------------------|
| SHOP-001 | As a shop admin, I want to configure basic shop information | Given I am in settings, When I update shop name/logo/contact, Then changes are saved and reflected in storefront |
| SHOP-002 | As a shop admin, I want to enable/disable payment methods | Given available payment methods, When I toggle them, Then customers see only enabled methods at checkout |
| SHOP-003 | As a shop admin, I want to configure notification settings | Given notification settings, When I enable/disable email/SMS, Then notifications are sent accordingly |
| SHOP-004 | As a shop admin, I want to set business hours | Given business hours settings, When I configure them, Then they are displayed on the storefront |

**Frontend Pages to Create:**
```
shop-admin/src/views/settings/
  general/index.vue          # Basic shop info
  payment/index.vue          # Payment methods
  shipping/index.vue         # Shipping configuration
  notifications/index.vue    # Notification settings
```

---

## Implementation Priorities Summary

### P0 - Critical (v1.1)
1. Dashboard Real Statistics API
2. Order Cancel API
3. Payment Reminder API

**Estimated Effort:** 2 sprints (4 weeks)
**Risk:** Low (well-defined requirements)

### P1 - High Priority (v1.2)
1. Role Management API
2. Warehouse Management Page

**Estimated Effort:** 2 sprints (4 weeks)
**Risk:** Medium (new frontend page for warehouse)

### P2 - Medium Priority (v1.3)
1. Category Advanced Features Integration
2. Brand Advanced Features Integration

**Estimated Effort:** 2 sprints (4 weeks)
**Risk:** Low (APIs exist, frontend only)

### P3 - Low Priority (v1.4)
1. Shop/Tenant Settings APIs
2. Review Product Stats Integration

**Estimated Effort:** 2 sprints (4 weeks)
**Risk:** Medium (new APIs and pages)

---

## Technical Debt Items

### Backend
1. Add caching layer for dashboard statistics (Redis)
2. Implement event publishing for order cancellation (MQ)
3. Add audit logging for warehouse changes
4. Create database indexes for new query patterns

### Frontend
1. Extract common table components for reusability
2. Implement optimistic updates for toggle operations
3. Add error boundary components
4. Create shared filter components

---

## Success Metrics

### Business Metrics
| Metric | Current | Target | Timeline |
|--------|---------|--------|----------|
| Order cancellation handling time | 24hrs | 5min | v1.1 |
| Payment reminder conversion | N/A | +15% | v1.1 |
| Warehouse operation efficiency | Manual | 25% faster | v1.2 |
| Category management time | 2hrs/week | 30min/week | v1.3 |

### Technical Metrics
| Metric | Target |
|--------|--------|
| API response time (p95) | < 200ms |
| Frontend bundle size | < 500KB |
| Test coverage | > 80% |
| Zero critical bugs in production | 100% |

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Dashboard API performance under load | Medium | High | Implement caching, consider materialized views |
| Order cancellation refund failures | Medium | High | Implement retry mechanism, manual fallback |
| Multi-market visibility complexity | Low | Medium | Clear documentation, phased rollout |
| Role permission conflicts | Low | High | Permission audit before deployment |

---

## Appendix A: Error Code Allocations

### New Error Codes for v1.1
```
Dashboard Module (200xxx):
  200001 - Dashboard data unavailable
  200002 - Statistics calculation failed

Order Module (40xxx):
  40010 - Order cannot be cancelled
  40011 - Cancellation reason required
  40012 - Order is not pending payment
  40013 - Reminder limit exceeded
  40014 - Reminder too frequent
```

### New Error Codes for v1.2
```
Role Module (100xxx):
  100001 - Role not found
  100002 - Role already exists
  100003 - Cannot delete system role

Warehouse Module (inherited from Inventory 170xxx):
  170001 - Warehouse not found
  170002 - Warehouse has inventory
  170003 - Cannot delete default warehouse
```

---

## Appendix B: File Structure Reference

### Backend Files to Create
```
admin/desc/
  dashboard.api
  role.api
  shop.api

admin/internal/logic/
  dashboard/
    stats_logic.go
    sales_trend_logic.go
    order_status_logic.go
    hot_products_logic.go
    activities_logic.go
  roles/
    list_roles_logic.go
  shop/
    get_settings_logic.go
    update_settings_logic.go
```

### Frontend Files to Create/Modify
```
shop-admin/src/api/
  dashboard.ts (new)
  role.ts (new)
  shop.ts (new)

shop-admin/src/views/
  dashboard/index.vue (modify)
  inventory/warehouses/ (new directory)
  settings/ (expand)
  categories/index.vue (modify)
  brands/index.vue (modify)
```

---

## Appendix C: API Integration Status

### Fully Integrated APIs
- Products API
- Orders API (basic operations)
- Users API (basic operations)
- Fulfillment API
- Reviews API

### Partially Integrated APIs
- Categories API (missing advanced features)
- Brands API (missing advanced features)
- Inventory API (missing warehouse management page)
- Users API (missing enhanced stats usage)

### Not Integrated (Need Implementation)
- Dashboard statistics (backend missing)
- Role management (backend missing)
- Shop settings (backend missing)

---

*Document Version: 1.0*
*Last Updated: 2026-03-25*
*Author: Product Management Team*