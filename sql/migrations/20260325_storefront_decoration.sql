-- ============================================
-- Storefront Decoration Module Migration
-- Date: 2026-03-25
-- Description: Add theme management, page decoration, version history, and SEO configuration
-- ============================================

-- ============================================
-- 1. Extend themes table
-- ============================================

-- Add new columns to themes table
ALTER TABLE `themes`
    ADD COLUMN `preview_image` VARCHAR(500) DEFAULT '' COMMENT 'Preview image URL' AFTER `thumbnail`,
    ADD COLUMN `config_schema` TEXT COMMENT 'JSON schema for configurable fields' AFTER `config`,
    ADD COLUMN `default_config` TEXT COMMENT 'JSON default configuration' AFTER `config_schema`,
    ADD COLUMN `is_preset` TINYINT NOT NULL DEFAULT 1 COMMENT '1=preset theme, 0=custom theme' AFTER `is_custom`;

-- Update existing themes to be preset themes
UPDATE `themes` SET `is_preset` = 1 WHERE `is_preset` IS NULL OR `is_preset` = 0;

-- ============================================
-- 2. Extend shops table
-- ============================================

ALTER TABLE `shops`
    ADD COLUMN `current_theme_id` BIGINT DEFAULT NULL COMMENT 'Current active theme ID' AFTER `seo_keywords`,
    ADD COLUMN `theme_config` TEXT COMMENT 'JSON theme customization config' AFTER `current_theme_id`;

-- ============================================
-- 3. Extend pages table
-- ============================================

ALTER TABLE `pages`
    ADD COLUMN `is_published` TINYINT NOT NULL DEFAULT 0 COMMENT 'Whether page is published' AFTER `sort`,
    ADD COLUMN `published_at` BIGINT DEFAULT NULL COMMENT 'Published timestamp (Unix)' AFTER `is_published`,
    ADD COLUMN `version` INT NOT NULL DEFAULT 1 COMMENT 'Current version number' AFTER `published_at`,
    ADD COLUMN `deleted_at` BIGINT DEFAULT NULL COMMENT 'Soft delete timestamp' AFTER `updated_by`;

-- Add index for soft delete queries
ALTER TABLE `pages` ADD INDEX `idx_deleted_at` (`deleted_at`);

-- ============================================
-- 4. Create decorations table
-- ============================================

CREATE TABLE IF NOT EXISTS `decorations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'Decoration block ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_id` BIGINT NOT NULL COMMENT 'Page ID',
    `block_type` VARCHAR(50) NOT NULL COMMENT 'Block type: banner, product_grid, rich_text, image_carousel, featured_products, categories, divider',
    `block_config` TEXT NOT NULL COMMENT 'JSON block configuration',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order within page',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Whether block is active',
    `created_at` BIGINT NOT NULL COMMENT 'Created timestamp (Unix UTC)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated timestamp (Unix UTC)',
    PRIMARY KEY (`id`),
    INDEX `idx_page_sort` (`page_id`, `sort_order`),
    INDEX `idx_tenant_page` (`tenant_id`, `page_id`),
    INDEX `idx_block_type` (`block_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Page decoration blocks';

-- ============================================
-- 5. Create page_versions table
-- ============================================

CREATE TABLE IF NOT EXISTS `page_versions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'Version ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_id` BIGINT NOT NULL COMMENT 'Page ID',
    `version` INT NOT NULL COMMENT 'Version number',
    `blocks` TEXT NOT NULL COMMENT 'JSON snapshot of decoration blocks',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT 'User who created this version',
    `created_at` BIGINT NOT NULL COMMENT 'Created timestamp (Unix UTC)',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_tenant_page_ver` (`tenant_id`, `page_id`, `version`),
    INDEX `idx_page_version` (`page_id`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Page version history';

-- ============================================
-- 6. Create seo_configs table
-- ============================================

CREATE TABLE IF NOT EXISTS `seo_configs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'SEO config ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_type` VARCHAR(30) NOT NULL COMMENT 'Page type: global, home, category, product, custom',
    `page_id` BIGINT DEFAULT NULL COMMENT 'Page ID for custom pages (NULL for global/page type defaults)',
    `title` VARCHAR(200) NOT NULL DEFAULT '' COMMENT 'SEO title',
    `description` TEXT NOT NULL COMMENT 'SEO description',
    `keywords` VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'SEO keywords',
    `created_at` BIGINT NOT NULL COMMENT 'Created timestamp (Unix UTC)',
    `updated_at` BIGINT NOT NULL COMMENT 'Updated timestamp (Unix UTC)',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_tenant_page_type` (`tenant_id`, `page_type`, `page_id`),
    INDEX `idx_page_type` (`page_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SEO configurations';

-- ============================================
-- 7. Insert preset themes
-- ============================================

-- Remove existing preset themes data (will be replaced with new preset themes)
DELETE FROM `themes` WHERE `is_preset` = 1;

-- Insert new preset themes (tenant_id = 0 means it's a system preset available to all tenants)
INSERT INTO `themes` (`id`, `tenant_id`, `name`, `code`, `description`, `thumbnail`, `preview_image`, `config`, `config_schema`, `default_config`, `is_active`, `is_custom`, `is_preset`, `created_at`, `updated_at`) VALUES
(1001, 0, 'Classic', 'classic', 'A timeless design with clean lines and professional appearance. Perfect for businesses looking for a traditional e-commerce feel.',
 'https://cdn.shopjoy.com/themes/classic-thumb.png', 'https://cdn.shopjoy.com/themes/classic-preview.png',
 '{"colors":{"primary":"#3B82F6","secondary":"#1E40AF"},"fonts":{"heading":"Inter","body":"Inter"},"layout":"standard"}',
 '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#3B82F6"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#1E40AF"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"inter","label":"Inter"},{"value":"roboto","label":"Roboto"},{"value":"opensans","label":"Open Sans"}],"default":"inter"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"square","label":"Square"}],"default":"rounded"}]',
 '{"primary_color":"#3B82F6","secondary_color":"#1E40AF","font_family":"inter","button_style":"rounded"}',
 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1002, 0, 'Modern', 'modern', 'A sleek and contemporary design with bold colors and dynamic layouts. Ideal for fashion and lifestyle brands.',
 'https://cdn.shopjoy.com/themes/modern-thumb.png', 'https://cdn.shopjoy.com/themes/modern-preview.png',
 '{"colors":{"primary":"#10B981","secondary":"#059669"},"fonts":{"heading":"Poppins","body":"Inter"},"layout":"modern"}',
 '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#10B981"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#059669"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"poppins","label":"Poppins"},{"value":"montserrat","label":"Montserrat"}],"default":"poppins"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"pill","label":"Pill"},{"value":"rounded","label":"Rounded"}],"default":"pill"}]',
 '{"primary_color":"#10B981","secondary_color":"#059669","font_family":"poppins","button_style":"pill"}',
 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1003, 0, 'Minimal', 'minimal', 'A clean and minimalist design that puts your products front and center. Great for luxury and high-end brands.',
 'https://cdn.shopjoy.com/themes/minimal-thumb.png', 'https://cdn.shopjoy.com/themes/minimal-preview.png',
 '{"colors":{"primary":"#000000","secondary":"#6B7280"},"fonts":{"heading":"Helvetica Neue","body":"Helvetica Neue"},"layout":"minimal"}',
 '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#000000"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#6B7280"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"helvetica","label":"Helvetica Neue"},{"value":"arial","label":"Arial"}],"default":"helvetica"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"underline","label":"Underline"},{"value":"solid","label":"Solid"}],"default":"underline"}]',
 '{"primary_color":"#000000","secondary_color":"#6B7280","font_family":"helvetica","button_style":"underline"}',
 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1004, 0, 'Bold', 'bold', 'A vibrant and energetic design with strong visual elements. Perfect for creative brands and youth-oriented products.',
 'https://cdn.shopjoy.com/themes/bold-thumb.png', 'https://cdn.shopjoy.com/themes/bold-preview.png',
 '{"colors":{"primary":"#8B5CF6","secondary":"#6D28D9"},"fonts":{"heading":"DM Sans","body":"DM Sans"},"layout":"bold"}',
 '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#8B5CF6"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#6D28D9"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"dmsans","label":"DM Sans"},{"value":"nunito","label":"Nunito"}],"default":"dmsans"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"pill","label":"Pill"}],"default":"rounded"}]',
 '{"primary_color":"#8B5CF6","secondary_color":"#6D28D9","font_family":"dmsans","button_style":"rounded"}',
 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1005, 0, 'Nature', 'nature', 'An organic and eco-friendly design with natural colors and earthy tones. Ideal for sustainable and wellness brands.',
 'https://cdn.shopjoy.com/themes/nature-thumb.png', 'https://cdn.shopjoy.com/themes/nature-preview.png',
 '{"colors":{"primary":"#059669","secondary":"#047857"},"fonts":{"heading":"Merriweather","body":"Open Sans"},"layout":"nature"}',
 '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#059669"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#047857"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"merriweather","label":"Merriweather"},{"value":"lora","label":"Lora"}],"default":"merriweather"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"leaf","label":"Leaf"}],"default":"rounded"}]',
 '{"primary_color":"#059669","secondary_color":"#047857","font_family":"merriweather","button_style":"rounded"}',
 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- ============================================
-- 8. Create default pages for existing shops
-- ============================================

-- Create default home pages for existing shops that don't have one
INSERT INTO `pages` (`tenant_id`, `name`, `slug`, `type`, `content`, `seo_title`, `seo_description`, `seo_keywords`, `status`, `sort`, `is_published`, `version`, `created_at`, `updated_at`, `created_by`, `updated_by`)
SELECT
    s.tenant_id,
    'Home',
    'home',
    0,
    '{}',
    CONCAT(s.name, ' - Home'),
    CONCAT('Welcome to ', s.name),
    '',
    1,
    0,
    0,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0,
    0
FROM `shops` s
WHERE NOT EXISTS (
    SELECT 1 FROM `pages` p WHERE p.tenant_id = s.tenant_id AND p.slug = 'home' AND (p.deleted_at IS NULL OR p.deleted_at = 0)
);

-- ============================================
-- 9. Create default global SEO configs
-- ============================================

INSERT INTO `seo_configs` (`tenant_id`, `page_type`, `page_id`, `title`, `description`, `keywords`, `created_at`, `updated_at`)
SELECT
    s.tenant_id,
    'global',
    NULL,
    s.seo_title,
    s.seo_description,
    s.seo_keywords,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP()
FROM `shops` s;

-- Create default page type SEO configs
INSERT INTO `seo_configs` (`tenant_id`, `page_type`, `page_id`, `title`, `description`, `keywords`, `created_at`, `updated_at`)
SELECT
    s.tenant_id,
    'home',
    NULL,
    CONCAT(s.name, ' - Home'),
    CONCAT('Welcome to ', s.name, ' online store'),
    '',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP()
FROM `shops` s;

INSERT INTO `seo_configs` (`tenant_id`, `page_type`, `page_id`, `title`, `description`, `keywords`, `created_at`, `updated_at`)
SELECT
    s.tenant_id,
    'product',
    NULL,
    CONCAT('{{product_name}} - ', s.name),
    '{{product_description}}',
    '',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP()
FROM `shops` s;

INSERT INTO `seo_configs` (`tenant_id`, `page_type`, `page_id`, `title`, `description`, `keywords`, `created_at`, `updated_at`)
SELECT
    s.tenant_id,
    'category',
    NULL,
    CONCAT('{{category_name}} - ', s.name),
    'Browse our {{category_name}} collection',
    '',
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP()
FROM `shops` s;

-- ============================================
-- 10. Set default themes for existing shops
-- ============================================

UPDATE `shops` SET `current_theme_id` = 1001 WHERE `current_theme_id` IS NULL;

-- ============================================
-- Migration completed
-- ============================================