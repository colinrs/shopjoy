<template>
  <div class="product-detail-page">
    <!-- Page Header -->
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="page-header">
        <div class="header-left">
          <el-button
            link
            @click="handleBack"
          >
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('products.backToList') }}
          </el-button>
          <el-divider direction="vertical" />
          <h2 class="product-title">
            {{ product?.name || $t('common.loading') }}
            <el-tag
              v-if="product"
              :type="getStatusType(product.status)"
              size="small"
            >
              {{ getStatusText(product.status) }}
            </el-tag>
          </h2>
        </div>
        <div class="header-right">
          <el-button
            v-if="product?.status === 'on_sale'"
            :loading="statusLoading"
            @click="handleTakeOffSale"
          >
            <el-icon><Hide /></el-icon>
            {{ $t('products.offSale') }}
          </el-button>
          <el-button
            v-else-if="product?.status === 'off_sale' || product?.status === 'draft'"
            type="success"
            :loading="statusLoading"
            @click="handlePutOnSale"
          >
            <el-icon><View /></el-icon>
            {{ $t('products.onSale') }}
          </el-button>
          <el-button
            type="primary"
            :loading="saveLoading"
            @click="handleSave"
          >
            <el-icon><Check /></el-icon>
            {{ $t('products.saveChanges') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Loading State -->
    <el-skeleton
      v-if="loading"
      :rows="10"
      animated
    />

    <!-- Tab Layout -->
    <el-card
      v-else
      class="tabs-card"
      shadow="never"
    >
      <el-tabs
        v-model="activeTab"
        class="product-tabs"
      >
        <!-- Basic Info Tab -->
        <el-tab-pane
          :label="$t('products.basicInfo')"
          name="basic"
        >
          <ProductInfoTab
            :product="product"
            :product-form="productForm"
            :form-ref="formRef"
            :loading="loading"
            @update:product-form="handleProductFormUpdate"
            @save="handleShowAddImageDialog"
          />
        </el-tab-pane>

        <!-- Markets Tab -->
        <el-tab-pane
          :label="$t('products.markets')"
          name="markets"
        >
          <ProductMarketsTab
            :product-id="productId"
            :product-markets="productMarkets"
            :markets="markets"
            :loading="marketsLoading"
            @update:product-markets="handleProductMarketsUpdate"
            @refresh="loadProductMarkets"
            @show-push-to-market="handleShowPushToMarketDialog"
          />
        </el-tab-pane>

        <!-- Variants Tab -->
        <el-tab-pane
          :label="$t('products.variants')"
          name="variants"
        >
          <ProductVariantsTab
            :product-id="productId"
            :default-price="productForm.price"
            :default-currency="productForm.currency"
            :loading="variantsLoading"
            @variants-change="handleShowVariantDialog"
          />
        </el-tab-pane>

        <!-- Pricing Tab -->
        <el-tab-pane
          :label="$t('products.pricing')"
          name="pricing"
        >
          <ProductPricingTab
            :product-id="productId"
            :product-markets="productMarkets"
            :loading="marketsLoading"
            @refresh="loadProductMarkets"
          />
        </el-tab-pane>

        <!-- Localization Tab -->
        <el-tab-pane
          :label="$t('products.localization')"
          name="localization"
        >
          <ProductLocalizationTab
            :product-id="productId"
            :product-name="productForm.name"
            :product-description="productForm.description"
            :loading="localizationsLoading"
            @localizations-change="loadLocalizations"
          />
        </el-tab-pane>

        <!-- Inventory Tab -->
        <el-tab-pane
          :label="$t('products.stock')"
          name="inventory"
        >
          <ProductInventoryTab
            :product-id="productId"
            :sku="productForm.sku"
            :loading="inventoryLoading"
            @inventory-change="handleShowAdjustStockDialog"
          />
        </el-tab-pane>

        <!-- Reviews Tab -->
        <el-tab-pane
          :label="$t('products.reviewStats')"
          name="reviews"
        >
          <ProductReviewsTab
            :product-id="productId"
            :loading="reviewsLoading"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Dialogs -->
    <PushToMarketDialog
      v-model:visible="pushToMarketDialogVisible"
      :product-id="productId"
      :product-markets="productMarkets"
      :markets="markets"
      :product-price="productForm.price"
      :loading="pushToMarketLoading"
      @success="handlePushToMarketSuccess"
    />

    <AddImageDialog
      v-model:visible="addImageDialogVisible"
      @add="handleAddImage"
    />

    <AdjustStockDialog
      v-model:visible="adjustStockDialogVisible"
      :sku="productForm.sku"
      :warehouses="warehouses"
      :loading="false"
      @success="handleAdjustStockSuccess"
    />

    <VariantDialog
      v-model:visible="variantDialogVisible"
      :product-id="productId"
      :is-edit="isEditVariant"
      :variant="variantForm"
      :loading="false"
      @success="handleVariantSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Check, Hide, View } from '@element-plus/icons-vue'
import { getProductStatusType } from '@/utils/status'
import {
  getProduct,
  updateProduct,
  getProductMarkets,
  putOnSale,
  takeOffSale,
  type Product,
  type ProductMarket
} from '@/api/product'
import { getMarkets, type Market } from '@/api/market'
import {
  getWarehouses,
  type Warehouse
} from '@/api/inventory'
import { t } from '@/plugins/i18n'
import {
  ProductInfoTab,
  ProductMarketsTab,
  ProductVariantsTab,
  ProductPricingTab,
  ProductLocalizationTab,
  ProductInventoryTab,
  ProductReviewsTab,
  PushToMarketDialog,
  AddImageDialog,
  AdjustStockDialog,
  VariantDialog
} from './components'
import type { ProductFormData, VariantFormData } from './types'
import { useErrorHandler } from '@/composables/useErrorHandler'

const route = useRoute()
const router = useRouter()
const { handleError } = useErrorHandler()

const productId = computed(() => Number(route.params.id))

// State
const loading = ref(false)
const saveLoading = ref(false)
const statusLoading = ref(false)
const marketsLoading = ref(false)
const pushToMarketLoading = ref(false)
const activeTab = ref('basic')
const product = ref<Product | null>(null)
const productMarkets = ref<ProductMarket[]>([])
const markets = ref<Market[]>([])
const pushToMarketDialogVisible = ref(false)
const addImageDialogVisible = ref(false)
const adjustStockDialogVisible = ref(false)
const formRef = ref()

// Inventory state
const inventoryLoading = ref(false)
const warehouses = ref<Warehouse[]>([])

// Reviews state
const reviewsLoading = ref(false)

// Variants state
const variantsLoading = ref(false)
const variantDialogVisible = ref(false)
const isEditVariant = ref(false)
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

// Localization state
const localizationsLoading = ref(false)

// Form
const productForm = reactive<ProductFormData>({
  name: '',
  description: '',
  price: '0',
  currency: 'USD',
  cost_price: '0',
  stock: 0,
  status: 'draft',
  category_id: 0,
  sku: '',
  brand: '',
  tags: [],
  images: [],
  is_matrix_product: false,
  hs_code: '',
  coo: '',
  weight: '',
  weight_unit: 'kg',
  length: '',
  width: '',
  height: '',
  dangerous_goods: []
})

// Helper functions
const getStatusType = getProductStatusType

const getStatusText = (status: string) => {
  const statusKeyMap: Record<string, string> = {
    on_sale: 'products.onSale',
    off_sale: 'products.offSale',
    draft: 'products.draft'
  }
  const key = statusKeyMap[status]
  if (key) {
    const translated = t(key)
    return translated !== key ? translated : status
  }
  return status
}

// Navigation
const handleBack = () => {
  router.push('/products')
}

// Load product data
const loadProduct = async () => {
  loading.value = true
  try {
    const data = await getProduct(productId.value)
    product.value = data
    // Populate form - price is already in yuan as string from API
    Object.assign(productForm, {
      name: data.name || '',
      description: data.description || '',
      price: data.price || '0',
      currency: data.currency || 'USD',
      cost_price: data.cost_price || '0',
      stock: data.stock || 0,
      status: data.status || 'draft',
      category_id: data.category_id || 0,
      sku: data.sku || '',
      brand: data.brand || '',
      tags: data.tags || [],
      images: data.images || [],
      is_matrix_product: data.is_matrix_product || false,
      hs_code: data.hs_code || '',
      coo: data.coo || '',
      weight: data.weight || '',
      weight_unit: data.weight_unit || 'kg',
      length: data.length || '',
      width: data.width || '',
      height: data.height || '',
      dangerous_goods: data.dangerous_goods || []
    })
  } catch (error) {
    handleError(error, t('products.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Load product markets
const loadProductMarkets = async () => {
  marketsLoading.value = true
  try {
    const response = await getProductMarkets(productId.value)
    productMarkets.value = response.list || []
  } catch (error) {
    handleError(error, t('products.loadMarketsFailed'))
  } finally {
    marketsLoading.value = false
  }
}

// Load all markets
const loadMarkets = async () => {
  try {
    const response = await getMarkets()
    markets.value = response.list || []
  } catch (error) {
    handleError(error, t('products.loadMarketsFailed'))
  }
}

// Put on sale
const handlePutOnSale = async () => {
  statusLoading.value = true
  try {
    await putOnSale(productId.value)
    ElMessage.success(t('products.onSaleSuccess'))
    loadProduct()
  } catch (error) {
    handleError(error, t('products.onSaleFailed'))
  } finally {
    statusLoading.value = false
  }
}

// Take off sale
const handleTakeOffSale = async () => {
  statusLoading.value = true
  try {
    await takeOffSale(productId.value)
    ElMessage.success(t('products.offSaleSuccess'))
    loadProduct()
  } catch (error) {
    handleError(error, t('products.offSaleFailed'))
  } finally {
    statusLoading.value = false
  }
}

// Save product
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saveLoading.value = true
  try {
    await updateProduct({
      id: productId.value,
      name: productForm.name,
      description: productForm.description,
      price: productForm.price,
      currency: productForm.currency,
      category_id: productForm.category_id,
      sku: productForm.sku,
      brand: productForm.brand,
      tags: productForm.tags,
      images: productForm.images,
      is_matrix_product: productForm.is_matrix_product,
      hs_code: productForm.hs_code,
      coo: productForm.coo,
      weight: productForm.weight,
      weight_unit: productForm.weight_unit,
      length: productForm.length,
      width: productForm.width,
      height: productForm.height,
      dangerous_goods: productForm.dangerous_goods
    })
    ElMessage.success(t('products.updateSuccess'))
    loadProduct()
  } catch (error) {
    handleError(error, t('products.updateFailed'))
  } finally {
    saveLoading.value = false
  }
}

// Handle product form updates from child components
const handleProductFormUpdate = (form: ProductFormData) => {
  Object.assign(productForm, form)
}

const handleProductMarketsUpdate = (markets: ProductMarket[]) => {
  productMarkets.value = markets
}

// Show add image dialog
const handleShowAddImageDialog = () => {
  addImageDialogVisible.value = true
}

// Handle add image
const handleAddImage = (url: string) => {
  productForm.images.push(url)
}

// Show push to market dialog
const handleShowPushToMarketDialog = () => {
  pushToMarketDialogVisible.value = true
}

// Handle push to market success
const handlePushToMarketSuccess = () => {
  loadProductMarkets()
}

// Load warehouses
const loadWarehouses = async () => {
  try {
    const response = await getWarehouses()
    warehouses.value = response.list || []
  } catch (error) {
    handleError(error)
  }
}

// Show adjust stock dialog
const handleShowAdjustStockDialog = () => {
  adjustStockDialogVisible.value = true
}

// Handle adjust stock success
const handleAdjustStockSuccess = () => {
  loadProduct()
}

// Load localizations
const loadLocalizations = async () => {
  localizationsLoading.value = true
  try {
    // Re-import is handled by the child component
  } finally {
    localizationsLoading.value = false
  }
}

// Show variant dialog
const handleShowVariantDialog = () => {
  isEditVariant.value = false
  variantForm.id = 0
  variantForm.code = ''
  variantForm.price = parseFloat(productForm.price) || 0
  variantForm.currency = productForm.currency
  variantForm.stock = 0
  variantForm.safety_stock = 0
  variantForm.pre_sale_enabled = false
  variantForm.attributes = {}
  variantDialogVisible.value = true
}

// Handle variant success
const handleVariantSuccess = () => {
  loadProductMarkets()
}

// Initialize
onMounted(() => {
  loadProduct()
  loadProductMarkets()
  loadMarkets()
  loadWarehouses()
})
</script>

<style scoped>
.product-detail-page {
  padding: 0;
}

/* Header */
.header-card {
  margin-bottom: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* Tabs Card */
.tabs-card {
  min-height: 500px;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-left {
    flex-direction: column;
    align-items: flex-start;
  }

  .product-title {
    font-size: 18px;
  }
}
</style>
