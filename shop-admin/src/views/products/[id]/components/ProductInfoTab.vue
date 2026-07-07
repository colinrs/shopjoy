<template>
  <el-form
    ref="formRef"
    :model="localForm"
    label-width="140px"
    :rules="formRules"
  >
    <!-- Basic Information Section -->
    <div class="form-section">
      <h3 class="section-title">
        {{ $t('products.basicInfo') }}
      </h3>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item
            :label="$t('products.productName')"
            prop="name"
          >
            <el-input
              v-model="localForm.name"
              :placeholder="$t('products.enterProductName')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item
            label="SKU"
            prop="sku"
          >
            <el-input
              v-model="localForm.sku"
              :placeholder="$t('products.enterSKU')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.brand')">
            <el-select
              v-model="localForm.brand"
              filterable
              clearable
              :loading="brandsLoading"
              :placeholder="$t('products.selectBrand')"
              style="width: 100%"
            >
              <el-option
                v-for="brand in props.brands"
                :key="brand.id"
                :label="brand.name"
                :value="brand.name"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.categoryPath')">
            <el-tooltip
              :content="`${$t('products.deepestCategoryId')}: ${localForm.category_id || '-'}`"
              placement="top"
              :disabled="!localForm.category_id"
            >
              <el-cascader
                v-model="localForm.category_id"
                :options="categoryOptions"
                :props="categoryCascaderProps"
                :loading="categoriesLoading"
                :placeholder="$t('products.selectCategory')"
                style="width: 100%"
              />
            </el-tooltip>
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item :label="$t('products.productDescription')">
            <el-input
              v-model="localForm.description"
              type="textarea"
              :rows="4"
              :placeholder="$t('products.enterProductDescription')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item
            :label="$t('products.price')"
            prop="price"
          >
            <el-input-number
              :model-value="localForm.price === '' ? 0 : Number(localForm.price)"
              :min="0"
              :precision="2"
              style="width: 100%"
              @update:model-value="localForm.price = String($event)"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.currency')">
            <el-select
              v-model="localForm.currency"
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
        <el-col :span="12">
          <el-form-item :label="$t('products.costPrice')">
            <el-input-number
              :model-value="localForm.cost_price === '' ? 0 : Number(localForm.cost_price)"
              :min="0"
              :precision="2"
              style="width: 100%"
              @update:model-value="localForm.cost_price = String($event)"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.stock')">
            <el-input-number
              v-model="localForm.stock"
              :min="0"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('common.status')">
            <el-select
              v-model="localForm.status"
              style="width: 100%"
            >
              <el-option
                :label="$t('products.draft')"
                value="draft"
              />
              <el-option
                :label="$t('products.onSale')"
                value="on_sale"
              />
              <el-option
                :label="$t('products.offSale')"
                value="off_sale"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.isMatrixProduct')">
            <el-switch v-model="localForm.is_matrix_product" />
          </el-form-item>
        </el-col>
      </el-row>
    </div>

    <!-- Compliance Section -->
    <div class="form-section">
      <h3 class="section-title">
        {{ $t('products.complianceInfo') }}
      </h3>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('products.hsCode')">
            <el-input
              v-model="localForm.hs_code"
              :placeholder="$t('products.enterHsCode')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('products.coo')">
            <el-input
              v-model="localForm.coo"
              :placeholder="$t('products.enterCoo')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('products.weight')">
            <el-input
              v-model="localForm.weight"
              :placeholder="$t('products.weightPlaceholder')"
            >
              <template #append>
                <el-select
                  v-model="localForm.weight_unit"
                  style="width: 80px"
                >
                  <el-option
                    label="kg"
                    value="kg"
                  />
                  <el-option
                    label="g"
                    value="g"
                  />
                  <el-option
                    label="lb"
                    value="lb"
                  />
                </el-select>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="16">
          <el-form-item :label="$t('products.dimensions')">
            <el-row :gutter="8">
              <el-col :span="8">
                <el-input
                  v-model="localForm.length"
                  :placeholder="$t('products.length')"
                >
                  <template #append>
                    cm
                  </template>
                </el-input>
              </el-col>
              <el-col :span="8">
                <el-input
                  v-model="localForm.width"
                  :placeholder="$t('products.width')"
                >
                  <template #append>
                    cm
                  </template>
                </el-input>
              </el-col>
              <el-col :span="8">
                <el-input
                  v-model="localForm.height"
                  :placeholder="$t('products.height')"
                >
                  <template #append>
                    cm
                  </template>
                </el-input>
              </el-col>
            </el-row>
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item :label="$t('products.dangerousGoods')">
            <el-checkbox-group v-model="localForm.dangerous_goods">
              <el-checkbox label="battery">
                {{ $t('products.battery') }}
              </el-checkbox>
              <el-checkbox label="liquid">
                {{ $t('products.liquid') }}
              </el-checkbox>
              <el-checkbox label="flammable">
                {{ $t('products.flammable') }}
              </el-checkbox>
              <el-checkbox label="magnetic">
                {{ $t('products.magnetic') }}
              </el-checkbox>
              <el-checkbox label="fragile">
                {{ $t('products.fragile') }}
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </el-col>
      </el-row>
    </div>

    <!-- Images Section -->
    <div class="form-section">
      <h3 class="section-title">
        {{ $t('products.productImage') }}
      </h3>
      <el-form-item :label="$t('products.imageUrl')">
        <div class="image-list">
          <div
            v-for="(img, index) in localForm.images"
            :key="index"
            class="image-item"
          >
            <el-image
              :src="img"
              fit="cover"
              class="product-image"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <el-button
              type="danger"
              size="small"
              circle
              class="remove-btn"
              @click="handleRemoveImage(index)"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
          <div
            class="add-image"
            @click="handleShowAddImage"
          >
            <el-icon><Plus /></el-icon>
            <span>{{ $t('products.addImage') }}</span>
          </div>
        </div>
      </el-form-item>
      <div class="save-actions">
        <el-button
          type="primary"
          :disabled="!isDirty"
          :loading="saveLoading"
          @click="handleSave"
        >
          <el-icon><Check /></el-icon>
          {{ $t('products.saveChanges') }}
        </el-button>
      </div>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Picture, Plus, Close, Check } from '@element-plus/icons-vue'
import { t } from '@/plugins/i18n'
import type { FormInstance } from 'element-plus'
import type { ProductInfoTabProps, ProductInfoTabEmits, ProductFormData } from '../types'
import type { CategoryTree } from '@/api/category'

const props = defineProps<ProductInfoTabProps>()
const emit = defineEmits<ProductInfoTabEmits>()

const formRef = ref<FormInstance>()

const localForm = ref<ProductFormData>({ ...props.productForm })

const categoryOptions = computed(() => {
  const transform = (nodes: CategoryTree[]): CategoryTree[] => {
    return nodes.map(node => {
      const children = node.children?.length ? transform(node.children) : undefined
      return {
        ...node,
        children
      } as CategoryTree
    })
  }
  return transform(props.categories)
})

const categoryCascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  emitPath: false
}

watch(() => props.productForm, (newVal) => {
  localForm.value = { ...newVal }
}, { deep: true })

watch(localForm, (newVal) => {
  emit('update:productForm', newVal)
}, { deep: true })

const formRules = {
  name: [{ required: true, message: () => t('products.enterProductName'), trigger: 'blur' }],
  price: [{ required: true, message: () => t('products.enterPrice'), trigger: 'blur' }]
}

const handleShowAddImage = () => {
  emit('show-add-image')
}

const handleSave = () => {
  emit('save')
}

const handleRemoveImage = (index: number) => {
  localForm.value.images.splice(index, 1)
}

// Expose form ref for parent validation
defineExpose({
  formRef: () => formRef.value,
  validate: async () => {
    if (formRef.value) {
      return formRef.value.validate()
    }
    return true
  }
})
</script>

<style scoped>
.form-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #E5E7EB;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.image-item {
  position: relative;
  width: 120px;
  height: 120px;
}

.product-image {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #9CA3AF;
}

.remove-btn {
  position: absolute;
  top: -8px;
  right: -8px;
}

.add-image {
  width: 120px;
  height: 120px;
  border: 2px dashed #D1D5DB;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  color: #9CA3AF;
  transition: all 0.2s;
}

.add-image:hover {
  border-color: #409EFF;
  color: #409EFF;
}

.category-path-display {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #374151;
  font-size: 13px;
  line-height: 32px;
}

.save-actions {
  display: flex;
  justify-content: center;
  padding-top: 16px;
}
</style>
