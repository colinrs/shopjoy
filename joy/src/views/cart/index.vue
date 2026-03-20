<template>
  <div class="cart-page">
    <div class="page-container">
      <!-- Page Header -->
      <div class="page-header">
        <h1 class="page-title">购物车</h1>
        <p class="page-subtitle">共 {{ cartItems.length }} 件商品</p>
      </div>

      <div v-if="cartItems.length === 0" class="empty-cart">
        <div class="empty-icon">
          <ShoppingCartIcon class="icon" />
        </div>
        <h2>购物车是空的</h2>
        <p>快去挑选心仪的商品吧</p>
        <router-link to="/products" class="btn-primary">
          去购物
          <ArrowRightIcon class="icon" />
        </router-link>
      </div>

      <div v-else class="cart-content">
        <!-- Cart Items -->
        <div class="cart-main">
          <!-- Select All Header -->
          <div class="cart-header">
            <label class="checkbox-label">
              <input 
                type="checkbox" 
                :checked="isAllSelected"
                @change="toggleSelectAll"
              />
              <span class="checkmark"></span>
              <span class="label-text">全选</span>
            </label>
            <button class="btn-text" @click="clearCart">清空购物车</button>
          </div>

          <!-- Cart Items List -->
          <div class="cart-items">
            <div 
              v-for="item in cartItems" 
              :key="item.id"
              class="cart-item"
              :class="{ 'out-of-stock': item.stock === 0 }"
            >
              <label class="checkbox-label item-checkbox">
                <input 
                  type="checkbox" 
                  v-model="item.selected"
                  :disabled="item.stock === 0"
                />
                <span class="checkmark"></span>
              </label>

              <router-link :to="`/products/${item.productId}`" class="item-image">
                <img :src="item.image" :alt="item.name" />
                <span v-if="item.stock === 0" class="stock-badge">缺货</span>
              </router-link>

              <div class="item-details">
                <router-link :to="`/products/${item.productId}`" class="item-name">
                  {{ item.name }}
                </router-link>
                <p class="item-variant">{{ item.variant }}</p>
                <div class="item-price-mobile">
                  <span class="current-price">¥{{ item.price }}</span>
                  <span v-if="item.originalPrice" class="original-price">
                    ¥{{ item.originalPrice }}
                  </span>
                </div>
              </div>

              <div class="item-quantity">
                <div class="quantity-control">
                  <button 
                    @click="decreaseQuantity(item)"
                    :disabled="item.quantity <= 1"
                  >
                    <MinusIcon class="icon" />
                  </button>
                  <input 
                    v-model.number="item.quantity"
                    type="number"
                    min="1"
                    :max="item.stock"
                    @change="updateQuantity(item)"
                  />
                  <button 
                    @click="increaseQuantity(item)"
                    :disabled="item.quantity >= item.stock"
                  >
                    <PlusIcon class="icon" />
                  </button>
                </div>
                <p v-if="item.stock < 10 && item.stock > 0" class="stock-warning">
                  仅剩 {{ item.stock }} 件
                </p>
              </div>

              <div class="item-price">
                <span class="subtotal">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
              </div>

              <button class="remove-btn" @click="removeItem(item)">
                <TrashIcon class="icon" />
              </button>
            </div>
          </div>
        </div>

        <!-- Cart Summary -->
        <div class="cart-sidebar">
          <div class="summary-card">
            <h3 class="summary-title">订单 summary</h3>

            <div class="promo-section">
              <div class="promo-input">
                <input 
                  v-model="promoCode"
                  type="text"
                  placeholder="输入优惠码"
                  :disabled="promoApplied"
                />
                <button 
                  @click="applyPromoCode"
                  :disabled="!promoCode || promoApplied"
                >
                  {{ promoApplied ? '已应用' : '应用' }}
                </button>
              </div>
              <p v-if="promoError" class="promo-error">{{ promoError }}</p>
            </div>

            <div class="summary-rows">
              <div class="summary-row">
                <span>商品小计</span>
                <span>¥{{ subtotal.toFixed(2) }}</span>
              </div>
              <div class="summary-row">
                <span>运费</span>
                <span class="free" v-if="shipping === 0">免运费</span>
                <span v-else>¥{{ shipping.toFixed(2) }}</span>
              </div>
              <div v-if="discount > 0" class="summary-row discount">
                <span>优惠折扣</span>
                <span>-¥{{ discount.toFixed(2) }}</span>
              </div>
            </div>

            <div class="summary-total">
              <span>合计</span>
              <span class="total-price">¥{{ total.toFixed(2) }}</span>
            </div>

            <button 
              class="btn-checkout"
              :disabled="selectedItems.length === 0"
              @click="checkout"
            >
              去结算
              <span v-if="selectedItems.length > 0">({{ selectedItems.length }})</span>
            </button>

            <div class="secure-payment">
              <ShieldCheckIcon class="icon" />
              <span>安全支付保障</span>
            </div>
          </div>

          <!-- Recommendations -->
          <div class="recommendations">
            <h4>猜你喜欢</h4>
            <div class="rec-list">
              <router-link 
                v-for="item in recommendations" 
                :key="item.id"
                :to="`/products/${item.id}`"
                class="rec-item"
              >
                <img :src="item.image" :alt="item.name" />
                <div class="rec-info">
                  <p class="rec-name">{{ item.name }}</p>
                  <p class="rec-price">¥{{ item.price }}</p>
                </div>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  ShoppingCartIcon,
  ArrowRightIcon,
  MinusIcon,
  PlusIcon,
  TrashIcon,
  ShieldCheckIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()

const cartItems = ref([
  {
    id: 1,
    productId: 1,
    name: '无线蓝牙耳机 Pro',
    variant: '深邃黑',
    price: 299,
    originalPrice: 399,
    quantity: 1,
    stock: 156,
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=200&fit=crop',
    selected: true
  },
  {
    id: 2,
    productId: 2,
    name: '智能手表 Series 7',
    variant: '银色/44mm',
    price: 1999,
    originalPrice: 2299,
    quantity: 1,
    stock: 89,
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=200&h=200&fit=crop',
    selected: true
  },
  {
    id: 3,
    productId: 3,
    name: '便携充电宝 20000mAh',
    variant: '白色',
    price: 129,
    originalPrice: 159,
    quantity: 2,
    stock: 234,
    image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=200&h=200&fit=crop',
    selected: false
  }
])

const promoCode = ref('')
const promoApplied = ref(false)
const promoError = ref('')

const isAllSelected = computed(() => {
  const availableItems = cartItems.value.filter(item => item.stock > 0)
  return availableItems.length > 0 && availableItems.every(item => item.selected)
})

const selectedItems = computed(() => cartItems.value.filter(item => item.selected && item.stock > 0))

const subtotal = computed(() => {
  return selectedItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const shipping = computed(() => subtotal.value >= 199 ? 0 : 10)

const discount = computed(() => {
  if (promoApplied.value) {
    return subtotal.value * 0.1 // 10% discount
  }
  return 0
})

const total = computed(() => subtotal.value + shipping.value - discount.value)

const recommendations = ref([
  { id: 4, name: '机械键盘 RGB', price: 459, image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=150&h=150&fit=crop' },
  { id: 5, name: '无线充电器', price: 99, image: 'https://images.unsplash.com/photo-1586816879360-004f5b0c51e3?w=150&h=150&fit=crop' },
  { id: 6, name: '游戏鼠标', price: 259, image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=150&h=150&fit=crop' }
])

const toggleSelectAll = () => {
  const newValue = !isAllSelected.value
  cartItems.value.forEach(item => {
    if (item.stock > 0) {
      item.selected = newValue
    }
  })
}

const increaseQuantity = (item: any) => {
  if (item.quantity < item.stock) {
    item.quantity++
  }
}

const decreaseQuantity = (item: any) => {
  if (item.quantity > 1) {
    item.quantity--
  }
}

const updateQuantity = (item: any) => {
  if (item.quantity < 1) item.quantity = 1
  if (item.quantity > item.stock) item.quantity = item.stock
}

const removeItem = (item: any) => {
  const index = cartItems.value.findIndex(i => i.id === item.id)
  if (index > -1) {
    cartItems.value.splice(index, 1)
  }
}

const clearCart = () => {
  if (confirm('确定要清空购物车吗？')) {
    cartItems.value = []
  }
}

const applyPromoCode = () => {
  if (promoCode.value.toUpperCase() === 'SAVE10') {
    promoApplied.value = true
    promoError.value = ''
  } else {
    promoError.value = '无效的优惠码'
  }
}

const checkout = () => {
  router.push('/checkout')
}
</script>

<style scoped>
.cart-page {
  min-height: calc(100vh - 72px);
  background: #F9FAFB;
  padding: 40px 0;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.page-header {
  margin-bottom: 32px;
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

/* Empty Cart */
.empty-cart {
  text-align: center;
  padding: 80px 24px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  width: 120px;
  height: 120px;
  background: #F3F4F6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24px;
}

.empty-icon .icon {
  width: 60px;
  height: 60px;
  color: #9CA3AF;
}

.empty-cart h2 {
  font-size: 24px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-cart p {
  font-size: 16px;
  color: #6B7280;
  margin: 0 0 24px 0;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 16px 32px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

.btn-primary .icon {
  width: 20px;
  height: 20px;
}

/* Cart Content */
.cart-content {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
}

/* Cart Main */
.cart-main {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #E5E7EB;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.checkbox-label input {
  display: none;
}

.checkmark {
  width: 20px;
  height: 20px;
  border: 2px solid #D1D5DB;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.checkbox-label input:checked + .checkmark {
  background: #059669;
  border-color: #059669;
}

.checkbox-label input:checked + .checkmark::after {
  content: '';
  width: 6px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.label-text {
  font-size: 15px;
  font-weight: 500;
  color: #374151;
}

.btn-text {
  font-size: 14px;
  color: #6B7280;
  background: none;
  border: none;
  cursor: pointer;
  transition: color 0.2s;
}

.btn-text:hover {
  color: #EF4444;
}

/* Cart Items */
.cart-items {
  padding: 0 24px;
}

.cart-item {
  display: grid;
  grid-template-columns: auto 100px 1fr 140px 100px auto;
  align-items: center;
  gap: 20px;
  padding: 24px 0;
  border-bottom: 1px solid #E5E7EB;
}

.cart-item:last-child {
  border-bottom: none;
}

.cart-item.out-of-stock {
  opacity: 0.6;
}

.item-checkbox .checkmark {
  border-color: #D1D5DB;
}

.item-checkbox input:disabled + .checkmark {
  opacity: 0.5;
  cursor: not-allowed;
}

.item-image {
  position: relative;
  aspect-ratio: 1;
  border-radius: 12px;
  overflow: hidden;
  background: #F3F4F6;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.stock-badge {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.item-details {
  min-width: 0;
}

.item-name {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  text-decoration: none;
  margin-bottom: 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-name:hover {
  color: #059669;
}

.item-variant {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.item-price-mobile {
  display: none;
  margin-top: 8px;
}

.item-quantity {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.quantity-control {
  display: flex;
  align-items: center;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  overflow: hidden;
}

.quantity-control button {
  width: 36px;
  height: 36px;
  border: none;
  background: #F9FAFB;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.quantity-control button:hover:not(:disabled) {
  background: #F3F4F6;
}

.quantity-control button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity-control button .icon {
  width: 16px;
  height: 16px;
  color: #374151;
}

.quantity-control input {
  width: 50px;
  height: 36px;
  border: none;
  text-align: center;
  font-size: 14px;
  font-weight: 600;
  outline: none;
}

.stock-warning {
  font-size: 12px;
  color: #F59E0B;
  margin: 0;
}

.item-price {
  text-align: right;
}

.subtotal {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.remove-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.remove-btn:hover {
  background: #FEF2F2;
}

.remove-btn .icon {
  width: 20px;
  height: 20px;
  color: #9CA3AF;
}

.remove-btn:hover .icon {
  color: #EF4444;
}

/* Cart Sidebar */
.cart-sidebar {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.summary-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.summary-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 20px 0;
}

.promo-section {
  margin-bottom: 24px;
}

.promo-input {
  display: flex;
  gap: 8px;
}

.promo-input input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #E5E7EB;
  border-radius: 10px;
  font-size: 14px;
  outline: none;
}

.promo-input input:focus {
  border-color: #059669;
}

.promo-input input:disabled {
  background: #F9FAFB;
}

.promo-input button {
  padding: 12px 20px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.promo-input button:hover:not(:disabled) {
  background: #374151;
}

.promo-input button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.promo-error {
  font-size: 13px;
  color: #EF4444;
  margin: 8px 0 0 0;
}

.summary-rows {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #E5E7EB;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  font-size: 15px;
  color: #4B5563;
}

.summary-row .free {
  color: #059669;
  font-weight: 600;
}

.summary-row.discount {
  color: #EF4444;
}

.summary-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.summary-total span:first-child {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.total-price {
  font-size: 28px;
  font-weight: 800;
  color: #EF4444;
}

.btn-checkout {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-checkout:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(5, 150, 105, 0.5);
}

.btn-checkout:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.secure-payment {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  font-size: 13px;
  color: #6B7280;
}

.secure-payment .icon {
  width: 16px;
  height: 16px;
  color: #059669;
}

/* Recommendations */
.recommendations {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.recommendations h4 {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 16px 0;
}

.rec-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rec-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  transition: background 0.2s;
  text-decoration: none;
}

.rec-item:hover {
  background: #F9FAFB;
}

.rec-item img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
}

.rec-info {
  flex: 1;
  min-width: 0;
}

.rec-name {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin: 0 0 4px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.rec-price {
  font-size: 15px;
  font-weight: 700;
  color: #EF4444;
  margin: 0;
}

/* Responsive */
@media (max-width: 1024px) {
  .cart-content {
    grid-template-columns: 1fr;
  }

  .cart-sidebar {
    flex-direction: row;
  }

  .summary-card {
    flex: 1;
  }

  .recommendations {
    display: none;
  }
}

@media (max-width: 768px) {
  .cart-item {
    grid-template-columns: auto 80px 1fr auto;
    grid-template-rows: auto auto;
    gap: 12px;
  }

  .item-details {
    grid-column: 3 / -1;
  }

  .item-quantity {
    grid-column: 2;
    grid-row: 2;
  }

  .item-price {
    display: none;
  }

  .item-price-mobile {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .current-price {
    font-size: 16px;
    font-weight: 700;
    color: #EF4444;
  }

  .original-price {
    font-size: 13px;
    color: #9CA3AF;
    text-decoration: line-through;
  }

  .remove-btn {
    grid-column: 4;
    grid-row: 2;
  }

  .cart-sidebar {
    flex-direction: column;
  }
}
</style>