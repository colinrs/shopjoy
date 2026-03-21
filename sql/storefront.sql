-- ============================================
-- 店铺表 (shops)
-- ============================================

CREATE TABLE IF NOT EXISTS `shops` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '店铺ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '店铺名称',
    `description` TEXT COMMENT '店铺描述',
    `logo` VARCHAR(500) DEFAULT '' COMMENT 'Logo URL',
    `banner` VARCHAR(500) DEFAULT '' COMMENT 'Banner URL',
    `contact_phone` VARCHAR(20) DEFAULT '' COMMENT '联系电话',
    `contact_email` VARCHAR(255) DEFAULT '' COMMENT '联系邮箱',
    `address` TEXT COMMENT '地址',
    `social_links` JSON COMMENT '社交链接',
    `seo_title` VARCHAR(255) DEFAULT '' COMMENT 'SEO标题',
    `seo_description` TEXT COMMENT 'SEO描述',
    `seo_keywords` VARCHAR(500) DEFAULT '' COMMENT 'SEO关键词',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='店铺表';

-- ============================================
-- 主题表 (themes)
-- ============================================

CREATE TABLE IF NOT EXISTS `themes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主题ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '主题名称',
    `code` VARCHAR(100) NOT NULL COMMENT '主题代码',
    `description` TEXT COMMENT '描述',
    `thumbnail` VARCHAR(500) DEFAULT '' COMMENT '缩略图',
    `config` JSON COMMENT '配置',
    `is_active` TINYINT NOT NULL DEFAULT 0 COMMENT '是否激活',
    `is_custom` TINYINT NOT NULL DEFAULT 0 COMMENT '是否自定义',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_code` (`code`),
    KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主题表';

-- ============================================
-- 页面表 (pages)
-- ============================================

CREATE TABLE IF NOT EXISTS `pages` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '页面ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '页面名称',
    `slug` VARCHAR(255) NOT NULL COMMENT 'URL别名',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-首页, 1-商品页, 2-集合页, 3-自定义页',
    `content` LONGTEXT COMMENT '内容',
    `seo_title` VARCHAR(255) DEFAULT '' COMMENT 'SEO标题',
    `seo_description` TEXT COMMENT 'SEO描述',
    `seo_keywords` VARCHAR(500) DEFAULT '' COMMENT 'SEO关键词',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_slug` (`tenant_id`, `slug`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='页面表';

-- ============================================
-- 导航表 (navigations)
-- ============================================

CREATE TABLE IF NOT EXISTS `navigations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '导航ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '导航名称',
    `position` VARCHAR(50) DEFAULT '' COMMENT '位置: header, footer, sidebar',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导航表';

-- ============================================
-- 导航项表 (nav_items)
-- ============================================

CREATE TABLE IF NOT EXISTS `nav_items` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '导航项ID',
    `nav_id` BIGINT NOT NULL COMMENT '导航ID',
    `parent_id` BIGINT NOT NULL DEFAULT 0 COMMENT '父ID',
    `name` VARCHAR(100) NOT NULL COMMENT '名称',
    `link` VARCHAR(500) DEFAULT '' COMMENT '链接',
    `type` VARCHAR(50) DEFAULT '' COMMENT '类型: page, category, product, custom',
    `target_id` BIGINT DEFAULT NULL COMMENT '目标ID',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    PRIMARY KEY (`id`),
    KEY `idx_nav_id` (`nav_id`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导航项表';

-- ============================================
-- 测试数据
-- ============================================

-- 店铺数据
INSERT INTO `shops` (`id`, `tenant_id`, `name`, `description`, `logo`, `banner`, `contact_phone`, `contact_email`, `address`, `social_links`, `seo_title`, `seo_description`, `seo_keywords`, `status`, `created_at`, `updated_at`) VALUES
(1, 1, 'Demo Shop 官方店铺', 'Demo Shop 是一家专注于高品质运动装备的在线商店', 'https://cdn.example.com/shop1-logo.png', 'https://cdn.example.com/shop1-banner.jpg', '400-888-8888', 'support@demoshop.com', '北京市朝阳区建国路88号', '{"weibo": "https://weibo.com/demoshop", "wechat": "demoshop_official"}', 'Demo Shop - 高品质运动装备', 'Demo Shop 提供Nike、Adidas等品牌运动鞋服，正品保障，全国包邮', '运动鞋,跑步鞋,Nike,Adidas', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 2, 'Test Store', 'Test Store - 精选好物，品质生活', 'https://cdn.example.com/shop2-logo.png', 'https://cdn.example.com/shop2-banner.jpg', '400-999-9999', 'support@teststore.com', '上海市浦东新区陆家嘴', '{"weibo": "https://weibo.com/teststore"}', 'Test Store - 精选好物', 'Test Store 致力于为您提供高品质的生活用品', '生活用品,家居,电子', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 3, 'Enterprise Corp 官方商城', 'Enterprise Corp 企业级电商平台', 'https://cdn.example.com/shop3-logo.png', 'https://cdn.example.com/shop3-banner.jpg', '400-666-6666', 'support@enterprisecorp.com', '广州市天河区珠江新城', '{"weibo": "https://weibo.com/enterprisecorp", "wechat": "ent_official", "twitter": "https://twitter.com/entcorp"}', 'Enterprise Corp - 企业级电商', 'Enterprise Corp 提供跨境购物服务，正品保障', '跨境电商,企业采购,进口商品', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 主题数据
INSERT INTO `themes` (`id`, `tenant_id`, `name`, `code`, `description`, `thumbnail`, `config`, `is_active`, `is_custom`, `created_at`, `updated_at`) VALUES
(1, 1, '简约白', 'minimal-white', '简洁大气的白色主题', 'https://cdn.example.com/theme1.png', '{"colors": {"primary": "#1890ff", "secondary": "#52c41a"}, "fonts": {"title": "PingFang SC", "body": "PingFang SC"}, "layout": "full-width"}', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 1, '深色模式', 'dark-mode', '深色主题，护眼模式', 'https://cdn.example.com/theme2.png', '{"colors": {"primary": "#177ddc", "secondary": "#49aa19"}, "fonts": {"title": "PingFang SC", "body": "PingFang SC"}, "layout": "full-width"}', 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 2, '清新绿', 'fresh-green', '清新的绿色主题', 'https://cdn.example.com/theme3.png', '{"colors": {"primary": "#52c41a", "secondary": "#faad14"}, "fonts": {"title": "Microsoft YaHei", "body": "Microsoft YaHei"}, "layout": "boxed"}', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 3, '企业蓝', 'enterprise-blue', '专业企业风格', 'https://cdn.example.com/theme4.png', '{"colors": {"primary": "#1d39c4", "secondary": "#fa8c16"}, "fonts": {"title": "Helvetica Neue", "body": "Helvetica Neue"}, "layout": "full-width"}', 1, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 页面数据 (Demo Shop)
INSERT INTO `pages` (`id`, `tenant_id`, `name`, `slug`, `type`, `content`, `seo_title`, `seo_description`, `seo_keywords`, `status`, `sort`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '首页', 'home', 0, '<section class="hero"><h1>欢迎来到Demo Shop</h1></section>', 'Demo Shop - 首页', 'Demo Shop 官方首页', '首页', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '关于我们', 'about-us', 3, '<h1>关于我们</h1><p>Demo Shop 成立于2020年...</p>', '关于我们 - Demo Shop', '了解Demo Shop的故事', '关于我们', 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '联系我们', 'contact', 3, '<h1>联系我们</h1><p>客服电话: 400-888-8888</p>', '联系我们 - Demo Shop', '联系Demo Shop客服', '联系方式', 1, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '隐私政策', 'privacy', 3, '<h1>隐私政策</h1><p>我们重视您的隐私...</p>', '隐私政策', 'Demo Shop隐私政策', '隐私政策', 1, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2);

-- 导航数据 (Demo Shop)
INSERT INTO `navigations` (`id`, `tenant_id`, `name`, `position`, `status`) VALUES
(1, 1, '顶部导航', 'header', 1),
(2, 1, '底部导航', 'footer', 1);

-- 导航项数据 (Demo Shop)
INSERT INTO `nav_items` (`id`, `nav_id`, `parent_id`, `name`, `link`, `type`, `target_id`, `sort`) VALUES
(1, 1, 0, '首页', '/', 'custom', NULL, 1),
(2, 1, 0, '商品分类', '/categories', 'custom', NULL, 2),
(3, 1, 2, '服装', '/category/clothing', 'category', 1, 1),
(4, 1, 2, '电子产品', '/category/electronics', 'category', 4, 2),
(5, 1, 2, '家居', '/category/home', 'category', 7, 3),
(6, 1, 0, '关于我们', '/about-us', 'page', 2, 3),
(7, 1, 0, '联系我们', '/contact', 'page', 3, 4),
(8, 2, 0, '关于我们', '/about-us', 'page', 2, 1),
(9, 2, 0, '联系我们', '/contact', 'page', 3, 2),
(10, 2, 0, '隐私政策', '/privacy', 'page', 4, 3);