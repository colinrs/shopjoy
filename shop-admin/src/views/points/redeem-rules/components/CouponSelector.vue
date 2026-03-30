<template>
  <div class="coupon-selector">
    <el-select
      :model-value="modelValue"
      @update:model-value="handleSelect"
      :placeholder="$t('points.selectCoupon')"
      filterable
      remote
      :remote-method="searchCoupons"
      :loading="loading"
      class="coupon-select"
    >
      <el-option
        v-for="coupon in couponList"
        :key="coupon.id"
        :label="coupon.name"
        :value="coupon.id"
      >
        <div class="coupon-option">
          <span class="coupon-name">{{ coupon.name }}</span>
          <span class="coupon-code">{{ coupon.code }}</span>
        </div>
      </el-option>
    </el-select>

    <!-- Selected Coupon Preview -->
    <div v-if="selectedCoupon" class="selected-coupon">
      <div class="coupon-preview">
        <div class="coupon-preview-icon">
          <el-icon><Ticket /></el-icon>
        </div>
        <div class="coupon-preview-info">
          <p class="preview-name">{{ selectedCoupon.name }}</p>
          <p class="preview-code">{{ $t('points.couponCode') }} {{ selectedCoupon.code }}</p>
          <div class="preview-details">
            <el-tag size="small" :type="selectedCoupon.type === 'fixed_amount' ? 'success' : 'warning'">
              {{ selectedCoupon.type === 'fixed_amount' ? $t('points.fixedAmountLabel') : $t('points.percentageLabel') }}
            </el-tag>
            <span class="preview-value">
              {{ selectedCoupon.type === 'fixed_amount' ? '$' : '' }}{{ selectedCoupon.discount_value }}{{ selectedCoupon.type === 'percentage' ? '%' : '' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Ticket } from '@element-plus/icons-vue'
import { getAvailableCoupons, type AvailableCoupon } from '@/api/points'

type Coupon = AvailableCoupon & {
  type: 'fixed_amount' | 'percentage'
}

const props = defineProps<{
  modelValue: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
  change: [coupon: Coupon]
}>()

const loading = ref(false)
const couponList = ref<Coupon[]>([])
const selectedCoupon = ref<Coupon | null>(null)

const loadCoupons = async (keyword: string = '') => {
  loading.value = true
  try {
    const res = await getAvailableCoupons({ page: 1, page_size: 20, name: keyword })
    // Map API response to include required type field
    couponList.value = (res?.list || []).map(c => ({
      ...c,
      type: c.type || 'fixed_amount' as const
    }))
  } catch (error) {
    console.error('Failed to load coupons:', error)
    ElMessage.error('加载可用优惠券失败')
  } finally {
    loading.value = false
  }
}

const searchCoupons = (keyword: string) => {
  loadCoupons(keyword)
}

const handleSelect = (value: number) => {
  emit('update:modelValue', value)
  const coupon = couponList.value.find(c => c.id === value)
  if (coupon) {
    selectedCoupon.value = coupon
    emit('change', coupon)
  }
}

// Watch modelValue to find selected coupon
watch(() => props.modelValue, (val) => {
  if (val && couponList.value.length > 0) {
    const coupon = couponList.value.find(c => c.id === val)
    if (coupon) {
      selectedCoupon.value = coupon
    }
  }
}, { immediate: true })

onMounted(() => {
  loadCoupons()
})
</script>

<style scoped>
.coupon-selector {
  width: 100%;
}

.coupon-select {
  width: 100%;
}

.coupon-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.coupon-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.coupon-option .coupon-name {
  font-weight: 500;
}

.coupon-option .coupon-code {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  color: #9CA3AF;
}

/* Selected Coupon Preview */
.selected-coupon {
  margin-top: 12px;
}

.coupon-preview {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  border: 1px solid rgba(99, 102, 241, 0.1);
}

.coupon-preview-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.coupon-preview-info {
  flex: 1;
}

.preview-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 2px 0;
}

.preview-code {
  font-size: 12px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
  margin: 0 0 8px 0;
}

.preview-details {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preview-value {
  font-weight: 600;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}
</style>