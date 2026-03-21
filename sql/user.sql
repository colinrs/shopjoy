-- ============================================
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
    `birthday` BIGINT DEFAULT NULL COMMENT '生日(Unix时间戳)',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未激活, 1-正常, 2-暂停, 3-已删除',
    `last_login` BIGINT DEFAULT NULL COMMENT '最后登录时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
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
(1, 1, 'user1@example.com', '13800000001', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '小明', 'https://cdn.example.com/u1.png', 1, UNIX_TIMESTAMP('1990-05-15'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(2, 1, 'user2@example.com', '13800000002', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '小红', 'https://cdn.example.com/u2.png', 2, UNIX_TIMESTAMP('1995-08-20'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(3, 1, 'user3@example.com', '13800000003', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '小刚', 'https://cdn.example.com/u3.png', 1, UNIX_TIMESTAMP('1988-03-10'), 1, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),

-- Test Store 用户
(4, 2, 'user4@example.com', '13800000004', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '大卫', 'https://cdn.example.com/u4.png', 1, UNIX_TIMESTAMP('1992-11-25'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(5, 2, 'user5@example.com', '13800000005', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '艾米', 'https://cdn.example.com/u5.png', 2, UNIX_TIMESTAMP('1998-01-05'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),

-- Enterprise Corp 用户
(6, 3, 'user6@example.com', '13800000006', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '约翰', 'https://cdn.example.com/u6.png', 1, UNIX_TIMESTAMP('1985-07-12'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(7, 3, 'user7@example.com', '13800000007', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '玛丽', 'https://cdn.example.com/u7.png', 2, UNIX_TIMESTAMP('1993-04-30'), 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0),
(8, 3, 'user8@example.com', '13800000008', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt8k0US', '汤姆', 'https://cdn.example.com/u8.png', 1, UNIX_TIMESTAMP('1990-09-18'), 2, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0, 0);