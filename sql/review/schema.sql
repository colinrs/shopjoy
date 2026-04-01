-- ============================================
-- 评价表 (reviews)
-- ============================================

CREATE TABLE `reviews` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `order_id` VARCHAR(64) NOT NULL,
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
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
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
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
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
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_tenant_product` (`tenant_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 测试数据
-- ============================================

-- 评价数据 (Demo Shop)
INSERT INTO `reviews` (`id`, `tenant_id`, `order_id`, `product_id`, `sku_code`, `user_id`, `user_name`, `quality_rating`, `value_rating`, `overall_rating`, `content`, `images`, `status`, `is_anonymous`, `is_verified`, `is_featured`, `helpful_count`, `created_at`) VALUES
-- ORD202503010001 订单的评价
(1, 1, 'ORD202503010001', 1, 'SKU-001-BLK-42', 1, '小明', 5, 5, 5.00, '鞋子非常舒适，尺码标准，物流也很快！', '["https://cdn.example.com/review1-1.jpg", "https://cdn.example.com/review1-2.jpg"]', 1, FALSE, TRUE, TRUE, 15, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY))),
(2, 1, 'ORD202503010001', 3, 'SKU-003-BLK', 1, '小明', 4, 4, 4.00, '手机壳质量不错，就是颜色有点色差', '["https://cdn.example.com/review2-1.jpg"]', 1, FALSE, TRUE, FALSE, 8, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 24 DAY))),
(3, 1, 'ORD202503010001', 5, 'SKU-005', 1, '小明', 5, 5, 5.00, '台灯很漂亮，装点效果很好', NULL, 1, FALSE, TRUE, FALSE, 5, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 23 DAY))),

-- ORD202503100001 订单的评价
(4, 1, 'ORD202503100001', 2, 'SKU-002-BLK-42', 2, '小红', 5, 5, 5.00, 'Adidas的鞋就是好，穿着很舒服，尺码也很准', NULL, 1, FALSE, TRUE, TRUE, 12, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY))),

-- ORD202503120001 订单的评价 (Enterprise)
(5, 3, 'ORD202503120001', 4, 'SKU-004', 6, '约翰', 4, 5, 4.50, '充电器做工不错，就是发热有点大', NULL, 1, FALSE, TRUE, FALSE, 3, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY)));

-- 评价回复数据
INSERT INTO `review_replies` (`id`, `review_id`, `tenant_id`, `admin_id`, `admin_name`, `content`, `created_at`) VALUES
(1, 1, 1, 2, 'Demo管理员', '感谢您的好评！我们会继续努力提供更好的服务！', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 24 DAY))),
(2, 4, 1, 2, 'Demo管理员', '感谢您的认可，欢迎下次再来！', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)));

-- 评价统计数据
INSERT INTO `review_stats` (`id`, `tenant_id`, `product_id`, `total_reviews`, `average_rating`, `quality_avg_rating`, `value_avg_rating`, `rating_1_count`, `rating_2_count`, `rating_3_count`, `rating_4_count`, `rating_5_count`, `with_image_count`, `created_at`) VALUES
-- Demo Shop 产品统计
(1, 1, 1, 1, 5.00, 5.00, 5.00, 0, 0, 0, 0, 1, 1, UNIX_TIMESTAMP()),
(2, 1, 2, 1, 5.00, 5.00, 5.00, 0, 0, 0, 0, 1, 0, UNIX_TIMESTAMP()),
(3, 1, 3, 1, 4.00, 4.00, 4.00, 0, 0, 0, 1, 0, 1, UNIX_TIMESTAMP()),
(4, 1, 5, 1, 5.00, 5.00, 5.00, 0, 0, 0, 0, 1, 0, UNIX_TIMESTAMP()),
-- Enterprise 产品统计
(5, 3, 4, 1, 4.50, 4.00, 5.00, 0, 0, 0, 1, 0, 0, UNIX_TIMESTAMP());
