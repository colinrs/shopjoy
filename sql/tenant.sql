-- ============================================
-- 租户表 (tenants)
-- 多租户支持，支持不同套餐计划
-- ============================================

CREATE TABLE IF NOT EXISTS `tenants` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '租户名称',
    `code` VARCHAR(100) NOT NULL COMMENT '租户代码',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待审核, 1-正常, 2-暂停, 3-过期',
    `plan` TINYINT NOT NULL DEFAULT 0 COMMENT '套餐: 0-免费, 1-基础, 2-专业, 3-企业',
    `domain` VARCHAR(255) DEFAULT '' COMMENT '系统域名',
    `custom_domain` VARCHAR(255) DEFAULT '' COMMENT '自定义域名',
    `logo` VARCHAR(500) DEFAULT '' COMMENT 'Logo URL',
    `contact_name` VARCHAR(100) DEFAULT '' COMMENT '联系人',
    `contact_phone` VARCHAR(20) DEFAULT '' COMMENT '联系电话',
    `contact_email` VARCHAR(255) DEFAULT '' COMMENT '联系邮箱',
    `address` TEXT COMMENT '地址',
    `expire_at` BIGINT DEFAULT NULL COMMENT '过期时间(Unix时间戳)',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_domain` (`domain`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- ============================================
-- 测试数据
-- ============================================

INSERT INTO `tenants` (`id`, `name`, `code`, `status`, `plan`, `domain`, `custom_domain`, `logo`, `contact_name`, `contact_phone`, `contact_email`, `address`, `expire_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 'Demo Shop', 'demo', 1, 2, 'demo.shopjoy.com', 'www.demoshop.com', 'https://cdn.example.com/logo.png', '张三', '13800138000', 'admin@demoshop.com', '北京市朝阳区', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 YEAR)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),
(2, 'Test Store', 'test', 1, 1, 'test.shopjoy.com', '', 'https://cdn.example.com/logo2.png', '李四', '13900139000', 'admin@teststore.com', '上海市浦东新区', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 6 MONTH)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),
(3, 'Enterprise Corp', 'enterprise', 1, 3, 'enterprise.shopjoy.com', 'shop.enterprisecorp.com', 'https://cdn.example.com/logo3.png', '王五', '13700137000', 'admin@enterprisecorp.com', '广州市天河区', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 2 YEAR)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1);