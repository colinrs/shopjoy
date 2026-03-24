package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserAddressesLogic {
	return GetUserAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAddressesLogic) GetUserAddresses(req *types.GetUserRequest) (resp *types.UserAddressListResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	addressList, err := l.svcCtx.UserService.GetAddresses(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	addresses := make([]*types.UserAddressResponse, 0, len(addressList.List))
	for _, addr := range addressList.List {
		addresses = append(addresses, &types.UserAddressResponse{
			ID:         addr.ID,
			Name:       addr.Name,
			Phone:      addr.Phone,
			Country:    addr.Country,
			Province:   addr.Province,
			City:       addr.City,
			District:   addr.District,
			Address:    addr.Address,
			PostalCode: addr.PostalCode,
			IsDefault:  addr.IsDefault,
			CreatedAt:  addr.CreatedAt,
		})
	}

	return &types.UserAddressListResponse{
		List:  addresses,
		Total: addressList.Total,
	}, nil
}
