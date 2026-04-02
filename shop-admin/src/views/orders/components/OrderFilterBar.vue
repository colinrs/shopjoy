<template>
  <el-card
    class="filter-card"
    shadow="never"
  >
    <div class="filter-bar">
      <div class="filter-left">
        <el-input
          v-model="localSearchQuery"
          :placeholder="$t('orders.searchPlaceholder')"
          class="search-input"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="localStatusFilter"
          :placeholder="$t('orders.filterOrderStatus')"
          clearable
          class="filter-select"
          @change="handleSearch"
        >
          <el-option
            :label="$t('common.all')"
            value=""
          />
          <el-option
            :label="$t('orders.pendingPayment')"
            value="pending_payment"
          />
          <el-option
            :label="$t('orders.paid')"
            value="paid"
          />
          <el-option
            :label="$t('orders.shipped')"
            value="shipped"
          />
          <el-option
            :label="$t('orders.delivered')"
            value="delivered"
          />
          <el-option
            :label="$t('orders.cancelled')"
            value="cancelled"
          />
          <el-option
            :label="$t('orders.refunded')"
            value="refunded"
          />
        </el-select>
        <el-select
          v-model="localFulfillmentFilter"
          :placeholder="$t('orders.filterFulfillment')"
          clearable
          class="filter-select"
          @change="handleSearch"
        >
          <el-option
            :label="$t('common.all')"
            value=""
          />
          <el-option
            :label="$t('orders.unshipped')"
            value="0"
          />
          <el-option
            :label="$t('orders.partialShipped')"
            value="1"
          />
          <el-option
            :label="$t('orders.shipped')"
            value="2"
          />
          <el-option
            :label="$t('orders.delivered')"
            value="3"
          />
        </el-select>
        <el-date-picker
          v-model="localDateRange"
          type="daterange"
          :range-separator="$t('common.to')"
          :start-placeholder="$t('orders.startDate')"
          :end-placeholder="$t('orders.endDate')"
          class="date-picker"
          value-format="YYYY-MM-DD"
          @change="handleSearch"
        />
      </div>
      <div class="filter-right">
        <el-button
          :loading="exporting"
          @click="handleExport"
        >
          <el-icon><Download /></el-icon>
          {{ $t('common.export') }}
        </el-button>
        <el-button
          type="primary"
          @click="handleRefresh"
        >
          <el-icon><Refresh /></el-icon>
          {{ $t('common.refresh') }}
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Search, Download, Refresh } from '@element-plus/icons-vue'
import type { OrderStatus, FulfillmentStatus } from '@/api/order'

const props = defineProps<{
  searchQuery: string
  statusFilter: OrderStatus | ''
  fulfillmentFilter: FulfillmentStatus | ''
  dateRange: [string, string] | null
  exporting: boolean
}>()

const emit = defineEmits<{
  'update:searchQuery': [value: string]
  'update:statusFilter': [value: OrderStatus | '']
  'update:fulfillmentFilter': [value: FulfillmentStatus | '']
  'update:dateRange': [value: [string, string] | null]
  search: []
  refresh: []
  export: []
}>()

const localSearchQuery = computed({
  get: () => props.searchQuery,
  set: (val) => emit('update:searchQuery', val)
})

const localStatusFilter = computed({
  get: () => props.statusFilter,
  set: (val) => emit('update:statusFilter', val)
})

const localFulfillmentFilter = computed({
  get: () => props.fulfillmentFilter,
  set: (val) => emit('update:fulfillmentFilter', val)
})

const localDateRange = computed({
  get: () => props.dateRange,
  set: (val) => emit('update:dateRange', val)
})

const handleSearch = () => {
  emit('search')
}

const handleRefresh = () => {
  emit('refresh')
}

const handleExport = () => {
  emit('export')
}
</script>

<style scoped>
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 220px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.date-picker {
  width: 260px;
}

.date-picker :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select,
  .date-picker {
    width: 100%;
  }
}
</style>
