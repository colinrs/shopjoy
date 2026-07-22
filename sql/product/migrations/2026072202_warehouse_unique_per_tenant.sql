-- ============================================
-- FU-2.4: warehouses unique code per tenant (soft-delete aware)
-- ============================================
--
-- Context:
--   The `warehouses` table historically carried a UNIQUE key
--   `idx_tenant_code (tenant_id, code)`. This is per-tenant unique, which is
--   correct, BUT it does not include `deleted_at`. Consequences:
--     1. A soft-deleted warehouse (deleted_at IS NOT NULL) still occupies the
--        (tenant_id, code) slot, so an operator cannot re-create a warehouse
--        with the same code after deleting the old one.
--     2. The index name diverges from the project convention
--        `uk_tenant_code (tenant_id, code, deleted_at)` used by the
--        fulfillment domain and elsewhere.
--
-- This migration aligns warehouses with the canonical pattern:
--     UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`, `deleted_at`)
--
-- Safety / idempotency:
--   MySQL has no "DROP INDEX IF EXISTS", so we guard each DDL with
--   information_schema checks executed via a prepared statement. Running this
--   migration twice is a no-op.
--
-- Ordering note:
--   Run AFTER 2026072201_report_warehouse_orphan_tenant.sql. Any orphan
--   (tenant_id = 0) rows are reported there and must be re-parented by the
--   operator; they do NOT block this index change because the composite key
--   still permits distinct codes under tenant_id = 0.

-- --------------------------------------------------------------------------
-- 1. Drop the legacy UNIQUE key `idx_tenant_code` if present.
-- --------------------------------------------------------------------------
SET @drop_idx_tenant_code := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA = DATABASE()
               AND TABLE_NAME   = 'warehouses'
               AND INDEX_NAME   = 'idx_tenant_code'
        ),
        'ALTER TABLE `warehouses` DROP INDEX `idx_tenant_code`',
        'SELECT 1'
    )
);
PREPARE stmt FROM @drop_idx_tenant_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------------------------
-- 2. Drop any legacy global-unique `uk_code` (defensive; may not exist).
-- --------------------------------------------------------------------------
SET @drop_uk_code := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA = DATABASE()
               AND TABLE_NAME   = 'warehouses'
               AND INDEX_NAME   = 'uk_code'
        ),
        'ALTER TABLE `warehouses` DROP INDEX `uk_code`',
        'SELECT 1'
    )
);
PREPARE stmt FROM @drop_uk_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- --------------------------------------------------------------------------
-- 3. Add the canonical soft-delete-aware unique key if not already present.
-- --------------------------------------------------------------------------
SET @add_uk_tenant_code := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM information_schema.STATISTICS
             WHERE TABLE_SCHEMA = DATABASE()
               AND TABLE_NAME   = 'warehouses'
               AND INDEX_NAME   = 'uk_tenant_code'
        ),
        'SELECT 1',
        'ALTER TABLE `warehouses` ADD UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`, `deleted_at`)'
    )
);
PREPARE stmt FROM @add_uk_tenant_code;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
