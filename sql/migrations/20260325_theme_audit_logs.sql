-- ============================================
-- Theme Audit Logs Migration
-- Date: 2026-03-25
-- Description: Add theme audit logging for tracking theme changes
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
    `created_at` BIGINT NOT NULL COMMENT 'Created timestamp (Unix UTC)',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_action` (`tenant_id`, `action`),
    INDEX `idx_tenant_theme` (`tenant_id`, `theme_id`),
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Theme change audit logs';