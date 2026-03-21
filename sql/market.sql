-- ============================================
-- 市场表 (markets)
-- 跨境电商市场配置
-- ============================================

CREATE TABLE IF NOT EXISTS `markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '市场ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    `code` VARCHAR(10) NOT NULL COMMENT '市场代码: US, UK, DE, FR, AU',
    `name` VARCHAR(100) NOT NULL COMMENT '市场名称',
    `currency` VARCHAR(10) NOT NULL COMMENT '货币: USD, GBP, EUR, AUD',
    `default_language` VARCHAR(10) DEFAULT 'en' COMMENT '默认语言',
    `flag` VARCHAR(255) DEFAULT '' COMMENT '旗帜图标',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否主市场',
    `tax_rules` JSON COMMENT '税务配置',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    KEY `idx_code` (`code`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场表';

-- ============================================
-- 测试数据
-- ============================================

INSERT INTO `markets` (`id`, `tenant_id`, `code`, `name`, `currency`, `default_language`, `flag`, `is_active`, `is_default`, `tax_rules`, `created_at`, `updated_at`) VALUES
-- Demo Shop 市场
(1, 1, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', NOW(), NOW()),
(2, 1, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false, "GSTRate": 0}', NOW(), NOW()),
(3, 1, 'UK', '英国', 'GBP', 'en-GB', '🇬🇧', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', NOW(), NOW()),
(4, 1, 'DE', '德国', 'EUR', 'de', '🇩🇪', 1, 0, '{"IncludeTax": true, "VATRate": 0.19, "IOSSEnabled": true}', NOW(), NOW()),
(5, 1, 'AU', '澳大利亚', 'AUD', 'en-AU', '🇦🇺', 1, 0, '{"IncludeTax": true, "GSTRate": 0.10}', NOW(), NOW()),

-- Test Store 市场
(6, 2, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', NOW(), NOW()),
(7, 2, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false}', NOW(), NOW()),

-- Enterprise Corp 市场
(8, 3, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', NOW(), NOW()),
(9, 3, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false}', NOW(), NOW()),
(10, 3, 'UK', '英国', 'GBP', 'en-GB', '🇬🇧', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', NOW(), NOW()),
(11, 3, 'DE', '德国', 'EUR', 'de', '🇩🇪', 1, 0, '{"IncludeTax": true, "VATRate": 0.19, "IOSSEnabled": true}', NOW(), NOW()),
(12, 3, 'FR', '法国', 'EUR', 'fr', '🇫🇷', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', NOW(), NOW()),
(13, 3, 'AU', '澳大利亚', 'AUD', 'en-AU', '🇦🇺', 0, 0, '{"IncludeTax": true, "GSTRate": 0.10}', NOW(), NOW());