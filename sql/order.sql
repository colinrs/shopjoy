-- ============================================
-- 订单表 (orders)
-- ============================================

CREATE TABLE IF NOT EXISTS `orders` (
    `id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `order_no` VARCHAR(64) NOT NULL COMMENT '订单号',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待支付, 1-已支付, 2-待发货, 3-已发货, 4-已完成, 5-已取消, 6-退款中, 7-已退款',
    `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '商品总额(分)',
    `discount_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '优惠金额(分)',
    `freight_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '运费(分)',
    `pay_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '实付金额(分)',
    `currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 收货地址
    `address_name` VARCHAR(100) DEFAULT '' COMMENT '收货人',
    `address_phone` VARCHAR(20) DEFAULT '' COMMENT '收货电话',
    `address_province` VARCHAR(50) DEFAULT '' COMMENT '省份',
    `address_city` VARCHAR(50) DEFAULT '' COMMENT '城市',
    `address_district` VARCHAR(50) DEFAULT '' COMMENT '区县',
    `address_detail` TEXT COMMENT '详细地址',
    `address_zipcode` VARCHAR(20) DEFAULT '' COMMENT '邮编',

    -- 物流信息
    `tracking_no` VARCHAR(100) DEFAULT '' COMMENT '快递单号',
    `carrier` VARCHAR(50) DEFAULT '' COMMENT '快递公司',

    `remark` TEXT COMMENT '备注',
    `expire_at` BIGINT NOT NULL DEFAULT 0 COMMENT '过期时间',
    `paid_at` BIGINT DEFAULT NULL COMMENT '支付时间',
    `shipped_at` BIGINT DEFAULT NULL COMMENT '发货时间',
    `completed_at` BIGINT DEFAULT NULL COMMENT '完成时间',
    `cancelled_at` BIGINT DEFAULT NULL COMMENT '取消时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- ============================================
-- 订单商品表 (order_items)
-- ============================================

CREATE TABLE IF NOT EXISTS `order_items` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
    `product_name` VARCHAR(255) NOT NULL COMMENT '商品名称',
    `sku_name` VARCHAR(255) DEFAULT '' COMMENT 'SKU名称',
    `image` VARCHAR(500) DEFAULT '' COMMENT '图片',
    `price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价(分)',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
    `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '小计(分)',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_sku_id` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单商品表';

-- ============================================
-- 测试数据
-- ============================================

-- 订单数据 (Demo Shop)
INSERT INTO `orders` (`id`, `tenant_id`, `user_id`, `order_no`, `status`, `total_amount`, `discount_amount`, `freight_amount`, `pay_amount`, `currency`, `address_name`, `address_phone`, `address_province`, `address_city`, `address_district`, `address_detail`, `address_zipcode`, `tracking_no`, `carrier`, `remark`, `expire_at`, `paid_at`, `shipped_at`, `completed_at`, `cancelled_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
-- 已完成订单
('ORD202503010001', 1, 1, 'ORD202503010001', 4, 259800, 5000, 0, 254800, 'CNY', '小明', '13800000001', '北京市', '北京市', '朝阳区', '建国路88号院1号楼101', '100022', 'SF1234567890', '顺丰速运', '', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 29 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 28 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), 1, 1),

-- 已发货订单
('ORD202503100001', 1, 1, 'ORD202503100001', 3, 169800, 0, 0, 169800, 'CNY', '小明', '13800000001', '北京市', '北京市', '朝阳区', '建国路88号院1号楼101', '100022', 'YT9876543210', '圆通快递', '尽快发货', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), NULL, NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), 1, 1),

-- 已支付待发货订单
('ORD202503150001', 1, 2, 'ORD202503150001', 2, 9900, 0, 0, 9900, 'CNY', '小红', '13800000002', '上海市', '上海市', '浦东新区', '陆家嘴环路1000号', '200120', '', '', '', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY)), NULL, NULL, NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY)), 2, 2),

-- 待支付订单
('ORD202503200001', 1, 3, 'ORD202503200001', 0, 129900, 12990, 0, 116910, 'CNY', '小刚', '13800000003', '广东省', '深圳市', '南山区', '科技园南区A栋', '518000', '', '', '希望快点发货', UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 DAY)), NULL, NULL, NULL, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 3, 3),

-- 已取消订单
('ORD202503050001', 1, 1, 'ORD202503050001', 5, 29900, 0, 0, 29900, 'CNY', '小明', '13800000001', '北京市', '北京市', '朝阳区', '建国路88号院1号楼101', '100022', '', '', '', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), NULL, NULL, NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 19 DAY)), 1, 1),

-- Enterprise Corp 订单
('ORD202503120001', 3, 6, 'ORD202503120001', 4, 79900, 0, 1500, 81400, 'CNY', '约翰', '13800000006', '广东省', '广州市', '天河区', '珠江新城花城大道', '510600', 'EMS2025031201', 'EMS', '', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 7 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY)), NULL, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY)), 6, 6);

-- 订单商品数据
INSERT INTO `order_items` (`id`, `order_id`, `product_id`, `sku_id`, `product_name`, `sku_name`, `image`, `price`, `quantity`, `total_amount`, `created_at`) VALUES
-- ORD202503010001 商品
(1, 'ORD202503010001', 1, 1, 'Nike Air Max 270', '黑色 42码', 'https://cdn.example.com/p1-1.jpg', 129900, 1, 129900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY))),
(2, 'ORD202503010001', 3, 7, 'iPhone 15 手机壳', '黑色', 'https://cdn.example.com/p3-1.jpg', 9900, 3, 29700, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY))),
(3, 'ORD202503010001', 5, 0, '简约台灯', '', 'https://cdn.example.com/p5-1.jpg', 29900, 1, 29900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY))),

-- ORD202503100001 商品
(4, 'ORD202503100001', 2, 5, 'Adidas Ultraboost 22', '黑色 42码', 'https://cdn.example.com/p2-1.jpg', 159900, 1, 159900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY))),

-- ORD202503150001 商品
(5, 'ORD202503150001', 3, 8, 'iPhone 15 手机壳', '白色', 'https://cdn.example.com/p3-1.jpg', 9900, 1, 9900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY))),

-- ORD202503200001 商品
(6, 'ORD202503200001', 1, 1, 'Nike Air Max 270', '黑色 42码', 'https://cdn.example.com/p1-1.jpg', 129900, 1, 129900, UNIX_TIMESTAMP()),

-- ORD202503050001 商品 (已取消)
(7, 'ORD202503050001', 5, 0, '简约台灯', '', 'https://cdn.example.com/p5-1.jpg', 29900, 1, 29900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY))),

-- ORD202503120001 商品
(8, 'ORD202503120001', 4, 0, 'MacBook 充电器', '', 'https://cdn.example.com/p4-1.jpg', 79900, 1, 79900, UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY)));