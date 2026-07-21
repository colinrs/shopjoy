package shipping

import "math"

// WeightUnit 重量单位
type WeightUnit string

const (
	// WeightUnitGram 克
	WeightUnitGram WeightUnit = "g"
	// WeightUnitKilogram 千克
	WeightUnitKilogram WeightUnit = "kg"
	// WeightUnitPound 磅
	WeightUnitPound WeightUnit = "lb"
	// WeightUnitOunce 盎司
	WeightUnitOunce WeightUnit = "oz"
)

// 转换系数（1 单位 = 多少克）
const (
	gramsPerKilogram = 1000.0
	gramsPerPound    = 453.59237
	gramsPerOunce    = 28.349523125
)

// gramsPerUnit 返回 1 个指定单位对应的克数；未知单位按克处理。
func gramsPerUnit(unit WeightUnit) float64 {
	switch unit {
	case WeightUnitKilogram:
		return gramsPerKilogram
	case WeightUnitPound:
		return gramsPerPound
	case WeightUnitOunce:
		return gramsPerOunce
	case WeightUnitGram:
		return 1.0
	default:
		return 1.0
	}
}

// WeightConverter 重量单位转换器，绑定一个源/目标单位。
type WeightConverter struct {
	unit WeightUnit
}

// NewWeightConverter 创建绑定指定单位的转换器。
func NewWeightConverter(unit WeightUnit) WeightConverter {
	return WeightConverter{unit: unit}
}

// ToGrams 将当前单位的值转换为克（四舍五入到整克）。
func (c WeightConverter) ToGrams(value float64) int {
	grams := value * gramsPerUnit(c.unit)
	return int(math.Round(grams))
}

// FromGrams 将克转换为当前单位的值（浮点）。
func (c WeightConverter) FromGrams(grams int) float64 {
	return float64(grams) / gramsPerUnit(c.unit)
}
