package shipping

import (
	"slices"
	"sort"
)

// MatchInput 区域匹配输入
type MatchInput struct {
	CountryCode  string // ISO 3166-1 alpha-2
	ProvinceCode string // ISO 3166-2 或旧城市码
	CityCode     string // 兼容旧中国城市码
	PostalCode   string
}

// ZoneMatcher 多级区域匹配器（province > city > country fallback）
type ZoneMatcher struct {
	zones []*ShippingZone
}

func NewZoneMatcher(zones []*ShippingZone) *ZoneMatcher {
	return &ZoneMatcher{zones: zones}
}

// Match 按优先级匹配：精确 > 城市 > 国家
func (m *ZoneMatcher) Match(in MatchInput) *ShippingZone {
	// 按 sort 升序遍历，第一个匹配胜出
	sorted := make([]*ShippingZone, len(m.zones))
	copy(sorted, m.zones)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Sort < sorted[j].Sort
	})

	var countryFallback *ShippingZone
	for _, z := range sorted {
		if !z.Regions.Contains(in.ProvinceCode) && !z.Regions.Contains(in.CityCode) {
			// 不是精确匹配
			if z.Regions.Contains(in.CountryCode) {
				if countryFallback == nil || z.Sort < countryFallback.Sort {
					countryFallback = z
				}
			}
			continue
		}
		return z
	}
	return countryFallback
}

// Contains 检查 regions 是否包含指定 code
func (r Regions) Contains(code string) bool {
	return slices.Contains(r, code)
}
