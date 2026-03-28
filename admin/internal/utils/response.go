// Package utils provides shared utilities for the admin service
package utils

import (
	"time"

	"github.com/colinrs/shopjoy/admin/internal/types"
)

// FormatTime formats a time.Time to RFC3339 string
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// FormatDateTime formats a time.Time to a human-readable datetime string
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// FormatDate formats a time.Time to a date string
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateOnly)
}

// Pagination represents common pagination parameters
type Pagination struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// Validate ensures pagination values are within acceptable bounds
func (p *Pagination) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}
}

// Offset returns the calculated offset for database queries
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the page size for database queries
func (p *Pagination) Limit() int {
	return p.PageSize
}

// ToPaginationResp creates a common pagination response
func ToPaginationResp(page, pageSize int, total int64) PaginationResponse {
	return PaginationResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

// PaginationResponse is a generic pagination response structure
type PaginationResponse struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// Status constants
const (
	StatusDisabled = 0
	StatusEnabled  = 1
)

// BoolToStatus converts a boolean to status int8
func BoolToStatus(enabled bool) int8 {
	if enabled {
		return StatusEnabled
	}
	return StatusDisabled
}

// StatusToBool converts status int8 to boolean
func StatusToBool(status int8) bool {
	return status == StatusEnabled
}

// BuildCategoryTreeResp builds a category tree response from category data
func BuildCategoryTreeResp(
	id, parentID int64,
	name, code string,
	level, sort int,
	icon, image, seoTitle, seoDescription string,
	status int8,
	productCount int64,
	createdAt, updatedAt time.Time,
	children []*types.CategoryTreeResp,
) *types.CategoryTreeResp {
	return &types.CategoryTreeResp{
		ID:             id,
		ParentID:       parentID,
		Name:           name,
		Code:           code,
		Level:          level,
		Sort:           sort,
		Icon:           icon,
		Image:          image,
		SeoTitle:       seoTitle,
		SeoDescription: seoDescription,
		Status:         status,
		ProductCount:   productCount,
		CreatedAt:      FormatDateTime(createdAt),
		UpdatedAt:      FormatDateTime(updatedAt),
		Children:       children,
	}
}

// BuildBrandDetailResp builds a brand detail response
func BuildBrandDetailResp(
	id int64,
	name, logo, description, website, trademarkNumber, trademarkCountry string,
	enablePage bool,
	sort int,
	status int8,
	productCount int64,
	createdAt, updatedAt time.Time,
) *types.BrandDetailResp {
	return &types.BrandDetailResp{
		ID:               id,
		Name:             name,
		Logo:             logo,
		Description:      description,
		Website:          website,
		TrademarkNumber:  trademarkNumber,
		TrademarkCountry: trademarkCountry,
		EnablePage:       enablePage,
		Sort:             sort,
		Status:           status,
		ProductCount:     productCount,
		CreatedAt:        FormatDateTime(createdAt),
		UpdatedAt:        FormatDateTime(updatedAt),
	}
}

// BuildWarehouseDetailResp builds a warehouse detail response
func BuildWarehouseDetailResp(
	id int64,
	code, name, country, address string,
	isDefault bool,
	status int8,
	createdAt, updatedAt time.Time,
) *types.WarehouseDetailResp {
	return &types.WarehouseDetailResp{
		ID:        id,
		Code:      code,
		Name:      name,
		Country:   country,
		Address:   address,
		IsDefault: isDefault,
		Status:    status,
		CreatedAt: FormatDateTime(createdAt),
		UpdatedAt: FormatDateTime(updatedAt),
	}
}

// BuildMarketResponse builds a market response
func BuildMarketResponse(
	id int64,
	code, name, currency, defaultLanguage, flag string,
	isActive, isDefault bool,
	taxRules types.TaxConfig,
	createdAt, updatedAt time.Time,
) *types.MarketResponse {
	return &types.MarketResponse{
		ID:              id,
		Code:            code,
		Name:            name,
		Currency:        currency,
		DefaultLanguage: defaultLanguage,
		Flag:            flag,
		IsActive:        isActive,
		IsDefault:       isDefault,
		TaxRules:        taxRules,
		CreatedAt:       FormatDateTime(createdAt),
		UpdatedAt:       FormatDateTime(updatedAt),
	}
}
