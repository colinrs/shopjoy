package products

import (
	"context"
	"slices"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

func buildCategoryPath(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) string {
	if categoryID <= 0 {
		return ""
	}

	repo := persistence.NewCategoryRepository()
	categories, err := repo.FindAll(ctx, db, tenantID)
	if err != nil {
		return ""
	}

	categoryMap := make(map[int64]*product.Category, len(categories))
	for _, c := range categories {
		categoryMap[c.Model.ID] = c
	}

	var path []string
	for id := categoryID; id > 0; id = categoryMap[id].ParentID {
		category, ok := categoryMap[id]
		if !ok {
			break
		}
		path = append(path, category.Name)
		if category.ParentID == 0 {
			break
		}
	}

	slices.Reverse(path)
	return strings.Join(path, " -> ")
}
