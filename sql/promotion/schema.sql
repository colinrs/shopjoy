-- ============================================
-- 促销活动表 (promotions) — 含优惠券 merged 2026-07-19
-- ============================================

CREATE TABLE IF NOT EXISTS `promotions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '促销ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' COMMENT '所有者类型',
    `name` VARCHAR(255) NOT NULL COMMENT '促销名称',
    `code` VARCHAR(100) NULL COMMENT '优惠券代码（仅 COUPON 使用）',
    `description` TEXT COMMENT '描述',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-折扣, 1-限时抢购, 2-捆绑销售, 3-买X送Y',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待生效, 1-生效中, 2-已暂停, 3-已结束',
    `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级',
    `market_id` BIGINT NULL COMMENT '市场ID',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `scope_type` VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' COMMENT '范围类型',
    `scope_ids` JSON DEFAULT NULL COMMENT '范围ID列表',
    `exclude_ids` JSON DEFAULT NULL COMMENT '排除ID列表',
    `usage_limit` INT NOT NULL DEFAULT 0 COMMENT '0 = unlimited',
    `total_count` INT NULL COMMENT '发放总数（COUPON 使用）',
    `used_count` INT NULL COMMENT '已使用数量（COUPON 使用）',
    `per_user_limit` INT NOT NULL DEFAULT 1 COMMENT 'per-user cap; 0 = unlimited',
    `tags` JSON DEFAULT NULL COMMENT 'tags',
    `start_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `end_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `code_unique` VARCHAR(100) GENERATED ALWAYS AS (IF(`kind` = 'COUPON', `code`, NULL)) VIRTUAL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_promotion_code` (`code_unique`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_type` (`type`),
    KEY `idx_priority` (`priority`),
    KEY `idx_start_end` (`start_at`, `end_at`),
    KEY `idx_active` (`status`, `currency`, `start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销活动表（含优惠券）';

-- usage_limit / per_user_limit / tags added 2026-07-18.
-- kind / code / market_id / total_count / used_count / code_unique / uk_promotion_code added 2026-07-19 during promotion+coupon merge.

-- ============================================
-- 促销规则表 (promotion_rules) — owner_kind/owner_id added 2026-07-19
-- ============================================

CREATE TABLE IF NOT EXISTS `promotion_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `promotion_id` BIGINT NULL COMMENT '促销ID（兼容旧促销规则；COUPON 规则此列为 NULL）',
    `owner_kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' COMMENT '所有者类型',
    `owner_id` BIGINT NOT NULL DEFAULT 0 COMMENT '所有者ID（指向 promotions.id）',
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
    KEY `idx_owner` (`owner_kind`, `owner_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销规则表（含优惠券规则）';

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
-- 用户优惠券表 (user_coupons) — coupon_id 仍然指向 archived coupons 表的 ID
-- ============================================

CREATE TABLE IF NOT EXISTS `user_coupons` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `coupon_id` BIGINT NOT NULL COMMENT '优惠券ID（指向 _archived_coupons_20260719.id）',
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
    KEY `idx_coupon_id_active` (`coupon_id`, `status`),
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
(1, 1, 0, 10000, 1, 10, 5000, 'CNY'),
(2, 1, 0, 30000, 1, 15, 10000, 'CNY'),
(3, 1, 0, 50000, 1, 20, 20000, 'CNY'),
(4, 2, 0, 0, 1, 30, 10000, 'CNY'),
(5, 3, 1, 2, 0, 0, 0, 'CNY'),
(6, 4, 0, 20000, 1, 25, 15000, 'CNY'),
(7, 5, 0, 0, 1, 40, 20000, 'CNY'),
(8, 6, 0, 50000, 1, 15, 30000, 'CNY'),
(9, 6, 0, 100000, 1, 20, 50000, 'CNY'),
(10, 7, 0, 30000, 1, 10, 20000, 'CNY');

-- 原 coupons 表已于 2026-07-19 RENAME 为 _archived_coupons_20260719（保留历史追溯）；新优惠券数据已合并入 promotions / promotion_rules 表。

