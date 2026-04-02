<template>
  <el-card class="table-card" shadow="never">
    <EmptyState
      v-if="orderList.length === 0 && !loading"
      :title="$t('orders.noOrders')"
      :description="$t('orders.noOrdersDesc')"
    />
    <Table
      v-else
      ref="tableRef"
      :data="orderList"
      :loading="loading"
      @selection-change="handleSelectionChange"
    >
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
                <p><strong>{{ $t('orders.paymentMethod') }}:</strong> {{ row.payment_method_text || '-' }}</p>
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
            <el-tag
              v-if="row.refund_status !== 'none' && row.refund_status !== undefined"
              size="small"
              type="warning"
              effect="dark"
            >
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
            v-if="row.fulfillment_status !== undefined && row.fulfillment_status !== 0"
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
                <el-dropdown-item v-if="row.status === 'pending_payment'" command="adjust">
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
    </Table>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Picture, Edit, ArrowDown, PriceTag, Close } from '@element-plus/icons-vue'
import { getOrderStatusType } from '@/utils/status'
import type { Order, FulfillmentStatus } from '@/api/order'
import Table from '@/components/common/Table.vue'
import EmptyState from '@/components/common/EmptyState.vue'

defineProps<{
  orderList: Order[]
  loading: boolean
}>()

const emit = defineEmits<{
  'selection-change': [selection: Order[]]
  detail: [order: Order]
  ship: [order: Order]
  remind: [order: Order]
  command: [command: string, order: Order]
  'clear-selection': []
}>()

const tableRef = ref()

const formatAmount = (amount: string | undefined) => {
  if (!amount) return '0.00'
  return parseFloat(amount).toFixed(2)
}

const formatTime = (time: string | undefined | null) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

const getStatusType = getOrderStatusType

const getFulfillmentType = (status: FulfillmentStatus) => {
  const types: Record<FulfillmentStatus, string> = {
    0: 'warning',
    1: 'primary',
    2: 'info',
    3: 'success'
  }
  return types[status] || 'info'
}

const canShip = (order: Order) => {
  return order.status === 'paid' &&
    (order.fulfillment_status === 0 || order.fulfillment_status === 1)
}

const canCancel = (order: Order) => {
  return ['pending_payment', 'paid'].includes(order.status)
}

const handleSelectionChange = (selection: Order[]) => {
  emit('selection-change', selection)
}

const handleDetail = (row: Order) => {
  emit('detail', row)
}

const handleShip = (row: Order) => {
  emit('ship', row)
}

const handleRemind = (row: Order) => {
  emit('remind', row)
}

const handleCommand = (command: string, row: Order) => {
  emit('command', command, row)
}

// Expose for clearSelection
defineExpose({
  clearSelection: () => tableRef.value?.clearSelection()
})
</script>

<style scoped>
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

@media (max-width: 768px) {
  .goods-preview {
    flex-wrap: wrap;
  }

  .table-card {
    border-radius: 14px;
  }
}
</style>
