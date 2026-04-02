<template>
  <div class="admin-user-detail-page" v-loading="loading">
    <!-- Back Button -->
    <div class="page-header">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        {{ $t('adminUsers.backToUserList') }}
      </el-button>
      <h2 class="page-title">{{ $t('adminUsers.adminDetail') }}</h2>
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
              {{ adminUser?.type_text || $t('adminUsers.unknownType') }}
            </el-tag>
          </div>
          <div class="meta-item">
            <el-tag :type="adminUser?.status === 1 ? 'success' : 'danger'" size="large">
              {{ adminUser?.status === 1 ? $t('adminUsers.enabled') : $t('adminUsers.disabled') }}
            </el-tag>
          </div>
        </div>
        <div class="user-actions">
          <el-button type="primary" @click="handleEdit">{{ $t('adminUsers.edit') }}</el-button>
          <el-button @click="handleAssignRoles">{{ $t('adminUsers.assignRoles') }}</el-button>
          <el-button type="warning" @click="handleResetPassword">{{ $t('adminUsers.resetPassword') }}</el-button>
          <el-button
            :type="adminUser?.status === 1 ? 'danger' : 'success'"
            @click="handleStatusToggle"
          >
            {{ adminUser?.status === 1 ? $t('adminUsers.disabled2') : $t('adminUsers.enabled2') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Info Card -->
    <el-card class="info-card" shadow="never">
      <template #header>
        <span class="card-title">{{ $t('adminUsers.basicInfo') }}</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('adminUsers.userId')">{{ adminUser?.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.realName')">{{ adminUser?.real_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.email')">{{ adminUser?.email || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.mobile')">{{ adminUser?.mobile || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.type')">
          <el-tag :type="getTypeTagType(adminUser?.type)">
            {{ adminUser?.type_text || '-' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.status')">
          <el-tag :type="adminUser?.status === 1 ? 'success' : 'danger'">
            {{ adminUser?.status === 1 ? $t('adminUsers.enabled') : $t('adminUsers.disabled') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.createdAt')">{{ formatDateTime(adminUser?.created_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('adminUsers.updatedAt')">{{ formatDateTime(adminUser?.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Roles Card -->
    <el-card class="roles-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ $t('adminUsers.assignedRoles') }}</span>
          <el-button type="primary" size="small" @click="handleAssignRoles">{{ $t('adminUsers.assignRoles') }}</el-button>
        </div>
      </template>
      <div class="role-list" v-if="adminUser?.roles?.length">
        <el-tag v-for="role in adminUser.roles" :key="role.id" class="role-tag">
          {{ role.name }}
        </el-tag>
      </div>
      <el-empty v-else :description="$t('adminUsers.noRolesAssigned')" />
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog v-model="editDialogVisible" :title="$t('adminUsers.editAdmin')" width="500px" destroy-on-close>
      <AdminUserForm
        ref="editFormRef"
        v-model="editFormData"
        :is-edit="true"
        :admin-user="adminUser"
      />
      <template #footer>
        <el-button @click="editDialogVisible = false">{{ $t('adminUsers.cancel') }}</el-button>
        <el-button type="primary" @click="confirmEdit" :loading="editLoading">{{ $t('adminUsers.confirm') }}</el-button>
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
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n()
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
    ElMessage.error(t('adminUsers.failedToLoadAdminDetail'))
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
  return date.toLocaleString(undefined, {
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
    ElMessage.success(t('adminUsers.updateSuccess'))
    editDialogVisible.value = false
    loadAdminUser()
  } catch (error) {
    console.error('Failed to update:', error)
    ElMessage.error(t('adminUsers.updateFailed'))
  } finally {
    editLoading.value = false
  }
}

const handleAssignRoles = () => {
  roleDialogVisible.value = true
}

const handleResetPassword = async () => {
  try {
    await ElMessageBox.confirm(t('adminUsers.confirmResetPassword'), t('adminUsers.prompt'), {
      confirmButtonText: t('adminUsers.confirm'),
      cancelButtonText: t('adminUsers.cancel'),
      type: 'warning'
    })
    const res = await resetAdminPassword(userId())
    ElMessage.success(t('adminUsers.passwordResetSuccess') + res.temporary_password)
  } catch {
    // User cancelled
  }
}

const handleStatusToggle = async () => {
  if (!adminUser.value) return

  const isDisable = adminUser.value.status === 1
  try {
    await ElMessageBox.confirm(
      isDisable ? t('adminUsers.confirmDisable') : t('adminUsers.confirmEnable'),
      t('adminUsers.prompt'),
      {
        confirmButtonText: t('adminUsers.confirm'),
        cancelButtonText: t('adminUsers.cancel'),
        type: 'warning'
      }
    )

    if (isDisable) {
      await disableAdminUser(adminUser.value.id)
      ElMessage.success(t('adminUsers.updateSuccess'))
    } else {
      await enableAdminUser(adminUser.value.id)
      ElMessage.success(t('adminUsers.updateSuccess'))
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
