<template>
  <el-dialog
    v-model="visible"
    :title="$t('orders.shipOrder')"
    width="600px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <!-- Order Info -->
      <div class="order-info-section">
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
            {{ order?.order_no }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('orders.items')">
            {{ $t('orders.itemsCount', { count: order?.items?.length || 0 }) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Logistics Info -->
      <div class="form-section">
        <div class="section-header">
          <el-icon><Van /></el-icon>
          <span>{{ $t('orders.logisticsInformation') }}</span>
        </div>

        <el-form-item
          :label="$t('orders.carrier')"
          prop="carrier_code"
        >
          <el-select
            v-model="form.carrier_code"
            :placeholder="$t('orders.selectCarrier')"
            style="width: 100%"
            @change="handleCarrierChange"
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
          v-if="form.carrier_code === 'OTHER'"
          :label="$t('orders.customCarrier')"
          prop="carrier_name"
        >
          <el-input
            v-model="form.carrier_name"
            :placeholder="$t('orders.enterCarrierName')"
          />
        </el-form-item>

        <el-form-item
          :label="$t('orders.trackingNo')"
          prop="tracking_no"
        >
          <el-input
            v-model="form.tracking_no"
            :placeholder="$t('orders.enterTrackingNumber')"
          >
            <template #prefix>
              <el-icon><Tickets /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('orders.shippingCost')">
              <el-input-number
                v-model="form.shipping_cost"
                :min="0"
                :precision="2"
                style="width: 100%"
                :placeholder="$t('orders.optional')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('orders.weight')">
              <el-input-number
                v-model="form.weight"
                :min="0"
                :precision="2"
                style="width: 100%"
                :placeholder="$t('orders.optional')"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- Items Selection -->
      <div class="items-section">
        <div class="section-header">
          <el-icon><Goods /></el-icon>
          <span>{{ $t('orders.itemsToShip') }}</span>
          <el-tag
            size="small"
            type="info"
          >
            {{ selectedItems.length }} {{ $t('orders.selected') }}
          </el-tag>
        </div>

        <div class="items-list">
          <div
            v-for="item in availableItems"
            :key="item.order_item_id"
            class="item-row"
            :class="{ selected: isItemSelected(item.order_item_id) }"
            @click="toggleItem(item)"
          >
            <el-checkbox :model-value="isItemSelected(item.order_item_id)" />
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
            <div class="item-quantity">
              <el-input-number
                v-if="isItemSelected(item.order_item_id)"
                :model-value="getItemQuantity(item.order_item_id)"
                :min="1"
                :max="item.pending_qty || item.quantity"
                size="small"
                style="width: 100px"
                @update:model-value="(val: number) => setItemQuantity(item.order_item_id, val)"
                @click.stop
              />
              <span
                v-else
                class="quantity-text"
              >x {{ item.quantity }}</span>
            </div>
          </div>
        </div>

        <div class="select-actions">
          <el-button
            size="small"
            @click="selectAllItems"
          >
            {{ $t('orders.selectAll') }}
          </el-button>
          <el-button
            size="small"
            @click="clearItems"
          >
            {{ $t('orders.clear') }}
          </el-button>
        </div>
      </div>

      <!-- Remark -->
      <el-form-item
        :label="$t('orders.remark')"
        class="remark-item"
      >
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="3"
          :placeholder="$t('orders.optionalNotes')"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          {{ $t('orders.confirmShipment') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Van, Goods, Tickets, Picture } from '@element-plus/icons-vue'
import { shipOrder, type Carrier } from '@/api/order'
import { t } from '@/plugins/i18n'

interface OrderItem {
  order_item_id: number
  product_name: string
  sku_name: string
  image: string
  quantity: number
  pending_qty?: number
}

interface Order {
  order_id: number
  order_no: string
  items: OrderItem[]
}

const props = defineProps<{
  modelValue: boolean
  order: Order | null
  carriers: Carrier[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref()
const submitting = ref(false)
const selectedItems = ref<number[]>([])

const form = reactive({
  carrier_code: '',
  carrier_name: '',
  tracking_no: '',
  shipping_cost: '0',
  weight: '0',
  remark: '',
  items: [] as { order_item_id: number; quantity: number }[]
})

const rules = {
  carrier_code: [{ required: true, message: t('orders.selectCarrier'), trigger: 'change' }],
  carrier_name: [{ required: true, message: t('orders.enterCarrierName'), trigger: 'blur' }],
  tracking_no: [{ required: true, message: t('orders.enterTrackingNumber'), trigger: 'blur' }]
}

const availableItems = computed(() => {
  return props.order?.items || []
})

const handleCarrierChange = (code: string) => {
  if (code !== 'OTHER') {
    form.carrier_name = ''
  }
}

const isItemSelected = (id: number) => {
  return selectedItems.value.includes(id)
}

const getItemQuantity = (id: number) => {
  const item = form.items.find(i => i.order_item_id === id)
  return item?.quantity || 1
}

const setItemQuantity = (id: number, quantity: number) => {
  const item = form.items.find(i => i.order_item_id === id)
  if (item) {
    item.quantity = quantity
  }
}

const toggleItem = (item: OrderItem) => {
  const index = selectedItems.value.indexOf(item.order_item_id)
  if (index === -1) {
    selectedItems.value.push(item.order_item_id)
    form.items.push({ order_item_id: item.order_item_id, quantity: item.pending_qty || item.quantity })
  } else {
    selectedItems.value.splice(index, 1)
    form.items = form.items.filter(i => i.order_item_id !== item.order_item_id)
  }
}

const selectAllItems = () => {
  selectedItems.value = availableItems.value.map(item => item.order_item_id)
  form.items = availableItems.value.map(item => ({
    order_item_id: item.order_item_id,
    quantity: item.pending_qty || item.quantity
  }))
}

const clearItems = () => {
  selectedItems.value = []
  form.items = []
}

const handleSubmit = async () => {
  if (!formRef.value) return
  const order = props.order
  if (!order) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    if (form.items.length === 0) {
      ElMessage.warning(t('orders.pleaseSelectItem'))
      return
    }

    submitting.value = true
    try {
      await shipOrder(order.order_id, {
        carrier_code: form.carrier_code,
        carrier_name: form.carrier_code === 'OTHER' ? form.carrier_name : undefined,
        tracking_no: form.tracking_no,
        shipping_cost: form.shipping_cost && form.shipping_cost !== '0' ? form.shipping_cost : undefined,
        weight: form.weight && form.weight !== '0' ? form.weight : undefined,
        remark: form.remark || undefined,
        items: form.items
      })
      ElMessage.success(t('orders.shipSuccess'))
      emit('success')
      visible.value = false
    } catch (error: unknown) {
      ElMessage.error((error as Error)?.message || t('orders.shipFailed'))
    } finally {
      submitting.value = false
    }
  })
}

// Reset form when dialog opens
watch(visible, (val) => {
  if (val && props.order) {
    form.carrier_code = ''
    form.carrier_name = ''
    form.tracking_no = ''
    form.shipping_cost = '0'
    form.weight = '0'
    form.remark = ''
    selectedItems.value = props.order.items.map(item => item.order_item_id)
    form.items = props.order.items.map(item => ({
      order_item_id: item.order_item_id,
      quantity: item.pending_qty || item.quantity
    }))
  }
})
</script>

<style scoped>
.order-info-section,
.form-section,
.items-section {
  margin-bottom: 24px;
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

/* Descriptions */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}

/* Items List */
.items-list {
  max-height: 280px;
  overflow-y: auto;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 8px;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.item-row:hover {
  background: #F5F3FF;
}

.item-row.selected {
  background: #EEF2FF;
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.item-image {
  width: 48px;
  height: 48px;
  border-radius: 8px;
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
  font-size: 14px;
}

.item-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
}

.item-quantity {
  flex-shrink: 0;
}

.quantity-text {
  color: #6366F1;
  font-weight: 600;
}

.select-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

/* Remark */
.remark-item {
  margin-top: 16px;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>