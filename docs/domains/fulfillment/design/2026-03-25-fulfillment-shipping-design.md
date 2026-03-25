# Shipping Module Design Document

**Version:** 1.0
**Date:** 2026-03-25
**Author:** Claude Code
**Status:** Draft

---

## 1. Overview

### 1.1 Module Positioning

The Shipping module is responsible for **shipping fee rule configuration and calculation**, providing shipping fee amounts for order creation and checkout pages. This is a critical component after product pricing.

**Key Distinction from Fulfillment:**

| Module | Core Responsibility | Key Entities |
|--------|---------------------|--------------|
| **Shipping** | Shipping fee rule configuration | Shipping Template, Shipping Zone, Fee Calculation Rules |
| **Fulfillment** | Shipment execution | Shipment, Tracking Number, Delivery Status, Refund |

Shipping = Pre-order rules (how much to charge for delivery)
Fulfillment = Post-order execution (delivering the package)

### 1.2 Business Goals

Enable merchants to:
- Set different shipping fee rules based on delivery regions
- Support common shipping templates (fixed fee, free shipping threshold, per-item/weight pricing)
- Configure differentiated shipping fees for different products and regions
- Easily manage and test shipping fee rules in the admin panel
- Platform admins can manage global shipping templates

---

## 2. Core Concepts

### 2.1 Shipping Template

A collection of shipping fee rules. Each tenant can have multiple templates, with one designated as the default.

### 2.2 Shipping Zone

A delivery area defined by city-level administrative divisions, with specific fee calculation rules.

### 2.3 Fee Types

| Type | Code | Description |
|------|------|-------------|
| Fixed Fee | `fixed` | Flat rate regardless of quantity/weight |
| Per-Item | `by_count` | Fee based on item quantity |
| Per-Weight | `by_weight` | Fee based on total weight (grams) |
| Free Shipping | `free` | No shipping fee |

### 2.4 Free Shipping Threshold

- **By Amount:** Free shipping when order total reaches threshold
- **By Count:** Free shipping when item quantity reaches threshold

### 2.5 Template Association

Products and categories can be associated with specific shipping templates. Priority order:
1. Product-specific template
2. Category template
3. Shop default template

---

## 3. Data Model

### 3.1 Entity Relationship

```
shipping_templates (1) ──< (N) shipping_zones
shipping_templates (1) ──< (N) shipping_template_mappings
```

### 3.2 Table Definitions

#### shipping_templates

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| name | VARCHAR(100) | Template name |
| is_default | BOOLEAN | Is default template |
| is_active | BOOLEAN | Is enabled |
| created_at | TIMESTAMP | Created timestamp |
| updated_at | TIMESTAMP | Updated timestamp |
| deleted_at | TIMESTAMP | Soft delete timestamp |

**Indexes:**
- `idx_tenant_id` on `tenant_id`
- `idx_tenant_default` on `(tenant_id, is_default)` for quick default lookup

#### shipping_zones

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| template_id | BIGINT | Parent template ID |
| name | VARCHAR(100) | Zone name (e.g., "East China") |
| regions | JSONB | City code array, e.g., `["330100","330200"]` |
| fee_type | VARCHAR(20) | Fee type: fixed / by_count / by_weight / free |
| first_unit | INT | First unit (count or grams) |
| first_fee | BIGINT | First fee (cents) |
| additional_unit | INT | Additional unit |
| additional_fee | BIGINT | Additional fee (cents) |
| free_threshold_amount | BIGINT | Free shipping threshold amount (cents), 0 = disabled |
| free_threshold_count | INT | Free shipping threshold count, 0 = disabled |
| sort | INT | Sort order (matching priority) |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

**Indexes:**
- `idx_template_id` on `template_id`
- `idx_tenant_id` on `tenant_id`

**JSONB regions format:**
```json
["330100", "330200", "330300"]
```
City codes follow Chinese administrative division standards (6-digit codes).

#### shipping_template_mappings

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| template_id | BIGINT | Shipping template ID |
| target_type | VARCHAR(20) | product / category |
| target_id | BIGINT | Product ID or Category ID |

**Indexes:**
- `idx_template_id` on `template_id`
- `idx_target` on `(target_type, target_id)` for quick lookup

### 3.3 Reused Tables

The `carriers` table from Fulfillment module is reused for shipping company management. Shipping module provides an independent management entry point.

---

## 4. API Design

### 4.1 Admin API

#### Shipping Template Management

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/shipping-templates` | List templates |
| POST | `/api/v1/shipping-templates` | Create template |
| GET | `/api/v1/shipping-templates/:id` | Get template detail (with zones and mappings) |
| PUT | `/api/v1/shipping-templates/:id` | Update template |
| DELETE | `/api/v1/shipping-templates/:id` | Delete template |
| PUT | `/api/v1/shipping-templates/:id/set-default` | Set as default |

**List Templates Response:**
```json
{
  "list": [
    {
      "id": 1,
      "name": "National Free Shipping Template",
      "is_default": true,
      "is_active": true,
      "zone_count": 3,
      "product_count": 0,
      "category_count": 2,
      "created_at": "2026-03-25T10:00:00Z"
    }
  ],
  "total": 5
}
```

#### Shipping Zone Management

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/shipping-templates/:id/zones` | Add zone |
| PUT | `/api/v1/shipping-zones/:id` | Update zone |
| DELETE | `/api/v1/shipping-zones/:id` | Delete zone |
| PUT | `/api/v1/shipping-templates/:id/zones/reorder` | Reorder zones |

**Zone Request/Response:**
```json
{
  "id": 1,
  "template_id": 1,
  "name": "East China",
  "regions": ["330100", "330200", "310100"],
  "fee_type": "by_weight",
  "first_unit": 1000,
  "first_fee": "1000",
  "additional_unit": 500,
  "additional_fee": "200",
  "free_threshold_amount": "9900",
  "free_threshold_count": 0,
  "sort": 1
}
```

#### Template Mapping Management

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/shipping-templates/:id/mappings` | Get mappings |
| POST | `/api/v1/shipping-template-mappings` | Create mapping |
| DELETE | `/api/v1/shipping-template-mappings/:id` | Delete mapping |

**Mapping Request:**
```json
{
  "template_id": 1,
  "target_type": "product",
  "target_id": 123
}
```

#### Shipping Fee Calculator (Admin Testing Tool)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/shipping/calculate` | Calculate shipping fee |

**Calculate Request:**
```json
{
  "address": {
    "province_code": "330000",
    "city_code": "330100",
    "district_code": "330102"
  },
  "items": [
    {
      "product_id": 123,
      "sku_id": 456,
      "quantity": 2,
      "weight": 500,
      "price": "99.00"
    }
  ]
}
```

**Calculate Response:**
```json
{
  "shipping_fee": "1200",
  "currency": "CNY",
  "template_id": 1,
  "template_name": "National Free Shipping Template",
  "zone_name": "East China",
  "fee_detail": {
    "fee_type": "by_weight",
    "first_unit": 1000,
    "first_fee": "1000",
    "additional_unit": 500,
    "additional_fee": "200",
    "calculated_weight": 1000,
    "calculated_units": 1
  }
}
```

#### Carrier Management (Shared with Fulfillment)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/carriers` | List carriers |
| POST | `/api/v1/carriers` | Create carrier |
| PUT | `/api/v1/carriers/:id` | Update carrier |
| DELETE | `/api/v1/carriers/:id` | Delete carrier |

---

## 5. Frontend Design

### 5.1 Page Structure

```
/settings/shipping
├── /templates           # Template list
├── /templates/:id       # Template detail (edit zones and mappings)
└── /test                # Shipping fee calculator tool
```

### 5.2 Page 1: Shipping Template List

**Features:**
- Display all templates as cards
- Identify default template
- Create, edit, delete, set as default
- Show associated product/category count

**List Fields:**

| Field | Description |
|-------|-------------|
| Template Name | e.g., "National Free Shipping Template" |
| Default Badge | Is default template |
| Zone Count | Number of zones configured |
| Associated Products | Product/category count |
| Status | Enabled/Disabled |
| Actions | Edit, Delete, Set Default |

### 5.3 Page 2: Template Detail Page

**Layout:**
- **Header:** Template name, default toggle
- **Tab 1 - Zones:** Shipping zone configuration
- **Tab 2 - Associations:** Product/category associations

**Zone Configuration:**
- Zone name
- Region selector (city-level multi-select)
- Fee type selector
- Fee configuration fields (dynamic based on fee type)
- Free shipping threshold settings

**Region Selector Interaction:**
- Left panel: Province tree (expandable to city level)
- Right panel: Selected regions list
- "Select all province" quick action

**Association Configuration:**
- Product selector (search and multi-select)
- Category selector (tree multi-select)

### 5.4 Page 3: Shipping Fee Calculator

**Features:**
- Address selector (province-city-district cascader)
- Add test items (product selector, quantity input)
- Real-time fee calculation
- Display matched template and zone info

---

## 6. Business Rules

### 6.1 Template Rules

1. Each tenant can have multiple templates
2. Each tenant must have exactly one default template
3. Cannot delete default template (must set another as default first)
4. Deleting a template removes all associated zones and mappings

### 6.2 Zone Matching Rules

1. Zones are matched in sort order (lower sort = higher priority)
2. First matching zone wins (based on city code)
3. Unmatched addresses use default template's "default zone"

### 6.3 Fee Calculation Logic

**Fixed Fee:**
```
shipping_fee = first_fee
```

**Per-Item:**
```
if quantity <= first_unit:
    shipping_fee = first_fee
else:
    additional_units = ceil((quantity - first_unit) / additional_unit)
    shipping_fee = first_fee + additional_units * additional_fee
```

**Per-Weight:**
```
if weight <= first_unit:
    shipping_fee = first_fee
else:
    additional_units = ceil((weight - first_unit) / additional_unit)
    shipping_fee = first_fee + additional_units * additional_fee
```

**Free Shipping:**
```
shipping_fee = 0
```

**Free Shipping Threshold Check:**
```
if order_amount >= free_threshold_amount OR item_count >= free_threshold_count:
    shipping_fee = 0
```

### 6.4 Template Association Priority

1. Product-specific template (highest priority)
2. Category template
3. Shop default template (lowest priority)

---

## 7. Implementation Phases

### Phase 1: Backend Core (MVP)

- [ ] Create database migrations
- [ ] Implement domain entities
- [ ] Implement repository layer
- [ ] Implement application services
- [ ] Create Admin API endpoints
- [ ] Implement fee calculation logic

### Phase 2: Frontend Management

- [ ] Template list page
- [ ] Template detail page
- [ ] Zone configuration (with region selector)
- [ ] Association management
- [ ] Shipping fee calculator tool

### Phase 3: Shop API Integration

- [ ] Shop API for checkout calculation
- [ ] Cart preview integration
- [ ] Checkout page integration

---

## 8. Error Codes

Error code range for Shipping module: **180xxx**

| Code | HTTP Status | Description |
|------|-------------|-------------|
| 180001 | 400 | Template name required |
| 180002 | 400 | Template not found |
| 180003 | 400 | Cannot delete default template |
| 180004 | 400 | Zone not found |
| 180005 | 400 | Invalid region codes |
| 180006 | 400 | Fee configuration invalid |
| 180007 | 400 | No matching zone for address |
| 180008 | 400 | Product already associated |

---

## 9. Technical Notes

### 9.1 Region Data

City-level administrative division data should be stored in a reference table or loaded from a static JSON file. Recommended source: National Bureau of Statistics administrative division codes.

### 9.2 Performance Considerations

- Cache default template per tenant
- Cache zone configurations per template
- Consider batch fee calculation for cart scenarios

### 9.3 Future Enhancements

- Shipping fee promotions (free shipping coupons)
- Time-based shipping fee rules
- Multi-warehouse shipping fee calculation
- International shipping support