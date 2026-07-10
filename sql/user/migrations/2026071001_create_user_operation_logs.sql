-- Create user_operation_logs table
CREATE TABLE IF NOT EXISTS `user_operation_logs` (
    `id` BIGINT NOT NULL COMMENT '日志ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `action` VARCHAR(64) NOT NULL COMMENT '操作类型',
    `operator_id` BIGINT NOT NULL DEFAULT 0 COMMENT '操作人ID',
    `operator_name` VARCHAR(64) NOT NULL DEFAULT 'system' COMMENT '操作人名称',
    `reason` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '操作原因',
    `ip_address` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '操作IP',
    `user_agent` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '客户端UA',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',

    INDEX `idx_uol_user_id` (`user_id`, `created_at`),
    INDEX `idx_uol_tenant` (`tenant_id`, `created_at`),
    INDEX `idx_uol_action` (`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户操作日志表';