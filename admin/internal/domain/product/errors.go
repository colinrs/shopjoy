package product

import "errors"

var (
	ErrEmptyName               = errors.New("商品名称不能为空")
	ErrInvalidPrice            = errors.New("商品价格必须大于0")
	ErrCurrencyMismatch        = errors.New("币种不匹配")
	ErrInsufficientAmount      = errors.New("金额不足")
	ErrProductDeleted          = errors.New("商品已删除")
	ErrInvalidStatusTransition = errors.New("无效的状态转换")
	ErrNoStock                 = errors.New("库存不能为0")
	ErrNegativeStock           = errors.New("库存不能为负数")
	ErrProductNotOnSale        = errors.New("商品未上架")
	ErrInvalidQuantity         = errors.New("无效的数量")
	ErrInsufficientStock       = errors.New("库存不足")
	ErrProductNotFound         = errors.New("商品不存在")
)
