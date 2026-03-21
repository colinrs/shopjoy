-- ============================================
-- 购物车表 (carts)
-- ============================================

CREATE TABLE IF NOT EXISTS `carts` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT DEFAULT NULL COMMENT '用户ID',
    `session_id` VARCHAR(255) DEFAULT '' COMMENT '会话ID(未登录用户)',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user` (`tenant_id`, `user_id`),
    UNIQUE KEY `uk_tenant_session` (`tenant_id`, `session_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';

-- ============================================
-- 购物车商品表 (cart_items)
-- ============================================

CREATE TABLE IF NOT EXISTS `cart_items` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `cart_id` BIGINT NOT NULL COMMENT '购物车ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
    `product_name` VARCHAR(255) NOT NULL COMMENT '商品名称',
    `sku_name` VARCHAR(255) DEFAULT '' COMMENT 'SKU名称',
    `image` VARCHAR(500) DEFAULT '' COMMENT '图片',
    `price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价(分)',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
    `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '小计(分)',
    `selected` TINYINT NOT NULL DEFAULT 1 COMMENT '是否选中',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_cart_id` (`cart_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_sku_id` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车商品表';

-- ============================================
-- 测试数据
-- ============================================

-- 购物车数据
INSERT INTO `carts` (`id`, `tenant_id`, `user_id`, `session_id`, `updated_at`) VALUES
(1, 1, 1, '', UNIX_TIMESTAMP()),
(2, 1, 2, '', UNIX_TIMESTAMP()),
(3, 1, 3, '', UNIX_TIMESTAMP()),
(4, 2, 4, '', UNIX_TIMESTAMP()),
(5, 3, 6, '', UNIX_TIMESTAMP()),
(6, 3, 7, '', UNIX_TIMESTAMP());

-- 购物车商品数据 (Demo Shop - 用户1)
INSERT INTO `cart_items` (`id`, `tenant_id`, `user_id`, `cart_id`, `product_id`, `sku_id`, `product_name`, `sku_name`, `image`, `price`, `quantity`, `total_amount`, `selected`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, 1, 1, 1, 2, 'Nike Air Max 270', '黑色 43码', 'https://cdn.example.com/p1-1.jpg', 129900, 1, 129900, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),
(2, 1, 1, 1, 3, 7, 'iPhone 15 手机壳', '黑色', 'https://cdn.example.com/p3-1.jpg', 9900, 2, 19800, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),
(3, 1, 1, 1, 4, 0, 'MacBook 充电器', '', 'https://cdn.example.com/p4-1.jpg', 79900, 1, 79900, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),

-- 购物车商品数据 (Demo Shop - 用户2)
(4, 1, 2, 2, 2, 5, 'Adidas Ultraboost 22', '黑色 42码', 'https://cdn.example.com/p2-1.jpg', 159900, 1, 159900, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, 2, 2, 5, 0, '简约台灯', '', 'https://cdn.example.com/p5-1.jpg', 29900, 2, 59800, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- 购物车商品数据 (Demo Shop - 用户3)
(6, 1, 3, 3, 1, 3, 'Nike Air Max 270', '白色 42码', 'https://cdn.example.com/p1-1.jpg', 129900, 1, 129900, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 3, 3),

-- 购物车商品数据 (Enterprise Corp - 用户6)
(7, 3, 6, 5, 3, 8, 'iPhone 15 手机壳', '白色', 'https://cdn.example.com/p3-1.jpg', 9900, 1, 9900, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 6, 6),
(8, 3, 6, 5, 4, 0, 'MacBook 充电器', '', 'https://cdn.example.com/p4-1.jpg', 79900, 1, 79900, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 6, 6);