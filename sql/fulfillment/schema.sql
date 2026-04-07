-- ============================================
-- 物流发货表 (shipments)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '发货ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_id` BIGINT NOT NULL COMMENT '订单ID',
    `shipment_no` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '发货单号',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待发货, 1-已发货, 2-运输中, 3-已送达, 4-配送失败',
    `carrier` VARCHAR(50) DEFAULT '' COMMENT '快递公司',
    `carrier_code` VARCHAR(20) DEFAULT '' COMMENT '快递公司代码',
    `tracking_no` VARCHAR(100) DEFAULT '' COMMENT '快递单号',
    `weight` DECIMAL(10,3) DEFAULT 0.000 COMMENT '重量(kg)',
    `cost_amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '运费成本',
    `cost_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `remark` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '备注',
    `shipped_at` TIMESTAMP NULL COMMENT '发货时间',
    `delivered_at` TIMESTAMP NULL COMMENT '送达时间',
    `cancelled_at` TIMESTAMP NULL COMMENT '取消时间',
    `cancelled_by` BIGINT NOT NULL DEFAULT 0 COMMENT '取消人',
    `cancelled_reason` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '取消原因',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_shipment_no` (`shipment_no`),
    KEY `idx_tracking_no` (`tracking_no`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `uk_shipment_no` (`tenant_id`, `shipment_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='物流发货表';

-- ============================================
-- 发货商品表 (shipment_items)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipment_items` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `shipment_id` BIGINT NOT NULL COMMENT '发货ID',
    `order_item_id` BIGINT NOT NULL COMMENT '订单商品ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
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
    `order_id` BIGINT NOT NULL COMMENT '订单ID',
    `refund_no` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '退款单号',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '退款类型: 1-全额退款, 2-部分退款',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待审核, 1-已批准, 2-已拒绝, 3-已完成, 4-已取消',
    `reason_type` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '退款原因类型',
    `reason` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '退款原因',
    `description` TEXT COMMENT '详细描述',
    `images` JSON COMMENT '凭证图片',
    `amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '退款金额',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `reject_reason` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '拒绝原因',
    `approved_at` TIMESTAMP NULL COMMENT '批准时间',
    `approved_by` BIGINT NOT NULL DEFAULT 0 COMMENT '批准人',
    `completed_at` TIMESTAMP NULL COMMENT '完成时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_refund_no` (`refund_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `uk_refund_no` (`tenant_id`, `refund_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款表';

-- ============================================
-- 运费模板表 (shipping_templates)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipping_templates` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `name` VARCHAR(100) NOT NULL COMMENT 'Template name',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT 'Is default template (0=no, 1=yes)',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Is active (0=inactive, 1=active)',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping templates table';

-- ============================================
-- 运费区域表 (shipping_zones)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipping_zones` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `template_id` BIGINT NOT NULL COMMENT 'Template ID',
    `name` VARCHAR(100) NOT NULL COMMENT 'Zone name',
    `regions` JSON NOT NULL COMMENT 'Region codes (city codes array)',
    `fee_type` VARCHAR(20) NOT NULL COMMENT 'Fee type: fixed, by_count, by_weight, free',
    `first_unit` INT NOT NULL DEFAULT 1 COMMENT 'First unit (count or grams)',
    `first_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'First fee',
    `additional_unit` INT NOT NULL DEFAULT 1 COMMENT 'Additional unit',
    `additional_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Additional fee',
    `free_threshold_amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Free shipping threshold amount, 0=disabled',
    `free_threshold_count` INT NOT NULL DEFAULT 0 COMMENT 'Free shipping threshold count, 0=disabled',
    `sort` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping zones table';

-- ============================================
-- 运费模板映射表 (shipping_template_mappings)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipping_template_mappings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `template_id` BIGINT NOT NULL COMMENT 'Template ID',
    `target_type` VARCHAR(20) NOT NULL COMMENT 'Target type: product, category',
    `target_id` BIGINT NOT NULL COMMENT 'Target ID (product_id or category_id)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_template_id` (`template_id`),
    INDEX `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping template mappings table';

-- ============================================
-- 运费区域城市关联表 (shipping_zone_regions)
-- ============================================

CREATE TABLE IF NOT EXISTS `shipping_zone_regions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `zone_id` BIGINT NOT NULL COMMENT 'Zone ID',
    `city_code` VARCHAR(20) NOT NULL COMMENT 'City code',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_zone_city` (`zone_id`, `city_code`),
    INDEX `idx_city_code` (`city_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping zone regions junction table';

-- ============================================
-- 测试数据
-- ============================================

-- 发货数据
INSERT IGNORE INTO `shipments` (`id`, `tenant_id`, `order_id`, `status`, `carrier`, `tracking_no`, `weight`, `cost_amount`, `cost_currency`, `shipped_at`, `delivered_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- Order 1 发货
(1, 1, 1, 3, '顺丰速运', 'SF1234567890', 1.50, 1200, 'CNY', DATE_SUB(NOW(), INTERVAL 28 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 29 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY), 2, 2),

-- Order 2 发货
(2, 1, 2, 2, '圆通快递', 'YT9876543210', 0.45, 800, 'CNY', DATE_SUB(NOW(), INTERVAL 2 DAY), NULL, DATE_SUB(NOW(), INTERVAL 4 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 2, 2),

-- Order 6 发货 (Enterprise)
(3, 3, 6, 3, 'EMS', 'EMS2025031201', 0.55, 1500, 'CNY', DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 9 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 5, 5);

-- 发货商品数据
INSERT IGNORE INTO `shipment_items` (`id`, `tenant_id`, `shipment_id`, `order_item_id`, `product_id`, `sku_id`, `quantity`) VALUES
(1, 1, 1, 1, 1, 1, 1),
(2, 1, 1, 2, 3, 7, 3),
(3, 1, 1, 3, 5, 11, 1),
(4, 1, 2, 4, 2, 5, 1),
(5, 3, 3, 8, 4, 10, 1);

-- 退款数据
INSERT IGNORE INTO `refunds` (`id`, `tenant_id`, `order_id`, `refund_no`, `user_id`, `type`, `status`, `reason_type`, `reason`, `description`, `images`, `amount`, `currency`, `approved_at`, `completed_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- Order 1 手机壳退款 (已完成)
(1, 1, 1, 'REF202503010001', 1, 2, 3, 'quality', '商品质量问题', '手机壳有划痕，申请退款', '["https://cdn.example.com/refund1.jpg", "https://cdn.example.com/refund2.jpg"]', 29700, 'CNY', DATE_SUB(NOW(), INTERVAL 26 DAY), DATE_SUB(NOW(), INTERVAL 25 DAY), DATE_SUB(NOW(), INTERVAL 27 DAY), DATE_SUB(NOW(), INTERVAL 25 DAY), 1, 1),

-- Order 5 全额退款 (已完成)
(2, 1, 5, 'REF202503050001', 1, 1, 3, 'cancel', '不想要了', '取消订单，申请退款', '[]', 29900, 'CNY', DATE_SUB(NOW(), INTERVAL 19 DAY), DATE_SUB(NOW(), INTERVAL 19 DAY), DATE_SUB(NOW(), INTERVAL 19 DAY), DATE_SUB(NOW(), INTERVAL 19 DAY), 1, 1);
