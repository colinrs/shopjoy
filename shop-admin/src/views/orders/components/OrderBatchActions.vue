<template>
  <transition name="slide-down">
    <div v-if="selectedCount > 0" class="batch-actions-bar">
      <div class="batch-info">
        <el-icon><Check /></el-icon>
        <span>{{ $t('orders.ordersSelected', { count: selectedCount }) }}</span>
      </div>
      <div class="batch-buttons">
        <el-button type="primary" @click="handleBatchShip">
          <el-icon><Van /></el-icon>
          {{ $t('orders.batchShip') }}
        </el-button>
        <el-button type="danger" @click="handleBatchCancel">
          <el-icon><Close /></el-icon>
          {{ $t('orders.batchCancel') }}
        </el-button>
        <el-button @click="handleBatchRemark">
          <el-icon><Edit /></el-icon>
          {{ $t('orders.batchRemark') }}
        </el-button>
        <el-button @click="handleClearSelection">{{ $t('orders.clearSelection') }}</el-button>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { Check, Van, Close, Edit } from '@element-plus/icons-vue'

defineProps<{
  selectedCount: number
}>()

const emit = defineEmits<{
  'batch-ship': []
  'batch-cancel': []
  'batch-remark': []
  'clear-selection': []
}>()

const handleBatchShip = () => {
  emit('batch-ship')
}

const handleBatchCancel = () => {
  emit('batch-cancel')
}

const handleBatchRemark = () => {
  emit('batch-remark')
}

const handleClearSelection = () => {
  emit('clear-selection')
}
</script>

<style scoped>
.batch-actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  margin-bottom: 16px;
}

.batch-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #4F46E5;
}

.batch-buttons {
  display: flex;
  gap: 8px;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .batch-actions-bar {
    flex-direction: column;
    gap: 12px;
  }

  .batch-buttons {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
