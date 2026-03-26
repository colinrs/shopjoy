<template>
  <div class="shop-settings" v-loading="loading">
    <el-row :gutter="20">
      <!-- Left Column - Basic Info -->
      <el-col :xs="24" :lg="16">
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Shop /></el-icon>
                店铺基本信息
              </span>
            </div>
          </template>

          <el-form :model="shopForm" label-width="120px" :rules="shopRules" ref="shopFormRef">
            <el-form-item label="店铺Logo">
              <el-upload
                class="logo-uploader"
                action="#"
                :show-file-list="false"
                :auto-upload="false"
                :on-change="handleLogoChange"
              >
                <img v-if="shopForm.logo" :src="shopForm.logo" class="logo-preview" />
                <div v-else class="logo-placeholder">
                  <el-icon size="32"><Plus /></el-icon>
                  <span>上传Logo</span>
                </div>
              </el-upload>
              <span class="upload-tip">建议尺寸 200x200px，支持 JPG、PNG 格式</span>
            </el-form-item>

            <el-form-item label="店铺名称" prop="name">
              <el-input v-model="shopForm.name" placeholder="请输入店铺名称" maxlength="50" show-word-limit />
            </el-form-item>

            <el-form-item label="店铺简介" prop="description">
              <el-input
                v-model="shopForm.description"
                type="textarea"
                :rows="4"
                placeholder="请输入店铺简介，这将展示在店铺首页"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item label="联系人" prop="contact_name">
              <el-input v-model="shopForm.contact_name" placeholder="请输入联系人姓名" />
            </el-form-item>

            <el-form-item label="客服电话" prop="contact_phone">
              <el-input v-model="shopForm.contact_phone" placeholder="请输入客服电话">
                <template #prefix>
                  <el-icon><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="客服邮箱" prop="contact_email">
              <el-input v-model="shopForm.contact_email" placeholder="请输入客服邮箱">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="店铺地址">
              <el-input v-model="shopForm.address" placeholder="请输入店铺地址">
                <template #prefix>
                  <el-icon><Location /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="自定义域名" prop="custom_domain">
              <el-input v-model="shopForm.custom_domain" placeholder="请输入自定义域名，如 shop.example.com">
                <template #prefix>
                  <el-icon><Link /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Business Hours Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Clock /></el-icon>
                营业时间设置
              </span>
            </div>
          </template>

          <div class="business-hours">
            <div v-for="(hour, index) in businessHours" :key="index" class="hours-row">
              <span class="day-name">{{ getDayName(hour.day_of_week) }}</span>
              <el-switch
                v-model="hour.is_closed"
                :active-text="'营业'"
                :inactive-text="'休息'"
                :active-value="false"
                :inactive-value="true"
                style="margin: 0 16px"
              />
              <template v-if="!hour.is_closed">
                <el-time-select
                  v-model="hour.open_time"
                  :max-time="hour.close_time"
                  placeholder="开始时间"
                  start="00:00"
                  step="00:30"
                  end="23:30"
                  style="width: 120px"
                />
                <span style="margin: 0 8px">至</span>
                <el-time-select
                  v-model="hour.close_time"
                  :min-time="hour.open_time"
                  placeholder="结束时间"
                  start="00:00"
                  step="00:30"
                  end="23:30"
                  style="width: 120px"
                />
              </template>
            </div>
          </div>
        </el-card>

        <!-- Shipping Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Van /></el-icon>
                运费设置
              </span>
            </div>
          </template>

          <el-form :model="shippingForm" label-width="120px">
            <el-form-item label="默认运费">
              <el-input-number
                v-model="shippingForm.default_shipping_fee"
                :precision="2"
                :min="0"
                placeholder="默认运费"
                style="width: 200px"
              />
              <span style="margin-left: 8px; color: #909399">{{ shippingForm.currency || 'CNY' }}</span>
            </el-form-item>

            <el-form-item label="包邮门槛">
              <el-input-number
                v-model="shippingForm.free_shipping_threshold"
                :precision="2"
                :min="0"
                placeholder="满额包邮"
                style="width: 200px"
              />
              <span style="margin-left: 8px; color: #909399">{{ shippingForm.currency || 'CNY' }}</span>
              <div class="form-tip">订单金额达到此数值时免运费，0 表示不设包邮门槛</div>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- Right Column - Status & Settings -->
      <el-col :xs="24" :lg="8">
        <!-- Shop Status -->
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><InfoFilled /></el-icon>
                店铺状态
              </span>
            </div>
          </template>

          <div class="status-info">
            <div class="info-row">
              <span class="info-label">店铺状态</span>
              <el-tag :type="shopSettings?.status === 1 ? 'success' : 'info'">
                {{ shopSettings?.status_text || '-' }}
              </el-tag>
            </div>
            <div class="info-row">
              <span class="info-label">当前套餐</span>
              <el-tag type="warning">{{ shopSettings?.plan_text || '-' }}</el-tag>
            </div>
            <div class="info-row" v-if="shopSettings?.expire_at">
              <span class="info-label">到期时间</span>
              <span class="info-value">{{ shopSettings?.expire_at }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">默认货币</span>
              <span class="info-value">{{ shopSettings?.default_currency || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">时区</span>
              <span class="info-value">{{ shopSettings?.timezone || '-' }}</span>
            </div>
          </div>
        </el-card>

        <!-- Notification Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Bell /></el-icon>
                通知设置
              </span>
            </div>
          </template>

          <el-form :model="notificationForm" label-position="top">
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_created">新订单通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_paid">订单支付通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_shipped">订单发货通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_cancelled">订单取消通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.refund_requested">退款申请通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.new_review">新评价通知</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.low_stock_alert">低库存预警</el-checkbox>
            </div>
            <el-form-item label="库存预警阈值" v-if="notificationForm.low_stock_alert" style="margin-top: 12px">
              <el-input-number v-model="notificationForm.low_stock_threshold" :min="1" :max="1000" style="width: 100%" />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Payment Settings -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><CreditCard /></el-icon>
                支付设置
              </span>
            </div>
          </template>

          <el-form :model="paymentForm" label-position="top">
            <div class="notify-item">
              <el-checkbox v-model="paymentForm.stripe_enabled">启用 Stripe 支付</el-checkbox>
            </div>
            <el-form-item label="Stripe Public Key" v-if="paymentForm.stripe_enabled" style="margin-top: 12px">
              <el-input v-model="paymentForm.stripe_public_key" placeholder="pk_live_xxx 或 pk_test_xxx" />
            </el-form-item>
            <el-form-item label="Stripe Secret Key" v-if="paymentForm.stripe_enabled">
              <el-input
                v-model="paymentForm.stripe_secret_key"
                type="password"
                placeholder="留空表示不修改现有密钥"
                show-password
                @change="paymentForm.stripe_secret_key_changed = true"
              />
              <div class="form-tip">密钥将被加密存储，留空则保持现有密钥不变</div>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Quick Actions -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Operation /></el-icon>
                快捷操作
              </span>
            </div>
          </template>

          <div class="quick-actions">
            <el-button plain style="width: 100%; margin-bottom: 12px" @click="viewShopPage">
              <el-icon><Document /></el-icon>
              查看店铺主页
            </el-button>
            <el-button plain style="width: 100%; margin-bottom: 12px" @click="copyShopLink">
              <el-icon><Share /></el-icon>
              复制店铺链接
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Save Button -->
    <div class="save-bar">
      <el-button type="primary" size="large" @click="handleSave" :loading="saving">
        <el-icon><Check /></el-icon>
        保存设置
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Shop, Plus, Phone, Message, Location, Link,
  InfoFilled, Bell, Operation, Document, Share, Check,
  Clock, Van, CreditCard
} from '@element-plus/icons-vue'
import { uploadImage } from '@/api/upload'
import {
  getShopSettings,
  updateShopSettings,
  getBusinessHours,
  updateBusinessHours,
  getNotificationSettings,
  updateNotificationSettings,
  getPaymentSettings,
  updatePaymentSettings,
  getShippingSettings,
  updateShippingSettings,
  type ShopSettings,
  type BusinessHours,
  type NotificationSettings,
  type PaymentSettings,
  type ShippingSettings
} from '@/api/shop'

const loading = ref(false)
const saving = ref(false)
const shopFormRef = ref()

// Shop settings data
const shopSettings = ref<ShopSettings | null>(null)

const shopForm = reactive({
  name: '',
  description: '',
  logo: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  address: '',
  custom_domain: ''
})

const businessHours = ref<BusinessHours[]>([])

const notificationForm = reactive<NotificationSettings>({
  order_created: true,
  order_paid: true,
  order_shipped: true,
  order_cancelled: true,
  low_stock_alert: true,
  low_stock_threshold: 10,
  refund_requested: true,
  new_review: true
})

const paymentForm = reactive({
  stripe_enabled: false,
  stripe_public_key: '',
  stripe_secret_key: '',
  stripe_secret_key_changed: false  // Track if secret key was modified
})

// Shipping form uses numbers for el-input-number, converts to/from strings for API
const shippingForm = reactive({
  free_shipping_threshold: 0,  // Number for UI
  default_shipping_fee: 0,      // Number for UI
  currency: 'CNY'
})

const shopRules = {
  name: [
    { required: true, message: '请输入店铺名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '最多 500 个字符', trigger: 'blur' }
  ],
  contact_phone: [
    { pattern: /^$|^1[3-9]\d{9}$|^400-\d{3}-\d{4}$/, message: '请输入正确的电话号码', trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  custom_domain: [
    {
      pattern: /^$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/,
      message: '请输入有效的域名，如 shop.example.com',
      trigger: 'blur'
    }
  ]
}

// Day names mapping
const dayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

const getDayName = (dayOfWeek: number): string => {
  return dayNames[dayOfWeek] || ''
}

// Load all settings
const loadSettings = async () => {
  loading.value = true
  try {
    // Load shop settings
    const shopData = await getShopSettings()
    shopSettings.value = shopData
    shopForm.name = shopData.name || ''
    shopForm.description = shopData.description || ''
    shopForm.logo = shopData.logo || ''
    shopForm.contact_name = shopData.contact_name || ''
    shopForm.contact_phone = shopData.contact_phone || ''
    shopForm.contact_email = shopData.contact_email || ''
    shopForm.address = shopData.address || ''
    shopForm.custom_domain = shopData.custom_domain || ''

    // Load business hours - sort by day of week
    const hoursData = await getBusinessHours()
    if (hoursData && hoursData.length > 0) {
      businessHours.value = hoursData.sort((a, b) => a.day_of_week - b.day_of_week)
    } else {
      // Initialize default business hours
      businessHours.value = dayNames.map((_, index) => ({
        day_of_week: index,
        open_time: '09:00',
        close_time: '18:00',
        is_closed: index === 0 // Sunday closed by default
      }))
    }

    // Load notification settings
    const notifyData = await getNotificationSettings()
    Object.assign(notificationForm, notifyData)

    // Load payment settings - secret key not returned for security
    const paymentData = await getPaymentSettings()
    paymentForm.stripe_enabled = paymentData.stripe_enabled
    paymentForm.stripe_public_key = paymentData.stripe_public_key || ''
    paymentForm.stripe_secret_key = ''  // Never returned from backend
    paymentForm.stripe_secret_key_changed = false

    // Load shipping settings - convert strings to numbers for el-input-number
    const shippingData = await getShippingSettings()
    shippingForm.free_shipping_threshold = parseFloat(shippingData.free_shipping_threshold) || 0
    shippingForm.default_shipping_fee = parseFloat(shippingData.default_shipping_fee) || 0
    shippingForm.currency = shippingData.currency || 'CNY'
  } catch (error) {
    console.error('Failed to load settings:', error)
    ElMessage.error('加载店铺设置失败')
  } finally {
    loading.value = false
  }
}

// Handle logo upload
const handleLogoChange = async (file: any) => {
  try {
    const response = await uploadImage(file.raw, 'banner')
    shopForm.logo = response.url
    ElMessage.success('Logo上传成功')
  } catch (error) {
    console.error('Upload failed:', error)
    ElMessage.error('Logo上传失败')
  }
}

// View shop page
const viewShopPage = () => {
  if (shopSettings.value?.domain) {
    window.open(`https://${shopSettings.value.domain}`, '_blank')
  } else if (shopSettings.value?.custom_domain) {
    window.open(`https://${shopSettings.value.custom_domain}`, '_blank')
  } else {
    ElMessage.warning('店铺域名未配置')
  }
}

// Copy shop link
const copyShopLink = async () => {
  const link = shopSettings.value?.custom_domain || shopSettings.value?.domain
  if (link) {
    try {
      await navigator.clipboard.writeText(`https://${link}`)
      ElMessage.success('店铺链接已复制到剪贴板')
    } catch {
      ElMessage.error('复制失败')
    }
  } else {
    ElMessage.warning('店铺域名未配置')
  }
}

// Validate business hours
const validateBusinessHours = (hours: BusinessHours[]): boolean => {
  for (const hour of hours) {
    if (!hour.is_closed) {
      if (!hour.open_time || !hour.close_time) {
        ElMessage.error(`${getDayName(hour.day_of_week)} 营业时间未设置`)
        return false
      }
      if (hour.open_time >= hour.close_time) {
        ElMessage.error(`${getDayName(hour.day_of_week)} 结束时间必须晚于开始时间`)
        return false
      }
    }
  }
  return true
}

// Save all settings
const handleSave = async () => {
  if (!shopFormRef.value) return

  try {
    await shopFormRef.value.validate()
  } catch {
    ElMessage.warning('请检查表单填写是否正确')
    return
  }

  // Validate business hours
  if (!validateBusinessHours(businessHours.value)) {
    return
  }

  saving.value = true
  try {
    // Save shop settings
    // Sanitize custom domain: remove protocol and trim whitespace
    const sanitizedDomain = shopForm.custom_domain
      ? shopForm.custom_domain.replace(/^https?:\/\//, '').trim()
      : undefined

    await updateShopSettings({
      name: shopForm.name,
      description: shopForm.description,
      logo: shopForm.logo || undefined,
      contact_name: shopForm.contact_name || undefined,
      contact_phone: shopForm.contact_phone || undefined,
      contact_email: shopForm.contact_email || undefined,
      address: shopForm.address || undefined,
      custom_domain: sanitizedDomain
    })

    // Save business hours
    await updateBusinessHours({ hours: businessHours.value })

    // Save notification settings
    await updateNotificationSettings(notificationForm)

    // Save payment settings - only send secret key if it was changed
    await updatePaymentSettings({
      stripe_enabled: paymentForm.stripe_enabled,
      stripe_secret_key: paymentForm.stripe_secret_key_changed ? paymentForm.stripe_secret_key : undefined
    })

    // Save shipping settings - convert numbers to strings for API
    await updateShippingSettings({
      free_shipping_threshold: shippingForm.free_shipping_threshold.toString(),
      default_shipping_fee: shippingForm.default_shipping_fee.toString()
    })

    ElMessage.success('店铺设置保存成功')

    // Reset secret key change flag
    paymentForm.stripe_secret_key_changed = false

    // Reload settings to get updated values
    await loadSettings()
  } catch (error: unknown) {
    console.error('Failed to save settings:', error)
    const err = error as { response?: { data?: { message?: string } } }
    ElMessage.error(err?.response?.data?.message || '保存失败，请重试')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.shop-settings {
  padding: 0;
}

.settings-card {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Logo Uploader */
.logo-uploader {
  border: 2px dashed #D1D5DB;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
  width: 120px;
  height: 120px;
}

.logo-uploader:hover {
  border-color: #059669;
}

.logo-preview {
  width: 120px;
  height: 120px;
  object-fit: cover;
  display: block;
}

.logo-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9CA3AF;
  gap: 8px;
}

.logo-placeholder span {
  font-size: 12px;
}

.upload-tip {
  margin-left: 16px;
  font-size: 12px;
  color: #9CA3AF;
}

/* Business Hours */
.business-hours {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hours-row {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.hours-row:last-child {
  border-bottom: none;
}

.day-name {
  width: 60px;
  font-weight: 500;
  color: #374151;
}

/* Status Info */
.status-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: #6B7280;
}

.info-value {
  font-size: 14px;
  color: #111827;
}

/* Notify Items */
.notify-item {
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.notify-item:last-child {
  border-bottom: none;
}

/* Form Tips */
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

/* Quick Actions */
.quick-actions {
  display: flex;
  flex-direction: column;
}

.quick-actions .el-button {
  justify-content: flex-start;
  gap: 8px;
}

/* Save Bar */
.save-bar {
  position: fixed;
  bottom: 0;
  left: 220px;
  right: 0;
  padding: 16px 24px;
  background: #fff;
  border-top: 1px solid #E5E7EB;
  display: flex;
  justify-content: center;
  z-index: 100;
}

.save-bar .el-button {
  min-width: 160px;
}

/* Responsive */
@media (max-width: 768px) {
  .save-bar {
    left: 0;
  }
}
</style>