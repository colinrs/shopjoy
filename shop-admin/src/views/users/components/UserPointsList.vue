<template>
  <div class="user-points-list">
    <el-table v-loading="loading" :data="txns" stripe>
      <el-table-column :label="$t('users.pointsColumns.time')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.pointsColumns.type')" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.type)" size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.pointsColumns.change')" width="120" align="right">
        <template #default="{ row }">
          <span :class="row.points >= 0 ? 'positive' : 'negative'">
            {{ row.points >= 0 ? '+' : '' }}{{ row.points }}
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.pointsColumns.balance')" width="100" align="right">
        <template #default="{ row }">
          <span class="balance">{{ row.balance_after }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.pointsColumns.description')" min-width="180">
        <template #default="{ row }">
          <span>{{ row.description || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.pointsColumns.reference')" width="180">
        <template #default="{ row }">
          <el-button v-if="row.reference_type === 'ORDER' && row.reference_id" link type="primary" size="small" @click="goOrder(row.reference_id)">
            {{ $t('users.pointsColumns.reference') }}: {{ row.reference_id }}
          </el-button>
          <span v-else-if="row.reference_id">{{ row.reference_type }}:{{ row.reference_id }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
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

    <el-empty v-if="!loading && txns.length === 0" :description="$t('users.noPoints')" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getUserPointsTransactions, type UserPointsTransaction } from '@/api/user'
import { ElMessage } from 'element-plus'
import { t } from '@/plugins/i18n'

const props = defineProps<{ userId?: string }>()
const router = useRouter()

const loading = ref(false)
const txns = ref<UserPointsTransaction[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserPointsTransactions(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    txns.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    console.error('Failed to load points:', err)
    ElMessage.error(t('users.loadPointsFailed'))
    txns.value = []
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

const getTypeTag = (t: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (t === 'EARN') return 'success'
  if (t === 'REDEEM') return 'primary'
  if (t === 'EXPIRE') return 'warning'
  return 'info'
}

const goOrder = (id: string) => router.push(`/orders/${id}`)

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-points-list { padding: 0; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.positive { color: #10B981; font-weight: 600; }
.negative { color: #EF4444; font-weight: 600; }
.balance { font-weight: 600; color: #1E1B4B; }
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>