package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePageLogic {
	return CreatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePageLogic) CreatePage(req *types.CreatePageRequest) (resp *types.CreatePageResponse, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("page name is required")
	}

	// Generate slug if not provided
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Validate page type
	pageType := req.PageType
	if pageType == "" {
		pageType = "custom"
	}
	if !isValidPageType(pageType) {
		pageType = "custom"
	}

	// Create the page
	pageDTO, err := l.svcCtx.PageService.CreatePage(l.ctx, shared.TenantID(tenantID), req.Name, slug, pageType)
	if err != nil {
		return nil, err
	}

	return &types.CreatePageResponse{
		PageID: pageDTO.ID,
	}, nil
}

// generateSlug creates a URL-friendly slug from the page name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove characters that are not alphanumeric or hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// isValidPageType checks if the page type is valid
func isValidPageType(pageType string) bool {
	validTypes := map[string]bool{
		"home":       true,
		"product":    true,
		"collection": true,
		"custom":    true,
	}
	return validTypes[pageType]
}
