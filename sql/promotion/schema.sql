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
    `start_at` BIGINT NOT NULL DEFAULT 0 COMMENT '开始时间',
    `end_at` BIGINT NOT NULL DEFAULT 0 COMMENT '结束时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_type` (`type`),
    KEY `idx_start_end` (`start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销活动表';

-- ============================================
-- 促销规则表 (promotion_rules)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotion_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `promotion_id` BIGINT NOT NULL COMMENT '促销ID',
    `condition_type` TINYINT NOT NULL DEFAULT 0 COMMENT '条件类型: 0-最低金额, 1-最低数量',
    `condition_value` BIGINT NOT NULL DEFAULT 0 COMMENT '条件值',
    `action_type` TINYINT NOT NULL DEFAULT 0 COMMENT '动作类型: 0-固定金额, 1-百分比',
    `action_value` BIGINT NOT NULL DEFAULT 0 COMMENT '动作值',
    `max_discount_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '最大优惠金额(分)',
    `max_discount_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    PRIMARY KEY (`id`),
    KEY `idx_promotion_id` (`promotion_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销规则表';

-- ============================================
-- 测试数据
-- ============================================

-- 促销活动数据 (Demo Shop)
INSERT INTO `promotions` (`id`, `tenant_id`, `name`, `description`, `type`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '春季大促', '春季全场促销活动', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '限时抢购-鞋类专场', '运动鞋限时抢购', 1, 1, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 8 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '买二送一', '购买任意两件商品送一件同等价值商品', 3, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 60 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '会员日特惠', '每月18号会员专享折扣', 0, 2, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 40 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, '清仓大甩卖', '库存清仓特价', 1, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 15 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- Enterprise Corp 促销活动
(6, 3, '企业采购季', '企业客户采购优惠季', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(7, 3, '跨境特惠', '跨境电商专属优惠', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 促销规则数据
INSERT INTO `promotion_rules` (`id`, `promotion_id`, `condition_type`, `condition_value`, `action_type`, `action_value`, `max_discount_amount`, `max_discount_currency`) VALUES
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
(10, 7, 0, 30000, 1, 10, 20000, 'CNY'); -- 满300减10%, 最多减200-- ============================================
-- 优惠券表 (coupons)
-- ============================================

CREATE TABLE IF NOT EXISTS `coupons` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '名称',
    `code` VARCHAR(100) NOT NULL COMMENT '优惠券代码',
    `description` TEXT COMMENT '描述',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-固定金额, 1-百分比, 2-免邮',
    `value` BIGINT NOT NULL DEFAULT 0 COMMENT '优惠值(分或百分比)',
    `min_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '最低消费(分)',
    `max_discount` BIGINT NOT NULL DEFAULT 0 COMMENT '最大优惠(分)',
    `total_count` INT NOT NULL DEFAULT 0 COMMENT '发放总数',
    `used_count` INT NOT NULL DEFAULT 0 COMMENT '已使用数量',
    `per_user_limit` INT NOT NULL DEFAULT 1 COMMENT '每用户限领',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未激活, 1-激活, 2-过期, 3-用完',
    `start_at` BIGINT NOT NULL DEFAULT 0 COMMENT '开始时间',
    `end_at` BIGINT NOT NULL DEFAULT 0 COMMENT '结束时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
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
    `used_at` BIGINT DEFAULT NULL COMMENT '使用时间',
    `order_id` VARCHAR(64) DEFAULT '' COMMENT '订单ID',
    `received_at` BIGINT NOT NULL DEFAULT 0 COMMENT '领取时间',
    `expire_at` BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间',
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

-- 优惠券数据 (Demo Shop)
INSERT INTO `coupons` (`id`, `tenant_id`, `name`, `code`, `description`, `type`, `value`, `min_amount`, `max_discount`, `total_count`, `used_count`, `per_user_limit`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '新用户专享', 'NEWUSER50', '新用户首单立减50元', 0, 5000, 10000, 0, 1000, 150, 1, 1, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '全场9折', 'SAVE10', '全场商品9折优惠', 1, 10, 5000, 10000, 500, 80, 3, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '满200减30', 'SAVE30', '满200元减30元', 0, 3000, 20000, 0, 200, 45, 2, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 60 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '免邮券', 'FREESHIP', '全场免运费', 2, 0, 5000, 1500, 300, 20, 1, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, 'VIP专享8折', 'VIP20', 'VIP会员专享8折', 1, 20, 10000, 20000, 100, 10, 1, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- Enterprise Corp 优惠券
(6, 3, '企业专属优惠', 'ENT50', '企业客户专享50元优惠', 0, 5000, 20000, 0, 500, 30, 5, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(7, 3, '跨境免邮', 'GLOBALSHIP', '跨境电商免邮', 2, 0, 30000, 5000, 200, 15, 3, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 用户优惠券数据 (Demo Shop)
INSERT INTO `user_coupons` (`id`, `tenant_id`, `user_id`, `coupon_id`, `status`, `used_at`, `order_id`, `received_at`, `expire_at`) VALUES
(1, 1, 1, 1, 1, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY)), 'ORD202503160001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 300 DAY))),
(2, 1, 1, 2, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY))),
(3, 1, 1, 3, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 60 DAY))),
(4, 1, 2, 1, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY))),
(5, 1, 2, 2, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY))),
(6, 1, 3, 1, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY))),
(7, 1, 3, 4, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 30 DAY))),

-- Enterprise Corp 用户优惠券
(8, 3, 6, 6, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY))),
(9, 3, 6, 7, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY))),
(10, 3, 7, 6, 0, NULL, '', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 365 DAY)));