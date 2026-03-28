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
    `current_theme_id` BIGINT DEFAULT NULL COMMENT 'Current active theme ID',
    `theme_config` TEXT COMMENT 'JSON theme customization config',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
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
    `preview_image` VARCHAR(500) DEFAULT '' COMMENT 'Preview image URL',
    `config` JSON COMMENT '配置',
    `config_schema` TEXT COMMENT 'JSON schema for configurable fields',
    `default_config` TEXT COMMENT 'JSON default configuration',
    `is_active` TINYINT NOT NULL DEFAULT 0 COMMENT '是否激活',
    `is_custom` TINYINT NOT NULL DEFAULT 0 COMMENT '是否自定义',
    `is_preset` TINYINT NOT NULL DEFAULT 1 COMMENT '1=preset theme, 0=custom theme',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_code` (`code`),
    KEY `idx_is_active` (`is_active`),
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
    `is_published` TINYINT NOT NULL DEFAULT 0 COMMENT 'Whether page is published',
    `published_at` TIMESTAMP NULL COMMENT 'Published timestamp',
    `version` INT NOT NULL DEFAULT 1 COMMENT 'Current version number',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_slug` (`tenant_id`, `slug`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
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
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_position` (`position`),
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
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    KEY `idx_nav_id` (`nav_id`),
    KEY `idx_parent_id` (`parent_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导航项表';

-- ============================================
-- 页面装修表 (decorations)
-- ============================================

CREATE TABLE IF NOT EXISTS `decorations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'Decoration block ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_id` BIGINT NOT NULL COMMENT 'Page ID',
    `block_type` VARCHAR(50) NOT NULL COMMENT 'Block type: banner, product_grid, rich_text, image_carousel, featured_products, categories, divider',
    `block_config` TEXT NOT NULL COMMENT 'JSON block configuration',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order within page',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Whether block is active',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created timestamp',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated timestamp',
    PRIMARY KEY (`id`),
    INDEX `idx_page_sort` (`page_id`, `sort_order`),
    INDEX `idx_tenant_page` (`tenant_id`, `page_id`),
    INDEX `idx_block_type` (`block_type`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Page decoration blocks';

-- ============================================
-- 页面版本表 (page_versions)
-- ============================================

CREATE TABLE IF NOT EXISTS `page_versions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'Version ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_id` BIGINT NOT NULL COMMENT 'Page ID',
    `version` INT NOT NULL COMMENT 'Version number',
    `blocks` TEXT NOT NULL COMMENT 'JSON snapshot of decoration blocks',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT 'User who created this version',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created timestamp',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated timestamp',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_tenant_page_ver` (`tenant_id`, `page_id`, `version`),
    INDEX `idx_page_version` (`page_id`, `version`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Page version history';

-- ============================================
-- SEO配置表 (seo_configs)
-- ============================================

CREATE TABLE IF NOT EXISTS `seo_configs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'SEO config ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `page_type` VARCHAR(30) NOT NULL COMMENT 'Page type: global, home, category, product, custom',
    `page_id` BIGINT DEFAULT NULL COMMENT 'Page ID for custom pages (NULL for global/page type defaults)',
    `title` VARCHAR(200) NOT NULL DEFAULT '' COMMENT 'SEO title',
    `description` TEXT NOT NULL COMMENT 'SEO description',
    `keywords` VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'SEO keywords',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created timestamp',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated timestamp',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_tenant_page_type` (`tenant_id`, `page_type`, `page_id`),
    INDEX `idx_page_type` (`page_type`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SEO configurations';

-- ============================================
-- 主题审计日志表 (theme_audit_logs)
-- ============================================

CREATE TABLE IF NOT EXISTS `theme_audit_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'Audit log ID',
    `tenant_id` BIGINT NOT NULL COMMENT 'Tenant ID',
    `action` VARCHAR(30) NOT NULL COMMENT 'Action: switch_theme, update_config',
    `theme_id` BIGINT NOT NULL COMMENT 'Theme ID',
    `theme_name` VARCHAR(100) NOT NULL COMMENT 'Theme name at the time of action',
    `theme_code` VARCHAR(50) NOT NULL COMMENT 'Theme code',
    `old_config` TEXT COMMENT 'Previous configuration (JSON)',
    `new_config` TEXT COMMENT 'New configuration (JSON)',
    `user_id` BIGINT NOT NULL COMMENT 'User who performed the action',
    `user_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'User name',
    `ip_address` VARCHAR(45) DEFAULT '' COMMENT 'IP address',
    `user_agent` VARCHAR(500) DEFAULT '' COMMENT 'User agent',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created timestamp',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated timestamp',
    `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_action` (`tenant_id`, `action`),
    INDEX `idx_tenant_theme` (`tenant_id`, `theme_id`),
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_created_at` (`created_at`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Theme change audit logs';

-- ============================================
-- 测试数据
-- ============================================

-- 店铺数据
INSERT INTO `shops` (`id`, `tenant_id`, `name`, `description`, `logo`, `banner`, `contact_phone`, `contact_email`, `address`, `social_links`, `seo_title`, `seo_description`, `seo_keywords`, `current_theme_id`, `status`, `created_at`, `updated_at`) VALUES
(1, 1, 'Demo Shop 官方店铺', 'Demo Shop 是一家专注于高品质运动装备的在线商店', 'https://cdn.example.com/shop1-logo.png', 'https://cdn.example.com/shop1-banner.jpg', '400-888-8888', 'support@demoshop.com', '北京市朝阳区建国路88号', '{"weibo": "https://weibo.com/demoshop", "wechat": "demoshop_official"}', 'Demo Shop - 高品质运动装备', 'Demo Shop 提供Nike、Adidas等品牌运动鞋服，正品保障，全国包邮', '运动鞋,跑步鞋,Nike,Adidas', 1001, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 2, 'Test Store', 'Test Store - 精选好物，品质生活', 'https://cdn.example.com/shop2-logo.png', 'https://cdn.example.com/shop2-banner.jpg', '400-999-9999', 'support@teststore.com', '上海市浦东新区陆家嘴', '{"weibo": "https://weibo.com/teststore"}', 'Test Store - 精选好物', 'Test Store 致力于为您提供高品质的生活用品', '生活用品,家居,电子', 1001, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 3, 'Enterprise Corp 官方商城', 'Enterprise Corp 企业级电商平台', 'https://cdn.example.com/shop3-logo.png', 'https://cdn.example.com/shop3-banner.jpg', '400-666-6666', 'support@enterprisecorp.com', '广州市天河区珠江新城', '{"weibo": "https://weibo.com/enterprisecorp", "wechat": "ent_official", "twitter": "https://twitter.com/entcorp"}', 'Enterprise Corp - 企业级电商', 'Enterprise Corp 提供跨境购物服务，正品保障', '跨境电商,企业采购,进口商品', 1001, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 主题数据
INSERT INTO `themes` (`id`, `tenant_id`, `name`, `code`, `description`, `thumbnail`, `preview_image`, `config`, `config_schema`, `default_config`, `is_active`, `is_custom`, `is_preset`, `created_at`, `updated_at`) VALUES
(1, 1, '简约白', 'minimal-white', '简洁大气的白色主题', 'https://cdn.example.com/theme1.png', '', '{"colors": {"primary": "#1890ff", "secondary": "#52c41a"}, "fonts": {"title": "PingFang SC", "body": "PingFang SC"}, "layout": "full-width"}', NULL, NULL, 1, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 1, '深色模式', 'dark-mode', '深色主题，护眼模式', 'https://cdn.example.com/theme2.png', '', '{"colors": {"primary": "#177ddc", "secondary": "#49aa19"}, "fonts": {"title": "PingFang SC", "body": "PingFang SC"}, "layout": "full-width"}', NULL, NULL, 0, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 2, '清新绿', 'fresh-green', '清新的绿色主题', 'https://cdn.example.com/theme3.png', '', '{"colors": {"primary": "#52c41a", "secondary": "#faad14"}, "fonts": {"title": "Microsoft YaHei", "body": "Microsoft YaHei"}, "layout": "boxed"}', NULL, NULL, 1, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 3, '企业蓝', 'enterprise-blue', '专业企业风格', 'https://cdn.example.com/theme4.png', '', '{"colors": {"primary": "#1d39c4", "secondary": "#fa8c16"}, "fonts": {"title": "Helvetica Neue", "body": "Helvetica Neue"}, "layout": "full-width"}', NULL, NULL, 1, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- Preset themes
(1001, 0, 'Classic', 'classic', 'A timeless design with clean lines and professional appearance. Perfect for businesses looking for a traditional e-commerce feel.', 'https://cdn.shopjoy.com/themes/classic-thumb.png', 'https://cdn.shopjoy.com/themes/classic-preview.png', '{"colors":{"primary":"#3B82F6","secondary":"#1E40AF"},"fonts":{"heading":"Inter","body":"Inter"},"layout":"standard"}', '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#3B82F6"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#1E40AF"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"inter","label":"Inter"},{"value":"roboto","label":"Roboto"},{"value":"opensans","label":"Open Sans"}],"default":"inter"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"square","label":"Square"}],"default":"rounded"}]', '{"primary_color":"#3B82F6","secondary_color":"#1E40AF","font_family":"inter","button_style":"rounded"}', 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1002, 0, 'Modern', 'modern', 'A sleek and contemporary design with bold colors and dynamic layouts. Ideal for fashion and lifestyle brands.', 'https://cdn.shopjoy.com/themes/modern-thumb.png', 'https://cdn.shopjoy.com/themes/modern-preview.png', '{"colors":{"primary":"#10B981","secondary":"#059669"},"fonts":{"heading":"Poppins","body":"Inter"},"layout":"modern"}', '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#10B981"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#059669"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"poppins","label":"Poppins"},{"value":"montserrat","label":"Montserrat"}],"default":"poppins"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"pill","label":"Pill"},{"value":"rounded","label":"Rounded"}],"default":"pill"}]', '{"primary_color":"#10B981","secondary_color":"#059669","font_family":"poppins","button_style":"pill"}', 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1003, 0, 'Minimal', 'minimal', 'A clean and minimalist design that puts your products front and center. Great for luxury and high-end brands.', 'https://cdn.shopjoy.com/themes/minimal-thumb.png', 'https://cdn.shopjoy.com/themes/minimal-preview.png', '{"colors":{"primary":"#000000","secondary":"#6B7280"},"fonts":{"heading":"Helvetica Neue","body":"Helvetica Neue"},"layout":"minimal"}', '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#000000"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#6B7280"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"helvetica","label":"Helvetica Neue"},{"value":"arial","label":"Arial"}],"default":"helvetica"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"underline","label":"Underline"},{"value":"solid","label":"Solid"}],"default":"underline"}]', '{"primary_color":"#000000","secondary_color":"#6B7280","font_family":"helvetica","button_style":"underline"}', 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1004, 0, 'Bold', 'bold', 'A vibrant and energetic design with strong visual elements. Perfect for creative brands and youth-oriented products.', 'https://cdn.shopjoy.com/themes/bold-thumb.png', 'https://cdn.shopjoy.com/themes/bold-preview.png', '{"colors":{"primary":"#8B5CF6","secondary":"#6D28D9"},"fonts":{"heading":"DM Sans","body":"DM Sans"},"layout":"bold"}', '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#8B5CF6"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#6D28D9"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"dmsans","label":"DM Sans"},{"value":"nunito","label":"Nunito"}],"default":"dmsans"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"pill","label":"Pill"}],"default":"rounded"}]', '{"primary_color":"#8B5CF6","secondary_color":"#6D28D9","font_family":"dmsans","button_style":"rounded"}', 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

(1005, 0, 'Nature', 'nature', 'An organic and eco-friendly design with natural colors and earthy tones. Ideal for sustainable and wellness brands.', 'https://cdn.shopjoy.com/themes/nature-thumb.png', 'https://cdn.shopjoy.com/themes/nature-preview.png', '{"colors":{"primary":"#059669","secondary":"#047857"},"fonts":{"heading":"Merriweather","body":"Open Sans"},"layout":"nature"}', '[{"key":"primary_color","label":"Primary Color","type":"color","default":"#059669"},{"key":"secondary_color","label":"Secondary Color","type":"color","default":"#047857"},{"key":"font_family","label":"Font Family","type":"select","options":[{"value":"merriweather","label":"Merriweather"},{"value":"lora","label":"Lora"}],"default":"merriweather"},{"key":"button_style","label":"Button Style","type":"select","options":[{"value":"rounded","label":"Rounded"},{"value":"leaf","label":"Leaf"}],"default":"rounded"}]', '{"primary_color":"#059669","secondary_color":"#047857","font_family":"merriweather","button_style":"rounded"}', 0, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 页面数据 (Demo Shop)
INSERT INTO `pages` (`id`, `tenant_id`, `name`, `slug`, `type`, `content`, `seo_title`, `seo_description`, `seo_keywords`, `status`, `sort`, `is_published`, `version`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '首页', 'home', 0, '<section class="hero"><h1>欢迎来到Demo Shop</h1></section>', 'Demo Shop - 首页', 'Demo Shop 官方首页', '首页', 1, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '关于我们', 'about-us', 3, '<h1>关于我们</h1><p>Demo Shop 成立于2020年...</p>', '关于我们 - Demo Shop', '了解Demo Shop的故事', '关于我们', 1, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '联系我们', 'contact', 3, '<h1>联系我们</h1><p>客服电话: 400-888-8888</p>', '联系我们 - Demo Shop', '联系Demo Shop客服', '联系方式', 1, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '隐私政策', 'privacy', 3, '<h1>隐私政策</h1><p>我们重视您的隐私...</p>', '隐私政策', 'Demo Shop隐私政策', '隐私政策', 1, 4, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2);

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

-- SEO配置数据 (Demo Shop)
INSERT INTO `seo_configs` (`tenant_id`, `page_type`, `page_id`, `title`, `description`, `keywords`, `created_at`, `updated_at`) VALUES
(1, 'global', NULL, 'Demo Shop - 高品质运动装备', 'Demo Shop 提供Nike、Adidas等品牌运动鞋服，正品保障，全国包邮', '运动鞋,跑步鞋,Nike,Adidas', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'home', NULL, 'Demo Shop - 首页', 'Demo Shop 官方首页', '首页', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'product', NULL, '{{product_name}} - Demo Shop', '{{product_description}}', '', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'category', NULL, '{{category_name}} - Demo Shop', 'Browse our {{category_name}} collection', '', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
