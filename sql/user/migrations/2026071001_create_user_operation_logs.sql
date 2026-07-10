-- Create user_operation_logs table
CREATE TABLE user_operation_logs (
    id              BIGINT       PRIMARY KEY,
    tenant_id       BIGINT       NOT NULL,
    user_id         BIGINT       NOT NULL,
    action          VARCHAR(64)  NOT NULL,
    operator_id     BIGINT       NOT NULL DEFAULT 0,
    operator_name   VARCHAR(64)  NOT NULL DEFAULT 'system',
    reason          VARCHAR(500) NOT NULL DEFAULT '',
    ip_address      VARCHAR(64)  NOT NULL DEFAULT '',
    user_agent      VARCHAR(500) NOT NULL DEFAULT '',
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP    NULL,

    INDEX idx_uol_user_id (user_id, created_at),
    INDEX idx_uol_tenant  (tenant_id, created_at),
    INDEX idx_uol_action  (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
