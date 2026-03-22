<template>
  <div class="users-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card primary">
          <p class="stat-value">{{ stats.total }}</p>
          <p class="stat-label">用户总数</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card success">
          <p class="stat-value">{{ stats.active }}</p>
          <p class="stat-label">活跃用户</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card warning">
          <p class="stat-value">{{ stats.newToday }}</p>
          <p class="stat-label">今日新增</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card info">
          <p class="stat-value">{{ stats.vip }}</p>
          <p class="stat-label">VIP用户</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card danger">
          <p class="stat-value">{{ stats.blocked }}</p>
          <p class="stat-label">已禁用</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card purple">
          <p class="stat-value">{{ stats.online }}</p>
          <p class="stat-label">当前在线</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名/邮箱/手机号"
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="账号状态" clearable class="filter-select">
            <el-option label="全部" value="" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
          <el-select v-model="roleFilter" placeholder="用户角色" clearable class="filter-select">
            <el-option label="普通用户" value="user" />
            <el-option label="VIP用户" value="vip" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增用户
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Users Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="userList" v-loading="loading" stripe>
        <el-table-column type="selection" width="50" />
        <el-table-column label="用户信息" min-width="250">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="48" :src="row.avatar" class="user-avatar">
                {{ row.name.charAt(0) }}
              </el-avatar>
              <div class="user-details">
                <p class="user-name">
                  {{ row.name }}
                  <el-tag v-if="row.is_vip" size="small" type="warning" effect="dark">VIP</el-tag>
                  <el-tag v-if="row.is_online" size="small" type="success" effect="plain">在线</el-tag>
                </p>
                <p class="user-email">{{ row.email }}</p>
                <p class="user-phone">{{ row.phone }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="消费统计" width="180" align="right">
          <template #default="{ row }">
            <div class="stats-cell">
              <p class="order-count">{{ row.order_count }} 单</p>
              <p class="total-spent">¥{{ (row.total_spent / 100).toFixed(2) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="150">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="(val) => handleStatusChange(row, val)"
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
                  <el-dropdown-item @click="handleOrders(row)">查看订单</el-dropdown-item>
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

    <!-- User Detail Dialog -->
    <el-dialog v-model="detailDialogVisible" title="用户详情" width="700px">
      <el-descriptions :column="2" border v-if="currentUser">
        <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentUser.name }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentUser.email }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentUser.phone }}</el-descriptions-item>
        <el-descriptions-item label="角色">{{ getRoleText(currentUser.role) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentUser.status === 1 ? 'success' : 'danger'">
            {{ currentUser.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ currentUser.created_at }}</el-descriptions-item>
        <el-descriptions-item label="最后登录">{{ currentUser.last_login || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单数量">{{ currentUser.order_count }}</el-descriptions-item>
        <el-descriptions-item label="累计消费">¥{{ (currentUser.total_spent / 100).toFixed(2) }}</el-descriptions-item>
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
import { Search, Download, Plus, ArrowDown } from '@element-plus/icons-vue'

const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const roleFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(100)
const detailDialogVisible = ref(false)
const currentUser = ref<any>(null)

const stats = ref({
  total: 3256,
  active: 2847,
  newToday: 23,
  vip: 156,
  blocked: 12,
  online: 89
})

const userList = ref([
  {
    id: 1,
    name: '张三',
    email: 'zhangsan@example.com',
    phone: '13800138001',
    avatar: '',
    role: 'user',
    status: 1,
    is_vip: true,
    is_online: true,
    order_count: 23,
    total_spent: 1258000,
    created_at: '2024-01-15',
    last_login: '2024-03-18 14:30'
  },
  {
    id: 2,
    name: '李四',
    email: 'lisi@example.com',
    phone: '13900139002',
    avatar: '',
    role: 'vip',
    status: 1,
    is_vip: true,
    is_online: false,
    order_count: 56,
    total_spent: 2890000,
    created_at: '2024-01-10',
    last_login: '2024-03-17 20:15'
  },
  {
    id: 3,
    name: '王五',
    email: 'wangwu@example.com',
    phone: '13700137003',
    avatar: '',
    role: 'user',
    status: 1,
    is_vip: false,
    is_online: true,
    order_count: 8,
    total_spent: 234000,
    created_at: '2024-02-01',
    last_login: '2024-03-18 10:20'
  },
  {
    id: 4,
    name: '赵六',
    email: 'zhaoliu@example.com',
    phone: '13600136004',
    avatar: '',
    role: 'admin',
    status: 1,
    is_vip: false,
    is_online: true,
    order_count: 12,
    total_spent: 567000,
    created_at: '2023-12-20',
    last_login: '2024-03-18 09:00'
  },
  {
    id: 5,
    name: '钱七',
    email: 'qianqi@example.com',
    phone: '13500135005',
    avatar: '',
    role: 'user',
    status: 0,
    is_vip: false,
    is_online: false,
    order_count: 3,
    total_spent: 89000,
    created_at: '2024-03-01',
    last_login: '2024-03-10 18:45'
  }
])

const getRoleType = (role: string) => {
  const types: Record<string, string> = {
    'user': 'info',
    'vip': 'warning',
    'admin': 'danger'
  }
  return types[role] || 'info'
}

const getRoleText = (role: string) => {
  const texts: Record<string, string> = {
    'user': '普通用户',
    'vip': 'VIP用户',
    'admin': '管理员'
  }
  return texts[role] || role
}

const handleExport = () => {
  ElMessage.success('用户数据导出成功')
}

const handleAdd = () => {
  ElMessage.info('新增用户功能')
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑用户: ' + row.name)
}

const handleDetail = (row: any) => {
  currentUser.value = row
  detailDialogVisible.value = true
}

const handleOrders = (row: any) => {
  ElMessage.info('查看用户订单: ' + row.name)
}

const handleResetPassword = (row: any) => {
  ElMessageBox.confirm(`确认重置用户 "${row.name}" 的密码?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('密码重置成功，新密码已发送至用户邮箱')
  })
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确认删除用户 "${row.name}"? 此操作不可恢复！`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(() => {
    ElMessage.success('删除成功')
  })
}

const handleStatusChange = (row: any, val: number) => {
  const statusText = val === 1 ? '启用' : '禁用'
  ElMessage.success(`用户已${statusText}`)
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

onMounted(() => {
  // Load users
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

.stat-card.info {
  border-left-color: #3B82F6;
}

.stat-card.danger {
  border-left-color: #EF4444;
}

.stat-card.purple {
  border-left-color: #8B5CF6;
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
  display: flex;
  align-items: center;
  gap: 8px;
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

/* Stats Cell */
.stats-cell {
  text-align: right;
}

.order-count {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 4px 0;
}

.total-spent {
  font-size: 16px;
  font-weight: 700;
  color: #10B981;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
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
