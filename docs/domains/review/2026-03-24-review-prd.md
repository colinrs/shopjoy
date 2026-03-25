# Review System Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Review System PRD |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-23 |
| Author | Product Team |
| Module Location | admin/internal/domain/review/ |

---

## Product Overview

### Summary

This PRD describes the Review System for ShopJoy, enabling merchants to manage product reviews from buyers. The review system is order-based, meaning buyers can only review products they have purchased, ensuring review authenticity and reducing fake reviews. The system includes multi-dimensional ratings, automatic moderation, merchant reply capability, and helpful voting features.

### Problem Statement

Merchants currently lack a systematic way to collect and manage product reviews. Without a review system:

- **No social proof**: Products lack customer feedback to help potential buyers make purchasing decisions
- **No quality feedback loop**: Merchants cannot gather structured feedback on product quality and value
- **No engagement opportunity**: Merchants cannot respond to reviews to address concerns or thank customers
- **Moderation burden**: Without automatic moderation rules, managing reviews becomes time-consuming
- **Trust deficit**: Buyers cannot verify if reviews are from actual purchasers

### Solution Overview

Implement a comprehensive review management system with the following key features:

1. **Order-based reviews** - Buyers can only review after order completion (verified purchases)
2. **Multi-dimensional ratings** - Quality (1-5) and Value for Money (1-5) ratings
3. **Automatic moderation** - 4-5 star reviews auto-approve; 1-3 stars go to pending for review
4. **Merchant reply** - Single reply per review with edit/delete capability
5. **Anonymous option** - Buyers can choose to hide their identity
6. **Helpful voting** - Simple thumbs-up mechanism for review quality
7. **Featured reviews** - Merchants can highlight valuable reviews
8. **Product statistics** - Aggregated rating counts and averages per product

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | Increase conversion rate | Reviews displayed on product pages | 30% of products have reviews |
| BG-002 | Reduce moderation workload | Auto-approval rate | 70% of reviews auto-approved |
| BG-003 | Improve merchant engagement | Reply rate | 80% of reviews have merchant reply |
| BG-004 | Build trust | Verified purchase badge display | 100% of reviews are verified |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | Manage pending reviews efficiently | Merchant Admin |
| UG-002 | Respond to customer feedback | Merchant Admin / Customer Service |
| UG-003 | View review statistics | Merchant Admin / Store Manager |
| UG-004 | Feature helpful reviews | Merchant Admin |
| UG-005 | Moderate inappropriate reviews | Merchant Admin |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | Review editing by buyer | Prevents review manipulation, maintains integrity |
| NG-002 | Buyer-side review submission UI | Belongs to shop service |
| NG-003 | Review reply notification | Phase 2 feature |
| NG-004 | Review report/flagging | Phase 2 feature |
| NG-005 | AI-powered moderation | Phase 2 feature |
| NG-006 | Review sentiment analysis | Phase 2 feature |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | Store owner or manager with full access | Approve/hide reviews, reply to reviews, feature reviews, view statistics |
| Customer Service | Staff handling customer interactions | Reply to reviews, add notes |
| Operations Manager | Staff managing day-to-day operations | View statistics, batch operations |

### Role-Based Access

| Role | View Reviews | Approve/Hide | Reply | Feature | Batch Operations | View Stats |
|------|--------------|--------------|-------|---------|------------------|------------|
| Tenant Admin | Yes | Yes | Yes | Yes | Yes | Yes |
| Tenant Operations Manager | Yes | Yes | Yes | Yes | Yes | Yes |
| Tenant Customer Service | Yes | No | Yes | No | No | Yes |

---

## Functional Requirements

### Review Management

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | List reviews | GET /api/v1/reviews with pagination and filters | P0 |
| FR-002 | Get review detail | GET /api/v1/reviews/:id with full information | P0 |
| FR-003 | Approve review | PUT /api/v1/reviews/:id/approve for pending reviews | P0 |
| FR-004 | Hide review | PUT /api/v1/reviews/:id/hide to hide from storefront | P0 |
| FR-005 | Show review | PUT /api/v1/reviews/:id/show to unhide hidden review | P0 |
| FR-006 | Delete review | DELETE /api/v1/reviews/:id for soft delete | P0 |
| FR-007 | Feature review | PUT /api/v1/reviews/:id/featured to toggle featured status | P1 |
| FR-008 | Batch approve | POST /api/v1/reviews/batch-approve for multiple reviews | P1 |
| FR-009 | Batch hide | POST /api/v1/reviews/batch-hide for multiple reviews | P1 |

### Merchant Reply

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-010 | Create reply | POST /api/v1/reviews/:id/reply with content | P0 |
| FR-011 | Edit reply | PUT /api/v1/reviews/:id/reply to update content | P0 |
| FR-012 | Delete reply | DELETE /api/v1/reviews/:id/reply to remove reply | P0 |
| FR-013 | Single reply constraint | Only one reply allowed per review | P0 |
| FR-014 | Reply length limit | Maximum 500 characters | P0 |

### Statistics and Reporting

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | Review statistics | GET /api/v1/reviews/stats for overall summary | P0 |
| FR-021 | Product statistics | GET /api/v1/reviews/product/:product_id/stats | P0 |
| FR-022 | Rating distribution | Count of reviews per rating level (1-5 stars) | P0 |
| FR-023 | Reply rate calculation | Percentage of reviews with merchant reply | P1 |

### Filtering and Search

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-030 | Filter by status | pending, approved, hidden | P0 |
| FR-031 | Filter by product | Filter reviews for specific product | P0 |
| FR-032 | Filter by rating | Rating range filter (min/max) | P0 |
| FR-033 | Filter by images | Show only reviews with images | P1 |
| FR-034 | Keyword search | Search in review content | P1 |
| FR-035 | Date range filter | Filter by creation date | P1 |

---

## User Experience

### Entry Points

1. **Main navigation** - Reviews menu item in admin sidebar
2. **Product detail page** - Link to product's reviews
3. **Dashboard** - Pending reviews count widget
4. **Order detail page** - Link to review if exists

### Core Experience: Review List Page

1. Admin navigates to Reviews section from sidebar
2. System displays statistics cards at top:
   - Total Reviews count
   - Pending Approval count
   - Average Rating
   - Reviews with Images count
3. Admin can see filter bar with options:
   - Search keyword
   - Product filter
   - Status filter (All / Pending / Approved / Hidden)
   - Rating filter (min-max range)
   - Has image checkbox
   - Date range picker
4. Review table shows:
   - Product name and image
   - Reviewer name (or "Anonymous")
   - Quality and Value ratings (star display)
   - Content preview (truncated)
   - Image thumbnails
   - Status badge
   - Helpful count
   - Reply status
   - Created time
   - Actions dropdown
5. Actions available:
   - View Detail
   - Approve (for pending)
   - Hide/Show
   - Feature/Unfeature
   - Reply

### Core Experience: Review Detail Dialog

1. Admin clicks on a review row or "View Detail" action
2. Dialog displays full review information:
   - Order information (Order ID, Product name, SKU code)
   - Reviewer information (Name, Anonymous badge, Verified badge)
   - Rating display:
     - Quality: 5 stars
     - Value for Money: 4 stars
     - Overall: 4.5
   - Full review content
   - Images gallery (clickable for lightbox)
   - Statistics (Helpful count, Created time)
   - Status badge
   - Featured badge (if applicable)
3. Merchant reply section:
   - If no reply: Reply form with character counter
   - If has reply: Display reply with edit/delete options
4. Status management buttons:
   - Approve (for pending)
   - Hide/Show toggle
   - Delete

### Core Experience: Reply to Review

1. Admin clicks "Reply" button on review
2. Reply dialog opens:
   - Display review content (reference)
   - Text area for reply (max 500 chars)
   - Character counter
   - Cancel and Submit buttons
3. Admin types reply and submits
4. System validates:
   - Content not empty
   - Content within limit
   - Review not hidden
   - No existing reply (or editing existing)
5. Success: Reply saved, displayed in review detail
6. Error: Display appropriate error message

### Advanced Features

**Batch Operations**

1. Admin selects multiple reviews via checkboxes
2. Batch action buttons appear: "Approve Selected", "Hide Selected"
3. Admin clicks desired action
4. Confirmation dialog shows count of affected reviews
5. Admin confirms, batch operation executes
6. Success toast displays operation result

**Featured Reviews**

1. Admin identifies a high-quality, helpful review
2. Clicks "Feature" action
3. Review gets featured badge
4. Featured reviews are highlighted in storefront

---

## Narrative

Merchant admin Sarah opens the ShopJoy admin panel on Monday morning. She sees a notification badge on the Reviews menu - 8 new pending reviews from the weekend.

Clicking into Reviews, she sees the statistics dashboard: 1,250 total reviews, 15 pending approval, 4.3 average rating. She notices the pending reviews and filters to show only "Pending" status.

Sarah reviews each pending review carefully. Most are 3-star reviews with constructive feedback. She reads a review mentioning the product arrived slightly damaged during shipping. She clicks "Reply" and writes a response apologizing and offering a discount on the next purchase. She approves the review.

Another review has inappropriate language. Sarah clicks "Hide" to remove it from public view while preserving it for records.

Moving to approved reviews, she notices a 5-star review with detailed photos and helpful feedback. She clicks "Feature" to highlight this review on the product page.

Before lunch, Sarah uses the batch approve function to approve 5 pending reviews that all look legitimate, improving her efficiency.

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Moderation time | Average time to approve pending reviews | < 24 hours |
| Reply rate | Percentage of reviews with merchant reply | 80% |
| Featured reviews | Percentage of products with featured reviews | 20% |

### Business Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Review submission rate | Percentage of completed orders with reviews | 30% |
| Average rating | Overall average rating across products | > 4.0 |
| Auto-approval rate | Percentage of reviews auto-approved | 70% |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API response time | List reviews endpoint P95 | < 200ms |
| API response time | Review detail endpoint P95 | < 100ms |
| Error rate | All review operations | < 0.1% |

---

## Technical Considerations

### Integration Points

| Integration | Description |
|-------------|-------------|
| Order Service | Verify order completion status, retrieve order items |
| Product Service | Retrieve product information, update product statistics |
| User Service | Retrieve user information for display |
| File Storage | Store review images |

### Data Storage and Privacy

- Reviews are stored with soft delete support (deleted_at field)
- Anonymous reviews hide user identity but retain user_id for fraud detection
- Images stored in object storage with CDN delivery
- Statistics are denormalized for performance

### Scalability and Performance

- Review statistics table (review_stats) provides fast aggregated queries
- Indexes on tenant_id, product_id, status for efficient filtering
- Pagination required for all list endpoints
- Consider caching product statistics with TTL

### Potential Challenges

1. **Concurrent replies**: Multiple admins attempting to reply simultaneously
   - Solution: Database constraint ensures single reply per review
2. **Statistics consistency**: Stats can become out of sync
   - Solution: Event-driven recalculation or periodic sync job
3. **Large image uploads**: Performance impact
   - Solution: Client-side compression, size limits, async upload

---

## API Endpoints

### Review Management

#### GET /api/v1/reviews - List Reviews

**Query Parameters:**
```
page: 1                          // Page number (default: 1)
page_size: 20                    // Items per page (default: 20, max: 100)
product_id: 123                  // Filter by product ID (optional)
status: pending                  // Filter by status: pending/approved/hidden (optional)
rating_min: 1                    // Minimum rating (optional)
rating_max: 5                    // Maximum rating (optional)
has_image: true                  // Only reviews with images (optional)
keyword: great                   // Search in content (optional)
start_time: 2026-03-01T00:00:00Z // Start date filter (optional)
end_time: 2026-03-23T23:59:59Z   // End date filter (optional)
```

**Response:**
```json
{
  "list": [
    {
      "id": 1,
      "order_id": 123,
      "product_id": 456,
      "product_name": "Wireless Headphones",
      "sku_code": "WH-BLACK",
      "user_name": "John D.",
      "is_anonymous": false,
      "is_verified": true,
      "quality_rating": 5,
      "value_rating": 4,
      "overall_rating": "4.50",
      "content": "Great product! Sound quality is excellent.",
      "images": ["https://cdn.example.com/review/1.jpg", "https://cdn.example.com/review/2.jpg"],
      "status": "approved",
      "is_featured": false,
      "helpful_count": 12,
      "has_reply": true,
      "created_at": "2026-03-18T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

#### GET /api/v1/reviews/:id - Get Review Detail

**Response:**
```json
{
  "id": 1,
  "tenant_id": 1,
  "order_id": 123,
  "product_id": 456,
  "product_name": "Wireless Headphones",
  "sku_code": "WH-BLACK",
  "user_id": 789,
  "user_name": "John Doe",
  "is_anonymous": false,
  "is_verified": true,
  "quality_rating": 5,
  "value_rating": 4,
  "overall_rating": "4.50",
  "content": "Great product! Sound quality is excellent. Battery life is amazing, lasting over 20 hours on a single charge. The noise cancellation feature works perfectly for my commute.",
  "images": ["https://cdn.example.com/review/1.jpg", "https://cdn.example.com/review/2.jpg"],
  "status": "approved",
  "is_featured": false,
  "helpful_count": 12,
  "created_at": "2026-03-18T10:00:00Z",
  "updated_at": "2026-03-18T10:00:00Z",
  "reply": {
    "id": 1,
    "content": "Thank you for your feedback! We're glad you enjoyed our product.",
    "admin_name": "Store Manager",
    "created_at": "2026-03-18T12:00:00Z",
    "updated_at": "2026-03-18T12:00:00Z"
  }
}
```

#### PUT /api/v1/reviews/:id/approve - Approve Review

**Response:**
```json
{
  "id": 1,
  "status": "approved",
  "updated_at": "2026-03-23T14:00:00Z"
}
```

#### PUT /api/v1/reviews/:id/hide - Hide Review

**Request:**
```json
{
  "reason": "Inappropriate content"
}
```

**Response:**
```json
{
  "id": 1,
  "status": "hidden",
  "updated_at": "2026-03-23T14:00:00Z"
}
```

#### PUT /api/v1/reviews/:id/show - Show Hidden Review

**Response:**
```json
{
  "id": 1,
  "status": "approved",
  "updated_at": "2026-03-23T14:00:00Z"
}
```

#### DELETE /api/v1/reviews/:id - Soft Delete Review

**Response:**
```json
{
  "id": 1,
  "status": "deleted",
  "deleted_at": "2026-03-23T14:00:00Z"
}
```

#### PUT /api/v1/reviews/:id/featured - Toggle Featured Status

**Request:**
```json
{
  "is_featured": true
}
```

**Response:**
```json
{
  "id": 1,
  "is_featured": true,
  "updated_at": "2026-03-23T14:00:00Z"
}
```

### Merchant Reply

#### POST /api/v1/reviews/:id/reply - Create Reply

**Request:**
```json
{
  "content": "Thank you for your feedback! We're glad you enjoyed our product."
}
```

**Response:**
```json
{
  "id": 1,
  "review_id": 123,
  "content": "Thank you for your feedback! We're glad you enjoyed our product.",
  "admin_id": 1,
  "admin_name": "Store Manager",
  "created_at": "2026-03-23T14:00:00Z"
}
```

#### PUT /api/v1/reviews/:id/reply - Edit Reply

**Request:**
```json
{
  "content": "Updated reply content here."
}
```

**Response:**
```json
{
  "id": 1,
  "review_id": 123,
  "content": "Updated reply content here.",
  "admin_id": 1,
  "admin_name": "Store Manager",
  "created_at": "2026-03-23T14:00:00Z",
  "updated_at": "2026-03-23T15:00:00Z"
}
```

#### DELETE /api/v1/reviews/:id/reply - Delete Reply

**Response:**
```json
{
  "success": true
}
```

### Statistics

#### GET /api/v1/reviews/stats - Overall Statistics

**Response:**
```json
{
  "total_reviews": 1250,
  "pending_reviews": 15,
  "approved_reviews": 1200,
  "hidden_reviews": 35,
  "average_rating": "4.30",
  "quality_avg_rating": "4.25",
  "value_avg_rating": "4.35",
  "five_star_count": 800,
  "four_star_count": 300,
  "three_star_count": 100,
  "two_star_count": 30,
  "one_star_count": 20,
  "with_image_count": 450,
  "reply_rate": 0.85,
  "featured_count": 50
}
```

#### GET /api/v1/reviews/product/:product_id/stats - Product Statistics

**Response:**
```json
{
  "product_id": 456,
  "total_reviews": 150,
  "average_rating": "4.50",
  "quality_avg_rating": "4.55",
  "value_avg_rating": "4.45",
  "rating_distribution": {
    "5": 100,
    "4": 35,
    "3": 10,
    "2": 3,
    "1": 2
  },
  "with_image_count": 45,
  "reply_count": 120,
  "reply_rate": 0.80
}
```

### Batch Operations

#### POST /api/v1/reviews/batch-approve - Batch Approve

**Request:**
```json
{
  "ids": [1, 2, 3, 4, 5]
}
```

**Response:**
```json
{
  "success_count": 5,
  "failed_count": 0,
  "errors": []
}
```

#### POST /api/v1/reviews/batch-hide - Batch Hide

**Request:**
```json
{
  "ids": [1, 2, 3],
  "reason": "Inappropriate content"
}
```

**Response:**
```json
{
  "success_count": 3,
  "failed_count": 0,
  "errors": []
}
```

---

## Database Schema

### reviews Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO | Primary key |
| tenant_id | BIGINT | NO | - | Tenant isolation |
| order_id | BIGINT | NO | - | Linked order |
| product_id | BIGINT | NO | - | Product being reviewed |
| sku_code | VARCHAR(64) | NO | '' | Specific SKU code |
| user_id | BIGINT | NO | - | Reviewer user ID |
| quality_rating | TINYINT | NO | - | Quality rating (1-5) |
| value_rating | TINYINT | NO | - | Value for money rating (1-5) |
| overall_rating | DECIMAL(3,2) | NO | - | Calculated average rating |
| content | TEXT | NO | '' | Review text content |
| images | JSON | YES | NULL | Array of image URLs |
| status | TINYINT | NO | 0 | 0=pending, 1=approved, 2=hidden, 3=deleted |
| is_anonymous | BOOLEAN | NO | FALSE | Anonymous review flag |
| is_verified | BOOLEAN | NO | FALSE | Verified purchase flag |
| is_featured | BOOLEAN | NO | FALSE | Featured review flag |
| helpful_count | INT | NO | 0 | Helpful votes count |
| created_at | BIGINT | NO | - | Unix timestamp |
| updated_at | BIGINT | NO | - | Unix timestamp |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**
- `idx_tenant_product` (tenant_id, product_id)
- `idx_tenant_user` (tenant_id, user_id)
- `idx_order_id` (order_id)
- `idx_status` (status)
- `idx_product_status` (product_id, status)

**Migration SQL:**
```sql
CREATE TABLE `reviews` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `order_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(64) NOT NULL DEFAULT '',
    `user_id` BIGINT NOT NULL,
    `quality_rating` TINYINT NOT NULL,
    `value_rating` TINYINT NOT NULL,
    `overall_rating` DECIMAL(3,2) NOT NULL,
    `content` TEXT NOT NULL,
    `images` JSON NULL,
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending,1=approved,2=hidden,3=deleted',
    `is_anonymous` BOOLEAN NOT NULL DEFAULT FALSE,
    `is_verified` BOOLEAN NOT NULL DEFAULT FALSE,
    `is_featured` BOOLEAN NOT NULL DEFAULT FALSE,
    `helpful_count` INT NOT NULL DEFAULT 0,
    `created_at` BIGINT NOT NULL,
    `updated_at` BIGINT NOT NULL,
    `deleted_at` BIGINT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_product` (`tenant_id`, `product_id`),
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_order_id` (`order_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_product_status` (`product_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### review_replies Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO | Primary key |
| review_id | BIGINT | NO | - | Linked review |
| tenant_id | BIGINT | NO | - | Tenant isolation |
| admin_id | BIGINT | NO | - | Admin who replied |
| content | TEXT | NO | - | Reply text content |
| created_at | BIGINT | NO | - | Unix timestamp |
| updated_at | BIGINT | NO | - | Unix timestamp |

**Indexes:**
- `idx_review_id` (review_id) - UNIQUE

**Migration SQL:**
```sql
CREATE TABLE `review_replies` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `review_id` BIGINT NOT NULL,
    `tenant_id` BIGINT NOT NULL,
    `admin_id` BIGINT NOT NULL,
    `content` TEXT NOT NULL,
    `created_at` BIGINT NOT NULL,
    `updated_at` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_review_id` (`review_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### review_stats Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO | Primary key |
| tenant_id | BIGINT | NO | - | Tenant isolation |
| product_id | BIGINT | NO | - | Product ID |
| total_reviews | INT | NO | 0 | Total review count |
| average_rating | DECIMAL(3,2) | NO | 0.00 | Overall average rating |
| quality_avg_rating | DECIMAL(3,2) | NO | 0.00 | Quality dimension average |
| value_avg_rating | DECIMAL(3,2) | NO | 0.00 | Value dimension average |
| rating_1_count | INT | NO | 0 | 1-star review count |
| rating_2_count | INT | NO | 0 | 2-star review count |
| rating_3_count | INT | NO | 0 | 3-star review count |
| rating_4_count | INT | NO | 0 | 4-star review count |
| rating_5_count | INT | NO | 0 | 5-star review count |
| with_image_count | INT | NO | 0 | Reviews with images count |
| last_updated_at | BIGINT | NO | - | Last stats update timestamp |

**Unique Index:** `uk_tenant_product` (tenant_id, product_id)

**Migration SQL:**
```sql
CREATE TABLE `review_stats` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `total_reviews` INT NOT NULL DEFAULT 0,
    `average_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `quality_avg_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `value_avg_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `rating_1_count` INT NOT NULL DEFAULT 0,
    `rating_2_count` INT NOT NULL DEFAULT 0,
    `rating_3_count` INT NOT NULL DEFAULT 0,
    `rating_4_count` INT NOT NULL DEFAULT 0,
    `rating_5_count` INT NOT NULL DEFAULT 0,
    `with_image_count` INT NOT NULL DEFAULT 0,
    `last_updated_at` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_tenant_product` (`tenant_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

---

## Business Rules

### Review Creation Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | Order completion required | User can only review after order status is `completed` |
| BR-002 | One review per order item | Each product/SKU in an order can only be reviewed once |
| BR-003 | Rating required | Both Quality and Value ratings are required (1-5) |
| BR-004 | Content optional | Text content is optional but recommended |
| BR-005 | Image limit | Maximum 5 images per review |
| BR-006 | Content length | Maximum 1000 characters for review content |
| BR-007 | Verified purchase | All reviews are marked as verified purchase |

### Moderation Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-010 | Auto-approve positive | Reviews with overall rating >= 4 are auto-approved |
| BR-011 | Pending for negative | Reviews with overall rating < 4 go to `pending` status |
| BR-012 | Hidden reviews not visible | Reviews with `hidden` status are not shown on storefront |
| BR-013 | Deleted is soft delete | Deleted reviews are retained in database with `deleted_at` set |
| BR-014 | Status transition | Pending -> Approved/Hidden; Approved -> Hidden; Hidden -> Approved |

### Reply Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-020 | Single reply | Only one merchant reply allowed per review |
| BR-021 | Edit allowed | Merchant can edit their existing reply |
| BR-022 | Delete allowed | Merchant can delete their reply |
| BR-023 | Content length | Maximum 500 characters for reply content |
| BR-024 | Cannot reply to hidden | Merchant cannot reply to hidden reviews |
| BR-025 | Reply un-hides review | If review was hidden, replying does not auto-unhide |

### Featured Review Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-030 | Approved only | Only approved reviews can be featured |
| BR-031 | Toggle operation | Featured status can be toggled on/off |
| BR-032 | No limit | No limit on number of featured reviews per product |

### Helpful Count Rules

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-040 | Increment only | Helpful count can only be incremented (from shop service) |
| BR-041 | No decrement | Users cannot remove their helpful vote |

---

## Error Codes

Review module uses error code range 210xxx:

| Constant | HTTP Status | Code | Message |
|----------|-------------|------|---------|
| ErrReviewNotFound | 404 | 210001 | review not found |
| ErrReviewAlreadyReplied | 400 | 210002 | review already has reply |
| ErrReviewCannotReplyHidden | 400 | 210003 | cannot reply to hidden review |
| ErrReviewInvalidStatus | 400 | 210004 | invalid review status |
| ErrReviewContentTooLong | 400 | 210005 | review content exceeds limit |
| ErrReplyContentTooLong | 400 | 210006 | reply content exceeds limit |
| ErrReplyNotFound | 404 | 210007 | reply not found |
| ErrReviewCannotApprove | 400 | 210008 | cannot approve review in current status |
| ErrReviewCannotHide | 400 | 210009 | cannot hide review in current status |
| ErrReviewCannotShow | 400 | 210010 | cannot show review in current status |
| ErrReviewCannotFeature | 400 | 210011 | can only feature approved reviews |
| ErrReviewAlreadyDeleted | 400 | 210012 | review already deleted |
| ErrReviewReplyEmpty | 400 | 210013 | reply content cannot be empty |
| ErrReviewInvalidRating | 400 | 210014 | rating must be between 1 and 5 |
| ErrReviewBatchEmpty | 400 | 210015 | batch operation requires at least one review id |
| ErrReviewBatchLimitExceeded | 400 | 210016 | batch operation limited to 100 reviews |

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Description |
|-------|----------|-------------|
| Phase 1: Database setup | 0.5 day | Create tables, indexes, migration scripts |
| Phase 2: Backend API | 2 days | Domain entities, repository, service, handlers |
| Phase 3: Admin Frontend | 2 days | List page, detail dialog, reply dialog, batch operations |
| Phase 4: Testing & Integration | 0.5 day | Unit tests, integration tests, bug fixes |

**Total Estimate: 5 working days**

### Suggested Phases

#### Phase 1: Database Setup (Day 1 morning)

- Create migration scripts for reviews, review_replies, review_stats tables
- Define indexes and constraints
- Run migrations on development environment

#### Phase 2: Backend API (Day 1 afternoon - Day 2)

- Create domain entities (Review, ReviewReply, ReviewStats)
- Define repository interfaces
- Implement infrastructure layer (repository implementation)
- Implement application service layer
- Create API handlers for all endpoints
- Add error codes to pkg/code/code.go
- Run `make build` to verify compilation

#### Phase 3: Admin Frontend (Day 3 - Day 4)

- Create review list page with statistics cards
- Implement filter bar and data table
- Create review detail dialog
- Implement merchant reply dialog
- Add batch operations UI
- Style and polish

#### Phase 4: Testing & Integration (Day 5 morning)

- Write unit tests for domain logic
- Write integration tests for API endpoints
- Perform E2E testing
- Bug fixes and refinements

---

## User Stories

### US-001: View Review List

**ID**: US-001

**Title**: View Review List

**Description**: As a merchant admin, I want to see a list of all reviews so that I can manage customer feedback efficiently.

**Acceptance Criteria**:
- Given I am logged in as a merchant admin
- When I navigate to the Reviews section
- Then I see a paginated list of reviews
- And I see statistics cards showing total reviews, pending count, average rating, and image count
- And I can see each review's product name, rating, content preview, status, and creation time

---

### US-002: Filter Reviews by Status

**ID**: US-002

**Title**: Filter Reviews by Status

**Description**: As a merchant admin, I want to filter reviews by status so that I can focus on specific types of reviews.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I select "Pending" from the status filter
- Then I only see reviews with pending status
- When I select "Approved" from the status filter
- Then I only see reviews with approved status
- When I select "Hidden" from the status filter
- Then I only see reviews with hidden status

---

### US-003: Search Reviews by Keyword

**ID**: US-003

**Title**: Search Reviews by Keyword

**Description**: As a merchant admin, I want to search reviews by keyword so that I can find specific feedback.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I type "quality" in the search box
- And I press Enter or click Search
- Then I see only reviews containing "quality" in their content

---

### US-004: View Review Detail

**ID**: US-004

**Title**: View Review Detail

**Description**: As a merchant admin, I want to see the full details of a review so that I can understand the complete feedback.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I click on a review row
- Then a detail dialog opens
- And I see the full review content
- And I see all images in a gallery
- And I see the order information (order ID, product name, SKU)
- And I see the reviewer information (name, anonymous status, verified badge)
- And I see both rating dimensions (Quality, Value) and overall rating
- And I see helpful count and creation time

---

### US-005: Approve Pending Review

**ID**: US-005

**Title**: Approve Pending Review

**Description**: As a merchant admin, I want to approve pending reviews so that they become visible on the storefront.

**Acceptance Criteria**:
- Given I am viewing a pending review
- When I click the "Approve" button
- Then the review status changes to "approved"
- And the review becomes visible on the storefront
- And I see a success message

---

### US-006: Hide Inappropriate Review

**ID**: US-006

**Title**: Hide Inappropriate Review

**Description**: As a merchant admin, I want to hide inappropriate reviews so that they are not visible to customers.

**Acceptance Criteria**:
- Given I am viewing an approved review
- When I click the "Hide" button
- And optionally enter a reason
- Then the review status changes to "hidden"
- And the review is no longer visible on the storefront
- And I see a success message

---

### US-007: Show Hidden Review

**ID**: US-007

**Title**: Show Hidden Review

**Description**: As a merchant admin, I want to restore hidden reviews so that they become visible again.

**Acceptance Criteria**:
- Given I am viewing a hidden review
- When I click the "Show" button
- Then the review status changes to "approved"
- And the review becomes visible on the storefront
- And I see a success message

---

### US-008: Delete Review

**ID**: US-008

**Title**: Delete Review

**Description**: As a merchant admin, I want to delete reviews so that they are removed from the system.

**Acceptance Criteria**:
- Given I am viewing a review
- When I click the "Delete" button
- And I confirm the deletion in the confirmation dialog
- Then the review is soft-deleted (status = deleted)
- And the review is removed from the list
- And I see a success message

---

### US-009: Reply to Review

**ID**: US-009

**Title**: Reply to Review

**Description**: As a merchant admin, I want to reply to reviews so that I can address customer feedback.

**Acceptance Criteria**:
- Given I am viewing a review without a reply
- When I click the "Reply" button
- And I enter my reply text
- And I click "Submit"
- Then the reply is saved and associated with the review
- And the review detail shows my reply
- And I see a success message

---

### US-010: Edit Merchant Reply

**ID**: US-010

**Title**: Edit Merchant Reply

**Description**: As a merchant admin, I want to edit my reply so that I can correct or update my response.

**Acceptance Criteria**:
- Given I am viewing a review with my reply
- When I click the "Edit" button on my reply
- And I modify the reply text
- And I click "Save"
- Then the reply is updated
- And the updated_at timestamp is refreshed
- And I see a success message

---

### US-011: Delete Merchant Reply

**ID**: US-011

**Title**: Delete Merchant Reply

**Description**: As a merchant admin, I want to delete my reply so that it is removed from the review.

**Acceptance Criteria**:
- Given I am viewing a review with my reply
- When I click the "Delete" button on my reply
- And I confirm the deletion
- Then the reply is removed
- And the review no longer shows a reply
- And I see a success message

---

### US-012: Cannot Reply to Hidden Review

**ID**: US-012

**Title**: Cannot Reply to Hidden Review

**Description**: As a merchant admin, I should not be able to reply to hidden reviews.

**Acceptance Criteria**:
- Given I am viewing a hidden review
- When I try to click the "Reply" button
- Then the button is disabled or not shown
- Or if I attempt to reply via API
- Then I receive error code 210003 "cannot reply to hidden review"

---

### US-013: Single Reply Constraint

**ID**: US-013

**Title**: Single Reply Constraint

**Description**: As a system, ensure only one reply exists per review.

**Acceptance Criteria**:
- Given a review already has a reply
- When I attempt to create another reply
- Then I receive error code 210002 "review already has reply"
- And the existing reply remains unchanged

---

### US-014: Reply Length Validation

**ID**: US-014

**Title**: Reply Length Validation

**Description**: As a merchant admin, I want to be prevented from entering replies that are too long.

**Acceptance Criteria**:
- Given I am writing a reply
- When I enter more than 500 characters
- Then I see a character counter showing I exceeded the limit
- And the "Submit" button is disabled
- And if I submit via API
- Then I receive error code 210006 "reply content exceeds limit"

---

### US-015: Feature a Review

**ID**: US-015

**Title**: Feature a Review

**Description**: As a merchant admin, I want to mark a review as featured so that it is highlighted on the storefront.

**Acceptance Criteria**:
- Given I am viewing an approved review
- When I click the "Feature" button
- Then the review is marked as featured (is_featured = true)
- And a featured badge appears on the review
- And I see a success message

---

### US-016: Unfeature a Review

**ID**: US-016

**Title**: Unfeature a Review

**Description**: As a merchant admin, I want to remove the featured status from a review.

**Acceptance Criteria**:
- Given I am viewing a featured review
- When I click the "Unfeature" button
- Then the featured status is removed (is_featured = false)
- And the featured badge is removed
- And I see a success message

---

### US-017: Batch Approve Reviews

**ID**: US-017

**Title**: Batch Approve Reviews

**Description**: As a merchant admin, I want to approve multiple reviews at once to save time.

**Acceptance Criteria**:
- Given I am on the Reviews list page with multiple pending reviews
- When I select multiple reviews using checkboxes
- And I click the "Batch Approve" button
- And I confirm the action
- Then all selected reviews are approved
- And I see a success message with the count of approved reviews

---

### US-018: Batch Hide Reviews

**ID**: US-018

**Title**: Batch Hide Reviews

**Description**: As a merchant admin, I want to hide multiple reviews at once to save time.

**Acceptance Criteria**:
- Given I am on the Reviews list page with multiple reviews
- When I select multiple reviews using checkboxes
- And I click the "Batch Hide" button
- And I optionally enter a reason
- And I confirm the action
- Then all selected reviews are hidden
- And I see a success message with the count of hidden reviews

---

### US-019: View Review Statistics

**ID**: US-019

**Title**: View Review Statistics

**Description**: As a merchant admin, I want to see overall review statistics to understand customer sentiment.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- Then I see statistics cards showing:
  - Total reviews count
  - Pending reviews count
  - Average rating (overall)
  - Reviews with images count
- And I can see rating distribution (1-5 stars)

---

### US-020: View Product-Specific Statistics

**ID**: US-020

**Title**: View Product-Specific Statistics

**Description**: As a merchant admin, I want to see review statistics for a specific product.

**Acceptance Criteria**:
- Given I request GET /api/v1/reviews/product/:product_id/stats
- Then I see product-specific statistics including:
  - Total reviews for the product
  - Average rating (overall, quality, value)
  - Rating distribution (1-5 stars)
  - Reviews with images count
  - Reply rate

---

### US-021: Filter Reviews by Product

**ID**: US-021

**Title**: Filter Reviews by Product

**Description**: As a merchant admin, I want to filter reviews for a specific product.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I select a product from the product filter dropdown
- Then I only see reviews for that product
- And the statistics update to reflect the filtered results

---

### US-022: Filter Reviews by Rating

**ID**: US-022

**Title**: Filter Reviews by Rating

**Description**: As a merchant admin, I want to filter reviews by rating range.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I set rating min to 4
- Then I only see reviews with overall rating >= 4
- When I set rating max to 2
- Then I only see reviews with overall rating <= 2
- When I set rating min to 3 and max to 4
- Then I only see reviews with overall rating between 3 and 4

---

### US-023: Filter Reviews with Images

**ID**: US-023

**Title**: Filter Reviews with Images

**Description**: As a merchant admin, I want to filter reviews that have images.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I check the "Has Images" checkbox
- Then I only see reviews that have at least one image

---

### US-024: Filter Reviews by Date Range

**ID**: US-024

**Title**: Filter Reviews by Date Range

**Description**: As a merchant admin, I want to filter reviews by creation date range.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I set start date to 2026-03-01
- And I set end date to 2026-03-15
- Then I only see reviews created between those dates

---

### US-025: Review Not Found

**ID**: US-025

**Title**: Review Not Found

**Description**: As a system, handle the case when a requested review does not exist.

**Acceptance Criteria**:
- Given I request a review with ID 9999 that does not exist
- When I make the API request
- Then I receive error code 210001 "review not found"
- And the HTTP status is 404

---

### US-026: Invalid Status Transition

**ID**: US-026

**Title**: Invalid Status Transition

**Description**: As a system, prevent invalid status transitions.

**Acceptance Criteria**:
- Given a review with status "hidden"
- When I try to approve it directly
- Then I should first show it (change to approved)
- Or I receive an appropriate error
- Given a review with status "deleted"
- When I try to approve or hide it
- Then I receive error code 210004 "invalid review status"

---

### US-027: Anonymous Review Display

**ID**: US-027

**Title**: Anonymous Review Display

**Description**: As a merchant admin, I want to see that a review is anonymous while still having access to reviewer info for fraud detection.

**Acceptance Criteria**:
- Given an anonymous review
- When I view the review in the list
- Then I see "Anonymous" or masked name for the reviewer
- When I view the review detail as admin
- Then I can see the is_anonymous flag is true
- And I can see the actual user_id for fraud detection purposes

---

### US-028: Verified Purchase Badge

**ID**: US-028

**Title**: Verified Purchase Badge

**Description**: As a merchant admin, I want to see that a review is from a verified purchase.

**Acceptance Criteria**:
- Given a review from a completed order
- When I view the review
- Then I see a "Verified Purchase" badge
- And the is_verified flag is true in the API response

---

### US-029: View Review Images

**ID**: US-029

**Title**: View Review Images

**Description**: As a merchant admin, I want to view images attached to reviews in a gallery format.

**Acceptance Criteria**:
- Given a review with images
- When I view the review detail
- Then I see thumbnail images
- When I click on an image
- Then a lightbox opens showing the full-size image
- And I can navigate between images

---

### US-030: Access Control for Review Operations

**ID**: US-030

**Title**: Access Control for Review Operations

**Description**: As a system, enforce role-based access control for review operations.

**Acceptance Criteria**:
- Given I am logged in as Tenant Customer Service
- When I try to hide a review
- Then the action is denied (Customer Service cannot hide)
- And I see an appropriate permission error
- Given I am logged in as Tenant Admin
- When I perform any review operation
- Then the action is allowed

---

### US-031: Pagination for Review List

**ID**: US-031

**Title**: Pagination for Review List

**Description**: As a merchant admin, I want to paginate through reviews to handle large numbers of reviews.

**Acceptance Criteria**:
- Given there are more than 20 reviews
- When I view the Reviews list
- Then I see pagination controls
- And I see 20 reviews per page by default
- When I click "Next"
- Then I see the next page of reviews
- When I change page_size to 50
- Then I see 50 reviews per page

---

### US-032: Reply Empty Validation

**ID**: US-032

**Title**: Reply Empty Validation

**Description**: As a system, prevent empty reply submissions.

**Acceptance Criteria**:
- Given I am replying to a review
- When I submit a reply with empty content
- Then I receive error code 210013 "reply content cannot be empty"
- And the reply is not saved

---

### US-033: Batch Operation Limit

**ID**: US-033

**Title**: Batch Operation Limit

**Description**: As a system, limit the number of reviews in a batch operation to prevent performance issues.

**Acceptance Criteria**:
- Given I select more than 100 reviews for batch operation
- When I submit the batch request
- Then I receive error code 210016 "batch operation limited to 100 reviews"
- And no reviews are modified

---

### US-034: Feature Only Approved Reviews

**ID**: US-034

**Title**: Feature Only Approved Reviews

**Description**: As a system, prevent featuring non-approved reviews.

**Acceptance Criteria**:
- Given a review with status "pending"
- When I try to feature it
- Then I receive error code 210011 "can only feature approved reviews"
- And the review is not featured

---

### US-035: Review Sorting

**ID**: US-035

**Title**: Review Sorting

**Description**: As a merchant admin, I want to sort reviews by different criteria.

**Acceptance Criteria**:
- Given I am on the Reviews list page
- When I click the "Created Time" column header
- Then reviews are sorted by creation time (descending by default)
- When I click the column header again
- Then the sort order is reversed
- When I click the "Rating" column header
- Then reviews are sorted by overall rating

---

### US-036: Rating Calculation

**ID**: US-036

**Title**: Rating Calculation

**Description**: As a system, calculate the overall rating as the average of quality and value ratings.

**Acceptance Criteria**:
- Given a review with quality_rating = 5 and value_rating = 4
- When the review is created
- Then overall_rating is calculated as 4.50
- Given a review with quality_rating = 3 and value_rating = 3
- When the review is created
- Then overall_rating is calculated as 3.00

---

### US-037: Auto-Approve High Rating Reviews

**ID**: US-037

**Title**: Auto-Approve High Rating Reviews

**Description**: As a system, automatically approve reviews with high ratings.

**Acceptance Criteria**:
- Given a new review with overall_rating >= 4.0
- When the review is created
- Then the review status is automatically set to "approved"
- Given a new review with overall_rating < 4.0
- When the review is created
- Then the review status is set to "pending"

---

### US-038: Product Statistics Update

**ID**: US-038

**Title**: Product Statistics Update

**Description**: As a system, update product review statistics when reviews are created, approved, or deleted.

**Acceptance Criteria**:
- Given a product with 10 reviews
- When a new review is approved
- Then the product's total_reviews increments
- And the product's average_rating is recalculated
- And the rating distribution is updated
- When a review is deleted
- Then the product's total_reviews decrements
- And all statistics are recalculated

---

### US-039: Reply Rate Calculation

**ID**: US-039

**Title**: Reply Rate Calculation

**Description**: As a merchant admin, I want to see the reply rate to understand merchant engagement.

**Acceptance Criteria**:
- Given 100 approved reviews
- And 85 reviews have merchant replies
- When I view the statistics
- Then I see reply_rate = 0.85 (85%)

---

### US-040: Concurrent Reply Prevention

**ID**: US-040

**Title**: Concurrent Reply Prevention

**Description**: As a system, prevent concurrent reply creation through database constraints.

**Acceptance Criteria**:
- Given two admins try to reply to the same review simultaneously
- When both requests reach the database
- Then only one reply is saved (due to unique constraint on review_id)
- And the other request receives error 210002 "review already has reply"