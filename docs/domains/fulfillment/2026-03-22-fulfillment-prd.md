# Fulfillment Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Fulfillment Module PRD |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-22 |
| Author | Product Team |
| Module Location | admin/internal/domain/fulfillment/ |

---

## Product Overview

### Summary

The Fulfillment module is an independent bounded context within the ShopJoy e-commerce SaaS platform that handles the complete order fulfillment lifecycle, from shipment creation and logistics tracking to refund processing and after-sales statistics. This module bridges the gap between order placement and customer delivery, ensuring merchants can efficiently manage shipments, track logistics status, and process refund requests in a multi-tenant environment.

The MVP release focuses on core fulfillment operations: shipment creation for pending orders, logistics information entry and tracking, and basic refund application processing. The system supports both single and split shipments, enabling merchants to handle complex scenarios like multi-warehouse fulfillment and partial order shipping.

### Problem Statement

Merchants on the ShopJoy platform currently face several challenges in the fulfillment process:

- **Manual shipment management**: No centralized system to create shipments, enter tracking numbers, and update order status
- **Lack of logistics visibility**: Customers cannot track their packages, leading to increased customer service inquiries
- **Complex refund processing**: No structured workflow for buyers to request refunds and merchants to review and process them
- **No fulfillment analytics**: Merchants lack visibility into refund rates, return reasons, and problematic products
- **Multi-warehouse complexity**: Large items and cross-border shipments require split shipment capabilities

### Solution Overview

Implement a domain-driven Fulfillment bounded context that provides:

1. Shipment management with support for single and split shipments
2. Logistics information entry with carrier and tracking number management
3. Order fulfillment status tracking (pending, partial_shipped, shipped, delivered)
4. Refund application and processing workflow with approval/rejection
5. After-sales statistics and analytics dashboard
6. Integration with inventory management for stock updates

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | Reduce order processing time | Average time from payment to shipment | <24 hours |
| BG-002 | Decrease customer service inquiries about shipping | Support ticket reduction | -30% within 3 months |
| BG-003 | Streamline refund processing | Average refund resolution time | <48 hours |
| BG-004 | Improve customer satisfaction | CSAT score for delivery experience | >4.5/5.0 |
| BG-005 | Enable data-driven supply chain optimization | Refund rate monitoring accuracy | 100% |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | Quickly ship pending orders with logistics info | Merchant Admin |
| UG-002 | View and manage all shipments in one place | Merchant Admin |
| UG-003 | Process refund requests efficiently | Merchant Admin |
| UG-004 | Track my order status and delivery | Buyer |
| UG-005 | Apply for refund when needed | Buyer |
| UG-006 | View after-sales statistics and insights | Merchant Admin |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | Third-party logistics API integration | Phase 2 feature (requires carrier partnerships) |
| NG-002 | Return merchandise authorization (RMA) | Phase 2 feature |
| NG-003 | Automated shipping label generation | Phase 2 feature (requires carrier integration) |
| NG-004 | Multi-warehouse intelligent routing | Phase 2 feature |
| NG-005 | Cross-border customs documentation | Phase 2 feature |
| NG-006 | Partial refund (line item level) | Phase 2 feature (MVP focuses on full refund) |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | Store owner or operations manager | Create shipments, enter tracking, process refunds |
| Operations Staff | Warehouse or fulfillment staff | Ship orders, print packing lists |
| Buyer | End customer who placed the order | Track order, apply for refund |

### Persona Details

#### Merchant Admin (Primary)

- **Demographics**: Small to medium business owners, e-commerce operations managers
- **Technical Proficiency**: Moderate
- **Goals**: Efficient order fulfillment, reduce customer complaints, maintain financial control
- **Pain Points**: Manual tracking entry, unclear refund policies, lack of fulfillment visibility
- **Frequency**: Daily shipment processing, frequent refund handling

#### Operations Staff (Secondary)

- **Demographics**: Warehouse workers, packing staff, shipping clerks
- **Technical Proficiency**: Basic to moderate
- **Goals**: Quickly process orders, minimize shipping errors
- **Pain Points**: Complex multi-item orders, unclear shipping instructions
- **Frequency**: Hourly during operational hours

#### Buyer (Tertiary)

- **Demographics**: Online shoppers of all ages
- **Technical Proficiency**: Varies widely
- **Goals**: Know when package arrives, easy refund process if needed
- **Pain Points**: Unclear delivery status, complicated refund process
- **Frequency**: Per order (tracking) or when issues arise (refund)

### Role-Based Access

| Role | Permissions |
|------|-------------|
| Tenant Admin | Full access to shipments and refunds |
| Tenant Operations Manager | Create/edit shipments, process refunds |
| Tenant Operations Staff | View shipments, update shipment status |
| Buyer | View own orders and shipments, apply for refunds |

---

## Functional Requirements

### Priority 1: Order Fulfillment Status

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | View pending shipment orders | List orders with status "paid" or "pending_shipment" | P0 |
| FR-002 | Update order fulfillment status | Track fulfillment_status separately from order status | P0 |
| FR-003 | Status transitions | pending -> partial_shipped -> shipped -> delivered | P0 |
| FR-004 | Automatic status update | Fulfillment status updates based on shipment records | P0 |

### Priority 2: Shipment Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-005 | Create shipment | Create shipment record for an order | P0 |
| FR-006 | Enter logistics info | Input carrier name and tracking number | P0 |
| FR-007 | Single shipment | One shipment for entire order | P0 |
| FR-008 | Split shipment | Multiple shipments for one order (partial fulfillment) | P1 |
| FR-009 | Shipment items tracking | Track which items are in which shipment | P0 |
| FR-010 | Shipment status management | Update shipment status (pending, shipped, in_transit, delivered, failed) | P0 |
| FR-011 | Batch shipment creation | Create shipments for multiple orders at once | P1 |

### Priority 3: Logistics Entry and Tracking

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-012 | Carrier selection | Select from predefined carrier list | P0 |
| FR-013 | Custom carrier entry | Enter custom carrier name if not in list | P1 |
| FR-014 | Tracking number validation | Validate tracking number format per carrier | P1 |
| FR-015 | Manual status update | Mark shipment as delivered manually | P0 |
| FR-016 | Shipping cost tracking | Record actual shipping cost per shipment | P1 |
| FR-017 | Weight recording | Record package weight for shipping | P1 |

### Priority 4: Refund Application and Processing

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-018 | Buyer applies for refund | Buyer submits refund request with reason | P0 |
| FR-019 | Refund reason selection | Select from predefined reasons or enter custom | P0 |
| FR-020 | Evidence upload | Upload images as refund evidence | P1 |
| FR-021 | Merchant reviews refund | View pending refund requests | P0 |
| FR-022 | Approve refund | Merchant approves refund request | P0 |
| FR-023 | Reject refund | Merchant rejects refund with reason | P0 |
| FR-024 | Full refund only | MVP supports only full order refund | P0 |
| FR-025 | Refund amount calculation | System calculates refund amount from order | P0 |
| FR-026 | Refund status tracking | Track refund status through workflow | P0 |

### Priority 5: After-Sales Statistics

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-027 | Refund rate calculation | Calculate refund rate by time period | P1 |
| FR-028 | Refund reason analysis | Aggregate refund reasons | P1 |
| FR-029 | Problem product identification | Identify products with high refund rates | P1 |
| FR-030 | Delivery success rate | Track delivery success/failure rates | P1 |
| FR-031 | Carrier performance | Compare carrier delivery times | P2 |

---

## User Experience

### Entry Points

| Entry Point | Description | User Flow |
|-------------|-------------|-----------|
| Admin Orders Page | "Pending Shipment" filter and "Ship" button | Orders > Pending Shipment > Ship |
| Admin Shipments Page | Shipments list with all status | Orders > Shipments > List |
| Admin Refunds Page | Refunds list and processing | Orders > Refunds > Process |
| Shop Order Detail | Buyer tracking and refund request | My Orders > Order Detail > Track/Refund |

### Core Experience

#### Shipping an Order (Merchant Admin)

1. Navigate to Orders page
2. Filter by "Pending Shipment" status
3. Select order(s) to ship
4. Click "Ship" button
5. Enter logistics information:
   - Select carrier from dropdown
   - Enter tracking number
   - Optionally enter shipping cost and weight
6. Confirm shipment
7. Order status changes to "Shipped"
8. Shipment record is created
9. Buyer receives notification with tracking info

#### Applying for Refund (Buyer)

1. Navigate to "My Orders"
2. Select order to refund
3. Click "Apply for Refund"
4. Select refund reason from dropdown
5. Enter detailed description
6. Upload evidence images (optional)
7. Submit refund application
8. Order status changes to "Refunding"
9. Merchant reviews and processes

#### Processing Refund (Merchant Admin)

1. Navigate to Refunds page
2. View pending refund requests
3. Click on refund to view details
4. Review buyer's reason and evidence
5. Check order details and items
6. Click "Approve" or "Reject"
7. If rejecting, enter rejection reason
8. If approving, refund is processed
9. Order status updates to "Refunded"

### Advanced Features

| Feature | Description |
|---------|-------------|
| Split Shipment | Ship partial items, track remaining items to ship |
| Batch Shipping | Select multiple orders and ship together |
| Shipment History | View all shipments for an order |
| Refund History | View all refund attempts for an order |
| Analytics Dashboard | Visualize fulfillment metrics |

### UI/UX Highlights

- Clear order status indicators with color coding
- One-click shipping action from order list
- Inline tracking number entry
- Mobile-responsive admin interface
- Real-time order count by fulfillment status
- Visual timeline for order fulfillment journey
- Bulk action toolbar for batch operations

---

## Narrative

As a Merchant Admin for an online electronics store, I receive daily orders that need to be shipped. Every morning, I log into ShopJoy admin and navigate to the Orders page. I filter by "Pending Shipment" and see 15 orders waiting to be shipped.

I select all orders going to the same region, click "Batch Ship", and enter the carrier (SF Express) and the starting tracking number. The system auto-generates tracking numbers in sequence for all selected orders. I confirm, and all 15 orders are now marked as "Shipped".

Later that day, I receive a notification about a new refund request. A customer claims the wireless earbuds they received are defective. I navigate to the Refunds page, open the request, and see the customer has uploaded photos showing the issue. I check the product's refund rate and see it's higher than average. I approve the refund, and the system automatically initiates the refund process through the payment gateway. The order status updates to "Refunded", and the customer receives their refund notification.

At the end of the week, I check the After-Sales Statistics dashboard. I see our overall refund rate is 3.2%, with "Product Defective" being the most common reason. I identify two products with unusually high refund rates and decide to remove them from sale pending quality review.

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Shipment Creation Time | Average time to create a shipment | <2 minutes |
| Refund Resolution Time | Average time from application to resolution | <48 hours |
| Tracking Adoption | Percentage of shipments with tracking info | >95% |
| Buyer Self-Service | Percentage of buyers checking tracking without contacting support | >80% |

### Business Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Order Processing Time | Time from payment to shipment | <24 hours average |
| Refund Rate | Percentage of orders refunded | <5% |
| Delivery Success Rate | Percentage of successful deliveries | >98% |
| Support Ticket Reduction | Decrease in shipping-related tickets | -30% in 3 months |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API Response Time | Shipment creation endpoint | <200ms p95 |
| System Uptime | Fulfillment service availability | 99.9% |
| Error Rate | Failed shipment/refund operations | <0.1% |

---

## Technical Considerations

### Integration Points

| Integration | Description | Method |
|-------------|-------------|--------|
| Order Service | Update order status on shipment/refund | Domain Event |
| Inventory Service | Release/deduct stock on shipment/refund | Domain Event |
| Payment Service | Process refund through payment gateway | Internal API |
| Notification Service | Send shipment/refund notifications | Message Queue |
| Product Service | Fetch product info for refund analysis | Internal API |

### Data Storage and Privacy

| Aspect | Consideration |
|--------|---------------|
| Monetary Values | Use int64 (cents) for all amounts, no floating-point |
| Time Storage | All timestamps in UTC, stored as Unix BIGINT |
| Multi-Tenancy | All queries filtered by tenant_id |
| Soft Delete | Support deleted_at for data recovery |
| Tracking Numbers | Not PII, no special handling needed |
| Refund Evidence | Images stored in CDN, URLs in database |

### Scalability and Performance

| Consideration | Strategy |
|---------------|----------|
| High Order Volume | Partition shipment/refund tables by tenant_id |
| Concurrent Shipment Creation | Optimistic locking on order status |
| Large Order Items | Limit items per order, pagination for display |
| Analytics Queries | Pre-aggregate statistics, cache results |

### Potential Challenges

| Challenge | Mitigation |
|-----------|------------|
| Race Conditions on Order Status | Database-level constraints, optimistic locking |
| Partial Shipment Complexity | Clear item-level tracking, UI guidance |
| Refund Fraud | Validation rules, manual review for high-value |
| Carrier Data Inconsistency | Validate carrier codes, support custom entries |
| Timezone Confusion | Always store/display in UTC, frontend converts |

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Team Size | Description |
|-------|----------|-----------|-------------|
| Phase 1: Foundation | 1 week | 2 backend | Core entities, repositories, database schema |
| Phase 2: Shipment MVP | 1 week | 2 backend, 1 frontend | Basic shipment creation and listing |
| Phase 3: Refund MVP | 1 week | 2 backend, 1 frontend | Refund application and processing |
| Phase 4: Integration | 1 week | 2 backend | Order/inventory integration, events |
| Phase 5: Statistics | 0.5 week | 1 backend, 1 frontend | Analytics dashboard |

**Total Estimate: 4.5 weeks**

### Suggested Phases

#### Phase 1: Foundation (Week 1)

- Implement/enhance fulfillment domain entities
- Create shipment and refund repository interfaces
- Implement GORM repositories
- Create database migrations (if needed)
- Define error codes and constants

#### Phase 2: Shipment MVP (Week 2)

- Build shipment CRUD APIs
- Implement "Ship Order" use case
- Create shipment list and detail views
- Add carrier selection support
- Integrate with order status update
- Admin UI for shipment management

#### Phase 3: Refund MVP (Week 3)

- Build refund CRUD APIs
- Implement "Apply for Refund" use case (buyer side)
- Implement "Process Refund" use case (merchant side)
- Create refund list and detail views
- Add refund reason management
- Shop UI for refund application
- Admin UI for refund processing

#### Phase 4: Integration (Week 4)

- Integrate with inventory service (stock release)
- Integrate with payment service (refund execution)
- Add domain events for shipment/refund
- Implement notification triggers
- Test end-to-end flows

#### Phase 5: Statistics (Week 4.5)

- Build refund rate calculation
- Implement refund reason aggregation
- Create problem product identification
- Build analytics dashboard UI
- Add export functionality

---

## Database Schema

### shipments

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| order_id | VARCHAR(64) | NO | - | Order ID |
| shipment_no | VARCHAR(32) | NO | '' | Shipment number (auto-generated) |
| status | TINYINT | NO | 0 | 0=pending, 1=shipped, 2=in_transit, 3=delivered, 4=failed |
| carrier | VARCHAR(50) | NO | '' | Carrier name |
| carrier_code | VARCHAR(20) | NO | '' | Carrier code (SF, YT, ZT, etc.) |
| tracking_no | VARCHAR(100) | NO | '' | Tracking number |
| shipping_cost | BIGINT | NO | 0 | Shipping cost (cents) |
| shipping_currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| weight | DECIMAL(10,2) | NO | 0.00 | Weight (kg) |
| shipped_at | BIGINT | YES | NULL | Shipment timestamp |
| delivered_at | BIGINT | YES | NULL | Delivery timestamp |
| remark | VARCHAR(500) | NO | '' | Remarks |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- UNIQUE KEY `uk_shipment_no` (`tenant_id`, `shipment_no`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_order_id` (`order_id`)
- KEY `idx_tracking_no` (`tracking_no`)
- KEY `idx_status` (`status`)
- KEY `idx_deleted_at` (`deleted_at`)

### shipment_items

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| shipment_id | BIGINT | NO | - | Shipment ID |
| order_item_id | BIGINT | NO | - | Order item ID |
| product_id | BIGINT | NO | - | Product ID |
| sku_id | BIGINT | NO | - | SKU ID |
| product_name | VARCHAR(255) | NO | '' | Product name (snapshot) |
| sku_name | VARCHAR(255) | NO | '' | SKU name (snapshot) |
| image | VARCHAR(500) | NO | '' | Product image (snapshot) |
| quantity | INT | NO | 1 | Quantity shipped |
| created_at | BIGINT | NO | 0 | Creation timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_shipment_id` (`shipment_id`)
- KEY `idx_order_item_id` (`order_item_id`)
- KEY `idx_product_id` (`product_id`)

### refunds

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| order_id | VARCHAR(64) | NO | - | Order ID |
| refund_no | VARCHAR(32) | NO | '' | Refund number (auto-generated) |
| user_id | BIGINT | NO | - | Buyer user ID |
| type | TINYINT | NO | 1 | 1=full_refund, 2=partial_refund (future) |
| status | TINYINT | NO | 0 | 0=pending, 1=approved, 2=rejected, 3=completed, 4=cancelled |
| reason_type | VARCHAR(50) | NO | '' | Reason type code |
| reason | VARCHAR(500) | NO | '' | Refund reason |
| description | TEXT | YES | NULL | Detailed description |
| images | JSON | YES | NULL | Evidence image URLs |
| amount | BIGINT | NO | 0 | Refund amount (cents) |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| reject_reason | VARCHAR(500) | NO | '' | Rejection reason |
| approved_at | BIGINT | YES | NULL | Approval timestamp |
| approved_by | BIGINT | YES | NULL | Approver ID |
| completed_at | BIGINT | YES | NULL | Completion timestamp |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- UNIQUE KEY `uk_refund_no` (`tenant_id`, `refund_no`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_order_id` (`order_id`)
- KEY `idx_user_id` (`user_id`)
- KEY `idx_status` (`status`)
- KEY `idx_deleted_at` (`deleted_at`)

### refund_reasons (Reference Table)

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| code | VARCHAR(50) | NO | - | Reason code |
| name | VARCHAR(100) | NO | - | Reason display name |
| sort | INT | NO | 0 | Sort order |
| is_active | TINYINT | NO | 1 | Active status |
| created_at | BIGINT | NO | 0 | Creation timestamp |

**Predefined Refund Reasons:**

| Code | Name | Description |
|------|------|-------------|
| DEFECTIVE | Product Defective | Product has quality issues |
| WRONG_ITEM | Wrong Item Received | Received different product |
| NOT_AS_DESCRIBED | Not As Described | Product differs from description |
| DAMAGED | Damaged in Transit | Product damaged during shipping |
| NO_LONGER_NEEDED | No Longer Needed | Changed mind |
| LATE_DELIVERY | Late Delivery | Delivery took too long |
| OTHER | Other | Other reasons |

### carriers (Reference Table)

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| code | VARCHAR(20) | NO | - | Carrier code |
| name | VARCHAR(100) | NO | - | Carrier display name |
| tracking_url | VARCHAR(255) | NO | '' | Tracking URL template |
| is_active | TINYINT | NO | 1 | Active status |
| sort | INT | NO | 0 | Sort order |
| created_at | BIGINT | NO | 0 | Creation timestamp |

**Predefined Carriers:**

| Code | Name | Tracking URL Template |
|------|------|----------------------|
| SF | SF Express | https://www.sf-express.com/track?id={tracking_no} |
| YT | YTO Express | https://www.yto.net.cn/query.html?id={tracking_no} |
| ZT | ZTO Express | https://www.zto.com/track?id={tracking_no} |
| ST | STO Express | https://www.sto.cn/track?id={tracking_no} |
| YD | Yunda Express | https://www.yundaex.com/track?id={tracking_no} |
| EMS | EMS | https://www.ems.com.cn/track?id={tracking_no} |
| JD | JD Logistics | https://www.jdl.com/track?id={tracking_no} |
| OTHER | Other | - |

### Order Table Extensions

Add the following columns to the `orders` table:

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| fulfillment_status | TINYINT | NO | 0 | 0=pending, 1=partial_shipped, 2=shipped, 3=delivered |
| refund_status | TINYINT | NO | 0 | 0=none, 1=pending, 2=approved, 3=rejected, 4=completed |

---

## API Endpoints

### Shipment Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/shipments | Create shipment |
| POST | /api/v1/shipments/batch | Batch create shipments |
| GET | /api/v1/shipments | List shipments |
| GET | /api/v1/shipments/:id | Get shipment details |
| PUT | /api/v1/shipments/:id | Update shipment |
| PUT | /api/v1/shipments/:id/status | Update shipment status |
| GET | /api/v1/orders/:order_id/shipments | Get order shipments |
| GET | /api/v1/carriers | List available carriers |

### Refund Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/refunds | List refunds |
| GET | /api/v1/refunds/:id | Get refund details |
| PUT | /api/v1/refunds/:id/approve | Approve refund |
| PUT | /api/v1/refunds/:id/reject | Reject refund |
| GET | /api/v1/refund-reasons | List refund reasons |
| GET | /api/v1/refunds/statistics | Get refund statistics |

### Order Fulfillment (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/orders | List orders (with fulfillment_status filter) |
| PUT | /api/v1/orders/:id/ship | Ship order (create shipment) |
| GET | /api/v1/orders/:id | Get order with fulfillment info |
| GET | /api/v1/orders/fulfillment-summary | Get fulfillment summary counts |

### Refund Application (Shop - Buyer)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/orders/:order_id/refunds | Apply for refund |
| GET | /api/v1/orders/:order_id/refunds | Get order refunds |
| GET | /api/v1/refunds/:id | Get refund details |
| POST | /api/v1/refunds/:id/cancel | Cancel refund application |

### Order Tracking (Shop - Buyer)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/orders/:id/tracking | Get order tracking info |
| GET | /api/v1/shipments/:id/tracking | Get shipment tracking |

---

## Business Rules

### Shipment Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | Order must be paid | Cannot ship unpaid orders |
| BR-002 | Shipment number generation | Auto-generate unique shipment_no |
| BR-003 | Item quantity validation | Cannot ship more items than ordered |
| BR-004 | Order status update | Update order status on shipment creation |
| BR-005 | Fulfillment status calculation | Calculate based on shipped item quantities |
| BR-006 | Carrier validation | Carrier must be active if selected from list |
| BR-007 | Tracking number uniqueness | Tracking number should be unique per carrier |

### Refund Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-008 | Refund eligibility | Only paid orders can be refunded |
| BR-009 | Single active refund | Only one pending refund per order |
| BR-010 | Refund amount | Cannot exceed order pay_amount |
| BR-011 | Refund time limit | Must apply within N days of delivery (configurable) |
| BR-012 | Refund number generation | Auto-generate unique refund_no |
| BR-013 | Stock release | Approved refund releases locked stock |

### Status Transition Rules

**Order Fulfillment Status:**

| From | To | Trigger |
|------|-----|---------|
| pending | partial_shipped | First shipment created, items remaining |
| pending | shipped | All items shipped in one shipment |
| partial_shipped | shipped | All remaining items shipped |
| shipped | delivered | All shipments delivered |

**Refund Status:**

| From | To | Trigger |
|------|-----|---------|
| pending | approved | Merchant approves |
| pending | rejected | Merchant rejects |
| pending | cancelled | Buyer cancels |
| approved | completed | Payment refund processed |

---

## User Stories

### US-001: View Pending Shipment Orders

**Description**: As a Merchant Admin, I want to view all orders pending shipment so that I can process them efficiently.

**Acceptance Criteria**:
- Given I am logged in as a Merchant Admin
- When I navigate to Orders and filter by "Pending Shipment"
- Then I see all orders with status "paid" and fulfillment_status "pending"
- And each order shows order_no, buyer info, items, total amount, payment time
- And I can sort by payment time (oldest first by default)
- And I can search by order_no or buyer name

---

### US-002: Ship Single Order

**Description**: As a Merchant Admin, I want to ship a single order so that it can be delivered to the buyer.

**Acceptance Criteria**:
- Given an order with status "paid" and fulfillment_status "pending"
- When I click "Ship" on the order
- And I select carrier "SF Express"
- And I enter tracking number "SF1234567890"
- And I confirm the shipment
- Then a shipment record is created with status "shipped"
- And the order fulfillment_status changes to "shipped"
- And the order status changes to "shipped"
- And shipped_at timestamp is set

---

### US-003: Ship Order with Partial Items

**Description**: As a Merchant Admin, I want to ship some items of an order separately so that I can handle multi-warehouse or backorder scenarios.

**Acceptance Criteria**:
- Given an order with 3 items
- When I click "Ship" and select only 2 items to ship now
- And I enter carrier and tracking number
- And I confirm the shipment
- Then a shipment record is created with 2 items
- And the order fulfillment_status changes to "partial_shipped"
- And I can ship the remaining item later
- And the order shows "2/3 items shipped"

---

### US-4: Batch Ship Orders

**Description**: As a Merchant Admin, I want to ship multiple orders at once so that I can save time on bulk fulfillment.

**Acceptance Criteria**:
- Given I have selected 10 orders pending shipment
- When I click "Batch Ship"
- And I select carrier "ZTO Express"
- And I enter the starting tracking number
- And I confirm
- Then 10 shipment records are created
- And tracking numbers are auto-incremented
- And all orders' fulfillment_status change to "shipped"

---

### US-005: Update Shipment Status

**Description**: As a Merchant Admin, I want to update shipment status so that the tracking is accurate.

**Acceptance Criteria**:
- Given a shipment with status "shipped"
- When I update status to "delivered"
- And I confirm
- Then shipment status changes to "delivered"
- And delivered_at timestamp is set
- If all order shipments are delivered, order fulfillment_status changes to "delivered"

---

### US-006: View Shipment Details

**Description**: As a Merchant Admin, I want to view shipment details so that I can verify logistics information.

**Acceptance Criteria**:
- Given I am on the Shipments page
- When I click on a shipment
- Then I see shipment_no, carrier, tracking_no, status
- And I see all items in the shipment with quantities
- And I see shipping cost and weight
- And I see shipped_at and delivered_at timestamps

---

### US-007: View Order Tracking (Buyer)

**Description**: As a Buyer, I want to track my order so that I know when it will arrive.

**Acceptance Criteria**:
- Given I am logged in and have an order with status "shipped"
- When I navigate to My Orders and click on the order
- Then I see the tracking information: carrier name, tracking number
- And I see a tracking URL link to carrier website
- And I see the shipping address
- And I see the estimated delivery date (if available)

---

### US-008: Apply for Refund

**Description**: As a Buyer, I want to apply for a refund so that I can get my money back when there's an issue.

**Acceptance Criteria**:
- Given I have a paid order
- When I navigate to order detail and click "Apply for Refund"
- And I select reason "Product Defective"
- And I enter description "The product has a scratch on the screen"
- And I upload 2 evidence images
- And I submit
- Then a refund record is created with status "pending"
- And the order status changes to "refunding"
- And I see the refund amount equals the order pay_amount

---

### US-009: View Pending Refunds

**Description**: As a Merchant Admin, I want to view all pending refund requests so that I can process them.

**Acceptance Criteria**:
- Given I am logged in as a Merchant Admin
- When I navigate to Refunds page
- And I filter by status "pending"
- Then I see all refund requests with pending status
- And each shows refund_no, order_no, buyer, amount, reason, created time
- And I can sort by created time (oldest first by default)

---

### US-010: Approve Refund

**Description**: As a Merchant Admin, I want to approve a refund request so that the buyer can receive their money.

**Acceptance Criteria**:
- Given a refund request with status "pending"
- When I click on the refund
- And I review the order details, reason, and evidence
- And I click "Approve"
- And I confirm
- Then the refund status changes to "approved"
- And approved_at and approved_by are set
- If payment refund is processed, status changes to "completed"
- And order status changes to "refunded"

---

### US-011: Reject Refund

**Description**: As a Merchant Admin, I want to reject a refund request with a reason so that the buyer understands why.

**Acceptance Criteria**:
- Given a refund request with status "pending"
- When I click on the refund
- And I click "Reject"
- And I enter rejection reason "Evidence insufficient, please provide clearer photos"
- And I confirm
- Then the refund status changes to "rejected"
- And reject_reason is saved
- And the order status reverts to its previous status

---

### US-012: Cancel Refund Application

**Description**: As a Buyer, I want to cancel my refund application so that I can keep the order.

**Acceptance Criteria**:
- Given I have a pending refund application
- When I navigate to the refund and click "Cancel"
- And I confirm
- Then the refund status changes to "cancelled"
- And the order status reverts to its previous status

---

### US-013: View Refund Details

**Description**: As a Merchant Admin, I want to view refund details so that I can make an informed decision.

**Acceptance Criteria**:
- Given I click on a refund from the list
- Then I see refund_no, order_no, amount, currency
- And I see buyer info and order items
- And I see reason type and description
- And I see uploaded evidence images
- And I see status timeline with timestamps
- And I see approve/reject actions if status is pending

---

### US-014: View Refund Statistics

**Description**: As a Merchant Admin, I want to view refund statistics so that I can identify issues.

**Acceptance Criteria**:
- Given I navigate to After-Sales Statistics
- Then I see refund rate for selected time period
- And I see refund rate trend chart
- And I see breakdown by refund reason
- And I see top products with highest refund rates
- And I can filter by time range (7 days, 30 days, 90 days)

---

### US-015: Select Carrier from List

**Description**: As a Merchant Admin, I want to select a carrier from a predefined list so that data is consistent.

**Acceptance Criteria**:
- Given I am creating a shipment
- When I click the carrier dropdown
- Then I see a list of active carriers
- And each carrier shows name and code
- And I can select one carrier
- And the carrier code is stored for consistency

---

### US-016: Enter Custom Carrier

**Description**: As a Merchant Admin, I want to enter a custom carrier name so that I can use carriers not in the list.

**Acceptance Criteria**:
- Given I am creating a shipment
- When I select "Other" from the carrier dropdown
- Then I can enter a custom carrier name in the carrier field
- And the carrier_code is set to "OTHER"
- And the custom carrier name is stored in the carrier field

---

### US-017: Validate Tracking Number

**Description**: As a Merchant Admin, I want validation on tracking numbers so that errors are caught early.

**Acceptance Criteria**:
- Given I am creating a shipment with carrier "SF Express"
- When I enter an empty tracking number
- Then I see an error "Tracking number is required"
- When I enter a tracking number that already exists for the same carrier
- Then I see a warning "This tracking number already exists"

---

### US-018: Record Shipping Cost

**Description**: As a Merchant Admin, I want to record actual shipping cost so that I can track fulfillment expenses.

**Acceptance Criteria**:
- Given I am creating a shipment
- When I enter shipping cost "1200" (cents)
- And I confirm the shipment
- Then the shipping cost is saved
- And I can view total shipping costs in reports

---

### US-019: View Order Fulfillment Summary

**Description**: As a Merchant Admin, I want to see a summary of order fulfillment status so that I can prioritize work.

**Acceptance Criteria**:
- Given I am on the Orders page
- Then I see a summary bar with counts:
  - Pending Shipment: X orders
  - Partial Shipped: X orders
  - Shipped: X orders
  - Delivered: X orders
- And I can click on each to filter the list

---

### US-020: Create Split Shipment

**Description**: As a Merchant Admin, I want to create multiple shipments for one order so that items can be delivered from different locations.

**Acceptance Criteria**:
- Given an order with 5 items
- When I create the first shipment with 3 items
- Then fulfillment_status changes to "partial_shipped"
- And I see "3/5 items shipped" on the order
- When I create a second shipment with remaining 2 items
- Then fulfillment_status changes to "shipped"
- And I see "5/5 items shipped" on the order

---

### US-021: View Shipment History

**Description**: As a Merchant Admin, I want to view all shipments for an order so that I can track the complete fulfillment journey.

**Acceptance Criteria**:
- Given an order with multiple shipments
- When I view the order detail
- Then I see a "Shipments" section
- And I see all shipments with shipment_no, carrier, tracking, status
- And I see which items are in each shipment
- And I can click to view shipment details

---

### US-022: Automatic Order Status Update

**Description**: As the system, I want to automatically update order status when shipment is created so that consistency is maintained.

**Acceptance Criteria**:
- Given an order with status "paid" and fulfillment_status "pending"
- When a shipment is created
- Then the order status changes to "shipped"
- And shipped_at timestamp is set on the order
- And the fulfillment_status is calculated based on items shipped

---

### US-023: Refund Time Limit Validation

**Description**: As the system, I want to validate refund time limits so that expired refund requests are rejected.

**Acceptance Criteria**:
- Given the system has a refund_time_limit of 30 days
- And an order was delivered 45 days ago
- When a buyer tries to apply for refund
- Then they see an error "Refund period has expired (30 days after delivery)"
- And the refund application is not created

---

### US-024: Prevent Duplicate Pending Refund

**Description**: As the system, I want to prevent duplicate refund applications so that there's only one active refund per order.

**Acceptance Criteria**:
- Given an order has a pending refund application
- When the buyer tries to apply for another refund
- Then they see an error "This order already has a pending refund application"
- And the new refund is not created

---

### US-025: User Authentication for Admin Actions

**Description**: As a Merchant Admin, I must be authenticated to access fulfillment management so that only authorized users can modify data.

**Acceptance Criteria**:
- Given I am not logged in
- When I try to access /api/v1/shipments
- Then I receive a 401 Unauthorized error
- Given I am logged in as a tenant admin
- When I access /api/v1/shipments
- Then I receive the shipments list for my tenant only

---

### US-026: Tenant Isolation for Shipments

**Description**: As the system, I want to ensure shipments are isolated by tenant so that data is secure.

**Acceptance Criteria**:
- Given I am logged in as Tenant A admin
- When I request /api/v1/shipments
- Then I only see Tenant A's shipments
- And I cannot access Tenant B's shipments by ID
- And all queries are filtered by tenant_id

---

### US-027: Soft Delete Shipment

**Description**: As the system, I want shipments to be soft-deleted so that data can be recovered if needed.

**Acceptance Criteria**:
- Given I delete a shipment
- When the delete API is called
- Then deleted_at is set to current timestamp
- And the shipment no longer appears in normal queries
- But the data is preserved in the database

---

### US-028: Pagination for Large Lists

**Description**: As a Merchant Admin, I want shipment and refund lists to be paginated so that the interface remains performant.

**Acceptance Criteria**:
- Given there are 500 shipments
- When I view the shipments list
- Then I see 20 shipments per page (default)
- And I can navigate to other pages
- And I can change the page size (max 100)
- And the API returns total count for pagination

---

### US-029: Filter Shipments by Status

**Description**: As a Merchant Admin, I want to filter shipments by status so that I can focus on specific groups.

**Acceptance Criteria**:
- Given I am on the shipments list page
- When I select "Pending" from the status filter
- Then I see only pending shipments
- When I select "Delivered"
- Then I see only delivered shipments
- When I select multiple statuses
- Then I see shipments matching any selected status

---

### US-030: Filter Refunds by Status

**Description**: As a Merchant Admin, I want to filter refunds by status so that I can process pending requests.

**Acceptance Criteria**:
- Given I am on the refunds list page
- When I select "Pending" from the status filter
- Then I see only pending refund requests
- And the count badge shows the number of pending refunds

---

### US-031: Search Orders for Shipment

**Description**: As a Merchant Admin, I want to search orders by order number or buyer name so that I can quickly find specific orders.

**Acceptance Criteria**:
- Given I am on the pending shipment orders page
- When I type "ORD2026" in the search box
- Then I see only orders with "ORD2026" in the order number
- When I type "John" in the search box
- Then I see orders where buyer name contains "John"
- And results are filtered in real-time

---

### US-032: View Buyer Evidence Images

**Description**: As a Merchant Admin, I want to view buyer-uploaded evidence images so that I can assess refund requests.

**Acceptance Criteria**:
- Given a refund request with uploaded images
- When I view the refund details
- Then I see thumbnail previews of the images
- And I can click to view full-size images
- And I can download images if needed

---

### US-033: Concurrent Shipment Creation

**Description**: As the system, I want to handle concurrent shipment creation correctly so that over-shipping is prevented.

**Acceptance Criteria**:
- Given an order with 1 item remaining to ship
- When two admins try to create a shipment simultaneously
- Then only one shipment succeeds
- And the other receives an error "Items already shipped"
- And the database transaction ensures consistency

---

### US-034: Audit Trail for Refund Changes

**Description**: As a Merchant Admin, I want to see who processed a refund and when so that I can track decisions.

**Acceptance Criteria**:
- Given I view a refund's details
- When I look at the audit section
- Then I see created_by and created_at
- And I see updated_by and updated_at
- And I see approved_by and approved_at (if approved)
- And I can see the last modifier's name

---

### US-035: Refund Amount Cannot Exceed Order Total

**Description**: As the system, I want to ensure refund amount never exceeds the order pay_amount so that financial integrity is maintained.

**Acceptance Criteria**:
- Given an order with pay_amount of 10000 cents
- When a refund is created
- Then the refund amount is automatically set to 10000 cents
- And the refund amount cannot be modified to exceed 10000

---

### US-036: Order Must Be Paid for Refund

**Description**: As the system, I want to validate that only paid orders can have refund applications so that unpaid orders are not refunded.

**Acceptance Criteria**:
- Given an order with status "pending_payment"
- When a buyer tries to apply for refund
- Then they see an error "Cannot apply for refund on unpaid order"
- And the refund is not created

---

### US-037: Notify Buyer on Shipment

**Description**: As the system, I want to notify buyers when their order is shipped so that they are informed.

**Acceptance Criteria**:
- Given a shipment is created for an order
- When the shipment is confirmed
- Then a notification is sent to the buyer
- And the notification includes carrier and tracking number
- And the notification includes a tracking link

---

### US-038: Notify Buyer on Refund Status Change

**Description**: As the system, I want to notify buyers when their refund status changes so that they are informed.

**Acceptance Criteria**:
- Given a refund status changes to "approved"
- Then the buyer receives a notification "Your refund request has been approved"
- Given a refund status changes to "rejected"
- Then the buyer receives a notification with the rejection reason

---

### US-039: Problem Product Identification

**Description**: As a Merchant Admin, I want to identify products with high refund rates so that I can investigate quality issues.

**Acceptance Criteria**:
- Given I am on the After-Sales Statistics page
- When I view the "Problem Products" section
- Then I see products with refund rate above threshold (e.g., 10%)
- And each product shows refund count and total refund amount
- And I can click to view detailed refund breakdown

---

### US-040: Export Refund Report

**Description**: As a Merchant Admin, I want to export refund data so that I can analyze it externally.

**Acceptance Criteria**:
- Given I am on the refunds list page
- When I click "Export"
- Then a CSV file is downloaded
- And the file contains all refunds for the selected period
- And includes refund_no, order_no, amount, reason, status, dates

---

### US-041: Delivery Success Rate Tracking

**Description**: As a Merchant Admin, I want to see delivery success rates so that I can evaluate carrier performance.

**Acceptance Criteria**:
- Given I am on the After-Sales Statistics page
- When I view the "Delivery Statistics" section
- Then I see overall delivery success rate
- And I see breakdown by carrier
- And I see average delivery time by carrier

---

### US-042: Refund Reason Analysis

**Description**: As a Merchant Admin, I want to see refund reason breakdown so that I can identify common issues.

**Acceptance Criteria**:
- Given I am on the After-Sales Statistics page
- When I view the "Refund Reasons" section
- Then I see a pie chart of refund reasons
- And I see percentage for each reason type
- And I can filter by time period

---

### US-043: API Error Handling for Invalid Shipment

**Description**: As the system, I want to return clear error messages when shipment creation fails so that users understand the issue.

**Acceptance Criteria**:
- Given an invalid shipment request
- When the API is called
- Then a 400 error is returned
- And the error message clearly states the reason:
  - "Order not found" for non-existent order
  - "Order already shipped" for duplicate shipment
  - "Carrier is required" for missing carrier
  - "Tracking number is required" for missing tracking

---

### US-044: API Error Handling for Invalid Refund

**Description**: As the system, I want to return clear error messages when refund application fails so that users understand the issue.

**Acceptance Criteria**:
- Given an invalid refund request
- When the API is called
- Then a 400 error is returned
- And the error message clearly states the reason:
  - "Order not found" for non-existent order
  - "Order already has pending refund" for duplicate
  - "Refund period expired" for time limit exceeded
  - "Reason is required" for missing reason

---

### US-045: Real-time Order Count Updates

**Description**: As a Merchant Admin, I want to see real-time order counts by fulfillment status so that I know my workload.

**Acceptance Criteria**:
- Given I am on the Orders page
- When I view the summary bar
- Then I see accurate counts of orders by fulfillment status
- And counts update when I refresh the page
- And pending shipment count is highlighted if non-zero

---

### US-046: Shipment Number Auto-Generation

**Description**: As the system, I want to auto-generate shipment numbers so that they are unique and traceable.

**Acceptance Criteria**:
- Given a shipment is created
- When the system generates shipment_no
- Then it follows format "SHP{YYYYMMDD}{sequence}"
- And it is unique within the tenant
- And the sequence resets daily

---

### US-047: Refund Number Auto-Generation

**Description**: As the system, I want to auto-generate refund numbers so that they are unique and traceable.

**Acceptance Criteria**:
- Given a refund is created
- When the system generates refund_no
- Then it follows format "REF{YYYYMMDD}{sequence}"
- And it is unique within the tenant
- And the sequence resets daily

---

### US-048: Order Item Shipment Tracking

**Description**: As a Merchant Admin, I want to see which items have been shipped so that I know what remains.

**Acceptance Criteria**:
- Given an order with partial shipments
- When I view the order detail
- Then each item shows "Shipped" or "Pending" status
- And I see quantity shipped vs ordered
- And I can see which shipment each item is in

---

### US-049: Mobile Responsive Shipment Interface

**Description**: As a Merchant Admin, I want the shipment interface to work on mobile so that I can ship orders from anywhere.

**Acceptance Criteria**:
- Given I access the admin on a mobile device
- When I view the pending shipments page
- Then the layout is responsive and usable
- And I can create shipments with the mobile interface
- And key actions are easily tappable

---

### US-050: Bulk Reject Refunds

**Description**: As a Merchant Admin, I want to reject multiple refund requests at once so that I can process them efficiently.

**Acceptance Criteria**:
- Given I have selected 5 pending refund requests
- When I click "Batch Reject"
- And I enter a common rejection reason
- And I confirm
- Then all 5 refunds are rejected
- And the same reject_reason is applied to all
- And buyers are notified

---

### US-051: Refund Workflow Permission

**Description**: As a Tenant Admin, I want to control who can approve refunds so that financial controls are maintained.

**Acceptance Criteria**:
- Given I am configuring roles
- When I assign "Refund Manager" role to a user
- Then that user can approve/reject refunds
- And users without the role can only view refunds
- And high-value refunds may require additional approval (future)

---

### US-052: Shipment Draft Status

**Description**: As a Merchant Admin, I want to save shipment drafts so that I can prepare shipments before confirming.

**Acceptance Criteria**:
- Given I am creating a shipment
- When I click "Save Draft"
- Then the shipment is saved with status "pending"
- And the order fulfillment_status is not updated
- And I can later edit and confirm the shipment

---

### US-053: Cancel Pending Shipment

**Description**: As a Merchant Admin, I want to cancel a pending shipment so that I can correct mistakes before it ships.

**Acceptance Criteria**:
- Given a shipment with status "pending"
- When I click "Cancel"
- And I confirm
- Then the shipment status changes to "cancelled"
- And the order fulfillment_status is recalculated
- And the items are available for shipping again

---

### US-054: View Shipping Address

**Description**: As a Merchant Admin, I want to view the shipping address clearly so that I can verify it matches the shipment.

**Acceptance Criteria**:
- Given I am creating a shipment
- Then I see the complete shipping address
- And I see recipient name and phone
- And I can copy the address for label printing

---

### US-055: Refund Revert Order Status

**Description**: As the system, I want to revert order status when refund is rejected or cancelled so that the order continues normally.

**Acceptance Criteria**:
- Given an order with status "refunding" due to a pending refund
- When the refund is rejected
- Then the order status reverts to "shipped" (if previously shipped)
- Or reverts to "paid" (if never shipped)
- And the order can continue through normal fulfillment

---

### US-056: Inventory Release on Refund Approval

**Description**: As the system, I want to release inventory when refund is approved so that stock is accurate.

**Acceptance Criteria**:
- Given a refund is approved
- When the refund status changes to "approved"
- Then an event is published to inventory service
- And the inventory service releases the locked stock
- And the stock becomes available for other orders

---

### US-057: Payment Refund Integration

**Description**: As the system, I want to integrate with payment service for refund processing so that buyers receive their money.

**Acceptance Criteria**:
- Given a refund is approved
- When the refund status changes to "approved"
- Then the system calls payment service to process refund
- And if payment refund succeeds, refund status changes to "completed"
- And if payment refund fails, an error is logged and retry is scheduled

---

### US-058: Carrier Performance Dashboard

**Description**: As a Merchant Admin, I want to compare carrier performance so that I can choose better carriers.

**Acceptance Criteria**:
- Given I am on the After-Sales Statistics page
- When I view "Carrier Performance"
- Then I see delivery time by carrier
- And I see delivery success rate by carrier
- And I see total shipments by carrier
- And I can filter by time period

---

### US-059: Refund Rate Threshold Alert

**Description**: As a Merchant Admin, I want to be alerted when refund rate exceeds threshold so that I can investigate issues.

**Acceptance Criteria**:
- Given daily refund rate calculation runs
- When the refund rate exceeds 5%
- Then an alert is sent to admin users
- And the alert includes the current rate and trend

---

### US-060: Order Fulfillment Timeline

**Description**: As a Merchant Admin, I want to see a timeline of order fulfillment events so that I can understand the journey.

**Acceptance Criteria**:
- Given I view an order detail
- When I look at the "Fulfillment Timeline" section
- Then I see events in chronological order:
  - Order placed
  - Payment received
  - Shipment created
  - Shipment delivered
  - Refund applied (if any)
- And each event shows timestamp and actor