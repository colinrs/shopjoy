-- 创建区域目录表（多层级：国家/省/市/区）
-- 用于国际运费配置的多级区域匹配（Phase 2 配套）
CREATE TABLE regions (
    id              BIGINT       PRIMARY KEY,
    tenant_id       BIGINT       NOT NULL DEFAULT 0 COMMENT '0=平台预置数据',
    code            VARCHAR(20)  NOT NULL COMMENT '区域编码（如 US、US-CA、110000）',
    name            VARCHAR(100) NOT NULL,
    level           INT          NOT NULL DEFAULT 1 COMMENT '1=country, 2=province, 3=city, 4=district',
    parent_code     VARCHAR(20)  NOT NULL DEFAULT '' COMMENT '父级区域code',
    country_code    VARCHAR(2)   NOT NULL DEFAULT 'CN' COMMENT 'ISO 3166-1 alpha-2',
    postal_pattern  VARCHAR(255) NULL COMMENT '邮编正则模式',
    sort            INT          NOT NULL DEFAULT 0,
    is_active       TINYINT(1)   NOT NULL DEFAULT 1,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP    NULL,
    UNIQUE KEY uk_tenant_code (tenant_id, code, deleted_at),
    INDEX idx_country (country_code),
    INDEX idx_country_parent (country_code, parent_code),
    INDEX idx_parent (parent_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区域目录表（多层级国家/省/市）';