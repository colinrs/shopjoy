<template>
  <el-dialog
    v-model="visible"
    title="选择配送地区"
    width="600px"
    :close-on-click-modal="false"
  >
    <div class="region-selector">
      <!-- Search -->
      <el-input
        v-model="searchText"
        placeholder="搜索城市名称"
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <!-- Province List -->
      <div class="region-list" v-loading="loading">
        <div v-if="!searchText" class="region-level">
          <div class="level-header">
            <span class="level-title">省份/直辖市</span>
            <el-checkbox
              v-model="allProvincesSelected"
              :indeterminate="provinceIndeterminate"
              @change="handleSelectAllProvinces"
            >
              全选
            </el-checkbox>
          </div>
          <div class="region-grid">
            <div
              v-for="province in provinces"
              :key="province.code"
              class="region-item"
              :class="{ selected: isProvinceSelected(province.code) }"
              @click="handleProvinceClick(province)"
            >
              <el-checkbox
                :model-value="isProvinceSelected(province.code)"
                @change="(val: boolean) => handleProvinceCheck(province.code, val)"
                @click.stop
              />
              <span class="region-name">{{ province.name }}</span>
              <el-icon v-if="province.children?.length" class="expand-icon"><ArrowRight /></el-icon>
            </div>
          </div>
        </div>

        <!-- City List (shown when province is expanded) -->
        <div v-if="expandedProvince && !searchText" class="region-level">
          <div class="level-header">
            <el-button link @click="expandedProvince = null">
              <el-icon><ArrowLeft /></el-icon>
              返回省份列表
            </el-button>
            <span class="level-title">{{ expandedProvinceName }} - 城市</span>
            <el-checkbox
              v-model="allCitiesSelected"
              :indeterminate="cityIndeterminate"
              @change="handleSelectAllCities"
            >
              全选
            </el-checkbox>
          </div>
          <div class="region-grid">
            <div
              v-for="city in currentCities"
              :key="city.code"
              class="region-item"
              :class="{ selected: selectedRegions.includes(city.code) }"
              @click="handleCityClick(city)"
            >
              <el-checkbox
                :model-value="selectedRegions.includes(city.code)"
                @change="(val: boolean) => handleCityCheck(city.code, val)"
                @click.stop
              />
              <span class="region-name">{{ city.name }}</span>
            </div>
          </div>
        </div>

        <!-- Search Results -->
        <div v-if="searchText" class="region-level">
          <div class="level-header">
            <span class="level-title">搜索结果</span>
          </div>
          <div class="region-grid">
            <div
              v-for="city in filteredCities"
              :key="city.code"
              class="region-item"
              :class="{ selected: selectedRegions.includes(city.code) }"
            >
              <el-checkbox
                :model-value="selectedRegions.includes(city.code)"
                @change="(val: boolean) => handleCityCheck(city.code, val)"
              />
              <span class="region-name">{{ city.name }} ({{ city.provinceName }})</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Selected Summary -->
      <div class="selected-summary" v-if="selectedRegions.length > 0">
        <span class="summary-text">已选择 {{ selectedRegions.length }} 个城市</span>
        <el-button type="primary" link size="small" @click="clearSelection">清空</el-button>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleConfirm" :disabled="selectedRegions.length === 0">
        确认选择
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Search, ArrowRight, ArrowLeft } from '@element-plus/icons-vue'
import { getRegions, type Region } from '@/api/shipping'

interface CityWithProvince extends Region {
  provinceName: string
}

const props = defineProps<{
  modelValue: boolean
  selected?: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: [regions: string[]]
}>()

// State
const loading = ref(false)
const searchText = ref('')
const provinces = ref<Region[]>([])
const citiesMap = ref<Record<string, Region[]>>({})
const selectedRegions = ref<string[]>([...props.selected || []])
const expandedProvince = ref<string | null>(null)

// Computed
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const expandedProvinceName = computed(() => {
  if (!expandedProvince.value) return ''
  const province = provinces.value.find(p => p.code === expandedProvince.value)
  return province?.name || ''
})

const currentCities = computed(() => {
  if (!expandedProvince.value) return []
  return citiesMap.value[expandedProvince.value] || []
})

const filteredCities = computed<CityWithProvince[]>(() => {
  if (!searchText.value) return []
  const result: CityWithProvince[] = []
  const keyword = searchText.value.toLowerCase()

  for (const [provinceCode, cities] of Object.entries(citiesMap.value)) {
    const province = provinces.value.find(p => p.code === provinceCode)
    for (const city of cities) {
      if (city.name.toLowerCase().includes(keyword)) {
        result.push({ ...city, provinceName: province?.name || '' })
      }
    }
  }
  return result
})

const allProvincesSelected = computed(() => {
  // Check if all province capital cities are selected
  for (const province of provinces.value) {
    const cities = citiesMap.value[province.code] || []
    const capitalCity = cities[0]?.code // First city is usually the capital
    if (capitalCity && !selectedRegions.value.includes(capitalCity)) {
      return false
    }
  }
  return provinces.value.length > 0
})

const provinceIndeterminate = computed(() => {
  let selectedCount = 0
  for (const province of provinces.value) {
    const cities = citiesMap.value[province.code] || []
    if (cities.length > 0 && selectedRegions.value.includes(cities[0].code)) {
      selectedCount++
    }
  }
  return selectedCount > 0 && selectedCount < provinces.value.length
})

const allCitiesSelected = computed(() => {
  if (!expandedProvince.value) return false
  const cities = citiesMap.value[expandedProvince.value] || []
  return cities.length > 0 && cities.every(c => selectedRegions.value.includes(c.code))
})

const cityIndeterminate = computed(() => {
  if (!expandedProvince.value) return false
  const cities = citiesMap.value[expandedProvince.value] || []
  const selectedCount = cities.filter(c => selectedRegions.value.includes(c.code)).length
  return selectedCount > 0 && selectedCount < cities.length
})

// Methods
const loadProvinces = async () => {
  loading.value = true
  try {
    const data = await getRegions()
    provinces.value = data
  } catch (error) {
    console.error('Failed to load provinces:', error)
  } finally {
    loading.value = false
  }
}

const loadCities = async (provinceCode: string) => {
  if (citiesMap.value[provinceCode]) return // Already loaded

  try {
    const data = await getRegions(provinceCode)
    citiesMap.value[provinceCode] = data
  } catch (error) {
    console.error('Failed to load cities:', error)
  }
}

const isProvinceSelected = (provinceCode: string) => {
  const cities = citiesMap.value[provinceCode] || []
  if (cities.length === 0) return false
  // Consider province selected if all its cities are selected
  return cities.every(c => selectedRegions.value.includes(c.code))
}

const handleProvinceClick = async (province: Region) => {
  await loadCities(province.code)
  expandedProvince.value = province.code
}

const handleProvinceCheck = async (provinceCode: string, checked: boolean) => {
  await loadCities(provinceCode)
  const cities = citiesMap.value[provinceCode] || []

  if (checked) {
    // Add all cities
    for (const city of cities) {
      if (!selectedRegions.value.includes(city.code)) {
        selectedRegions.value.push(city.code)
      }
    }
  } else {
    // Remove all cities
    const cityCodes = cities.map(c => c.code)
    selectedRegions.value = selectedRegions.value.filter(code => !cityCodes.includes(code))
  }
}

const handleCityClick = (city: Region) => {
  handleCityCheck(city.code, !selectedRegions.value.includes(city.code))
}

const handleCityCheck = (cityCode: string, checked: boolean) => {
  if (checked) {
    if (!selectedRegions.value.includes(cityCode)) {
      selectedRegions.value.push(cityCode)
    }
  } else {
    selectedRegions.value = selectedRegions.value.filter(code => code !== cityCode)
  }
}

const handleSelectAllProvinces = async (checked: boolean) => {
  for (const province of provinces.value) {
    await loadCities(province.code)
    const cities = citiesMap.value[province.code] || []
    if (checked) {
      for (const city of cities) {
        if (!selectedRegions.value.includes(city.code)) {
          selectedRegions.value.push(city.code)
        }
      }
    } else {
      const cityCodes = cities.map(c => c.code)
      selectedRegions.value = selectedRegions.value.filter(code => !cityCodes.includes(code))
    }
  }
}

const handleSelectAllCities = (checked: boolean) => {
  if (!expandedProvince.value) return
  const cities = citiesMap.value[expandedProvince.value] || []

  if (checked) {
    for (const city of cities) {
      if (!selectedRegions.value.includes(city.code)) {
        selectedRegions.value.push(city.code)
      }
    }
  } else {
    const cityCodes = cities.map(c => c.code)
    selectedRegions.value = selectedRegions.value.filter(code => !cityCodes.includes(code))
  }
}

const clearSelection = () => {
  selectedRegions.value = []
}

const handleCancel = () => {
  visible.value = false
}

const handleConfirm = () => {
  emit('confirm', [...selectedRegions.value])
  visible.value = false
}

// Watch for dialog open
watch(visible, (val) => {
  if (val) {
    loadProvinces()
    selectedRegions.value = [...props.selected || []]
    expandedProvince.value = null
    searchText.value = ''
  }
})
</script>

<style scoped>
.region-selector {
  min-height: 400px;
}

.search-input {
  margin-bottom: 16px;
}

.region-list {
  max-height: 400px;
  overflow-y: auto;
}

.region-level {
  margin-bottom: 16px;
}

.level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #EBEEF5;
}

.level-title {
  font-weight: 600;
  color: #303133;
}

.region-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.region-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.region-item:hover {
  background: #F5F7FA;
}

.region-item.selected {
  background: rgba(64, 158, 255, 0.1);
}

.region-name {
  flex: 1;
  font-size: 13px;
  color: #606266;
}

.expand-icon {
  color: #C0C4CC;
}

.selected-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #F5F7FA;
  border-radius: 6px;
  margin-top: 16px;
}

.summary-text {
  font-size: 14px;
  color: #409EFF;
  font-weight: 500;
}

@media (max-width: 600px) {
  .region-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>