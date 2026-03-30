import { createI18n } from 'vue-i18n'
import type { Ref } from 'vue'
import en from '@/locales/en.json'
import zh from '@/locales/zh.json'

export type MessageSchema = typeof en

const i18n = createI18n<[MessageSchema], 'en' | 'zh'>({
  legacy: false,
  locale: localStorage.getItem('locale') || 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    zh
  }
})

export default i18n

// Typed translation function for use in Vue components
// Usage: import { t } from '@/plugins/i18n' then t('common.save')
// Supports interpolation: t('key', { name: 'John' })
export function t(key: string, params?: Record<string, string | number>): string {
  if (params) {
    return i18n.global.t(key, params) as string
  }
  return i18n.global.t(key) as string
}

// Composable for组件中使用
export function useLocale(): {
  locale: Ref<string>
  toggleLocale: () => void
} {
  // In Composition API mode with legacy:false, locale is a Ref<string>
  const locale = i18n.global.locale as unknown as Ref<string>

  const toggleLocale = () => {
    const current = locale.value
    locale.value = current === 'en' ? 'zh' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return {
    locale,
    toggleLocale
  }
}
