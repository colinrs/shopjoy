<template>
  <el-select
    :model-value="modelValue"
    :placeholder="placeholder"
    :multiple="multiple"
    :clearable="clearable"
    filterable
    remote
    :remote-method="debouncedSearch"
    :loading="loading"
    :style="style"
    class="sku-search-select"
    @update:model-value="handleSelect"
  >
    <el-option
      v-for="item in options"
      :key="item.sku_code"
      :label="item.sku_code"
      :value="item.sku_code"
    >
      <div class="sku-option">
        <span class="sku-code">{{ item.sku_code }}</span>
        <span class="sku-product">{{ item.product_name }}</span>
      </div>
    </el-option>
  </el-select>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { searchSKUs, type SearchSKUItem } from '@/api/product'

const { t } = useI18n()

const props = withDefaults(
  defineProps<{
    modelValue: string | string[]
    multiple?: boolean
    placeholder?: string
    clearable?: boolean
    style?: string
  }>(),
  {
    multiple: false,
    clearable: true
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]]
  change: [items: SearchSKUItem[]]
}>()

const loading = ref(false)
const options = ref<SearchSKUItem[]>([])

let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadOptions = async (keyword: string) => {
  loading.value = true
  try {
    const res = await searchSKUs({ keyword, page: 1, page_size: 20 })
    options.value = res.list || []
  } catch (error) {
    console.error('Failed to search SKUs:', error)
    ElMessage.error(t('common.searchFailed') ?? 'Search failed')
    options.value = []
  } finally {
    loading.value = false
  }
}

const debouncedSearch = (keyword: string) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    if (keyword && keyword.length >= 1) {
      loadOptions(keyword)
    } else {
      options.value = []
    }
  }, 300)
}

const handleSelect = (value: string | string[]) => {
  emit('update:modelValue', value)
  const selectedCodes = Array.isArray(value) ? value : value ? [value] : []
  const selectedItems = options.value.filter(o => selectedCodes.includes(o.sku_code))
  emit('change', selectedItems)
}
</script>

<style scoped>
.sku-search-select {
  width: 100%;
}

.sku-search-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.sku-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.sku-option .sku-code {
  font-family: 'Fira Code', monospace;
  font-weight: 500;
}

.sku-option .sku-product {
  font-size: 12px;
  color: #9CA3AF;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 240px;
}
</style>
