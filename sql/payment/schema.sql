-- ============================================
-- 支付表 (payments) - 已废弃，请使用 order_payments
-- ============================================

CREATE TABLE IF NOT EXISTS `payments` (
    `id` VARCHAR(64) NOT NULL COMMENT '支付ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '支付金额',
    `currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待支付, 1-处理中, 2-成功, 3-失败, 4-取消, 5-已退款',
    `method` TINYINT NOT NULL DEFAULT 0 COMMENT '支付方式: 0-支付宝, 1-微信, 2-信用卡, 3-银行转账, 4-货到付款',
    `transaction_id` VARCHAR(255) DEFAULT '' COMMENT '第三方交易号',
    `paid_at` TIMESTAMP NULL COMMENT '支付时间',
    `expire_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '过期时间',
    `notify_url` VARCHAR(500) DEFAULT '' COMMENT '回调URL',
    `return_url` VARCHAR(500) DEFAULT '' COMMENT '返回URL',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_transaction_id` (`transaction_id`),
    KEY `idx_status` (`status`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付表';

-- ============================================
-- 订单支付表 (order_payments)
-- ============================================

CREATE TABLE IF NOT EXISTS `order_payments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT 'Tenant ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT 'Order ID',
    `payment_no` VARCHAR(32) NOT NULL COMMENT 'Payment number',
    `payment_method` VARCHAR(32) NOT NULL COMMENT 'Payment method: stripe, alipay, wechat',
    `channel_intent_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel PaymentIntent ID (Stripe: pi_xxx)',
    `channel_payment_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel Charge ID (Stripe: ch_xxx)',
    `amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Payment amount',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=processing, 2=success, 3=failed, 4=cancelled, 5=refunded, 6=partially_refunded, 7=requires_action',
    `transaction_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Transaction fee',
    `fee_currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Fee currency',
    `paid_at` TIMESTAMP NULL COMMENT 'Payment success time',
    `failed_at` TIMESTAMP NULL COMMENT 'Payment failure time',
    `failed_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Failure reason',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_payment_no` (`payment_no`),
    INDEX `idx_channel_payment_id` (`channel_payment_id`),
    INDEX `idx_channel_intent_id` (`channel_intent_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Order payments table';

-- ============================================
-- 支付交易记录表 (payment_transactions)
-- ============================================

CREATE TABLE IF NOT EXISTS `payment_transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT 'Tenant ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT 'Order ID',
    `payment_id` BIGINT UNSIGNED NOT NULL COMMENT 'Payment ID',
    `transaction_id` VARCHAR(64) NOT NULL COMMENT 'Transaction ID',
    `payment_method` VARCHAR(32) NOT NULL COMMENT 'Payment method: stripe, alipay, wechat',
    `channel_transaction_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel transaction ID',
    `amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Transaction amount',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=succeeded, 2=failed',
    `transaction_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Transaction fee',
    `paid_at` TIMESTAMP NULL COMMENT 'Payment success time',
    `failed_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Failure reason',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_transaction_id` (`transaction_id`),
    INDEX `idx_channel_transaction_id` (`channel_transaction_id`),
    INDEX `idx_payment_id` (`payment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Payment transactions table';

-- ============================================
-- 支付退款详情表 (payment_refunds_detail)
-- ============================================

CREATE TABLE IF NOT EXISTS `payment_refunds` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT 'Tenant ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT 'Order ID',
    `payment_id` BIGINT UNSIGNED NOT NULL COMMENT 'Payment ID',
    `fulfillment_refund_id` BIGINT UNSIGNED DEFAULT NULL COMMENT 'Fulfillment refund ID (link to fulfillment.refunds)',
    `refund_no` VARCHAR(32) NOT NULL COMMENT 'Refund number',
    `idempotency_key` VARCHAR(64) NOT NULL COMMENT 'Idempotency key for deduplication',
    `channel_refund_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel refund ID (Stripe: re_xxx)',
    `amount` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Refund amount',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `refund_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT 'Refund fee',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=succeeded, 2=failed',
    `reason_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Refund reason type',
    `reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Refund reason details',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `refunded_at` TIMESTAMP NULL COMMENT 'Refund completed time',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT 'Created by user ID',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_refund_no` (`refund_no`),
    UNIQUE INDEX `uk_idempotency_key` (`idempotency_key`),
    INDEX `idx_channel_refund_id` (`channel_refund_id`),
    INDEX `idx_payment_id` (`payment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Payment refunds table';

-- ============================================
-- Webhook 事件表 (webhook_events)
-- ============================================

CREATE TABLE IF NOT EXISTS `webhook_events` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Tenant ID (0 for platform-level events)',
    `event_id` VARCHAR(64) NOT NULL COMMENT 'Event ID (Stripe: evt_xxx)',
    `event_type` VARCHAR(64) NOT NULL COMMENT 'Event type',
    `resource_id` VARCHAR(64) DEFAULT '' COMMENT 'Resource ID (PaymentIntent/Charge ID)',
    `processed` TINYINT NOT NULL DEFAULT 0 COMMENT 'Processed status: 0=pending, 1=processed, 2=failed',
    `raw_payload` TEXT COMMENT 'Raw event JSON payload',
    `error_message` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Processing error message',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    `processed_at` TIMESTAMP NULL COMMENT 'Processed time',
    `deleted_at` TIMESTAMP NULL COMMENT 'Deleted at',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_event_id` (`event_id`),
    INDEX `idx_tenant_event` (`tenant_id`, `event_type`),
    INDEX `idx_resource` (`resource_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook events table for deduplication';

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
