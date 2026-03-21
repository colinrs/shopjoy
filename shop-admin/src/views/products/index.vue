<template>
  <div class="products-page">
    <!-- Market Filter Bar -->
    <el-card class="market-filter-card" shadow="never">
      <div class="market-filter-bar">
        <el-radio-group v-model="selectedMarket" @change="handleMarketChange">
          <el-radio-button value="">All Markets</el-radio-button>
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
            placeholder="搜索商品名称/编号"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="filterStatus" placeholder="商品状态" clearable class="filter-select">
            <el-option label="全部" value="" />
            <el-option label="在售" value="on_sale" />
            <el-option label="下架" value="off_sale" />
            <el-option label="草稿" value="draft" />
          </el-select>
          <el-select v-model="filterCategory" placeholder="商品分类" clearable class="filter-select">
            <el-option label="数码电子" value="electronics" />
            <el-option label="服装配饰" value="clothing" />
            <el-option label="家居生活" value="home" />
            <el-option label="运动户外" value="sports" />
          </el-select>
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增商品
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Bulk Actions -->
    <div class="bulk-actions" v-if="selectedProducts.length > 0">
      <span class="selected-count">已选择 {{ selectedProducts.length }} 项</span>
      <el-button size="small" @click="handleBatchOnSale">批量上架</el-button>
      <el-button size="small" @click="handleBatchOffSale">批量下架</el-button>
      <el-button size="small" type="success" @click="handleBatchPushToMarket">Push to Market</el-button>
      <el-button size="small" type="danger" @click="handleBatchDelete">批量删除</el-button>
    </div>

    <!-- Products Table -->
    <el-card class="table-card" shadow="never">
      <el-table
        :data="productList"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        stripe
      >
        <el-table-column type="selection" width="50" />
        <el-table-column label="商品信息" min-width="300">
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
                <p class="product-sku">SKU: {{ row.sku || '暂无' }}</p>
                <div class="product-tags">
                  <el-tag v-if="row.tags?.includes('hot')" size="small" type="danger" effect="plain">热销</el-tag>
                  <el-tag v-if="row.tags?.includes('new')" size="small" type="success" effect="plain">新品</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="价格" width="150" align="right">
          <template #default="{ row }">
            <div class="price-cell">
              <p class="sale-price">{{ row.currency || 'USD' }}{{ formatPrice(row.price) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStockType(row.stock)" size="small">
              {{ row.stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Markets" min-width="150" align="center">
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
                No markets
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="分类ID" width="100" align="center">
          <template #default="{ row }">
            {{ row.category_id || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'on_sale'"
              :inactive-value="'off_sale'"
              @change="(val) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" link size="small" @click="handleView(row)">
              预览
            </el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button type="primary" link size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="copy">复制</el-dropdown-item>
                  <el-dropdown-item command="top" divided>置顶</el-dropdown-item>
                  <el-dropdown-item command="delete" type="danger">删除</el-dropdown-item>
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
      :title="isEdit ? '编辑商品' : '新增商品'"
      width="800px"
      destroy-on-close
    >
      <el-form :model="productForm" label-width="100px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="商品名称" prop="name">
              <el-input v-model="productForm.name" placeholder="请输入商品名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商品价格" prop="price">
              <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="原始价格">
              <el-input-number v-model="productForm.original_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="库存数量" prop="stock">
              <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商品分类" prop="category">
              <el-select v-model="productForm.category" placeholder="请选择分类" style="width: 100%">
                <el-option label="数码电子" value="electronics" />
                <el-option label="服装配饰" value="clothing" />
                <el-option label="家居生活" value="home" />
                <el-option label="运动户外" value="sports" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="商品图片">
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
            <el-form-item label="商品描述">
              <el-input v-model="productForm.description" type="textarea" rows="4" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">保存</el-button>
      </template>
    </el-dialog>

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
              v-for="market in availableMarkets"
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
          <div class="price-note">Note: The price is in base currency (USD). It will be applied to all selected markets regardless of their local currency.</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushToMarketDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleConfirmPushToMarket" :loading="pushToMarketLoading">
          Push to Market
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Download, Picture, ArrowDown } from '@element-plus/icons-vue'
import { getProductList, pushToMarket, type Product, type ListProductsParams } from '@/api/product'
import { getMarkets, type Market } from '@/api/market'

const loading = ref(false)
const saveLoading = ref(false)
const pushToMarketLoading = ref(false)
const dialogVisible = ref(false)
const pushToMarketDialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const filterStatus = ref('')
const filterCategory = ref('')
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
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入商品价格', trigger: 'blur' }],
  stock: [{ required: true, message: '请输入库存数量', trigger: 'blur' }],
  category: [{ required: true, message: '请选择商品分类', trigger: 'change' }]
}

const formatPrice = (price: number) => {
  return (price / 100).toFixed(2)
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

const handleExport = () => {
  ElMessage.success('导出成功')
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
  isEdit.value = true
  Object.assign(productForm, { ...row })
  dialogVisible.value = true
}

const handleView = (row: any) => {
  ElMessage.info('预览商品: ' + row.name)
}

const handleCommand = (cmd: string, row: any) => {
  switch (cmd) {
    case 'copy':
      ElMessage.success('已复制商品')
      break
    case 'top':
      ElMessage.success('已置顶商品')
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确认删除商品 "${row.name}"?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
  })
}

const handleStatusChange = (row: any, val: string) => {
  const statusText = val === 'on_sale' ? '上架' : '下架'
  ElMessage.success(`商品已${statusText}`)
}

const handleSelectionChange = (selection: any[]) => {
  selectedProducts.value = selection
}

const handleBatchOnSale = () => {
  ElMessage.success(`已批量上架 ${selectedProducts.value.length} 个商品`)
}

const handleBatchOffSale = () => {
  ElMessage.success(`已批量下架 ${selectedProducts.value.length} 个商品`)
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确认删除选中的 ${selectedProducts.value.length} 个商品?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('批量删除成功')
  })
}

const handleBatchPushToMarket = () => {
  if (selectedProducts.value.length === 0) {
    ElMessage.warning('Please select at least one product')
    return
  }

  // Reset form
  pushToMarketForm.selectedMarkets = []
  pushToMarketForm.price = 0
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
    ElMessage.success(`Push to market completed. Success: ${successCount}, Failed: ${failCount}`)
    loadProducts()
  } catch (error) {
    console.error('Failed to push to market:', error)
    ElMessage.error('Failed to push to market')
  } finally {
    pushToMarketLoading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate((valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      setTimeout(() => {
        saveLoading.value = false
        dialogVisible.value = false
        ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
      }, 1000)
    }
  })
}

const handleImageChange = (file: any) => {
  productForm.image = URL.createObjectURL(file.raw)
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
    ElMessage.error('Failed to load markets')
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
    if (selectedMarket.value) {
      params.market_id = selectedMarket.value
    }

    const response = await getProductList(params)
    productList.value = response.list || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load products:', error)
    ElMessage.error('Failed to load products')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMarkets()
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
}

.market-filter-bar {
  display: flex;
  justify-content: flex-start;
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
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

.filter-select {
  width: 140px;
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
  padding: 12px 16px;
  background: #ECFDF5;
  border-radius: 8px;
  margin-bottom: 16px;
}

.selected-count {
  font-size: 14px;
  color: #059669;
  font-weight: 500;
}

/* Table Card */
.table-card {
  margin-bottom: 20px;
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
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
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

.product-details {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-weight: 500;
  color: #111827;
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
  font-weight: 600;
  color: #EF4444;
  margin: 0;
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
  border-top: 1px solid #E5E7EB;
  margin-top: 20px;
}

/* Upload */
.avatar-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
  width: 178px;
  height: 178px;
}

.avatar-uploader:hover {
  border-color: var(--el-color-primary);
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
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
}
</style>
