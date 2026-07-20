import { computed, type Ref } from 'vue'
import { useI18n } from 'vue-i18n'

export interface CurrencyConfig {
  code: string
  symbol: string
  decimals: number          // JPY=0, USD=2, KRW=0
  position: 'before' | 'after'
}

const CURRENCY_CONFIG: Record<string, CurrencyConfig> = {
  CNY: { code: 'CNY', symbol: '¥',  decimals: 2, position: 'before' },
  USD: { code: 'USD', symbol: '$',  decimals: 2, position: 'before' },
  EUR: { code: 'EUR', symbol: '€',  decimals: 2, position: 'before' },
  GBP: { code: 'GBP', symbol: '£',  decimals: 2, position: 'before' },
  JPY: { code: 'JPY', symbol: '¥',  decimals: 0, position: 'before' },
  KRW: { code: 'KRW', symbol: '₩',  decimals: 0, position: 'before' },
  SGD: { code: 'SGD', symbol: 'S$', decimals: 2, position: 'before' },
  MYR: { code: 'MYR', symbol: 'RM', decimals: 2, position: 'before' },
  THB: { code: 'THB', symbol: '฿',  decimals: 2, position: 'before' },
  IDR: { code: 'IDR', symbol: 'Rp', decimals: 0, position: 'before' },
  PHP: { code: 'PHP', symbol: '₱',  decimals: 2, position: 'before' },
  VND: { code: 'VND', symbol: '₫',  decimals: 0, position: 'after' },
}

export function useCurrency(currencyRef: Ref<string>) {
  const { locale } = useI18n()

  const config = computed(() =>
    CURRENCY_CONFIG[currencyRef.value] || CURRENCY_CONFIG.CNY
  )

  /**
   * 格式化金额：
   * - input: "1.99" -> "$1.99"
   * - JPY/KRW/IDR/VND 无小数
   * - VND 符号在右侧
   */
  const format = (amount: string | number): string => {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount
    if (isNaN(num)) return ''
    const fixed = num.toFixed(config.value.decimals)
    const formatted = config.value.decimals > 0
      ? Number(fixed).toLocaleString(locale.value, {
          minimumFractionDigits: config.value.decimals,
          maximumFractionDigits: config.value.decimals,
        })
      : Number(fixed).toLocaleString(locale.value)
    return config.value.position === 'before'
      ? `${config.value.symbol}${formatted}`
      : `${formatted}${config.value.symbol}`
  }

  return { config, format }
}