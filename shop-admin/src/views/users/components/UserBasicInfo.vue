<template>
  <div class="basic-info">
    <el-descriptions :column="2" border>
      <el-descriptions-item label="用户ID">{{ user?.id }}</el-descriptions-item>
      <el-descriptions-item label="用户名">{{ user?.name || '-' }}</el-descriptions-item>
      <el-descriptions-item label="邮箱">{{ user?.email || '-' }}</el-descriptions-item>
      <el-descriptions-item label="手机号">{{ user?.phone || '-' }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="user?.status === 1 ? 'success' : 'danger'">
          {{ user?.status === 1 ? '正常' : '已禁用' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="会员等级">
        <el-tag type="info">普通会员</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="注册时间">{{ formatDateTime(user?.created_at) }}</el-descriptions-item>
      <el-descriptions-item label="最后登录">{{ formatDateTime(user?.last_login) }}</el-descriptions-item>
    </el-descriptions>

    <div class="section-title">积分信息</div>
    <el-descriptions :column="4" border>
      <el-descriptions-item label="可用积分">
        <span class="points-value">{{ user?.points_balance?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="冻结积分">
        <span class="frozen-value">{{ user?.frozen_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="累计获得">
        <span class="earned-value">{{ user?.total_earned_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="累计兑换">
        <span class="redeemed-value">{{ user?.total_redeemed_points?.toLocaleString() || 0 }}</span>
      </el-descriptions-item>
    </el-descriptions>

    <div class="section-title">消费统计</div>
    <el-descriptions :column="3" border>
      <el-descriptions-item label="订单总数">{{ user?.order_count || 0 }} 单</el-descriptions-item>
      <el-descriptions-item label="累计消费">
        <span class="spent-value">¥{{ formatPrice(user?.total_spent) }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="最后下单">{{ formatDateTime(user?.last_order_at) }}</el-descriptions-item>
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