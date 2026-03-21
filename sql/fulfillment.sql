-- ============================================
-- 物流发货表 (shipments)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '发货ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待发货, 1-已发货, 2-运输中, 3-已送达, 4-配送失败',
    `carrier` VARCHAR(50) DEFAULT '' COMMENT '快递公司',
    `tracking_no` VARCHAR(100) DEFAULT '' COMMENT '快递单号',
    `weight` DECIMAL(10,2) DEFAULT 0.00 COMMENT '重量(kg)',
    `cost_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '运费成本(分)',
    `cost_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `shipped_at` BIGINT DEFAULT NULL COMMENT '发货时间',
    `delivered_at` BIGINT DEFAULT NULL COMMENT '送达时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_tracking_no` (`tracking_no`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='物流发货表';

-- ============================================
-- 发货商品表 (shipment_items)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipment_items` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `shipment_id` BIGINT NOT NULL COMMENT '发货ID',
    `order_item_id` BIGINT NOT NULL COMMENT '订单商品ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
    PRIMARY KEY (`id`),
    KEY `idx_shipment_id` (`shipment_id`),
    KEY `idx_order_item_id` (`order_item_id`),
    KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发货商品表';

-- ============================================
-- 退款表 (refunds)
-- ============================================

CREATE TABLE IF NOT EXISTS `refunds` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '退款ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待审核, 1-已批准, 2-已拒绝, 3-已完成',
    `reason` VARCHAR(500) DEFAULT '' COMMENT '退款原因',
    `description` TEXT COMMENT '详细描述',
    `images` JSON COMMENT '凭证图片',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '退款金额(分)',
    `currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `approved_at` BIGINT DEFAULT NULL COMMENT '批准时间',
    `completed_at` BIGINT DEFAULT NULL COMMENT '完成时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款表';

-- ============================================
-- 测试数据
-- ============================================

-- 发货数据
INSERT INTO `shipments` (`id`, `tenant_id`, `order_id`, `status`, `carrier`, `tracking_no`, `weight`, `cost_amount`, `cost_currency`, `shipped_at`, `delivered_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- ORD202503010001 发货
(1, 1, 'ORD202503010001', 3, '顺丰速运', 'SF1234567890', 1.50, 1200, 'CNY', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 28 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 29 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), 2, 2),

-- ORD202503100001 发货
(2, 1, 'ORD202503100001', 2, '圆通快递', 'YT9876543210', 0.45, 800, 'CNY', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), 2, 2),

-- ORD202503120001 发货 (Enterprise)
(3, 3, 'ORD202503120001', 3, 'EMS', 'EMS2025031201', 0.55, 1500, 'CNY', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 7 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY)), 5, 5);

-- 发货商品数据
INSERT INTO `shipment_items` (`id`, `shipment_id`, `order_item_id`, `product_id`, `sku_id`, `quantity`) VALUES
(1, 1, 1, 1, 1, 1),
(2, 1, 2, 3, 7, 3),
(3, 1, 3, 5, 0, 1),
(4, 2, 4, 2, 5, 1),
(5, 3, 8, 4, 0, 1);

-- 退款数据
INSERT INTO `refunds` (`id`, `tenant_id`, `order_id`, `user_id`, `status`, `reason`, `description`, `images`, `amount`, `currency`, `approved_at`, `completed_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- ORD202503010001 手机壳退款 (已完成)
(1, 1, 'ORD202503010001', 1, 3, '商品质量问题', '手机壳有划痕，申请退款', '["https://cdn.example.com/refund1.jpg", "https://cdn.example.com/refund2.jpg"]', 29700, 'CNY', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 26 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 27 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY)), 1, 1),

-- ORD202503050001 全额退款 (已完成)
(2, 1, 'ORD202503050001', 1, 3, '不想要了', '取消订单，申请退款', '[]', 29900, 'CNY', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), 1, 1);