-- ============================================
-- 评价表 (reviews)
-- ============================================

CREATE TABLE `reviews` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `order_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(64) NOT NULL DEFAULT '',
    `user_id` BIGINT NOT NULL,
    `user_name` VARCHAR(100) NOT NULL,
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

-- ============================================
-- 评价回复表 (review_replies)
-- ============================================

CREATE TABLE `review_replies` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `review_id` BIGINT NOT NULL,
    `tenant_id` BIGINT NOT NULL,
    `admin_id` BIGINT NOT NULL,
    `admin_name` VARCHAR(100) NOT NULL,
    `content` TEXT NOT NULL,
    `created_at` BIGINT NOT NULL,
    `updated_at` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_review_id` (`review_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 评价统计表 (review_stats)
-- ============================================

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
