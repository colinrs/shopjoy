<template>
  <!-- Card Mode (for list display) -->
  <el-card
    v-if="!isDialog && zone"
    class="zone-card"
    shadow="never"
  >
    <div class="zone-header">
      <div class="zone-title-row">
        <el-icon class="drag-handle">
          <Rank />
        </el-icon>
        <h4 class="zone-name">
          {{ zone.name }}
        </h4>
        <el-tag
          size="small"
          type="info"
        >
          {{ getFeeTypeLabel(zone.fee_type) }}
        </el-tag>
      </div>
      <div class="zone-actions">
        <el-button
          type="primary"
          link
          size="small"
          @click="startEdit"
        >
          <el-icon><Edit /></el-icon>
          {{ $t('common.edit') }}
        </el-button>
        <el-popconfirm
          :title="$t('shipping.confirmDelete')"
          @confirm="$emit('delete', zone.id)"
        >
          <template #reference>
            <el-button
              type="danger"
              link
              size="small"
            >
              <el-icon><Delete /></el-icon>
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-popconfirm>
      </div>
    </div>

    <div class="zone-content">
      <div class="zone-row">
        <span class="label">{{ $t('shipping.deliveryArea') }}：</span>
        <span class="value">{{ formatRegions(zone.regions) }}</span>
      </div>
      <div class="zone-row">
        <span class="label">{{ $t('shipping.shippingRule') }}：</span>
        <span class="value">{{ formatFeeConfig(zone) }}</span>
      </div>
      <div
        v-if="hasFreeThreshold"
        class="zone-row"
      >
        <span class="label">{{ $t('shipping.freeShippingCondition') }}：</span>
        <span class="value">{{ formatFreeThreshold(zone) }}</span>
      </div>
    </div>
  </el-card>

  <!-- Dialog/Form Mode -->
  <el-form
    v-else
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item
      :label="$t('shipping.zoneName')"
      prop="name"
    >
      <el-input
        v-model="form.name"
        :placeholder="$t('shipping.enterZoneName')"
      />
    </el-form-item>

    <el-form-item
      :label="$t('shipping.deliveryArea')"
      prop="regions"
      required
    >
      <el-input
        v-model="selectedRegionsText"
        readonly
        :placeholder="$t('shipping.selectDeliveryArea')"
        @click="showRegionSelector"
      >
        <template #suffix>
          <el-icon><Location /></el-icon>
        </template>
      </el-input>
      <div
        v-if="form.regions.length > 0"
        class="selected-regions"
      >
        <el-tag
          v-for="region in displayRegions"
          :key="region"
          size="small"
          class="region-tag"
        >
          {{ region }}
        </el-tag>
        <span
          v-if="form.regions.length > 10"
          class="more-regions"
        >
          {{ $t('shipping.andMoreCount', { count: form.regions.length }) }}
        </span>
      </div>
    </el-form-item>

    <el-form-item
      :label="$t('shipping.feeType')"
      prop="fee_type"
    >
      <FeeTypeSelector
        :fee-type="form.fee_type"
        :first-unit="form.first_unit"
        :first-fee="form.first_fee"
        :additional-unit="form.additional_unit"
        :additional-fee="form.additional_fee"
        :volumetric-divisor="form.volumetric_divisor"
        @update="(patch) => Object.assign(form, patch)"
      />
    </el-form-item>

    <!-- P1-6: Tax Settings -->
    <el-divider content-position="left">
      <span class="section-label">{{ $t('shipping.taxSettings') }}</span>
    </el-divider>
    <el-row :gutter="16">
      <el-col :span="6">
        <el-form-item :label="$t('shipping.taxable')">
          <el-switch v-model="form.taxable" />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item :label="$t('shipping.taxRate')">
          <el-input
            v-model="form.tax_rate"
            :placeholder="$t('shipping.taxRatePlaceholder')"
            :disabled="!form.taxable"
          />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item :label="$t('shipping.taxIncluded')">
          <el-switch v-model="form.tax_included" :disabled="!form.taxable" />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item :label="$t('shipping.iossApplicable')">
          <el-switch v-model="form.ioss_applicable" :disabled="!form.taxable" />
        </el-form-item>
      </el-col>
    </el-row>

    <!-- P1-7: Remote Area Surcharge -->
    <el-divider content-position="left">
      <span class="section-label">{{ $t('shipping.remoteArea') }}</span>
    </el-divider>
    <el-row :gutter="16">
      <el-col :span="8">
        <el-form-item :label="$t('shipping.remoteSurcharge')">
          <el-input
            v-model="form.remote_surcharge"
            placeholder="0.00"
          >
            <template #prefix>
              ¥
            </template>
          </el-input>
        </el-form-item>
      </el-col>
      <el-col :span="16">
        <el-form-item :label="$t('shipping.remoteZipPatterns')">
          <el-input
            v-model="remoteZipText"
            type="textarea"
            :rows="2"
            :placeholder="$t('shipping.remoteZipPatternsPlaceholder')"
            @blur="syncRemoteZipPatterns"
          />
        </el-form-item>
      </el-col>
    </el-row>

    <!-- P1-8: Fuel Surcharge -->
    <el-divider content-position="left">
      <span class="section-label">{{ $t('shipping.fuelSurcharge') }}</span>
    </el-divider>
    <el-row :gutter="16">
      <el-col :span="8">
        <el-form-item :label="$t('shipping.fuelSurchargePct')">
          <el-input
            v-model="form.fuel_surcharge_pct"
            placeholder="0"
          >
            <template #append>
              %
            </template>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>

    <!-- P1-10: Multilingual Name -->
    <el-divider content-position="left">
      <span class="section-label">{{ $t('shipping.nameI18n') }}</span>
    </el-divider>
    <div
      v-for="(entry, idx) in form.name_i18n"
      :key="idx"
      class="i18n-row"
    >
      <el-row :gutter="8" align="middle">
        <el-col :span="8">
          <el-input
            v-model="entry.locale"
            :placeholder="$t('shipping.localePlaceholder')"
          />
        </el-col>
        <el-col :span="14">
          <el-input
            v-model="entry.name"
            :placeholder="$t('shipping.zoneName')"
          />
        </el-col>
        <el-col :span="2">
          <el-button
            link
            type="danger"
            @click="removeNameI18nEntry(idx)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </el-col>
      </el-row>
    </div>
    <el-form-item>
      <el-button
        link
        type="primary"
        @click="addNameI18nEntry"
      >
        <el-icon><Plus /></el-icon>
        {{ $t('shipping.addNameI18n') }}
      </el-button>
    </el-form-item>

    <el-form-item :label="$t('shipping.freeShippingCondition')">
      <el-row :gutter="12">
        <el-col :span="12">
          <el-checkbox v-model="form.enable_amount_threshold">
            {{ $t('shipping.freeShippingAmount') }}
          </el-checkbox>
          <el-input-number
            v-model="form.free_threshold_amount"
            :disabled="!form.enable_amount_threshold"
            :min="0"
            :precision="2"
            style="width: 100%; margin-top: 8px"
          >
            <template #prefix>
              ¥
            </template>
          </el-input-number>
        </el-col>
        <el-col :span="12">
          <el-checkbox v-model="form.enable_count_threshold">
            {{ $t('shipping.freeShippingCount') }}
          </el-checkbox>
          <el-input-number
            v-model="form.free_threshold_count"
            :disabled="!form.enable_count_threshold"
            :min="0"
            style="width: 100%; margin-top: 8px"
          />
        </el-col>
      </el-row>
    </el-form-item>

    <el-form-item class="form-actions">
      <el-button @click="$emit('cancel')">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        :loading="submitting"
        @click="handleSubmit"
      >
        {{ zone ? $t('shipping.update') : $t('shipping.add') }}
      </el-button>
    </el-form-item>
  </el-form>

  <!-- Region Selector Dialog -->
  <RegionSelector
    v-model="regionSelectorVisible"
    :selected="form.regions"
    @confirm="handleRegionConfirm"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { FormItemRule } from 'element-plus'
import { Rank, Edit, Delete, Location, Plus } from '@element-plus/icons-vue'
import type { ShippingZone, CreateZoneRequest, NameI18nEntry, Region } from '@/api/shipping'
import { getRegions } from '@/api/shipping'
import FeeTypeSelector from './FeeTypeSelector.vue'
import RegionSelector from './RegionSelector.vue'

const { t } = useI18n()

const props = defineProps<{
  zone?: ShippingZone
  index?: number
  isDialog?: boolean
}>()

const emit = defineEmits<{
  edit: [zone: ShippingZone]
  delete: [zoneId: string]
  save: [zone: CreateZoneRequest]
  cancel: []
}>()

// State
const formRef = ref()
const submitting = ref(false)
const selectedRegionsText = ref('')
const regionSelectorVisible = ref(false)
const remoteZipText = ref('')
// Track selected region objects for display; codes are stored in form.regions
const selectedRegionObjects = ref<Region[]>([])
// Map of region code → name for display of pre-existing zones (loaded on demand)
const regionNameMap = ref<Record<string, string>>({})
// Caches for lazy region-name resolution
const provincesCache = ref<Region[]>([])
const childrenCache = ref<Record<string, Region[]>>({})

interface ZoneForm extends CreateZoneRequest {
  enable_amount_threshold: boolean
  enable_count_threshold: boolean
  name_i18n: NameI18nEntry[]
}

const form = reactive<ZoneForm>({
  name: '',
  regions: [],
  fee_type: 'fixed',
  first_unit: 1,
  first_fee: '0',
  additional_unit: 1,
  additional_fee: '0',
  free_threshold_amount: '0',
  free_threshold_count: 0,
  sort: 0,
  enable_amount_threshold: false,
  enable_count_threshold: false,
  // P1-6 税费
  taxable: false,
  tax_rate: '0',
  tax_included: false,
  ioss_applicable: false,
  // P1-7 偏远地区
  remote_surcharge: '0',
  remote_zip_patterns: [],
  // P1-8 燃油附加费
  fuel_surcharge_pct: '0',
  // P1-9 体积重除数
  volumetric_divisor: 5000,
  // P1-10 多语言
  name_i18n: []
})

const rules = {
  name: [
    { required: true, message: t('shipping.zoneNameRequired'), trigger: 'blur' },
    { min: 2, max: 50, message: t('shipping.zoneNameLength'), trigger: 'blur' }
  ],
  regions: [
    {
      // eslint-disable-next-line no-unused-vars
      validator: (_rule: FormItemRule, value: string[], callback: (error?: Error) => void) => {
        if (!value || value.length === 0) {
          callback(new Error(t('shipping.selectDeliveryArea')))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  fee_type: [
    { required: true, message: t('shipping.selectFeeType'), trigger: 'change' }
  ]
}

// Computed
const hasFreeThreshold = computed(() => {
  if (!props.zone) return false
  const amount = parseFloat(props.zone.free_threshold_amount || '0')
  const count = props.zone.free_threshold_count || 0
  return amount > 0 || count > 0
})

const displayRegions = computed(() => {
  // Return first 10 region names for display (in dialog mode use the
  // currently selected Region objects; in card mode use props.zone.regions
  // resolved through the regionNameMap).
  if (selectedRegionObjects.value.length > 0) {
    return selectedRegionObjects.value.slice(0, 10).map(r => r.name)
  }
  if (props.zone?.regions) {
    return props.zone.regions.slice(0, 10).map(code => regionNameMap.value[code] || code)
  }
  return []
})

// Methods
const getFeeTypeLabel = (feeType: string) => {
  switch (feeType) {
    case 'fixed': return t('shipping.fixed')
    case 'by_count': return t('shipping.byCount')
    case 'by_weight': return t('shipping.byWeight')
    case 'by_volume': return t('shipping.byVolume')
    case 'free': return t('shipping.free')
    default: return feeType
  }
}

const formatRegions = (regions: string[]) => {
  if (!regions || regions.length === 0) return t('shipping.areaNotSet')
  // Prefer names from selectedRegionObjects (just confirmed), then regionNameMap,
  // then fall back to the code itself.
  const names = regions.map(code => {
    const obj = selectedRegionObjects.value.find(r => r.code === code)
    return obj?.name || regionNameMap.value[code] || code
  })
  if (names.length <= 3) {
    return names.join(', ')
  }
  return `${names.slice(0, 3).join(', ')} ${t('shipping.andMoreCount', { count: names.length })}`
}

const formatFeeConfig = (zone: ShippingZone) => {
  switch (zone.fee_type) {
    case 'fixed':
      return t('shipping.feeConfigFixed', { fee: zone.first_fee })
    case 'by_count':
      return t('shipping.feeConfigByCount', { first: zone.first_unit, fee: zone.first_fee, add: zone.additional_unit, addFee: zone.additional_fee })
    case 'by_weight':
      return t('shipping.feeConfigByWeight', { first: zone.first_unit, fee: zone.first_fee, add: zone.additional_unit, addFee: zone.additional_fee })
    case 'by_volume':
      return t('shipping.feeConfigByVolume', { first: zone.first_unit, fee: zone.first_fee, add: zone.additional_unit, addFee: zone.additional_fee })
    case 'free':
      return t('shipping.free')
    default:
      return zone.fee_type
  }
}

const formatFreeThreshold = (zone: ShippingZone) => {
  const parts: string[] = []
  const amount = parseFloat(zone.free_threshold_amount || '0')
  const count = zone.free_threshold_count || 0

  if (amount > 0) {
    parts.push(t('shipping.freeThresholdAmount', { amount }))
  }
  if (count > 0) {
    parts.push(t('shipping.freeThresholdCount', { count }))
  }
  return parts.join(' or ') + ' ' + t('shipping.freeShippingApplied')
}

const startEdit = () => {
  emit('edit', props.zone!)
}

const showRegionSelector = () => {
  regionSelectorVisible.value = true
}

const handleRegionConfirm = (regions: Region[]) => {
  // Store full Region objects for display
  selectedRegionObjects.value = [...regions]
  // Store codes for API submission
  form.regions = regions.map(r => r.code)
  // Update name map so future renders use names
  for (const r of regions) {
    regionNameMap.value[r.code] = r.name
  }
  selectedRegionsText.value = formatRegions(form.regions)
}

// Lazily resolve names for any codes we don't yet know.
// For codes already in regionNameMap this is a no-op; otherwise we
// load provinces and (for cities) the parent's children to populate the map.
const resolveRegionNames = async (codes: string[]) => {
  const missing = codes.filter(c => !regionNameMap.value[c])
  if (missing.length === 0) return
  try {
    // Load provinces first; many missing codes will resolve at the province level
    if (provincesCache.value.length === 0) {
      provincesCache.value = await getRegions()
    }
    for (const p of provincesCache.value) {
      regionNameMap.value[p.code] = p.name
    }
    // For any codes still missing, try fetching their parent's children
    const stillMissing = missing.filter(c => !regionNameMap.value[c])
    for (const code of stillMissing) {
      // Assume the code's parent is a province we've loaded; pick the first
      // province that has this code as a child.
      for (const p of provincesCache.value) {
        // Cache children to avoid repeated calls
        if (!childrenCache.value[p.code]) {
          try {
            childrenCache.value[p.code] = await getRegions(p.code)
          } catch {
            childrenCache.value[p.code] = []
          }
        }
        const child = childrenCache.value[p.code].find(c => c.code === code)
        if (child) {
          regionNameMap.value[code] = child.name
          break
        }
      }
    }
  } catch (err) {
    // Non-fatal: display will fall back to showing the code
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        // Sync textarea → array before submit
        syncRemoteZipPatterns()
        // Filter empty i18n entries
        const cleanNameI18n = (form.name_i18n || []).filter(
          (e) => e.locale && e.locale.trim() && e.name && e.name.trim()
        )
        const data: CreateZoneRequest = {
          name: form.name,
          name_i18n: cleanNameI18n.length > 0 ? cleanNameI18n : undefined,
          regions: form.regions,
          fee_type: form.fee_type,
          first_unit: form.first_unit,
          first_fee: form.first_fee,
          additional_unit: form.additional_unit,
          additional_fee: form.additional_fee,
          free_threshold_amount: form.enable_amount_threshold ? form.free_threshold_amount : '0',
          free_threshold_count: form.enable_count_threshold ? form.free_threshold_count : 0,
          // P1-6
          taxable: form.taxable,
          tax_rate: form.taxable ? (form.tax_rate || '0') : undefined,
          tax_included: form.taxable ? form.tax_included : undefined,
          ioss_applicable: form.taxable ? form.ioss_applicable : undefined,
          // P1-7
          remote_surcharge: form.remote_surcharge || '0',
          remote_zip_patterns: form.remote_zip_patterns || [],
          // P1-8
          fuel_surcharge_pct: form.fuel_surcharge_pct || '0',
          // P1-9
          volumetric_divisor: form.fee_type === 'by_volume' ? (form.volumetric_divisor || 5000) : undefined,
          sort: form.sort
        }
        emit('save', data)
      } finally {
        submitting.value = false
      }
    }
  })
}

// Multilingual name helpers
const addNameI18nEntry = () => {
  if (!form.name_i18n) form.name_i18n = []
  form.name_i18n.push({ locale: '', name: '' })
}

const removeNameI18nEntry = (idx: number) => {
  if (!form.name_i18n) return
  form.name_i18n.splice(idx, 1)
}

// Remote ZIP patterns: textarea (one per line) → string[]
const syncRemoteZipPatterns = () => {
  const lines = remoteZipText.value
    .split(/[\n,]+/)
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
  form.remote_zip_patterns = lines
}

// Initialize form when editing
watch(() => props.zone, (newZone) => {
  if (newZone && props.isDialog) {
    form.name = newZone.name
    form.regions = newZone.regions || []
    form.fee_type = newZone.fee_type
    form.first_unit = newZone.first_unit
    form.first_fee = newZone.first_fee
    form.additional_unit = newZone.additional_unit
    form.additional_fee = newZone.additional_fee
    form.sort = newZone.sort
    const amount = parseFloat(newZone.free_threshold_amount || '0')
    const count = newZone.free_threshold_count || 0
    form.free_threshold_amount = newZone.free_threshold_amount
    form.free_threshold_count = count
    form.enable_amount_threshold = amount > 0
    form.enable_count_threshold = count > 0
    // P1-6
    form.taxable = newZone.taxable
    form.tax_rate = newZone.tax_rate || '0'
    form.tax_included = newZone.tax_included
    form.ioss_applicable = newZone.ioss_applicable
    // P1-7
    form.remote_surcharge = newZone.remote_surcharge || '0'
    form.remote_zip_patterns = newZone.remote_zip_patterns || []
    remoteZipText.value = (newZone.remote_zip_patterns || []).join('\n')
    // P1-8
    form.fuel_surcharge_pct = newZone.fuel_surcharge_pct || '0'
    // P1-9
    form.volumetric_divisor = newZone.volumetric_divisor || 5000
    // P1-10
    form.name_i18n = newZone.name_i18n ? newZone.name_i18n.map((e) => ({ ...e })) : []
    // Reset in-memory selection so display reflects the zone being edited
    selectedRegionObjects.value = []
    selectedRegionsText.value = formatRegions(form.regions)
    // Lazily fetch region names so display shows "Beijing" instead of "110000"
    if (form.regions.length > 0) {
      resolveRegionNames(form.regions)
    }
  } else if (newZone && !props.isDialog) {
    // Card mode: just resolve names for display
    if (newZone.regions && newZone.regions.length > 0) {
      resolveRegionNames(newZone.regions)
    }
  }
}, { immediate: true })
</script>

<style scoped>
.zone-card {
  border-radius: 12px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  margin-bottom: 12px;
}

.zone-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.zone-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.drag-handle {
  cursor: move;
  color: #9CA3AF;
}

.zone-name {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.zone-actions {
  display: flex;
  gap: 8px;
}

.zone-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.zone-row {
  display: flex;
  align-items: flex-start;
}

.zone-row .label {
  width: 80px;
  flex-shrink: 0;
  font-size: 13px;
  color: #6B7280;
}

.zone-row .value {
  flex: 1;
  font-size: 13px;
  color: #1E1B4B;
}

.selected-regions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 8px;
}

.region-tag {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
}

.more-regions {
  font-size: 12px;
  color: #6B7280;
  line-height: 24px;
}

.form-actions {
  margin-top: 24px;
  text-align: right;
}

/* Tags */
:deep(.el-tag--info) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: #6366F1;
}

.i18n-row {
  margin-bottom: 8px;
}
</style>
