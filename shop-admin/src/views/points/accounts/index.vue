<template>
  <div class="accounts-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="账户总数"
          :value="accountStats.total"
          :icon="User"
          icon-color="primary"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="总余额"
          :value="accountStats.total_balance"
          :icon="Star"
          icon-color="warning"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="活跃用户"
          :value="accountStats.active"
          :icon="CircleCheck"
          icon-color="success"
        />
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户ID或邮箱"
          clearable
          class="search-input"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">
          搜索
        </el-button>
      </div>
    </el-card>

    <!-- Accounts Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="accountList" v-loading="loading" stripe @row-click="handleRowClick">
        <el-table-column label="用户ID" width="100" align="center">
          <template #default="{ row }">
            <span class="user-id">U{{ row.user_id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="邮箱" min-width="200">
          <template #default="{ row }">
            <span class="email-text">{{ row.user_email || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="可用余额" width="120" align="right">
          <template #default="{ row }">
            <span class="balance-value">{{ row.balance.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column label="冻结积分" width="100" align="right">
          <template #default="{ row }">
            <span class="frozen-value">{{ row.frozen_balance.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column label="累计获得" width="120" align="right">
          <template #default="{ row }">
            <span class="earned-value">{{ row.total_earned.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column label="累计兑换" width="120" align="right">
          <template #default="{ row }">
            <span class="redeemed-value">{{ row.total_redeemed.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="viewDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <TablePagination
        :page="currentPage"
        :page-size="pageSize"
        :total="total"
        @change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Star, CircleCheck, Search } from '@element-plus/icons-vue'
import PointsStatsCard from '../components/PointsStatsCard.vue'
import TablePagination from '@/components/common/TablePagination.vue'
import { getPointsAccounts, type PointsAccount } from '@/api/points'

const router = useRouter()

// State
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const accountList = ref<PointsAccount[]>([])

const accountStats = ref({
  total: 0,
  total_balance: 0,
  active: 0
})

// Load functions
const loadAccounts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    // Parse search query
    if (searchQuery.value) {
      const query = searchQuery.value.trim()
      if (/^\d+$/.test(query)) {
        params.user_id = parseInt(query)
      } else if (query.includes('@')) {
        params.email = query
      }
    }

    const res = await getPointsAccounts(params)
    accountList.value = res.data.list || []
    total.value = res.data.total || 0
    accountStats.value = res.data.stats
  } catch (error) {
    console.error('Failed to load accounts:', error)
    // Mock data
    accountList.value = [
      {
        id: 1,
        user_id: 12345,
        user_email: 'user1@example.com',
        balance: 5000,
        frozen_balance: 0,
        total_earned: 10000,
        total_redeemed: 4500,
        total_expired: 500,
        created_at: '2026-01-15T08:00:00Z',
        updated_at: '2026-03-24T10:30:00Z'
      },
      {
        id: 2,
        user_id: 12346,
        user_email: 'user2@example.com',
        balance: 3500,
        frozen_balance: 0,
        total_earned: 8000,
        total_redeemed: 4500,
        total_expired: 0,
        created_at: '2026-02-01T10:00:00Z',
        updated_at: '2026-03-24T09:00:00Z'
      },
      {
        id: 3,
        user_id: 12347,
        user_email: 'user3@example.com',
        balance: 1200,
        frozen_balance: 500,
        total_earned: 5000,
        total_redeemed: 3300,
        total_expired: 0,
        created_at: '2026-02-15T14:00:00Z',
        updated_at: '2026-03-23T16:00:00Z'
      }
    ]
    total.value = 3
    accountStats.value = { total: 3, total_balance: 9700, active: 3 }
  } finally {
    loading.value = false
  }
}

// Handlers
const handleSearch = () => {
  currentPage.value = 1
  loadAccounts()
}

const handlePageChange = (page: number, size: number) => {
  currentPage.value = page
  pageSize.value = size
  loadAccounts()
}

const handleRowClick = (row: PointsAccount) => {
  router.push(`/points/accounts/${row.id}`)
}

const viewDetail = (row: PointsAccount) => {
  router.push(`/points/accounts/${row.id}`)
}

// Initialize
onMounted(() => {
  loadAccounts()
})
</script>

<style scoped>
.accounts-page {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  gap: 12px;
}

.search-input {
  flex: 1;
  max-width: 300px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.user-id {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.email-text {
  color: #6B7280;
}

.balance-value {
  font-weight: 600;
  color: #F59E0B;
  font-family: 'Fira Sans', sans-serif;
}

.frozen-value {
  color: #6B7280;
}

.earned-value {
  color: #10B981;
  font-family: 'Fira Sans', sans-serif;
}

.redeemed-value {
  color: #3B82F6;
  font-family: 'Fira Sans', sans-serif;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
  }

  .search-input {
    max-width: 100%;
  }
}
</style>