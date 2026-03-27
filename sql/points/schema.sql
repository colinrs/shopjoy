-- ============================================
-- 积分获取规则表 (earn_rules)
-- 定义积分获取方式、计算规则和条件
-- ============================================

CREATE TABLE IF NOT EXISTS `earn_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '规则名称',
    `description` TEXT COMMENT '描述',
    `scenario` VARCHAR(50) NOT NULL COMMENT '场景: ORDER_PAYMENT-订单支付, SIGN_IN-签到, PRODUCT_REVIEW-商品评价, FIRST_ORDER-首单',
    `calculation_type` VARCHAR(20) NOT NULL COMMENT '计算类型: FIXED-固定积分, RATIO-比例, TIERED-阶梯',
    `fixed_points` BIGINT NOT NULL DEFAULT 0 COMMENT '固定积分数(当calculation_type=FIXED时使用)',
    `ratio` VARCHAR(20) DEFAULT '' COMMENT '积分比例(当calculation_type=RATIO时使用,如"1:100"表示1元=100积分)',
    `tiers` TEXT COMMENT '阶梯配置JSON(当calculation_type=TIERED时使用)',
    `condition_type` VARCHAR(50) NOT NULL DEFAULT 'NONE' COMMENT '条件类型: NONE-无条件, NEW_USER-新用户, FIRST_ORDER-首单, SPECIFIC_PRODUCTS-指定商品, MIN_AMOUNT-最低金额',
    `condition_value` TEXT COMMENT '条件值JSON',
    `expiration_months` INT NOT NULL DEFAULT 12 COMMENT '积分过期月数',
    `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '状态: draft-草稿, active-生效, inactive-停用',
    `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级(数字越大优先级越高)',
    `start_at` BIGINT DEFAULT NULL COMMENT '开始时间(Unix时间戳(秒))',
    `end_at` BIGINT DEFAULT NULL COMMENT '结束时间(Unix时间戳(秒))',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(Unix时间戳(秒))',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(Unix时间戳(秒))',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间(Unix时间戳(秒))',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_scenario` (`scenario`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_start_end` (`start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分获取规则表';

-- ============================================
-- 积分兑换规则表 (redeem_rules)
-- 定义积分兑换优惠券的规则
-- ============================================

CREATE TABLE IF NOT EXISTS `redeem_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '规则名称',
    `description` TEXT COMMENT '描述',
    `coupon_id` BIGINT NOT NULL COMMENT '关联优惠券ID',
    `points_required` BIGINT NOT NULL COMMENT '所需积分',
    `total_stock` BIGINT NOT NULL DEFAULT 0 COMMENT '总库存',
    `used_stock` BIGINT NOT NULL DEFAULT 0 COMMENT '已使用库存',
    `per_user_limit` INT NOT NULL DEFAULT 1 COMMENT '每用户限兑次数',
    `status` VARCHAR(20) NOT NULL DEFAULT 'inactive' COMMENT '状态: inactive-未激活, active-激活',
    `start_at` BIGINT DEFAULT NULL COMMENT '开始时间(Unix时间戳(秒))',
    `end_at` BIGINT DEFAULT NULL COMMENT '结束时间(Unix时间戳(秒))',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(Unix时间戳(秒))',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(Unix时间戳(秒))',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间(Unix时间戳(秒))',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_coupon_id` (`coupon_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_start_end` (`start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分兑换规则表';

-- ============================================
-- 积分账户表 (points_accounts)
-- 用户积分账户,记录余额和统计信息
-- ============================================

CREATE TABLE IF NOT EXISTS `points_accounts` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '账户ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `balance` BIGINT NOT NULL DEFAULT 0 COMMENT '可用积分余额',
    `frozen_balance` BIGINT NOT NULL DEFAULT 0 COMMENT '冻结积分',
    `total_earned` BIGINT NOT NULL DEFAULT 0 COMMENT '累计获得积分',
    `total_redeemed` BIGINT NOT NULL DEFAULT 0 COMMENT '累计兑换积分',
    `total_expired` BIGINT NOT NULL DEFAULT 0 COMMENT '累计过期积分',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(Unix时间戳(秒))',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间(Unix时间戳(秒))',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间(Unix时间戳(秒))',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user` (`tenant_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分账户表';

-- ============================================
-- 积分交易记录表 (points_transactions)
-- 记录所有积分变动明细
-- ============================================

CREATE TABLE IF NOT EXISTS `points_transactions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '交易ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `account_id` BIGINT NOT NULL COMMENT '账户ID',
    `points` BIGINT NOT NULL COMMENT '积分变动(正数为获得,负数为扣除)',
    `balance_after` BIGINT NOT NULL COMMENT '变动后余额',
    `type` VARCHAR(20) NOT NULL COMMENT '交易类型: EARN-获得, REDEEM-兑换, ADJUST-调整, EXPIRE-过期, FREEZE-冻结, UNFREEZE-解冻',
    `reference_type` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '关联类型: ORDER, EARN_RULE, REDEEM_RULE, MANUAL_ADJUST等',
    `reference_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '关联ID',
    `description` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '描述',
    `expires_at` BIGINT DEFAULT NULL COMMENT '积分过期时间(Unix时间戳(秒))',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间(Unix时间戳(秒))',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间(Unix时间戳(秒))',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_account_id` (`account_id`),
    KEY `idx_type` (`type`),
    KEY `idx_reference` (`reference_type`, `reference_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_expires_at` (`expires_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分交易记录表';

-- ============================================
-- 积分兑换记录表 (points_redemptions)
-- 记录用户积分兑换优惠券的详情
-- ============================================

CREATE TABLE IF NOT EXISTS `points_redemptions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '兑换记录ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `redeem_rule_id` BIGINT NOT NULL COMMENT '兑换规则ID',
    `coupon_id` BIGINT NOT NULL COMMENT '优惠券ID',
    `user_coupon_id` BIGINT DEFAULT NULL COMMENT '用户优惠券ID(兑换成功后生成)',
    `points_used` BIGINT NOT NULL COMMENT '消耗积分',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待处理, completed-已完成, cancelled-已取消',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `completed_at` BIGINT DEFAULT NULL COMMENT '完成时间(Unix时间戳(秒))',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间(Unix时间戳(秒))',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_redeem_rule_id` (`redeem_rule_id`),
    KEY `idx_coupon_id` (`coupon_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分兑换记录表';

-- ============================================
-- 测试数据
-- ============================================

-- 积分获取规则数据 (Demo Shop)
INSERT INTO `earn_rules` (`id`, `tenant_id`, `name`, `description`, `scenario`, `calculation_type`, `fixed_points`, `ratio`, `tiers`, `condition_type`, `condition_value`, `expiration_months`, `status`, `priority`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '订单支付返积分', '每消费1元获得1积分', 'ORDER_PAYMENT', 'RATIO', 0, '1:1', NULL, 'NONE', NULL, 12, 'active', 10, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '签到送积分', '每日签到获得5积分', 'SIGN_IN', 'FIXED', 5, '', NULL, 'NONE', NULL, 12, 'active', 5, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '商品评价送积分', '评价商品获得10积分', 'PRODUCT_REVIEW', 'FIXED', 10, '', NULL, 'NONE', NULL, 6, 'active', 5, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '首单双倍积分', '首单支付获得双倍积分', 'FIRST_ORDER', 'RATIO', 0, '2:1', NULL, 'FIRST_ORDER', NULL, 12, 'active', 15, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, '阶梯返积分', '消费满不同金额获得不同比例积分', 'ORDER_PAYMENT', 'TIERED', 0, '', '[{"threshold":10000,"ratio":"1:1"},{"threshold":50000,"ratio":"1.5:1"},{"threshold":100000,"ratio":"2:1"}]', 'MIN_AMOUNT', '{"min_amount":10000}', 12, 'active', 8, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- Enterprise Corp 积分规则
(6, 3, '企业客户积分', '企业客户每消费1元获得2积分', 'ORDER_PAYMENT', 'RATIO', 0, '2:1', NULL, 'NONE', NULL, 24, 'active', 10, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(7, 3, '大额订单奖励', '订单满500元额外获得100积分', 'ORDER_PAYMENT', 'FIXED', 100, '', NULL, 'MIN_AMOUNT', '{"min_amount":50000}', 12, 'active', 5, UNIX_TIMESTAMP(NOW()), NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 积分兑换规则数据 (Demo Shop)
INSERT INTO `redeem_rules` (`id`, `tenant_id`, `name`, `description`, `coupon_id`, `points_required`, `total_stock`, `used_stock`, `per_user_limit`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '满50减10优惠券', '500积分兑换满50减10优惠券', 3, 500, 1000, 150, 3, 'active', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '满100减20优惠券', '800积分兑换满100减20优惠券', 2, 800, 500, 80, 2, 'active', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '新用户专享券', '300积分兑换新用户专享50元优惠券', 1, 300, 2000, 300, 1, 'active', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '免邮券', '200积分兑换免邮券', 4, 200, 500, 50, 1, 'active', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 60 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- Enterprise Corp 兑换规则
(5, 3, '企业专属优惠券', '1000积分兑换企业专属50元优惠券', 6, 1000, 300, 45, 5, 'active', UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 积分账户数据 (Demo Shop 用户)
INSERT INTO `points_accounts` (`id`, `tenant_id`, `user_id`, `balance`, `frozen_balance`, `total_earned`, `total_redeemed`, `total_expired`, `created_at`, `updated_at`) VALUES
(1, 1, 1, 2500, 0, 5000, 2000, 500, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 1, 2, 1800, 100, 3500, 1600, 200, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 1, 3, 3200, 0, 6000, 2500, 300, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- Enterprise Corp 用户
(4, 3, 6, 8500, 0, 12000, 3000, 500, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, 3, 7, 4200, 200, 8000, 3600, 400, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 积分交易记录数据 (Demo Shop)
INSERT INTO `points_transactions` (`id`, `tenant_id`, `user_id`, `account_id`, `points`, `balance_after`, `type`, `reference_type`, `reference_id`, `description`, `expires_at`, `created_at`) VALUES
(1, 1, 1, 1, 1500, 1500, 'EARN', 'ORDER', 'ORD202503100001', '订单支付获得积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 12 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 15 DAY))),
(2, 1, 1, 1, 500, 2000, 'EARN', 'EARN_RULE', '2', '签到获得积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 12 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY))),
(3, 1, 1, 1, -500, 1500, 'REDEEM', 'REDEEM_RULE', '1', '兑换优惠券', NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 8 DAY))),
(4, 1, 1, 1, 1000, 2500, 'EARN', 'ORDER', 'ORD202503160001', '订单支付获得积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 12 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY))),
(5, 1, 2, 2, 2000, 2000, 'EARN', 'ORDER', 'ORD202503120001', '订单支付获得积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 12 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 12 DAY))),
(6, 1, 2, 2, -200, 1800, 'REDEEM', 'REDEEM_RULE', '4', '兑换免邮券', NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 6 DAY))),
(7, 1, 3, 3, 3000, 3000, 'EARN', 'ORDER', 'ORD202503140001', '首单双倍积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 12 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY))),
(8, 1, 3, 3, 200, 3200, 'EARN', 'EARN_RULE', '3', '商品评价获得积分', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 6 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY))),

-- Enterprise Corp 交易记录
(9, 3, 6, 4, 5000, 5000, 'EARN', 'ORDER', 'ORD202503150001', '订单支付获得积分(企业双倍)', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 24 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 8 DAY))),
(10, 3, 6, 4, -1000, 4000, 'REDEEM', 'REDEEM_RULE', '5', '兑换企业专属优惠券', NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY))),
(11, 3, 7, 5, 4500, 4500, 'EARN', 'ORDER', 'ORD202503130001', '订单支付获得积分(企业双倍)', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 24 MONTH)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 7 DAY))),
(12, 3, 7, 5, -300, 4200, 'REDEEM', 'REDEEM_RULE', '5', '兑换企业专属优惠券', NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)));

-- 积分兑换记录数据 (Demo Shop)
INSERT INTO `points_redemptions` (`id`, `tenant_id`, `user_id`, `redeem_rule_id`, `coupon_id`, `user_coupon_id`, `points_used`, `status`, `created_at`, `completed_at`) VALUES
(1, 1, 1, 1, 3, 1, 500, 'completed', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 8 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 8 DAY))),
(2, 1, 2, 4, 4, 7, 200, 'completed', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 6 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 6 DAY))),
(3, 1, 3, 2, 2, NULL, 800, 'pending', UNIX_TIMESTAMP(NOW()), NULL),

-- Enterprise Corp 兑换记录
(4, 3, 6, 5, 6, 8, 1000, 'completed', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY))),
(5, 3, 7, 5, 6, 9, 1000, 'completed', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)));