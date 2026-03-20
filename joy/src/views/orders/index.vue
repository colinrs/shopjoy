<template>
  <div class="orders-page">
    <div class="page-container">
      <!-- Page Header -->
      <div class="page-header">
        <h1 class="page-title">我的订单</h1>
        <div class="order-stats">
          <div class="stat-item">
            <span class="stat-value">{{ orderStats.pending }}</span>
            <span class="stat-label">待付款</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ orderStats.shipping }}</span>
            <span class="stat-label">待发货</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ orderStats.receiving }}</span>
            <span class="stat-label">待收货</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ orderStats.review }}</span>
            <span class="stat-label">待评价</span>
          </div>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="filter-tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.value"
          class="tab-btn"
          :class="{ active: activeTab === tab.value }"
          @click="activeTab = tab.value"
        >
          {{ tab.label }}
          <span v-if="tab.count" class="tab-badge">{{ tab.count }}</span>
        </button>
      </div>

      <!-- Orders List -->
      <div class="orders-list">
        <div v-for="order in filteredOrders" :key="order.id" class="order-card">
          <div class="order-header">
            <div class="order-info">
              <span class="order-no">订单号: {{ order.orderNo }}</span>
              <span class="order-date">{{ order.date }}</span>
            </div>
            <div class="order-status" :class="order.status">
              {{ getStatusText(order.status) }}
            </div>
          </div>

          <div class="order-items">
            <div v-for="item in order.items" :key="item.id" class="order-item">
              <img :src="item.image" :alt="item.name" class="item-image" />
              <div class="item-details">
                <h4 class="item-name">{{ item.name }}</h4>
                <p class="item-variant">{{ item.variant }}</p>
                <p class="item-quantity">x{{ item.quantity }}</p>
              </div>
              <span class="item-price">¥{{ item.price }}</span>
            </div>
          </div>

          <div class="order-footer">
            <div class="order-total">
              <span>共 {{ order.items.length }} 件商品</span>
              <span class="total-amount">
                实付款: <strong>¥{{ order.total }}</strong>
              </span>
            </div>
            <div class="order-actions">
              <button 
                v-if="order.status === 'pending'"
                class="btn-primary"
                @click="payOrder(order)"
              >
                立即支付
              </button>
              <button 
                v-if="order.status === 'shipped'"
                class="btn-primary"
                @click="confirmReceive(order)"
              >
                确认收货
              </button>
              <button 
                v-if="order.status === 'completed' && !order.reviewed"
                class="btn-outline"
                @click="reviewOrder(order)"
              >
                评价
              </button>
              <button 
                v-if="['pending', 'paid'].includes(order.status)"
                class="btn-text"
                @click="cancelOrder(order)"
              >
                取消订单
              </button>
              <button 
                class="btn-text"
                @click="viewDetail(order)"
              >
                查看详情
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="filteredOrders.length === 0" class="empty-state">
        <div class="empty-icon">
          <ClipboardDocumentListIcon class="icon" />
        </div>
        <h3>暂无订单</h3>
        <p>快去选购心仪的商品吧</p>
        <router-link to="/products" class="btn-primary">去购物</router-link>
      </div>

      <!-- Pagination -->
      <div v-if="filteredOrders.length > 0" class="pagination">
        <button class="page-btn" :disabled="currentPage === 1" @click="currentPage--">
          <ChevronLeftIcon class="icon" />
        </button>
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="currentPage === totalPages" @click="currentPage++">
          <ChevronRightIcon class="icon" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  ClipboardDocumentListIcon,
  ChevronLeftIcon,
  ChevronRightIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()

const orderStats = ref({
  pending: 2,
  shipping: 1,
  receiving: 1,
  review: 3
})

const tabs = [
  { label: '全部', value: 'all' },
  { label: '待付款', value: 'pending', count: 2 },
  { label: '待发货', value: 'paid', count: 1 },
  { label: '待收货', value: 'shipped', count: 1 },
  { label: '待评价', value: 'completed', count: 3 }
]

const activeTab = ref('all')
const currentPage = ref(1)
const totalPages = ref(5)

const orders = ref([
  {
    id: 1,
    orderNo: 'ORD20240318001',
    date: '2024-03-18 14:30:25',
    status: 'pending',
    total: 2298,
    items: [
      { id: 1, name: '无线蓝牙耳机 Pro', variant: '深邃黑', price: 299, quantity: 1, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=100&h=100&fit=crop' },
      { id: 2, name: '智能手表 Series 7', variant: '银色/44mm', price: 1999, quantity: 1, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=100&h=100&fit=crop' }
    ]
  },
  {
    id: 2,
    orderNo: 'ORD20240315002',
    date: '2024-03-15 10:20:18',
    status: 'shipped',
    total: 129,
    items: [
      { id: 3, name: '便携充电宝 20000mAh', variant: '白色', price: 129, quantity: 1, image: 'https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=100&h=100&fit=crop' }
    ]
  },
  {
    id: 3,
    orderNo: 'ORD20240310003',
    date: '2024-03-10 16:45:30',
    status: 'completed',
    total: 459,
    reviewed: false,
    items: [
      { id: 4, name: '机械键盘 RGB', variant: '青轴', price: 459, quantity: 1, image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=100&h=100&fit=crop' }
    ]
  }
])

const filteredOrders = computed(() => {
  if (activeTab.value === 'all') return orders.value
  return orders.value.filter(order => order.status === activeTab.value)
})

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待付款',
    paid: '待发货',
    shipped: '待收货',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

const payOrder = (order: any) => {
  router.push('/checkout')
}

const confirmReceive = (order: any) => {
  if (confirm('确认已收到商品？')) {
    order.status = 'completed'
    alert('确认收货成功')
  }
}

const reviewOrder = (order: any) => {
  alert('前往评价页面')
}

const cancelOrder = (order: any) => {
  if (confirm('确定要取消该订单吗？')) {
    order.status = 'cancelled'
    alert('订单已取消')
  }
}

const viewDetail = (order: any) => {
  alert('查看订单详情: ' + order.orderNo)
}
</script>

<style scoped>
.orders-page {
  min-height: calc(100vh - 72px);
  background: #F9FAFB;
  padding: 40px 0;
}

.page-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Page Header */
.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 24px 0;
}

.order-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  flex: 1;
  background: white;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: #059669;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
}

/* Filter Tabs */
.filter-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.tab-btn {
  padding: 12px 20px;
  background: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #6B7280;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-btn:hover {
  background: #ECFDF5;
  color: #059669;
}

.tab-btn.active {
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
}

.tab-badge {
  padding: 2px 8px;
  background: #EF4444;
  color: white;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 600;
}

/* Orders List */
.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #F9FAFB;
  border-bottom: 1px solid #E5E7EB;
}

.order-info {
  display: flex;
  gap: 16px;
}

.order-no {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.order-date {
  font-size: 13px;
  color: #9CA3AF;
}

.order-status {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.order-status.pending {
  background: #FEF3C7;
  color: #B45309;
}

.order-status.paid {
  background: #DBEAFE;
  color: #1D4ED8;
}

.order-status.shipped {
  background: #ECFDF5;
  color: #059669;
}

.order-status.completed {
  background: #F3F4F6;
  color: #6B7280;
}

.order-status.cancelled {
  background: #FEF2F2;
  color: #DC2626;
}

/* Order Items */
.order-items {
  padding: 20px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #F3F4F6;
}

.order-item:last-child {
  border-bottom: none;
}

.item-image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
}

.item-details {
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
  margin: 0 0 4px 0;
}

.item-quantity {
  font-size: 13px;
  color: #9CA3AF;
}

.item-price {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

/* Order Footer */
.order-footer {
  padding: 16px 20px;
  border-top: 1px solid #E5E7EB;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-total {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 14px;
  color: #6B7280;
}

.total-amount {
  font-size: 16px;
}

.total-amount strong {
  font-size: 20px;
  color: #EF4444;
}

.order-actions {
  display: flex;
  gap: 12px;
}

.btn-primary {
  padding: 10px 20px;
  background: linear-gradient(135deg, #059669 0%, #10B981 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px -8px rgba(5, 150, 105, 0.5);
}

.btn-outline {
  padding: 10px 20px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-outline:hover {
  border-color: #059669;
  color: #059669;
}

.btn-text {
  padding: 10px 16px;
  background: transparent;
  border: none;
  font-size: 14px;
  color: #6B7280;
  cursor: pointer;
  transition: color 0.2s;
}

.btn-text:hover {
  color: #374151;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 80px 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.empty-icon {
  width: 100px;
  height: 100px;
  background: #F3F4F6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24px;
}

.empty-icon .icon {
  width: 48px;
  height: 48px;
  color: #9CA3AF;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 15px;
  color: #6B7280;
  margin: 0 0 24px 0;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 40px;
}

.page-btn {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border: 1px solid #E5E7EB;
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

.page-info {
  font-size: 14px;
  color: #6B7280;
}

/* Responsive */
@media (max-width: 640px) {
  .order-stats {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }

  .order-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .order-footer {
    flex-direction: column;
    gap: 16px;
  }

  .order-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
