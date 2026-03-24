<template>
  <div class="redemptions-page">
    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-select v-model="searchParams.status" placeholder="兑换状态" clearable class="filter-select" @change="loadRedemptions">
            <el-option label="全部" value="" />
            <el-option label="待处理" value="pending" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            class="date-picker"
            @change="handleDateChange"
          />
        </div>
      </div>
    </el-card>

    <!-- Redemptions Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="redemptionList" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" align="center">
          <template #default="{ row }">
            <span class="id-text">#{{ row.id }}</span>
          </template>
        </el-table-column>

        <el-table-column label="用户" width="100" align="center">
          <template #default="{ row }">
            <span class="user-id">U{{ row.user_id }}</span>
          </template>
        </el-table-column>

        <el-table-column label="优惠券" min-width="150">
          <template #default="{ row }">
            <div class="coupon-cell">
              <div class="coupon-icon">
                <el-icon><Ticket /></el-icon>
              </div>
              <span class="coupon-name">{{ row.coupon_name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="使用积分" width="100" align="right">
          <template #default="{ row }">
            <span class="points-value">{{ row.points_used.toLocaleString() }}</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" effect="light" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="兑换时间" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="完成时间" width="180">
          <template #default="{ row }">
            <span v-if="row.completed_at" class="time-text">{{ formatDateTime(row.completed_at) }}</span>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <TablePagination
        :page="searchParams.page"
        :page-size="searchParams.page_size"
        :total="total"
        @change="handlePageChange"
      />
    </el-card>

    <!-- Detail Dialog -->
    <el-dialog v-model="detailDialogVisible" title="兑换详情" width="500px">
      <div v-if="currentRedemption" class="detail-content">
        <div class="detail-row">
          <span class="detail-label">兑换ID:</span>
          <span class="detail-value">#{{ currentRedemption.id }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">用户ID:</span>
          <span class="detail-value">U{{ currentRedemption.user_id }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">优惠券:</span>
          <span class="detail-value">{{ currentRedemption.coupon_name }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">使用积分:</span>
          <span class="detail-value points">{{ currentRedemption.points_used.toLocaleString() }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">状态:</span>
          <el-tag :type="getStatusTagType(currentRedemption.status)" size="small">
            {{ getStatusText(currentRedemption.status) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="detail-label">兑换时间:</span>
          <span class="detail-value">{{ formatDateTime(currentRedemption.created_at) }}</span>
        </div>
        <div v-if="currentRedemption.completed_at" class="detail-row">
          <span class="detail-label">完成时间:</span>
          <span class="detail-value">{{ formatDateTime(currentRedemption.completed_at) }}</span>
        </div>
        <div v-if="currentRedemption.user_coupon_id" class="detail-row">
          <span class="detail-label">用户优惠券ID:</span>
          <span class="detail-value">{{ currentRedemption.user_coupon_id }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Ticket } from '@element-plus/icons-vue'
import TablePagination from '@/components/common/TablePagination.vue'
import { getRedemptions, type PointsRedemption, type ListRedemptionsParams } from '@/api/points'

// State
const loading = ref(false)
const dateRange = ref<[string, string] | null>(null)
const detailDialogVisible = ref(false)
const currentRedemption = ref<PointsRedemption | null>(null)

const redemptionList = ref<PointsRedemption[]>([])
const total = ref(0)

const searchParams = reactive<ListRedemptionsParams>({
  page: 1,
  page_size: 20,
  status: '',
  start_time: '',
  end_time: ''
})

// Load functions
const loadRedemptions = async () => {
  loading.value = true
  try {
    const params: ListRedemptionsParams = {
      page: searchParams.page,
      page_size: searchParams.page_size
    }
    if (searchParams.status) params.status = searchParams.status
    if (searchParams.start_time) params.start_time = searchParams.start_time
    if (searchParams.end_time) params.end_time = searchParams.end_time

    const res = await getRedemptions(params)
    redemptionList.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to load redemptions:', error)
    // Mock data
    redemptionList.value = [
      {
        id: 101,
        user_id: 12345,
        redeem_rule_id: 1,
        coupon_id: 1,
        coupon_name: 'SAVE10',
        user_coupon_id: 10001,
        points_used: 500,
        status: 'completed',
        created_at: '2026-03-24T10:30:00Z',
        completed_at: '2026-03-24T10:30:05Z'
      },
      {
        id: 102,
        user_id: 12346,
        redeem_rule_id: 2,
        coupon_id: 2,
        coupon_name: 'SAVE20',
        user_coupon_id: null,
        points_used: 1000,
        status: 'pending',
        created_at: '2026-03-24T09:15:00Z',
        completed_at: null
      },
      {
        id: 103,
        user_id: 12347,
        redeem_rule_id: 3,
        coupon_id: 3,
        coupon_name: 'FREESHIP',
        user_coupon_id: 10002,
        points_used: 200,
        status: 'completed',
        created_at: '2026-03-23T15:00:00Z',
        completed_at: '2026-03-23T15:00:05Z'
      },
      {
        id: 104,
        user_id: 12345,
        redeem_rule_id: 1,
        coupon_id: 1,
        coupon_name: 'SAVE10',
        user_coupon_id: null,
        points_used: 500,
        status: 'cancelled',
        created_at: '2026-03-22T14:00:00Z',
        completed_at: null
      }
    ]
    total.value = 4
  } finally {
    loading.value = false
  }
}

// Helper functions
const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'completed': 'success',
    'cancelled': 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'pending': '待处理',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return texts[status] || status
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

// Handlers
const handleDateChange = () => {
  if (dateRange.value) {
    searchParams.start_time = dateRange.value[0]
    searchParams.end_time = dateRange.value[1]
  } else {
    searchParams.start_time = ''
    searchParams.end_time = ''
  }
  loadRedemptions()
}

const handlePageChange = (page: number, pageSize: number) => {
  searchParams.page = page
  searchParams.page_size = pageSize
  loadRedemptions()
}

const viewDetail = (row: PointsRedemption) => {
  currentRedemption.value = row
  detailDialogVisible.value = true
}

// Initialize
onMounted(() => {
  loadRedemptions()
})
</script>

<style scoped>
.redemptions-page {
  padding: 0;
}

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

.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
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

.coupon-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.coupon-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.coupon-name {
  font-family: 'Fira Code', monospace;
  color: #1E1B4B;
}

.points-value {
  font-weight: 600;
  color: #F59E0B;
  font-family: 'Fira Sans', sans-serif;
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

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

/* Detail Dialog */
.detail-content {
  padding: 0;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #F3F4F6;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  color: #6B7280;
}

.detail-value {
  font-weight: 500;
  color: #1E1B4B;
}

.detail-value.points {
  color: #F59E0B;
  font-family: 'Fira Sans', sans-serif;
}

:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-left {
    flex-direction: column;
  }

  .filter-select,
  .date-picker {
    width: 100%;
  }
}
</style>