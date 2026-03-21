-- ============================================
-- Data Migration: Migrate tenant_id and brand_id
-- Date: 2026-03-21
-- Description: Migrate tenant_id from markets to products/skus,
--              migrate brand string to brand_id, and
--              migrate stock to available_stock
-- Prerequisites: Run 20260321_category_brand_inventory.sql first
-- ============================================

-- ============================================
-- 1. Migrate tenant_id from markets to products
-- ============================================
-- Update products.tenant_id from product_markets -> markets
UPDATE products p
    INNER JOIN product_markets pm ON p.id = pm.product_id
    INNER JOIN markets m ON pm.market_id = m.id
    SET p.tenant_id = m.tenant_id
    WHERE p.tenant_id = 0;

-- Verify migration
SELECT 'Products tenant_id migration:' AS step;
SELECT COUNT(*) AS migrated_count FROM products WHERE tenant_id > 0;
SELECT COUNT(*) AS unmigrated_count FROM products WHERE tenant_id = 0;

-- ============================================
-- 2. Migrate tenant_id from products to skus
-- ============================================
UPDATE skus s
    INNER JOIN products p ON s.product_id = p.id
    SET s.tenant_id = p.tenant_id
    WHERE s.tenant_id = 0;

-- Verify migration
SELECT 'SKUs tenant_id migration:' AS step;
SELECT COUNT(*) AS migrated_count FROM skus WHERE tenant_id > 0;
SELECT COUNT(*) AS unmigrated_count FROM skus WHERE tenant_id = 0;

-- ============================================
-- 3. Migrate brand string to brand_id
-- ============================================
-- Create brands from distinct brand names in products
INSERT INTO brands (tenant_id, name, sort, status, created_at, updated_at, created_by, updated_by)
SELECT DISTINCT
    p.tenant_id,
    TRIM(p.brand) AS name,
    0 AS sort,
    1 AS status,
    UNIX_TIMESTAMP() AS created_at,
    UNIX_TIMESTAMP() AS updated_at,
    0 AS created_by,
    0 AS updated_by
FROM products p
WHERE p.brand IS NOT NULL
  AND p.brand != ''
  AND NOT EXISTS (
      SELECT 1 FROM brands b
      WHERE b.name = TRIM(p.brand)
        AND b.tenant_id = p.tenant_id
        AND b.deleted_at IS NULL
  );

-- Verify brands created
SELECT 'New brands created:' AS step;
SELECT COUNT(*) AS new_brands FROM brands WHERE created_at >= UNIX_TIMESTAMP('2026-03-21');

-- Update products to link to brand_id
UPDATE products p
    INNER JOIN brands b ON b.name = TRIM(p.brand) AND b.tenant_id = p.tenant_id
    SET p.brand_id = b.id
    WHERE p.brand IS NOT NULL AND p.brand != '';

-- Verify migration
SELECT 'Products with brand_id:' AS step;
SELECT COUNT(*) AS products_with_brand FROM products WHERE brand_id IS NOT NULL;
SELECT COUNT(*) AS products_without_brand FROM products WHERE brand_id IS NULL;

-- ============================================
-- 4. Migrate SKU stock to available_stock
-- ============================================
UPDATE skus SET available_stock = stock WHERE available_stock = 0;

-- Verify migration
SELECT 'SKU stock migration:' AS step;
SELECT id, code, stock, available_stock FROM skus LIMIT 10;

-- ============================================
-- 5. Final verification
-- ============================================
SELECT '=== Final Migration Summary ===' AS report;

SELECT 'Products' AS table_name,
       COUNT(*) AS total,
       SUM(CASE WHEN tenant_id > 0 THEN 1 ELSE 0 END) AS with_tenant,
       SUM(CASE WHEN brand_id IS NOT NULL THEN 1 ELSE 0 END) AS with_brand
FROM products;

SELECT 'SKUs' AS table_name,
       COUNT(*) AS total,
       SUM(CASE WHEN tenant_id > 0 THEN 1 ELSE 0 END) AS with_tenant,
       SUM(CASE WHEN available_stock > 0 THEN 1 ELSE 0 END) AS with_available_stock
FROM skus;

SELECT 'Brands' AS table_name,
       COUNT(*) AS total
FROM brands
WHERE deleted_at IS NULL;

SELECT 'Categories' AS table_name,
       COUNT(*) AS total
FROM categories
WHERE deleted_at IS NULL;