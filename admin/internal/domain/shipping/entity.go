package shipping

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

// FeeType 运费计费类型
type FeeType string

const (
	FeeTypeFixed    FeeType = "fixed"     // 固定运费
	FeeTypeByCount  FeeType = "by_count"  // 按件计费
	FeeTypeByWeight FeeType = "by_weight" // 按重量计费
	FeeTypeFree     FeeType = "free"      // 免运费
)

// IsValid 检查计费类型是否有效
func (f FeeType) IsValid() bool {
	switch f {
	case FeeTypeFixed, FeeTypeByCount, FeeTypeByWeight, FeeTypeFree:
		return true
	default:
		return false
	}
}

// String 返回计费类型字符串
func (f FeeType) String() string {
	return string(f)
}

// Regions 区域代码列表（JSONB）
type Regions []string

// Value 实现 driver.Valuer 接口
func (r Regions) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan 实现 sql.Scanner 接口
func (r *Regions) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, r)
}

// ShippingTemplate 运费模板实体
type ShippingTemplate struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	TenantID  int64          `gorm:"column:tenant_id;not null;index"`
	Name      string         `gorm:"column:name;size:100;not null"`
	IsDefault bool           `gorm:"column:is_default;not null;default:false;index"`
	IsActive  bool           `gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 返回表名
func (t *ShippingTemplate) TableName() string {
	return "shipping_templates"
}

// CanDelete 检查是否可以删除
func (t *ShippingTemplate) CanDelete() error {
	if t.IsDefault {
		return code.ErrShippingTemplateIsDefault
	}
	return nil
}

// SetDefault 设置为默认模板
func (t *ShippingTemplate) SetDefault() {
	t.IsDefault = true
}

// UnsetDefault 取消默认模板
func (t *ShippingTemplate) UnsetDefault() {
	t.IsDefault = false
}

// Activate 激活模板
func (t *ShippingTemplate) Activate() {
	t.IsActive = true
}

// Deactivate 停用模板
func (t *ShippingTemplate) Deactivate() {
	t.IsActive = false
}

// ShippingZone 配送区域实体
type ShippingZone struct {
	ID                  int64     `gorm:"column:id;primaryKey"`
	TenantID            int64     `gorm:"column:tenant_id;not null;index"`
	TemplateID          int64     `gorm:"column:template_id;not null;index"`
	Name                string    `gorm:"column:name;size:100;not null"`
	Regions             Regions   `gorm:"column:regions;type:jsonb;not null"`
	FeeType             FeeType   `gorm:"column:fee_type;size:20;not null"`
	FirstUnit           int       `gorm:"column:first_unit;not null;default:1"`
	FirstFee            int64     `gorm:"column:first_fee;not null;default:0"` // 分为单位
	AdditionalUnit      int       `gorm:"column:additional_unit;not null;default:1"`
	AdditionalFee       int64     `gorm:"column:additional_fee;not null;default:0"`
	FreeThresholdAmount int64     `gorm:"column:free_threshold_amount;not null;default:0"` // 满额包邮，0表示不限制
	FreeThresholdCount  int       `gorm:"column:free_threshold_count;not null;default:0"`  // 满件包邮，0表示不限制
	Sort                int       `gorm:"column:sort;not null;default:0"`
	CreatedAt           time.Time `gorm:"column:created_at;not null"`
	UpdatedAt           time.Time `gorm:"column:updated_at;not null"`
}

// TableName 返回表名
func (z *ShippingZone) TableName() string {
	return "shipping_zones"
}

// Validate 验证区域配置
func (z *ShippingZone) Validate() error {
	if z.Name == "" {
		return code.ErrShippingZoneNameRequired
	}
	if len(z.Regions) == 0 {
		return code.ErrShippingZoneRegionsRequired
	}
	if !z.FeeType.IsValid() {
		return code.ErrShippingZoneInvalidFeeType
	}
	// 非免运费类型需要配置运费
	if z.FeeType != FeeTypeFree {
		if z.FirstUnit <= 0 {
			return code.ErrShippingZoneFeeConfigRequired
		}
		if z.FirstFee < 0 {
			return code.ErrShippingZoneFeeConfigRequired
		}
		if z.AdditionalUnit <= 0 {
			return code.ErrShippingZoneFeeConfigRequired
		}
		if z.AdditionalFee < 0 {
			return code.ErrShippingZoneFeeConfigRequired
		}
	}
	return nil
}

// MatchesRegion 检查是否匹配指定区域
func (z *ShippingZone) MatchesRegion(cityCode string) bool {
	for _, region := range z.Regions {
		if region == cityCode {
			return true
		}
	}
	return false
}

// TargetType 关联目标类型
type TargetType string

const (
	TargetTypeProduct  TargetType = "product"
	TargetTypeCategory TargetType = "category"
)

// IsValid 检查目标类型是否有效
func (t TargetType) IsValid() bool {
	switch t {
	case TargetTypeProduct, TargetTypeCategory:
		return true
	default:
		return false
	}
}

// ShippingZoneRegion 配送区域城市关联实体（用于索引查询）
type ShippingZoneRegion struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	ZoneID    int64     `gorm:"column:zone_id;not null;uniqueIndex:idx_zone_city"`
	CityCode  string    `gorm:"column:city_code;size:20;not null;uniqueIndex:idx_zone_city;index"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

// TableName 返回表名
func (r *ShippingZoneRegion) TableName() string {
	return "shipping_zone_regions"
}

// ShippingTemplateMapping 模板关联实体
type ShippingTemplateMapping struct {
	ID         int64      `gorm:"column:id;primaryKey"`
	TenantID   int64      `gorm:"column:tenant_id;not null;index"`
	TemplateID int64      `gorm:"column:template_id;not null;index"`
	TargetType TargetType `gorm:"column:target_type;size:20;not null"`
	TargetID   int64      `gorm:"column:target_id;not null;index"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null"`
}

// TableName 返回表名
func (m *ShippingTemplateMapping) TableName() string {
	return "shipping_template_mappings"
}

// Validate 验证关联配置
func (m *ShippingTemplateMapping) Validate() error {
	if !m.TargetType.IsValid() {
		return code.ErrShippingMappingInvalidTarget
	}
	return nil
}

// CalculateShippingFeeRequest 计算运费请求
type CalculateShippingFeeRequest struct {
	ProvinceCode string
	CityCode     string
	DistrictCode string
	Items        []CalculateItem
}

// CalculateItem 计算项
type CalculateItem struct {
	ProductID int64
	SKUID     int64
	Quantity  int
	Weight    int    // 克
	Price     int64  // 分为单位
}

// CalculateShippingFeeResult 计算运费结果
type CalculateShippingFeeResult struct {
	ShippingFee    int64
	Currency       string
	TemplateID     int64
	TemplateName   string
	ZoneName       string
	FeeType        FeeType
	FirstUnit      int
	FirstFee       int64
	AdditionalUnit int
	AdditionalFee  int64
	TotalWeight    int
	TotalUnits     int
}

// CalculateFee 计算运费
func (z *ShippingZone) CalculateFee(items []CalculateItem, orderAmount int64, itemCount int) int64 {
	// 检查免运费条件
	if z.FreeThresholdAmount > 0 && orderAmount >= z.FreeThresholdAmount {
		return 0
	}
	if z.FreeThresholdCount > 0 && itemCount >= z.FreeThresholdCount {
		return 0
	}

	// 免运费类型
	if z.FeeType == FeeTypeFree {
		return 0
	}

	// 固定运费
	if z.FeeType == FeeTypeFixed {
		return z.FirstFee
	}

	// 计算总单位数
	var totalUnits int
	switch z.FeeType {
	case FeeTypeByCount:
		for _, item := range items {
			totalUnits += item.Quantity
		}
	case FeeTypeByWeight:
		for _, item := range items {
			totalUnits += item.Weight
		}
	}

	// 计算运费
	if totalUnits <= z.FirstUnit {
		return z.FirstFee
	}

	// 计算续费单位数
	additionalUnits := (totalUnits - z.FirstUnit + z.AdditionalUnit - 1) / z.AdditionalUnit
	return z.FirstFee + int64(additionalUnits)*z.AdditionalFee
}