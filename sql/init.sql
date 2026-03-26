-- ============================================
-- ShopJoy 数据库初始化脚本
--
-- 执行顺序说明:
-- 1. user/        - 用户、租户、角色表
-- 2. product/     - 商品相关表
-- 3. promotion/   - 促销、优惠券表
-- 4. order/       - 订单、购物车表
-- 5. payment/     - 支付相关表
-- 6. fulfillment/ - 履约、物流表
-- 7. storefront/  - 店铺装修表
-- 8. points/      - 积分表
-- 9. shop/        - 店铺设置表
-- 10. review/     - 评价表
-- ============================================

-- 使用方法:
-- 方式一: 在MySQL客户端中逐个执行
-- source sql/user/schema.sql;
-- source sql/product/schema.sql;
-- ...

-- 方式二: 使用命令行一次性导入
-- mysql -u username -p database_name < sql/user/schema.sql
-- mysql -u username -p database_name < sql/product/schema.sql
-- ...

-- ============================================
-- Schema 文件路径（按领域组织）
-- ============================================

-- 用户与权限领域
-- source user/schema.sql

-- 商品目录领域
-- source product/schema.sql

-- 促销领域
-- source promotion/schema.sql

-- 订单领域
-- source order/schema.sql

-- 支付领域
-- source payment/schema.sql

-- 履约领域
-- source fulfillment/schema.sql

-- 店铺装修领域
-- source storefront/schema.sql

-- 积分领域
-- source points/schema.sql

-- 店铺设置领域
-- source shop/schema.sql

-- 评价领域
-- source review/schema.sql