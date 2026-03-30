<template>
  <div class="basic-info">
    <el-descriptions :column="2" border>
      <el-descriptions-item :label="$t('users.userId')">{{ user?.id }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.username')">{{ user?.name || '-' }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.email')">{{ user?.email || '-' }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.mobile')">{{ user?.phone || '-' }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.status')">
        <el-tag :type="user?.status === 1 ? 'success' : 'danger'">
          {{ user?.status === 1 ? $t('users.enabled') : $t('users.disabled') }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.memberLevel')">
        <el-tag type="info">{{ $t('users.regularMember') }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.createdAt')">{{ formatDateTime(user?.created_at) }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.lastLogin')">{{ formatDateTime(user?.last_login) }}</el-descriptions-item>
    </el-descriptions>

    <div class="section-title">{{ $t('users.pointsInfo') }}</div>
    <el-descriptions :column="4" border>
      <el-descriptions-item :label="$t('users.availablePoints')">
        <span class="points-value">{{ user?.points_balance?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.frozenPoints')">
        <span class="frozen-value">{{ user?.frozen_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.totalEarnedPoints')">
        <span class="earned-value">{{ user?.total_earned_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.totalRedeemedPoints')">
        <span class="redeemed-value">{{ user?.total_redeemed_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
    </el-descriptions>

    <div class="section-title">{{ $t('users.consumptionStats') }}</div>
    <el-descriptions :column="3" border>
      <el-descriptions-item :label="$t('users.totalOrders')">{{ user?.order_count || 0 }}</el-descriptions-item>
      <el-descriptions-item :label="$t('users.totalSpent')">
        <span class="spent-value">¥{{ formatPrice(user?.total_spent) }}</span>
      </el-descriptions-item>
      <el-descriptions-item :label="$t('users.lastOrder')">{{ formatDateTime(user?.last_order_at) }}</el-descriptions-item>
    </el-descriptions>
  </div>
</template>

<script setup lang="ts">
import type { UserDetail } from '@/api/user'

defineProps<{
  user: UserDetail | null
}>()

defineEmits<{
  refresh: []
}>()

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

const formatPrice = (price: string | undefined) => {
  if (!price) return '0.00'
  return parseFloat(price).toFixed(2)
}
</script>

<style scoped>
.basic-info {
  padding: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 24px 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #F3F4F6;
}

.section-title:first-of-type {
  margin-top: 0;
}

:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-descriptions__label) {
  width: 100px;
  font-weight: 500;
}

.points-value {
  font-size: 18px;
  font-weight: 700;
  color: #F59E0B;
}

.frozen-value {
  font-size: 16px;
  font-weight: 600;
  color: #6B7280;
}

.earned-value {
  font-size: 16px;
  font-weight: 600;
  color: #10B981;
}

.redeemed-value {
  font-size: 16px;
  font-weight: 600;
  color: #3B82F6;
}

.spent-value {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}
</style>
