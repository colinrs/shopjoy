-- Add shipment cancellation fields
ALTER TABLE shipments
    ADD COLUMN `cancelled_at` TIMESTAMP NULL COMMENT '取消时间' AFTER `delivered_at`,
    ADD COLUMN `cancelled_by` BIGINT NOT NULL DEFAULT 0 COMMENT '取消人' AFTER `cancelled_at`,
    ADD COLUMN `cancelled_reason` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '取消原因' AFTER `cancelled_by`;
