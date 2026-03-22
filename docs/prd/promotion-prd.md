# Promotion Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Promotion Module PRD |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-22 |
| Author | Product Team |
| Module Location | admin/internal/domain/promotion/ |

---

## Product Overview

### Summary

The Promotion module is an independent bounded context within the ShopJoy e-commerce SaaS platform that enables merchants to create, manage, and apply various promotional rules to drive sales and customer engagement. This module handles discount calculations, coupon management, and promotion scope configuration in a multi-tenant environment.

The MVP release focuses on essential promotion types: order threshold discounts (tiered full-reduction), fixed amount discounts, percentage discounts, and coupon management. The system supports real-time discount calculation at checkout with proper conflict resolution based on priority rules.

### Problem Statement

Merchants on the ShopJoy platform currently lack a unified promotion management system. They need:

- Flexible discount rules to incentivize purchases (e.g., spend $100 get $10 off)
- Coupon generation and tracking capabilities
- Control over which products/categories promotions apply to
- Ability to schedule promotions with time ranges
- Clear conflict resolution when multiple promotions overlap

### Solution Overview

Implement a domain-driven Promotion bounded context that provides:

1. Multiple promotion types with configurable rules
2. Promotion scope targeting (products, categories, brands, storewide)
3. Time-based activation and priority-based conflict resolution
4. Coupon generation, distribution, and redemption tracking
5. Real-time discount calculation engine

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | Increase average order value | AOV growth | +15% within 3 months |
| BG-002 | Improve merchant retention | Feature adoption rate | 60% of merchants using promotions within 6 months |
| BG-003 | Reduce development time for new promotion types | Time to implement new type | <2 weeks |
| BG-004 | Support multi-tenant scalability | Concurrent tenants | 1,000+ tenants |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | Create and manage promotions easily | Merchant Admin |
| UG-002 | View all active promotions and their status | Merchant Admin |
| UG-003 | Generate and distribute coupons to customers | Merchant Admin |
| UG-004 | See applicable discounts at checkout | Shopper |
| UG-005 | Apply coupons and see updated totals | Shopper |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | Buy X Get Y promotions | Phase 2 feature |
| NG-002 | Flash sale/Limited-time offers with inventory locks | Phase 2 feature |
| NG-003 | Member-tier exclusive promotions | Phase 2 feature |
| NG-004 | Bundle promotions | Phase 2 feature |
| NG-005 | Complex rule engine with scripting/DSL | Phase 2 feature |
| NG-006 | Automatic promotion recommendation | Out of scope |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | Store owner or marketing manager | Create promotions, manage coupons, view reports |
| Shopper | End customer making purchases | Apply discounts at checkout, use coupons |

### Persona Details

#### Merchant Admin (Primary)

- **Demographics**: Small to medium business owners, e-commerce managers
- **Technical Proficiency**: Moderate
- **Goals**: Drive sales, clear inventory, reward loyal customers
- **Pain Points**: Complex promotion setup, unclear discount application
- **Frequency**: Daily to weekly promotion management

#### Shopper (Secondary)

- **Demographics**: Online shoppers of all ages
- **Technical Proficiency**: Varies widely
- **Goals**: Get the best deal, easy coupon application
- **Pain Points**: Unclear why discount was not applied, expired coupons
- **Frequency**: Per shopping session

### Role-Based Access

| Role | Permissions |
|------|-------------|
| Tenant Admin | Full CRUD on promotions and coupons |
| Tenant Marketing Manager | Create/Edit promotions, manage coupons |
| Tenant Staff | View promotions (read-only) |
| Shopper | Apply coupons, view applicable promotions |

---

## Functional Requirements

### Priority 1: Order Threshold Discount (Full Reduction)

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | Create threshold discount | Merchant can create "spend X, save Y" promotions | P0 |
| FR-002 | Tiered thresholds | Support multiple tiers (e.g., $50 save $5, $100 save $15) | P0 |
| FR-003 | Maximum discount cap | Optional cap on discount amount | P1 |
| FR-004 | Minimum order amount | Required threshold before discount applies | P0 |

### Priority 2: Fixed and Percentage Discounts

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-005 | Fixed amount discount | Create promotions that deduct a fixed amount | P0 |
| FR-006 | Percentage discount | Create promotions that deduct a percentage (e.g., 20% off) | P0 |
| FR-007 | Maximum discount limit | Cap percentage discounts at a maximum amount | P1 |

### Priority 3: Coupon Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-008 | Create coupon templates | Define coupon type, value, and constraints | P0 |
| FR-009 | Fixed amount coupons | Coupons that deduct a fixed amount | P0 |
| FR-010 | Percentage coupons | Coupons that deduct a percentage | P0 |
| FR-011 | Minimum spend requirement | Set minimum order amount for coupon use | P0 |
| FR-012 | Usage limits | Total usage count and per-user limit | P0 |
| FR-013 | Coupon codes | Generate unique or batch codes | P0 |
| FR-014 | Coupon expiration | Set validity period for coupons | P0 |

### Priority 4: Promotion Scope

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-015 | Storewide promotions | Apply to all products | P0 |
| FR-016 | Category-specific promotions | Target specific categories | P0 |
| FR-017 | Brand-specific promotions | Target specific brands | P1 |
| FR-018 | Product-specific promotions | Target specific products | P0 |
| FR-019 | Exclude items | Exclude certain products from promotions | P1 |

### Priority 5: Time Range and Priority

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | Schedule promotions | Set start and end datetime | P0 |
| FR-021 | Priority ordering | Assign priority for conflict resolution | P0 |
| FR-022 | Status management | Draft, Active, Inactive, Ended states | P0 |
| FR-023 | Immediate activation | Toggle promotion on/off | P1 |

### Priority 6: Conflict Resolution

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-024 | Priority-based selection | Higher priority promotion wins conflicts | P0 |
| FR-025 | Mutually exclusive option | Configure promotions as non-combinable | P1 |
| FR-026 | Best discount automatic | Automatically apply best applicable discount | P1 |

### Priority 7: Real-time Calculation

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-027 | Calculate cart discount | Compute applicable discounts for cart | P0 |
| FR-028 | Apply coupon discount | Calculate and apply coupon discount | P0 |
| FR-029 | Return discount breakdown | Itemized discount details | P1 |
| FR-030 | Validate promotion eligibility | Check if cart meets promotion conditions | P0 |

---

## User Experience

### Entry Points

| Entry Point | Description | User Flow |
|-------------|-------------|-----------|
| Admin Dashboard | Promotions menu in sidebar | Dashboard > Promotions > List |
| Product Edit Page | Quick promotion assignment | Product > Edit > Assign Promotion |
| Checkout Page | Apply coupon code | Cart > Checkout > Enter Code |

### Core Experience

#### Creating a Promotion (Merchant Admin)

1. Navigate to Promotions section
2. Click "Create Promotion"
3. Select promotion type (Threshold Discount / Fixed Discount / Percentage Discount)
4. Configure rule details:
   - For Threshold: Add tiers (min amount, discount amount)
   - For Fixed/Percentage: Set value, max discount
5. Set scope (products, categories, brands, or storewide)
6. Set time range and priority
7. Save as draft or activate immediately
8. View in promotions list with status indicator

#### Applying a Coupon (Shopper)

1. Add items to cart
2. Navigate to cart or checkout
3. Enter coupon code in designated field
4. Click "Apply"
5. See updated cart total with discount shown
6. Proceed to checkout with discount applied

### Advanced Features

| Feature | Description |
|---------|-------------|
| Bulk Coupon Generation | Generate multiple unique codes for distribution |
| Promotion Preview | Preview how promotion affects specific products |
| Usage Analytics | View coupon usage statistics |
| Conflict Preview | See which promotions would conflict |

### UI/UX Highlights

- Clear visual hierarchy for promotion rules
- Real-time validation for coupon codes
- Color-coded status indicators (Active/Inactive/Draft)
- Tiered threshold visual representation
- Inline help for promotion configuration
- Mobile-responsive admin interface

---

## Narrative

As a Merchant Admin, I want to create a "Spend $100, Save $15" promotion for my summer sale. I navigate to the Promotions section in my admin dashboard and click "Create Promotion". I select "Threshold Discount" as the type and add a single tier: minimum spend $100, discount $15. I choose "Storewide" as the scope so all products are eligible. I set the promotion to run from July 1st to July 31st with priority 10. After reviewing the details, I save and activate the promotion.

When a shopper adds items worth $120 to their cart during the sale period, they see the promotion automatically applied at checkout. The cart shows "Summer Sale: -$15.00" in the discount breakdown. The shopper is pleased with the savings and completes the purchase.

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Promotion Creation Time | Average time to create a promotion | <5 minutes |
| Coupon Redemption Rate | Percentage of issued coupons redeemed | >20% |
| Cart Abandonment Rate | Reduction after showing applicable promotions | -10% |

### Business Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Feature Adoption | Percentage of merchants using promotions | 60% in 6 months |
| Average Order Value | AOV for orders with promotions vs without | +15% lift |
| Promotion ROI | Revenue generated vs discount cost | >3:1 ratio |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API Response Time | Discount calculation endpoint | <100ms p95 |
| System Uptime | Promotion service availability | 99.9% |
| Error Rate | Failed promotion applications | <0.1% |

---

## Technical Considerations

### Integration Points

| Integration | Description | Method |
|-------------|-------------|--------|
| Product Service | Fetch product details for scope validation | Internal API |
| Cart Service | Apply discounts to cart | Domain Event |
| Order Service | Record promotion usage | Domain Event |
| User Service | Validate user-specific limits | Internal API |

### Data Storage and Privacy

| Aspect | Consideration |
|--------|---------------|
| Monetary Values | Use int64 (cents) for all amounts, no floating-point |
| Time Storage | All timestamps in UTC, stored as Unix BIGINT |
| Multi-Tenancy | All queries filtered by tenant_id |
| Soft Delete | Support deleted_at for data recovery |
| PII Handling | Coupon codes are not PII, no special handling needed |

### Scalability and Performance

| Consideration | Strategy |
|---------------|----------|
| High Query Volume | Cache active promotions per tenant (Redis) |
| Concurrent Promotions | Limit active promotions per tenant (configurable) |
| Calculation Load | Stateless calculation service, horizontal scaling |
| Database Growth | Partition promotion_usage by tenant_id |

### Potential Challenges

| Challenge | Mitigation |
|-----------|------------|
| Complex Rule Conflicts | Clear priority system, conflict preview UI |
| Race Conditions on Coupon Use | Database-level constraints, optimistic locking |
| Timezone Confusion | Always store/display in UTC, frontend converts |
| Currency Mismatch | Validate currency on promotion application |

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Team Size | Description |
|-------|----------|-----------|-------------|
| Phase 1: Foundation | 2 weeks | 2 backend, 1 frontend | Core promotion types, basic scope |
| Phase 2: Coupons | 2 weeks | 2 backend, 1 frontend | Coupon management, redemption |
| Phase 3: Calculation Engine | 1.5 weeks | 2 backend | Real-time discount calculation |
| Phase 4: Testing & Polish | 1 week | Full team | E2E testing, performance tuning |

**Total Estimate: 6.5 weeks**

### Suggested Phases

#### Phase 1: Foundation (Week 1-2)

- Implement promotion entity and repository
- Create promotion types: FULL_REDUCE, FIXED_DISCOUNT, PERCENT_DISCOUNT
- Build CRUD APIs for promotions
- Implement scope targeting (products, categories, brands)
- Add time range and priority support
- Basic admin UI for promotion management

#### Phase 2: Coupon System (Week 3-4)

- Design coupon entity and user_coupon relationship
- Implement coupon CRUD APIs
- Build coupon code generation (single and batch)
- Add usage tracking and limits
- Implement coupon redemption flow
- Admin UI for coupon management

#### Phase 3: Calculation Engine (Week 5-6)

- Build discount calculation service
- Implement priority-based conflict resolution
- Create real-time cart discount API
- Integrate with cart/checkout flow
- Add promotion usage recording

#### Phase 4: Testing and Launch (Week 7)

- End-to-end testing
- Performance testing and optimization
- Documentation
- Staging deployment and QA
- Production deployment

---

## Database Schema

### promotions

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| name | VARCHAR(200) | NO | - | Promotion name |
| description | TEXT | YES | NULL | Promotion description |
| type | VARCHAR(32) | NO | - | FULL_REDUCE, FIXED_DISCOUNT, PERCENT_DISCOUNT |
| status | TINYINT | NO | 0 | 0=draft, 1=active, 2=inactive |
| priority | INT | NO | 0 | Higher = more priority |
| start_at | BIGINT | NO | - | Start timestamp (UTC) |
| end_at | BIGINT | NO | - | End timestamp (UTC) |
| scope_type | VARCHAR(32) | NO | 'storewide' | storewide, products, categories, brands |
| scope_ids | JSON | YES | NULL | Array of IDs for scope |
| exclude_ids | JSON | YES | NULL | Array of excluded product IDs |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_status` (`status`)
- KEY `idx_tenant_status_time` (`tenant_id`, `status`, `start_at`, `end_at`)
- KEY `idx_deleted_at` (`deleted_at`)

### promotion_rules

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| promotion_id | BIGINT | NO | - | Parent promotion ID |
| condition_type | VARCHAR(32) | NO | - | min_amount, min_quantity |
| condition_value | BIGINT | NO | 0 | Threshold value (cents) |
| action_type | VARCHAR(32) | NO | - | fixed_amount, percentage |
| action_value | BIGINT | NO | 0 | Discount value (cents or basis points) |
| max_discount | BIGINT | NO | 0 | Maximum discount (cents), 0 = no limit |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| sort_order | INT | NO | 0 | Sort order for tiered rules |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_promotion_id` (`promotion_id`)

### coupons

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| name | VARCHAR(200) | NO | - | Coupon name |
| code | VARCHAR(50) | NO | - | Unique coupon code |
| description | TEXT | YES | NULL | Coupon description |
| type | VARCHAR(32) | NO | - | fixed_amount, percentage |
| value | BIGINT | NO | 0 | Discount value (cents or basis points) |
| min_amount | BIGINT | NO | 0 | Minimum order amount (cents) |
| max_discount | BIGINT | NO | 0 | Maximum discount (cents), 0 = no limit |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| total_count | INT | NO | 0 | Total available, 0 = unlimited |
| used_count | INT | NO | 0 | Number of times used |
| per_user_limit | INT | NO | 0 | Max uses per user, 0 = unlimited |
| status | TINYINT | NO | 0 | 0=inactive, 1=active |
| start_at | BIGINT | NO | - | Start timestamp (UTC) |
| end_at | BIGINT | NO | - | End timestamp (UTC) |
| scope_type | VARCHAR(32) | NO | 'storewide' | storewide, products, categories, brands |
| scope_ids | JSON | YES | NULL | Array of IDs for scope |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_status` (`status`)
- KEY `idx_tenant_status_time` (`tenant_id`, `status`, `start_at`, `end_at`)
- KEY `idx_deleted_at` (`deleted_at`)

### user_coupons

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| user_id | BIGINT | NO | - | User ID |
| coupon_id | BIGINT | NO | - | Coupon ID |
| status | TINYINT | NO | 0 | 0=unused, 1=used, 2=expired |
| used_at | BIGINT | YES | NULL | Usage timestamp |
| order_id | VARCHAR(50) | NO | '' | Order ID when used |
| received_at | BIGINT | NO | - | Receipt timestamp |
| expire_at | BIGINT | NO | - | Expiration timestamp |
| created_at | BIGINT | NO | 0 | Creation timestamp |
| updated_at | BIGINT | NO | 0 | Update timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_user_id` (`user_id`)
- KEY `idx_coupon_id` (`coupon_id`)
- KEY `idx_user_status` (`user_id`, `status`)

### promotion_usage

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| promotion_id | BIGINT | NO | - | Promotion ID |
| rule_id | BIGINT | YES | NULL | Rule ID (for tiered promotions) |
| order_id | VARCHAR(50) | NO | - | Order ID |
| user_id | BIGINT | NO | - | User ID |
| discount_amount | BIGINT | NO | 0 | Discount applied (cents) |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| original_amount | BIGINT | NO | 0 | Original order amount |
| final_amount | BIGINT | NO | 0 | Final order amount |
| coupon_id | BIGINT | YES | NULL | Coupon ID if used |
| created_at | BIGINT | NO | 0 | Creation timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_promotion_id` (`promotion_id`)
- KEY `idx_order_id` (`order_id`)
- KEY `idx_user_id` (`user_id`)
- KEY `idx_created_at` (`created_at`)

---

## API Endpoints

### Promotion Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/promotions | Create promotion |
| PUT | /api/v1/promotions/:id | Update promotion |
| GET | /api/v1/promotions/:id | Get promotion details |
| GET | /api/v1/promotions | List promotions |
| DELETE | /api/v1/promotions/:id | Delete promotion |
| POST | /api/v1/promotions/:id/activate | Activate promotion |
| POST | /api/v1/promotions/:id/deactivate | Deactivate promotion |
| GET | /api/v1/promotions/:id/rules | Get promotion rules |
| POST | /api/v1/promotions/:id/rules | Add promotion rule |
| PUT | /api/v1/promotion-rules/:id | Update promotion rule |
| DELETE | /api/v1/promotion-rules/:id | Delete promotion rule |

### Coupon Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/coupons | Create coupon |
| PUT | /api/v1/coupons/:id | Update coupon |
| GET | /api/v1/coupons/:id | Get coupon details |
| GET | /api/v1/coupons | List coupons |
| DELETE | /api/v1/coupons/:id | Delete coupon |
| POST | /api/v1/coupons/generate | Generate batch coupon codes |
| GET | /api/v1/coupons/:id/usage | Get coupon usage history |
| POST | /api/v1/user-coupons | Issue coupon to user |
| GET | /api/v1/user-coupons | List user's coupons |

### Discount Calculation

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/cart/calculate-discount | Calculate applicable discounts |
| POST | /api/v1/cart/apply-coupon | Apply coupon to cart |
| POST | /api/v1/cart/remove-coupon | Remove coupon from cart |
| GET | /api/v1/promotions/active | Get active promotions for display |

---

## Business Rules

### Promotion Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | Single active scope | A product can only have one active promotion of the same type at a time |
| BR-002 | Priority resolution | When conflicts occur, higher priority promotion wins |
| BR-003 | Time validity | Promotion must be within start_at and end_at to be active |
| BR-004 | Currency match | Promotion currency must match order currency |
| BR-005 | Amount validation | Discount cannot exceed order subtotal |
| BR-006 | Tier selection | For tiered promotions, highest applicable tier is selected |

### Coupon Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-007 | Code uniqueness | Coupon code must be unique per tenant |
| BR-008 | Usage limit | Coupon cannot be used more than total_count |
| BR-009 | Per-user limit | User cannot exceed per_user_limit for coupon |
| BR-010 | Minimum spend | Order must meet min_amount to use coupon |
| BR-011 | Time validity | Coupon must be within start_at and end_at |
| BR-012 | Single use per order | Only one coupon per order |
| BR-013 | First apply promotion, then coupon | Order: automatic promotions first, then coupon |

### Conflict Resolution Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-014 | Priority order | Apply promotions in descending priority order |
| BR-015 | Coupon combines | Coupons stack with automatic promotions |
| BR-016 | Manual selection | If multiple equal priority, let user choose (Phase 2) |

---

## User Stories

### US-001: Create Threshold Discount Promotion

**Description**: As a Merchant Admin, I want to create a threshold discount promotion so that customers are incentivized to spend more.

**Acceptance Criteria**:
- Given I am logged in as a Merchant Admin
- When I navigate to Promotions and click "Create Promotion"
- And I select "Threshold Discount" type
- And I enter name, description
- And I add at least one tier with minimum amount and discount
- And I set time range
- And I save the promotion
- Then the promotion is created with status "Draft"
- And I can view the promotion in the promotions list

---

### US-002: Create Tiered Threshold Discount

**Description**: As a Merchant Admin, I want to create multiple tiers for threshold discounts so that higher spenders get better discounts.

**Acceptance Criteria**:
- Given I am creating a threshold discount promotion
- When I add multiple tiers (e.g., $50 save $5, $100 save $15, $200 save $40)
- And I save the promotion
- Then all tiers are saved with correct sort order
- When a customer's cart meets tier 2 threshold ($100)
- Then the tier 2 discount ($15) is applied, not tier 1

---

### US-003: Create Fixed Amount Discount

**Description**: As a Merchant Admin, I want to create a fixed amount discount promotion so that customers get a specific amount off.

**Acceptance Criteria**:
- Given I am creating a promotion
- When I select "Fixed Discount" type
- And I enter discount amount (e.g., $10)
- And I optionally set maximum discount
- And I save the promotion
- Then the promotion is created successfully
- When applied to a $50 order, the discount is $10

---

### US-004: Create Percentage Discount

**Description**: As a Merchant Admin, I want to create a percentage discount promotion so that customers get a percentage off their order.

**Acceptance Criteria**:
- Given I am creating a promotion
- When I select "Percentage Discount" type
- And I enter percentage (e.g., 20%)
- And I set maximum discount cap (e.g., $50 max)
- And I save the promotion
- Then the promotion is created successfully
- When applied to a $100 order, the discount is $20
- When applied to a $500 order, the discount is capped at $50

---

### US-005: Set Promotion Scope to Products

**Description**: As a Merchant Admin, I want to limit a promotion to specific products so that I can target my discount strategy.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I select scope type "Products"
- And I search and select specific products
- And I save the promotion
- Then the promotion only applies to selected products
- And the scope_ids field contains the product IDs

---

### US-006: Set Promotion Scope to Categories

**Description**: As a Merchant Admin, I want to limit a promotion to specific categories so that I can run category-wide sales.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I select scope type "Categories"
- And I select one or more categories
- And I save the promotion
- Then the promotion applies to all products in selected categories

---

### US-007: Set Promotion Scope to Brands

**Description**: As a Merchant Admin, I want to limit a promotion to specific brands so that I can run brand-specific campaigns.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I select scope type "Brands"
- And I select one or more brands
- And I save the promotion
- Then the promotion applies to all products of selected brands

---

### US-008: Set Storewide Promotion

**Description**: As a Merchant Admin, I want to create a storewide promotion so that all products are eligible for the discount.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I select scope type "Storewide"
- And I save the promotion
- Then the promotion applies to all products in the store
- And no scope_ids are required

---

### US-009: Exclude Products from Promotion

**Description**: As a Merchant Admin, I want to exclude certain products from a promotion so that high-margin items are not discounted.

**Acceptance Criteria**:
- Given I am creating a storewide or category promotion
- When I add products to the exclusion list
- And I save the promotion
- Then excluded products do not receive the discount
- And the exclude_ids field contains the excluded product IDs

---

### US-010: Schedule Promotion with Time Range

**Description**: As a Merchant Admin, I want to set start and end times for a promotion so that it automatically activates and deactivates.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I set start_at to "2026-07-01T00:00:00Z"
- And I set end_at to "2026-07-31T23:59:59Z"
- And I set status to "Active"
- And I save the promotion
- Then the promotion is only valid within the time range
- Before start_at, the promotion is not applicable
- After end_at, the promotion is not applicable

---

### US-011: Set Promotion Priority

**Description**: As a Merchant Admin, I want to set priority for promotions so that conflicts are resolved correctly.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I set priority to 10
- And another promotion for the same product has priority 5
- And both promotions are active
- Then the priority 10 promotion is applied

---

### US-012: Activate Promotion

**Description**: As a Merchant Admin, I want to activate a draft promotion so that it becomes effective.

**Acceptance Criteria**:
- Given a promotion with status "Draft"
- When I click "Activate"
- Then the promotion status changes to "Active"
- And the promotion is now eligible for discount calculation
- If current time is within the time range, the promotion applies immediately

---

### US-013: Deactivate Promotion

**Description**: As a Merchant Admin, I want to deactivate an active promotion so that it stops applying immediately.

**Acceptance Criteria**:
- Given an active promotion
- When I click "Deactivate"
- Then the promotion status changes to "Inactive"
- And the promotion no longer applies to new orders
- Existing orders with the promotion are not affected

---

### US-014: Edit Promotion

**Description**: As a Merchant Admin, I want to edit an existing promotion so that I can adjust its settings.

**Acceptance Criteria**:
- Given an existing promotion
- When I click "Edit"
- Then I can modify name, description, rules, scope, time range, priority
- And changes are saved successfully
- If promotion is active, changes take effect immediately

---

### US-015: Delete Promotion

**Description**: As a Merchant Admin, I want to delete a promotion so that it is removed from the system.

**Acceptance Criteria**:
- Given a promotion with status "Draft" or "Inactive"
- When I click "Delete"
- Then the promotion is soft-deleted (deleted_at set)
- And the promotion no longer appears in the list
- Active promotions cannot be deleted (must deactivate first)

---

### US-016: View Promotion List

**Description**: As a Merchant Admin, I want to view all promotions so that I can manage them effectively.

**Acceptance Criteria**:
- Given I am on the Promotions page
- Then I see a list of all promotions for my tenant
- And each promotion shows name, type, status, time range
- And I can filter by status and type
- And I can search by name
- And pagination is available for large lists

---

### US-017: View Promotion Details

**Description**: As a Merchant Admin, I want to view detailed information about a promotion so that I can verify its configuration.

**Acceptance Criteria**:
- Given I click on a promotion
- Then I see all promotion details
- And I see all configured rules/tiers
- And I see the scope (products, categories, brands)
- And I see usage statistics (times applied, total discount given)

---

### US-018: Create Fixed Amount Coupon

**Description**: As a Merchant Admin, I want to create a fixed amount coupon so that customers can get a specific discount.

**Acceptance Criteria**:
- Given I navigate to Coupons and click "Create Coupon"
- When I select "Fixed Amount" type
- And I enter coupon name, code, discount value
- And I set minimum spend and usage limits
- And I set validity period
- And I save the coupon
- Then the coupon is created with status "Active" (if start time is now or past)

---

### US-019: Create Percentage Coupon

**Description**: As a Merchant Admin, I want to create a percentage coupon so that customers can get a percentage off.

**Acceptance Criteria**:
- Given I am creating a coupon
- When I select "Percentage" type
- And I enter percentage value (e.g., 15)
- And I optionally set maximum discount
- And I save the coupon
- Then the coupon is created successfully

---

### US-020: Generate Batch Coupon Codes

**Description**: As a Merchant Admin, I want to generate multiple unique coupon codes so that I can distribute them to customers.

**Acceptance Criteria**:
- Given I am creating a coupon
- When I select "Generate Multiple Codes"
- And I specify the number of codes (e.g., 100)
- And I specify code format (prefix, length)
- And I save
- Then the system generates 100 unique coupon codes
- And all codes are linked to the same coupon template
- And I can export the codes as CSV

---

### US-021: Set Coupon Usage Limits

**Description**: As a Merchant Admin, I want to set usage limits for coupons so that I can control distribution costs.

**Acceptance Criteria**:
- Given I am creating or editing a coupon
- When I set total_count to 1000
- And I set per_user_limit to 1
- And I save the coupon
- Then the coupon can only be used 1000 times total
- And each user can only use it once
- When limits are reached, the coupon becomes invalid

---

### US-022: Set Coupon Minimum Spend

**Description**: As a Merchant Admin, I want to set a minimum order amount for coupons so that discounts apply only to larger orders.

**Acceptance Criteria**:
- Given I am creating a coupon
- When I set min_amount to $50
- And I save the coupon
- Then the coupon can only be applied to orders of $50 or more
- When a customer tries to use it on a $40 order, an error is shown

---

### US-023: View Coupon List

**Description**: As a Merchant Admin, I want to view all coupons so that I can manage them effectively.

**Acceptance Criteria**:
- Given I am on the Coupons page
- Then I see a list of all coupons for my tenant
- And each coupon shows name, code, type, value, usage count, status
- And I can filter by status and type
- And I can search by name or code

---

### US-024: View Coupon Usage History

**Description**: As a Merchant Admin, I want to view coupon usage history so that I can track campaign effectiveness.

**Acceptance Criteria**:
- Given I click on a coupon
- When I navigate to the "Usage" tab
- Then I see a list of orders where the coupon was used
- And each entry shows order ID, user, discount amount, timestamp
- And I can export the usage history

---

### US-025: Deactivate Coupon

**Description**: As a Merchant Admin, I want to deactivate a coupon so that it can no longer be used.

**Acceptance Criteria**:
- Given an active coupon
- When I click "Deactivate"
- Then the coupon status changes to "Inactive"
- And the coupon can no longer be redeemed
- Existing user coupons are not affected

---

### US-026: Calculate Cart Discount

**Description**: As a Shopper, I want to see applicable discounts on my cart so that I know my final price.

**Acceptance Criteria**:
- Given I have items in my cart totaling $150
- And there is an active threshold promotion: $100 save $15
- When I view my cart
- Then I see the discount applied: "-$15.00"
- And I see the promotion name: "Summer Sale"
- And I see the updated total: $135.00

---

### US-027: Apply Coupon to Cart

**Description**: As a Shopper, I want to apply a coupon code to my cart so that I can get additional discount.

**Acceptance Criteria**:
- Given I have items in my cart totaling $100
- And I have a valid coupon code "SAVE10" for $10 off
- When I enter the code and click "Apply"
- Then the coupon discount is shown: "-$10.00"
- And the total is updated to $90.00
- And the coupon code is shown as applied

---

### US-028: Apply Invalid Coupon

**Description**: As a Shopper, I want to receive clear feedback when a coupon is invalid so that I understand why it did not work.

**Acceptance Criteria**:
- Given I have items in my cart totaling $30
- And I have a coupon code "SAVE10" with minimum spend $50
- When I enter the code and click "Apply"
- Then I see an error message: "This coupon requires a minimum order of $50"
- And no discount is applied

---

### US-029: Apply Expired Coupon

**Description**: As a Shopper, I want to receive clear feedback when a coupon is expired so that I understand why it did not work.

**Acceptance Criteria**:
- Given I have a coupon code "OLDSALE" that expired yesterday
- When I enter the code and click "Apply"
- Then I see an error message: "This coupon has expired"
- And no discount is applied

---

### US-030: Apply Already Used Coupon

**Description**: As a Shopper, I want to receive feedback when a coupon has reached its usage limit so that I understand why it did not work.

**Acceptance Criteria**:
- Given I have a coupon code "LIMITED" with per_user_limit=1
- And I have already used this coupon once
- When I enter the code and click "Apply"
- Then I see an error message: "You have already used this coupon"
- And no discount is applied

---

### US-031: Remove Coupon from Cart

**Description**: As a Shopper, I want to remove an applied coupon from my cart so that I can use a different one.

**Acceptance Criteria**:
- Given I have applied coupon "SAVE10" to my cart
- When I click "Remove" on the coupon
- Then the coupon discount is removed
- And my cart total is updated to the original amount
- And I can now apply a different coupon

---

### US-032: View Discount Breakdown

**Description**: As a Shopper, I want to see a breakdown of all discounts applied to my order so that I understand my final price.

**Acceptance Criteria**:
- Given I have a promotion discount and a coupon applied
- When I view my cart or checkout
- Then I see an itemized breakdown:
  - Subtotal: $100.00
  - Promotion (Summer Sale): -$15.00
  - Coupon (SAVE10): -$10.00
  - Total: $75.00

---

### US-033: Promotion and Coupon Stack

**Description**: As a Shopper, I want promotions and coupons to both apply so that I get maximum savings.

**Acceptance Criteria**:
- Given I have a $100 cart
- And an active threshold promotion gives $10 off
- And I have a coupon for $5 off
- When I apply the coupon
- Then both discounts apply:
  - Promotion: -$10.00
  - Coupon: -$5.00
  - Total: $85.00

---

### US-034: Promotion Conflict Resolution by Priority

**Description**: As a Merchant Admin, I want conflicting promotions to resolve by priority so that only one applies.

**Acceptance Criteria**:
- Given two active promotions for the same product
- Promotion A: 20% off, priority 10
- Promotion B: $10 off, priority 5
- When a customer orders the product
- Then only Promotion A (20% off) is applied
- Because it has higher priority

---

### US-035: Promotion Scope Filtering at Checkout

**Description**: As a Shopper, I want only applicable promotions to show on my cart so that I see relevant discounts.

**Acceptance Criteria**:
- Given a promotion scoped to "Electronics" category
- And my cart has items from "Electronics" and "Clothing" categories
- When I view my cart
- Then only the Electronics items get the discount
- And the discount breakdown shows which items were discounted

---

### US-036: Issue Coupon to User

**Description**: As a Merchant Admin, I want to issue a coupon to a specific user so that I can provide personalized offers.

**Acceptance Criteria**:
- Given I have created a coupon
- When I navigate to "Issue Coupon"
- And I select a user by email or ID
- And I click "Issue"
- Then a user_coupon record is created
- And the user can see the coupon in their account

---

### US-037: View My Coupons

**Description**: As a Shopper, I want to view my available coupons so that I can use them at checkout.

**Acceptance Criteria**:
- Given I am logged in
- When I navigate to "My Coupons"
- Then I see all coupons issued to me
- And each coupon shows value, minimum spend, expiration date
- And I see status (Available, Used, Expired)

---

### US-038: Record Promotion Usage

**Description**: As the system, I want to record promotion usage when an order is placed so that we have accurate analytics.

**Acceptance Criteria**:
- Given an order is placed with a promotion applied
- When the order is confirmed
- Then a promotion_usage record is created
- And the record contains order_id, user_id, discount_amount
- And usage statistics are updated for the promotion

---

### US-039: Record Coupon Usage

**Description**: As the system, I want to record coupon usage when an order is placed so that limits are enforced.

**Acceptance Criteria**:
- Given an order is placed with a coupon
- When the order is confirmed
- Then the coupon's used_count is incremented
- And the user_coupon status is updated to "Used"
- And the used_at timestamp is set

---

### US-040: User Authentication for Admin Actions

**Description**: As a Merchant Admin, I must be authenticated to access promotion management so that only authorized users can modify promotions.

**Acceptance Criteria**:
- Given I am not logged in
- When I try to access /api/v1/promotions
- Then I receive a 401 Unauthorized error
- Given I am logged in as a tenant admin
- When I access /api/v1/promotions
- Then I receive the promotions list for my tenant only

---

### US-041: Tenant Isolation for Promotions

**Description**: As the system, I want to ensure promotions are isolated by tenant so that data is secure.

**Acceptance Criteria**:
- Given I am logged in as Tenant A admin
- When I request /api/v1/promotions
- Then I only see Tenant A's promotions
- And I cannot access Tenant B's promotions by ID
- And all queries are filtered by tenant_id

---

### US-042: Bulk Create Promotion Rules

**Description**: As a Merchant Admin, I want to add multiple rules at once when creating a tiered promotion so that setup is faster.

**Acceptance Criteria**:
- Given I am creating a threshold discount promotion
- When I add multiple tiers in one form
- And I click "Save"
- Then all rules are created atomically
- If any rule fails validation, none are saved

---

### US-043: Copy Promotion

**Description**: As a Merchant Admin, I want to duplicate an existing promotion so that I can quickly create similar promotions.

**Acceptance Criteria**:
- Given an existing promotion
- When I click "Duplicate"
- Then a new promotion is created with the same settings
- And the new promotion has status "Draft"
- And the name is appended with " (Copy)"

---

### US-044: Export Coupon Codes

**Description**: As a Merchant Admin, I want to export generated coupon codes so that I can distribute them.

**Acceptance Criteria**:
- Given I have generated batch coupon codes
- When I click "Export Codes"
- Then a CSV file is downloaded
- And the file contains all codes with their status

---

### US-045: Preview Promotion Impact

**Description**: As a Merchant Admin, I want to preview how a promotion affects specific products so that I can verify the discount.

**Acceptance Criteria**:
- Given I am creating or editing a promotion
- When I enter a product ID or select a product
- Then I see the original price and discounted price
- And I see the discount amount

---

### US-046: Search Promotions

**Description**: As a Merchant Admin, I want to search promotions by name so that I can find specific promotions quickly.

**Acceptance Criteria**:
- Given I am on the promotions list page
- When I type "Summer" in the search box
- Then I see only promotions with "Summer" in the name
- And results are filtered in real-time

---

### US-047: Filter Promotions by Status

**Description**: As a Merchant Admin, I want to filter promotions by status so that I can focus on active or draft promotions.

**Acceptance Criteria**:
- Given I am on the promotions list page
- When I select "Active" from the status filter
- Then I see only active promotions
- When I select "Draft"
- Then I see only draft promotions

---

### US-048: Filter Coupons by Status

**Description**: As a Merchant Admin, I want to filter coupons by status so that I can manage them effectively.

**Acceptance Criteria**:
- Given I am on the coupons list page
- When I select "Active" from the status filter
- Then I see only active coupons
- When I select "Inactive"
- Then I see only inactive coupons

---

### US-049: Automatic Promotion Status Update

**Description**: As the system, I want promotions to automatically update status based on time so that merchants don't have to manually manage them.

**Acceptance Criteria**:
- Given an active promotion with end_at in the past
- When the time passes end_at
- Then the promotion status changes to "Ended"
- And the promotion no longer applies to orders
- Given a draft promotion with start_at in the past and status is changed to active
- When the time is within the promotion period
- Then the promotion applies to orders

---

### US-050: Validate Promotion Currency

**Description**: As the system, I want to validate that promotion currency matches order currency so that discounts are correctly calculated.

**Acceptance Criteria**:
- Given a promotion with currency "USD"
- And an order with currency "EUR"
- When the system calculates discounts
- Then the promotion is not applied
- And no error is thrown (promotion is simply skipped)

---

### US-051: API Error Handling for Invalid Coupon

**Description**: As the system, I want to return clear error messages when coupon application fails so that users understand the issue.

**Acceptance Criteria**:
- Given an invalid coupon code is submitted
- When the API is called
- Then a 400 error is returned
- And the error message clearly states the reason:
  - "Coupon not found" for non-existent code
  - "Coupon expired" for past end_at
  - "Coupon not yet valid" for future start_at
  - "Coupon usage limit reached" for depleted coupon
  - "Minimum spend not met" for insufficient order amount

---

### US-052: Rate Limiting for Coupon Validation

**Description**: As the system, I want to rate limit coupon validation requests so that brute-force attacks are prevented.

**Acceptance Criteria**:
- Given a user makes more than 10 coupon validation requests per minute
- When the limit is exceeded
- Then a 429 Too Many Requests error is returned
- And the user must wait before trying again

---

### US-053: Audit Trail for Promotion Changes

**Description**: As a Merchant Admin, I want to see who changed a promotion and when so that I can track modifications.

**Acceptance Criteria**:
- Given I view a promotion's details
- When I look at the audit section
- Then I see created_by and created_at
- And I see updated_by and updated_at
- And I can see the last modifier's name

---

### US-054: Soft Delete Promotion

**Description**: As the system, I want promotions to be soft-deleted so that data can be recovered if needed.

**Acceptance Criteria**:
- Given I delete a promotion
- When the delete API is called
- Then deleted_at is set to current timestamp
- And the promotion no longer appears in normal queries
- But the data is preserved in the database
- And admin users can filter to see deleted items if needed

---

### US-055: Pagination for Large Lists

**Description**: As a Merchant Admin, I want promotion and coupon lists to be paginated so that the interface remains performant.

**Acceptance Criteria**:
- Given there are 500 promotions
- When I view the promotions list
- Then I see 20 promotions per page (default)
- And I can navigate to other pages
- And I can change the page size (max 100)
- And the API returns total count for pagination

---

### US-056: Promotion Activation Validation

**Description**: As the system, I want to validate promotion settings before activation so that incomplete promotions are not activated.

**Acceptance Criteria**:
- Given a promotion with no rules configured
- When I try to activate it
- Then an error is returned: "Cannot activate promotion without rules"
- Given a promotion with end_at before start_at
- When I try to activate it
- Then an error is returned: "End time must be after start time"

---

### US-057: Concurrent Coupon Redemption

**Description**: As the system, I want to handle concurrent coupon redemption correctly so that usage limits are not exceeded.

**Acceptance Criteria**:
- Given a coupon with total_count=1 remaining
- When two users try to redeem simultaneously
- Then only one redemption succeeds
- And the other receives an error: "Coupon usage limit reached"
- And the database transaction ensures consistency

---

### US-058: Discount Cannot Exceed Order Total

**Description**: As the system, I want to ensure discounts never exceed the order subtotal so that the total is never negative.

**Acceptance Criteria**:
- Given an order with subtotal $20
- And a coupon for $30 off
- When the coupon is applied
- Then the discount is capped at $20
- And the final total is $0
- And the order can be placed

---

### US-059: Real-time Promotion Eligibility Check

**Description**: As a Shopper, I want to see which promotions my cart is eligible for so that I know what discounts I can get.

**Acceptance Criteria**:
- Given I am on my cart page
- When I view the promotions section
- Then I see active promotions I qualify for
- And I see promotions I am $X away from qualifying for
- Example: "Spend $10 more to get $15 off!"

---

### US-060: Admin Dashboard Promotion Summary

**Description**: As a Merchant Admin, I want to see a summary of promotion activity on my dashboard so that I can monitor performance.

**Acceptance Criteria**:
- Given I am on the admin dashboard
- Then I see:
  - Number of active promotions
  - Number of active coupons
  - Total discount given this month
  - Top performing promotions
  - Coupon redemption rate