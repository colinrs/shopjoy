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

	// ErrUserInvalidEmail ==================== User Module (11xxx) ====================
	ErrUserInvalidEmail     = &Err{HTTPCode: http.StatusBadRequest, Code: 11001, Msg: "invalid email format"}
	ErrUserInvalidPhone     = &Err{HTTPCode: http.StatusBadRequest, Code: 11002, Msg: "invalid phone format"}
	ErrUserPasswordTooWeak  = &Err{HTTPCode: http.StatusBadRequest, Code: 11003, Msg: "password too weak"}
	ErrUserNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 11004, Msg: "user not found"}
	ErrUserDuplicateUser    = &Err{HTTPCode: http.StatusConflict, Code: 11005, Msg: "duplicate user"}
	ErrUserWrongPassword    = &Err{HTTPCode: http.StatusUnauthorized, Code: 11006, Msg: "wrong password"}
	ErrUserAlreadyDeleted   = &Err{HTTPCode: http.StatusBadRequest, Code: 11007, Msg: "user already deleted"}
	ErrUserPasswordMismatch = &Err{HTTPCode: http.StatusBadRequest, Code: 11008, Msg: "passwords do not match"}

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
	ErrCategoryNotFound    = &Err{HTTPCode: http.StatusNotFound, Code: 30101, Msg: "category not found"}
	ErrCategoryDuplicate   = &Err{HTTPCode: http.StatusConflict, Code: 30102, Msg: "duplicate category"}
	ErrCategoryInvalid     = &Err{HTTPCode: http.StatusBadRequest, Code: 30103, Msg: "invalid category"}
	ErrCategoryHasChildren = &Err{HTTPCode: http.StatusBadRequest, Code: 30104, Msg: "category has children"}
	ErrCategoryHasProducts = &Err{HTTPCode: http.StatusBadRequest, Code: 30105, Msg: "category has products"}

	// ErrOrderNotFound ==================== Order Module (40xxx) ====================
	ErrOrderNotFound          = &Err{HTTPCode: http.StatusNotFound, Code: 40001, Msg: "order not found"}
	ErrOrderInvalidStatus     = &Err{HTTPCode: http.StatusBadRequest, Code: 40002, Msg: "invalid order status"}
	ErrOrderAlreadyPaid       = &Err{HTTPCode: http.StatusBadRequest, Code: 40003, Msg: "order already paid"}
	ErrOrderNotPaid           = &Err{HTTPCode: http.StatusBadRequest, Code: 40004, Msg: "order not paid"}
	ErrOrderExpired           = &Err{HTTPCode: http.StatusBadRequest, Code: 40005, Msg: "order expired"}
	ErrOrderInsufficientStock = &Err{HTTPCode: http.StatusBadRequest, Code: 40006, Msg: "insufficient stock"}
	ErrOrderInvalidAmount     = &Err{HTTPCode: http.StatusBadRequest, Code: 40007, Msg: "invalid amount"}
	ErrOrderCartEmpty         = &Err{HTTPCode: http.StatusBadRequest, Code: 40008, Msg: "cart is empty"}

	// ErrPaymentNotFound ==================== Payment Module (50xxx) ====================
	ErrPaymentNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 50001, Msg: "payment not found"}
	ErrPaymentInvalidAmount = &Err{HTTPCode: http.StatusBadRequest, Code: 50002, Msg: "invalid payment amount"}
	ErrPaymentFailed        = &Err{HTTPCode: http.StatusPaymentRequired, Code: 50003, Msg: "payment failed"}
	ErrPaymentAlreadyPaid   = &Err{HTTPCode: http.StatusBadRequest, Code: 50004, Msg: "payment already completed"}
	ErrPaymentExpired       = &Err{HTTPCode: http.StatusBadRequest, Code: 50005, Msg: "payment expired"}

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

	// ErrPromotionNotFound ==================== Promotion Module (80xxx) ====================
	ErrPromotionNotFound   = &Err{HTTPCode: http.StatusNotFound, Code: 80001, Msg: "promotion not found"}
	ErrPromotionInvalid    = &Err{HTTPCode: http.StatusBadRequest, Code: 80002, Msg: "invalid promotion"}
	ErrPromotionExpired    = &Err{HTTPCode: http.StatusBadRequest, Code: 80003, Msg: "promotion expired"}
	ErrPromotionNotStarted = &Err{HTTPCode: http.StatusBadRequest, Code: 80004, Msg: "promotion not started"}

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
	ErrRoleNotFound  = &Err{HTTPCode: http.StatusNotFound, Code: 100001, Msg: "role not found"}
	ErrRoleDuplicate = &Err{HTTPCode: http.StatusConflict, Code: 100002, Msg: "duplicate role"}
	ErrRoleInvalid   = &Err{HTTPCode: http.StatusBadRequest, Code: 100003, Msg: "invalid role"}

	// ErrShopNotFound ==================== Storefront Module (110xxx) ====================
	ErrShopNotFound = &Err{HTTPCode: http.StatusNotFound, Code: 110001, Msg: "shop not found"}
	ErrShopInvalid  = &Err{HTTPCode: http.StatusBadRequest, Code: 110002, Msg: "invalid shop"}

	// ErrShipmentNotFound ==================== Fulfillment Module (120xxx) ====================
	ErrShipmentNotFound        = &Err{HTTPCode: http.StatusNotFound, Code: 120001, Msg: "shipment not found"}
	ErrShipmentInvalidTracking = &Err{HTTPCode: http.StatusBadRequest, Code: 120002, Msg: "invalid tracking number"}
	ErrShipmentAlreadyShipped  = &Err{HTTPCode: http.StatusBadRequest, Code: 120003, Msg: "order already shipped"}

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
)
