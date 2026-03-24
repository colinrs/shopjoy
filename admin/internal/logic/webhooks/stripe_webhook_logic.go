package webhooks

import (
	"context"
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type StripeWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStripeWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) StripeWebhookLogic {
	return StripeWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StripeWebhookLogic) StripeWebhook(r *http.Request) error {
	// todo: add your logic here and delete this line

	return nil
}
