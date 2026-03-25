-- Add fulfillment and price adjustment fields to orders table
-- Migration: 20260323_add_order_fields

ALTER TABLE `orders`
ADD COLUMN `fulfillment_status` TINYINT NOT NULL DEFAULT 0 COMMENT '履约状态: 0-待发货, 1-部分发货, 2-已发货, 3-已送达' AFTER `status`,
ADD COLUMN `refund_status` TINYINT NOT NULL DEFAULT 0 COMMENT '退款状态: 0-无, 1-待处理, 2-已批准, 3-已拒绝, 4-已完成' AFTER `fulfillment_status`,
ADD COLUMN `merchant_remark` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '商家内部备注' AFTER `remark`,
ADD COLUMN `original_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '改价前原金额(分)' AFTER `pay_amount`,
ADD COLUMN `adjust_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '改价金额(分)' AFTER `original_amount`,
ADD COLUMN `adjust_reason` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '改价原因' AFTER `adjust_amount`,
ADD COLUMN `adjusted_by` BIGINT NOT NULL DEFAULT 0 COMMENT '改价操作人ID' AFTER `adjust_reason`,
ADD COLUMN `adjusted_at` BIGINT DEFAULT NULL COMMENT '改价时间' AFTER `adjusted_by`,
ADD COLUMN `version` INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号' AFTER `adjusted_at`,
ADD COLUMN `payment_method` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '支付方式' AFTER `version`,
ADD COLUMN `source` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '订单来源' AFTER `payment_method`;

-- Migrate existing data: set original_amount = pay_amount for existing orders
UPDATE `orders` SET `original_amount` = `pay_amount` WHERE `original_amount` = 0;

-- Add indexes for new fields
ALTER TABLE `orders` ADD INDEX `idx_fulfillment_status` (`fulfillment_status`);
ALTER TABLE `orders` ADD INDEX `idx_refund_status` (`refund_status`);