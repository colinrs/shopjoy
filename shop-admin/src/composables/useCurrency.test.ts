import { describe, expect, it } from 'vitest'
import { createApp, defineComponent, h } from 'vue'
import { createI18n } from 'vue-i18n'
import { ref } from 'vue'
import { useCurrency } from './useCurrency'

// useCurrency() calls useI18n(), which relies on the current Vue instance
// (it injects the i18n symbol registered when createI18n() is installed on
// the app). Calling useI18n from a plain function throws "must be called at
// the top of a `setup` function". The cleanest way to test such a composable
// is to call it inside a component's setup() — Vue's reactivity helpers
// only resolve inject()s during a setup call.
// locale defaults to 'en' so toLocaleString output uses comma thousands
// separators — matching the assertions below.
function withCurrency<T>(currencyCode: string, cb: (api: ReturnType<typeof useCurrency>) => T): T {
  let result!: T
  let caughtError: unknown = null

  const i18n = createI18n({
    legacy: false,
    locale: 'en',
    fallbackLocale: 'en',
    messages: { en: {}, zh: {} }
  })

  const Comp = defineComponent({
    setup() {
      const api = useCurrency(ref(currencyCode))
      try {
        result = cb(api)
      } catch (e) {
        caughtError = e
      }
      return () => h('div')
    }
  })

  const app = createApp(Comp)
  app.use(i18n)
  const root = document.createElement('div')
  document.body.appendChild(root)
  try {
    app.mount(root)
  } finally {
    app.unmount()
    document.body.removeChild(root)
  }

  if (caughtError) throw caughtError
  return result
}

describe('useCurrency', () => {
  it('formats USD with 2 decimals', () => {
    const { format } = withCurrency('USD', (api) => api)
    expect(format('1.99')).toBe('$1.99')
    expect(format('1234.5')).toBe('$1,234.50')
  })

  it('formats JPY with 0 decimals', () => {
    const { format } = withCurrency('JPY', (api) => api)
    expect(format('1234.56')).toBe('¥1,235')
  })

  it('formats VND with symbol after', () => {
    const { format } = withCurrency('VND', (api) => api)
    expect(format('100000')).toBe('100,000₫')
  })
})