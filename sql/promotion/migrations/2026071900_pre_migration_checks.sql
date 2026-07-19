-- ============================================
-- Promotion × Coupon merge: pre-migration checks
-- Run BEFORE 2026071901_merge_promotion_coupon.sql
-- Any non-zero row count = BLOCK migration
-- ============================================

SELECT 'CHECK 1: coupon code duplicates (BLOCK if > 0)' AS check_name;
SELECT code, COUNT(*) AS c
FROM coupons
WHERE deleted_at IS NULL
GROUP BY code
HAVING c > 1;

SELECT 'CHECK 2: orphan promotion_rules (BLOCK if > 0)' AS check_name;
SELECT r.*
FROM promotion_rules r
LEFT JOIN promotions p ON p.id = r.promotion_id
WHERE p.id IS NULL;

-- CHECK 3 (was: market_ids content) removed — coupons.market_ids column
-- does not exist in the live schema. Per spec §2.7 the "discard market_ids"
-- strategy is moot because no data exists to discard. Promotion.market_id is
-- a NEW column that will be created by the merge migration, not populated
-- from existing data.

SELECT 'CHECK 4: existing row counts (baseline for post-migration comparison)' AS check_name;
SELECT 'promotions' AS tbl, COUNT(*) AS n FROM promotions UNION ALL
SELECT 'promotion_rules', COUNT(*) FROM promotion_rules UNION ALL
SELECT 'coupons', COUNT(*) FROM coupons WHERE deleted_at IS NULL UNION ALL
SELECT 'user_coupons', COUNT(*) FROM user_coupons;