<template>
  <div class="product-detail-page">
    <!-- Breadcrumb -->
    <div class="breadcrumb-container">
      <div class="breadcrumb">
        <router-link to="/">首页</router-link>
        <ChevronRightIcon class="separator" />
        <router-link to="/products">商品列表</router-link>
        <ChevronRightIcon class="separator" />
        <span>{{ product.name }}</span>
      </div>
    </div>

    <div class="page-container">
      <!-- Product Gallery -->
      <div class="product-gallery">
        <div class="main-image">
          <img :src="currentImage" :alt="product.name" />
          <div class="image-actions">
            <button class="action-btn" @click="zoomImage">
              <MagnifyingGlassPlusIcon class="icon" />
            </button>
            <button class="action-btn" @click="shareProduct">
              <ShareIcon class="icon" />
            </button>
          </div>
        </div>
        <div class="thumbnail-list">
          <button 
            v-for="(image, index) in product.images" 
            :key="index"
            class="thumbnail"
            :class="{ active: currentImage === image }"
            @click="currentImage = image"
          >
            <img :src="image" :alt="`${product.name} - ${index + 1}`" />
          </button>
        </div>
      </div>

      <!-- Product Info -->
      <div class="product-info">
        <!-- Tags -->
        <div class="product-tags">
          <span v-if="product.isNew" class="tag new">新品</span>
          <span v-if="product.isHot" class="tag hot">热卖</span>
          <span v-if="product.discount" class="tag discount">-{{ product.discount }}%</span>
        </div>

        <h1 class="product-title">{{ product.name }}</h1>
        
        <div class="product-meta">
          <div class="rating">
            <div class="stars">
              <StarIcon v-for="n in 5" :key="n" class="star" :class="{ filled: n <= Math.floor(product.rating) }" />
            </div>
            <span class="rating-score">{{ product.rating }}</span>
            <span class="reviews">({{ product.reviews }} 条评价)</span>
          </div>
          <div class="sales">已售 {{ product.sales }}+</div>
        </div>

        <p class="product-description">{{ product.description }}</p>

        <!-- Price -->
        <div class="price-section">
          <div class="price-row">
            <span class="current-price">¥{{ currentPrice }}</span>
            <span v-if="product.originalPrice" class="original-price">¥{{ product.originalPrice }}</span>
            <span v-if="product.discount" class="discount-badge">省 ¥{{ product.originalPrice - currentPrice }}</span>
          </div>
          <div class="promo-tags">
            <span class="promo-tag">满199减30</span>
            <span class="promo-tag">包邮</span>
            <span class="promo-tag">7天无理由</span>
          </div>
        </div>

        <!-- Variants -->
        <div v-if="product.variants" class="variants-section">
          <div v-for="variant in product.variants" :key="variant.name" class="variant-group">
            <h4>{{ variant.name }}</h4>
            <div class="variant-options">
              <button 
                v-for="option in variant.options" 
                :key="option.value"
                class="variant-btn"
                :class="{ 
                  active: selectedVariants[variant.name] === option.value,
                  disabled: option.disabled 
                }"
                :disabled="option.disabled"
                @click="selectVariant(variant.name, option.value)"
              >
                <img v-if="option.image" :src="option.image" :alt="option.label" />
                <span>{{ option.label }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Quantity -->
        <div class="quantity-section">
          <h4>数量</h4>
          <div class="quantity-selector">
            <button @click="quantity > 1 && quantity--" :disabled="quantity <= 1">
              <MinusIcon class="icon" />
            </button>
            <input v-model.number="quantity" type="number" min="1" :max="product.stock" />
            <button @click="quantity < product.stock && quantity++" :disabled="quantity >= product.stock">
              <PlusIcon class="icon" />
            </button>
          </div>
          <span class="stock-hint">库存 {{ product.stock }} 件</span>
        </div>

        <!-- Actions -->
        <div class="action-buttons">
          <button class="btn-wishlist" @click="toggleWishlist" :class="{ active: isWishlisted }">
            <HeartIcon class="icon" />
            {{ isWishlisted ? '已收藏' : '收藏' }}
          </button>
          <button class="btn-cart" @click="addToCart">
            <ShoppingCartIcon class="icon" />
            加入购物车
          </button>
          <button class="btn-buy" @click="buyNow">
            立即购买
          </button>
        </div>

        <!-- Services -->
        <div class="services-section">
          <div class="service-item">
            <ShieldCheckIcon class="icon" />
            <span>正品保证</span>
          </div>
          <div class="service-item">
            <TruckIcon class="icon" />
            <span>极速发货</span>
          </div>
          <div class="service-item">
            <ArrowPathIcon class="icon" />
            <span>7天退换</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Product Details Tabs -->
    <div class="details-section">
      <div class="tabs-container">
        <div class="tabs-header">
          <button 
            v-for="tab in tabs" 
            :key="tab.id"
            class="tab-btn"
            :class="{ active: activeTab === tab.id }"
            @click="activeTab = tab.id"
          >
            {{ tab.name }}
          </button>
        </div>

        <div class="tab-content">
          <!-- Description Tab -->
          <div v-if="activeTab === 'description'" class="description-content">
            <div class="feature-image">
              <img :src="product.detailImage" :alt="product.name" />
            </div>
            <div class="detail-text">
              <h3>商品详情</h3>
              <p>{{ product.detailDescription }}</p>
              <h4>产品特点</h4>
              <ul>
                <li v-for="(feature, index) in product.features" :key="index">
                  <CheckCircleIcon class="icon" />
                  {{ feature }}
                </li>
              </ul>
              <h4>规格参数</h4>
              <table class="specs-table">
                <tr v-for="(spec, index) in product.specs" :key="index">
                  <th>{{ spec.name }}</th>
                  <td>{{ spec.value }}</td>
                </tr>
              </table>
            </div>
          </div>

          <!-- Reviews Tab -->
          <div v-if="activeTab === 'reviews'" class="reviews-content">
            <div class="reviews-summary">
              <div class="rating-overview">
                <div class="big-rating">{{ product.rating }}</div>
                <div class="stars">
                  <StarIcon v-for="n in 5" :key="n" class="star" :class="{ filled: n <= Math.round(product.rating) }" />
                </div>
                <div class="total-reviews">{{ product.reviews }} 条评价</div>
              </div>
              <div class="rating-bars">
                <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="rating-bar">
                  <span class="star-label">{{ star }}星</span>
                  <div class="progress-bar">
                    <div class="progress" :style="{ width: getRatingPercent(star) + '%' }"></div>
                  </div>
                  <span class="percent">{{ getRatingPercent(star) }}%</span>
                </div>
              </div>
            </div>

            <div class="reviews-list">
              <div v-for="review in reviews" :key="review.id" class="review-item">
                <div class="review-header">
                  <div class="reviewer">
                    <img :src="review.avatar" :alt="review.name" class="avatar" />
                    <div class="info">
                      <span class="name">{{ review.name }}</span>
                      <div class="stars">
                        <StarIcon v-for="n in 5" :key="n" class="star" :class="{ filled: n <= review.rating }" />
                      </div>
                    </div>
                  </div>
                  <span class="date">{{ review.date }}</span>
                </div>
                <p class="review-content">{{ review.content }}</p>
                <div v-if="review.images" class="review-images">
                  <img v-for="(img, idx) in review.images" :key="idx" :src="img" @click="previewImage(img)" />
                </div>
                <div class="review-footer">
                  <span class="variant">{{ review.variant }}</span>
                  <button class="helpful-btn" @click="likeReview(review)">
                    <HandThumbUpIcon class="icon" />
                    有用 ({{ review.helpful }})
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- FAQ Tab -->
          <div v-if="activeTab === 'faq'" class="faq-content">
            <div class="faq-list">
              <div v-for="(item, index) in faqs" :key="index" class="faq-item">
                <button class="faq-question" @click="toggleFaq(index)">
                  <span>Q: {{ item.question }}</span>
                  <ChevronDownIcon class="icon" :class="{ rotated: openFaq === index }" />
                </button>
                <div class="faq-answer" :class="{ open: openFaq === index }">
                  <p>A: {{ item.answer }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Related Products -->
      <div class="related-section">
        <h3 class="section-title">相关推荐</h3>
        <div class="related-grid">
          <div v-for="item in relatedProducts" :key="item.id" class="related-card">
            <router-link :to="`/products/${item.id}`">
              <div class="related-image">
                <img :src="item.image" :alt="item.name" />
              </div>
              <h4 class="related-name">{{ item.name }}</h4>
              <div class="related-price">
                <span class="price">¥{{ item.price }}</span>
                <span v-if="item.originalPrice" class="original">¥{{ item.originalPrice }}</span>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ChevronRightIcon,
  StarIcon,
  HeartIcon,
  ShoppingCartIcon,
  MinusIcon,
  PlusIcon,
  ShieldCheckIcon,
  TruckIcon,
  ArrowPathIcon,
  CheckCircleIcon,
  HandThumbUpIcon,
  ChevronDownIcon,
  MagnifyingGlassPlusIcon,
  ShareIcon
} from '@heroicons/vue/24/solid'

const route = useRoute()
const router = useRouter()

// Product data
const product = ref({
  id: 1,
  name: '无线蓝牙耳机 Pro - 主动降噪版',
  description: '采用先进的主动降噪技术，带来沉浸式音频体验。40mm大动圈单元，Hi-Fi音质，支持AAC/SBC高清解码。30小时超长续航，快充15分钟可使用3小时。',
  price: 299,
  originalPrice: 399,
  discount: 25,
  rating: 4.8,
  reviews: 2341,
  sales: 5678,
  stock: 156,
  isNew: true,
  isHot: true,
  images: [
    'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=800&h=800&fit=crop',
    'https://images.unsplash.com/photo-1484704849700-f032a568e944?w=800&h=800&fit=crop',
    'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=800&h=800&fit=crop',
    'https://images.unsplash.com/photo-1583394838336-acd977736f90?w=800&h=800&fit=crop'
  ],
  detailImage: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=1200&h=600&fit=crop',
  detailDescription: '这款无线蓝牙耳机采用业界领先的主动降噪技术，能够有效消除环境噪音，让您沉浸在纯粹的音乐世界中。精选40mm大动圈单元，配合专业调音，呈现出丰富细腻的音质表现。',
  features: [
    '主动降噪技术，降噪深度达35dB',
    '40mm大动圈单元，Hi-Fi级音质',
    '蓝牙5.0稳定连接，低延迟',
    '30小时超长续航，支持快充',
    '佩戴舒适，人体工学设计',
    '触控操作，便捷智能'
  ],
  specs: [
    { name: '品牌', value: 'ShopJoy' },
    { name: '型号', value: 'BT-Pro-2024' },
    { name: '蓝牙版本', value: '5.0' },
    { name: '驱动单元', value: '40mm' },
    { name: '频响范围', value: '20Hz-20kHz' },
    { name: '阻抗', value: '32Ω' },
    { name: '续航时间', value: '30小时' },
    { name: '充电时间', value: '2小时' },
    { name: '重量', value: '250g' }
  ],
  variants: [
    {
      name: '颜色',
      options: [
        { value: 'black', label: '深邃黑', image: 'https://via.placeholder.com/40/000000' },
        { value: 'white', label: '纯净白', image: 'https://via.placeholder.com/40/FFFFFF' },
        { value: 'blue', label: '天际蓝', image: 'https://via.placeholder.com/40/3B82F6' }
      ]
    }
  ]
})

const currentImage = ref(product.value.images[0])
const quantity = ref(1)
const isWishlisted = ref(false)
const selectedVariants = ref<Record<string, string>>({})
const activeTab = ref('description')
const openFaq = ref<number | null>(null)

const tabs = [
  { id: 'description', name: '商品详情' },
  { id: 'reviews', name: `用户评价 (${product.value.reviews})` },
  { id: 'faq', name: '常见问题' }
]

const currentPrice = computed(() => {
  return product.value.price
})

const reviews = ref([
  {
    id: 1,
    name: '张三',
    avatar: 'https://via.placeholder.com/40',
    rating: 5,
    date: '2024-03-15',
    content: '音质非常好，降噪效果惊艳！在地铁上完全听不到噪音，佩戴也很舒适，连续戴了4个小时耳朵也不疼。强烈推荐！',
    variant: '深邃黑',
    helpful: 128,
    images: [
      'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=200&fit=crop',
      'https://images.unsplash.com/photo-1484704849700-f032a568e944?w=200&h=200&fit=crop'
    ]
  },
  {
    id: 2,
    name: '李四',
    avatar: 'https://via.placeholder.com/40',
    rating: 5,
    date: '2024-03-10',
    content: '做工精致，音质清晰，低音有力。连接速度很快，续航能力也很强，充一次电可以用好几天。',
    variant: '纯净白',
    helpful: 86
  },
  {
    id: 3,
    name: '王五',
    avatar: 'https://via.placeholder.com/40',
    rating: 4,
    date: '2024-03-05',
    content: '总体不错，降噪效果还可以，音质对得起这个价格。就是白色版本容易脏，建议买深色。',
    variant: '纯净白',
    helpful: 45
  }
])

const faqs = ref([
  { question: '这款耳机支持哪些设备？', answer: '支持所有带有蓝牙功能的设备，包括iOS、Android手机、平板、电脑等。' },
  { question: '降噪功能可以关闭吗？', answer: '可以，您可以通过耳机上的按键或配套APP切换降噪模式、通透模式或关闭降噪。' },
  { question: '支持多设备连接吗？', answer: '支持，可以同时连接两台设备，自动切换音频来源。' },
  { question: '质保多久？', answer: '提供一年质保，非人为损坏免费维修或更换。' }
])

const relatedProducts = ref([
  { id: 2, name: '智能手表 Series 7', price: 1999, originalPrice: 2299, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=300&fit=crop' },
  { id: 3, name: '便携充电宝 20000mAh', price: 129, originalPrice: 159, image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=300&h=300&fit=crop' },
  { id: 4, name: '无线充电器', price: 99, originalPrice: 149, image: 'https://images.unsplash.com/photo-1586816879360-004f5b0c51e3?w=300&h=300&fit=crop' }
])

const selectVariant = (name: string, value: string) => {
  selectedVariants.value[name] = value
}

const toggleWishlist = () => {
  isWishlisted.value = !isWishlisted.value
}

const addToCart = () => {
  alert(`已将 ${quantity.value} 件商品加入购物车`)
}

const buyNow = () => {
  router.push('/checkout')
}

const zoomImage = () => {
  // Open image zoom modal
}

const shareProduct = () => {
  if (navigator.share) {
    navigator.share({
      title: product.value.name,
      text: product.value.description,
      url: window.location.href
    })
  } else {
    // Copy to clipboard
    navigator.clipboard.writeText(window.location.href)
    alert('链接已复制到剪贴板')
  }
}

const getRatingPercent = (star: number) => {
  const distribution = { 5: 75, 4: 18, 3: 5, 2: 1, 1: 1 }
  return distribution[star as keyof typeof distribution] || 0
}

const likeReview = (review: any) => {
  review.helpful++
}

const toggleFaq = (index: number) => {
  openFaq.value = openFaq.value === index ? null : index
}

const previewImage = (img: string) => {
  // Open image preview modal
}

onMounted(() => {
  // Load product data based on route params
  const productId = route.params.id
  console.log('Loading product:', productId)
})
</script>

<style scoped>
.product-detail-page {
  padding-bottom: 80px;
}

/* Breadcrumb */
.breadcrumb-container {
  background: #F9FAFB;
  border-bottom: 1px solid #E5E7EB;
  padding: 16px 0;
}

.breadcrumb {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6B7280;
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

/* Page Container */
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 24px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 64px;
}

/* Product Gallery */
.product-gallery {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.main-image {
  position: relative;
  aspect-ratio: 1;
  background: #F3F4F6;
  border-radius: 20px;
  overflow: hidden;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-actions {
  position: absolute;
  top: 16px;
  right: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-btn {
  width: 44px;
  height: 44px;
  background: white;
  border: none;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.action-btn:hover {
  transform: scale(1.05);
}

.action-btn .icon {
  width: 20px;
  height: 20px;
  color: #374151;
}

.thumbnail-list {
  display: flex;
  gap: 12px;
}

.thumbnail {
  width: 80px;
  height: 80px;
  border: 2px solid transparent;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background: #F3F4F6;
}

.thumbnail.active {
  border-color: #059669;
}

.thumbnail:hover {
  border-color: #10B981;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Product Info */
.product-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.product-tags {
  display: flex;
  gap: 8px;
}

.tag {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.tag.new {
  background: #ECFDF5;
  color: #059669;
}

.tag.hot {
  background: #FEF2F2;
  color: #EF4444;
}

.tag.discount {
  background: #FEF3C7;
  color: #D97706;
}

.product-title {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  margin: 0;
  line-height: 1.3;
}

.product-meta {
  display: flex;
  align-items: center;
  gap: 24px;
}

.rating {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stars {
  display: flex;
  gap: 4px;
}

.star {
  width: 18px;
  height: 18px;
  color: #D1D5DB;
}

.star.filled {
  color: #F59E0B;
}

.rating-score {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.reviews {
  font-size: 14px;
  color: #059669;
  cursor: pointer;
}

.sales {
  font-size: 14px;
  color: #6B7280;
}

.product-description {
  font-size: 15px;
  color: #4B5563;
  line-height: 1.7;
  margin: 0;
}

/* Price Section */
.price-section {
  padding: 24px;
  background: linear-gradient(135deg, #FEF2F2 0%, #FDF2F8 100%);
  border-radius: 16px;
}

.price-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.current-price {
  font-size: 36px;
  font-weight: 800;
  color: #EF4444;
}

.original-price {
  font-size: 20px;
  color: #9CA3AF;
  text-decoration: line-through;
}

.discount-badge {
  padding: 6px 12px;
  background: #EF4444;
  color: white;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.promo-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.promo-tag {
  padding: 6px 12px;
  background: white;
  border: 1px solid #FECACA;
  border-radius: 6px;
  font-size: 12px;
  color: #DC2626;
}

/* Variants */
.variants-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.variant-group h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.variant-options {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.variant-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: 2px solid #E5E7EB;
  background: white;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.variant-btn:hover:not(.disabled) {
  border-color: #10B981;
}

.variant-btn.active {
  border-color: #059669;
  background: #ECFDF5;
}

.variant-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.variant-btn img {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  object-fit: cover;
}

.variant-btn span {
  font-size: 14px;
  font-weight: 500;
}

/* Quantity Section */
.quantity-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.quantity-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0;
}

.quantity-selector {
  display: flex;
  align-items: center;
  border: 1px solid #E5E7EB;
  border-radius: 10px;
  overflow: hidden;
}

.quantity-selector button {
  width: 44px;
  height: 44px;
  border: none;
  background: #F3F4F6;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.quantity-selector button:hover:not(:disabled) {
  background: #E5E7EB;
}

.quantity-selector button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity-selector button .icon {
  width: 16px;
  height: 16px;
  color: #374151;
}

.quantity-selector input {
  width: 60px;
  height: 44px;
  border: none;
  text-align: center;
  font-size: 16px;
  font-weight: 600;
  outline: none;
}

.stock-hint {
  font-size: 13px;
  color: #6B7280;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 12px;
}

.btn-wishlist {
  padding: 16px 24px;
  border: 2px solid #E5E7EB;
  background: white;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-wishlist:hover {
  border-color: #EF4444;
  color: #EF4444;
}

.btn-wishlist.active {
  border-color: #EF4444;
  background: #FEF2F2;
  color: #EF4444;
}

.btn-wishlist.active .icon {
  fill: #EF4444;
}

.btn-wishlist .icon {
  width: 20px;
  height: 20px;
}

.btn-cart {
  flex: 1;
  padding: 16px 32px;
  background: white;
  border: 2px solid #059669;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #059669;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cart:hover {
  background: #ECFDF5;
}

.btn-cart .icon {
  width: 20px;
  height: 20px;
}

.btn-buy {
  flex: 1;
  padding: 16px 32px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-buy:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

/* Services Section */
.services-section {
  display: flex;
  gap: 24px;
  padding: 20px;
  background: #F9FAFB;
  border-radius: 12px;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #4B5563;
}

.service-item .icon {
  width: 20px;
  height: 20px;
  color: #059669;
}

/* Details Section */
.details-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.tabs-container {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.tabs-header {
  display: flex;
  border-bottom: 1px solid #E5E7EB;
}

.tab-btn {
  flex: 1;
  padding: 20px;
  background: none;
  border: none;
  font-size: 15px;
  font-weight: 500;
  color: #6B7280;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.tab-btn:hover {
  color: #059669;
  background: #F9FAFB;
}

.tab-btn.active {
  color: #059669;
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #059669, #10B981);
}

.tab-content {
  padding: 40px;
}

/* Description Content */
.description-content {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.feature-image img {
  width: 100%;
  border-radius: 16px;
}

.detail-text h3 {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 20px 0;
}

.detail-text h4 {
  font-size: 18px;
  font-weight: 600;
  color: #374151;
  margin: 32px 0 16px 0;
}

.detail-text p {
  font-size: 15px;
  color: #4B5563;
  line-height: 1.8;
  margin: 0 0 20px 0;
}

.detail-text ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-text li {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 15px;
  color: #374151;
}

.detail-text li .icon {
  width: 20px;
  height: 20px;
  color: #059669;
}

.specs-table {
  width: 100%;
  border-collapse: collapse;
}

.specs-table th,
.specs-table td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid #E5E7EB;
}

.specs-table th {
  width: 200px;
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}

.specs-table td {
  color: #374151;
}

/* Reviews Content */
.reviews-summary {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 64px;
  padding: 32px;
  background: #F9FAFB;
  border-radius: 16px;
  margin-bottom: 40px;
}

.rating-overview {
  text-align: center;
}

.big-rating {
  font-size: 64px;
  font-weight: 800;
  color: #111827;
  line-height: 1;
}

.rating-overview .stars {
  justify-content: center;
  margin: 12px 0;
}

.rating-overview .star {
  width: 24px;
  height: 24px;
}

.total-reviews {
  font-size: 14px;
  color: #6B7280;
}

.rating-bars {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rating-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.star-label {
  width: 40px;
  font-size: 14px;
  color: #6B7280;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.progress {
  height: 100%;
  background: linear-gradient(90deg, #F59E0B, #FBBF24);
  border-radius: 4px;
  transition: width 0.3s;
}

.percent {
  width: 50px;
  text-align: right;
  font-size: 14px;
  color: #6B7280;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.review-item {
  padding: 24px;
  background: #F9FAFB;
  border-radius: 16px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.reviewer {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.info .stars {
  gap: 2px;
}

.info .star {
  width: 14px;
  height: 14px;
}

.date {
  font-size: 13px;
  color: #9CA3AF;
}

.review-content {
  font-size: 15px;
  color: #374151;
  line-height: 1.7;
  margin: 0 0 16px 0;
}

.review-images {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.review-images img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.2s;
}

.review-images img:hover {
  transform: scale(1.05);
}

.review-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.variant {
  font-size: 13px;
  color: #9CA3AF;
}

.helpful-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  font-size: 13px;
  color: #6B7280;
  cursor: pointer;
  transition: all 0.2s;
}

.helpful-btn:hover {
  border-color: #059669;
  color: #059669;
}

.helpful-btn .icon {
  width: 16px;
  height: 16px;
}

/* FAQ Content */
.faq-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.faq-item {
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  overflow: hidden;
}

.faq-question {
  width: 100%;
  padding: 20px 24px;
  background: white;
  border: none;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 15px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.faq-question:hover {
  background: #F9FAFB;
}

.faq-question .icon {
  width: 20px;
  height: 20px;
  color: #6B7280;
  transition: transform 0.3s;
}

.faq-question .icon.rotated {
  transform: rotate(180deg);
}

.faq-answer {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
}

.faq-answer.open {
  max-height: 200px;
}

.faq-answer p {
  padding: 0 24px 20px;
  margin: 0;
  font-size: 15px;
  color: #6B7280;
  line-height: 1.7;
}

/* Related Section */
.related-section {
  margin-top: 64px;
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 24px 0;
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.related-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

.related-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.related-card a {
  text-decoration: none;
  color: inherit;
  display: block;
}

.related-image {
  aspect-ratio: 1;
  overflow: hidden;
}

.related-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.related-card:hover .related-image img {
  transform: scale(1.05);
}

.related-name {
  font-size: 15px;
  font-weight: 500;
  color: #111827;
  margin: 16px 16px 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.related-price {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 16px 16px;
}

.related-price .price {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.related-price .original {
  font-size: 14px;
  color: #9CA3AF;
  text-decoration: line-through;
}

/* Responsive */
@media (max-width: 1024px) {
  .page-container {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .product-gallery {
    max-width: 600px;
    margin: 0 auto;
  }

  .reviews-summary {
    grid-template-columns: 1fr;
    gap: 32px;
  }

  .related-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .product-title {
    font-size: 24px;
  }

  .current-price {
    font-size: 28px;
  }

  .action-buttons {
    flex-direction: column;
  }

  .services-section {
    flex-wrap: wrap;
  }

  .tab-content {
    padding: 20px;
  }

  .related-grid {
    grid-template-columns: 1fr;
  }
}
</style>
