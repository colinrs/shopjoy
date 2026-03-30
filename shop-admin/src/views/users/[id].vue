<template>
  <div class="user-detail-page">
    <!-- Back Button -->
    <div class="page-header">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        {{ $t('users.backToUserList') }}
      </el-button>
      <h2 class="page-title">{{ $t('users.userDetail') }}</h2>
    </div>

    <!-- User Summary Card -->
    <el-card class="summary-card" shadow="never" v-loading="loading">
      <div class="user-summary">
        <div class="user-avatar-section">
          <el-avatar :size="80" :src="user?.avatar" class="user-avatar">
            {{ user?.name ? user.name.charAt(0) : 'U' }}
          </el-avatar>
          <div class="user-info">
            <h3 class="user-name">{{ user?.name || '-' }}</h3>
            <p class="user-email">{{ user?.email || '-' }}</p>
            <p class="user-phone">{{ user?.phone || '-' }}</p>
          </div>
        </div>
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-value">{{ user?.order_count || 0 }}</span>
            <span class="stat-label">{{ $t('users.orderCount') }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ user?.points_balance?.toLocaleString() || 0 }}</span>
            <span class="stat-label">{{ $t('users.pointsBalance') }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">¥{{ formatPrice(user?.total_spent) }}</span>
            <span class="stat-label">{{ $t('users.totalSpent') }}</span>
          </div>
          <div class="stat-item">
            <el-tag :type="user?.status === 1 ? 'success' : 'danger'" size="large">
              {{ user?.status === 1 ? $t('users.enabled') : $t('users.disabled') }}
            </el-tag>
          </div>
        </div>
        <div class="user-actions">
          <el-button type="primary" @click="handleEdit">{{ $t('common.edit') }}</el-button>
          <el-button
            :type="user?.status === 1 ? 'danger' : 'success'"
            @click="handleStatusToggle"
          >
            {{ user?.status === 1 ? $t('users.disabled') : $t('users.enabled') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Tabs -->
    <UserDetailTabs :user="user" @refresh="loadUser" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { t } from '@/plugins/i18n'
import UserDetailTabs from './components/UserDetailTabs.vue'
import {
  getUserDetail,
  suspendUser,
  activateUser,
  type UserDetail
} from '@/api/user'

const route = useRoute()
const router = useRouter()

const user = ref<UserDetail | null>(null)
const loading = ref(false)

const userId = () => parseInt(route.params.id as string)

const loadUser = async () => {
  loading.value = true
  try {
    const res = await getUserDetail(userId())
    user.value = res
  } catch (error) {
    console.error('Failed to load user:', error)
    ElMessage.error(t('users.loadUserDetailFailed'))
  } finally {
    loading.value = false
  }
}

const formatPrice = (price: string | undefined) => {
  if (!price) return '0.00'
  return parseFloat(price).toFixed(2)
}

const goBack = () => {
  router.push('/users')
}

const handleEdit = () => {
  ElMessage.info(t('users.editFunctionComingSoon'))
}

const handleStatusToggle = async () => {
  if (!user.value) return

  const isDisable = user.value.status === 1
  const confirmMsg = isDisable 
    ? t('users.confirmDisableUser')
    : t('users.confirmEnableUser')
  
  try {
    await ElMessageBox.confirm(confirmMsg, t('common.tips'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    if (isDisable) {
      await suspendUser(user.value.id)
      ElMessage.success(t('users.disabled'))
    } else {
      await activateUser(user.value.id)
      ElMessage.success(t('users.enabled'))
    }
    loadUser()
  } catch {
    // User cancelled
  }
}

onMounted(() => {
  loadUser()
})
</script>

<style scoped>
.user-detail-page {
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

.user-stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  font-family: 'Fira Sans', sans-serif;
}

.stat-label {
  display: block;
  font-size: 13px;
  color: #6B7280;
  margin-top: 4px;
}

.user-actions {
  display: flex;
  gap: 12px;
}

@media (max-width: 768px) {
  .user-summary {
    flex-direction: column;
    align-items: flex-start;
  }

  .user-stats {
    width: 100%;
    justify-content: space-between;
  }

  .user-actions {
    width: 100%;
  }

  .user-actions .el-button {
    flex: 1;
  }
}
</style>
