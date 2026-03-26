<template>
  <el-container class="main-layout">
    <!-- Sidebar -->
    <el-aside width="220px" class="sidebar">
      <div class="logo-section">
        <div class="logo">
          <el-icon size="28" color="#fff"><ShoppingBag /></el-icon>
          <span class="logo-text">ShopJoy</span>
        </div>
        <div class="version">v2.0.0</div>
      </div>
      
      <el-menu
        :default-active="$route.path"
        router
        class="sidebar-menu"
        :collapse="isCollapse"
        :collapse-transition="false"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <template #title>
            <span>数据概览</span>
            <el-tag v-if="hasNewData" size="small" type="danger" effect="dark" class="menu-badge">NEW</el-tag>
          </template>
        </el-menu-item>
        
        <el-sub-menu index="/products">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="/products">
            <span>商品列表</span>
          </el-menu-item>
          <el-menu-item index="/categories">
            <span>商品分类</span>
          </el-menu-item>
          <el-menu-item index="/brands">
            <span>品牌管理</span>
          </el-menu-item>
          <el-menu-item index="/inventory">
            <span>库存管理</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/orders">
          <template #title>
            <el-icon><List /></el-icon>
            <span>订单管理</span>
          </template>
          <el-menu-item index="/orders">
            <span>全部订单</span>
            <el-badge :value="pendingOrders" v-if="pendingOrders > 0" class="menu-badge" />
          </el-menu-item>
          <el-menu-item index="/fulfillment/shipments">
            <span>发货管理</span>
          </el-menu-item>
          <el-menu-item index="/fulfillment/refunds">
            <span>退款管理</span>
            <el-badge :value="pendingRefunds" v-if="pendingRefunds > 0" class="menu-badge" type="warning" />
          </el-menu-item>
          <el-menu-item index="/fulfillment/statistics">
            <span>售后统计</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <template #title>
            <span>顾客管理</span>
          </template>
        </el-menu-item>

        <el-menu-item index="/reviews">
          <el-icon><ChatDotRound /></el-icon>
          <template #title>
            <span>评价管理</span>
            <el-badge :value="pendingReviews" v-if="pendingReviews > 0" class="menu-badge" type="warning" />
          </template>
        </el-menu-item>

        <el-menu-item index="/admin-users">
          <el-icon><UserFilled /></el-icon>
          <template #title>
            <span>用户管理</span>
          </template>
        </el-menu-item>
        
        <el-menu-item index="/promotions">
          <el-icon><Ticket /></el-icon>
          <template #title>
            <span>营销推广</span>
          </template>
        </el-menu-item>
        
        <el-sub-menu index="/shop">
          <template #title>
            <el-icon><Shop /></el-icon>
            <span>店铺设置</span>
          </template>
          <el-menu-item index="/shop">
            <span>基本设置</span>
          </el-menu-item>
          <el-menu-item index="/settings/markets">
            <span>市场管理</span>
          </el-menu-item>
          <el-menu-item index="/storefront/themes">
            <span>主题管理</span>
          </el-menu-item>
          <el-menu-item index="/storefront/pages">
            <span>页面装修</span>
          </el-menu-item>
          <el-menu-item index="/storefront/seo">
            <span>SEO设置</span>
          </el-menu-item>
          <el-menu-item index="/shipping">
            <span>物流配置</span>
          </el-menu-item>
          <el-menu-item index="/payment">
            <span>支付设置</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/analytics">
          <el-icon><TrendCharts /></el-icon>
          <template #title>
            <span>数据分析</span>
          </template>
        </el-menu-item>
      </el-menu>
      
      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <div class="storage-info">
          <div class="storage-header">
            <span>存储空间</span>
            <span>75%</span>
          </div>
          <el-progress :percentage="75" :stroke-width="6" :show-text="false" color="#10B981" />
          <p class="storage-detail">7.5GB / 10GB 已使用</p>
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
            @click="toggleCollapse"
            text
          />
          <breadcrumb />
        </div>
        
        <div class="header-right">
          <!-- Search -->
          <div class="header-search">
            <el-input
              v-model="searchQuery"
              placeholder="搜索..."
              prefix-icon="Search"
              clearable
              class="search-input"
            />
          </div>
          
          <!-- Quick Actions -->
          <el-tooltip :content="isDarkMode ? '切换到浅色模式' : '切换到深色模式'" placement="bottom">
            <el-button text @click="toggleTheme" class="theme-toggle-btn">
              <el-icon size="18">
                <Sunny v-if="isDarkMode" />
                <Moon v-else />
              </el-icon>
            </el-button>
          </el-tooltip>

          <el-tooltip content="全屏" placement="bottom">
            <el-button text :icon="FullScreen" @click="toggleFullscreen" />
          </el-tooltip>
          
          <el-tooltip content="通知" placement="bottom">
            <el-badge :value="3" class="notification-badge">
              <el-button text :icon="Bell" @click="showNotifications" />
            </el-badge>
          </el-tooltip>
          
          <!-- User Menu -->
          <el-dropdown trigger="click" class="user-dropdown">
            <div class="user-info">
              <el-avatar :size="36" :src="userStore.userInfo?.avatar" class="user-avatar">
                {{ userStore.userInfo?.name?.charAt(0) || 'A' }}
              </el-avatar>
              <div class="user-meta">
                <span class="user-name">{{ userStore.userInfo?.name || 'Admin' }}</span>
                <span class="user-role">超级管理员</span>
              </div>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToProfile">
                  <el-icon><User /></el-icon>个人中心
                </el-dropdown-item>
                <el-dropdown-item @click="goToSettings">
                  <el-icon><Setting /></el-icon>账号设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- Main Content -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
  
  <!-- Notification Drawer -->
  <el-drawer
    v-model="notificationVisible"
    title="通知中心"
    size="400px"
    :with-header="true"
  >
    <div class="notification-list">
      <div v-for="(notice, index) in notifications" :key="index" class="notification-item" :class="{ unread: !notice.read }">
        <div class="notice-icon" :class="notice.type">
          <el-icon size="20">
            <component :is="notice.icon" />
          </el-icon>
        </div>
        <div class="notice-content">
          <h4 class="notice-title">{{ notice.title }}</h4>
          <p class="notice-desc">{{ notice.content }}</p>
          <span class="notice-time">{{ notice.time }}</span>
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
    title: '新订单提醒',
    content: '您有一个新订单 #ORD20240318010 待处理',
    time: '10分钟前',
    type: 'order',
    icon: 'Box',
    read: false
  },
  {
    title: '库存预警',
    content: '商品 "无线蓝牙耳机 Pro" 库存不足，请及时补货',
    time: '30分钟前',
    type: 'warning',
    icon: 'Goods',
    read: false
  },
  {
    title: '支付成功',
    content: '用户 张三 已完成支付，订单金额 ¥299.00',
    time: '1小时前',
    type: 'success',
    icon: 'CreditCard',
    read: false
  },
  {
    title: '系统通知',
    content: '系统将于今晚 02:00 进行例行维护',
    time: '2小时前',
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
  ElMessage.info('个人中心')
}

const goToSettings = () => {
  ElMessage.info('账号设置')
}

const logout = () => {
  userStore.clearToken()
  router.push('/login')
  ElMessage.success('退出登录成功')
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
