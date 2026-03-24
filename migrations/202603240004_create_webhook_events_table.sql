-- ============================================
-- Migration: Create webhook_events table
-- Date: 2026-03-24
-- Description: Create webhook_events table for
--              webhook event deduplication
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
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    `processed_at` BIGINT DEFAULT NULL COMMENT 'Processed time (UTC timestamp)',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_event_id` (`event_id`),
    INDEX `idx_tenant_event` (`tenant_id`, `event_type`),
    INDEX `idx_resource` (`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook events table for deduplication';
