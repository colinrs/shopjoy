<template>
  <div class="brands-page">
    <!-- Search & Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索品牌名称"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="filterStatus" placeholder="状态" clearable class="filter-select" @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>查询
          </el-button>
        </div>
        <div class="filter-right">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增品牌
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Brands Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="brandList" v-loading="loading" stripe>
        <el-table-column label="品牌信息" min-width="250">
          <template #default="{ row }">
            <div class="brand-cell">
              <el-image
                v-if="row.logo"
                :src="row.logo"
                fit="cover"
                class="brand-logo"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div v-else class="brand-logo-placeholder">
                <el-icon><Picture /></el-icon>
              </div>
              <div class="brand-details">
                <p class="brand-name">{{ row.name }}</p>
                <p class="brand-website" v-if="row.website">{{ row.website }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="product_count" label="商品数量" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.product_count > 0 ? 'success' : 'info'" size="small">
              {{ row.product_count || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enable_page" label="品牌专区" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.enable_page"
              @change="(val: boolean) => handleTogglePage(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="(val: number) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
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
      :title="isEdit ? '编辑品牌' : '新增品牌'"
      width="700px"
      destroy-on-close
    >
      <el-form :model="brandForm" label-width="120px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="品牌名称" prop="name">
              <el-input v-model="brandForm.name" placeholder="请输入品牌名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="brandForm.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="Logo URL">
              <el-input v-model="brandForm.logo" placeholder="品牌Logo URL" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="品牌官网">
              <el-input v-model="brandForm.website" placeholder="品牌官网URL" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="品牌描述">
              <el-input v-model="brandForm.description" type="textarea" rows="3" placeholder="品牌描述" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商标号">
              <el-input v-model="brandForm.trademark_number" placeholder="商标注册号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商标注册国">
              <el-input v-model="brandForm.trademark_country" placeholder="如: CN, US" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用品牌专区">
              <el-switch v-model="brandForm.enable_page" />
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
import { Search, Plus, Picture } from '@element-plus/icons-vue'
import {
  getBrands,
  createBrand,
  updateBrand,
  deleteBrand,
  updateBrandStatus,
  toggleBrandPage,
  type Brand
} from '@/api/brand'

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const brandList = ref<Brand[]>([])
const formRef = ref()

const brandForm = reactive({
  id: 0,
  name: '',
  logo: '',
  description: '',
  website: '',
  trademark_number: '',
  trademark_country: '',
  enable_page: false,
  sort: 0
})

const formRules = {
  name: [{ required: true, message: '请输入品牌名称', trigger: 'blur' }]
}

const handleSearch = () => {
  currentPage.value = 1
  loadBrands()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(brandForm, {
    id: 0,
    name: '',
    logo: '',
    description: '',
    website: '',
    trademark_number: '',
    trademark_country: '',
    enable_page: false,
    sort: 0
  })
  dialogVisible.value = true
}

const handleEdit = (row: Brand) => {
  isEdit.value = true
  Object.assign(brandForm, {
    id: row.id,
    name: row.name,
    logo: row.logo || '',
    description: row.description || '',
    website: row.website || '',
    trademark_number: row.trademark_number || '',
    trademark_country: row.trademark_country || '',
    enable_page: row.enable_page,
    sort: row.sort || 0
  })
  dialogVisible.value = true
}

const handleDelete = (row: Brand) => {
  if (row.product_count > 0) {
    ElMessage.warning('该品牌下存在商品，无法删除')
    return
  }

  ElMessageBox.confirm(`确认删除品牌 "${row.name}"?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteBrand(row.id)
      ElMessage.success('删除成功')
      loadBrands()
    } catch (error) {
      console.error('Failed to delete brand:', error)
      ElMessage.error('删除失败')
    }
  })
}

const handleStatusChange = async (row: Brand, status: number) => {
  try {
    await updateBrandStatus(row.id, status)
    ElMessage.success(status === 1 ? '已启用' : '已禁用')
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error('更新状态失败')
    row.status = status === 1 ? 0 : 1
  }
}

const handleTogglePage = async (row: Brand, enabled: boolean) => {
  try {
    await toggleBrandPage(row.id, enabled)
    ElMessage.success(enabled ? '已启用品牌专区' : '已禁用品牌专区')
  } catch (error) {
    console.error('Failed to toggle page:', error)
    ElMessage.error('操作失败')
    row.enable_page = !enabled
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        if (isEdit.value) {
          await updateBrand({
            id: brandForm.id,
            name: brandForm.name,
            logo: brandForm.logo,
            description: brandForm.description,
            website: brandForm.website,
            trademark_number: brandForm.trademark_number,
            trademark_country: brandForm.trademark_country,
            enable_page: brandForm.enable_page,
            sort: brandForm.sort
          })
          ElMessage.success('更新成功')
        } else {
          await createBrand({
            name: brandForm.name,
            logo: brandForm.logo,
            description: brandForm.description,
            website: brandForm.website,
            trademark_number: brandForm.trademark_number,
            trademark_country: brandForm.trademark_country,
            enable_page: brandForm.enable_page,
            sort: brandForm.sort
          })
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadBrands()
      } catch (error) {
        console.error('Failed to save brand:', error)
        ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
      } finally {
        saveLoading.value = false
      }
    }
  })
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadBrands()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadBrands()
}

const loadBrands = async () => {
  loading.value = true
  try {
    const response = await getBrands({
      page: currentPage.value,
      page_size: pageSize.value,
      name: searchQuery.value || undefined,
      status: filterStatus.value !== '' ? filterStatus.value : undefined
    })
    brandList.value = response.list || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load brands:', error)
    ElMessage.error('加载品牌失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadBrands()
})
</script>

<style scoped>
.brands-page {
  padding: 0;
}

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

.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.brand-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo,
.brand-logo-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid #E5E7EB;
  transition: transform 0.2s ease;
}

.brand-logo:hover {
  transform: scale(1.05);
}

.brand-logo-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  color: #6366F1;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  color: #6366F1;
}

.brand-details {
  flex: 1;
  min-width: 0;
}

.brand-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-size: 14px;
}

.brand-website {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
  font-family: 'Fira Code', monospace;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Dialog */
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

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

@media (max-width: 768px) {
  .search-input {
    width: 100%;
  }
}
</style>