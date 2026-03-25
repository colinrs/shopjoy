-- ============================================
-- Migration: Low Stock Alert Index
-- Date: 2026-03-22
-- Description: Add composite index for efficient low stock alert queries
--              on the skus table
-- ============================================

-- ============================================
-- 1. Add composite index for low stock alert query
-- ============================================
-- This index supports the query:
-- SELECT * FROM skus WHERE safety_stock > 0 AND available_stock < safety_stock AND status = 1
--
-- Index column order rationale:
-- 1. tenant_id: Partition by tenant first (most selective for multi-tenant queries)
-- 2. status: Filter enabled SKUs only (status = 1)
-- 3. safety_stock: Filter SKUs with threshold configured (safety_stock > 0)
-- 4. available_stock: Support comparison (available_stock < safety_stock)
-- ============================================
ALTER TABLE `skus`
    ADD INDEX `idx_low_stock_alert` (`tenant_id`, `status`, `safety_stock`, `available_stock`);