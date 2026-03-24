<template>
  <div class="redeem-rules-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="规则总数"
          :value="ruleStats.total"
          :icon="Document"
          icon-color="primary"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="已激活"
          :value="ruleStats.active"
          :icon="CircleCheck"
          icon-color="success"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          title="已兑换次数"
          :value="ruleStats.total_redeemed"
          :icon="Present"
          icon-color="warning"
        />
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchParams.name"
            placeholder="搜索规则名称"
            clearable
            class="search-input"
            @clear="loadRules"
            @keyup.enter="loadRules"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchParams.status" placeholder="状态" clearable class="filter-select" @change="loadRules">
            <el-option label="全部" value="" />
            <el-option label="已激活" value="active" />
            <el-option label="已停用" value="inactive" />
          </el-select>
        </div>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建规则
        </el-button>
      </div>
    </el-card>

    <!-- Rules Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="ruleList" v-loading="loading" stripe>
        <el-table-column label="规则信息" min-width="250">
          <template #default="{ row }">
            <div class="rule-cell">
              <div class="rule-icon">
                <el-icon size="24"><Ticket /></el-icon>
              </div>
              <div class="rule-details">
                <p class="rule-name">{{ row.name }}</p>
                <p class="rule-desc">{{ row.description || '无描述' }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="关联优惠券" width="150" align="center">
          <template #default="{ row }">
            <span class="coupon-name">{{ row.coupon_name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所需积分" width="100" align="center">
          <template #default="{ row }">
            <span class="points-value">{{ row.points_required.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column label="库存" width="180" align="center">
          <template #default="{ row }">
            <div class="stock-cell">
              <el-progress
                :percentage="getStockPercentage(row)"
                :status="getStockStatus(row)"
                :stroke-width="6"
                :show-text="false"
              />
              <span class="stock-text">
                {{ row.used_stock }}/{{ row.total_stock > 0 ? row.total_stock : '∞' }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="每人限兑" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.per_user_limit > 0">{{ row.per_user_limit }}次</span>
            <span v-else class="no-limit">不限</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" effect="light" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              v-if="row.status === 'inactive'"
              type="success"
              link
              size="small"
              @click="handleActivate(row)"
            >
              激活
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              link
              size="small"
              @click="handleDeactivate(row)"
            >
              停用
            </el-button>
            <el-button
              v-if="row.status !== 'active'"
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <TablePagination
        :page="searchParams.page"
        :page-size="searchParams.page_size"
        :total="total"
        @change="handlePageChange"
      />
    </el-card>

    <!-- Rule Form Dialog -->
    <RedeemRuleForm
      v-model:visible="formDialogVisible"
      :rule="currentRule"
      :loading="saveLoading"
      @submit="handleSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, CircleCheck, Present, Ticket, Search, Plus } from '@element-plus/icons-vue'
import PointsStatsCard from '../components/PointsStatsCard.vue'
import TablePagination from '@/components/common/TablePagination.vue'
import RedeemRuleForm from './components/RedeemRuleForm.vue'
import {
  getRedeemRules,
  createRedeemRule,
  updateRedeemRule,
  deleteRedeemRule,
  activateRedeemRule,
  deactivateRedeemRule,
  type RedeemRule,
  type ListRedeemRulesParams,
  type CreateRedeemRuleParams
} from '@/api/points'

// State
const loading = ref(false)
const saveLoading = ref(false)
const formDialogVisible = ref(false)
const currentRule = ref<RedeemRule | null>(null)

const ruleList = ref<RedeemRule[]>([])
const total = ref(0)

const ruleStats = ref({
  total: 0,
  active: 0,
  total_redeemed: 0
})

const searchParams = reactive<ListRedeemRulesParams>({
  page: 1,
  page_size: 10,
  name: '',
  status: ''
})

// Load functions
const loadRules = async () => {
  loading.value = true
  try {
    const res = await getRedeemRules(searchParams)
    ruleList.value = res.data.list || []
    total.value = res.data.total || 0
    ruleStats.value = res.data.stats
  } catch (error) {
    console.error('Failed to load redeem rules:', error)
    // Mock data
    ruleList.value = [
      {
        id: 1,
        name: '$10 优惠券兑换',
        description: '使用500积分兑换$10优惠券',
        coupon_id: 1,
        coupon_name: 'SAVE10',
        points_required: 500,
        total_stock: 100,
        used_stock: 45,
        per_user_limit: 5,
        status: 'active',
        start_at: null,
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      },
      {
        id: 2,
        name: '$20 优惠券兑换',
        description: '使用1000积分兑换$20优惠券',
        coupon_id: 2,
        coupon_name: 'SAVE20',
        points_required: 1000,
        total_stock: 50,
        used_stock: 20,
        per_user_limit: 3,
        status: 'active',
        start_at: null,
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      },
      {
        id: 3,
        name: '免邮券兑换',
        description: '使用200积分兑换免邮券',
        coupon_id: 3,
        coupon_name: 'FREESHIP',
        points_required: 200,
        total_stock: 200,
        used_stock: 80,
        per_user_limit: 10,
        status: 'inactive',
        start_at: null,
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      }
    ]
    total.value = 3
    ruleStats.value = { total: 3, active: 2, total_redeemed: 145 }
  } finally {
    loading.value = false
  }
}

// Helper functions
const getStockPercentage = (rule: RedeemRule) => {
  if (rule.total_stock === 0) return 0
  return Math.round((rule.used_stock / rule.total_stock) * 100)
}

const getStockStatus = (rule: RedeemRule) => {
  const pct = getStockPercentage(rule)
  if (pct >= 90) return 'exception'
  if (pct >= 70) return 'warning'
  return ''
}

const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'inactive': 'warning'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': '已激活',
    'inactive': '已停用'
  }
  return texts[status] || status
}

// Handlers
const handlePageChange = (page: number, pageSize: number) => {
  searchParams.page = page
  searchParams.page_size = pageSize
  loadRules()
}

const handleCreate = () => {
  currentRule.value = null
  formDialogVisible.value = true
}

const handleEdit = (row: RedeemRule) => {
  currentRule.value = { ...row }
  formDialogVisible.value = true
}

const handleActivate = async (row: RedeemRule) => {
  try {
    await activateRedeemRule(row.id)
    ElMessage.success('激活成功')
    loadRules()
  } catch (error) {
    console.error('Failed to activate:', error)
  }
}

const handleDeactivate = async (row: RedeemRule) => {
  try {
    await deactivateRedeemRule(row.id)
    ElMessage.success('停用成功')
    loadRules()
  } catch (error) {
    console.error('Failed to deactivate:', error)
  }
}

const handleDelete = async (row: RedeemRule) => {
  try {
    await ElMessageBox.confirm(`确定要删除规则 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteRedeemRule(row.id)
    ElMessage.success('删除成功')
    loadRules()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete:', error)
    }
  }
}

const handleSave = async (data: CreateRedeemRuleParams) => {
  saveLoading.value = true
  try {
    if (currentRule.value) {
      await updateRedeemRule({ id: currentRule.value.id, ...data })
      ElMessage.success('更新成功')
    } else {
      await createRedeemRule(data)
      ElMessage.success('创建成功')
    }
    formDialogVisible.value = false
    loadRules()
  } catch (error) {
    console.error('Failed to save:', error)
  } finally {
    saveLoading.value = false
  }
}

// Initialize
onMounted(() => {
  loadRules()
})
</script>

<style scoped>
.redeem-rules-page {
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
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 200px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 120px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.rule-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.rule-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.rule-details {
  flex: 1;
  min-width: 0;
}

.rule-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.rule-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.coupon-name {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.points-value {
  font-weight: 600;
  color: #F59E0B;
  font-family: 'Fira Sans', sans-serif;
}

.stock-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stock-text {
  font-size: 12px;
  color: #6B7280;
}

.no-limit {
  color: #9CA3AF;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

/* Progress */
:deep(.el-progress-bar__inner) {
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>