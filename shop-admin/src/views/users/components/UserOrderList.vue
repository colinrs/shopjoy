<template>
  <div class="user-order-list">
    <el-table
      v-loading="loading"
      :data="orders"
      stripe
    >
      <el-table-column
        :label="$t('users.orderNo')"
        min-width="180"
      >
        <template #default="{ row }">
          <span class="order-no">{{ row.order_no }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('users.orderStatus')"
        width="120"
        align="center"
      >
        <template #default="{ row }">
          <el-tag
            :type="getStatusType(row.status)"
            size="small"
          >
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('users.orderItemCount')"
        width="100"
        align="center"
      >
        <template #default="{ row }">
          <span>{{ row.item_count ?? '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('users.orderTotalAmount')"
        width="140"
        align="right"
      >
        <template #default="{ row }">
          <span class="amount">{{ row.currency }} {{ row.total_amount }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('users.orderCreatedAt')"
        width="180"
      >
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.actions')"
        width="120"
        align="center"
        fixed="right"
      >
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            @click="goDetail(row)"
          >
            {{ $t('users.viewDetail') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div
      v-if="total > pageSize"
      class="pagination-wrapper"
    >
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty
      v-if="!loading && orders.length === 0 && !error"
      :description="$t('users.noOrders')"
    />
    <el-empty
      v-if="!loading && error"
      :description="error"
      image-error
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getUserOrders, type UserOrderListItem } from '@/api/user'
import { t } from '@/plugins/i18n'

const props = defineProps<{ userId?: string }>()
const router = useRouter()

const loading = ref(false)
const orders = ref<UserOrderListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const error = ref<string>('')

const load = async () => {
  if (!props.userId) return
  loading.value = true
  error.value = ''
  try {
    const res = await getUserOrders(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    orders.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    console.error('Failed to load orders:', err)
    error.value = t('users.loadOrdersFailed')
    orders.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (status === 'completed') return 'success'
  if (status === 'cancelled' || status === 'refunded') return 'info'
  if (status === 'paid' || status === 'shipped') return 'primary'
  if (status === 'pending_payment' || status === 'pending_shipment') return 'warning'
  return 'danger'
}

const goDetail = (row: UserOrderListItem) => {
  router.push(`/orders/${row.order_id}`)
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-order-list { padding: 0; }
.order-no { font-family: 'Fira Code', monospace; font-weight: 500; }
.amount { font-family: 'Fira Sans', sans-serif; font-weight: 600; color: #EF4444; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>