package uploads_test

import (
	"context"
	"errors"
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/logic/uploads"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/snowflake"
)

// fakeSigner implements both storage.Storage and storage.Signer so it can be
// assigned to ServiceContext.Storage. Sign behaviour is configurable per test.
type fakeSigner struct {
	sig    storage.Signature
	err    error
	gotCtx context.Context
	gotP   storage.SignParams
}

func (f *fakeSigner) Sign(ctx context.Context, p storage.SignParams) (storage.Signature, error) {
	f.gotCtx = ctx
	f.gotP = p
	return f.sig, f.err
}

func (f *fakeSigner) Save(_ context.Context, _ storage.AssetDraft) (*storage.Asset, error) {
	return nil, nil
}

func (f *fakeSigner) RegisterAsset(_ context.Context, _ storage.RemoteAsset) (*storage.Asset, error) {
	return nil, nil
}

func (f *fakeSigner) Delete(_ context.Context, _ string) error                { return nil }
func (f *fakeSigner) Get(_ context.Context, _ string) (*storage.Asset, error) { return nil, nil }
func (f *fakeSigner) GetURL(_ context.Context, _ string) (string, error)      { return "", nil }

// fakeIDGen is a deterministic snowflake.Snowflake for tests.
type fakeIDGen struct{ id int64 }

func (f *fakeIDGen) GetNodeID() int64                                    { return 0 }
func (f *fakeIDGen) NextID(_ context.Context) (int64, error)             { return f.id, nil }
func (f *fakeIDGen) NextIDs(_ context.Context, n int64) ([]int64, error) { return []int64{f.id}, nil }

// nonSignerStorage implements storage.Storage but NOT storage.Signer.
type nonSignerStorage struct{}

func (nonSignerStorage) Save(_ context.Context, _ storage.AssetDraft) (*storage.Asset, error) {
	return nil, nil
}
func (nonSignerStorage) RegisterAsset(_ context.Context, _ storage.RemoteAsset) (*storage.Asset, error) {
	return nil, nil
}
func (nonSignerStorage) Delete(_ context.Context, _ string) error                { return nil }
func (nonSignerStorage) Get(_ context.Context, _ string) (*storage.Asset, error) { return nil, nil }
func (nonSignerStorage) GetURL(_ context.Context, _ string) (string, error)      { return "", nil }

func newSignServiceContext(s storage.Storage, idGen snowflake.Snowflake) *svc.ServiceContext {
	return &svc.ServiceContext{Storage: s, IDGen: idGen}
}

func TestUploadSignLogic_OK(t *testing.T) {
	signer := &fakeSigner{sig: storage.Signature{
		CloudName:    "demo",
		APIKey:       "k",
		Timestamp:    "100",
		Signature:    "abc",
		Folder:       "dev/1/product",
		PublicID:     "pid",
		UploadPreset: "preset",
	}}
	idGen := &fakeIDGen{id: 12345}
	svcCtx := newSignServiceContext(signer, idGen)

	ctx := contextx.SetTenantID(context.Background(), 7)
	ctx = contextx.SetUserID(ctx, 42)

	l := uploads.NewUploadSignLogic(ctx, svcCtx)
	got, err := l.UploadSign(&types.UploadSignRequest{Category: "product", Filename: "a.jpg"})
	if err != nil {
		t.Fatalf("UploadSign() unexpected err: %v", err)
	}

	if got.CloudName != "demo" {
		t.Errorf("CloudName = %q, want %q", got.CloudName, "demo")
	}
	if got.PublicID != "pid" {
		t.Errorf("PublicID = %q, want %q", got.PublicID, "pid")
	}
	if got.Signature != "abc" {
		t.Errorf("Signature = %q, want %q", got.Signature, "abc")
	}
	if got.Folder != "dev/1/product" {
		t.Errorf("Folder = %q, want %q", got.Folder, "dev/1/product")
	}
	if got.APIKey != "k" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "k")
	}
	if got.Timestamp != "100" {
		t.Errorf("Timestamp = %q, want %q", got.Timestamp, "100")
	}
	if got.UploadPreset != "preset" {
		t.Errorf("UploadPreset = %q, want %q", got.UploadPreset, "preset")
	}
	if got.AssetID != "12345" {
		t.Errorf("AssetID = %q, want %q", got.AssetID, "12345")
	}
	wantURL := "https://api.cloudinary.com/v1_1/demo/image/upload"
	if got.UploadURL != wantURL {
		t.Errorf("UploadURL = %q, want %q", got.UploadURL, wantURL)
	}

	// Verify Signer received correct SignParams.
	if signer.gotP.TenantID != 7 {
		t.Errorf("Sign got TenantID = %d, want 7", signer.gotP.TenantID)
	}
	if signer.gotP.Category != storage.CategoryProduct {
		t.Errorf("Sign got Category = %q, want %q", signer.gotP.Category, storage.CategoryProduct)
	}
	if signer.gotP.Filename != "a.jpg" {
		t.Errorf("Sign got Filename = %q, want %q", signer.gotP.Filename, "a.jpg")
	}
	if signer.gotP.Timestamp <= 0 {
		t.Errorf("Sign got Timestamp = %d, want > 0", signer.gotP.Timestamp)
	}
}

func TestUploadSignLogic_DefaultCategory(t *testing.T) {
	signer := &fakeSigner{sig: storage.Signature{
		CloudName: "demo", APIKey: "k", Timestamp: "1",
		Signature: "s", Folder: "f", PublicID: "p",
	}}
	svcCtx := newSignServiceContext(signer, &fakeIDGen{id: 1})

	ctx := contextx.SetTenantID(context.Background(), 1)
	ctx = contextx.SetUserID(ctx, 1)
	l := uploads.NewUploadSignLogic(ctx, svcCtx)
	if _, err := l.UploadSign(&types.UploadSignRequest{}); err != nil {
		t.Fatalf("UploadSign() unexpected err: %v", err)
	}
	if signer.gotP.Category != storage.CategoryProduct {
		t.Errorf("default Category = %q, want %q", signer.gotP.Category, storage.CategoryProduct)
	}
}

func TestUploadSignLogic_BadStorage(t *testing.T) {
	signer := &fakeSigner{err: errors.New("boom")}
	svcCtx := newSignServiceContext(signer, &fakeIDGen{id: 1})

	ctx := contextx.SetTenantID(context.Background(), 1)
	ctx = contextx.SetUserID(ctx, 1)
	l := uploads.NewUploadSignLogic(ctx, svcCtx)
	_, err := l.UploadSign(&types.UploadSignRequest{Category: "product"})

	if !errors.Is(err, code.ErrUploadSignFailed) {
		t.Fatalf("UploadSign() err = %v, want %v", err, code.ErrUploadSignFailed)
	}
}

func TestUploadSignLogic_NonSignerStorage(t *testing.T) {
	// A storage implementation that doesn't support Sign (e.g., local-only)
	// must surface a provider error.
	svcCtx := newSignServiceContext(nonSignerStorage{}, &fakeIDGen{id: 1})

	ctx := contextx.SetTenantID(context.Background(), 1)
	ctx = contextx.SetUserID(ctx, 1)
	l := uploads.NewUploadSignLogic(ctx, svcCtx)
	_, err := l.UploadSign(&types.UploadSignRequest{Category: "product"})
	if !errors.Is(err, code.ErrUploadProviderError) {
		t.Fatalf("UploadSign() err = %v, want %v", err, code.ErrUploadProviderError)
	}
}

// TestUploadSignLogic_MissingAuthContext verifies that the cross-tenant guard
// rejects forged (unauthenticated) requests, even before signing.
func TestUploadSignLogic_MissingAuthContext(t *testing.T) {
	signer := &fakeSigner{sig: storage.Signature{
		CloudName: "demo", APIKey: "k", Timestamp: "1",
		Signature: "s", Folder: "f", PublicID: "p",
	}}
	svcCtx := newSignServiceContext(signer, &fakeIDGen{id: 1})

	cases := []struct {
		name string
		ctx  context.Context
	}{
		{"no-tenant-no-user", context.Background()},
		{"tenant-only", contextx.SetTenantID(context.Background(), 1)},
		{"user-only", contextx.SetUserID(context.Background(), 1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := uploads.NewUploadSignLogic(tc.ctx, svcCtx)
			_, err := l.UploadSign(&types.UploadSignRequest{Category: "product"})
			if !errors.Is(err, code.ErrUploadCrossTenantAccess) {
				t.Fatalf("err = %v, want %v", err, code.ErrUploadCrossTenantAccess)
			}
		})
	}
}