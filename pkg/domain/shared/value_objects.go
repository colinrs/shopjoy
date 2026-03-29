// Package shared 共享内核
// 包含所有限界上下文共享的值对象、错误和事件
package shared

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// ==================== 值对象 ====================

// TenantID 租户ID（值对象）
type TenantID int64

func (t TenantID) Int64() int64 {
	return int64(t)
}

func (t TenantID) String() string {
	return fmt.Sprintf("%d", t)
}

func (t TenantID) IsValid() bool {
	return t > 0
}

// Money 金额（值对象）
type Money struct {
	Amount   decimal.Decimal // 单位为元，避免浮点数精度问题
	Currency string         // 币种，如 CNY, USD
}

// NewMoney 创建金额
func NewMoney(amount decimal.Decimal, currency string) Money {
	if currency == "" {
		currency = "CNY"
	}
	return Money{Amount: amount, Currency: currency}
}

// NewMoneyFromString 从字符串创建金额（元）
// 例如 "1.99" -> Money{Amount: decimal.NewFromString("1.99"), Currency: "CNY"}
func NewMoneyFromString(amountStr string, currency string) (Money, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount: %w", err)
	}
	return NewMoney(amount, currency), nil
}

// Add 金额相加
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return NewMoney(m.Amount.Add(other.Amount), m.Currency), nil
}

// Subtract 金额相减
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	if m.Amount.LessThan(other.Amount) {
		return Money{}, fmt.Errorf("insufficient amount: %s < %s", m.Amount.String(), other.Amount.String())
	}
	return NewMoney(m.Amount.Sub(other.Amount), m.Currency), nil
}

// Multiply 金额乘以整数
func (m Money) Multiply(multiplier int) Money {
	return NewMoney(m.Amount.Mul(decimal.NewFromInt(int64(multiplier))), m.Currency)
}

// MultiplyFloat 金额乘以浮点数（注意：可能存在精度问题，仅在必要时使用）
func (m Money) MultiplyFloat(multiplier float64) Money {
	return NewMoney(m.Amount.Mul(decimal.NewFromFloat(multiplier)), m.Currency)
}

// Equals 金额相等
func (m Money) Equals(other Money) bool {
	return m.Amount.Equal(other.Amount) && m.Currency == other.Currency
}

// IsZero 金额是否为零
func (m Money) IsZero() bool {
	return m.Amount.IsZero()
}

// IsPositive 金额是否为正
func (m Money) IsPositive() bool {
	return m.Amount.IsPositive()
}

// IsNegative 金额是否为负
func (m Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

// String 返回格式化字符串
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(2), m.Currency)
}

// FormatMoneyOnly 将金额格式化为字符串（元），不带货币单位
func FormatMoneyOnly(m Money) string {
	return m.Amount.StringFixed(2)
}

// ParseMoneyFromString 从字符串解析为 decimal.Decimal（元）
// 例如 "1.99" -> decimal.NewFromString("1.99")
func ParseMoneyFromString(amountStr string) (decimal.Decimal, error) {
	return decimal.NewFromString(amountStr)
}

// FormatMoneyToStringOnly 将 decimal.Decimal 格式化为字符串（元），不带货币单位
// 例如 decimal.NewFromString("1.99") -> "1.99"
func FormatMoneyToStringOnly(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}

// FormatMoneyToString 将 decimal.Decimal 格式化为字符串（元），带货币单位
// 例如 decimal.NewFromString("1.99"), "CNY" -> "1.99 CNY"
func FormatMoneyToString(amount decimal.Decimal, currency string) string {
	return fmt.Sprintf("%s %s", amount.StringFixed(2), currency)
}

// ==================== 状态枚举 ====================

// Status 通用状态
type Status int

const (
	StatusDisabled Status = iota // 禁用
	StatusEnabled                // 启用
)

func (s Status) String() string {
	switch s {
	case StatusDisabled:
		return "disabled"
	case StatusEnabled:
		return "enabled"
	default:
		return "unknown"
	}
}

func (s Status) IsValid() bool {
	return s >= StatusDisabled && s <= StatusEnabled
}

// AuditInfo 审计信息（值对象）
type AuditInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int64
	UpdatedBy int64
}

// NewAuditInfo 创建审计信息
func NewAuditInfo(createdBy int64) AuditInfo {
	now := time.Now().UTC()
	return AuditInfo{
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
}

// Update 更新审计信息
func (a *AuditInfo) Update(updatedBy int64) {
	a.UpdatedAt = time.Now().UTC()
	a.UpdatedBy = updatedBy
}

// ==================== 分页相关 ====================

// PageQuery 分页查询
type PageQuery struct {
	Page     int
	PageSize int
}

// Validate 验证分页参数
func (p *PageQuery) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset 计算偏移量
func (p PageQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 返回限制数
func (p PageQuery) Limit() int {
	return p.PageSize
}

// PageResult 分页结果
type PageResult struct {
	Total    int64
	Page     int
	PageSize int
}

// ==================== 时间范围 ====================

// TimeRange 时间范围
type TimeRange struct {
	StartAt time.Time
	EndAt   time.Time
}

// IsValid 验证时间范围
func (t TimeRange) IsValid() bool {
	return !t.StartAt.IsZero() && !t.EndAt.IsZero() && t.StartAt.Before(t.EndAt)
}

// Contains 检查时间点是否在范围内
func (t TimeRange) Contains(tm time.Time) bool {
	return !tm.Before(t.StartAt) && !tm.After(t.EndAt)
}

// IsActive 检查当前时间是否在范围内
func (t TimeRange) IsActive() bool {
	now := time.Now().UTC()
	return t.Contains(now)
}

// ==================== Unix时间戳类型 ====================

// UnixTime 存储为Unix时间戳的时间类型
// 用于数据库中存储为BIGINT的时间字段
type UnixTime struct {
	time.Time
}

// NewUnixTime 从time.Time创建UnixTime
func NewUnixTime(t time.Time) UnixTime {
	return UnixTime{Time: t}
}

// NewUnixTimeFromInt64 从Unix时间戳创建UnixTime
func NewUnixTimeFromInt64(ts int64) UnixTime {
	return UnixTime{Time: time.Unix(ts, 0).UTC()}
}

// Scan 实现sql.Scanner接口，从数据库读取
func (ut *UnixTime) Scan(value interface{}) error {
	if value == nil {
		ut.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case int64:
		ut.Time = time.Unix(v, 0).UTC()
		return nil
	case float64:
		ut.Time = time.Unix(int64(v), 0).UTC()
		return nil
	default:
		return fmt.Errorf("cannot scan %T into UnixTime", value)
	}
}

// Value 实现driver.Valuer接口，写入数据库
func (ut UnixTime) Value() (driver.Value, error) {
	if ut.Time.IsZero() {
		return nil, nil
	}
	return ut.Time.Unix(), nil
}

// UnixTimePtr 将*UnixTime转为*time.Time的辅助函数
func UnixTimePtr(t *time.Time) *UnixTime {
	if t == nil {
		return nil
	}
	return &UnixTime{Time: *t}
}

// TimePtr 将*UnixTime转为*time.Time的辅助函数
func (ut *UnixTime) TimePtr() *time.Time {
	if ut == nil || ut.Time.IsZero() {
		return nil
	}
	return &ut.Time
}
