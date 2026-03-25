<template>
  <div class="shipping-templates-page">
    <!-- Page Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2>运费模板</h2>
          <p class="header-desc">管理不同地区的运费规则</p>
        </div>
        <div class="header-right">
          <el-button @click="goToCalculator">
            <el-icon><Coin /></el-icon>
            运费计算器
          </el-button>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建模板
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索模板名称..."
          class="search-input"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select v-model="filterStatus" placeholder="状态" clearable class="filter-select" @change="handleSearch">
          <el-option label="全部" value="" />
          <el-option label="已启用" value="true" />
          <el-option label="已禁用" value="false" />
        </el-select>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </div>
    </el-card>

    <!-- Templates Grid -->
    <div class="templates-grid" v-loading="loading">
      <el-empty v-if="!loading && templates.length === 0" description="暂无运费模板">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建模板
        </el-button>
      </el-empty>

      <TemplateCard
        v-for="template in templates"
        :key="template.id"
        :template="template"
        @edit="handleEdit"
        @delete="handleDelete"
        @set-default="handleSetDefault"
      />
    </div>

    <!-- Pagination -->
    <el-card v-if="total > 0" class="pagination-card" shadow="never">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="loadTemplates"
        @current-change="loadTemplates"
      />
    </el-card>

    <!-- Create Dialog -->
    <el-dialog
      v-model="createDialogVisible"
      title="新建运费模板"
      width="500px"
      destroy-on-close
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="createForm.name" placeholder="例如：全国包邮模板" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="createForm.is_default" />
          <div class="form-tip">默认模板将用于未指定模板的商品</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmCreate" :loading="createLoading">
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Coin } from '@element-plus/icons-vue'
import {
  getShippingTemplates,
  createShippingTemplate,
  deleteShippingTemplate,
  setDefaultTemplate,
  type ShippingTemplate
} from '@/api/shipping'
import TemplateCard from './components/TemplateCard.vue'

const router = useRouter()

// State
const loading = ref(false)
const createLoading = ref(false)
const createDialogVisible = ref(false)
const templates = ref<ShippingTemplate[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const searchQuery = ref('')
const filterStatus = ref('')
const createFormRef = ref()

const createForm = reactive({
  name: '',
  is_default: false
})

const createRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度为2-50个字符', trigger: 'blur' }
  ]
}

// Methods
const loadTemplates = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchQuery.value) {
      params.name = searchQuery.value
    }
    if (filterStatus.value !== '') {
      params.is_active = filterStatus.value === 'true'
    }
    const data = await getShippingTemplates(params)
    templates.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to load templates:', error)
    ElMessage.error('加载模板失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadTemplates()
}

const handleCreate = () => {
  createForm.name = ''
  createForm.is_default = false
  createDialogVisible.value = true
}

const handleConfirmCreate = async () => {
  if (!createFormRef.value) return

  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      createLoading.value = true
      try {
        await createShippingTemplate({
          name: createForm.name,
          is_default: createForm.is_default
        })
        ElMessage.success('创建成功')
        createDialogVisible.value = false
        loadTemplates()
      } catch (error) {
        console.error('Failed to create template:', error)
        ElMessage.error('创建失败')
      } finally {
        createLoading.value = false
      }
    }
  })
}

const handleEdit = (id: number) => {
  router.push(`/shipping/${id}`)
}

const handleDelete = async (template: ShippingTemplate) => {
  if (template.is_default) {
    ElMessage.warning('无法删除默认模板，请先设置其他模板为默认')
    return
  }

  try {
    await ElMessageBox.confirm(`确认删除模板 "${template.name}"？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteShippingTemplate(template.id)
    ElMessage.success('删除成功')
    loadTemplates()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete template:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleSetDefault = async (template: ShippingTemplate) => {
  try {
    await setDefaultTemplate(template.id)
    ElMessage.success('已设置为默认模板')
    loadTemplates()
  } catch (error) {
    console.error('Failed to set default:', error)
    ElMessage.error('设置失败')
  }
}

const goToCalculator = () => {
  router.push('/shipping/calculator')
}

// Lifecycle
onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.shipping-templates-page {
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

.header-desc {
  margin: 4px 0 0;
  font-size: 13px;
  color: #6B7280;
}

.header-right {
  display: flex;
  gap: 12px;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 300px;
}

.filter-select {
  width: 120px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;
  min-height: 200px;
}

.pagination-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.pagination-card :deep(.el-card__body) {
  display: flex;
  justify-content: flex-end;
}

.form-tip {
  font-size: 12px;
  color: #6B7280;
  margin-top: 4px;
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

/* Responsive */
@media (max-width: 1200px) {
  .templates-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .header-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-right {
    width: 100%;
  }

  .header-right .el-button {
    flex: 1;
  }

  .templates-grid {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-direction: column;
  }

  .search-input {
    width: 100%;
  }

  .filter-select {
    width: 100%;
  }
}
</style>