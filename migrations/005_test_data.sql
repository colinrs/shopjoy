-- ShopJoy 测试数据
-- 用于多市场商品管理系统测试

-- ============================================
-- 商品测试数据 (带合规字段)
-- ============================================

INSERT INTO `products` (`id`, `tenant_id`, `sku`, `name`, `description`, `price`, `cost_price`, `currency`, `stock`, `status`, `category_id`, `brand`, `tags`, `images`, `is_matrix_product`, `hs_code`, `coo`, `weight`, `weight_unit`, `length`, `width`, `height`, `dangerous_goods`, `created_at`, `updated_at`) VALUES
-- 商品1: 无线耳机 (多规格商品)
(1, 1, 'WEP-001', '无线耳机 Pro', '高品质无线蓝牙耳机，支持主动降噪，续航30小时', 29900, 15000, 'CNY', 500, 1, 1, 'SoundMax', '["电子", "耳机", "蓝牙"]', '["https://example.com/images/wep-001-1.jpg", "https://example.com/images/wep-001-2.jpg"]', 1, '8518.30', 'CN', 250.00, 'g', 8.50, 6.20, 3.00, '["battery"]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品2: 智能手表 (多规格商品)
(2, 1, 'SW-007', '智能手表 S7', '智能运动手表，支持心率监测、GPS定位、防水50米', 199900, 80000, 'CNY', 300, 1, 1, 'TechFit', '["电子", "手表", "运动"]', '["https://example.com/images/sw-007-1.jpg", "https://example.com/images/sw-007-2.jpg"]', 1, '9102.11', 'CN', 85.00, 'g', 4.50, 4.50, 1.20, '["battery", "magnet"]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品3: 真皮钱包 (单规格商品)
(3, 1, 'WL-101', '真皮商务钱包', '意大利进口小牛皮，简约商务风格，多卡位设计', 59900, 25000, 'CNY', 200, 1, 2, 'LuxeLeather', '["配饰", "钱包", "真皮"]', '["https://example.com/images/wl-101-1.jpg"]', 0, '4202.31', 'IT', 120.00, 'g', 11.00, 9.50, 2.00, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品4: 运动水壶 (单规格商品)
(4, 1, 'BK-202', '保温运动水壶 750ml', '316不锈钢内胆，保温12小时，防漏设计', 15900, 6000, 'CNY', 800, 1, 3, 'HydroMax', '["运动", "水壶", "保温"]', '["https://example.com/images/bk-202-1.jpg", "https://example.com/images/bk-202-2.jpg"]', 0, '9617.00', 'CN', 350.00, 'g', 7.50, 7.50, 24.00, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品5: 有机茶叶礼盒 (合规要求高)
(5, 1, 'TEA-050', '有机龙井茶礼盒 250g', '西湖产区有机龙井，明前采摘，礼盒装', 38800, 15000, 'CNY', 150, 1, 4, 'TeaMaster', '["食品", "茶叶", "礼盒", "有机"]', '["https://example.com/images/tea-050-1.jpg"]', 0, '0902.10', 'CN', 320.00, 'g', 20.00, 15.00, 8.00, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品6: 儿童玩具车 (含电池，需合规)
(6, 1, 'TOY-301', '电动遥控越野车', '1:16比例遥控越野车，配备可充电电池，适合6岁以上儿童', 25900, 10000, 'CNY', 400, 1, 5, 'PlayJoy', '["玩具", "遥控", "儿童"]', '["https://example.com/images/toy-301-1.jpg"]', 0, '9503.00', 'CN', 580.00, 'g', 25.00, 15.00, 12.00, '["battery", "magnet"]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品7: 化妆品套装 (液体，有运输限制)
(7, 1, 'BEA-888', '护肤精华套装', '包含精华液30ml + 面霜50ml + 眼霜15ml', 45900, 18000, 'CNY', 250, 1, 6, 'GlowSkin', '["美妆", "护肤", "套装"]', '["https://example.com/images/bea-888-1.jpg"]', 0, '3304.99', 'KR', 200.00, 'g', 15.00, 12.00, 8.00, '["liquid"]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品8: 瑜伽垫 (大件商品)
(8, 1, 'YM-100', '加厚防滑瑜伽垫 6mm', 'TPE环保材质，防滑纹理设计，附带收纳带', 12900, 4000, 'CNY', 600, 1, 3, 'YogaPro', '["运动", "瑜伽", "健身"]', '["https://example.com/images/ym-100-1.jpg"]', 0, '9506.91', 'CN', 1200.00, 'g', 183.00, 61.00, 0.60, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品9: 草稿状态商品 (未上架)
(9, 1, 'NEW-001', '新品智能音箱', 'AI语音助手，支持智能家居控制', 49900, 20000, 'CNY', 0, 0, 1, 'TechFit', '["电子", "智能家居"]', '["https://example.com/images/new-001-1.jpg"]', 0, '8518.29', 'CN', 450.00, 'g', 10.00, 10.00, 8.00, '["battery"]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 商品10: 缺少合规信息的商品
(10, 1, 'TEMP-999', '临时测试商品', '这是一个测试商品，缺少合规信息', 9900, 3000, 'CNY', 100, 1, 1, 'TestBrand', '["测试"]', '[]', 0, NULL, NULL, NULL, 'g', NULL, NULL, NULL, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());


-- ============================================
-- 商品市场关联测试数据
-- ============================================

-- 商品1: 无线耳机 - 全市场销售
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 1, NULL, 1, 1, 29.99, 39.99, 50, NOW(), NOW(), NOW()), -- US: $29.99
(1, 1, NULL, 2, 1, 24.99, 32.99, 50, NOW(), NOW(), NOW()), -- UK: £24.99
(1, 1, NULL, 3, 1, 29.99, 39.99, 50, NOW(), NOW(), NOW()), -- DE: €29.99
(1, 1, NULL, 4, 1, 29.99, 39.99, 50, NOW(), NOW(), NOW()), -- FR: €29.99
(1, 1, NULL, 5, 1, 44.99, 54.99, 50, NOW(), NOW(), NOW()); -- AU: A$44.99

-- 商品2: 智能手表 - 全市场销售
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 2, NULL, 1, 1, 199.00, 249.00, 30, NOW(), NOW(), NOW()), -- US: $199.00
(1, 2, NULL, 2, 1, 169.00, 209.00, 30, NOW(), NOW(), NOW()), -- UK: £169.00
(1, 2, NULL, 3, 1, 199.00, 249.00, 30, NOW(), NOW(), NOW()), -- DE: €199.00
(1, 2, NULL, 4, 1, 199.00, 249.00, 30, NOW(), NOW(), NOW()), -- FR: €199.00
(1, 2, NULL, 5, 1, 299.00, 349.00, 30, NOW(), NOW(), NOW()); -- AU: A$299.00

-- 商品3: 真皮钱包 - 仅US/UK市场
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 3, NULL, 1, 1, 79.00, NULL, 20, NOW(), NOW(), NOW()), -- US: $79.00
(1, 3, NULL, 2, 1, 65.00, NULL, 20, NOW(), NOW(), NOW()); -- UK: £65.00

-- 商品4: 运动水壶 - US/UK/DE市场
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 4, NULL, 1, 1, 19.99, 24.99, 100, NOW(), NOW(), NOW()), -- US: $19.99
(1, 4, NULL, 2, 1, 16.99, 21.99, 100, NOW(), NOW(), NOW()), -- UK: £16.99
(1, 4, NULL, 3, 1, 19.99, 24.99, 100, NOW(), NOW(), NOW()); -- DE: €19.99

-- 商品5: 有机茶叶 - 仅中国出口市场 (US)
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 5, NULL, 1, 1, 58.00, NULL, 20, NOW(), NOW(), NOW()); -- US: $58.00

-- 商品6: 儿童玩具车 - 全市场销售 (需合规)
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 6, NULL, 1, 1, 35.99, NULL, 50, NOW(), NOW(), NOW()), -- US: $35.99
(1, 6, NULL, 2, 1, 29.99, NULL, 50, NOW(), NOW(), NOW()), -- UK: £29.99
(1, 6, NULL, 3, 1, 35.99, NULL, 50, NOW(), NOW(), NOW()), -- DE: €35.99
(1, 6, NULL, 4, 1, 35.99, NULL, 50, NOW(), NOW(), NOW()), -- FR: €35.99
(1, 6, NULL, 5, 1, 49.99, NULL, 50, NOW(), NOW(), NOW()); -- AU: A$49.99

-- 商品7: 化妆品套装 - EU市场有特殊限制，仅US/UK/AU
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 7, NULL, 1, 1, 69.00, NULL, 30, NOW(), NOW(), NOW()), -- US: $69.00
(1, 7, NULL, 2, 1, 55.00, NULL, 30, NOW(), NOW(), NOW()), -- UK: £55.00
(1, 7, NULL, 5, 1, 99.00, NULL, 30, NOW(), NOW(), NOW()); -- AU: A$99.00

-- 商品8: 瑜伽垫 - 全市场销售
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 8, NULL, 1, 1, 15.99, 19.99, 80, NOW(), NOW(), NOW()), -- US: $15.99
(1, 8, NULL, 2, 1, 12.99, 16.99, 80, NOW(), NOW(), NOW()), -- UK: £12.99
(1, 8, NULL, 3, 1, 15.99, 19.99, 80, NOW(), NOW(), NOW()), -- DE: €15.99
(1, 8, NULL, 4, 1, 15.99, 19.99, 80, NOW(), NOW(), NOW()), -- FR: €15.99
(1, 8, NULL, 5, 1, 24.99, 29.99, 80, NOW(), NOW(), NOW()); -- AU: A$24.99

-- 商品9: 草稿商品 - 未发布到任何市场
-- (无 product_markets 记录)

-- 商品10: 测试商品 - 仅US市场，未启用
INSERT INTO `product_markets` (`tenant_id`, `product_id`, `variant_id`, `market_id`, `is_enabled`, `price`, `compare_at_price`, `stock_alert_threshold`, `published_at`, `created_at`, `updated_at`) VALUES
(1, 10, NULL, 1, 0, 12.99, NULL, 10, NULL, NOW(), NOW()); -- US: $12.99 (未启用)


-- ============================================
-- 分类测试数据 (如果不存在)
-- ============================================

INSERT IGNORE INTO `categories` (`id`, `tenant_id`, `name`, `parent_id`, `level`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 1, '数码电子', 0, 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 1, '服饰配饰', 0, 1, 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 1, '运动户外', 0, 1, 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 1, '食品饮料', 0, 1, 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, 1, '母婴玩具', 0, 1, 5, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(6, 1, '美妆护肤', 0, 1, 6, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());


-- ============================================
-- 查询验证
-- ============================================

-- 查看商品市场分布
SELECT
    p.id,
    p.sku,
    p.name,
    p.hs_code,
    p.coo,
    COUNT(pm.id) as market_count,
    GROUP_CONCAT(m.code ORDER BY m.code) as markets
FROM products p
LEFT JOIN product_markets pm ON p.id = pm.product_id AND pm.is_enabled = 1
LEFT JOIN markets m ON pm.market_id = m.id
GROUP BY p.id, p.sku, p.name, p.hs_code, p.coo;

-- 查看合规状态
SELECT
    id,
    sku,
    name,
    hs_code,
    coo,
    weight,
    CASE
        WHEN hs_code IS NOT NULL AND coo IS NOT NULL AND weight IS NOT NULL
        THEN '合规'
        ELSE '不完整'
    END as compliance_status
FROM products
WHERE status != 3;