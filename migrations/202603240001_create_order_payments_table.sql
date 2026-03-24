-- ============================================
-- Migration: Create order_payments table
-- Date: 2026-03-24
-- Description: Create order_payments table for
--              Stripe payment integration
-- ============================================

CREATE TABLE IF NOT EXISTS `order_payments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT 'Tenant ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT 'Order ID',
    `payment_no` VARCHAR(32) NOT NULL COMMENT 'Payment number',
    `payment_method` VARCHAR(32) NOT NULL COMMENT 'Payment method: stripe, alipay, wechat',
    `channel_intent_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel PaymentIntent ID (Stripe: pi_xxx)',
    `channel_payment_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel Charge ID (Stripe: ch_xxx)',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT 'Payment amount in cents',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=processing, 2=success, 3=failed, 4=cancelled, 5=refunded, 6=partially_refunded, 7=requires_action',
    `transaction_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'Transaction fee in cents',
    `fee_currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Fee currency',
    `paid_at` BIGINT DEFAULT NULL COMMENT 'Payment success time (UTC timestamp)',
    `failed_at` BIGINT DEFAULT NULL COMMENT 'Payment failure time (UTC timestamp)',
    `failed_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Failure reason',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    `deleted_at` BIGINT DEFAULT NULL COMMENT 'Deleted at (UTC timestamp, soft delete)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_payment_no` (`payment_no`),
    INDEX `idx_channel_payment_id` (`channel_payment_id`),
    INDEX `idx_channel_intent_id` (`channel_intent_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Order payments table';