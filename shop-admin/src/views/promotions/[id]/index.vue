<template>
  <div class="promotion-detail-page" v-loading="loading">
    <PageHeader
      :title="isEdit ? $t('promotions.editPromotion') : $t('promotions.createPromotion')"
      :subtitle="promotionForm.name || $t('promotions.configurePromotionInfo')"
      @back="handleBack"
    />

    <el-row :gutter="20">
      <el-col :xs="24" :lg="16">
        <el-card class="form-card" shadow="never">
          <el-form
            :model="promotionForm"
            :rules="rules"
            ref="formRef"
            label-width="120px"
            label-position="top"
          >
            <!-- Basic Info -->
            <div class="section-title">{{ $t('promotions.basicInfo') }}</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="$t('promotions.promotionName')" prop="name">
                  <el-input v-model="promotionForm.name" :placeholder="$t('promotions.enterPromotionName')" maxlength="100" show-word-limit />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item :label="$t('promotions.promotionTypeSelect')" prop="type">
                  <el-select v-model="promotionForm.type" :placeholder="$t('promotions.selectPromotionType')" style="width: 100%">
                    <el-option :label="$t('promotions.discount')" value="discount" />
                    <el-option :label="$t('promotions.flashSale')" value="flash_sale" />
                    <el-option :label="$t('promotions.bundle')" value="bundle" />
                    <el-option :label="$t('promotions.buyXGetY')" value="buy_x_get_y" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item :label="$t('promotions.promotionDescription')">
              <el-input
                v-model="promotionForm.description"
                type="textarea"
                :rows="3"
                :placeholder="$t('promotions.enterPromotionDesc')"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <!-- Discount Settings -->
            <div class="section-title">{{ $t('promotions.discountSettings') }}</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="8">
                <el-form-item :label="$t('promotions.discountType')" prop="discount_type">
                  <el-select v-model="promotionForm.discount_type" :placeholder="$t('promotions.selectDiscountType')" style="width: 100%">
                    <el-option :label="$t('promotions.fixedAmount')" value="fixed_amount" />
                    <el-option :label="$t('promotions.percentage')" value="percentage" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8">
                <el-form-item :label="promotionForm.discount_type === 'fixed_amount' ? $t('promotions.discountAmountLabel') : $t('promotions.discountRatioLabel')" prop="discount_value">
                  <el-input-number
                    v-model="promotionForm.discount_value_num"
                    :min="0"
                    :max="promotionForm.discount_type === 'percentage' ? 100 : 99999"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8">
                <el-form-item :label="$t('promotions.lowestConsume')">
                  <el-input-number
                    v-model="promotionForm.min_order_amount_num"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20" v-if="promotionForm.discount_type === 'percentage'">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="$t('promotions.maxDiscountAmount')">
                  <el-input-number
                    v-model="promotionForm.max_discount_num"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <!-- Time Range -->
            <div class="section-title">{{ $t('promotions.activityTime') }}</div>
            <el-form-item :label="$t('promotions.promotionPeriod')" prop="dateRange">
              <el-date-picker
                v-model="promotionForm.dateRange"
                type="datetimerange"
                :range-separator="$t('promotions.to')"
                :start-placeholder="$t('promotions.startPlaceholder')"
                :end-placeholder="$t('promotions.endPlaceholder')"
                value-format="YYYY-MM-DDTHH:mm:ss[Z]"
                style="width: 100%"
              />
            </el-form-item>

            <!-- Scope Settings -->
            <div class="section-title">{{ $t('promotions.scopeSettings') }}</div>
            <el-form-item :label="$t('promotions.activityScope')">
              <el-radio-group v-model="promotionForm.scope_type">
                <el-radio label="storewide">{{ $t('promotions.storewide') }}</el-radio>
                <el-radio label="products">{{ $t('promotions.specificProductsScope') }}</el-radio>
                <el-radio label="categories">{{ $t('promotions.specificCategories') }}</el-radio>
                <el-radio label="brands">{{ $t('promotions.specificBrands') }}</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item v-if="promotionForm.scope_type === 'products'" :label="$t('promotions.selectProducts')">
              <el-select
                v-model="promotionForm.product_ids"
                multiple
                filterable
                remote
                reserve-keyword
                :placeholder="$t('promotions.searchProducts')"
                :remote-method="searchProducts"
                style="width: 100%"
              >
                <el-option
                  v-for="item in productOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item v-if="promotionForm.scope_type === 'categories'" :label="$t('promotions.selectCategories')">
              <el-select v-model="promotionForm.category_ids" multiple :placeholder="$t('promotions.selectCategoryPlaceholder')" style="width: 100%">
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <!-- Usage Limits -->
            <div class="section-title">{{ $t('promotions.usageLimits') }}</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="$t('promotions.totalUsageLimit')">
                  <el-input-number v-model="promotionForm.usage_limit" :min="0" style="width: 100%" />
                  <div class="form-tip">{{ $t('promotions.zeroUnlimitedUsage') }}</div>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item :label="$t('promotions.perUserUsageLimit')">
                  <el-input-number v-model="promotionForm.per_user_limit" :min="0" style="width: 100%" />
                  <div class="form-tip">{{ $t('promotions.zeroUnlimitedUsage') }}</div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <!-- Status Card -->
        <el-card v-if="isEdit" class="status-card" shadow="never">
          <template #header>
            <span class="card-title">{{ $t('promotions.activityStatus') }}</span>
          </template>
          <div class="status-info">
            <el-tag :type="getStatusType(promotionForm.status)" size="large">
              {{ getStatusText(promotionForm.status) }}
            </el-tag>
            <div class="status-actions">
              <el-button
                v-if="promotionForm.status === 'pending' || promotionForm.status === 'paused'"
                type="success"
                @click="handleActivate"
                :loading="activating"
              >
                {{ $t('promotions.activateActivity') }}
              </el-button>
              <el-button
                v-if="promotionForm.status === 'active'"
                type="warning"
                @click="handleDeactivate"
                :loading="activating"
              >
                {{ $t('promotions.deactivateActivity') }}
              </el-button>
            </div>
          </div>
        </el-card>

        <!-- Preview Card -->
        <el-card class="preview-card" shadow="never">
          <template #header>
            <span class="card-title">{{ $t('promotions.activityPreview') }}</span>
          </template>
          <div class="preview-content">
            <h4 class="preview-name">{{ promotionForm.name || $t('promotions.activityName') }}</h4>
            <p class="preview-desc">{{ promotionForm.description || $t('promotions.activityDesc') }}</p>
            <div class="preview-discount">
              <span class="discount-label">{{ $t('promotions.previewDiscount') }}</span>
              <span class="discount-value">
                <template v-if="promotionForm.discount_type === 'fixed_amount'">
                  ¥{{ promotionForm.discount_value_num || 0 }}
                </template>
                <template v-else>
                  {{ promotionForm.discount_value_num || 0 }}%
                </template>
              </span>
            </div>
            <div class="preview-scope">
              <span class="scope-label">{{ $t('promotions.previewScope') }}</span>
              <span class="scope-value">{{ getScopeText(promotionForm.scope_type) }}</span>
            </div>
            <div class="preview-time" v-if="promotionForm.dateRange?.length === 2">
              <span class="time-label">{{ $t('promotions.previewTime') }}</span>
              <span class="time-value">
                {{ formatDateTime(promotionForm.dateRange[0]) }} {{ $t('promotions.to') }} {{ formatDateTime(promotionForm.dateRange[1]) }}
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Rules Section -->
    <el-card v-if="isEdit" class="rules-card" shadow="never">
      <template #header>
        <div class="rules-header">
          <span class="card-title">{{ $t('promotions.promotionRules') }}</span>
          <el-button type="primary" size="small" @click="showAddRuleDialog">
            <el-icon><Plus /></el-icon>
            {{ $t('promotions.addRule') }}
          </el-button>
        </div>
      </template>
      <el-table :data="rulesList" v-loading="rulesLoading" stripe>
        <el-table-column :label="$t('promotions.ruleType')" width="150">
          <template #default="{ row }">
            {{ getRuleTypeText(row.rule_type) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('promotions.operator')" width="120">
          <template #default="{ row }">
            {{ getOperatorText(row.operator) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('promotions.ruleValue')" width="150">
          <template #default="{ row }">
            {{ row.value }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('promotions.ruleDiscountType')" width="150">
          <template #default="{ row }">
            {{ row.discount_type ? (row.discount_type === 'percentage' ? $t('promotions.percentage') : $t('promotions.fixedAmount')) : '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('promotions.ruleDiscountValue')" width="120">
          <template #default="{ row }">
            {{ row.discount_value || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('promotions.rulePriority')" width="100" align="center">
          <template #default="{ row }">
            {{ row.priority }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEditRule(row)">
              {{ $t('common.edit') }}
            </el-button>
            <el-button type="danger" link size="small" @click="handleDeleteRule(row)">
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!rulesLoading && rulesList.length === 0" :description="$t('promotions.noRules')" />
    </el-card>

    <!-- Add/Edit Rule Dialog -->
    <el-dialog v-model="ruleDialogVisible" :title="editingRule ? $t('promotions.editRule') : $t('promotions.addRule')" width="500px" destroy-on-close>
      <el-form :model="ruleForm" label-width="120px">
        <el-form-item :label="$t('promotions.ruleType')">
          <el-select v-model="ruleForm.rule_type" style="width: 100%">
            <el-option :label="$t('promotions.ruleTypeProduct')" value="product" />
            <el-option :label="$t('promotions.ruleTypeCategory')" value="category" />
            <el-option :label="$t('promotions.ruleTypeAmount')" value="amount" />
            <el-option :label="$t('promotions.ruleTypeQuantity')" value="quantity" />
            <el-option :label="$t('promotions.ruleTypeUserGroup')" value="user_group" />
            <el-option :label="$t('promotions.ruleTypeMarket')" value="market" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promotions.operator')">
          <el-select v-model="ruleForm.operator" style="width: 100%">
            <el-option :label="$t('promotions.operatorEq')" value="eq" />
            <el-option :label="$t('promotions.operatorNe')" value="ne" />
            <el-option :label="$t('promotions.operatorGt')" value="gt" />
            <el-option :label="$t('promotions.operatorGte')" value="gte" />
            <el-option :label="$t('promotions.operatorLt')" value="lt" />
            <el-option :label="$t('promotions.operatorLte')" value="lte" />
            <el-option :label="$t('promotions.operatorIn')" value="in" />
            <el-option :label="$t('promotions.operatorNotIn')" value="not_in" />
            <el-option :label="$t('promotions.operatorContains')" value="contains" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promotions.ruleValue')">
          <el-input v-model="ruleForm.value" />
        </el-form-item>
        <el-form-item :label="$t('promotions.ruleDiscountType')">
          <el-select v-model="ruleForm.discount_type" style="width: 100%">
            <el-option :label="$t('promotions.fixedAmount')" value="fixed_amount" />
            <el-option :label="$t('promotions.percentage')" value="percentage" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promotions.ruleDiscountValue')">
          <el-input-number v-model="ruleForm.discount_value_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.rulePriority')">
          <el-input-number v-model="ruleForm.priority" :min="0" :max="100" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveRule" :loading="ruleSaving">
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Footer Actions -->
    <div class="footer-actions">
      <el-button @click="handleBack">{{ $t('promotions.cancelAction') }}</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        {{ isEdit ? $t('promotions.saveChanges') : $t('promotions.createActivity') }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getPromotionStatusType } from '@/utils/status'
import {
  getPromotion, createPromotion, updatePromotion, activatePromotion, deactivatePromotion,
  getPromotionRules, createPromotionRules, updatePromotionRule, deletePromotionRule,
  type PromotionRule
} from '@/api/promotion'
import { getProductList } from '@/api/product'
import { getCategories } from '@/api/category'
import PageHeader from '@/components/common/PageHeader.vue'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const route = useRoute()
const { handleError } = useErrorHandler()

const formRef = ref()
const loading = ref(false)
const saving = ref(false)
const activating = ref(false)

const promotionId = computed(() => Number(route.params.id))
const isEdit = computed(() => promotionId.value > 0)

const promotionForm = reactive({
  name: '',
  description: '',
  type: 'discount' as 'discount' | 'flash_sale' | 'bundle' | 'buy_x_get_y',
  discount_type: 'fixed_amount' as 'fixed_amount' | 'percentage' | 'buy_x_get_y',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  dateRange: [] as string[],
  scope_type: 'storewide' as 'storewide' | 'products' | 'categories' | 'brands',
  product_ids: [] as number[],
  category_ids: [] as number[],
  usage_limit: 0,
  per_user_limit: 0,
  status: 'pending' as string
})

const rules = {
  name: [{ required: true, message: t('promotions.enterPromotionName'), trigger: 'blur' }],
  type: [{ required: true, message: t('promotions.selectPromotionType'), trigger: 'change' }],
  discount_type: [{ required: true, message: t('promotions.selectDiscountType'), trigger: 'change' }],
  dateRange: [{ required: true, message: t('promotions.selectTimePeriod'), trigger: 'change' }]
}

const productOptions = ref<{ id: number; name: string }[]>([])
const categoryOptions = ref<{ id: number; name: string }[]>([])

const loadPromotion = async () => {
  if (!promotionId.value) return

  loading.value = true
  try {
    const res = await getPromotion(promotionId.value)
    promotionForm.name = res.name
    promotionForm.description = res.description || ''
    promotionForm.type = res.type as 'discount' | 'flash_sale' | 'bundle' | 'buy_x_get_y'
    promotionForm.discount_type = (res.discount_type || 'fixed_amount') as 'fixed_amount' | 'percentage'
    promotionForm.discount_value_num = parseFloat(res.discount_value) || 0
    promotionForm.min_order_amount_num = parseFloat(res.min_order_amount) || 0
    promotionForm.max_discount_num = parseFloat(res.max_discount) || 0
    promotionForm.dateRange = [res.start_time, res.end_time]
    promotionForm.product_ids = res.product_ids || []
    promotionForm.category_ids = res.category_ids || []
    promotionForm.usage_limit = res.usage_limit || 0
    promotionForm.per_user_limit = res.per_user_limit || 0
    promotionForm.status = res.status
  } catch (error) {
    handleError(error, t('promotions.loadPromotionFailed'))
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await getCategories()
    categoryOptions.value = (res.list || []).map(c => ({ id: c.id, name: c.name }))
  } catch (error) {
    handleError(error, t('promotions.loadCategoriesFailed'))
  }
}

const searchProducts = async (query: string) => {
  if (!query) return
  try {
    const res = await getProductList({ page: 1, page_size: 20, name: query })
    productOptions.value = (res.list || []).map(p => ({ id: p.id, name: p.name }))
  } catch (error) {
    handleError(error, t('promotions.searchProductsFailed'))
  }
}

const handleBack = () => {
  router.push('/promotions')
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saving.value = true
    try {
      const data = {
        name: promotionForm.name,
        description: promotionForm.description,
        type: promotionForm.type,
        discount_type: promotionForm.discount_type,
        discount_value: String(promotionForm.discount_value_num),
        min_order_amount: String(promotionForm.min_order_amount_num),
        max_discount: String(promotionForm.max_discount_num),
        start_time: promotionForm.dateRange[0],
        end_time: promotionForm.dateRange[1],
        usage_limit: promotionForm.usage_limit,
        per_user_limit: promotionForm.per_user_limit,
        product_ids: promotionForm.product_ids,
        category_ids: promotionForm.category_ids
      }

      if (isEdit.value) {
        await updatePromotion({ id: promotionId.value, ...data })
        ElMessage.success(t('promotions.updateSuccess'))
      } else {
        await createPromotion(data)
        ElMessage.success(t('promotions.createSuccess'))
        router.push('/promotions')
      }
    } catch (error) {
      handleError(error, t('promotions.savePromotionFailed'))
    } finally {
      saving.value = false
    }
  })
}

const handleActivate = async () => {
  activating.value = true
  try {
    await activatePromotion(promotionId.value)
    ElMessage.success(t('promotions.activateSuccess'))
    promotionForm.status = 'active'
  } catch (error) {
    handleError(error, t('promotions.activatePromotionFailed'))
  } finally {
    activating.value = false
  }
}

const handleDeactivate = async () => {
  activating.value = true
  try {
    await deactivatePromotion(promotionId.value)
    ElMessage.success(t('promotions.deactivateSuccess'))
    promotionForm.status = 'paused'
  } catch (error) {
    handleError(error, t('promotions.deactivatePromotionFailed'))
  } finally {
    activating.value = false
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

const getStatusType = getPromotionStatusType

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': t('promotions.activeStatus'),
    'paused': t('promotions.paused'),
    'pending': t('promotions.pending'),
    'ended': t('promotions.ended')
  }
  return texts[status] || status
}

const getScopeText = (scope: string) => {
  const texts: Record<string, string> = {
    'storewide': t('promotions.scopeStorewide'),
    'products': t('promotions.scopeSpecificProducts'),
    'categories': t('promotions.scopeSpecificCategories'),
    'brands': t('promotions.scopeSpecificBrands')
  }
  return texts[scope] || scope
}

// Rules section
const rulesList = ref<PromotionRule[]>([])
const rulesLoading = ref(false)
const ruleDialogVisible = ref(false)
const editingRule = ref<PromotionRule | null>(null)
const ruleSaving = ref(false)
const ruleForm = reactive({
  rule_type: 'product',
  operator: 'eq',
  value: '',
  discount_type: 'fixed_amount',
  discount_value_num: 0,
  priority: 0
})

const loadRules = async () => {
  if (!isEdit.value) return
  rulesLoading.value = true
  try {
    const res = await getPromotionRules(promotionId.value)
    rulesList.value = res.list || []
  } catch (error) {
    handleError(error, t('promotions.loadRulesFailed'))
  } finally {
    rulesLoading.value = false
  }
}

const showAddRuleDialog = () => {
  editingRule.value = null
  ruleForm.rule_type = 'product'
  ruleForm.operator = 'eq'
  ruleForm.value = ''
  ruleForm.discount_type = 'fixed_amount'
  ruleForm.discount_value_num = 0
  ruleForm.priority = 0
  ruleDialogVisible.value = true
}

const handleEditRule = (rule: PromotionRule) => {
  editingRule.value = rule
  ruleForm.rule_type = rule.rule_type
  ruleForm.operator = rule.operator
  ruleForm.value = rule.value
  ruleForm.discount_type = rule.discount_type || 'fixed_amount'
  ruleForm.discount_value_num = parseFloat(rule.discount_value) || 0
  ruleForm.priority = rule.priority
  ruleDialogVisible.value = true
}

const handleDeleteRule = async (rule: PromotionRule) => {
  try {
    await ElMessageBox.confirm(t('promotions.deleteRuleConfirm'), t('common.warning'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await deletePromotionRule(rule.id)
    ElMessage.success(t('promotions.ruleDeletedSuccess'))
    loadRules()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('promotions.deleteRuleFailed'))
    }
  }
}

const handleSaveRule = async () => {
  ruleSaving.value = true
  try {
    const data = {
      rule_type: ruleForm.rule_type,
      operator: ruleForm.operator,
      value: ruleForm.value,
      discount_type: ruleForm.discount_type,
      discount_value: String(ruleForm.discount_value_num),
      priority: ruleForm.priority
    }

    if (editingRule.value) {
      await updatePromotionRule({ id: editingRule.value.id, ...data })
      ElMessage.success(t('promotions.ruleUpdatedSuccess'))
    } else {
      await createPromotionRules(promotionId.value, [data])
      ElMessage.success(t('promotions.ruleCreatedSuccess'))
    }
    ruleDialogVisible.value = false
    loadRules()
  } catch (error) {
    handleError(error, editingRule.value ? t('promotions.updateRuleFailed') : t('promotions.createRuleFailed'))
  } finally {
    ruleSaving.value = false
  }
}

const getRuleTypeText = (type: string) => {
  const texts: Record<string, string> = {
    'product': t('promotions.ruleTypeProduct'),
    'category': t('promotions.ruleTypeCategory'),
    'amount': t('promotions.ruleTypeAmount'),
    'quantity': t('promotions.ruleTypeQuantity'),
    'user_group': t('promotions.ruleTypeUserGroup'),
    'market': t('promotions.ruleTypeMarket')
  }
  return texts[type] || type
}

const getOperatorText = (operator: string) => {
  const texts: Record<string, string> = {
    'eq': t('promotions.operatorEq'),
    'ne': t('promotions.operatorNe'),
    'gt': t('promotions.operatorGt'),
    'gte': t('promotions.operatorGte'),
    'lt': t('promotions.operatorLt'),
    'lte': t('promotions.operatorLte'),
    'in': t('promotions.operatorIn'),
    'not_in': t('promotions.operatorNotIn'),
    'contains': t('promotions.operatorContains')
  }
  return texts[operator] || operator
}

onMounted(() => {
  loadCategories()
  if (isEdit.value) {
    loadPromotion()
    loadRules()
  }
})
</script>

<style scoped>
.promotion-detail-page {
  padding: 0;
}

.form-card, .status-card, .preview-card, .rules-card {
  border-radius: 16px;
  margin-bottom: 20px;
}

.rules-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #F3F4F6;
}

.card-title {
  font-weight: 600;
  color: #1E1B4B;
}

.status-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.status-actions {
  width: 100%;
}

.status-actions .el-button {
  width: 100%;
}

.preview-content {
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
}

.preview-name {
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 8px 0;
}

.preview-desc {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.preview-discount, .preview-scope, .preview-time {
  font-size: 14px;
  margin-bottom: 8px;
}

.discount-label, .scope-label, .time-label {
  color: #6B7280;
}

.discount-value {
  color: #EF4444;
  font-weight: 600;
}

.scope-value, .time-value {
  color: #1E1B4B;
}

.form-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

.footer-actions {
  position: fixed;
  bottom: 0;
  left: 200px;
  right: 0;
  padding: 16px 24px;
  background: #fff;
  border-top: 1px solid #F3F4F6;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  z-index: 100;
}

@media (max-width: 768px) {
  .footer-actions {
    left: 0;
  }
}
</style>