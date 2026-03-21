-- ============================================
-- ShopJoy 数据库初始化脚本
--
-- 执行顺序说明:
-- 1. tenant.sql     - 租户表 (基础表，其他表依赖)
-- 2. admin_user.sql - 管理员表
-- 3. user.sql       - 用户表
-- 4. role.sql       - 角色权限表
-- 5. market.sql     - 市场表
-- 6. product.sql    - 商品相关表 (依赖分类、市场)
-- 7. storefront.sql - 店铺相关表
-- 8. coupon.sql     - 优惠券表
-- 9. promotion.sql  - 促销活动表
-- 10. cart.sql      - 购物车表
-- 11. order.sql     - 订单表
-- 12. payment.sql   - 支付表
-- 13. fulfillment.sql - 物流表
-- ============================================

-- 使用方法:
-- 方式一: 在MySQL客户端中逐个执行
-- source sql/tenant.sql;
-- source sql/admin_user.sql;
-- ...

-- 方式二: 使用命令行一次性导入
-- mysql -u username -p database_name < sql/init.sql

-- ============================================
-- 以下为所有SQL文件的合并内容
-- ============================================

-- 租户表
source tenant.sql;

-- 管理员表
source admin_user.sql;

-- 用户表
source user.sql;

-- 角色权限表
source role.sql;

-- 市场表
source market.sql;

-- 商品相关表
source product.sql;

-- 店铺相关表
source storefront.sql;

-- 优惠券表
source coupon.sql;

-- 促销活动表
source promotion.sql;

-- 购物车表
source cart.sql;

-- 订单表
source order.sql;

-- 支付表
source payment.sql;

-- 物流表
source fulfillment.sql;