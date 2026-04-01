<template>
  <div class="shipments-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-item pending" shadow="hover" @click="handleStatusFilter('0')">
          <template #header>
            <div class="stat-header">
              <span>{{ $t('fulfillment.pendingShipment') }}</span>
              <el-tooltip :content="$t('fulfillment.pendingShipmentTip')">
                <el-icon class="stat-tooltip"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <p class="stat-number">{{ stats.pending }}</p>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-item shipped" shadow="hover" @click="handleStatusFilter('1')">
          <template #header>
            <div class="stat-header">
              <span>{{ $t('fulfillment.shippedShipment') }}</span>
              <el-tooltip :content="$t('fulfillment.shippedShipmentTip')">
                <el-icon class="stat-tooltip"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <p class="stat-number">{{ stats.shipped }}</p>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-item delivered" shadow="hover" @click="handleStatusFilter('3')">
          <template #header>
            <div class="stat-header">
              <span>{{ $t('fulfillment.deliveredShipment') }}</span>
              <el-tooltip :content="$t('fulfillment.deliveredShipmentTip')">
                <el-icon class="stat-tooltip"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <p class="stat-number">{{ stats.delivered }}</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('fulfillment.searchTrackingNo')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" :placeholder="$t('fulfillment.status')" clearable class="filter-select">
            <el-option :label="$t('common.all')" value="" />
            <el-option :label="$t('fulfillment.pendingShipment')" value="0" />
            <el-option :label="$t('fulfillment.shippedShipment')" value="1" />
            <el-option :label="$t('fulfillment.inTransitShipment')" value="2" />
            <el-option :label="$t('fulfillment.deliveredShipment')" value="3" />
            <el-option :label="$t('fulfillment.failedShipment')" value="4" />
          </el-select>
          <el-select v-model="carrierFilter" :placeholder="$t('fulfillment.carrier')" clearable class="filter-select">
            <el-option :label="$t('common.all')" value="" />
            <el-option v-for="carrier in carriers" :key="carrier.code" :label="carrier.name" :value="carrier.code" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="to"
            :start-placeholder="$t('fulfillment.startDate')"
            :end-placeholder="$t('fulfillment.endDate')"
            class="date-picker"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>{{ $t('common.export') }}
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>{{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Batch Actions Bar -->
    <transition name="slide-down">
      <div v-if="selectedRows.length > 0" class="batch-actions-bar">
        <div class="batch-info">
          <el-icon><InfoFilled /></el-icon>
          <span>{{ $t('fulfillment.selectedRows', { n: selectedRows.length }) }}</span>
        </div>
        <div class="batch-actions">
          <el-button type="primary" @click="handleBatchShip">
            <el-icon><Van /></el-icon>{{ $t('fulfillment.batchShipment') }}
          </el-button>
          <el-button type="warning" @click="handleBatchUpdateTracking">
            <el-icon><Edit /></el-icon>{{ $t('fulfillment.batchUpdateTracking') }}
          </el-button>
          <el-button @click="handleImportTracking">
            <el-icon><Upload /></el-icon>{{ $t('fulfillment.importTracking') }}
          </el-button>
          <el-button @click="clearSelection">{{ $t('fulfillment.clearSelection') }}</el-button>
        </div>
      </div>
    </transition>

    <!-- Shipments Table -->
    <el-card class="table-card" shadow="never">
      <EmptyState
        v-if="shipmentList.length === 0 && !loading"
        :title="$t('fulfillment.noShipments')"
        :description="$t('fulfillment.noShipmentsDesc')"
      />
      <el-table
        v-else
        ref="tableRef"
        :data="shipmentList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" :selectable="isSelectable" />
        <el-table-column prop="shipment_no" :label="$t('fulfillment.shipmentNo')" min-width="150">
          <template #default="{ row }">
            <div class="shipment-no-cell">
              <span class="shipment-no">{{ row.shipment_no }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" :label="$t('fulfillment.orderNo')" min-width="150">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewOrder(row.order_id)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('fulfillment.logistics')" min-width="180">
          <template #default="{ row }">
            <div class="logistics-info">
              <p class="carrier-name">{{ row.carrier }}</p>
              <p class="tracking-no" v-if="row.tracking_no">
                <el-icon><Location /></el-icon>
                {{ row.tracking_no }}
              </p>
              <p class="tracking-no" v-else>
                <el-tag type="warning" size="small">{{ $t('fulfillment.toBeEntered') }}</el-tag>
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.items')" min-width="160">
          <template #default="{ row }">
            <div class="items-preview">
              <el-image
                v-for="(item, idx) in row.items.slice(0, 3)"
                :key="idx"
                :src="item.image"
                class="item-thumb"
                fit="cover"
              >
                <template #error>
                  <div class="thumb-placeholder">
                    <el-icon><Goods /></el-icon>
                  </div>
                </template>
              </el-image>
              <span v-if="row.items.length > 3" class="more-items">+{{ row.items.length - 3 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="110" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="statusTypeMap" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('fulfillment.shipmentTime')" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.shipped_at || row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              {{ $t('common.detail') }}
            </el-button>
            <el-button
              v-if="row.status === ShipmentStatusMap.PENDING"
              type="success"
              link
              size="small"
              @click="openShipDialog(row)"
            >
              {{ $t('orders.ship') }}
            </el-button>
            <el-button
              v-if="row.status === ShipmentStatusMap.IN_TRANSIT"
              type="primary"
              link
              size="small"
              @click="markDelivered(row)"
            >
              {{ $t('fulfillment.confirmDeliveryAction') }}
            </el-button>
            <el-button
              v-if="row.tracking_no"
              type="primary"
              link
              size="small"
              @click="copyTracking(row.tracking_no)"
            >
              {{ $t('fulfillment.copyTracking') }}
            </el-button>
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

    <!-- Ship Dialog -->
    <ship-dialog
      v-model="shipDialogVisible"
      :shipment="currentShipment"
      :carriers="carriers"
      @success="handleShipSuccess"
    />

    <!-- Batch Ship Dialog -->
    <batch-ship-dialog
      v-model="batchShipDialogVisible"
      :shipments="selectedRows"
      :carriers="carriers"
      @success="handleBatchShipSuccess"
    />

    <!-- Import Tracking Dialog -->
    <import-tracking-dialog
      v-model="importTrackingDialogVisible"
      @success="handleImportSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Download, Refresh, InfoFilled, Van, Location, Goods, QuestionFilled, Upload, Edit } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ShipDialog from '@/components/fulfillment/ShipDialog.vue'
import BatchShipDialog from '@/components/fulfillment/BatchShipDialog.vue'
import ImportTrackingDialog from '@/components/fulfillment/ImportTrackingDialog.vue'
import { t } from '@/plugins/i18n'
import {
  getShipmentList,
  getCarrierList,
  updateShipmentStatus,
  getFulfillmentSummary,
  exportShipmentsUrl,
  batchUpdateTracking,
  type Shipment,
  type Carrier,
  type ShipmentListParams,
  ShipmentStatusMap,
  type BatchUpdateTrackingRequest
} from '@/api/fulfillment'
import { downloadFile } from '@/utils/download'

const router = useRouter()

// State
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<string | ''>('')
const carrierFilter = ref('')
const dateRange = ref<[string, string] | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const shipDialogVisible = ref(false)
const batchShipDialogVisible = ref(false)
const importTrackingDialogVisible = ref(false)
const currentShipment = ref<Shipment | null>(null)
const selectedRows = ref<Shipment[]>([])
const tableRef = ref()

const carriers = ref<Carrier[]>([])

const stats = ref({
  pending: 0,
  shipped: 0,
  delivered: 0
})

const statusTypeMap: Record<string, { type: 'warning' | 'primary' | 'info' | 'success' | 'danger', text: string }> = {
  '0': { type: 'warning', text: 'Pending' },
  '1': { type: 'primary', text: 'Shipped' },
  '2': { type: 'info', text: 'In Transit' },
  '3': { type: 'success', text: 'Delivered' },
  '4': { type: 'danger', text: 'Failed' },
  '5': { type: 'info', text: 'Cancelled' }
}

// Shipment list
const shipmentList = ref<Shipment[]>([])

// Methods
const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res
  } catch (error) {
    ElMessage.error(t('fulfillment.loadCarriersFailed'))
  }
}

const loadStats = async () => {
  try {
    const res = await getFulfillmentSummary()
    stats.value = {
      pending: res.pending_shipment || 0,
      shipped: res.shipped || 0,
      delivered: res.delivered || 0
    }
  } catch (error) {
    ElMessage.error(t('fulfillment.loadFailed'))
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const params: ShipmentListParams = {
      page: currentPage.value,
      page_size: pageSize.value,
      carrier_code: carrierFilter.value || undefined,
      tracking_no: searchQuery.value || undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    }
    if (statusFilter.value !== '') {
      params.status = statusFilter.value as import('@/api/fulfillment').ShipmentStatus
    }
    const res = await getShipmentList(params)
    shipmentList.value = res.list
    total.value = res.total
  } catch (error) {
    ElMessage.error(t('fulfillment.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleStatusFilter = (status: string) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = async () => {
  try {
    const { url, params } = exportShipmentsUrl({
      tracking_no: searchQuery.value || undefined,
      status: statusFilter.value !== '' ? statusFilter.value : undefined,
      carrier_code: carrierFilter.value || undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    })

    await downloadFile(url, params)
  } catch (error) {
    console.error('Export failed:', error)
    // Error message is handled by downloadFile utility
  }
}

const handleRefresh = () => {
  loadData()
}

const handleSelectionChange = (rows: Shipment[]) => {
  selectedRows.value = rows
}

const isSelectable = (row: Shipment) => {
  return row.status === ShipmentStatusMap.PENDING
}

const clearSelection = () => {
  tableRef.value?.clearSelection()
}

const handleBatchShip = () => {
  batchShipDialogVisible.value = true
}

const viewDetail = (row: Shipment) => {
  router.push(`/fulfillment/shipments/${row.id}`)
}

const viewOrder = (orderId: string) => {
  router.push(`/orders?id=${orderId}`)
}

const openShipDialog = (row: Shipment) => {
  currentShipment.value = row
  shipDialogVisible.value = true
}

const handleShipSuccess = () => {
  shipDialogVisible.value = false
  loadData()
  ElMessage.success(t('fulfillment.shipmentConfirmed'))
}

const handleBatchShipSuccess = () => {
  batchShipDialogVisible.value = false
  clearSelection()
  loadData()
  ElMessage.success(t('fulfillment.batchShipmentCompleted'))
}

const handleImportTracking = () => {
  importTrackingDialogVisible.value = true
}

const handleImportSuccess = () => {
  importTrackingDialogVisible.value = false
  loadData()
}

// Batch update tracking for selected shipments
const handleBatchUpdateTracking = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('fulfillment.selectAtLeastOneShipment'))
    return
  }

  try {
    const { value: trackingNo } = await ElMessageBox.prompt(
      t('fulfillment.enterTrackingNumber'),
      t('fulfillment.batchUpdateTracking'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /^.{5,50}$/,
        inputErrorMessage: t('fulfillment.trackingNumberLengthError')
      }
    )

    // Use first selected row's carrier as default
    const defaultCarrier = selectedRows.value[0]?.carrier_code || ''

    const { value: carrierCode } = await ElMessageBox.confirm(
      t('fulfillment.selectCarrierForBatch'),
      t('fulfillment.batchUpdateTracking'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        showInput: true,
        inputPlaceholder: t('fulfillment.enterCarrierCode'),
        inputValue: defaultCarrier
      }
    )

    const data: BatchUpdateTrackingRequest = {
      shipment_ids: selectedRows.value.map(s => s.id),
      carrier_code: carrierCode,
      tracking_no: trackingNo
    }

    const res = await batchUpdateTracking(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('fulfillment.batchUpdateTrackingPartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('fulfillment.batchUpdateTrackingSuccess', { count: successCount }))
    }

    clearSelection()
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch update tracking:', error)
      ElMessage.error(t('fulfillment.batchUpdateTrackingFailed'))
    }
  }
}

const markDelivered = async (row: Shipment) => {
  try {
    await ElMessageBox.confirm(
      t('fulfillment.markDelivered'),
      t('fulfillment.confirmDelivery'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'success'
      }
    )
    await updateShipmentStatus(row.id, ShipmentStatusMap.DELIVERED)
    ElMessage.success(t('fulfillment.shipmentMarkedDelivered'))
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('fulfillment.loadFailed'))
    }
  }
}

const copyTracking = (trackingNo: string) => {
  navigator.clipboard.writeText(trackingNo)
  ElMessage.success(t('fulfillment.trackingNumberCopied'))
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

onMounted(() => {
  loadCarriers()
  loadStats()
  loadData()
})
</script>

<style scoped>
.shipments-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-item {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-item:hover {
  transform: translateY(-2px);
}

.stat-item.pending {
  border-left: 4px solid #F59E0B;
}

.stat-item.shipped {
  border-left: 4px solid #6366F1;
}

.stat-item.transit {
  border-left: 4px solid #3B82F6;
}

.stat-item.delivered {
  border-left: 4px solid #10B981;
}

.stat-item :deep(.el-card__header) {
  padding: 12px 20px;
  border-bottom: none;
}

.stat-item :deep(.el-card__body) {
  padding: 0 20px 20px;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-tooltip {
  color: #9CA3AF;
  cursor: help;
  transition: color 0.2s;
}

.stat-tooltip:hover {
  color: #6366F1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
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
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  padding: 12px 20px;
  margin-bottom: 16px;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.batch-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6366F1;
  font-size: 14px;
}

.batch-info strong {
  font-weight: 700;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

/* Table */
.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Shipment No */
.shipment-no {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
}

/* Logistics Info */
.logistics-info {
  line-height: 1.6;
}

.carrier-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.tracking-no {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'Fira Code', monospace;
}

/* Items Preview */
.items-preview {
  display: flex;
  align-items: center;
  gap: 4px;
}

.item-thumb {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #E5E7EB;
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

.more-items {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 8px;
  font-size: 12px;
  color: #6366F1;
  font-weight: 600;
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

/* Transitions */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
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

  .stat-item {
    border-radius: 14px;
    padding: 20px;
  }

  .stat-number {
    font-size: 28px;
  }
}
</style>
