<template>
  <div
    v-loading="loading"
    class="variants-section"
  >
    <!-- Guided empty state when no variants -->
    <div
      v-if="variants.length === 0 && !loading"
      class="guided-empty"
    >
      <el-empty
        :description="$t('products.noVariants')"
        :image-size="120"
      >
        <p class="guided-empty-hint">
          {{ $t('products.noVariantsHint') }}
        </p>
        <el-button
          type="primary"
          size="large"
          @click="handleShowAddVariantDialog"
        >
          <el-icon><Plus /></el-icon>
          {{ $t('products.createFirstVariant') }}
        </el-button>
      </el-empty>
    </div>

    <!-- Variants table (only when there are variants) -->
    <template v-else>
      <div class="section-header">
        <h3 class="section-title">
          {{ $t('products.productVariants') }}
        </h3>
        <el-button
          type="primary"
          size="small"
          @click="handleShowAddVariantDialog"
        >
          <el-icon><Plus /></el-icon>
          {{ $t('products.addVariant') }}
        </el-button>
      </div>
      <el-table
        :data="variants"
        stripe
      >
        <el-table-column
          :label="$t('products.skuCode')"
          prop="code"
          min-width="150"
        />
        <el-table-column
          :label="$t('products.attributes')"
          min-width="200"
        >
          <template #default="{ row }">
            <div class="attribute-tags">
              <el-tag
                v-for="(value, key) in row.attributes"
                :key="key"
                size="small"
                class="attribute-tag"
              >
                {{ key }}: {{ value }}
              </el-tag>
              <span
                v-if="Object.keys(row.attributes || {}).length === 0"
                class="text-muted"
              >-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.price')"
          width="120"
          align="right"
        >
          <template #default="{ row }">
            {{ row.currency }} {{ row.price }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.stock')"
          prop="stock"
          width="100"
          align="center"
        />
        <el-table-column
          :label="$t('products.availableStock')"
          prop="available_stock"
          width="100"
          align="center"
        />
        <el-table-column
          :label="$t('products.safetyStock')"
          prop="safety_stock"
          width="100"
          align="center"
        />
        <el-table-column
          :label="$t('common.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              size="small"
              class="status-tag"
              @click="handleToggleStatus(row)"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.stockAlert')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              v-if="row.is_low_stock"
              type="danger"
              size="small"
            >
              {{ $t('products.lowStock') }}
            </el-tag>
            <span
              v-else
              class="text-muted"
            >-</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="120"
          align="center"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEditVariant(row)"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDeleteVariant(row)"
            >
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getSKUsByProduct, deleteSKU, updateSKU, type SKU } from '@/api/product'
import { t } from '@/plugins/i18n'
import type { ProductVariantsTabProps, ProductVariantsTabEmits, VariantFormData } from '../types'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<ProductVariantsTabProps>()
const { handleError } = useErrorHandler()
const emit = defineEmits<ProductVariantsTabEmits>()

const variants = ref<SKU[]>([])
const loading = ref(false)

const loadVariants = async () => {
  loading.value = true
  try {
    const response = await getSKUsByProduct(props.productId)
    variants.value = response.list || []
  } catch (error) {
    handleError(error, t('products.loadSKUsFailed'))
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  return status === '1' || status === 'enabled' ? 'success' : 'info'
}

const getStatusText = (status: string) => {
  return status === '1' || status === 'enabled' ? t('common.enabled') : t('common.disabled')
}

const skuToVariantForm = (sku: SKU): VariantFormData => ({
  id: sku.id,
  code: sku.code,
  price: parseFloat(sku.price) || 0,
  currency: sku.currency,
  stock: sku.stock,
  safety_stock: sku.safety_stock,
  pre_sale_enabled: sku.pre_sale_enabled,
  attributes: { ...sku.attributes }
})

const handleShowAddVariantDialog = () => {
  emit('variants-change')
}

const handleEditVariant = (row: SKU) => {
  emit('edit-variant', skuToVariantForm(row))
}

const handleToggleStatus = async (row: SKU) => {
  const newStatus = getStatusType(row.status) === 'success' ? 'disabled' : 'enabled'
  try {
    await updateSKU({ id: row.id, status: newStatus })
    ElMessage.success(t('products.statusUpdateSuccess'))
    loadVariants()
    emit('variants-change')
  } catch (error) {
    handleError(error, t('products.statusUpdateFailed'))
  }
}

const handleDeleteVariant = async (row: SKU) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDeleteVariant', { code: row.code }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteSKU(row.id)
    ElMessage.success(t('products.variantDeleteSuccess'))
    loadVariants()
    emit('variants-change')
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('products.variantDeleteFailed'))
    }
  }
}

onMounted(() => {
  loadVariants()
})

defineExpose({
  loadVariants,
  getVariants: () => variants.value
})
</script>

<style scoped>
.variants-section {
  padding: 0;
}

.guided-empty {
  padding: 40px 20px;
  display: flex;
  justify-content: center;
}

.guided-empty-hint {
  color: #6B7280;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 16px 0;
  max-width: 420px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.attribute-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.attribute-tag {
  margin: 2px;
}

.text-muted {
  color: #9CA3AF;
  font-size: 12px;
}

.status-tag {
  cursor: pointer;
  user-select: none;
}

.status-tag:hover {
  opacity: 0.85;
}
</style>
