<template>
  <div class="transaction-table">
    <!-- Filters -->
    <div v-if="showFilters" class="filter-bar">
      <el-select
        v-model="filterType"
        placeholder="交易类型"
        clearable
        class="filter-select"
        @change="handleFilterChange"
      >
        <el-option label="全部" value="" />
        <el-option label="获得" value="EARN" />
        <el-option label="兑换" value="REDEEM" />
        <el-option label="调整" value="ADJUST" />
        <el-option label="过期" value="EXPIRE" />
        <el-option label="冻结" value="FREEZE" />
        <el-option label="解冻" value="UNFREEZE" />
      </el-select>
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        class="date-picker"
        @change="handleFilterChange"
      />
    </div>

    <!-- Table -->
    <el-table :data="transactions" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" align="center">
        <template #default="{ row }">
          <span class="id-text">#{{ row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column v-if="showUser" label="用户" width="100" align="center">
        <template #default="{ row }">
          <span class="user-id">U{{ row.user_id }}</span>
        </template>
      </el-table-column>

      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)" effect="light" size="small">
            {{ getTypeText(row.type) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="积分" width="120" align="right">
        <template #default="{ row }">
          <span class="points-value" :class="{ positive: row.points > 0, negative: row.points < 0 }">
            {{ row.points > 0 ? '+' : '' }}{{ row.points.toLocaleString() }}
          </span>
        </template>
      </el-table-column>

      <el-table-column label="变动后余额" width="120" align="right">
        <template #default="{ row }">
          <span class="balance-text">{{ row.balance_after.toLocaleString() }}</span>
        </template>
      </el-table-column>

      <el-table-column label="描述" min-width="200">
        <template #default="{ row }">
          <div class="description-cell">
            <p class="description-text">{{ row.description }}</p>
            <p v-if="row.reference_type" class="reference-text">
              {{ row.reference_type }}: {{ row.reference_id }}
            </p>
          </div>
        </template>
      </el-table-column>

      <el-table-column v-if="showExpires" label="过期时间" width="160">
        <template #default="{ row }">
          <span v-if="row.expires_at" class="time-text">{{ formatDateTime(row.expires_at) }}</span>
          <span v-else class="no-data">-</span>
        </template>
      </el-table-column>

      <el-table-column label="时间" width="160">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
    </el-table>

    <!-- Pagination -->
    <div v-if="showPagination && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { PointsTransaction } from '@/api/points'

const props = withDefaults(defineProps<{
  transactions: PointsTransaction[]
  loading?: boolean
  showFilters?: boolean
  showUser?: boolean
  showExpires?: boolean
  showPagination?: boolean
  total?: number
  page?: number
  pageSize?: number
}>(), {
  loading: false,
  showFilters: false,
  showUser: false,
  showExpires: false,
  showPagination: false,
  total: 0,
  page: 1,
  pageSize: 10
})

const emit = defineEmits<{
  filter: [filters: { type: string; startDate: string; endDate: string }]
  'page-change': [page: number, pageSize: number]
}>()

// Local state
const filterType = ref('')
const dateRange = ref<[string, string] | null>(null)
const currentPage = ref(props.page)
const pageSize = ref(props.pageSize)

// Watch props changes
watch(() => props.page, (val) => {
  currentPage.value = val
})

watch(() => props.pageSize, (val) => {
  pageSize.value = val
})

// Methods
const getTypeTagType = (type: string) => {
  const types: Record<string, string> = {
    'EARN': 'success',
    'REDEEM': 'primary',
    'ADJUST': 'warning',
    'EXPIRE': 'danger',
    'FREEZE': 'info',
    'UNFREEZE': 'info'
  }
  return types[type] || 'info'
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    'EARN': '获得',
    'REDEEM': '兑换',
    'ADJUST': '调整',
    'EXPIRE': '过期',
    'FREEZE': '冻结',
    'UNFREEZE': '解冻'
  }
  return texts[type] || type
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleFilterChange = () => {
  emit('filter', {
    type: filterType.value,
    startDate: dateRange.value?.[0] || '',
    endDate: dateRange.value?.[1] || ''
  })
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  emit('page-change', currentPage.value, val)
}

const handlePageChange = (val: number) => {
  currentPage.value = val
  emit('page-change', val, pageSize.value)
}
</script>

<style scoped>
.transaction-table {
  background: #fff;
  border-radius: 16px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
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

.id-text {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.user-id {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  color: #6B7280;
}

.points-value {
  font-weight: 600;
  font-family: 'Fira Sans', sans-serif;
}

.points-value.positive {
  color: #10B981;
}

.points-value.negative {
  color: #EF4444;
}

.balance-text {
  color: #1E1B4B;
  font-family: 'Fira Sans', sans-serif;
}

.description-cell {
  max-width: 100%;
}

.description-text {
  margin: 0;
  color: #1E1B4B;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reference-text {
  margin: 4px 0 0 0;
  font-size: 12px;
  color: #9CA3AF;
  font-family: 'Fira Code', monospace;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.no-data {
  color: #9CA3AF;
  font-size: 13px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid #F3F4F6;
  margin-top: 16px;
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--primary) {
  background-color: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.2);
  color: #3B82F6;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
  }

  .filter-select,
  .date-picker {
    width: 100%;
  }
}
</style>