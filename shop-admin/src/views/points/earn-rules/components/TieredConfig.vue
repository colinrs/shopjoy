<template>
  <div class="tiered-config">
    <div class="tier-list">
      <div
        v-for="(tier, index) in modelValue"
        :key="index"
        class="tier-row"
      >
        <div class="tier-index">{{ index + 1 }}</div>

        <!-- Threshold -->
        <div class="tier-field">
          <label>门槛 ($)</label>
          <el-input-number
            v-if="tier.threshold !== null"
            v-model="tier.threshold"
            :min="0"
            :step="100"
            :precision="0"
            placeholder="金额"
            style="width: 140px"
          />
          <el-input
            v-else
            value="无上限"
            disabled
            style="width: 140px"
          />
        </div>

        <!-- Ratio -->
        <div class="tier-field">
          <label>积分/$1</label>
          <el-input-number
            :model-value="parseFloat(tier.ratio) || 0"
            @update:model-value="(val: number) => updateRatio(index, val)"
            :min="0.1"
            :max="100"
            :precision="2"
            :step="0.5"
            style="width: 140px"
          />
        </div>

        <!-- Actions -->
        <div class="tier-actions">
          <el-button
            v-if="modelValue.length > 1"
            type="danger"
            link
            size="small"
            @click="removeTier(index)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>

    <el-button type="primary" link @click="addTier">
      <el-icon><Plus /></el-icon>
      添加阶梯
    </el-button>

    <!-- Preview -->
    <div class="preview-section">
      <div class="preview-title">计算预览</div>
      <div class="preview-examples">
        <div v-for="example in previewExamples" :key="example.amount" class="preview-item">
          <span class="preview-label">订单 ${{ example.amount }}:</span>
          <span class="preview-value">{{ example.points }} 积分</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { TierConfig } from '@/api/points'

const props = defineProps<{
  modelValue: TierConfig[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: TierConfig[]]
}>()

const addTier = () => {
  const tiers = [...props.modelValue]
  // Find the last tier with threshold and make it unlimited
  const lastThresholdIndex = tiers.length - 1
  const newThreshold = tiers[lastThresholdIndex].threshold
    ? tiers[lastThresholdIndex].threshold! + 10000
    : 10000

  tiers.splice(lastThresholdIndex, 0, {
    threshold: newThreshold,
    ratio: '1.0'
  })

  emit('update:modelValue', tiers)
}

const removeTier = (index: number) => {
  const tiers = [...props.modelValue]
  tiers.splice(index, 1)
  emit('update:modelValue', tiers)
}

// Update ratio value (convert number to string for API compatibility)
const updateRatio = (index: number, value: number) => {
  const tiers = [...props.modelValue]
  tiers[index] = { ...tiers[index], ratio: String(value) }
  emit('update:modelValue', tiers)
}

// Calculate points for tiered config (incremental calculation)
// Example: $600 with tiers [$100, $500, null] at ratios [1.0, 1.5, 2.0]
// = (100 * 1.0) + (400 * 1.5) + (100 * 2.0) = 100 + 600 + 200 = 900 points
const calculatePoints = (amount: number): number => {
  const tiers = props.modelValue
  let points = 0
  let remaining = amount
  let prevThreshold = 0

  for (let i = 0; i < tiers.length; i++) {
    const tier = tiers[i]
    const ratio = parseFloat(tier.ratio) || 0
    const threshold = tier.threshold

    if (remaining <= 0) break

    if (threshold === null) {
      // Last tier - unlimited, apply to all remaining amount
      points += Math.floor(remaining * ratio)
      break
    }

    // Calculate the amount in this tier's range
    // Tier covers (prevThreshold, threshold]
    const tierAmount = Math.min(remaining, threshold - prevThreshold)
    points += Math.floor(tierAmount * ratio)
    remaining -= tierAmount
    prevThreshold = threshold
  }

  return points
}

const previewExamples = computed(() => {
  return [
    { amount: 50, points: calculatePoints(50) },
    { amount: 200, points: calculatePoints(200) },
    { amount: 600, points: calculatePoints(600) }
  ]
})
</script>

<style scoped>
.tiered-config {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.tier-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
}

.tier-row {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  padding: 12px;
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.tier-index {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.tier-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tier-field label {
  font-size: 12px;
  color: #6B7280;
}

.tier-actions {
  display: flex;
  align-items: center;
}

/* Preview Section */
.preview-section {
  margin-top: 16px;
  padding: 12px;
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.preview-title {
  font-size: 13px;
  font-weight: 500;
  color: #1E1B4B;
  margin-bottom: 8px;
}

.preview-examples {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.preview-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.preview-label {
  font-size: 12px;
  color: #6B7280;
}

.preview-value {
  font-size: 12px;
  font-weight: 600;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

/* Responsive */
@media (max-width: 576px) {
  .tier-row {
    flex-wrap: wrap;
  }

  .preview-examples {
    flex-direction: column;
    gap: 8px;
  }
}
</style>