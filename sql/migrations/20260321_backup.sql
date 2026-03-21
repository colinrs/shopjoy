-- ============================================
-- Backup Script
-- Date: 2026-03-21
-- Description: Create backup tables before running migration
-- Run this BEFORE running the migration scripts
-- ============================================

-- ============================================
-- Create backup tables
-- ============================================

-- Backup products table
CREATE TABLE IF NOT EXISTS `products_backup_20260321` AS SELECT * FROM `products`;

-- Backup skus table
CREATE TABLE IF NOT EXISTS `skus_backup_20260321` AS SELECT * FROM `skus`;

-- Backup categories table
CREATE TABLE IF NOT EXISTS `categories_backup_20260321` AS SELECT * FROM `categories`;

-- Backup brands table
CREATE TABLE IF NOT EXISTS `brands_backup_20260321` AS SELECT * FROM `brands`;

-- Backup product_markets table
CREATE TABLE IF NOT EXISTS `product_markets_backup_20260321` AS SELECT * FROM `product_markets`;

-- ============================================
-- Verify backups
-- ============================================
SELECT 'Backup verification:' AS step;
SELECT 'products_backup_20260321' AS table_name, COUNT(*) AS count FROM products_backup_20260321
UNION ALL
SELECT 'skus_backup_20260321', COUNT(*) FROM skus_backup_20260321
UNION ALL
SELECT 'categories_backup_20260321', COUNT(*) FROM categories_backup_20260321
UNION ALL
SELECT 'brands_backup_20260321', COUNT(*) FROM brands_backup_20260321
UNION ALL
SELECT 'product_markets_backup_20260321', COUNT(*) FROM product_markets_backup_20260321;

-- ============================================
-- Compare with original counts
-- ============================================
SELECT 'Original vs Backup comparison:' AS step;
SELECT
    'products' AS table_name,
    (SELECT COUNT(*) FROM products) AS original,
    (SELECT COUNT(*) FROM products_backup_20260321) AS backup
UNION ALL
SELECT
    'skus',
    (SELECT COUNT(*) FROM skus),
    (SELECT COUNT(*) FROM skus_backup_20260321)
UNION ALL
SELECT
    'categories',
    (SELECT COUNT(*) FROM categories),
    (SELECT COUNT(*) FROM categories_backup_20260321)
UNION ALL
SELECT
    'brands',
    (SELECT COUNT(*) FROM brands),
    (SELECT COUNT(*) FROM brands_backup_20260321);

-- ============================================
-- Success message
-- ============================================
SELECT '========================================' AS message;
SELECT 'Backup completed successfully!' AS message;
SELECT 'If migration fails, run restore script:' AS message;
SELECT '  source sql/migrations/20260321_restore_from_backup.sql' AS message;
SELECT '========================================' AS message;