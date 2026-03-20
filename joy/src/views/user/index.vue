<template>
  <div class="user-page">
    <div class="page-container">
      <!-- User Profile Card -->
      <div class="profile-card">
        <div class="profile-header">
          <div class="avatar">
            <img v-if="user.avatar" :src="user.avatar" :alt="user.name" />
            <span v-else class="avatar-text">{{ user.name.charAt(0) }}</span>
          </div>
          <div class="profile-info">
            <h2 class="user-name">{{ user.name }}</h2>
            <p class="user-level">
              <SparklesIcon class="icon" />
              {{ user.level }}
            </p>
          </div>
          <router-link to="/user/settings" class="settings-btn">
            <Cog6ToothIcon class="icon" />
          </router-link>
        </div>

        <div class="profile-stats">
          <div class="stat-item">
            <span class="stat-value">{{ user.coupons }}</span>
            <span class="stat-label">优惠券</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ user.points }}</span>
            <span class="stat-label">积分</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ user.balance }}</span>
            <span class="stat-label">余额</span>
          </div>
        </div>
      </div>

      <!-- Order Status -->
      <div class="section-card">
        <div class="section-header">
          <h3>我的订单</h3>
          <router-link to="/orders" class="view-all">
            查看全部
            <ChevronRightIcon class="icon" />
          </router-link>
        </div>
        <div class="order-status-grid">
          <router-link to="/orders?status=pending" class="status-item">
            <div class="status-icon">
              <CreditCardIcon class="icon" />
              <span v-if="orderStats.pending > 0" class="badge">{{ orderStats.pending }}</span>
            </div>
            <span class="status-label">待付款</span>
          </router-link>
          <router-link to="/orders?status=paid" class="status-item">
            <div class="status-icon">
              <CubeIcon class="icon" />
              <span v-if="orderStats.shipping > 0" class="badge">{{ orderStats.shipping }}</span>
            </div>
            <span class="status-label">待发货</span>
          </router-link>
          <router-link to="/orders?status=shipped" class="status-item">
            <div class="status-icon">
              <TruckIcon class="icon" />
              <span v-if="orderStats.receiving > 0" class="badge">{{ orderStats.receiving }}</span>
            </div>
            <span class="status-label">待收货</span>
          </router-link>
          <router-link to="/orders?status=completed" class="status-item">
            <div class="status-icon">
              <ChatBubbleLeftRightIcon class="icon" />
              <span v-if="orderStats.review > 0" class="badge">{{ orderStats.review }}</span>
            </div>
            <span class="status-label">待评价</span>
          </router-link>
          <router-link to="/orders?status=refund" class="status-item">
            <div class="status-icon">
              <ArrowPathIcon class="icon" />
            </div>
            <span class="status-label">退款/售后</span>
          </router-link>
        </div>
      </div>

      <!-- Service Menu -->
      <div class="section-card">
        <h3 class="section-title">我的服务</h3>
        <div class="service-grid">
          <router-link to="/user/address" class="service-item">
            <div class="service-icon address">
              <MapPinIcon class="icon" />
            </div>
            <span class="service-label">收货地址</span>
          </router-link>
          <router-link to="/user/coupons" class="service-item">
            <div class="service-icon coupon">
              <TicketIcon class="icon" />
            </div>
            <span class="service-label">优惠券</span>
          </router-link>
          <router-link to="/user/favorites" class="service-item">
            <div class="service-icon favorite">
              <HeartIcon class="icon" />
            </div>
            <span class="service-label">我的收藏</span>
          </router-link>
          <router-link to="/user/history" class="service-item">
            <div class="service-icon history">
              <ClockIcon class="icon" />
            </div>
            <span class="service-label">浏览记录</span>
          </router-link>
          <router-link to="/user/points" class="service-item">
            <div class="service-icon points">
              <GiftIcon class="icon" />
            </div>
            <span class="service-label">积分商城</span>
          </router-link>
          <router-link to="/user/vip" class="service-item">
            <div class="service-icon vip">
              <CrownIcon class="icon" />
            </div>
            <span class="service-label">会员中心</span>
          </router-link>
          <router-link to="/user/help" class="service-item">
            <div class="service-icon help">
              <QuestionMarkCircleIcon class="icon" />
            </div>
            <span class="service-label">帮助中心</span>
          </router-link>
          <router-link to="/user/service" class="service-item">
            <div class="service-icon customer">
              <HeadphonesIcon class="icon" />
            </div>
            <span class="service-label">客服中心</span>
          </router-link>
        </div>
      </div>

      <!-- Account Settings -->
      <div class="section-card">
        <h3 class="section-title">账号设置</h3>
        <div class="settings-list">
          <router-link to="/user/profile" class="setting-item">
            <div class="setting-left">
              <UserCircleIcon class="icon" />
              <span>个人资料</span>
            </div>
            <ChevronRightIcon class="arrow" />
          </router-link>
          <router-link to="/user/security" class="setting-item">
            <div class="setting-left">
              <ShieldCheckIcon class="icon" />
              <span>账号安全</span>
            </div>
            <ChevronRightIcon class="arrow" />
          </router-link>
          <router-link to="/user/notification" class="setting-item">
            <div class="setting-left">
              <BellIcon class="icon" />
              <span>消息通知</span>
            </div>
            <ChevronRightIcon class="arrow" />
          </router-link>
          <router-link to="/user/privacy" class="setting-item">
            <div class="setting-left">
              <LockClosedIcon class="icon" />
              <span>隐私设置</span>
            </div>
            <ChevronRightIcon class="arrow" />
          </router-link>
        </div>
      </div>

      <!-- Logout Button -->
      <button class="logout-btn" @click="logout">
        <ArrowRightOnRectangleIcon class="icon" />
        退出登录
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import {
  SparklesIcon,
  Cog6ToothIcon,
  ChevronRightIcon,
  CreditCardIcon,
  CubeIcon,
  TruckIcon,
  ChatBubbleLeftRightIcon,
  ArrowPathIcon,
  MapPinIcon,
  TicketIcon,
  HeartIcon,
  ClockIcon,
  GiftIcon,
  CrownIcon,
  QuestionMarkCircleIcon,
  HeadphonesIcon,
  UserCircleIcon,
  ShieldCheckIcon,
  BellIcon,
  LockClosedIcon,
  ArrowRightOnRectangleIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()

const user = {
  name: '张三',
  avatar: '',
  level: '黄金会员',
  coupons: 3,
  points: 2580,
  balance: '¥0.00'
}

const orderStats = {
  pending: 2,
  shipping: 1,
  receiving: 1,
  review: 3
}

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    router.push('/')
  }
}
</script>

<style scoped>
.user-page {
  min-height: calc(100vh - 72px);
  background: #F9FAFB;
  padding: 24px 0;
}

.page-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Profile Card */
.profile-card {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border-radius: 20px;
  padding: 24px;
  color: white;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: white;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  font-size: 32px;
  font-weight: 700;
  color: #059669;
}

.profile-info {
  flex: 1;
}

.user-name {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.user-level {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.user-level .icon {
  width: 18px;
  height: 18px;
}

.settings-btn {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.settings-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.settings-btn .icon {
  width: 24px;
  height: 24px;
  color: white;
}

.profile-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
}

.stat-label {
  font-size: 13px;
  opacity: 0.8;
}

/* Section Card */
.section-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3,
.section-title {
  font-size: 17px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.view-all {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #6B7280;
  text-decoration: none;
  transition: color 0.2s;
}

.view-all:hover {
  color: #059669;
}

.view-all .icon {
  width: 18px;
  height: 18px;
}

/* Order Status Grid */
.order-status-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px;
}

.status-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.2s;
}

.status-item:hover {
  background: #F9FAFB;
}

.status-icon {
  position: relative;
  width: 48px;
  height: 48px;
  background: #F3F4F6;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-item:hover .status-icon {
  background: #ECFDF5;
}

.status-icon .icon {
  width: 24px;
  height: 24px;
  color: #6B7280;
}

.status-item:hover .status-icon .icon {
  color: #059669;
}

.badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: #EF4444;
  color: white;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-label {
  font-size: 13px;
  color: #4B5563;
}

/* Service Grid */
.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.service-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.2s;
}

.service-item:hover {
  background: #F9FAFB;
}

.service-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-icon .icon {
  width: 24px;
  height: 24px;
}

.service-icon.address {
  background: #DBEAFE;
  color: #3B82F6;
}

.service-icon.coupon {
  background: #FEF3C7;
  color: #F59E0B;
}

.service-icon.favorite {
  background: #FCE7F3;
  color: #EC4899;
}

.service-icon.history {
  background: #E0E7FF;
  color: #6366F1;
}

.service-icon.points {
  background: #D1FAE5;
  color: #10B981;
}

.service-icon.vip {
  background: #FEF9C3;
  color: #EAB308;
}

.service-icon.help {
  background: #F3E8FF;
  color: #A855F7;
}

.service-icon.customer {
  background: #ECFDF5;
  color: #059669;
}

.service-label {
  font-size: 13px;
  color: #4B5563;
}

/* Settings List */
.settings-list {
  display: flex;
  flex-direction: column;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid #F3F4F6;
  text-decoration: none;
  transition: all 0.2s;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item:hover {
  padding-left: 8px;
}

.setting-left {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 15px;
  color: #374151;
}

.setting-left .icon {
  width: 22px;
  height: 22px;
  color: #6B7280;
}

.arrow {
  width: 20px;
  height: 20px;
  color: #9CA3AF;
}

/* Logout Button */
.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  color: #EF4444;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn:hover {
  background: #FEF2F2;
  border-color: #FECACA;
}

.logout-btn .icon {
  width: 20px;
  height: 20px;
}

/* Responsive */
@media (max-width: 640px) {
  .order-status-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .service-grid {
    grid-template-columns: repeat(4, 1fr);
  }

  .profile-stats {
    gap: 16px;
  }
}
</style>
