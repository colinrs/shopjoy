<template>
  <div class="products-page">
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
                :src="row.image || defaultImage"
                fit="cover"
                class="product-thumb"
                :preview-src-list="[row.image]"
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
                  <el-tag v-if="row.is_hot" size="small" type="danger" effect="plain">热销</el-tag>
                  <el-tag v-if="row.is_new" size="small" type="success" effect="plain">新品</el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="价格" width="150" align="right">
          <template #default="{ row }">
            <div class="price-cell">
              <p class="sale-price">¥{{ formatPrice(row.price) }}</p>
              <p v-if="row.original_price" class="original-price">
                ¥{{ formatPrice(row.original_price) }}
              </p>
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
        <el-table-column prop="category" label="分类" width="120" />
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
        <el-table-column prop="sales" label="销量" width="100" align="center" sortable />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Download, Picture, ArrowDown } from '@element-plus/icons-vue'

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const filterStatus = ref('')
const filterCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(100)
const selectedProducts = ref<any[]>([])
const formRef = ref()
const defaultImage = 'https://via.placeholder.com/80'

const productForm = reactive({
  id: null,
  name: '',
  price: 0,
  original_price: 0,
  stock: 0,
  category: '',
  image: '',
  description: ''
})

const formRules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入商品价格', trigger: 'blur' }],
  stock: [{ required: true, message: '请输入库存数量', trigger: 'blur' }],
  category: [{ required: true, message: '请选择商品分类', trigger: 'change' }]
}

// Mock data
const productList = ref([
  {
    id: 1,
    name: '无线蓝牙耳机 Pro',
    sku: 'BT-2024-001',
    price: 29900,
    original_price: 39900,
    stock: 156,
    category: '数码电子',
    status: 'on_sale',
    sales: 328,
    is_hot: true,
    is_new: false,
    image: ''
  },
  {
    id: 2,
    name: '智能手表 Series 7',
    sku: 'SW-2024-002',
    price: 199900,
    original_price: 229900,
    stock: 89,
    category: '数码电子',
    status: 'on_sale',
    sales: 156,
    is_hot: true,
    is_new: true,
    image: ''
  },
  {
    id: 3,
    name: '便携充电宝 20000mAh',
    sku: 'PB-2024-003',
    price: 12900,
    original_price: 15900,
    stock: 234,
    category: '数码电子',
    status: 'on_sale',
    sales: 892,
    is_hot: false,
    is_new: false,
    image: ''
  },
  {
    id: 4,
    name: '机械键盘 RGB 青轴',
    sku: 'KB-2024-004',
    price: 45900,
    original_price: 59900,
    stock: 12,
    category: '数码电子',
    status: 'on_sale',
    sales: 67,
    is_hot: false,
    is_new: true,
    image: ''
  },
  {
    id: 5,
    name: '4K 高清显示器 27寸',
    sku: 'MN-2024-005',
    price: 159900,
    stock: 0,
    category: '数码电子',
    status: 'off_sale',
    sales: 23,
    is_hot: false,
    is_new: false,
    image: ''
  }
])

const formatPrice = (price: number) => {
  return (price / 100).toFixed(2)
}

const getStockType = (stock: number) => {
  if (stock === 0) return 'danger'
  if (stock < 20) return 'warning'
  return 'success'
}

const handleSearch = () => {
  ElMessage.info('搜索: ' + searchQuery.value)
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

const loadProducts = async () => {
  loading.value = true
  try {
    // TODO: Call API
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.products-page {
  padding: 0;
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
