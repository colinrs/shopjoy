-- Create promotion_usage table
CREATE TABLE IF NOT EXISTS promotion_usage (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id       BIGINT NOT NULL,
    promotion_id    BIGINT NOT NULL,
    rule_id         BIGINT DEFAULT NULL,
    order_id        VARCHAR(50) NOT NULL,
    user_id         BIGINT NOT NULL,
    discount_amount BIGINT NOT NULL DEFAULT 0,
    currency        VARCHAR(10) NOT NULL DEFAULT 'CNY',
    original_amount BIGINT NOT NULL DEFAULT 0,
    final_amount    BIGINT NOT NULL DEFAULT 0,
    coupon_id       BIGINT DEFAULT NULL,
    created_at      BIGINT NOT NULL,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_promotion_id (promotion_id),
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;