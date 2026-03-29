<template>
  <div class="shipments-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-item pending" @click="handleStatusFilter(0)">
          <p class="stat-number">{{ stats.pending }}</p>
          <p class="stat-label">待发货</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item shipped" @click="handleStatusFilter(1)">
          <p class="stat-number">{{ stats.shipped }}</p>
          <p class="stat-label">已发货</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item transit" @click="handleStatusFilter(2)">
          <p class="stat-number">{{ stats.in_transit }}</p>
          <p class="stat-label">运输中</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item delivered" @click="handleStatusFilter(3)">
          <p class="stat-number">{{ stats.delivered }}</p>
          <p class="stat-label">已送达</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索运单号/发货单号"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="发货状态" clearable class="filter-select">
            <el-option label="全部" value="" />
            <el-option label="待发货" :value="0" />
            <el-option label="已发货" :value="1" />
            <el-option label="运输中" :value="2" />
            <el-option label="已送达" :value="3" />
            <el-option label="配送失败" :value="4" />
          </el-select>
          <el-select v-model="carrierFilter" placeholder="物流公司" clearable class="filter-select">
            <el-option label="全部" value="" />
            <el-option v-for="carrier in carriers" :key="carrier.code" :label="carrier.name" :value="carrier.code" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            class="date-picker"
            value-format="YYYY-MM-DD"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Batch Actions Bar -->
    <transition name="slide-down">
      <div v-if="selectedRows.length > 0" class="batch-actions-bar">
        <div class="batch-info">
          <el-icon><InfoFilled /></el-icon>
          <span>已选择 <strong>{{ selectedRows.length }}</strong> 条发货单</span>
        </div>
        <div class="batch-actions">
          <el-button type="primary" @click="handleBatchShip">
            <el-icon><Van /></el-icon>批量发货
          </el-button>
          <el-button @click="clearSelection">取消选择</el-button>
        </div>
      </div>
    </transition>

    <!-- Shipments Table -->
    <el-card class="table-card" shadow="never">
      <el-table
        ref="tableRef"
        :data="shipmentList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" :selectable="isSelectable" />
        <el-table-column prop="shipment_no" label="发货单号" min-width="150">
          <template #default="{ row }">
            <div class="shipment-no-cell">
              <span class="shipment-no">{{ row.shipment_no }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="订单号" min-width="150">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewOrder(row.order_id)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="物流信息" min-width="180">
          <template #default="{ row }">
            <div class="logistics-info">
              <p class="carrier-name">{{ row.carrier }}</p>
              <p class="tracking-no" v-if="row.tracking_no">
                <el-icon><Location /></el-icon>
                {{ row.tracking_no }}
              </p>
              <p class="tracking-no" v-else>
                <el-tag type="warning" size="small">待录入</el-tag>
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="商品" min-width="160">
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
        <el-table-column prop="status" label="状态" width="110" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="statusTypeMap" />
          </template>
        </el-table-column>
        <el-table-column label="发货时间" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.shipped_at || row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              详情
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              size="small"
              @click="openShipDialog(row)"
            >
              发货
            </el-button>
            <el-button
              v-if="row.status === 2"
              type="primary"
              link
              size="small"
              @click="markDelivered(row)"
            >
              确认送达
            </el-button>
            <el-button
              v-if="row.tracking_no"
              type="primary"
              link
              size="small"
              @click="copyTracking(row.tracking_no)"
            >
              复制单号
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Download, Refresh, InfoFilled, Van, Location, Goods } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ShipDialog from '@/components/fulfillment/ShipDialog.vue'
import BatchShipDialog from '@/components/fulfillment/BatchShipDialog.vue'
import {
  getShipmentList,
  getCarrierList,
  updateShipmentStatus,
  type Shipment,
  type Carrier,
  type ShipmentListParams
} from '@/api/fulfillment'

const router = useRouter()

// State
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<number | ''>('')
const carrierFilter = ref('')
const dateRange = ref<[string, string] | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const shipDialogVisible = ref(false)
const batchShipDialogVisible = ref(false)
const currentShipment = ref<Shipment | null>(null)
const selectedRows = ref<Shipment[]>([])
const tableRef = ref()

const carriers = ref<Carrier[]>([])

const stats = ref({
  pending: 15,
  shipped: 42,
  in_transit: 28,
  delivered: 156
})

const statusTypeMap = {
  0: { type: 'warning' as const, text: '待发货' },
  1: { type: 'primary' as const, text: '已发货' },
  2: { type: 'info' as const, text: '运输中' },
  3: { type: 'success' as const, text: '已送达' },
  4: { type: 'danger' as const, text: '配送失败' }
}

// Mock data for development
const shipmentList = ref<Shipment[]>([
  {
    id: 1,
    shipment_no: 'SHP20260322001',
    order_id: 'ORD001',
    order_no: 'ORD2026031800100',
    status: 1,
    status_text: '已发货',
    carrier: 'SF Express',
    carrier_code: 'SF',
    tracking_no: 'SF1234567890',
    tracking_url: 'https://www.sf-express.com/sf-service-ow/f、梁.operation.entrega_ar?trackingNo=SF1234567890',
    shipping_cost: "12.00",
    shipping_currency: 'CNY',
    weight: "1.5",
    shipped_at: '2026-03-22T14:30:25Z',
    delivered_at: null,
    remark: '',
    created_at: '2026-03-22T10:00:00Z',
    updated_at: '2026-03-22T10:00:00Z',
    created_by: 1,
    created_by_name: '管理员',
    items: [
      { id: 1, product_id: 1, sku_id: 1, product_name: 'Wireless Bluetooth Earphones Pro', sku_name: 'Black', image: '', quantity: 1 }
    ]
  },
  {
    id: 2,
    shipment_no: 'SHP20260322002',
    order_id: 'ORD002',
    order_no: 'ORD2026031800099',
    status: 0,
    status_text: '待发货',
    carrier: '',
    carrier_code: '',
    tracking_no: '',
    tracking_url: '',
    shipping_cost: "0.00",
    shipping_currency: 'CNY',
    weight: "0",
    shipped_at: null,
    delivered_at: null,
    remark: '',
    created_at: '2026-03-22T11:00:00Z',
    updated_at: '2026-03-22T11:00:00Z',
    created_by: 1,
    created_by_name: '管理员',
    items: [
      { id: 2, product_id: 2, sku_id: 2, product_name: 'Smart Watch Series 7', sku_name: 'Silver 42mm', image: '', quantity: 1 },
      { id: 3, product_id: 3, sku_id: 3, product_name: 'Portable Power Bank', sku_name: '20000mAh', image: '', quantity: 2 }
    ]
  },
  {
    id: 3,
    shipment_no: 'SHP20260322003',
    order_id: 'ORD003',
    order_no: 'ORD2026031800098',
    status: 2,
    status_text: '运输中',
    carrier: 'ZTO Express',
    carrier_code: 'ZT',
    tracking_no: 'ZT9876543210',
    tracking_url: 'https://www.zto.com/express/express/search?num=ZT9876543210',
    shipping_cost: "8.00",
    shipping_currency: 'CNY',
    weight: "0.8",
    shipped_at: '2026-03-21T16:45:30Z',
    delivered_at: null,
    remark: 'Fragile item',
    created_at: '2026-03-21T15:00:00Z',
    updated_at: '2026-03-21T15:00:00Z',
    created_by: 1,
    created_by_name: '管理员',
    items: [
      { id: 4, product_id: 4, sku_id: 4, product_name: 'Mechanical Keyboard RGB', sku_name: 'Blue Switch', image: '', quantity: 1 }
    ]
  },
  {
    id: 4,
    shipment_no: 'SHP20260322004',
    order_id: 'ORD004',
    order_no: 'ORD2026031800097',
    status: 3,
    status_text: '已送达',
    carrier: 'SF Express',
    carrier_code: 'SF',
    tracking_no: 'SF1122334455',
    tracking_url: 'https://www.sf-express.com/sf-service-ow/f、梁.operation.entrega_ar?trackingNo=SF1122334455',
    shipping_cost: "15.00",
    shipping_currency: 'CNY',
    weight: "3.2",
    shipped_at: '2026-03-20T09:00:00Z',
    delivered_at: '2026-03-22T15:30:00Z',
    remark: '',
    created_at: '2026-03-20T08:00:00Z',
    updated_at: '2026-03-22T15:30:00Z',
    created_by: 1,
    created_by_name: '管理员',
    items: [
      { id: 5, product_id: 5, sku_id: 5, product_name: '4K HD Monitor 27 inch', sku_name: 'Black', image: '', quantity: 1 }
    ]
  },
  {
    id: 5,
    shipment_no: 'SHP20260322005',
    order_id: 'ORD005',
    order_no: 'ORD2026031800096',
    status: 4,
    status_text: '配送失败',
    carrier: 'YTO Express',
    carrier_code: 'YT',
    tracking_no: 'YT5566778899',
    tracking_url: 'https://www.yto.net.cn/express/express/search?num=YT5566778899',
    shipping_cost: "6.00",
    shipping_currency: 'CNY',
    weight: "0.5",
    shipped_at: '2026-03-19T14:00:00Z',
    delivered_at: null,
    remark: 'Delivery failed - wrong address',
    created_at: '2026-03-19T13:00:00Z',
    updated_at: '2026-03-19T14:00:00Z',
    created_by: 1,
    created_by_name: '管理员',
    items: [
      { id: 6, product_id: 6, sku_id: 6, product_name: 'Phone Case', sku_name: 'iPhone 15 Pro', image: '', quantity: 1 }
    ]
  }
])

// Methods
const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res
  } catch (error) {
    // Use mock data
    carriers.value = [
      { code: 'SF', name: 'SF Express', tracking_url: 'https://www.sf-express.com/track?id={tracking_no}', is_active: true },
      { code: 'ZT', name: 'ZTO Express', tracking_url: 'https://www.zto.com/track?id={tracking_no}', is_active: true },
      { code: 'YT', name: 'YTO Express', tracking_url: 'https://www.yto.net.cn/query.html?id={tracking_no}', is_active: true },
      { code: 'ST', name: 'STO Express', tracking_url: 'https://www.sto.cn/track?id={tracking_no}', is_active: true },
      { code: 'YD', name: 'Yunda Express', tracking_url: 'https://www.yundaex.com/track?id={tracking_no}', is_active: true },
      { code: 'EMS', name: 'EMS', tracking_url: 'https://www.ems.com.cn/track?id={tracking_no}', is_active: true },
      { code: 'JD', name: 'JD Logistics', tracking_url: 'https://www.jdl.com/track?id={tracking_no}', is_active: true },
      { code: 'OTHER', name: 'Other', tracking_url: '', is_active: true }
    ]
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
      params.status = statusFilter.value
    }
    const res = await getShipmentList(params)
    shipmentList.value = res.list
    total.value = res.total
  } catch (error) {
    // Mock data already set
  } finally {
    loading.value = false
  }
}

const handleStatusFilter = (status: number) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = () => {
  ElMessage.success('Export successful')
}

const handleRefresh = () => {
  loadData()
}

const handleSelectionChange = (rows: Shipment[]) => {
  selectedRows.value = rows
}

const isSelectable = (row: Shipment) => {
  return row.status === 0
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
  ElMessage.success('Shipment created successfully')
}

const handleBatchShipSuccess = () => {
  batchShipDialogVisible.value = false
  clearSelection()
  loadData()
  ElMessage.success('Batch shipment completed')
}

const markDelivered = async (row: Shipment) => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to mark this shipment as delivered?',
      'Confirm Delivery',
      {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'success'
      }
    )
    await updateShipmentStatus(row.id, 3)
    ElMessage.success('Shipment marked as delivered')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to update status')
    }
  }
}

const copyTracking = (trackingNo: string) => {
  navigator.clipboard.writeText(trackingNo)
  ElMessage.success('Tracking number copied')
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
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
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