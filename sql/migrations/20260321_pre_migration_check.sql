-- ============================================
-- Pre-Migration Validation Script
-- Date: 2026-03-21
-- Description: Run this BEFORE starting the migration to identify potential issues
-- ============================================

-- ============================================
-- 1. Check for products without market assignment
--    (These will have tenant_id = 0 after migration)
-- ============================================
SELECT 'Products without market assignment:' AS check_name;
SELECT p.id, p.name, p.sku
FROM products p
LEFT JOIN product_markets pm ON p.id = pm.product_id
WHERE pm.id IS NULL;

-- ============================================
-- 2. Check for products with markets from different tenants
--    (Data integrity issue - products should only be in one tenant's markets)
-- ============================================
SELECT 'Products with cross-tenant markets:' AS check_name;
SELECT p.id, p.name, COUNT(DISTINCT m.tenant_id) as tenant_count
FROM products p
JOIN product_markets pm ON p.id = pm.product_id
JOIN markets m ON pm.market_id = m.id
GROUP BY p.id
HAVING COUNT(DISTINCT m.tenant_id) > 1;

-- ============================================
-- 3. Count records to be migrated
-- ============================================
SELECT 'Record counts:' AS check_name;
SELECT 'products' AS table_name, COUNT(*) AS count FROM products
UNION ALL
SELECT 'skus', COUNT(*) FROM skus
UNION ALL
SELECT 'categories', COUNT(*) FROM categories
UNION ALL
SELECT 'brands', COUNT(*) FROM brands;

-- ============================================
-- 4. Check for duplicate brand names within tenants
--    (Migration might fail if same brand name exists in same tenant)
-- ============================================
SELECT 'Duplicate brand names within tenants:' AS check_name;
SELECT p.tenant_id, TRIM(p.brand) AS brand_name, COUNT(*) AS product_count
FROM products p
WHERE p.brand IS NOT NULL AND p.brand != ''
GROUP BY p.tenant_id, TRIM(p.brand)
HAVING COUNT(*) > 1;

-- ============================================
-- 5. Verify backup tables don't already exist
-- ============================================
SELECT 'Checking for existing backup tables:' AS check_name;
SELECT TABLE_NAME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
AND TABLE_NAME LIKE '%_backup_20260321';

-- ============================================
-- 6. Check MySQL version (for compatibility)
-- ============================================
SELECT 'MySQL Version:' AS check_name;
SELECT VERSION();

-- ============================================
-- 7. Check for NULL brand values that should be handled
-- ============================================
SELECT 'Products with NULL or empty brand:' AS check_name;
SELECT COUNT(*) AS count
FROM products
WHERE brand IS NULL OR brand = '';

-- ============================================
-- 8. Check for categories with NULL values
-- ============================================
SELECT 'Categories with missing required fields:' AS check_name;
SELECT id, name
FROM categories
WHERE name IS NULL OR name = '';

-- ============================================
-- Instructions
-- ============================================
SELECT '========================================' AS instructions;
SELECT 'Pre-Migration Checklist:' AS instructions;
SELECT '========================================' AS instructions;
SELECT '1. If products without markets exist, assign them to a market first' AS instructions;
SELECT '2. If cross-tenant products exist, investigate and fix data integrity' AS instructions;
SELECT '3. Run backup script before migration' AS instructions;
SELECT '4. After verification, run the migration scripts in order:' AS instructions;
SELECT '   - 20260321_category_brand_inventory.sql' AS instructions;
SELECT '   - 20260321_data_migration.sql' AS instructions;
SELECT '========================================' AS instructions;