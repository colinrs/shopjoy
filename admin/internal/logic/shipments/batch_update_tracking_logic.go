package shipments

import (
	"context"
	"errors"
	"strconv"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/utils"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateTrackingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateTrackingLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchUpdateTrackingLogic {
	return BatchUpdateTrackingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateTrackingLogic) BatchUpdateTracking(req *types.BatchUpdateTrackingReq) (resp *types.BatchUpdateTrackingResp, err error) {
	// Get tenantID from context

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Validate tracking number
	if req.TrackingNo == "" {
		return nil, code.ErrShipmentTrackingRequired
	}

	// Initialize response
	resp = &types.BatchUpdateTrackingResp{
		Success: make([]string, 0, len(req.ShipmentIDs)),
		Failed:  make([]types.BatchTrackingFail, 0),
	}

	// Parse shipment IDs (string -> int64)
	shipmentIDs, err := utils.ParseInt64Slice(req.ShipmentIDs)
	if err != nil {
		return nil, code.ErrParam
	}

	// Process each shipment
	for _, shipmentID := range shipmentIDs {
		// Build update request
		var weight decimal.Decimal
		if req.Weight != nil {
			var err error
			weight, err = decimal.NewFromString(*req.Weight)
			if err != nil {
				resp.Failed = append(resp.Failed, types.BatchTrackingFail{
					ShipmentID: shipmentID,
					Code:       code.ErrSharedInvalidParam.Code,
					Message:    "invalid weight format",
				})
				continue
			}
		}

		updateReq := appfulfillment.UpdateShipmentRequest{
			ID:          shipmentID,
			CarrierCode: req.CarrierCode,
			TrackingNo:  req.TrackingNo,
			Weight:      weight,
		}

		// Update shipment - carrier validation is handled internally by UpdateShipment
		_, err := l.svcCtx.ShipmentApp.UpdateShipment(l.ctx, userID, updateReq)
		if err != nil {
			// Extract error code and message
			var codeErr *code.Err
			if errors.As(err, &codeErr) {
				resp.Failed = append(resp.Failed, types.BatchTrackingFail{
					ShipmentID: shipmentID,
					Code:       codeErr.Code,
					Message:    codeErr.Msg,
				})
			} else {
				resp.Failed = append(resp.Failed, types.BatchTrackingFail{
					ShipmentID: shipmentID,
					Code:       code.UnknownErr.Code,
					Message:    err.Error(),
				})
			}
		} else {
			resp.Success = append(resp.Success, strconv.FormatInt(shipmentID, 10))
		}
	}

	return resp, nil
}
