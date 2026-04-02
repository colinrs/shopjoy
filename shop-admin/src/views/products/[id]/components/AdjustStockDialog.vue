<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="emit('update:visible', $event)"
    :title="$t('products.adjustStock')"
    width="450px"
  >
    <el-form :model="adjustStockForm" label-width="100px">
      <el-form-item :label="$t('products.warehouse')">
        <el-select v-model="adjustStockForm.warehouse_id" :placeholder="$t('products.selectWarehouse')" style="width: 100%">
          <el-option
            v-for="warehouse in warehouses"
            :key="warehouse.id"
            :label="warehouse.name"
            :value="warehouse.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('products.adjustQuantity')">
        <el-input-number
          v-model="adjustStockForm.quantity"
          :step="1"
          style="width: 100%"
        />
        <div class="adjust-tip">{{ $t('products.stockInOutTip') }}</div>
      </el-form-item>
      <el-form-item :label="$t('products.remark')">
        <el-input v-model="adjustStockForm.remark" type="textarea" :rows="2" :placeholder="$t('products.enterRemark')" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="loading">{{ $t('products.confirmAdjustment') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { adjustStock } from '@/api/inventory'
import { t } from '@/plugins/i18n'
import type { AdjustStockDialogProps, AdjustStockDialogEmits, AdjustStockFormData } from '../types'

const props = defineProps<AdjustStockDialogProps>()
const emit = defineEmits<AdjustStockDialogEmits>()

const adjustStockForm = reactive<AdjustStockFormData>({
  warehouse_id: 0,
  quantity: 0,
  remark: ''
})

watch(() => props.visible, (newVal) => {
  if (newVal) {
    adjustStockForm.warehouse_id = props.warehouses.find(w => w.is_default)?.id || props.warehouses[0]?.id || 0
    adjustStockForm.quantity = 0
    adjustStockForm.remark = ''
  }
})

const handleClose = () => {
  emit('update:visible', false)
}

const handleConfirm = async () => {
  if (!props.sku) {
    ElMessage.warning(t('products.skuNotExist'))
    return
  }
  if (adjustStockForm.quantity === 0) {
    ElMessage.warning(t('products.enterAdjustQuantity'))
    return
  }

  try {
    await adjustStock({
      sku_code: props.sku,
      warehouse_id: adjustStockForm.warehouse_id,
      quantity: adjustStockForm.quantity,
      remark: adjustStockForm.remark
    })
    ElMessage.success(t('products.adjustSuccess'))
    emit('success')
    handleClose()
  } catch (error) {
    console.error('Failed to adjust stock:', error)
    ElMessage.error(t('products.adjustFailed'))
  }
}
</script>

<style scoped>
.adjust-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}
</style>
