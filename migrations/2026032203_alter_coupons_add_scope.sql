-- Add new columns to coupons table for scope support
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'CNY';
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS scope_type VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE';
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS scope_ids JSON DEFAULT NULL;
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS deleted_at BIGINT DEFAULT NULL;

-- Backfill existing records
UPDATE coupons SET currency = 'CNY', scope_type = 'STOREWIDE' WHERE currency = '' OR currency IS NULL;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_coupons_deleted_at ON coupons(deleted_at);