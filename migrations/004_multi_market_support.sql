-- ShopJoy Multi-Market Support Migration
-- Creates markets, product_markets tables and adds compliance fields to products

-- Markets table: Market configuration entity
CREATE TABLE IF NOT EXISTS `markets` (
    `id`               bigint(20)      NOT NULL AUTO_INCREMENT COMMENT 'Market ID',
    `tenant_id`        bigint(20)      NOT NULL DEFAULT 0 COMMENT 'Tenant ID, 0 = global markets',
    `code`             varchar(10)     NOT NULL COMMENT 'Market code: US, UK, DE, FR, AU',
    `name`             varchar(64)     NOT NULL COMMENT 'Market name: United States',
    `currency`         varchar(10)     NOT NULL COMMENT 'Currency: USD, GBP, EUR, AUD',
    `default_language` varchar(10)     NOT NULL DEFAULT 'en' COMMENT 'Default language: en, de, fr',
    `flag`             varchar(32)     NULL COMMENT 'Flag emoji or image URL',
    `is_active`        tinyint(1)      NOT NULL DEFAULT 1 COMMENT '1 = active, 0 = inactive',
    `is_default`       tinyint(1)      NOT NULL DEFAULT 0 COMMENT '1 = primary market',
    `tax_rules`        json            NULL COMMENT 'Tax configuration: VAT, GST, IOSS',
    `created_at`       datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`       datetime        NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Markets table';

-- Product markets table: Market-specific product data
CREATE TABLE IF NOT EXISTS `product_markets` (
    `id`                    bigint(20)      NOT NULL AUTO_INCREMENT COMMENT 'Record ID',
    `tenant_id`             bigint(20)      NOT NULL COMMENT 'Tenant ID',
    `product_id`            bigint(20)      NOT NULL COMMENT 'Product ID',
    `variant_id`            bigint(20)      NULL COMMENT 'Variant ID, references skus.id, NULL for base product',
    `market_id`             bigint(20)      NOT NULL COMMENT 'Market ID',
    `is_enabled`            tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'Product visible in this market',
    `status_override`       tinyint(4)      NULL COMMENT 'Override product status per market',
    `price`                 decimal(19,4)   NOT NULL DEFAULT 0.0000 COMMENT 'Market-specific price',
    `compare_at_price`      decimal(19,4)   NULL COMMENT 'Compare at price for sales',
    `stock_alert_threshold` int(11)         NOT NULL DEFAULT 0 COMMENT 'Low stock alert threshold',
    `published_at`          datetime        NULL COMMENT 'Published timestamp in this market',
    `created_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_product_variant_market` (`product_id`, `variant_id`, `market_id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_market_id` (`market_id`),
    KEY `idx_is_enabled` (`is_enabled`),
    CONSTRAINT `fk_product_markets_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_product_markets_market_id` FOREIGN KEY (`market_id`) REFERENCES `markets` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_product_markets_variant_id` FOREIGN KEY (`variant_id`) REFERENCES `skus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Product-Market association table';

-- Add compliance fields to products table
-- Note: images column already exists in products table from 001_init_schema.sql
ALTER TABLE `products`
    ADD COLUMN `sku` varchar(64) NULL COMMENT 'SKU code' AFTER `name`,
    ADD COLUMN `brand` varchar(64) NULL COMMENT 'Brand name' AFTER `category_id`,
    ADD COLUMN `tags` json NULL COMMENT 'Product tags' AFTER `brand`,
    ADD COLUMN `is_matrix_product` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Has variants' AFTER `status`,
    ADD COLUMN `hs_code` varchar(20) NULL COMMENT 'Harmonized System Code' AFTER `is_matrix_product`,
    ADD COLUMN `coo` varchar(10) NULL COMMENT 'Country of Origin' AFTER `hs_code`,
    ADD COLUMN `weight` decimal(10,2) NULL COMMENT 'Weight' AFTER `coo`,
    ADD COLUMN `weight_unit` varchar(10) NULL DEFAULT 'g' COMMENT 'Weight unit: g, kg' AFTER `weight`,
    ADD COLUMN `length` decimal(10,2) NULL COMMENT 'Package length (cm)' AFTER `weight_unit`,
    ADD COLUMN `width` decimal(10,2) NULL COMMENT 'Package width (cm)' AFTER `length`,
    ADD COLUMN `height` decimal(10,2) NULL COMMENT 'Package height (cm)' AFTER `width`,
    ADD COLUMN `dangerous_goods` json NULL COMMENT 'Dangerous goods flags' AFTER `height`;

-- Add unique index on SKU
ALTER TABLE `products` ADD UNIQUE KEY `uk_sku` (`sku`);

-- Insert default markets for MVP
INSERT INTO `markets` (`tenant_id`, `code`, `name`, `currency`, `default_language`, `flag`, `is_active`, `is_default`, `tax_rules`) VALUES
(0, 'US', 'United States', 'USD', 'en', '🇺🇸', 1, 1, '{"IncludeTax": false}'),
(0, 'UK', 'United Kingdom', 'GBP', 'en', '🇬🇧', 1, 0, '{"VATRate": "20", "IncludeTax": true}'),
(0, 'DE', 'Germany', 'EUR', 'de', '🇩🇪', 1, 0, '{"VATRate": "19", "IOSSEnabled": true, "IncludeTax": true}'),
(0, 'FR', 'France', 'EUR', 'fr', '🇫🇷', 1, 0, '{"VATRate": "20", "IOSSEnabled": true, "IncludeTax": true}'),
(0, 'AU', 'Australia', 'AUD', 'en', '🇦🇺', 1, 0, '{"GSTRate": "10", "IncludeTax": true}')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);