-- ============================================
-- 促销活动表 (promotions)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotions` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '促销ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '促销名称',
    `description` TEXT COMMENT '描述',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型: 0-折扣, 1-限时抢购, 2-捆绑销售, 3-买X送Y',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待生效, 1-生效中, 2-已暂停, 3-已结束',
    `start_at` BIGINT NOT NULL DEFAULT 0 COMMENT '开始时间',
    `end_at` BIGINT NOT NULL DEFAULT 0 COMMENT '结束时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_type` (`type`),
    KEY `idx_start_end` (`start_at`, `end_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销活动表';

-- ============================================
-- 促销规则表 (promotion_rules)
-- ============================================

CREATE TABLE IF NOT EXISTS `promotion_rules` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '规则ID',
    `promotion_id` BIGINT NOT NULL COMMENT '促销ID',
    `condition_type` TINYINT NOT NULL DEFAULT 0 COMMENT '条件类型: 0-最低金额, 1-最低数量',
    `condition_value` BIGINT NOT NULL DEFAULT 0 COMMENT '条件值',
    `action_type` TINYINT NOT NULL DEFAULT 0 COMMENT '动作类型: 0-固定金额, 1-百分比',
    `action_value` BIGINT NOT NULL DEFAULT 0 COMMENT '动作值',
    `max_discount_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '最大优惠金额(分)',
    `max_discount_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    PRIMARY KEY (`id`),
    KEY `idx_promotion_id` (`promotion_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='促销规则表';

-- ============================================
-- 测试数据
-- ============================================

-- 促销活动数据 (Demo Shop)
INSERT INTO `promotions` (`id`, `tenant_id`, `name`, `description`, `type`, `status`, `start_at`, `end_at`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, '春季大促', '春季全场促销活动', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 30 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, '限时抢购-鞋类专场', '运动鞋限时抢购', 1, 1, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 1 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 8 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, '买二送一', '购买任意两件商品送一件同等价值商品', 3, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 60 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, '会员日特惠', '每月18号会员专享折扣', 0, 2, UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 10 DAY)), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 40 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, '清仓大甩卖', '库存清仓特价', 1, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 15 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),

-- Enterprise Corp 促销活动
(6, 3, '企业采购季', '企业客户采购优惠季', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 90 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5),
(7, 3, '跨境特惠', '跨境电商专属优惠', 0, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(DATE_ADD(NOW(), INTERVAL 180 DAY)), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 5, 5);

-- 促销规则数据
INSERT INTO `promotion_rules` (`id`, `promotion_id`, `condition_type`, `condition_value`, `action_type`, `action_value`, `max_discount_amount`, `max_discount_currency`) VALUES
-- 春季大促规则
(1, 1, 0, 10000, 1, 10, 5000, 'CNY'),   -- 满100减10%, 最多减50
(2, 1, 0, 30000, 1, 15, 10000, 'CNY'),  -- 满300减15%, 最多减100
(3, 1, 0, 50000, 1, 20, 20000, 'CNY'),  -- 满500减20%, 最多减200

-- 限时抢购规则
(4, 2, 0, 0, 1, 30, 10000, 'CNY'),      -- 无门槛7折, 最多减100

-- 买二送一规则
(5, 3, 1, 2, 0, 0, 0, 'CNY'),           -- 买2件送1件

-- 会员日特惠规则
(6, 4, 0, 20000, 1, 25, 15000, 'CNY'),  -- 满200减25%, 最多减150

-- 清仓大甩卖规则
(7, 5, 0, 0, 1, 40, 20000, 'CNY'),      -- 无门槛6折, 最多减200

-- 企业采购季规则
(8, 6, 0, 50000, 1, 15, 30000, 'CNY'),  -- 满500减15%, 最多减300
(9, 6, 0, 100000, 1, 20, 50000, 'CNY'), -- 满1000减20%, 最多减500

-- 跨境特惠规则
(10, 7, 0, 30000, 1, 10, 20000, 'CNY'); -- 满300减10%, 最多减200