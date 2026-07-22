// Package region 区域目录领域层
//
// 用于多层级区域匹配（国家/省/市/区）。
// 数据来源包括：
//   - 平台预置数据（tenant_id=0）：中国省份/城市/区，海外国家/州省
//   - 租户自定义数据（tenant_id>0）：覆盖或新增
package region

import (
	"github.com/colinrs/shopjoy/pkg/application"
)

// Region 区域
type Region struct {
	application.Model
	TenantID      int64  `gorm:"column:tenant_id;not null;default:0;index"`
	Code          string `gorm:"column:code;size:20;not null"`
	Name          string `gorm:"column:name;size:100;not null"`
	Level         int    `gorm:"column:level;not null;default:1"`
	ParentCode    string `gorm:"column:parent_code;size:20;not null;default:''"`
	CountryCode   string `gorm:"column:country_code;size:2;not null;default:'CN';index"`
	PostalPattern string `gorm:"column:postal_pattern;size:255;not null;default:''"`
	Sort          int    `gorm:"column:sort;not null;default:0"`
	IsActive      bool   `gorm:"column:is_active;not null;default:true"`
}

func (r *Region) TableName() string {
	return "regions"
}

// RegionLevel 区域层级
const (
	RegionLevelCountry  = 1
	RegionLevelProvince = 2
	RegionLevelCity     = 3
	RegionLevelDistrict = 4
)
