-- ============================================
-- 支付表 (payments)
-- ============================================

CREATE TABLE IF NOT EXISTS `payments` (
    `id` VARCHAR(64) NOT NULL COMMENT '支付ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '支付金额(分)',
    `currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待支付, 1-处理中, 2-成功, 3-失败, 4-取消, 5-已退款',
    `method` TINYINT NOT NULL DEFAULT 0 COMMENT '支付方式: 0-支付宝, 1-微信, 2-信用卡, 3-银行转账, 4-货到付款',
    `transaction_id` VARCHAR(255) DEFAULT '' COMMENT '第三方交易号',
    `paid_at` BIGINT DEFAULT NULL COMMENT '支付时间',
    `expire_at` BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间',
    `notify_url` VARCHAR(500) DEFAULT '' COMMENT '回调URL',
    `return_url` VARCHAR(500) DEFAULT '' COMMENT '返回URL',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_transaction_id` (`transaction_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付表';

-- ============================================
-- 支付退款表 (payment_refunds)
-- ============================================

CREATE TABLE IF NOT EXISTS `payment_refunds` (
    `id` VARCHAR(64) NOT NULL COMMENT '退款ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `payment_id` VARCHAR(64) NOT NULL COMMENT '支付ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '退款金额(分)',
    `currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `reason` VARCHAR(500) DEFAULT '' COMMENT '退款原因',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待处理, 1-处理中, 2-完成, 3-失败',
    `transaction_id` VARCHAR(255) DEFAULT '' COMMENT '第三方交易号',
    `refunded_at` BIGINT DEFAULT NULL COMMENT '退款时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_payment_id` (`payment_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付退款表';

-- ============================================
-- 测试数据
-- ============================================

-- 支付数据
INSERT INTO `payments` (`id`, `tenant_id`, `order_id`, `user_id`, `amount`, `currency`, `status`, `method`, `transaction_id`, `paid_at`, `expire_at`, `notify_url`, `return_url`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- Demo Shop 支付
('PAY202503010001', 1, 'ORD202503010001', 1, 254800, 'CNY', 2, 0, 'ALI202503010001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 29 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 29 DAY)), 'https://api.demoshop.com/payment/notify', 'https://www.demoshop.com/order/success', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 29 DAY)), 1, 1),
('PAY202503100001', 1, 'ORD202503100001', 1, 169800, 'CNY', 2, 1, 'WX202503100001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), 'https://api.demoshop.com/payment/notify', 'https://www.demoshop.com/order/success', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), 1, 1),
('PAY202503150001', 1, 'ORD202503150001', 2, 9900, 'CNY', 2, 0, 'ALI202503150001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY)), 'https://api.demoshop.com/payment/notify', 'https://www.demoshop.com/order/success', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY)), 2, 2),
('PAY202503200001', 1, 'ORD202503200001', 3, 116910, 'CNY', 0, 0, '', NULL, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 DAY)), 'https://api.demoshop.com/payment/notify', 'https://www.demoshop.com/order/success', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 3, 3),

-- Enterprise Corp 支付
('PAY202503120001', 3, 'ORD202503120001', 6, 81400, 'CNY', 2, 2, 'CC202503120001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 DAY)), 'https://api.enterprisecorp.com/payment/notify', 'https://shop.enterprisecorp.com/order/success', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 DAY)), 6, 6);

-- 支付退款数据
INSERT INTO `payment_refunds` (`id`, `tenant_id`, `payment_id`, `order_id`, `user_id`, `amount`, `currency`, `reason`, `status`, `transaction_id`, `refunded_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- ORD202503010001 部分退款 (手机壳)
('REF202503050001', 1, 'PAY202503010001', 'ORD202503010001', 1, 29700, 'CNY', '商品质量问题', 2, 'ALIREF202503050001', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 26 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY)), 1, 1);