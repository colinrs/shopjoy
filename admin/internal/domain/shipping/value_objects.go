package shipping

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// StringI18n 多语言字符串映射（如 {"en-US": "Standard", "ja-JP": "通常"}）
type StringI18n map[string]string

// Value 实现 driver.Valuer 接口（GORM JSONB 写入）
func (s StringI18n) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口（GORM JSONB 读取）
func (s *StringI18n) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("StringI18n: cannot scan %T", src)
	}
	return json.Unmarshal(data, s)
}

// Get 按 locale 取值，三级回退：精确 → fallback → 第一个非空值。
// locale/fallback 为 BCP-47 标签（如 "zh-CN"），fallback 通常为该区域的首选语言标签。
// 注意：空字符串 locale 不视为有效匹配，避免与空键碰撞。
func (s StringI18n) Get(locale, fallback string) string {
	if s == nil {
		return ""
	}
	if locale != "" {
		if v, ok := s[locale]; ok && v != "" {
			return v
		}
	}
	if fallback != "" && fallback != locale {
		if v, ok := s[fallback]; ok && v != "" {
			return v
		}
	}
	// 最后回退到第一个非空值
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

// StringArray 字符串数组值对象（用于邮编模式、远程区域等列表字段）。
type StringArray []string

// Value 实现 driver.Valuer 接口（GORM JSONB 写入）
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口（GORM JSONB 读取）
func (a *StringArray) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("StringArray: unsupported scan type")
	}
	return json.Unmarshal(data, a)
}

// Contains 检查数组中是否包含给定字符串。
func (a StringArray) Contains(s string) bool {
	return slices.Contains(a, s)
}