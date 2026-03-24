<template>
  <div class="admin-user-detail-page" v-loading="loading">
    <!-- Back Button -->
    <div class="page-header">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回用户列表
      </el-button>
      <h2 class="page-title">管理员详情</h2>
    </div>

    <!-- User Summary Card -->
    <el-card class="summary-card" shadow="never">
      <div class="user-summary">
        <div class="user-avatar-section">
          <el-avatar :size="80" :src="adminUser?.avatar" class="user-avatar">
            {{ getAvatarText() }}
          </el-avatar>
          <div class="user-info">
            <h3 class="user-name">{{ adminUser?.real_name || adminUser?.email || '-' }}</h3>
            <p class="user-email">{{ adminUser?.email || '-' }}</p>
            <p class="user-phone">{{ adminUser?.mobile || '-' }}</p>
          </div>
        </div>
        <div class="user-meta">
          <div class="meta-item">
            <el-tag :type="getTypeTagType(adminUser?.type)" size="large">
              {{ adminUser?.type_text || '未知类型' }}
            </el-tag>
          </div>
          <div class="meta-item">
            <el-tag :type="adminUser?.status === 1 ? 'success' : 'danger'" size="large">
              {{ adminUser?.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </div>
        </div>
        <div class="user-actions">
          <el-button type="primary" @click="handleEdit">编辑</el-button>
          <el-button @click="handleAssignRoles">分配角色</el-button>
          <el-button type="warning" @click="handleResetPassword">重置密码</el-button>
          <el-button
            :type="adminUser?.status === 1 ? 'danger' : 'success'"
            @click="handleStatusToggle"
          >
            {{ adminUser?.status === 1 ? '禁用' : '启用' }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Info Card -->
    <el-card class="info-card" shadow="never">
      <template #header>
        <span class="card-title">基本信息</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户ID">{{ adminUser?.id }}</el-descriptions-item>
        <el-descriptions-item label="真实姓名">{{ adminUser?.real_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ adminUser?.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ adminUser?.mobile || '-' }}</el-descriptions-item>
        <el-descriptions-item label="用户类型">
          <el-tag :type="getTypeTagType(adminUser?.type)">
            {{ adminUser?.type_text || '-' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="adminUser?.status === 1 ? 'success' : 'danger'">
            {{ adminUser?.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(adminUser?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDateTime(adminUser?.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Roles Card -->
    <el-card class="roles-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">已分配角色</span>
          <el-button type="primary" size="small" @click="handleAssignRoles">分配角色</el-button>
        </div>
      </template>
      <div class="role-list" v-if="adminUser?.roles?.length">
        <el-tag v-for="role in adminUser.roles" :key="role.id" class="role-tag">
          {{ role.name }}
        </el-tag>
      </div>
      <el-empty v-else description="暂未分配角色" />
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog v-model="editDialogVisible" title="编辑管理员" width="500px" destroy-on-close>
      <AdminUserForm
        ref="editFormRef"
        v-model="editFormData"
        :is-edit="true"
        :admin-user="adminUser"
      />
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit" :loading="editLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- Role Assign Dialog -->
    <RoleAssignDialog
      v-model:visible="roleDialogVisible"
      :admin-user="adminUser"
      @assigned="loadAdminUser"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import AdminUserForm from './components/AdminUserForm.vue'
import RoleAssignDialog from './components/RoleAssignDialog.vue'
import {
  getAdminUserDetail,
  updateAdminUser,
  disableAdminUser,
  enableAdminUser,
  resetAdminPassword,
  type AdminUserDetail,
  type UpdateAdminUserParams
} from '@/api/admin-user'

const route = useRoute()
const router = useRouter()

const adminUser = ref<AdminUserDetail | null>(null)
const loading = ref(false)

const editDialogVisible = ref(false)
const roleDialogVisible = ref(false)
const editLoading = ref(false)
const editFormRef = ref()
const editFormData = ref<UpdateAdminUserParams>({})

const userId = () => parseInt(route.params.id as string)

const loadAdminUser = async () => {
  loading.value = true
  try {
    const res = await getAdminUserDetail(userId())
    adminUser.value = res
  } catch (error) {
    console.error('Failed to load admin user:', error)
    // Mock data
    adminUser.value = {
      id: userId(),
      email: 'admin@example.com',
      mobile: '13800138000',
      real_name: '管理员',
      avatar: '',
      type: 2,
      type_text: '商家管理员',
      status: 1,
      created_at: '2024-01-15T10:00:00Z',
      updated_at: '2024-03-24T15:30:00Z',
      roles: [
        { id: 1, name: '超级管理员', code: 'super_admin' },
        { id: 2, name: '运营管理', code: 'operator' }
      ]
    }
  } finally {
    loading.value = false
  }
}

const getAvatarText = () => {
  const name = adminUser.value?.real_name || adminUser.value?.email || ''
  return name.charAt(0).toUpperCase()
}

const getTypeTagType = (type: number | undefined) => {
  const types: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: 'info'
  }
  return types[type || 0] || 'info'
}

const formatDateTime = (dateStr: string | undefined) => {
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

const goBack = () => {
  router.push('/admin-users')
}

const handleEdit = () => {
  editFormData.value = {
    email: adminUser.value?.email || '',
    mobile: adminUser.value?.mobile || '',
    real_name: adminUser.value?.real_name || ''
  }
  editDialogVisible.value = true
}

const confirmEdit = async () => {
  const valid = await editFormRef.value?.validate()
  if (!valid) return

  editLoading.value = true
  try {
    await updateAdminUser(userId(), editFormData.value)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    loadAdminUser()
  } catch (error) {
    console.error('Failed to update:', error)
    ElMessage.error('更新失败')
  } finally {
    editLoading.value = false
  }
}

const handleAssignRoles = () => {
  roleDialogVisible.value = true
}

const handleResetPassword = async () => {
  try {
    await ElMessageBox.confirm('确认重置该管理员的密码?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await resetAdminPassword(userId())
    ElMessage.success(`密码重置成功，临时密码: ${res.temporary_password}`)
  } catch {
    // User cancelled
  }
}

const handleStatusToggle = async () => {
  if (!adminUser.value) return

  const action = adminUser.value.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确认${action}该管理员?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    if (adminUser.value.status === 1) {
      await disableAdminUser(adminUser.value.id)
      ElMessage.success('已禁用')
    } else {
      await enableAdminUser(adminUser.value.id)
      ElMessage.success('已启用')
    }
    loadAdminUser()
  } catch {
    // User cancelled
  }
}

onMounted(() => {
  loadAdminUser()
})
</script>

<style scoped>
.admin-user-detail-page {
  padding: 0;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
}

.summary-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.user-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.user-avatar-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  font-size: 32px;
  font-weight: 600;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.user-email,
.user-phone {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

.user-meta {
  display: flex;
  gap: 12px;
}

.user-actions {
  display: flex;
  gap: 12px;
}

.info-card,
.roles-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.card-title {
  font-weight: 600;
  color: #1E1B4B;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

.role-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.role-tag {
  font-size: 14px;
}

@media (max-width: 768px) {
  .user-summary {
    flex-direction: column;
    align-items: flex-start;
  }

  .user-actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .user-actions .el-button {
    flex: 1;
  }
}
</style>