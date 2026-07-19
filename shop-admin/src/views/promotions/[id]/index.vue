<template>
  <div
    v-loading="loading"
    class="promotion-detail-page"
  >
    <PageHeader
      :title="isEdit ? (kind === 'coupon' ? $t('promotions.editCoupon') : $t('promotions.editPromotion')) : (kind === 'coupon' ? $t('promotions.createCoupon') : $t('promotions.createPromotion'))"
      :subtitle="promotionForm.name || (kind === 'coupon' ? $t('promotions.configureCouponInfo') : $t('promotions.configurePromotionInfo'))"
      @back="handleBack"
    />

    <el-row :gutter="20">
      <el-col
        :xs="24"
        :lg="16"
      >
        <el-card
          class="form-card"
          shadow="never"
        >
          <el-form
            ref="formRef"
            :model="promotionForm"
            :rules="rules"
            label-width="120px"
            label-position="top"
          >
            <!-- Basic Info -->
            <div class="section-title">
              {{ $t('promotions.basicInfo') }}
            </div>
            <el-row :gutter="20">
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item
                  :label="$t('promotions.promotionName')"
                  prop="name"
                >
                  <el-input
                    v-model="promotionForm.name"
                    :placeholder="kind === 'coupon' ? $t('promotions.enterCouponName') : $t('promotions.enterPromotionName')"
                    maxlength="100"
                    show-word-limit
                  />
                </el-form-item>
              </el-col>
              <el-col
                v-if="kind === 'coupon'"
                :xs="24"
                :sm="12"
              >
                <el-form-item
                  :label="$t('promotions.couponCode')"
                  prop="code"
                >
                  <el-input
                    v-model="promotionForm.code"
                    :placeholder="$t('promotions.enterCouponCode')"
                    maxlength="50"
                  >
                    <template #append>
                      <el-button @click="generateCouponCode">
                        {{ $t('promotions.generate') }}
                      </el-button>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item
                  :label="$t('promotions.promotionTypeSelect')"
                  prop="type"
                >
                  <el-select
                    v-model="promotionForm.type"
                    :placeholder="$t('promotions.selectPromotionType')"
                    style="width: 100%"
                  >
                    <template v-if="kind === 'coupon'">
                      <el-option
                        :label="$t('promotions.fixedAmount')"
                        value="fixed_amount"
                      />
                      <el-option
                        :label="$t('promotions.percentage')"
                        value="percentage"
                      />
                      <el-option
                        :label="$t('promotions.freeShipping')"
                        value="free_shipping"
                      />
                    </template>
                    <template v-else>
                      <el-option
                        :label="$t('promotions.discount')"
                        value="discount"
                      />
                      <el-option
                        :label="$t('promotions.flashSale')"
                        value="flash_sale"
                      />
                      <el-option
                        :label="$t('promotions.bundle')"
                        value="bundle"
                      />
                      <el-option
                        :label="$t('promotions.buyXGetY')"
                        value="buy_x_get_y"
                      />
                    </template>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item :label="$t('promotions.currency')">
                  <el-select
                    v-model="promotionForm.currency"
                    style="width: 100%"
                  >
                    <el-option label="CNY" value="CNY" />
                    <el-option label="USD" value="USD" />
                    <el-option label="EUR" value="EUR" />
                    <el-option label="JPY" value="JPY" />
                    <el-option label="GBP" value="GBP" />
                    <el-option label="SGD" value="SGD" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item :label="$t('promotions.market')">
                  <el-select
                    v-model="promotionForm.market_id"
                    :placeholder="$t('promotions.selectMarket')"
                    clearable
                    style="width: 100%"
                  >
                    <el-option
                      v-for="m in marketOptions"
                      :key="m.id"
                      :label="m.name"
                      :value="m.id"
                    />
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

            <!-- Rules -->
            <div class="section-title">
              <div class="rules-header">
                <span>{{ kind === 'coupon' ? $t('promotions.couponRules') : $t('promotions.promotionRules') }}</span>
                <el-button
                  type="primary"
                  size="small"
                  @click="showAddRuleDialog"
                >
                  <el-icon><Plus /></el-icon>
                  {{ $t('promotions.addRule') }}
                </el-button>
              </div>
            </div>
            <el-table
              v-loading="rulesLoading"
              :data="rulesList"
              stripe
              style="margin-bottom: 20px;"
            >
              <el-table-column
                :label="$t('promotions.conditionType')"
                width="150"
              >
                <template #default="{ row }">
                  {{ formatConditionType(row.condition_type) }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('promotions.conditionValue')"
                width="120"
              >
                <template #default="{ row }">
                  {{ row.condition_value }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('promotions.actionType')"
                width="150"
              >
                <template #default="{ row }">
                  {{ formatActionType(row.action_type) }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('promotions.actionValue')"
                width="120"
              >
                <template #default="{ row }">
                  <template v-if="row.action_type === 'fixed_amount'">
                    {{ currencySymbol(promotionForm.currency) }}{{ row.action_value }}
                  </template>
                  <template v-else-if="row.action_type === 'percentage'">
                    {{ row.action_value }}%
                  </template>
                  <template v-else>
                    -
                  </template>
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('promotions.maxDiscount')"
                width="120"
              >
                <template #default="{ row }">
                  {{ row.max_discount && parseFloat(row.max_discount) > 0 ? `${currencySymbol(promotionForm.currency)}${row.max_discount}` : '-' }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('promotions.sortOrder')"
                width="80"
                align="center"
              >
                <template #default="{ row }">
                  {{ row.sort_order ?? '-' }}
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('common.actions')"
                width="150"
                fixed="right"
              >
                <template #default="{ row }">
                  <el-button
                    type="primary"
                    link
                    size="small"
                    @click="handleEditRule(row)"
                  >
                    {{ $t('common.edit') }}
                  </el-button>
                  <el-button
                    type="danger"
                    link
                    size="small"
                    @click="handleDeleteRule(row)"
                  >
                    {{ $t('common.delete') }}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty
              v-if="!rulesLoading && rulesList.length === 0"
              :description="$t('promotions.noRules')"
            />

            <!-- Time Range -->
            <div class="section-title">
              {{ $t('promotions.activityTime') }}
            </div>
            <el-form-item
              :label="$t('promotions.promotionPeriod')"
              prop="dateRange"
            >
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
            <div class="section-title">
              {{ $t('promotions.scopeSettings') }}
            </div>
            <el-form-item :label="$t('promotions.activityScope')">
              <el-radio-group v-model="promotionForm.scope_type">
                <el-radio label="storewide">
                  {{ $t('promotions.storewide') }}
                </el-radio>
                <el-radio label="products">
                  {{ $t('promotions.specificProductsScope') }}
                </el-radio>
                <el-radio label="categories">
                  {{ $t('promotions.specificCategories') }}
                </el-radio>
                <el-radio label="brands">
                  {{ $t('promotions.specificBrands') }}
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item
              v-if="promotionForm.scope_type === 'products'"
              :label="$t('promotions.selectProducts')"
            >
              <el-select
                v-model="promotionForm.scope_ids"
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

            <el-form-item
              v-if="promotionForm.scope_type === 'categories'"
              :label="$t('promotions.selectCategories')"
            >
              <el-select
                v-model="promotionForm.scope_ids"
                multiple
                :placeholder="$t('promotions.selectCategoryPlaceholder')"
                style="width: 100%"
              >
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              v-if="promotionForm.scope_type === 'brands'"
              :label="$t('promotions.selectBrand')"
            >
              <el-select
                v-model="promotionForm.scope_ids"
                multiple
                :placeholder="$t('promotions.selectBrand')"
                style="width: 100%"
              >
                <el-option
                  v-for="item in brandOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <!-- Usage Limits -->
            <div class="section-title">
              {{ $t('promotions.usageLimits') }}
            </div>
            <el-row :gutter="20">
              <el-col
                v-if="kind === 'coupon'"
                :xs="24"
                :sm="12"
              >
                <el-form-item :label="$t('promotions.totalCount')">
                  <el-input-number
                    v-model="promotionForm.total_count"
                    :min="0"
                    style="width: 100%"
                  />
                  <div class="form-tip">
                    {{ $t('promotions.totalCountTip') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item :label="$t('promotions.totalUsageLimit')">
                  <el-input-number
                    v-model="promotionForm.usage_limit"
                    :min="0"
                    style="width: 100%"
                  />
                  <div class="form-tip">
                    {{ $t('promotions.zeroUnlimitedUsage') }}
                  </div>
                </el-form-item>
              </el-col>
              <el-col
                :xs="24"
                :sm="12"
              >
                <el-form-item :label="$t('promotions.perUserUsageLimit')">
                  <el-input-number
                    v-model="promotionForm.per_user_limit"
                    :min="0"
                    style="width: 100%"
                  />
                  <div class="form-tip">
                    {{ $t('promotions.zeroUnlimitedUsage') }}
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :lg="8"
      >
        <!-- Status Card -->
        <el-card
          v-if="isEdit"
          class="status-card"
          shadow="never"
        >
          <template #header>
            <span class="card-title">{{ $t('promotions.activityStatus') }}</span>
          </template>
          <div class="status-info">
            <el-tag
              :type="getStatusType(promotionForm.status)"
              size="large"
            >
              {{ getStatusText(promotionForm.status) }}
            </el-tag>
            <div class="status-actions">
              <el-button
                v-if="promotionForm.status === 'pending' || promotionForm.status === 'paused'"
                type="success"
                :loading="activating"
                @click="handleActivate"
              >
                {{ $t('promotions.activateActivity') }}
              </el-button>
              <el-button
                v-if="promotionForm.status === 'active'"
                type="warning"
                :loading="activating"
                @click="handleDeactivate"
              >
                {{ $t('promotions.deactivateActivity') }}
              </el-button>
            </div>
          </div>
        </el-card>

        <!-- Preview Card -->
        <el-card
          class="preview-card"
          shadow="never"
        >
          <template #header>
            <span class="card-title">{{ $t('promotions.activityPreview') }}</span>
          </template>
          <div class="preview-content">
            <h4 class="preview-name">
              {{ promotionForm.name || $t('promotions.activityName') }}
            </h4>
            <p class="preview-desc">
              {{ promotionForm.description || $t('promotions.activityDesc') }}
            </p>
            <div
              v-if="rulesList.length > 0"
              class="preview-discount"
            >
              <span class="discount-label">{{ $t('promotions.previewDiscount') }}</span>
              <span class="discount-value">
                <template v-if="rulesList[0].action_type === 'fixed_amount'">
                  {{ currencySymbol(promotionForm.currency) }}{{ rulesList[0].action_value || 0 }}
                </template>
                <template v-else-if="rulesList[0].action_type === 'percentage'">
                  {{ rulesList[0].action_value || 0 }}%
                </template>
                <template v-else>
                  {{ $t('promotions.freeShipping') }}
                </template>
              </span>
            </div>
            <div class="preview-scope">
              <span class="scope-label">{{ $t('promotions.previewScope') }}</span>
              <span class="scope-value">{{ getScopeText(promotionForm.scope_type) }}</span>
            </div>
            <div
              v-if="promotionForm.dateRange?.length === 2"
              class="preview-time"
            >
              <span class="time-label">{{ $t('promotions.previewTime') }}</span>
              <span class="time-value">
                {{ formatDateTime(promotionForm.dateRange[0]) }} {{ $t('promotions.to') }} {{ formatDateTime(promotionForm.dateRange[1]) }}
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Add/Edit Rule Dialog -->
    <el-dialog
      v-model="ruleDialogVisible"
      :title="editingRule ? $t('promotions.editRule') : $t('promotions.addRule')"
      width="500px"
      destroy-on-close
    >
      <el-form
        :model="ruleForm"
        label-width="120px"
      >
        <el-form-item :label="$t('promotions.conditionType')">
          <el-select
            v-model="ruleForm.condition_type"
            style="width: 100%"
          >
            <el-option
              :label="$t('promotions.minAmount')"
              value="min_amount"
            />
            <el-option
              :label="$t('promotions.minQuantity')"
              value="min_quantity"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promotions.conditionValue')">
          <el-input-number
            v-model="ruleForm.condition_value_num"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.actionType')">
          <el-select
            v-model="ruleForm.action_type"
            style="width: 100%"
          >
            <el-option
              :label="$t('promotions.fixedAmountType')"
              value="fixed_amount"
            />
            <el-option
              :label="$t('promotions.percentageType')"
              value="percentage"
            />
            <el-option
              :label="$t('promotions.freeShipping')"
              value="free_shipping"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="ruleForm.action_type !== 'free_shipping'"
          :label="$t('promotions.actionValue')"
        >
          <el-input-number
            v-model="ruleForm.action_value_num"
            :min="0"
            :max="ruleForm.action_type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item
          v-if="ruleForm.action_type === 'percentage'"
          :label="$t('promotions.maxDiscount')"
        >
          <el-input-number
            v-model="ruleForm.max_discount_num"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.sortOrder')">
          <el-input-number
            v-model="ruleForm.sort_order"
            :min="0"
            :max="100"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="ruleSaving"
          @click="handleSaveRule"
        >
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Footer Actions -->
    <div class="footer-actions">
      <el-button @click="handleBack">
        {{ $t('promotions.cancelAction') }}
      </el-button>
      <el-button
        type="primary"
        :loading="saving"
        @click="handleSave"
      >
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
  type PromotionKind, type PromotionRule, type PromotionRuleRequest
} from '@/api/promotion'
import { getProductList } from '@/api/product'
import { getCategories } from '@/api/category'
import { getBrands } from '@/api/brand'
import { getMarkets, type Market } from '@/api/market'
import PageHeader from '@/components/common/PageHeader.vue'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'
import { currencySymbol } from '@/utils/currency'

const router = useRouter()
const route = useRoute()
const { handleError } = useErrorHandler()

const formRef = ref()
const loading = ref(false)
const saving = ref(false)
const activating = ref(false)

const promotionId = computed(() => route.params.id as string)
const isEdit = computed(() => !!promotionId.value)
const kind = ref<PromotionKind>('promotion')

const promotionForm = reactive({
  name: '',
  description: '',
  code: '',
  type: 'discount' as string,
  currency: 'CNY',
  market_id: '',
  dateRange: [] as string[],
  scope_type: 'storewide' as 'storewide' | 'products' | 'categories' | 'brands',
  scope_ids: [] as string[],
  exclude_ids: [] as string[],
  usage_limit: 0,
  per_user_limit: 0,
  total_count: 0,
  status: 'pending' as string
})

const rules = {
  name: [{ required: true, message: t('promotions.enterPromotionName'), trigger: 'blur' }],
  type: [{ required: true, message: t('promotions.selectPromotionType'), trigger: 'change' }],
  dateRange: [{ required: true, message: t('promotions.selectTimePeriod'), trigger: 'change' }]
}

const productOptions = ref<{ id: string; name: string }[]>([])
const categoryOptions = ref<{ id: string; name: string }[]>([])
const brandOptions = ref<{ id: string; name: string }[]>([])
const marketOptions = ref<Market[]>([])

const loadPromotion = async () => {
  if (!promotionId.value) return

  loading.value = true
  try {
    const res = await getPromotion(promotionId.value)
    kind.value = res.kind
    promotionForm.name = res.name
    promotionForm.description = res.description || ''
    promotionForm.code = res.code || ''
    promotionForm.type = res.type || (res.kind === 'coupon' ? 'fixed_amount' : 'discount')
    promotionForm.currency = res.currency || 'CNY'
    promotionForm.market_id = res.market_id || ''
    promotionForm.dateRange = [res.start_time, res.end_time]
    if (res.scope_type) {
      promotionForm.scope_type = res.scope_type
    }
    promotionForm.scope_ids = res.scope_ids || []
    promotionForm.exclude_ids = res.exclude_ids || []
    promotionForm.usage_limit = res.usage_limit || 0
    promotionForm.per_user_limit = res.per_user_limit || 0
    promotionForm.total_count = res.total_count || 0
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

const loadBrands = async () => {
  try {
    const res = await getBrands({ page: 1, page_size: 200 })
    brandOptions.value = (res.list || []).map(b => ({ id: b.id, name: b.name }))
  } catch (error) {
    handleError(error, t('promotions.loadBrandsFailed'))
  }
}

const loadMarkets = async () => {
  try {
    const res = await getMarkets()
    marketOptions.value = res.list || []
  } catch (error) {
    handleError(error, t('promotions.loadMarketsFailed'))
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
      const baseData = {
        kind: kind.value,
        name: promotionForm.name,
        description: promotionForm.description,
        type: promotionForm.type,
        currency: promotionForm.currency,
        market_id: promotionForm.market_id || undefined,
        start_time: promotionForm.dateRange[0],
        end_time: promotionForm.dateRange[1],
        usage_limit: promotionForm.usage_limit,
        per_user_limit: promotionForm.per_user_limit,
        scope_type: promotionForm.scope_type,
        scope_ids: promotionForm.scope_ids,
        exclude_ids: promotionForm.exclude_ids
      }

      const couponExtra = kind.value === 'coupon'
        ? { code: promotionForm.code, total_count: promotionForm.total_count }
        : {}

      const data = { ...baseData, ...couponExtra }

      if (isEdit.value) {
        await updatePromotion({ id: promotionId.value, ...data })
        ElMessage.success(t('promotions.updateSuccess'))
      } else {
        await createPromotion(data as Parameters<typeof createPromotion>[0])
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

const generateCouponCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 10; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  promotionForm.code = code
}

const getStatusType = getPromotionStatusType

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': t('promotions.activeStatus'),
    'paused': t('promotions.paused'),
    'pending': t('promotions.pending'),
    'ended': t('promotions.ended'),
    'expired': t('promotions.expiredStatus'),
    'depleted': t('promotions.depletedStatus')
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

const formatConditionType = (val: string) => {
  return val === 'min_amount' ? t('promotions.minAmount') : t('promotions.minQuantity')
}

const formatActionType = (val: string) => {
  const map: Record<string, string> = {
    'fixed_amount': t('promotions.fixedAmountType'),
    'percentage': t('promotions.percentageType'),
    'free_shipping': t('promotions.freeShipping')
  }
  return map[val] || val
}

// Rules section
const rulesList = ref<PromotionRule[]>([])
const rulesLoading = ref(false)
const ruleDialogVisible = ref(false)
const editingRule = ref<PromotionRule | null>(null)
const ruleSaving = ref(false)
const ruleForm = reactive({
  condition_type: 'min_amount' as 'min_amount' | 'min_quantity',
  condition_value_num: 0,
  action_type: 'fixed_amount' as 'fixed_amount' | 'percentage' | 'free_shipping',
  action_value_num: 0,
  max_discount_num: 0,
  sort_order: 0
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
  ruleForm.condition_type = 'min_amount'
  ruleForm.condition_value_num = 0
  ruleForm.action_type = 'fixed_amount'
  ruleForm.action_value_num = 0
  ruleForm.max_discount_num = 0
  ruleForm.sort_order = 0
  ruleDialogVisible.value = true
}

const handleEditRule = (rule: PromotionRule) => {
  editingRule.value = rule
  ruleForm.condition_type = rule.condition_type
  ruleForm.condition_value_num = parseFloat(rule.condition_value) || 0
  ruleForm.action_type = rule.action_type
  ruleForm.action_value_num = parseFloat(rule.action_value) || 0
  ruleForm.max_discount_num = parseFloat(rule.max_discount || '0') || 0
  ruleForm.sort_order = rule.sort_order || 0
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
    const data: PromotionRuleRequest = {
      condition_type: ruleForm.condition_type,
      condition_value: String(ruleForm.condition_value_num),
      action_type: ruleForm.action_type,
      action_value: ruleForm.action_type === 'free_shipping' ? '0' : String(ruleForm.action_value_num),
      max_discount: ruleForm.action_type === 'percentage' ? String(ruleForm.max_discount_num) : undefined,
      sort_order: ruleForm.sort_order
    }

    if (editingRule.value) {
      await updatePromotionRule({ ...data, id: editingRule.value.id })
      ElMessage.success(t('promotions.ruleUpdatedSuccess'))
    } else {
      await createPromotionRules(promotionId.value, kind.value, [data])
      ElMessage.success(t('promotions.ruleCreatedSuccess'))
    }
    ruleDialogVisible.value = false
    loadRules()
  } catch (error) {
    handleError(
      error,
      editingRule.value
        ? t('promotions.updateRuleFailed')
        : t('promotions.createRuleFailed')
    )
  } finally {
    ruleSaving.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadBrands()
  loadMarkets()
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
