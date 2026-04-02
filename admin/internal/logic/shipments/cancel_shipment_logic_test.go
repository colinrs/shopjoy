package shipments

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	domain "github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// mockShipmentApp is a mock implementation of fulfillment.ShipmentApp
type mockShipmentApp struct {
	getShipmentFunc    func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error)
	cancelShipmentFunc func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error)
}

func (m *mockShipmentApp) CreateShipment(ctx context.Context, tenantID shared.TenantID, userID int64, req fulfillment.CreateShipmentRequest) (*fulfillment.ShipmentResponse, error) {
	return nil, nil
}

func (m *mockShipmentApp) BatchCreateShipments(ctx context.Context, tenantID shared.TenantID, userID int64, carrierCode, carrierName string, shipments []fulfillment.BatchShipmentItem) (*fulfillment.BatchShipmentResult, error) {
	return nil, nil
}

func (m *mockShipmentApp) GetShipment(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
	if m.getShipmentFunc != nil {
		return m.getShipmentFunc(ctx, tenantID, id)
	}
	return nil, nil
}

func (m *mockShipmentApp) ListShipments(ctx context.Context, tenantID shared.TenantID, req fulfillment.QueryShipmentRequest) (*fulfillment.ShipmentListResponse, error) {
	return nil, nil
}

func (m *mockShipmentApp) UpdateShipment(ctx context.Context, tenantID shared.TenantID, userID int64, req fulfillment.UpdateShipmentRequest) (*fulfillment.ShipmentResponse, error) {
	return nil, nil
}

func (m *mockShipmentApp) UpdateShipmentStatus(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, status domain.ShipmentStatus) (*fulfillment.ShipmentResponse, error) {
	return nil, nil
}

func (m *mockShipmentApp) GetOrderShipments(ctx context.Context, tenantID shared.TenantID, orderID int64) ([]*fulfillment.ShipmentResponse, error) {
	return nil, nil
}

func (m *mockShipmentApp) CancelShipment(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
	if m.cancelShipmentFunc != nil {
		return m.cancelShipmentFunc(ctx, tenantID, userID, id, reason)
	}
	return nil, nil
}

// mockServiceContext creates a minimal service context for testing
func mockServiceContext(mockApp *mockShipmentApp) *svc.ServiceContext {
	return &svc.ServiceContext{
		ShipmentApp: mockApp,
	}
}

func TestCancelShipmentLogic_CancelShipment(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func(ctx context.Context) context.Context
		mockApp   *mockShipmentApp
		req       *types.CancelShipmentReq
		wantErr   error
		wantResp  bool
		respCheck func(t *testing.T, resp *types.CancelShipmentResp)
	}{
		{
			name: "successful cancellation",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 100)
				ctx = contextx.SetUserType(ctx, 0) // not platform admin
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusPending),
					}, nil
				},
				cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusCancelled),
					}, nil
				},
			},
			req: &types.CancelShipmentReq{
				ID:     1,
				Reason: "Customer requested cancellation",
			},
			wantErr:  nil,
			wantResp: true,
			respCheck: func(t *testing.T, resp *types.CancelShipmentResp) {
				if resp.ID != 1 {
					t.Errorf("ID = %v, want 1", resp.ID)
				}
				if resp.ShipmentNo != "SHP001" {
					t.Errorf("ShipmentNo = %v, want SHP001", resp.ShipmentNo)
				}
				if resp.Status != "cancelled" {
					t.Errorf("Status = %v, want cancelled", resp.Status)
				}
				if resp.StatusText != "cancelled" {
					t.Errorf("StatusText = %v, want cancelled", resp.StatusText)
				}
				if resp.CancelledAt == "" {
					t.Error("CancelledAt should not be empty")
				}
				if resp.Reason != "Customer requested cancellation" {
					t.Errorf("Reason = %v, want 'Customer requested cancellation'", resp.Reason)
				}
			},
		},
		{
			name: "shipment not found",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 100)
				ctx = contextx.SetUserType(ctx, 0)
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					return nil, code.ErrShipmentNotFound
				},
			},
			req: &types.CancelShipmentReq{
				ID:     999,
				Reason: "Test reason",
			},
			wantErr:  code.ErrShipmentNotFound,
			wantResp: false,
		},
		{
			name: "already cancelled - returns error from app layer",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 100)
				ctx = contextx.SetUserType(ctx, 0)
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusCancelled),
					}, nil
				},
				cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
					return nil, code.ErrShipmentAlreadyCancelled
				},
			},
			req: &types.CancelShipmentReq{
				ID:     1,
				Reason: "Trying to cancel again",
			},
			wantErr:  code.ErrShipmentAlreadyCancelled,
			wantResp: false,
		},
		{
			name: "delivered shipment cannot be cancelled",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 100)
				ctx = contextx.SetUserType(ctx, 0)
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusDelivered),
					}, nil
				},
				cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
					return nil, code.ErrShipmentCannotCancelDelivered
				},
			},
			req: &types.CancelShipmentReq{
				ID:     1,
				Reason: "Cancel after delivery",
			},
			wantErr:  code.ErrShipmentCannotCancelDelivered,
			wantResp: false,
		},
		{
			name: "platform admin can cancel shipment",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 999) // platform admin user
				ctx = contextx.SetUserType(ctx, 1) // platform admin type = 1
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					// Platform admin passes tenantID=0 to app layer
					if tenantID != 0 {
						t.Errorf("tenantID = %v, want 0 for platform admin", tenantID)
					}
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusPending),
					}, nil
				},
				cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusCancelled),
					}, nil
				},
			},
			req: &types.CancelShipmentReq{
				ID:     1,
				Reason: "Platform admin cancellation",
			},
			wantErr:  nil,
			wantResp: true,
		},
		{
			name: "cancel shipment returns generic error",
			setupCtx: func(ctx context.Context) context.Context {
				ctx = contextx.SetTenantID(ctx, 1)
				ctx = contextx.SetUserID(ctx, 100)
				ctx = contextx.SetUserType(ctx, 0)
				return ctx
			},
			mockApp: &mockShipmentApp{
				getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
					return &fulfillment.ShipmentResponse{
						ID:         1,
						ShipmentNo: "SHP001",
						Status:     int(domain.ShipmentStatusPending),
					}, nil
				},
				cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
					return nil, errors.New("database error")
				},
			},
			req: &types.CancelShipmentReq{
				ID:     1,
				Reason: "Test reason",
			},
			wantErr:  errors.New("database error"),
			wantResp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.setupCtx != nil {
				ctx = tt.setupCtx(ctx)
			}

			svcCtx := mockServiceContext(tt.mockApp)
			logic := NewCancelShipmentLogic(ctx, svcCtx)

			resp, err := logic.CancelShipment(tt.req)

			if tt.wantErr == nil && err != nil {
				t.Errorf("CancelShipment() error = %v, want nil", err)
				return
			}
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CancelShipment() error = nil, want %v", tt.wantErr)
					return
				}
				// For errors.New case, just check error message
				if tt.wantErr.Error() != err.Error() {
					t.Errorf("CancelShipment() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantResp {
				if resp == nil {
					t.Errorf("CancelShipment() resp = nil, want non-nil")
					return
				}
				if tt.respCheck != nil {
					tt.respCheck(t, resp)
				}
			}
		})
	}
}

func TestCancelShipmentLogic_CancelShipment_CancelledAtIsRFC3339Format(t *testing.T) {
	ctx := context.Background()
	ctx = contextx.SetTenantID(ctx, 1)
	ctx = contextx.SetUserID(ctx, 100)
	ctx = contextx.SetUserType(ctx, 0)

	mockApp := &mockShipmentApp{
		getShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, id int64) (*fulfillment.ShipmentResponse, error) {
			return &fulfillment.ShipmentResponse{
				ID:         1,
				ShipmentNo: "SHP001",
				Status:     int(domain.ShipmentStatusPending),
			}, nil
		},
		cancelShipmentFunc: func(ctx context.Context, tenantID shared.TenantID, userID int64, id int64, reason string) (*fulfillment.ShipmentResponse, error) {
			return &fulfillment.ShipmentResponse{
				ID:         1,
				ShipmentNo: "SHP001",
				Status:     int(domain.ShipmentStatusCancelled),
			}, nil
		},
	}

	svcCtx := mockServiceContext(mockApp)
	logic := NewCancelShipmentLogic(ctx, svcCtx)

	resp, err := logic.CancelShipment(&types.CancelShipmentReq{
		ID:     1,
		Reason: "Test reason",
	})

	if err != nil {
		t.Fatalf("CancelShipment() error = %v, want nil", err)
	}

	if resp.CancelledAt == "" {
		t.Fatal("CancelledAt is empty")
	}

	// Verify it's valid RFC3339 format
	_, err = time.Parse(time.RFC3339, resp.CancelledAt)
	if err != nil {
		t.Errorf("CancelledAt is not valid RFC3339 format: %v", err)
	}
}
