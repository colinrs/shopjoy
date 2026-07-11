<template>
  <div class="user-operation-log">
    <div class="filter-bar">
      <el-select
        v-model="filterAction"
        :placeholder="$t('users.logFilterAll')"
        clearable
        class="action-filter"
        @change="reload"
      >
        <el-option :label="$t('users.logFilterAll')" value="" />
        <el-option
          v-for="opt in actionOptions"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>
    </div>

    <el-table v-loading="loading" :data="logs" stripe>
      <el-table-column :label="$t('users.logColumns.time')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.logColumns.action')" width="160" align="center">
        <template #default="{ row }">
          <el-tag :type="getActionType(row.action)" size="small">
            {{ row.action_text }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.operator')" width="120">
        <template #default="{ row }">
          <span>{{ row.operator_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.logColumns.ip')" width="140">
        <template #default="{ row }">
          <span class="mono">{{ row.ip_address || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.logColumns.reason')" min-width="180">
        <template #default="{ row }">
          <span>{{ row.reason || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.logColumns.userAgent')" width="80" align="center">
        <template #default="{ row }">
          <el-tooltip v-if="row.user_agent" :content="row.user_agent" placement="top">
            <el-button link type="primary" size="small">{{ $t('common.view') }}</el-button>
          </el-tooltip>
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

    <el-empty v-if="!loading && logs.length === 0 && !error" :description="$t('users.noLogs')" />
    <el-empty v-if="!loading && error" :description="error" image-error />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getUserOperationLogs, type UserOperationLogItem } from '@/api/user'

const props = defineProps<{ userId?: string }>()
const { t } = useI18n()

const loading = ref(false)
const logs = ref<UserOperationLogItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const filterAction = ref<string>('')
const error = ref<string>('')

const actionOptions = [
  { value: 'CREATE_USER', label: t('users.logActions.CREATE_USER') },
  { value: 'UPDATE_USER', label: t('users.logActions.UPDATE_USER') },
  { value: 'SUSPEND_USER', label: t('users.logActions.SUSPEND_USER') },
  { value: 'SUSPEND_WITH_REASON', label: t('users.logActions.SUSPEND_WITH_REASON') },
  { value: 'ACTIVATE_USER', label: t('users.logActions.ACTIVATE_USER') },
  { value: 'DELETE_USER', label: t('users.logActions.DELETE_USER') },
  { value: 'RESET_PASSWORD', label: t('users.logActions.RESET_PASSWORD') }
]

const load = async () => {
  if (!props.userId) return
  loading.value = true
  error.value = ''
  try {
    const res = await getUserOperationLogs(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value,
      action: filterAction.value || undefined
    })
    logs.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    console.error('Failed to load operation logs:', err)
    error.value = t('users.loadLogsFailed')
    logs.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const reload = () => {
  currentPage.value = 1
  load()
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

const getActionType = (action: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (action.startsWith('SUSPEND') || action === 'DELETE_USER') return 'danger'
  if (action === 'ACTIVATE_USER') return 'success'
  if (action === 'RESET_PASSWORD') return 'warning'
  if (action === 'UPDATE_USER') return 'primary'
  return 'info'
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-operation-log { padding: 0; }
.filter-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }
.action-filter { width: 200px; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.mono { font-family: 'Fira Code', monospace; font-size: 12px; color: #374151; }
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>