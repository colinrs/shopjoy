-- Create carriers table for logistics tracking
CREATE TABLE IF NOT EXISTS `carriers` (
    `id` BIGINT NOT NULL COMMENT 'ID',
    `code` VARCHAR(20) NOT NULL COMMENT '物流公司代码',
    `name` VARCHAR(100) NOT NULL COMMENT '物流公司名称',
    `tracking_url` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '物流跟踪URL模板',
    `is_active` INT NOT NULL DEFAULT 1 COMMENT '是否启用: 0=禁用, 1=启用',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='物流公司表';

-- Insert default carriers
INSERT IGNORE INTO `carriers` (`id`, `code`, `name`, `tracking_url`, `is_active`, `sort`, `created_at`) VALUES
(1, 'SF', '顺丰速运', 'https://www.sf-express.com/cn/sc/query/{tracking_no}', 1, 10, NOW()),
(2, 'YTO', '圆通速递', 'https://www.yto.net.cn/cn/query/{tracking_no}', 1, 20, NOW()),
(3, 'ZTO', '中通快递', 'https://www.zto.com/express/query/{tracking_no}', 1, 30, NOW()),
(4, 'STO', '申通快递', 'https://www.sto.cn/query/{tracking_no}', 1, 40, NOW()),
(5, 'YD', '韵达快递', 'https://www.yunda56.com/cms/query/{tracking_no}', 1, 50, NOW()),
(6, 'EMS', 'EMS', 'https://www.ems.com.cn/query/{tracking_no}', 1, 60, NOW()),
(7, 'JD', '京东物流', 'https://order.jd.com/center/search.action?keyword={tracking_no}', 1, 70, NOW()),
(8, 'DBL', '德邦快递', 'https://www.deppon.com/track/queryByTrackNo/?trackNo={tracking_no}', 1, 80, NOW());
