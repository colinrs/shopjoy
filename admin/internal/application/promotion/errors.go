package promotion

import "errors"

var (
	// Promotion errors
	ErrPromotionIsActive     = errors.New("promotion is active, cannot modify or delete")
	ErrPromotionAlreadyEnded = errors.New("promotion has already ended")

	// Coupon errors
	ErrCouponIsActive   = errors.New("coupon is active, cannot modify or delete")
	ErrCouponNotActive  = errors.New("coupon is not active")
	ErrCouponUsedUp     = errors.New("coupon usage limit reached")
)