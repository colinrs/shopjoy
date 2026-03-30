# i18n Setup Instructions

## Step 1: Install vue-i18n

```bash
cd /Users/dengyichuan/workspace/go/src/github.com/colinrs/shopjoy/shop-admin
npm install vue-i18n@9
```

## Step 2: Update main.ts

Add the following import and app.use() call:

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import i18n from './plugins/i18n'  // ADD THIS LINE
import './styles/global.css'

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(i18n)  // ADD THIS LINE

app.mount('#app')
```

## Step 3: Verify

After installation, the `$t()` function will be available in all Vue components.

## Usage in Components

```vue
<template>
  <el-button>{{ $t('common.save') }}</el-button>
  <el-message>{{ $t('orders.loadFailed') }}</el-message>
</template>

<script setup>
// In script, you can use:
const message = i18n.global.t('orders.loadFailed')
</script>
```

## Files Created

- `src/plugins/i18n.ts` - i18n configuration plugin
- `src/locales/en.json` - English translations
- `src/locales/zh.json` - Chinese translations

## Migration Status

Currently, no Vue files have been migrated to use i18n. The infrastructure is ready.

To migrate a file, replace hardcoded strings:
- `"Save"` → `"{{ $t('common.save') }}"`
- `"加载失败"` → `"{{ $t('common.loadFailed') }}"`
