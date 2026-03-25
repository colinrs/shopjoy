# ShopJoy Admin System Feature Completion PRD

## Product Overview

### Document Information
- **Document Title**: ShopJoy Admin System Feature Completion PRD
- **Version**: 1.0
- **Created Date**: 2026-03-25
- **Last Updated**: 2026-03-25

### Product Summary

ShopJoy is a multi-tenant, multi-market e-commerce admin system built with Go (go-zero framework) backend and Vue 3 + TypeScript + Element Plus frontend. This PRD addresses feature gaps identified between the backend API layer and frontend implementation, aiming to achieve full feature parity and operational readiness.

### Current State Analysis

#### Backend APIs Status
- Total endpoints: ~149 across 15 modules
- All core modules well-implemented (Products, Orders, Users, Inventory, Fulfillment, etc.)
- Dashboard APIs fully defined but not connected to frontend
- Role management APIs exist but response format mismatch with frontend expectations

#### Frontend Pages Status
- 25 page modules implemented
- 14 API service files
- Dashboard uses hardcoded mock data instead of real APIs
- Inventory page includes warehouse management in a tab (not standalone page)

### Problem Statement

1. **Dashboard Mock Data**: Dashboard page displays hardcoded mock data instead of connecting to existing backend APIs, making real-time monitoring impossible.

2. **Role Assignment Flow Broken**: Frontend `RoleAssignDialog` expects `GET /api/v1/roles` to return a simple role list for assignment, but backend returns paginated response with different structure.

3. **Feature Underutilization**: Backend provides advanced features (market visibility, product count, batch operations) for categories and brands, but frontend UI does not expose these capabilities.

4. **Incomplete Order Management**: Order cancellation API exists but frontend does not provide UI to trigger it.

## Goals

### Business Goals
- Achieve 100% API-to-Frontend connectivity for core admin operations
- Enable real-time dashboard monitoring for business decision-making
- Complete role-based access control implementation
- Provide full order lifecycle management capabilities

### User Goals
- Admin users can view real-time business statistics on dashboard
- Admin managers can assign roles to admin users seamlessly
- Operations staff can cancel orders when necessary
- Merchandisers can manage category/brand market visibility

### Non-goals
- Mobile responsive optimization (separate initiative)
- New feature development beyond completion of existing APIs
- Performance optimization for high-traffic scenarios
- Multi-language support implementation

## User Personas

### Key User Types

#### Admin Manager
- **Role**: System administrator with full access
- **Responsibilities**: User management, role assignment, system configuration
- **Pain Points**: Cannot assign roles due to API response format mismatch
- **Access Level**: Full system access

#### Operations Staff
- **Role**: Daily operations management
- **Responsibilities**: Order processing, inventory management, shipment handling
- **Pain Points**: Dashboard shows fake data, cannot see real metrics
- **Access Level**: Order, Inventory, Fulfillment modules

#### Merchandiser
- **Role**: Product and catalog management
- **Responsibilities**: Product management, category/brand configuration
- **Pain Points**: Advanced category/brand features not accessible in UI
- **Access Level**: Products, Categories, Brands modules

### Role-based Access

| Role | Dashboard | Users | Orders | Products | Inventory | Settings |
|------|-----------|-------|--------|----------|-----------|----------|
| Admin Manager | Full | Full | Full | Full | Full | Full |
| Operations Staff | Read | - | Full | Read | Full | - |
| Merchandiser | Read | - | Read | Full | Read | - |
| Customer Service | Read | Read | Read | Read | - | - |

## Functional Requirements

### Priority 1: Critical (P1)

#### FR-001: Dashboard Real-time Data Integration
**Priority**: P1 - Critical
**Description**: Connect dashboard page to backend APIs for real-time statistics display.

**Current State**: Dashboard uses hardcoded mock data:
```typescript
const stats = ref({
  todayOrders: 128,
  todaySales: 15860,
  totalProducts: 486,
  totalUsers: 3256
})
```

**Required Changes**:
- Replace mock data with API calls to existing dashboard endpoints
- Implement proper error handling and loading states
- Add data refresh mechanism

**Backend APIs Available**:
- `GET /api/v1/dashboard/overview` - Overview statistics
- `GET /api/v1/dashboard/sales-trend` - Sales trend data
- `GET /api/v1/dashboard/order-status` - Order status distribution
- `GET /api/v1/dashboard/top-products` - Top selling products
- `GET /api/v1/dashboard/pending-orders` - Pending orders
- `GET /api/v1/dashboard/activities` - Recent activities
- `GET /api/v1/dashboard` - Aggregated data (all-in-one)

#### FR-002: Role Assignment API Response Format
**Priority**: P1 - Critical
**Description**: Fix role list API to return format expected by frontend for role assignment.

**Current Issue**:
- Frontend expects: `AdminRole[]` with simple `{ id, name, code }` structure
- Backend returns: Paginated `ListRolesResponse` with `{ list, total, page, page_size }`

**Frontend Code (RoleAssignDialog.vue)**:
```typescript
export function getAvailableRoles() {
  return request<AdminRole[]>({
    url: '/api/v1/roles',
    method: 'get'
  })
}
```

**Solution Options**:
1. Add new endpoint `GET /api/v1/roles/available` returning simple list
2. Modify frontend to handle paginated response and extract list
3. Add query parameter `simple=true` to existing endpoint

**Recommended**: Option 2 - Modify frontend to handle paginated response (minimal backend change)

### Priority 2: Important (P2)

#### FR-003: Order Cancellation UI
**Priority**: P2 - Important
**Description**: Add order cancellation functionality to order management interface.

**Backend API Available**:
```
PUT /api/v1/orders/:id/cancel
Request: { reason: string }
Response: { order_id, order_no, status, status_text, cancelled_at, reason }
```

**Error Codes Already Defined**:
- `ErrOrderCannotCancel` (40014): Order cannot be cancelled in current status
- `ErrOrderCancelReasonRequired` (40015): Cancel reason is required

**Required Frontend Changes**:
- Add "Cancel Order" button in order detail view
- Implement cancellation dialog with reason input
- Handle error cases with proper user feedback

#### FR-004: Payment Reminder Feature
**Priority**: P2 - Important
**Description**: Allow admin to send payment reminder notifications for unpaid orders.

**Backend API Available**:
```
POST /api/v1/orders/:id/remind-payment
Response: { order_id, order_no, reminded_at, message }
```

**Error Codes**:
- `ErrPaymentReminderSent` (40016): Payment reminder already sent recently
- `ErrOrderAlreadyPaidForRemind` (40017): Order already paid

**Required Frontend Changes**:
- Add "Send Reminder" button for unpaid orders
- Display last reminder time
- Implement rate limiting feedback

#### FR-005: Shop Settings Page
**Priority**: P2 - Important
**Description**: Implement shop settings management page with full backend integration.

**Backend APIs Available**:
- `GET /api/v1/shop` - Get shop settings
- `PUT /api/v1/shop` - Update shop settings
- `GET /api/v1/shop/business-hours` - Business hours
- `PUT /api/v1/shop/business-hours` - Update hours
- `GET /api/v1/shop/notifications` - Notification settings
- `PUT /api/v1/shop/notifications` - Update notifications
- `GET /api/v1/shop/payment` - Payment settings
- `PUT /api/v1/shop/payment` - Update payment
- `GET /api/v1/shop/shipping` - Shipping settings
- `PUT /api/v1/shop/shipping` - Update shipping

### Priority 3: Enhancement (P3)

#### FR-006: Category Advanced Features UI
**Priority**: P3 - Enhancement
**Description**: Expose advanced category management features in the UI.

**Backend APIs Available (not used in frontend)**:
- `PUT /api/v1/categories/sort` - Batch update category sort order
- `PUT /api/v1/categories/:id/move` - Move category to new parent
- `GET /api/v1/categories/:id/product-count` - Get product count
- `PUT /api/v1/categories/:id/market-visibility` - Set market visibility
- `GET /api/v1/categories/:id/market-visibility` - Get market visibility

**Frontend API Already Defined**: `category.ts` has all functions implemented

**Required Changes**:
- Add drag-and-drop sorting UI
- Add "Move Category" dialog
- Add market visibility configuration
- Display product count in category list

#### FR-007: Brand Advanced Features UI
**Priority**: P3 - Enhancement
**Description**: Expose advanced brand management features in the UI.

**Backend APIs Available (not used in frontend)**:
- `PUT /api/v1/brands/:id/toggle-page` - Toggle brand page
- `GET /api/v1/brands/:id/product-count` - Get product count
- `PUT /api/v1/brands/:id/market-visibility` - Set market visibility
- `GET /api/v1/brands/:id/market-visibility` - Get market visibility

**Frontend API Already Defined**: `brand.ts` has all functions implemented

**Required Changes**:
- Add brand page toggle switch
- Display product count in brand list
- Add market visibility configuration

## User Experience

### Entry Points

#### Dashboard
- **Route**: `/dashboard`
- **Access**: All authenticated users
- **First Impression**: Real-time statistics cards with key metrics
- **Key Actions**: View sales trend, check pending orders, review activities

#### Admin User Management
- **Route**: `/admin-users`
- **Access**: Admin Manager only
- **Entry**: User list with role assignment capability
- **Key Actions**: Create user, assign roles, enable/disable users

#### Order Management
- **Route**: `/orders`
- **Access**: Operations Staff, Admin Manager
- **Entry**: Order list with status filters
- **Key Actions**: View details, ship order, cancel order, send reminders

### Core Experience

#### Dashboard Real-time Data Flow
1. User navigates to dashboard
2. Page makes parallel API calls for:
   - Overview statistics
   - Sales trend (default: week)
   - Order status distribution
   - Pending orders (limit: 5)
   - Top products (limit: 5)
   - Recent activities (limit: 10)
3. Charts render with real data
4. User can switch time range for sales trend
5. Data refreshes automatically every 5 minutes (optional)

#### Role Assignment Flow
1. Admin Manager views admin user list
2. Clicks "Assign Roles" button for a user
3. Dialog fetches available roles via `GET /api/v1/roles?page_size=100&status=1`
4. Dialog displays roles with checkboxes
5. User selects/deselects roles
6. Clicks "Confirm" to call `PUT /api/v1/admin-users/:id/roles`
7. Success message displayed, list refreshed

#### Order Cancellation Flow
1. Operations staff views order detail
2. Order is in "pending_payment" status
3. Staff clicks "Cancel Order" button
4. Dialog appears requiring cancellation reason
5. Staff enters reason (min 10 characters)
6. Clicks "Confirm" to call `PUT /api/v1/orders/:id/cancel`
7. Order status changes to "cancelled"
8. Audit log records the action

### UI/UX Highlights

#### Dashboard Design
- **Statistics Cards**: Gradient backgrounds, clear metrics, trend indicators
- **Sales Trend Chart**: Interactive ECharts with time range selector
- **Order Status Pie Chart**: Color-coded segments with percentages
- **Pending Orders Table**: Quick action buttons, status tags
- **Recent Activities Timeline**: Chronological activity feed

#### Error Handling
- **API Errors**: Toast notifications with specific error messages
- **Network Errors**: Retry button with exponential backoff
- **Loading States**: Skeleton screens for better perceived performance
- **Empty States**: Helpful illustrations and guidance text

## Narrative

As an operations manager at a growing e-commerce company using ShopJoy, I need to monitor business performance through a dashboard that shows real data, not placeholder numbers. When I log in each morning, I want to see yesterday's sales, pending orders that need attention, and inventory alerts - all based on actual system data.

Currently, my dashboard shows mock data that doesn't reflect reality. I also need to assign specific roles to new team members, but the role assignment dialog fails silently because the API returns data in an unexpected format. When customers call to cancel orders, I have no way to process the cancellation through the admin panel, forcing me to work around the system.

After this feature completion, I will be able to see accurate real-time metrics, assign roles without errors, and handle order cancellations directly in the system - making the admin panel truly operational for daily business management.

## Success Metrics

### User-centric Metrics
- **Dashboard Engagement**: Time spent on dashboard (target: 3+ minutes per session)
- **Task Completion Rate**: Role assignment success rate (target: 95%+)
- **Error Reduction**: API-related errors in role assignment (target: <1%)

### Business Metrics
- **Operational Efficiency**: Time to process order cancellations (target: <2 minutes)
- **Feature Adoption**: Usage of advanced category/brand features (target: 50% of merchants)
- **Support Tickets**: Reduction in "how do I..." tickets (target: 30% decrease)

### Technical Metrics
- **API Coverage**: Frontend-to-backend API utilization (target: 90%+)
- **Error Rate**: API error rate for dashboard endpoints (target: <0.5%)
- **Response Time**: Dashboard initial load time (target: <2 seconds)

## Technical Considerations

### Integration Points

#### Backend API Endpoints to Connect

**Dashboard Module** (`/api/v1/dashboard/*`):
- All 7 endpoints defined in `admin/desc/dashboard.api`
- Need to create frontend API service file `dashboard.ts`
- Implement proper TypeScript interfaces

**Role Module** (`/api/v1/roles`):
- Existing endpoint returns paginated response
- Frontend should handle pagination or request all active roles
- Consider caching role list for better UX

**Order Module** (`/api/v1/orders/*`):
- Cancellation endpoint ready with proper error handling
- Payment reminder endpoint ready with rate limiting
- Both return appropriate error codes

### Data Storage/Privacy

**Dashboard Data**:
- Read-only aggregation queries
- No sensitive data stored locally
- Cache with 5-minute TTL acceptable

**Role Data**:
- Role assignments stored in `admin_user_roles` table
- Audit trail required for role changes
- GDPR: Role data is operational, not PII

### Scalability/Performance

**Dashboard Performance**:
- Aggregation queries should use proper indexes
- Consider caching dashboard stats for 1-5 minutes
- Lazy load non-critical widgets

**Role List Performance**:
- Typically small dataset (<100 roles)
- Can be cached in frontend for session duration
- No pagination needed for assignment dialog

### Potential Challenges

1. **Dashboard API Implementation**: Backend endpoints may need business logic implementation (aggregating data from multiple tables)

2. **Real-time Updates**: Dashboard should reflect near real-time data, consider WebSocket for live updates in future

3. **Role Response Format**: Changing backend response format could affect other consumers; frontend adaptation is safer

4. **Order Cancellation Business Rules**: Ensure cancellation rules are properly enforced (timing, payment status, fulfillment status)

## Milestones & Sequencing

### Project Estimate
- **Total Effort**: 3-4 weeks
- **Team Size**: 2-3 developers (1 backend, 1-2 frontend)

### Phase 1: Critical Fixes (Week 1)

**Objective**: Fix blocking issues for core workflows

**Deliverables**:
1. Dashboard API service integration (Frontend)
2. Dashboard page connected to real APIs (Frontend)
3. Role assignment API response handling fix (Frontend)

**Milestones**:
- Day 1-2: Create dashboard.ts API service
- Day 3-4: Update dashboard page to use real APIs
- Day 5: Fix role assignment dialog

### Phase 2: Order Management Completion (Week 2)

**Objective**: Complete order lifecycle management

**Deliverables**:
1. Order cancellation UI implementation
2. Payment reminder UI implementation
3. Error handling and validation

**Milestones**:
- Day 1-2: Add cancel order dialog and API integration
- Day 3-4: Add payment reminder feature
- Day 5: Testing and edge case handling

### Phase 3: Shop Settings & Advanced Features (Week 3-4)

**Objective**: Complete remaining feature gaps

**Deliverables**:
1. Shop settings page implementation
2. Category advanced features UI
3. Brand advanced features UI

**Milestones**:
- Week 3, Day 1-3: Shop settings page
- Week 3, Day 4-5: Category advanced features
- Week 4, Day 1-2: Brand advanced features
- Week 4, Day 3-5: Testing, bug fixes, documentation

## User Stories

### US-001: Dashboard Real-time Statistics
**ID**: US-001
**Title**: View real-time business statistics on dashboard
**Priority**: P1 - Critical

**As an** admin user
**I want to** see real-time business statistics on the dashboard
**So that** I can make informed decisions about business operations

**Acceptance Criteria**:
1. Dashboard displays actual today's orders count from backend API
2. Dashboard displays actual today's sales amount with correct currency
3. Dashboard shows growth percentage compared to yesterday
4. Total products count reflects database count
5. Total users count reflects registered users
6. New users today count is accurate
7. All statistics load within 2 seconds
8. Loading skeleton shown while data fetches
9. Error state shown if API fails with retry option

---

### US-002: Dashboard Sales Trend Chart
**ID**: US-002
**Title**: View sales trend visualization
**Priority**: P1 - Critical

**As an** admin user
**I want to** view sales trend over time with different time ranges
**So that** I can identify business patterns and trends

**Acceptance Criteria**:
1. Sales trend chart displays data from backend API `GET /api/v1/dashboard/sales-trend`
2. Default time range is "week" (7 days)
3. User can switch between "week", "month", and "year" views
4. Chart shows both sales amount and order count
5. Chart is interactive with hover tooltips
6. Currency is displayed based on shop settings
7. Empty state shown when no data available for period

---

### US-003: Dashboard Order Status Distribution
**ID**: US-003
**Title**: View order status distribution
**Priority**: P1 - Critical

**As an** admin user
**I want to** see the distribution of orders by status
**So that** I can understand the current order pipeline

**Acceptance Criteria**:
1. Pie chart displays order distribution from `GET /api/v1/dashboard/order-status`
2. Each status segment shows count and percentage
3. Color coding matches order status conventions
4. Total order count displayed
5. Legend shows all status types
6. Click on segment shows orders in that status (future enhancement)

---

### US-004: Dashboard Pending Orders
**ID**: US-004
**Title**: View pending orders requiring attention
**Priority**: P1 - Critical

**As an** operations staff
**I want to** see a list of pending orders on the dashboard
**So that** I can quickly identify orders that need action

**Acceptance Criteria**:
1. Table displays top 5 pending orders from `GET /api/v1/dashboard/pending-orders`
2. Shows order number, amount, status, and creation time
3. "Process" button links to order detail
4. Total pending count displayed
5. "View All" link navigates to orders page with pending filter
6. Refresh button to reload pending orders

---

### US-005: Dashboard Top Products
**ID**: US-005
**Title**: View top selling products
**Priority**: P1 - Critical

**As an** admin user
**I want to** see the top selling products on the dashboard
**So that** I can identify best performers and manage inventory

**Acceptance Criteria**:
1. Table displays top 5 products from `GET /api/v1/dashboard/top-products`
2. Shows product image, name, sales count, and revenue
3. Ranking number displayed for each product
4. Top 3 products highlighted with special styling
5. Time period selector (week/month/all time)
6. "View All" link to products page

---

### US-006: Dashboard Recent Activities
**ID**: US-006
**Title**: View recent system activities
**Priority**: P1 - Critical

**As an** admin user
**I want to** see recent system activities on the dashboard
**So that** I can stay informed about important events

**Acceptance Criteria**:
1. Timeline displays recent activities from `GET /api/v1/dashboard/activities`
2. Activities sorted by time (most recent first)
3. Activity type indicated by icon and color
4. Shows activity content, time, and operator (if applicable)
5. Different activity types: order_created, payment_received, low_stock, etc.
6. Load more button for additional activities

---

### US-007: Role Assignment Dialog
**ID**: US-007
**Title**: Assign roles to admin users
**Priority**: P1 - Critical

**As an** admin manager
**I want to** assign roles to admin users
**So that** users have appropriate access permissions

**Acceptance Criteria**:
1. Role assignment dialog fetches roles from `GET /api/v1/roles?page_size=100&status=1`
2. Handles paginated response correctly (extracts `list` array)
3. Shows role name and code for each available role
4. Current user's roles are pre-selected
5. Can select/deselect multiple roles
6. "Confirm" button calls `PUT /api/v1/admin-users/:id/roles`
7. Success message shown after assignment
8. Dialog closes and user list refreshes
9. Error message shown if assignment fails
10. System roles cannot be modified (shown with indicator)

---

### US-008: Order Cancellation
**ID**: US-008
**Title**: Cancel an order
**Priority**: P2 - Important

**As an** operations staff
**I want to** cancel an order with a reason
**So that** orders that cannot be fulfilled are properly closed

**Acceptance Criteria**:
1. "Cancel Order" button visible for orders in cancellable status
2. Button disabled for non-cancellable orders (already shipped, delivered, etc.)
3. Clicking button opens cancellation dialog
4. Dialog requires cancellation reason (min 10 characters)
5. Reason field has character counter
6. Cancel button calls `PUT /api/v1/orders/:id/cancel`
7. Success: Order status changes to "cancelled"
8. Success: Audit log records who cancelled and why
9. Error: Specific error message if cancellation not allowed
10. Error: Validation error if reason too short
11. Order detail page refreshes after cancellation

---

### US-009: Payment Reminder
**ID**: US-009
**Title**: Send payment reminder for unpaid orders
**Priority**: P2 - Important

**As an** operations staff
**I want to** send a payment reminder notification to customers
**So that** customers are prompted to complete their purchase

**Acceptance Criteria**:
1. "Send Reminder" button visible for unpaid orders only
2. Button shows last reminder time if previously sent
3. Rate limit: Only one reminder per hour allowed
4. Clicking button calls `POST /api/v1/orders/:id/remind-payment`
5. Success: Confirmation message with reminder timestamp
6. Success: Last reminder time updated on UI
7. Error: "Already sent recently" message with wait time
8. Error: "Order already paid" message if status changed
9. Button disabled during cooldown period

---

### US-010: Shop Basic Settings
**ID**: US-010
**Title**: Manage shop basic information
**Priority**: P2 - Important

**As a** shop admin
**I want to** update my shop's basic information
**So that** customers see correct branding and contact info

**Acceptance Criteria**:
1. Settings page displays current shop info from `GET /api/v1/shop`
2. Editable fields: name, logo, description, contact info
3. Domain settings: view system domain, configure custom domain
4. Branding settings: primary color, secondary color, favicon
5. Business settings: default currency, language, timezone
6. Save button calls `PUT /api/v1/shop`
7. Validation errors shown inline
8. Success message after save
9. Changes reflected immediately in UI

---

### US-011: Shop Business Hours
**ID**: US-011
**Title**: Configure business hours
**Priority**: P2 - Important

**As a** shop admin
**I want to** set my business hours
**So that** customers know when support is available

**Acceptance Criteria**:
1. Business hours section shows 7 days (Sunday-Saturday)
2. Each day has open time, close time, and closed toggle
3. Closed toggle disables time inputs for that day
4. Save button calls `PUT /api/v1/shop/business-hours`
5. Validation: open time must be before close time
6. Success message after save

---

### US-012: Shop Notification Settings
**ID**: US-012
**Title**: Configure notification preferences
**Priority**: P2 - Important

**As a** shop admin
**I want to** configure which notifications I receive
**So that** I'm alerted about important events

**Acceptance Criteria**:
1. Toggle switches for each notification type:
   - Order created
   - Order paid
   - Order shipped
   - Order cancelled
   - Low stock alert
   - Refund requested
   - New review
2. Low stock threshold number input
3. Save button calls `PUT /api/v1/shop/notifications`
4. Success message after save
5. Settings persist across sessions

---

### US-013: Category Drag-and-Drop Sort
**ID**: US-013
**Title**: Reorder categories with drag and drop
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** reorder categories by dragging
**So that** I can control the display order in the storefront

**Acceptance Criteria**:
1. Category list supports drag-and-drop reordering
2. Drag handle visible on hover
3. Visual indicator during drag operation
4. Drop target highlighted
5. On drop, calls `PUT /api/v1/categories/sort`
6. Success: Order persists after page refresh
7. Error: Revert to previous order with error message
8. Only categories at same level can be reordered

---

### US-014: Category Move
**ID**: US-014
**Title**: Move category to different parent
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** move a category to a different parent
**So that** I can reorganize the category hierarchy

**Acceptance Criteria**:
1. "Move" action in category row menu
2. Move dialog shows category tree selector
3. Cannot move to self or descendants
4. Cannot exceed max depth (3 levels)
5. Confirm button calls `PUT /api/v1/categories/:id/move`
6. Success: Category appears under new parent
7. Error: Clear message if move not allowed
8. Tree view refreshes after move

---

### US-015: Category Product Count
**ID**: US-015
**Title**: View product count per category
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** see how many products are in each category
**So that** I can identify empty or overpopulated categories

**Acceptance Criteria**:
1. Product count column in category list
2. Count fetched from `GET /api/v1/categories/:id/product-count`
3. Shows "0" for empty categories
4. Includes products from subcategories (optional toggle)
5. Clickable count links to products filtered by category
6. Count updates when products are added/removed

---

### US-016: Category Market Visibility
**ID**: US-016
**Title**: Configure category visibility per market
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** control which markets can see each category
**So that** I can customize catalog per market

**Acceptance Criteria**:
1. "Market Visibility" action in category menu
2. Dialog shows list of active markets
3. Toggle switch for each market visibility
4. Save button calls `PUT /api/v1/categories/:id/market-visibility`
5. Bulk toggle for all markets
6. Success message after save
7. Visibility status shown in category list

---

### US-017: Brand Page Toggle
**ID**: US-017
**Title**: Enable/disable brand showcase page
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** toggle whether a brand has a dedicated page
**So that** important brands get special showcase treatment

**Acceptance Criteria**:
1. "Enable Page" toggle in brand list
2. Toggle calls `PUT /api/v1/brands/:id/toggle-page`
3. Toggle state persists
4. Visual indicator for brands with pages enabled
5. Success toast after toggle
6. Error handling for failed toggle

---

### US-018: Brand Product Count
**ID**: US-018
**Title**: View product count per brand
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** see how many products are associated with each brand
**So that** I can understand brand coverage

**Acceptance Criteria**:
1. Product count column in brand list
2. Count fetched from `GET /api/v1/brands/:id/product-count`
3. Shows "0" for brands with no products
4. Clickable count links to products filtered by brand
5. Count updates when products are reassigned

---

### US-019: Brand Market Visibility
**ID**: US-019
**Title**: Configure brand visibility per market
**Priority**: P3 - Enhancement

**As a** merchandiser
**I want to** control which markets can see each brand
**So that** I can customize brand availability per region

**Acceptance Criteria**:
1. "Market Visibility" action in brand menu
2. Dialog shows list of active markets
3. Toggle switch for each market visibility
4. Save button calls `PUT /api/v1/brands/:id/market-visibility`
5. Bulk toggle for all markets
6. Success message after save
7. Visibility status shown in brand list

---

### US-020: Dashboard Auto-refresh
**ID**: US-020
**Title**: Dashboard data auto-refresh
**Priority**: P3 - Enhancement

**As an** admin user
**I want the** dashboard to automatically refresh data
**So that** I always see current information without manual refresh

**Acceptance Criteria**:
1. Dashboard refreshes data every 5 minutes
2. Refresh interval configurable in user preferences
3. Visual indicator showing last refresh time
4. Manual refresh button available
5. Auto-refresh pauses when tab is hidden (optional optimization)
6. No page flicker during refresh (smooth update)

---

### US-021: Secure Access - Dashboard
**ID**: US-021
**Title**: Dashboard access requires authentication
**Priority**: P1 - Critical

**As a** system administrator
**I want to** ensure dashboard data is only accessible to authenticated users
**So that** sensitive business data is protected

**Acceptance Criteria**:
1. All dashboard endpoints require AuthMiddleware
2. Unauthenticated requests return 401 Unauthorized
3. User must have at least one role to access
4. Tenant isolation enforced (users see only their tenant's data)
5. API tokens have appropriate scope for dashboard access

---

### US-022: Secure Access - Role Management
**ID**: US-022
**Title**: Role management access restricted to admins
**Priority**: P1 - Critical

**As a** system administrator
**I want to** ensure only admin managers can manage roles
**So that** access control integrity is maintained

**Acceptance Criteria**:
1. Role list endpoint requires admin role
2. Role creation requires admin role
3. Role update requires admin role
4. Role deletion requires admin role
5. System roles cannot be deleted
6. Roles in use cannot be deleted
7. Audit log records all role modifications

---

### US-023: Error Handling - Dashboard
**ID**: US-023
**Title**: Graceful error handling on dashboard
**Priority**: P1 - Critical

**As an** admin user
**I want to** see helpful error messages when dashboard fails to load
**So that** I understand what went wrong and can take action

**Acceptance Criteria**:
1. API errors show specific error messages (not generic)
2. Network errors show "Connection lost" with retry button
3. Timeout errors show "Request timed out" with retry button
4. Partial failures show which widgets failed
5. Retry button available for failed sections
6. Error state visually distinct from loading state
7. Errors logged to monitoring system

---

### US-024: Audit Log - Order Operations
**ID**: US-024
**Title**: Audit trail for order operations
**Priority**: P2 - Important

**As a** system administrator
**I want to** track who performed order cancellations and sent reminders
**So that** there is accountability for order operations

**Acceptance Criteria**:
1. Order cancellation records operator ID and reason
2. Payment reminder records operator ID and timestamp
3. Audit entries include timestamp, operator, action, and details
4. Audit log viewable in order detail page
5. Audit entries cannot be modified or deleted
6. Audit data retained per compliance requirements

---

### US-025: Dashboard Loading States
**ID**: US-025
**Title**: Skeleton loading for dashboard widgets
**Priority**: P2 - Important

**As an** admin user
**I want to** see skeleton placeholders while dashboard loads
**So that** the interface feels responsive and polished

**Acceptance Criteria**:
1. Statistics cards show skeleton while loading
2. Charts show loading placeholder
3. Tables show skeleton rows while loading
4. Skeleton animations are smooth (not jarring)
5. Content fades in when loaded
6. No layout shift during transition