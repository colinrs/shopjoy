-- ============================================
-- 媒体资源表 (media_assets)
-- ============================================
-- Backend-managed asset metadata for the multi-driver (Local + Cloudinary)
-- image upload pipeline. See docs/superpowers/specs/2026-07-11-cloudinary-image-storage-design.md.

CREATE TABLE IF NOT EXISTS `media_assets` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `public_id`  VARCHAR(255) NOT NULL                COMMENT '存储驱动返回的资产ID (Cloudinary public_id 或 本地相对路径)',
    `url`        VARCHAR(1024) NOT NULL               COMMENT '可访问的资产URL',
    `filename`   VARCHAR(255) NOT NULL DEFAULT ''     COMMENT '原始文件名',
    `size_bytes` BIGINT       NOT NULL DEFAULT 0      COMMENT '文件大小 (字节)',
    `mime_type`  VARCHAR(64)  NOT NULL DEFAULT ''     COMMENT 'MIME 类型',
    `width`      INT          NOT NULL DEFAULT 0      COMMENT '图片宽度 (px)',
    `height`     INT          NOT NULL DEFAULT 0      COMMENT '图片高度 (px)',
    `format`     VARCHAR(32)  NOT NULL DEFAULT ''     COMMENT '图片格式 (jpg/png/webp/...)',
    `category`   VARCHAR(32)  NOT NULL DEFAULT 'common' COMMENT '业务分类: product/banner/avatar/common',
    `provider`   VARCHAR(16)  NOT NULL DEFAULT 'local' COMMENT '存储驱动: local/cloudinary',
    `tenant_id`  BIGINT       NOT NULL                COMMENT '租户ID',
    `created_by` BIGINT       NOT NULL DEFAULT 0      COMMENT '创建人',
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP    NULL                    COMMENT '删除时间 (软删除)',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_provider_public` (`provider`, `public_id`),
    KEY `idx_tenant_category` (`tenant_id`, `category`, `deleted_at`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='媒体资源元数据表';
