package persistence

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/colinrs/shopjoy/admin/internal/domain/media"
)

type mediaAssetRepo struct {
	db *gorm.DB
}

func NewMediaAssetRepository(db *gorm.DB) media.Repository {
	return &mediaAssetRepo{db: db}
}

func (r *mediaAssetRepo) Insert(ctx context.Context, a *media.Asset) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errors.New("duplicate media asset")
		}
		// Fallback: production MySQL driver wraps duplicate errors with this
		// message in some driver versions / mock drivers. Detect by substring
		// to keep behavior consistent across both real and mocked paths.
		if strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("duplicate media asset")
		}
		return err
	}
	return nil
}

func (r *mediaAssetRepo) FindByID(ctx context.Context, id int64) (*media.Asset, error) {
	var a media.Asset
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *mediaAssetRepo) FindByPublicID(ctx context.Context, provider, publicID string) (*media.Asset, error) {
	var a media.Asset
	if err := r.db.WithContext(ctx).
		Where("provider = ? AND public_id = ? AND deleted_at IS NULL", provider, publicID).
		First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *mediaAssetRepo) SoftDelete(ctx context.Context, id int64) error {
	// GORM's Delete on a soft-deletable model never returns ErrRecordNotFound
	// for a missing row — it returns nil with RowsAffected == 0. Use that signal
	// to distinguish missing-row from successfully-soft-deleted.
	res := r.db.WithContext(ctx).Delete(&media.Asset{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("asset not found")
	}
	return nil
}
