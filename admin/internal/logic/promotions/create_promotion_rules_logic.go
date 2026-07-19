package promotions

import (
	"context"
	"strconv"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionRulesLogic {
	return CreatePromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreatePromotionRules takes the owner kind + owner id from the path
// (no Promotion lookup required) and forwards the rule set to
// PromotionApp.CreateRules. The app layer expects a typed Kind
// (PROMOTION / COUPON) so it can route the row to the correct rules
// table slice / owner_kind column.
func (l *CreatePromotionRulesLogic) CreatePromotionRules(req *types.CreatePromotionRulesReq) (resp *types.CreatePromotionRulesResp, err error) {
	ownerKind := pkgpromotion.Kind(strings.ToUpper(req.OwnerKind))
	rules := convertRuleReqsToDomain(req.Rules)
	out, err := l.svcCtx.PromotionApp.CreateRules(l.ctx, ownerKind, req.OwnerID, rules)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(out))
	for i, r := range out {
		ids[i] = strconv.FormatInt(r.ID, 10)
	}
	return &types.CreatePromotionRulesResp{IDs: ids}, nil
}