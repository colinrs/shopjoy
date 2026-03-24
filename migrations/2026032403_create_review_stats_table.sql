-- Create review_stats table
CREATE TABLE `review_stats` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `total_reviews` INT NOT NULL DEFAULT 0,
    `average_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `quality_avg_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `value_avg_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    `rating_1_count` INT NOT NULL DEFAULT 0,
    `rating_2_count` INT NOT NULL DEFAULT 0,
    `rating_3_count` INT NOT NULL DEFAULT 0,
    `rating_4_count` INT NOT NULL DEFAULT 0,
    `rating_5_count` INT NOT NULL DEFAULT 0,
    `with_image_count` INT NOT NULL DEFAULT 0,
    `last_updated_at` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_tenant_product` (`tenant_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;