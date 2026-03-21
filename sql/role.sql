-- ============================================
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
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_code` (`code`),
    KEY `idx_status` (`status`)
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
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_code` (`code`),
    KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- ============================================
-- 用户角色关联表 (user_roles)
-- ============================================

CREATE TABLE IF NOT EXISTS `user_roles` (
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    PRIMARY KEY (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 角色权限关联表 (role_permissions)
-- ============================================

CREATE TABLE IF NOT EXISTS `role_permissions` (
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT NOT NULL COMMENT '权限ID',
    PRIMARY KEY (`role_id`, `permission_id`),
    KEY `idx_permission_id` (`permission_id`)
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