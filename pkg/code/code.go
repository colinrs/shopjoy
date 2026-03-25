package code

import "net/http"

var (
	OKErr              = &Err{HTTPCode: http.StatusOK, Code: 0, Msg: "success"}
	ErrParam           = &Err{HTTPCode: http.StatusOK, Code: 10003, Msg: "参数有误"}
	UnknownErr         = &Err{HTTPCode: http.StatusOK, Code: 10004, Msg: "unknown error"}
	ErrValidation      = &Err{HTTPCode: http.StatusOK, Code: 20001, Msg: "Validation failed."}
	ErrDatabase        = &Err{HTTPCode: http.StatusOK, Code: 20002, Msg: "Database error."}
	BizTagAlreadyExist = &Err{HTTPCode: http.StatusOK, Code: 20003, Msg: "业务已存在"}
	BizTagNotExist     = &Err{HTTPCode: http.StatusOK, Code: 20004, Msg: "业务不存在"}
	EtcdKeyNotExist    = &Err{HTTPCode: http.StatusOK, Code: 20005, Msg: "Etcd key 不存在"}
	HTTPClientErr      = &Err{HTTPCode: http.StatusOK, Code: 20006, Msg: "http client error"}

	// ErrUnauthorized Standard HTTP errors with proper status codes
	ErrUnauthorized       = &Err{HTTPCode: http.StatusUnauthorized, Code: 40100, Msg: "未授权，请先登录"}
	ErrTokenExpired       = &Err{HTTPCode: http.StatusUnauthorized, Code: 40101, Msg: "Token 已过期，请重新登录"}
	ErrTokenInvalid       = &Err{HTTPCode: http.StatusUnauthorized, Code: 40102, Msg: "无效的 Token"}
	ErrForbidden          = &Err{HTTPCode: http.StatusForbidden, Code: 40300, Msg: "没有权限访问该资源"}
	ErrNotFound           = &Err{HTTPCode: http.StatusNotFound, Code: 40400, Msg: "资源不存在"}
	ErrMethodNotAllowed   = &Err{HTTPCode: http.StatusMethodNotAllowed, Code: 40500, Msg: "请求方法不允许"}
	ErrTooManyRequests    = &Err{HTTPCode: http.StatusTooManyRequests, Code: 42900, Msg: "请求过于频繁，请稍后再试"}
	ErrInternalServer     = &Err{HTTPCode: http.StatusInternalServerError, Code: 50000, Msg: "服务内部错误"}
	ErrServiceUnavailable = &Err{HTTPCode: http.StatusServiceUnavailable, Code: 50300, Msg: "服务暂时不可用"}

	// ErrAdminInvalidEmail ==================== Admin User Module (10xxx) ====================
	ErrAdminInvalidEmail     = &Err{HTTPCode: http.StatusBadRequest, Code: 10001, Msg: "invalid email format"}
	ErrAdminInvalidPhone     = &Err{HTTPCode: http.StatusBadRequest, Code: 10002, Msg: "invalid phone format"}
	ErrAdminPasswordTooWeak  = &Err{HTTPCode: http.StatusBadRequest, Code: 10003, Msg: "password too weak"}
	ErrAdminUserNotFound     = &Err{HTTPCode: http.StatusNotFound, Code: 10004, Msg: "admin user not found"}
	ErrAdminDuplicateUser    = &Err{HTTPCode: http.StatusConflict, Code: 10005, Msg: "duplicate admin user"}
	ErrAdminWrongPassword    = &Err{HTTPCode: http.StatusUnauthorized, Code: 10006, Msg: "wrong password"}
	ErrAdminCannotDeleteSelf = &Err{HTTPCode: http.StatusBadRequest, Code: 10007, Msg: "cannot delete yourself"}
	ErrAdminAlreadyDeleted   = &Err{HTTPCode: http.StatusBadRequest, Code: 10008, Msg: "user already deleted"}
	ErrAdminAccountDisabled  = &Err{HTTPCode: http.StatusForbidden, Code: 10009, Msg: "account disabled or deleted"}
	ErrAdminPasswordMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 10010, Msg: "passwords do not match"}
	ErrAdminPermissionDenied = &Err{HTTPCode: http.StatusForbidden, Code: 10011, Msg: "permission denied"}

	// Admin additional errors
	ErrAdminCannotDisableSelf = &Err{HTTPCode: http.StatusForbidden, Code: 10012, Msg: "cannot disable yourself"}
	ErrAdminInvalidType       = &Err{HTTPCode: http.StatusBadRequest, Code: 10013, Msg: "invalid admin type"}
	ErrAdminUsernameExists    = &Err{HTTPCode: http.StatusBadRequest, Code: 10014, Msg: "username already exists"}
	ErrAdminEmailExists       = &Err{HTTPCode: http.StatusBadRequest, Code: 10015, Msg: "email already exists"}
	ErrAdminMainAccountExists = &Err{HTTPCode: http.StatusBadRequest, Code: 10016, Msg: "tenant already has a main account"}

	// ErrUserInvalidEmail ==================== User Module (11xxx) ====================
	ErrUserInvalidEmail     = &Err{HTTPCode: http.StatusBadRequest, Code: 11001, Msg: "invalid email format"}
	ErrUserInvalidPhone     = &Err{HTTPCode: http.StatusBadRequest, Code: 11002, Msg: "invalid phone format"}
	ErrUserPasswordTooWeak  = &Err{HTTPCode: http.StatusBadRequest, Code: 11003, Msg: "password too weak"}
	ErrUserNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 11004, Msg: "user not found"}
	ErrUserDuplicateUser    = &Err{HTTPCode: http.StatusConflict, Code: 11005, Msg: "duplicate user"}
	ErrUserWrongPassword    = &Err{HTTPCode: http.StatusUnauthorized, Code: 11006, Msg: "wrong password"}
	ErrUserAlreadyDeleted   = &Err{HTTPCode: http.StatusBadRequest, Code: 11007, Msg: "user already deleted"}
	ErrUserPasswordMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 11008, Msg: "passwords do not match"}

	// User additional errors
	ErrUserSuspended          = &Err{HTTPCode: http.StatusBadRequest, Code: 11009, Msg: "user already suspended"}
	ErrUserCannotSuspendSelf  = &Err{HTTPCode: http.StatusForbidden, Code: 11010, Msg: "cannot suspend yourself"}
	ErrUserCannotDeleteSelf   = &Err{HTTPCode: http.StatusForbidden, Code: 11011, Msg: "cannot delete yourself"}
	ErrAddressNotFound        = &Err{HTTPCode: http.StatusNotFound, Code: 11012, Msg: "address not found"}
	ErrUserExportLimitExceed  = &Err{HTTPCode: http.StatusBadRequest, Code: 11013, Msg: "export limit exceeded, maximum 10000 records"}

	// ErrProductEmptyName ==================== Product Module (30xxx) ====================
	ErrProductEmptyName               = &Err{HTTPCode: http.StatusBadRequest, Code: 30001, Msg: "商品名称不能为空"}
	ErrProductInvalidPrice            = &Err{HTTPCode: http.StatusBadRequest, Code: 30002, Msg: "商品价格必须大于0"}
	ErrProductCurrencyMismatch        = &Err{HTTPCode: http.StatusBadRequest, Code: 30003, Msg: "币种不匹配"}
	ErrProductInsufficientAmount      = &Err{HTTPCode: http.StatusBadRequest, Code: 30004, Msg: "金额不足"}
	ErrProductDeleted                 = &Err{HTTPCode: http.StatusBadRequest, Code: 30005, Msg: "商品已删除"}
	ErrProductInvalidStatusTransition = &Err{HTTPCode: http.StatusBadRequest, Code: 30006, Msg: "无效的状态转换"}
	ErrProductNoStock                 = &Err{HTTPCode: http.StatusBadRequest, Code: 30007, Msg: "库存不能为0"}
	ErrProductNegativeStock           = &Err{HTTPCode: http.StatusBadRequest, Code: 30008, Msg: "库存不能为负数"}
	ErrProductNotOnSale               = &Err{HTTPCode: http.StatusBadRequest, Code: 30009, Msg: "商品未上架"}
	ErrProductInvalidQuantity         = &Err{HTTPCode: http.StatusBadRequest, Code: 30010, Msg: "无效的数量"}
	ErrProductInsufficientStock       = &Err{HTTPCode: http.StatusBadRequest, Code: 30011, Msg: "库存不足"}
	ErrProductNotFound                = &Err{HTTPCode: http.StatusNotFound, Code: 30012, Msg: "商品不存在"}
	ErrProductInvalidID               = &Err{HTTPCode: http.StatusBadRequest, Code: 30013, Msg: "invalid product id"}

	//  ErrCategoryNotFound Category errors (30xxx - 31xxx)
	ErrCategoryNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 30101, Msg: "category not found"}
	ErrCategoryDuplicate        = &Err{HTTPCode: http.StatusConflict, Code: 30102, Msg: "duplicate category"}
	ErrCategoryInvalid          = &Err{HTTPCode: http.StatusBadRequest, Code: 30103, Msg: "invalid category"}
	ErrCategoryHasChildren      = &Err{HTTPCode: http.StatusBadRequest, Code: 30104, Msg: "category has children"}
	ErrCategoryHasProducts      = &Err{HTTPCode: http.StatusBadRequest, Code: 30105, Msg: "category has products"}
	ErrCategoryMaxLevelExceeded = &Err{HTTPCode: http.StatusBadRequest, Code: 30106, Msg: "category level cannot exceed 3"}

	// ErrOrderNotFound ==================== Order Module (40xxx) ====================
	ErrOrderNotFound          = &Err{HTTPCode: http.StatusNotFound, Code: 40001, Msg: "order not found"}
	ErrOrderInvalidStatus     = &Err{HTTPCode: http.StatusBadRequest, Code: 40002, Msg: "invalid order status"}
	ErrOrderAlreadyPaid       = &Err{HTTPCode: http.StatusBadRequest, Code: 40003, Msg: "order already paid"}
	ErrOrderNotPaid           = &Err{HTTPCode: http.StatusBadRequest, Code: 40004, Msg: "order not paid"}
	ErrOrderExpired           = &Err{HTTPCode: http.StatusBadRequest, Code: 40005, Msg: "order expired"}
	ErrOrderInsufficientStock = &Err{HTTPCode: http.StatusBadRequest, Code: 40006, Msg: "insufficient stock"}
	ErrOrderInvalidAmount     = &Err{HTTPCode: http.StatusBadRequest, Code: 40007, Msg: "invalid amount"}
	ErrOrderCartEmpty         = &Err{HTTPCode: http.StatusBadRequest, Code: 40008, Msg: "cart is empty"}
	// Order adjustment errors
	ErrOrderCannotAdjustPrice    = &Err{HTTPCode: http.StatusBadRequest, Code: 40009, Msg: "当前状态无法改价"}
	ErrOrderAdjustAmountExceed   = &Err{HTTPCode: http.StatusBadRequest, Code: 40010, Msg: "改价金额超出限制"}
	ErrOrderExportLimitExceed    = &Err{HTTPCode: http.StatusBadRequest, Code: 40011, Msg: "导出数量超出限制"}
	ErrOrderAdjustReasonRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 40012, Msg: "改价原因不能为空"}
	ErrOrderVersionConflict      = &Err{HTTPCode: http.StatusConflict, Code: 40013, Msg: "订单已被修改，请刷新后重试"}
	// Order cancellation and reminder errors
	ErrOrderCannotCancel         = &Err{HTTPCode: http.StatusBadRequest, Code: 40014, Msg: "order cannot be cancelled in current status"}
	ErrOrderCancelReasonRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 40015, Msg: "cancel reason is required"}
	ErrPaymentReminderSent       = &Err{HTTPCode: http.StatusTooManyRequests, Code: 40016, Msg: "payment reminder already sent recently"}
	ErrOrderAlreadyPaidForRemind = &Err{HTTPCode: http.StatusBadRequest, Code: 40017, Msg: "order already paid, cannot send reminder"}
	ErrOrderCannotRemind         = &Err{HTTPCode: http.StatusBadRequest, Code: 40018, Msg: "order cannot be reminded in current status"}

	// ErrPaymentNotFound ==================== Payment Module (50xxx) ====================
	ErrPaymentNotFound              = &Err{HTTPCode: http.StatusNotFound, Code: 50001, Msg: "payment not found"}
	ErrPaymentInvalidAmount         = &Err{HTTPCode: http.StatusBadRequest, Code: 50002, Msg: "invalid payment amount"}
	ErrPaymentFailed                = &Err{HTTPCode: http.StatusPaymentRequired, Code: 50003, Msg: "payment failed"}
	ErrPaymentAlreadyPaid           = &Err{HTTPCode: http.StatusBadRequest, Code: 50004, Msg: "payment already completed"}
	ErrPaymentExpired               = &Err{HTTPCode: http.StatusBadRequest, Code: 50005, Msg: "payment expired"}
	ErrPaymentOrderNotPaid          = &Err{HTTPCode: http.StatusBadRequest, Code: 50006, Msg: "order not paid, cannot refund"}
	ErrPaymentRefundAmountExceeded = &Err{HTTPCode: http.StatusBadRequest, Code: 50007, Msg: "refund amount exceeds refundable"}
	ErrPaymentRefundReasonRequired  = &Err{HTTPCode: http.StatusBadRequest, Code: 50008, Msg: "refund reason is required"}
	ErrChannelRefundFailed          = &Err{HTTPCode: http.StatusInternalServerError, Code: 50009, Msg: "channel refund failed"}
	ErrPaymentOrderNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 50010, Msg: "order not found"}
	ErrPaymentOrderAlreadyRefunded  = &Err{HTTPCode: http.StatusBadRequest, Code: 50011, Msg: "order already fully refunded"}
	ErrTransactionNotFound           = &Err{HTTPCode: http.StatusNotFound, Code: 50012, Msg: "transaction not found"}
	ErrPaymentChannelUnavailable     = &Err{HTTPCode: http.StatusServiceUnavailable, Code: 50013, Msg: "payment channel unavailable"}
	ErrRefundNotSupported           = &Err{HTTPCode: http.StatusBadRequest, Code: 50014, Msg: "refund not supported for this channel"}
	ErrCurrencyNotSupported         = &Err{HTTPCode: http.StatusBadRequest, Code: 50015, Msg: "currency not supported by channel"}
	ErrPaymentRefundNotFound        = &Err{HTTPCode: http.StatusNotFound, Code: 50016, Msg: "payment refund not found"}
	ErrPaymentRequiresAction        = &Err{HTTPCode: http.StatusAccepted, Code: 50017, Msg: "payment requires additional action"}
	ErrRefundCurrencyMismatch       = &Err{HTTPCode: http.StatusBadRequest, Code: 50018, Msg: "refund currency must match payment currency"}
	ErrIdempotencyKeyConflict        = &Err{HTTPCode: http.StatusConflict, Code: 50019, Msg: "duplicate idempotency key"}
	ErrDisputeCreated                = &Err{HTTPCode: http.StatusConflict, Code: 50020, Msg: "dispute created for this charge"}

	// ErrCartItemNotFound ==================== Cart Module (60xxx) ====================
	ErrCartItemNotFound    = &Err{HTTPCode: http.StatusNotFound, Code: 60001, Msg: "cart item not found"}
	ErrCartInvalidQuantity = &Err{HTTPCode: http.StatusBadRequest, Code: 60002, Msg: "invalid quantity"}
	ErrCartEmpty           = &Err{HTTPCode: http.StatusBadRequest, Code: 60003, Msg: "cart is empty"}

	// ErrCouponNotFound ==================== Coupon Module (70xxx) ====================
	ErrCouponNotFound       = &Err{HTTPCode: http.StatusNotFound, Code: 70001, Msg: "coupon not found"}
	ErrCouponExpired        = &Err{HTTPCode: http.StatusBadRequest, Code: 70002, Msg: "coupon expired"}
	ErrCouponUsedUp         = &Err{HTTPCode: http.StatusBadRequest, Code: 70003, Msg: "coupon used up"}
	ErrCouponNotStarted     = &Err{HTTPCode: http.StatusBadRequest, Code: 70004, Msg: "coupon not started"}
	ErrCouponAlreadyUsed    = &Err{HTTPCode: http.StatusBadRequest, Code: 70005, Msg: "coupon already used"}
	ErrCouponInvalidCode    = &Err{HTTPCode: http.StatusBadRequest, Code: 70006, Msg: "invalid coupon code"}
	ErrCouponAmountBelowMin = &Err{HTTPCode: http.StatusBadRequest, Code: 70007, Msg: "cart amount below minimum"}

	// Coupon additional errors
	ErrCouponUserLimitReached = &Err{HTTPCode: http.StatusBadRequest, Code: 70008, Msg: "user coupon usage limit reached"}
	ErrCouponCurrencyMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 70009, Msg: "coupon currency mismatch"}
	ErrCouponScopeInvalid     = &Err{HTTPCode: http.StatusBadRequest, Code: 70010, Msg: "invalid coupon scope"}
	ErrCouponTypeInvalid      = &Err{HTTPCode: http.StatusBadRequest, Code: 70011, Msg: "invalid coupon type"}
	ErrCouponCannotDelete     = &Err{HTTPCode: http.StatusBadRequest, Code: 70012, Msg: "cannot delete active coupon"}
	ErrCouponNameRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 70013, Msg: "coupon name is required"}
	ErrCouponValueRequired    = &Err{HTTPCode: http.StatusBadRequest, Code: 70014, Msg: "coupon value must be positive"}
	ErrCouponTimeRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 70015, Msg: "coupon start time and end time are required"}
	ErrCouponInvalidTimeRange = &Err{HTTPCode: http.StatusBadRequest, Code: 70016, Msg: "coupon start time must be before end time"}
	ErrCouponNotActive        = &Err{HTTPCode: http.StatusBadRequest, Code: 70017, Msg: "coupon is not active"}

	// UserCoupon errors (70xxx - 701xx)
	ErrUserCouponNotFound    = &Err{HTTPCode: http.StatusNotFound, Code: 70101, Msg: "user coupon not found"}
	ErrUserCouponExpired     = &Err{HTTPCode: http.StatusBadRequest, Code: 70102, Msg: "user coupon expired"}
	ErrUserCouponAlreadyUsed = &Err{HTTPCode: http.StatusBadRequest, Code: 70103, Msg: "user coupon already used"}

	// ErrPromotionNotFound ==================== Promotion Module (80xxx) ====================
	ErrPromotionNotFound   = &Err{HTTPCode: http.StatusNotFound, Code: 80001, Msg: "promotion not found"}
	ErrPromotionInvalid    = &Err{HTTPCode: http.StatusBadRequest, Code: 80002, Msg: "invalid promotion"}
	ErrPromotionExpired    = &Err{HTTPCode: http.StatusBadRequest, Code: 80003, Msg: "promotion expired"}
	ErrPromotionNotStarted = &Err{HTTPCode: http.StatusBadRequest, Code: 80004, Msg: "promotion not started"}

	// Promotion additional errors
	ErrPromotionNotActive        = &Err{HTTPCode: http.StatusBadRequest, Code: 80005, Msg: "promotion is not active"}
	ErrPromotionCurrencyMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 80006, Msg: "promotion currency mismatch"}
	ErrPromotionRuleNotFound     = &Err{HTTPCode: http.StatusNotFound, Code: 80007, Msg: "promotion rule not found"}
	ErrPromotionRuleInvalid      = &Err{HTTPCode: http.StatusBadRequest, Code: 80008, Msg: "invalid promotion rule"}
	ErrPromotionCannotDelete     = &Err{HTTPCode: http.StatusBadRequest, Code: 80009, Msg: "cannot delete active promotion"}
	ErrPromotionScopeInvalid     = &Err{HTTPCode: http.StatusBadRequest, Code: 80010, Msg: "invalid promotion scope"}
	ErrPromotionTypeInvalid      = &Err{HTTPCode: http.StatusBadRequest, Code: 80011, Msg: "invalid promotion type"}
	ErrPromotionNameRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 80012, Msg: "promotion name is required"}
	ErrPromotionCurrencyRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 80013, Msg: "promotion currency is required"}
	ErrPromotionTimeRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 80014, Msg: "promotion start time and end time are required"}
	ErrPromotionInvalidTimeRange = &Err{HTTPCode: http.StatusBadRequest, Code: 80015, Msg: "promotion start time must be before end time"}

	// ErrTenantNotFound ==================== Tenant Module (90xxx) ====================
	ErrTenantNotFound             = &Err{HTTPCode: http.StatusNotFound, Code: 90001, Msg: "tenant not found"}
	ErrTenantDuplicate            = &Err{HTTPCode: http.StatusConflict, Code: 90002, Msg: "duplicate tenant"}
	ErrTenantInvalidDomain        = &Err{HTTPCode: http.StatusBadRequest, Code: 90003, Msg: "invalid domain"}
	ErrTenantInactive             = &Err{HTTPCode: http.StatusForbidden, Code: 90004, Msg: "tenant is inactive"}
	ErrTenantCannotSuspendExpired = &Err{HTTPCode: http.StatusBadRequest, Code: 90005, Msg: "cannot suspend expired tenant"}
	ErrTenantInvalidID            = &Err{HTTPCode: http.StatusBadRequest, Code: 90006, Msg: "invalid tenant id"}
	ErrTenantNameRequired         = &Err{HTTPCode: http.StatusBadRequest, Code: 90007, Msg: "tenant name is required"}
	ErrTenantCodeRequired         = &Err{HTTPCode: http.StatusBadRequest, Code: 90008, Msg: "tenant code is required"}

	// ErrRoleNotFound ==================== Role Module (100xxx) ====================
	ErrRoleNotFound           = &Err{HTTPCode: http.StatusNotFound, Code: 100001, Msg: "role not found"}
	ErrRoleDuplicate          = &Err{HTTPCode: http.StatusConflict, Code: 100002, Msg: "duplicate role"}
	ErrRoleInvalid            = &Err{HTTPCode: http.StatusBadRequest, Code: 100003, Msg: "invalid role"}
	ErrRoleCannotModifySystem = &Err{HTTPCode: http.StatusBadRequest, Code: 100004, Msg: "cannot modify system role"}
	ErrRoleInUse              = &Err{HTTPCode: http.StatusBadRequest, Code: 100005, Msg: "role is in use by users"}

	// ErrShopNotFound ==================== Storefront Module (110xxx) ====================
	ErrShopNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 110001, Msg: "shop not found"}
	ErrShopInvalid       = &Err{HTTPCode: http.StatusBadRequest, Code: 110002, Msg: "invalid shop"}
	ErrShopInvalidDomain = &Err{HTTPCode: http.StatusBadRequest, Code: 110003, Msg: "invalid domain format"}
	ErrShopDomainInUse   = &Err{HTTPCode: http.StatusConflict, Code: 110004, Msg: "custom domain already in use"}
	ErrShopPlanRestricted = &Err{HTTPCode: http.StatusForbidden, Code: 110005, Msg: "feature not available in current plan"}

	// Theme errors (110101-110199)
	ErrThemeNotFound       = &Err{HTTPCode: http.StatusNotFound, Code: 110101, Msg: "theme not found"}
	ErrThemeAlreadyActive  = &Err{HTTPCode: http.StatusBadRequest, Code: 110102, Msg: "theme is already active"}
	ErrThemeConfigInvalid  = &Err{HTTPCode: http.StatusBadRequest, Code: 110103, Msg: "invalid theme configuration"}
	ErrThemePresetOnly     = &Err{HTTPCode: http.StatusBadRequest, Code: 110104, Msg: "preset theme cannot be modified"}

	// Page errors (110201-110299)
	ErrPageNotFound        = &Err{HTTPCode: http.StatusNotFound, Code: 110201, Msg: "page not found"}
	ErrPageSlugDuplicate   = &Err{HTTPCode: http.StatusBadRequest, Code: 110202, Msg: "page slug already exists"}
	ErrPagePublishFailed   = &Err{HTTPCode: http.StatusInternalServerError, Code: 110203, Msg: "failed to publish page"}
	ErrPageAlreadyPublished = &Err{HTTPCode: http.StatusBadRequest, Code: 110204, Msg: "page is already published"}
	ErrPageNotPublished    = &Err{HTTPCode: http.StatusBadRequest, Code: 110205, Msg: "page is not published"}

	// Decoration errors (110301-110399)
	ErrDecorationNotFound  = &Err{HTTPCode: http.StatusNotFound, Code: 110301, Msg: "decoration block not found"}
	ErrInvalidBlockType    = &Err{HTTPCode: http.StatusBadRequest, Code: 110302, Msg: "invalid block type"}
	ErrInvalidBlockConfig  = &Err{HTTPCode: http.StatusBadRequest, Code: 110303, Msg: "invalid block configuration"}
	ErrBlockLimitExceeded  = &Err{HTTPCode: http.StatusBadRequest, Code: 110304, Msg: "block limit exceeded"}
	ErrDecorationPageMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 110305, Msg: "decoration does not belong to this page"}

	// Version errors (110401-110499)
	ErrVersionNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 110401, Msg: "version not found"}
	ErrVersionRestoreFailed = &Err{HTTPCode: http.StatusInternalServerError, Code: 110402, Msg: "failed to restore version"}
	ErrVersionCannotRestore = &Err{HTTPCode: http.StatusBadRequest, Code: 110403, Msg: "cannot restore to the same version"}

	// SEO errors (110501-110599)
	ErrSEOConfigNotFound   = &Err{HTTPCode: http.StatusNotFound, Code: 110501, Msg: "SEO config not found"}
	ErrSEOPageTypeInvalid  = &Err{HTTPCode: http.StatusBadRequest, Code: 110502, Msg: "invalid SEO page type"}
	ErrSEOTitleTooLong     = &Err{HTTPCode: http.StatusBadRequest, Code: 110503, Msg: "SEO title exceeds maximum length"}
	ErrSEODescriptionTooLong = &Err{HTTPCode: http.StatusBadRequest, Code: 110504, Msg: "SEO description exceeds maximum length"}

	// ErrShipmentNotFound ==================== Fulfillment Module (120xxx) ====================
	ErrShipmentNotFound              = &Err{HTTPCode: http.StatusNotFound, Code: 120001, Msg: "shipment not found"}
	ErrShipmentInvalidTracking       = &Err{HTTPCode: http.StatusBadRequest, Code: 120002, Msg: "invalid tracking number"}
	ErrShipmentAlreadyShipped        = &Err{HTTPCode: http.StatusBadRequest, Code: 120003, Msg: "order already shipped"}
	ErrShipmentCarrierRequired       = &Err{HTTPCode: http.StatusBadRequest, Code: 120004, Msg: "carrier is required"}
	ErrShipmentTrackingRequired      = &Err{HTTPCode: http.StatusBadRequest, Code: 120005, Msg: "tracking number is required"}
	ErrShipmentInvalidStatusTransition = &Err{HTTPCode: http.StatusBadRequest, Code: 120006, Msg: "invalid shipment status transition"}
	ErrShipmentCannotCancelDelivered = &Err{HTTPCode: http.StatusBadRequest, Code: 120007, Msg: "cannot cancel delivered shipment"}
	ErrShipmentItemNotFound          = &Err{HTTPCode: http.StatusNotFound, Code: 120008, Msg: "shipment item not found"}
	ErrShipmentOrderNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 120009, Msg: "order not found"}
	ErrShipmentOrderCannotShip       = &Err{HTTPCode: http.StatusBadRequest, Code: 120010, Msg: "order cannot be shipped"}
	ErrShipmentDuplicateTracking     = &Err{HTTPCode: http.StatusConflict, Code: 120011, Msg: "tracking number already exists"}
	ErrShipmentItemsRequired         = &Err{HTTPCode: http.StatusBadRequest, Code: 120012, Msg: "shipment items are required"}
	ErrShipmentItemQuantityExceeded  = &Err{HTTPCode: http.StatusBadRequest, Code: 120013, Msg: "shipment item quantity exceeded order quantity"}
	ErrShipmentInvalidItems         = &Err{HTTPCode: http.StatusBadRequest, Code: 120014, Msg: "invalid shipment items"}

	// Refund errors (1201xx)
	ErrRefundNotFound              = &Err{HTTPCode: http.StatusNotFound, Code: 120101, Msg: "refund not found"}
	ErrRefundInvalidStatus         = &Err{HTTPCode: http.StatusBadRequest, Code: 120102, Msg: "invalid refund status"}
	ErrRefundCannotCancel          = &Err{HTTPCode: http.StatusBadRequest, Code: 120103, Msg: "cannot cancel refund in current status"}
	ErrRefundAlreadyPending        = &Err{HTTPCode: http.StatusConflict, Code: 120104, Msg: "order already has pending refund"}
	ErrRefundOrderNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 120105, Msg: "order not found"}
	ErrRefundOrderCannotRefund     = &Err{HTTPCode: http.StatusBadRequest, Code: 120106, Msg: "order cannot be refunded"}
	ErrRefundTimeExpired           = &Err{HTTPCode: http.StatusBadRequest, Code: 120107, Msg: "refund period has expired"}
	ErrRefundReasonRequired        = &Err{HTTPCode: http.StatusBadRequest, Code: 120108, Msg: "refund reason is required"}
	ErrRefundRejectReasonRequired  = &Err{HTTPCode: http.StatusBadRequest, Code: 120109, Msg: "reject reason is required"}
	ErrRefundAmountExceeded        = &Err{HTTPCode: http.StatusBadRequest, Code: 120110, Msg: "refund amount exceeds order amount"}
	ErrRefundOrderNotPaid          = &Err{HTTPCode: http.StatusBadRequest, Code: 120111, Msg: "order is not paid"}

	// Carrier errors (1202xx)
	ErrCarrierNotFound   = &Err{HTTPCode: http.StatusNotFound, Code: 120201, Msg: "carrier not found"}
	ErrCarrierDuplicate  = &Err{HTTPCode: http.StatusConflict, Code: 120202, Msg: "carrier code already exists"}
	ErrCarrierInactive   = &Err{HTTPCode: http.StatusBadRequest, Code: 120203, Msg: "carrier is inactive"}

	// RefundReason errors (1203xx)
	ErrRefundReasonNotFound  = &Err{HTTPCode: http.StatusNotFound, Code: 120301, Msg: "refund reason not found"}
	ErrRefundReasonDuplicate = &Err{HTTPCode: http.StatusConflict, Code: 120302, Msg: "refund reason code already exists"}

	// ErrSharedCurrencyMismatch ==================== Shared Module (200xxx) ====================
	ErrSharedCurrencyMismatch   = &Err{HTTPCode: http.StatusBadRequest, Code: 200001, Msg: "currency mismatch"}
	ErrSharedInsufficientAmount = &Err{HTTPCode: http.StatusBadRequest, Code: 200002, Msg: "insufficient amount"}
	ErrSharedInvalidAmount      = &Err{HTTPCode: http.StatusBadRequest, Code: 200003, Msg: "invalid amount"}
	ErrSharedNotFound           = &Err{HTTPCode: http.StatusNotFound, Code: 200004, Msg: "resource not found"}
	ErrSharedDuplicate          = &Err{HTTPCode: http.StatusConflict, Code: 200005, Msg: "duplicate resource"}
	ErrSharedInvalidParam       = &Err{HTTPCode: http.StatusBadRequest, Code: 200006, Msg: "invalid parameter"}

	// ErrAuthTokenInvalid ==================== Auth Module (130xxx) ====================
	ErrAuthTokenInvalid = &Err{HTTPCode: http.StatusUnauthorized, Code: 130001, Msg: "invalid token"}
	ErrAuthTokenExpired = &Err{HTTPCode: http.StatusUnauthorized, Code: 130002, Msg: "token expired"}

	// ErrCacheNotFound ==================== Cache Module (140xxx) ====================
	ErrCacheNotFound   = &Err{HTTPCode: http.StatusNotFound, Code: 140001, Msg: "not found"}
	ErrCacheFromSource = &Err{HTTPCode: http.StatusInternalServerError, Code: 140002, Msg: "from source err"}

	// ErrMarketNotFound ==================== Market Module (150xxx) ====================
	ErrMarketNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 150001, Msg: "market not found"}
	ErrMarketDuplicate        = &Err{HTTPCode: http.StatusConflict, Code: 150002, Msg: "duplicate market code"}
	ErrMarketCodeRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 150003, Msg: "market code is required"}
	ErrMarketNameRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 150004, Msg: "market name is required"}
	ErrMarketCurrencyRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 150005, Msg: "market currency is required"}
	ErrMarketInactive         = &Err{HTTPCode: http.StatusBadRequest, Code: 150006, Msg: "market is inactive"}
	ErrMarketAlreadyDefault   = &Err{HTTPCode: http.StatusBadRequest, Code: 150007, Msg: "market is already default"}
	ErrMarketCannotDelete     = &Err{HTTPCode: http.StatusBadRequest, Code: 150008, Msg: "cannot delete default market"}

	// ErrProductMarketNotFound ==================== ProductMarket Module (160xxx) ====================
	ErrProductMarketNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 160001, Msg: "product market not found"}
	ErrProductMarketAlreadyExists = &Err{HTTPCode: http.StatusConflict, Code: 160002, Msg: "product already in market"}
	ErrProductMarketPriceRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 160003, Msg: "price is required for market"}

	// ErrInventoryInsufficientStock ==================== Inventory Module (170xxx) ====================
	ErrInventoryInsufficientStock       = &Err{HTTPCode: http.StatusBadRequest, Code: 170001, Msg: "insufficient available stock"}
	ErrInventoryInsufficientLockedStock = &Err{HTTPCode: http.StatusBadRequest, Code: 170002, Msg: "insufficient locked stock"}
	ErrInventorySKUNotFound             = &Err{HTTPCode: http.StatusNotFound, Code: 170003, Msg: "sku not found"}
	ErrInventoryWarehouseNotFound       = &Err{HTTPCode: http.StatusNotFound, Code: 170004, Msg: "warehouse not found"}
	ErrInventoryDuplicateWarehouseCode  = &Err{HTTPCode: http.StatusConflict, Code: 170005, Msg: "duplicate warehouse code"}

	// ErrBrandNotFound ==================== Brand Module (180xxx) ====================
	ErrBrandNotFound    = &Err{HTTPCode: http.StatusNotFound, Code: 180001, Msg: "brand not found"}
	ErrBrandDuplicate   = &Err{HTTPCode: http.StatusConflict, Code: 180002, Msg: "brand name already exists"}
	ErrBrandHasProducts = &Err{HTTPCode: http.StatusBadRequest, Code: 180003, Msg: "brand has products"}

	// ErrSKUPrefixInvalid ==================== SKU Module (190xxx) ====================
	ErrSKUPrefixInvalid          = &Err{HTTPCode: http.StatusBadRequest, Code: 190001, Msg: "SKU前缀格式无效"}
	ErrSKUPrefixTooLong          = &Err{HTTPCode: http.StatusBadRequest, Code: 190002, Msg: "SKU前缀长度超过限制"}
	ErrSKUGenerateFailed         = &Err{HTTPCode: http.StatusInternalServerError, Code: 190003, Msg: "SKU生成失败"}
	ErrSKUParseFailed            = &Err{HTTPCode: http.StatusBadRequest, Code: 190004, Msg: "SKU编码解析失败"}
	ErrSKUCodeTooLong            = &Err{HTTPCode: http.StatusBadRequest, Code: 190005, Msg: "SKU编码长度超过限制"}
	ErrSKUPrefixStartsWithNumber = &Err{HTTPCode: http.StatusBadRequest, Code: 190006, Msg: "SKU前缀不能以数字开头"}

	// ErrReviewNotFound ==================== Review Module (210xxx) ====================
	ErrReviewNotFound            = &Err{HTTPCode: http.StatusNotFound, Code: 210001, Msg: "review not found"}
	ErrReviewAlreadyReplied      = &Err{HTTPCode: http.StatusBadRequest, Code: 210002, Msg: "review already has reply"}
	ErrReviewCannotReplyHidden   = &Err{HTTPCode: http.StatusBadRequest, Code: 210003, Msg: "cannot reply to hidden review"}
	ErrReviewInvalidStatus       = &Err{HTTPCode: http.StatusBadRequest, Code: 210004, Msg: "invalid review status"}
	ErrReviewContentTooLong      = &Err{HTTPCode: http.StatusBadRequest, Code: 210005, Msg: "review content exceeds limit"}
	ErrReplyContentTooLong       = &Err{HTTPCode: http.StatusBadRequest, Code: 210006, Msg: "reply content exceeds limit"}
	ErrReplyNotFound             = &Err{HTTPCode: http.StatusNotFound, Code: 210007, Msg: "reply not found"}
	ErrReviewCannotApprove       = &Err{HTTPCode: http.StatusBadRequest, Code: 210008, Msg: "cannot approve review in current status"}
	ErrReviewCannotHide          = &Err{HTTPCode: http.StatusBadRequest, Code: 210009, Msg: "cannot hide review in current status"}
	ErrReviewCannotShow          = &Err{HTTPCode: http.StatusBadRequest, Code: 210010, Msg: "cannot show review in current status"}
	ErrReviewCannotFeature       = &Err{HTTPCode: http.StatusBadRequest, Code: 210011, Msg: "can only feature approved reviews"}
	ErrReviewAlreadyDeleted      = &Err{HTTPCode: http.StatusBadRequest, Code: 210012, Msg: "review already deleted"}
	ErrReviewReplyEmpty          = &Err{HTTPCode: http.StatusBadRequest, Code: 210013, Msg: "reply content cannot be empty"}
	ErrReviewInvalidRating       = &Err{HTTPCode: http.StatusBadRequest, Code: 210014, Msg: "rating must be between 1 and 5"}
	ErrReviewBatchEmpty          = &Err{HTTPCode: http.StatusBadRequest, Code: 210015, Msg: "batch operation requires at least one review id"}
	ErrReviewBatchLimitExceeded  = &Err{HTTPCode: http.StatusBadRequest, Code: 210016, Msg: "batch operation limited to 100 reviews"}

	// ErrDashboardDataUnavailable ==================== Dashboard Module (220xxx) ====================
	ErrDashboardDataUnavailable = &Err{HTTPCode: http.StatusServiceUnavailable, Code: 220001, Msg: "dashboard data temporarily unavailable"}

	// Inventory Transfer errors (additional 170xxx)
	ErrInventoryTransferFailed       = &Err{HTTPCode: http.StatusBadRequest, Code: 170006, Msg: "stock transfer failed"}
	ErrInsufficientStockForTransfer  = &Err{HTTPCode: http.StatusBadRequest, Code: 170007, Msg: "insufficient stock for transfer"}
	ErrSameWarehouseTransfer         = &Err{HTTPCode: http.StatusBadRequest, Code: 170008, Msg: "cannot transfer to same warehouse"}

	// ==================== Shipping Module (230xxx) ====================
	ErrShippingTemplateNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 230001, Msg: "shipping template not found"}
	ErrShippingTemplateNameRequired  = &Err{HTTPCode: http.StatusBadRequest, Code: 230002, Msg: "template name is required"}
	ErrShippingTemplateHasZones      = &Err{HTTPCode: http.StatusBadRequest, Code: 230003, Msg: "cannot delete template with zones"}
	ErrShippingTemplateIsDefault     = &Err{HTTPCode: http.StatusBadRequest, Code: 230004, Msg: "cannot delete default template"}
	ErrShippingTemplateDuplicate     = &Err{HTTPCode: http.StatusConflict, Code: 230005, Msg: "template name already exists"}

	// Shipping Zone errors (2301xx)
	ErrShippingZoneNotFound          = &Err{HTTPCode: http.StatusNotFound, Code: 230101, Msg: "shipping zone not found"}
	ErrShippingZoneNameRequired      = &Err{HTTPCode: http.StatusBadRequest, Code: 230102, Msg: "zone name is required"}
	ErrShippingZoneRegionsRequired   = &Err{HTTPCode: http.StatusBadRequest, Code: 230103, Msg: "zone regions are required"}
	ErrShippingZoneInvalidFeeType    = &Err{HTTPCode: http.StatusBadRequest, Code: 230104, Msg: "invalid fee type"}
	ErrShippingZoneFeeConfigRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 230105, Msg: "fee configuration is required"}
	ErrShippingZoneDuplicateRegion   = &Err{HTTPCode: http.StatusBadRequest, Code: 230106, Msg: "region already assigned to another zone"}

	// Shipping Mapping errors (2302xx)
	ErrShippingMappingNotFound       = &Err{HTTPCode: http.StatusNotFound, Code: 230201, Msg: "shipping mapping not found"}
	ErrShippingMappingAlreadyExists  = &Err{HTTPCode: http.StatusConflict, Code: 230202, Msg: "mapping already exists"}
	ErrShippingMappingInvalidTarget  = &Err{HTTPCode: http.StatusBadRequest, Code: 230203, Msg: "invalid target type"}

	// Shipping Calculator errors (2303xx)
	ErrShippingCalcNoMatchZone       = &Err{HTTPCode: http.StatusBadRequest, Code: 230301, Msg: "no matching zone for address"}
	ErrShippingCalcNoDefaultTemplate = &Err{HTTPCode: http.StatusBadRequest, Code: 230302, Msg: "no default shipping template configured"}
	ErrShippingCalcItemsRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 230303, Msg: "items are required"}
	ErrShippingCalcAddressRequired   = &Err{HTTPCode: http.StatusBadRequest, Code: 230304, Msg: "address is required"}
	ErrShippingCalcInvalidQuantity   = &Err{HTTPCode: http.StatusBadRequest, Code: 230305, Msg: "invalid quantity in items"}
	ErrShippingCalcInvalidWeight     = &Err{HTTPCode: http.StatusBadRequest, Code: 230306, Msg: "invalid weight in items"}
	ErrShippingCalcInvalidPrice      = &Err{HTTPCode: http.StatusBadRequest, Code: 230307, Msg: "invalid price in items"}

	// ==================== Upload Module (20xxx) ====================
	ErrUploadUnsupportedFileType = &Err{HTTPCode: http.StatusBadRequest, Code: 200001, Msg: "unsupported file type"}
	ErrUploadFileSizeExceeded    = &Err{HTTPCode: http.StatusBadRequest, Code: 200002, Msg: "file size exceeded"}
	ErrUploadInvalidCategory     = &Err{HTTPCode: http.StatusBadRequest, Code: 200003, Msg: "invalid category"}
	ErrUploadFailed              = &Err{HTTPCode: http.StatusInternalServerError, Code: 200004, Msg: "upload failed"}
	ErrUploadNotFound            = &Err{HTTPCode: http.StatusNotFound, Code: 200005, Msg: "file not found"}
)
