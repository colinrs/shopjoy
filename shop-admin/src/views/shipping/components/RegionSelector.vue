<template>
  <el-dialog
    v-model="visible"
    :title="$t('shipping.selectArea')"
    width="700px"
    :close-on-click-modal="false"
  >
    <div class="region-selector">
      <!-- Country Selector -->
      <div class="country-row">
        <span class="row-label">{{ $t('shipping.country') }}</span>
        <el-select
          v-model="selectedCountryCode"
          :placeholder="$t('shipping.selectCountry')"
          filterable
          style="width: 100%"
          @change="handleCountryChange"
        >
          <el-option
            v-for="c in countries"
            :key="c.code"
            :label="c.name"
            :value="c.code"
          />
        </el-select>
      </div>

      <!-- Search (only after country is chosen) -->
      <el-input
        v-if="selectedCountryCode"
        v-model="searchText"
        :placeholder="$t('shipping.searchCity')"
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <!-- Selected Summary -->
      <div
        v-if="selectedRegions.length > 0"
        class="selected-summary"
      >
        <span class="summary-text">{{ $t('shipping.selectedCount', { count: selectedRegions.length }) }}</span>
        <el-button
          type="primary"
          link
          size="small"
          @click="clearSelection"
        >
          {{ $t('shipping.clear') }}
        </el-button>
      </div>

      <div
        v-loading="loading"
        class="region-list"
      >
        <!-- Province list (only when no search and a country is chosen) -->
        <div
          v-if="selectedCountryCode && !searchText && currentProvinces.length > 0"
          class="region-level"
        >
          <div class="level-header">
            <span class="level-title">{{ $t('shipping.province') }}</span>
          </div>
          <div
            v-for="province in currentProvinces"
            :key="province.code"
            class="province-block"
          >
            <div
              class="province-row"
              :class="{ expanded: expandedProvinces.has(province.code) }"
            >
              <el-checkbox
                v-model="provinceChecked[province.code]"
                :indeterminate="provinceIndeterminate[province.code]"
                @change="(val: boolean) => handleProvinceCheck(province, val)"
              />
              <span
                class="province-name"
                @click="toggleExpand(province)"
              >{{ province.name }}</span>
              <el-icon
                v-if="loadingCitiesFor.has(province.code)"
                class="loading-icon"
              >
                <Loading />
              </el-icon>
              <el-button
                link
                size="small"
                @click="toggleExpand(province)"
              >
                <el-icon>
                  <component :is="expandedProvinces.has(province.code) ? ArrowDown : ArrowRight" />
                </el-icon>
              </el-button>
            </div>

            <!-- City grid (shown when province is expanded) -->
            <div
              v-if="expandedProvinces.has(province.code)"
              class="city-grid"
            >
              <div
                v-for="city in citiesMap[province.code] || []"
                :key="city.code"
                class="city-item"
                :class="{ selected: isCitySelected(city.code) }"
                @click="handleCityClick(city)"
              >
                <el-checkbox
                  :model-value="isCitySelected(city.code)"
                  @change="(val: boolean) => handleCityCheck(city, val)"
                  @click.stop
                />
                <span class="city-name">{{ city.name }}</span>
              </div>
              <div
                v-if="(citiesMap[province.code] || []).length === 0 && !loadingCitiesFor.has(province.code)"
                class="city-empty"
              >
                {{ $t('shipping.noData') }}
              </div>
            </div>
          </div>
        </div>

        <!-- Empty province list -->
        <el-empty
          v-else-if="selectedCountryCode && !searchText && currentProvinces.length === 0 && !loading"
          :description="$t('shipping.noProvinces')"
          :image-size="60"
        />

        <!-- Search Results -->
        <div
          v-if="searchText"
          class="region-level"
        >
          <div class="level-header">
            <span class="level-title">{{ $t('shipping.searchResults') }}</span>
          </div>
          <div
            v-if="searchResults.length > 0"
            class="city-grid"
          >
            <div
              v-for="city in searchResults"
              :key="city.code"
              class="city-item"
              :class="{ selected: isCitySelected(city.code) }"
              @click="handleCityClick(city)"
            >
              <el-checkbox
                :model-value="isCitySelected(city.code)"
                @change="(val: boolean) => handleCityCheck(city, val)"
                @click.stop
              />
              <span class="city-name">{{ city.name }} ({{ city.provinceName }})</span>
            </div>
          </div>
          <el-empty
            v-else-if="!loading"
            :description="$t('shipping.noResults')"
            :image-size="60"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleCancel">
        {{ $t('shipping.cancel') }}
      </el-button>
      <el-button
        type="primary"
        :disabled="selectedRegions.length === 0"
        @click="handleConfirm"
      >
        {{ $t('shipping.confirmSelection') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Search, ArrowRight, ArrowDown, Loading } from '@element-plus/icons-vue'
import { getRegions, type Region } from '@/api/shipping'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()

interface CityWithProvince extends Region {
  provinceName: string
}

const props = defineProps<{
  modelValue: boolean
  selected?: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: [regions: Region[]]
}>()

// ── State ──
// Level 1: countries (level 1)
const countries = ref<Region[]>([])
const selectedCountryCode = ref<string>('')

// Level 2: provinces for the selected country
const provincesMap = ref<Record<string, Region[]>>({})

// Level 3: cities for each loaded province
const citiesMap = ref<Record<string, Region[]>>({})
const loadedProvinces = ref<Set<string>>(new Set())
const loadingCitiesFor = ref<Set<string>>(new Set())

// UI state
const loading = ref(false)
const searchText = ref('')
const expandedProvinces = ref<Set<string>>(new Set())
const selectedRegions = ref<Region[]>([])
// Per-province checkbox state for v-model (immediate visual feedback)
// and per-province indeterminate state.
const provinceChecked = reactive<Record<string, boolean>>({})
const provinceIndeterminate = reactive<Record<string, boolean>>({})

// ── Computed ──
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const currentProvinces = computed<Region[]>(() => {
  if (!selectedCountryCode.value) return []
  return provincesMap.value[selectedCountryCode.value] || []
})

const searchResults = computed<CityWithProvince[]>(() => {
  if (!searchText.value) return []
  const keyword = searchText.value.toLowerCase()
  const result: CityWithProvince[] = []
  const provinces = currentProvinces.value
  for (const province of provinces) {
    const cities = citiesMap.value[province.code] || []
    for (const city of cities) {
      if (city.name.toLowerCase().includes(keyword)) {
        result.push({ ...city, provinceName: province.name })
      }
    }
  }
  return result
})

// ── Methods ──
const loadCountries = async () => {
  loading.value = true
  try {
    const data = await getRegions()
    countries.value = data
  } catch (error) {
    handleError(error, t('shipping.loadCountriesFailed'))
  } finally {
    loading.value = false
  }
}

const loadProvinces = async (countryCode: string) => {
  if (provincesMap.value[countryCode]) return
  loading.value = true
  try {
    const data = await getRegions(countryCode)
    provincesMap.value[countryCode] = data
  } catch (error) {
    handleError(error, t('shipping.loadProvincesFailed'))
  } finally {
    loading.value = false
  }
}

const loadCities = async (province: Region) => {
  if (loadedProvinces.value.has(province.code)) return
  loadingCitiesFor.value.add(province.code)
  try {
    const data = await getRegions(province.code)
    citiesMap.value[province.code] = data
    loadedProvinces.value.add(province.code)
  } catch (error) {
    handleError(error, t('shipping.loadCitiesFailed'))
  } finally {
    loadingCitiesFor.value.delete(province.code)
  }
}

const isCitySelected = (cityCode: string) =>
  selectedRegions.value.some(r => r.code === cityCode)

// Keep the v-model state maps in sync with the actual selection state.
// v-model needs reactive state that updates synchronously on click;
// deriving from selectedRegions alone had timing issues with el-checkbox.
const recomputeProvinceState = () => {
  const provinces = currentProvinces.value
  for (const province of provinces) {
    const children = citiesMap.value[province.code]
    // Province hasn't been expanded yet — leave state as the user set it
    if (!children) continue
    if (children.length === 0) {
      // Leaf province (e.g. US states): checked iff itself is in selectedRegions
      provinceChecked[province.code] = selectedRegions.value.some(r => r.code === province.code)
      provinceIndeterminate[province.code] = false
      continue
    }
    const count = children.filter(c => selectedRegions.value.some(r => r.code === c.code)).length
    provinceChecked[province.code] = count === children.length
    provinceIndeterminate[province.code] = count > 0 && count < children.length
  }
}

// React to changes in cities / selection → recompute checkbox state.
watch(
  [() => Object.keys(citiesMap.value).length, selectedRegions, () => currentProvinces.value.length],
  () => recomputeProvinceState(),
  { deep: true }
)

const handleCountryChange = async (countryCode: string) => {
  expandedProvinces.value = new Set()
  // Clear per-province checkbox state for the previous country
  for (const key of Object.keys(provinceChecked)) delete provinceChecked[key]
  for (const key of Object.keys(provinceIndeterminate)) delete provinceIndeterminate[key]
  await loadProvinces(countryCode)
}

const toggleExpand = async (province: Region) => {
  if (expandedProvinces.value.has(province.code)) {
    expandedProvinces.value.delete(province.code)
    return
  }
  expandedProvinces.value.add(province.code)
  if (!loadedProvinces.value.has(province.code)) {
    await loadCities(province)
  }
}

const handleProvinceCheck = async (province: Region, checked: boolean) => {
  // The v-model has already updated provinceChecked[province.code] synchronously,
  // giving the user immediate visual feedback. Now finalize the selection
  // by loading children (if needed) and updating selectedRegions.
  if (!loadedProvinces.value.has(province.code)) {
    await loadCities(province)
  }

  const children = citiesMap.value[province.code] || []

  // If this province has no children (e.g. US states — level 3 counties
  // are not populated), select the province itself.
  if (children.length === 0) {
    if (checked) {
      if (!selectedRegions.value.some(r => r.code === province.code)) {
        selectedRegions.value.push({ ...province })
      }
    } else {
      selectedRegions.value = selectedRegions.value.filter(r => r.code !== province.code)
    }
    return
  }

  // Auto-expand so the user sees what was selected
  if (checked) {
    expandedProvinces.value.add(province.code)
  }

  if (checked) {
    for (const child of children) {
      if (!selectedRegions.value.some(r => r.code === child.code)) {
        selectedRegions.value.push({ ...child })
      }
    }
  } else {
    const codes = new Set(children.map(c => c.code))
    selectedRegions.value = selectedRegions.value.filter(r => !codes.has(r.code))
  }
}

const handleCityClick = (city: Region) => {
  handleCityCheck(city, !isCitySelected(city.code))
}

const handleCityCheck = (city: Region, checked: boolean) => {
  if (checked) {
    if (!selectedRegions.value.some(r => r.code === city.code)) {
      selectedRegions.value.push({ ...city })
    }
  } else {
    selectedRegions.value = selectedRegions.value.filter(r => r.code !== city.code)
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

// Reset on dialog open
watch(visible, (val) => {
  if (val) {
    loadCountries()
    selectedCountryCode.value = ''
    provincesMap.value = {}
    citiesMap.value = {}
    loadedProvinces.value = new Set()
    expandedProvinces.value = new Set()
    for (const key of Object.keys(provinceChecked)) delete provinceChecked[key]
    for (const key of Object.keys(provinceIndeterminate)) delete provinceIndeterminate[key]
    selectedRegions.value = []
    searchText.value = ''
  }
})
</script>

<style scoped>
.region-selector {
  min-height: 400px;
}

.country-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.row-label {
  flex-shrink: 0;
  font-size: 14px;
  color: #606266;
  width: 56px;
}

.search-input {
  margin-bottom: 12px;
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

.province-block {
  margin-bottom: 4px;
}

.province-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  background: #F9FAFB;
  transition: background 0.2s;
}

.province-row:hover {
  background: #F3F4F6;
}

.province-row.expanded {
  background: rgba(64, 158, 255, 0.08);
}

.province-name {
  flex: 1;
  font-size: 14px;
  color: #1F2937;
  font-weight: 500;
  cursor: pointer;
}

.city-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 4px;
  padding: 8px 12px 12px 36px;
}

.city-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.15s;
}

.city-item:hover {
  background: #F3F4F6;
}

.city-item.selected {
  background: rgba(64, 158, 255, 0.1);
}

.city-name {
  font-size: 13px;
  color: #4B5563;
}

.city-empty {
  grid-column: 1 / -1;
  text-align: center;
  color: #9CA3AF;
  font-size: 13px;
  padding: 12px;
}

.selected-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: rgba(64, 158, 255, 0.06);
  border-radius: 6px;
  margin-bottom: 12px;
}

.summary-text {
  font-size: 13px;
  color: #409EFF;
  font-weight: 500;
}

@media (max-width: 600px) {
  .city-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>