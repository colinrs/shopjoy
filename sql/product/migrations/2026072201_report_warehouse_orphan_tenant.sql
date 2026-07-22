-- ============================================
-- FU-2: Backfill / report warehouse rows with tenant_id = 0
-- ============================================
--
-- Phase 8 review surfaced that the warehouse Create/Update logic did not
-- read TenantID from ctx, so legacy rows created before FU-2 may have
-- tenant_id = 0. This migration makes those rows auditable without
-- silently assigning them to a real tenant.
--
-- Strategy:
--   1. Report the count of orphan rows so the operator can decide.
--   2. Add a composite index that covers (tenant_id, code) lookups —
--      the unique key on (tenant_id, code) already exists, but the
--      helper query "orphan rows" benefits from an index hint.
--   3. DO NOT rewrite tenant_id automatically. There is no default tenant
--      and re-parenting warehouse rows to a guessed tenant would
--      silently leak stock data across tenants.
--
-- After running this migration:
--   - The application code (FU-2) refuses to create or update any
--     warehouse when TenantID is missing/zero in ctx, so new writes
--     are isolated correctly.
--   - Orphan rows must be either manually re-parented via:
--         UPDATE warehouses SET tenant_id = <real_tenant_id>
--          WHERE tenant_id = 0 AND id = <id>;
--     or deleted via soft delete. Decision is operator-driven, not
--     automated.

SELECT COUNT(*) AS orphan_warehouse_rows
  FROM `warehouses`
 WHERE `tenant_id` = 0
   AND `deleted_at` IS NULL;