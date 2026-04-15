<template>
  <div class="shipment-detail-page">
    <!-- Skeleton Loading -->
    <div
      v-if="pageLoading"
      class="skeleton-container"
    >
      <el-skeleton
        :rows="1"
        animated
      />
      <el-row :gutter="20">
        <el-col
          :xs="24"
          :lg="16"
        >
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <el-skeleton
                :rows="1"
                animated
              />
            </template>
            <el-skeleton
              :rows="4"
              animated
            />
          </el-card>
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <el-skeleton
                :rows="1"
                animated
              />
            </template>
            <el-skeleton
              :rows="4"
              animated
            />
          </el-card>
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <el-skeleton
                :rows="1"
                animated
              />
            </template>
            <el-skeleton
              :rows="3"
              animated
            />
          </el-card>
        </el-col>
        <el-col
          :xs="24"
          :lg="8"
        >
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <el-skeleton
                :rows="1"
                animated
              />
            </template>
            <el-skeleton
              :rows="5"
              animated
            />
          </el-card>
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <el-skeleton
                :rows="1"
                animated
              />
            </template>
            <el-skeleton
              :rows="3"
              animated
            />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- Actual Content -->
    <template v-else>
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-left">
          <el-button
            link
            @click="goBack"
          >
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('common.back') }}
          </el-button>
          <div class="title-section">
            <h1 class="page-title">
              {{ $t('fulfillment.shipmentDetails') }}
            </h1>
            <p class="shipment-no">
              {{ shipment?.shipment_no }}
            </p>
          </div>
        </div>
        <div class="header-right">
          <status-tag
            :status="shipment?.status"
            :type-map="statusTypeMap"
            size="large"
          />
          <el-button
            v-if="shipment?.tracking_no"
            type="primary"
            @click="copyTracking"
          >
            <el-icon><Link /></el-icon>
            {{ $t('fulfillment.copyTracking') }}
          </el-button>
        </div>
      </div>

      <el-row :gutter="20">
        <!-- Left Column -->
        <el-col
          :xs="24"
          :lg="16"
        >
          <!-- Basic Info Card -->
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Document /></el-icon>
                  {{ $t('fulfillment.basicInformation') }}
                </span>
              </div>
            </template>
            <el-descriptions
              :column="2"
              border
            >
              <el-descriptions-item :label="$t('fulfillment.shipmentNo')">
                <span class="value-text">{{ shipment?.shipment_no }}</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.orderNo')">
                <el-button
                  type="primary"
                  link
                  @click="viewOrder"
                >
                  {{ shipment?.order_no }}
                </el-button>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('common.status')">
                <status-tag
                  :status="shipment?.status"
                  :type-map="statusTypeMap"
                />
              </el-descriptions-item>
              <el-descriptions-item :label="$t('common.createdAt')">
                {{ shipment?.created_at }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- Logistics Info Card -->
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Van /></el-icon>
                  {{ $t('fulfillment.logisticsInformation') }}
                </span>
                <el-button
                  v-if="shipment?.status === '0'"
                  type="primary"
                  size="small"
                  @click="openEditLogistics"
                >
                  {{ $t('fulfillment.editLogistics') }}
                </el-button>
              </div>
            </template>
            <el-descriptions
              :column="2"
              border
            >
              <el-descriptions-item :label="$t('fulfillment.carrier')">
                <span class="value-text">{{ shipment?.carrier || '-' }}</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.trackingNo')">
                <span
                  v-if="shipment?.tracking_no"
                  class="tracking-value"
                >
                  {{ shipment.tracking_no }}
                  <el-button
                    link
                    size="small"
                    @click="copyTracking"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
                <span
                  v-else
                  class="no-data"
                >{{ $t('fulfillment.notEntered') }}</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.shippingCost')">
                <span v-if="shipment?.shipping_cost">
                  {{ shipment.currency }} {{ shipment.shipping_cost }}
                </span>
                <span
                  v-else
                  class="no-data"
                >-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.weight')">
                <span v-if="shipment?.weight">{{ shipment.weight }} {{ $t('fulfillment.kg') }}</span>
                <span
                  v-else
                  class="no-data"
                >-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.shippedAt')">
                {{ shipment?.shipped_at || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="$t('fulfillment.deliveredAt')">
                {{ shipment?.delivered_at || '-' }}
              </el-descriptions-item>
            </el-descriptions>
            <div
              v-if="shipment?.remark"
              class="remark-section"
            >
              <span class="remark-label">{{ $t('fulfillment.remark') }}:</span>
              <span class="remark-text">{{ shipment.remark }}</span>
            </div>
          </el-card>

          <!-- Items Card -->
          <el-card
            class="info-card"
            shadow="never"
          >
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Goods /></el-icon>
                  {{ $t('fulfillment.shipmentItems') }}
                </span>
                <span class="item-count">{{ $t('fulfillment.totalItems', { n: shipment?.items?.length || 0 }) }}</span>
              </div>
            </template>
            <div class="items-list">
              <div
                v-for="item in shipment?.items"
                :key="item.id"
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
                    {{ $t('fulfillment.sku') }}: {{ item.sku_name }}
                  </p>
                </div>
                <div class="item-quantity">
                  x {{ item.quantity }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- Right Column -->
        <el-col
          :xs="24"
          :lg="8"
        >
          <!-- Timeline Card -->
          <el-card
            class="timeline-card"
            shadow="never"
          >
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Clock /></el-icon>
                  {{ $t('fulfillment.statusTimeline') }}
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
                  <p class="timeline-title">
                    {{ event.title }}
                  </p>
                  <p
                    v-if="event.description"
                    class="timeline-desc"
                  >
                    {{ event.description }}
                  </p>
                </div>
              </el-timeline-item>
            </el-timeline>
          </el-card>

          <!-- Actions Card -->
          <el-card
            class="actions-card"
            shadow="never"
          >
            <template #header>
              <div class="card-header">
                <span class="card-title">
                  <el-icon><Operation /></el-icon>
                  {{ $t('fulfillment.actions') }}
                </span>
              </div>
            </template>
            <div class="action-buttons">
              <el-button
                v-if="shipment?.status === '0'"
                type="success"
                class="action-btn"
                @click="confirmShip"
              >
                <el-icon><Van /></el-icon>
                {{ $t('fulfillment.confirmShip') }}
              </el-button>
              <el-button
                v-if="shipment?.status === '2'"
                type="primary"
                class="action-btn"
                @click="markDelivered"
              >
                <el-icon><CircleCheck /></el-icon>
                {{ $t('fulfillment.markAsDelivered') }}
              </el-button>
              <el-button
                v-if="shipment?.tracking_no"
                class="action-btn"
                @click="trackShipment"
              >
                <el-icon><Location /></el-icon>
                {{ $t('fulfillment.trackPackage') }}
              </el-button>
              <el-button
                v-if="shipment?.status === '0'"
                type="danger"
                class="action-btn"
                @click="cancelShipment"
              >
                <el-icon><Close /></el-icon>
                {{ $t('fulfillment.cancelShipment') }}
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- Edit Logistics Dialog -->
      <el-dialog
        v-model="editDialogVisible"
        :title="$t('fulfillment.editLogisticsInformation')"
        width="500px"
      >
        <el-form
          :model="editForm"
          label-width="100px"
        >
          <el-form-item
            :label="$t('fulfillment.carrier')"
            required
          >
            <el-select
              v-model="editForm.carrier_code"
              :placeholder="$t('fulfillment.selectCarrier')"
              style="width: 100%"
            >
              <el-option
                v-for="carrier in carriers"
                :key="carrier.code"
                :label="carrier.name"
                :value="carrier.code"
              />
            </el-select>
          </el-form-item>
          <el-form-item
            :label="$t('fulfillment.trackingNo')"
            required
          >
            <el-input
              v-model="editForm.tracking_no"
              :placeholder="$t('fulfillment.enterTrackingNumber')"
            />
          </el-form-item>
          <el-form-item :label="$t('fulfillment.shippingCost')">
            <el-input-number
              v-model="editForm.shipping_cost"
              :min="0"
              :precision="2"
              style="width: 100%"
            >
              <template #prefix>
                {{ editForm.currency }}
              </template>
            </el-input-number>
          </el-form-item>
          <el-form-item :label="$t('fulfillment.weight')">
            <el-input-number
              v-model="editForm.weight"
              :min="0"
              :precision="2"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('fulfillment.remark')">
            <el-input
              v-model="editForm.remark"
              type="textarea"
              :rows="3"
              :placeholder="$t('fulfillment.optional')"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="editDialogVisible = false">
            {{ $t('common.cancel') }}
          </el-button>
          <el-button
            type="primary"
            @click="saveLogistics"
          >
            {{ $t('common.save') }}
          </el-button>
        </template>
      </el-dialog>
    </template>
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
import { t } from '@/plugins/i18n'
import {
  getShipmentDetail,
  updateShipment,
  updateShipmentStatus,
  getCarrierList,
  cancelShipment as cancelShipmentApi,
  type Shipment,
  type Carrier,
  type UpdateShipmentRequest
} from '@/api/fulfillment'

interface TimelineEvent {
  title: string
  time: string | null
  type: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  active: boolean
  description?: string
}

const route = useRoute()
const router = useRouter()

const shipment = ref<Shipment | null>(null)
const carriers = ref<Carrier[]>([])
const editDialogVisible = ref(false)
const pageLoading = ref(true)

const editForm = reactive({
  carrier_code: '',
  tracking_no: '',
  shipping_cost: '0',
  weight: '0',
  remark: '',
  currency: 'CNY'
})

const statusTypeMap: Record<string, { type: 'warning' | 'primary' | 'info' | 'success' | 'danger', text: string }> = {
  '0': { type: 'warning', text: 'Pending' },
  '1': { type: 'primary', text: 'Shipped' },
  '2': { type: 'info', text: 'In Transit' },
  '3': { type: 'success', text: 'Delivered' },
  '4': { type: 'danger', text: 'Failed' }
}

const timeline = computed<TimelineEvent[]>(() => {
  if (!shipment.value) return []

  const events: TimelineEvent[] = [
    { title: t('fulfillment.shipmentCreated'), time: shipment.value.created_at, type: 'primary', active: true, description: `${t('fulfillment.orderNo')}: ${shipment.value.order_no}` }
  ]

  if (shipment.value.shipped_at) {
    events.push({
      title: t('fulfillment.shippedStatus'),
      time: shipment.value.shipped_at,
      type: 'primary',
      active: ['1', '2', '3', '4'].includes(shipment.value.status),
      description: `${t('fulfillment.carrier')}: ${shipment.value.carrier}, ${t('fulfillment.trackingNo')}: ${shipment.value.tracking_no}`
    })
  }

  if (['1', '2', '3'].includes(shipment.value.status) && shipment.value.shipped_at) {
    events.push({
      title: t('fulfillment.inTransitStatus'),
      time: shipment.value.shipped_at,
      type: 'info',
      active: ['2', '3'].includes(shipment.value.status),
      description: t('fulfillment.packageOnTheWay')
    })
  }

  if (shipment.value.delivered_at) {
    events.push({
      title: t('fulfillment.deliveredStatus'),
      time: shipment.value.delivered_at,
      type: 'success',
      active: shipment.value.status === '3',
      description: t('fulfillment.packageDelivered')
    })
  }

  if (shipment.value.status === '4') {
    events.push({
      title: t('fulfillment.deliveryFailed'),
      time: shipment.value.delivered_at || shipment.value.shipped_at || '',
      type: 'danger',
      active: true,
      description: shipment.value.remark || t('fulfillment.deliveryFailed')
    })
  }

  return events
})

const loadShipment = async () => {
  const id = route.params.id
  try {
    const res = await getShipmentDetail(Number(id))
    shipment.value = res
  } catch (error) {
    ElMessage.error(t('fulfillment.loadShipmentDetailsFailed'))
  } finally {
    pageLoading.value = false
  }
}

const loadCarriers = async () => {
  try {
    const res = await getCarrierList()
    carriers.value = res.list
  } catch (error) {
    ElMessage.error(t('fulfillment.loadCarriersFailed'))
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
    ElMessage.success(t('fulfillment.trackingNumberCopied'))
  }
}

const openEditLogistics = () => {
  if (shipment.value) {
    editForm.carrier_code = shipment.value.carrier_code
    editForm.tracking_no = shipment.value.tracking_no
    editForm.shipping_cost = shipment.value.shipping_cost
    editForm.weight = shipment.value.weight
    editForm.remark = shipment.value.remark || ''
    editDialogVisible.value = true
  }
}

const saveLogistics = async () => {
  if (!shipment.value) return
  try {
    const data: UpdateShipmentRequest = {
      carrier_code: editForm.carrier_code,
      tracking_no: editForm.tracking_no,
      shipping_cost: editForm.shipping_cost,
      weight: editForm.weight,
      remark: editForm.remark
    }
    await updateShipment(shipment.value.id, data)
    ElMessage.success(t('fulfillment.logisticsUpdated'))
    editDialogVisible.value = false
    loadShipment()
  } catch (error) {
    ElMessage.error(t('fulfillment.updateShipmentFailed'))
  }
}

const confirmShip = async () => {
  try {
    await ElMessageBox.confirm(
      t('fulfillment.confirmToShip'),
      t('fulfillment.confirmShip'),
      { type: 'success' }
    )
    await updateShipmentStatus(shipment.value!.id, '1')
    ElMessage.success(t('fulfillment.shipmentConfirmed'))
    loadShipment()
  } catch (error) {
    // Cancelled or error
  }
}

const markDelivered = async () => {
  try {
    await ElMessageBox.confirm(
      t('fulfillment.markDelivered'),
      t('fulfillment.confirmDelivery'),
      { type: 'success' }
    )
    await updateShipmentStatus(shipment.value!.id, '3')
    ElMessage.success(t('fulfillment.shipmentMarkedDelivered'))
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
      ElMessage.info(`${t('fulfillment.trackingNo')}: ${shipment.value.tracking_no}`)
    }
  }
}

const cancelShipment = async () => {
  try {
    const { value: reason } = await ElMessageBox.prompt(
      t('fulfillment.cancelReasonRequired'),
      t('fulfillment.cancelShipment'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        inputPattern: /^.{5,200}$/,
        inputErrorMessage: t('fulfillment.cancelReasonLengthError')
      }
    )

    if (shipment.value?.id) {
      await cancelShipmentApi(shipment.value.id, reason)
      ElMessage.success(t('fulfillment.shipmentCancelled'))
      router.push('/fulfillment/shipments')
    }
  } catch (error) {
    // User cancelled or validation failed
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

/* Skeleton Container */
.skeleton-container {
  padding: 0;
}

.skeleton-container .info-card {
  margin-bottom: 20px;
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
