START TRANSACTION;

-- 1) Backup old coupons (idempotent: TRUNCATE then re-populate so re-runs are clean)
CREATE TABLE IF NOT EXISTS _deprecated_coupons LIKE coupons;
TRUNCATE TABLE _deprecated_coupons;
INSERT INTO _deprecated_coupons SELECT * FROM coupons;

-- 2) promotions: add new columns (idempotent via information_schema dynamic SQL;
--    MySQL 9.0 doesn't support ADD COLUMN IF NOT EXISTS inside ALTER TABLE statements)
--    Note: scope_type / scope_ids / exclude_ids / usage_limit / per_user_limit / tags already exist on promotions
--    from the previous SDD run (2026-07-18 promotion-usage-limit-tags merge). Only the truly-new columns are added here.
SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotions' AND column_name='kind');
SET @stmt := IF(@c=0, "ALTER TABLE promotions ADD COLUMN `kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `tenant_id`", 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotions' AND column_name='code');
SET @stmt := IF(@c=0, 'ALTER TABLE promotions ADD COLUMN `code` VARCHAR(100) NULL AFTER `name`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotions' AND column_name='market_id');
SET @stmt := IF(@c=0, 'ALTER TABLE promotions ADD COLUMN `market_id` BIGINT NULL AFTER `priority`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotions' AND column_name='total_count');
SET @stmt := IF(@c=0, 'ALTER TABLE promotions ADD COLUMN `total_count` INT NULL AFTER `usage_limit`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotions' AND column_name='used_count');
SET @stmt := IF(@c=0, 'ALTER TABLE promotions ADD COLUMN `used_count` INT NULL AFTER `total_count`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 3) Migrate coupons → promotions
--    Note: coupons table has NO usage_limit column (only per_user_limit), so we
--    emit NULL for usage_limit during the merge.
INSERT INTO promotions (
  tenant_id, kind, name, description, code, type, status, priority, market_id, currency,
  total_count, used_count, usage_limit, per_user_limit, scope_type, scope_ids, exclude_ids,
  start_at, end_at, created_at, updated_at, created_by, updated_by
)
SELECT
  tenant_id,
  'COUPON'                                                          AS kind,
  name, description, code,
  0                                                                 AS type,
  CASE status WHEN 0 THEN 0 WHEN 1 THEN 1 WHEN 2 THEN 3 WHEN 3 THEN 3 END AS status,
  0                                                                 AS priority,
  NULL                                                              AS market_id,
  currency,
  total_count, used_count,
  0                                                                 AS usage_limit,
  per_user_limit,
  COALESCE(scope_type, 'STOREWIDE'),
  scope_ids, exclude_ids,
  start_at, end_at,
  created_at, updated_at, created_by, updated_by
FROM coupons
WHERE deleted_at IS NULL;

-- 4) Backup promotion_rules (idempotent: TRUNCATE then re-populate)
CREATE TABLE IF NOT EXISTS _deprecated_promotion_rules LIKE promotion_rules;
TRUNCATE TABLE _deprecated_promotion_rules;
INSERT INTO _deprecated_promotion_rules SELECT * FROM promotion_rules;

-- 5) Rebuild promotion_rules with owner_kind + owner_id
--    sort_order already exists from previous SDD run.
--    MySQL 9.0 does NOT support ADD COLUMN IF NOT EXISTS, so use dynamic SQL per column.
SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotion_rules' AND column_name='owner_kind');
SET @stmt := IF(@c=0, "ALTER TABLE promotion_rules ADD COLUMN `owner_kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `promotion_id`", 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotion_rules' AND column_name='owner_id');
SET @stmt := IF(@c=0, 'ALTER TABLE promotion_rules ADD COLUMN `owner_id` BIGINT NOT NULL DEFAULT 0 AFTER `owner_kind`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @c := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='promotion_rules' AND column_name='promotion_id' AND is_nullable='YES');
SET @stmt := IF(@c=0, 'ALTER TABLE promotion_rules MODIFY COLUMN `promotion_id` BIGINT NULL', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 6) Tag existing PROMOTION rules
UPDATE promotion_rules SET owner_kind = 'PROMOTION' WHERE promotion_id IS NOT NULL;

-- 7) Backfill owner_id from promotion_id for PROMOTION rows
UPDATE promotion_rules SET owner_id = promotion_id WHERE owner_kind = 'PROMOTION' AND owner_id = 0;

-- 8) Convert each COUPON's (type/value/min_amount/max_discount) into one rule row
INSERT INTO promotion_rules (
  owner_kind, owner_id, condition_type, condition_value, action_type, action_value,
  max_discount_amount, max_discount_currency, currency, sort_order,
  created_at, updated_at
)
SELECT
  'COUPON'                                                              AS owner_kind,
  p.id                                                                  AS owner_id,
  0                                                                     AS condition_type,
  c.min_amount                                                          AS condition_value,
  CASE c.type WHEN 0 THEN 0 WHEN 1 THEN 1 WHEN 2 THEN 0 END             AS action_type,
  CASE c.type WHEN 0 THEN c.value WHEN 1 THEN c.value WHEN 2 THEN 0 END AS action_value,
  c.max_discount                                                        AS max_discount_amount,
  c.currency                                                            AS max_discount_currency,
  c.currency,
  0                                                                     AS sort_order,
  NOW(), NOW()
FROM coupons c
JOIN promotions p ON p.kind = 'COUPON' AND p.code = c.code AND p.tenant_id = c.tenant_id
WHERE c.deleted_at IS NULL;

-- 9) Index rebuild (idempotent via information_schema checks + dynamic SQL)
SET @idx_owner_exists := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'promotion_rules' AND index_name = 'idx_owner');
SET @stmt := IF(@idx_owner_exists = 0, 'CREATE INDEX `idx_owner` ON `promotion_rules` (`owner_kind`, `owner_id`, `sort_order`)', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx_old1 := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'promotion_rules' AND index_name = 'idx_promotion_id');
SET @stmt := IF(@idx_old1 > 0, 'DROP INDEX `idx_promotion_id` ON `promotion_rules`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx_old2 := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'promotion_rules' AND index_name = 'idx_promotion_rules_sort_order');
SET @stmt := IF(@idx_old2 > 0, 'DROP INDEX `idx_promotion_rules_sort_order` ON `promotion_rules`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 10) COUPON.code partial-unique via generated column (idempotent)
SET @col_exists := (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'promotions' AND column_name = 'code_unique');
SET @stmt := IF(@col_exists = 0,
  'ALTER TABLE `promotions` ADD COLUMN `code_unique` VARCHAR(100) GENERATED ALWAYS AS (IF(kind = ''COUPON'', code, NULL)) VIRTUAL',
  'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @uk_exists := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'promotions' AND index_name = 'uk_promotion_code');
SET @stmt := IF(@uk_exists = 0, 'ALTER TABLE `promotions` ADD UNIQUE KEY `uk_promotion_code` (`code_unique`)', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 11) user_coupons aid index (idempotent)
SET @idx_uc := (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'user_coupons' AND index_name = 'idx_coupon_id_active');
SET @stmt := IF(@idx_uc = 0, 'ALTER TABLE `user_coupons` ADD INDEX `idx_coupon_id_active` (`coupon_id`, `status`)', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 12) Archive old coupons table (only if not already archived)
SET @archived_exists := (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '_archived_coupons_20260719');
SET @stmt := IF(@archived_exists = 0, 'RENAME TABLE `coupons` TO `_archived_coupons_20260719`', 'SELECT 1');
PREPARE stmt FROM @stmt; EXECUTE stmt; DEALLOCATE PREPARE stmt;

COMMIT;