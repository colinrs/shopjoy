<template>
  <div class="earn-rules-page">
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
            <el-option label="草稿" value="draft" />
            <el-option label="已激活" value="active" />
            <el-option label="已停用" value="inactive" />
          </el-select>
          <el-select v-model="searchParams.scenario" placeholder="场景" clearable class="filter-select" @change="loadRules">
            <el-option label="全部" value="" />
            <el-option label="订单支付" value="ORDER_PAYMENT" />
            <el-option label="每日签到" value="SIGN_IN" />
            <el-option label="商品评价" value="PRODUCT_REVIEW" />
            <el-option label="首单奖励" value="FIRST_ORDER" />
          </el-select>
          <el-select v-model="searchParams.calculation_type" placeholder="计算类型" clearable class="filter-select" @change="loadRules">
            <el-option label="全部" value="" />
            <el-option label="固定积分" value="FIXED" />
            <el-option label="比例" value="RATIO" />
            <el-option label="阶梯" value="TIERED" />
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
        <el-table-column label="规则信息" min-width="280">
          <template #default="{ row }">
            <div class="rule-cell">
              <div class="rule-icon" :class="getScenarioClass(row.scenario)">
                <el-icon size="24"><Document /></el-icon>
              </div>
              <div class="rule-details">
                <p class="rule-name">{{ row.name }}</p>
                <p class="rule-desc">{{ getCalculationPreview(row) }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="场景" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getScenarioTagType(row.scenario)" effect="light" size="small">
              {{ getScenarioText(row.scenario) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="计算类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getCalcTagType(row.calculation_type)" effect="plain" size="small">
              {{ getCalcText(row.calculation_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="条件" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.condition_type === 'NONE'" type="info" size="small">无</el-tag>
            <el-tag v-else type="warning" size="small">{{ getConditionText(row.condition_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="过期" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.expiration_months === 0" class="no-expire">永不过期</span>
            <span v-else class="expire-months">{{ row.expiration_months }}个月</span>
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
              v-if="row.status === 'draft' || row.status === 'inactive'"
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
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        @change="handlePageChange"
      />
    </el-card>

    <!-- Rule Form Dialog -->
    <EarnRuleForm
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
import { Document, CircleCheck, Search, Plus } from '@element-plus/icons-vue'
import PointsStatsCard from '../components/PointsStatsCard.vue'
import TablePagination from '@/components/common/TablePagination.vue'
import EarnRuleForm from './components/EarnRuleForm.vue'
import {
  getEarnRules,
  createEarnRule,
  updateEarnRule,
  deleteEarnRule,
  activateEarnRule,
  deactivateEarnRule,
  type EarnRule,
  type CreateEarnRuleParams
} from '@/api/points'

// State
const loading = ref(false)
const saveLoading = ref(false)
const formDialogVisible = ref(false)
const currentRule = ref<EarnRule | null>(null)

const ruleList = ref<EarnRule[]>([])
const total = ref(0)

const ruleStats = ref({
  total: 0,
  active: 0
})

const searchParams = reactive({
  name: '',
  status: '',
  scenario: '',
  calculation_type: ''
})

const currentPage = ref(1)
const pageSize = ref(10)

// Load functions
const loadRules = async () => {
  loading.value = true
  try {
    const res = await getEarnRules({
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchParams
    })
    ruleList.value = res.list || []
    total.value = res.total || 0
    ruleStats.value = res.stats
  } catch (error) {
    console.error('Failed to load earn rules:', error)
    // Mock data
    ruleList.value = [
      {
        id: 1,
        name: '订单积分奖励',
        description: '每笔订单可获得积分奖励',
        scenario: 'ORDER_PAYMENT',
        calculation_type: 'TIERED',
        fixed_points: 0,
        ratio: '0',
        tiers: [
          { threshold: 10000, ratio: '1.0' },
          { threshold: 50000, ratio: '1.5' },
          { threshold: null, ratio: '2.0' }
        ],
        condition_type: 'NONE',
        condition_value: null,
        expiration_months: 12,
        status: 'active',
        priority: 10,
        start_at: '2026-04-01T00:00:00Z',
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      },
      {
        id: 2,
        name: '每日签到奖励',
        description: '每日签到可获得5积分',
        scenario: 'SIGN_IN',
        calculation_type: 'FIXED',
        fixed_points: 5,
        ratio: '0',
        tiers: null,
        condition_type: 'NONE',
        condition_value: null,
        expiration_months: 6,
        status: 'active',
        priority: 5,
        start_at: null,
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      },
      {
        id: 3,
        name: '商品评价奖励',
        description: '评价商品可获得10积分',
        scenario: 'PRODUCT_REVIEW',
        calculation_type: 'FIXED',
        fixed_points: 10,
        ratio: '0',
        tiers: null,
        condition_type: 'NONE',
        condition_value: null,
        expiration_months: 12,
        status: 'draft',
        priority: 3,
        start_at: null,
        end_at: null,
        created_at: '2026-03-24T10:00:00Z',
        updated_at: '2026-03-24T10:00:00Z'
      }
    ]
    total.value = 3
    ruleStats.value = { total: 3, active: 2 }
  } finally {
    loading.value = false
  }
}

// Helper functions
const getScenarioClass = (scenario: string) => {
  const classes: Record<string, string> = {
    'ORDER_PAYMENT': 'order',
    'SIGN_IN': 'signin',
    'PRODUCT_REVIEW': 'review',
    'FIRST_ORDER': 'first'
  }
  return classes[scenario] || ''
}

const getScenarioTagType = (scenario: string) => {
  const types: Record<string, string> = {
    'ORDER_PAYMENT': 'primary',
    'SIGN_IN': 'success',
    'PRODUCT_REVIEW': 'warning',
    'FIRST_ORDER': 'danger'
  }
  return types[scenario] || 'info'
}

const getScenarioText = (scenario: string) => {
  const texts: Record<string, string> = {
    'ORDER_PAYMENT': '订单支付',
    'SIGN_IN': '每日签到',
    'PRODUCT_REVIEW': '商品评价',
    'FIRST_ORDER': '首单奖励'
  }
  return texts[scenario] || scenario
}

const getCalcTagType = (type: string) => {
  const types: Record<string, string> = {
    'FIXED': 'success',
    'RATIO': 'primary',
    'TIERED': 'warning'
  }
  return types[type] || 'info'
}

const getCalcText = (type: string) => {
  const texts: Record<string, string> = {
    'FIXED': '固定',
    'RATIO': '比例',
    'TIERED': '阶梯'
  }
  return texts[type] || type
}

const getConditionText = (type: string) => {
  const texts: Record<string, string> = {
    'NEW_USER': '新用户',
    'FIRST_ORDER': '首单',
    'SPECIFIC_PRODUCTS': '指定商品',
    'MIN_AMOUNT': '最低金额'
  }
  return texts[type] || type
}

const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    'draft': 'info',
    'active': 'success',
    'inactive': 'warning'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'draft': '草稿',
    'active': '已激活',
    'inactive': '已停用'
  }
  return texts[status] || status
}

const getCalculationPreview = (rule: EarnRule) => {
  if (rule.calculation_type === 'FIXED') {
    return `固定 ${rule.fixed_points} 积分`
  } else if (rule.calculation_type === 'RATIO') {
    return `${rule.ratio} 积分/$1`
  } else if (rule.calculation_type === 'TIERED' && rule.tiers) {
    return `阶梯: ${rule.tiers.length} 档`
  }
  return '-'
}

// Handlers
const handlePageChange = () => {
  loadRules()
}

const handleCreate = () => {
  currentRule.value = null
  formDialogVisible.value = true
}

const handleEdit = (row: EarnRule) => {
  currentRule.value = { ...row }
  formDialogVisible.value = true
}

const handleActivate = async (row: EarnRule) => {
  try {
    await activateEarnRule(row.id)
    ElMessage.success('激活成功')
    loadRules()
  } catch (error) {
    console.error('Failed to activate:', error)
  }
}

const handleDeactivate = async (row: EarnRule) => {
  try {
    await deactivateEarnRule(row.id)
    ElMessage.success('停用成功')
    loadRules()
  } catch (error) {
    console.error('Failed to deactivate:', error)
  }
}

const handleDelete = async (row: EarnRule) => {
  try {
    await ElMessageBox.confirm(`确定要删除规则 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteEarnRule(row.id)
    ElMessage.success('删除成功')
    loadRules()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete:', error)
    }
  }
}

const handleSave = async (data: CreateEarnRuleParams) => {
  saveLoading.value = true
  try {
    if (currentRule.value) {
      await updateEarnRule({ id: currentRule.value.id, ...data })
      ElMessage.success('更新成功')
    } else {
      await createEarnRule(data)
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
.earn-rules-page {
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
  background: #F5F3FF;
  color: #6366F1;
}

.rule-icon.order {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.rule-icon.signin {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  color: white;
}

.rule-icon.review {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.rule-icon.first {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
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

.no-expire {
  font-size: 12px;
  color: #10B981;
}

.expire-months {
  font-size: 12px;
  color: #6B7280;
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

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
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