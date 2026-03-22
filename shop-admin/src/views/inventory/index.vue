<template>
  <div class="inventory-page">
    <!-- Low Stock Alert Card -->
    <el-card class="alert-card" shadow="never" v-if="lowStockList.length > 0">
      <template #header>
        <div class="card-header">
          <span><el-icon><Warning /></el-icon> 库存预警</span>
          <el-badge :value="lowStockTotal" type="danger" />
        </div>
      </template>
      <el-table :data="lowStockList" stripe size="small">
        <el-table-column prop="sku_code" label="SKU编码" width="150" />
        <el-table-column prop="product_name" label="商品名称" min-width="200" />
        <el-table-column prop="available_stock" label="可用库存" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="danger" size="small">{{ row.available_stock }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="safety_stock" label="安全库存" width="100" align="center" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAdjustStock(row)">
              补货
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Main Content Tabs -->
    <el-card class="main-card" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- Inventory Logs Tab -->
        <el-tab-pane label="库存记录" name="logs">
          <div class="filter-bar">
            <el-input
              v-model="logFilter.sku_code"
              placeholder="SKU编码"
              clearable
              style="width: 180px"
            />
            <el-select v-model="logFilter.type" placeholder="变更类型" clearable style="width: 140px">
              <el-option label="手动调整" value="manual" />
              <el-option label="订单扣减" value="order" />
              <el-option label="退货入库" value="return" />
              <el-option label="库存盘点" value="adjustment" />
            </el-select>
            <el-button type="primary" @click="loadInventoryLogs">查询</el-button>
          </div>
          <el-table :data="logList" v-loading="logLoading" stripe>
            <el-table-column prop="sku_code" label="SKU编码" width="150" />
            <el-table-column prop="change_type" label="变更类型" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="getChangeTypeTag(row.change_type)">
                  {{ getChangeTypeLabel(row.change_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="change_quantity" label="变更数量" width="100" align="center">
              <template #default="{ row }">
                <span :class="row.change_quantity >= 0 ? 'text-success' : 'text-danger'">
                  {{ row.change_quantity >= 0 ? '+' : '' }}{{ row.change_quantity }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="before_stock" label="变更前" width="100" align="center" />
            <el-table-column prop="after_stock" label="变更后" width="100" align="center" />
            <el-table-column prop="order_no" label="关联订单" width="150" />
            <el-table-column prop="remark" label="备注" min-width="150" />
            <el-table-column prop="created_at" label="时间" width="180" />
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
        <el-tab-pane label="仓库管理" name="warehouses">
          <div class="filter-bar">
            <el-button type="primary" @click="handleAddWarehouse">
              <el-icon><Plus /></el-icon>新增仓库
            </el-button>
          </div>
          <el-table :data="warehouseList" v-loading="warehouseLoading" stripe>
            <el-table-column prop="code" label="仓库编码" width="120" />
            <el-table-column prop="name" label="仓库名称" min-width="150" />
            <el-table-column prop="country" label="国家" width="100" />
            <el-table-column prop="address" label="地址" min-width="200" />
            <el-table-column prop="is_default" label="默认仓库" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="(val: number) => handleWarehouseStatusChange(row, val)"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditWarehouse(row)">
                  编辑
                </el-button>
                <el-button
                  type="success"
                  link
                  size="small"
                  @click="handleSetDefaultWarehouse(row)"
                  v-if="!row.is_default"
                >
                  设为默认
                </el-button>
                <el-button
                  type="danger"
                  link
                  size="small"
                  @click="handleDeleteWarehouse(row)"
                  v-if="!row.is_default"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Stock Adjustment Tab -->
        <el-tab-pane label="库存调整" name="adjust">
          <el-form :model="adjustForm" label-width="120px" class="adjust-form">
            <el-form-item label="SKU编码" required>
              <el-input v-model="adjustForm.sku_code" placeholder="请输入SKU编码" style="width: 300px" />
            </el-form-item>
            <el-form-item label="仓库">
              <el-select v-model="adjustForm.warehouse_id" placeholder="选择仓库" style="width: 300px">
                <el-option
                  v-for="w in warehouseList"
                  :key="w.id"
                  :label="w.name"
                  :value="w.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="调整数量" required>
              <el-input-number
                v-model="adjustForm.quantity"
                :min="-9999"
                :max="9999"
                style="width: 200px"
              />
              <span class="form-tip">正数增加库存，负数减少库存</span>
            </el-form-item>
            <el-form-item label="备注">
              <el-input
                v-model="adjustForm.remark"
                type="textarea"
                rows="3"
                placeholder="请输入调整原因"
                style="width: 400px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleAdjustSubmit" :loading="adjustLoading">
                确认调整
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Warehouse Dialog -->
    <el-dialog
      v-model="warehouseDialogVisible"
      :title="isEditWarehouse ? '编辑仓库' : '新增仓库'"
      width="500px"
      destroy-on-close
    >
      <el-form :model="warehouseForm" label-width="100px" :rules="warehouseRules" ref="warehouseFormRef">
        <el-form-item label="仓库编码" prop="code" v-if="!isEditWarehouse">
          <el-input v-model="warehouseForm.code" placeholder="如: WH-001" />
        </el-form-item>
        <el-form-item label="仓库名称" prop="name">
          <el-input v-model="warehouseForm.name" placeholder="请输入仓库名称" />
        </el-form-item>
        <el-form-item label="国家">
          <el-input v-model="warehouseForm.country" placeholder="如: CN, US" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="warehouseForm.address" type="textarea" rows="2" placeholder="详细地址" />
        </el-form-item>
        <el-form-item label="默认仓库">
          <el-switch v-model="warehouseForm.is_default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="warehouseDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveWarehouse" :loading="warehouseSaveLoading">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- Stock Adjust Dialog (from low stock alert) -->
    <el-dialog v-model="adjustDialogVisible" title="库存补货" width="400px" destroy-on-close>
      <el-form :model="quickAdjustForm" label-width="100px">
        <el-form-item label="SKU编码">
          <el-input :value="quickAdjustForm.sku_code" disabled />
        </el-form-item>
        <el-form-item label="补货数量">
          <el-input-number v-model="quickAdjustForm.quantity" :min="1" :max="9999" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="quickAdjustForm.remark" placeholder="补货备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleQuickAdjust" :loading="quickAdjustLoading">
          确认补货
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning } from '@element-plus/icons-vue'
import {
  getWarehouses,
  createWarehouse,
  updateWarehouse,
  updateWarehouseStatus,
  setDefaultWarehouse,
  deleteWarehouse,
  getInventoryLogs,
  getLowStockSKUs,
  adjustStock,
  type Warehouse,
  type InventoryLog,
  type LowStockSKU
} from '@/api/inventory'

// Low stock alerts
const lowStockList = ref<LowStockSKU[]>([])
const lowStockTotal = ref(0)

// Inventory logs
const logList = ref<InventoryLog[]>([])
const logLoading = ref(false)
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
  code: [{ required: true, message: '请输入仓库编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入仓库名称', trigger: 'blur' }]
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
  remark: '低库存补货'
})

// Load low stock alerts
const loadLowStockAlerts = async () => {
  try {
    const res = await getLowStockSKUs({ page: 1, page_size: 10 })
    lowStockList.value = res.list || []
    lowStockTotal.value = res.total || 0
  } catch (error) {
    console.error('Failed to load low stock alerts:', error)
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
    console.error('Failed to load inventory logs:', error)
    ElMessage.error('加载库存记录失败')
  } finally {
    logLoading.value = false
  }
}

// Load warehouses
const loadWarehouses = async () => {
  warehouseLoading.value = true
  try {
    const res = await getWarehouses()
    warehouseList.value = res || []
  } catch (error) {
    console.error('Failed to load warehouses:', error)
    ElMessage.error('加载仓库失败')
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
          ElMessage.success('更新成功')
        } else {
          await createWarehouse({
            code: warehouseForm.code,
            name: warehouseForm.name,
            country: warehouseForm.country,
            address: warehouseForm.address,
            is_default: warehouseForm.is_default
          })
          ElMessage.success('创建成功')
        }
        warehouseDialogVisible.value = false
        loadWarehouses()
      } catch (error) {
        console.error('Failed to save warehouse:', error)
        ElMessage.error(isEditWarehouse.value ? '更新失败' : '创建失败')
      } finally {
        warehouseSaveLoading.value = false
      }
    }
  })
}

const handleWarehouseStatusChange = async (row: Warehouse, status: number) => {
  try {
    await updateWarehouseStatus(row.id, status)
    ElMessage.success(status === 1 ? '已启用' : '已禁用')
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error('更新状态失败')
    row.status = status === 1 ? 0 : 1
  }
}

const handleSetDefaultWarehouse = async (row: Warehouse) => {
  try {
    await setDefaultWarehouse(row.id)
    ElMessage.success('已设为默认仓库')
    loadWarehouses()
  } catch (error) {
    console.error('Failed to set default:', error)
    ElMessage.error('操作失败')
  }
}

const handleDeleteWarehouse = (row: Warehouse) => {
  ElMessageBox.confirm(`确认删除仓库 "${row.name}"?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteWarehouse(row.id)
      ElMessage.success('删除成功')
      loadWarehouses()
    } catch (error) {
      console.error('Failed to delete warehouse:', error)
      ElMessage.error('删除失败')
    }
  })
}

// Stock adjustment handlers
const handleAdjustSubmit = async () => {
  if (!adjustForm.sku_code || !adjustForm.warehouse_id || adjustForm.quantity === 0) {
    ElMessage.warning('请填写完整信息')
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
    ElMessage.success('库存调整成功')
    adjustForm.sku_code = ''
    adjustForm.quantity = 0
    adjustForm.remark = ''
    loadInventoryLogs()
    loadLowStockAlerts()
  } catch (error) {
    console.error('Failed to adjust stock:', error)
    ElMessage.error('调整失败')
  } finally {
    adjustLoading.value = false
  }
}

// Quick adjust from low stock alert
const handleAdjustStock = (row: LowStockSKU) => {
  quickAdjustForm.sku_code = row.sku_code
  quickAdjustForm.quantity = Math.max(row.safety_stock * 2, 100)
  quickAdjustForm.remark = '低库存补货'
  adjustDialogVisible.value = true
}

const handleQuickAdjust = async () => {
  if (!quickAdjustForm.sku_code || quickAdjustForm.quantity <= 0) {
    ElMessage.warning('请填写有效数量')
    return
  }

  // Use default warehouse
  const defaultWarehouse = warehouseList.value.find(w => w.is_default)
  if (!defaultWarehouse) {
    ElMessage.warning('请先设置默认仓库')
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
    ElMessage.success('补货成功')
    adjustDialogVisible.value = false
    loadLowStockAlerts()
    loadInventoryLogs()
  } catch (error) {
    console.error('Failed to quick adjust:', error)
    ElMessage.error('补货失败')
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
    manual: '手动调整',
    order: '订单扣减',
    return: '退货入库',
    adjustment: '库存盘点'
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