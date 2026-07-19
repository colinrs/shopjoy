package user_coupons

import (
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
)

func mapUserCouponStatus(statusStr string) pkgpromotion.UserCouponStatus {
	switch statusStr {
	case "available":
		return pkgpromotion.UserCouponStatusUnused
	case "used":
		return pkgpromotion.UserCouponStatusUsed
	case "expired":
		return pkgpromotion.UserCouponStatusExpired
	default:
		return pkgpromotion.UserCouponStatusUnused
	}
}

func mapUserCouponStatusToString(status pkgpromotion.UserCouponStatus) string {
	switch status {
	case pkgpromotion.UserCouponStatusUnused:
		return "available"
	case pkgpromotion.UserCouponStatusUsed:
		return "used"
	case pkgpromotion.UserCouponStatusExpired:
		return "expired"
	default:
		return "available"
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// convertUserCouponToDetailResp maps the unified UserCouponResponse to the
// legacy wire shape. CouponCode / CouponName / DiscountType are not available
// in the unified App response (would require an extra JOIN); Task 8 will
// regenerate wire types and these fields will be sourced from the joined
// PromotionResponse. Until then, leave them empty/default.
func convertUserCouponToDetailResp(uc *apppromotion.UserCouponResponse) *types.UserCouponDetailResp {
	if uc == nil {
		return nil
	}
	return &types.UserCouponDetailResp{
		ID:           uc.ID,
		CouponID:     uc.CouponID,
		DiscountType: "fixed_amount", // default; real type from joined Promotion in Task 8
		Status:       mapUserCouponStatusToString(uc.Status),
		EndTime:      formatTime(uc.ExpireAt),
		UsedAt:       formatTimePtr(uc.UsedAt),
		OrderID:      uc.OrderID,
		CreatedAt:    formatTime(uc.ReceivedAt),
	}
}