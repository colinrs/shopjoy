<template>
  <div class="orders-page">
    <!-- Statistics Cards -->
    <OrderStatsCards :stats="orderStats" />

    <!-- Filter Bar -->
    <OrderFilterBar
      v-model:searchQuery="searchQuery"
      v-model:statusFilter="statusFilter"
      v-model:fulfillmentFilter="fulfillmentFilter"
      v-model:dateRange="dateRange"
      :exporting="exporting"
      @search="handleSearch"
      @refresh="handleRefresh"
      @export="handleExport"
    />

    <!-- Batch Actions Bar -->
    <OrderBatchActions
      :selected-count="selectedOrders.length"
      @batch-ship="handleBatchShip"
      @batch-cancel="handleBatchCancel"
      @batch-remark="handleBatchRemark"
      @clear-selection="clearSelection"
    />

    <!-- Orders Table -->
    <OrderTable
      ref="tableRef"
      :order-list="orderList"
      :loading="loading"
      @selection-change="handleSelectionChange"
      @detail="handleDetail"
      @ship="handleShip"
      @remind="handleRemind"
      @command="handleCommand"
    />

    <!-- Pagination -->
    <OrderPagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      @page-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />

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
  getOrderList,
  remindPayment,
  cancelOrder,
  exportOrders,
  getCarrierList,
  batchCancelOrders,
  type Order,
  type Carrier,
  type OrderStatus,
  type FulfillmentStatus,
  type OrderListParams,
  type ExportOrdersParams,
  type BatchCancelOrderRequest
} from '@/api/order'
import { getFulfillmentSummary } from '@/api/fulfillment'
import { getOrderStatusDistribution } from '@/api/dashboard'
import {
  OrderDetailDrawer,
  ShipDialog,
  AdjustPriceDialog,
  RemarkDialog,
  BatchShipDialog,
  OrderStatsCards,
  OrderFilterBar,
  OrderBatchActions,
  OrderTable,
  OrderPagination
} from './components'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

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

// Order statistics type
interface OrderStats {
  pending_payment: number
  partial_shipped: number
  shipped: number
  delivered: number
}

// Statistics
// Note: FulfillmentSummary API (getFulfillmentSummary) does not include pending_payment count.
// pending_payment is fetched separately via getOrderStatusDistribution() API.
const orderStats = ref<OrderStats>({
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
const currentOrderId = ref<number>(0)
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
    handleError(error, t('orders.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Load statistics
// Note: FulfillmentSummary API does not return pending_payment count, so we fetch it separately
// via getOrderStatusDistribution() which provides counts for all order statuses.
const loadStats = async () => {
  try {
    const [fulfillmentRes, statusDistRes] = await Promise.all([
      getFulfillmentSummary(),
      getOrderStatusDistribution()
    ])

    // Find pending_payment count from status distribution
    const pendingPaymentItem = statusDistRes.list.find(item => item.status === 'pending_payment')

    orderStats.value = {
      pending_payment: pendingPaymentItem?.count || 0,
      partial_shipped: fulfillmentRes.partial_shipped || 0,
      shipped: fulfillmentRes.shipped || 0,
      delivered: fulfillmentRes.delivered || 0
    }
  } catch (error) {
    handleError(error, t('orders.loadStatsFailed'))
  }
}

// Load carriers
const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res.filter(c => c.is_active)
  } catch (error) {
    handleError(error, t('orders.loadCarriersFailed'))
  }
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
  selectedOrders.value = []
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
        inputErrorMessage: t('orders.cancelReasonLengthError')
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
  const shippableOrders = selectedOrders.value.filter(order =>
    order.status === 'paid' &&
    (order.fulfillment_status === 0 || order.fulfillment_status === 1)
  )
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

// Batch cancel orders
const handleBatchCancel = async () => {
  const cancellableOrders = selectedOrders.value.filter(order =>
    ['pending_payment', 'paid'].includes(order.status)
  )
  if (cancellableOrders.length === 0) {
    ElMessage.warning(t('orders.noCancellableOrders'))
    return
  }
  if (cancellableOrders.length < selectedOrders.value.length) {
    ElMessage.warning(t('orders.cannotCancelCount', { count: selectedOrders.value.length - cancellableOrders.length }))
  }

  try {
    const { value } = await ElMessageBox.prompt(
      t('orders.cancelReasonRequired'),
      t('orders.batchCancel'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /^.{5,200}$/,
        inputErrorMessage: t('orders.cancelReasonLengthError')
      }
    )

    const data: BatchCancelOrderRequest = {
      order_ids: cancellableOrders.map(o => o.order_id),
      reason: value
    }

    const res = await batchCancelOrders(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('orders.batchCancelPartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('orders.batchCancelSuccess', { count: successCount }))
    }

    clearSelection()
    loadOrders()
    loadStats()
  } catch {
    // User cancelled
  }
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
</style>
