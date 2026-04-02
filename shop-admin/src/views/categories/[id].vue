<template>
  <div class="category-detail-page" v-loading="loading">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <el-button text @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('common.back') }}
          </el-button>
          <el-divider direction="vertical" />
          <span class="page-title">{{ category?.name || t('categories.categoryDetail') }}</span>
          <el-tag v-if="category?.status === 1" type="success" size="small">{{ $t('categories.enabled') }}</el-tag>
          <el-tag v-else type="info" size="small">{{ $t('categories.disabled') }}</el-tag>
        </div>
        <div class="header-right">
          <el-button @click="handleEdit">
            <el-icon><Edit /></el-icon>
            {{ $t('categories.edit') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Category Info -->
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('categories.categoryInfo') }}</span>
            </div>
          </template>
          <el-descriptions :column="2" border v-if="category">
            <el-descriptions-item :label="$t('categories.categoryName')">
              {{ category.name }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.categoryCode')">
              {{ category.code || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.level')">
              <el-tag size="small">L{{ category.level }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.sort')">
              {{ category.sort }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.icon')" v-if="category.icon">
              <el-image :src="category.icon" fit="cover" class="icon-image">
                <template #error>
                  <div class="image-placeholder"><el-icon><Picture /></el-icon></div>
                </template>
              </el-image>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.image')" v-if="category.image">
              <el-image :src="category.image" fit="cover" class="image-preview">
                <template #error>
                  <div class="image-placeholder"><el-icon><Picture /></el-icon></div>
                </template>
              </el-image>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.seoTitle')" :span="2">
              {{ category.seo_title || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.seoDescription')" :span="2">
              {{ category.seo_description || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('categories.createdAt')">
              {{ category.created_at }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- Market Visibility -->
        <el-card class="info-card market-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('categories.marketVisibility') }}</span>
              <el-button size="small" type="primary" @click="showMarketDialog = true">
                {{ $t('categories.configureMarket') }}
              </el-button>
            </div>
          </template>
          <div class="market-tags" v-if="marketVisibility.length > 0">
            <el-tag
              v-for="m in marketVisibility"
              :key="m.market_id"
              :type="m.is_visible ? 'success' : 'info'"
              class="market-tag"
            >
              {{ getMarketName(m.market_id) }}
            </el-tag>
          </div>
          <el-empty v-else :description="$t('categories.noMarketVisibility')" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <!-- Stats Card -->
        <el-card class="stats-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('categories.statistics') }}</span>
            </div>
          </template>
          <div class="stat-item">
            <div class="stat-icon products">
              <el-icon><Goods /></el-icon>
            </div>
            <div class="stat-content">
              <p class="stat-value">{{ category?.product_count || 0 }}</p>
              <p class="stat-label">{{ $t('categories.productCount') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Edit Dialog -->
    <el-dialog
      v-model="editDialogVisible"
      :title="$t('categories.editCategory')"
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
        <el-button @click="editDialogVisible = false">{{ $t('categories.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('categories.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Market Visibility Dialog -->
    <el-dialog
      v-model="showMarketDialog"
      :title="$t('categories.marketVisibility')"
      width="500px"
      destroy-on-close
    >
      <div class="market-selector" v-loading="marketsLoading">
        <el-checkbox
          v-for="market in availableMarkets"
          :key="market.id"
          v-model="market.selected"
          class="market-checkbox"
        >
          <span class="market-name">{{ market.name }}</span>
          <span class="market-code">({{ market.code }})</span>
        </el-checkbox>
      </div>
      <template #footer>
        <el-button @click="showMarketDialog = false">{{ $t('categories.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveMarketVisibility" :loading="savingMarket">
          {{ $t('categories.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Edit, Picture, Goods } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getCategory,
  updateCategory,
  getCategoryMarketVisibility,
  setCategoryMarketVisibility,
  type Category
} from '@/api/category'
import { getMarkets, type Market } from '@/api/market'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()
const route = useRoute()
const router = useRouter()

const categoryId = Number(route.params.id)

const loading = ref(false)
const saveLoading = ref(false)
const savingMarket = ref(false)
const marketsLoading = ref(false)
const formRef = ref()
const category = ref<Category | null>(null)
const editDialogVisible = ref(false)
const showMarketDialog = ref(false)
const marketVisibility = ref<{ market_id: number; is_visible: boolean }[]>([])
const availableMarkets = ref<(Market & { selected: boolean })[]>([])

const categoryForm = reactive({
  id: 0,
  name: '',
  code: '',
  sort: 0,
  icon: '',
  image: '',
  seo_title: '',
  seo_description: ''
})

const formRules = {
  name: [{ required: true, message: t('categories.enterCategoryName'), trigger: 'blur' }]
}

const goBack = () => {
  router.push('/categories')
}

const loadCategory = async () => {
  loading.value = true
  try {
    category.value = await getCategory(categoryId)
  } catch (error) {
    ElMessage.error(t('categories.loadFailed'))
    router.push('/categories')
  } finally {
    loading.value = false
  }
}

const loadMarketVisibility = async () => {
  try {
    const res = await getCategoryMarketVisibility(categoryId)
    marketVisibility.value = res.markets || []
  } catch (error) {
    handleError(error)
  }
}

const loadMarkets = async () => {
  marketsLoading.value = true
  try {
    const res = await getMarkets()
    const markets = res.list || []

    // Mark selected markets based on visibility
    availableMarkets.value = markets.map(m => ({
      ...m,
      selected: marketVisibility.value.some(v => v.market_id === m.id && v.is_visible)
    }))
  } catch (error) {
    handleError(error)
  } finally {
    marketsLoading.value = false
  }
}

const getMarketName = (marketId: number) => {
  const market = availableMarkets.value.find(m => m.id === marketId)
  return market?.name || `Market ${marketId}`
}

const handleEdit = () => {
  if (!category.value) return
  Object.assign(categoryForm, {
    id: category.value.id,
    name: category.value.name,
    code: category.value.code || '',
    sort: category.value.sort || 0,
    icon: category.value.icon || '',
    image: category.value.image || '',
    seo_title: category.value.seo_title || '',
    seo_description: category.value.seo_description || ''
  })
  editDialogVisible.value = true
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
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
        editDialogVisible.value = false
        loadCategory()
      } catch (error) {
        handleError(error, t('categories.updateFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

const handleSaveMarketVisibility = async () => {
  savingMarket.value = true
  try {
    const selectedIds = availableMarkets.value.filter(m => m.selected).map(m => m.id)
    await setCategoryMarketVisibility(categoryId, selectedIds, true)
    ElMessage.success(t('categories.updateSuccess'))
    showMarketDialog.value = false
    loadMarketVisibility()
  } catch (error) {
    handleError(error, t('categories.updateFailed'))
  } finally {
    savingMarket.value = false
  }
}

const openMarketDialog = async () => {
  await loadMarkets()
  showMarketDialog.value = true
}

watch(showMarketDialog, (val) => {
  if (val) {
    openMarketDialog()
  }
})

onMounted(() => {
  loadCategory()
  loadMarketVisibility()
})
</script>

<script lang="ts">
export default {
  name: 'CategoryDetail'
}
</script>

<style scoped>
.category-detail-page {
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

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
}

.info-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #1E1B4B;
}

.icon-image,
.image-preview {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.image-placeholder {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  color: #6366F1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.market-card :deep(.el-card__header) {
  border-bottom: none;
}

.market-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.market-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.stats-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #fff;
}

.stat-icon.products {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

.market-selector {
  max-height: 400px;
  overflow-y: auto;
}

.market-checkbox {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 8px 0;
}

.market-name {
  font-weight: 500;
  color: #1E1B4B;
}

.market-code {
  color: #9CA3AF;
  font-size: 12px;
  margin-left: 4px;
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
