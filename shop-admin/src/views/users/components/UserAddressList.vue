<template>
  <div class="address-list">
    <div class="address-header">
      <span class="header-title">收货地址列表</span>
    </div>

    <el-table :data="addresses" v-loading="loading" stripe>
      <el-table-column label="收货人" width="120">
        <template #default="{ row }">
          <span class="recipient-name">{{ row.recipient_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="联系电话" width="140">
        <template #default="{ row }">
          <span class="phone-text">{{ row.phone }}</span>
        </template>
      </el-table-column>
      <el-table-column label="地址" min-width="300">
        <template #default="{ row }">
          <span class="address-text">
            {{ row.province }} {{ row.city }} {{ row.district }} {{ row.detail_address }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="默认" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
          <span v-else class="non-default">-</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrapper" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="loadAddresses"
        @current-change="loadAddresses"
      />
    </div>

    <el-empty v-if="!loading && addresses.length === 0" description="暂无收货地址" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getUserAddresses, type UserAddress } from '@/api/user'

const props = defineProps<{
  userId?: number
}>()

const loading = ref(false)
const addresses = ref<UserAddress[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const loadAddresses = async () => {
  if (!props.userId) return

  loading.value = true
  try {
    const res = await getUserAddresses(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    addresses.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to load addresses:', error)
    // Mock data
    addresses.value = [
      {
        id: 1,
        user_id: props.userId,
        recipient_name: '张三',
        phone: '13800138001',
        province: '广东省',
        city: '深圳市',
        district: '南山区',
        detail_address: '科技园南区xxx大厦A座1001',
        is_default: true,
        created_at: '2024-01-15T10:00:00Z'
      },
      {
        id: 2,
        user_id: props.userId,
        recipient_name: '张三',
        phone: '13800138002',
        province: '广东省',
        city: '广州市',
        district: '天河区',
        detail_address: 'xxx路xxx号',
        is_default: false,
        created_at: '2024-02-20T14:30:00Z'
      }
    ]
    total.value = 2
  } finally {
    loading.value = false
  }
}

const formatDateTime = (dateStr: string) => {
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

watch(() => props.userId, () => {
  loadAddresses()
}, { immediate: true })

onMounted(() => {
  loadAddresses()
})
</script>

<style scoped>
.address-list {
  padding: 0;
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
}

.recipient-name {
  font-weight: 500;
  color: #1E1B4B;
}

.phone-text {
  font-family: 'Fira Code', monospace;
  color: #6B7280;
}

.address-text {
  font-size: 14px;
  color: #374151;
  line-height: 1.5;
}

.non-default {
  color: #9CA3AF;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}
</style>