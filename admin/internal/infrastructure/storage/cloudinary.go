package storage

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/colinrs/shopjoy/admin/internal/domain/media"
	snowflake "github.com/colinrs/shopjoy/pkg/snowflake"
)

type cloudinaryStorage struct {
	cloud *cloudinary.Cloudinary
	cfg   CloudinaryConfig
	repo  media.Repository
	idGen snowflake.Snowflake
}

func newCloudinary(cfg Config, repo media.Repository, idGen snowflake.Snowflake) (*cloudinaryStorage, error) {
	if cfg.Cloudinary.CloudName == "" || cfg.Cloudinary.APIKey == "" || cfg.Cloudinary.APISecret == "" {
		return nil, errors.New("cloudinary config incomplete")
	}
	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary client: %w", err)
	}
	return &cloudinaryStorage{
		cloud: cld,
		cfg:   cfg.Cloudinary,
		repo:  repo,
		idGen: idGen,
	}, nil
}

// Sign produces a Cloudinary-compatible signed bundle for browser direct upload.
func (s *cloudinaryStorage) Sign(ctx context.Context, p SignParams) (Signature, error) {
	if s.cfg.APISecret == "" {
		return Signature{}, errors.New("cloudinary api_secret missing")
	}
	ts := p.Timestamp
	if ts == 0 {
		ts = time.Now().UTC().Unix()
	}
	folder := s.folder(p.Category, p.TenantID)
	publicID := uuidV4()

	params := url.Values{}
	params.Set("folder", folder)
	params.Set("public_id", publicID)
	params.Set("timestamp", strconv.FormatInt(ts, 10))

	sig := sha1Hex(canonicalParams(params) + s.cfg.APISecret)

	return Signature{
		CloudName:    s.cfg.CloudName,
		APIKey:       s.cfg.APIKey,
		Timestamp:    strconv.FormatInt(ts, 10),
		Signature:    sig,
		Folder:       folder,
		PublicID:     publicID,
		UploadPreset: s.cfg.UploadPreset,
	}, nil
}

// Save performs a server-side proxy upload via the Cloudinary SDK. The primary
// flow uses direct browser uploads + RegisterAsset; Save remains for callers
// that prefer server-side uploads.
func (s *cloudinaryStorage) Save(ctx context.Context, draft AssetDraft) (*Asset, error) {
	if draft.Reader == nil {
		return nil, errors.New("reader nil")
	}
	folder := s.folder(draft.Category, draft.TenantID)
	publicID := uuidV4()
	res, err := s.cloud.Upload.Upload(ctx, draft.Reader, uploader.UploadParams{
		Folder:   folder,
		PublicID: publicID,
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload: %w", err)
	}
	assetID, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate asset id: %w", err)
	}
	a := &media.Asset{
		PublicID:  res.PublicID,
		URL:       res.SecureURL,
		Filename:  draft.Filename,
		SizeBytes: int64(res.Bytes),
		MimeType:  draft.MimeType,
		Width:     res.Width,
		Height:    res.Height,
		Format:    res.Format,
		Category:  string(draft.Category),
		Provider:  "cloudinary",
		TenantID:  draft.TenantID,
		CreatedBy: draft.CreatedBy,
	}
	a.ID = assetID
	if err := s.repo.Insert(ctx, a); err != nil {
		return nil, fmt.Errorf("insert asset: %w", err)
	}
	return s.assetFromDomain(a), nil
}

// RegisterAsset records a Cloudinary asset whose bytes were uploaded directly
// from the browser. The folder prefix is validated to prevent callers from
// registering assets under a foreign tenant path.
func (s *cloudinaryStorage) RegisterAsset(ctx context.Context, remote RemoteAsset) (*Asset, error) {
	if !strings.HasPrefix(remote.PublicID, s.folder(remote.Category, remote.TenantID)) {
		return nil, errors.New("public_id folder mismatch")
	}
	assetID, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate asset id: %w", err)
	}
	a := &media.Asset{
		PublicID:  remote.PublicID,
		URL:       remote.URL,
		Filename:  remote.Filename,
		SizeBytes: remote.Size,
		MimeType:  remote.MimeType,
		Width:     remote.Width,
		Height:    remote.Height,
		Format:    remote.Format,
		Category:  string(remote.Category),
		Provider:  "cloudinary",
		TenantID:  remote.TenantID,
		CreatedBy: remote.CreatedBy,
	}
	a.ID = assetID
	if err := s.repo.Insert(ctx, a); err != nil {
		return nil, fmt.Errorf("insert asset: %w", err)
	}
	return s.assetFromDomain(a), nil
}

func (s *cloudinaryStorage) Delete(ctx context.Context, id string) error {
	a, err := s.lookupByID(ctx, id)
	if err != nil {
		return err
	}
	if _, derr := s.cloud.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: a.PublicID}); derr != nil {
		// Best-effort remote cleanup; DB soft delete still happens below.
		// Failed destroys are surfaced via logx for Operations to act on.
		logx.WithContext(ctx).Errorf("cloudinary destroy failed: public_id=%s err=%v", a.PublicID, derr)
	}
	return s.repo.SoftDelete(ctx, a.ID)
}

func (s *cloudinaryStorage) Get(ctx context.Context, id string) (*Asset, error) {
	a, err := s.lookupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.assetFromDomain(a), nil
}

func (s *cloudinaryStorage) GetURL(ctx context.Context, id string) (string, error) {
	a, err := s.lookupByID(ctx, id)
	if err != nil {
		return "", err
	}
	return a.URL, nil
}

func (s *cloudinaryStorage) folder(c Category, tenantID int64) string {
	return strings.Join([]string{s.cfg.Environment, strconv.FormatInt(tenantID, 10), string(c)}, "/")
}

func (s *cloudinaryStorage) lookupByID(ctx context.Context, id string) (*media.Asset, error) {
	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(ctx, idInt)
}

func (s *cloudinaryStorage) assetFromDomain(a *media.Asset) *Asset {
	return &Asset{
		ID:        strconv.FormatInt(a.ID, 10),
		PublicID:  a.PublicID,
		URL:       a.URL,
		Filename:  a.Filename,
		Size:      a.SizeBytes,
		MimeType:  a.MimeType,
		Width:     a.Width,
		Height:    a.Height,
		Format:    a.Format,
		Category:  Category(a.Category),
		Provider:  a.Provider,
		TenantID:  a.TenantID,
		CreatedBy: a.CreatedBy,
		CreatedAt: a.CreatedAt,
	}
}

// canonicalParams joins a sorted url.Values as "k=v&k=v".
func canonicalParams(v url.Values) string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+v.Get(k))
	}
	return strings.Join(parts, "&")
}

func sha1Hex(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// uuidV4 returns a fresh UUID v4 string.
func uuidV4() string { return uuid.NewString() }
