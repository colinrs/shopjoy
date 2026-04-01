package fulfillment

import (
	"testing"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

func TestShipment_CancelShipment(t *testing.T) {
	tests := []struct {
		name        string
		shipment    *Shipment
		reason      string
		cancelledBy int64
		wantErr     error
		validate    func(t *testing.T, s *Shipment)
	}{
		{
			name: "successful cancellation from pending status",
			shipment: &Shipment{
				Model:     application.Model{ID: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:  shared.TenantID(1),
				OrderID:   "ORDER001",
				ShipmentNo: "SHP001",
				Status:    ShipmentStatusPending,
			},
			reason:      "Customer requested cancellation",
			cancelledBy: 100,
			wantErr:     nil,
			validate: func(t *testing.T, s *Shipment) {
				if s.Status != ShipmentStatusCancelled {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusCancelled)
				}
				if s.CancelledAt == nil {
					t.Error("CancelledAt should not be nil")
				}
				if s.CancelledBy != 100 {
					t.Errorf("CancelledBy = %v, want 100", s.CancelledBy)
				}
				if s.CancelledReason != "Customer requested cancellation" {
					t.Errorf("CancelledReason = %v, want 'Customer requested cancellation'", s.CancelledReason)
				}
			},
		},
		{
			name: "successful cancellation from shipped status",
			shipment: &Shipment{
				Model:      application.Model{ID: 2, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:   shared.TenantID(1),
				OrderID:    "ORDER002",
				ShipmentNo: "SHP002",
				Status:     ShipmentStatusShipped,
				Carrier:    "SF Express",
				TrackingNo: "SF123456789",
			},
			reason:      "Wrong address provided",
			cancelledBy: 200,
			wantErr:     nil,
			validate: func(t *testing.T, s *Shipment) {
				if s.Status != ShipmentStatusCancelled {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusCancelled)
				}
				if s.CancelledAt == nil {
					t.Error("CancelledAt should not be nil")
				}
				if s.CancelledBy != 200 {
					t.Errorf("CancelledBy = %v, want 200", s.CancelledBy)
				}
			},
		},
		{
			name: "successful cancellation from in_transit status",
			shipment: &Shipment{
				Model:      application.Model{ID: 3, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:   shared.TenantID(1),
				OrderID:    "ORDER003",
				ShipmentNo: "SHP003",
				Status:     ShipmentStatusInTransit,
			},
			reason:      "Package lost",
			cancelledBy: 300,
			wantErr:     nil,
			validate: func(t *testing.T, s *Shipment) {
				if s.Status != ShipmentStatusCancelled {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusCancelled)
				}
			},
		},
		{
			name: "successful cancellation from failed status",
			shipment: &Shipment{
				Model:      application.Model{ID: 4, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:   shared.TenantID(1),
				OrderID:    "ORDER004",
				ShipmentNo: "SHP004",
				Status:     ShipmentStatusFailed,
			},
			reason:      "Order cancelled by customer",
			cancelledBy: 400,
			wantErr:     nil,
			validate: func(t *testing.T, s *Shipment) {
				if s.Status != ShipmentStatusCancelled {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusCancelled)
				}
			},
		},
		{
			name: "already cancelled - returns error",
			shipment: &Shipment{
				Model:         application.Model{ID: 5, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:      shared.TenantID(1),
				OrderID:       "ORDER005",
				ShipmentNo:    "SHP005",
				Status:        ShipmentStatusCancelled,
				CancelledAt:   func() *time.Time { t := time.Now().UTC(); return &t }(),
				CancelledBy:   500,
				CancelledReason: "Already cancelled reason",
			},
			reason:      "Trying to cancel again",
			cancelledBy: 600,
			wantErr:     code.ErrShipmentAlreadyCancelled,
			validate: func(t *testing.T, s *Shipment) {
				// Status should remain cancelled
				if s.Status != ShipmentStatusCancelled {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusCancelled)
				}
				// CancelledBy should not be updated
				if s.CancelledBy != 500 {
					t.Errorf("CancelledBy = %v, want 500 (original value)", s.CancelledBy)
				}
			},
		},
		{
			name: "delivered shipment - returns error",
			shipment: &Shipment{
				Model:       application.Model{ID: 6, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
				TenantID:    shared.TenantID(1),
				OrderID:     "ORDER006",
				ShipmentNo:  "SHP006",
				Status:      ShipmentStatusDelivered,
				DeliveredAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
			},
			reason:      "Customer wants to cancel after delivery",
			cancelledBy: 700,
			wantErr:     code.ErrShipmentCannotCancelDelivered,
			validate: func(t *testing.T, s *Shipment) {
				// Status should remain delivered
				if s.Status != ShipmentStatusDelivered {
					t.Errorf("Status = %v, want %v", s.Status, ShipmentStatusDelivered)
				}
				// Cancelled fields should not be set
				if s.CancelledAt != nil {
					t.Error("CancelledAt should remain nil")
				}
				if s.CancelledBy != 0 {
					t.Errorf("CancelledBy = %v, want 0", s.CancelledBy)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shipment.CancelShipment(tt.reason, tt.cancelledBy)

			if tt.wantErr == nil && err != nil {
				t.Errorf("CancelShipment() error = %v, want nil", err)
				return
			}
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("CancelShipment() error = nil, want %v", tt.wantErr)
					return
				}
				// Check if it's the expected error
				if err != tt.wantErr {
					t.Errorf("CancelShipment() error = %v, want %v", err, tt.wantErr)
				}
			}

			if tt.validate != nil {
				tt.validate(t, tt.shipment)
			}
		})
	}
}

func TestShipment_CancelShipment_UpdatedAtIsSet(t *testing.T) {
	shipment := &Shipment{
		Model:      application.Model{ID: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:   shared.TenantID(1),
		OrderID:    "ORDER001",
		ShipmentNo: "SHP001",
		Status:     ShipmentStatusPending,
	}
	originalUpdatedAt := shipment.UpdatedAt

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	err := shipment.CancelShipment("Test reason", 100)
	if err != nil {
		t.Fatalf("CancelShipment() error = %v, want nil", err)
	}

	if shipment.UpdatedAt.Equal(originalUpdatedAt) || shipment.UpdatedAt.Before(originalUpdatedAt) {
		t.Errorf("UpdatedAt was not updated: original=%v, current=%v", originalUpdatedAt, shipment.UpdatedAt)
	}
}

func TestShipment_CancelShipment_CancelledAtIsUTC(t *testing.T) {
	shipment := &Shipment{
		Model:      application.Model{ID: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:   shared.TenantID(1),
		OrderID:    "ORDER001",
		ShipmentNo: "SHP001",
		Status:     ShipmentStatusPending,
	}

	err := shipment.CancelShipment("Test reason", 100)
	if err != nil {
		t.Fatalf("CancelShipment() error = %v, want nil", err)
	}

	if shipment.CancelledAt == nil {
		t.Fatal("CancelledAt is nil")
	}

	// Check that CancelledAt is in UTC
	loc := shipment.CancelledAt.Location()
	if loc != time.UTC {
		t.Errorf("CancelledAt location = %v, want UTC", loc)
	}
}
