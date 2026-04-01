package user_coupons

import (
	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgcoupon "github.com/colinrs/shopjoy/pkg/domain/promotion"
)

func mapUserCouponStatus(statusStr string) pkgcoupon.UserCouponStatus {
	switch statusStr {
	case "available":
		return pkgcoupon.UserCouponStatusUnused
	case "used":
		return pkgcoupon.UserCouponStatusUsed
	case "expired":
		return pkgcoupon.UserCouponStatusExpired
	default:
		return pkgcoupon.UserCouponStatusUnused
	}
}

func mapUserCouponStatusToString(status int) string {
	switch pkgcoupon.UserCouponStatus(status) {
	case pkgcoupon.UserCouponStatusUnused:
		return "available"
	case pkgcoupon.UserCouponStatusUsed:
		return "used"
	case pkgcoupon.UserCouponStatusExpired:
		return "expired"
	default:
		return "available"
	}
}

func convertUserCouponToDetailResp(uc *apppromotion.UserCouponResponse) *types.UserCouponDetailResp {
	return &types.UserCouponDetailResp{
		ID:           uc.ID,
		UserID:       0, // Not available in response
		CouponID:     uc.CouponID,
		CouponCode:   uc.CouponCode,
		CouponName:   uc.CouponName,
		DiscountType: "fixed_amount", // Default, should be fetched from coupon
		Status:       mapUserCouponStatusToString(uc.Status),
		StartTime:    "", // Not available in response
		EndTime:      uc.ExpireAt,
		UsedAt:       uc.UsedAt,
		OrderID:      uc.OrderID,
		CreatedAt:    uc.ReceivedAt,
	}
}