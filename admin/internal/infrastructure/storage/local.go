package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/media"
	snowflake "github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type localStorage struct {
	basePath string
	repo     media.Repository
	idGen    snowflake.Snowflake
}

func newLocal(cfg Config, repo media.Repository, idGen snowflake.Snowflake) *localStorage {
	base := cfg.Local.BasePath
	if base == "" {
		base = "./uploads"
	}
	return &localStorage{basePath: base, repo: repo, idGen: idGen}
}

func (s *localStorage) Save(ctx context.Context, draft AssetDraft) (*Asset, error) {
	if draft.Reader == nil {
		return nil, errors.New("reader is nil")
	}
	ext := strings.ToLower(filepath.Ext(draft.Filename))
	if ext == "" || !strings.Contains(".jpg.jpeg.png.gif.webp", ext) {
		return nil, errors.New("unsupported extension")
	}

	now := time.Now().UTC()
	dir := filepath.Join(s.basePath, string(draft.Category),
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		fmt.Sprintf("%02d", now.Day()),
	)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	assetID, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate asset id: %w", err)
	}
	assetIDStr := fmt.Sprintf("img_%d", assetID)
	filePath := filepath.Join(dir, assetIDStr+ext)

	absBase, _ := filepath.Abs(s.basePath)
	absFile, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFile, absBase) {
		return nil, errors.New("path traversal blocked")
	}

	buf, err := io.ReadAll(draft.Reader)
	if err != nil {
		return nil, fmt.Errorf("buffer reader: %w", err)
	}

	dst, err := os.Create(filePath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	width, height := 0, 0
	if ext != ".gif" && ext != ".webp" {
		if cfg, _, derr := image.DecodeConfig(bytes.NewReader(buf)); derr == nil {
			width, height = cfg.Width, cfg.Height
		}
	}
	written, err := dst.Write(buf)
	if err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	relPath := fmt.Sprintf("/uploads/%s/%d/%02d/%02d/%s%s",
		draft.Category, now.Year(), int(now.Month()), now.Day(), assetIDStr, ext)

	asset := &media.Asset{
		PublicID:  relPath,
		URL:       relPath,
		Filename:  draft.Filename,
		SizeBytes: int64(written), // populate from actual writer result
		MimeType:  draft.MimeType,
		Width:     width,
		Height:    height,
		Format:    strings.TrimPrefix(ext, "."),
		Category:  string(draft.Category),
		Provider:  "local",
		TenantID:  draft.TenantID,
		CreatedBy: draft.CreatedBy,
	}
	asset.ID = assetID
	if err := s.repo.Insert(ctx, asset); err != nil {
		return nil, fmt.Errorf("insert asset: %w", err)
	}
	return &Asset{
		ID:        fmt.Sprintf("%d", assetID),
		PublicID:  relPath,
		URL:       relPath,
		Filename:  draft.Filename,
		Size:      int64(written),
		MimeType:  draft.MimeType,
		Width:     width,
		Height:    height,
		Format:    asset.Format,
		Category:  draft.Category,
		Provider:  "local",
		TenantID:  draft.TenantID,
		CreatedBy: draft.CreatedBy,
		CreatedAt: now,
	}, nil
}

func (s *localStorage) RegisterAsset(ctx context.Context, remote RemoteAsset) (*Asset, error) {
	// Local path: file is already on disk via prior Save() call; this method
	// for the local driver just records the metadata if not already recorded.
	assetID, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate asset id: %w", err)
	}
	asset := &media.Asset{
		PublicID:  remote.PublicID,
		URL:       remote.URL,
		Filename:  remote.Filename,
		SizeBytes: remote.Size,
		MimeType:  remote.MimeType,
		Width:     remote.Width,
		Height:    remote.Height,
		Format:    remote.Format,
		Category:  string(remote.Category),
		Provider:  "local",
		TenantID:  remote.TenantID,
		CreatedBy: remote.CreatedBy,
	}
	asset.ID = assetID
	if err := s.repo.Insert(ctx, asset); err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &Asset{
		ID:        fmt.Sprintf("%d", assetID),
		PublicID:  remote.PublicID,
		URL:       remote.URL,
		Filename:  remote.Filename,
		Size:      remote.Size,
		MimeType:  remote.MimeType,
		Width:     remote.Width,
		Height:    remote.Height,
		Format:    remote.Format,
		Category:  remote.Category,
		Provider:  "local",
		TenantID:  remote.TenantID,
		CreatedBy: remote.CreatedBy,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *localStorage) Delete(ctx context.Context, id string) error {
	asset, err := s.lookupByID(ctx, id)
	if err != nil {
		return err
	}
	// Best-effort unlink; ignore not-found errors on disk.
	full := filepath.Join(s.basePath, asset.PublicID)
	rel, err := filepath.Rel(s.basePath, full)
	if err != nil {
		logx.WithContext(ctx).Errorf("compute relative path for %s: %v", full, err)
		return err
	}
	if err := os.Remove(filepath.Join(s.basePath, rel)); err != nil && !errors.Is(err, os.ErrNotExist) {
		logx.WithContext(ctx).Errorf("local delete file failed: id=%s path=%s: %v", id, rel, err)
		// continue — best-effort. DB soft-delete still happens.
	}
	return s.repo.SoftDelete(ctx, asset.ID)
}

// DeleteByTenant soft-deletes the asset when both id and tenantID match.
// The repository's DeleteByTenant collapses missing-row and cross-tenant
// into one ErrMediaAssetNotFound signal to prevent IDOR-style existence
// leaks. The local driver then unlinks the file as a best-effort step.
func (s *localStorage) DeleteByTenant(ctx context.Context, id string, tenantID int64) error {
	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		return errors.New("invalid id")
	}
	if err := s.repo.DeleteByTenant(ctx, idInt, tenantID); err != nil {
		return err
	}
	// File path includes upload-date components we don't have without a
	// lookup, so do a follow-up read for unlink only (ignore lookup errors
	// — DB already soft-deleted; orphan file on disk is acceptable cleanup).
	asset, lookupErr := s.repo.FindByID(ctx, idInt)
	if lookupErr == nil {
		full := filepath.Join(s.basePath, asset.PublicID)
		if err := os.Remove(full); err != nil && !errors.Is(err, os.ErrNotExist) {
			logx.WithContext(ctx).Errorf("local DeleteByTenant file unlink failed: id=%d path=%s: %v", idInt, asset.PublicID, err)
		}
	}
	return nil
}

func (s *localStorage) Get(ctx context.Context, id string) (*Asset, error) {
	asset, err := s.lookupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Asset{
		ID:        fmt.Sprintf("%d", asset.ID),
		PublicID:  asset.PublicID,
		URL:       asset.URL,
		Filename:  asset.Filename,
		Size:      asset.SizeBytes,
		MimeType:  asset.MimeType,
		Width:     asset.Width,
		Height:    asset.Height,
		Format:    asset.Format,
		Category:  Category(asset.Category),
		Provider:  asset.Provider,
		TenantID:  asset.TenantID,
		CreatedBy: asset.CreatedBy,
		CreatedAt: asset.CreatedAt,
	}, nil
}

func (s *localStorage) GetURL(ctx context.Context, id string) (string, error) {
	asset, err := s.lookupByID(ctx, id)
	if err != nil {
		return "", err
	}
	return asset.URL, nil
}

// lookupByID parses id as int64 and fetches via repository.
func (s *localStorage) lookupByID(ctx context.Context, id string) (*media.Asset, error) {
	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(ctx, idInt)
}
