<template>
  <el-container class="main-layout">
    <!-- Sidebar -->
    <el-aside
      width="220px"
      class="sidebar"
    >
      <div class="logo-section">
        <div class="logo">
          <el-icon
            size="28"
            color="#fff"
          >
            <ShoppingBag />
          </el-icon>
          <span class="logo-text">ShopJoy</span>
        </div>
        <div class="version">
          v2.0.0
        </div>
      </div>
      
      <el-menu
        :default-active="$route.path"
        class="sidebar-menu"
        :collapse="isCollapse"
        :collapse-transition="false"
      >
        <el-menu-item index="/dashboard" @click="router.push('/dashboard')">
          <el-icon><DataLine /></el-icon>
          <template #title>
            <span>{{ $t('dashboard.title') }}</span>
            <el-tag
              v-if="hasNewData"
              size="small"
              type="danger"
              effect="dark"
              class="menu-badge"
            >
              NEW
            </el-tag>
          </template>
        </el-menu-item>

        <el-sub-menu index="/products">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>{{ $t('products.title') }}</span>
          </template>
          <el-menu-item index="/products" @click="router.push('/products')">
            <span>{{ $t('products.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/categories" @click="router.push('/categories')">
            <span>{{ $t('categories.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/brands" @click="router.push('/brands')">
            <span>{{ $t('brands.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/inventory" @click="router.push('/inventory')">
            <span>{{ $t('inventory.title') }}</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="/orders">
          <template #title>
            <el-icon><List /></el-icon>
            <span>{{ $t('orders.title') }}</span>
          </template>
          <el-menu-item index="/orders" @click="router.push('/orders')">
            <span>{{ $t('orders.title') }}</span>
            <el-badge
              v-if="pendingOrders > 0"
              :value="pendingOrders"
              class="menu-badge"
            />
          </el-menu-item>
          <el-menu-item index="/fulfillment/shipments" @click="router.push('/fulfillment/shipments')">
            <span>{{ $t('fulfillment.shipments') }}</span>
          </el-menu-item>
          <el-menu-item index="/fulfillment/refunds" @click="router.push('/fulfillment/refunds')">
            <span>{{ $t('fulfillment.refunds') }}</span>
            <el-badge
              v-if="pendingRefunds > 0"
              :value="pendingRefunds"
              class="menu-badge"
              type="warning"
            />
          </el-menu-item>
          <el-menu-item index="/fulfillment/statistics" @click="router.push('/fulfillment/statistics')">
            <span>{{ $t('fulfillment.statistics') }}</span>
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/users" @click="router.push('/users')">
          <el-icon><User /></el-icon>
          <template #title>
            <span>{{ $t('users.title') }}</span>
          </template>
        </el-menu-item>

        <el-menu-item index="/reviews" @click="router.push('/reviews')">
          <el-icon><ChatDotRound /></el-icon>
          <template #title>
            <span>{{ $t('reviews.title') }}</span>
            <el-badge
              v-if="pendingReviews > 0"
              :value="pendingReviews"
              class="menu-badge"
              type="warning"
            />
          </template>
        </el-menu-item>

        <el-menu-item index="/admin-users" @click="router.push('/admin-users')">
          <el-icon><UserFilled /></el-icon>
          <template #title>
            <span>{{ $t('adminUsers.title') }}</span>
          </template>
        </el-menu-item>

        <el-menu-item index="/promotions" @click="router.push('/promotions')">
          <el-icon><Ticket /></el-icon>
          <template #title>
            <span>{{ $t('promotions.title') }}</span>
          </template>
        </el-menu-item>

        <el-sub-menu index="/shop">
          <template #title>
            <el-icon><Shop /></el-icon>
            <span>{{ $t('shop.title') }}</span>
          </template>
          <el-menu-item index="/shop" @click="router.push('/shop')">
            <span>{{ $t('shop.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/settings/markets" @click="router.push('/settings/markets')">
            <span>{{ $t('settings.markets.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/storefront/themes" @click="router.push('/storefront/themes')">
            <span>{{ $t('storefront.themes') }}</span>
          </el-menu-item>
          <el-menu-item index="/storefront/pages" @click="router.push('/storefront/pages')">
            <span>{{ $t('storefront.pages') }}</span>
          </el-menu-item>
          <el-menu-item index="/storefront/seo" @click="router.push('/storefront/seo')">
            <span>{{ $t('storefront.seo') }}</span>
          </el-menu-item>
          <el-menu-item index="/shipping" @click="router.push('/shipping')">
            <span>{{ $t('shipping.title') }}</span>
          </el-menu-item>
          <el-menu-item index="/payments" @click="router.push('/payments')">
            <span>{{ $t('payments.title') }}</span>
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/analytics">
          <el-icon><TrendCharts /></el-icon>
          <template #title>
            <span>{{ $t('dashboard.title') }}</span>
          </template>
        </el-menu-item>
      </el-menu>
      
      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <div class="storage-info">
          <div class="storage-header">
            <span>{{ $t('common.storageSpace') }}</span>
            <span>75%</span>
          </div>
          <el-progress
            :percentage="75"
            :stroke-width="6"
            :show-text="false"
            color="#10B981"
          />
          <p class="storage-detail">
            7.5GB / 10GB {{ $t('common.used') }}
          </p>
        </div>
      </div>
    </el-aside>
    
    <el-container>
      <!-- Header -->
      <el-header class="header">
        <div class="header-left">
          <el-button
            class="collapse-btn"
            :icon="isCollapse ? Expand : Fold"
            text
            @click="toggleCollapse"
          />
          <breadcrumb />
        </div>
        
        <div class="header-right">
          <!-- Search -->
          <div class="header-search">
            <el-input
              v-model="searchQuery"
              :placeholder="$t('common.search')"
              prefix-icon="Search"
              clearable
              class="search-input"
            />
          </div>
          
          <!-- Quick Actions -->
          <el-tooltip
            :content="isDarkMode ? $t('common.switchToLightMode') : $t('common.switchToDarkMode')"
            placement="bottom"
          >
            <el-button
              text
              class="theme-toggle-btn"
              @click="toggleTheme"
            >
              <el-icon size="18">
                <Sunny v-if="isDarkMode" />
                <Moon v-else />
              </el-icon>
            </el-button>
          </el-tooltip>

          <el-tooltip
            :content="$t('common.fullscreen')"
            placement="bottom"
          >
            <el-button
              text
              :icon="FullScreen"
              @click="toggleFullscreen"
            />
          </el-tooltip>
          
          <el-tooltip
            :content="$t('common.notifications')"
            placement="bottom"
          >
            <el-badge
              :value="3"
              class="notification-badge"
            >
              <el-button
                text
                :icon="Bell"
                @click="showNotifications"
              />
            </el-badge>
          </el-tooltip>

          <!-- Language Switcher -->
          <el-dropdown
            trigger="click"
            @command="handleLocaleChange"
          >
            <div class="locale-switch">
              <span class="locale-text">{{ $t(locale === 'en' ? 'common.languageEN' : 'common.languageZH') }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  command="en"
                  :disabled="locale === 'en'"
                >
                  {{ $t('common.languageEN') }}
                </el-dropdown-item>
                <el-dropdown-item
                  command="zh"
                  :disabled="locale === 'zh'"
                >
                  {{ $t('common.languageZH') }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <!-- User Menu -->
          <el-dropdown
            trigger="click"
            class="user-dropdown"
          >
            <div class="user-info">
              <el-avatar
                :size="36"
                :src="userStore.userInfo?.avatar"
                class="user-avatar"
              >
                {{ userStore.userInfo?.real_name?.charAt(0) || 'A' }}
              </el-avatar>
              <div class="user-meta">
                <span class="user-name">{{ userStore.userInfo?.real_name || $t('common.defaultAdmin') }}</span>
                <span class="user-role">{{ $t('common.superAdmin') }}</span>
              </div>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToProfile">
                  <el-icon><User /></el-icon>{{ $t('common.profile') }}
                </el-dropdown-item>
                <el-dropdown-item @click="goToSettings">
                  <el-icon><Setting /></el-icon>{{ $t('common.accountSettings') }}
                </el-dropdown-item>
                <el-dropdown-item
                  divided
                  @click="logout"
                >
                  <el-icon><SwitchButton /></el-icon>{{ $t('common.logout') }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- Main Content -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition
            name="fade"
            mode="out-in"
          >
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
  
  <!-- Notification Drawer -->
  <el-drawer
    v-model="notificationVisible"
    :title="$t('notifications.title')"
    size="400px"
    :with-header="true"
  >
    <div class="notification-list">
      <div
        v-for="(notice, index) in notifications"
        :key="index"
        class="notification-item"
        :class="{ unread: !notice.read }"
      >
        <div
          class="notice-icon"
          :class="notice.type"
        >
          <el-icon size="20">
            <component :is="notice.icon" />
          </el-icon>
        </div>
        <div class="notice-content">
          <h4 class="notice-title">
            {{ $t(notice.titleKey) }}
          </h4>
          <p class="notice-desc">
            {{ $t(notice.contentKey) }}
          </p>
          <span class="notice-time">{{ $t(notice.timeKey, notice.timeParams) }}</span>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import Breadcrumb from '@/components/Breadcrumb.vue'
import {
  ShoppingBag, DataLine, Goods, List, User, UserFilled, Ticket, Shop,
  TrendCharts, Expand, Fold, FullScreen, Bell, ArrowDown,
  Setting, SwitchButton, Moon, Sunny, ChatDotRound
} from '@element-plus/icons-vue'
import { useLocale, t } from '@/plugins/i18n'

const { locale } = useLocale()

const handleLocaleChange = (lang: string) => {
  locale.value = lang
  localStorage.setItem('locale', lang)
}

const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const searchQuery = ref('')
const notificationVisible = ref(false)
const hasNewData = ref(true)
const pendingOrders = ref(5)
const pendingRefunds = ref(8)
const pendingReviews = ref(12)
const isDarkMode = ref(false)

const notifications = ref([
  {
    titleKey: 'notifications.newOrderTitle',
    contentKey: 'notifications.newOrderContent',
    timeKey: 'notifications.minutesAgo',
    timeParams: { n: 10 },
    type: 'order',
    icon: 'Box',
    read: false
  },
  {
    titleKey: 'notifications.lowStockTitle',
    contentKey: 'notifications.lowStockContent',
    timeKey: 'notifications.minutesAgo',
    timeParams: { n: 30 },
    type: 'warning',
    icon: 'Goods',
    read: false
  },
  {
    titleKey: 'notifications.paymentSuccessTitle',
    contentKey: 'notifications.paymentSuccessContent',
    timeKey: 'notifications.hoursAgo',
    timeParams: { n: 1 },
    type: 'success',
    icon: 'CreditCard',
    read: false
  },
  {
    titleKey: 'notifications.systemNoticeTitle',
    contentKey: 'notifications.systemNoticeContent',
    timeKey: 'notifications.hoursAgo',
    timeParams: { n: 2 },
    type: 'info',
    icon: 'Setting',
    read: true
  }
])

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
    localStorage.setItem('theme', 'light')
  }
}

const initTheme = () => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
}

const showNotifications = () => {
  notificationVisible.value = true
}

const goToProfile = () => {
  ElMessage.info(t('common.profile'))
}

const goToSettings = () => {
  ElMessage.info(t('common.accountSettings'))
}

const logout = () => {
  userStore.clearToken()
  router.push('/login')
  ElMessage.success(t('common.logoutSuccess'))
}

onMounted(() => {
  initTheme()
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
  overflow: hidden;
}

/* Sidebar */
.sidebar {
  background: linear-gradient(180deg, #1E1B4B 0%, #312E81 100%);
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.logo-section {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-text {
  color: #fff;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.5px;
  font-family: 'Fira Sans', sans-serif;
}

.version {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
  margin-top: 4px;
  margin-left: 40px;
  font-family: 'Fira Code', monospace;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
  background: transparent;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  color: rgba(255, 255, 255, 0.7);
  height: 50px;
  line-height: 50px;
  margin: 4px 8px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.8) 0%, rgba(99, 102, 241, 0.2) 100%);
  color: #fff;
  border-left: 3px solid #A5B4FC;
  border-radius: 0 8px 8px 0;
  margin-left: 0;
  padding-left: 17px;
}

.sidebar-menu :deep(.el-icon) {
  color: inherit;
}

/* Sub-menu styles */
.sidebar-menu :deep(.el-sub-menu .el-menu) {
  background: transparent;
}

.sidebar-menu :deep(.el-sub-menu .el-menu-item) {
  color: rgba(255, 255, 255, 0.6);
  height: 46px;
  line-height: 46px;
  padding-left: 54px !important;
}

.sidebar-menu :deep(.el-sub-menu .el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.sidebar-menu :deep(.el-sub-menu .el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.6) 0%, transparent 100%);
  color: #fff;
  border-left: 3px solid #A5B4FC;
}

.menu-badge {
  margin-left: 8px;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.storage-info {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  padding: 14px;
}

.storage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}

.storage-detail {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  margin: 8px 0 0 0;
}

/* Header */
.header {
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
}

/* Dark mode header */
[data-theme="dark"] .header {
  background: var(--color-bg-card);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 18px;
  color: var(--color-text-secondary);
  transition: color 0.2s;
}

.collapse-btn:hover {
  color: var(--color-primary);
}

.theme-toggle-btn {
  transition: all 0.3s ease;
}

.theme-toggle-btn:hover {
  color: var(--color-primary);
}

[data-theme="dark"] .theme-toggle-btn {
  color: var(--color-warning);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-search {
  margin-right: 16px;
}

.search-input {
  width: 260px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.notification-badge :deep(.el-badge__content) {
  top: 8px;
  right: 8px;
  background-color: var(--color-primary);
}

/* Language Switcher */
.locale-switch {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--color-text-secondary);
}

.locale-switch:hover {
  background: #F5F3FF;
  color: var(--color-primary);
}

[data-theme="dark"] .locale-switch:hover {
  background: rgba(99, 102, 241, 0.1);
  color: var(--color-primary-light);
}

.locale-text {
  font-size: 14px;
  font-weight: 500;
}

/* User Dropdown */
.user-dropdown {
  margin-left: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 6px 10px;
  border-radius: 10px;
  transition: background 0.2s ease;
}

.user-info:hover {
  background: #F5F3FF;
}

[data-theme="dark"] .user-info:hover {
  background: rgba(99, 102, 241, 0.1);
}

.user-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  font-weight: 600;
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #1E1B4B;
}

[data-theme="dark"] .user-name {
  color: var(--color-text-primary);
}

.user-role {
  font-size: 12px;
  color: #6B7280;
}

/* Main Content */
.main-content {
  background: #F5F3FF;
  padding: 20px;
  overflow-y: auto;
}

[data-theme="dark"] .main-content {
  background: var(--color-bg-base);
}

/* Notification Drawer */
.notification-list {
  padding: 16px 0;
}

.notification-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 12px;
  background: #F9FAFB;
  transition: all 0.2s ease;
  cursor: pointer;
}

[data-theme="dark"] .notification-item {
  background: var(--color-bg-elevated);
}

.notification-item.unread {
  background: #EEF2FF;
  border-left: 3px solid var(--color-primary);
}

[data-theme="dark"] .notification-item.unread {
  background: rgba(99, 102, 241, 0.15);
  border-left: 3px solid var(--color-primary-light);
}

.notification-item:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.notice-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notice-icon.order {
  background: #EEF2FF;
  color: #6366F1;
}

[data-theme="dark"] .notice-icon.order {
  background: rgba(99, 102, 241, 0.2);
  color: var(--color-primary-light);
}

.notice-icon.warning {
  background: #FEF3C7;
  color: #F59E0B;
}

[data-theme="dark"] .notice-icon.warning {
  background: rgba(245, 158, 11, 0.2);
  color: var(--color-warning);
}

.notice-icon.success {
  background: #D1FAE5;
  color: #10B981;
}

[data-theme="dark"] .notice-icon.success {
  background: rgba(16, 185, 129, 0.2);
  color: var(--color-cta);
}

.notice-icon.info {
  background: #F3F4F6;
  color: #6B7280;
}

[data-theme="dark"] .notice-icon.info {
  background: rgba(107, 114, 128, 0.2);
  color: var(--color-text-muted);
}

.notice-content {
  flex: 1;
}

.notice-title {
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

[data-theme="dark"] .notice-title {
  color: var(--color-text-primary);
}

.notice-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

[data-theme="dark"] .notice-desc {
  color: var(--color-text-secondary);
}

.notice-time {
  font-size: 12px;
  color: #9CA3AF;
  font-family: 'Fira Code', monospace;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: -220px;
    z-index: 1000;
    height: 100vh;
  }

  .sidebar.is-open {
    left: 0;
  }

  .header-search {
    display: none;
  }

  .user-meta {
    display: none;
  }
}
</style>
