package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchUpdateProductLogic {
	return BatchUpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateProductLogic) BatchUpdateProduct(req *types.BatchUpdateProductReq) (resp *types.BatchUpdateProductResp, err error) {
	// 从 context 获取 tenantID
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// 平台管理员设置 tenantID = 0 以访问所有数据
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// 验证 UpdateFields 至少提供一个字段
	if req.UpdateFields.Price == nil && req.UpdateFields.Stock == nil &&
		req.UpdateFields.Status == nil && req.UpdateFields.CategoryID == nil {
		return nil, code.ErrParam
	}

	// 解析价格
	var price *decimal.Decimal
	if req.UpdateFields.Price != nil {
		p, err := decimal.NewFromString(*req.UpdateFields.Price)
		if err != nil {
			return &types.BatchUpdateProductResp{
				Success: []int64{},
				Failed: []types.BatchProductFail{
					{ProductID: 0, Code: 30002, Message: "无效的价格格式"},
				},
			}, nil
		}
		price = &p
	}

	// 解析状态
	var status *product.Status
	if req.UpdateFields.Status != nil {
		s := parseBatchStatus(*req.UpdateFields.Status)
		if s == nil {
			return &types.BatchUpdateProductResp{
				Success: []int64{},
				Failed: []types.BatchProductFail{
					{ProductID: 0, Code: 30006, Message: "无效的状态值"},
				},
			}, nil
		}
		status = s
	}

	// 构建批量更新请求
	batchReq := appProduct.BatchUpdateProductRequest{
		ProductIDs: req.ProductIDs,
		Fields: appProduct.BatchProductFields{
			Price:      price,
			Stock:      req.UpdateFields.Stock,
			Status:     status,
			CategoryID: req.UpdateFields.CategoryID,
		},
	}

	// 执行批量更新
	successIDs, failed, err := l.svcCtx.ProductService.BatchUpdateProduct(l.ctx, shared.TenantID(tenantID), batchReq)
	if err != nil {
		return nil, err
	}

	// 转换失败结果
	failedResult := make([]types.BatchProductFail, len(failed))
	for i, f := range failed {
		failedResult[i] = types.BatchProductFail{
			ProductID: f.ProductID,
			Code:      f.Code,
			Message:   f.Message,
		}
	}

	return &types.BatchUpdateProductResp{
		Success: successIDs,
		Failed:  failedResult,
	}, nil
}

// parseBatchStatus 解析批量更新时的状态字符串
func parseBatchStatus(s string) *product.Status {
	switch s {
	case "draft":
		status := product.StatusDraft
		return &status
	case "on_sale":
		status := product.StatusOnSale
		return &status
	case "off_sale":
		status := product.StatusOffSale
		return &status
	default:
		return nil
	}
}
