<template>
  <div class="products-page">
    <!-- Breadcrumb & Header -->
    <div class="page-header">
      <div class="header-container">
        <nav class="breadcrumb">
          <router-link to="/">首页</router-link>
          <ChevronRightIcon class="separator" />
          <span>全部商品</span>
        </nav>
        <h1 class="page-title">全部商品</h1>
        <p class="page-subtitle">共 {{ totalProducts }} 件精选好物</p>
      </div>
    </div>

    <div class="page-container">
      <!-- Sidebar Filters -->
      <aside class="sidebar" :class="{ 'mobile-open': mobileFiltersOpen }">
        <div class="sidebar-header">
          <h3>筛选条件</h3>
          <button class="close-filters" @click="mobileFiltersOpen = false">
            <XMarkIcon class="icon" />
          </button>
        </div>

        <!-- Categories -->
        <div class="filter-section">
          <h4 class="filter-title">商品分类</h4>
          <div class="filter-options">
            <label v-for="cat in categories" :key="cat.id" class="filter-option">
              <input 
                type="checkbox" 
                v-model="selectedCategories" 
                :value="cat.slug"
              />
              <span class="checkbox"></span>
              <span class="label">{{ cat.name }}</span>
              <span class="count">({{ cat.count }})</span>
            </label>
          </div>
        </div>

        <!-- Price Range -->
        <div class="filter-section">
          <h4 class="filter-title">价格区间</h4>
          <div class="price-range">
            <input 
              v-model.number="priceRange.min" 
              type="number" 
              placeholder="最低价"
              class="price-input"
            />
            <span class="range-separator">-</span>
            <input 
              v-model.number="priceRange.max" 
              type="number" 
              placeholder="最高价"
              class="price-input"
            />
          </div>
          <div class="price-presets">
            <button 
              v-for="preset in pricePresets" 
              :key="preset.label"
              class="preset-btn"
              :class="{ active: isPricePresetActive(preset) }"
              @click="applyPricePreset(preset)"
            >
              {{ preset.label }}
            </button>
          </div>
        </div>

        <!-- Rating Filter -->
        <div class="filter-section">
          <h4 class="filter-title">评分</h4>
          <div class="filter-options">
            <label v-for="rating in ratings" :key="rating.value" class="filter-option">
              <input 
                type="radio" 
                v-model="selectedRating" 
                :value="rating.value"
                name="rating"
              />
              <span class="radio"></span>
              <div class="rating-label">
                <div class="stars">
                  <StarIcon 
                    v-for="n in 5" 
                    :key="n"
                    class="star"
                    :class="{ filled: n <= rating.stars }"
                  />
                </div>
                <span class="and-up">{{ rating.label }}</span>
              </div>
            </label>
          </div>
        </div>

        <!-- Apply Filters -->
        <div class="filter-actions">
          <button class="btn-primary" @click="applyFilters">应用筛选</button>
          <button class="btn-text" @click="resetFilters">重置</button>
        </div>
      </aside>

      <!-- Main Content -->
      <main class="main-content">
        <!-- Toolbar -->
        <div class="toolbar">
          <div class="toolbar-left">
            <button class="filter-toggle" @click="mobileFiltersOpen = true">
              <FunnelIcon class="icon" />
              筛选
            </button>
            <div class="active-filters" v-if="activeFiltersCount > 0">
              <span class="filters-count">{{ activeFiltersCount }} 个筛选</span>
              <button class="clear-filters" @click="resetFilters">清除</button>
            </div>
          </div>
          <div class="toolbar-right">
            <div class="sort-dropdown">
              <select v-model="sortBy" @change="handleSort">
                <option value="default">默认排序</option>
                <option value="price-asc">价格从低到高</option>
                <option value="price-desc">价格从高到低</option>
                <option value="sales">销量优先</option>
                <option value="rating">评分最高</option>
                <option value="newest">最新上架</option>
              </select>
              <ChevronDownIcon class="dropdown-icon" />
            </div>
            <div class="view-toggle">
              <button 
                class="view-btn" 
                :class="{ active: viewMode === 'grid' }"
                @click="viewMode = 'grid'"
              >
                <Squares2X2Icon class="icon" />
              </button>
              <button 
                class="view-btn" 
                :class="{ active: viewMode === 'list' }"
                @click="viewMode = 'list'"
              >
                <Bars3Icon class="icon" />
              </button>
            </div>
          </div>
        </div>

        <!-- Products Grid -->
        <div class="products-grid" :class="viewMode">
          <div 
            v-for="product in filteredProducts" 
            :key="product.id" 
            class="product-card"
            :class="viewMode"
          >
            <router-link :to="`/products/${product.id}`" class="product-link">
              <div class="product-image">
                <img :src="product.image" :alt="product.name" loading="lazy" />
                <div v-if="product.discount" class="discount-badge">-{{ product.discount }}%</div>
                <div v-if="product.isNew" class="badge new">新品</div>
                <div v-if="product.isHot" class="badge hot">热卖</div>
                <button 
                  class="wishlist-btn" 
                  @click.prevent="toggleWishlist(product)"
                  :class="{ active: product.isWishlisted }"
                >
                  <HeartIcon class="icon" />
                </button>
                <button class="quick-view-btn" @click.prevent="openQuickView(product)">
                  <EyeIcon class="icon" />
                  快速查看
                </button>
              </div>
              <div class="product-info">
                <div class="product-category">{{ product.category }}</div>
                <h3 class="product-name">{{ product.name }}</h3>
                <p class="product-description">{{ product.description }}</p>
                <div class="product-meta">
                  <div class="product-rating">
                    <StarIcon class="star" />
                    <span>{{ product.rating }}</span>
                    <span class="reviews">({{ product.reviews }})</span>
                  </div>
                  <div class="product-sales">已售 {{ product.sales }}</div>
                </div>
                <div class="product-price-row">
                  <div class="price-wrapper">
                    <span class="current-price">¥{{ product.price }}</span>
                    <span v-if="product.originalPrice" class="original-price">
                      ¥{{ product.originalPrice }}
                    </span>
                  </div>
                  <button class="add-cart-btn" @click.prevent="addToCart(product)">
                    <ShoppingCartIcon class="icon" />
                  </button>
                </div>
              </div>
            </router-link>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="filteredProducts.length === 0" class="empty-state">
          <img src="/empty-products.svg" alt="No products" class="empty-image" />
          <h3>没有找到相关商品</h3>
          <p>试试调整筛选条件或搜索其他关键词</p>
          <button class="btn-primary" @click="resetFilters">清除筛选</button>
        </div>

        <!-- Pagination -->
        <div class="pagination" v-if="totalPages > 1">
          <button 
            class="page-btn" 
            :disabled="currentPage === 1"
            @click="currentPage--"
          >
            <ChevronLeftIcon class="icon" />
          </button>
          <div class="page-numbers">
            <button 
              v-for="page in displayedPages" 
              :key="page"
              class="page-number"
              :class="{ active: page === currentPage }"
              @click="currentPage = page"
            >
              {{ page }}
            </button>
          </div>
          <button 
            class="page-btn" 
            :disabled="currentPage === totalPages"
            @click="currentPage++"
          >
            <ChevronRightIcon class="icon" />
          </button>
        </div>
      </main>
    </div>

    <!-- Quick View Modal -->
    <div v-if="quickViewProduct" class="modal-overlay" @click="closeQuickView">
      <div class="quick-view-modal" @click.stop>
        <button class="close-modal" @click="closeQuickView">
          <XMarkIcon class="icon" />
        </button>
        <div class="modal-content">
          <div class="modal-image">
            <img :src="quickViewProduct.image" :alt="quickViewProduct.name" />
          </div>
          <div class="modal-info">
            <span class="modal-category">{{ quickViewProduct.category }}</span>
            <h2 class="modal-title">{{ quickViewProduct.name }}</h2>
            <p class="modal-description">{{ quickViewProduct.description }}</p>
            <div class="modal-rating">
              <div class="stars">
                <StarIcon 
                  v-for="n in 5" 
                  :key="n"
                  class="star"
                  :class="{ filled: n <= Math.floor(quickViewProduct.rating) }"
                />
              </div>
              <span class="rating-text">{{ quickViewProduct.rating }} ({{ quickViewProduct.reviews }} 评价)</span>
            </div>
            <div class="modal-price">
              <span class="current-price">¥{{ quickViewProduct.price }}</span>
              <span v-if="quickViewProduct.originalPrice" class="original-price">
                ¥{{ quickViewProduct.originalPrice }}
              </span>
            </div>
            <div class="modal-actions">
              <div class="quantity-selector">
                <button @click="quantity > 1 && quantity--">-</button>
                <input v-model.number="quantity" type="number" min="1" max="99" />
                <button @click="quantity++">+</button>
              </div>
              <button class="btn-primary add-to-cart" @click="addToCartWithQuantity">
                <ShoppingCartIcon class="icon" />
                加入购物车
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ChevronRightIcon,
  XMarkIcon,
  FunnelIcon,
  ChevronDownIcon,
  Squares2X2Icon,
  Bars3Icon,
  StarIcon,
  HeartIcon,
  EyeIcon,
  ShoppingCartIcon,
  ChevronLeftIcon
} from '@heroicons/vue/24/solid'
import { StarIcon as StarOutlineIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const router = useRouter()

// State
const mobileFiltersOpen = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')
const currentPage = ref(1)
const itemsPerPage = 12
const totalProducts = ref(156)
const sortBy = ref('default')
const quickViewProduct = ref<any>(null)
const quantity = ref(1)

// Filters
const selectedCategories = ref<string[]>([])
const priceRange = ref({ min: null as number | null, max: null as number | null })
const selectedRating = ref<number | null>(null)

const categories = ref([
  { id: 1, name: '数码电子', slug: 'electronics', count: 45 },
  { id: 2, name: '服装配饰', slug: 'clothing', count: 38 },
  { id: 3, name: '家居生活', slug: 'home', count: 27 },
  { id: 4, name: '运动户外', slug: 'sports', count: 19 },
  { id: 5, name: '美妆护肤', slug: 'beauty', count: 15 },
  { id: 6, name: '食品饮料', slug: 'food', count: 12 }
])

const pricePresets = ref([
  { label: '¥0-100', min: 0, max: 100 },
  { label: '¥100-500', min: 100, max: 500 },
  { label: '¥500-1000', min: 500, max: 1000 },
  { label: '¥1000+', min: 1000, max: null }
])

const ratings = ref([
  { value: 4.5, stars: 5, label: '4.5分以上' },
  { value: 4.0, stars: 4, label: '4.0分以上' },
  { value: 3.0, stars: 3, label: '3.0分以上' }
])

// Mock products data
const products = ref([
  {
    id: 1,
    name: '无线蓝牙耳机 Pro',
    description: '主动降噪，Hi-Fi音质，超长续航',
    price: 299,
    originalPrice: 399,
    discount: 25,
    category: '数码电子',
    rating: 4.8,
    reviews: 2341,
    sales: 5678,
    isNew: true,
    isHot: true,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=400&fit=crop'
  },
  {
    id: 2,
    name: '智能手表 Series 7',
    description: '健康监测，运动追踪，时尚设计',
    price: 1999,
    originalPrice: 2299,
    discount: 13,
    category: '数码电子',
    rating: 4.9,
    reviews: 1856,
    sales: 3421,
    isNew: true,
    isHot: false,
    isWishlisted: true,
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400&h=400&fit=crop'
  },
  {
    id: 3,
    name: '便携充电宝 20000mAh',
    description: '超大容量，快充支持，轻薄便携',
    price: 129,
    originalPrice: 159,
    discount: 19,
    category: '数码电子',
    rating: 4.7,
    reviews: 3421,
    sales: 8921,
    isNew: false,
    isHot: true,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=400&h=400&fit=crop'
  },
  {
    id: 4,
    name: '机械键盘 RGB',
    description: '青轴手感，炫酷背光，游戏办公两相宜',
    price: 459,
    originalPrice: 599,
    discount: 23,
    category: '数码电子',
    rating: 4.6,
    reviews: 892,
    sales: 2134,
    isNew: false,
    isHot: false,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=400&h=400&fit=crop'
  },
  {
    id: 5,
    name: '4K 高清显示器 27寸',
    description: 'IPS面板，色彩精准，护眼模式',
    price: 1599,
    originalPrice: 1999,
    discount: 20,
    category: '数码电子',
    rating: 4.8,
    reviews: 567,
    sales: 1234,
    isNew: true,
    isHot: false,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400&h=400&fit=crop'
  },
  {
    id: 6,
    name: '智能音箱 Mini',
    description: '语音助手，智能家居控制，高品质音响',
    price: 299,
    originalPrice: 399,
    discount: 25,
    category: '数码电子',
    rating: 4.5,
    reviews: 1234,
    sales: 4567,
    isNew: false,
    isHot: true,
    isWishlisted: true,
    image: 'https://images.unsplash.com/photo-1558089687-f282ffcbc126?w=400&h=400&fit=crop'
  },
  {
    id: 7,
    name: '无线充电器',
    description: '快充协议，安全保护，兼容多设备',
    price: 99,
    originalPrice: 149,
    discount: 34,
    category: '数码电子',
    rating: 4.4,
    reviews: 2156,
    sales: 6789,
    isNew: false,
    isHot: false,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1586816879360-004f5b0c51e3?w=400&h=400&fit=crop'
  },
  {
    id: 8,
    name: '游戏鼠标',
    description: '高精度传感器，RGB灯效，人体工学设计',
    price: 259,
    originalPrice: 359,
    discount: 28,
    category: '数码电子',
    rating: 4.7,
    reviews: 1567,
    sales: 3456,
    isNew: true,
    isHot: true,
    isWishlisted: false,
    image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400&h=400&fit=crop'
  }
])

// Computed
const activeFiltersCount = computed(() => {
  let count = 0
  if (selectedCategories.value.length > 0) count++
  if (priceRange.value.min !== null || priceRange.value.max !== null) count++
  if (selectedRating.value !== null) count++
  return count
})

const filteredProducts = computed(() => {
  let result = [...products.value]

  // Filter by category
  if (selectedCategories.value.length > 0) {
    result = result.filter(p => selectedCategories.value.includes(p.category))
  }

  // Filter by price
  if (priceRange.value.min !== null) {
    result = result.filter(p => p.price >= priceRange.value.min!)
  }
  if (priceRange.value.max !== null) {
    result = result.filter(p => p.price <= priceRange.value.max!)
  }

  // Filter by rating
  if (selectedRating.value !== null) {
    result = result.filter(p => p.rating >= selectedRating.value!)
  }

  // Sort
  switch (sortBy.value) {
    case 'price-asc':
      result.sort((a, b) => a.price - b.price)
      break
    case 'price-desc':
      result.sort((a, b) => b.price - a.price)
      break
    case 'sales':
      result.sort((a, b) => b.sales - a.sales)
      break
    case 'rating':
      result.sort((a, b) => b.rating - a.rating)
      break
  }

  return result
})

const totalPages = computed(() => Math.ceil(filteredProducts.value.length / itemsPerPage))

const displayedPages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  
  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
const isPricePresetActive = (preset: any) => {
  return priceRange.value.min === preset.min && priceRange.value.max === preset.max
}

const applyPricePreset = (preset: any) => {
  priceRange.value = { min: preset.min, max: preset.max }
}

const applyFilters = () => {
  currentPage.value = 1
  mobileFiltersOpen.value = false
}

const resetFilters = () => {
  selectedCategories.value = []
  priceRange.value = { min: null, max: null }
  selectedRating.value = null
  currentPage.value = 1
}

const handleSort = () => {
  currentPage.value = 1
}

const toggleWishlist = (product: any) => {
  product.isWishlisted = !product.isWishlisted
}

const openQuickView = (product: any) => {
  quickViewProduct.value = product
  quantity.value = 1
}

const closeQuickView = () => {
  quickViewProduct.value = null
}

const addToCart = (product: any) => {
  alert(`已将 "${product.name}" 加入购物车`)
}

const addToCartWithQuantity = () => {
  if (quickViewProduct.value) {
    alert(`已将 ${quantity.value} 件 "${quickViewProduct.value.name}" 加入购物车`)
    closeQuickView()
  }
}

// Watch for route query changes
watch(() => route.query, (newQuery) => {
  if (newQuery.category) {
    selectedCategories.value = [newQuery.category as string]
  }
  if (newQuery.search) {
    // Handle search
  }
}, { immediate: true })
</script>

<style scoped>
.products-page {
  min-height: 100vh;
  background: #F9FAFB;
}

/* Page Header */
.page-header {
  background: white;
  border-bottom: 1px solid #E5E7EB;
  padding: 32px 0;
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 16px;
}

.breadcrumb a {
  color: #059669;
  text-decoration: none;
}

.breadcrumb a:hover {
  text-decoration: underline;
}

.separator {
  width: 16px;
  height: 16px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 16px;
  color: #6B7280;
  margin: 0;
}

/* Page Container */
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 24px;
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 32px;
}

/* Sidebar */
.sidebar {
  background: white;
  border-radius: 16px;
  padding: 24px;
  height: fit-content;
  position: sticky;
  top: 96px;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.sidebar-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.close-filters {
  display: none;
  width: 32px;
  height: 32px;
  border: none;
  background: #F3F4F6;
  border-radius: 8px;
  cursor: pointer;
}

.filter-section {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #E5E7EB;
}

.filter-section:last-of-type {
  border-bottom: none;
}

.filter-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 16px 0;
}

.filter-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filter-option {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  font-size: 14px;
  color: #4B5563;
}

.filter-option input {
  display: none;
}

.checkbox, .radio {
  width: 18px;
  height: 18px;
  border: 2px solid #D1D5DB;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.radio {
  border-radius: 50%;
}

.filter-option input:checked + .checkbox {
  background: #059669;
  border-color: #059669;
}

.filter-option input:checked + .checkbox::after {
  content: '';
  width: 8px;
  height: 8px;
  background: white;
  border-radius: 2px;
}

.filter-option input:checked + .radio {
  background: #059669;
  border-color: #059669;
}

.filter-option input:checked + .radio::after {
  content: '';
  width: 8px;
  height: 8px;
  background: white;
  border-radius: 50%;
}

.filter-option .label {
  flex: 1;
}

.filter-option .count {
  color: #9CA3AF;
  font-size: 12px;
}

.rating-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stars {
  display: flex;
  gap: 2px;
}

.star {
  width: 16px;
  height: 16px;
  color: #D1D5DB;
}

.star.filled {
  color: #F59E0B;
}

.and-up {
  font-size: 12px;
  color: #6B7280;
}

/* Price Range */
.price-range {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.price-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
}

.price-input:focus {
  border-color: #059669;
}

.range-separator {
  color: #6B7280;
}

.price-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preset-btn {
  padding: 8px 12px;
  border: 1px solid #E5E7EB;
  background: white;
  border-radius: 6px;
  font-size: 12px;
  color: #4B5563;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-btn:hover, .preset-btn.active {
  border-color: #059669;
  color: #059669;
  background: #ECFDF5;
}

/* Filter Actions */
.filter-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.btn-primary {
  flex: 1;
  padding: 12px 24px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

.btn-text {
  padding: 12px 16px;
  background: transparent;
  color: #6B7280;
  border: none;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.2s;
}

.btn-text:hover {
  color: #059669;
}

/* Main Content */
.main-content {
  min-width: 0;
}

/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 24px;
  background: white;
  border-radius: 12px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.filter-toggle {
  display: none;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #F3F4F6;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
}

.filter-toggle .icon {
  width: 18px;
  height: 18px;
}

.active-filters {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filters-count {
  font-size: 14px;
  color: #059669;
  font-weight: 500;
}

.clear-filters {
  font-size: 13px;
  color: #6B7280;
  background: none;
  border: none;
  cursor: pointer;
  text-decoration: underline;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.sort-dropdown {
  position: relative;
}

.sort-dropdown select {
  appearance: none;
  padding: 10px 40px 10px 16px;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  font-size: 14px;
  color: #374151;
  background: white;
  cursor: pointer;
  outline: none;
}

.sort-dropdown select:focus {
  border-color: #059669;
}

.dropdown-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 16px;
  color: #6B7280;
  pointer-events: none;
}

.view-toggle {
  display: flex;
  gap: 4px;
}

.view-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #E5E7EB;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.view-btn.active {
  background: #059669;
  border-color: #059669;
}

.view-btn.active .icon {
  color: white;
}

.view-btn .icon {
  width: 20px;
  height: 20px;
  color: #6B7280;
}

/* Products Grid */
.products-grid {
  display: grid;
  gap: 24px;
}

.products-grid.grid {
  grid-template-columns: repeat(3, 1fr);
}

.products-grid.list {
  grid-template-columns: 1fr;
}

.product-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.product-card.list {
  display: grid;
  grid-template-columns: 240px 1fr;
}

.product-link {
  text-decoration: none;
  color: inherit;
}

.product-image {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
  background: #F3F4F6;
}

.product-card.list .product-image {
  aspect-ratio: auto;
  height: 100%;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.discount-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: #EF4444;
  color: white;
  padding: 6px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
  z-index: 2;
}

.badge {
  position: absolute;
  top: 12px;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  z-index: 2;
}

.badge.new {
  left: 12px;
  background: #059669;
  color: white;
}

.badge.hot {
  right: 12px;
  background: #F59E0B;
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
  opacity: 0;
  transform: translateY(-10px);
  transition: all 0.3s ease;
  z-index: 3;
}

.product-card:hover .wishlist-btn {
  opacity: 1;
  transform: translateY(0);
}

.wishlist-btn.active {
  opacity: 1;
  transform: translateY(0);
}

.wishlist-btn .icon {
  width: 18px;
  height: 18px;
  color: #9CA3AF;
  transition: all 0.2s;
}

.wishlist-btn.active .icon,
.wishlist-btn:hover .icon {
  color: #EF4444;
  fill: #EF4444;
}

.quick-view-btn {
  position: absolute;
  bottom: 12px;
  left: 50%;
  transform: translateX(-50%) translateY(20px);
  padding: 10px 20px;
  background: white;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  opacity: 0;
  transition: all 0.3s ease;
  z-index: 3;
}

.product-card:hover .quick-view-btn {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.quick-view-btn .icon {
  width: 16px;
  height: 16px;
}

.product-info {
  padding: 20px;
}

.product-category {
  font-size: 12px;
  color: #059669;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-description {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 12px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.product-rating .star {
  width: 14px;
  height: 14px;
  color: #F59E0B;
}

.product-rating span {
  font-size: 13px;
  font-weight: 600;
  color: #111827;
}

.reviews {
  font-size: 12px;
  color: #9CA3AF;
  font-weight: 400;
}

.product-sales {
  font-size: 12px;
  color: #6B7280;
}

.product-price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-price {
  font-size: 20px;
  font-weight: 700;
  color: #EF4444;
}

.original-price {
  font-size: 14px;
  color: #9CA3AF;
  text-decoration: line-through;
}

.add-cart-btn {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border: none;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.add-cart-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 16px -4px rgba(5, 150, 105, 0.4);
}

.add-cart-btn .icon {
  width: 20px;
  height: 20px;
  color: white;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 24px;
}

.empty-image {
  width: 200px;
  height: 200px;
  margin-bottom: 24px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 24px 0;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 48px;
  padding: 24px;
  background: white;
  border-radius: 12px;
}

.page-btn {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #E5E7EB;
  background: white;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: #059669;
  color: #059669;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-btn .icon {
  width: 20px;
  height: 20px;
}

.page-numbers {
  display: flex;
  gap: 8px;
}

.page-number {
  min-width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #E5E7EB;
  background: white;
  border-radius: 10px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.page-number:hover {
  border-color: #059669;
  color: #059669;
}

.page-number.active {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border-color: #059669;
  color: white;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
  backdrop-filter: blur(4px);
}

.quick-view-modal {
  background: white;
  border-radius: 24px;
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  overflow: hidden;
  position: relative;
  animation: modalIn 0.3s ease;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.close-modal {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  background: #F3F4F6;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.2s;
}

.close-modal:hover {
  background: #E5E7EB;
}

.close-modal .icon {
  width: 20px;
  height: 20px;
  color: #6B7280;
}

.modal-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  min-height: 500px;
}

.modal-image {
  background: #F3F4F6;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.modal-image img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.modal-info {
  padding: 40px;
  display: flex;
  flex-direction: column;
}

.modal-category {
  font-size: 13px;
  color: #059669;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 12px;
}

.modal-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 16px 0;
  line-height: 1.3;
}

.modal-description {
  font-size: 15px;
  color: #6B7280;
  line-height: 1.6;
  margin: 0 0 24px 0;
}

.modal-rating {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.modal-rating .stars {
  display: flex;
  gap: 4px;
}

.modal-rating .star {
  width: 20px;
  height: 20px;
  color: #D1D5DB;
}

.modal-rating .star.filled {
  color: #F59E0B;
}

.rating-text {
  font-size: 14px;
  color: #6B7280;
}

.modal-price {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #E5E7EB;
}

.modal-price .current-price {
  font-size: 36px;
  font-weight: 700;
  color: #EF4444;
}

.modal-price .original-price {
  font-size: 20px;
  color: #9CA3AF;
  text-decoration: line-through;
}

.modal-actions {
  display: flex;
  gap: 16px;
  margin-top: auto;
}

.quantity-selector {
  display: flex;
  align-items: center;
  border: 1px solid #E5E7EB;
  border-radius: 10px;
  overflow: hidden;
}

.quantity-selector button {
  width: 48px;
  height: 52px;
  border: none;
  background: #F3F4F6;
  font-size: 20px;
  color: #374151;
  cursor: pointer;
  transition: background 0.2s;
}

.quantity-selector button:hover {
  background: #E5E7EB;
}

.quantity-selector input {
  width: 60px;
  height: 52px;
  border: none;
  text-align: center;
  font-size: 16px;
  font-weight: 600;
  outline: none;
}

.modal-actions .add-to-cart {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 16px;
}

.modal-actions .add-to-cart .icon {
  width: 22px;
  height: 22px;
}

/* Responsive */
@media (max-width: 1024px) {
  .page-container {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: -100%;
    width: 320px;
    height: 100vh;
    z-index: 1000;
    border-radius: 0;
    transition: left 0.3s ease;
    overflow-y: auto;
  }

  .sidebar.mobile-open {
    left: 0;
  }

  .close-filters {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .filter-toggle {
    display: flex;
  }

  .products-grid.grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .modal-content {
    grid-template-columns: 1fr;
  }

  .modal-image {
    height: 300px;
  }
}

@media (max-width: 640px) {
  .products-grid.grid {
    grid-template-columns: 1fr;
  }

  .product-card.list {
    grid-template-columns: 1fr;
  }

  .toolbar {
    flex-direction: column;
    gap: 16px;
  }

  .toolbar-left, .toolbar-right {
    width: 100%;
    justify-content: space-between;
  }

  .modal-actions {
    flex-direction: column;
  }

  .quantity-selector {
    justify-content: center;
  }
}
</style>
