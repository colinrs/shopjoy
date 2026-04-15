-- Create refund_reasons table
CREATE TABLE IF NOT EXISTS `refund_reasons` (
    `id` BIGINT NOT NULL COMMENT 'ID',
    `code` VARCHAR(50) NOT NULL COMMENT '退款原因代码',
    `name` VARCHAR(100) NOT NULL COMMENT '退款原因名称',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `is_active` INT NOT NULL DEFAULT 1 COMMENT '是否启用: 0=禁用, 1=启用',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款原因表';

-- Insert default refund reasons
INSERT IGNORE INTO `refund_reasons` (`id`, `code`, `name`, `sort`, `is_active`, `created_at`) VALUES
(1, 'QUALITY', '商品质量问题', 10, 1, NOW()),
(2, 'WRONG_ITEM', '商品与描述不符', 20, 1, NOW()),
(3, 'NOT_RECEIVED', '未收到商品', 30, 1, NOW()),
(4, 'CHANGED_MIND', '买家改变主意', 40, 1, NOW()),
(5, 'DUPLICATE_ORDER', '重复下单', 50, 1, NOW()),
(6, 'OTHER', '其他原因', 100, 1, NOW());
