<template>
  <!-- Card Mode (for list display) -->
  <el-card v-if="!isDialog" class="zone-card" shadow="never">
    <div class="zone-header">
      <div class="zone-title-row">
        <el-icon class="drag-handle"><Rank /></el-icon>
        <h4 class="zone-name">{{ zone.name }}</h4>
        <el-tag size="small" type="info">{{ getFeeTypeLabel(zone.fee_type) }}</el-tag>
      </div>
      <div class="zone-actions">
        <el-button type="primary" link size="small" @click="startEdit">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
        <el-popconfirm title="确认删除该配送区域？" @confirm="$emit('delete', zone.id)">
          <template #reference>
            <el-button type="danger" link size="small">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-popconfirm>
      </div>
    </div>

    <div class="zone-content">
      <div class="zone-row">
        <span class="label">配送地区：</span>
        <span class="value">{{ formatRegions(zone.regions) }}</span>
      </div>
      <div class="zone-row">
        <span class="label">运费规则：</span>
        <span class="value">{{ formatFeeConfig(zone) }}</span>
      </div>
      <div class="zone-row" v-if="hasFreeThreshold">
        <span class="label">包邮条件：</span>
        <span class="value">{{ formatFreeThreshold(zone) }}</span>
      </div>
    </div>
  </el-card>

  <!-- Dialog/Form Mode -->
  <el-form v-else :model="form" :rules="rules" ref="formRef" label-width="100px">
    <el-form-item label="区域名称" prop="name">
      <el-input v-model="form.name" placeholder="例如：华东地区、偏远地区" />
    </el-form-item>

    <el-form-item label="配送地区" prop="regions" required>
      <el-input
        v-model="selectedRegionsText"
        readonly
        placeholder="点击选择配送地区"
        @click="showRegionSelector"
      >
        <template #suffix>
          <el-icon><Location /></el-icon>
        </template>
      </el-input>
      <div class="selected-regions" v-if="form.regions.length > 0">
        <el-tag
          v-for="region in displayRegions"
          :key="region"
          size="small"
          class="region-tag"
        >
          {{ region }}
        </el-tag>
        <span v-if="form.regions.length > 10" class="more-regions">
          等共 {{ form.regions.length }} 个城市
        </span>
      </div>
    </el-form-item>

    <el-form-item label="计费方式" prop="fee_type">
      <FeeTypeSelector v-model="form" />
    </el-form-item>

    <el-form-item label="包邮设置">
      <el-row :gutter="12">
        <el-col :span="12">
          <el-checkbox v-model="form.enable_amount_threshold">
            满额包邮
          </el-checkbox>
          <el-input-number
            v-model="form.free_threshold_amount"
            :disabled="!form.enable_amount_threshold"
            :min="0"
            :precision="2"
            style="width: 100%; margin-top: 8px"
          >
            <template #prefix>¥</template>
          </el-input-number>
        </el-col>
        <el-col :span="12">
          <el-checkbox v-model="form.enable_count_threshold">
            满件包邮
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
      <el-button @click="$emit('cancel')">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        {{ zone ? '更新区域' : '添加区域' }}
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Rank, Edit, Delete, Location } from '@element-plus/icons-vue'
import type { ShippingZone, CreateZoneRequest } from '@/api/shipping'
import FeeTypeSelector from './FeeTypeSelector.vue'

const props = defineProps<{
  zone?: ShippingZone
  index?: number
  isDialog?: boolean
}>()

const emit = defineEmits<{
  update: [zone: ShippingZone]
  delete: [zoneId: number]
  save: [zone: CreateZoneRequest]
  cancel: []
}>()

// State
const formRef = ref()
const submitting = ref(false)
const selectedRegionsText = ref('')

interface ZoneForm extends CreateZoneRequest {
  enable_amount_threshold: boolean
  enable_count_threshold: boolean
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
  enable_count_threshold: false
})

const rules = {
  name: [
    { required: true, message: '请输入区域名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度为2-50个字符', trigger: 'blur' }
  ],
  regions: [
    {
      validator: (_rule: any, value: string[], callback: any) => {
        if (!value || value.length === 0) {
          callback(new Error('请选择配送地区'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  fee_type: [
    { required: true, message: '请选择计费方式', trigger: 'change' }
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
  // Return first 10 regions for display
  return props.zone?.regions?.slice(0, 10) || []
})

// Methods
const getFeeTypeLabel = (feeType: string) => {
  const labels: Record<string, string> = {
    fixed: '固定运费',
    by_count: '按件计费',
    by_weight: '按重量计费',
    free: '免运费'
  }
  return labels[feeType] || feeType
}

const formatRegions = (regions: string[]) => {
  if (!regions || regions.length === 0) return '未设置'
  if (regions.length <= 3) {
    return regions.join('、')
  }
  return `${regions.slice(0, 3).join('、')} 等 ${regions.length} 个城市`
}

const formatFeeConfig = (zone: ShippingZone) => {
  switch (zone.fee_type) {
    case 'fixed':
      return `固定运费：¥${zone.first_fee}`
    case 'by_count':
      return `首${zone.first_unit}件 ¥${zone.first_fee}，续${zone.additional_unit}件 ¥${zone.additional_fee}`
    case 'by_weight':
      return `首${zone.first_unit}g ¥${zone.first_fee}，续${zone.additional_unit}g ¥${zone.additional_fee}`
    case 'free':
      return '免运费'
    default:
      return zone.fee_type
  }
}

const formatFreeThreshold = (zone: ShippingZone) => {
  const parts: string[] = []
  const amount = parseFloat(zone.free_threshold_amount || '0')
  const count = zone.free_threshold_count || 0

  if (amount > 0) {
    parts.push(`满 ¥${amount}`)
  }
  if (count > 0) {
    parts.push(`满 ${count} 件`)
  }
  return parts.join(' 或 ') + ' 包邮'
}

const startEdit = () => {
  emit('update', props.zone!)
}

const showRegionSelector = () => {
  // TODO: Implement region selector dialog
  ElMessage.info('地区选择器功能开发中')
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitting.value = true
      try {
        const data: CreateZoneRequest = {
          name: form.name,
          regions: form.regions,
          fee_type: form.fee_type,
          first_unit: form.first_unit,
          first_fee: form.first_fee,
          additional_unit: form.additional_unit,
          additional_fee: form.additional_fee,
          free_threshold_amount: form.enable_amount_threshold ? form.free_threshold_amount : '0',
          free_threshold_count: form.enable_count_threshold ? form.free_threshold_count : 0,
          sort: form.sort
        }
        emit('save', data)
      } finally {
        submitting.value = false
      }
    }
  })
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
    selectedRegionsText.value = formatRegions(form.regions)
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
</style>