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
import { ElMessage } from 'element-plus'
import {
  getPromotion, createPromotion, updatePromotion, activatePromotion, deactivatePromotion
} from '@/api/promotion'
import { getProductList } from '@/api/product'
import { getCategories } from '@/api/category'
import PageHeader from '@/components/common/PageHeader.vue'
import { t } from '@/plugins/i18n'

const router = useRouter()
const route = useRoute()

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
    console.error('Failed to load promotion:', error)
    ElMessage.error(t('promotions.loadPromotionFailed'))
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await getCategories()
    categoryOptions.value = (res.list || []).map(c => ({ id: c.id, name: c.name }))
  } catch (error) {
    console.error('Failed to load categories:', error)
    ElMessage.error(t('promotions.loadCategoriesFailed'))
  }
}

const searchProducts = async (query: string) => {
  if (!query) return
  try {
    const res = await getProductList({ page: 1, page_size: 20, name: query })
    productOptions.value = (res.list || []).map(p => ({ id: p.id, name: p.name }))
  } catch (error) {
    console.error('Failed to search products:', error)
    ElMessage.error(t('promotions.searchProductsFailed'))
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
      console.error('Failed to save promotion:', error)
      ElMessage.error(t('promotions.savePromotionFailed'))
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
    console.error('Failed to activate:', error)
    ElMessage.error(t('promotions.activatePromotionFailed'))
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
    console.error('Failed to deactivate:', error)
    ElMessage.error(t('promotions.deactivatePromotionFailed'))
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

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'paused': 'warning',
    'pending': 'info',
    'ended': 'info'
  }
  return types[status] || 'info'
}

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

onMounted(() => {
  loadCategories()
  if (isEdit.value) {
    loadPromotion()
  }
})
</script>

<style scoped>
.promotion-detail-page {
  padding: 0;
}

.form-card, .status-card, .preview-card {
  border-radius: 16px;
  margin-bottom: 20px;
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