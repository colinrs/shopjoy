<template>
  <div class="shipping-templates-page">
    <!-- Page Header -->
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="header-bar">
        <div class="header-left">
          <h2>{{ $t('shipping.templatesTitle') }}</h2>
          <p class="header-desc">
            {{ $t('shipping.templatesDesc') }}
          </p>
        </div>
        <div class="header-right">
          <el-button @click="goToCalculator">
            <el-icon><Coin /></el-icon>
            {{ $t('shipping.calculatorBtn') }}
          </el-button>
          <el-button
            type="primary"
            @click="handleCreate"
          >
            <el-icon><Plus /></el-icon>
            {{ $t('shipping.createTemplateBtn') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Filter Bar -->
    <el-card
      class="filter-card"
      shadow="never"
    >
      <div class="filter-bar">
        <el-input
          v-model="searchQuery"
          :placeholder="$t('shipping.searchPlaceholder')"
          class="search-input"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="filterStatus"
          :placeholder="$t('shipping.status')"
          clearable
          class="filter-select"
          @change="handleSearch"
        >
          <el-option
            :label="$t('shipping.all')"
            value=""
          />
          <el-option
            :label="$t('shipping.enabled')"
            value="true"
          />
          <el-option
            :label="$t('shipping.disabled')"
            value="false"
          />
        </el-select>
        <el-button
          type="primary"
          @click="handleSearch"
        >
          <el-icon><Search /></el-icon>
          {{ $t('shipping.search') }}
        </el-button>
      </div>
    </el-card>

    <!-- Templates Grid -->
    <div
      v-loading="loading"
      class="templates-grid"
    >
      <el-empty
        v-if="!loading && templates.length === 0"
        :description="$t('shipping.noTemplates')"
      >
        <el-button
          type="primary"
          @click="handleCreate"
        >
          <el-icon><Plus /></el-icon>
          {{ $t('shipping.createTemplate') }}
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
    <el-card
      v-if="total > 0"
      class="pagination-card"
      shadow="never"
    >
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
      :title="$t('shipping.createDialogTitle')"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="100px"
      >
        <el-form-item
          :label="$t('shipping.templateName')"
          prop="name"
        >
          <el-input
            v-model="createForm.name"
            :placeholder="$t('shipping.templateNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="$t('shipping.setAsDefault')">
          <el-switch v-model="createForm.is_default" />
          <div class="form-tip">
            {{ $t('shipping.defaultTip') }}
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="createLoading"
          @click="handleConfirmCreate"
        >
          {{ $t('shipping.create') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Coin } from '@element-plus/icons-vue'
import {
  getShippingTemplates,
  createShippingTemplate,
  deleteShippingTemplate,
  setDefaultTemplate,
  type ShippingTemplate,
  type TemplateListParams
} from '@/api/shipping'
import TemplateCard from './components/TemplateCard.vue'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { t } = useI18n()
const { handleError } = useErrorHandler()

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
    { required: true, message: 'shipping.templateNameRequired', trigger: 'blur' },
    { min: 2, max: 50, message: 'shipping.templateNameLength', trigger: 'blur' }
  ]
}

// Methods
const loadTemplates = async () => {
  loading.value = true
  try {
    const params: TemplateListParams = {
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
    handleError(error, t('shipping.loadTemplateFailed'))
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
        ElMessage.success(t('shipping.createSuccess'))
        createDialogVisible.value = false
        loadTemplates()
      } catch (error) {
        handleError(error, t('shipping.createFailed'))
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
    ElMessage.warning(t('shipping.cannotDeleteDefault'))
    return
  }

  try {
    await ElMessageBox.confirm(t('shipping.confirmDelete', { name: template.name }), t('shipping.warning'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await deleteShippingTemplate(template.id)
    ElMessage.success(t('shipping.deleteSuccess'))
    loadTemplates()
  } catch (error: unknown) {
    if (error !== 'cancel') {
      handleError(error, t('shipping.deleteFailed'))
    }
  }
}

const handleSetDefault = async (template: ShippingTemplate) => {
  try {
    await setDefaultTemplate(template.id)
    ElMessage.success(t('shipping.setAsDefaultSuccess'))
    loadTemplates()
  } catch (error) {
    handleError(error, t('shipping.setAsDefaultFailed'))
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