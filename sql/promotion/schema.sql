-- ============================================
-- 促销活动表 (promotions)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '促销ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '促销名称',
    `description` TEXT COMMENT '描述',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-折扣, 1-限时抢购, 2-捆绑销售, 3-买X送Y',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待生效, 1-生效中, 2-已暂停, 3-已结束',
    `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `scope_type` VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' COMMENT '范围类型',
    `scope_ids` JSON DEFAULT NULL COMMENT '范围ID列表',
    `exclude_ids` JSON DEFAULT NULL COMMENT '排除ID列表',
    `start_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_type` (`type`),
    KEY `idx_priority` (`priority`),
    KEY `idx_start_end` (`start_at`, `end_at`),
    KEY `idx_active` (`status`, `currency`, `start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销活动表';

-- ============================================
-- 促销规则表 (promotion_rules)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotion_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `promotion_id` BIGINT NOT NULL COMMENT '促销ID',
    `condition_type` TINYINT NOT NULL DEFAULT 0 COMMENT '条件类型: 0-最低金额, 1-最低数量',
    `condition_value` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '条件值',
    `action_type` TINYINT NOT NULL DEFAULT 0 COMMENT '动作类型: 0-固定金额, 1-百分比',
    `action_value` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '动作值',
    `max_discount_amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '最大优惠金额',
    `max_discount_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_promotion_id` (`promotion_id`),
    KEY `idx_promotion_rules_sort_order` (`promotion_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销规则表';

-- ============================================
-- 促销使用记录表 (promotion_usage)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotion_usage` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `promotion_id` BIGINT NOT NULL,
    `rule_id` BIGINT DEFAULT NULL,
    `order_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `discount_amount` DECIMAL(19,4) NOT NULL DEFAULT 0,
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY',
    `original_amount` DECIMAL(19,4) NOT NULL DEFAULT 0,
    `final_amount` DECIMAL(19,4) NOT NULL DEFAULT 0,
    `coupon_id` BIGINT DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_promotion_id` (`promotion_id`),
    INDEX `idx_order_id` (`order_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销使用记录表';

-- ============================================
-- 优惠券表 (coupons)
-- ============================================

CREATE TABLE IF NOT EXISTS `coupons` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '名称',
    `code` VARCHAR(100) NOT NULL COMMENT '优惠券代码',
    `description` TEXT COMMENT '描述',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-固定金额, 1-百分比, 2-免邮',
    `value` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '优惠值',
    `min_amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '最低消费',
    `max_discount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '最大优惠',
    `total_count` INT NOT NULL DEFAULT 0 COMMENT '发放总数',
    `used_count` INT NOT NULL DEFAULT 0 COMMENT '已使用数量',
    `per_user_limit` INT NOT NULL DEFAULT 1 COMMENT '每用户限领',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未激活, 1-激活, 2-过期, 3-用完',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `scope_type` VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' COMMENT '范围类型',
    `scope_ids` JSON DEFAULT NULL COMMENT '范围ID列表',
    `exclude_ids` JSON DEFAULT NULL COMMENT '排除ID列表',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `start_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_start_end` (`start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券表';

-- ============================================
-- 用户优惠券表 (user_coupons)
-- ============================================

CREATE TABLE IF NOT EXISTS `user_coupons` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `coupon_id` BIGINT NOT NULL COMMENT '优惠券ID',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未使用, 1-已使用, 2-已过期',
    `used_at` TIMESTAMP NULL COMMENT '使用时间',
    `order_id` BIGINT DEFAULT 0 COMMENT '订单ID',
    `received_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
    `expire_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '过期时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_coupon_id` (`coupon_id`),
    KEY `idx_status` (`status`),
    KEY `idx_expire_at` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券表';

-- ============================================
-- 测试数据
-- ============================================

-- 促销活动数据 (Demo Shop)
INSERT IGNORE INTO `promotions` (`id`, `tenant_id`, `name`, `description`, `type`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '春季大促', '春季全场促销活动', 0, 1, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), NOW(), NOW(), 2, 2),
(2, 1, '限时抢购-鞋类专场', '运动鞋限时抢购', 1, 1, DATE_ADD(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 8 DAY), NOW(), NOW(), 2, 2),
(3, 1, '买二送一', '购买任意两件商品送一件同等价值商品', 3, 1, NOW(), DATE_ADD(NOW(), INTERVAL 60 DAY), NOW(), NOW(), 2, 2),
(4, 1, '会员日特惠', '每月18号会员专享折扣', 0, 2, DATE_ADD(NOW(), INTERVAL 10 DAY), DATE_ADD(NOW(), INTERVAL 40 DAY), NOW(), NOW(), 2, 2),
(5, 1, '清仓大甩卖', '库存清仓特价', 1, 1, NOW(), DATE_ADD(NOW(), INTERVAL 15 DAY), NOW(), NOW(), 2, 2),

-- Enterprise Corp 促销活动
(6, 3, '企业采购季', '企业客户采购优惠季', 0, 1, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY), NOW(), NOW(), 5, 5),
(7, 3, '跨境特惠', '跨境电商专属优惠', 0, 1, NOW(), DATE_ADD(NOW(), INTERVAL 180 DAY), NOW(), NOW(), 5, 5);

-- 促销规则数据
INSERT IGNORE INTO `promotion_rules` (`id`, `promotion_id`, `condition_type`, `condition_value`, `action_type`, `action_value`, `max_discount_amount`, `max_discount_currency`) VALUES
-- 春季大促规则
(1, 1, 0, 10000, 1, 10, 5000, 'CNY'),   -- 满100减10%, 最多减50
(2, 1, 0, 30000, 1, 15, 10000, 'CNY'),  -- 满300减15%, 最多减100
(3, 1, 0, 50000, 1, 20, 20000, 'CNY'),  -- 满500减20%, 最多减200

-- 限时抢购规则
(4, 2, 0, 0, 1, 30, 10000, 'CNY'),      -- 无门槛7折, 最多减100

-- 买二送一规则
(5, 3, 1, 2, 0, 0, 0, 'CNY'),           -- 买2件送1件

-- 会员日特惠规则
(6, 4, 0, 20000, 1, 25, 15000, 'CNY'),  -- 满200减25%, 最多减150

-- 清仓大甩卖规则
(7, 5, 0, 0, 1, 40, 20000, 'CNY'),      -- 无门槛6折, 最多减200

-- 企业采购季规则
(8, 6, 0, 50000, 1, 15, 30000, 'CNY'),  -- 满500减15%, 最多减300
(9, 6, 0, 100000, 1, 20, 50000, 'CNY'), -- 满1000减20%, 最多减500

-- 跨境特惠规则
(10, 7, 0, 30000, 1, 10, 20000, 'CNY'); -- 满300减10%, 最多减200

-- 优惠券数据 (Demo Shop)
INSERT IGNORE INTO `coupons` (`id`, `tenant_id`, `name`, `code`, `description`, `type`, `value`, `min_amount`, `max_discount`, `total_count`, `used_count`, `per_user_limit`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '新用户专享', 'NEWUSER50', '新用户首单立减50元', 0, 5000, 10000, 0, 1000, 150, 1, 1, DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_ADD(NOW(), INTERVAL 365 DAY), NOW(), NOW(), 2, 2),
(2, 1, '全场9折', 'SAVE10', '全场商品9折优惠', 1, 10, 5000, 10000, 500, 80, 3, 1, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY), NOW(), NOW(), 2, 2),
(3, 1, '满200减30', 'SAVE30', '满200元减30元', 0, 3000, 20000, 0, 200, 45, 2, 1, NOW(), DATE_ADD(NOW(), INTERVAL 60 DAY), NOW(), NOW(), 2, 2),
(4, 1, '免邮券', 'FREESHIP', '全场免运费', 2, 0, 5000, 1500, 300, 20, 1, 1, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), NOW(), NOW(), 2, 2),
(5, 1, 'VIP专享8折', 'VIP20', 'VIP会员专享8折', 1, 20, 10000, 20000, 100, 10, 1, 1, NOW(), DATE_ADD(NOW(), INTERVAL 180 DAY), NOW(), NOW(), 2, 2),

-- Enterprise Corp 优惠券
(6, 3, '企业专属优惠', 'ENT50', '企业客户专享50元优惠', 0, 5000, 20000, 0, 500, 30, 5, 1, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY), NOW(), NOW(), 5, 5),
(7, 3, '跨境免邮', 'GLOBALSHIP', '跨境电商免邮', 2, 0, 30000, 5000, 200, 15, 3, 1, NOW(), DATE_ADD(NOW(), INTERVAL 180 DAY), NOW(), NOW(), 5, 5);

-- 用户优惠券数据 (Demo Shop)
INSERT IGNORE INTO `user_coupons` (`id`, `tenant_id`, `user_id`, `coupon_id`, `status`, `used_at`, `order_id`, `received_at`, `expire_at`) VALUES
(1, 1, 1, 1, 1, DATE_SUB(NOW(), INTERVAL 5 DAY), 1, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_ADD(NOW(), INTERVAL 300 DAY)),
(2, 1, 1, 2, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY)),
(3, 1, 1, 3, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 60 DAY)),
(4, 1, 2, 1, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY)),
(5, 1, 2, 2, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY)),
(6, 1, 3, 1, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY)),
(7, 1, 3, 4, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY)),

-- Enterprise Corp 用户优惠券
(8, 3, 6, 6, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY)),
(9, 3, 6, 7, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 180 DAY)),
(10, 3, 7, 6, 0, NULL, 0, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY));
