<template>
  <div class="products-page">
    <!-- Market Filter Bar -->
    <el-card class="market-filter-card" shadow="never">
      <div class="market-filter-bar">
        <el-radio-group v-model="selectedMarket" @change="handleMarketChange">
          <el-radio-button value="">{{ $t('common.all') }}</el-radio-button>
          <el-radio-button
            v-for="market in markets"
            :key="market.id"
            :value="market.id"
          >
            {{ market.code }}
          </el-radio-button>
        </el-radio-group>
      </div>
    </el-card>

    <!-- Search & Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('products.searchPlaceholder')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="filterStatus" :placeholder="$t('products.status')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('common.all')" value="" />
            <el-option :label="$t('products.onSale')" value="on_sale" />
            <el-option :label="$t('products.offSale')" value="off_sale" />
            <el-option :label="$t('products.draft')" value="draft" />
          </el-select>
          <el-cascader
            v-model="filterCategory"
            :options="categoryOptions"
            :props="{ checkStrictly: true, emitPath: false, value: 'id', label: 'name' }"
            :placeholder="$t('products.category')"
            clearable
            class="filter-select"
            @change="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>{{ $t('common.search') }}
          </el-button>
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>{{ $t('common.export') }}
          </el-button>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>{{ $t('products.addProduct') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Bulk Actions -->
    <div class="bulk-actions" v-if="selectedProducts.length > 0">
      <span class="selected-count">{{ $t('products.selectedCount', { count: selectedProducts.length }) }}</span>
      <el-button size="small" @click="handleBatchOnSale">{{ $t('products.batchOnSale') }}</el-button>
      <el-button size="small" @click="handleBatchOffSale">{{ $t('products.batchOffSale') }}</el-button>
      <el-button size="small" type="success" @click="handleBatchPushToMarket">{{ $t('products.pushToMarket') }}</el-button>
      <el-button size="small" type="danger" @click="handleBatchDelete">{{ $t('products.batchDelete') }}</el-button>
    </div>

    <!-- Products Table -->
    <el-card class="table-card" shadow="never">
      <!-- Skeleton loading -->
      <TableSkeleton v-if="loading && productList.length === 0" :rows="10" :columns="7" />

      <!-- Actual table -->
      <el-table
        v-else
        :data="productList"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        stripe
      >
        <el-table-column type="selection" width="50" />
        <el-table-column :label="$t('products.productInfo')" min-width="300">
          <template #default="{ row }">
            <div class="product-cell">
              <el-image
                :src="row.images?.[0] || defaultImage"
                fit="cover"
                class="product-thumb"
                :preview-src-list="row.images"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="product-details">
                <p class="product-name">{{ row.name }}</p>
                <p class="product-sku">SKU: {{ row.sku || $t('common.noData') }}</p>
                <div class="product-tags">
                  <el-tag v-if="row.tags?.includes('hot')" size="small" type="danger" effect="plain">{{ $t('products.hot') }}</el-tag>
                  <el-tag v-if="row.tags?.includes('new')" size="small" type="success" effect="plain">{{ $t('products.new') }}</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.price')" width="150" align="right">
          <template #default="{ row }">
            <div class="price-cell">
              <p class="sale-price">{{ row.currency || 'USD' }}{{ formatPrice(row.price) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="stock" :label="$t('products.stock')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStockType(row.stock)" size="small">
              {{ row.stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.markets')" min-width="150" align="center">
          <template #default="{ row }">
            <div class="market-tags">
              <el-tag
                v-for="market in row.markets"
                :key="market.market_id"
                :type="market.is_enabled ? 'success' : 'info'"
                size="small"
                class="market-tag"
              >
                {{ market.market_code }}
              </el-tag>
              <span v-if="!row.markets || row.markets.length === 0" class="no-markets">
                {{ $t('products.noMarkets') }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('products.categoryId')" width="100" align="center">
          <template #default="{ row }">
            {{ row.category_id || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('products.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'on_sale'"
              :inactive-value="'off_sale'"
              @change="(val: string) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              {{ $t('common.edit') }}
            </el-button>
            <el-button type="primary" link size="small" @click="handleView(row)">
              {{ $t('products.preview') }}
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button type="primary" link size="small">
                {{ $t('common.more') }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="copy">{{ $t('products.copy') }}</el-dropdown-item>
                  <el-dropdown-item command="top" divided>{{ $t('products.setTop') }}</el-dropdown-item>
                  <el-dropdown-item command="delete" type="danger">{{ $t('common.delete') }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? $t('products.editProduct') : $t('products.addProduct')"
      width="800px"
      destroy-on-close
    >
      <el-form :model="productForm" label-width="100px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item :label="$t('products.productName')" prop="name">
              <el-input v-model="productForm.name" :placeholder="$t('products.enterProductName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.productPrice')" prop="price">
              <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.originalPrice')">
              <el-input-number v-model="productForm.original_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.stockQuantity')" prop="stock">
              <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.productCategory')" prop="category_id">
              <el-cascader
                v-model="productForm.category_id"
                :options="categoryOptions"
                :props="{ checkStrictly: true, emitPath: false, value: 'id', label: 'name' }"
                :placeholder="$t('products.selectCategory')"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('products.productImage')">
              <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :auto-upload="false"
                :on-change="handleImageChange"
              >
                <img v-if="productForm.image" :src="productForm.image" class="avatar" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('products.productDescription')">
              <el-input v-model="productForm.description" type="textarea" rows="4" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Push to Market Dialog -->
    <el-dialog
      v-model="pushToMarketDialogVisible"
      :title="$t('products.pushToMarket')"
      width="500px"
      destroy-on-close
    >
      <el-form :model="pushToMarketForm" label-width="120px" ref="pushToMarketFormRef">
        <el-form-item :label="$t('products.markets')" prop="markets" required>
          <el-checkbox-group v-model="pushToMarketForm.selectedMarkets">
            <el-checkbox
              v-for="market in availableMarkets"
              :key="market.id"
              :value="market.id"
              :label="market.id"
            >
              {{ market.code }} - {{ market.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="$t('products.priceUSD')" prop="price" required>
          <el-input-number
            v-model="pushToMarketForm.price"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
          <div class="price-note">{{ $t('products.priceNote') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushToMarketDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleConfirmPushToMarket" :loading="pushToMarketLoading">
          {{ $t('products.pushToMarket') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Download, Picture, ArrowDown } from '@element-plus/icons-vue'
import { getProductList, pushToMarket, putOnSale, takeOffSale, createProduct, deleteProduct, type Product, type ListProductsParams } from '@/api/product'
import { getMarkets, type Market } from '@/api/market'
import { getCategoryTree, type CategoryTree } from '@/api/category'
import { uploadImage } from '@/api/upload'
import { TableSkeleton } from '@/components/skeleton'
import { t } from '@/plugins/i18n'

const router = useRouter()

const loading = ref(false)
const saveLoading = ref(false)
const pushToMarketLoading = ref(false)
const dialogVisible = ref(false)
const pushToMarketDialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const filterStatus = ref('')
const filterCategory = ref<number | ''>('')
const selectedMarket = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedProducts = ref<Product[]>([])
const formRef = ref()
const pushToMarketFormRef = ref()
const defaultImage = 'https://via.placeholder.com/80'

// Market data
const markets = ref<Market[]>([])
const productList = ref<Product[]>([])

// Category data
const categoryOptions = ref<CategoryTree[]>([])

// Push to market form
const pushToMarketForm = reactive({
  selectedMarkets: [] as number[],
  price: 0
})

// Available markets for push to market dialog
const availableMarkets = computed(() => markets.value.filter(m => m.is_active))

const productForm = reactive({
  id: null as number | null,
  name: '',
  price: 0,
  original_price: 0,
  stock: 0,
  category: '',
  category_id: 0,
  image: '',
  images: [] as string[],
  description: '',
  sku: '',
  brand: '',
  tags: [] as string[],
  is_matrix_product: false,
  hs_code: '',
  coo: '',
  weight: '',
  weight_unit: 'kg',
  length: '',
  width: '',
  height: '',
  dangerous_goods: [] as string[],
  currency: 'USD',
  cost_price: 0,
  status: 'draft'
})

const formRules = {
  name: [{ required: true, message: '', trigger: 'blur' }],
  price: [{ required: true, message: '', trigger: 'blur' }],
  stock: [{ required: true, message: '', trigger: 'blur' }],
  category: [{ required: true, message: '', trigger: 'change' }]
}

const formatPrice = (price: string) => {
  return price
}

const getStockType = (stock: number) => {
  if (stock === 0) return 'danger'
  if (stock < 20) return 'warning'
  return 'success'
}

const handleSearch = () => {
  currentPage.value = 1
  loadProducts()
}

const handleExport = async () => {
  try {
    loading.value = true
    // Build export params from current filters
    const params: Record<string, any> = {
      page: 1,
      page_size: 10000 // Export all matching records
    }
    if (searchQuery.value) {
      params.name = searchQuery.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterCategory.value) {
      params.category_id = filterCategory.value
    }
    if (selectedMarket.value) {
      params.market_id = selectedMarket.value
    }

    // Use window.open for export since there's no dedicated export API
    const queryString = new URLSearchParams(params).toString()
    const exportUrl = `/api/v1/products/export?${queryString}`
    window.open(exportUrl, '_blank')
    ElMessage.success(t('products.exporting'))
  } catch (error) {
    ElMessage.error(t('products.exportFailed'))
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(productForm, {
    id: null,
    name: '',
    price: 0,
    original_price: 0,
    stock: 0,
    category: '',
    image: '',
    description: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  router.push(`/products/${row.id}`)
}

const handleView = (row: any) => {
  router.push(`/products/${row.id}`)
}

const handleCommand = async (cmd: string, row: any) => {
  switch (cmd) {
    case 'copy':
      try {
        await navigator.clipboard.writeText(JSON.stringify(row, null, 2))
        ElMessage.success(t('products.copiedToClipboard'))
      } catch (error) {
        console.error('Failed to copy:', error)
        ElMessage.error(t('products.copyFailed'))
      }
      break
    case 'top':
      ElMessage.success(t('products.setTopSuccess'))
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDelete', { name: row.name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteProduct(row.id)
    ElMessage.success(t('products.deleteSuccess'))
    loadProducts()
  } catch (error) {
    console.error('Failed to delete product:', error)
    if (error !== 'cancel') {
      ElMessage.error(t('products.deleteFailed'))
    }
  }
}

const handleStatusChange = async (row: Product, val: string) => {
  const statusKey = val === 'on_sale' ? 'products.onSaleSuccess' : 'products.offSaleSuccess'
  try {
    if (val === 'on_sale') {
      await putOnSale(row.id)
    } else {
      await takeOffSale(row.id)
    }
    ElMessage.success(t(statusKey))
    loadProducts()
  } catch (error) {
    console.error('Failed to update status:', error)
    const errorKey = val === 'on_sale' ? 'products.onSaleFailed' : 'products.offSaleFailed'
    ElMessage.error(t(errorKey))
    // Revert the switch
    row.status = val === 'on_sale' ? 'off_sale' : 'on_sale'
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedProducts.value = selection
}

const handleBatchOnSale = async () => {
  try {
    const promises = selectedProducts.value.map(p => putOnSale(p.id))
    await Promise.all(promises)
    ElMessage.success(t('products.batchOnSaleSuccess', { count: selectedProducts.value.length }))
    loadProducts()
  } catch (error) {
    console.error('Failed to batch put on sale:', error)
    ElMessage.error(t('products.batchOnSaleFailed'))
  }
}

const handleBatchOffSale = async () => {
  try {
    const promises = selectedProducts.value.map(p => takeOffSale(p.id))
    await Promise.all(promises)
    ElMessage.success(t('products.batchOffSaleSuccess', { count: selectedProducts.value.length }))
    loadProducts()
  } catch (error) {
    console.error('Failed to batch take off sale:', error)
    ElMessage.error(t('products.batchOffSaleFailed'))
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmBatchDelete', { count: selectedProducts.value.length }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const deletePromises = selectedProducts.value.map(p => deleteProduct(p.id))
    const results = await Promise.allSettled(deletePromises)

    const failedCount = results.filter(r => r.status === 'rejected').length
    const successCount = results.filter(r => r.status === 'fulfilled').length

    if (failedCount > 0) {
      ElMessage.error(t('products.batchDeletePartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('products.batchDeleteSuccess', { count: successCount }))
    }

    selectedProducts.value = []
    loadProducts()
  } catch (error) {
    console.error('Failed to batch delete:', error)
    if ((error as any)?.message !== 'cancel') {
      ElMessage.error(t('products.batchDeleteFailed'))
    }
  }
}

const handleBatchPushToMarket = () => {
  if (selectedProducts.value.length === 0) {
    ElMessage.warning(t('products.selectAtLeastOneProduct'))
    return
  }

  // Reset form
  pushToMarketForm.selectedMarkets = []
  pushToMarketForm.price = 0
  pushToMarketDialogVisible.value = true
}

const handleConfirmPushToMarket = async () => {
  if (pushToMarketForm.selectedMarkets.length === 0) {
    ElMessage.warning(t('products.selectAtLeastOneMarket'))
    return
  }
  if (pushToMarketForm.price <= 0) {
    ElMessage.warning(t('products.enterValidPrice'))
    return
  }

  pushToMarketLoading.value = true
  try {
    const prices = pushToMarketForm.selectedMarkets.map(() =>
      pushToMarketForm.price.toFixed(2)
    )

    let successCount = 0
    let failCount = 0

    for (const product of selectedProducts.value) {
      try {
        const result = await pushToMarket(product.id, {
          market_ids: pushToMarketForm.selectedMarkets,
          prices
        })
        successCount += result.success?.length || 0
        failCount += result.failed?.length || 0
      } catch (error) {
        failCount += pushToMarketForm.selectedMarkets.length
      }
    }

    pushToMarketDialogVisible.value = false
    ElMessage.success(t('products.pushToMarketSuccess', { success: successCount, failed: failCount }))
    loadProducts()
  } catch (error) {
    console.error('Failed to push to market:', error)
    ElMessage.error(t('products.pushToMarketFailed'))
  } finally {
    pushToMarketLoading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        await createProduct({
          name: productForm.name,
          description: productForm.description,
          price: String(productForm.price),
          currency: 'USD',
          category_id: productForm.category_id || 1,
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
        dialogVisible.value = false
        ElMessage.success(t('products.addSuccess'))
        loadProducts()
      } catch (error) {
        console.error('Failed to create product:', error)
        ElMessage.error(t('products.addFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

const handleImageChange = async (file: any) => {
  try {
    const response = await uploadImage(file.raw, 'product')
    productForm.image = response.url
  } catch (error) {
    console.error('Upload failed:', error)
    ElMessage.error(t('products.imageUploadFailed'))
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadProducts()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadProducts()
}

const handleMarketChange = () => {
  currentPage.value = 1
  loadProducts()
}

const loadMarkets = async () => {
  try {
    const response = await getMarkets()
    markets.value = response.list || []
  } catch (error) {
    console.error('Failed to load markets:', error)
    ElMessage.error(t('products.loadMarketsFailed'))
  }
}

const loadCategories = async () => {
  try {
    const response = await getCategoryTree()
    categoryOptions.value = response || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    ElMessage.error(t('products.loadCategoriesFailed'))
  }
}

const loadProducts = async () => {
  loading.value = true
  try {
    const params: ListProductsParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (searchQuery.value) {
      params.name = searchQuery.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterCategory.value) {
      params.category_id = filterCategory.value
    }
    if (selectedMarket.value) {
      params.market_id = selectedMarket.value
    }

    const response = await getProductList(params)
    productList.value = response.list || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load products:', error)
    ElMessage.error(t('products.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMarkets()
  loadCategories()
  loadProducts()
})
</script>

<style scoped>
.products-page {
  padding: 0;
}

/* Market Filter Bar */
.market-filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.market-filter-bar {
  display: flex;
  justify-content: flex-start;
}

.market-filter-bar :deep(.el-radio-button__inner) {
  border-radius: 10px !important;
  border: 1px solid #E5E7EB;
  padding: 8px 18px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.market-filter-bar :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  border-color: #6366F1;
  box-shadow: 0 4px 10px -2px rgba(99, 102, 241, 0.3);
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Bulk Actions */
.bulk-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  margin-bottom: 16px;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.selected-count {
  font-size: 14px;
  color: #6366F1;
  font-weight: 600;
}

/* Table Card */
.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Product Cell */
.product-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.product-thumb {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid #E5E7EB;
  transition: transform 0.2s ease;
}

.product-thumb:hover {
  transform: scale(1.05);
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.product-details {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-size: 14px;
  line-height: 1.4;
}

.product-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0 0 8px 0;
  font-family: 'Fira Code', monospace;
}

.product-tags {
  display: flex;
  gap: 8px;
}

/* Price Cell */
.price-cell {
  text-align: right;
}

.sale-price {
  font-size: 16px;
  font-weight: 700;
  color: #EF4444;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.original-price {
  font-size: 12px;
  color: #9CA3AF;
  text-decoration: line-through;
  margin: 4px 0 0 0;
}

/* Market Tags */
.market-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  justify-content: center;
}

.market-tag {
  font-size: 11px;
  border-radius: 6px;
}

.no-markets {
  color: #9CA3AF;
  font-size: 12px;
}

.price-note {
  font-size: 12px;
  color: #F59E0B;
  margin-top: 8px;
  line-height: 1.4;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Dialog Styling */
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

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Upload */
.avatar-uploader {
  border: 2px dashed #E5E7EB;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.2s ease;
  width: 178px;
  height: 178px;
  background: #F9FAFB;
}

.avatar-uploader:hover {
  border-color: #6366F1;
  background: #F5F3FF;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #9CA3AF;
  width: 178px;
  height: 178px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: cover;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .product-cell {
    flex-direction: column;
    align-items: flex-start;
  }

  .product-thumb {
    width: 60px;
    height: 60px;
  }

  .market-filter-card,
  .filter-card,
  .table-card {
    border-radius: 14px;
  }
}
</style>
