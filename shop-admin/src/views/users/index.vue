<template>
  <div class="users-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card primary">
          <p class="stat-value">{{ stats.total }}</p>
          <p class="stat-label">用户总数</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card success">
          <p class="stat-value">{{ stats.active }}</p>
          <p class="stat-label">活跃用户</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card warning">
          <p class="stat-value">{{ stats.new_today }}</p>
          <p class="stat-label">今日新增</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card danger">
          <p class="stat-value">{{ stats.suspended }}</p>
          <p class="stat-label">已禁用</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名/邮箱"
            class="search-input"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="账号状态" clearable class="filter-select" @change="handleSearch">
            <el-option label="全部" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="已禁用" :value="2" />
          </el-select>
        </div>
        <div class="filter-right">
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Users Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="userList" v-loading="loading" stripe>
        <el-table-column label="用户信息" min-width="250">
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
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_login" label="最后登录" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ row.last_login || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.status === 1"
              @change="(val: boolean) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" link size="small" @click="handleDetail(row)">
              详情
            </el-button>
            <el-dropdown trigger="click">
              <el-button type="primary" link size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleResetPassword(row)">重置密码</el-dropdown-item>
                  <el-dropdown-item divided type="danger" @click="handleDelete(row)">
                    删除账号
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

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
    <el-dialog v-model="editDialogVisible" title="编辑用户" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="头像">
          <el-input v-model="editForm.avatar" placeholder="头像URL（选填）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit" :loading="editLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- User Detail Dialog -->
    <el-dialog v-model="detailDialogVisible" title="用户详情" width="600px">
      <el-descriptions :column="2" border v-if="currentUser">
        <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentUser.name }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentUser.email }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentUser.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentUser.status === 1 ? 'success' : 'danger'">
            {{ currentUser.status === 1 ? '正常' : '已禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ currentUser.created_at }}</el-descriptions-item>
        <el-descriptions-item label="最后登录">{{ currentUser.last_login || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, ArrowDown } from '@element-plus/icons-vue'
import {
  getUserList,
  getUserStats,
  updateUser,
  suspendUser,
  activateUser,
  deleteUser,
  resetUserPassword,
  type User,
  type UserStats
} from '@/api/user'

const loading = ref(false)
const editLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentUser = ref<User | null>(null)

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
    const res = await getUserList({
      page: currentPage.value,
      page_size: pageSize.value,
      name: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    })
    userList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to load users:', error)
    ElMessage.error('加载用户列表失败')
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
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
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
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    loadUsers()
  } catch (error) {
    console.error('Failed to update user:', error)
    ElMessage.error('更新失败')
  } finally {
    editLoading.value = false
  }
}

const handleDetail = (row: User) => {
  currentUser.value = row
  detailDialogVisible.value = true
}

const handleStatusChange = async (row: User, val: boolean) => {
  try {
    if (val) {
      await activateUser(row.id)
      ElMessage.success('用户已启用')
    } else {
      await suspendUser(row.id)
      ElMessage.success('用户已禁用')
    }
    loadUsers()
    loadStats()
  } catch (error) {
    console.error('Failed to change status:', error)
    ElMessage.error('操作失败')
  }
}

const handleResetPassword = (row: User) => {
  ElMessageBox.confirm(`确认重置用户 "${row.name}" 的密码?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await resetUserPassword(row.id)
      ElMessage.success(`密码重置成功，临时密码: ${res.temporary_password}`)
    } catch (error) {
      console.error('Failed to reset password:', error)
      ElMessage.error('密码重置失败')
    }
  })
}

const handleDelete = (row: User) => {
  ElMessageBox.confirm(`确认删除用户 "${row.name}"? 此操作不可恢复！`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success('删除成功')
      loadUsers()
      loadStats()
    } catch (error) {
      console.error('Failed to delete user:', error)
      ElMessage.error('删除失败')
    }
  })
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

/* Table */
.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
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