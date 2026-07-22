package shipping

import "testing"

// TestCarrierRegistry 验证 Register / Get / 未注册返回 false。
func TestCarrierRegistry(t *testing.T) {
	r := &CarrierRegistry{carriers: map[string]Carrier{}}

	// 未注册返回 false
	if _, ok := r.Get("standard"); ok {
		t.Error("expected Get on empty registry to return false")
	}

	// 注册后可取
	r.Register(StandardCarrier{})
	c, ok := r.Get("standard")
	if !ok {
		t.Fatal("expected Get('standard') to return true after Register")
	}
	if c.Code() != "standard" {
		t.Errorf("expected carrier code 'standard', got %q", c.Code())
	}

	// 未知 code 返回 false
	if _, ok := r.Get("fedex"); ok {
		t.Error("expected Get('fedex') to return false")
	}
}

// TestCarrierRegistry_DefaultsStandard 验证 NewCarrierRegistry 默认含 standard。
func TestCarrierRegistry_DefaultsStandard(t *testing.T) {
	r := NewCarrierRegistry()
	c, ok := r.Get("standard")
	if !ok {
		t.Fatal("expected NewCarrierRegistry to include 'standard'")
	}
	if c.Name() != "Standard Shipping" {
		t.Errorf("expected Name='Standard Shipping', got %q", c.Name())
	}
}
