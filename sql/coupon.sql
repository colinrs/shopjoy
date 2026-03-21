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