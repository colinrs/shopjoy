<template>
  <div class="operation-log">
    <div class="log-header">
      <span class="header-title">{{ $t('users.operationLogs') }}</span>
      <el-button size="small" @click="loadLogs">
        <el-icon><Refresh /></el-icon>
        {{ $t('common.refresh') }}
      </el-button>
    </div>

    <el-timeline v-if="logs.length > 0">
      <el-timeline-item
        v-for="log in logs"
        :key="log.id"
        :timestamp="formatDateTime(log.created_at)"
        placement="top"
        :type="getLogType(log.action)"
      >
        <el-card class="log-card" shadow="never">
          <div class="log-content">
            <span class="log-action">{{ log.action_text }}</span>
            <span class="log-operator">{{ $t('users.operator') }}: {{ log.operator_name || $t('users.system') }}</span>
          </div>
          <div class="log-detail" v-if="log.detail">
            <span class="detail-label">{{ $t('users.detail') }}: </span>
            <span class="detail-text">{{ log.detail }}</span>
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>

    <el-empty v-if="!loading && logs.length === 0" :description="$t('users.noLogs')" />

    <div class="pagination-wrapper" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="loadLogs"
        @current-change="loadLogs"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { t } from '@/plugins/i18n'
import request from '@/utils/request'

interface OperationLog {
  id: number
  user_id: number
  action: string
  action_text: string
  detail: string
  operator_id: number
  operator_name: string
  created_at: string
}

const props = defineProps<{
  userId?: number
}>()

const loading = ref(false)
const logs = ref<OperationLog[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const loadLogs = async () => {
  if (!props.userId) return

  loading.value = true
  try {
    const res = await request<{ list: OperationLog[]; total: number }>({
      url: `/api/v1/users/${props.userId}/logs`,
      method: 'get',
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    logs.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to load logs:', error)
    ElMessage.error(t('users.loadLogsFailed'))
  } finally {
    loading.value = false
  }
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

const getLogType = (action: string) => {
  const types: Record<string, string> = {
    'login': 'primary',
    'register': 'success',
    'profile_update': 'warning',
    'order_create': 'primary',
    'suspend': 'danger',
    'activate': 'success'
  }
  return types[action] || 'info'
}

watch(() => props.userId, () => {
  loadLogs()
}, { immediate: true })

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.operation-log {
  padding: 0;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
}

.log-card {
  border-radius: 12px;
}

.log-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-action {
  font-weight: 600;
  color: #1E1B4B;
}

.log-operator {
  font-size: 13px;
  color: #6B7280;
}

.log-detail {
  margin-top: 8px;
  font-size: 13px;
  color: #6B7280;
}

.detail-label {
  color: #9CA3AF;
}

.detail-text {
  color: #374151;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}
</style>
