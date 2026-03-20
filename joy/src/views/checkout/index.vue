<template>
  <div class="checkout-page">
    <div class="page-container">
      <!-- Checkout Steps -->
      <div class="checkout-steps">
        <div class="step active">
          <div class="step-number">1</div>
          <span class="step-label">确认订单</span>
        </div>
        <div class="step-line"></div>
        <div class="step">
          <div class="step-number">2</div>
          <span class="step-label">支付</span>
        </div>
        <div class="step-line"></div>
        <div class="step">
          <div class="step-number">3</div>
          <span class="step-label">完成</span>
        </div>
      </div>

      <div class="checkout-content">
        <!-- Left Column -->
        <div class="checkout-main">
          <!-- Address Section -->
          <section class="checkout-section">
            <div class="section-header">
              <h3>收货地址</h3>
              <button class="btn-text" @click="showAddressModal = true">
                <PlusIcon class="icon" />
                新增地址
              </button>
            </div>
            <div class="address-list">
              <div 
                v-for="addr in addresses" 
                :key="addr.id"
                class="address-card"
                :class="{ active: selectedAddress?.id === addr.id }"
                @click="selectedAddress = addr"
              >
                <div class="address-radio">
                  <div class="radio-indicator" :class="{ checked: selectedAddress?.id === addr.id }"></div>
                </div>
                <div class="address-info">
                  <div class="address-header">
                    <span class="name">{{ addr.name }}</span>
                    <span class="phone">{{ addr.phone }}</span>
                    <span v-if="addr.isDefault" class="default-tag">默认</span>
                  </div>
                  <p class="address-detail">{{ addr.province }} {{ addr.city }} {{ addr.district }} {{ addr.detail }}</p>
                </div>
                <button class="edit-btn" @click.stop="editAddress(addr)">
                  <PencilIcon class="icon" />
                </button>
              </div>
            </div>
          </section>

          <!-- Order Items -->
          <section class="checkout-section">
            <h3>商品清单</h3>
            <div class="order-items">
              <div v-for="item in orderItems" :key="item.id" class="order-item">
                <img :src="item.image" :alt="item.name" class="item-image" />
                <div class="item-info">
                  <h4 class="item-name">{{ item.name }}</h4>
                  <p class="item-variant">{{ item.variant }}</p>
                </div>
                <div class="item-price">
                  <span class="price">¥{{ item.price }}</span>
                  <span class="quantity">x{{ item.quantity }}</span>
                </div>
                <span class="item-subtotal">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
              </div>
            </div>
          </section>

          <!-- Shipping Method -->
          <section class="checkout-section">
            <h3>配送方式</h3>
            <div class="shipping-options">
              <label 
                v-for="method in shippingMethods" 
                :key="method.id"
                class="shipping-card"
                :class="{ active: selectedShipping === method.id }"
              >
                <input type="radio" v-model="selectedShipping" :value="method.id" />
                <div class="shipping-content">
                  <div class="shipping-header">
                    <span class="method-name">{{ method.name }}</span>
                    <span class="method-price">
                      {{ method.price === 0 ? '免运费' : '¥' + method.price }}
                    </span>
                  </div>
                  <p class="method-desc">{{ method.description }}</p>
                </div>
              </label>
            </div>
          </section>

          <!-- Coupon -->
          <section class="checkout-section">
            <h3>优惠券</h3>
            <div class="coupon-section">
              <div v-if="appliedCoupon" class="applied-coupon">
                <div class="coupon-info">
                  <span class="coupon-name">{{ appliedCoupon.name }}</span>
                  <span class="coupon-discount">-¥{{ appliedCoupon.discount }}</span>
                </div>
                <button class="remove-coupon" @click="appliedCoupon = null">
                  <XMarkIcon class="icon" />
                </button>
              </div>
              <button v-else class="select-coupon-btn" @click="showCouponModal = true">
                <TicketIcon class="icon" />
                选择优惠券
                <ChevronRightIcon class="arrow" />
              </button>
            </div>
          </section>

          <!-- Note -->
          <section class="checkout-section">
            <h3>订单备注</h3>
            <textarea 
              v-model="orderNote"
              class="order-note"
              rows="3"
              placeholder="请输入订单备注，如配送时间要求等"
            ></textarea>
          </section>
        </div>

        <!-- Right Column - Summary -->
        <div class="checkout-sidebar">
          <div class="summary-card">
            <h4 class="summary-title">订单 summary</h4>
            
            <div class="summary-rows">
              <div class="summary-row">
                <span>商品总额</span>
                <span>¥{{ subtotal.toFixed(2) }}</span>
              </div>
              <div class="summary-row">
                <span>运费</span>
                <span>{{ shippingPrice === 0 ? '免运费' : '¥' + shippingPrice.toFixed(2) }}</span>
              </div>
              <div v-if="couponDiscount > 0" class="summary-row discount">
                <span>优惠券</span>
                <span>-¥{{ couponDiscount.toFixed(2) }}</span>
              </div>
            </div>

            <div class="summary-total">
              <div class="total-row">
                <span>应付总额</span>
                <span class="total-price">¥{{ total.toFixed(2) }}</span>
              </div>
              <p class="total-saving" v-if="totalSaving > 0">
                已优惠 ¥{{ totalSaving.toFixed(2) }}
              </p>
            </div>

            <button class="btn-submit" @click="submitOrder" :disabled="!canSubmit">
              提交订单
            </button>

            <div class="payment-methods">
              <p class="payment-title">支付方式</p>
              <div class="payment-icons">
                <div class="payment-icon wechat">
                  <svg viewBox="0 0 24 24" fill="currentColor"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 01.213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 00.167-.054l1.903-1.114a.864.864 0 01.717-.098 10.16 10.16 0 002.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178A1.17 1.17 0 014.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178 1.17 1.17 0 01-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 01.598.082l1.584.926a.272.272 0 00.14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 01-.023-.156.49.49 0 01.201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-7.062-6.122zM14.51 13.88a.94.94 0 01.939-.944c.519 0 .94.423.94.944a.94.94 0 01-.94.943.94.94 0 01-.94-.943zm4.963 0a.94.94 0 01.94-.944c.519 0 .939.423.939.944a.94.94 0 01-.94.943.94.94 0 01-.939-.943z"/></svg>
                </div>
                <div class="payment-icon alipay">
                  <svg viewBox="0 0 24 24" fill="currentColor"><path d="M5.5 2C3.567 2 2 3.567 2 5.5v13C2 20.433 3.567 22 5.5 22h13c1.933 0 3.5-1.567 3.5-3.5v-13C22 3.567 20.433 2 18.5 2h-13zm8.182 4h1.529c.103 0 .186.083.186.186v8.628a.186.186 0 01-.186.186h-1.529a.186.186 0 01-.186-.186V6.186c0-.103.083-.186.186-.186zm-4.364 0h1.529c.103 0 .186.083.186.186v8.628a.186.186 0 01-.186.186H9.318a.186.186 0 01-.186-.186V6.186c0-.103.083-.186.186-.186zm8.728 0h1.529c.103 0 .186.083.186.186v8.628a.186.186 0 01-.186.186h-1.529a.186.186 0 01-.186-.186V6.186c0-.103.083-.186.186-.186z"/></svg>
                </div>
                <div class="payment-icon card">
                  <CreditCardIcon class="icon" />
                </div>
              </div>
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
  PlusIcon,
  PencilIcon,
  TicketIcon,
  ChevronRightIcon,
  XMarkIcon,
  CreditCardIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()

// Addresses
const addresses = ref([
  {
    id: 1,
    name: '张三',
    phone: '138****8001',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    detail: '科技园南区A栋1201',
    isDefault: true
  },
  {
    id: 2,
    name: '张三',
    phone: '138****8001',
    province: '广东省',
    city: '广州市',
    district: '天河区',
    detail: '珠江新城华夏路100号',
    isDefault: false
  }
])

const selectedAddress = ref(addresses.value[0])
const showAddressModal = ref(false)

// Order Items
const orderItems = ref([
  { id: 1, name: '无线蓝牙耳机 Pro', variant: '深邃黑', price: 299, quantity: 1, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=100&h=100&fit=crop' },
  { id: 2, name: '智能手表 Series 7', variant: '银色/44mm', price: 1999, quantity: 1, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=100&h=100&fit=crop' }
])

// Shipping
const shippingMethods = ref([
  { id: 'express', name: '快递配送', price: 0, description: '预计3-5个工作日送达，满199免运费' },
  { id: 'same-day', name: '当日达', price: 25, description: '当日下单，当日送达（仅限部分城市）' },
  { id: 'pickup', name: '门店自提', price: 0, description: '下单后次日可到指定门店自提' }
])
const selectedShipping = ref('express')

// Coupon
const appliedCoupon = ref<any>(null)
const showCouponModal = ref(false)

// Note
const orderNote = ref('')

// Computed
const subtotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const shippingPrice = computed(() => {
  const method = shippingMethods.value.find(m => m.id === selectedShipping.value)
  return method?.price || 0
})

const couponDiscount = computed(() => {
  return appliedCoupon.value?.discount || 0
})

const total = computed(() => subtotal.value + shippingPrice.value - couponDiscount.value)

const totalSaving = computed(() => {
  const originalTotal = orderItems.value.reduce((sum, item) => {
    const originalPrice = item.price * 1.2 // Assume 20% discount
    return sum + originalPrice * item.quantity
  }, 0)
  return originalTotal - total.value
})

const canSubmit = computed(() => {
  return selectedAddress.value && orderItems.value.length > 0
})

// Methods
const editAddress = (addr: any) => {
  console.log('Edit address:', addr)
}

const submitOrder = () => {
  router.push('/payment')
}
</script>

<style scoped>
.checkout-page {
  min-height: calc(100vh - 72px);
  background: #F9FAFB;
  padding: 40px 0;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Checkout Steps */
.checkout-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 40px;
}

.step {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #E5E7EB;
  color: #6B7280;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.step.active .step-number {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
}

.step-label {
  font-size: 15px;
  font-weight: 500;
  color: #6B7280;
}

.step.active .step-label {
  color: #059669;
  font-weight: 600;
}

.step-line {
  width: 80px;
  height: 2px;
  background: #E5E7EB;
}

/* Checkout Content */
.checkout-content {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
}

/* Checkout Main */
.checkout-main {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.checkout-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.checkout-section h3 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 20px 0;
}

.btn-text {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #059669;
  background: none;
  border: none;
  cursor: pointer;
  font-weight: 500;
}

.btn-text .icon {
  width: 18px;
  height: 18px;
}

/* Address List */
.address-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.address-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 2px solid #E5E7EB;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.address-card:hover {
  border-color: #10B981;
}

.address-card.active {
  border-color: #059669;
  background: #ECFDF5;
}

.address-radio {
  flex-shrink: 0;
}

.radio-indicator {
  width: 20px;
  height: 20px;
  border: 2px solid #D1D5DB;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.radio-indicator.checked {
  border-color: #059669;
  background: #059669;
}

.radio-indicator.checked::after {
  content: '';
  width: 8px;
  height: 8px;
  background: white;
  border-radius: 50%;
}

.address-info {
  flex: 1;
  min-width: 0;
}

.address-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.phone {
  font-size: 14px;
  color: #6B7280;
}

.default-tag {
  padding: 4px 8px;
  background: #059669;
  color: white;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.address-detail {
  font-size: 14px;
  color: #4B5563;
  margin: 0;
}

.edit-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.edit-btn:hover {
  background: #F3F4F6;
}

.edit-btn .icon {
  width: 18px;
  height: 18px;
  color: #6B7280;
}

/* Order Items */
.order-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 12px;
}

.item-image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 4px 0;
}

.item-variant {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.item-price {
  text-align: center;
}

.item-price .price {
  display: block;
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.item-price .quantity {
  font-size: 13px;
  color: #6B7280;
}

.item-subtotal {
  font-size: 16px;
  font-weight: 700;
  color: #EF4444;
  min-width: 100px;
  text-align: right;
}

/* Shipping Options */
.shipping-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.shipping-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 2px solid #E5E7EB;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.shipping-card:hover {
  border-color: #10B981;
}

.shipping-card.active {
  border-color: #059669;
  background: #ECFDF5;
}

.shipping-card input {
  display: none;
}

.shipping-content {
  flex: 1;
}

.shipping-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.method-name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.method-price {
  font-size: 15px;
  font-weight: 700;
  color: #059669;
}

.method-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Coupon Section */
.applied-coupon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  border-radius: 10px;
}

.coupon-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.coupon-name {
  font-size: 14px;
  font-weight: 500;
  color: #92400E;
}

.coupon-discount {
  font-size: 16px;
  font-weight: 700;
  color: #B45309;
}

.remove-coupon {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-coupon .icon {
  width: 18px;
  height: 18px;
  color: #92400E;
}

.select-coupon-btn {
  width: 100%;
  padding: 16px;
  border: 2px dashed #D1D5DB;
  background: white;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 15px;
  color: #6B7280;
  cursor: pointer;
  transition: all 0.2s;
}

.select-coupon-btn:hover {
  border-color: #059669;
  color: #059669;
}

.select-coupon-btn .icon {
  width: 20px;
  height: 20px;
}

.select-coupon-btn .arrow {
  width: 18px;
  height: 18px;
  margin-left: auto;
}

/* Order Note */
.order-note {
  width: 100%;
  padding: 16px;
  border: 1px solid #E5E7EB;
  border-radius: 10px;
  font-size: 14px;
  resize: vertical;
  outline: none;
  transition: border-color 0.2s;
}

.order-note:focus {
  border-color: #059669;
}

/* Checkout Sidebar */
.checkout-sidebar {
  position: sticky;
  top: 96px;
  height: fit-content;
}

.summary-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.summary-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 20px 0;
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

.summary-row.discount {
  color: #059669;
}

.summary-total {
  margin-bottom: 24px;
}

.total-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.total-row span:first-child {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.total-price {
  font-size: 32px;
  font-weight: 800;
  color: #EF4444;
}

.total-saving {
  text-align: right;
  font-size: 13px;
  color: #059669;
  margin: 8px 0 0 0;
}

.btn-submit {
  width: 100%;
  padding: 18px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 17px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 20px;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px -8px rgba(5, 150, 105, 0.5);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.payment-methods {
  text-align: center;
}

.payment-title {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 12px 0;
}

.payment-icons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.payment-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.payment-icon.wechat {
  background: #22C55E;
  color: white;
}

.payment-icon.alipay {
  background: #3B82F6;
  color: white;
}

.payment-icon.card {
  background: #F3F4F6;
}

.payment-icon svg,
.payment-icon .icon {
  width: 24px;
  height: 24px;
}

/* Responsive */
@media (max-width: 1024px) {
  .checkout-content {
    grid-template-columns: 1fr;
  }

  .checkout-sidebar {
    position: static;
  }
}

@media (max-width: 640px) {
  .checkout-steps {
    gap: 8px;
  }

  .step-label {
    display: none;
  }

  .step-line {
    width: 40px;
  }

  .order-item {
    flex-wrap: wrap;
  }

  .item-subtotal {
    width: 100%;
    text-align: left;
    margin-top: 8px;
  }
}
</style>
