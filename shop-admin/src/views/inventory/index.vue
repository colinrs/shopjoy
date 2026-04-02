<template>
  <div class="inventory-page">
    <!-- Low Stock Alert Card -->
    <el-card
      v-if="lowStockList.length > 0"
      class="alert-card"
      shadow="never"
    >
      <template #header>
        <div class="card-header">
          <span><el-icon><Warning /></el-icon> {{ $t('inventory.lowStockAlert') }}</span>
          <el-badge
            :value="lowStockTotal"
            type="danger"
          />
        </div>
      </template>
      <el-table
        :data="lowStockList"
        stripe
        size="small"
      >
        <el-table-column
          prop="sku_code"
          :label="$t('inventory.skuCode')"
          width="150"
        />
        <el-table-column
          prop="product_name"
          :label="$t('inventory.productName')"
          min-width="200"
        />
        <el-table-column
          prop="available_stock"
          :label="$t('inventory.availableStock')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              type="danger"
              size="small"
            >
              {{ row.available_stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="safety_stock"
          :label="$t('inventory.safetyStock')"
          width="100"
          align="center"
        />
        <el-table-column
          :label="$t('inventory.actions')"
          width="120"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleAdjustStock(row)"
            >
              {{ $t('inventory.replenish') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Main Content Tabs -->
    <el-card
      class="main-card"
      shadow="never"
    >
      <el-tabs v-model="activeTab">
        <!-- Inventory Logs Tab -->
        <el-tab-pane
          :label="$t('inventory.logs')"
          name="logs"
        >
          <div class="filter-bar">
            <el-input
              v-model="logFilter.sku_code"
              :placeholder="$t('inventory.skuCode')"
              clearable
              style="width: 180px"
            />
            <el-select
              v-model="logFilter.type"
              :placeholder="$t('inventory.changeType')"
              clearable
              style="width: 140px"
            >
              <el-option
                :label="$t('inventory.manual')"
                value="manual"
              />
              <el-option
                :label="$t('inventory.order')"
                value="order"
              />
              <el-option
                :label="$t('inventory.return')"
                value="return"
              />
              <el-option
                :label="$t('inventory.adjustment')"
                value="adjustment"
              />
            </el-select>
            <el-button
              type="primary"
              @click="loadInventoryLogs"
            >
              {{ $t('inventory.search') }}
            </el-button>
            <el-button
              :loading="logExporting"
              @click="handleExportInventoryLogs"
            >
              <el-icon v-if="!logExporting">
                <Download />
              </el-icon>
              {{ $t('inventory.exportLogs') }}
            </el-button>
          </div>
          <el-table
            v-loading="logLoading"
            :data="logList"
            stripe
          >
            <el-table-column
              prop="sku_code"
              :label="$t('inventory.skuCode')"
              width="150"
            />
            <el-table-column
              prop="change_type"
              :label="$t('inventory.changeType')"
              width="100"
            >
              <template #default="{ row }">
                <el-tag
                  size="small"
                  :type="getChangeTypeTag(row.change_type)"
                >
                  {{ getChangeTypeLabel(row.change_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="change_quantity"
              :label="$t('inventory.changeQuantity')"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <span :class="row.change_quantity >= 0 ? 'text-success' : 'text-danger'">
                  {{ row.change_quantity >= 0 ? '+' : '' }}{{ row.change_quantity }}
                </span>
              </template>
            </el-table-column>
            <el-table-column
              prop="before_stock"
              :label="$t('inventory.beforeStock')"
              width="100"
              align="center"
            />
            <el-table-column
              prop="after_stock"
              :label="$t('inventory.afterStock')"
              width="100"
              align="center"
            />
            <el-table-column
              prop="order_no"
              :label="$t('inventory.relatedOrder')"
              width="150"
            />
            <el-table-column
              prop="remark"
              :label="$t('inventory.remark')"
              min-width="150"
            />
            <el-table-column
              prop="created_at"
              :label="$t('inventory.time')"
              width="180"
            />
          </el-table>
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="logPage"
              v-model:page-size="logPageSize"
              :total="logTotal"
              :page-sizes="[20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @size-change="loadInventoryLogs"
              @current-change="loadInventoryLogs"
            />
          </div>
        </el-tab-pane>

        <!-- Warehouses Tab -->
        <el-tab-pane
          :label="$t('inventory.warehouseManagement')"
          name="warehouses"
        >
          <div class="filter-bar">
            <el-button
              type="primary"
              @click="handleAddWarehouse"
            >
              <el-icon><Plus /></el-icon>{{ $t('inventory.addWarehouse') }}
            </el-button>
          </div>
          <el-table
            v-loading="warehouseLoading"
            :data="warehouseList"
            stripe
          >
            <el-table-column
              prop="code"
              :label="$t('inventory.warehouseCode')"
              width="120"
            />
            <el-table-column
              prop="name"
              :label="$t('inventory.warehouseName')"
              min-width="150"
            />
            <el-table-column
              prop="country"
              :label="$t('inventory.country')"
              width="100"
            />
            <el-table-column
              prop="address"
              :label="$t('inventory.address')"
              min-width="200"
            />
            <el-table-column
              prop="is_default"
              :label="$t('inventory.defaultWarehouse')"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <el-tag
                  v-if="row.is_default"
                  type="success"
                  size="small"
                >
                  {{ $t('inventory.isDefault') }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              prop="status"
              :label="$t('inventory.status')"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="(val: number) => handleWarehouseStatusChange(row, val)"
                />
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('inventory.actions')"
              width="200"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="handleEditWarehouse(row)"
                >
                  {{ $t('inventory.edit') }}
                </el-button>
                <el-button
                  v-if="!row.is_default"
                  type="success"
                  link
                  size="small"
                  @click="handleSetDefaultWarehouse(row)"
                >
                  {{ $t('inventory.setAsDefault') }}
                </el-button>
                <el-button
                  v-if="!row.is_default"
                  type="danger"
                  link
                  size="small"
                  @click="handleDeleteWarehouse(row)"
                >
                  {{ $t('inventory.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Stock Adjustment Tab -->
        <el-tab-pane
          :label="$t('inventory.stockAdjustment')"
          name="adjust"
        >
          <el-form
            :model="adjustForm"
            label-width="120px"
            class="adjust-form"
          >
            <el-form-item
              :label="$t('inventory.skuCode')"
              required
            >
              <el-input
                v-model="adjustForm.sku_code"
                :placeholder="$t('inventory.skuCode')"
                style="width: 300px"
              />
            </el-form-item>
            <el-form-item :label="$t('inventory.selectWarehouse')">
              <el-select
                v-model="adjustForm.warehouse_id"
                :placeholder="$t('inventory.selectWarehouse')"
                style="width: 300px"
              >
                <el-option
                  v-for="w in warehouseList"
                  :key="w.id"
                  :label="w.name"
                  :value="w.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item
              :label="$t('inventory.adjustQuantity')"
              required
            >
              <el-input-number
                v-model="adjustForm.quantity"
                :min="-9999"
                :max="9999"
                style="width: 200px"
              />
              <span class="form-tip">{{ $t('inventory.quantityPositiveIncrease') }}</span>
            </el-form-item>
            <el-form-item :label="$t('inventory.remark')">
              <el-input
                v-model="adjustForm.remark"
                type="textarea"
                rows="3"
                :placeholder="$t('inventory.enterRemark')"
                style="width: 400px"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="adjustLoading"
                @click="handleAdjustSubmit"
              >
                {{ $t('inventory.confirmAdjustment') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Batch Safety Stock Tab -->
        <el-tab-pane
          :label="$t('inventory.batchSafetyStock')"
          name="batchSafetyStock"
        >
          <div class="filter-bar">
            <el-button
              type="primary"
              @click="showBatchSafetyStockDialog"
            >
              <el-icon><Plus /></el-icon>{{ $t('inventory.batchUpdateSafetyStock') }}
            </el-button>
          </div>
          <el-table
            :data="lowStockList"
            stripe
            size="small"
          >
            <el-table-column
              prop="sku_code"
              :label="$t('inventory.skuCode')"
              width="150"
            />
            <el-table-column
              prop="product_name"
              :label="$t('inventory.productName')"
              min-width="200"
            />
            <el-table-column
              prop="safety_stock"
              :label="$t('inventory.safetyStock')"
              width="100"
              align="center"
            />
            <el-table-column
              :label="$t('inventory.actions')"
              width="120"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="handleEditSafetyStock(row)"
                >
                  {{ $t('inventory.edit') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Warehouse Dialog -->
    <el-dialog
      v-model="warehouseDialogVisible"
      :title="isEditWarehouse ? $t('inventory.editWarehouse') : $t('inventory.newWarehouse')"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="warehouseFormRef"
        :model="warehouseForm"
        label-width="100px"
        :rules="warehouseRules"
      >
        <el-form-item
          v-if="!isEditWarehouse"
          :label="$t('inventory.warehouseCode')"
          prop="code"
        >
          <el-input
            v-model="warehouseForm.code"
            :placeholder="$t('inventory.warehouseCode')"
          />
        </el-form-item>
        <el-form-item
          :label="$t('inventory.warehouseName')"
          prop="name"
        >
          <el-input
            v-model="warehouseForm.name"
            :placeholder="$t('inventory.enterWarehouseName')"
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.country')">
          <el-input
            v-model="warehouseForm.country"
            :placeholder="$t('inventory.countryPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.address')">
          <el-input
            v-model="warehouseForm.address"
            type="textarea"
            rows="2"
            :placeholder="$t('inventory.addressPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.defaultWarehouse')">
          <el-switch v-model="warehouseForm.is_default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="warehouseDialogVisible = false">
          {{ $t('inventory.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="warehouseSaveLoading"
          @click="handleSaveWarehouse"
        >
          {{ $t('inventory.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Stock Adjust Dialog (from low stock alert) -->
    <el-dialog
      v-model="adjustDialogVisible"
      :title="$t('inventory.stockReplenishment')"
      width="400px"
      destroy-on-close
    >
      <el-form
        :model="quickAdjustForm"
        label-width="100px"
      >
        <el-form-item :label="$t('inventory.skuCode')">
          <el-input
            :value="quickAdjustForm.sku_code"
            disabled
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.replenishQuantity')">
          <el-input-number
            v-model="quickAdjustForm.quantity"
            :min="1"
            :max="9999"
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.remark')">
          <el-input
            v-model="quickAdjustForm.remark"
            :placeholder="$t('inventory.replenishRemark')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustDialogVisible = false">
          {{ $t('inventory.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="quickAdjustLoading"
          @click="handleQuickAdjust"
        >
          {{ $t('inventory.confirmReplenish') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Batch Safety Stock Dialog -->
    <el-dialog
      v-model="batchSafetyStockDialogVisible"
      :title="$t('inventory.batchUpdateSafetyStock')"
      width="700px"
      destroy-on-close
    >
      <div class="batch-safety-stock-form">
        <el-form
          :model="batchSafetyStockForm"
          label-width="120px"
        >
          <el-form-item :label="$t('inventory.selectSKU')">
            <el-select
              v-model="batchSafetyStockForm.sku_codes"
              multiple
              filterable
              allow-create
              default-first-option
              :placeholder="$t('inventory.enterSKUCode')"
              style="width: 100%"
            >
              <el-option
                v-for="item in lowStockList"
                :key="item.sku_code"
                :label="item.sku_code"
                :value="item.sku_code"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('inventory.safetyStockValue')">
            <el-input-number
              v-model="batchSafetyStockForm.safety_stock"
              :min="0"
              :max="999999"
              style="width: 200px"
            />
          </el-form-item>
        </el-form>

        <div
          v-if="batchSafetyStockForm.sku_codes.length > 0"
          class="selected-items"
        >
          <h4>{{ $t('inventory.selectedItems') }} ({{ batchSafetyStockForm.sku_codes.length }})</h4>
          <el-table
            :data="getSelectedSkuDetails()"
            stripe
            size="small"
            max-height="200"
          >
            <el-table-column
              prop="sku_code"
              :label="$t('inventory.skuCode')"
              width="150"
            />
            <el-table-column
              prop="product_name"
              :label="$t('inventory.productName')"
              min-width="200"
            />
            <el-table-column
              prop="safety_stock"
              :label="$t('inventory.currentSafetyStock')"
              width="120"
              align="center"
            />
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="batchSafetyStockDialogVisible = false">
          {{ $t('inventory.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="batchSafetyStockLoading"
          @click="handleBatchUpdateSafetyStock"
        >
          {{ $t('inventory.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Edit Safety Stock Dialog -->
    <el-dialog
      v-model="editSafetyStockDialogVisible"
      :title="$t('inventory.safetyStock')"
      width="400px"
      destroy-on-close
    >
      <el-form
        :model="editSafetyStockForm"
        label-width="120px"
      >
        <el-form-item :label="$t('inventory.skuCode')">
          <el-input
            :value="editSafetyStockForm.sku_code"
            disabled
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.productName')">
          <el-input
            :value="editSafetyStockForm.product_name"
            disabled
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.currentSafetyStock')">
          <el-input
            :value="editSafetyStockForm.current_safety_stock"
            disabled
          />
        </el-form-item>
        <el-form-item :label="$t('inventory.newSafetyStock')">
          <el-input-number
            v-model="editSafetyStockForm.safety_stock"
            :min="0"
            :max="999999"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editSafetyStockDialogVisible = false">
          {{ $t('inventory.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="editSafetyStockLoading"
          @click="handleUpdateSafetyStock"
        >
          {{ $t('inventory.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning, Download } from '@element-plus/icons-vue'
import {
  getWarehouses,
  createWarehouse,
  updateWarehouse,
  updateWarehouseStatus,
  setDefaultWarehouse,
  deleteWarehouse,
  getInventoryLogs,
  exportInventoryLogsUrl,
  getLowStockSKUs,
  adjustStock,
  batchUpdateSafetyStock,
  type Warehouse,
  type InventoryLog,
  type LowStockSKU
} from '@/api/inventory'
import { downloadFile } from '@/utils/download'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

// Low stock alerts
const lowStockList = ref<LowStockSKU[]>([])
const lowStockTotal = ref(0)

// Inventory logs
const logList = ref<InventoryLog[]>([])
const logLoading = ref(false)
const logExporting = ref(false)
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const logFilter = reactive({
  sku_code: '',
  type: ''
})

// Warehouses
const warehouseList = ref<Warehouse[]>([])
const warehouseLoading = ref(false)
const warehouseDialogVisible = ref(false)
const isEditWarehouse = ref(false)
const warehouseSaveLoading = ref(false)
const warehouseFormRef = ref()
const warehouseForm = reactive({
  id: 0,
  code: '',
  name: '',
  country: '',
  address: '',
  is_default: false
})
const warehouseRules = {
  code: [{ required: true, message: t('inventory.enterWarehouseCode'), trigger: 'blur' }],
  name: [{ required: true, message: t('inventory.enterWarehouseName'), trigger: 'blur' }]
}

// Stock adjustment
const activeTab = ref('logs')
const adjustForm = reactive({
  sku_code: '',
  warehouse_id: 0,
  quantity: 0,
  remark: ''
})
const adjustLoading = ref(false)

// Quick adjust dialog
const adjustDialogVisible = ref(false)
const quickAdjustLoading = ref(false)
const quickAdjustForm = reactive({
  sku_code: '',
  quantity: 100,
  remark: t('inventory.lowStockReplenish')
})

// Batch safety stock dialog
const batchSafetyStockDialogVisible = ref(false)
const batchSafetyStockLoading = ref(false)
const batchSafetyStockForm = reactive({
  sku_codes: [] as string[],
  safety_stock: 0
})

// Edit safety stock dialog
const editSafetyStockDialogVisible = ref(false)
const editSafetyStockLoading = ref(false)
const editSafetyStockForm = reactive({
  sku_code: '',
  product_name: '',
  current_safety_stock: 0,
  safety_stock: 0
})

// Show batch safety stock dialog
const showBatchSafetyStockDialog = () => {
  batchSafetyStockForm.sku_codes = []
  batchSafetyStockForm.safety_stock = 0
  batchSafetyStockDialogVisible.value = true
}

// Get selected SKU details
const getSelectedSkuDetails = () => {
  return lowStockList.value.filter(item => batchSafetyStockForm.sku_codes.includes(item.sku_code))
}

// Handle batch update safety stock
const handleBatchUpdateSafetyStock = async () => {
  if (batchSafetyStockForm.sku_codes.length === 0) {
    ElMessage.warning(t('inventory.pleaseSelectSKU'))
    return
  }
  if (batchSafetyStockForm.safety_stock < 0) {
    ElMessage.warning(t('inventory.pleaseEnterValidStock'))
    return
  }

  batchSafetyStockLoading.value = true
  try {
    const items = batchSafetyStockForm.sku_codes.map(sku_code => ({
      sku_code,
      safety_stock: batchSafetyStockForm.safety_stock
    }))
    await batchUpdateSafetyStock(items)
    ElMessage.success(t('inventory.batchUpdateSuccess'))
    batchSafetyStockDialogVisible.value = false
    loadLowStockAlerts()
  } catch (error) {
    handleError(error, t('inventory.batchUpdateFailed'))
  } finally {
    batchSafetyStockLoading.value = false
  }
}

// Handle edit safety stock
const handleEditSafetyStock = (row: LowStockSKU) => {
  editSafetyStockForm.sku_code = row.sku_code
  editSafetyStockForm.product_name = row.product_name
  editSafetyStockForm.current_safety_stock = row.safety_stock
  editSafetyStockForm.safety_stock = row.safety_stock
  editSafetyStockDialogVisible.value = true
}

// Handle single update safety stock
const handleUpdateSafetyStock = async () => {
  if (editSafetyStockForm.safety_stock < 0) {
    ElMessage.warning(t('inventory.pleaseEnterValidStock'))
    return
  }

  editSafetyStockLoading.value = true
  try {
    await batchUpdateSafetyStock([{
      sku_code: editSafetyStockForm.sku_code,
      safety_stock: editSafetyStockForm.safety_stock
    }])
    ElMessage.success(t('inventory.batchUpdateSuccess'))
    editSafetyStockDialogVisible.value = false
    loadLowStockAlerts()
  } catch (error) {
    handleError(error, t('inventory.batchUpdateFailed'))
  } finally {
    editSafetyStockLoading.value = false
  }
}

// Load low stock alerts
const loadLowStockAlerts = async () => {
  try {
    const res = await getLowStockSKUs({ page: 1, page_size: 10 })
    lowStockList.value = res.list || []
    lowStockTotal.value = res.total || 0
  } catch (error) {
    handleError(error, t('inventory.loadLowStockFailed'))
  }
}

// Load inventory logs
const loadInventoryLogs = async () => {
  logLoading.value = true
  try {
    const res = await getInventoryLogs({
      page: logPage.value,
      page_size: logPageSize.value,
      sku_code: logFilter.sku_code || undefined,
      type: logFilter.type || undefined
    })
    logList.value = res.list || []
    logTotal.value = res.total || 0
  } catch (error) {
    handleError(error, t('inventory.loadLogsFailed'))
  } finally {
    logLoading.value = false
  }
}

// Export inventory logs
const handleExportInventoryLogs = async () => {
  logExporting.value = true
  try {
    const { url, params } = exportInventoryLogsUrl({
      sku_code: logFilter.sku_code || undefined,
      type: logFilter.type || undefined
    })
    await downloadFile(url, params)
  } catch (error) {
    handleError(error)
  } finally {
    logExporting.value = false
  }
}

// Load warehouses
const loadWarehouses = async () => {
  warehouseLoading.value = true
  try {
    const res = await getWarehouses()
    warehouseList.value = res || []
  } catch (error) {
    handleError(error, t('inventory.loadWarehouseFailed'))
  } finally {
    warehouseLoading.value = false
  }
}

// Warehouse handlers
const handleAddWarehouse = () => {
  isEditWarehouse.value = false
  Object.assign(warehouseForm, {
    id: 0,
    code: '',
    name: '',
    country: '',
    address: '',
    is_default: false
  })
  warehouseDialogVisible.value = true
}

const handleEditWarehouse = (row: Warehouse) => {
  isEditWarehouse.value = true
  Object.assign(warehouseForm, {
    id: row.id,
    code: row.code,
    name: row.name,
    country: row.country || '',
    address: row.address || '',
    is_default: row.is_default
  })
  warehouseDialogVisible.value = true
}

const handleSaveWarehouse = async () => {
  if (!warehouseFormRef.value) return

  await warehouseFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      warehouseSaveLoading.value = true
      try {
        if (isEditWarehouse.value) {
          await updateWarehouse({
            id: warehouseForm.id,
            name: warehouseForm.name,
            country: warehouseForm.country,
            address: warehouseForm.address,
            is_default: warehouseForm.is_default
          })
          ElMessage.success(t('inventory.updateSuccess'))
        } else {
          await createWarehouse({
            code: warehouseForm.code,
            name: warehouseForm.name,
            country: warehouseForm.country,
            address: warehouseForm.address,
            is_default: warehouseForm.is_default
          })
          ElMessage.success(t('inventory.createSuccess'))
        }
        warehouseDialogVisible.value = false
        loadWarehouses()
      } catch (error) {
        handleError(error, isEditWarehouse.value ? t('inventory.updateFailed') : t('inventory.createFailed'))
      } finally {
        warehouseSaveLoading.value = false
      }
    }
  })
}

const handleWarehouseStatusChange = async (row: Warehouse, status: number) => {
  try {
    await updateWarehouseStatus(row.id, status)
    ElMessage.success(status === 1 ? t('inventory.enabledSuccess') : t('inventory.disabledSuccess'))
  } catch (error) {
    handleError(error, t('inventory.updateStatusFailed'))
    row.status = status === 1 ? 0 : 1
  }
}

const handleSetDefaultWarehouse = async (row: Warehouse) => {
  try {
    await setDefaultWarehouse(row.id)
    ElMessage.success(t('inventory.setDefaultSuccess'))
    loadWarehouses()
  } catch (error) {
    handleError(error, t('inventory.operationFailed'))
  }
}

const handleDeleteWarehouse = (row: Warehouse) => {
  ElMessageBox.confirm(t('inventory.confirmDeleteWarehouse', { name: row.name }), t('inventory.warning'), {
    confirmButtonText: t('inventory.confirm'),
    cancelButtonText: t('inventory.cancel'),
    type: 'warning'
  }).then(async () => {
    try {
      await deleteWarehouse(row.id)
      ElMessage.success(t('inventory.deleteSuccess'))
      loadWarehouses()
    } catch (error) {
      handleError(error, t('inventory.deleteFailed'))
    }
  })
}

// Stock adjustment handlers
const handleAdjustSubmit = async () => {
  if (!adjustForm.sku_code || !adjustForm.warehouse_id || adjustForm.quantity === 0) {
    ElMessage.warning(t('inventory.fillCompleteInfo'))
    return
  }

  adjustLoading.value = true
  try {
    await adjustStock({
      sku_code: adjustForm.sku_code,
      warehouse_id: adjustForm.warehouse_id,
      quantity: adjustForm.quantity,
      remark: adjustForm.remark
    })
    ElMessage.success(t('inventory.adjustSuccess'))
    adjustForm.sku_code = ''
    adjustForm.quantity = 0
    adjustForm.remark = ''
    loadInventoryLogs()
    loadLowStockAlerts()
  } catch (error) {
    handleError(error, t('inventory.adjustFailed'))
  } finally {
    adjustLoading.value = false
  }
}

// Quick adjust from low stock alert
const handleAdjustStock = (row: LowStockSKU) => {
  quickAdjustForm.sku_code = row.sku_code
  quickAdjustForm.quantity = Math.max(row.safety_stock * 2, 100)
  quickAdjustForm.remark = t('inventory.lowStockReplenish')
  adjustDialogVisible.value = true
}

const handleQuickAdjust = async () => {
  if (!quickAdjustForm.sku_code || quickAdjustForm.quantity <= 0) {
    ElMessage.warning(t('inventory.fillValidQuantity'))
    return
  }

  // Use default warehouse
  const defaultWarehouse = warehouseList.value.find(w => w.is_default)
  if (!defaultWarehouse) {
    ElMessage.warning(t('inventory.setDefaultWarehouseFirst'))
    return
  }

  quickAdjustLoading.value = true
  try {
    await adjustStock({
      sku_code: quickAdjustForm.sku_code,
      warehouse_id: defaultWarehouse.id,
      quantity: quickAdjustForm.quantity,
      remark: quickAdjustForm.remark
    })
    ElMessage.success(t('inventory.replenishSuccess'))
    adjustDialogVisible.value = false
    loadLowStockAlerts()
    loadInventoryLogs()
  } catch (error) {
    handleError(error, t('inventory.replenishFailed'))
  } finally {
    quickAdjustLoading.value = false
  }
}

// Helpers
const getChangeTypeTag = (type: string) => {
  const map: Record<string, string> = {
    manual: 'info',
    order: 'warning',
    return: 'success',
    adjustment: ''
  }
  return map[type] || 'info'
}

const getChangeTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    manual: t('inventory.manual'),
    order: t('inventory.order'),
    return: t('inventory.return'),
    adjustment: t('inventory.adjustment')
  }
  return map[type] || type
}

onMounted(() => {
  loadLowStockAlerts()
  loadInventoryLogs()
  loadWarehouses()
})
</script>

<style scoped>
.inventory-page {
  padding: 0;
}

.alert-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(245, 158, 11, 0.2);
}

.alert-card :deep(.el-card__header) {
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  border-bottom: none;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-header .el-icon {
  color: #F59E0B;
}

.main-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-bar :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-bar :deep(.el-select .el-input__wrapper) {
  border-radius: 12px;
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

.adjust-form {
  max-width: 600px;
}

.adjust-form :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.form-tip {
  margin-left: 12px;
  color: #6B7280;
  font-size: 12px;
}

/* Batch Safety Stock Form */
.batch-safety-stock-form {
  padding: 10px 0;
}

.selected-items {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
}

.selected-items h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.text-success {
  color: #10B981;
  font-weight: 600;
}

.text-danger {
  color: #EF4444;
  font-weight: 600;
}

/* Tags */
:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Dialog */
:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Tabs */
:deep(.el-tabs__item.is-active) {
  color: #6366F1;
  font-weight: 600;
}

:deep(.el-tabs__active-bar) {
  background-color: #6366F1;
}
</style>