<template>
  <el-card class="tabs-card" shadow="never">
    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="基本信息" name="basic">
        <UserBasicInfo :user="user" @refresh="$emit('refresh')" />
      </el-tab-pane>
      <el-tab-pane label="收货地址" name="addresses">
        <UserAddressList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane label="订单记录" name="orders">
        <div class="coming-soon">
          <el-empty description="订单记录功能开发中" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="积分记录" name="points">
        <div class="coming-soon">
          <el-empty description="积分记录功能开发中" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="评价记录" name="reviews">
        <div class="coming-soon">
          <el-empty description="评价记录功能开发中" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="操作日志" name="logs">
        <UserOperationLog :user-id="user?.id" />
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import UserBasicInfo from './UserBasicInfo.vue'
import UserAddressList from './UserAddressList.vue'
import UserOperationLog from './UserOperationLog.vue'
import type { UserDetail } from '@/api/user'

defineProps<{
  user: UserDetail | null
}>()

defineEmits<{
  refresh: []
}>()

const activeTab = ref('basic')
</script>

<style scoped>
.tabs-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.detail-tabs :deep(.el-tabs__item) {
  font-size: 15px;
  font-weight: 500;
}

.detail-tabs :deep(.el-tabs__item.is-active) {
  color: #6366F1;
}

.detail-tabs :deep(.el-tabs__active-bar) {
  background-color: #6366F1;
}

.coming-soon {
  padding: 40px 0;
}
</style>