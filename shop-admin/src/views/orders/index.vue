<template>
  <div class="orders">
    <el-card>
      <template #header>
        <span>订单管理</span>
      </template>
      
      <el-table :data="orderList" v-loading="loading">
        <el-table-column prop="order_no" label="订单号" />
        <el-table-column prop="pay_amount" label="金额">
          <template #default="{ row }">
            ¥{{ (row.pay_amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button size="small" @click="handleDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const orderList = ref([])

onMounted(() => {
  loadOrders()
})

const loadOrders = async () => {
  loading.value = true
  try {
    orderList.value = []
  } finally {
    loading.value = false
  }
}

const handleDetail = (row: any) => {
  ElMessage.info('订单详情: ' + row.order_no)
}

const getStatusType = (status: number) => {
  const types: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: number) => {
  const texts: Record<number, string> = {
    0: '待支付',
    1: '已支付',
    2: '待发货',
    3: '已发货',
    4: '已完成',
    5: '已取消'
  }
  return texts[status] || '未知'
}
</script>
