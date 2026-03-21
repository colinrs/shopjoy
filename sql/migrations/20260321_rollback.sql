-- ============================================
-- Rollback: Category, Brand, and Inventory Management
-- Date: 2026-03-21
-- Description: Revert all schema changes from 20260321 migration
-- WARNING: This will DROP columns and tables. Data will be lost.
-- ============================================

-- ============================================
-- 1. Remove inventory fields from skus table
-- ============================================
ALTER TABLE `skus`
    DROP COLUMN `available_stock`,
    DROP COLUMN `locked_stock`,
    DROP COLUMN `safety_stock`,
    DROP COLUMN `presale_enabled`;

-- ============================================
-- 2. Remove brand_id from products table
-- ============================================
ALTER TABLE `products`
    DROP INDEX `idx_brand_id`,
    DROP COLUMN `brand_id`;

-- ============================================
-- 3. Remove compliance fields from brands table
-- ============================================
ALTER TABLE `brands`
    DROP COLUMN `enable_page`,
    DROP COLUMN `trademark_number`,
    DROP COLUMN `trademark_country`,
    DROP COLUMN `deleted_at`,
    DROP INDEX `idx_deleted_at`;

-- ============================================
-- 4. Remove SEO fields from categories table
-- ============================================
ALTER TABLE `categories`
    DROP COLUMN `seo_title`,
    DROP COLUMN `seo_description`,
    DROP COLUMN `deleted_at`,
    DROP INDEX `idx_deleted_at`;

-- ============================================
-- 5. Remove tenant_id from skus table
-- ============================================
ALTER TABLE `skus`
    DROP INDEX `idx_tenant_id`,
    DROP COLUMN `tenant_id`;

-- ============================================
-- 6. Remove tenant_id from products table
-- ============================================
ALTER TABLE `products`
    DROP INDEX `idx_tenant_id`,
    DROP COLUMN `tenant_id`;

-- ============================================
-- 7. Drop new tables
-- ============================================
DROP TABLE IF EXISTS `inventory_logs`;
DROP TABLE IF EXISTS `warehouse_inventories`;
DROP TABLE IF EXISTS `warehouses`;
DROP TABLE IF EXISTS `brand_markets`;
DROP TABLE IF EXISTS `category_markets`;

-- ============================================
-- 8. Delete newly created brands (from migration)
-- ============================================
DELETE FROM brands WHERE created_at >= UNIX_TIMESTAMP('2026-03-21');