-- Add new columns to promotions table for scope and priority support
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 0;
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'CNY';
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS scope_type VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE';
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS scope_ids JSON DEFAULT NULL;
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS exclude_ids JSON DEFAULT NULL;
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS deleted_at BIGINT DEFAULT NULL;

-- Backfill existing records with default values
UPDATE promotions SET currency = 'CNY', scope_type = 'STOREWIDE' WHERE currency = '' OR currency IS NULL;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_promotions_priority ON promotions(priority);
CREATE INDEX IF NOT EXISTS idx_promotions_deleted_at ON promotions(deleted_at);