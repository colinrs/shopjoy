<template>
  <div class="users-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card primary">
          <p class="stat-value">{{ stats.total }}</p>
          <p class="stat-label">{{ $t('users.totalUsers') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card success">
          <p class="stat-value">{{ stats.active }}</p>
          <p class="stat-label">{{ $t('users.activeUsers') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card warning">
          <p class="stat-value">{{ stats.new_today }}</p>
          <p class="stat-label">{{ $t('users.newUsersToday') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card danger">
          <p class="stat-value">{{ stats.suspended }}</p>
          <p class="stat-label">{{ $t('users.suspendedUsers') }}</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('users.searchPlaceholder')"
            class="search-input"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" :placeholder="$t('users.accountStatus')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('users.all')" :value="0" />
            <el-option :label="$t('users.normal')" :value="1" />
            <el-option :label="$t('users.disabled')" :value="2" />
          </el-select>
          <el-button @click="showAdvancedFilters = !showAdvancedFilters">
            <el-icon><Filter /></el-icon>
            {{ $t('users.advancedFilters') }}
            <el-icon class="filter-arrow" :class="{ 'is-expanded': showAdvancedFilters }"><ArrowDown /></el-icon>
          </el-button>
        </div>
        <div class="filter-right">
          <el-button @click="handleExport" :loading="exportLoading">
            <el-icon><Download /></el-icon>{{ $t('common.export') }}
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>{{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
      <!-- Advanced Filters -->
      <div v-if="showAdvancedFilters" class="advanced-filters">
        <el-row :gutter="16">
          <el-col :span="6">
            <el-input
              v-model="pointsMin"
              :placeholder="$t('users.pointsMin')"
              type="number"
              clearable
              @clear="clearAdvancedFilters"
            >
              <template #prefix>{{ $t('users.points') }}</template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="pointsMax"
              :placeholder="$t('users.pointsMax')"
              type="number"
              clearable
              @clear="clearAdvancedFilters"
            >
              <template #prefix>{{ $t('users.points') }}</template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="orderCountMin"
              :placeholder="$t('users.orderCountMin')"
              type="number"
              clearable
              @clear="clearAdvancedFilters"
            >
              <template #prefix>{{ $t('users.orders') }}</template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="orderCountMax"
              :placeholder="$t('users.orderCountMax')"
              type="number"
              clearable
              @clear="clearAdvancedFilters"
            >
              <template #prefix>{{ $t('users.orders') }}</template>
            </el-input>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- Bulk Actions -->
    <div class="bulk-actions" v-if="selectedUsers.length > 0">
      <span class="selected-count">{{ $t('users.selectedCount', { count: selectedUsers.length }) }}</span>
      <el-button size="small" type="success" @click="handleBatchActivate">{{ $t('users.batchActivate') }}</el-button>
      <el-button size="small" type="warning" @click="handleBatchSuspend">{{ $t('users.batchSuspend') }}</el-button>
      <el-button size="small" @click="handleClearSelection">{{ $t('common.clearSelection') }}</el-button>
    </div>

    <!-- Users Table -->
    <el-card class="table-card" shadow="never">
      <Table ref="tableRef" :data="userList" :loading="loading" @selection-change="handleSelectionChange">
        <el-table-column :label="$t('users.userInfo')" min-width="250">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="48" :src="row.avatar" class="user-avatar">
                {{ row.name ? row.name.charAt(0) : 'U' }}
              </el-avatar>
              <div class="user-details">
                <p class="user-name">{{ row.name }}</p>
                <p class="user-email">{{ row.email }}</p>
                <p class="user-phone">{{ row.phone || '-' }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('users.registrationTime')" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_login" :label="$t('users.lastLogin')" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ row.last_login || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('users.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.status === 1"
              @change="(val: boolean) => confirmStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column v-if="isEnhanced" :label="$t('users.pointsBalance')" width="120" align="center">
          <template #default="{ row }">
            <span class="enhanced-value">{{ (row as ExtendedUser).points_balance || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isEnhanced" :label="$t('users.orderCount')" width="100" align="center">
          <template #default="{ row }">
            <span class="enhanced-value">{{ (row as ExtendedUser).order_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="isEnhanced" :label="$t('users.totalSpent')" width="120" align="center">
          <template #default="{ row }">
            <span class="enhanced-value">{{ (row as ExtendedUser).total_spent || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="handleEdit(row)">
              {{ $t('common.edit') }}
            </el-button>
            <el-button type="primary" link size="small" @click.stop="handleDetail(row)">
              {{ $t('users.viewDetail') }}
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button type="primary" link size="small">
                {{ $t('users.more') }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="resetPassword">{{ $t('users.resetPassword') }}</el-dropdown-item>
                  <el-dropdown-item divided command="delete" style="color: #EF4444;">{{ $t('users.deleteAccount') }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </Table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog v-model="editDialogVisible" :title="$t('users.editUser')" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item :label="$t('users.username')">
          <el-input v-model="editForm.name" :placeholder="$t('users.enterUsername')" />
        </el-form-item>
        <el-form-item :label="$t('common.avatar')">
          <el-input v-model="editForm.avatar" :placeholder="$t('common.avatarUrl')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmEdit" :loading="editLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- User Detail Dialog -->
    <el-dialog v-model="detailDialogVisible" :title="$t('users.userDetail')" width="600px">
      <el-descriptions :column="2" border v-if="currentUser">
        <el-descriptions-item :label="$t('users.userId')">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('users.username')">{{ currentUser.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('users.email')">{{ currentUser.email }}</el-descriptions-item>
        <el-descriptions-item :label="$t('users.mobile')">{{ currentUser.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('users.status')">
          <el-tag :type="currentUser.status === 1 ? 'success' : 'danger'">
            {{ currentUser.status === 1 ? $t('users.enabled') : $t('users.disabled') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('users.createdAt')">{{ currentUser.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('users.lastLogin')">{{ currentUser.last_login || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, ArrowDown, Download, Filter } from '@element-plus/icons-vue'
import { t } from '@/plugins/i18n'
import {
  getUserList,
  getUserListEnhanced,
  getUserStats,
  updateUser,
  suspendUser,
  activateUser,
  deleteUser,
  resetUserPassword,
  exportUsers,
  batchUpdateUserStatus,
  type User,
  type UserStatus,
  type UserStats,
  type ExtendedUser,
  type BatchUpdateUserStatusRequest
} from '@/api/user'
import Table from '@/components/common/Table.vue'

const router = useRouter()
const loading = ref(false)
const editLoading = ref(false)
const exportLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<UserStatus | undefined>(undefined)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentUser = ref<User | null>(null)

// Advanced filters
const showAdvancedFilters = ref(false)

// Batch selection
const selectedUsers = ref<User[]>([])
const tableRef = ref()

const handleSelectionChange = (selection: User[]) => {
  selectedUsers.value = selection
}

const handleClearSelection = () => {
  tableRef.value?.clearSelection()
  selectedUsers.value = []
}

const pointsMin = ref<number | null>(null)
const pointsMax = ref<number | null>(null)
const orderCountMin = ref<number | null>(null)
const orderCountMax = ref<number | null>(null)

// Check if enhanced filters are applied
const isEnhanced = computed(() => {
  return showAdvancedFilters.value ||
    pointsMin.value != null ||
    pointsMax.value != null ||
    orderCountMin.value != null ||
    orderCountMax.value != null
})

const stats = ref<UserStats>({
  total: 0,
  active: 0,
  suspended: 0,
  new_today: 0
})

const userList = ref<User[]>([])

const editForm = reactive({
  id: 0,
  name: '',
  avatar: ''
})

// Load user list
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchQuery.value || undefined,
      status: statusFilter.value || undefined,
      points_min: pointsMin.value ?? undefined,
      points_max: pointsMax.value ?? undefined,
      order_count_min: orderCountMin.value ?? undefined,
      order_count_max: orderCountMax.value ?? undefined
    }

    // Use enhanced API when advanced filters are applied
    if (isEnhanced.value) {
      const res = await getUserListEnhanced(params)
      userList.value = res.list || []
      total.value = res.total || 0
    } else {
      const res = await getUserList(params)
      userList.value = res.list || []
      total.value = res.total || 0
    }
  } catch (error) {
    console.error('Failed to load users:', error)
    ElMessage.error(t('users.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Load user stats
const loadStats = async () => {
  try {
    const res = await getUserStats()
    stats.value = res
  } catch (error) {
    console.error('Failed to load stats:', error)
    ElMessage.error(t('users.loadStatsFailed'))
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

// Clear advanced filters - set to null instead of empty string
const clearAdvancedFilters = () => {
  pointsMin.value = null
  pointsMax.value = null
  orderCountMin.value = null
  orderCountMax.value = null
  handleSearch()
}

const handleRefresh = () => {
  loadUsers()
  loadStats()
}

const handleEdit = (row: User) => {
  editForm.id = row.id
  editForm.name = row.name
  editForm.avatar = row.avatar || ''
  editDialogVisible.value = true
}

const confirmEdit = async () => {
  editLoading.value = true
  try {
    await updateUser(editForm.id, {
      name: editForm.name,
      avatar: editForm.avatar || undefined
    })
    ElMessage.success(t('users.updateSuccess'))
    editDialogVisible.value = false
    loadUsers()
  } catch (error) {
    console.error('Failed to update user:', error)
    ElMessage.error(t('users.updateFailed'))
  } finally {
    editLoading.value = false
  }
}

const handleDetail = (row: User) => {
  router.push(`/users/${row.id}`)
}

const handleExport = async () => {
  exportLoading.value = true
  try {
    const blob = await exportUsers({
      page: 1,
      page_size: 10000,
      keyword: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `users_${new Date().toISOString().split('T')[0]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success(t('users.exportSuccess'))
  } catch (error) {
    console.error('Failed to export users:', error)
    ElMessage.error(t('users.exportFailed'))
  } finally {
    exportLoading.value = false
  }
}

const handleCommand = (command: string, row: User) => {
  if (command === 'resetPassword') {
    handleResetPassword(row)
  } else if (command === 'delete') {
    handleDelete(row)
  }
}

const confirmStatusChange = (row: User, val: boolean) => {
  const message = val ? t('users.confirmEnableUser') : t('users.confirmDisableUser')
  ElMessageBox.confirm(message, t('common.tips'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  }).then(() => {
    handleStatusChange(row, val)
  }).catch(() => {
    // User cancelled, switch will revert automatically due to :model-value
  })
}

const handleStatusChange = async (row: User, val: boolean) => {
  try {
    if (val) {
      await activateUser(row.id)
      ElMessage.success(t('users.userEnabled'))
    } else {
      await suspendUser(row.id)
      ElMessage.success(t('users.userDisabled'))
    }
    loadUsers()
    loadStats()
  } catch (error) {
    console.error('Failed to change status:', error)
    ElMessage.error(t('users.operationFailed'))
  }
}

const handleResetPassword = (row: User) => {
  ElMessageBox.confirm(
    t('users.confirmResetPassword') + ` "${row.name}"?`,
    t('common.tips'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await resetUserPassword(row.id)
      ElMessage.success(`${t('users.passwordResetSuccess')}: ${res.temporary_password}`)
    } catch (error) {
      console.error('Failed to reset password:', error)
      ElMessage.error(t('users.passwordResetFailed'))
    }
  })
}

const handleDelete = (row: User) => {
  ElMessageBox.confirm(
    `${t('users.confirmDeleteUser')} "${row.name}"? ${t('users.deleteWarning')}`,
    t('common.warningConfirm'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success(t('users.deleteSuccess'))
      loadUsers()
      loadStats()
    } catch (error) {
      console.error('Failed to delete user:', error)
      ElMessage.error(t('users.deleteFailed'))
    }
  })
}

// Batch activate users
const handleBatchActivate = async () => {
  try {
    const data: BatchUpdateUserStatusRequest = {
      user_ids: selectedUsers.value.map(u => u.id),
      status: 1
    }

    const res = await batchUpdateUserStatus(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('users.batchActivatePartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('users.batchActivateSuccess', { count: successCount }))
    }

    handleClearSelection()
    loadUsers()
    loadStats()
  } catch (error) {
    console.error('Failed to batch activate users:', error)
    ElMessage.error(t('users.batchActivateFailed'))
  }
}

// Batch suspend users
const handleBatchSuspend = async () => {
  try {
    const { value } = await ElMessageBox.prompt(
      t('users.suspendReasonRequired'),
      t('users.batchSuspend'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /^.{5,200}$/,
        inputErrorMessage: t('users.suspendReasonLengthError')
      }
    )

    const data: BatchUpdateUserStatusRequest = {
      user_ids: selectedUsers.value.map(u => u.id),
      status: 2,
      reason: value
    }

    const res = await batchUpdateUserStatus(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('users.batchSuspendPartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('users.batchSuspendSuccess', { count: successCount }))
    }

    handleClearSelection()
    loadUsers()
    loadStats()
  } catch {
    // User cancelled
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadUsers()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadUsers()
}

onMounted(() => {
  loadUsers()
  loadStats()
})
</script>

<style scoped>
.users-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  border-left: 4px solid;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-card.primary {
  border-left-color: #6366F1;
}

.stat-card.success {
  border-left-color: #10B981;
}

.stat-card.warning {
  border-left-color: #F59E0B;
}

.stat-card.danger {
  border-left-color: #EF4444;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
  font-weight: 500;
}

/* Filter Bar */
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
  width: 260px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

.filter-arrow {
  transition: transform 0.3s;
  margin-left: 4px;
}

.filter-arrow.is-expanded {
  transform: rotate(180deg);
}

/* Advanced Filters */
.advanced-filters {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #F3F4F6;
}

.advanced-filters :deep(.el-input__wrapper) {
  border-radius: 12px;
}

/* Table */
.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Bulk Actions */
.bulk-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  margin-bottom: 16px;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.selected-count {
  font-size: 14px;
  color: #6366F1;
  font-weight: 600;
}

/* Table row hover */
:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* User Cell */
.user-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  font-weight: 600;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.user-email {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 2px 0;
}

.user-phone {
  font-size: 12px;
  color: #9CA3AF;
  margin: 0;
  font-family: 'Fira Code', monospace;
}

/* Time Text */
.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Dialog */
:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

/* Descriptions */
:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

/* Enhanced Values */
.enhanced-value {
  font-weight: 600;
  color: #1E1B4B;
  font-family: 'Fira Code', monospace;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .stat-card {
    border-radius: 14px;
  }
}
</style>
