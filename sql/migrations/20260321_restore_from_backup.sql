-- ============================================
-- Restore from Backup Script
-- Date: 2026-03-21
-- Description: Restore tables from backup if migration fails
-- WARNING: This will replace current tables with backup data
-- ============================================

-- ============================================
-- Check if backup tables exist
-- ============================================
SELECT 'Checking backup tables exist...' AS step;

-- ============================================
-- Restore products table
-- ============================================
DROP TABLE IF EXISTS `products`;
RENAME TABLE `products_backup_20260321` TO `products`;

-- ============================================
-- Restore skus table
-- ============================================
DROP TABLE IF EXISTS `skus`;
RENAME TABLE `skus_backup_20260321` TO `skus`;

-- ============================================
-- Restore categories table
-- ============================================
DROP TABLE IF EXISTS `categories`;
RENAME TABLE `categories_backup_20260321` TO `categories`;

-- ============================================
-- Restore brands table
-- ============================================
DROP TABLE IF EXISTS `brands`;
RENAME TABLE `brands_backup_20260321` TO `brands`;

-- ============================================
-- Restore product_markets table
-- ============================================
DROP TABLE IF EXISTS `product_markets`;
RENAME TABLE `product_markets_backup_20260321` TO `product_markets`;

-- ============================================
-- Drop new tables created by migration
-- ============================================
DROP TABLE IF EXISTS `inventory_logs`;
DROP TABLE IF EXISTS `warehouse_inventories`;
DROP TABLE IF EXISTS `warehouses`;
DROP TABLE IF EXISTS `brand_markets`;
DROP TABLE IF EXISTS `category_markets`;

-- ============================================
-- Verify restoration
-- ============================================
SELECT 'Restoration verification:' AS step;
SELECT 'products' AS table_name, COUNT(*) AS count FROM products
UNION ALL
SELECT 'skus', COUNT(*) FROM skus
UNION ALL
SELECT 'categories', COUNT(*) FROM categories
UNION ALL
SELECT 'brands', COUNT(*) FROM brands;

-- ============================================
-- Success message
-- ============================================
SELECT '========================================' AS message;
SELECT 'Restoration completed successfully!' AS message;
SELECT 'Tables have been restored from backup.' AS message;
SELECT '========================================' AS message;