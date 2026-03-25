-- migrations/20260325160000_add_zone_regions_index.sql
-- Add zone regions junction table for efficient region lookup
-- This replaces JSON_CONTAINS queries with indexed lookups
-- Date: 2026-03-25

-- +migrate Up

-- Create junction table for zone regions (normalized for efficient queries)
CREATE TABLE IF NOT EXISTS `shipping_zone_regions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `zone_id` BIGINT NOT NULL COMMENT 'Zone ID',
    `city_code` VARCHAR(20) NOT NULL COMMENT 'City code',
    `created_at` BIGINT NOT NULL COMMENT 'Created at (UTC timestamp)',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_zone_city` (`zone_id`, `city_code`),
    INDEX `idx_city_code` (`city_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Shipping zone regions junction table';

-- +migrate Down
DROP TABLE IF EXISTS `shipping_zone_regions`;