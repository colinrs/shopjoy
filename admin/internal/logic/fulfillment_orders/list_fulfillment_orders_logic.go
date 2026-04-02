package fulfillment_orders

import (
	"context"
	"strconv"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFulfillmentOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFulfillmentOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListFulfillmentOrdersLogic {
	return ListFulfillmentOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFulfillmentOrdersLogic) ListFulfillmentOrders(req *types.ListFulfillmentOrdersReq) (resp *types.ListFulfillmentOrdersResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request
	queryReq := appfulfillment.QueryOrderRequest{
		Page:              req.Page,
		PageSize:          req.PageSize,
		OrderNo:           req.OrderNo,
		UserID:            req.UserID,
		UserName:          req.UserName,
		Status:            req.Status,
	}

	// Parse fulfillment status - convert from string to int8
	if req.FulfillmentStatus != "" {
		v, _ := strconv.ParseInt(req.FulfillmentStatus, 10, 8)
		queryReq.FulfillmentStatus = int8(v)
	}

	// Parse refund status - convert from string to int8
	if req.RefundStatus != "" {
		v, _ := strconv.ParseInt(req.RefundStatus, 10, 8)
		queryReq.RefundStatus = int8(v)
	}

	// Parse time range
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			queryReq.StartTime = startTime
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			queryReq.EndTime = endTime
		}
	}

	// List orders
	listResp, err := l.svcCtx.OrderFulfillmentApp.ListOrders(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return nil, err
	}

	// Build response
	resp = &types.ListFulfillmentOrdersResp{
		List:     make([]*types.OrderFulfillmentDetailResp, len(listResp.List)),
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}

	for i, o := range listResp.List {
		resp.List[i] = toOrderFulfillmentDetailResp(o)
	}

	return resp, nil
}

// toOrderFulfillmentDetailResp converts OrderFulfillmentDetail to types.OrderFulfillmentDetailResp
func toOrderFulfillmentDetailResp(o *appfulfillment.OrderFulfillmentDetail) *types.OrderFulfillmentDetailResp {
	items := make([]*types.OrderFulfillmentItemResp, len(o.Items))
	for i, item := range o.Items {
		items[i] = &types.OrderFulfillmentItemResp{
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			SKUID:       item.SKUID,
			ProductName: item.ProductName,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Quantity:    item.Quantity,
			ShippedQty:  item.ShippedQty,
			PendingQty:  item.PendingQty,
			UnitPrice:   utils.FormatAmount(item.UnitPrice),
			Currency:    item.Currency,
		}
	}

	shipments := make([]*types.ShipmentDetailResp, len(o.Shipments))
	for i, s := range o.Shipments {
		shipments[i] = toShipmentDetailResp(s)
	}

	var refund *types.RefundDetailResp
	if o.Refund != nil {
		refund = toRefundDetailResp(o.Refund)
	}

	var shippingAddress *types.OrderShippingAddress
	if o.ShippingAddress != nil {
		shippingAddress = &types.OrderShippingAddress{
			ReceiverName:  o.ShippingAddress.ReceiverName,
			ReceiverPhone: o.ShippingAddress.ReceiverPhone,
			Province:      o.ShippingAddress.Province,
			City:          o.ShippingAddress.City,
			District:      o.ShippingAddress.District,
			Address:       o.ShippingAddress.Address,
			FullAddress:   o.ShippingAddress.FullAddress,
		}
	}

	return &types.OrderFulfillmentDetailResp{
		OrderID:           o.OrderID,
		OrderNo:           o.OrderNo,
		Status:            o.Status,
		FulfillmentStatus: fulfillment.FulfillmentStatus(o.FulfillmentStatus).String(),
		FulfillmentText:   o.FulfillmentText,
		RefundStatus:      fulfillment.RefundStatus(o.RefundStatus).String(),
		RefundText:        o.RefundText,
		TotalAmount:       utils.FormatAmount(o.TotalAmount),
		Currency:          o.Currency,
		UserID:            o.UserID,
		UserName:          o.UserName,
		UserPhone:         o.UserPhone,
		ShippingAddress:   shippingAddress,
		Items:             items,
		Shipments:         shipments,
		Refund:            refund,
		PaidAt:            utils.FormatTimeToRFC3339(o.PaidAt),
		ShippedAt:         utils.FormatTimeToRFC3339(o.ShippedAt),
		DeliveredAt:       utils.FormatTimeToRFC3339(o.DeliveredAt),
		CreatedAt:         o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         o.UpdatedAt.Format(time.RFC3339),
	}
}