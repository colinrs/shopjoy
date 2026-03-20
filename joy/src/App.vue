<template>
  <div class="app-layout">
    <!-- Navigation -->
    <nav class="navbar" :class="{ 'scrolled': isScrolled }">
      <div class="nav-container">
        <!-- Logo -->
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <ShoppingBagIcon class="icon" />
          </div>
          <span class="logo-text">ShopJoy</span>
        </router-link>

        <!-- Search Bar -->
        <div class="search-bar" :class="{ 'active': searchActive }">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索商品..."
            @focus="searchActive = true"
            @blur="searchActive = false"
            @keyup.enter="handleSearch"
          />
          <button class="search-btn" @click="handleSearch">
            <MagnifyingGlassIcon class="icon" />
          </button>
        </div>

        <!-- Navigation Links -->
        <div class="nav-links" :class="{ 'mobile-open': mobileMenuOpen }">
          <router-link to="/" class="nav-link" active-class="active">首页</router-link>
          <router-link to="/products" class="nav-link" active-class="active">全部商品</router-link>
          <router-link to="/products?category=new" class="nav-link" active-class="active">新品上市</router-link>
          <router-link to="/products?category=sale" class="nav-link sale" active-class="active">限时特惠</router-link>
        </div>

        <!-- Right Actions -->
        <div class="nav-actions">
          <!-- Cart -->
          <router-link to="/cart" class="action-btn cart-btn">
            <ShoppingCartIcon class="icon" />
            <span v-if="cartCount > 0" class="cart-badge">{{ cartCount }}</span>
          </router-link>

          <!-- User Menu -->
          <div class="user-menu" v-if="isLoggedIn">
            <button class="action-btn user-btn" @click="toggleUserMenu">
              <UserCircleIcon class="icon" />
            </button>
            <div class="dropdown-menu" v-show="userMenuOpen">
              <router-link to="/user" class="dropdown-item">
                <UserIcon class="dropdown-icon" />
                个人中心
              </router-link>
              <router-link to="/orders" class="dropdown-item">
                <ClipboardDocumentListIcon class="dropdown-icon" />
                我的订单
              </router-link>
              <div class="dropdown-divider"></div>
              <button class="dropdown-item" @click="logout">
                <ArrowRightOnRectangleIcon class="dropdown-icon" />
                退出登录
              </button>
            </div>
          </div>
          <router-link v-else to="/login" class="login-btn">
            登录
          </router-link>

          <!-- Mobile Menu Toggle -->
          <button class="mobile-toggle" @click="toggleMobileMenu">
            <Bars3Icon v-if="!mobileMenuOpen" class="icon" />
            <XMarkIcon v-else class="icon" />
          </button>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Footer -->
    <footer class="footer">
      <div class="footer-container">
        <div class="footer-grid">
          <!-- Brand -->
          <div class="footer-brand">
            <div class="footer-logo">
              <ShoppingBagIcon class="icon" />
              <span>ShopJoy</span>
            </div>
            <p class="footer-desc">品质生活，从这里开始。我们致力于为您提供优质的商品和卓越的购物体验。</p>
            <div class="social-links">
              <a href="#" class="social-link"><svg class="icon" viewBox="0 0 24 24" fill="currentColor"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg></a>
              <a href="#" class="social-link"><svg class="icon" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.163 22 16.418 22 12c0-5.523-4.477-10-10-10z"/></svg></a>
              <a href="#" class="social-link"><svg class="icon" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z"/></svg></a>
            </div>
          </div>

          <!-- Links -->
          <div class="footer-links">
            <h4>购物指南</h4>
            <ul>
              <li><a href="#">购物流程</a></li>
              <li><a href="#">会员介绍</a></li>
              <li><a href="#">常见问题</a></li>
              <li><a href="#">联系客服</a></li>
            </ul>
          </div>

          <div class="footer-links">
            <h4>配送服务</h4>
            <ul>
              <li><a href="#">配送方式</a></li>
              <li><a href="#">运费标准</a></li>
              <li><a href="#">物流跟踪</a></li>
              <li><a href="#">签收须知</a></li>
            </ul>
          </div>

          <div class="footer-links">
            <h4>售后服务</h4>
            <ul>
              <li><a href="#">退换货政策</a></li>
              <li><a href="#">退款说明</a></li>
              <li><a href="#">取消订单</a></li>
              <li><a href="#">意见反馈</a></li>
            </ul>
          </div>

          <!-- Contact -->
          <div class="footer-contact">
            <h4>联系我们</h4>
            <p class="contact-item">
              <PhoneIcon class="contact-icon" />
              400-888-8888
            </p>
            <p class="contact-item">
              <EnvelopeIcon class="contact-icon" />
              support@shopjoy.com
            </p>
            <p class="contact-item">
              <MapPinIcon class="contact-icon" />
              深圳市南山区科技园
            </p>
            <p class="contact-time">周一至周日 9:00-21:00</p>
          </div>
        </div>

        <div class="footer-bottom">
          <p>&copy; 2024 ShopJoy. All rights reserved.</p>
          <div class="footer-legal">
            <a href="#">隐私政策</a>
            <span class="divider">|</span>
            <a href="#">服务条款</a>
            <span class="divider">|</span>
            <a href="#">网站地图</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ShoppingBagIcon,
  MagnifyingGlassIcon,
  ShoppingCartIcon,
  UserCircleIcon,
  UserIcon,
  ClipboardDocumentListIcon,
  ArrowRightOnRectangleIcon,
  Bars3Icon,
  XMarkIcon,
  PhoneIcon,
  EnvelopeIcon,
  MapPinIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const isScrolled = ref(false)
const searchQuery = ref('')
const searchActive = ref(false)
const mobileMenuOpen = ref(false)
const userMenuOpen = ref(false)
const isLoggedIn = ref(false) // TODO: Replace with actual auth state
const cartCount = ref(3)

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/products', query: { search: searchQuery.value } })
  }
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const toggleUserMenu = () => {
  userMenuOpen.value = !userMenuOpen.value
}

const logout = () => {
  isLoggedIn.value = false
  userMenuOpen.value = false
  router.push('/')
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Navbar */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: #fff;
  transition: all 0.3s ease;
}

.navbar.scrolled {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.nav-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 32px;
}

/* Logo */
.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  flex-shrink: 0;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon .icon {
  width: 24px;
  height: 24px;
  color: white;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #059669;
  letter-spacing: -0.5px;
}

/* Search Bar */
.search-bar {
  flex: 1;
  max-width: 500px;
  position: relative;
}

.search-bar input {
  width: 100%;
  height: 44px;
  padding: 0 48px 0 20px;
  border: 2px solid #E5E7EB;
  border-radius: 22px;
  font-size: 15px;
  transition: all 0.3s ease;
  background: #F9FAFB;
}

.search-bar input:focus {
  outline: none;
  border-color: #059669;
  background: #fff;
  box-shadow: 0 0 0 4px rgba(5, 150, 105, 0.1);
}

.search-btn {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  background: #059669;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.search-btn:hover {
  background: #047857;
  transform: translateY(-50%) scale(1.05);
}

.search-btn .icon {
  width: 18px;
  height: 18px;
  color: white;
}

/* Nav Links */
.nav-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-link {
  padding: 10px 16px;
  color: #4B5563;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.2s;
}

.nav-link:hover {
  color: #059669;
  background: #ECFDF5;
}

.nav-link.active {
  color: #059669;
}

.nav-link.sale {
  color: #EF4444;
}

.nav-link.sale:hover {
  background: #FEF2F2;
}

/* Nav Actions */
.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.action-btn {
  width: 44px;
  height: 44px;
  border: none;
  background: transparent;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #F3F4F6;
}

.action-btn .icon {
  width: 24px;
  height: 24px;
  color: #4B5563;
}

.cart-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  background: #EF4444;
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* User Menu */
.user-menu {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  padding: 8px 0;
  min-width: 180px;
  z-index: 1001;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #374151;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s;
  border: none;
  background: none;
  width: 100%;
  cursor: pointer;
}

.dropdown-item:hover {
  background: #F3F4F6;
  color: #059669;
}

.dropdown-icon {
  width: 18px;
  height: 18px;
}

.dropdown-divider {
  height: 1px;
  background: #E5E7EB;
  margin: 8px 0;
}

.login-btn {
  padding: 10px 24px;
  background: #059669;
  color: white;
  text-decoration: none;
  font-size: 14px;
  font-weight: 600;
  border-radius: 8px;
  transition: all 0.2s;
}

.login-btn:hover {
  background: #047857;
}

/* Mobile Toggle */
.mobile-toggle {
  display: none;
  width: 44px;
  height: 44px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
}

.mobile-toggle .icon {
  width: 24px;
  height: 24px;
  color: #4B5563;
}

/* Main Content */
.main-content {
  flex: 1;
  margin-top: 72px;
}

/* Footer */
.footer {
  background: #111827;
  color: #9CA3AF;
  padding: 64px 0 32px;
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1.5fr;
  gap: 48px;
  margin-bottom: 48px;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.footer-logo .icon {
  width: 32px;
  height: 32px;
  color: #10B981;
}

.footer-logo span {
  font-size: 24px;
  font-weight: 700;
  color: white;
}

.footer-desc {
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 24px;
}

.social-links {
  display: flex;
  gap: 12px;
}

.social-link {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.social-link:hover {
  background: #059669;
}

.social-link .icon {
  width: 20px;
  height: 20px;
  color: white;
}

.footer-links h4 {
  color: white;
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 20px 0;
}

.footer-links ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links li {
  margin-bottom: 12px;
}

.footer-links a {
  color: #9CA3AF;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: #10B981;
}

.footer-contact h4 {
  color: white;
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 20px 0;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 0 12px 0;
  font-size: 14px;
}

.contact-icon {
  width: 18px;
  height: 18px;
  color: #10B981;
  flex-shrink: 0;
}

.contact-time {
  font-size: 13px;
  color: #6B7280;
  margin: 16px 0 0 0;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 32px;
  border-top: 1px solid #374151;
}

.footer-bottom p {
  margin: 0;
  font-size: 14px;
}

.footer-legal {
  display: flex;
  align-items: center;
  gap: 16px;
}

.footer-legal a {
  color: #9CA3AF;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.footer-legal a:hover {
  color: #10B981;
}

.divider {
  color: #4B5563;
}

/* Page Transitions */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* Responsive */
@media (max-width: 1024px) {
  .nav-links {
    display: none;
    position: absolute;
    top: 72px;
    left: 0;
    right: 0;
    background: white;
    flex-direction: column;
    padding: 16px 24px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }

  .nav-links.mobile-open {
    display: flex;
  }

  .nav-link {
    width: 100%;
    padding: 16px;
  }

  .mobile-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .footer-grid {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }
}

@media (max-width: 640px) {
  .nav-container {
    padding: 0 16px;
  }

  .search-bar {
    display: none;
  }

  .footer-grid {
    grid-template-columns: 1fr;
  }

  .footer-bottom {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
}
</style>
