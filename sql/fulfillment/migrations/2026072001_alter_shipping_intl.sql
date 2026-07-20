-- sql/fulfillment/migrations/2026072001_alter_shipping_intl.sql
-- 一次性为 shipping_templates 和 shipping_zones 增加国际化所需全部字段

ALTER TABLE shipping_templates
    ADD COLUMN market_id      BIGINT       NOT NULL DEFAULT 0 COMMENT '市场ID，0=全市场通用',
    ADD COLUMN currency       VARCHAR(3)   NOT NULL DEFAULT 'CNY' COMMENT 'ISO 4217',
    ADD COLUMN carrier_code   VARCHAR(50)  NOT NULL DEFAULT 'standard' COMMENT '物流商代码',
    ADD COLUMN warehouse_id   BIGINT       NOT NULL DEFAULT 0 COMMENT '发货仓库ID',
    ADD INDEX idx_market_default (market_id, is_default);

ALTER TABLE shipping_zones
    ADD COLUMN market_id              BIGINT       NOT NULL DEFAULT 0 COMMENT '市场ID',
    ADD COLUMN currency               VARCHAR(3)   NOT NULL DEFAULT 'CNY',
    ADD COLUMN name_i18n              JSON         NULL COMMENT '多语言名称',
    ADD COLUMN taxable                TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否计税',
    ADD COLUMN tax_rate               DECIMAL(5,4) NOT NULL DEFAULT 0 COMMENT '税率 0.0000-1.0000',
    ADD COLUMN tax_included           TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '价格含税',
    ADD COLUMN ioss_applicable        TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '欧盟IOSS申报',
    ADD COLUMN remote_surcharge       DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '偏远地区附加费',
    ADD COLUMN remote_zip_patterns    JSON         NULL COMMENT '偏远邮编正则',
    ADD COLUMN fuel_surcharge_pct     DECIMAL(5,4) NOT NULL DEFAULT 0 COMMENT '燃油附加费%',
    ADD COLUMN volumetric_divisor     INT          NOT NULL DEFAULT 5000 COMMENT '体积重除数 cm³/kg',
    ADD INDEX idx_zone_market (market_id);