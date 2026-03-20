-- admin_users 表（管理后台用户：平台超管 + 商家管理员 + 商家子账号）
CREATE TABLE IF NOT EXISTS `admin_users` (
    `id`              bigint(20)      NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
    `tenant_id`       bigint(20)      NOT NULL DEFAULT 0 COMMENT '租户ID，0表示平台超管',
    `username`        varchar(64)     NULL UNIQUE COMMENT '登录用户名',
    `email`           varchar(128)    NOT NULL UNIQUE COMMENT '邮箱（商家管理员常用）',
    `mobile`          varchar(20)     NULL UNIQUE COMMENT '手机号',
    `password`        varchar(255)    NOT NULL COMMENT 'bcrypt加密密码',
    `real_name`       varchar(32)     NULL COMMENT '真实姓名',
    `avatar`          varchar(255)    NULL COMMENT '头像URL',
    `type`            tinyint(4)      NOT NULL DEFAULT 1 COMMENT '1=平台超管 2=商家管理员 3=商家子账号',
    `status`          tinyint(4)      NOT NULL DEFAULT 1 COMMENT '1=正常 2=禁用 3=已删除',
    `last_login_at`   datetime        NULL COMMENT '最后登录时间',
    `last_login_ip`   varchar(45)     NULL COMMENT '最后登录IP',
    `created_at`      datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`      datetime        NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_email` (`email`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_mobile` (`mobile`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理后台用户表';

-- users 表（C端买家用户）
CREATE TABLE IF NOT EXISTS `users` (
    `id`              bigint(20)      NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `tenant_id`       bigint(20)      NOT NULL COMMENT '所属租户ID（必填）',
    `nickname`        varchar(64)     NOT NULL COMMENT '昵称',
    `avatar`          varchar(255)    NULL COMMENT '头像URL',
    `mobile`          varchar(20)     NULL UNIQUE COMMENT '手机号',
    `email`           varchar(128)    NULL UNIQUE COMMENT '邮箱',
    `gender`          tinyint(4)      NOT NULL DEFAULT 0 COMMENT '0=未知 1=男 2=女',
    `birthday`        date            NULL COMMENT '生日',
    `status`          tinyint(4)      NOT NULL DEFAULT 1 COMMENT '1=正常 2=禁用 3=注销',
    `source`          tinyint(4)      NOT NULL DEFAULT 1 COMMENT '注册来源：1=直接注册 2=微信 3=其他',
    `total_amount`    decimal(19,4)   NOT NULL DEFAULT 0.0000 COMMENT '累计消费金额',
    `order_count`     int(11)         NOT NULL DEFAULT 0 COMMENT '订单总数',
    `last_login_at`   datetime        NULL COMMENT '最后登录时间',
    `created_at`      datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`      datetime        NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_mobile` (`mobile`),
    KEY `idx_email` (`email`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='C端用户表';

-- 插入平台超级管理员（默认账号）
-- 密码为 bcrypt 加密的 "admin123"
INSERT INTO `admin_users` (`id`, `tenant_id`, `username`, `email`, `mobile`, `password`, `real_name`, `type`, `status`, `created_at`, `updated_at`)
VALUES (1, 0, 'superadmin', 'admin@shopjoy.com', '13800138000', '$2a$10$N9qo8uLOickgx2ZMRZoMye1jG0J0qF6zHq3Xq8u0yqF6zHq3Xq8u', '平台管理员', 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `email` = 'admin@shopjoy.com';
