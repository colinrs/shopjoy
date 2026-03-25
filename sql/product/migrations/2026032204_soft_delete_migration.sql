-- ============================================
-- Migration: Add deleted_at for Soft Delete
-- Date: 2026-03-22
-- Description: Add deleted_at column to products and users tables
--              for unified soft delete mechanism
-- ============================================

-- ============================================
-- 1. Add deleted_at to products table
-- ============================================
ALTER TABLE `products`
    ADD COLUMN `deleted_at` BIGINT NULL COMMENT '删除时间' AFTER `dangerous_goods`,
    ADD INDEX `idx_deleted_at` (`deleted_at`);

-- ============================================
-- 2. Add deleted_at to users table
-- ============================================
ALTER TABLE `users`
    ADD COLUMN `deleted_at` BIGINT NULL COMMENT '删除时间' AFTER `last_login`,
    ADD INDEX `idx_deleted_at` (`deleted_at`);

-- ============================================
-- 3. Migrate existing deleted records
-- ============================================
-- Products with status = deleted (3) should have deleted_at set
UPDATE `products`
SET `deleted_at` = UNIX_TIMESTAMP()
WHERE `status` = 3 AND `deleted_at` IS NULL;

-- Users with status = deleted (3) should have deleted_at set
UPDATE `users`
SET `deleted_at` = UNIX_TIMESTAMP()
WHERE `status` = 3 AND `deleted_at` IS NULL;