-- ============================================
-- 管理员用户表 (admin_users)
-- 后台管理系统用户，支持超管、租户管理员、子账号
-- ============================================

CREATE TABLE IF NOT EXISTS `admin_users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID，0表示平台超管',
    `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
    `email` VARCHAR(128) NOT NULL COMMENT '邮箱',
    `mobile` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `password` VARCHAR(255) NOT NULL COMMENT '密码(bcrypt)',
    `real_name` VARCHAR(32) DEFAULT '' COMMENT '真实姓名',
    `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型: 1-平台超管, 2-租户管理员, 3-租户子账号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-正常, 2-禁用, 3-已删除',
    `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT '' COMMENT '最后登录IP',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    UNIQUE KEY `uk_mobile` (`mobile`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';

-- ============================================
-- 测试数据
-- 密码均为: password123,$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy
-- ============================================

INSERT INTO `admin_users` (`id`, `tenant_id`, `username`, `email`, `mobile`, `password`, `real_name`, `avatar`, `type`, `status`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
-- 平台超级管理员
(1, 0, 'superadmin', 'superadmin@shopjoy.com', '13600000001', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '超级管理员', 'https://cdn.example.com/avatar1.png', 1, 1, NOW(), '127.0.0.1', NOW(), NOW()),

-- Demo Shop 管理员
(2, 1, 'demo_admin', 'admin@demoshop.com', '13600000002', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', 'Demo管理员', 'https://cdn.example.com/avatar2.png', 2, 1, NOW(), '127.0.0.1', NOW(), NOW()),
(3, 1, 'demo_sub1', 'sub1@demoshop.com', '13600000003', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', 'Demo子账号1', 'https://cdn.example.com/avatar3.png', 3, 1, NULL, '', NOW(), NOW()),

-- Test Store 管理员
(4, 2, 'test_admin', 'admin@teststore.com', '13600000004', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', 'Test管理员', 'https://cdn.example.com/avatar4.png', 2, 1, NOW(), '192.168.1.100', NOW(), NOW()),

-- Enterprise Corp 管理员
(5, 3, 'ent_admin', 'admin@enterprisecorp.com', '13600000005', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '企业管理员', 'https://cdn.example.com/avatar5.png', 2, 1, NOW(), '10.0.0.1', NOW(), NOW()),
(6, 3, 'ent_sub1', 'sub1@enterprisecorp.com', '13600000006', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '企业子账号1', '', 3, 1, NULL, '', NOW(), NOW()),
(7, 3, 'ent_sub2', 'sub2@enterprisecorp.com', '13600000007', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '企业子账号2', '', 3, 2, NULL, '', NOW(), NOW());