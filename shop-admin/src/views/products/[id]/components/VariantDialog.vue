<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? $t('products.editVariant') : $t('products.addVariant')"
    width="550px"
    @update:model-value="emit('update:visible', $event)"
  >
    <el-form
      ref="variantFormRef"
      :model="variantForm"
      label-width="100px"
    >
      <el-form-item
        :label="$t('products.skuCode')"
        required
      >
        <el-input
          v-model="variantForm.code"
          :placeholder="$t('products.enterSkuCode')"
        />
      </el-form-item>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('products.price')">
            <el-input-number
              v-model="variantForm.price"
              :min="0"
              :precision="2"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.currency')">
            <el-select
              v-model="variantForm.currency"
              style="width: 100%"
            >
              <el-option
                label="USD"
                value="USD"
              />
              <el-option
                label="EUR"
                value="EUR"
              />
              <el-option
                label="GBP"
                value="GBP"
              />
              <el-option
                label="CNY"
                value="CNY"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('products.stock')">
            <el-input-number
              v-model="variantForm.stock"
              :min="0"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.safetyStock')">
            <el-input-number
              v-model="variantForm.safety_stock"
              :min="0"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item :label="$t('products.enablePreSale')">
        <el-switch v-model="variantForm.pre_sale_enabled" />
      </el-form-item>
      <el-form-item :label="$t('products.attributes')">
        <div class="attributes-section">
          <div
            v-for="(value, key) in variantForm.attributes"
            :key="key"
            class="attribute-item"
          >
            <span class="attribute-text">{{ key }}: {{ value }}</span>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleRemoveAttribute(key as string)"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
          <div class="add-attribute-row">
            <el-input
              v-model="newAttributeKey"
              :placeholder="$t('products.attributeName')"
              style="width: 120px"
            />
            <el-input
              v-model="newAttributeValue"
              :placeholder="$t('products.attributeValue')"
              style="width: 120px"
            />
            <el-button
              type="primary"
              size="small"
              @click="handleAddAttribute"
            >
              {{ $t('common.add') }}
            </el-button>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSave"
      >
        {{ $t('common.save') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Close } from '@element-plus/icons-vue'
import { createSKU, updateSKU, type CreateSKURequest } from '@/api/product'
import { t } from '@/plugins/i18n'
import type { VariantDialogProps, VariantDialogEmits, VariantFormData } from '../types'

const props = defineProps<VariantDialogProps>()
const emit = defineEmits<VariantDialogEmits>()

const variantFormRef = ref()
const variantForm = reactive<VariantFormData>({
  id: 0,
  code: '',
  price: 0,
  currency: 'USD',
  stock: 0,
  safety_stock: 0,
  pre_sale_enabled: false,
  attributes: {}
})
const newAttributeKey = ref('')
const newAttributeValue = ref('')

watch(() => props.visible, (newVal) => {
  if (newVal) {
    // Reset form when dialog opens
    variantForm.id = props.variant.id
    variantForm.code = props.variant.code
    variantForm.price = props.variant.price
    variantForm.currency = props.variant.currency
    variantForm.stock = props.variant.stock
    variantForm.safety_stock = props.variant.safety_stock
    variantForm.pre_sale_enabled = props.variant.pre_sale_enabled
    variantForm.attributes = { ...props.variant.attributes }
    newAttributeKey.value = ''
    newAttributeValue.value = ''
  }
})

const isEdit = props.isEdit

const handleClose = () => {
  emit('update:visible', false)
}

const handleAddAttribute = () => {
  if (newAttributeKey.value && newAttributeValue.value) {
    variantForm.attributes[newAttributeKey.value] = newAttributeValue.value
    newAttributeKey.value = ''
    newAttributeValue.value = ''
  }
}

const handleRemoveAttribute = (key: string) => {
  delete variantForm.attributes[key]
}

const handleSave = async () => {
  if (!variantForm.code) {
    ElMessage.warning(t('products.enterSkuCode'))
    return
  }
  if (variantForm.price <= 0) {
    ElMessage.warning(t('products.enterValidPrice'))
    return
  }

  try {
    const data: CreateSKURequest = {
      product_id: props.productId,
      code: variantForm.code,
      price: variantForm.price.toFixed(2),
      currency: variantForm.currency,
      stock: variantForm.stock,
      safety_stock: variantForm.safety_stock,
      pre_sale_enabled: variantForm.pre_sale_enabled,
      attributes: variantForm.attributes
    }

    if (isEdit) {
      await updateSKU({ ...data, id: variantForm.id })
      ElMessage.success(t('products.variantUpdateSuccess'))
    } else {
      await createSKU(data)
      ElMessage.success(t('products.variantCreateSuccess'))
    }

    emit('success')
    handleClose()
  } catch (error) {
    console.error('Failed to save variant:', error)
    ElMessage.error(isEdit ? t('products.variantUpdateFailed') : t('products.variantCreateFailed'))
  }
}
</script>

<style scoped>
.attributes-section {
  width: 100%;
}

.attribute-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.attribute-text {
  font-size: 14px;
}

.add-attribute-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-top: 8px;
}
</style>
