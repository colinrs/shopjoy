-- ============================================
-- Migration: Create payment_refunds table
-- Date: 2026-03-24
-- Description: Create payment_refunds table for
--              payment refund management with idempotency
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
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT 'Refund amount in cents',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD' COMMENT 'Currency (ISO 4217)',
    `refund_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'Refund fee in cents',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status: 0=pending, 1=succeeded, 2=failed',
    `reason_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Refund reason type',
    `reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Refund reason details',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    `refunded_at` BIGINT DEFAULT NULL COMMENT 'Refund completed time (UTC timestamp)',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT 'Created by user ID',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_order` (`tenant_id`, `order_id`),
    UNIQUE INDEX `uk_refund_no` (`refund_no`),
    UNIQUE INDEX `uk_idempotency_key` (`idempotency_key`),
    INDEX `idx_channel_refund_id` (`channel_refund_id`),
    INDEX `idx_payment_id` (`payment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Payment refunds table';