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

// fakeStorage implements storage.Storage with no signing capability. It
// returns success on RegisterAsset so the happy path can be exercised.
type fakeStorage struct {
	gotRemote storage.RemoteAsset
	gotCtx    context.Context
	ret       *storage.Asset
	retErr    error
}

func (f *fakeStorage) Save(_ context.Context, _ storage.AssetDraft) (*storage.Asset, error) {
	return nil, nil
}
func (f *fakeStorage) RegisterAsset(ctx context.Context, remote storage.RemoteAsset) (*storage.Asset, error) {
	f.gotCtx = ctx
	f.gotRemote = remote
	return f.ret, f.retErr
}
func (f *fakeStorage) Delete(_ context.Context, _ string) error                  { return nil }
func (f *fakeStorage) DeleteByTenant(_ context.Context, _ string, _ int64) error { return nil }
func (f *fakeStorage) Get(_ context.Context, _ string) (*storage.Asset, error) {
	return nil, nil
}
func (f *fakeStorage) GetURL(_ context.Context, _ string) (string, error) { return "", nil }

// folderMismatchStorage pretends the storage layer rejected the asset because
// its public_id folder didn't match the signed folder.
type folderMismatchStorage struct{ fakeStorage }

func (f *folderMismatchStorage) RegisterAsset(ctx context.Context, remote storage.RemoteAsset) (*storage.Asset, error) {
	f.gotCtx = ctx
	f.gotRemote = remote
	return nil, errors.New("public_id folder mismatch")
}

func newConfirmServiceContext(s storage.Storage) *svc.ServiceContext {
	return &svc.ServiceContext{Storage: s}
}

func TestConfirmLogic_BadMime(t *testing.T) {
	svcCtx := newConfirmServiceContext(&fakeStorage{})
	l := uploads.NewUploadConfirmLogic(context.Background(), svcCtx)

	got, err := l.UploadConfirm(&types.UploadConfirmRequest{
		PublicID: "dev/1/p/pid",
		URL:      "u",
		MimeType: "application/zip",
		Category: "product",
	})
	if got != nil {
		t.Fatalf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadUnsupportedFileType) {
		t.Fatalf("expected ErrUploadUnsupportedFileType, got %v", err)
	}
}

func TestConfirmLogic_EmptyPublicID(t *testing.T) {
	svcCtx := newConfirmServiceContext(&fakeStorage{})
	l := uploads.NewUploadConfirmLogic(context.Background(), svcCtx)

	got, err := l.UploadConfirm(&types.UploadConfirmRequest{
		PublicID: "",
		URL:      "u",
		MimeType: "image/jpeg",
		Category: "product",
	})
	if got != nil {
		t.Fatalf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadConfirmFailed) {
		t.Fatalf("expected ErrUploadConfirmFailed, got %v", err)
	}
}

func TestConfirmLogic_EmptyURL(t *testing.T) {
	svcCtx := newConfirmServiceContext(&fakeStorage{})
	l := uploads.NewUploadConfirmLogic(context.Background(), svcCtx)

	got, err := l.UploadConfirm(&types.UploadConfirmRequest{
		PublicID: "dev/1/p/pid",
		URL:      "",
		MimeType: "image/jpeg",
		Category: "product",
	})
	if got != nil {
		t.Fatalf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadConfirmFailed) {
		t.Fatalf("expected ErrUploadConfirmFailed, got %v", err)
	}
}

func TestConfirmLogic_PathTraversal(t *testing.T) {
	svcCtx := newConfirmServiceContext(&fakeStorage{})
	l := uploads.NewUploadConfirmLogic(context.Background(), svcCtx)

	cases := map[string]string{
		"dotdot": "../etc/passwd",
		"backsl": `dev\1\p\pid`,
		"mixed":  "dev/../../pid",
	}
	for name, publicID := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := l.UploadConfirm(&types.UploadConfirmRequest{
				PublicID: publicID,
				URL:      "u",
				MimeType: "image/jpeg",
				Category: "product",
			})
			if got != nil {
				t.Fatalf("expected nil resp, got %+v", got)
			}
			if !errors.Is(err, code.ErrUploadConfirmFailed) {
				t.Fatalf("expected ErrUploadConfirmFailed, got %v", err)
			}
		})
	}
}

func TestConfirmLogic_PublicIdFolderMismatch(t *testing.T) {
	s := &folderMismatchStorage{}
	svcCtx := newConfirmServiceContext(s)
	l := uploads.NewUploadConfirmLogic(context.Background(), svcCtx)

	got, err := l.UploadConfirm(&types.UploadConfirmRequest{
		PublicID: "different/folder/pid",
		URL:      "u",
		MimeType: "image/jpeg",
		Category: "product",
		Width:    1,
		Height:   1,
		Size:     1,
	})
	if got != nil {
		t.Fatalf("expected nil resp, got %+v", got)
	}
	if !errors.Is(err, code.ErrUploadConfirmFailed) {
		t.Fatalf("expected ErrUploadConfirmFailed, got %v", err)
	}
}

func TestConfirmLogic_OK(t *testing.T) {
	fixed := time.Date(2026, 7, 11, 8, 0, 0, 0, time.UTC)
	fs := &fakeStorage{ret: &storage.Asset{
		ID:        "asset-1",
		URL:       "https://cdn/x.jpg",
		Filename:  "x.jpg",
		Category:  storage.CategoryProduct,
		Size:      1234,
		MimeType:  "image/jpeg",
		Width:     100,
		Height:    200,
		CreatedAt: fixed,
	}}
	svcCtx := newConfirmServiceContext(fs)

	ctx := contextx.SetTenantID(context.Background(), 7)
	ctx = contextx.SetUserID(ctx, 42)

	l := uploads.NewUploadConfirmLogic(ctx, svcCtx)
	got, err := l.UploadConfirm(&types.UploadConfirmRequest{
		PublicID: "dev/7/product/pid",
		URL:      "https://cdn/x.jpg",
		Filename: "x.jpg",
		Size:     1234,
		MimeType: "image/jpeg",
		Width:    100,
		Height:   200,
		Format:   "jpg",
		Category: "product",
	})
	if err != nil {
		t.Fatalf("UploadConfirm() unexpected err: %v", err)
	}
	if got == nil {
		t.Fatalf("UploadConfirm() returned nil resp")
	}
	if got.ID != "asset-1" {
		t.Errorf("ID = %q, want %q", got.ID, "asset-1")
	}
	if got.CreatedAt != "2026-07-11T08:00:00Z" {
		t.Errorf("CreatedAt = %q, want %q", got.CreatedAt, "2026-07-11T08:00:00Z")
	}
	if fs.gotRemote.TenantID != 7 {
		t.Errorf("RegisterAsset got TenantID = %d, want 7", fs.gotRemote.TenantID)
	}
	if fs.gotRemote.CreatedBy != 42 {
		t.Errorf("RegisterAsset got CreatedBy = %d, want 42", fs.gotRemote.CreatedBy)
	}
	if fs.gotRemote.PublicID != "dev/7/product/pid" {
		t.Errorf("RegisterAsset got PublicID = %q", fs.gotRemote.PublicID)
	}
}
