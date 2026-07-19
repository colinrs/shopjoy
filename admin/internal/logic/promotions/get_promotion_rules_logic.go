package promotions

import (
	"context"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPromotionRulesLogic {
	return GetPromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetPromotionRules takes the owner kind + owner id from the path
// (no Promotion lookup required) and asks the unified PromotionApp
// for its rules. The wire response uses the unified rule shape
// (condition_type / condition_value / action_type / action_value).
func (l *GetPromotionRulesLogic) GetPromotionRules(req *types.GetPromotionRulesReq) (resp *types.ListPromotionRulesResp, err error) {
	ownerKind := pkgpromotion.Kind(strings.ToUpper(req.OwnerKind))
	rules, err := l.svcCtx.PromotionApp.GetRules(l.ctx, ownerKind, req.OwnerID)
	if err != nil {
		return nil, err
	}

	return &types.ListPromotionRulesResp{
		List:  convertRulesPtrToResp(rules),
		Total: int64(len(rules)),
	}, nil
}