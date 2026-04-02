<template>
  <div class="earn-rules-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.totalRules')"
          :value="ruleStats.total"
          :icon="Document"
          icon-color="primary"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.activated')"
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
            :placeholder="$t('points.searchRuleName')"
            clearable
            class="search-input"
            @clear="loadRules"
            @keyup.enter="loadRules"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchParams.status" :placeholder="$t('points.statusFilter')" clearable class="filter-select" @change="loadRules">
            <el-option :label="$t('points.all')" value="" />
            <el-option :label="$t('points.draft')" value="draft" />
            <el-option :label="$t('points.active')" value="active" />
            <el-option :label="$t('points.inactive')" value="inactive" />
          </el-select>
          <el-select v-model="searchParams.scenario" :placeholder="$t('points.scenarioFilter')" clearable class="filter-select" @change="loadRules">
            <el-option :label="$t('points.all')" value="" />
            <el-option :label="$t('points.orderPayment')" value="ORDER_PAYMENT" />
            <el-option :label="$t('points.signIn')" value="SIGN_IN" />
            <el-option :label="$t('points.productReview')" value="PRODUCT_REVIEW" />
            <el-option :label="$t('points.firstOrder')" value="FIRST_ORDER" />
          </el-select>
          <el-select v-model="searchParams.calculation_type" :placeholder="$t('points.calculationTypeFilter')" clearable class="filter-select" @change="loadRules">
            <el-option :label="$t('points.all')" value="" />
            <el-option :label="$t('points.fixed')" value="FIXED" />
            <el-option :label="$t('points.ratio')" value="RATIO" />
            <el-option :label="$t('points.tiered')" value="TIERED" />
          </el-select>
        </div>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ $t('points.createRule') }}
        </el-button>
      </div>
    </el-card>

    <!-- Rules Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="ruleList" v-loading="loading" stripe>
        <el-table-column :label="$t('points.ruleName')" min-width="280">
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
        <el-table-column :label="$t('points.scenario')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getScenarioTagType(row.scenario)" effect="light" size="small">
              {{ getScenarioText(row.scenario) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('points.calculationType')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getCalcTagType(row.calculation_type)" effect="plain" size="small">
              {{ getCalcText(row.calculation_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('points.condition')" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.condition_type === 'NONE'" type="info" size="small">{{ $t('points.unconditional') }}</el-tag>
            <el-tag v-else type="warning" size="small">{{ getConditionText(row.condition_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('points.expiration')" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.expiration_months === 0" class="no-expire">{{ $t('points.noExpire') }}</span>
            <span v-else class="expire-months">{{ $t('points.expireMonths', { months: row.expiration_months }) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('points.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" effect="light" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('points.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              {{ $t('points.edit') }}
            </el-button>
            <el-button
              v-if="row.status === 'draft' || row.status === 'inactive'"
              type="success"
              link
              size="small"
              @click="handleActivate(row)"
            >
              {{ $t('points.activate') }}
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              link
              size="small"
              @click="handleDeactivate(row)"
            >
              {{ $t('points.deactivate') }}
            </el-button>
            <el-button
              v-if="row.status !== 'active'"
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
            >
              {{ $t('points.delete') }}
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
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

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
    handleError(error, t('points.loadEarnRulesFailed'))
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
    'ORDER_PAYMENT': t('points.orderPayment'),
    'SIGN_IN': t('points.signIn'),
    'PRODUCT_REVIEW': t('points.productReview'),
    'FIRST_ORDER': t('points.firstOrder')
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
    'FIXED': t('points.fixed'),
    'RATIO': t('points.ratio'),
    'TIERED': t('points.tiered')
  }
  return texts[type] || type
}

const getConditionText = (type: string) => {
  const texts: Record<string, string> = {
    'NEW_USER': t('points.newUser'),
    'FIRST_ORDER': t('points.firstPurchase'),
    'SPECIFIC_PRODUCTS': t('points.specificProducts'),
    'MIN_AMOUNT': t('points.minAmount')
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
    'draft': t('points.draft'),
    'active': t('points.active'),
    'inactive': t('points.inactive')
  }
  return texts[status] || status
}

const getCalculationPreview = (rule: EarnRule) => {
  if (rule.calculation_type === 'FIXED') {
    return `${t('points.fixed')} ${rule.fixed_points} ${t('points.points')}`
  } else if (rule.calculation_type === 'RATIO') {
    return `${rule.ratio} ${t('points.pointsPerDollar')}`
  } else if (rule.calculation_type === 'TIERED' && rule.tiers) {
    return `${t('points.tiered')}: ${rule.tiers.length} tiers`
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
    ElMessage.success(t('points.activateSuccess'))
    loadRules()
  } catch (error) {
    handleError(error, t('points.activateRuleFailed'))
  }
}

const handleDeactivate = async (row: EarnRule) => {
  try {
    await deactivateEarnRule(row.id)
    ElMessage.success(t('points.deactivateSuccess'))
    loadRules()
  } catch (error) {
    handleError(error, t('points.deactivateRuleFailed'))
  }
}

const handleDelete = async (row: EarnRule) => {
  try {
    await ElMessageBox.confirm(t('points.deleteConfirm', { name: row.name }), t('common.warning'), {
      confirmButtonText: t('points.confirmDelete'),
      cancelButtonText: t('points.cancelDelete'),
      type: 'warning'
    })
    await deleteEarnRule(row.id)
    ElMessage.success(t('points.deleteSuccess'))
    loadRules()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('points.deleteRuleFailed'))
    }
  }
}

const handleSave = async (data: CreateEarnRuleParams) => {
  saveLoading.value = true
  try {
    if (currentRule.value) {
      await updateEarnRule({ id: currentRule.value.id, ...data })
      ElMessage.success(t('points.updateRuleSuccess'))
    } else {
      await createEarnRule(data)
      ElMessage.success(t('points.createRuleSuccess'))
    }
    formDialogVisible.value = false
    loadRules()
  } catch (error) {
    handleError(error, t('points.saveRuleFailed'))
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