<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-container">
        <div class="hero-content">
          <span class="hero-badge">新品上市</span>
          <h1 class="hero-title">
            发现生活
            <span class="highlight">无限可能</span>
          </h1>
          <p class="hero-description">
            精选全球好物，品质生活从这里开始。新用户注册即享 <strong>¥50</strong> 优惠券
          </p>
          <div class="hero-actions">
            <router-link to="/products" class="btn-primary">
              立即选购
              <ArrowRightIcon class="icon" />
            </router-link>
            <router-link to="/products?category=new" class="btn-secondary">
              新品首发
            </router-link>
          </div>
          <div class="hero-stats">
            <div class="stat">
              <span class="stat-value">10万+</span>
              <span class="stat-label">精选商品</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <span class="stat-value">50万+</span>
              <span class="stat-label">满意用户</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <span class="stat-value">24h</span>
              <span class="stat-label">极速发货</span>
            </div>
          </div>
        </div>
        <div class="hero-image">
          <div class="hero-image-wrapper">
            <img src="https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=600&h=600&fit=crop" alt="Shopping" />
            <div class="floating-card card-1">
              <CheckCircleIcon class="icon" />
              <span>正品保证</span>
            </div>
            <div class="floating-card card-2">
              <TruckIcon class="icon" />
              <span>免费配送</span>
            </div>
            <div class="floating-card card-3">
              <ShieldCheckIcon class="icon" />
              <span>7天退换</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories Section -->
    <section class="categories">
      <div class="section-container">
        <h2 class="section-title">热门分类</h2>
        <div class="categories-grid">
          <router-link 
            v-for="category in categories" 
            :key="category.id"
            :to="`/products?category=${category.slug}`"
            class="category-card"
            :style="{ background: category.gradient }"
          >
            <div class="category-icon">
              <component :is="category.icon" class="icon" />
            </div>
            <h3 class="category-name">{{ category.name }}</h3>
            <p class="category-count">{{ category.count }} 件商品</p>
          </router-link>
        </div>
      </div>
    </section>

    <!-- Flash Sale Section -->
    <section class="flash-sale">
      <div class="section-container">
        <div class="sale-header">
          <div class="sale-title">
            <BoltIcon class="icon" />
            <h2>限时秒杀</h2>
          </div>
          <div class="countdown">
            <span class="countdown-label">距结束</span>
            <div class="countdown-time">
              <span class="time-block">{{ countdown.hours }}</span>
              <span class="time-separator">:</span>
              <span class="time-block">{{ countdown.minutes }}</span>
              <span class="time-separator">:</span>
              <span class="time-block">{{ countdown.seconds }}</span>
            </div>
          </div>
          <router-link to="/products?type=flash" class="view-all">
            查看全部
            <ChevronRightIcon class="icon" />
          </router-link>
        </div>
        <div class="products-scroll">
          <div v-for="product in flashProducts" :key="product.id" class="product-card flash">
            <div class="product-badge">-{{ product.discount }}%</div>
            <div class="product-image">
              <img :src="product.image" :alt="product.name" />
            </div>
            <div class="product-info">
              <h3 class="product-name">{{ product.name }}</h3>
              <div class="product-price">
                <span class="current-price">¥{{ product.price }}</span>
                <span class="original-price">¥{{ product.originalPrice }}</span>
              </div>
              <div class="flash-progress">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: product.soldPercent + '%' }"></div>
                </div>
                <span class="progress-text">已售 {{ product.soldPercent }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Featured Products Section -->
    <section class="featured">
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">精选推荐</h2>
          <div class="filter-tabs">
            <button 
              v-for="tab in filterTabs" 
              :key="tab.value"
              class="tab-btn"
              :class="{ active: activeTab === tab.value }"
              @click="activeTab = tab.value"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>
        <div class="products-grid">
          <div v-for="product in filteredProducts" :key="product.id" class="product-card">
            <router-link :to="`/products/${product.id}`" class="product-link">
              <div class="product-image">
                <img :src="product.image" :alt="product.name" />
                <div v-if="product.isNew" class="product-tag new">新品</div>
                <div v-if="product.isHot" class="product-tag hot">热卖</div>
                <button class="wishlist-btn" @click.prevent="toggleWishlist(product.id)">
                  <HeartIcon class="icon" :class="{ filled: product.isWishlisted }" />
                </button>
              </div>
              <div class="product-info">
                <h3 class="product-name">{{ product.name }}</h3>
                <p class="product-desc">{{ product.description }}</p>
                <div class="product-meta">
                  <div class="product-rating">
                    <StarIcon class="icon" />
                    <span>{{ product.rating }}</span>
                    <span class="review-count">({{ product.reviews }})</span>
                  </div>
                  <div class="product-price">
                    <span class="current-price">¥{{ product.price }}</span>
                    <span v-if="product.originalPrice" class="original-price">
                      ¥{{ product.originalPrice }}
                    </span>
                  </div>
                </div>
              </div>
            </router-link>
            <button class="add-cart-btn" @click="addToCart(product)">
              <ShoppingCartIcon class="icon" />
              加入购物车
            </button>
          </div>
        </div>
        <div class="section-footer">
          <router-link to="/products" class="btn-outline">
            查看全部商品
            <ArrowRightIcon class="icon" />
          </router-link>
        </div>
      </div>
    </section>

    <!-- Benefits Section -->
    <section class="benefits">
      <div class="section-container">
        <div class="benefits-grid">
          <div v-for="benefit in benefits" :key="benefit.id" class="benefit-card">
            <div class="benefit-icon" :style="{ background: benefit.bgColor }">
              <component :is="benefit.icon" class="icon" />
            </div>
            <h3 class="benefit-title">{{ benefit.title }}</h3>
            <p class="benefit-desc">{{ benefit.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Newsletter Section -->
    <section class="newsletter">
      <div class="newsletter-container">
        <div class="newsletter-content">
          <EnvelopeIcon class="newsletter-icon" />
          <h2>订阅获取最新优惠</h2>
          <p>第一时间获取新品资讯、独家优惠和限时活动</p>
          <div class="newsletter-form">
            <input 
              v-model="email" 
              type="email" 
              placeholder="请输入您的邮箱"
              @keyup.enter="subscribe"
            />
            <button @click="subscribe" :disabled="subscribing">
              {{ subscribing ? '订阅中...' : '立即订阅' }}
            </button>
          </div>
          <p class="privacy-note">订阅即表示您同意我们的隐私政策</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  ArrowRightIcon,
  CheckCircleIcon,
  TruckIcon,
  ShieldCheckIcon,
  DevicePhoneMobileIcon,
  HomeIcon,
  ShoppingBagIcon,
  SparklesIcon,
  BoltIcon,
  ChevronRightIcon,
  HeartIcon,
  StarIcon,
  ShoppingCartIcon,
  EnvelopeIcon,
  ClockIcon,
  ArrowPathIcon,
  CreditCardIcon
} from '@heroicons/vue/24/outline'

// Countdown
const countdown = ref({ hours: '02', minutes: '45', seconds: '30' })
let countdownInterval: number

const updateCountdown = () => {
  const endTime = new Date()
  endTime.setHours(endTime.getHours() + 2)
  endTime.setMinutes(endTime.getMinutes() + 45)
  
  const update = () => {
    const now = new Date()
    const diff = endTime.getTime() - now.getTime()
    
    if (diff <= 0) {
      clearInterval(countdownInterval)
      return
    }
    
    const hours = Math.floor(diff / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    const seconds = Math.floor((diff % (1000 * 60)) / 1000)
    
    countdown.value = {
      hours: hours.toString().padStart(2, '0'),
      minutes: minutes.toString().padStart(2, '0'),
      seconds: seconds.toString().padStart(2, '0')
    }
  }
  
  update()
  countdownInterval = window.setInterval(update, 1000)
}

// Categories
const categories = ref([
  { id: 1, name: '数码电子', slug: 'electronics', count: 1256, icon: DevicePhoneMobileIcon, gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { id: 2, name: '家居生活', slug: 'home', count: 892, icon: HomeIcon, gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { id: 3, name: '时尚服饰', slug: 'fashion', count: 2156, icon: ShoppingBagIcon, gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { id: 4, name: '美妆护肤', slug: 'beauty', count: 678, icon: SparklesIcon, gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }
])

// Flash Sale Products
const flashProducts = ref([
  { id: 1, name: '无线蓝牙耳机 Pro', price: 199, originalPrice: 399, discount: 50, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=300&fit=crop', soldPercent: 78 },
  { id: 2, name: '智能手表 Series 7', price: 1299, originalPrice: 1999, discount: 35, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=300&fit=crop', soldPercent: 65 },
  { id: 3, name: '便携充电宝 20000mAh', price: 89, originalPrice: 159, discount: 44, image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=300&h=300&fit=crop', soldPercent: 92 },
  { id: 4, name: '机械键盘 RGB', price: 299, originalPrice: 459, discount: 35, image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=300&h=300&fit=crop', soldPercent: 45 },
  { id: 5, name: '4K 高清显示器', price: 1299, originalPrice: 1899, discount: 32, image: 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=300&h=300&fit=crop', soldPercent: 58 }
])

// Featured Products
const filterTabs = ref([
  { label: '综合', value: 'all' },
  { label: '新品', value: 'new' },
  { label: '热卖', value: 'hot' },
  { label: '特惠', value: 'sale' }
])
const activeTab = ref('all')

const featuredProducts = ref([
  { id: 1, name: '无线蓝牙耳机 Pro', description: '主动降噪，Hi-Fi音质', price: 299, originalPrice: 399, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=400&fit=crop', rating: 4.8, reviews: 2341, isNew: true, isHot: true, isWishlisted: false },
  { id: 2, name: '智能手表 Series 7', description: '健康监测，运动追踪', price: 1999, originalPrice: 2299, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400&h=400&fit=crop', rating: 4.9, reviews: 1856, isNew: true, isHot: false, isWishlisted: true },
  { id: 3, name: '便携充电宝 20000mAh', description: '超大容量，快充支持', price: 129, originalPrice: 159, image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=400&h=400&fit=crop', rating: 4.7, reviews: 3421, isNew: false, isHot: true, isWishlisted: false },
  { id: 4, name: '机械键盘 RGB', description: '青轴手感，炫酷背光', price: 459, originalPrice: 599, image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=400&h=400&fit=crop', rating: 4.6, reviews: 892, isNew: false, isHot: false, isWishlisted: false },
  { id: 5, name: '4K 高清显示器 27寸', description: 'IPS面板，色彩精准', price: 1599, originalPrice: 1999, image: 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400&h=400&fit=crop', rating: 4.8, reviews: 567, isNew: true, isHot: false, isWishlisted: false },
  { id: 6, name: '智能音箱 Mini', description: '语音助手，智能家居', price: 299, originalPrice: 399, image: 'https://images.unsplash.com/photo-1558089687-f282ffcbc126?w=400&h=400&fit=crop', rating: 4.5, reviews: 1234, isNew: false, isHot: true, isWishlisted: true },
  { id: 7, name: '无线充电器', description: '快充协议，安全保护', price: 99, originalPrice: 149, image: 'https://images.unsplash.com/photo-1586816879360-004f5b0c51e3?w=400&h=400&fit=crop', rating: 4.4, reviews: 2156, isNew: false, isHot: false, isWishlisted: false },
  { id: 8, name: '游戏鼠标', description: '高精度传感器，RGB灯效', price: 259, originalPrice: 359, image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400&h=400&fit=crop', rating: 4.7, reviews: 1567, isNew: true, isHot: true, isWishlisted: false }
])

const filteredProducts = computed(() => {
  if (activeTab.value === 'all') return featuredProducts.value
  if (activeTab.value === 'new') return featuredProducts.value.filter(p => p.isNew)
  if (activeTab.value === 'hot') return featuredProducts.value.filter(p => p.isHot)
  if (activeTab.value === 'sale') return featuredProducts.value.filter(p => p.originalPrice > p.price)
  return featuredProducts.value
})

// Benefits
const benefits = ref([
  { id: 1, title: '正品保障', description: '100%正品，假一赔十', icon: ShieldCheckIcon, bgColor: '#ECFDF5' },
  { id: 2, title: '极速发货', description: '24小时内发货', icon: ClockIcon, bgColor: '#EFF6FF' },
  { id: 3, title: '7天退换', description: '无理由退换货', icon: ArrowPathIcon, bgColor: '#FEF3C7' },
  { id: 4, title: '安全支付', description: '多种支付方式，安全便捷', icon: CreditCardIcon, bgColor: '#F3E8FF' }
])

// Newsletter
const email = ref('')
const subscribing = ref(false)

const toggleWishlist = (id: number) => {
  const product = featuredProducts.value.find(p => p.id === id)
  if (product) {
    product.isWishlisted = !product.isWishlisted
  }
}

const addToCart = (product: any) => {
  alert(`已将 "${product.name}" 加入购物车`)
}

const subscribe = () => {
  if (!email.value) return
  subscribing.value = true
  setTimeout(() => {
    subscribing.value = false
    email.value = ''
    alert('订阅成功！感谢您的关注。')
  }, 1500)
}

onMounted(() => {
  updateCountdown()
})

onUnmounted(() => {
  clearInterval(countdownInterval)
})
</script>

<style scoped>
.home-page {
  padding-bottom: 0;
}

/* Hero Section */
.hero {
  background: linear-gradient(135deg, #ECFDF5 0%, #D1FAE5 50%, #A7F3D0 100%);
  padding: 80px 0;
  position: relative;
  overflow: hidden;
}

.hero-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 64px;
  align-items: center;
}

.hero-content {
  max-width: 560px;
}

.hero-badge {
  display: inline-block;
  background: #059669;
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 24px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.hero-title {
  font-size: 56px;
  font-weight: 800;
  color: #111827;
  line-height: 1.1;
  margin: 0 0 24px 0;
}

.hero-title .highlight {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-description {
  font-size: 18px;
  color: #4B5563;
  line-height: 1.7;
  margin: 0 0 32px 0;
}

.hero-description strong {
  color: #059669;
  font-size: 24px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  margin-bottom: 48px;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px -10px rgba(5, 150, 105, 0.5);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 40px -10px rgba(5, 150, 105, 0.6);
}

.btn-primary .icon {
  width: 20px;
  height: 20px;
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  background: white;
  color: #059669;
  padding: 16px 32px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
  border: 2px solid #059669;
}

.btn-secondary:hover {
  background: #ECFDF5;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 24px;
}

.stat {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: #059669;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: #D1D5DB;
}

/* Hero Image */
.hero-image {
  display: flex;
  justify-content: center;
  align-items: center;
}

.hero-image-wrapper {
  position: relative;
  width: 500px;
  height: 500px;
}

.hero-image-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 30px;
  box-shadow: 0 30px 60px -20px rgba(5, 150, 105, 0.3);
}

.floating-card {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 8px;
  background: white;
  padding: 12px 20px;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  font-size: 14px;
  font-weight: 500;
  color: #111827;
  animation: float 3s ease-in-out infinite;
}

.floating-card .icon {
  width: 20px;
  height: 20px;
  color: #059669;
}

.card-1 {
  top: 20px;
  left: -20px;
  animation-delay: 0s;
}

.card-2 {
  bottom: 80px;
  left: -40px;
  animation-delay: 1s;
}

.card-3 {
  top: 50%;
  right: -30px;
  animation-delay: 2s;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* Categories Section */
.categories {
  padding: 80px 0;
  background: white;
}

.section-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.section-title {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  text-align: center;
  margin: 0 0 48px 0;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.category-card {
  padding: 32px 24px;
  border-radius: 20px;
  text-decoration: none;
  color: white;
  text-align: center;
  transition: all 0.3s ease;
}

.category-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.category-icon {
  width: 64px;
  height: 64px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.category-icon .icon {
  width: 32px;
  height: 32px;
}

.category-name {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.category-count {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

/* Flash Sale Section */
.flash-sale {
  padding: 80px 0;
  background: linear-gradient(135deg, #FEF2F2 0%, #FDF2F8 100%);
}

.sale-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.sale-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sale-title .icon {
  width: 32px;
  height: 32px;
  color: #EF4444;
}

.sale-title h2 {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.countdown {
  display: flex;
  align-items: center;
  gap: 12px;
}

.countdown-label {
  font-size: 14px;
  color: #6B7280;
}

.countdown-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.time-block {
  background: #111827;
  color: white;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 700;
  font-family: 'Fira Code', monospace;
}

.time-separator {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.view-all {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #059669;
  text-decoration: none;
  font-weight: 500;
}

.view-all .icon {
  width: 20px;
  height: 20px;
}

.products-scroll {
  display: flex;
  gap: 20px;
  overflow-x: auto;
  padding-bottom: 16px;
  scrollbar-width: thin;
}

.products-scroll::-webkit-scrollbar {
  height: 6px;
}

.products-scroll::-webkit-scrollbar-track {
  background: #F3F4F6;
  border-radius: 3px;
}

.products-scroll::-webkit-scrollbar-thumb {
  background: #D1D5DB;
  border-radius: 3px;
}

/* Product Card */
.product-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  flex-shrink: 0;
  width: 220px;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.product-card.flash {
  position: relative;
}

.product-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: #EF4444;
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  z-index: 1;
}

.product-image {
  position: relative;
  height: 180px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-tag {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.product-tag.new {
  background: #059669;
  color: white;
}

.product-tag.hot {
  background: #EF4444;
  color: white;
}

.wishlist-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 36px;
  height: 36px;
  background: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.wishlist-btn:hover {
  transform: scale(1.1);
}

.wishlist-btn .icon {
  width: 20px;
  height: 20px;
  color: #9CA3AF;
  transition: all 0.2s;
}

.wishlist-btn .icon.filled {
  color: #EF4444;
  fill: #EF4444;
}

.product-info {
  padding: 16px;
}

.product-link {
  text-decoration: none;
  color: inherit;
}

.product-name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 12px 0;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.product-rating .icon {
  width: 16px;
  height: 16px;
  color: #F59E0B;
  fill: #F59E0B;
}

.product-rating span {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.review-count {
  font-size: 12px;
  color: #9CA3AF;
  font-weight: 400;
}

.product-price {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-price {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.original-price {
  font-size: 13px;
  color: #9CA3AF;
  text-decoration: line-through;
}

.flash-progress {
  margin-top: 12px;
}

.progress-bar {
  height: 6px;
  background: #E5E7EB;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #EF4444 0%, #F87171 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #6B7280;
  margin-top: 6px;
  display: block;
}

.add-cart-btn {
  width: calc(100% - 32px);
  margin: 0 16px 16px;
  padding: 12px;
  background: #059669;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-cart-btn:hover {
  background: #047857;
}

.add-cart-btn .icon {
  width: 18px;
  height: 18px;
}

/* Featured Section */
.featured {
  padding: 80px 0;
  background: white;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  background: #F3F4F6;
  padding: 4px;
  border-radius: 10px;
}

.tab-btn {
  padding: 10px 20px;
  border: none;
  background: transparent;
  color: #6B7280;
  font-size: 14px;
  font-weight: 500;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: white;
  color: #059669;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.section-footer {
  text-align: center;
  margin-top: 48px;
}

.btn-outline {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 28px;
  border: 2px solid #059669;
  color: #059669;
  text-decoration: none;
  border-radius: 10px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-outline:hover {
  background: #059669;
  color: white;
}

.btn-outline .icon {
  width: 18px;
  height: 18px;
}

/* Benefits Section */
.benefits {
  padding: 80px 0;
  background: #F9FAFB;
}

.benefits-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.benefit-card {
  text-align: center;
  padding: 32px 24px;
}

.benefit-icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
}

.benefit-icon .icon {
  width: 32px;
  height: 32px;
  color: #059669;
}

.benefit-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.benefit-desc {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

/* Newsletter Section */
.newsletter {
  padding: 80px 0;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
}

.newsletter-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 0 24px;
  text-align: center;
}

.newsletter-icon {
  width: 64px;
  height: 64px;
  color: white;
  margin-bottom: 24px;
}

.newsletter h2 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0 0 16px 0;
}

.newsletter p {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 32px 0;
}

.newsletter-form {
  display: flex;
  gap: 12px;
  max-width: 480px;
  margin: 0 auto;
}

.newsletter-form input {
  flex: 1;
  padding: 16px 24px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  outline: none;
}

.newsletter-form input::placeholder {
  color: #9CA3AF;
}

.newsletter-form button {
  padding: 16px 32px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.newsletter-form button:hover:not(:disabled) {
  background: #1F2937;
  transform: translateY(-2px);
}

.newsletter-form button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.privacy-note {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 16px;
}

/* Responsive */
@media (max-width: 1024px) {
  .hero-container {
    grid-template-columns: 1fr;
    text-align: center;
  }

  .hero-content {
    max-width: 100%;
  }

  .hero-stats {
    justify-content: center;
  }

  .hero-image {
    display: none;
  }

  .categories-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .products-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .benefits-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .hero {
    padding: 60px 0;
  }

  .hero-title {
    font-size: 36px;
  }

  .hero-actions {
    flex-direction: column;
  }

  .categories-grid {
    grid-template-columns: 1fr;
  }

  .sale-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .products-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .section-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .benefits-grid {
    grid-template-columns: 1fr;
  }

  .newsletter-form {
    flex-direction: column;
  }
}
</style>
