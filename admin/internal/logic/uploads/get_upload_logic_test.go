package uploads_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/logic/uploads"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
)

// getUploadStorage is a configurable fake of storage.Storage used by the
// GetUpload tests. It records the most recent call so tests can verify the
// logic hands the right ID to the storage layer.
type getUploadStorage struct {
	gotID  string
	ret    *storage.Asset
	retErr error
}

func (g *getUploadStorage) Save(_ context.Context, _ storage.AssetDraft) (*storage.Asset, error) {
	return nil, nil
}
func (g *getUploadStorage) RegisterAsset(_ context.Context, _ storage.RemoteAsset) (*storage.Asset, error) {
	return nil, nil
}
func (g *getUploadStorage) Delete(_ context.Context, _ string) error       { return nil }
func (g *getUploadStorage) DeleteByTenant(_ context.Context, _ string, _ int64) error {
	return nil
}
func (g *getUploadStorage) Get(_ context.Context, id string) (*storage.Asset, error) {
	g.gotID = id
	return g.ret, g.retErr
}
func (g *getUploadStorage) GetURL(_ context.Context, _ string) (string, error) {
	return "", nil
}

func newGetServiceContext(s storage.Storage) *svc.ServiceContext {
	return &svc.ServiceContext{Storage: s}
}

func TestGetUploadLogic_OK(t *testing.T) {
	fixed := time.Date(2026, 7, 11, 8, 0, 0, 0, time.UTC)
	store := &getUploadStorage{
		ret: &storage.Asset{
			ID:        "asset-1",
			URL:       "https://cdn/x.jpg",
			Filename:  "x.jpg",
			Category:  storage.CategoryProduct,
			Size:      1234,
			MimeType:  "image/jpeg",
			Width:     100,
			Height:    200,
			TenantID:  7,
			CreatedAt: fixed,
		},
	}
	svcCtx := newGetServiceContext(store)

	// Same tenant on context and asset — happy path.
	ctx := contextx.SetTenantID(context.Background(), 7)

	l := uploads.NewGetUploadLogic(ctx, svcCtx)
	got, err := l.GetUpload(&types.GetUploadReq{ID: "asset-1"})
	if err != nil {
		t.Fatalf("GetUpload() unexpected err: %v", err)
	}
	if got == nil {
		t.Fatalf("GetUpload() returned nil resp")
	}
	if got.ID != "asset-1" {
		t.Errorf("ID = %q, want %q", got.ID, "asset-1")
	}
	if got.URL != "https://cdn/x.jpg" {
		t.Errorf("URL = %q, want %q", got.URL, "https://cdn/x.jpg")
	}
	if got.Filename != "x.jpg" {
		t.Errorf("Filename = %q, want %q", got.Filename, "x.jpg")
	}
	if got.Category != "product" {
		t.Errorf("Category = %q, want %q", got.Category, "product")
	}
	if got.Size != 1234 {
		t.Errorf("Size = %d, want 1234", got.Size)
	}
	if got.MimeType != "image/jpeg" {
		t.Errorf("MimeType = %q, want %q", got.MimeType, "image/jpeg")
	}
	if got.Width != 100 {
		t.Errorf("Width = %d, want 100", got.Width)
	}
	if got.Height != 200 {
		t.Errorf("Height = %d, want 200", got.Height)
	}
	if got.CreatedAt != "2026-07-11T08:00:00Z" {
		t.Errorf("CreatedAt = %q, want %q", got.CreatedAt, "2026-07-11T08:00:00Z")
	}
	if store.gotID != "asset-1" {
		t.Errorf("Storage.Get got ID = %q, want %q", store.gotID, "asset-1")
	}
}

func TestGetUploadLogic_CrossTenant(t *testing.T) {
	// asset owned by tenant 7, request issued under tenant 1 — should be denied.
	store := &getUploadStorage{
		ret: &storage.Asset{
			ID:       "asset-1",
			TenantID: 7,
		},
	}
	svcCtx := newGetServiceContext(store)

	ctx := contextx.SetTenantID(context.Background(), 1)

	l := uploads.NewGetUploadLogic(ctx, svcCtx)
	got, err := l.GetUpload(&types.GetUploadReq{ID: "asset-1"})
	if got != nil {
		t.Errorf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadCrossTenantAccess) {
		t.Fatalf("GetUpload() err = %v, want %v", err, code.ErrUploadCrossTenantAccess)
	}
}

func TestGetUploadLogic_NotFound(t *testing.T) {
	// Storage returns an error → GetUpload should surface ErrUploadNotFound
	// regardless of the underlying cause.
	store := &getUploadStorage{retErr: errors.New("missing")}
	svcCtx := newGetServiceContext(store)

	l := uploads.NewGetUploadLogic(context.Background(), svcCtx)
	got, err := l.GetUpload(&types.GetUploadReq{ID: "missing"})
	if got != nil {
		t.Errorf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadNotFound) {
		t.Fatalf("GetUpload() err = %v, want %v", err, code.ErrUploadNotFound)
	}
}
