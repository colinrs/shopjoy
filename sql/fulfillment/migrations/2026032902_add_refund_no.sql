-- Add refund_no field to refunds table
-- Migration ID: 2026032902
-- Created: 2026-03-29

ALTER TABLE refunds ADD COLUMN refund_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT '退款单号' AFTER order_id;

-- Add unique index for refund_no
ALTER TABLE refunds ADD UNIQUE INDEX uk_refund_no (tenant_id, refund_no);
