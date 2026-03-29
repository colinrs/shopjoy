<template>
  <div class="promotions-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <el-icon class="stat-icon" size="32"><Ticket /></el-icon>
          <div class="stat-info">
            <p class="stat-value">{{ stats.totalCoupons }}</p>
            <p class="stat-label">优惠券总数</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <el-icon class="stat-icon active" size="32"><CircleCheck /></el-icon>
          <div class="stat-info">
            <p class="stat-value">{{ stats.activeCoupons }}</p>
            <p class="stat-label">进行中</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <el-icon class="stat-icon used" size="32"><User /></el-icon>
          <div class="stat-info">
            <p class="stat-value">{{ stats.totalUsed }}</p>
            <p class="stat-label">已使用</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <el-icon class="stat-icon promotions" size="32"><Present /></el-icon>
          <div class="stat-info">
            <p class="stat-value">{{ stats.totalPromotions }}</p>
            <p class="stat-label">促销活动</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-card class="tabs-card" shadow="never">
      <el-tabs v-model="activeTab" class="promotion-tabs" @tab-change="handleTabChange">
        <!-- Coupons Tab -->
        <el-tab-pane label="优惠券管理" name="coupon">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="couponParams.name"
                placeholder="搜索优惠券名称"
                clearable
                class="search-input"
                @clear="loadCoupons"
                @keyup.enter="loadCoupons"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="couponParams.status" placeholder="状态" clearable class="filter-select" @change="loadCoupons">
                <el-option label="全部" value="" />
                <el-option label="未激活" value="inactive" />
                <el-option label="已激活" value="active" />
              </el-select>
              <el-select v-model="couponParams.type" placeholder="类型" clearable class="filter-select" @change="loadCoupons">
                <el-option label="全部" value="" />
                <el-option label="固定金额" value="fixed_amount" />
                <el-option label="百分比" value="percentage" />
              </el-select>
            </div>
            <el-button type="primary" @click="handleAddCoupon">
              <el-icon><Plus /></el-icon>创建优惠券
            </el-button>
          </div>

          <el-table :data="couponList" v-loading="couponLoading" stripe>
            <el-table-column label="优惠券信息" min-width="250">
              <template #default="{ row }">
                <div class="coupon-cell">
                  <div class="coupon-icon" :class="row.type">
                    <el-icon size="24"><Ticket /></el-icon>
                  </div>
                  <div class="coupon-details">
                    <p class="coupon-name">{{ row.name }}</p>
                    <p class="coupon-code">优惠码: {{ row.code }}</p>
                    <div class="coupon-tags">
                      <el-tag size="small" :type="row.type === 'fixed_amount' ? 'success' : 'warning'">
                        {{ row.type === 'fixed_amount' ? '固定金额' : '百分比' }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="优惠力度" width="150" align="center">
              <template #default="{ row }">
                <div class="discount-value">
                  <span v-if="row.type === 'fixed_amount'">¥{{ row.discount_value }}</span>
                  <span v-else>{{ row.discount_value }}%</span>
                </div>
                <div class="min-order" v-if="row.min_order_amount && parseFloat(row.min_order_amount) > 0">
                  满¥{{ row.min_order_amount }}可用
                </div>
              </template>
            </el-table-column>
            <el-table-column label="使用统计" width="180" align="center">
              <template #default="{ row }">
                <div class="usage-stats">
                  <el-progress
                    :percentage="getUsagePercentage(row)"
                    :status="getProgressStatus(row)"
                    :stroke-width="8"
                  />
                  <p class="usage-text">{{ row.used_count }}/{{ row.usage_limit > 0 ? row.usage_limit : '∞' }} 已使用</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="有效期" width="200">
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">至</p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getCouponStatusType(row.status)" effect="light" size="small">
                  {{ getCouponStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditCoupon(row)">
                  编辑
                </el-button>
                <el-button type="primary" link size="small" @click="handleCouponUsage(row)">
                  数据
                </el-button>
                <el-button type="danger" link size="small" @click="handleDeleteCoupon(row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <TablePagination
            v-model:current-page="couponParams.page"
            v-model:page-size="couponParams.page_size"
            :total="couponTotal"
            @change="handleCouponPageChange"
          />
        </el-tab-pane>

        <!-- Promotions Tab -->
        <el-tab-pane label="促销活动" name="promotion">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="promotionParams.name"
                placeholder="搜索活动"
                clearable
                class="search-input"
                @clear="loadPromotions"
                @keyup.enter="loadPromotions"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="promotionParams.status" placeholder="状态" clearable class="filter-select" @change="loadPromotions">
                <el-option label="全部" value="" />
                <el-option label="待开始" value="pending" />
                <el-option label="进行中" value="active" />
                <el-option label="已暂停" value="paused" />
                <el-option label="已结束" value="ended" />
              </el-select>
              <el-select v-model="promotionParams.type" placeholder="类型" clearable class="filter-select" @change="loadPromotions">
                <el-option label="全部" value="" />
                <el-option label="折扣" value="discount" />
                <el-option label="限时秒杀" value="flash_sale" />
                <el-option label="捆绑销售" value="bundle" />
                <el-option label="买X送Y" value="buy_x_get_y" />
              </el-select>
            </div>
            <el-button type="primary" @click="handleAddPromotion">
              <el-icon><Plus /></el-icon>创建活动
            </el-button>
          </div>

          <el-table :data="promotionList" v-loading="promotionLoading" stripe>
            <el-table-column label="活动信息" min-width="250">
              <template #default="{ row }">
                <div class="promo-cell">
                  <div class="promo-icon" :class="row.type">
                    <el-icon size="24"><Present /></el-icon>
                  </div>
                  <div class="promo-details">
                    <p class="promo-name">{{ row.name }}</p>
                    <p class="promo-desc">{{ row.description || '暂无描述' }}</p>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="类型" width="120" align="center">
              <template #default="{ row }">
                <el-tag size="small" :type="row.type === 'full_reduce' ? 'danger' : 'primary'">
                  {{ getPromotionTypeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="优惠内容" width="150" align="center">
              <template #default="{ row }">
                <div class="discount-value">
                  <span v-if="row.discount_type === 'fixed_amount'">¥{{ row.discount_value }}</span>
                  <span v-else-if="row.discount_type === 'percentage'">{{ row.discount_value }}%</span>
                  <span v-else>-</span>
                </div>
                <div class="min-order" v-if="row.min_order_amount && parseFloat(row.min_order_amount) > 0">
                  满¥{{ row.min_order_amount }}可用
                </div>
              </template>
            </el-table-column>
            <el-table-column label="有效期" width="200">
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">至</p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getPromoStatusType(row.status)" effect="light" size="small">
                  {{ getPromoStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditPromotion(row)">
                  编辑
                </el-button>
                <el-button
                  v-if="row.status === 'pending' || row.status === 'paused'"
                  type="success"
                  link
                  size="small"
                  @click="handleActivatePromotion(row)"
                >
                  激活
                </el-button>
                <el-button
                  v-if="row.status === 'active'"
                  type="warning"
                  link
                  size="small"
                  @click="handleDeactivatePromotion(row)"
                >
                  停用
                </el-button>
                <el-button
                  v-if="row.status !== 'active'"
                  type="danger"
                  link
                  size="small"
                  @click="handleDeletePromotion(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <TablePagination
            v-model:current-page="promotionParams.page"
            v-model:page-size="promotionParams.page_size"
            :total="promotionTotal"
            @change="handlePromotionPageChange"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Coupon Dialog -->
    <el-dialog v-model="couponDialogVisible" :title="isEditCoupon ? '编辑优惠券' : '创建优惠券'" width="600px" destroy-on-close>
      <el-form :model="couponForm" label-width="100px" :rules="couponRules" ref="couponFormRef">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="couponForm.name" placeholder="例如: 新用户专享券" maxlength="100" />
        </el-form-item>
        <el-form-item label="优惠码" prop="code">
          <el-input v-model="couponForm.code" placeholder="例如: NEWUSER2024" maxlength="50">
            <template #append>
              <el-button @click="generateCouponCode">生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="couponForm.description" type="textarea" :rows="2" placeholder="优惠券描述" />
        </el-form-item>
        <el-form-item label="优惠类型" prop="type">
          <el-radio-group v-model="couponForm.type">
            <el-radio label="fixed_amount">固定金额</el-radio>
            <el-radio label="percentage">百分比折扣</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="couponForm.type === 'fixed_amount' ? '优惠金额' : '折扣比例'" prop="discount_value">
          <el-input-number
            v-model="couponForm.discount_value_num"
            :min="0"
            :max="couponForm.type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="最低消费">
          <el-input-number v-model="couponForm.min_order_amount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="couponForm.type === 'percentage'" label="最大优惠">
          <el-input-number v-model="couponForm.max_discount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="发放总量">
          <el-input-number v-model="couponForm.usage_limit" :min="0" style="width: 100%" />
          <div class="form-tip">0 表示不限制</div>
        </el-form-item>
        <el-form-item label="每人限领">
          <el-input-number v-model="couponForm.per_user_limit" :min="0" style="width: 100%" />
          <div class="form-tip">0 表示不限制</div>
        </el-form-item>
        <el-form-item label="有效期" prop="dateRange">
          <el-date-picker
            v-model="couponForm.dateRange"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="couponDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCoupon" :loading="saveLoading">保存</el-button>
      </template>
    </el-dialog>

    <!-- Promotion Dialog -->
    <el-dialog v-model="promotionDialogVisible" :title="isEditPromotion ? '编辑促销活动' : '创建促销活动'" width="700px" destroy-on-close>
      <el-form :model="promotionForm" label-width="100px" :rules="promotionRules" ref="promotionFormRef">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="promotionForm.name" placeholder="例如: 春季大促" maxlength="100" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="promotionForm.description" type="textarea" :rows="2" placeholder="活动描述" />
        </el-form-item>
        <el-form-item label="活动类型" prop="type">
          <el-radio-group v-model="promotionForm.type">
            <el-radio label="discount">折扣</el-radio>
            <el-radio label="flash_sale">限时秒杀</el-radio>
            <el-radio label="bundle">捆绑销售</el-radio>
            <el-radio label="buy_x_get_y">买X送Y</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="优惠类型" prop="discount_type">
          <el-radio-group v-model="promotionForm.discount_type">
            <el-radio label="fixed_amount">固定金额</el-radio>
            <el-radio label="percentage">百分比折扣</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="promotionForm.discount_type === 'fixed_amount' ? '优惠金额' : '折扣比例'" prop="discount_value">
          <el-input-number
            v-model="promotionForm.discount_value_num"
            :min="0"
            :max="promotionForm.discount_type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="最低消费">
          <el-input-number v-model="promotionForm.min_order_amount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="promotionForm.discount_type === 'percentage'" label="最大优惠">
          <el-input-number v-model="promotionForm.max_discount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="有效期" prop="dateRange">
          <el-date-picker
            v-model="promotionForm.dateRange"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="promotionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePromotion" :loading="saveLoading">保存</el-button>
      </template>
    </el-dialog>

    <!-- Coupon Usage Dialog -->
    <el-dialog v-model="usageDialogVisible" title="优惠券使用记录" width="800px">
      <el-table :data="usageList" v-loading="usageLoading" stripe>
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="order_id" label="订单ID" width="120" />
        <el-table-column prop="discount_amount" label="优惠金额" width="120">
          <template #default="{ row }">
            ¥{{ row.discount_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="used_at" label="使用时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.used_at) }}
          </template>
        </el-table-column>
      </el-table>
      <TablePagination
        v-model:current-page="usageParams.page"
        v-model:page-size="usageParams.page_size"
        :total="usageTotal"
        @change="handleUsagePageChange"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Ticket, CircleCheck, User, Search, Plus, Present } from '@element-plus/icons-vue'
import {
  getCouponList, createCoupon, updateCoupon, deleteCoupon,
  getPromotionList, createPromotion, updatePromotion, deletePromotion,
  activatePromotion, deactivatePromotion, getCouponUsage,
  type Coupon, type Promotion, type CouponUsage,
  type CouponStatus, type CouponType, type PromotionStatus, type PromotionType
} from '@/api/promotion'
import TablePagination from '@/components/common/TablePagination.vue'

const router = useRouter()

// Tab state
const activeTab = ref('coupon')

// Stats
const stats = ref({
  totalCoupons: 0,
  activeCoupons: 0,
  totalUsed: 0,
  totalPromotions: 0
})

// Coupon list state
const couponList = ref<Coupon[]>([])
const couponLoading = ref(false)
const couponTotal = ref(0)
const couponParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as CouponStatus | undefined,
  type: undefined as CouponType | undefined
})

// Promotion list state
const promotionList = ref<Promotion[]>([])
const promotionLoading = ref(false)
const promotionTotal = ref(0)
const promotionParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as PromotionStatus | undefined,
  type: undefined as PromotionType | undefined
})

// Usage dialog state
const usageDialogVisible = ref(false)
const usageList = ref<CouponUsage[]>([])
const usageLoading = ref(false)
const usageTotal = ref(0)
const usageParams = reactive({
  id: 0,
  page: 1,
  page_size: 10
})

// Coupon form state
const couponDialogVisible = ref(false)
const isEditCoupon = ref(false)
const saveLoading = ref(false)
const couponFormRef = ref()
const couponForm = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  type: 'fixed_amount' as 'fixed_amount' | 'percentage',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  usage_limit: 0,
  per_user_limit: 0,
  dateRange: [] as string[]
})

const couponRules = {
  name: [{ required: true, message: '请输入优惠券名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入优惠码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择优惠类型', trigger: 'change' }],
  discount_value: [{ required: true, message: '请输入优惠值', trigger: 'blur' }],
  dateRange: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

// Promotion form state
const promotionDialogVisible = ref(false)
const isEditPromotion = ref(false)
const promotionFormRef = ref()
const promotionForm = reactive({
  id: 0,
  name: '',
  description: '',
  type: 'discount' as 'discount' | 'flash_sale' | 'bundle' | 'buy_x_get_y',
  discount_type: 'fixed_amount' as 'fixed_amount' | 'percentage' | 'buy_x_get_y',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  dateRange: [] as string[]
})

const promotionRules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择活动类型', trigger: 'change' }],
  discount_type: [{ required: true, message: '请选择优惠类型', trigger: 'change' }],
  dateRange: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

// Load functions
const loadCoupons = async () => {
  couponLoading.value = true
  try {
    const res = await getCouponList(couponParams)
    couponList.value = res.list || []
    couponTotal.value = res.total || 0
    // Update stats
    stats.value.totalCoupons = couponTotal.value
    stats.value.activeCoupons = couponList.value.filter(c => c.status === 'active').length
    stats.value.totalUsed = couponList.value.reduce((sum, c) => sum + c.used_count, 0)
  } catch (error) {
    console.error('Failed to load coupons:', error)
  } finally {
    couponLoading.value = false
  }
}

const loadPromotions = async () => {
  promotionLoading.value = true
  try {
    const res = await getPromotionList(promotionParams)
    promotionList.value = res.list || []
    promotionTotal.value = res.total || 0
    stats.value.totalPromotions = promotionTotal.value
  } catch (error) {
    console.error('Failed to load promotions:', error)
  } finally {
    promotionLoading.value = false
  }
}

const loadCouponUsage = async () => {
  usageLoading.value = true
  try {
    const res = await getCouponUsage(usageParams.id, { page: usageParams.page, page_size: usageParams.page_size })
    usageList.value = res.list || []
    usageTotal.value = res.total || 0
  } catch (error) {
    console.error('Failed to load coupon usage:', error)
  } finally {
    usageLoading.value = false
  }
}

// Helper functions
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

const getUsagePercentage = (row: Coupon) => {
  if (row.usage_limit === 0) return 0
  return Math.round((row.used_count / row.usage_limit) * 100)
}

const getProgressStatus = (row: Coupon) => {
  const percentage = getUsagePercentage(row)
  if (percentage >= 90) return 'exception'
  if (percentage >= 70) return 'warning'
  return ''
}

const getCouponStatusType = (status: string) => {
  const types: Record<string, string> = {
    'inactive': 'info',
    'active': 'success',
    'expired': 'warning',
    'depleted': 'danger'
  }
  return types[status] || 'info'
}

const getCouponStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'inactive': '未激活',
    'active': '已激活',
    'expired': '已过期',
    'depleted': '已用完'
  }
  return texts[status] || status
}

const getPromoStatusType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'paused': 'warning',
    'pending': 'info',
    'ended': 'info'
  }
  return types[status] || 'info'
}

const getPromoStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': '进行中',
    'paused': '已暂停',
    'pending': '待开始',
    'ended': '已结束'
  }
  return texts[status] || status
}

const getPromotionTypeText = (type: string) => {
  const texts: Record<string, string> = {
    'discount': '折扣',
    'flash_sale': '限时秒杀',
    'bundle': '捆绑销售',
    'buy_x_get_y': '买X送Y'
  }
  return texts[type] || type
}

const generateCouponCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 10; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  couponForm.code = code
}

// Page change handlers
const handleCouponPageChange = () => {
  loadCoupons()
}

const handlePromotionPageChange = () => {
  loadPromotions()
}

const handleUsagePageChange = () => {
  loadCouponUsage()
}

const handleTabChange = (tab: string) => {
  if (tab === 'coupon') {
    loadCoupons()
  } else if (tab === 'promotion') {
    loadPromotions()
  }
}

// Coupon actions
const handleAddCoupon = () => {
  isEditCoupon.value = false
  Object.assign(couponForm, {
    id: 0,
    name: '',
    code: '',
    description: '',
    type: 'fixed_amount',
    discount_value_num: 0,
    min_order_amount_num: 0,
    max_discount_num: 0,
    usage_limit: 100,
    per_user_limit: 0,
    dateRange: []
  })
  couponDialogVisible.value = true
}

const handleEditCoupon = (row: Coupon) => {
  isEditCoupon.value = true
  Object.assign(couponForm, {
    id: row.id,
    name: row.name,
    code: row.code,
    description: row.description,
    type: row.type,
    discount_value_num: parseFloat(row.discount_value) || 0,
    min_order_amount_num: parseFloat(row.min_order_amount) || 0,
    max_discount_num: parseFloat(row.max_discount) || 0,
    usage_limit: row.usage_limit,
    per_user_limit: row.per_user_limit,
    dateRange: [row.start_time, row.end_time]
  })
  couponDialogVisible.value = true
}

const handleCouponUsage = (row: Coupon) => {
  usageParams.id = row.id
  usageParams.page = 1
  usageDialogVisible.value = true
  loadCouponUsage()
}

const handleDeleteCoupon = async (row: Coupon) => {
  try {
    await ElMessageBox.confirm(`确定要删除优惠券 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteCoupon(row.id)
    ElMessage.success('删除成功')
    loadCoupons()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete coupon:', error)
    }
  }
}

const saveCoupon = async () => {
  if (!couponFormRef.value) return

  await couponFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saveLoading.value = true
    try {
      const data = {
        name: couponForm.name,
        code: couponForm.code,
        description: couponForm.description,
        type: couponForm.type,
        discount_value: String(couponForm.discount_value_num),
        min_order_amount: String(couponForm.min_order_amount_num),
        max_discount: String(couponForm.max_discount_num),
        usage_limit: couponForm.usage_limit,
        per_user_limit: couponForm.per_user_limit,
        start_time: couponForm.dateRange[0],
        end_time: couponForm.dateRange[1]
      }

      if (isEditCoupon.value) {
        await updateCoupon({ id: couponForm.id, ...data })
      } else {
        await createCoupon(data)
      }

      ElMessage.success(isEditCoupon.value ? '编辑成功' : '创建成功')
      couponDialogVisible.value = false
      loadCoupons()
    } catch (error) {
      console.error('Failed to save coupon:', error)
    } finally {
      saveLoading.value = false
    }
  })
}

// Promotion actions
const handleAddPromotion = () => {
  isEditPromotion.value = false
  Object.assign(promotionForm, {
    id: 0,
    name: '',
    description: '',
    type: 'discount',
    discount_type: 'fixed_amount',
    discount_value_num: 0,
    min_order_amount_num: 0,
    max_discount_num: 0,
    dateRange: []
  })
  promotionDialogVisible.value = true
}

const handleEditPromotion = (row: Promotion) => {
  router.push(`/promotions/${row.id}`)
}

const handleActivatePromotion = async (row: Promotion) => {
  try {
    await activatePromotion(row.id)
    ElMessage.success('激活成功')
    loadPromotions()
  } catch (error) {
    console.error('Failed to activate promotion:', error)
  }
}

const handleDeactivatePromotion = async (row: Promotion) => {
  try {
    await deactivatePromotion(row.id)
    ElMessage.success('停用成功')
    loadPromotions()
  } catch (error) {
    console.error('Failed to deactivate promotion:', error)
  }
}

const handleDeletePromotion = async (row: Promotion) => {
  try {
    await ElMessageBox.confirm(`确定要删除促销活动 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deletePromotion(row.id)
    ElMessage.success('删除成功')
    loadPromotions()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete promotion:', error)
    }
  }
}

const savePromotion = async () => {
  if (!promotionFormRef.value) return

  await promotionFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saveLoading.value = true
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
        end_time: promotionForm.dateRange[1]
      }

      if (isEditPromotion.value) {
        await updatePromotion({ id: promotionForm.id, ...data })
      } else {
        await createPromotion(data)
      }

      ElMessage.success(isEditPromotion.value ? '编辑成功' : '创建成功')
      promotionDialogVisible.value = false
      loadPromotions()
    } catch (error) {
      console.error('Failed to save promotion:', error)
    } finally {
      saveLoading.value = false
    }
  })
}

// Initialize
onMounted(() => {
  loadCoupons()
  loadPromotions()
})
</script>

<style scoped>
.promotions-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #6B7280;
}

.stat-icon.active {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.stat-icon.used {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  color: white;
}

.stat-icon.promotions {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Tabs Card */
.tabs-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Tabs */
:deep(.el-tabs__item.is-active) {
  color: #6366F1;
  font-weight: 600;
}

:deep(.el-tabs__active-bar) {
  background-color: #6366F1;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.tab-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 220px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

/* Coupon Cell */
.coupon-cell, .promo-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.coupon-icon, .promo-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F5F3FF;
  color: #6366F1;
}

.coupon-icon.fixed_amount, .promo-icon.discount {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.coupon-icon.percentage, .promo-icon.full_reduce {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.coupon-details, .promo-details {
  flex: 1;
  min-width: 0;
}

.coupon-name, .promo-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.coupon-code {
  font-size: 12px;
  color: #6B7280;
  margin: 0 0 8px 0;
  font-family: 'Fira Code', monospace;
}

.coupon-tags {
  display: flex;
  gap: 8px;
}

.promo-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Discount Value */
.discount-value {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.min-order {
  font-size: 12px;
  color: #6B7280;
  margin-top: 4px;
}

/* Usage Stats */
.usage-stats {
  text-align: center;
}

.usage-text {
  font-size: 12px;
  color: #6B7280;
  margin: 8px 0 0 0;
}

/* Validity Period */
.validity-period {
  text-align: center;
  font-size: 13px;
  color: #4B5563;
}

.to-text {
  margin: 2px 0;
  color: #9CA3AF;
  font-size: 12px;
}

/* Form */
.form-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

/* Dialog */
:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Progress */
:deep(.el-progress-bar__inner) {
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
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
  .tab-header {
    flex-direction: column;
    align-items: stretch;
  }

  .tab-filters {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>