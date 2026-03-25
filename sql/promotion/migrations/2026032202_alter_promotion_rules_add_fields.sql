-- Add new columns to promotion_rules table
ALTER TABLE promotion_rules ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'CNY';
ALTER TABLE promotion_rules ADD COLUMN IF NOT EXISTS sort_order INT NOT NULL DEFAULT 0;

-- Backfill existing records
UPDATE promotion_rules SET currency = 'CNY' WHERE currency = '' OR currency IS NULL;

-- Add index
CREATE INDEX IF NOT EXISTS idx_promotion_rules_sort_order ON promotion_rules(promotion_id, sort_order);