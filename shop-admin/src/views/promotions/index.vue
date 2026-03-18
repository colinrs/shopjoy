<template>
  <div class="promotions">
    <el-card>
      <template #header>
        <span>促销管理</span>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="优惠券" name="coupon">
          <el-table :data="couponList" v-loading="loading">
            <el-table-column prop="name" label="优惠券名称" />
            <el-table-column prop="code" label="优惠券码" />
            <el-table-column prop="type" label="类型">
              <template #default="{ row }">
                {{ row.type === 0 ? '固定金额' : '百分比' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="促销活动" name="promotion">
          <p>促销活动功能开发中...</p>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const activeTab = ref('coupon')
const loading = ref(false)
const couponList = ref([])

onMounted(() => {
  loadCoupons()
})

const loadCoupons = async () => {
  loading.value = true
  try {
    couponList.value = []
  } finally {
    loading.value = false
  }
}
</script>
