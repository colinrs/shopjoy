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
          <el-select v-model="filterStatus" placeholder="商品状态" clearable class="filter-select" @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="在售" value="on_sale" />
            <el-option label="下架" value="off_sale" />
            <el-option label="草稿" value="draft" />
          </el-select>
          <el-select v-model="filterCategory" placeholder="商品分类" clearable class="filter-select" @change="handleSearch">
            <el-option label="数码电子" value="electronics" />
            <el-option label="服装配饰" value="clothing" />
            <el-option label="家居生活" value="home" />
            <el-option label="运动户外" value="sports" />
          </el-select>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>查询
          </el-button>
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Download, Picture, ArrowDown } from '@element-plus/icons-vue'
import { getProductList, pushToMarket, putOnSale, takeOffSale, createProduct, type Product, type ListProductsParams } from '@/api/product'
import { getMarkets, type Market } from '@/api/market'
import { TableSkeleton } from '@/components/skeleton'

const router = useRouter()

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
  router.push(`/products/${row.id}`)
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

const handleStatusChange = async (row: Product, val: string) => {
  const statusText = val === 'on_sale' ? '上架' : '下架'
  try {
    if (val === 'on_sale') {
      await putOnSale(row.id)
    } else {
      await takeOffSale(row.id)
    }
    ElMessage.success(`商品已${statusText}`)
    loadProducts()
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error(`商品${statusText}失败`)
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
    ElMessage.success(`已批量上架 ${selectedProducts.value.length} 个商品`)
    loadProducts()
  } catch (error) {
    console.error('Failed to batch put on sale:', error)
    ElMessage.error('批量上架失败')
  }
}

const handleBatchOffSale = async () => {
  try {
    const promises = selectedProducts.value.map(p => takeOffSale(p.id))
    await Promise.all(promises)
    ElMessage.success(`已批量下架 ${selectedProducts.value.length} 个商品`)
    loadProducts()
  } catch (error) {
    console.error('Failed to batch take off sale:', error)
    ElMessage.error('批量下架失败')
  }
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

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        const categoryMap: Record<string, number> = {
          electronics: 1,
          clothing: 2,
          home: 3,
          sports: 4
        }

        // Price in cents (backend expects int64)
        const priceInCents = Math.round(productForm.price * 100)

        await createProduct({
          name: productForm.name,
          description: productForm.description,
          price: priceInCents,
          currency: 'USD',
          category_id: categoryMap[productForm.category] || 1,
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
        ElMessage.success('添加成功')
        loadProducts()
      } catch (error) {
        console.error('Failed to create product:', error)
        ElMessage.error('添加失败')
      } finally {
        saveLoading.value = false
      }
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
    if (filterCategory.value) {
      // Map category string to ID
      const categoryMap: Record<string, number> = {
        electronics: 1,
        clothing: 2,
        home: 3,
        sports: 4
      }
      params.category_id = categoryMap[filterCategory.value] || 0
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
