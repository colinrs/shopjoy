-- migrations/20260325150000_create_shipping_tables.sql
-- Create shipping templates, zones, and mappings tables
-- Date: 2026-03-25

-- +migrate Up

-- Shipping templates table
CREATE TABLE IF NOT EXISTS `shipping_templates` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `name` VARCHAR(100) NOT NULL COMMENT 'Template name',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT 'Is default template (0=no, 1=yes)',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Is active (0=inactive, 1=active)',
    `deleted_at` BIGINT NULL DEFAULT NULL COMMENT 'Deleted at (UTC timestamp)',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_is_default` (`is_default`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping templates table';

-- Shipping zones table
CREATE TABLE IF NOT EXISTS `shipping_zones` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `template_id` BIGINT NOT NULL COMMENT 'Template ID',
    `name` VARCHAR(100) NOT NULL COMMENT 'Zone name',
    `regions` JSON NOT NULL COMMENT 'Region codes (city codes array)',
    `fee_type` VARCHAR(20) NOT NULL COMMENT 'Fee type: fixed, by_count, by_weight, free',
    `first_unit` INT NOT NULL DEFAULT 1 COMMENT 'First unit (count or grams)',
    `first_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'First fee in cents',
    `additional_unit` INT NOT NULL DEFAULT 1 COMMENT 'Additional unit',
    `additional_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'Additional fee in cents',
    `free_threshold_amount` BIGINT NOT NULL DEFAULT 0 COMMENT 'Free shipping threshold amount in cents, 0=disabled',
    `free_threshold_count` INT NOT NULL DEFAULT 0 COMMENT 'Free shipping threshold count, 0=disabled',
    `sort` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping zones table';

-- Shipping template mappings table
CREATE TABLE IF NOT EXISTS `shipping_template_mappings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `template_id` BIGINT NOT NULL COMMENT 'Template ID',
    `target_type` VARCHAR(20) NOT NULL COMMENT 'Target type: product, category',
    `target_id` BIGINT NOT NULL COMMENT 'Target ID (product_id or category_id)',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated at (UTC timestamp)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_template_id` (`template_id`),
    INDEX `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping template mappings table';

-- +migrate Down
DROP TABLE IF EXISTS `shipping_template_mappings`;
DROP TABLE IF EXISTS `shipping_zones`;
DROP TABLE IF EXISTS `shipping_templates`;