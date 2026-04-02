<template>
  <div class="redeem-rules-page">
    <!-- Stats Cards -->
    <el-row
      :gutter="16"
      class="stats-row"
    >
      <el-col
        :xs="12"
        :sm="6"
      >
        <PointsStatsCard
          :title="$t('points.totalRules')"
          :value="ruleStats.total"
          :icon="Document"
          icon-color="primary"
        />
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <PointsStatsCard
          :title="$t('points.activated')"
          :value="ruleStats.active"
          :icon="CircleCheck"
          icon-color="success"
        />
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <PointsStatsCard
          :title="$t('points.redeemedTimes')"
          :value="ruleStats.total_redeemed"
          :icon="Present"
          icon-color="warning"
        />
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card
      class="filter-card"
      shadow="never"
    >
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
          <el-select
            v-model="searchParams.status"
            :placeholder="$t('points.statusFilter')"
            clearable
            class="filter-select"
            @change="loadRules"
          >
            <el-option
              :label="$t('points.all')"
              value=""
            />
            <el-option
              :label="$t('points.active')"
              value="active"
            />
            <el-option
              :label="$t('points.inactive')"
              value="inactive"
            />
          </el-select>
        </div>
        <el-button
          type="primary"
          @click="handleCreate"
        >
          <el-icon><Plus /></el-icon>
          {{ $t('points.createRule') }}
        </el-button>
      </div>
    </el-card>

    <!-- Rules Table -->
    <el-card
      class="table-card"
      shadow="never"
    >
      <el-table
        v-loading="loading"
        :data="ruleList"
        stripe
      >
        <el-table-column
          :label="$t('points.ruleName')"
          min-width="250"
        >
          <template #default="{ row }">
            <div class="rule-cell">
              <div class="rule-icon">
                <el-icon size="24">
                  <Ticket />
                </el-icon>
              </div>
              <div class="rule-details">
                <p class="rule-name">
                  {{ row.name }}
                </p>
                <p class="rule-desc">
                  {{ row.description || $t('points.noDescription') }}
                </p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('points.coupon')"
          width="150"
          align="center"
        >
          <template #default="{ row }">
            <span class="coupon-name">{{ row.coupon_name }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('points.pointsRequired')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <span class="points-value">{{ row.points_required.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('points.stock')"
          width="180"
          align="center"
        >
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
        <el-table-column
          :label="$t('points.limitPerUser')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <span v-if="row.per_user_limit > 0">{{ row.per_user_limit }}{{ $t('points.timesUnit') }}</span>
            <span
              v-else
              class="no-limit"
            >{{ $t('points.unlimitedText') }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('points.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              effect="light"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('points.actions')"
          width="200"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              {{ $t('points.edit') }}
            </el-button>
            <el-button
              v-if="row.status === 'inactive'"
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
  type CreateRedeemRuleParams
} from '@/api/points'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

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

const searchParams = reactive({
  name: '',
  status: ''
})

const currentPage = ref(1)
const pageSize = ref(10)

// Load functions
const loadRules = async () => {
  loading.value = true
  try {
    const res = await getRedeemRules({
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchParams
    })
    ruleList.value = res.list || []
    total.value = res.total || 0
    ruleStats.value = res.stats
  } catch (error) {
    handleError(error, t('points.loadRedeemRulesFailed'))
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
    'active': t('points.active'),
    'inactive': t('points.inactive')
  }
  return texts[status] || status
}

// Handlers
const handlePageChange = () => {
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
    ElMessage.success(t('points.activateSuccess'))
    loadRules()
  } catch (error) {
    handleError(error, t('points.activateRuleFailed'))
  }
}

const handleDeactivate = async (row: RedeemRule) => {
  try {
    await deactivateRedeemRule(row.id)
    ElMessage.success(t('points.deactivateSuccess'))
    loadRules()
  } catch (error) {
    handleError(error, t('points.deactivateRuleFailed'))
  }
}

const handleDelete = async (row: RedeemRule) => {
  try {
    await ElMessageBox.confirm(t('points.deleteConfirm', { name: row.name }), t('common.warning'), {
      confirmButtonText: t('points.confirmDelete'),
      cancelButtonText: t('points.cancelDelete'),
      type: 'warning'
    })
    await deleteRedeemRule(row.id)
    ElMessage.success(t('points.deleteSuccess'))
    loadRules()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('points.deleteRuleFailed'))
    }
  }
}

const handleSave = async (data: CreateRedeemRuleParams) => {
  saveLoading.value = true
  try {
    if (currentRule.value) {
      await updateRedeemRule({ id: currentRule.value.id, ...data })
      ElMessage.success(t('points.updateRuleSuccess'))
    } else {
      await createRedeemRule(data)
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