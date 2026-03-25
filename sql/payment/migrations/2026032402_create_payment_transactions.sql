-- ============================================
-- Migration: Create payment_transactions table
-- Date: 2026-03-24
-- Description: Create payment_transactions table for
--              payment transaction records
-- ============================================

CREATE TABLE IF NOT EXISTS `payment_transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT 'Tenant ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT 'Order ID',
    `payment_id` BIGINT UNSIGNED NOT NULL COMMENT 'Payment ID',
    `transaction_id` VARCHAR(64) NOT NULL COMMENT 'Transaction ID',
    `payment_method` VARCHAR(32) NOT NULL COMMENT 'Payment method: stripe, alipay, wechat',
    `channel_transaction_id` VARCHAR(64) DEFAULT '' COMMENT 'Channel transaction ID',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT 'Transaction amount in cents',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=succeeded, 2=failed',
    `transaction_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'Transaction fee in cents',
    `paid_at` BIGINT DEFAULT NULL COMMENT 'Payment success time (UTC timestamp)',
    `failed_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Failure reason',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    `deleted_at` BIGINT DEFAULT NULL COMMENT 'Deleted at (UTC timestamp, soft delete)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_transaction_id` (`transaction_id`),
    INDEX `idx_channel_transaction_id` (`channel_transaction_id`),
    INDEX `idx_payment_id` (`payment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Payment transactions table';
