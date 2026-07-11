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
              class="status-tag"
              :loading="statusLoading"
              @click="handleToggleStatus"
            >
              {{ getStatusText(product.status) }}
            </el-tag>
          </h2>
        </div>
        <div class="header-right" />
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
            ref="productInfoTabRef"
            :product="product"
            :product-form="productForm"
            :loading="loading"
            :is-dirty="isDirty"
            :save-loading="saveLoading"
            :categories="categories"
            :brands="brands"
            :categories-loading="categoriesLoading"
            :brands-loading="brandsLoading"
            @update:product-form="handleProductFormUpdate"
            @save="handleSave"
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
            ref="variantsTabRef"
            :product-id="productId"
            :default-price="productForm.price"
            :default-currency="productForm.currency"
            :loading="variantsLoading"
            @variants-change="handleShowVariantDialog"
            @edit-variant="handleEditVariant"
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
            ref="inventoryTabRef"
            :product-id="productId"
            :sku="productForm.sku"
            :loading="inventoryLoading"
            @inventory-change="handleShowAdjustStockDialog"
            @go-to-variants="handleGoToVariants"
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
import { ref, reactive, onMounted, computed, nextTick, onBeforeUnmount } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
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
import { getCategoryTree, type CategoryTree } from '@/api/category'
import { getBrands, type Brand } from '@/api/brand'
import {
  ProductInfoTab,
  ProductMarketsTab,
  ProductVariantsTab,
  ProductPricingTab,
  ProductLocalizationTab,
  ProductInventoryTab,
  ProductReviewsTab,
  PushToMarketDialog,
  AdjustStockDialog,
  VariantDialog
} from './components'
import type { ProductFormData, VariantFormData } from './types'
import { useErrorHandler } from '@/composables/useErrorHandler'

const route = useRoute()
const router = useRouter()
const { handleError } = useErrorHandler()

const productId = computed(() => route.params.id as string)

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
const categories = ref<CategoryTree[]>([])
const brands = ref<Brand[]>([])
const categoriesLoading = ref(false)
const brandsLoading = ref(false)
const pushToMarketDialogVisible = ref(false)
const adjustStockDialogVisible = ref(false)
const productInfoTabRef = ref()
const originalForm = ref<ProductFormData | null>(null)

const deepEqual = (a: unknown, b: unknown): boolean => {
  if (a === b) return true
  if (typeof a !== 'object' || typeof b !== 'object' || a == null || b == null) return false
  const keysA = Object.keys(a as object)
  const keysB = Object.keys(b as object)
  if (keysA.length !== keysB.length) return false
  return keysA.every(key => deepEqual((a as Record<string, unknown>)[key], (b as Record<string, unknown>)[key]))
}

const isDirty = computed(() => {
  if (!originalForm.value) return false
  return !deepEqual(productForm, originalForm.value)
})

// Inventory state
const inventoryLoading = ref(false)
const warehouses = ref<Warehouse[]>([])

// Reviews state
const reviewsLoading = ref(false)

// Variants state
const variantsLoading = ref(false)
const variantDialogVisible = ref(false)
const isEditVariant = ref(false)
const variantsTabRef = ref()
const inventoryTabRef = ref()
const variantForm = reactive<VariantFormData>({
  id: '',
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
  category_id: '',
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
      category_id: data.category_id || '',
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
    originalForm.value = JSON.parse(JSON.stringify(productForm))
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

// Load category tree
const loadCategories = async () => {
  categoriesLoading.value = true
  try {
    const data = await getCategoryTree()
    categories.value = data || []
  } catch (error) {
    handleError(error, t('products.loadCategoriesFailed'))
  } finally {
    categoriesLoading.value = false
  }
}

// Load enabled brands
const loadBrands = async () => {
  brandsLoading.value = true
  try {
    const response = await getBrands({ status: 1, page_size: 1000 })
    brands.value = response.list || []
  } catch (error) {
    handleError(error, t('products.loadBrandsFailed'))
  } finally {
    brandsLoading.value = false
  }
}

const handleToggleStatus = async () => {
  if (!product.value || statusLoading.value) return
  if (product.value.status === 'on_sale') {
    await handleTakeOffSale()
  } else {
    await handlePutOnSale()
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
  try {
    await productInfoTabRef.value?.validate()
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
      weight: parseFloat(productForm.weight) || 0,
      weight_unit: productForm.weight_unit,
      length: parseFloat(productForm.length) || 0,
      width: parseFloat(productForm.width) || 0,
      height: parseFloat(productForm.height) || 0,
      dangerous_goods: productForm.dangerous_goods
    })
    ElMessage.success(t('products.updateSuccess'))
    await loadProduct()
    originalForm.value = JSON.parse(JSON.stringify(productForm))
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

// Switch to Variants tab (triggered from Inventory tab when no SKU)
const handleGoToVariants = () => {
  activeTab.value = 'variants'
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
  variantForm.id = ''
  variantForm.code = ''
  variantForm.price = parseFloat(productForm.price) || 0
  variantForm.currency = productForm.currency
  variantForm.stock = 0
  variantForm.safety_stock = 0
  variantForm.pre_sale_enabled = false
  variantForm.attributes = {}
  variantDialogVisible.value = true
}

// Edit existing variant
const handleEditVariant = (variant: VariantFormData) => {
  isEditVariant.value = true
  variantForm.id = variant.id
  variantForm.code = variant.code
  variantForm.price = variant.price
  variantForm.currency = variant.currency
  variantForm.stock = variant.stock
  variantForm.safety_stock = variant.safety_stock
  variantForm.pre_sale_enabled = variant.pre_sale_enabled
  variantForm.attributes = { ...variant.attributes }
  variantDialogVisible.value = true
}

// Handle variant success — refresh product (updates productForm.sku), variants list, and inventory
const handleVariantSuccess = async () => {
  await loadProduct()
  variantsTabRef.value?.loadVariants()
  await nextTick()
  inventoryTabRef.value?.loadInventoryData()
}

const beforeUnloadHandler = (e: BeforeUnloadEvent) => {
  if (isDirty.value) {
    e.preventDefault()
    e.returnValue = ''
  }
}

// Initialize
onMounted(() => {
  loadProduct()
  loadProductMarkets()
  loadMarkets()
  loadCategories()
  loadBrands()
  loadWarehouses()
  window.addEventListener('beforeunload', beforeUnloadHandler)
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', beforeUnloadHandler)
})

onBeforeRouteLeave(async (_to, _from, next) => {
  if (!isDirty.value) {
    next()
    return
  }
  try {
    await ElMessageBox.confirm(
      t('products.unsavedChangesConfirm'),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    next()
  } catch {
    next(false)
  }
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

.status-tag {
  cursor: pointer;
  user-select: none;
}

.status-tag:hover {
  opacity: 0.85;
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
