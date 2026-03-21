/**
 * Shared validation utilities for the admin frontend
 */

/**
 * Email validation regex
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

/**
 * Phone validation regex (supports Chinese mobile numbers)
 */
const PHONE_REGEX = /^1[3-9]\d{9}$/

/**
 * Password strength requirements:
 * - At least 8 characters
 * - Contains uppercase and lowercase letters
 * - Contains numbers
 */
const PASSWORD_STRONG_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/

/**
 * URL validation regex
 */
const URL_REGEX = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  return EMAIL_REGEX.test(email)
}

/**
 * Validate phone number (Chinese mobile)
 */
export function isValidPhone(phone: string): boolean {
  return PHONE_REGEX.test(phone)
}

/**
 * Validate password strength
 */
export function isStrongPassword(password: string): boolean {
  return PASSWORD_STRONG_REGEX.test(password)
}

/**
 * Validate URL format
 */
export function isValidUrl(url: string): boolean {
  return URL_REGEX.test(url)
}

/**
 * Validate required field (non-empty string)
 */
export function isRequired(value: string | null | undefined): boolean {
  if (value === null || value === undefined) return false
  return value.trim().length > 0
}

/**
 * Validate minimum length
 */
export function minLength(value: string, min: number): boolean {
  return value.length >= min
}

/**
 * Validate maximum length
 */
export function maxLength(value: string, max: number): boolean {
  return value.length <= max
}

/**
 * Validate numeric range
 */
export function inRange(value: number, min: number, max: number): boolean {
  return value >= min && value <= max
}

/**
 * Validate positive number
 */
export function isPositive(value: number): boolean {
  return value > 0
}

/**
 * Validate non-negative number (including zero)
 */
export function isNonNegative(value: number): boolean {
  return value >= 0
}

/**
 * Validate integer
 */
export function isInteger(value: number): boolean {
  return Number.isInteger(value)
}

/**
 * Validate SKU code format (alphanumeric with dashes and underscores)
 */
export function isValidSkuCode(sku: string): boolean {
  return /^[A-Za-z0-9_-]+$/.test(sku)
}

/**
 * Validate category code format
 */
export function isValidCategoryCode(code: string): boolean {
  return /^[A-Za-z0-9_-]+$/.test(code)
}

/**
 * Validate market code (ISO country code format)
 */
export function isValidMarketCode(code: string): boolean {
  return /^[A-Z]{2,3}$/.test(code)
}

/**
 * Validate currency code (ISO 4217)
 */
export function isValidCurrencyCode(currency: string): boolean {
  return /^[A-Z]{3}$/.test(currency)
}

/**
 * Form validation rules for Element Plus
 */
export const formRules = {
  required: (message = '此字段不能为空') => ({
    required: true,
    message,
    trigger: 'blur' as const
  }),

  email: (message = '请输入有效的邮箱地址') => ({
    validator: (_rule: any, value: string, callback: (error?: Error) => void) => {
      if (!value || isValidEmail(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    },
    trigger: 'blur' as const
  }),

  phone: (message = '请输入有效的手机号码') => ({
    validator: (_rule: any, value: string, callback: (error?: Error) => void) => {
      if (!value || isValidPhone(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    },
    trigger: 'blur' as const
  }),

  minLength: (min: number, message?: string) => ({
    min,
    message: message || `长度不能少于 ${min} 个字符`,
    trigger: 'blur' as const
  }),

  maxLength: (max: number, message?: string) => ({
    max,
    message: message || `长度不能超过 ${max} 个字符`,
    trigger: 'blur' as const
  }),

  password: (message = '密码至少8位，包含大小写字母和数字') => ({
    validator: (_rule: any, value: string, callback: (error?: Error) => void) => {
      if (!value || isStrongPassword(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    },
    trigger: 'blur' as const
  }),

  url: (message = '请输入有效的URL') => ({
    validator: (_rule: any, value: string, callback: (error?: Error) => void) => {
      if (!value || isValidUrl(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    },
    trigger: 'blur' as const
  }),

  positiveNumber: (message = '请输入大于0的数字') => ({
    validator: (_rule: any, value: number, callback: (error?: Error) => void) => {
      if (value === undefined || value === null || isPositive(value)) {
        callback()
      } else {
        callback(new Error(message))
      }
    },
    trigger: 'blur' as const
  })
}

/**
 * Common validation result type
 */
export interface ValidationResult {
  valid: boolean
  message?: string
}

/**
 * Validate form data against rules
 */
export function validateForm(data: Record<string, any>, rules: Record<string, any[]>): ValidationResult {
  for (const [field, fieldRules] of Object.entries(rules)) {
    const value = data[field]
    for (const rule of fieldRules) {
      if (rule.required && !isRequired(value)) {
        return { valid: false, message: rule.message || `${field} 不能为空` }
      }
      if (rule.min && typeof value === 'string' && !minLength(value, rule.min)) {
        return { valid: false, message: rule.message || `${field} 长度不足` }
      }
      if (rule.max && typeof value === 'string' && !maxLength(value, rule.max)) {
        return { valid: false, message: rule.message || `${field} 长度超出限制` }
      }
    }
  }
  return { valid: true }
}