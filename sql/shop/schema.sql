-- ============================================
-- 店铺设置表 (shop_settings)
-- ============================================

CREATE TABLE IF NOT EXISTS `shop_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '设置ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '店铺名称',
    `code` VARCHAR(50) NOT NULL COMMENT '店铺代码',
    `logo` VARCHAR(500) DEFAULT '' COMMENT 'Logo URL',
    `description` VARCHAR(500) DEFAULT '' COMMENT '店铺描述',
    `contact_name` VARCHAR(50) DEFAULT '' COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20) DEFAULT '' COMMENT '联系电话',
    `contact_email` VARCHAR(100) DEFAULT '' COMMENT '联系邮箱',
    `address` VARCHAR(255) DEFAULT '' COMMENT '地址',
    `domain` VARCHAR(100) DEFAULT '' COMMENT '店铺域名',
    `custom_domain` VARCHAR(100) DEFAULT '' COMMENT '自定义域名',
    `primary_color` VARCHAR(7) DEFAULT '#1890ff' COMMENT '主题色',
    `secondary_color` VARCHAR(7) DEFAULT '#52c41a' COMMENT '辅助色',
    `favicon` VARCHAR(255) DEFAULT '' COMMENT 'Favicon URL',
    `default_currency` VARCHAR(3) DEFAULT 'CNY' COMMENT '默认货币',
    `default_language` VARCHAR(10) DEFAULT 'zh-CN' COMMENT '默认语言',
    `timezone` VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '时区',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0=inactive, 1=active, 2=suspended',
    `plan` TINYINT NOT NULL DEFAULT 0 COMMENT '套餐: 0=basic, 1=standard, 2=premium',
    `expire_at` VARCHAR(50) DEFAULT '' COMMENT '过期时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='店铺设置表';

-- ============================================
-- 营业时间表 (shop_business_hours)
-- ============================================

CREATE TABLE IF NOT EXISTS `shop_business_hours` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '营业时间ID',
    `shop_id` BIGINT NOT NULL COMMENT '店铺ID(关联shop_settings)',
    `day_of_week` TINYINT NOT NULL COMMENT '星期: 0=Sunday, 1=Monday...6=Saturday',
    `open_time` VARCHAR(5) DEFAULT '09:00' COMMENT '开门时间',
    `close_time` VARCHAR(5) DEFAULT '18:00' COMMENT '关门时间',
    `is_closed` TINYINT NOT NULL DEFAULT 0 COMMENT '是否休息: 0=营业, 1=休息',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_shop_day` (`shop_id`, `day_of_week`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='营业时间表';

-- ============================================
-- 通知设置表 (shop_notification_settings)
-- ============================================

CREATE TABLE IF NOT EXISTS `shop_notification_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '通知设置ID',
    `shop_id` BIGINT NOT NULL COMMENT '店铺ID(关联shop_settings)',
    `order_created` TINYINT NOT NULL DEFAULT 1 COMMENT '订单创建通知',
    `order_paid` TINYINT NOT NULL DEFAULT 1 COMMENT '订单支付通知',
    `order_shipped` TINYINT NOT NULL DEFAULT 1 COMMENT '订单发货通知',
    `order_cancelled` TINYINT NOT NULL DEFAULT 1 COMMENT '订单取消通知',
    `low_stock_alert` TINYINT NOT NULL DEFAULT 1 COMMENT '低库存提醒',
    `low_stock_threshold` INT NOT NULL DEFAULT 10 COMMENT '低库存阈值',
    `refund_requested` TINYINT NOT NULL DEFAULT 1 COMMENT '退款申请通知',
    `new_review` TINYINT NOT NULL DEFAULT 1 COMMENT '新评价通知',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_shop_id` (`shop_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知设置表';

-- ============================================
-- 支付设置表 (shop_payment_settings)
-- ============================================

CREATE TABLE IF NOT EXISTS `shop_payment_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '支付设置ID',
    `shop_id` BIGINT NOT NULL COMMENT '店铺ID(关联shop_settings)',
    `stripe_enabled` TINYINT NOT NULL DEFAULT 0 COMMENT 'Stripe启用状态',
    `stripe_public_key` VARCHAR(255) DEFAULT '' COMMENT 'Stripe公钥',
    `stripe_secret_key` VARCHAR(255) DEFAULT '' COMMENT 'Stripe私钥(加密存储)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_shop_id` (`shop_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付设置表';

-- ============================================
-- 运费设置表 (shop_shipping_settings)
-- ============================================

CREATE TABLE IF NOT EXISTS `shop_shipping_settings` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '运费设置ID',
    `shop_id` BIGINT NOT NULL COMMENT '店铺ID(关联shop_settings)',
    `free_shipping_threshold` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '免运费门槛金额(元)',
    `default_shipping_fee` DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '默认运费金额(元)',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_shop_id` (`shop_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运费设置表';

-- ============================================
-- 测试数据
-- ============================================

-- 店铺设置数据
INSERT INTO `shop_settings` (`id`, `tenant_id`, `name`, `code`, `logo`, `description`, `contact_name`, `contact_phone`, `contact_email`, `address`, `domain`, `custom_domain`, `primary_color`, `secondary_color`, `favicon`, `default_currency`, `default_language`, `timezone`, `status`, `plan`, `expire_at`, `created_at`, `updated_at`) VALUES
(1, 1, 'Demo Shop', 'demo-shop', 'https://cdn.example.com/shop-logo.png', 'Demo Shop 是一家专注于高品质运动装备的在线商店', '张三', '400-888-8888', 'support@demoshop.com', '北京市朝阳区建国路88号', 'demo.myshopjoy.com', '', '#1890ff', '#52c41a', 'https://cdn.example.com/favicon.ico', 'CNY', 'zh-CN', 'Asia/Shanghai', 1, 1, '2027-12-31', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 2, 'Test Store', 'test-store', 'https://cdn.example.com/store-logo.png', 'Test Store - 精选好物', '李四', '400-999-9999', 'support@teststore.com', '上海市浦东新区陆家嘴', 'test.myshopjoy.com', '', '#1d39c4', '#fa8c16', '', 'CNY', 'zh-CN', 'Asia/Shanghai', 1, 0, '2026-06-30', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 营业时间数据
INSERT INTO `shop_business_hours` (`shop_id`, `day_of_week`, `open_time`, `close_time`, `is_closed`) VALUES
(1, 0, '10:00', '18:00', 1),  -- Sunday closed
(1, 1, '09:00', '22:00', 0),  -- Monday
(1, 2, '09:00', '22:00', 0),  -- Tuesday
(1, 3, '09:00', '22:00', 0),  -- Wednesday
(1, 4, '09:00', '22:00', 0),  -- Thursday
(1, 5, '09:00', '22:00', 0),  -- Friday
(1, 6, '09:00', '20:00', 0);  -- Saturday

-- 通知设置数据
INSERT INTO `shop_notification_settings` (`shop_id`, `order_created`, `order_paid`, `order_shipped`, `order_cancelled`, `low_stock_alert`, `low_stock_threshold`, `refund_requested`, `new_review`) VALUES
(1, 1, 1, 1, 1, 1, 10, 1, 1),
(2, 1, 1, 1, 1, 1, 5, 1, 0);

-- 支付设置数据
INSERT INTO `shop_payment_settings` (`shop_id`, `stripe_enabled`, `stripe_public_key`, `stripe_secret_key`) VALUES
(1, 1, 'pk_live_xxxxx', 'sk_live_xxxxx'),
(2, 0, '', '');

-- 运费设置数据
INSERT INTO `shop_shipping_settings` (`shop_id`, `free_shipping_threshold`, `default_shipping_fee`, `currency`) VALUES
(1, 99.00, 10.00, 'CNY'),
(2, 199.00, 15.00, 'CNY');