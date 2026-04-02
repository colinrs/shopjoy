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
                {{ $t('shop.basicInfo') }}
              </span>
            </div>
          </template>

          <el-form :model="shopForm" label-width="120px" :rules="shopRules" ref="shopFormRef">
            <el-form-item :label="$t('shop.shopLogo')">
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
                  <span>{{ $t('shop.uploadLogo') }}</span>
                </div>
              </el-upload>
              <span class="upload-tip">{{ $t('shop.logoTip') }}</span>
            </el-form-item>

            <el-form-item :label="$t('shop.shopName')" prop="name">
              <el-input v-model="shopForm.name" :placeholder="$t('shop.shopNamePlaceholder')" maxlength="50" show-word-limit />
            </el-form-item>

            <el-form-item :label="$t('shop.shopDescription')" prop="description">
              <el-input
                v-model="shopForm.description"
                type="textarea"
                :rows="4"
                :placeholder="$t('shop.shopDescPlaceholder')"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item :label="$t('shop.contactPerson')" prop="contact_name">
              <el-input v-model="shopForm.contact_name" :placeholder="$t('shop.contactPersonPlaceholder')" />
            </el-form-item>

            <el-form-item :label="$t('shop.customerServicePhone')" prop="contact_phone">
              <el-input v-model="shopForm.contact_phone" :placeholder="$t('shop.phonePlaceholder')">
                <template #prefix>
                  <el-icon><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item :label="$t('shop.customerServiceEmail')" prop="contact_email">
              <el-input v-model="shopForm.contact_email" :placeholder="$t('shop.emailPlaceholder')">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item :label="$t('shop.shopAddress')">
              <el-input v-model="shopForm.address" :placeholder="$t('shop.addressPlaceholder')">
                <template #prefix>
                  <el-icon><Location /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item :label="$t('shop.customDomain')" prop="custom_domain">
              <el-input v-model="shopForm.custom_domain" :placeholder="$t('shop.domainPlaceholder')">
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
                {{ $t('shop.businessHoursSettings') }}
              </span>
            </div>
          </template>

          <div class="business-hours">
            <div v-for="(hour, index) in businessHours" :key="index" class="hours-row">
              <span class="day-name">{{ getDayName(hour.day_of_week) }}</span>
              <el-switch
                v-model="hour.is_closed"
                :active-text="$t('shop.open')"
                :inactive-text="$t('shop.closed')"
                :active-value="false"
                :inactive-value="true"
                style="margin: 0 16px"
              />
              <template v-if="!hour.is_closed">
                <el-time-select
                  v-model="hour.open_time"
                  :max-time="hour.close_time"
                  :placeholder="$t('shop.startTime')"
                  start="00:00"
                  step="00:30"
                  end="23:30"
                  style="width: 120px"
                />
                <span style="margin: 0 8px">{{ $t('shop.to') }}</span>
                <el-time-select
                  v-model="hour.close_time"
                  :min-time="hour.open_time"
                  :placeholder="$t('shop.endTime')"
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
                {{ $t('shop.shippingSettings') }}
              </span>
            </div>
          </template>

          <el-form :model="shippingForm" label-width="120px">
            <el-form-item :label="$t('shop.defaultShippingFee')">
              <el-input-number
                v-model="shippingForm.default_shipping_fee"
                :precision="2"
                :min="0"
                :placeholder="$t('shop.defaultShippingFee')"
                style="width: 200px"
              />
              <span style="margin-left: 8px; color: #909399">{{ shippingForm.currency || 'CNY' }}</span>
            </el-form-item>

            <el-form-item :label="$t('shop.freeShippingThreshold')">
              <el-input-number
                v-model="shippingForm.free_shipping_threshold"
                :precision="2"
                :min="0"
                :placeholder="$t('shop.freeShippingThreshold')"
                style="width: 200px"
              />
              <span style="margin-left: 8px; color: #909399">{{ shippingForm.currency || 'CNY' }}</span>
              <div class="form-tip">{{ $t('shop.freeShippingThresholdTip') }}</div>
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
                {{ $t('shop.shopStatus') }}
              </span>
            </div>
          </template>

          <div class="status-info">
            <div class="info-row">
              <span class="info-label">{{ $t('shop.shopStatus') }}</span>
              <el-tag :type="shopSettings?.status === 1 ? 'success' : 'info'">
                {{ shopSettings?.status_text || '-' }}
              </el-tag>
            </div>
            <div class="info-row">
              <span class="info-label">{{ $t('shop.currentPlan') }}</span>
              <el-tag type="warning">{{ shopSettings?.plan_text || '-' }}</el-tag>
            </div>
            <div class="info-row" v-if="shopSettings?.expire_at">
              <span class="info-label">{{ $t('shop.expireTime') }}</span>
              <span class="info-value">{{ shopSettings?.expire_at }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ $t('shop.defaultCurrency') }}</span>
              <span class="info-value">{{ shopSettings?.default_currency || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ $t('shop.timezone') }}</span>
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
                {{ $t('shop.notificationSettings') }}
              </span>
            </div>
          </template>

          <el-form :model="notificationForm" label-position="top">
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_created">{{ $t('shop.newOrderNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_paid">{{ $t('shop.orderPaidNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_shipped">{{ $t('shop.orderShippedNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.order_cancelled">{{ $t('shop.orderCancelledNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.refund_requested">{{ $t('shop.refundRequestNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.new_review">{{ $t('shop.newReviewNotification') }}</el-checkbox>
            </div>
            <div class="notify-item">
              <el-checkbox v-model="notificationForm.low_stock_alert">{{ $t('shop.lowStockAlert') }}</el-checkbox>
            </div>
            <el-form-item :label="$t('shop.lowStockThreshold')" v-if="notificationForm.low_stock_alert" style="margin-top: 12px">
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
                {{ $t('shop.paymentSettings') }}
              </span>
            </div>
          </template>

          <el-form :model="paymentForm" label-position="top">
            <div class="notify-item">
              <el-checkbox v-model="paymentForm.stripe_enabled">{{ $t('shop.enableStripe') }}</el-checkbox>
            </div>
            <el-form-item :label="$t('shop.stripePublicKey')" v-if="paymentForm.stripe_enabled" style="margin-top: 12px">
              <el-input v-model="paymentForm.stripe_public_key" :placeholder="$t('shop.stripePublicKeyPlaceholder')" />
            </el-form-item>
            <el-form-item :label="$t('shop.stripeSecretKey')" v-if="paymentForm.stripe_enabled">
              <el-input
                v-model="paymentForm.stripe_secret_key"
                type="password"
                :placeholder="$t('shop.stripeSecretKeyPlaceholder')"
                show-password
                @change="paymentForm.stripe_secret_key_changed = true"
              />
              <div class="form-tip">{{ $t('shop.stripeSecretKeyTip') }}</div>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- Quick Actions -->
        <el-card class="settings-card" shadow="never" style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Operation /></el-icon>
                {{ $t('shop.quickActions') }}
              </span>
            </div>
          </template>

          <div class="quick-actions">
            <el-button plain style="width: 100%; margin-bottom: 12px" @click="viewShopPage">
              <el-icon><Document /></el-icon>
              {{ $t('shop.viewShopHomepage') }}
            </el-button>
            <el-button plain style="width: 100%; margin-bottom: 12px" @click="copyShopLink">
              <el-icon><Share /></el-icon>
              {{ $t('shop.copyShopLink') }}
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Save Button -->
    <div class="save-bar">
      <el-button type="primary" size="large" @click="handleSave" :loading="saving">
        <el-icon><Check /></el-icon>
        {{ $t('shop.saveSettings') }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
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
  type NotificationSettings
} from '@/api/shop'

const { t } = useI18n({ useScope: 'global' })

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
    { required: true, message: t('shop.validationNameRequired'), trigger: 'blur' },
    { min: 2, max: 50, message: t('shop.validationNameLength'), trigger: 'blur' }
  ],
  description: [
    { max: 500, message: t('shop.validationDescLength'), trigger: 'blur' }
  ],
  contact_phone: [
    { pattern: /^$|^1[3-9]\d{9}$|^400-\d{3}-\d{4}$/, message: t('shop.validationPhone'), trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: t('shop.validationEmail'), trigger: 'blur' }
  ],
  custom_domain: [
    {
      pattern: /^$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/,
      message: t('shop.validationDomain'),
      trigger: 'blur'
    }
  ]
}

// Day names mapping
const dayKeys = ['daySunday', 'dayMonday', 'dayTuesday', 'dayWednesday', 'dayThursday', 'dayFriday', 'daySaturday']

const getDayName = (dayOfWeek: number): string => {
  return dayKeys[dayOfWeek] ? t(`shop.${dayKeys[dayOfWeek]}`) : ''
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
      businessHours.value = dayKeys.map((_, index) => ({
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
    ElMessage.error(t('shop.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Handle logo upload
const handleLogoChange = async (file: any) => {
  try {
    const response = await uploadImage(file.raw, 'banner')
    shopForm.logo = response.url
    ElMessage.success(t('shop.logoUploadSuccess'))
  } catch (error) {
    console.error('Upload failed:', error)
    ElMessage.error(t('shop.logoUploadFailed'))
  }
}

// View shop page
const viewShopPage = () => {
  if (shopSettings.value?.domain) {
    window.open(`https://${shopSettings.value.domain}`, '_blank')
  } else if (shopSettings.value?.custom_domain) {
    window.open(`https://${shopSettings.value.custom_domain}`, '_blank')
  } else {
    ElMessage.warning(t('shop.domainNotConfigured'))
  }
}

// Copy shop link
const copyShopLink = async () => {
  const link = shopSettings.value?.custom_domain || shopSettings.value?.domain
  if (link) {
    try {
      await navigator.clipboard.writeText(`https://${link}`)
      ElMessage.success(t('shop.linkCopied'))
    } catch {
      ElMessage.error(t('shop.copyFailed'))
    }
  } else {
    ElMessage.warning(t('shop.domainNotConfigured'))
  }
}

// Validate business hours
const validateBusinessHours = (hours: BusinessHours[]): boolean => {
  for (const hour of hours) {
    if (!hour.is_closed) {
      if (!hour.open_time || !hour.close_time) {
        ElMessage.error(`${getDayName(hour.day_of_week)}${t('shop.businessHoursNotSet')}`)
        return false
      }
      if (hour.open_time >= hour.close_time) {
        ElMessage.error(`${getDayName(hour.day_of_week)}${t('shop.endTimeMustBeAfterStart')}`)
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
    ElMessage.warning(t('shop.pleaseCheckForm'))
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

    ElMessage.success(t('shop.settingsSaved'))

    // Reset secret key change flag
    paymentForm.stripe_secret_key_changed = false

    // Reload settings to get updated values
    await loadSettings()
  } catch (error: unknown) {
    console.error('Failed to save settings:', error)
    const err = error as { response?: { data?: { message?: string } } }
    ElMessage.error(err?.response?.data?.message || t('common.actionFailed'))
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