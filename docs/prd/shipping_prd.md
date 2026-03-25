# Shipping Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Shipping Module PRD |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-25 |
| Author | Product Team |
| Module Location | admin/internal/domain/shipping/ |

---

## Product Overview

### Summary

The Shipping module is a core component of the ShopJoy e-commerce SaaS platform that handles shipping fee rule configuration and calculation. This module enables merchants to create flexible shipping templates with region-based pricing, supporting multiple fee calculation methods (fixed, per-item, per-weight, and free shipping). The shipping fee rules configured here power the checkout experience, cart preview, and order creation flows.

The MVP release focuses on three admin-facing pages: Shipping Template List, Template Detail Page (with zones and associations), and Shipping Fee Calculator (testing tool). These pages allow merchants to configure shipping rules that integrate seamlessly with the order checkout process.

### Problem Statement

Merchants on the ShopJoy platform currently face several challenges in shipping fee management:

- **No centralized shipping rule configuration**: Merchants cannot define different shipping fees for different regions
- **Limited fee calculation options**: No support for per-item or per-weight pricing models
- **No free shipping threshold configuration**: Cannot set minimum order amounts for free shipping
- **Lack of template association**: Cannot assign different shipping rules to specific products or categories
- **No testing tool for fee calculation**: Merchants cannot verify shipping rules before going live
- **Complex multi-region pricing**: Difficult to manage shipping costs across different delivery zones

### Solution Overview

Implement a shipping configuration module that provides:

1. Shipping template management with CRUD operations
2. Shipping zone configuration with region-based pricing
3. Multiple fee types: fixed, by_count, by_weight, and free
4. Free shipping threshold settings (by amount or quantity)
5. Template associations with products and categories
6. Shipping fee calculator tool for testing rules
7. Integration with checkout and cart preview flows

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | Reduce shipping configuration time | Average time to create template | <10 minutes |
| BG-002 | Increase shipping rule adoption | Percentage of merchants with configured templates | >80% within 1 month |
| BG-003 | Reduce shipping-related support inquiries | Support ticket reduction | -40% within 3 months |
| BG-004 | Improve checkout conversion | Cart abandonment rate due to shipping | Reduce by 15% |
| BG-005 | Enable flexible pricing strategies | Adoption of multiple fee types | >50% using non-fixed fees |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | Create and manage shipping templates easily | Merchant Admin |
| UG-002 | Configure region-based shipping fees | Merchant Admin |
| UG-003 | Set free shipping thresholds | Merchant Admin |
| UG-004 | Associate templates with products/categories | Merchant Admin |
| UG-005 | Test shipping fee calculation before going live | Merchant Admin |
| UG-006 | View correct shipping fees at checkout | Buyer |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | Carrier API integration | Phase 2 feature (requires carrier partnerships) |
| NG-002 | Real-time shipping rate calculation from carriers | Phase 2 feature |
| NG-003 | Multi-warehouse shipping logic | Phase 2 feature |
| NG-004 | International shipping with customs | Phase 2 feature |
| NG-005 | Shipping label generation | Part of Fulfillment module |
| NG-006 | Time-based shipping rules (peak season pricing) | Phase 2 feature |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | Store owner or operations manager | Create templates, configure zones, set free shipping |
| Operations Staff | Warehouse or fulfillment staff | View shipping rules, test fee calculation |
| Product Manager | Category or product managers | Associate products with shipping templates |
| Buyer | End customer placing orders | View shipping fees at checkout |

### Persona Details

#### Merchant Admin (Primary)

- **Demographics**: Small to medium business owners, e-commerce operations managers
- **Technical Proficiency**: Moderate
- **Goals**: Accurate shipping pricing, reduce shipping losses, competitive pricing
- **Pain Points**: Complex regional pricing, unclear fee calculations, no testing capability
- **Frequency**: Initial setup, occasional updates for promotions

#### Operations Staff (Secondary)

- **Demographics**: Warehouse workers, fulfillment coordinators
- **Technical Proficiency**: Basic to moderate
- **Goals**: Understand shipping rules, verify fees during fulfillment
- **Pain Points**: Unclear which template applies to which product
- **Frequency**: Daily reference during fulfillment

#### Product Manager (Secondary)

- **Demographics**: Category managers, merchandisers
- **Technical Proficiency**: Moderate
- **Goals**: Assign appropriate shipping templates to products
- **Pain Points**: Manual association process, bulk updates
- **Frequency**: When adding new products, during promotions

#### Buyer (Tertiary)

- **Demographics**: Online shoppers
- **Technical Proficiency**: Varies
- **Goals**: Transparent shipping costs, fair pricing
- **Pain Points**: Unexpected shipping fees, unclear free shipping eligibility
- **Frequency**: During checkout

### Role-Based Access

| Role | Permissions |
|------|-------------|
| Tenant Admin | Full access to shipping templates, zones, mappings |
| Tenant Operations Manager | View templates, use calculator, edit zones |
| Tenant Product Manager | View templates, manage product associations |
| Buyer | View shipping fees at checkout (no admin access) |

---

## Functional Requirements

### Priority 1: Shipping Template Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | Create shipping template | Create new template with name and default status | P0 |
| FR-002 | Edit shipping template | Modify template name and settings | P0 |
| FR-003 | Delete shipping template | Delete template (with validation) | P0 |
| FR-004 | Set default template | Designate one template as default for the shop | P0 |
| FR-005 | List shipping templates | View all templates with key statistics | P0 |
| FR-006 | Template activation toggle | Enable/disable template without deleting | P1 |
| FR-007 | Prevent deletion of default template | Must set another as default first | P0 |

### Priority 2: Shipping Zone Configuration

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-008 | Create shipping zone | Add new zone to template with regions | P0 |
| FR-009 | Edit shipping zone | Modify zone name, regions, fee settings | P0 |
| FR-010 | Delete shipping zone | Remove zone from template | P0 |
| FR-011 | Reorder shipping zones | Change zone matching priority | P1 |
| FR-012 | City-level region selection | Select regions by city administrative codes | P0 |
| FR-013 | Province quick-select | Select all cities in a province at once | P1 |
| FR-014 | Duplicate region validation | Warn if regions overlap with other zones | P1 |

### Priority 3: Fee Configuration

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-015 | Fixed fee type | Set flat shipping fee regardless of quantity/weight | P0 |
| FR-016 | Per-item fee type | Calculate fee based on item quantity | P0 |
| FR-017 | Per-weight fee type | Calculate fee based on total weight (grams) | P0 |
| FR-018 | Free shipping type | No shipping fee for this zone | P0 |
| FR-019 | First unit/fee configuration | Set first unit quantity and corresponding fee | P0 |
| FR-020 | Additional unit/fee configuration | Set additional unit increment and fee | P0 |
| FR-021 | Free shipping threshold by amount | Free shipping when order reaches amount | P0 |
| FR-022 | Free shipping threshold by count | Free shipping when item count reaches threshold | P1 |
| FR-023 | Combined threshold support | Support both amount and count thresholds | P1 |

### Priority 4: Template Associations

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-024 | Associate products with template | Assign specific template to products | P0 |
| FR-025 | Associate categories with template | Assign specific template to categories | P0 |
| FR-026 | View product associations | List all products using a template | P0 |
| FR-027 | View category associations | List all categories using a template | P0 |
| FR-028 | Remove associations | Unlink products/categories from template | P0 |
| FR-029 | Product search for association | Search products by name/SKU when associating | P1 |
| FR-030 | Bulk product association | Associate multiple products at once | P2 |

### Priority 5: Shipping Fee Calculator

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-031 | Address selector for calculation | Select province/city/district for test | P0 |
| FR-032 | Add test items | Add products with quantity for test | P0 |
| FR-033 | Calculate shipping fee | Execute fee calculation and display result | P0 |
| FR-034 | Display matched template/zone | Show which template and zone matched | P0 |
| FR-035 | Display fee breakdown | Show calculation details (units, fees) | P1 |
| FR-036 | Quick product selector | Search and select products for test | P1 |
| FR-037 | Manual item entry | Enter weight/price manually for test | P1 |

---

## User Experience

### Entry Points

| Entry Point | Description | User Flow |
|-------------|-------------|-----------|
| Settings Navigation | Shipping under Settings menu | Settings > Shipping > Templates |
| Product Edit Page | Link to associate shipping template | Product > Edit > Shipping Template |
| Category Management | Link to assign category template | Categories > Edit > Shipping Template |
| Order Page | View applied shipping template | Orders > Detail > Shipping Info |

### Page Structure

```
/settings/shipping
├── /templates           # Template list page
├── /templates/:id       # Template detail page
│   ├── Zones Tab        # Shipping zone configuration
│   └── Associations Tab # Product/category associations
└── /calculator          # Shipping fee calculator tool
```

### Page 1: Shipping Template List

**Purpose**: Display all shipping templates with key information and quick actions.

**Layout**:
- Header: Page title, description, "Create Template" button
- Stats Bar: Total templates, Active templates, Default template indicator
- Template Cards: Grid layout with card for each template

**Template Card Content**:
- Template name
- Default badge (if applicable)
- Zone count
- Associated products count
- Associated categories count
- Status toggle
- Actions: Edit, Delete, Set as Default

**Interactions**:
- Click card to navigate to detail page
- Toggle status inline
- Set as default with confirmation
- Delete with validation (cannot delete default)

**Empty State**:
- Illustration with message "No shipping templates yet"
- "Create your first template" CTA button

### Page 2: Template Detail Page

**Purpose**: Configure template zones and manage associations.

**Layout**:
- Header: Template name, Default toggle, Back button
- Tabs: Zones | Associations
- Footer: Save/Cancel actions

**Zones Tab**:
- Zone list with drag-and-drop reorder
- Each zone shows: name, regions preview, fee type, quick actions
- "Add Zone" button to create new zone
- Zone edit dialog with full configuration

**Zone Edit Dialog Fields**:
1. Zone name (required)
2. Region selector (required)
   - Left panel: Province tree with expand/collapse
   - Right panel: Selected regions list
   - "Select all province" checkbox per province
3. Fee type selector (required)
   - Fixed / Per-item / Per-weight / Free
4. Fee configuration (dynamic based on fee type)
   - Fixed: Single fee amount
   - Per-item/Per-weight: First unit, First fee, Additional unit, Additional fee
5. Free shipping thresholds (optional)
   - By amount: Minimum order amount
   - By count: Minimum item quantity

**Associations Tab**:
- Two sections: Products | Categories
- Product section: Search and multi-select products
- Category section: Tree multi-select categories
- Current associations list with remove option

### Page 3: Shipping Fee Calculator

**Purpose**: Test shipping fee calculation before going live.

**Layout**:
- Left panel: Test input form
- Right panel: Calculation result display

**Test Input Form**:
1. Delivery Address
   - Province selector (cascader)
   - City selector (cascader)
   - District selector (cascader)
2. Test Items
   - Product selector (searchable)
   - Quantity input
   - Or manual entry: Weight (grams), Price
   - Add multiple items
   - Remove items

**Calculation Result Display**:
- Total shipping fee (prominent)
- Matched template name
- Matched zone name
- Fee breakdown:
  - Fee type
  - First unit/fee
  - Additional units calculated
  - Additional fee
  - Free shipping eligibility (if applicable)

**Interactions**:
- Real-time calculation on input change
- Clear form button
- Test history (optional)

### User Flows

#### Create Shipping Template

1. Navigate to Settings > Shipping
2. Click "Create Template" button
3. Enter template name
4. System creates template and navigates to detail page
5. Add zones with regions and fee rules
6. Save template
7. Optionally set as default

#### Configure Shipping Zone

1. Navigate to template detail page
2. Click "Add Zone" on Zones tab
3. Enter zone name (e.g., "East China")
4. Select regions using region selector
5. Select fee type
6. Configure fee parameters
7. Optionally set free shipping thresholds
8. Save zone

#### Test Shipping Fee

1. Navigate to Shipping Fee Calculator
2. Select delivery address (province/city/district)
3. Add test items (select products or enter manually)
4. Click "Calculate"
5. View shipping fee result
6. Review matched template and zone
7. Adjust inputs and re-test as needed

---

## Narrative

### Scenario 1: Setting Up National Shipping

As a Merchant Admin for an online clothing store, I need to set up shipping rules for nationwide delivery. I navigate to Settings > Shipping and click "Create Template". I name it "National Standard Shipping" and click Create.

The system takes me to the template detail page. I click "Add Zone" to configure my first shipping zone. I name it "East China" and use the region selector to select Shanghai, Jiangsu, Zhejiang, and Anhui. I choose "Per-weight" as my fee type, set 1000g as the first unit with 12 CNY fee, and 500g as additional unit with 3 CNY additional fee. I also set 99 CNY as the free shipping threshold.

I repeat this process for other regions: "South China", "North China", "Remote Areas", each with different fee structures based on distance from my warehouse. I also create a "Free Shipping Zone" for my local city with free shipping enabled.

After configuring all zones, I click "Set as Default" to make this my default shipping template. Now when customers check out, the system will automatically match their address to the appropriate zone and calculate the shipping fee.

### Scenario 2: Testing Shipping Rules

Before going live, I want to verify my shipping rules work correctly. I navigate to the Shipping Fee Calculator page. I select "Hangzhou, Zhejiang" as the delivery address, then add a product that weighs 1.5kg with a price of 150 CNY.

I click "Calculate" and see that the shipping fee is 15 CNY, matched to my "East China" zone. The breakdown shows: first 1000g costs 12 CNY, and the additional 500g costs 3 CNY. I also notice that the free shipping threshold of 99 CNY is met (150 CNY > 99 CNY), so the final shipping fee is 0 CNY - free shipping!

This confirms my rules are working as expected. I make a note to adjust my free shipping threshold to 199 CNY to avoid giving free shipping on smaller orders.

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Template Creation Time | Average time to create a complete template | <10 minutes |
| Zone Configuration Success | Percentage of zones configured correctly on first try | >90% |
| Calculator Usage | Percentage of merchants using calculator before going live | >60% |
| Template Adoption Rate | Percentage of merchants with at least one active template | >80% within 1 month |

### Business Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Shipping Configuration Accuracy | Reduction in shipping fee disputes | -50% |
| Checkout Completion Rate | Improvement in checkout conversion | +10% |
| Support Ticket Reduction | Decrease in shipping-related inquiries | -40% in 3 months |
| Template Usage Distribution | Balanced usage of fee types | All fee types used |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API Response Time | Template CRUD operations | <200ms p95 |
| Fee Calculation Time | Shipping fee calculation endpoint | <100ms p95 |
| System Uptime | Shipping service availability | 99.9% |
| Error Rate | Failed template/zone operations | <0.1% |

---

## Technical Considerations

### Integration Points

| Integration | Description | Method |
|-------------|-------------|--------|
| Checkout Service | Calculate shipping fee during checkout | Internal API |
| Cart Service | Preview shipping fee in cart | Internal API |
| Product Service | Fetch product weight/price for calculation | Internal API |
| Order Service | Apply shipping template to order | Domain Event |
| Notification Service | Notify on template changes | Message Queue |

### Data Storage and Privacy

| Aspect | Consideration |
|--------|---------------|
| Monetary Values | Use BIGINT (cents) for all amounts, no floating-point |
| Time Storage | All timestamps in UTC, RFC3339 format |
| Multi-Tenancy | All queries filtered by tenant_id |
| Soft Delete | Support deleted_at for data recovery |
| Region Codes | Use Chinese administrative division 6-digit codes |
| JSONB Storage | Region arrays stored as JSONB for flexible queries |

### Scalability and Performance

| Consideration | Strategy |
|---------------|----------|
| Template Caching | Cache default template per tenant in Redis |
| Zone Caching | Cache zone configurations per template |
| Region Indexing | GIN index on JSONB regions column |
| Fee Calculation | Pre-calculate common scenarios, cache results |
| Concurrent Updates | Optimistic locking on template updates |

### Potential Challenges

| Challenge | Mitigation |
|-----------|------------|
| Region Overlap | Validation warnings, clear zone priority |
| Complex Fee Calculation | Comprehensive test coverage, calculator tool |
| Template Association Sync | Event-driven updates, cache invalidation |
| Large Region Lists | Pagination, search functionality |
| Currency Precision | Decimal type in Go, BIGINT in database |

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Team Size | Description |
|-------|----------|-----------|-------------|
| Phase 1: Backend Core | 1 week | 2 backend | Database schema, entities, repositories, services |
| Phase 2: Admin API | 1 week | 2 backend | CRUD endpoints, fee calculation logic |
| Phase 3: Frontend List Page | 0.5 week | 1 frontend | Template list with card layout |
| Phase 4: Frontend Detail Page | 1 week | 1 frontend | Zone config, region selector, associations |
| Phase 5: Calculator Page | 0.5 week | 1 frontend | Fee calculation testing tool |
| Phase 6: Integration & Testing | 0.5 week | 2 backend, 1 frontend | Checkout integration, E2E tests |

**Total Estimate: 4.5 weeks**

### Suggested Phases

#### Phase 1: Backend Core (Week 1)

- Implement shipping domain entities
- Create repository interfaces and implementations
- Implement fee calculation service
- Create database migrations
- Define error codes and constants

#### Phase 2: Admin API (Week 2)

- Build template CRUD APIs
- Build zone CRUD APIs
- Build mapping CRUD APIs
- Implement fee calculation API
- Add validation and business rules

#### Phase 3: Frontend List Page (Week 2.5)

- Create template list page with card layout
- Implement create/edit/delete template
- Add default template handling
- Style and polish

#### Phase 4: Frontend Detail Page (Week 3-3.5)

- Build template detail page with tabs
- Implement zone configuration with region selector
- Build fee type configuration forms
- Implement associations management
- Add drag-and-drop zone reordering

#### Phase 5: Calculator Page (Week 4)

- Build shipping fee calculator page
- Implement address selector (cascader)
- Add product/item input forms
- Display calculation results with breakdown
- Add test history (optional)

#### Phase 6: Integration & Testing (Week 4.5)

- Integrate with checkout service
- Integrate with cart preview
- Write E2E tests
- Performance testing
- Bug fixes and polish

---

## Database Schema

### shipping_templates

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| name | VARCHAR(100) | NO | - | Template name |
| is_default | BOOLEAN | NO | FALSE | Is default template |
| is_active | BOOLEAN | NO | TRUE | Is enabled |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update timestamp |
| deleted_at | TIMESTAMP | YES | NULL | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_tenant_default` (`tenant_id`, `is_default`)

### shipping_zones

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| template_id | BIGINT | NO | - | Parent template ID |
| name | VARCHAR(100) | NO | - | Zone name |
| regions | JSONB | NO | '[]' | City code array |
| fee_type | VARCHAR(20) | NO | 'fixed' | Fee type |
| first_unit | INT | NO | 1 | First unit (count or grams) |
| first_fee | BIGINT | NO | 0 | First fee (cents) |
| additional_unit | INT | NO | 1 | Additional unit |
| additional_fee | BIGINT | NO | 0 | Additional fee (cents) |
| free_threshold_amount | BIGINT | NO | 0 | Free shipping threshold amount |
| free_threshold_count | INT | NO | 0 | Free shipping threshold count |
| sort | INT | NO | 0 | Sort order |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Update timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_template_id` (`template_id`)
- KEY `idx_regions` (`regions`) USING GIN

### shipping_template_mappings

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| template_id | BIGINT | NO | - | Shipping template ID |
| target_type | VARCHAR(20) | NO | - | product / category |
| target_id | BIGINT | NO | - | Product ID or Category ID |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- KEY `idx_tenant_id` (`tenant_id`)
- KEY `idx_template_id` (`template_id`)
- KEY `idx_target` (`target_type`, `target_id`)

---

## API Endpoints

### Shipping Template Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/shipping-templates | List templates |
| POST | /api/v1/shipping-templates | Create template |
| GET | /api/v1/shipping-templates/:id | Get template detail |
| PUT | /api/v1/shipping-templates/:id | Update template |
| DELETE | /api/v1/shipping-templates/:id | Delete template |
| PUT | /api/v1/shipping-templates/:id/set-default | Set as default |

### Shipping Zone Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/shipping-templates/:id/zones | Add zone |
| PUT | /api/v1/shipping-zones/:id | Update zone |
| DELETE | /api/v1/shipping-zones/:id | Delete zone |
| PUT | /api/v1/shipping-templates/:id/zones/reorder | Reorder zones |

### Template Mapping Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/shipping-templates/:id/mappings | Get mappings |
| POST | /api/v1/shipping-template-mappings | Create mapping |
| DELETE | /api/v1/shipping-template-mappings/:id | Delete mapping |

### Shipping Fee Calculator (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/shipping/calculate | Calculate shipping fee |

---

## Business Rules

### Template Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | One default per tenant | Each tenant has exactly one default template |
| BR-002 | Cannot delete default | Must set another as default first |
| BR-003 | Cascade delete | Deleting template removes zones and mappings |
| BR-004 | Template name uniqueness | Name must be unique within tenant |
| BR-005 | First template is default | If no default exists, first created is default |

### Zone Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-006 | Zone matching priority | Lower sort = higher priority |
| BR-007 | First match wins | First zone matching city code is used |
| BR-008 | Regions required | Zone must have at least one region |
| BR-009 | Fee type validation | Fee parameters must match fee type |
| BR-010 | Positive values only | All fee values must be >= 0 |

### Fee Calculation Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-011 | Fixed fee | shipping_fee = first_fee |
| BR-012 | Per-item fee | First unit covers initial quantity, additional for remainder |
| BR-013 | Per-weight fee | First unit covers initial weight, additional for remainder |
| BR-014 | Free shipping | shipping_fee = 0 |
| BR-015 | Free threshold check | Free shipping if order_amount >= threshold OR count >= threshold |
| BR-016 | Unit ceiling | Additional units calculated with ceiling division |

### Template Association Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-017 | Priority order | Product > Category > Default |
| BR-018 | One template per product | Product can only have one template |
| BR-019 | One template per category | Category can only have one template |
| BR-020 | Delete cascade | Removing template removes mappings |

---

## User Stories

### US-001: Create Shipping Template

**Description**: As a Merchant Admin, I want to create a new shipping template so that I can configure shipping rules.

**Acceptance Criteria**:
- Given I am on the Shipping Templates page
- When I click "Create Template"
- And I enter a template name "Standard Shipping"
- And I click Create
- Then a new template is created
- And I am navigated to the template detail page
- And the template is active by default

---

### US-002: Set Default Template

**Description**: As a Merchant Admin, I want to set a template as default so that it applies when no specific template is assigned.

**Acceptance Criteria**:
- Given I have multiple templates
- And one is already the default
- When I click "Set as Default" on another template
- And I confirm the action
- Then the new template becomes the default
- And the previous default is no longer marked as default
- And only one template is marked as default at a time

---

### US-003: Delete Non-Default Template

**Description**: As a Merchant Admin, I want to delete a template so that I can remove unused configurations.

**Acceptance Criteria**:
- Given a template that is not the default
- When I click "Delete" on the template
- And I confirm the deletion
- Then the template is soft-deleted
- And all associated zones are deleted
- And all mappings are deleted
- And the template no longer appears in the list

---

### US-004: Prevent Deletion of Default Template

**Description**: As the system, I want to prevent deletion of the default template so that there is always a fallback shipping rule.

**Acceptance Criteria**:
- Given a template that is marked as default
- When I try to delete it
- Then I see an error message "Cannot delete default template. Please set another template as default first."
- And the delete action is blocked
- And the template is not deleted

---

### US-005: Create Shipping Zone

**Description**: As a Merchant Admin, I want to create a shipping zone so that I can define region-specific shipping rules.

**Acceptance Criteria**:
- Given I am on a template detail page
- When I click "Add Zone"
- And I enter zone name "East China"
- And I select regions: Shanghai, Hangzhou, Nanjing
- And I select fee type "Per-weight"
- And I set first unit 1000g with fee 1200 cents
- And I set additional unit 500g with fee 300 cents
- And I save the zone
- Then the zone is created
- And it appears in the zone list
- And the zone count on the template card updates

---

### US-006: Configure Fixed Fee

**Description**: As a Merchant Admin, I want to set a fixed shipping fee so that customers pay the same rate regardless of order size.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I select fee type "Fixed"
- And I enter first fee 1500 (cents)
- And I save
- Then the zone uses fixed fee calculation
- And any order to this zone will have 15.00 CNY shipping

---

### US-007: Configure Per-Item Fee

**Description**: As a Merchant Admin, I want to charge shipping based on item quantity so that larger orders have higher shipping.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I select fee type "Per-item"
- And I set first unit 2 with fee 1000
- And I set additional unit 1 with fee 300
- And I save
- Then the zone uses per-item calculation
- And an order with 3 items will have: 10.00 + 3.00 = 13.00 CNY shipping

---

### US-008: Configure Per-Weight Fee

**Description**: As a Merchant Admin, I want to charge shipping based on weight so that heavier packages cost more.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I select fee type "Per-weight"
- And I set first unit 1000 (grams) with fee 1200
- And I set additional unit 500 (grams) with fee 300
- And I save
- Then the zone uses per-weight calculation
- And an order with 1800g total weight will have: 12.00 + 0.30 = 12.30 CNY shipping

---

### US-009: Configure Free Shipping

**Description**: As a Merchant Admin, I want to offer free shipping for a zone so that customers in that region pay no shipping.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I select fee type "Free"
- And I save
- Then the zone has 0 shipping fee
- And any order to this zone will be free shipping

---

### US-010: Set Free Shipping Threshold

**Description**: As a Merchant Admin, I want to offer free shipping when order reaches a threshold so that I can encourage larger orders.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I set free threshold amount to 9900 (cents)
- And I set free threshold count to 3
- And I save
- Then orders >= 99 CNY OR with >= 3 items will be free shipping
- And orders below both thresholds use normal fee calculation

---

### US-011: Select Regions for Zone

**Description**: As a Merchant Admin, I want to select cities for a zone so that the zone applies to specific regions.

**Acceptance Criteria**:
- Given I am creating/editing a zone
- When I open the region selector
- Then I see a tree of provinces and cities
- And I can expand a province to see its cities
- And I can check a city to add it to the zone
- And I can check a province to add all its cities
- And I see a list of selected regions on the right
- And I can remove regions from the list

---

### US-012: Reorder Shipping Zones

**Description**: As a Merchant Admin, I want to reorder zones so that I can control matching priority.

**Acceptance Criteria**:
- Given I have multiple zones in a template
- When I drag a zone to a new position
- Then the zone order changes
- And zones are matched in the new order (lower position = higher priority)
- And the sort values are updated

---

### US-013: Associate Products with Template

**Description**: As a Merchant Admin, I want to assign a shipping template to specific products so that they use different shipping rules.

**Acceptance Criteria**:
- Given I am on a template's Associations tab
- When I click "Add Products"
- And I search for products by name or SKU
- And I select products "T-Shirt Blue" and "T-Shirt Red"
- And I confirm
- Then the products are associated with this template
- And they appear in the products list
- And the product count updates

---

### US-014: Associate Categories with Template

**Description**: As a Merchant Admin, I want to assign a shipping template to a category so that all products in it use the same rules.

**Acceptance Criteria**:
- Given I am on a template's Associations tab
- When I click "Add Categories"
- And I see a category tree
- And I select categories "Electronics" and "Home Appliances"
- And I confirm
- Then the categories are associated with this template
- And they appear in the categories list
- And the category count updates

---

### US-015: Remove Template Association

**Description**: As a Merchant Admin, I want to remove a product or category association so that it uses the default template instead.

**Acceptance Criteria**:
- Given a product is associated with a template
- When I click "Remove" on the product
- And I confirm
- Then the association is deleted
- And the product will use its category template or default template
- And the product count updates

---

### US-016: Calculate Shipping Fee (Calculator)

**Description**: As a Merchant Admin, I want to test shipping fee calculation so that I can verify my rules work correctly.

**Acceptance Criteria**:
- Given I am on the Shipping Fee Calculator page
- When I select delivery address "Hangzhou, Zhejiang"
- And I add a test item: Product A, quantity 2, weight 800g, price 15000
- And I click "Calculate"
- Then I see the calculated shipping fee
- And I see the matched template name
- And I see the matched zone name
- And I see the fee breakdown

---

### US-017: View Fee Breakdown

**Description**: As a Merchant Admin, I want to see how the shipping fee was calculated so that I understand the formula.

**Acceptance Criteria**:
- Given I have calculated a shipping fee
- Then I see a breakdown section with:
  - Fee type: "Per-weight"
  - First unit: 1000g
  - First fee: 12.00 CNY
  - Calculated weight: 1600g
  - Additional units: 1
  - Additional fee: 3.00 CNY
  - Subtotal: 15.00 CNY
  - Free threshold check: Not met (150 CNY < 199 CNY)
  - Final fee: 15.00 CNY

---

### US-018: Test Free Shipping Eligibility

**Description**: As a Merchant Admin, I want to verify free shipping thresholds work correctly.

**Acceptance Criteria**:
- Given a zone with free threshold 9900 cents
- When I test with order amount 10000 cents
- Then the calculated fee is 0 (free shipping)
- And the breakdown shows "Free shipping threshold met (100.00 CNY >= 99.00 CNY)"

---

### US-019: Handle Unmatched Address

**Description**: As a Merchant Admin, I want to see what happens when an address doesn't match any zone.

**Acceptance Criteria**:
- Given a template with zones for "East China" only
- When I test with address "Urumqi, Xinjiang"
- Then I see "No matching zone found" message
- And the fee uses the default template's default zone

---

### US-020: Template List with Statistics

**Description**: As a Merchant Admin, I want to see key statistics for each template on the list page.

**Acceptance Criteria**:
- Given I am on the Shipping Templates page
- Then each template card shows:
  - Template name
  - Default badge (if applicable)
  - Zone count: 5
  - Associated products: 12
  - Associated categories: 3
  - Status toggle

---

### US-021: Toggle Template Status

**Description**: As a Merchant Admin, I want to quickly enable/disable a template without deleting it.

**Acceptance Criteria**:
- Given an active template
- When I toggle the status switch to off
- Then the template becomes inactive
- And products using it will fall back to category/default template
- And I can toggle it back to active anytime

---

### US-022: Validate Zone Regions

**Description**: As the system, I want to validate that zone regions are valid city codes.

**Acceptance Criteria**:
- Given I am saving a zone
- When the regions contain an invalid city code
- Then I see an error "Invalid region codes detected"
- And the zone is not saved
- When the regions array is empty
- Then I see an error "At least one region is required"

---

### US-023: Warn on Overlapping Regions

**Description**: As the system, I want to warn when regions overlap between zones in the same template.

**Acceptance Criteria**:
- Given Zone A includes Shanghai
- And I am creating Zone B
- When I add Shanghai to Zone B
- Then I see a warning "Shanghai is already in Zone A. This may cause unexpected matching behavior."
- And I can still save Zone B
- And the zone with lower sort order will match first

---

### US-024: Search Products for Association

**Description**: As a Merchant Admin, I want to search products when associating them with a template.

**Acceptance Criteria**:
- Given I am adding products to a template
- When I type "shirt" in the search box
- Then I see products with "shirt" in name or SKU
- And I can select multiple products
- And products already associated are marked

---

### US-025: Bulk Remove Associations

**Description**: As a Merchant Admin, I want to remove multiple associations at once.

**Acceptance Criteria**:
- Given I have multiple products associated with a template
- When I select products using checkboxes
- And I click "Remove Selected"
- And I confirm
- Then all selected associations are deleted
- And the product count updates

---

### US-026: Copy Template

**Description**: As a Merchant Admin, I want to duplicate a template so that I can create variations easily.

**Acceptance Criteria**:
- Given I have a configured template
- When I click "Copy" on the template
- Then a new template is created with name "Copy of [original name]"
- And all zones are copied
- And associations are NOT copied
- And the new template is not default

---

### US-027: Export Template Configuration

**Description**: As a Merchant Admin, I want to export my template configuration for backup.

**Acceptance Criteria**:
- Given I am on a template detail page
- When I click "Export"
- Then a JSON file is downloaded
- And the file contains all zones and their configurations
- And the file can be imported to restore or share

---

### US-028: Import Template Configuration

**Description**: As a Merchant Admin, I want to import a template configuration so that I can quickly set up new templates.

**Acceptance Criteria**:
- Given I have a valid template JSON file
- When I click "Import" and select the file
- Then a new template is created from the import
- And all zones are created
- And I can edit the imported template

---

### US-029: View Template Audit Log

**Description**: As a Merchant Admin, I want to see when templates were created and modified.

**Acceptance Criteria**:
- Given I am viewing a template
- Then I see created_at timestamp
- And I see updated_at timestamp
- And I can see a history of changes (future enhancement)

---

### US-030: Multi-Language Template Names

**Description**: As the system, I want to support template names in multiple languages.

**Acceptance Criteria**:
- Given the system supports multiple locales
- When I create a template
- Then I can enter names in different languages
- And the appropriate name is displayed based on user locale
- Note: MVP supports single language only, this is Phase 2

---

## UI/UX Specifications

### Component Library Reference

Based on existing patterns in the shop-admin project:

| Component | Element Plus | Usage |
|-----------|--------------|-------|
| Cards | el-card | Template cards, stats cards |
| Forms | el-form, el-form-item | Zone configuration, template creation |
| Tables | el-table | Associations list |
| Dialogs | el-dialog | Zone edit, confirmations |
| Selects | el-select, el-cascader | Fee type, region selector |
| Switches | el-switch | Status toggle, default toggle |
| Buttons | el-button | Actions throughout |
| Tags | el-tag | Zone fee type display |
| Icons | @element-plus/icons-vue | UI icons |
| Pagination | el-pagination | Association lists |

### Color Scheme

Following existing design patterns:
- Primary: Indigo (#6366F1)
- Success: Green (#10B981)
- Warning: Amber (#F59E0B)
- Danger: Red (#EF4444)
- Info: Gray (#6B7280)

### Responsive Design

- Desktop: Full card grid (4 columns)
- Tablet: Card grid (2 columns)
- Mobile: Single column layout
- Dialog forms scrollable on mobile

---

## Error Handling

### User-Facing Errors

| Error Code | Message | Resolution |
|------------|---------|------------|
| 180001 | Template name is required | Enter a name |
| 180002 | Template not found | Check template ID |
| 180003 | Cannot delete default template | Set another as default first |
| 180004 | Zone not found | Check zone ID |
| 180005 | Invalid region codes | Select valid city codes |
| 180006 | Fee configuration invalid | Check fee type and values |
| 180007 | No matching zone for address | Add zone for this region |
| 180008 | Product already associated | Remove existing association first |

### System Errors

- Database errors: Log and show generic error message
- Validation errors: Show specific field errors
- Network errors: Retry with exponential backoff
- Concurrent modification: Show conflict resolution dialog

---

## Glossary

| Term | Definition |
|------|------------|
| Shipping Template | A collection of shipping fee rules for a shop |
| Shipping Zone | A delivery area with specific fee rules |
| Fee Type | Method for calculating shipping: fixed, by_count, by_weight, free |
| Free Shipping Threshold | Order amount or quantity that triggers free shipping |
| Template Mapping | Association between template and product/category |
| Region Code | 6-digit Chinese administrative division code (city level) |
| First Unit | Initial quantity or weight for base fee |
| Additional Unit | Increment for calculating additional fees |

---

## Appendix

### Fee Type Examples

**Fixed Fee Example:**
- Zone: Local City
- Fee type: Fixed
- First fee: 500 (5.00 CNY)
- Result: Always 5.00 CNY shipping

**Per-Item Example:**
- Zone: National
- Fee type: by_count
- First unit: 2 items, First fee: 1000 (10.00 CNY)
- Additional unit: 1 item, Additional fee: 300 (3.00 CNY)
- Order: 5 items
- Calculation: 10.00 + (5-2) / 1 * 3.00 = 19.00 CNY

**Per-Weight Example:**
- Zone: Remote Areas
- Fee type: by_weight
- First unit: 1000g, First fee: 1500 (15.00 CNY)
- Additional unit: 500g, Additional fee: 500 (5.00 CNY)
- Order: 2800g
- Calculation: 15.00 + ceil((2800-1000)/500) * 5.00 = 15.00 + 4 * 5.00 = 35.00 CNY

**Free Shipping with Threshold Example:**
- Zone: VIP Region
- Fee type: by_weight (as backup)
- First unit: 1000g, First fee: 1000
- Free threshold amount: 19900 (199.00 CNY)
- Order: 2500g, 250.00 CNY
- Result: 250.00 >= 199.00, Free shipping!

### Region Code Reference

City codes follow Chinese National Bureau of Statistics standards:

| Province | Code Prefix | Example Cities |
|----------|-------------|----------------|
| Beijing | 11 | 110100 (Beijing City) |
| Shanghai | 31 | 310100 (Shanghai City) |
| Zhejiang | 33 | 330100 (Hangzhou), 330200 (Ningbo) |
| Jiangsu | 32 | 320100 (Nanjing), 320500 (Suzhou) |
| Guangdong | 44 | 440100 (Guangzhou), 440300 (Shenzhen) |