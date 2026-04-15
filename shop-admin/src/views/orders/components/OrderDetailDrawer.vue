<template>
  <el-drawer
    v-model="visible"
    :title="$t('orders.orderDetails')"
    size="720px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div
      v-loading="loading"
      class="order-detail-drawer"
    >
      <!-- Order Status Timeline -->
      <div class="status-section">
        <div class="section-header">
          <el-icon><Clock /></el-icon>
          <span>{{ $t('orders.orderStatusTimeline') }}</span>
        </div>
        <el-timeline class="status-timeline">
          <el-timeline-item
            v-for="(step, index) in statusTimeline"
            :key="index"
            :type="step.type"
            :timestamp="step.time"
            placement="top"
          >
            <p class="timeline-title">
              {{ step.title }}
            </p>
            <p
              v-if="step.description"
              class="timeline-desc"
            >
              {{ step.description }}
            </p>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- Order Info -->
      <div class="info-section">
        <div class="section-header">
          <el-icon><Document /></el-icon>
          <span>{{ $t('orders.orderInformation') }}</span>
        </div>
        <el-descriptions
          :column="2"
          border
          size="small"
        >
          <el-descriptions-item :label="$t('orders.orderNo')">
            <span class="order-no-text">{{ order?.order_no }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.status')">
            <el-tag
              :type="getStatusTagType(order?.status)"
              size="small"
            >
              {{ order?.status_text }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.fulfillmentStatus')">
            <el-tag
              :type="getFulfillmentTagType(order?.fulfillment_status)"
              size="small"
            >
              {{ order?.fulfillment_text }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('common.createdAt')">
            {{ formatTime(order?.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.source')">
            {{ order?.order_source || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.buyerRemark')">
            {{ order?.buyer_remark || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Items List -->
      <div class="items-section">
        <div class="section-header">
          <el-icon><Goods /></el-icon>
          <span>{{ $t('orders.orderItems') }}</span>
          <el-tag
            size="small"
            type="info"
          >
            {{ $t('orders.itemsCount', { count: order?.items?.length || 0 }) }}
          </el-tag>
        </div>
        <div class="items-list">
          <div
            v-for="item in order?.items"
            :key="item.order_item_id"
            class="item-row"
          >
            <el-image
              :src="item.image"
              class="item-image"
              fit="cover"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="item-info">
              <p class="item-name">
                {{ item.product_name }}
              </p>
              <p class="item-sku">
                {{ item.sku_name }}
              </p>
            </div>
            <div class="item-price">
              <p class="unit-price">
                {{ order?.currency }} {{ formatAmount(item.unit_price) }}
              </p>
              <p class="quantity">
                x{{ item.quantity }}
              </p>
            </div>
            <div class="item-total">
              {{ order?.currency }} {{ formatAmount(String(parseFloat(item.unit_price) * item.quantity)) }}
            </div>
            <div
              v-if="item.shipped_qty > 0"
              class="item-ship-info"
            >
              <el-tag
                size="small"
                type="success"
              >
                {{ item.shipped_qty }} {{ $t('orders.shippedQty') }}
              </el-tag>
              <el-tag
                v-if="item.pending_qty > 0"
                size="small"
                type="warning"
              >
                {{ item.pending_qty }} {{ $t('orders.pendingQty') }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>

      <!-- Shipping Address -->
      <div class="address-section">
        <div class="section-header">
          <el-icon><Location /></el-icon>
          <span>{{ $t('orders.shippingAddress') }}</span>
        </div>
        <div class="address-card">
          <p class="receiver">
            <span class="name">{{ order?.shipping_address?.receiver_name }}</span>
            <span class="phone">{{ order?.shipping_address?.receiver_phone }}</span>
          </p>
          <p class="address">
            {{ order?.shipping_address?.full_address }}
          </p>
        </div>
      </div>

      <!-- Payment Info -->
      <div class="payment-section">
        <div class="section-header">
          <el-icon><Wallet /></el-icon>
          <span>{{ $t('orders.paymentInformation') }}</span>
        </div>
        <el-descriptions
          :column="2"
          border
          size="small"
        >
          <el-descriptions-item :label="$t('orders.paymentMethod')">
            {{ order?.payment_method_text || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.productAmount')">
            {{ order?.currency }} {{ formatAmount(order?.total_amount) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.amount')">
            <span class="pay-amount">{{ order?.currency }} {{ formatAmount(order?.pay_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.paidAt')">
            {{ formatTime(order?.paid_at) || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Shipments -->
      <div
        v-if="order?.shipments?.length"
        class="shipments-section"
      >
        <div class="section-header">
          <el-icon><Van /></el-icon>
          <span>{{ $t('orders.shipments') }}</span>
          <el-tag
            size="small"
            type="info"
          >
            {{ $t('orders.shipmentCount', { count: order.shipments.length }) }}
          </el-tag>
        </div>
        <div class="shipments-list">
          <div
            v-for="shipment in order.shipments"
            :key="shipment.shipment_id"
            class="shipment-card"
          >
            <div class="shipment-header">
              <span class="shipment-no">{{ shipment.shipment_no }}</span>
              <el-tag
                :type="getShipmentStatusTagType(shipment.status)"
                size="small"
              >
                {{ getShipmentStatusText(shipment.status) }}
              </el-tag>
            </div>
            <div class="shipment-info">
              <p><strong>{{ $t('orders.carrier') }}:</strong> {{ shipment.carrier_name }}</p>
              <p><strong>{{ $t('orders.trackingNo') }}:</strong> {{ shipment.tracking_no }}</p>
              <p><strong>{{ $t('orders.shippedAtStatus') }}:</strong> {{ formatTime(shipment.shipped_at) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Seller Remark -->
      <div class="remark-section">
        <div class="section-header">
          <el-icon><ChatDotSquare /></el-icon>
          <span>{{ $t('orders.sellerRemark') }}</span>
        </div>
        <div class="remark-content">
          <p v-if="order?.seller_remark">
            {{ order.seller_remark }}
          </p>
          <p
            v-else
            class="no-remark"
          >
            {{ $t('orders.noRemark') }}
          </p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <el-button @click="handleRemark">
          <el-icon><Edit /></el-icon>
          {{ $t('orders.editRemark') }}
        </el-button>
        <el-button
          v-if="canShip"
          type="primary"
          @click="handleShip"
        >
          <el-icon><Van /></el-icon>
          {{ $t('orders.shipOrder') }}
        </el-button>
        <el-button
          v-if="canAdjustPrice"
          @click="handleAdjustPrice"
        >
          <el-icon><PriceTag /></el-icon>
          {{ $t('orders.adjustPrice') }}
        </el-button>
      </div>
    </template>

    <!-- Ship Dialog -->
    <ShipDialog
      v-model="shipDialogVisible"
      :order="orderForShip"
      :carriers="carriers"
      @success="handleShipSuccess"
    />

    <!-- Remark Dialog -->
    <RemarkDialog
      v-model="remarkDialogVisible"
      :order-id="orderId"
      :current-remark="order?.seller_remark || ''"
      @success="handleRemarkSuccess"
    />

    <!-- Adjust Price Dialog -->
    <AdjustPriceDialog
      v-model="adjustPriceDialogVisible"
      :order-id="orderId"
      :current-amount="order?.pay_amount || '0'"
      :currency="order?.currency || 'CNY'"
      @success="handleAdjustPriceSuccess"
    />
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Clock,
  Document,
  Goods,
  Location,
  Wallet,
  Van,
  Picture,
  ChatDotSquare,
  Edit,
  PriceTag
} from '@element-plus/icons-vue'
import {
  getOrderDetail,
  getCarrierList,
  type OrderDetail,
  type Carrier,
  type OrderStatus,
  type FulfillmentStatus
} from '@/api/order'
import ShipDialog from './ShipDialog.vue'
import RemarkDialog from './RemarkDialog.vue'
import AdjustPriceDialog from './AdjustPriceDialog.vue'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

const props = defineProps<{
  modelValue: boolean
  orderId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  refresh: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const order = ref<OrderDetail | null>(null)
const carriers = ref<Carrier[]>([])

const shipDialogVisible = ref(false)
const remarkDialogVisible = ref(false)
const adjustPriceDialogVisible = ref(false)

// Order for ship dialog - transform OrderDetail to Order-like structure
const orderForShip = computed(() => {
  if (!order.value) return null
  return {
    order_id: order.value.order_id,
    order_no: order.value.order_no,
    items: order.value.items.map(item => ({
      order_item_id: item.order_item_id,
      product_name: item.product_name,
      sku_name: item.sku_name,
      image: item.image,
      quantity: item.quantity,
      pending_qty: item.pending_qty
    }))
  }
})

// Status timeline
const statusTimeline = computed(() => {
  if (!order.value) return []

  const timeline = []
  const o = order.value

  // Created
  timeline.push({
    title: t('orders.orderCreated'),
    time: formatTime(o.created_at),
    type: 'primary' as const,
    description: `${t('orders.orderNo')}: ${o.order_no}`
  })

  // Paid
  if (o.paid_at) {
    timeline.push({
      title: t('orders.paymentReceived'),
      time: formatTime(o.paid_at),
      type: 'success' as const,
      description: `${o.payment_method_text} - ${o.currency} ${formatAmount(o.pay_amount)}`
    })
  }

  // Shipped
  if (o.shipments?.length) {
    o.shipments.forEach(s => {
      timeline.push({
        title: t('orders.orderShippedStatus'),
        time: formatTime(s.shipped_at),
        type: 'primary' as const,
        description: `${s.carrier_name} - ${s.tracking_no}`
      })
    })
  }

  // Completed
  if (o.status === 'delivered') {
    timeline.push({
      title: t('orders.orderCompleted'),
      time: formatTime(o.shipments?.[0]?.delivered_at),
      type: 'success' as const,
      description: ''
    })
  }

  // Cancelled
  if (o.status === 'cancelled') {
    timeline.push({
      title: t('orders.orderCancelledStatus'),
      time: '-',
      type: 'danger' as const,
      description: ''
    })
  }

  return timeline
})

// Computed flags for actions
const canShip = computed(() => {
  if (!order.value) return false
  return order.value.status === 'paid' &&
    (order.value.fulfillment_status === 0 || order.value.fulfillment_status === 1)
})

const canAdjustPrice = computed(() => {
  if (!order.value) return false
  return order.value.status === 'pending_payment'
})

// Load order detail
const loadOrderDetail = async () => {
  if (!props.orderId) return

  loading.value = true
  try {
    const res = await getOrderDetail(props.orderId)
    order.value = res
  } catch (error) {
    ElMessage.error(t('orders.loadDetailFailed'))
  } finally {
    loading.value = false
  }
}

// Load carriers
const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res.list.filter(c => c.is_active)
  } catch (error) {
    handleError(error, t('orders.loadCarriersFailed'))
  }
}

// Format helpers
const formatTime = (time: string | undefined | null) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

const formatAmount = (amount: string | undefined) => {
  if (!amount) return '0.00'
  return parseFloat(amount).toFixed(2)
}

// Status helpers
const getStatusTagType = (status: OrderStatus | undefined) => {
  const types: Record<OrderStatus, string> = {
    pending_payment: 'warning',
    paid: 'success',
    shipped: 'info',
    delivered: 'success',
    cancelled: 'danger',
    refunded: 'info'
  }
  return status ? types[status] : 'info'
}

const getFulfillmentTagType = (status: FulfillmentStatus | undefined) => {
  const types: Record<FulfillmentStatus, string> = {
    0: 'warning',
    1: 'primary',
    2: 'info',
    3: 'success'
  }
  return status !== undefined ? types[status] : 'info'
}

const getShipmentStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'shipped': 'primary',
    'in_transit': 'info',
    'delivered': 'success',
    'failed': 'danger'
  }
  return types[status] || 'info'
}

const getShipmentStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'pending': t('orders.shipmentPending'),
    'shipped': t('orders.shipmentShipped'),
    'in_transit': t('orders.shipmentInTransit'),
    'delivered': t('orders.shipmentDelivered'),
    'failed': t('orders.shipmentFailed')
  }
  return texts[status] || t('orders.shipmentUnknown')
}

// Action handlers
const handleShip = () => {
  shipDialogVisible.value = true
}

const handleRemark = () => {
  remarkDialogVisible.value = true
}

const handleAdjustPrice = () => {
  adjustPriceDialogVisible.value = true
}

const handleShipSuccess = () => {
  ElMessage.success(t('orders.orderShipped'))
  loadOrderDetail()
  emit('refresh')
}

const handleRemarkSuccess = () => {
  ElMessage.success(t('orders.remarkSaveSuccess'))
  loadOrderDetail()
  emit('refresh')
}

const handleAdjustPriceSuccess = () => {
  ElMessage.success(t('orders.priceAdjustSuccess'))
  loadOrderDetail()
  emit('refresh')
}

// Watch for visibility changes
watch(visible, (val) => {
  if (val) {
    loadOrderDetail()
    loadCarriers()
  }
})
</script>

<style scoped>
.order-detail-drawer {
  padding: 0 4px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-weight: 600;
  color: #1E1B4B;
  font-size: 15px;
}

.section-header .el-icon {
  color: #6366F1;
}

/* Status Section */
.status-section {
  margin-bottom: 24px;
  padding: 20px;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 16px;
}

.status-timeline {
  padding-left: 0;
}

.timeline-title {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.timeline-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

/* Info Section */
.info-section {
  margin-bottom: 24px;
}

.order-no-text {
  font-family: 'Fira Code', monospace;
  color: #6366F1;
  font-weight: 500;
}

/* Items Section */
.items-section {
  margin-bottom: 24px;
}

.items-list {
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  overflow: hidden;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #F3F4F6;
}

.item-row:last-child {
  border-bottom: none;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  flex-shrink: 0;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.item-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
}

.item-price {
  text-align: right;
}

.unit-price {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.quantity {
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
  margin: 4px 0 0 0;
}

.item-total {
  font-weight: 600;
  color: #1E1B4B;
  min-width: 80px;
  text-align: right;
}

.item-ship-info {
  display: flex;
  gap: 4px;
}

/* Address Section */
.address-section {
  margin-bottom: 24px;
}

.address-card {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.receiver {
  display: flex;
  gap: 16px;
  margin: 0 0 8px 0;
}

.receiver .name {
  font-weight: 600;
  color: #1E1B4B;
}

.receiver .phone {
  color: #6B7280;
}

.address {
  color: #4B5563;
  margin: 0;
}

/* Payment Section */
.payment-section {
  margin-bottom: 24px;
}

.discount-text {
  color: #10B981;
}

.pay-amount {
  font-weight: 700;
  color: #EF4444;
  font-size: 16px;
}

/* Shipments Section */
.shipments-section {
  margin-bottom: 24px;
}

.shipments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.shipment-card {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.shipment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.shipment-no {
  font-family: 'Fira Code', monospace;
  color: #6366F1;
  font-weight: 500;
}

.shipment-info p {
  font-size: 13px;
  color: #4B5563;
  margin: 0 0 4px 0;
}

/* Remark Section */
.remark-section {
  margin-bottom: 24px;
}

.remark-content {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.remark-content p {
  margin: 0;
  color: #4B5563;
}

.no-remark {
  color: #9CA3AF;
  font-style: italic;
}

/* Drawer Footer */
.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Descriptions */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}
</style>