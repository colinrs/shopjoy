<template>
  <div class="promotion-detail-page" v-loading="loading">
    <PageHeader
      :title="isEdit ? '编辑促销活动' : '创建促销活动'"
      :subtitle="promotionForm.name || '配置促销活动信息'"
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
            <div class="section-title">基本信息</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item label="活动名称" prop="name">
                  <el-input v-model="promotionForm.name" placeholder="请输入活动名称" maxlength="100" show-word-limit />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="活动类型" prop="type">
                  <el-select v-model="promotionForm.type" placeholder="请选择活动类型" style="width: 100%">
                    <el-option label="折扣" value="discount" />
                    <el-option label="满减" value="full_reduce" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="活动描述">
              <el-input
                v-model="promotionForm.description"
                type="textarea"
                :rows="3"
                placeholder="请输入活动描述"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <!-- Discount Settings -->
            <div class="section-title">优惠设置</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="8">
                <el-form-item label="优惠类型" prop="discount_type">
                  <el-select v-model="promotionForm.discount_type" placeholder="请选择优惠类型" style="width: 100%">
                    <el-option label="固定金额" value="fixed_amount" />
                    <el-option label="百分比折扣" value="percentage" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8">
                <el-form-item :label="promotionForm.discount_type === 'fixed_amount' ? '优惠金额' : '折扣比例'" prop="discount_value">
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
                <el-form-item label="最低消费">
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
                <el-form-item label="最大优惠金额">
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
            <div class="section-title">活动时间</div>
            <el-form-item label="活动有效期" prop="dateRange">
              <el-date-picker
                v-model="promotionForm.dateRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                value-format="YYYY-MM-DDTHH:mm:ss[Z]"
                style="width: 100%"
              />
            </el-form-item>

            <!-- Scope Settings -->
            <div class="section-title">适用范围</div>
            <el-form-item label="活动范围">
              <el-radio-group v-model="promotionForm.scope_type">
                <el-radio label="storewide">全场</el-radio>
                <el-radio label="products">指定商品</el-radio>
                <el-radio label="categories">指定分类</el-radio>
                <el-radio label="brands">指定品牌</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item v-if="promotionForm.scope_type === 'products'" label="选择商品">
              <el-select
                v-model="promotionForm.product_ids"
                multiple
                filterable
                remote
                reserve-keyword
                placeholder="搜索商品"
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

            <el-form-item v-if="promotionForm.scope_type === 'categories'" label="选择分类">
              <el-select v-model="promotionForm.category_ids" multiple placeholder="请选择分类" style="width: 100%">
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <!-- Usage Limits -->
            <div class="section-title">使用限制</div>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item label="总使用次数">
                  <el-input-number v-model="promotionForm.usage_limit" :min="0" style="width: 100%" />
                  <div class="form-tip">0 表示不限制</div>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="每人限用次数">
                  <el-input-number v-model="promotionForm.per_user_limit" :min="0" style="width: 100%" />
                  <div class="form-tip">0 表示不限制</div>
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
            <span class="card-title">活动状态</span>
          </template>
          <div class="status-info">
            <el-tag :type="getStatusType(promotionForm.status)" size="large">
              {{ getStatusText(promotionForm.status) }}
            </el-tag>
            <div class="status-actions">
              <el-button
                v-if="promotionForm.status === 'draft' || promotionForm.status === 'inactive'"
                type="success"
                @click="handleActivate"
                :loading="activating"
              >
                激活活动
              </el-button>
              <el-button
                v-if="promotionForm.status === 'active'"
                type="warning"
                @click="handleDeactivate"
                :loading="activating"
              >
                停用活动
              </el-button>
            </div>
          </div>
        </el-card>

        <!-- Preview Card -->
        <el-card class="preview-card" shadow="never">
          <template #header>
            <span class="card-title">活动预览</span>
          </template>
          <div class="preview-content">
            <h4 class="preview-name">{{ promotionForm.name || '活动名称' }}</h4>
            <p class="preview-desc">{{ promotionForm.description || '活动描述' }}</p>
            <div class="preview-discount">
              <span class="discount-label">优惠：</span>
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
              <span class="scope-label">适用范围：</span>
              <span class="scope-value">{{ getScopeText(promotionForm.scope_type) }}</span>
            </div>
            <div class="preview-time" v-if="promotionForm.dateRange?.length === 2">
              <span class="time-label">活动时间：</span>
              <span class="time-value">
                {{ formatDateTime(promotionForm.dateRange[0]) }} 至 {{ formatDateTime(promotionForm.dateRange[1]) }}
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Footer Actions -->
    <div class="footer-actions">
      <el-button @click="handleBack">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">
        {{ isEdit ? '保存修改' : '创建活动' }}
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
  type: 'discount' as 'discount' | 'full_reduce',
  discount_type: 'fixed_amount' as 'fixed_amount' | 'percentage',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  dateRange: [] as string[],
  scope_type: 'storewide' as 'storewide' | 'products' | 'categories' | 'brands',
  product_ids: [] as number[],
  category_ids: [] as number[],
  usage_limit: 0,
  per_user_limit: 0,
  status: 'draft' as string
})

const rules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择活动类型', trigger: 'change' }],
  discount_type: [{ required: true, message: '请选择优惠类型', trigger: 'change' }],
  dateRange: [{ required: true, message: '请选择活动时间', trigger: 'change' }]
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
    promotionForm.type = res.type as 'discount' | 'full_reduce'
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
    ElMessage.error('加载促销活动失败')
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
  }
}

const searchProducts = async (query: string) => {
  if (!query) return
  try {
    const res = await getProductList({ page: 1, page_size: 20, name: query })
    productOptions.value = (res.list || []).map(p => ({ id: p.id, name: p.name }))
  } catch (error) {
    console.error('Failed to search products:', error)
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
        ElMessage.success('保存成功')
      } else {
        await createPromotion(data)
        ElMessage.success('创建成功')
        router.push('/promotions')
      }
    } catch (error) {
      console.error('Failed to save promotion:', error)
    } finally {
      saving.value = false
    }
  })
}

const handleActivate = async () => {
  activating.value = true
  try {
    await activatePromotion(promotionId.value)
    ElMessage.success('激活成功')
    promotionForm.status = 'active'
  } catch (error) {
    console.error('Failed to activate:', error)
  } finally {
    activating.value = false
  }
}

const handleDeactivate = async () => {
  activating.value = true
  try {
    await deactivatePromotion(promotionId.value)
    ElMessage.success('停用成功')
    promotionForm.status = 'inactive'
  } catch (error) {
    console.error('Failed to deactivate:', error)
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
    'inactive': 'warning',
    'draft': 'info',
    'expired': 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': '进行中',
    'inactive': '已停用',
    'draft': '草稿',
    'expired': '已过期'
  }
  return texts[status] || status
}

const getScopeText = (scope: string) => {
  const texts: Record<string, string> = {
    'storewide': '全场商品',
    'products': '指定商品',
    'categories': '指定分类',
    'brands': '指定品牌'
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