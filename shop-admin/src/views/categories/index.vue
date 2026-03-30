<template>
  <div class="categories-page">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2>{{ $t('categories.title') }}</h2>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleAddRoot">
            <el-icon><Plus /></el-icon>{{ $t('categories.addTopCategory') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Category Tree Table -->
    <el-card class="tree-card" shadow="never">
      <el-table
        :data="categoryTree"
        v-loading="loading"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="name" :label="$t('categories.name')" min-width="200" />
        <el-table-column prop="code" :label="$t('categories.categoryCode')" width="120" />
        <el-table-column prop="level" :label="$t('categories.level')" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small">L{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('categories.sort')" width="80" align="center" />
        <el-table-column prop="product_count" :label="$t('categories.productCount')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.product_count > 0 ? 'success' : 'info'" size="small">
              {{ row.product_count || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('categories.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="(val: number) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAddChild(row)" v-if="row.level < 3">
              {{ $t('categories.addChildCategory') }}
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              {{ $t('categories.edit') }}
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              {{ $t('categories.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? $t('categories.editCategory') : (parentCategory ? $t('categories.addChildTitle') : $t('categories.addTopTitle'))"
      width="600px"
      destroy-on-close
    >
      <el-form :model="categoryForm" label-width="100px" :rules="formRules" ref="formRef">
        <el-form-item :label="$t('categories.categoryName')" prop="name">
          <el-input v-model="categoryForm.name" :placeholder="$t('categories.enterCategoryName')" />
        </el-form-item>
        <el-form-item :label="$t('categories.categoryCode')">
          <el-input v-model="categoryForm.code" :placeholder="$t('categories.enterCategoryCode')" />
        </el-form-item>
        <el-form-item :label="$t('categories.sort')">
          <el-input-number v-model="categoryForm.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('categories.iconUrl')">
          <el-input v-model="categoryForm.icon" :placeholder="$t('categories.enterIconUrl')" />
        </el-form-item>
        <el-form-item :label="$t('categories.imageUrl')">
          <el-input v-model="categoryForm.image" :placeholder="$t('categories.enterImageUrl')" />
        </el-form-item>
        <el-form-item :label="$t('categories.seoTitle')">
          <el-input v-model="categoryForm.seo_title" :placeholder="$t('categories.enterSeoTitle')" />
        </el-form-item>
        <el-form-item :label="$t('categories.seoDescription')">
          <el-input v-model="categoryForm.seo_description" type="textarea" rows="2" :placeholder="$t('categories.enterSeoDescription')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('categories.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('categories.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
  updateCategoryStatus,
  type CategoryTree
} from '@/api/category'

const { t } = useI18n()

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const parentCategory = ref<CategoryTree | null>(null)
const categoryTree = ref<CategoryTree[]>([])
const formRef = ref()

const categoryForm = reactive({
  id: 0,
  name: '',
  code: '',
  sort: 0,
  icon: '',
  image: '',
  seo_title: '',
  seo_description: '',
  parent_id: 0
})

const formRules = {
  name: [{ required: true, message: t('categories.enterCategoryName'), trigger: 'blur' }]
}

const loadCategories = async () => {
  loading.value = true
  try {
    const data = await getCategoryTree()
    categoryTree.value = data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    ElMessage.error(t('categories.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleAddRoot = () => {
  isEdit.value = false
  parentCategory.value = null
  Object.assign(categoryForm, {
    id: 0,
    name: '',
    code: '',
    sort: 0,
    icon: '',
    image: '',
    seo_title: '',
    seo_description: '',
    parent_id: 0
  })
  dialogVisible.value = true
}

const handleAddChild = (row: CategoryTree) => {
  isEdit.value = false
  parentCategory.value = row
  Object.assign(categoryForm, {
    id: 0,
    name: '',
    code: '',
    sort: 0,
    icon: '',
    image: '',
    seo_title: '',
    seo_description: '',
    parent_id: row.id
  })
  dialogVisible.value = true
}

const handleEdit = (row: CategoryTree) => {
  isEdit.value = true
  parentCategory.value = null
  Object.assign(categoryForm, {
    id: row.id,
    name: row.name,
    code: row.code || '',
    sort: row.sort || 0,
    icon: row.icon || '',
    image: row.image || '',
    seo_title: row.seo_title || '',
    seo_description: row.seo_description || '',
    parent_id: row.parent_id
  })
  dialogVisible.value = true
}

const handleDelete = (row: CategoryTree) => {
  if (row.children && row.children.length > 0) {
    ElMessage.warning(t('categories.deleteChildFirst'))
    return
  }

  ElMessageBox.confirm(`confirmDelete: "${row.name}"?`, t('categories.confirmDelete'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCategory(row.id)
      ElMessage.success(t('categories.deleteSuccess'))
      loadCategories()
    } catch (error) {
      console.error('Failed to delete category:', error)
      ElMessage.error(t('categories.deleteFailed'))
    }
  })
}

const handleStatusChange = async (row: CategoryTree, status: number) => {
  try {
    await updateCategoryStatus(row.id, status)
    ElMessage.success(status === 1 ? t('categories.enabledSuccess') : t('categories.disabledSuccess'))
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error(t('categories.updateStatusFailed'))
    row.status = status === 1 ? 0 : 1
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        if (isEdit.value) {
          await updateCategory({
            id: categoryForm.id,
            name: categoryForm.name,
            code: categoryForm.code,
            sort: categoryForm.sort,
            icon: categoryForm.icon,
            image: categoryForm.image,
            seo_title: categoryForm.seo_title,
            seo_description: categoryForm.seo_description
          })
          ElMessage.success(t('categories.updateSuccess'))
        } else {
          await createCategory({
            name: categoryForm.name,
            parent_id: categoryForm.parent_id || undefined,
            code: categoryForm.code,
            sort: categoryForm.sort,
            icon: categoryForm.icon,
            image: categoryForm.image,
            seo_title: categoryForm.seo_title,
            seo_description: categoryForm.seo_description
          })
          ElMessage.success(t('categories.createSuccess'))
        }
        dialogVisible.value = false
        loadCategories()
      } catch (error) {
        console.error('Failed to save category:', error)
        ElMessage.error(isEdit.value ? t('categories.updateFailed') : t('categories.createFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.categories-page {
  padding: 0;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
}

.tree-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
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
</style>
