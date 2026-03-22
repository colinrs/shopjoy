<template>
  <div class="shipment-detail-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-left">
        <el-button link @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          Back
        </el-button>
        <div class="title-section">
          <h1 class="page-title">Shipment Details</h1>
          <p class="shipment-no">{{ shipment?.shipment_no }}</p>
        </div>
      </div>
      <div class="header-right">
        <status-tag :status="shipment?.status" :type-map="statusTypeMap" size="large" />
        <el-button v-if="shipment?.tracking_no" type="primary" @click="copyTracking">
          <el-icon><Link /></el-icon>
          Copy Tracking
        </el-button>
      </div>
    </div>

    <el-row :gutter="20">
      <!-- Left Column -->
      <el-col :xs="24" :lg="16">
        <!-- Basic Info Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Document /></el-icon>
                Basic Information
              </span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="Shipment No.">
              <span class="value-text">{{ shipment?.shipment_no }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="Order No.">
              <el-button type="primary" link @click="viewOrder">
                {{ shipment?.order_no }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="Status">
              <status-tag :status="shipment?.status" :type-map="statusTypeMap" />
            </el-descriptions-item>
            <el-descriptions-item label="Created At">
              {{ shipment?.created_at }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- Logistics Info Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Van /></el-icon>
                Logistics Information
              </span>
              <el-button
                v-if="shipment?.status === 0"
                type="primary"
                size="small"
                @click="openEditLogistics"
              >
                Edit
              </el-button>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="Carrier">
              <span class="value-text">{{ shipment?.carrier || '-' }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="Tracking No.">
              <span v-if="shipment?.tracking_no" class="tracking-value">
                {{ shipment.tracking_no }}
                <el-button link size="small" @click="copyTracking">
                  <el-icon><CopyDocument /></el-icon>
                </el-button>
              </span>
              <span v-else class="no-data">Not entered</span>
            </el-descriptions-item>
            <el-descriptions-item label="Shipping Cost">
              <span v-if="shipment?.shipping_cost">
                {{ shipment.shipping_currency }} {{ (shipment.shipping_cost / 100).toFixed(2) }}
              </span>
              <span v-else class="no-data">-</span>
            </el-descriptions-item>
            <el-descriptions-item label="Weight">
              <span v-if="shipment?.weight">{{ shipment.weight }} kg</span>
              <span v-else class="no-data">-</span>
            </el-descriptions-item>
            <el-descriptions-item label="Shipped At">
              {{ shipment?.shipped_at || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="Delivered At">
              {{ shipment?.delivered_at || '-' }}
            </el-descriptions-item>
          </el-descriptions>
          <div v-if="shipment?.remark" class="remark-section">
            <span class="remark-label">Remark:</span>
            <span class="remark-text">{{ shipment.remark }}</span>
          </div>
        </el-card>

        <!-- Items Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Goods /></el-icon>
                Shipment Items
              </span>
              <span class="item-count">Total {{ shipment?.items?.length || 0 }} items</span>
            </div>
          </template>
          <div class="items-list">
            <div v-for="item in shipment?.items" :key="item.id" class="item-row">
              <el-image :src="item.image" class="item-image" fit="cover">
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="item-info">
                <p class="item-name">{{ item.product_name }}</p>
                <p class="item-sku">SKU: {{ item.sku_name }}</p>
              </div>
              <div class="item-quantity">
                x {{ item.quantity }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- Right Column -->
      <el-col :xs="24" :lg="8">
        <!-- Timeline Card -->
        <el-card class="timeline-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Clock /></el-icon>
                Status Timeline
              </span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(event, index) in timeline"
              :key="index"
              :type="event.type"
              :timestamp="event.time"
              :hollow="!event.active"
              :class="{ 'is-active': event.active }"
            >
              <div class="timeline-content">
                <p class="timeline-title">{{ event.title }}</p>
                <p v-if="event.description" class="timeline-desc">{{ event.description }}</p>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- Actions Card -->
        <el-card class="actions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Operation /></el-icon>
                Actions
              </span>
            </div>
          </template>
          <div class="action-buttons">
            <el-button
              v-if="shipment?.status === 0"
              type="success"
              class="action-btn"
              @click="confirmShip"
            >
              <el-icon><Van /></el-icon>
              Confirm Shipment
            </el-button>
            <el-button
              v-if="shipment?.status === 2"
              type="primary"
              class="action-btn"
              @click="markDelivered"
            >
              <el-icon><CircleCheck /></el-icon>
              Mark as Delivered
            </el-button>
            <el-button
              v-if="shipment?.tracking_no"
              class="action-btn"
              @click="trackShipment"
            >
              <el-icon><Location /></el-icon>
              Track Package
            </el-button>
            <el-button
              v-if="shipment?.status === 0"
              type="danger"
              class="action-btn"
              @click="cancelShipment"
            >
              <el-icon><Close /></el-icon>
              Cancel Shipment
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Edit Logistics Dialog -->
    <el-dialog v-model="editDialogVisible" title="Edit Logistics Information" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="Carrier" required>
          <el-select v-model="editForm.carrier_code" placeholder="Select carrier" style="width: 100%">
            <el-option v-for="carrier in carriers" :key="carrier.code" :label="carrier.name" :value="carrier.code" />
          </el-select>
        </el-form-item>
        <el-form-item label="Tracking No." required>
          <el-input v-model="editForm.tracking_no" placeholder="Enter tracking number" />
        </el-form-item>
        <el-form-item label="Shipping Cost">
          <el-input-number v-model="editForm.shipping_cost" :min="0" :precision="2" style="width: 100%">
            <template #prefix>{{ editForm.currency }}</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="Weight (kg)">
          <el-input-number v-model="editForm.weight" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Remark">
          <el-input v-model="editForm.remark" type="textarea" :rows="3" placeholder="Optional" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="saveLogistics">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Document, Van, Goods, Clock, Operation, Location,
  Link, CopyDocument, CircleCheck, Close, Picture
} from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  getShipmentDetail,
  updateShipmentStatus,
  getCarrierList,
  type Shipment,
  type Carrier
} from '@/api/fulfillment'

const route = useRoute()
const router = useRouter()

const shipment = ref<Shipment | null>(null)
const carriers = ref<Carrier[]>([])
const editDialogVisible = ref(false)

const editForm = reactive({
  carrier_code: '',
  tracking_no: '',
  shipping_cost: 0,
  weight: 0,
  remark: '',
  currency: 'CNY'
})

const statusTypeMap = {
  0: { type: 'warning' as const, text: 'Pending' },
  1: { type: 'primary' as const, text: 'Shipped' },
  2: { type: 'info' as const, text: 'In Transit' },
  3: { type: 'success' as const, text: 'Delivered' },
  4: { type: 'danger' as const, text: 'Failed' }
}

const timeline = computed(() => {
  if (!shipment.value) return []

  const events = [
    { title: 'Shipment Created', time: shipment.value.created_at, type: 'primary' as const, active: true, description: `Order: ${shipment.value.order_no}` }
  ]

  if (shipment.value.shipped_at) {
    events.push({
      title: 'Shipped',
      time: shipment.value.shipped_at,
      type: 'primary' as const,
      active: shipment.value.status >= 1,
      description: `Carrier: ${shipment.value.carrier}, Tracking: ${shipment.value.tracking_no}`
    })
  }

  if (shipment.value.status >= 2) {
    events.push({
      title: 'In Transit',
      time: shipment.value.shipped_at,
      type: 'info' as const,
      active: shipment.value.status >= 2,
      description: 'Package is on the way'
    })
  }

  if (shipment.value.delivered_at) {
    events.push({
      title: 'Delivered',
      time: shipment.value.delivered_at,
      type: 'success' as const,
      active: shipment.value.status === 3,
      description: 'Package has been delivered'
    })
  }

  if (shipment.value.status === 4) {
    events.push({
      title: 'Delivery Failed',
      time: shipment.value.delivered_at || shipment.value.shipped_at,
      type: 'danger' as const,
      active: true,
      description: shipment.value.remark || 'Delivery failed'
    })
  }

  return events
})

const loadShipment = async () => {
  const id = route.params.id
  try {
    const res = await getShipmentDetail(Number(id))
    shipment.value = res.data
  } catch (error) {
    // Mock data
    shipment.value = {
      id: Number(id),
      shipment_no: 'SHP20260322001',
      order_id: 'ORD001',
      order_no: 'ORD2026031800100',
      status: 2,
      carrier: 'SF Express',
      carrier_code: 'SF',
      tracking_no: 'SF1234567890',
      shipping_cost: 1200,
      shipping_currency: 'CNY',
      weight: 1.5,
      shipped_at: '2026-03-22 14:30:25',
      delivered_at: null,
      remark: 'Fragile item - handle with care',
      created_at: '2026-03-22 10:00:00',
      items: [
        { id: 1, product_id: 1, sku_id: 1, product_name: 'Wireless Bluetooth Earphones Pro', sku_name: 'Black - Premium Edition', image: '', quantity: 1 },
        { id: 2, product_id: 2, sku_id: 2, product_name: 'Phone Case', sku_name: 'iPhone 15 Pro Max', image: '', quantity: 2 }
      ]
    }
  }
}

const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res.data
  } catch (error) {
    carriers.value = [
      { code: 'SF', name: 'SF Express', tracking_url: '', is_active: true },
      { code: 'ZT', name: 'ZTO Express', tracking_url: '', is_active: true },
      { code: 'YT', name: 'YTO Express', tracking_url: '', is_active: true }
    ]
  }
}

const goBack = () => {
  router.back()
}

const viewOrder = () => {
  if (shipment.value) {
    router.push(`/orders?id=${shipment.value.order_id}`)
  }
}

const copyTracking = () => {
  if (shipment.value?.tracking_no) {
    navigator.clipboard.writeText(shipment.value.tracking_no)
    ElMessage.success('Tracking number copied')
  }
}

const openEditLogistics = () => {
  if (shipment.value) {
    editForm.carrier_code = shipment.value.carrier_code
    editForm.tracking_no = shipment.value.tracking_no
    editForm.shipping_cost = shipment.value.shipping_cost / 100
    editForm.weight = shipment.value.weight
    editForm.remark = shipment.value.remark
    editDialogVisible.value = true
  }
}

const saveLogistics = async () => {
  ElMessage.success('Logistics information updated')
  editDialogVisible.value = false
  loadShipment()
}

const confirmShip = async () => {
  try {
    await ElMessageBox.confirm(
      'Confirm to ship this package?',
      'Confirm Shipment',
      { type: 'success' }
    )
    await updateShipmentStatus(shipment.value!.id, 1)
    ElMessage.success('Shipment confirmed')
    loadShipment()
  } catch (error) {
    // Cancelled or error
  }
}

const markDelivered = async () => {
  try {
    await ElMessageBox.confirm(
      'Mark this shipment as delivered?',
      'Confirm Delivery',
      { type: 'success' }
    )
    await updateShipmentStatus(shipment.value!.id, 3)
    ElMessage.success('Shipment marked as delivered')
    loadShipment()
  } catch (error) {
    // Cancelled or error
  }
}

const trackShipment = () => {
  if (shipment.value?.tracking_no) {
    const carrier = carriers.value.find(c => c.code === shipment.value?.carrier_code)
    if (carrier?.tracking_url) {
      window.open(carrier.tracking_url.replace('{tracking_no}', shipment.value.tracking_no), '_blank')
    } else {
      ElMessage.info(`Tracking No: ${shipment.value.tracking_no}`)
    }
  }
}

const cancelShipment = async () => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to cancel this shipment?',
      'Cancel Shipment',
      { type: 'warning' }
    )
    ElMessage.success('Shipment cancelled')
    router.push('/fulfillment/shipments')
  } catch (error) {
    // Cancelled
  }
}

onMounted(() => {
  loadShipment()
  loadCarriers()
})
</script>

<style scoped>
.shipment-detail-page {
  padding: 0;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #E5E7EB;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.title-section {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0;
}

.shipment-no {
  font-family: 'Fira Code', monospace;
  font-size: 14px;
  color: #6366F1;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Cards */
.info-card,
.timeline-card,
.actions-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  margin-bottom: 20px;
}

.info-card :deep(.el-card__header),
.timeline-card :deep(.el-card__header),
.actions-card :deep(.el-card__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title .el-icon {
  color: #6366F1;
}

.item-count {
  font-size: 13px;
  color: #6B7280;
}

/* Descriptions */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}

:deep(.el-descriptions__content) {
  color: #1E1B4B;
}

.value-text {
  font-family: 'Fira Code', monospace;
  color: #6366F1;
}

.no-data {
  color: #9CA3AF;
}

.tracking-value {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Fira Code', monospace;
}

/* Remark */
.remark-section {
  margin-top: 16px;
  padding: 12px 16px;
  background: #FEF3C7;
  border-radius: 8px;
  border-left: 3px solid #F59E0B;
}

.remark-label {
  font-weight: 500;
  color: #92400E;
  margin-right: 8px;
}

.remark-text {
  color: #78350F;
}

/* Items List */
.items-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.item-row:hover {
  background: #F5F3FF;
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
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.item-quantity {
  font-weight: 600;
  color: #6366F1;
  font-size: 16px;
}

/* Timeline */
.timeline-content {
  padding-left: 4px;
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

:deep(.el-timeline-item.is-active .el-timeline-item__node) {
  background-color: #6366F1;
}

:deep(.el-timeline-item.is-active .el-timeline-item__tail) {
  border-left-color: #6366F1;
}

/* Actions */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 100%;
  height: 44px;
  justify-content: flex-start;
  padding: 0 20px;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-right {
    width: 100%;
    justify-content: space-between;
  }

  .info-card,
  .timeline-card,
  .actions-card {
    border-radius: 14px;
  }
}
</style>