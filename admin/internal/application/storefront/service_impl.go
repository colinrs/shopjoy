package storefront

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

type pageService struct {
	db            *gorm.DB
	pageRepo      storefront.PageRepository
	decorationRepo storefront.DecorationRepository
	versionRepo   storefront.PageVersionRepository
	idGen         snowflake.Snowflake
}

func NewPageService(
	db *gorm.DB,
	pageRepo storefront.PageRepository,
	decorationRepo storefront.DecorationRepository,
	versionRepo storefront.PageVersionRepository,
	idGen snowflake.Snowflake,
) PageService {
	return &pageService{
		db:            db,
		pageRepo:      pageRepo,
		decorationRepo: decorationRepo,
		versionRepo:   versionRepo,
		idGen:         idGen,
	}
}

func (s *pageService) ListPages(ctx context.Context, tenantID shared.TenantID, page, pageSize int) (*PaginatedResult[*PageDTO], error) {
	pages, total, err := s.pageRepo.FindAll(ctx, s.db, tenantID, page, pageSize)
	if err != nil {
		return nil, err
	}

	dtos := make([]*PageDTO, len(pages))
	for i, p := range pages {
		dtos[i] = &PageDTO{
			ID:          p.Model.ID,
			PageType:    pageTypeToString(p.Type),
			Name:        p.Name,
			Slug:        p.Slug,
			IsPublished: p.IsPublished,
			Version:     p.Version,
		}
	}
	return &PaginatedResult[*PageDTO]{
		Items:    dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *pageService) GetPage(ctx context.Context, tenantID shared.TenantID, pageID int64) (*PageDetailDTO, error) {
	page, err := s.pageRepo.FindByID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, code.ErrPageNotFound
	}

	decorations, err := s.decorationRepo.FindByPageID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return nil, err
	}

	return &PageDetailDTO{
		Page: &PageDTO{
			ID:          page.Model.ID,
			PageType:    pageTypeToString(page.Type),
			Name:        page.Name,
			Slug:        page.Slug,
			IsPublished: page.IsPublished,
			Version:     page.Version,
		},
		Decorations: decorationsToDTOs(decorations),
	}, nil
}

func (s *pageService) GetPageBySlug(ctx context.Context, tenantID shared.TenantID, slug string) (*PageDetailDTO, error) {
	page, err := s.pageRepo.FindBySlug(ctx, s.db, tenantID, slug)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, code.ErrPageNotFound
	}

	return s.GetPage(ctx, tenantID, page.Model.ID)
}

func (s *pageService) SaveDraft(ctx context.Context, tenantID shared.TenantID, pageID int64, blocks []*DecorationDTO, userID int64) error {
	page, err := s.pageRepo.FindByID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return err
	}
	if page == nil {
		return code.ErrPageNotFound
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing decorations
		if err := s.decorationRepo.DeleteByPageID(ctx, tx, tenantID, pageID); err != nil {
			return err
		}

		// Create new decorations
		for i, block := range blocks {
			id, err := s.idGen.NextID(ctx)
			if err != nil {
				return err
			}

			sortOrder := block.SortOrder
			if sortOrder == 0 {
				sortOrder = i
			}

			decoration := &storefront.Decoration{
				Model:       application.Model{ID: id},
				TenantID:    tenantID,
				PageID:      pageID,
				BlockType:   block.BlockType,
				BlockConfig: block.BlockConfig,
				SortOrder:   sortOrder,
				IsActive:    true,
			}

			if err := s.decorationRepo.Create(ctx, tx, decoration); err != nil {
				return err
			}
		}

		// Increment version
		page.Version++
		if err := s.pageRepo.Update(ctx, tx, page); err != nil {
			return err
		}

		return nil
	})
}

func (s *pageService) PublishPage(ctx context.Context, tenantID shared.TenantID, pageID int64, userID int64) error {
	page, err := s.pageRepo.FindByID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return err
	}
	if page == nil {
		return code.ErrPageNotFound
	}

	if page.IsPublished {
		return code.ErrPageAlreadyPublished
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get current decorations for version snapshot
		decorations, err := s.decorationRepo.FindByPageID(ctx, tx, tenantID, pageID)
		if err != nil {
			return err
		}

		// Create version snapshot
		id, err := s.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		blocks := make([]storefront.BlockSnapshot, len(decorations))
		for i, d := range decorations {
			blocks[i] = storefront.BlockSnapshot{
				BlockType:   d.BlockType,
				BlockConfig: d.BlockConfig,
				SortOrder:   d.SortOrder,
			}
		}

		version := &storefront.PageVersion{
			Model:     application.Model{ID: id},
			TenantID:  tenantID,
			PageID:    pageID,
			Version:   page.Version,
			Blocks:    blocks,
			CreatedBy: userID,
		}

		if err := s.versionRepo.Create(ctx, tx, version); err != nil {
			return err
		}

		// Clean up old versions (keep last 20)
		if err := s.versionRepo.DeleteOldest(ctx, tx, tenantID, pageID, 20); err != nil {
			return err
		}

		// Update page status
		now := time.Now().UTC()
		page.IsPublished = true
		page.PublishedAt = &now
		if err := s.pageRepo.Update(ctx, tx, page); err != nil {
			return err
		}

		return nil
	})
}

func (s *pageService) UnpublishPage(ctx context.Context, tenantID shared.TenantID, pageID int64) error {
	page, err := s.pageRepo.FindByID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return err
	}
	if page == nil {
		return code.ErrPageNotFound
	}

	if !page.IsPublished {
		return code.ErrPageNotPublished
	}

	page.IsPublished = false
	page.PublishedAt = nil
	return s.pageRepo.Update(ctx, s.db, page)
}

func pageTypeToString(t storefront.PageType) string {
	switch t {
	case storefront.PageTypeHome:
		return "home"
	case storefront.PageTypeProduct:
		return "product"
	case storefront.PageTypeCollection:
		return "collection"
	case storefront.PageTypeCustom:
		return "custom"
	default:
		return "custom"
	}
}

func decorationsToDTOs(decorations []*storefront.Decoration) []*DecorationDTO {
	dtos := make([]*DecorationDTO, len(decorations))
	for i, d := range decorations {
		dtos[i] = &DecorationDTO{
			ID:          d.Model.ID,
			BlockType:   d.BlockType,
			BlockConfig: d.BlockConfig,
			SortOrder:   d.SortOrder,
		}
	}
	return dtos
}

type decorationService struct {
	db             *gorm.DB
	decorationRepo storefront.DecorationRepository
	idGen          snowflake.Snowflake
}

func NewDecorationService(
	db *gorm.DB,
	decorationRepo storefront.DecorationRepository,
	idGen snowflake.Snowflake,
) DecorationService {
	return &decorationService{
		db:             db,
		decorationRepo: decorationRepo,
		idGen:          idGen,
	}
}

func (s *decorationService) GetDecorations(ctx context.Context, tenantID shared.TenantID, pageID int64) ([]*DecorationDTO, error) {
	decorations, err := s.decorationRepo.FindByPageID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return nil, err
	}
	return decorationsToDTOs(decorations), nil
}

func (s *decorationService) AddDecoration(ctx context.Context, tenantID shared.TenantID, pageID int64, blockType string, blockConfig map[string]any, sortOrder int) (*DecorationDTO, error) {
	// Validate block type
	if !isValidBlockType(blockType) {
		return nil, code.ErrInvalidBlockType
	}

	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	decoration := &storefront.Decoration{
		Model:       application.Model{ID: id},
		TenantID:    tenantID,
		PageID:      pageID,
		BlockType:   blockType,
		BlockConfig: blockConfig,
		SortOrder:   sortOrder,
		IsActive:    true,
	}

	if err := s.decorationRepo.Create(ctx, s.db, decoration); err != nil {
		return nil, err
	}

	return &DecorationDTO{
		ID:          decoration.Model.ID,
		BlockType:   decoration.BlockType,
		BlockConfig: decoration.BlockConfig,
		SortOrder:   decoration.SortOrder,
	}, nil
}

func (s *decorationService) UpdateDecoration(ctx context.Context, tenantID shared.TenantID, decorationID int64, blockConfig map[string]any) error {
	decoration, err := s.decorationRepo.FindByID(ctx, s.db, tenantID, decorationID)
	if err != nil {
		return err
	}
	if decoration == nil {
		return code.ErrDecorationNotFound
	}

	decoration.BlockConfig = blockConfig
	return s.decorationRepo.Update(ctx, s.db, decoration)
}

func (s *decorationService) DeleteDecoration(ctx context.Context, tenantID shared.TenantID, decorationID int64) error {
	return s.decorationRepo.Delete(ctx, s.db, tenantID, decorationID)
}

func (s *decorationService) ReorderBlocks(ctx context.Context, tenantID shared.TenantID, pageID int64, orders []BlockOrderDTO) error {
	blockOrders := make([]storefront.BlockOrder, len(orders))
	for i, o := range orders {
		blockOrders[i] = storefront.BlockOrder{
			Model:    application.Model{ID: o.ID},
			SortOrder: o.SortOrder,
		}
	}
	return s.decorationRepo.Reorder(ctx, s.db, tenantID, blockOrders)
}

// isValidBlockType validates that the block type is one of the allowed types
func isValidBlockType(blockType string) bool {
	validBlockTypes := map[string]bool{
		"banner":            true,
		"product_grid":      true,
		"rich_text":         true,
		"image_carousel":    true,
		"featured_products": true,
		"categories":        true,
		"divider":           true,
		"video":             true,
		"spacer":            true,
		"custom_html":       true,
	}
	return validBlockTypes[blockType]
}

type versionService struct {
	db            *gorm.DB
	pageRepo      storefront.PageRepository
	versionRepo   storefront.PageVersionRepository
	decorationRepo storefront.DecorationRepository
	idGen         snowflake.Snowflake
}

func NewVersionService(
	db *gorm.DB,
	pageRepo storefront.PageRepository,
	versionRepo storefront.PageVersionRepository,
	decorationRepo storefront.DecorationRepository,
	idGen snowflake.Snowflake,
) VersionService {
	return &versionService{
		db:            db,
		pageRepo:      pageRepo,
		versionRepo:   versionRepo,
		decorationRepo: decorationRepo,
		idGen:         idGen,
	}
}

func (s *versionService) ListVersions(ctx context.Context, tenantID shared.TenantID, pageID int64, page, pageSize int) (*PaginatedResult[*VersionDTO], error) {
	versions, total, err := s.versionRepo.FindByPageID(ctx, s.db, tenantID, pageID, page, pageSize)
	if err != nil {
		return nil, err
	}

	dtos := make([]*VersionDTO, len(versions))
	for i, v := range versions {
		dtos[i] = &VersionDTO{
			ID:        v.Model.ID,
			Version:   v.Version,
			CreatedBy: v.CreatedBy,
			CreatedAt: v.Model.CreatedAt.Format(time.RFC3339),
		}
	}
	return &PaginatedResult[*VersionDTO]{
		Items:    dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *versionService) GetVersion(ctx context.Context, tenantID shared.TenantID, pageID int64, version int) (*VersionDetailDTO, error) {
	v, err := s.versionRepo.FindByVersion(ctx, s.db, tenantID, pageID, version)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, code.ErrVersionNotFound
	}

	blocks := make([]*DecorationDTO, len(v.Blocks))
	for i, b := range v.Blocks {
		blocks[i] = &DecorationDTO{
			BlockType:   b.BlockType,
			BlockConfig: b.BlockConfig,
			SortOrder:   b.SortOrder,
		}
	}

	return &VersionDetailDTO{
		Version: &VersionDTO{
			ID:        v.Model.ID,
			Version:   v.Version,
			CreatedBy: v.CreatedBy,
			CreatedAt: v.Model.CreatedAt.Format(time.RFC3339),
		},
		Blocks: blocks,
	}, nil
}

func (s *versionService) RestoreVersion(ctx context.Context, tenantID shared.TenantID, pageID int64, version int, userID int64) error {
	v, err := s.versionRepo.FindByVersion(ctx, s.db, tenantID, pageID, version)
	if err != nil {
		return err
	}
	if v == nil {
		return code.ErrVersionNotFound
	}

	// Get current page to increment version
	page, err := s.pageRepo.FindByID(ctx, s.db, tenantID, pageID)
	if err != nil {
		return err
	}
	if page == nil {
		return code.ErrPageNotFound
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing decorations
		if err := s.decorationRepo.DeleteByPageID(ctx, tx, tenantID, pageID); err != nil {
			return err
		}

		// Recreate decorations from version snapshot
		for _, block := range v.Blocks {
			id, err := s.idGen.NextID(ctx)
			if err != nil {
				return err
			}

			decoration := &storefront.Decoration{
				Model:       application.Model{ID: id},
				TenantID:    tenantID,
				PageID:      pageID,
				BlockType:   block.BlockType,
				BlockConfig: block.BlockConfig,
				SortOrder:   block.SortOrder,
				IsActive:    true,
			}

			if err := s.decorationRepo.Create(ctx, tx, decoration); err != nil {
				return err
			}
		}

		// Create new version for this restore action
		newVersionID, err := s.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		// Increment page version number
		newVersionNum := page.Version + 1

		newVersion := &storefront.PageVersion{
			Model:     application.Model{ID: newVersionID},
			TenantID:  tenantID,
			PageID:    pageID,
			Version:   newVersionNum,
			Blocks:    v.Blocks,
			CreatedBy: userID,
		}

		if err := s.versionRepo.Create(ctx, tx, newVersion); err != nil {
			return err
		}

		// Update page version number
		page.Version = newVersionNum
		return s.pageRepo.Update(ctx, tx, page)
	})
}

type seoService struct {
	db        *gorm.DB
	seoRepo   storefront.SEOConfigRepository
}

func NewSEOService(
	db *gorm.DB,
	seoRepo storefront.SEOConfigRepository,
) SEOService {
	return &seoService{
		db:      db,
		seoRepo: seoRepo,
	}
}

func (s *seoService) GetGlobalSEO(ctx context.Context, tenantID shared.TenantID) (*SEOConfigDTO, error) {
	config, err := s.seoRepo.FindByPageType(ctx, s.db, tenantID, "global", nil)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &SEOConfigDTO{}, nil
	}

	return &SEOConfigDTO{
		Title:       config.Title,
		Description: config.Description,
		Keywords:    config.Keywords,
	}, nil
}

func (s *seoService) UpdateGlobalSEO(ctx context.Context, tenantID shared.TenantID, config SEOConfigDTO) error {
	// Validate SEO fields
	if len(config.Title) > 70 {
		return code.ErrSEOTitleTooLong
	}
	if len(config.Description) > 160 {
		return code.ErrSEODescriptionTooLong
	}

	seoConfig := &storefront.SEOConfigEntity{
		TenantID:    tenantID,
		PageType:    "global",
		PageID:      nil,
		Title:       config.Title,
		Description: config.Description,
		Keywords:    config.Keywords,
	}
	return s.seoRepo.Save(ctx, s.db, seoConfig)
}

func (s *seoService) GetPageSEO(ctx context.Context, tenantID shared.TenantID, pageType string, pageID *int64) (*PageSEOConfigDTO, error) {
	config, err := s.seoRepo.FindByPageType(ctx, s.db, tenantID, pageType, pageID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &PageSEOConfigDTO{
			PageType: pageType,
			PageID:   pageID,
			Config:   SEOConfigDTO{},
		}, nil
	}

	return &PageSEOConfigDTO{
		PageType: config.PageType,
		PageID:   config.PageID,
		Config: SEOConfigDTO{
			Title:       config.Title,
			Description: config.Description,
			Keywords:    config.Keywords,
		},
	}, nil
}

func (s *seoService) UpdatePageSEO(ctx context.Context, tenantID shared.TenantID, pageType string, pageID *int64, config SEOConfigDTO) error {
	// Validate SEO fields
	if len(config.Title) > 70 {
		return code.ErrSEOTitleTooLong
	}
	if len(config.Description) > 160 {
		return code.ErrSEODescriptionTooLong
	}

	seoConfig := &storefront.SEOConfigEntity{
		TenantID:    tenantID,
		PageType:    pageType,
		PageID:      pageID,
		Title:       config.Title,
		Description: config.Description,
		Keywords:    config.Keywords,
	}
	return s.seoRepo.Save(ctx, s.db, seoConfig)
}

func (s *seoService) ListPageSEO(ctx context.Context, tenantID shared.TenantID, page, pageSize int) (*PaginatedResult[*PageSEOConfigDTO], error) {
	configs, total, err := s.seoRepo.FindAll(ctx, s.db, tenantID, page, pageSize)
	if err != nil {
		return nil, err
	}

	dtos := make([]*PageSEOConfigDTO, len(configs))
	for i, c := range configs {
		dtos[i] = &PageSEOConfigDTO{
			PageType: c.PageType,
			PageID:   c.PageID,
			Config: SEOConfigDTO{
				Title:       c.Title,
				Description: c.Description,
				Keywords:    c.Keywords,
			},
		}
	}
	return &PaginatedResult[*PageSEOConfigDTO]{
		Items:    dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}