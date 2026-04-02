<template>
  <el-card
    class="tabs-card"
    shadow="never"
  >
    <el-tabs
      v-model="activeTab"
      class="detail-tabs"
    >
      <el-tab-pane
        :label="$t('users.basicInfo')"
        name="basic"
      >
        <UserBasicInfo
          :user="user"
          @refresh="$emit('refresh')"
        />
      </el-tab-pane>
      <el-tab-pane
        :label="$t('users.addresses')"
        name="addresses"
      >
        <UserAddressList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane
        :label="$t('users.orderRecords')"
        name="orders"
      >
        <div class="coming-soon">
          <el-empty :description="$t('users.comingSoon')" />
        </div>
      </el-tab-pane>
      <el-tab-pane
        :label="$t('users.pointsRecords')"
        name="points"
      >
        <div class="coming-soon">
          <el-empty :description="$t('users.comingSoon')" />
        </div>
      </el-tab-pane>
      <el-tab-pane
        :label="$t('users.reviewRecords')"
        name="reviews"
      >
        <div class="coming-soon">
          <el-empty :description="$t('users.comingSoon')" />
        </div>
      </el-tab-pane>
      <el-tab-pane
        :label="$t('users.operationLogs')"
        name="logs"
      >
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
