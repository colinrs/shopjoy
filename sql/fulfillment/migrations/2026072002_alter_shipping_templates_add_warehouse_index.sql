-- 同步 schema.sql 与 live DB：shipping_templates 表的 warehouse_id 列索引
-- 详见 Task 1.4 review Concern #3
ALTER TABLE shipping_templates ADD INDEX idx_warehouse_id (warehouse_id);