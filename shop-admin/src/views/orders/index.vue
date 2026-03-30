<template>
  <div class="orders-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.pending_payment }}</p>
          <p class="stat-label">{{ $t('orders.pendingPayment') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.partial_shipped }}</p>
          <p class="stat-label">{{ $t('orders.partialShipped') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.shipped }}</p>
          <p class="stat-label">{{ $t('orders.shipped') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.delivered }}</p>
          <p class="stat-label">{{ $t('orders.delivered') }}</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('orders.searchPlaceholder')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" :placeholder="$t('orders.filterOrderStatus')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('common.all')" value="" />
            <el-option :label="$t('orders.pendingPayment')" value="pending_payment" />
            <el-option :label="$t('orders.paid')" value="paid" />
            <el-option :label="$t('orders.pendingShipment')" value="pending_shipment" />
            <el-option :label="$t('orders.shipped')" value="shipped" />
            <el-option :label="$t('orders.completed')" value="completed" />
            <el-option :label="$t('orders.cancelled')" value="cancelled" />
            <el-option :label="$t('orders.refunding')" value="refunding" />
            <el-option :label="$t('orders.refunded')" value="refunded" />
          </el-select>
          <el-select v-model="fulfillmentFilter" :placeholder="$t('orders.filterFulfillment')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('common.all')" value="" />
            <el-option :label="$t('orders.unshipped')" value="pending" />
            <el-option :label="$t('orders.partialShipped')" value="partial_shipped" />
            <el-option :label="$t('orders.shipped')" value="shipped" />
            <el-option :label="$t('orders.delivered')" value="delivered" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('orders.startDate')"
            :end-placeholder="$t('orders.endDate')"
            class="date-picker"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport" :loading="exporting">
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            {{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Batch Actions Bar -->
    <transition name="slide-down">
      <div v-if="selectedOrders.length > 0" class="batch-actions-bar">
        <div class="batch-info">
          <el-icon><Check /></el-icon>
          <span>{{ $t('orders.ordersSelected', { count: selectedOrders.length }) }}</span>
        </div>
        <div class="batch-buttons">
          <el-button type="primary" @click="handleBatchShip">
            <el-icon><Van /></el-icon>
            {{ $t('orders.batchShip') }}
          </el-button>
          <el-button @click="handleBatchRemark">
            <el-icon><Edit /></el-icon>
            {{ $t('orders.batchRemark') }}
          </el-button>
          <el-button @click="clearSelection">{{ $t('orders.clearSelection') }}</el-button>
        </div>
      </div>
    </transition>

    <!-- Orders Table -->
    <el-card class="table-card" shadow="never">
      <el-table
        ref="tableRef"
        :data="orderList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="order-detail">
              <el-row :gutter="20">
                <el-col :span="12">
                  <h4>{{ $t('orders.items') }}</h4>
                  <div v-for="item in row.items" :key="item.order_item_id" class="order-item">
                    <el-image :src="item.image" class="item-image" fit="cover">
                      <template #error>
                        <div class="image-placeholder">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <div class="item-info">
                      <p class="item-name">{{ item.product_name }}</p>
                      <p class="item-sku">{{ item.sku_name }}</p>
                      <p class="item-price">
                        {{ row.currency }} {{ formatAmount(item.unit_price) }} x {{ item.quantity }}
                      </p>
                    </div>
                  </div>
                </el-col>
                <el-col :span="12">
                  <h4>{{ $t('orders.shippingInfo') }}</h4>
                  <p><strong>{{ $t('orders.receiver') }}:</strong> {{ row.shipping_address?.receiver_name || '-' }}</p>
                  <p><strong>{{ $t('orders.phone') }}:</strong> {{ row.shipping_address?.receiver_phone || '-' }}</p>
                  <p><strong>{{ $t('orders.address') }}:</strong> {{ row.shipping_address?.full_address || '-' }}</p>
                  <h4 style="margin-top: 20px">{{ $t('orders.paymentInfo') }}</h4>
                  <p><strong>{{ $t('orders.paymentMethod') }}:</strong> {{ row.payment_method_name || '-' }}</p>
                  <p><strong>{{ $t('orders.paidAt') }}:</strong> {{ formatTime(row.paid_at) }}</p>
                </el-col>
              </el-row>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" :label="$t('orders.orderNo')" min-width="160">
          <template #default="{ row }">
            <div class="order-no-cell">
              <span class="order-no" @click="handleDetail(row)">{{ row.order_no }}</span>
              <el-tag v-if="row.refund_status !== 'none' && row.refund_status !== undefined" size="small" type="warning" effect="dark">
                {{ row.refund_text }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('orders.items')" min-width="200">
          <template #default="{ row }">
            <div class="goods-preview">
              <el-image
                v-for="(item, idx) in row.items.slice(0, 3)"
                :key="idx"
                :src="item.image"
                class="goods-thumb"
                fit="cover"
              >
                <template #error>
                  <div class="thumb-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <span v-if="row.items.length > 3" class="more-goods">+{{ row.items.length - 3 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('orders.buyer')" min-width="150">
          <template #default="{ row }">
            <div class="buyer-info">
              <p class="buyer-name">{{ row.user_name }}</p>
              <p class="buyer-phone">{{ row.user_phone }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('orders.amount')" width="140" align="right">
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="total-amount">{{ row.currency }} {{ formatAmount(row.pay_amount) }}</p>
              <p class="item-count">{{ $t('orders.itemsCount', { count: row.item_count }) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light" size="small">
              {{ row.status_text }}
            </el-tag>
            <el-tag
              v-if="row.fulfillment_status !== undefined && row.fulfillment_status !== 'pending'"
              :type="getFulfillmentType(row.fulfillment_status)"
              effect="plain"
              size="small"
              style="margin-top: 4px"
            >
              {{ row.fulfillment_text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ formatTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="canShip(row)"
              type="primary"
              size="small"
              @click="handleShip(row)"
            >
              {{ $t('orders.ship') }}
            </el-button>
            <el-button
              v-if="row.status === 'pending_payment'"
              type="warning"
              size="small"
              @click="handleRemind(row)"
            >
              {{ $t('orders.remind') }}
            </el-button>
            <el-button type="primary" link size="small" @click="handleDetail(row)">
              {{ $t('common.detail') }}
            </el-button>
            <el-dropdown @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button type="primary" link size="small">
                {{ $t('common.more') }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="remark">
                    <el-icon><Edit /></el-icon>{{ $t('orders.editRemark') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="row.status === 'pending'" command="adjust">
                    <el-icon><PriceTag /></el-icon>{{ $t('orders.adjustPrice') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="canCancel(row)" command="cancel" divided>
                    <el-icon><Close /></el-icon>{{ $t('orders.cancelOrder') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Order Detail Drawer -->
    <OrderDetailDrawer
      v-model="detailDrawerVisible"
      :order-id="currentOrderId"
      @refresh="loadOrders"
    />

    <!-- Ship Dialog -->
    <ShipDialog
      v-model="shipDialogVisible"
      :order="currentOrderForShip"
      :carriers="carriers"
      @success="handleShipSuccess"
    />

    <!-- Remark Dialog -->
    <RemarkDialog
      v-model="remarkDialogVisible"
      :order-id="currentOrderId"
      :current-remark="currentOrderRemark"
      @success="handleRemarkSuccess"
    />

    <!-- Adjust Price Dialog -->
    <AdjustPriceDialog
      v-model="adjustPriceDialogVisible"
      :order-id="currentOrderId"
      :current-amount="currentOrderAmount"
      :currency="currentOrderCurrency"
      @success="handleAdjustPriceSuccess"
    />

    <!-- Batch Ship Dialog -->
    <BatchShipDialog
      v-model="batchShipDialogVisible"
      :orders="selectedOrders"
      @success="handleBatchShipSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Download,
  Refresh,
  Picture,
  Check,
  Van,
  Edit,
  ArrowDown,
  PriceTag,
  Close
} from '@element-plus/icons-vue'
import {
  getOrderList,
  getFulfillmentSummary,
  remindPayment,
  cancelOrder,
  exportOrders,
  getCarrierList,
  type Order,
  type Carrier,
  type OrderStatus,
  type FulfillmentStatus,
  type OrderListParams,
  type ExportOrdersParams
} from '@/api/order'
import {
  OrderDetailDrawer,
  ShipDialog,
  AdjustPriceDialog,
  RemarkDialog,
  BatchShipDialog
} from './components'
import { t } from '@/plugins/i18n'

// Loading states
const loading = ref(false)
const exporting = ref(false)

// Pagination
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// Filters
const searchQuery = ref('')
const statusFilter = ref<OrderStatus | ''>('')
const fulfillmentFilter = ref<FulfillmentStatus | ''>('')
const dateRange = ref<[string, string] | null>(null)

// Data
const orderList = ref<Order[]>([])
const carriers = ref<Carrier[]>([])

// Statistics
const orderStats = ref({
  pending_payment: 0,
  partial_shipped: 0,
  shipped: 0,
  delivered: 0
})

// Selection
const tableRef = ref()
const selectedOrders = ref<Order[]>([])

// Dialog states
const detailDrawerVisible = ref(false)
const shipDialogVisible = ref(false)
const remarkDialogVisible = ref(false)
const adjustPriceDialogVisible = ref(false)
const batchShipDialogVisible = ref(false)

// Current order for actions
const currentOrderId = ref('')
const currentOrderForShip = ref<Order | null>(null)
const currentOrderRemark = ref('')
const currentOrderAmount = ref('0')
const currentOrderCurrency = ref('CNY')

// Load orders
const loadOrders = async () => {
  loading.value = true
  try {
    const params: OrderListParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (statusFilter.value) params.status = statusFilter.value as OrderStatus
    if (fulfillmentFilter.value !== '') params.fulfillment_status = fulfillmentFilter.value as FulfillmentStatus
    if (searchQuery.value) params.order_no = searchQuery.value
    if (dateRange.value && dateRange.value[0]) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }

    const res = await getOrderList(params)
    orderList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to load orders:', error)
    ElMessage.error(t('orders.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Load statistics
const loadStats = async () => {
  try {
    // Fetch fulfillment stats and pending_payment count in parallel
    const [fulfillmentRes, pendingPaymentRes] = await Promise.all([
      getFulfillmentSummary(),
      getOrderList({ status: 'pending_payment', page: 1, page_size: 1 })
    ])

    orderStats.value = {
      pending_payment: pendingPaymentRes.total || 0,
      partial_shipped: fulfillmentRes.partial_shipped || 0,
      shipped: fulfillmentRes.shipped || 0,
      delivered: fulfillmentRes.delivered || 0
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
    ElMessage.error(t('orders.loadStatsFailed'))
  }
}

// Load carriers
const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res.filter(c => c.is_active)
  } catch (error) {
    console.error('Failed to load carriers:', error)
    ElMessage.error(t('orders.loadCarriersFailed'))
  }
}

// Format helpers
const formatAmount = (amount: string | undefined) => {
  if (!amount) return '0.00'
  return parseFloat(amount).toFixed(2)
}

const formatTime = (time: string | undefined | null) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

// Status helpers
const getStatusType = (status: OrderStatus) => {
  const types: Record<OrderStatus, string> = {
    pending_payment: 'warning',
    paid: 'success',
    pending_shipment: 'warning',
    shipped: 'info',
    completed: 'success',
    cancelled: 'danger',
    refunding: 'warning',
    refunded: 'info'
  }
  return types[status] || 'info'
}

const getFulfillmentType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'partial_shipped': 'primary',
    'shipped': 'info',
    'delivered': 'success'
  }
  return types[status] || 'info'
}

// Action helpers
const canShip = (order: Order) => {
  return order.status === 'paid' &&
    (order.fulfillment_status === 'pending' || order.fulfillment_status === 'partial_shipped')
}

const canCancel = (order: Order) => {
  return ['pending_payment', 'paid'].includes(order.status)
}

// Event handlers
const handleSearch = () => {
  currentPage.value = 1
  loadOrders()
}

const handleRefresh = () => {
  loadOrders()
  loadStats()
}

const handleSizeChange = () => {
  currentPage.value = 1
  loadOrders()
}

const handleCurrentChange = () => {
  loadOrders()
}

const handleSelectionChange = (selection: Order[]) => {
  selectedOrders.value = selection
}

const clearSelection = () => {
  tableRef.value?.clearSelection()
}

// Export
const handleExport = async () => {
  exporting.value = true
  let url: string | null = null
  try {
    const params: ExportOrdersParams = {}
    if (statusFilter.value) params.status = statusFilter.value as OrderStatus
    if (searchQuery.value) params.order_no = searchQuery.value
    if (dateRange.value && dateRange.value[0]) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }

    const response = await exportOrders(params)
    // Create download link
    url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `orders_${new Date().toISOString().slice(0, 10)}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success(t('orders.exportCompleted'))
  } catch (error) {
    ElMessage.error(t('orders.exportFailed'))
  } finally {
    if (url) {
      window.URL.revokeObjectURL(url)
    }
    exporting.value = false
  }
}

// Order actions
const handleDetail = (row: Order) => {
  currentOrderId.value = row.order_id
  detailDrawerVisible.value = true
}

const handleShip = (row: Order) => {
  currentOrderForShip.value = row
  currentOrderId.value = row.order_id
  shipDialogVisible.value = true
}

const handleRemind = async (row: Order) => {
  try {
    await remindPayment(row.order_id)
    ElMessage.success(t('orders.reminderSent', { name: row.user_name }))
  } catch (error) {
    ElMessage.error(t('orders.reminderFailed'))
  }
}

const handleCommand = (command: string, row: Order) => {
  currentOrderId.value = row.order_id
  switch (command) {
    case 'remark':
      currentOrderRemark.value = row.seller_remark || ''
      remarkDialogVisible.value = true
      break
    case 'adjust':
      currentOrderAmount.value = row.pay_amount
      currentOrderCurrency.value = row.currency
      adjustPriceDialogVisible.value = true
      break
    case 'cancel':
      handleCancel(row)
      break
  }
}

const handleCancel = async (row: Order) => {
  try {
    const { value } = await ElMessageBox.prompt(
      t('orders.cancelReasonRequired'),
      t('orders.cancelOrder'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /^.{5,200}$/,
        inputErrorMessage: t('orders.cancelReasonPlaceholder')
      }
    )
    await cancelOrder(row.order_id, value)
    ElMessage.success(t('orders.orderCancelled'))
    loadOrders()
    loadStats()
  } catch {
    // User cancelled
  }
}

// Batch actions
const handleBatchShip = () => {
  const shippableOrders = selectedOrders.value.filter(canShip)
  if (shippableOrders.length === 0) {
    ElMessage.warning(t('orders.noShippableOrders'))
    return
  }
  if (shippableOrders.length < selectedOrders.value.length) {
    ElMessage.warning(t('orders.cannotShipCount', { count: selectedOrders.value.length - shippableOrders.length }))
  }
  selectedOrders.value = shippableOrders
  batchShipDialogVisible.value = true
}

const handleBatchRemark = () => {
  ElMessage.info(t('orders.batchRemarkComing'))
}

// Success handlers
const handleShipSuccess = () => {
  ElMessage.success(t('orders.orderShipped'))
  loadOrders()
  loadStats()
}

const handleRemarkSuccess = () => {
  ElMessage.success(t('orders.remarkUpdated'))
  loadOrders()
}

const handleAdjustPriceSuccess = () => {
  ElMessage.success(t('orders.priceAdjusted'))
  loadOrders()
}

const handleBatchShipSuccess = () => {
  ElMessage.success(t('orders.batchShipCompleted'))
  clearSelection()
  loadOrders()
  loadStats()
}

// Initialize
onMounted(() => {
  loadOrders()
  loadStats()
  loadCarriers()
})
</script>

<style scoped>
.orders-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-item {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 6px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 220px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.date-picker {
  width: 260px;
}

.date-picker :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Batch Actions Bar */
.batch-actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  margin-bottom: 16px;
}

.batch-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #4F46E5;
}

.batch-buttons {
  display: flex;
  gap: 8px;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Table */
.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Order Detail Expand */
.order-detail {
  padding: 20px;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
}

.order-detail h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 12px 0;
}

.order-detail p {
  font-size: 13px;
  color: #4B5563;
  margin: 0 0 8px 0;
}

.order-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid rgba(99, 102, 241, 0.1);
}

.order-item:last-child {
  border-bottom: none;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  overflow: hidden;
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
  margin: 0 0 4px 0;
}

.item-price {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Order No Cell */
.order-no-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-no {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
  cursor: pointer;
}

.order-no:hover {
  text-decoration: underline;
}

/* Goods Preview */
.goods-preview {
  display: flex;
  align-items: center;
  gap: 4px;
}

.goods-thumb {
  width: 50px;
  height: 50px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #E5E7EB;
  transition: transform 0.2s ease;
}

.goods-thumb:hover {
  transform: scale(1.05);
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.more-goods {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 10px;
  font-size: 12px;
  color: #6366F1;
  font-weight: 600;
}

/* Buyer Info */
.buyer-info {
  line-height: 1.5;
}

.buyer-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.buyer-phone {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
  font-family: 'Fira Code', monospace;
}

/* Amount Cell */
.amount-cell {
  text-align: right;
}

.total-amount {
  font-size: 16px;
  font-weight: 700;
  color: #EF4444;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.item-count {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

/* Time Text */
.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select,
  .date-picker {
    width: 100%;
  }

  .goods-preview {
    flex-wrap: wrap;
  }

  .stat-item {
    border-radius: 14px;
    padding: 20px;
  }

  .stat-number {
    font-size: 28px;
  }

  .filter-card,
  .table-card {
    border-radius: 14px;
  }

  .batch-actions-bar {
    flex-direction: column;
    gap: 12px;
  }

  .batch-buttons {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>