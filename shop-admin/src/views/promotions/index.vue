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
            <p class="stat-label">已领取</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <el-icon class="stat-icon savings" size="32"><Money /></el-icon>
          <div class="stat-info">
            <p class="stat-value">¥{{ formatNumber(stats.totalSavings) }}</p>
            <p class="stat-label">累计优惠</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-card class="tabs-card" shadow="never">
      <el-tabs v-model="activeTab" class="promotion-tabs">
        <!-- Coupons Tab -->
        <el-tab-pane label="优惠券管理" name="coupon">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="couponSearch"
                placeholder="搜索优惠券"
                clearable
                class="search-input"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="couponStatus" placeholder="状态" clearable class="filter-select">
                <el-option label="全部" value="" />
                <el-option label="未开始" value="pending" />
                <el-option label="进行中" value="active" />
                <el-option label="已结束" value="ended" />
              </el-select>
            </div>
            <el-button type="primary" @click="handleAddCoupon">
              <el-icon><Plus /></el-icon>创建优惠券
            </el-button>
          </div>

          <el-table :data="couponList" v-loading="loading" stripe>
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
                      <el-tag size="small" :type="row.type === 'fixed' ? 'success' : 'warning'">
                        {{ row.type === 'fixed' ? '固定金额' : '百分比' }}
                      </el-tag>
                      <el-tag v-if="row.is_exclusive" size="small" type="danger" effect="plain">专属</el-tag>
                    </div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="优惠力度" width="150" align="center">
              <template #default="{ row }">
                <div class="discount-value">
                  <span v-if="row.type === 'fixed'">¥{{ row.value }}</span>
                  <span v-else>{{ row.value }}%</span>
                </div>
                <div class="min-order" v-if="row.min_order">满¥{{ row.min_order }}可用</div>
              </template>
            </el-table-column>
            <el-table-column label="使用统计" width="180" align="center">
              <template #default="{ row }">
                <div class="usage-stats">
                  <el-progress 
                    :percentage="Math.round(row.used / row.total * 100)" 
                    :status="getProgressStatus(row)"
                    :stroke-width="8"
                  />
                  <p class="usage-text">{{ row.used }}/{{ row.total }} 已使用</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="有效期" width="200">
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ row.start_date }}</p>
                  <p class="to-text">至</p>
                  <p>{{ row.end_date }}</p>
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
                <el-button type="primary" link size="small" @click="handleStats(row)">
                  数据
                </el-button>
                <el-button 
                  v-if="row.status !== 'ended'" 
                  type="danger" 
                  link 
                  size="small" 
                  @click="handleEndCoupon(row)"
                >
                  结束
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Promotions Tab -->
        <el-tab-pane label="促销活动" name="promotion">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="promoSearch"
                placeholder="搜索活动"
                clearable
                class="search-input"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
            <el-button type="primary" @click="handleAddPromotion">
              <el-icon><Plus /></el-icon>创建活动
            </el-button>
          </div>

          <div class="promotion-list">
            <el-empty v-if="promotionList.length === 0" description="暂无促销活动" />
            <div v-else class="promotion-grid">
              <el-card 
                v-for="promo in promotionList" 
                :key="promo.id" 
                class="promotion-card"
                shadow="hover"
              >
                <div class="promo-header">
                  <div class="promo-icon" :class="promo.type">
                    <el-icon size="28"><Present /></el-icon>
                  </div>
                  <el-tag :type="getPromoStatusType(promo.status)" size="small">
                    {{ getPromoStatusText(promo.status) }}
                  </el-tag>
                </div>
                <h4 class="promo-name">{{ promo.name }}</h4>
                <p class="promo-desc">{{ promo.description }}</p>
                <div class="promo-meta">
                  <span><el-icon><Timer /></el-icon> {{ promo.start_date }}</span>
                  <span><el-icon><User /></el-icon> {{ promo.participants }} 人参与</span>
                </div>
                <div class="promo-actions">
                  <el-button type="primary" link size="small" @click="handleEditPromotion(promo)">
                    编辑
                  </el-button>
                  <el-button type="primary" link size="small" @click="handlePromoStats(promo)">
                    数据
                  </el-button>
                </div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>

        <!-- Flash Sale Tab -->
        <el-tab-pane label="限时秒杀" name="flash">
          <el-empty description="限时秒杀功能开发中...">
            <el-button type="primary" @click="handleCreateFlashSale">创建秒杀活动</el-button>
          </el-empty>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Add Coupon Dialog -->
    <el-dialog v-model="couponDialogVisible" :title="isEdit ? '编辑优惠券' : '创建优惠券'" width="600px">
      <el-form :model="couponForm" label-width="100px" :rules="couponRules" ref="couponFormRef">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="couponForm.name" placeholder="例如: 新用户专享券" />
        </el-form-item>
        <el-form-item label="优惠码" prop="code">
          <el-input v-model="couponForm.code" placeholder="例如: NEWUSER2024">
            <template #append>
              <el-button @click="generateCode">生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="优惠类型" prop="type">
          <el-radio-group v-model="couponForm.type">
            <el-radio label="fixed">固定金额</el-radio>
            <el-radio label="percentage">百分比折扣</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="couponForm.type === 'fixed' ? '优惠金额' : '折扣比例'" prop="value">
          <el-input-number v-model="couponForm.value" :min="0" :max="couponForm.type === 'percentage' ? 100 : 99999" style="width: 100%" />
        </el-form-item>
        <el-form-item label="最低消费">
          <el-input-number v-model="couponForm.min_order" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="发放总量" prop="total">
          <el-input-number v-model="couponForm.total" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="有效期" prop="dateRange">
          <el-date-picker
            v-model="couponForm.dateRange"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="couponDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCoupon" :loading="saveLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Ticket, CircleCheck, User, Money, Search, Plus, Present, Timer } from '@element-plus/icons-vue'

const activeTab = ref('coupon')
const loading = ref(false)
const couponSearch = ref('')
const couponStatus = ref('')
const promoSearch = ref('')
const couponDialogVisible = ref(false)
const isEdit = ref(false)
const saveLoading = ref(false)
const couponFormRef = ref()

const stats = ref({
  totalCoupons: 24,
  activeCoupons: 8,
  totalUsed: 1586,
  totalSavings: 45680
})

const couponForm = reactive({
  name: '',
  code: '',
  type: 'fixed',
  value: 0,
  min_order: 0,
  total: 100,
  dateRange: []
})

const couponRules = {
  name: [{ required: true, message: '请输入优惠券名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入优惠码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择优惠类型', trigger: 'change' }],
  value: [{ required: true, message: '请输入优惠值', trigger: 'blur' }],
  total: [{ required: true, message: '请输入发放总量', trigger: 'blur' }],
  dateRange: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const couponList = ref([
  {
    id: 1,
    name: '新用户专享券',
    code: 'NEWUSER2024',
    type: 'fixed',
    value: 50,
    min_order: 100,
    total: 1000,
    used: 328,
    start_date: '2024-03-01',
    end_date: '2024-04-01',
    status: 'active',
    is_exclusive: true
  },
  {
    id: 2,
    name: '满200减30',
    code: 'SAVE30',
    type: 'fixed',
    value: 30,
    min_order: 200,
    total: 500,
    used: 156,
    start_date: '2024-03-15',
    end_date: '2024-03-31',
    status: 'active',
    is_exclusive: false
  },
  {
    id: 3,
    name: '全场85折',
    code: 'SALE15',
    type: 'percentage',
    value: 15,
    min_order: 0,
    total: 2000,
    used: 892,
    start_date: '2024-02-01',
    end_date: '2024-02-28',
    status: 'ended',
    is_exclusive: false
  }
])

const promotionList = ref([
  {
    id: 1,
    name: '春季大促',
    description: '全场满299减50，更有惊喜好礼相送',
    type: 'seasonal',
    status: 'active',
    start_date: '2024-03-01',
    end_date: '2024-03-31',
    participants: 1256
  },
  {
    id: 2,
    name: '新品尝鲜',
    description: '新品上市限时特惠，立享9折优惠',
    type: 'new_product',
    status: 'active',
    start_date: '2024-03-15',
    end_date: '2024-04-15',
    participants: 568
  }
])

const formatNumber = (num: number) => {
  return num.toLocaleString()
}

const getProgressStatus = (row: any) => {
  const percentage = row.used / row.total
  if (percentage >= 0.9) return 'exception'
  if (percentage >= 0.7) return 'warning'
  return ''
}

const getCouponStatusType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'info',
    'active': 'success',
    'ended': 'info',
    'paused': 'warning'
  }
  return types[status] || 'info'
}

const getCouponStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'pending': '未开始',
    'active': '进行中',
    'ended': '已结束',
    'paused': '已暂停'
  }
  return texts[status] || status
}

const getPromoStatusType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'ended': 'info',
    'pending': 'warning'
  }
  return types[status] || 'info'
}

const getPromoStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': '进行中',
    'ended': '已结束',
    'pending': '未开始'
  }
  return texts[status] || status
}

const generateCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 10; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  couponForm.code = code
}

const handleAddCoupon = () => {
  isEdit.value = false
  Object.assign(couponForm, {
    name: '',
    code: '',
    type: 'fixed',
    value: 0,
    min_order: 0,
    total: 100,
    dateRange: []
  })
  couponDialogVisible.value = true
}

const handleEditCoupon = (row: any) => {
  isEdit.value = true
  Object.assign(couponForm, { ...row })
  couponDialogVisible.value = true
}

const handleStats = (row: any) => {
  ElMessage.info('查看优惠券数据: ' + row.name)
}

const handleEndCoupon = (row: any) => {
  ElMessage.success('优惠券已结束: ' + row.name)
}

const saveCoupon = async () => {
  if (!couponFormRef.value) return
  
  await couponFormRef.value.validate((valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      setTimeout(() => {
        saveLoading.value = false
        couponDialogVisible.value = false
        ElMessage.success(isEdit.value ? '编辑成功' : '创建成功')
      }, 1000)
    }
  })
}

const handleAddPromotion = () => {
  ElMessage.info('创建促销活动')
}

const handleEditPromotion = (promo: any) => {
  ElMessage.info('编辑活动: ' + promo.name)
}

const handlePromoStats = (promo: any) => {
  ElMessage.info('查看活动数据: ' + promo.name)
}

const handleCreateFlashSale = () => {
  ElMessage.info('创建秒杀活动')
}
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

.stat-icon.savings {
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
  font-family: 'Fira Sans', sans-serif;
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
.coupon-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.coupon-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F5F3FF;
  color: #6366F1;
}

.coupon-icon.fixed {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.coupon-icon.percentage {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.coupon-details {
  flex: 1;
  min-width: 0;
}

.coupon-name {
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

/* Discount Value */
.discount-value {
  font-size: 20px;
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

/* Promotion Grid */
.promotion-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.promotion-card {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.promotion-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px -8px rgba(99, 102, 241, 0.15);
}

.promo-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.promo-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F5F3FF;
  color: #6366F1;
}

.promo-icon.seasonal {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.promo-icon.new_product {
  background: linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%);
  color: white;
}

.promo-name {
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 8px 0;
}

.promo-desc {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.promo-meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: #9CA3AF;
  margin-bottom: 16px;
}

.promo-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.promo-actions {
  display: flex;
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid #F3F4F6;
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
  
  .promotion-grid {
    grid-template-columns: 1fr;
  }
}
</style>
