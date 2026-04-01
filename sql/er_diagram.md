# ShopJoy Database ER Diagram

```mermaid
erDiagram
    %% ==================== TENANT (Core) ====================
    tenants {
        bigint id PK "租户ID"
        varchar name "租户名称"
        varchar code "租户代码"
        tinyint status "状态"
        tinyint plan "套餐"
        varchar domain "系统域名"
        varchar custom_domain "自定义域名"
        varchar logo "Logo"
        varchar contact_name "联系人"
        varchar contact_phone "联系电话"
        varchar contact_email "联系邮箱"
        text address "地址"
        timestamp expire_at "过期时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    %% ==================== USER DOMAIN ====================
    users {
        bigint id PK "用户ID"
        bigint tenant_id FK "租户ID"
        varchar email "邮箱"
        varchar phone "手机号"
        varchar password "密码"
        varchar name "昵称"
        varchar avatar "头像"
        tinyint gender "性别"
        timestamp birthday "生日"
        tinyint status "状态"
        timestamp last_login "最后登录"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    admin_users {
        bigint id PK "用户ID"
        bigint tenant_id FK "租户ID (0=平台超管)"
        varchar username "用户名"
        varchar email "邮箱"
        varchar mobile "手机号"
        varchar password "密码"
        varchar real_name "真实姓名"
        varchar avatar "头像"
        tinyint type "类型"
        tinyint status "状态"
        timestamp last_login_at "最后登录"
        varchar last_login_ip "最后登录IP"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    roles {
        bigint id PK "角色ID"
        bigint tenant_id FK "租户ID"
        varchar name "角色名称"
        varchar code "角色代码"
        text description "描述"
        tinyint status "状态"
        tinyint is_system "是否系统角色"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    permissions {
        bigint id PK "权限ID"
        varchar name "权限名称"
        varchar code "权限代码"
        tinyint type "类型"
        bigint parent_id "父ID"
        varchar path "路径"
        varchar icon "图标"
        int sort "排序"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    user_roles {
        bigint user_id PK,FK "用户ID"
        bigint role_id PK,FK "角色ID"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    role_permissions {
        bigint role_id PK,FK "角色ID"
        bigint permission_id PK,FK "权限ID"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    user_addresses {
        bigint id PK "地址ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        varchar name "收货人"
        varchar phone "电话"
        varchar country "国家"
        varchar province "省份"
        varchar city "城市"
        varchar district "区县"
        varchar address "地址"
        varchar postal_code "邮编"
        tinyint is_default "是否默认"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== PRODUCT DOMAIN ====================
    categories {
        bigint id PK "分类ID"
        bigint tenant_id FK "租户ID"
        bigint parent_id "父分类ID"
        varchar name "名称"
        varchar code "代码"
        tinyint level "层级"
        int sort "排序"
        varchar icon "图标"
        varchar image "图片"
        varchar seo_title "SEO标题"
        varchar seo_description "SEO描述"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    brands {
        bigint id PK "品牌ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        varchar logo "Logo"
        text description "描述"
        varchar website "官网"
        int sort "排序"
        tinyint enable_page "启用品牌专区"
        varchar trademark_number "商标号"
        varchar trademark_country "商标国家"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    products {
        bigint id PK "商品ID"
        bigint tenant_id FK "租户ID"
        varchar sku "SKU代码"
        varchar name "名称"
        text description "描述"
        decimal price "售价"
        decimal cost_price "成本价"
        varchar currency "货币"
        int stock "库存"
        int status "状态"
        bigint category_id FK "分类ID"
        varchar brand "品牌"
        bigint brand_id "品牌ID (关联 brands.id, 非外键约束)"
        varchar sku_prefix "SKU前缀"
        json tags "标签"
        json images "图片"
        tinyint is_matrix_product "是否有变体"
        varchar hs_code "HS编码"
        varchar coo "原产国"
        decimal weight "重量"
        varchar weight_unit "重量单位"
        decimal length "长度"
        decimal width "宽度"
        decimal height "高度"
        json dangerous_goods "危险品"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    skus {
        bigint id PK "SKU ID"
        bigint tenant_id FK "租户ID"
        bigint product_id FK "商品ID"
        varchar code "代码"
        decimal price_amount "价格"
        varchar price_currency "货币"
        int stock "库存"
        int available_stock "可用库存"
        int locked_stock "锁定库存"
        int safety_stock "安全库存"
        tinyint presale_enabled "预售"
        json attributes "属性"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    markets {
        bigint id PK "市场ID"
        bigint tenant_id FK "租户ID"
        varchar code "代码"
        varchar name "名称"
        varchar currency "货币"
        varchar default_language "默认语言"
        varchar flag "旗帜"
        tinyint is_active "启用"
        tinyint is_default "主市场"
        json tax_rules "税务配置"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    product_markets {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint product_id FK "商品ID"
        bigint variant_id "变体ID (关联 skus.id, 非外键约束)"
        bigint market_id FK "市场ID"
        tinyint is_enabled "启用"
        int status_override "状态覆盖"
        decimal price "价格"
        decimal compare_at_price "对比价"
        int stock_alert_threshold "库存预警"
        timestamp published_at "发布时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    product_localizations {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint product_id FK "商品ID"
        varchar language_code "语言代码"
        varchar name "名称"
        text description "描述"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    category_markets {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint category_id FK "分类ID"
        bigint market_id FK "市场ID"
        tinyint is_visible "可见"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    brand_markets {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint brand_id FK "品牌ID"
        bigint market_id FK "市场ID"
        tinyint is_visible "可见"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    warehouses {
        bigint id PK "仓库ID"
        bigint tenant_id FK "租户ID"
        varchar code "代码"
        varchar name "名称"
        varchar country "国家"
        varchar address "地址"
        tinyint is_default "默认"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    warehouse_inventories {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        varchar sku_code "SKU代码 (非外键约束)"
        bigint warehouse_id FK "仓库ID"
        int available_stock "可用"
        int locked_stock "锁定"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    inventory_logs {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        varchar sku_code "SKU代码 (非外键约束)"
        bigint product_id FK "商品ID"
        bigint warehouse_id "仓库ID"
        varchar change_type "类型"
        int change_quantity "数量"
        int before_stock "前库存"
        int after_stock "后库存"
        varchar order_no "订单号"
        varchar remark "备注"
        bigint operator_id "操作人"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== ORDER DOMAIN ====================
    orders {
        bigint id PK "订单ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        varchar order_no "订单号"
        tinyint status "状态"
        tinyint fulfillment_status "履约状态"
        tinyint refund_status "退款状态"
        decimal total_amount "总额"
        decimal discount_amount "优惠"
        decimal freight_amount "运费"
        decimal pay_amount "实付"
        decimal original_amount "原价"
        decimal adjust_amount "改价"
        varchar adjust_reason "改价原因"
        bigint adjusted_by "改价人"
        timestamp adjusted_at "改价时间"
        int version "版本"
        varchar payment_method "支付方式"
        varchar source "来源"
        varchar currency "货币"
        varchar address_name "收货人"
        varchar address_phone "电话"
        varchar address_province "省份"
        varchar address_city "城市"
        varchar address_district "区县"
        text address_detail "地址"
        varchar address_zipcode "邮编"
        varchar tracking_no "快递号"
        varchar carrier "快递公司"
        text remark "备注"
        varchar merchant_remark "商家备注"
        timestamp expire_at "过期时间"
        timestamp paid_at "支付时间"
        timestamp shipped_at "发货时间"
        timestamp completed_at "完成时间"
        timestamp cancelled_at "取消时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    order_items {
        bigint id PK "ID"
        bigint order_id FK "订单ID"
        bigint product_id FK "商品ID"
        bigint sku_id FK "SKU ID"
        varchar product_name "商品名称"
        varchar sku_name "SKU名称"
        varchar image "图片"
        decimal price "单价"
        int quantity "数量"
        decimal total_amount "小计"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    carts {
        bigint id PK "购物车ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        varchar session_id "会话ID"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    cart_items {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        bigint cart_id FK "购物车ID"
        bigint product_id FK "商品ID"
        bigint sku_id FK "SKU ID"
        varchar product_name "商品名称"
        varchar sku_name "SKU名称"
        varchar image "图片"
        decimal price "单价"
        int quantity "数量"
        decimal total_amount "小计"
        tinyint selected "选中"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    %% ==================== PAYMENT DOMAIN ====================
    order_payments {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint order_id FK "订单ID"
        varchar payment_no "支付号"
        varchar payment_method "方式"
        varchar channel_intent_id "Channel PaymentIntent"
        varchar channel_payment_id "Channel Charge ID"
        decimal amount "金额"
        varchar currency "货币"
        tinyint status "状态"
        decimal transaction_fee "手续费"
        varchar fee_currency "手续费货币"
        timestamp paid_at "支付时间"
        timestamp failed_at "失败时间"
        varchar failed_reason "失败原因"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    payment_transactions {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint order_id FK "订单ID"
        bigint payment_id FK "支付ID"
        varchar transaction_id "交易ID"
        varchar payment_method "方式"
        varchar channel_transaction_id "渠道交易ID"
        decimal amount "金额"
        varchar currency "货币"
        tinyint status "状态"
        decimal transaction_fee "手续费"
        timestamp paid_at "支付时间"
        varchar failed_reason "失败原因"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    payment_refunds {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint order_id FK "订单ID"
        bigint payment_id FK "支付ID"
        bigint fulfillment_refund_id "履约退款ID (关联 fulfillment.refunds)"
        varchar refund_no "退款号"
        varchar idempotency_key "幂等Key"
        varchar channel_refund_id "渠道退款ID"
        decimal amount "金额"
        varchar currency "货币"
        decimal refund_fee "退款手续费"
        tinyint status "状态"
        varchar reason_type "原因类型"
        varchar reason "原因"
        timestamp refunded_at "退款时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
    }

    webhook_events {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        varchar event_id "事件ID"
        varchar event_type "事件类型"
        varchar resource_id "资源ID"
        tinyint processed "处理状态"
        text raw_payload "原始数据"
        varchar error_message "错误信息"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp processed_at "处理时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== FULFILLMENT DOMAIN ====================
    shipments {
        bigint id PK "发货ID"
        bigint tenant_id FK "租户ID"
        bigint order_id FK "订单ID"
        varchar shipment_no "发货号"
        tinyint status "状态"
        varchar carrier "快递公司"
        varchar carrier_code "快递代码"
        varchar tracking_no "快递号"
        decimal weight "重量"
        decimal cost_amount "成本"
        varchar cost_currency "货币"
        varchar remark "备注"
        timestamp shipped_at "发货时间"
        timestamp delivered_at "送达时间"
        timestamp cancelled_at "取消时间"
        bigint cancelled_by "取消人"
        varchar cancelled_reason "取消原因"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    shipment_items {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint shipment_id FK "发货ID"
        bigint order_item_id FK "订单商品ID"
        bigint product_id FK "商品ID"
        bigint sku_id FK "SKU ID"
        int quantity "数量"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    refunds {
        bigint id PK "退款ID"
        bigint tenant_id FK "租户ID"
        bigint order_id FK "订单ID"
        varchar refund_no "退款号"
        bigint user_id FK "用户ID"
        tinyint type "类型"
        tinyint status "状态"
        varchar reason_type "原因类型"
        varchar reason "原因"
        text description "描述"
        json images "凭证"
        decimal amount "金额"
        varchar currency "货币"
        varchar reject_reason "拒绝原因"
        timestamp approved_at "批准时间"
        bigint approved_by "批准人"
        timestamp completed_at "完成时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    shipping_templates {
        bigint id PK "模板ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        tinyint is_default "默认"
        tinyint is_active "启用"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    shipping_zones {
        bigint id PK "区域ID"
        bigint tenant_id FK "租户ID"
        bigint template_id FK "模板ID"
        varchar name "名称"
        json regions "区域"
        varchar fee_type "费用类型"
        int first_unit "首单位"
        decimal first_fee "首费"
        int additional_unit "续单位"
        decimal additional_fee "续费"
        decimal free_threshold_amount "免邮门槛金额"
        int free_threshold_count "免邮门槛数量"
        int sort "排序"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    shipping_template_mappings {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint template_id FK "模板ID"
        varchar target_type "目标类型"
        bigint target_id "目标ID"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    shipping_zone_regions {
        bigint id PK "ID"
        bigint zone_id FK "区域ID"
        varchar city_code "城市代码"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== PROMOTION DOMAIN ====================
    promotions {
        bigint id PK "促销ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        text description "描述"
        tinyint type "类型"
        tinyint status "状态"
        int priority "优先级"
        varchar currency "货币"
        varchar scope_type "范围类型"
        json scope_ids "范围ID"
        json exclude_ids "排除ID"
        timestamp start_at "开始时间"
        timestamp end_at "结束时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    promotion_rules {
        bigint id PK "规则ID"
        bigint promotion_id FK "促销ID"
        tinyint condition_type "条件类型"
        decimal condition_value "条件值"
        tinyint action_type "动作类型"
        decimal action_value "动作值"
        decimal max_discount_amount "最大优惠"
        varchar max_discount_currency "货币"
        varchar currency "货币"
        int sort_order "排序"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    promotion_usage {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint promotion_id FK "促销ID"
        bigint rule_id FK "规则ID"
        bigint order_id FK "订单ID"
        bigint user_id FK "用户ID"
        decimal discount_amount "优惠金额"
        varchar currency "货币"
        decimal original_amount "原价"
        decimal final_amount "实付"
        bigint coupon_id "优惠券ID (非外键约束)"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    coupons {
        bigint id PK "优惠券ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        varchar code "代码"
        text description "描述"
        tinyint type "类型"
        decimal value "优惠值"
        decimal min_amount "最低消费"
        decimal max_discount "最大优惠"
        int total_count "总数"
        int used_count "已用"
        int per_user_limit "限领"
        tinyint status "状态"
        varchar currency "货币"
        varchar scope_type "范围类型"
        json scope_ids "范围ID"
        json exclude_ids "排除ID"
        timestamp start_at "开始时间"
        timestamp end_at "结束时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    user_coupons {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        bigint coupon_id FK "优惠券ID"
        tinyint status "状态"
        timestamp used_at "使用时间"
        bigint order_id "订单ID (关联 orders.id, 非外键约束)"
        timestamp received_at "领取时间"
        timestamp expire_at "过期时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== POINTS DOMAIN ====================
    earn_rules {
        bigint id PK "规则ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        text description "描述"
        varchar scenario "场景"
        varchar calculation_type "计算类型"
        bigint fixed_points "固定积分"
        decimal ratio "比例"
        text tiers "阶梯"
        varchar condition_type "条件类型"
        text condition_value "条件值"
        int expiration_months "过期月数"
        tinyint status "状态"
        int priority "优先级"
        timestamp start_at "开始时间"
        timestamp end_at "结束时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    redeem_rules {
        bigint id PK "规则ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        text description "描述"
        bigint coupon_id FK "优惠券ID (关联 coupons.id)"
        bigint points_required "所需积分"
        bigint total_stock "总库存"
        bigint used_stock "已用"
        int per_user_limit "限兑"
        tinyint status "状态"
        timestamp start_at "开始时间"
        timestamp end_at "结束时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
        timestamp deleted_at "删除时间"
    }

    points_accounts {
        bigint id PK "账户ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        bigint balance "余额"
        bigint frozen_balance "冻结"
        bigint total_earned "累计获得"
        bigint total_redeemed "累计兑换"
        bigint total_expired "累计过期"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    points_transactions {
        bigint id PK "交易ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        bigint account_id FK "账户ID"
        bigint points "积分"
        bigint balance_after "变动后余额"
        varchar type "类型"
        varchar reference_type "关联类型"
        varchar reference_id "关联ID"
        varchar description "描述"
        timestamp expires_at "过期时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    points_redemptions {
        bigint id PK "兑换ID"
        bigint tenant_id FK "租户ID"
        bigint user_id FK "用户ID"
        bigint redeem_rule_id FK "规则ID"
        bigint coupon_id FK "优惠券ID"
        bigint user_coupon_id "用户优惠券ID (关联 user_coupons)"
        bigint points_used "消耗积分"
        varchar status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp completed_at "完成时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== REVIEW DOMAIN ====================
    reviews {
        bigint id PK "评价ID"
        bigint tenant_id FK "租户ID"
        bigint order_id "订单ID (关联 orders.id)"
        bigint product_id FK "商品ID"
        varchar sku_code "SKU代码"
        bigint user_id FK "用户ID"
        varchar user_name "用户名"
        tinyint quality_rating "质量评分"
        tinyint value_rating "价值评分"
        decimal overall_rating "综合评分"
        text content "内容"
        json images "图片"
        tinyint status "状态"
        boolean is_anonymous "匿名"
        boolean is_verified "已验证"
        boolean is_featured "精选"
        int helpful_count "有帮助"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    review_replies {
        bigint id PK "回复ID"
        bigint review_id FK "评价ID"
        bigint tenant_id FK "租户ID"
        bigint admin_id "管理员ID (关联 admin_users, 非外键约束)"
        varchar admin_name "管理员名"
        text content "内容"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    review_stats {
        bigint id PK "ID"
        bigint tenant_id FK "租户ID"
        bigint product_id FK "商品ID"
        int total_reviews "总数"
        decimal average_rating "平均分"
        decimal quality_avg_rating "质量均分"
        decimal value_avg_rating "价值均分"
        int rating_1_count "1星"
        int rating_2_count "2星"
        int rating_3_count "3星"
        int rating_4_count "4星"
        int rating_5_count "5星"
        int with_image_count "图评数"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== STOREFRONT DOMAIN ====================
    shops {
        bigint id PK "店铺ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        text description "描述"
        varchar logo "Logo"
        varchar banner "Banner"
        varchar contact_phone "电话"
        varchar contact_email "邮箱"
        text address "地址"
        json social_links "社交链接"
        varchar seo_title "SEO标题"
        text seo_description "SEO描述"
        varchar seo_keywords "SEO关键词"
        bigint current_theme_id "当前主题 (关联 themes.id)"
        text theme_config "主题配置"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    themes {
        bigint id PK "主题ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        varchar code "代码"
        text description "描述"
        varchar thumbnail "缩略图"
        varchar preview_image "预览图"
        json config "配置"
        text config_schema "配置Schema"
        text default_config "默认配置"
        tinyint is_active "激活"
        tinyint is_custom "自定义"
        tinyint is_preset "预设"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    pages {
        bigint id PK "页面ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        varchar slug "URL别名"
        tinyint type "类型"
        longtext content "内容"
        varchar seo_title "SEO标题"
        text seo_description "SEO描述"
        varchar seo_keywords "SEO关键词"
        tinyint status "状态"
        int sort "排序"
        tinyint is_published "已发布"
        timestamp published_at "发布时间"
        int version "版本"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
        bigint created_by "创建人"
        bigint updated_by "更新人"
    }

    navigations {
        bigint id PK "导航ID"
        bigint tenant_id FK "租户ID"
        varchar name "名称"
        varchar position "位置"
        tinyint status "状态"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    nav_items {
        bigint id PK "导航项ID"
        bigint nav_id FK "导航ID"
        bigint parent_id "父ID"
        varchar name "名称"
        varchar link "链接"
        varchar type "类型"
        bigint target_id "目标ID"
        int sort "排序"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    decorations {
        bigint id PK "装修ID"
        bigint tenant_id FK "租户ID"
        bigint page_id FK "页面ID"
        varchar block_type "块类型"
        text block_config "块配置"
        int sort_order "排序"
        tinyint is_active "激活"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    page_versions {
        bigint id PK "版本ID"
        bigint tenant_id FK "租户ID"
        bigint page_id FK "页面ID"
        int version "版本号"
        text blocks "块快照"
        bigint created_by "创建人"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    seo_configs {
        bigint id PK "SEO配置ID"
        bigint tenant_id FK "租户ID"
        varchar page_type "页面类型"
        bigint page_id "页面ID"
        varchar title "标题"
        text description "描述"
        varchar keywords "关键词"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    theme_audit_logs {
        bigint id PK "日志ID"
        bigint tenant_id FK "租户ID"
        varchar action "动作"
        bigint theme_id FK "主题ID"
        varchar theme_name "主题名"
        varchar theme_code "主题代码"
        text old_config "旧配置"
        text new_config "新配置"
        bigint user_id "用户ID (关联 admin_users)"
        varchar user_name "用户名"
        varchar ip_address "IP"
        varchar user_agent "UserAgent"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
        timestamp deleted_at "删除时间"
    }

    %% ==================== SHOP SETTINGS DOMAIN ====================
    shop_settings {
        bigint id PK "设置ID"
        bigint tenant_id FK "租户ID"
        varchar name "店铺名称"
        varchar code "代码"
        varchar logo "Logo"
        varchar description "描述"
        varchar contact_name "联系人"
        varchar contact_phone "电话"
        varchar contact_email "邮箱"
        varchar address "地址"
        varchar domain "域名"
        varchar custom_domain "自定义域名"
        varchar primary_color "主题色"
        varchar secondary_color "辅助色"
        varchar favicon "Favicon"
        varchar default_currency "默认货币"
        varchar default_language "默认语言"
        varchar timezone "时区"
        tinyint status "状态"
        tinyint plan "套餐"
        varchar expire_at "过期时间"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }

    shop_business_hours {
        bigint id PK "营业时间ID"
        bigint shop_id FK "店铺ID"
        tinyint day_of_week "星期"
        varchar open_time "开门时间"
        varchar close_time "关门时间"
        tinyint is_closed "休息"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }

    shop_notification_settings {
        bigint id PK "通知设置ID"
        bigint shop_id FK "店铺ID"
        tinyint order_created "订单创建"
        tinyint order_paid "订单支付"
        tinyint order_shipped "订单发货"
        tinyint order_cancelled "订单取消"
        tinyint low_stock_alert "低库存"
        int low_stock_threshold "阈值"
        tinyint refund_requested "退款申请"
        tinyint new_review "新评价"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }

    shop_payment_settings {
        bigint id PK "支付设置ID"
        bigint shop_id FK "店铺ID"
        tinyint stripe_enabled "Stripe启用"
        varchar stripe_public_key "公钥"
        varchar stripe_secret_key "私钥"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }

    shop_shipping_settings {
        bigint id PK "运费设置ID"
        bigint shop_id FK "店铺ID"
        decimal free_shipping_threshold "免运门槛"
        decimal default_shipping_fee "默认运费"
        varchar currency "货币"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }

    %% ==================== RELATIONSHIPS ====================

    %% Tenant relationships
    tenants ||--o{ users : "租户拥有用户"
    tenants ||--o{ admin_users : "租户拥有管理员"
    tenants ||--o{ roles : "租户拥有角色"
    tenants ||--o{ products : "租户拥有商品"
    tenants ||--o{ categories : "租户拥有分类"
    tenants ||--o{ brands : "租户拥有品牌"
    tenants ||--o{ orders : "租户拥有订单"
    tenants ||--o{ carts : "租户拥有购物车"
    tenants ||--o{ promotions : "租户拥有促销"
    tenants ||--o{ coupons : "租户拥有优惠券"
    tenants ||--o{ reviews : "租户拥有评价"
    tenants ||--o{ markets : "租户拥有市场"
    tenants ||--o{ warehouses : "租户拥有仓库"
    tenants ||--o{ points_accounts : "租户拥有积分账户"
    tenants ||--o{ earn_rules : "租户拥有积分规则"
    tenants ||--o{ redeem_rules : "租户拥有兑换规则"
    tenants ||--o{ shops : "租户拥有店铺"
    tenants ||--o{ shop_settings : "租户拥有店铺设置"
    tenants ||--o{ themes : "租户拥有主题"
    tenants ||--o{ pages : "租户拥有页面"
    tenants ||--o{ navigations : "租户拥有导航"
    tenants ||--o{ decorations : "租户拥有装修"
    tenants ||--o{ shipments : "租户拥有发货"
    tenants ||--o{ refunds : "租户拥有退款"
    tenants ||--o{ order_payments : "租户拥有支付"
    tenants ||--o{ payment_transactions : "租户拥有交易"
    tenants ||--o{ payment_refunds : "租户拥有退款"

    %% User relationships
    users ||--o{ user_addresses : "用户拥有地址"
    users ||--o{ orders : "用户下单"
    users ||--o{ carts : "用户拥有购物车"
    users ||--o{ cart_items : "用户拥有购物车项"
    users ||--o{ user_coupons : "用户领取优惠券"
    users ||--o{ points_accounts : "用户拥有积分账户"
    users ||--o{ points_transactions : "用户积分变动"
    users ||--o{ points_redemptions : "用户兑换积分"
    users ||--o{ reviews : "用户评价"
    users ||--o{ admin_users : "用户是管理员"

    %% Admin User relationships
    admin_users ||--o{ user_roles : "管理员拥有角色"
    admin_users ||--o{ review_replies : "管理员回复评价"

    %% Role relationships
    roles ||--o{ user_roles : "角色分配给用户"
    roles ||--o{ role_permissions : "角色拥有权限"

    %% Permission relationships
    permissions ||--o{ role_permissions : "权限分配给角色"

    %% Product relationships
    categories ||--o{ categories : "分类自关联"
    categories ||--o{ products : "分类包含商品"
    brands ||--o{ products : "品牌包含商品"
    products ||--o{ skus : "商品拥有SKU"
    products ||--o{ order_items : "商品在订单中"
    products ||--o{ cart_items : "商品在购物车中"
    products ||--o{ product_markets : "商品面向市场"
    products ||--o{ product_localizations : "商品多语言"
    products ||--o{ reviews : "商品被评价"
    products ||--o{ review_stats : "商品评价统计"
    products ||--o{ inventory_logs : "商品库存变动"

    skus ||--o{ order_items : "SKU在订单中"
    skus ||--o{ cart_items : "SKU在购物车中"
    skus ||--o{ product_markets : "SKU面向市场"

    markets ||--o{ product_markets : "市场包含商品"
    markets ||--o{ category_markets : "市场包含分类"
    markets ||--o{ brand_markets : "市场包含品牌"

    warehouses ||--o{ warehouse_inventories : "仓库有库存"
    warehouse_inventories }o--|| warehouses : "库存属于仓库"

    inventory_logs }o--|| products : "日志记录商品"
    inventory_logs }o--|| warehouses : "日志记录仓库"

    %% Order relationships
    orders ||--o{ order_items : "订单包含商品"
    orders ||--o{ shipments : "订单发货"
    orders ||--o{ refunds : "订单退款"
    orders ||--o{ order_payments : "订单支付"
    orders ||--o{ promotion_usage : "订单使用促销"

    order_items ||--o{ shipment_items : "订单商品发货"

    carts ||--o{ cart_items : "购物车包含商品"

    %% Payment relationships
    orders ||--o{ order_payments : "订单支付记录"
    order_payments ||--o{ payment_transactions : "支付包含交易"
    order_payments ||--o{ payment_refunds : "支付包含退款"

    refunds ||--o{ payment_refunds : "payment_refunds.fulfillment_refund_id → refunds.id"

    %% Promotion relationships
    promotions ||--o{ promotion_rules : "促销拥有规则"
    promotions ||--o{ promotion_usage : "促销被使用"
    promotion_rules ||--o{ promotion_usage : "规则被使用"

    coupons ||--o{ user_coupons : "优惠券被用户领取"
    coupons ||--o{ promotion_usage : "优惠券被使用"
    coupons ||--o{ redeem_rules : "优惠券可兑换"

    user_coupons ||--o{ points_redemptions : "用户优惠券兑换"

    %% Points relationships
    earn_rules ||--o{ points_transactions : "积分规则产生变动"
    redeem_rules ||--o{ points_redemptions : "兑换规则被使用"
    points_accounts ||--o{ points_transactions : "账户产生变动"

    points_transactions }o--|| points_accounts : "变动属于账户"
    points_transactions }o--|| users : "变动属于用户"

    points_redemptions }o--|| redeem_rules : "兑换依据规则"
    points_redemptions }o--|| coupons : "兑换获得优惠券"
    points_redemptions }o--|| user_coupons : "兑换获得用户优惠券"

    %% Review relationships
    reviews ||--o{ review_replies : "评价被回复"
    reviews }o--|| orders : "reviews.order_id → orders.id (应用层关联)"

    review_stats }o--|| products : "统计属于商品"

    %% Storefront relationships
    shops ||--o{ themes : "主题被店铺使用 (current_theme_id → themes.id)"
    themes ||--o{ theme_audit_logs : "主题有审计日志"

    pages ||--o{ decorations : "页面拥有装修"
    pages ||--o{ page_versions : "页面有版本"
    pages ||--o{ seo_configs : "SEO配置属于页面"

    navigations ||--o{ nav_items : "导航拥有项"

    shop_settings ||--o{ shop_business_hours : "设置包含营业时间"
    shop_settings ||--o{ shop_notification_settings : "设置包含通知"
    shop_settings ||--o{ shop_payment_settings : "设置包含支付"
    shop_settings ||--o{ shop_shipping_settings : "设置包含运费"

    %% 注意: shops 和 shop_settings 是两个独立表，都通过 tenant_id 关联到 tenants
    %% shop_settings 的子表 (shop_business_hours 等) 关联到 shop_settings.id，而非 shops.id

    %% Shipping relationships
    shipping_templates ||--o{ shipping_zones : "模板拥有区域"
    shipping_templates ||--o{ shipping_template_mappings : "模板拥有映射"
    shipping_zones ||--o{ shipping_zone_regions : "区域包含城市"
```

## Entity Relationship Summary

### Domain: User (用户与权限)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `tenants` | 租户表 | Root entity for multi-tenancy |
| `users` | C端用户表 | tenant_id → tenants |
| `admin_users` | 后台管理员表 | tenant_id → tenants |
| `roles` | 角色表 | tenant_id → tenants |
| `permissions` | 权限表 | Self-referencing hierarchy (parent_id) |
| `user_roles` | 用户角色关联 | user_id → admin_users, role_id → roles |
| `role_permissions` | 角色权限关联 | role_id → roles, permission_id → permissions |
| `user_addresses` | 用户收货地址 | tenant_id → tenants, user_id → users |

### Domain: Product (商品目录)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `categories` | 分类表 | tenant_id → tenants, self-referencing (parent_id) |
| `brands` | 品牌表 | tenant_id → tenants |
| `products` | 商品表 | tenant_id → tenants, category_id → categories, brand_id (关联 brands.id, 非外键约束) |
| `skus` | SKU表 | tenant_id → tenants, product_id → products |
| `markets` | 市场表 | tenant_id → tenants |
| `product_markets` | 商品市场关联 | product_id → products, variant_id (关联 skus.id, 非外键约束), market_id → markets |
| `product_localizations` | 产品多语言 | product_id → products |
| `category_markets` | 分类市场可见性 | category_id → categories, market_id → markets |
| `brand_markets` | 品牌市场可见性 | brand_id → brands, market_id → markets |
| `warehouses` | 仓库表 | tenant_id → tenants |
| `warehouse_inventories` | 仓库库存 | warehouse_id → warehouses, sku_code (非外键约束) |
| `inventory_logs` | 库存变更日志 | product_id → products, warehouse_id → warehouses, sku_code (非外键约束) |

### Domain: Order (订单)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `orders` | 订单表 | tenant_id → tenants, user_id → users |
| `order_items` | 订单商品表 | order_id → orders, product_id → products, sku_id → skus |
| `carts` | 购物车表 | tenant_id → tenants, user_id → users |
| `cart_items` | 购物车商品表 | cart_id → carts, product_id → products, sku_id → skus |

### Domain: Payment (支付)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `order_payments` | 订单支付表 | tenant_id → tenants, order_id → orders |
| `payment_transactions` | 支付交易记录 | order_id → orders, payment_id → order_payments |
| `payment_refunds` | 支付退款表 | order_id → orders, payment_id → order_payments, fulfillment_refund_id → fulfillment.refunds |
| `webhook_events` | Webhook事件表 | tenant_id → tenants |

### Domain: Fulfillment (履约)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `shipments` | 发货表 | tenant_id → tenants, order_id → orders |
| `shipment_items` | 发货商品表 | shipment_id → shipments, order_item_id → order_items |
| `refunds` | 退款表 | tenant_id → tenants, order_id → orders, user_id → users |
| `shipping_templates` | 运费模板 | tenant_id → tenants |
| `shipping_zones` | 运费区域 | template_id → shipping_templates |
| `shipping_template_mappings` | 运费模板映射 | template_id → shipping_templates |
| `shipping_zone_regions` | 运费区域城市 | zone_id → shipping_zones |

### Domain: Promotion (促销)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `promotions` | 促销活动表 | tenant_id → tenants |
| `promotion_rules` | 促销规则表 | promotion_id → promotions |
| `promotion_usage` | 促销使用记录 | promotion_id → promotions, rule_id → promotion_rules, order_id → orders, user_id → users, coupon_id (非外键约束) |
| `coupons` | 优惠券表 | tenant_id → tenants |
| `user_coupons` | 用户优惠券表 | user_id → users, coupon_id → coupons, order_id → orders |

### Domain: Points (积分)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `earn_rules` | 积分获取规则 | tenant_id → tenants |
| `redeem_rules` | 积分兑换规则 | tenant_id → tenants, coupon_id → coupons.id |
| `points_accounts` | 积分账户 | tenant_id → tenants, user_id → users |
| `points_transactions` | 积分交易记录 | account_id → points_accounts, user_id → users |
| `points_redemptions` | 积分兑换记录 | user_id → users, redeem_rule_id → redeem_rules, coupon_id → coupons, user_coupon_id (关联 user_coupons) |

### Domain: Review (评价)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `reviews` | 评价表 | tenant_id → tenants, order_id (关联 orders.id, 非外键约束), product_id → products, user_id → users |
| `review_replies` | 评价回复表 | review_id → reviews.id (UNIQUE), admin_id (关联 admin_users, 非外键约束) |
| `review_stats` | 评价统计表 | tenant_id → tenants, product_id → products |

### Domain: Storefront (店铺装修)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `shops` | 店铺表 | tenant_id → tenants, current_theme_id (关联 themes.id, 非外键约束) |
| `themes` | 主题表 | tenant_id → tenants |
| `pages` | 页面表 | tenant_id → tenants |
| `navigations` | 导航表 | tenant_id → tenants |
| `nav_items` | 导航项表 | nav_id → navigations, self-referencing (parent_id) |
| `decorations` | 页面装修表 | tenant_id → tenants, page_id → pages |
| `page_versions` | 页面版本表 | tenant_id → tenants, page_id → pages |
| `seo_configs` | SEO配置表 | tenant_id → tenants, page_id → pages |
| `theme_audit_logs` | 主题审计日志 | tenant_id → tenants, theme_id → themes, user_id (关联 admin_users) |

### Domain: Shop (店铺设置)
| Table | Description | Key Relationships |
|-------|-------------|-------------------|
| `shop_settings` | 店铺设置表 | tenant_id → tenants (UNIQUE), shops 和 shop_settings 是两个独立表 |
| `shop_business_hours` | 营业时间表 | shop_id → shop_settings.id |
| `shop_notification_settings` | 通知设置表 | shop_id → shop_settings.id |
| `shop_payment_settings` | 支付设置表 | shop_id → shop_settings.id |
| `shop_shipping_settings` | 运费设置表 | shop_id → shop_settings.id |

**重要**: `shops` 表和 `shop_settings` 表是两个独立的表，都通过 `tenant_id` 关联到 `tenants`。`shop_settings` 的子表（如 shop_business_hours）通过 `shop_id` 关联到 `shop_settings`，而不是 `shops`。

## Key Design Patterns

1. **Multi-tenancy**: All tables have `tenant_id` as a foreign key to `tenants`
2. **Soft Delete**: All tables include `deleted_at` timestamp field
3. **Audit Trail**: Tables include `created_at`, `updated_at`, `created_by`, `updated_by`
4. **Currency/Money**: Monetary values use `DECIMAL(19,4)` for precision
5. **JSON Fields**: Flexible data stored as JSON (tags, images, scope_ids, etc.)
6. **Hierarchical Data**: Categories and nav_items use self-referencing `parent_id`
