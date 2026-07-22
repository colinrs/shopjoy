package shipping

import (
	"golang.org/x/text/language"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
)

// ResolveZoneName 按 Accept-Language 解析 zone 在当前 locale 下应展示的名称。
//
// 解析顺序（三级回退 + 兜底）：
//  1. 精确匹配 zone.NameI18n[acceptLanguage]，命中且非空则返回；
//  2. 仅匹配语言基础部分（例如 acceptLanguage="en-GB" 时匹配
//     NameI18n["en-US"]，因为两者 base lang 都是 "en"），命中且非空则返回；
//  3. 取 NameI18n 中第一个非空条目；
//  4. 全部为空或 zone 为 nil → 返回 zone.Name 作为 fallback。
//
// 该函数不依赖任何外部上下文；调用方负责从 HTTP Accept-Language 头取字符串传入。
// acceptLanguage 既支持 BCP-47 标签（如 "en-US"、"zh-Hant-CN"），也支持空字符串
// （空字符串直接走 fallback，不做语言解析）。
func ResolveZoneName(zone *shipping.ShippingZone, acceptLanguage string) string {
	if zone == nil {
		return ""
	}

	// NameI18n 为空（nil 或 len==0）→ fallback
	if len(zone.NameI18n) == 0 {
		return zone.Name
	}

	// acceptLanguage 为空 → 跳过 1-3，直接走 fallback（zone.Name）
	// 见 brief "边界：acceptLanguage 为空 → 直接走 fallback"
	if acceptLanguage == "" {
		return zone.Name
	}

	// 1. 精确匹配
	if v, ok := zone.NameI18n[acceptLanguage]; ok && v != "" {
		return v
	}

	// 2. 语言基础部分匹配（en-US ↔ en-GB → 同 base "en"）
	// Note: x/text 的 Tag.Base() 返回 Base/Confidence 二元组；language.Base{}
	// 是 zero-value 即 undetermined。language.Und 是 Tag，不能与 Base 直接比较。
	if tag, err := language.Parse(acceptLanguage); err == nil {
		base, _ := tag.Base()
		if (language.Base{}) != base {
			for locale, name := range zone.NameI18n {
				if name == "" {
					continue
				}
				if lt, err := language.Parse(locale); err == nil {
					if b, _ := lt.Base(); b == base {
						return name
					}
				}
			}
		}
	}

	// 3. 首个非空
	for _, v := range zone.NameI18n {
		if v != "" {
			return v
		}
	}

	// 4. fallback
	return zone.Name
}
