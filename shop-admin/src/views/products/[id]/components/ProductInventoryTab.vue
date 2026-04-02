<template>
  <div class="inventory-section" v-loading="inventoryLoading">
    <!-- Inventory Overview -->
    <div class="inventory-overview">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic :title="$t('products.totalStock')" :value="skuInventory?.total_stock || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.availableStock')" :value="skuInventory?.available_stock || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.lockedStock')" :value="skuInventory?.locked_stock || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.safetyStock')" :value="skuInventory?.safety_stock || 0" />
        </el-col>
      </el-row>
    </div>

    <!-- Warehouse Inventory -->
    <div class="warehouse-inventory">
      <div class="section-header">
        <h3 class="section-title">{{ $t('products.warehouseInventory') }}</h3>
        <el-button type="primary" size="small" @click="handleShowAdjustStockDialog">
          <el-icon><Edit /></el-icon>
          {{ $t('products.adjustStock') }}
        </el-button>
      </div>
      <el-table :data="skuInventory?.warehouses || []" stripe>
        <el-table-column :label="$t('products.warehouse')" min-width="150">
          <template #default="{ row }">
            <span>{{ row.warehouse_name || `${$t('products.warehouse')} ${row.warehouse_id}` }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.availableStock')" prop="available_stock" width="120" align="center" />
        <el-table-column :label="$t('products.lockedStock')" prop="locked_stock" width="120" align="center" />
        <el-table-column :label="$t('products.totalStock')" width="120" align="center">
          <template #default="{ row }">
            {{ row.available_stock + row.locked_stock }}
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!skuInventory?.warehouses?.length" :description="$t('products.noWarehouseInventory')" />
    </div>

    <!-- Inventory Logs -->
    <div class="inventory-logs">
      <h3 class="section-title">{{ $t('products.inventoryChangeLog') }}</h3>
      <el-table :data="inventoryLogs" stripe>
        <el-table-column :label="$t('products.time')" width="180">
          <template #default="{ row }">
            {{ row.created_at }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.type')" width="100">
          <template #default="{ row }">
            <el-tag :type="getLogTypeStyle(row.change_type)" size="small">
              {{ getLogTypeText(row.change_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.changeQuantity')" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.change_quantity >= 0 ? 'text-success' : 'text-danger'">
              {{ row.change_quantity >= 0 ? '+' : '' }}{{ row.change_quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.beforeChange')" prop="before_stock" width="100" align="center" />
        <el-table-column :label="$t('products.afterChange')" prop="after_stock" width="100" align="center" />
        <el-table-column :label="$t('products.remark')" prop="remark" min-width="150" />
      </el-table>
      <el-pagination
        v-if="inventoryLogsTotal > 0"
        class="pagination"
        background
        layout="total, prev, pager, next"
        :total="inventoryLogsTotal"
        :page-size="inventoryLogsPageSize"
        :current-page="inventoryLogsPage"
        @current-change="handleInventoryLogsPageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Edit } from '@element-plus/icons-vue'
import { getSKUInventory, getInventoryLogs, type SKUInventory, type InventoryLog } from '@/api/inventory'
import { t } from '@/plugins/i18n'
import type { ProductInventoryTabProps, ProductInventoryTabEmits } from '../types'

const props = defineProps<ProductInventoryTabProps>()
const emit = defineEmits<ProductInventoryTabEmits>()

const inventoryLoading = ref(false)
const skuInventory = ref<SKUInventory | null>(null)
const inventoryLogs = ref<InventoryLog[]>([])
const inventoryLogsTotal = ref(0)
const inventoryLogsPage = ref(1)
const inventoryLogsPageSize = ref(10)

const loadInventoryData = async () => {
  if (!props.sku) return

  inventoryLoading.value = true
  try {
    // Load SKU inventory
    const inventory = await getSKUInventory(props.sku)
    skuInventory.value = inventory

    // Load inventory logs
    const logs = await getInventoryLogs({
      page: inventoryLogsPage.value,
      page_size: inventoryLogsPageSize.value,
      sku_code: props.sku
    })
    inventoryLogs.value = logs.list || []
    inventoryLogsTotal.value = logs.total || 0
  } catch (error) {
    console.error('Failed to load inventory:', error)
  } finally {
    inventoryLoading.value = false
  }
}

const handleInventoryLogsPageChange = (page: number) => {
  inventoryLogsPage.value = page
  loadInventoryData()
}

const handleShowAdjustStockDialog = () => {
  emit('inventory-change')
}

const getLogTypeStyle = (type: string) => {
  const styles: Record<string, string> = {
    manual: 'primary',
    order: 'warning',
    return: 'success',
    adjustment: 'info'
  }
  return styles[type] || 'info'
}

const getLogTypeText = (type: string) => {
  const texts: Record<string, string> = {
    manual: t('products.manual'),
    order: t('products.order'),
    return: t('products.return'),
    adjustment: t('products.adjustment')
  }
  return texts[type] || type
}

onMounted(() => {
  loadInventoryData()
})

defineExpose({
  loadInventoryData,
  getSkuInventory: () => skuInventory.value
})
</script>

<style scoped>
.inventory-section {
  padding: 0;
}

.inventory-overview {
  padding: 20px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.warehouse-inventory {
  margin-bottom: 24px;
}

.inventory-logs {
  margin-top: 24px;
}

.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}

.text-success {
  color: #10B981;
  font-weight: 500;
}

.text-danger {
  color: #EF4444;
  font-weight: 500;
}
</style>
