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
    `expire_at` TIMESTAMP NULL COMMENT '过期时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
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
(3, 'Enterprise Corp', 'enterprise', 1, 3, 'enterprise.shopjoy.com', 'shop.enterprisecorp.com', 'https://cdn.example.com/logo3.png', '王五', '13700137000', 'admin@enterprisecorp.com', '广州市天河区', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 2 YEAR)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1);-- ============================================
-- 用户表 (users)
-- C端用户/顾客
-- ============================================

CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `name` VARCHAR(100) NOT NULL COMMENT '昵称',
    `avatar` VARCHAR(500) DEFAULT '' COMMENT '头像',
    `gender` TINYINT NOT NULL DEFAULT 0 COMMENT '性别: 0-未知, 1-男, 2-女, 3-其他',
    `birthday` TIMESTAMP NULL COMMENT '生日',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未激活, 1-正常, 2-暂停, 3-已删除',
    `last_login` TIMESTAMP NULL COMMENT '最后登录时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_email` (`email`),
    KEY `idx_phone` (`phone`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 测试数据
-- 密码均为: user123456
-- ============================================

INSERT INTO `users` (`id`, `tenant_id`, `email`, `phone`, `password`, `name`, `avatar`, `gender`, `birthday`, `status`, `last_login`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- Demo Shop 用户
(1, 1, 'user1@example.com', '13800000001', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '小明', 'https://cdn.example.com/u1.png', 1, UNIX_TIMESTAMP('1990-05-15'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(2, 1, 'user2@example.com', '13800000002', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '小红', 'https://cdn.example.com/u2.png', 2, UNIX_TIMESTAMP('1995-08-20'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(3, 1, 'user3@example.com', '13800000003', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '小刚', 'https://cdn.example.com/u3.png', 1, UNIX_TIMESTAMP('1988-03-10'), 1, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),

-- Test Store 用户
(4, 2, 'user4@example.com', '13800000004', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '大卫', 'https://cdn.example.com/u4.png', 1, UNIX_TIMESTAMP('1992-11-25'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(5, 2, 'user5@example.com', '13800000005', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '艾米', 'https://cdn.example.com/u5.png', 2, UNIX_TIMESTAMP('1998-01-05'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),

-- Enterprise Corp 用户
(6, 3, 'user6@example.com', '13800000006', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '约翰', 'https://cdn.example.com/u6.png', 1, UNIX_TIMESTAMP('1985-07-12'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(7, 3, 'user7@example.com', '13800000007', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '玛丽', 'https://cdn.example.com/u7.png', 2, UNIX_TIMESTAMP('1993-04-30'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(8, 3, 'user8@example.com', '13800000008', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '汤姆', 'https://cdn.example.com/u8.png', 1, UNIX_TIMESTAMP('1990-09-18'), 2, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);-- ============================================
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
    `last_login_at` TIMESTAMP NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT '' COMMENT '最后登录IP',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    UNIQUE KEY `uk_mobile` (`mobile`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
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
(7, 3, 'ent_sub2', 'sub2@enterprisecorp.com', '13600000007', '$2a$10$Wqlk81.6vgogQadFe2le1.WP6KKG2dueb0n11pbzzNb5fPUZhHgyy', '企业子账号2', '', 3, 2, NULL, '', NOW(), NOW());-- ============================================
-- 角色表 (roles)
-- ============================================

CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(100) NOT NULL COMMENT '角色代码',
    `description` TEXT COMMENT '描述',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `is_system` TINYINT NOT NULL DEFAULT 0 COMMENT '是否系统角色',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_code` (`code`),
    KEY `idx_status` (`status`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================
-- 权限表 (permissions)
-- ============================================

CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `code` VARCHAR(100) NOT NULL COMMENT '权限代码',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-菜单, 1-按钮, 2-API',
    `parent_id` BIGINT NOT NULL DEFAULT 0 COMMENT '父ID',
    `path` VARCHAR(255) DEFAULT '' COMMENT '路径',
    `icon` VARCHAR(100) DEFAULT '' COMMENT '图标',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_code` (`code`),
    KEY `idx_type` (`type`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- ============================================
-- 用户角色关联表 (user_roles)
-- ============================================

CREATE TABLE IF NOT EXISTS `user_roles` (
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 角色权限关联表 (role_permissions)
-- ============================================

CREATE TABLE IF NOT EXISTS `role_permissions` (
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT NOT NULL COMMENT '权限ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`role_id`, `permission_id`),
    KEY `idx_permission_id` (`permission_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ============================================
-- 测试数据
-- ============================================

-- 角色数据
INSERT INTO `roles` (`id`, `tenant_id`, `name`, `code`, `description`, `status`, `is_system`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '管理员', 'admin', 'Demo Shop 管理员', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 1, 1),
(2, 1, '运营', 'operator', 'Demo Shop 运营人员', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '客服', 'service', 'Demo Shop 客服人员', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 2, '管理员', 'admin', 'Test Store 管理员', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 4, 4),
(5, 3, '管理员', 'admin', 'Enterprise 管理员', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(6, 3, '运营主管', 'operator_lead', 'Enterprise 运营主管', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(7, 3, '普通运营', 'operator', 'Enterprise 普通运营', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 权限数据
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `icon`, `sort`) VALUES
-- 菜单权限
(1, '仪表盘', 'dashboard', 0, 0, '/dashboard', 'dashboard', 1),
(2, '商品管理', 'product', 0, 0, '/product', 'shopping', 2),
(3, '商品列表', 'product:list', 0, 2, '/product/list', '', 1),
(4, '分类管理', 'category', 0, 2, '/product/category', '', 2),
(5, '品牌管理', 'brand', 0, 2, '/product/brand', '', 3),
(6, '订单管理', 'order', 0, 0, '/order', 'file-text', 3),
(7, '订单列表', 'order:list', 0, 6, '/order/list', '', 1),
(8, '退款管理', 'refund', 0, 6, '/order/refund', '', 2),
(9, '用户管理', 'user', 0, 0, '/user', 'user', 4),
(10, '用户列表', 'user:list', 0, 9, '/user/list', '', 1),
(11, '营销管理', 'marketing', 0, 0, '/marketing', 'gift', 5),
(12, '优惠券', 'coupon', 0, 11, '/marketing/coupon', '', 1),
(13, '促销活动', 'promotion', 0, 11, '/marketing/promotion', '', 2),
(14, '系统设置', 'system', 0, 0, '/system', 'setting', 6),
(15, '角色管理', 'role', 0, 14, '/system/role', '', 1),
-- 按钮权限
(16, '新增商品', 'product:create', 1, 3, '', '', 1),
(17, '编辑商品', 'product:update', 1, 3, '', '', 2),
(18, '删除商品', 'product:delete', 1, 3, '', '', 3),
(19, '上架商品', 'product:onsale', 1, 3, '', '', 4),
(20, '下架商品', 'product:offsale', 1, 3, '', '', 5),
(21, '查看订单', 'order:view', 1, 7, '', '', 1),
(22, '发货', 'order:ship', 1, 7, '', '', 2),
(23, '取消订单', 'order:cancel', 1, 7, '', '', 3);

-- 用户角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(3, 2),
(3, 3),
(6, 6),
(7, 7);

-- 角色权限关联 (Demo Shop 管理员)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16), (1, 17), (1, 18), (1, 19), (1, 20),
(1, 21), (1, 22), (1, 23);

-- 角色权限关联 (Demo Shop 运营)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(2, 1), (2, 2), (2, 3), (2, 6), (2, 7), (2, 16), (2, 17), (2, 21), (2, 22);

-- 角色权限关联 (Demo Shop 客服)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(3, 1), (3, 6), (3, 7), (3, 21);

-- ============================================
-- 用户地址表 (user_addresses)
-- ============================================

CREATE TABLE IF NOT EXISTS `user_addresses` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `user_id` BIGINT NOT NULL,
    `name` VARCHAR(100) NOT NULL COMMENT '收货人姓名',
    `phone` VARCHAR(20) NOT NULL COMMENT '收货人电话',
    `country` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '国家代码',
    `province` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '省份/州',
    `city` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '城市',
    `district` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '区/县',
    `address` VARCHAR(255) NOT NULL COMMENT '详细地址',
    `postal_code` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '邮编',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否默认地址',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收货地址表';
