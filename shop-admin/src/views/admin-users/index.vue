<template>
  <div class="admin-users-page">
    <el-tabs v-model="activeTab" class="user-tabs">
      <el-tab-pane label="管理员" name="admin" v-if="isPlatformAdmin" />
      <el-tab-pane label="顾客" name="customer" />
    </el-tabs>

    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名/邮箱/手机号"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="filterType"
            placeholder="用户类型"
            clearable
            class="filter-select"
            v-if="activeTab === 'admin' && isPlatformAdmin"
            @change="handleSearch"
          >
            <el-option label="全部" :value="0" />
            <el-option label="商家管理员" :value="2" />
            <el-option label="商家子账号" :value="3" />
          </el-select>
          <el-select
            v-model="filterStatus"
            placeholder="状态"
            clearable
            class="filter-select"
            v-if="activeTab === 'admin'"
            @change="handleSearch"
          >
            <el-option label="全部" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </div>
        <div class="filter-right">
          <el-button type="primary" @click="handleAdd" v-if="activeTab === 'admin' && isPlatformAdmin">
            <el-icon><Plus /></el-icon>新增管理员
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card class="table-card" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe @row-click="handleRowClick">
        <el-table-column label="用户信息" min-width="250">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="44" :src="row.avatar" class="user-avatar">
                {{ getAvatarText(row) }}
              </el-avatar>
              <div class="user-details">
                <p class="user-name">{{ row.real_name || row.name || '-' }}</p>
                <p class="user-email">{{ row.email || '-' }}</p>
                <p class="user-phone">{{ row.mobile || row.phone || '-' }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="130" align="center" v-if="activeTab === 'admin'">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ row.type_text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="100" align="center" v-if="activeTab === 'customer'">
          <template #default>
            <el-tag type="info" size="small">普通</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="订单数" width="100" align="center" v-if="activeTab === 'customer'">
          <template #default="{ row }">
            <span class="order-count">{{ row.order_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="累计消费" width="120" align="right" v-if="activeTab === 'customer'">
          <template #default="{ row }">
            <span class="total-spent">¥{{ formatPrice(row.total_spent) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              @change="(val: number) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" v-if="activeTab === 'admin'">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="handleView(row)">
              详情
            </el-button>
            <el-button type="primary" link size="small" @click.stop="handleAssignRoles(row)">
              角色
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button type="primary" link size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="resetPassword">重置密码</el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: #EF4444;">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right" v-if="activeTab === 'customer'">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="handleViewCustomer(row)">
              详情
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click.stop="handleEditCustomer(row)"
            >
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="formData.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="初始密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入初始密码" show-password />
        </el-form-item>
        <el-form-item label="用户类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择用户类型" style="width: 100%">
            <el-option label="商家管理员" :value="2" />
            <el-option label="商家子账号" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="租户ID" prop="tenant_id" v-if="isPlatformAdmin">
          <el-input-number v-model="formData.tenant_id" :min="0" placeholder="创建商家管理员时使用" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="用户详情" width="600px">
      <el-descriptions :column="2" border v-if="currentRow">
        <el-descriptions-item label="用户ID">{{ currentRow.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentRow.real_name || currentRow.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentRow.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentRow.mobile || currentRow.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="类型" v-if="currentRow.type_text">
          <el-tag :type="getTypeTagType(currentRow.type)">{{ currentRow.type_text }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentRow.status === 1 ? 'success' : 'danger'">
            {{ currentRow.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentRow.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间" v-if="currentRow.updated_at">{{ currentRow.updated_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Role Assign Dialog -->
    <RoleAssignDialog
      v-model:visible="roleDialogVisible"
      :admin-user="currentAdminUser"
      @assigned="loadData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import RoleAssignDialog from './components/RoleAssignDialog.vue'
import {
  getAdminUsers,
  disableAdminUser,
  enableAdminUser,
  createAdminUser,
  deleteAdminUser,
  resetAdminPassword,
  getCustomerList,
  type AdminUser,
  type AdminUserDetail,
  type CreateAdminUserParams
} from '@/api/admin-user'

const router = useRouter()
const userStore = useUserStore()

const isPlatformAdmin = computed(() => userStore.userInfo?.type === 1)

const activeTab = ref('admin')
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref(0)
const filterStatus = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref<any[]>([])

const dialogVisible = ref(false)
const detailVisible = ref(false)
const roleDialogVisible = ref(false)
const submitLoading = ref(false)
const currentRow = ref<any>(null)
const currentAdminUser = ref<AdminUserDetail | null>(null)
const formRef = ref()

const formData = reactive<CreateAdminUserParams>({
  email: '',
  mobile: '',
  real_name: '',
  password: '',
  type: 2,
  tenant_id: undefined
})

const formRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入初始密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择用户类型', trigger: 'change' }]
}

const dialogTitle = computed(() => '新增管理员')

const getAvatarText = (row: any) => {
  const name = row.real_name || row.name || row.username || ''
  return name.charAt(0).toUpperCase()
}

const getTypeTagType = (type: number) => {
  const types: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: 'info'
  }
  return types[type] || 'info'
}

const formatPrice = (price: number | undefined) => {
  if (!price) return '0.00'
  return (price / 100).toFixed(2)
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleAdd = () => {
  Object.assign(formData, {
    email: '',
    mobile: '',
    real_name: '',
    password: '',
    type: 2,
    tenant_id: undefined
  })
  dialogVisible.value = true
}

const handleView = (row: AdminUser) => {
  router.push(`/admin-users/${row.id}`)
}

const handleViewCustomer = (row: any) => {
  router.push(`/users/${row.id}`)
}

const handleRowClick = (row: any) => {
  if (activeTab.value === 'admin') {
    router.push(`/admin-users/${row.id}`)
  } else {
    router.push(`/users/${row.id}`)
  }
}

const handleEditCustomer = (row: any) => {
  router.push(`/users/${row.id}`)
}

const handleAssignRoles = (row: AdminUser) => {
  currentAdminUser.value = {
    ...row,
    roles: []
  } as AdminUserDetail
  roleDialogVisible.value = true
}

const handleCommand = (command: string, row: AdminUser) => {
  if (command === 'resetPassword') {
    handleResetPassword(row)
  } else if (command === 'delete') {
    handleDelete(row)
  }
}

const handleResetPassword = async (row: AdminUser) => {
  try {
    await ElMessageBox.confirm(`确认重置用户 "${row.real_name || row.email}" 的密码?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await resetAdminPassword(row.id)
    ElMessage.success(`密码重置成功，临时密码: ${res.temporary_password}`)
  } catch {
    // User cancelled
  }
}

const handleDelete = async (row: AdminUser) => {
  try {
    await ElMessageBox.confirm(`确认删除用户 "${row.real_name || row.email}"? 此操作不可恢复！`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAdminUser(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // User cancelled
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitLoading.value = true
      try {
        await createAdminUser(formData)
        ElMessage.success('新增成功')
        dialogVisible.value = false
        loadData()
      } catch (error) {
        console.error('Failed to create:', error)
        ElMessage.error('新增失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const handleStatusChange = async (row: any, val: number) => {
  const action = val === 1 ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(`确认${action}该用户?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    if (activeTab.value === 'admin') {
      if (val === 1) {
        await enableAdminUser(row.id)
      } else {
        await disableAdminUser(row.id)
      }
    }
    ElMessage.success(`${action}成功`)
  } catch {
    row.status = val === 1 ? 2 : 1
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const loadData = async () => {
  loading.value = true
  try {
    if (activeTab.value === 'admin') {
      const params: any = {
        page: currentPage.value,
        page_size: pageSize.value
      }
      if (searchQuery.value) params.keyword = searchQuery.value
      if (filterType.value) params.type = filterType.value
      if (filterStatus.value) params.status = filterStatus.value

      const res = await getAdminUsers(params)
      tableData.value = res.list || []
      total.value = res.total || 0
    } else {
      const params: any = {
        page: currentPage.value,
        page_size: pageSize.value
      }
      if (searchQuery.value) {
        params.name = searchQuery.value
        params.email = searchQuery.value
      }

      const res = await getCustomerList(params)
      tableData.value = res.list || []
      total.value = res.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(activeTab, () => {
  currentPage.value = 1
  searchQuery.value = ''
  filterType.value = 0
  filterStatus.value = 0
  tableData.value = []
  total.value = 0
  loadData()
})

watch(
  () => userStore.userInfo,
  (info) => {
    if (info && info.type !== 1) {
      activeTab.value = 'customer'
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.admin-users-page {
  padding: 0;
}

.user-tabs {
  margin-bottom: 16px;
}

.user-tabs :deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 500;
}

.user-tabs :deep(.el-tabs__item.is-active) {
  color: #6366F1;
}

.user-tabs :deep(.el-tabs__active-bar) {
  background-color: #6366F1;
}

.filter-card {
  margin-bottom: 16px;
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

.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  font-weight: 600;
  flex-shrink: 0;
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

.user-email,
.user-phone {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.order-count {
  font-size: 14px;
  color: #6B7280;
}

.total-spent {
  font-size: 16px;
  font-weight: 700;
  color: #10B981;
  font-family: 'Fira Sans', sans-serif;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Tags */
:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
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

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Descriptions */
:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
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
  .filter-select {
    width: 100%;
  }
}
</style>
