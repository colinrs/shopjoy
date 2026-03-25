-- Create review_replies table
CREATE TABLE `review_replies` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `review_id` BIGINT NOT NULL,
    `tenant_id` BIGINT NOT NULL,
    `admin_id` BIGINT NOT NULL,
    `admin_name` VARCHAR(100) NOT NULL,
    `content` TEXT NOT NULL,
    `created_at` BIGINT NOT NULL,
    `updated_at` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_review_id` (`review_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;