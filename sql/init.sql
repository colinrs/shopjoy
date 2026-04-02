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

-- ============================================
-- 推荐使用方式: init.sh 脚本 (自动连接数据库并初始化)
-- ============================================

-- 方式一: 使用 init.sh 脚本 (推荐)
-- cd sql
-- chmod +x init.sh
-- ./init.sh

-- 或设置环境变量后执行:
-- export MYSQL_HOST=localhost
-- export MYSQL_PORT=3306
-- export MYSQL_USER=root
-- export MYSQL_PASSWORD=your_password
-- export MYSQL_DATABASE=shopjoy
-- ./init.sh

-- 方式二: 使用命令行参数:
-- ./init.sh -h localhost -P 3306 -u root -p your_password -d shopjoy

-- 方式三: 在MySQL客户端中逐个执行 (手动方式)
-- source sql/user/schema.sql;
-- source sql/product/schema.sql;
-- ...

-- 方式四: 使用命令行一次性导入 (手动方式)
-- mysql -u username -p database_name < sql/user/schema.sql
-- mysql -u username -p database_name < sql/product/schema.sql
-- ...

-- init.sh 高级选项:
--   --only=user,product    仅初始化指定模块
--   --skip=points,review   跳过指定模块
--   --tables-only          仅创建表，不插入数据
--   --data-only            仅插入数据，不创建表 (需表已存在)
--   --all                  创建表并插入数据 (默认)
--   --dry-run              模拟运行

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