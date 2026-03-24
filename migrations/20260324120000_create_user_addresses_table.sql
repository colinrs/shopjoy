-- migrations/20260324120000_create_user_addresses_table.sql

-- +migrate Up
CREATE TABLE `user_addresses` (
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
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收货地址表';

-- +migrate Down
DROP TABLE IF EXISTS `user_addresses`;