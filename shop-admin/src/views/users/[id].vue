<template>
  <div class="user-detail-page">
    <!-- Back Button -->
    <div class="page-header">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回用户列表
      </el-button>
      <h2 class="page-title">用户详情</h2>
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
            <span class="stat-label">订单数</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ user?.points_balance?.toLocaleString() || 0 }}</span>
            <span class="stat-label">积分余额</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">¥{{ formatPrice(user?.total_spent) }}</span>
            <span class="stat-label">累计消费</span>
          </div>
          <div class="stat-item">
            <el-tag :type="user?.status === 1 ? 'success' : 'danger'" size="large">
              {{ user?.status === 1 ? '正常' : '已禁用' }}
            </el-tag>
          </div>
        </div>
        <div class="user-actions">
          <el-button type="primary" @click="handleEdit">编辑</el-button>
          <el-button
            :type="user?.status === 1 ? 'danger' : 'success'"
            @click="handleStatusToggle"
          >
            {{ user?.status === 1 ? '禁用' : '启用' }}
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
    // Mock data for development
    user.value = {
      id: userId(),
      tenant_id: 1,
      email: 'user@example.com',
      phone: '13800138000',
      name: '测试用户',
      avatar: '',
      gender: 1,
      gender_text: '男',
      birthday: '1990-01-01',
      status: 1,
      status_text: '正常',
      review_count: 5,
      created_at: '2024-01-15T10:00:00Z',
      updated_at: '2024-03-24T15:30:00Z',
      last_login: '2024-03-24T15:30:00Z',
      points_balance: 5000,
      frozen_points: 0,
      total_earned_points: 10000,
      total_redeemed_points: 5000,
      order_count: 25,
      total_spent: '12580.50',
      last_order_at: '2024-03-20T08:00:00Z',
      default_address: null
    }
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
  ElMessage.info('编辑功能开发中')
}

const handleStatusToggle = async () => {
  if (!user.value) return

  const action = user.value.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确认${action}该用户?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    if (user.value.status === 1) {
      await suspendUser(user.value.id)
      ElMessage.success('已禁用')
    } else {
      await activateUser(user.value.id)
      ElMessage.success('已启用')
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