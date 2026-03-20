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
      <el-table :data="tableData" v-loading="loading" stripe>
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
          <template #default="{ row }">
            <el-tag type="warning" size="small" v-if="row.is_vip">VIP</el-tag>
            <el-tag type="info" size="small" v-else>普通</el-tag>
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
              @change="(val) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">
              详情
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
              v-if="activeTab === 'customer'"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import {
  getAdminUsers,
  disableAdminUser,
  enableAdminUser,
  registerTenantAdmin,
  getCustomerList
} from '@/api/admin-user'

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
const submitLoading = ref(false)
const currentRow = ref<any>(null)
const formRef = ref()

const formData = reactive({
  email: '',
  mobile: '',
  real_name: '',
  password: '',
  tenant_id: undefined as number | undefined
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
  ]
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
    tenant_id: undefined
  })
  dialogVisible.value = true
}

const handleView = (row: any) => {
  currentRow.value = row
  detailVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑顾客: ' + row.name)
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid: boolean) => {
    if (valid) {
      submitLoading.value = true
      registerTenantAdmin({
        email: formData.email,
        mobile: formData.mobile || undefined,
        real_name: formData.real_name,
        password: formData.password,
        tenant_id: formData.tenant_id
      })
        .then(() => {
          ElMessage.success('新增成功')
          dialogVisible.value = false
          loadData()
        })
        .finally(() => {
          submitLoading.value = false
        })
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

.filter-card {
  margin-bottom: 16px;
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

.filter-select {
  width: 140px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

.table-card {
  margin-bottom: 20px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  font-weight: 600;
  flex-shrink: 0;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: 500;
  color: #111827;
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
  font-weight: 600;
  color: #EF4444;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #E5E7EB;
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
