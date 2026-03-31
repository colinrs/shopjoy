<template>
  <div class="categories-page">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2>{{ $t('categories.title') }}</h2>
        </div>
        <div class="header-right">
          <el-button type="info" @click="handleSaveSort" :loading="sortLoading" :disabled="!hasSortChanges">
            <el-icon><Rank /></el-icon>
            {{ $t('categories.saveSort') }}
          </el-button>
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
        @cell-mouse-enter="handleMouseEnter"
        @cell-mouse-leave="handleMouseLeave"
      >
        <el-table-column label="" width="50" align="center">
          <template #default="{ row }">
            <div
              class="drag-handle"
              :class="{ 'drag-handle-active': dragState.dragging && dragState.sourceId === row.id }"
              draggable="true"
              @dragstart="handleDragStart($event, row)"
              @dragend="handleDragEnd"
              @dragover="handleDragOver"
              @drop="handleDrop($event, row)"
            >
              <el-icon><Rank /></el-icon>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('categories.name')" min-width="200" />
        <el-table-column prop="code" :label="$t('categories.categoryCode')" width="120" />
        <el-table-column prop="level" :label="$t('categories.level')" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small">L{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="$t('categories.sort')" width="80" align="center">
          <template #default="{ row }">
            <span v-if="!dragState.dragging || dragState.sourceId !== row.id">{{ row.sort }}</span>
            <span v-else class="sort-placeholder">-</span>
          </template>
        </el-table-column>
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
        <el-table-column :label="$t('common.actions')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">
              {{ $t('categories.viewDetail') }}
            </el-button>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Rank } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
  updateCategoryStatus,
  updateCategorySort,
  type CategoryTree
} from '@/api/category'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const saveLoading = ref(false)
const sortLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const parentCategory = ref<CategoryTree | null>(null)
const categoryTree = ref<CategoryTree[]>([])
const originalTree = ref<CategoryTree[]>([])
const formRef = ref()

// Drag state
const dragState = reactive({
  dragging: false,
  sourceId: null as number | null,
  targetId: null as number | null
})

// Check if sort order has changed
const hasSortChanges = computed(() => {
  return JSON.stringify(categoryTree.value) !== JSON.stringify(originalTree.value)
})

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
    originalTree.value = JSON.parse(JSON.stringify(data || []))
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

const handleViewDetail = (row: CategoryTree) => {
  router.push(`/categories/${row.id}`)
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

// Drag and drop handlers
const handleDragStart = (event: DragEvent, row: CategoryTree) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', row.id.toString())
  }
  dragState.dragging = true
  dragState.sourceId = row.id
}

const handleDragEnd = () => {
  dragState.dragging = false
  dragState.sourceId = null
  dragState.targetId = null
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
}

const handleDrop = (event: DragEvent, targetRow: CategoryTree) => {
  event.preventDefault()
  if (!dragState.sourceId || dragState.sourceId === targetRow.id) {
    handleDragEnd()
    return
  }

  // Find source and target nodes
  const sourceNode = findNodeById(categoryTree.value, dragState.sourceId)
  if (!sourceNode) {
    handleDragEnd()
    return
  }

  // Only allow reordering within the same parent (same level)
  if (sourceNode.parent_id !== targetRow.parent_id) {
    ElMessage.warning(t('categories.dragWithinSameLevel'))
    handleDragEnd()
    return
  }

  // Find siblings and reorder
  const siblings = getSiblings(categoryTree.value, sourceNode.parent_id)
  const sourceIndex = siblings.findIndex(s => s.id === dragState.sourceId)
  const targetIndex = siblings.findIndex(s => s.id === targetRow.id)

  if (sourceIndex === -1 || targetIndex === -1) {
    handleDragEnd()
    return
  }

  // Remove source from original position
  siblings.splice(sourceIndex, 1)
  // Insert source at new position
  siblings.splice(targetIndex, 0, sourceNode)

  // Update sort values
  siblings.forEach((sibling, index) => {
    sibling.sort = index
  })

  // Trigger reactivity
  categoryTree.value = [...categoryTree.value]

  handleDragEnd()
}

const findNodeById = (nodes: CategoryTree[], id: number): CategoryTree | null => {
  for (const node of nodes) {
    if (node.id === id) {
      return node
    }
    if (node.children && node.children.length > 0) {
      const found = findNodeById(node.children, id)
      if (found) {
        return found
      }
    }
  }
  return null
}

const getSiblings = (nodes: CategoryTree[], parentId: number): CategoryTree[] => {
  if (parentId === 0) {
    return nodes
  }
  for (const node of nodes) {
    if (node.id === parentId) {
      return node.children || []
    }
    if (node.children && node.children.length > 0) {
      const found = getSiblings(node.children, parentId)
      if (found.length > 0) {
        return found
      }
    }
  }
  return []
}

const handleMouseEnter = (_row: CategoryTree) => {
  // Could be used for visual feedback
}

const handleMouseLeave = (_row: CategoryTree) => {
  // Could be used for visual feedback
}

// Save sort order
const handleSaveSort = async () => {
  sortLoading.value = true
  try {
    // Collect all sort changes
    const sorts: { id: number; sort: number }[] = []
    collectSorts(categoryTree.value, sorts)

    if (sorts.length > 0) {
      await updateCategorySort(sorts)
      ElMessage.success(t('categories.saveSortSuccess'))
      // Refresh to get latest data
      loadCategories()
    }
  } catch (error) {
    console.error('Failed to save sort:', error)
    ElMessage.error(t('categories.saveSortFailed'))
  } finally {
    sortLoading.value = false
  }
}

const collectSorts = (nodes: CategoryTree[], result: { id: number; sort: number }[]) => {
  for (const node of nodes) {
    result.push({ id: node.id, sort: node.sort })
    if (node.children && node.children.length > 0) {
      collectSorts(node.children, result)
    }
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

/* Drag handle */
.drag-handle {
  cursor: grab;
  opacity: 0.3;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.drag-handle:hover {
  opacity: 0.8;
}

.drag-handle-active {
  opacity: 1;
  cursor: grabbing;
}

.sort-placeholder {
  color: #9CA3AF;
}

/* Dragging row style */
:deep(.el-table__row.dragging) {
  opacity: 0.5;
  background-color: #F5F3FF;
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
