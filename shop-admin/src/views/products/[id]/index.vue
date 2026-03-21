<template>
  <div class="product-detail-page">
    <!-- Page Header -->
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-left">
          <el-button link @click="handleBack">
            <el-icon><ArrowLeft /></el-icon>
            Back to Products
          </el-button>
          <el-divider direction="vertical" />
          <h2 class="product-title">
            {{ product?.name || 'Loading...' }}
            <el-tag v-if="product" :type="getStatusType(product.status)" size="small">
              {{ getStatusText(product.status) }}
            </el-tag>
          </h2>
        </div>
        <div class="header-right">
          <el-button @click="handleSave" type="primary" :loading="saveLoading">
            <el-icon><Check /></el-icon>
            Save Changes
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Loading State -->
    <el-skeleton v-if="loading" :rows="10" animated />

    <!-- Tab Layout -->
    <el-card v-else class="tabs-card" shadow="never">
      <el-tabs v-model="activeTab" class="product-tabs">
        <!-- Basic Info Tab -->
        <el-tab-pane label="Basic Info" name="basic">
          <el-form :model="productForm" label-width="140px" :rules="formRules" ref="formRef">
            <!-- Basic Information Section -->
            <div class="form-section">
              <h3 class="section-title">Basic Information</h3>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Product Name" prop="name">
                    <el-input v-model="productForm.name" placeholder="Enter product name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="SKU" prop="sku">
                    <el-input v-model="productForm.sku" placeholder="Enter SKU" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Brand">
                    <el-input v-model="productForm.brand" placeholder="Enter brand name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Category ID">
                    <el-input-number v-model="productForm.category_id" :min="0" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="Description">
                    <el-input v-model="productForm.description" type="textarea" :rows="4" placeholder="Enter product description" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Price" prop="price">
                    <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Currency">
                    <el-select v-model="productForm.currency" style="width: 100%">
                      <el-option label="USD" value="USD" />
                      <el-option label="EUR" value="EUR" />
                      <el-option label="GBP" value="GBP" />
                      <el-option label="CNY" value="CNY" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Cost Price">
                    <el-input-number v-model="productForm.cost_price" :min="0" :precision="2" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Stock">
                    <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Status">
                    <el-select v-model="productForm.status" style="width: 100%">
                      <el-option label="Draft" value="draft" />
                      <el-option label="On Sale" value="on_sale" />
                      <el-option label="Off Sale" value="off_sale" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Is Matrix Product">
                    <el-switch v-model="productForm.is_matrix_product" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- Compliance Section -->
            <div class="form-section">
              <h3 class="section-title">Cross-Border Compliance</h3>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="HS Code">
                    <el-input v-model="productForm.hs_code" placeholder="Enter HS Code" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Country of Origin">
                    <el-input v-model="productForm.coo" placeholder="Enter country code (e.g., CN, US)" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="Weight">
                    <el-input v-model="productForm.weight" placeholder="e.g., 1.5">
                      <template #append>
                        <el-select v-model="productForm.weight_unit" style="width: 80px">
                          <el-option label="kg" value="kg" />
                          <el-option label="g" value="g" />
                          <el-option label="lb" value="lb" />
                        </el-select>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="16">
                  <el-form-item label="Dimensions (L x W x H)">
                    <el-row :gutter="8">
                      <el-col :span="8">
                        <el-input v-model="productForm.length" placeholder="Length">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                      <el-col :span="8">
                        <el-input v-model="productForm.width" placeholder="Width">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                      <el-col :span="8">
                        <el-input v-model="productForm.height" placeholder="Height">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                    </el-row>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="Dangerous Goods">
                    <el-checkbox-group v-model="productForm.dangerous_goods">
                      <el-checkbox label="battery">Battery</el-checkbox>
                      <el-checkbox label="liquid">Liquid</el-checkbox>
                      <el-checkbox label="flammable">Flammable</el-checkbox>
                      <el-checkbox label="magnetic">Magnetic</el-checkbox>
                      <el-checkbox label="fragile">Fragile</el-checkbox>
                    </el-checkbox-group>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- Images Section -->
            <div class="form-section">
              <h3 class="section-title">Product Images</h3>
              <el-form-item label="Image URLs">
                <div class="image-list">
                  <div v-for="(img, index) in productForm.images" :key="index" class="image-item">
                    <el-image :src="img" fit="cover" class="product-image">
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
                      @click="removeImage(index)"
                    >
                      <el-icon><Close /></el-icon>
                    </el-button>
                  </div>
                  <div class="add-image" @click="addImage">
                    <el-icon><Plus /></el-icon>
                    <span>Add Image</span>
                  </div>
                </div>
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <!-- Markets Tab -->
        <el-tab-pane label="Markets" name="markets">
          <div class="markets-section">
            <div class="section-header">
              <h3 class="section-title">Market Availability</h3>
              <el-button type="primary" @click="showPushToMarketDialog">
                <el-icon><Plus /></el-icon>
                Push to Market
              </el-button>
            </div>

            <el-table :data="productMarkets" v-loading="marketsLoading" stripe>
              <el-table-column label="Market" min-width="150">
                <template #default="{ row }">
                  <div class="market-cell">
                    <span class="market-code">{{ row.market_code }}</span>
                    <span class="market-name">{{ row.market_name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="Status" width="120" align="center">
                <template #default="{ row }">
                  <el-switch
                    v-model="row.is_enabled"
                    @change="(val) => handleMarketEnableChange(row, val)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="Price" width="180" align="right">
                <template #default="{ row }">
                  <div class="price-cell">
                    <el-input-number
                      v-model="row.price"
                      :min="0"
                      :precision="2"
                      :controls="false"
                      size="small"
                      style="width: 100px"
                      @change="() => handleMarketPriceChange(row)"
                    />
                    <span class="currency">{{ row.currency }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="Compare at Price" width="180" align="right">
                <template #default="{ row }">
                  <div class="price-cell">
                    <el-input-number
                      v-model="row.compare_at_price"
                      :min="0"
                      :precision="2"
                      :controls="false"
                      size="small"
                      style="width: 100px"
                      @change="() => handleMarketPriceChange(row)"
                    />
                    <span class="currency">{{ row.currency }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="Stock Alert" width="120" align="center">
                <template #default="{ row }">
                  <el-input-number
                    v-model="row.stock_alert_threshold"
                    :min="0"
                    :controls="false"
                    size="small"
                    style="width: 80px"
                    @change="() => handleMarketPriceChange(row)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="Published" width="120" align="center">
                <template #default="{ row }">
                  <span v-if="row.published_at">{{ formatDate(row.published_at) }}</span>
                  <span v-else class="text-muted">Not published</span>
                </template>
              </el-table-column>
              <el-table-column label="Actions" width="100" align="center">
                <template #default="{ row }">
                  <el-button
                    type="danger"
                    link
                    size="small"
                    @click="handleRemoveFromMarket(row)"
                  >
                    Remove
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <div v-if="productMarkets.length === 0 && !marketsLoading" class="empty-markets">
              <el-empty description="This product is not available in any market">
                <el-button type="primary" @click="showPushToMarketDialog">Push to Market</el-button>
              </el-empty>
            </div>
          </div>
        </el-tab-pane>

        <!-- Variants Tab (Placeholder) -->
        <el-tab-pane label="Variants" name="variants">
          <el-empty description="Variants management coming soon" />
        </el-tab-pane>

        <!-- Pricing Tab (Placeholder) -->
        <el-tab-pane label="Pricing" name="pricing">
          <el-empty description="Pricing management coming soon" />
        </el-tab-pane>

        <!-- Localization Tab (Placeholder) -->
        <el-tab-pane label="Localization" name="localization">
          <el-empty description="Localization management coming soon" />
        </el-tab-pane>

        <!-- Inventory Tab (Placeholder) -->
        <el-tab-pane label="Inventory" name="inventory">
          <el-empty description="Inventory management coming soon" />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Push to Market Dialog -->
    <el-dialog
      v-model="pushToMarketDialogVisible"
      title="Push to Market"
      width="500px"
      destroy-on-close
    >
      <el-form :model="pushToMarketForm" label-width="120px" ref="pushToMarketFormRef">
        <el-form-item label="Markets" prop="markets" required>
          <el-checkbox-group v-model="pushToMarketForm.selectedMarkets">
            <el-checkbox
              v-for="market in availableMarketsForPush"
              :key="market.id"
              :value="market.id"
              :label="market.id"
            >
              {{ market.code }} - {{ market.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="Price (USD)" prop="price" required>
          <el-input-number
            v-model="pushToMarketForm.price"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
          <div class="price-note">Note: The price is in base currency (USD). It will be applied to all selected markets.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushToMarketDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleConfirmPushToMarket" :loading="pushToMarketLoading">
          Push to Market
        </el-button>
      </template>
    </el-dialog>

    <!-- Add Image Dialog -->
    <el-dialog v-model="addImageDialogVisible" title="Add Image URL" width="400px">
      <el-input v-model="newImageUrl" placeholder="Enter image URL" />
      <template #footer>
        <el-button @click="addImageDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="confirmAddImage">Add</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Check, Plus, Picture, Close } from '@element-plus/icons-vue'
import {
  getProduct,
  updateProduct,
  getProductMarkets,
  updateProductMarket,
  pushToMarket,
  removeFromMarket,
  type Product,
  type ProductMarket
} from '@/api/product'
import { getMarkets, type Market } from '@/api/market'

const route = useRoute()
const router = useRouter()

const productId = computed(() => Number(route.params.id))

// State
const loading = ref(false)
const saveLoading = ref(false)
const marketsLoading = ref(false)
const pushToMarketLoading = ref(false)
const activeTab = ref('basic')
const product = ref<Product | null>(null)
const productMarkets = ref<ProductMarket[]>([])
const markets = ref<Market[]>([])
const pushToMarketDialogVisible = ref(false)
const addImageDialogVisible = ref(false)
const newImageUrl = ref('')
const formRef = ref()
const pushToMarketFormRef = ref()

// Form
const productForm = reactive({
  name: '',
  description: '',
  price: 0,
  currency: 'USD',
  cost_price: 0,
  stock: 0,
  status: 'draft',
  category_id: 0,
  sku: '',
  brand: '',
  tags: [] as string[],
  images: [] as string[],
  is_matrix_product: false,
  hs_code: '',
  coo: '',
  weight: '',
  weight_unit: 'kg',
  length: '',
  width: '',
  height: '',
  dangerous_goods: [] as string[]
})

// Push to market form
const pushToMarketForm = reactive({
  selectedMarkets: [] as number[],
  price: 0
})

// Available markets for push (excluding existing ones)
const availableMarketsForPush = computed(() => {
  const existingMarketIds = productMarkets.value.map(pm => pm.market_id)
  return markets.value.filter(m => m.is_active && !existingMarketIds.includes(m.id))
})

// Form rules
const formRules = {
  name: [{ required: true, message: 'Please enter product name', trigger: 'blur' }],
  price: [{ required: true, message: 'Please enter price', trigger: 'blur' }]
}

// Helper functions
const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    on_sale: 'success',
    off_sale: 'warning',
    draft: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    on_sale: 'On Sale',
    off_sale: 'Off Sale',
    draft: 'Draft'
  }
  return texts[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString()
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
    // Populate form
    Object.assign(productForm, {
      name: data.name || '',
      description: data.description || '',
      price: data.price || 0,
      currency: data.currency || 'USD',
      cost_price: data.cost_price || 0,
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
    console.error('Failed to load product:', error)
    ElMessage.error('Failed to load product')
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
    console.error('Failed to load product markets:', error)
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
    console.error('Failed to load markets:', error)
  }
}

// Save product
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
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
        ElMessage.success('Product updated successfully')
        loadProduct()
      } catch (error) {
        console.error('Failed to update product:', error)
        ElMessage.error('Failed to update product')
      } finally {
        saveLoading.value = false
      }
    }
  })
}

// Market actions
const handleMarketEnableChange = async (row: ProductMarket, enabled: boolean) => {
  try {
    await updateProductMarket(productId.value, row.market_id, { is_enabled: enabled })
    ElMessage.success(`Market ${enabled ? 'enabled' : 'disabled'}`)
    loadProductMarkets()
  } catch (error) {
    console.error('Failed to update market:', error)
    ElMessage.error('Failed to update market status')
    row.is_enabled = !enabled // Revert
  }
}

const handleMarketPriceChange = async (row: ProductMarket) => {
  try {
    await updateProductMarket(productId.value, row.market_id, {
      price: row.price,
      compare_at_price: row.compare_at_price,
      stock_alert_threshold: row.stock_alert_threshold
    })
    ElMessage.success('Market price updated')
  } catch (error) {
    console.error('Failed to update market price:', error)
    ElMessage.error('Failed to update market price')
  }
}

const handleRemoveFromMarket = async (row: ProductMarket) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to remove this product from ${row.market_name}?`,
      'Warning',
      {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
    await removeFromMarket(productId.value, row.market_id)
    ElMessage.success('Product removed from market')
    loadProductMarkets()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove from market:', error)
      ElMessage.error('Failed to remove from market')
    }
  }
}

// Push to market
const showPushToMarketDialog = () => {
  pushToMarketForm.selectedMarkets = []
  pushToMarketForm.price = productForm.price / 100 // Convert from cents
  pushToMarketDialogVisible.value = true
}

const handleConfirmPushToMarket = async () => {
  if (pushToMarketForm.selectedMarkets.length === 0) {
    ElMessage.warning('Please select at least one market')
    return
  }
  if (pushToMarketForm.price <= 0) {
    ElMessage.warning('Please enter a valid price')
    return
  }

  pushToMarketLoading.value = true
  try {
    const prices = pushToMarketForm.selectedMarkets.map(() =>
      pushToMarketForm.price.toFixed(2)
    )

    const result = await pushToMarket(productId.value, {
      market_ids: pushToMarketForm.selectedMarkets,
      prices
    })

    pushToMarketDialogVisible.value = false
    ElMessage.success(
      `Push to market completed. Success: ${result.success?.length || 0}, Failed: ${result.failed?.length || 0}`
    )
    loadProductMarkets()
  } catch (error) {
    console.error('Failed to push to market:', error)
    ElMessage.error('Failed to push to market')
  } finally {
    pushToMarketLoading.value = false
  }
}

// Image management
const addImage = () => {
  newImageUrl.value = ''
  addImageDialogVisible.value = true
}

const confirmAddImage = () => {
  if (newImageUrl.value.trim()) {
    productForm.images.push(newImageUrl.value.trim())
    addImageDialogVisible.value = false
  }
}

const removeImage = (index: number) => {
  productForm.images.splice(index, 1)
}

// Initialize
onMounted(() => {
  loadProduct()
  loadProductMarkets()
  loadMarkets()
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

/* Form Sections */
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

/* Markets Section */
.markets-section {
  padding: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.market-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.market-code {
  font-weight: 600;
  font-size: 14px;
}

.market-name {
  font-size: 12px;
  color: #6B7280;
}

.price-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

.currency {
  font-size: 12px;
  color: #6B7280;
  min-width: 36px;
}

.text-muted {
  color: #9CA3AF;
  font-size: 12px;
}

.empty-markets {
  padding: 40px 0;
}

.price-note {
  font-size: 12px;
  color: #F59E0B;
  margin-top: 8px;
  line-height: 1.4;
}

/* Image Management */
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