# Points System Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Points System PRD (Admin) |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-24 |
| Updated | 2026-03-24 |
| Author | Product Team |
| Module Location | admin/internal/domain/points/ |

---

## Product Overview

### Summary

This PRD describes the Points System module for the ShopJoy e-commerce SaaS platform. The Points System enables merchants to create and manage loyalty programs through points earning and redemption, driving customer engagement and repeat purchases.

**Core capabilities**:
- Points earning rules management (order payment, sign-in, product review, first order bonus)
- Points redemption for coupons
- User points account management (balance, transactions, manual adjustment)
- Points statistics and analytics

**Phase 1 scope**:
- Admin service only (shop-side earning logic deferred to Phase 2)
- All earn scenarios: Order payment, Sign-in/check-in, Product review, First order bonus
- Redemption limited to coupons only
- Event-driven auto-award + manual adjustment capability

### Problem Statement

Merchants on the ShopJoy platform currently lack a unified loyalty program system. They need:

- Ability to reward customers for purchases and engagement activities
- Flexible points calculation rules (tiered ratios, fixed amounts)
- Control over points expiration policies
- Visibility into customer points balances and transaction history
- Ability to manually adjust points for customer service scenarios

### Solution Overview

Implement a Points System bounded context within the admin service that provides:

1. **Earn Rules Management** - Create and configure points earning rules with tiered calculations
2. **Redeem Rules Management** - Configure points-to-coupon redemption options
3. **User Accounts Management** - View and manage customer points accounts
4. **Transaction History** - Immutable audit trail of all points movements
5. **Statistics Dashboard** - Insights into points distribution and usage

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | Increase customer retention | Repeat purchase rate | +20% within 6 months |
| BG-002 | Improve customer engagement | Points program adoption | 50% of active users earn points |
| BG-003 | Drive average order value | AOV for points users vs non-users | +15% lift |
| BG-004 | Reduce customer service overhead | Self-service points inquiry | 80% queries resolved without support |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | Create and manage earn rules | Merchant Admin |
| UG-002 | Configure points redemption options | Merchant Admin |
| UG-003 | View user points accounts and history | Merchant Admin / Customer Service |
| UG-004 | Manually adjust user points | Customer Service |
| UG-005 | View points statistics and trends | Merchant Admin / Finance |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | Shop-side points earning endpoints | Phase 2 feature (shop service) |
| NG-002 | Points redemption for products | Phase 2 feature |
| NG-003 | Points transfer between users | Phase 2 feature |
| NG-004 | Points tier/level system | Phase 2 feature |
| NG-005 | Points expiration notifications | Phase 2 feature |
| NG-006 | Referral bonus points | Phase 2 feature |
| NG-007 | Social sharing rewards | Phase 2 feature |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | Store owner or marketing manager | Create rules, view statistics, manage program |
| Customer Service | Customer support staff | View accounts, manual adjustments, resolve issues |
| Finance | Financial analyst | View transaction history, audit trail |

### Persona Details

#### Merchant Admin (Primary)

- **Demographics**: Small to medium business owners, marketing managers
- **Technical Proficiency**: Moderate
- **Goals**: Drive customer loyalty, increase repeat purchases, reward engagement
- **Pain Points**: Complex loyalty program setup, unclear points distribution
- **Frequency**: Weekly to monthly rule management, daily statistics review

#### Customer Service (Secondary)

- **Demographics**: Support team members
- **Technical Proficiency**: Moderate
- **Goals**: Resolve customer points inquiries, handle adjustment requests
- **Pain Points**: Lack of visibility into points history, manual processes
- **Frequency**: Daily account lookups and adjustments

#### Finance (Tertiary)

- **Demographics**: Finance and accounting staff
- **Technical Proficiency**: Moderate
- **Goals**: Audit points transactions, track liability
- **Pain Points**: No centralized transaction records
- **Frequency**: Weekly to monthly auditing

### Role-Based Access

| Role | Manage Rules | View Accounts | Adjust Points | View Statistics |
|------|--------------|---------------|---------------|-----------------|
| Tenant Admin | Full | Full | Full | Full |
| Tenant Marketing Manager | Full | Full | Limited | Full |
| Tenant Customer Service | Read-only | Full | Limited | Read-only |

---

## Functional Requirements

### Feature 1: Earn Rules Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | Create earn rule | Create points earning rule with calculation type | P0 |
| FR-002 | Calculation types | Support FIXED, RATIO, TIERED calculation types | P0 |
| FR-003 | Tiered configuration | Configure tiered ratios (e.g., 1 point/$1 up to $100, 1.5 points/$1 above) | P0 |
| FR-004 | Earn scenarios | Support ORDER_PAYMENT, SIGN_IN, PRODUCT_REVIEW, FIRST_ORDER scenarios | P0 |
| FR-005 | Condition types | Support NONE, NEW_USER, FIRST_ORDER, SPECIFIC_PRODUCTS, MIN_AMOUNT conditions | P0 |
| FR-006 | Points expiration | Configure expiration months per rule | P0 |
| FR-007 | Rule status | Draft, Active, Inactive states | P0 |
| FR-008 | Time range | Set start and end datetime for rules | P1 |
| FR-009 | Multiple active rules | Allow multiple active rules that stack | P0 |
| FR-010 | Rule priority | Set priority for display purposes | P1 |

### Feature 2: Redeem Rules Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-011 | Create redeem rule | Create points-to-coupon redemption rule | P0 |
| FR-012 | Coupon selection | Link rule to existing coupon template | P0 |
| FR-013 | Points required | Set points required for redemption | P0 |
| FR-014 | Stock management | Set redemption stock limit | P0 |
| FR-015 | Per-user limit | Set max redemptions per user | P0 |
| FR-016 | Rule status | Active, Inactive states | P0 |
| FR-017 | Time range | Set validity period | P1 |

### Feature 3: Points Accounts

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | Account list | List all user points accounts | P0 |
| FR-021 | Account detail | View single account with balance breakdown | P0 |
| FR-022 | Balance display | Show total balance, frozen balance, total earned, total redeemed | P0 |
| FR-023 | Account search | Search by user ID, email, phone | P0 |
| FR-024 | Manual adjustment | Add or deduct points with reason | P0 |
| FR-025 | Adjustment audit | Record adjustment reason and operator | P0 |

### Feature 4: Points Transactions

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-030 | Transaction list | View all points transactions | P0 |
| FR-031 | Transaction types | EARN, REDEEM, ADJUST, EXPIRE, FREEZE, UNFREEZE | P0 |
| FR-032 | Transaction filters | Filter by user, type, time range | P0 |
| FR-033 | Transaction detail | View transaction with reference info | P0 |
| FR-034 | Immutable records | Transactions cannot be modified or deleted | P0 |
| FR-035 | Balance snapshot | Record balance after each transaction | P0 |
| FR-036 | Expiration tracking | Track points expiration date per transaction | P0 |

### Feature 5: Points Redemptions

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-040 | Redemption list | View all redemption records | P0 |
| FR-041 | Redemption status | Pending, Completed, Cancelled states | P0 |
| FR-042 | Redemption filters | Filter by user, status, time range | P0 |
| FR-043 | Redemption detail | View points used, coupon issued, status | P0 |

### Feature 6: Statistics and Analytics

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-050 | Overview statistics | Total points issued, redeemed, expired, outstanding | P0 |
| FR-051 | Trend analysis | Daily/weekly/monthly points distribution trend | P1 |
| FR-052 | Top users | Top points earners by period | P1 |
| FR-053 | Rule performance | Points earned per rule | P1 |
| FR-054 | Redemption rate | Percentage of points redeemed vs issued | P1 |

---

## User Experience

### Entry Points

| Entry Point | Description | User Flow |
|-------------|-------------|-----------|
| Admin Sidebar | Points menu with sub-items | Dashboard > Points > Dashboard/Rules/Accounts |
| User Detail | Points balance quick view | Users > User Detail > Points Tab |
| Order Detail | Points earned indicator | Orders > Order Detail > Points Earned |

### Core Experience

#### Creating an Earn Rule (Merchant Admin)

1. Navigate to Points > Earn Rules
2. Click "Create Rule"
3. Select earn scenario (Order Payment, Sign-in, Product Review, First Order)
4. Configure calculation:
   - FIXED: Enter fixed points amount
   - RATIO: Enter points per currency unit
   - TIERED: Configure tier thresholds and ratios
5. Set conditions (NEW_USER, FIRST_ORDER, etc.)
6. Set expiration months (0 = no expiration)
7. Set time range and priority
8. Save as draft or activate immediately
9. View in rules list with status indicator

#### Configuring Tiered Points (Merchant Admin)

1. Select "Tiered" calculation type
2. Add tier rows:
   - Tier 1: Up to $100, 1 point per $1
   - Tier 2: $100-$500, 1.5 points per $1
   - Tier 3: Above $500, 2 points per $1
3. Preview calculation shows example totals
4. Save rule

#### Manual Points Adjustment (Customer Service)

1. Navigate to Points > Accounts
2. Search for user by email or ID
3. Click user to view account detail
4. Click "Adjust Points" button
5. Select adjustment type (Add/Deduct)
6. Enter points amount
7. Enter reason (required field)
8. Confirm adjustment
9. Transaction recorded immediately

#### Viewing Statistics (Merchant Admin)

1. Navigate to Points > Dashboard
2. See overview cards:
   - Total Points Issued
   - Total Points Redeemed
   - Total Points Expired
   - Outstanding Balance
3. View trend chart (daily/weekly toggle)
4. See top users table
5. Filter by time range

### Advanced Features

| Feature | Description |
|---------|-------------|
| Bulk Status Update | Activate/deactivate multiple rules |
| Rule Duplication | Copy existing rule as template |
| Export Transactions | Download transaction history as CSV |
| Expiration Preview | Preview points expiring in next 30 days |

### UI/UX Highlights

- Visual tiered configuration with drag-and-drop reordering
- Real-time points calculation preview
- Color-coded transaction types
- Account balance progress visualization
- Mobile-responsive admin interface
- Inline help for rule configuration

---

## Narrative

Merchant Alice wants to launch a loyalty program to increase customer retention. She navigates to Points > Earn Rules and clicks "Create Rule". She selects "Order Payment" scenario and chooses "Tiered" calculation. She configures three tiers: 1 point per $1 for orders up to $100, 1.5 points per $1 for orders between $100-$500, and 2 points per $1 for orders above $500. She sets points to expire in 12 months and activates the rule.

She also creates a "Sign-in" rule giving 5 points per day, with a limit of one sign-in per day. For product reviews, she creates a rule awarding 10 points per approved review.

After a week, Alice checks the Points Dashboard. She sees 50,000 points issued, 8,000 points redeemed, and 500 active users with points balances. The trend chart shows steady growth in points earning. She notices the top user has earned 2,500 points this month.

A customer contacts support about missing points from an order. Customer service rep Bob searches for the user's account, sees the transaction history, and confirms points were not awarded due to a rule condition. He manually adds 100 points with reason "Manual adjustment for order #12345 - rule condition not met" and the issue is resolved.

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Rule Creation Time | Average time to create an earn rule | <3 minutes |
| Account Query Time | Time to find user account | <30 seconds |
| Adjustment Resolution | CS adjustments per month | <5% of transactions |
| Program Adoption | % of tenants with active rules | 40% in 6 months |

### Business Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Points Redemption Rate | Points redeemed / Points issued | >30% |
| Repeat Purchase Lift | Repeat rate for points users | +20% |
| AOV Lift | AOV for orders with points earned | +15% |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API Response Time | Account balance query P95 | <100ms |
| Transaction Processing | Points earn/redemption P95 | <200ms |
| System Uptime | Points service availability | 99.9% |
| Error Rate | Failed transactions | <0.1% |

---

## Technical Considerations

### Integration Points

| Integration | Description | Method |
|-------------|-------------|--------|
| Order Service | Order payment events for earning | Domain Event |
| User Service | User account validation | Internal API |
| Coupon Service | Issue coupon on redemption | Internal API |
| Notification Service | Points notifications (Phase 2) | Domain Event |

### Data Storage and Privacy

| Aspect | Consideration |
|--------|---------------|
| Points Balance | Use BIGINT for integer points, no floating-point |
| Time Storage | All timestamps in UTC, stored as TIMESTAMP |
| Multi-Tenancy | All queries filtered by tenant_id |
| Soft Delete | Support deleted_at for rules and redemptions |
| Immutable Transactions | Transaction records never updated or deleted |
| PII Handling | User ID linking only, no direct PII in points module |

### Scalability and Performance

| Consideration | Strategy |
|---------------|----------|
| High Query Volume | Cache account balances (Redis) |
| Concurrent Transactions | Optimistic locking on account balance |
| Large Transaction History | Partition by tenant_id, time-based archiving |
| Rule Evaluation | Stateless calculation service |

### Potential Challenges

| Challenge | Mitigation |
|-----------|------------|
| Concurrent Balance Updates | Database row locking + retry logic |
| Points Expiration Batch | Scheduled job with batch processing |
| Large Tier Configuration | JSON field with validation |
| Rule Conflict Resolution | Clear priority system, rule stacking documentation |

---

## Database Schema

### points_accounts

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | | Tenant ID |
| user_id | BIGINT | NO | | User ID |
| balance | BIGINT | NO | 0 | Current available balance |
| frozen_balance | BIGINT | NO | 0 | Frozen/unavailable balance |
| total_earned | BIGINT | NO | 0 | Total points ever earned |
| total_redeemed | BIGINT | NO | 0 | Total points ever redeemed |
| total_expired | BIGINT | NO | 0 | Total points expired |
| version | INT | NO | 0 | Optimistic lock version |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation time (UTC) |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update time (UTC) |

**Indexes:**

- PRIMARY KEY (`id`)
- UNIQUE KEY `uk_tenant_user` (`tenant_id`, `user_id`)
- KEY `idx_user_id` (`user_id`)

### earn_rules

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | | Tenant ID |
| name | VARCHAR(200) | NO | | Rule name |
| description | TEXT | YES | NULL | Rule description |
| scenario | VARCHAR(32) | NO | | ORDER_PAYMENT, SIGN_IN, PRODUCT_REVIEW, FIRST_ORDER |
| calculation_type | VARCHAR(32) | NO | | FIXED, RATIO, TIERED |
| fixed_points | INT | NO | 0 | Fixed points for FIXED type |
| ratio | DECIMAL(10,4) | NO | 0 | Points per currency unit for RATIO type |
| tiers | JSON | YES | NULL | Tiered configuration [{threshold, ratio}] |
| condition_type | VARCHAR(32) | NO | 'NONE' | NONE, NEW_USER, FIRST_ORDER, SPECIFIC_PRODUCTS, MIN_AMOUNT |
| condition_value | JSON | YES | NULL | Condition configuration |
| expiration_months | INT | NO | 0 | Points expiration months, 0 = no expiration |
| status | TINYINT | NO | 0 | 0=draft, 1=active, 2=inactive |
| priority | INT | NO | 0 | Display priority |
| start_at | TIMESTAMP | YES | NULL | Rule start time (UTC) |
| end_at | TIMESTAMP | YES | NULL | Rule end time (UTC) |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation time (UTC) |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update time (UTC) |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | TIMESTAMP | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_tenant_status` (`tenant_id`, `status`)
- KEY `idx_scenario` (`scenario`)
- KEY `idx_deleted_at` (`deleted_at`)

### redeem_rules

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | | Tenant ID |
| name | VARCHAR(200) | NO | | Rule name |
| description | TEXT | YES | NULL | Rule description |
| coupon_id | BIGINT | NO | | Linked coupon template ID |
| points_required | INT | NO | | Points required for redemption |
| total_stock | INT | NO | 0 | Total redemption stock, 0 = unlimited |
| used_stock | INT | NO | 0 | Used stock count |
| per_user_limit | INT | NO | 0 | Max redemptions per user, 0 = unlimited |
| status | TINYINT | NO | 0 | 0=inactive, 1=active |
| start_at | TIMESTAMP | YES | NULL | Validity start (UTC) |
| end_at | TIMESTAMP | YES | NULL | Validity end (UTC) |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation time (UTC) |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update time (UTC) |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | TIMESTAMP | YES | NULL | Soft delete timestamp |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_tenant_status` (`tenant_id`, `status`)
- KEY `idx_coupon_id` (`coupon_id`)
- KEY `idx_deleted_at` (`deleted_at`)

### points_transactions

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | | Tenant ID |
| user_id | BIGINT | NO | | User ID |
| account_id | BIGINT | NO | | Points account ID |
| points | BIGINT | NO | | Points change (positive=earn, negative=deduct) |
| balance_after | BIGINT | NO | | Balance after transaction |
| type | VARCHAR(32) | NO | | EARN, REDEEM, ADJUST, EXPIRE, FREEZE, UNFREEZE |
| reference_type | VARCHAR(32) | NO | '' | ORDER, REDEEM_RULE, MANUAL, SYSTEM |
| reference_id | VARCHAR(64) | NO | '' | Related entity ID |
| earn_rule_id | BIGINT | YES | NULL | Earn rule ID if applicable |
| redeem_rule_id | BIGINT | YES | NULL | Redeem rule ID if applicable |
| description | VARCHAR(500) | NO | '' | Transaction description |
| expires_at | TIMESTAMP | YES | NULL | Points expiration time (UTC) |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation time (UTC) |
| created_by | BIGINT | NO | 0 | Operator ID (for manual adjustments) |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_user` (`tenant_id`, `user_id`)
- KEY `idx_account_id` (`account_id`)
- KEY `idx_type` (`type`)
- KEY `idx_created_at` (`created_at`)
- KEY `idx_expires_at` (`expires_at`)
- KEY `idx_reference` (`reference_type`, `reference_id`)

### points_redemptions

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | | Tenant ID |
| user_id | BIGINT | NO | | User ID |
| redeem_rule_id | BIGINT | NO | | Redeem rule ID |
| coupon_id | BIGINT | NO | | Coupon template ID |
| user_coupon_id | BIGINT | YES | NULL | Issued user coupon ID |
| points_used | INT | NO | | Points used for redemption |
| status | TINYINT | NO | 0 | 0=pending, 1=completed, 2=cancelled |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation time (UTC) |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update time (UTC) |
| completed_at | TIMESTAMP | YES | NULL | Completion time (UTC) |
| cancelled_at | TIMESTAMP | YES | NULL | Cancellation time (UTC) |
| cancel_reason | VARCHAR(255) | NO | '' | Cancellation reason |

**Indexes:**

- PRIMARY KEY (`id`)
- KEY `idx_tenant_user` (`tenant_id`, `user_id`)
- KEY `idx_redeem_rule` (`redeem_rule_id`)
- KEY `idx_status` (`status`)
- KEY `idx_created_at` (`created_at`)

---

## API Endpoints

### Earn Rules Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/points/earn-rules | Create earn rule |
| PUT | /api/v1/points/earn-rules/:id | Update earn rule |
| GET | /api/v1/points/earn-rules/:id | Get earn rule detail |
| GET | /api/v1/points/earn-rules | List earn rules |
| DELETE | /api/v1/points/earn-rules/:id | Delete earn rule |
| POST | /api/v1/points/earn-rules/:id/activate | Activate earn rule |
| POST | /api/v1/points/earn-rules/:id/deactivate | Deactivate earn rule |

### Redeem Rules Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/points/redeem-rules | Create redeem rule |
| PUT | /api/v1/points/redeem-rules/:id | Update redeem rule |
| GET | /api/v1/points/redeem-rules/:id | Get redeem rule detail |
| GET | /api/v1/points/redeem-rules | List redeem rules |
| DELETE | /api/v1/points/redeem-rules/:id | Delete redeem rule |
| POST | /api/v1/points/redeem-rules/:id/activate | Activate redeem rule |
| POST | /api/v1/points/redeem-rules/:id/deactivate | Deactivate redeem rule |

### Points Accounts

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/points/accounts | List points accounts |
| GET | /api/v1/points/accounts/:id | Get account detail |
| GET | /api/v1/points/accounts/:id/transactions | Get account transactions |
| POST | /api/v1/points/accounts/:id/adjust | Manual points adjustment |
| GET | /api/v1/points/accounts/by-user/:user_id | Get user's points account |

### Points Transactions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/points/transactions | List all transactions |
| GET | /api/v1/points/transactions/:id | Get transaction detail |

### Points Redemptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/points/redemptions | List redemptions |
| GET | /api/v1/points/redemptions/:id | Get redemption detail |

### Statistics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/points/stats | Get points statistics overview |
| GET | /api/v1/points/stats/trend | Get points trend data |
| GET | /api/v1/points/stats/top-users | Get top points earners |

### API Request/Response Examples

#### POST /api/v1/points/earn-rules - Create Earn Rule

**Request:**
```json
{
  "name": "Order Points Reward",
  "description": "Earn points on every order",
  "scenario": "ORDER_PAYMENT",
  "calculation_type": "TIERED",
  "tiers": [
    {"threshold": 10000, "ratio": "1.0"},
    {"threshold": 50000, "ratio": "1.5"},
    {"threshold": null, "ratio": "2.0"}
  ],
  "condition_type": "NONE",
  "expiration_months": 12,
  "status": 1,
  "priority": 10,
  "start_at": "2026-04-01T00:00:00Z",
  "end_at": "2026-12-31T23:59:59Z"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Order Points Reward",
  "scenario": "ORDER_PAYMENT",
  "calculation_type": "TIERED",
  "status": 1,
  "status_text": "Active",
  "created_at": "2026-03-24T10:00:00Z"
}
```

#### GET /api/v1/points/accounts/:id - Get Account Detail

**Response:**
```json
{
  "id": 1,
  "user_id": 12345,
  "balance": 5000,
  "frozen_balance": 0,
  "total_earned": 10000,
  "total_redeemed": 4500,
  "total_expired": 500,
  "created_at": "2026-01-15T08:00:00Z",
  "updated_at": "2026-03-24T10:30:00Z"
}
```

#### POST /api/v1/points/accounts/:id/adjust - Manual Adjustment

**Request:**
```json
{
  "adjustment_type": "ADD",
  "points": 100,
  "reason": "Compensation for delayed order #12345"
}
```

**Response:**
```json
{
  "transaction_id": 10001,
  "points": 100,
  "balance_after": 5100,
  "created_at": "2026-03-24T10:35:00Z"
}
```

#### GET /api/v1/points/stats - Statistics Overview

**Response:**
```json
{
  "total_issued": 1000000,
  "total_redeemed": 300000,
  "total_expired": 50000,
  "outstanding_balance": 650000,
  "redemption_rate": "30.0%",
  "active_users": 1500,
  "period_start": "2026-03-01T00:00:00Z",
  "period_end": "2026-03-24T23:59:59Z"
}
```

---

## Business Rules

### Earn Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| ER-001 | Calculation base | Points calculated on order subtotal (before discounts) |
| ER-002 | Tiered calculation | Uses incremental calculation within each tier |
| ER-003 | Rule stacking | Multiple active rules can stack, points are cumulative |
| ER-004 | NEW_USER definition | User account created within 7 days |
| ER-005 | FIRST_ORDER definition | User has no prior paid orders |
| ER-006 | Expiration | expiration_months > 0 means points expire X months from earn date |
| ER-007 | No expiration | expiration_months = 0 means points never expire |
| ER-008 | Calculation precision | Points are integers, round down for fractional results |
| ER-009 | Minimum earn | Minimum 1 point per qualifying transaction |

### Redeem Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| RR-001 | Balance check | User must have sufficient available balance |
| RR-002 | Stock limit | Redemption fails if stock exhausted |
| RR-003 | Per-user limit | User cannot exceed per_user_limit redemptions |
| RR-004 | Time validity | Redemption only valid within start_at and end_at |
| RR-005 | Coupon issuance | Coupon auto-issued to user on successful redemption |
| RR-006 | Points deduction | Points deducted immediately upon redemption |
| RR-007 | Non-refundable | Redeemed points cannot be returned |

### Account Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| AR-001 | Single account | One points account per user per tenant |
| AR-002 | Non-negative balance | Balance cannot go below zero |
| AR-003 | FIFO expiration | Oldest points expire first |
| AR-004 | Frozen balance | Frozen points cannot be used for redemption |
| AR-005 | Manual adjustment limits | Adjustment requires reason, logged with operator |

### Transaction Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| TR-001 | Immutable records | Transaction records cannot be modified or deleted |
| TR-002 | Balance consistency | balance_after must equal previous balance + points |
| TR-003 | Reference tracking | All transactions linked to source (order, rule, manual) |
| TR-004 | Timestamp UTC | All timestamps stored and displayed in UTC |

---

## Error Codes

Points module uses 220xxx error code range:

| Constant | HTTP Status | Code | Message |
|----------|-------------|------|---------|
| ErrPointsAccountNotFound | 404 | 220001 | points account not found |
| ErrPointsInsufficientBalance | 400 | 220002 | insufficient points balance |
| ErrPointsEarnRuleNotFound | 404 | 220003 | earn rule not found |
| ErrPointsRedeemRuleNotFound | 404 | 220004 | redeem rule not found |
| ErrPointsTransactionNotFound | 404 | 220005 | points transaction not found |
| ErrPointsRedemptionNotFound | 404 | 220006 | points redemption not found |
| ErrPointsEarnRuleInactive | 400 | 220007 | earn rule is inactive |
| ErrPointsRedeemRuleInactive | 400 | 220008 | redeem rule is inactive |
| ErrPointsRedemptionStockExhausted | 400 | 220009 | redemption stock exhausted |
| ErrPointsRedemptionLimitExceeded | 400 | 220010 | user redemption limit exceeded |
| ErrPointsInvalidAdjustment | 400 | 220011 | invalid points adjustment |
| ErrPointsAdjustmentReasonRequired | 400 | 220012 | adjustment reason is required |
| ErrPointsInvalidCalculationType | 400 | 220013 | invalid calculation type |
| ErrPointsInvalidScenario | 400 | 220014 | invalid earn scenario |
| ErrPointsInvalidCondition | 400 | 220015 | invalid condition type |
| ErrPointsTierConfigInvalid | 400 | 220016 | tiered configuration is invalid |
| ErrPointsCouponNotFound | 404 | 220017 | linked coupon not found |
| ErrPointsCouponInactive | 400 | 220018 | linked coupon is inactive |
| ErrPointsRuleTimeInvalid | 400 | 220019 | rule time range is invalid |
| ErrPointsCannotDeleteActiveRule | 400 | 220020 | cannot delete active rule |
| ErrPointsFrozenBalance | 400 | 220021 | points are frozen |
| ErrPointsVersionConflict | 409 | 220022 | account version conflict, please retry |

---

## User Stories

### US-001: Create Fixed Points Earn Rule

**Description**: As a Merchant Admin, I want to create a fixed points earn rule so that users earn a set amount of points for specific actions.

**Acceptance Criteria**:
- Given I am logged in as a Merchant Admin
- When I navigate to Points > Earn Rules and click "Create Rule"
- And I select scenario "Sign-in"
- And I select calculation type "Fixed"
- And I enter fixed points "5"
- And I set expiration months "12"
- And I save the rule
- Then the rule is created with status "Draft"
- And I can view the rule in the rules list

---

### US-002: Create Ratio Points Earn Rule

**Description**: As a Merchant Admin, I want to create a ratio-based earn rule so that points are proportional to order amount.

**Acceptance Criteria**:
- Given I am creating an earn rule
- When I select scenario "Order Payment"
- And I select calculation type "Ratio"
- And I enter ratio "1.5" (points per currency unit)
- And I save the rule
- Then the rule is created successfully
- When an order of $100 is placed
- Then 150 points are calculated

---

### US-003: Create Tiered Points Earn Rule

**Description**: As a Merchant Admin, I want to create tiered earn rules so that higher spending earns more points per unit.

**Acceptance Criteria**:
- Given I am creating an earn rule
- When I select calculation type "Tiered"
- And I configure tiers:
  - Up to $100: 1 point per $1
  - $100 to $500: 1.5 points per $1
  - Above $500: 2 points per $1
- And I save the rule
- Then the rule is created with tier configuration
- When an order of $200 is placed
- Then points are calculated as: (100 * 1) + (100 * 1.5) = 250 points

---

### US-004: Set Earn Rule Condition

**Description**: As a Merchant Admin, I want to set conditions on earn rules so that points are only awarded to qualifying users.

**Acceptance Criteria**:
- Given I am creating an earn rule
- When I select condition type "NEW_USER"
- And I save the rule
- Then only users created within 7 days qualify
- When a user created 5 days ago makes a qualifying purchase
- Then points are awarded
- When a user created 10 days ago makes a qualifying purchase
- Then no points are awarded

---

### US-005: Set Points Expiration

**Description**: As a Merchant Admin, I want to configure points expiration so that points don't remain valid indefinitely.

**Acceptance Criteria**:
- Given I am creating an earn rule
- When I set expiration months "12"
- And I save the rule
- Then points earned under this rule expire 12 months from earn date
- When I set expiration months "0"
- Then points never expire

---

### US-006: Activate Earn Rule

**Description**: As a Merchant Admin, I want to activate an earn rule so that points start being awarded.

**Acceptance Criteria**:
- Given an earn rule with status "Draft"
- When I click "Activate"
- Then the rule status changes to "Active"
- And the rule is now used for points calculation
- If current time is within the rule's time range, points are awarded immediately

---

### US-007: Deactivate Earn Rule

**Description**: As a Merchant Admin, I want to deactivate an earn rule so that points stop being awarded.

**Acceptance Criteria**:
- Given an active earn rule
- When I click "Deactivate"
- Then the rule status changes to "Inactive"
- And the rule no longer awards points
- Existing points already earned are not affected

---

### US-008: Create Redeem Rule

**Description**: As a Merchant Admin, I want to create a redeem rule so that users can exchange points for coupons.

**Acceptance Criteria**:
- Given I navigate to Points > Redeem Rules
- When I click "Create Rule"
- And I select a coupon template
- And I set points required "100"
- And I set total stock "1000"
- And I set per-user limit "5"
- And I save the rule
- Then the redeem rule is created
- And users can redeem 100 points for the coupon

---

### US-009: Set Redemption Stock Limit

**Description**: As a Merchant Admin, I want to limit redemption stock so that I can control costs.

**Acceptance Criteria**:
- Given I am creating a redeem rule
- When I set total stock "100"
- And the rule is activated
- When 100 redemptions have been made
- Then further redemption attempts fail with "Stock exhausted"
- When I set total stock "0"
- Then stock is unlimited

---

### US-010: Set Per-User Redemption Limit

**Description**: As a Merchant Admin, I want to limit redemptions per user so that points are distributed fairly.

**Acceptance Criteria**:
- Given a redeem rule with per-user limit "3"
- When a user redeems 3 times
- Then further redemption attempts by that user fail with "Limit exceeded"
- Other users can still redeem if they haven't reached their limit

---

### US-011: View Points Accounts List

**Description**: As a Merchant Admin, I want to view all points accounts so that I can see user balances.

**Acceptance Criteria**:
- Given I navigate to Points > Accounts
- Then I see a list of all user points accounts
- And each account shows user ID, balance, total earned, total redeemed
- And I can sort by balance or creation date
- And I can paginate through results

---

### US-012: Search Points Account

**Description**: As a Customer Service rep, I want to search for a user's points account so that I can assist them.

**Acceptance Criteria**:
- Given I am on the Accounts page
- When I enter a user ID, email, or phone number
- And I click "Search"
- Then I see matching accounts
- And I can click to view the account detail

---

### US-013: View Account Detail

**Description**: As a Customer Service rep, I want to view a user's account details so that I can see their points status.

**Acceptance Criteria**:
- Given I click on a points account
- Then I see:
  - Current balance
  - Frozen balance
  - Total earned
  - Total redeemed
  - Total expired
  - Recent transactions

---

### US-014: View Transaction History

**Description**: As a Customer Service rep, I want to view a user's transaction history so that I can explain their points movements.

**Acceptance Criteria**:
- Given I am viewing an account detail
- When I scroll to the transactions section
- Then I see a list of recent transactions
- Each transaction shows:
  - Points change (+/-)
  - Transaction type
  - Reference (order ID, rule name, etc.)
  - Timestamp
  - Expiration date (if applicable)
- And I can load more transactions

---

### US-015: Manually Add Points

**Description**: As a Customer Service rep, I want to manually add points to a user's account to resolve issues.

**Acceptance Criteria**:
- Given I am viewing an account detail
- When I click "Adjust Points"
- And I select "Add Points"
- And I enter points "100"
- And I enter reason "Compensation for system error on order #12345"
- And I click "Confirm"
- Then 100 points are added to the account
- And a transaction record is created with type "ADJUST"
- And the operator ID is recorded

---

### US-016: Manually Deduct Points

**Description**: As a Customer Service rep, I want to manually deduct points to correct errors.

**Acceptance Criteria**:
- Given I am viewing an account detail
- And the account has sufficient balance
- When I click "Adjust Points"
- And I select "Deduct Points"
- And I enter points "50"
- And I enter reason "Duplicate points award correction"
- And I click "Confirm"
- Then 50 points are deducted from the account
- And a transaction record is created

---

### US-017: Prevent Negative Balance

**Description**: As the system, I want to prevent balance from going negative so that data integrity is maintained.

**Acceptance Criteria**:
- Given a user has 50 points balance
- When a deduction of 100 points is attempted
- Then the operation fails with "Insufficient balance"
- And the account balance remains unchanged

---

### US-018: View Statistics Dashboard

**Description**: As a Merchant Admin, I want to view points statistics so that I can monitor program health.

**Acceptance Criteria**:
- Given I navigate to Points > Dashboard
- Then I see overview cards:
  - Total Points Issued
  - Total Points Redeemed
  - Total Points Expired
  - Outstanding Balance
  - Redemption Rate
  - Active Users with Points

---

### US-019: View Points Trend

**Description**: As a Merchant Admin, I want to view points trend over time so that I can understand patterns.

**Acceptance Criteria**:
- Given I am on the Dashboard
- When I view the trend chart
- Then I see a line chart showing daily points:
  - Points issued (green line)
  - Points redeemed (blue line)
  - Points expired (red line)
- And I can toggle between daily/weekly view
- And I can select date range

---

### US-020: View Top Users

**Description**: As a Merchant Admin, I want to see top points earners so that I can identify engaged customers.

**Acceptance Criteria**:
- Given I am on the Dashboard
- When I view the Top Users section
- Then I see a table of top 10 users by:
  - Points earned in period
  - User ID
  - Account creation date
- And I can select the time period

---

### US-021: View Redemptions List

**Description**: As a Merchant Admin, I want to view all redemptions so that I can track coupon issuance.

**Acceptance Criteria**:
- Given I navigate to Points > Redemptions
- Then I see a list of all redemption records
- Each record shows:
  - User ID
  - Redeem rule name
  - Points used
  - Coupon issued
  - Status
  - Timestamp
- And I can filter by status and date range

---

### US-022: View Redemption Detail

**Description**: As a Merchant Admin, I want to view redemption details so that I can verify the transaction.

**Acceptance Criteria**:
- Given I click on a redemption record
- Then I see:
  - User information
  - Redeem rule details
  - Points deducted
  - Coupon template information
  - User coupon ID (issued coupon)
  - Status and timestamps

---

### US-023: Multiple Rules Stack

**Description**: As the system, I want active rules to stack so that users can earn points from multiple sources.

**Acceptance Criteria**:
- Given two active earn rules:
  - Rule A: Order payment gives 1 point per $1
  - Rule B: First order gives 100 bonus points
- When a new user places their first order of $50
- Then they earn: (50 * 1) + 100 = 150 points
- And two separate transaction records are created

---

### US-024: FIFO Expiration

**Description**: As the system, I want to expire oldest points first so that users keep their most recently earned points.

**Acceptance Criteria**:
- Given a user has:
  - 100 points earned on 2025-01-01, expiring 2026-01-01
  - 50 points earned on 2025-06-01, expiring 2026-06-01
- When 100 points are redeemed on 2025-12-01
- Then the 100 points from 2025-01-01 are deducted first
- And the 50 points from 2025-06-01 remain

---

### US-025: Concurrent Balance Update

**Description**: As the system, I want to handle concurrent balance updates correctly so that data is consistent.

**Acceptance Criteria**:
- Given a user has 100 points balance (version 1)
- When two concurrent requests attempt:
  - Request A: Add 50 points (version 1)
  - Request B: Redeem 30 points (version 1)
- Then one request succeeds and updates version to 2
- And the other request fails with "Version conflict"
- And the client can retry with the new version

---

### US-026: Transaction Audit Trail

**Description**: As a Finance user, I want a complete audit trail of points transactions so that I can verify program integrity.

**Acceptance Criteria**:
- Given I query points transactions
- Then each record shows:
  - Transaction ID
  - User ID
  - Points change
  - Balance after
  - Transaction type
  - Reference (order/rule/manual)
  - Operator (if manual)
  - Timestamp
- And records are immutable

---

### US-027: Filter Transactions by Type

**Description**: As a Merchant Admin, I want to filter transactions by type so that I can focus on specific activities.

**Acceptance Criteria**:
- Given I am on the Transactions page
- When I select type filter "REDEEM"
- Then I see only redemption transactions
- When I select type filter "ADJUST"
- Then I see only manual adjustment transactions

---

### US-028: Filter Transactions by Date Range

**Description**: As a Merchant Admin, I want to filter transactions by date range so that I can review specific periods.

**Acceptance Criteria**:
- Given I am on the Transactions page
- When I set date range "2026-03-01" to "2026-03-31"
- Then I see only transactions created within that period
- And the total count reflects the filtered results

---

### US-029: Delete Draft Earn Rule

**Description**: As a Merchant Admin, I want to delete a draft earn rule so that I can remove unused rules.

**Acceptance Criteria**:
- Given an earn rule with status "Draft"
- When I click "Delete"
- Then the rule is soft-deleted
- And it no longer appears in the rules list
- Given an active earn rule
- When I try to delete it
- Then an error is returned: "Cannot delete active rule"

---

### US-030: Edit Earn Rule

**Description**: As a Merchant Admin, I want to edit an existing earn rule so that I can adjust its settings.

**Acceptance Criteria**:
- Given an existing earn rule
- When I click "Edit"
- Then I can modify name, description, calculation, conditions, expiration
- And changes are saved successfully
- If the rule is active, changes take effect immediately for new transactions

---

### US-031: User Authentication for Admin Actions

**Description**: As a Merchant Admin, I must be authenticated to access points management so that only authorized users can modify rules.

**Acceptance Criteria**:
- Given I am not logged in
- When I try to access /api/v1/points/earn-rules
- Then I receive a 401 Unauthorized error
- Given I am logged in as a tenant admin
- When I access /api/v1/points/earn-rules
- Then I receive the earn rules list for my tenant only

---

### US-032: Tenant Isolation for Points

**Description**: As the system, I want to ensure points data is isolated by tenant so that data is secure.

**Acceptance Criteria**:
- Given I am logged in as Tenant A admin
- When I request /api/v1/points/accounts
- Then I only see Tenant A's accounts
- And I cannot access Tenant B's accounts by ID
- And all queries are filtered by tenant_id

---

### US-033: Validate Tiered Configuration

**Description**: As the system, I want to validate tiered configuration so that invalid rules are not saved.

**Acceptance Criteria**:
- Given I am creating a tiered earn rule
- When I enter tiers with overlapping thresholds
- Then an error is returned: "Invalid tier configuration"
- When I enter tiers with missing required fields
- Then an error is returned with field details
- When I enter valid tiers with ascending thresholds
- Then the rule is saved successfully

---

### US-034: Calculate Tiered Points

**Description**: As the system, I want to calculate tiered points correctly so that users receive accurate points.

**Acceptance Criteria**:
- Given a tiered rule:
  - Tier 1: Up to $100, 1 point/$
  - Tier 2: $100-$500, 1.5 points/$
  - Tier 3: Above $500, 2 points/$
- When an order of $600 is placed
- Then points = (100 * 1) + (400 * 1.5) + (100 * 2) = 100 + 600 + 200 = 900 points
- When an order of $50 is placed
- Then points = 50 * 1 = 50 points

---

### US-035: Coupon Auto-Issued on Redemption

**Description**: As the system, I want to automatically issue a coupon when points are redeemed so that users receive their reward.

**Acceptance Criteria**:
- Given a user redeems points for a coupon
- When the redemption is successful
- Then a user_coupon record is created
- And the user_coupon_id is linked to the redemption
- And the user can use the coupon on their next order

---

### US-036: Redeem Rule Time Validity

**Description**: As the system, I want to validate redeem rule time range so that expired rules are not used.

**Acceptance Criteria**:
- Given a redeem rule with end_at in the past
- When a user tries to redeem
- Then an error is returned: "Redeem rule is no longer valid"
- Given a redeem rule with start_at in the future
- When a user tries to redeem
- Then an error is returned: "Redeem rule is not yet active"

---

### US-037: Account Version for Concurrency

**Description**: As the system, I want to use optimistic locking on accounts so that concurrent updates are handled correctly.

**Acceptance Criteria**:
- Given an account with version 1
- When a transaction is processed
- Then version is incremented to 2
- When another transaction attempts with version 1
- Then it fails with "Version conflict, please retry"
- And the transaction can be retried with current version

---

### US-038: Points Rounding

**Description**: As the system, I want to round down fractional points so that points are always integers.

**Acceptance Criteria**:
- Given a ratio rule of 1.5 points per $1
- When an order of $99 is placed
- Then points = floor(99 * 1.5) = floor(148.5) = 148 points
- And no fractional points are stored

---

### US-039: Export Transactions

**Description**: As a Finance user, I want to export transaction history so that I can perform offline analysis.

**Acceptance Criteria**:
- Given I am on the Transactions page
- When I click "Export" and select date range
- Then a CSV file is downloaded
- And the file contains all transactions within the range
- And maximum export is 10,000 records

---

### US-040: View Expiring Points

**Description**: As a Merchant Admin, I want to see points expiring soon so that I can plan retention campaigns.

**Acceptance Criteria**:
- Given I am on the Dashboard
- When I view the "Expiring Soon" section
- Then I see total points expiring in the next 30 days
- And I can see a breakdown by expiration date

---

### US-041: First Order Bonus

**Description**: As a Merchant Admin, I want to give bonus points on first order so that new customers are encouraged.

**Acceptance Criteria**:
- Given I create an earn rule with scenario "FIRST_ORDER"
- And I set fixed points "200"
- When a user places their first paid order
- Then they receive 200 bonus points
- When the same user places another order
- Then no bonus points are given (not first order anymore)

---

### US-042: Product Review Points

**Description**: As a Merchant Admin, I want to reward points for product reviews so that customers are motivated to share feedback.

**Acceptance Criteria**:
- Given I create an earn rule with scenario "PRODUCT_REVIEW"
- And I set fixed points "10"
- When a user submits an approved product review
- Then they receive 10 points
- When the same user reviews the same product again
- Then no additional points (one review per product per user)

---

### US-043: Sign-in Daily Limit

**Description**: As the system, I want to limit sign-in points to once per day so that the system is not abused.

**Acceptance Criteria**:
- Given a sign-in earn rule with 5 fixed points
- When a user signs in today
- Then they receive 5 points
- When the same user signs in again today
- Then no additional points
- When the same user signs in tomorrow
- Then they receive another 5 points

---

### US-044: Specific Products Condition

**Description**: As a Merchant Admin, I want to limit points to specific products so that I can promote certain items.

**Acceptance Criteria**:
- Given I create an earn rule with condition "SPECIFIC_PRODUCTS"
- And I select products [101, 102, 103]
- When a user orders product 101
- Then they earn points under this rule
- When a user orders product 201 (not in list)
- Then they do not earn points under this rule

---

### US-045: Minimum Amount Condition

**Description**: As a Merchant Admin, I want to set a minimum order amount for points so that small orders don't qualify.

**Acceptance Criteria**:
- Given I create an earn rule with condition "MIN_AMOUNT"
- And I set minimum amount "$50"
- When a user places an order of $60
- Then they earn points under this rule
- When a user places an order of $40
- Then they do not earn points under this rule

---

### US-046: Pagination for Large Lists

**Description**: As a Merchant Admin, I want paginated lists so that the interface remains performant.

**Acceptance Criteria**:
- Given there are 500 points accounts
- When I view the accounts list
- Then I see 20 accounts per page (default)
- And I can navigate to other pages
- And I can change page size (max 100)
- And the API returns total count for pagination

---

### US-047: Rule Priority Display

**Description**: As a Merchant Admin, I want to see rule priority so that I understand rule ordering.

**Acceptance Criteria**:
- Given multiple earn rules exist
- When I view the rules list
- Then rules are sorted by priority (descending)
- And I can see the priority number for each rule
- And I can filter by priority range

---

### US-048: Duplicate Earn Rule

**Description**: As a Merchant Admin, I want to duplicate an existing earn rule so that I can quickly create similar rules.

**Acceptance Criteria**:
- Given an existing earn rule
- When I click "Duplicate"
- Then a new rule is created with the same settings
- And the new rule has status "Draft"
- And the name is appended with " (Copy)"

---

### US-049: Bulk Activate Rules

**Description**: As a Merchant Admin, I want to activate multiple rules at once so that I can efficiently manage rules.

**Acceptance Criteria**:
- Given multiple draft earn rules selected
- When I click "Activate Selected"
- Then all selected rules change to "Active" status
- And each rule's activation is logged

---

### US-050: View Rule Performance

**Description**: As a Merchant Admin, I want to see points earned per rule so that I can evaluate rule effectiveness.

**Acceptance Criteria**:
- Given I am viewing an earn rule detail
- When I look at the performance section
- Then I see:
  - Total points awarded under this rule
  - Number of qualifying transactions
  - Average points per transaction
  - Period-over-period comparison

---

## Frontend Structure

### New Pages

```
src/views/points/
├── dashboard/
│   └── index.vue                    # Points statistics dashboard
├── earn-rules/
│   ├── index.vue                    # Earn rules list
│   ├── [id].vue                     # Earn rule detail
│   └── components/
│       ├── TieredConfig.vue         # Tiered configuration component
│       ├── RuleForm.vue             # Rule creation/edit form
│       └── ConditionConfig.vue      # Condition configuration
├── redeem-rules/
│   ├── index.vue                    # Redeem rules list
│   ├── [id].vue                     # Redeem rule detail
│   └── components/
│       ├── CouponSelector.vue       # Coupon selection component
│       └── RedeemRuleForm.vue       # Redeem rule form
├── accounts/
│   ├── index.vue                    # Points accounts list
│   └── [id].vue                     # Account detail with transactions
├── transactions/
│   └── index.vue                    # All transactions list
└── redemptions/
    ├── index.vue                    # Redemptions list
    └── [id].vue                     # Redemption detail
```

### Component Specifications

#### TieredConfig.vue

- Dynamic tier row management (add/remove)
- Drag-and-drop reordering
- Threshold input (currency)
- Ratio input (decimal)
- Preview calculation section
- Validation for overlapping thresholds

#### CouponSelector.vue

- Search coupon templates
- Display coupon details (type, value, status)
- Select and link coupon to redeem rule
- Show coupon preview card

#### RuleForm.vue

- Rule name and description
- Scenario selection (dropdown)
- Calculation type tabs (Fixed/Ratio/Tiered)
- Condition configuration section
- Expiration months input
- Time range picker
- Priority input
- Status toggle

#### AccountDetail.vue

- Balance breakdown card
- Statistics mini cards (earned, redeemed, expired)
- Transaction history table with filters
- Manual adjustment button/dialog
- Quick action buttons

#### PointsStatsCard.vue

- Key metrics display
- Period comparison arrows
- Mini sparkline chart
- Color-coded indicators

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Team Size | Description |
|-------|----------|-----------|-------------|
| Phase 1: Database & Domain | 2 days | 2 backend | Schema design, entities, repositories |
| Phase 2: Earn Rules API | 2 days | 2 backend | CRUD, activation, calculation logic |
| Phase 3: Redeem Rules API | 1.5 days | 2 backend | CRUD, redemption flow |
| Phase 4: Accounts & Transactions | 2 days | 2 backend | Account management, manual adjustment |
| Phase 5: Statistics API | 1 day | 1 backend | Stats, trends, top users |
| Phase 6: Frontend Pages | 4 days | 2 frontend | All admin pages and components |
| Phase 7: Testing & Polish | 2 days | Full team | E2E testing, bug fixes |

**Total Estimate: 14.5 working days**

### Suggested Phases

#### Phase 1: Foundation (Day 1-2)

- Design and create database schema
- Implement domain entities (EarnRule, RedeemRule, PointsAccount, PointsTransaction, PointsRedemption)
- Create repository interfaces and implementations
- Set up error codes (220xxx range)
- Write database migrations

#### Phase 2: Earn Rules (Day 3-4)

- Implement EarnRuleService
- Create CRUD APIs for earn rules
- Implement calculation logic (Fixed, Ratio, Tiered)
- Add condition validation logic
- Implement rule activation/deactivation

#### Phase 3: Redeem Rules (Day 5-6)

- Implement RedeemRuleService
- Create CRUD APIs for redeem rules
- Implement redemption flow with coupon issuance
- Add stock and per-user limit validation
- Handle redemption status transitions

#### Phase 4: Accounts & Transactions (Day 7-8)

- Implement PointsAccountService
- Create account listing and detail APIs
- Implement transaction history APIs
- Add manual adjustment capability
- Implement optimistic locking for concurrency

#### Phase 5: Statistics (Day 9)

- Implement statistics calculation service
- Create overview statistics API
- Implement trend data API
- Add top users query
- Create expiration preview

#### Phase 6: Frontend (Day 10-13)

- Build Points Dashboard page
- Create Earn Rules management pages
- Create Redeem Rules management pages
- Build Accounts list and detail pages
- Implement Transactions list page
- Create Redemptions list and detail pages
- Add TieredConfig component with drag-and-drop

#### Phase 7: Testing (Day 14-15)

- Write unit tests for services
- Create integration tests for APIs
- Perform E2E testing
- Fix identified bugs
- Performance testing
- Documentation review

---

## Appendix

### Earn Scenario Types

| Code | Name | Description | Trigger |
|------|------|-------------|---------|
| ORDER_PAYMENT | Order Payment | Points earned on successful order payment | Order paid event |
| SIGN_IN | Sign-in | Points for daily sign-in | User sign-in action |
| PRODUCT_REVIEW | Product Review | Points for approved product review | Review approved event |
| FIRST_ORDER | First Order Bonus | One-time bonus for first order | First paid order |

### Calculation Types

| Code | Name | Description | Example |
|------|------|-------------|---------|
| FIXED | Fixed | Fixed points per occurrence | 10 points per review |
| RATIO | Ratio | Points proportional to amount | 1.5 points per $1 |
| TIERED | Tiered | Different ratios by threshold | 1 point/$1 up to $100, 1.5 points/$1 above |

### Condition Types

| Code | Name | Description | Condition Value |
|------|------|-------------|-----------------|
| NONE | No Condition | All users qualify | null |
| NEW_USER | New User | User created within 7 days | null |
| FIRST_ORDER | First Order | User's first paid order | null |
| SPECIFIC_PRODUCTS | Specific Products | Only selected products | {"product_ids": [1, 2, 3]} |
| MIN_AMOUNT | Minimum Amount | Order subtotal meets minimum | {"min_amount": 5000} |

### Transaction Types

| Code | Name | Description | Points Direction |
|------|------|-------------|------------------|
| EARN | Earn | Points earned from rule | Positive (+) |
| REDEEM | Redeem | Points used for redemption | Negative (-) |
| ADJUST | Adjust | Manual adjustment | Positive or Negative |
| EXPIRE | Expire | Points expired | Negative (-) |
| FREEZE | Freeze | Points frozen/unavailable | Negative (-) from balance |
| UNFREEZE | Unfreeze | Points unfrozen | Positive (+) to balance |

### Redemption Status

| Status | Value | Description |
|--------|-------|-------------|
| pending | 0 | Redemption created, coupon not yet issued |
| completed | 1 | Coupon successfully issued |
| cancelled | 2 | Redemption cancelled, points returned |

### Rule Status

| Status | Value | Description |
|--------|-------|-------------|
| draft | 0 | Rule created but not active |
| active | 1 | Rule is active and processing |
| inactive | 2 | Rule deactivated |

### Tiered Configuration JSON Schema

```json
{
  "tiers": [
    {
      "threshold": 10000,    // Upper threshold in cents (e.g., $100.00)
      "ratio": "1.0"         // Points per currency unit
    },
    {
      "threshold": 50000,    // $100.01 to $500.00
      "ratio": "1.5"
    },
    {
      "threshold": null,     // Above $500.00 (last tier)
      "ratio": "2.0"
    }
  ]
}
```

### Calculation Examples

#### Tiered Calculation Example

Rule: 1 point/$1 up to $100, 1.5 points/$1 from $100-$500, 2 points/$1 above $500

| Order Amount | Points Calculation | Total Points |
|--------------|-------------------|--------------|
| $50 | 50 * 1.0 | 50 |
| $100 | 100 * 1.0 | 100 |
| $200 | 100 * 1.0 + 100 * 1.5 | 250 |
| $600 | 100 * 1.0 + 400 * 1.5 + 100 * 2.0 | 900 |

#### Points Expiration Example

- Points earned on 2026-01-15 with 12-month expiration
- Expiration date: 2027-01-15
- System batch job checks daily for expired points
- On 2027-01-15, points are marked as expired and deducted from balance