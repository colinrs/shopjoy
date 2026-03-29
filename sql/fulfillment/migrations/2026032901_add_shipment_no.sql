-- Add shipment_no field to shipments table
-- Migration ID: 2026032901
-- Created: 2026-03-29

ALTER TABLE shipments ADD COLUMN shipment_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT '发货单号' AFTER order_id;

-- Add unique index for shipment_no
ALTER TABLE shipments ADD UNIQUE INDEX uk_shipment_no (tenant_id, shipment_no);
