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
          <el-menu-item index="/returns">
            <span>退换货</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <template #title>
            <span>顾客管理</span>
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
          <el-menu-item index="/appearance">
            <span>外观装修</span>
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import Breadcrumb from '@/components/Breadcrumb.vue'
import {
  ShoppingBag, DataLine, Goods, List, User, UserFilled, Ticket, Shop,
  TrendCharts, Expand, Fold, FullScreen, Bell, ArrowDown,
  Setting, SwitchButton, Box, CreditCard, Van, Picture
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const searchQuery = ref('')
const notificationVisible = ref(false)
const hasNewData = ref(true)
const pendingOrders = ref(5)

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
</script>

<style scoped>
.main-layout {
  height: 100vh;
  overflow: hidden;
}

/* Sidebar */
.sidebar {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
}

.logo-section {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
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
  letter-spacing: 1px;
}

.version {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
  margin-top: 4px;
  margin-left: 40px;
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
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, #059669 0%, transparent 100%);
  color: #fff;
  border-left: 3px solid #10B981;
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
  background: linear-gradient(90deg, #059669 0%, transparent 100%);
  color: #fff;
  border-left: 3px solid #10B981;
}

.menu-badge {
  margin-left: 8px;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.storage-info {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 12px;
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
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 18px;
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
  width: 240px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 20px;
}

.notification-badge :deep(.el-badge__content) {
  top: 8px;
  right: 8px;
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
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-info:hover {
  background: #F3F4F6;
}

.user-avatar {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  font-weight: 600;
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #111827;
}

.user-role {
  font-size: 12px;
  color: #6B7280;
}

/* Main Content */
.main-content {
  background: #F3F4F6;
  padding: 20px;
  overflow-y: auto;
}

/* Notification Drawer */
.notification-list {
  padding: 16px 0;
}

.notification-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 12px;
  background: #F9FAFB;
  transition: all 0.2s;
}

.notification-item.unread {
  background: #ECFDF5;
}

.notification-item:hover {
  transform: translateX(4px);
}

.notice-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notice-icon.order {
  background: #DBEAFE;
  color: #3B82F6;
}

.notice-icon.warning {
  background: #FEF3C7;
  color: #F59E0B;
}

.notice-icon.success {
  background: #D1FAE5;
  color: #059669;
}

.notice-icon.info {
  background: #E5E7EB;
  color: #6B7280;
}

.notice-content {
  flex: 1;
}

.notice-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 4px 0;
}

.notice-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.notice-time {
  font-size: 12px;
  color: #9CA3AF;
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
