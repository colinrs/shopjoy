<template>
  <div class="brand-detail-page" v-loading="loading">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <el-button text @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('common.back') }}
          </el-button>
          <el-divider direction="vertical" />
          <span class="page-title">{{ brand?.name || t('brands.brandDetail') }}</span>
          <el-tag v-if="brand?.status === 1" type="success" size="small">{{ $t('brands.enabled') }}</el-tag>
          <el-tag v-else type="info" size="small">{{ $t('brands.disabled') }}</el-tag>
        </div>
        <div class="header-right">
          <el-button @click="handleEdit">
            <el-icon><Edit /></el-icon>
            {{ $t('brands.edit') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Brand Info -->
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('brands.brandInfo') }}</span>
            </div>
          </template>
          <el-descriptions :column="2" border v-if="brand">
            <el-descriptions-item :label="$t('brands.brandName')">
              <div class="brand-name-cell">
                <el-image v-if="brand.logo" :src="brand.logo" fit="cover" class="brand-logo">
                  <template #error>
                    <div class="image-placeholder"><el-icon><Picture /></el-icon></div>
                  </template>
                </el-image>
                <span>{{ brand.name }}</span>
              </div>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.website')">
              <a v-if="brand.website" :href="brand.website" target="_blank" class="website-link">
                {{ brand.website }}
              </a>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.brandDescription')" :span="2">
              <span v-if="brand.description">{{ brand.description }}</span>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.trademarkNumber')">
              {{ brand.trademark_number || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.trademarkCountry')">
              {{ brand.trademark_country || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.brandPage')">
              <el-switch
                v-model="brand.enable_page"
                :active-value="1"
                :inactive-value="0"
                @change="handleTogglePage"
              />
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.sort')">
              {{ brand.sort }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('brands.createdAt')">
              {{ brand.created_at }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- Market Visibility -->
        <el-card class="info-card market-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('brands.marketVisibility') }}</span>
              <el-button size="small" type="primary" @click="showMarketDialog = true">
                {{ $t('brands.configureMarket') }}
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
          <el-empty v-else :description="$t('brands.noMarketVisibility')" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <!-- Stats Card -->
        <el-card class="stats-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('brands.statistics') }}</span>
            </div>
          </template>
          <div class="stat-item">
            <div class="stat-icon products">
              <el-icon><Goods /></el-icon>
            </div>
            <div class="stat-content">
              <p class="stat-value">{{ brand?.product_count || 0 }}</p>
              <p class="stat-label">{{ $t('brands.productCount') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Edit Dialog -->
    <el-dialog
      v-model="editDialogVisible"
      :title="$t('brands.editBrand')"
      width="700px"
      destroy-on-close
    >
      <el-form :model="brandForm" label-width="120px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('brands.brandName')" prop="name">
              <el-input v-model="brandForm.name" :placeholder="$t('brands.enterBrandName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('brands.sort')">
              <el-input-number v-model="brandForm.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('brands.logoUrl')">
              <el-input v-model="brandForm.logo" :placeholder="$t('brands.brandLogoUrl')" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('brands.brandWebsite')">
              <el-input v-model="brandForm.website" :placeholder="$t('brands.websiteUrl')" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('brands.brandDescription')">
              <el-input v-model="brandForm.description" type="textarea" rows="3" :placeholder="$t('brands.enterDescription')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('brands.trademarkNumber')">
              <el-input v-model="brandForm.trademark_number" :placeholder="$t('brands.trademarkRegNumber')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('brands.trademarkCountry')">
              <el-input v-model="brandForm.trademark_country" :placeholder="$t('brands.countryExample')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('brands.enableBrandPage')">
              <el-switch v-model="brandForm.enable_page" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">{{ $t('brands.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('brands.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Market Visibility Dialog -->
    <el-dialog
      v-model="showMarketDialog"
      :title="$t('brands.marketVisibility')"
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
        <el-button @click="showMarketDialog = false">{{ $t('brands.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveMarketVisibility" :loading="savingMarket">
          {{ $t('brands.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Edit, Picture, Goods } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getBrand,
  updateBrand,
  toggleBrandPage,
  getBrandMarketVisibility,
  setBrandMarketVisibility,
  type Brand
} from '@/api/brand'
import { getMarkets, type Market } from '@/api/market'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()
const route = useRoute()
const router = useRouter()

const brandId = Number(route.params.id)

const loading = ref(false)
const saveLoading = ref(false)
const savingMarket = ref(false)
const marketsLoading = ref(false)
const formRef = ref()
const brand = ref<Brand | null>(null)
const editDialogVisible = ref(false)
const showMarketDialog = ref(false)
const marketVisibility = ref<{ market_id: number; is_visible: boolean }[]>([])
const availableMarkets = ref<(Market & { selected: boolean })[]>([])

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
  name: [{ required: true, message: t('brands.enterBrandName'), trigger: 'blur' }]
}

const goBack = () => {
  router.push('/brands')
}

const loadBrand = async () => {
  loading.value = true
  try {
    brand.value = await getBrand(brandId)
  } catch (error) {
    ElMessage.error(t('brands.loadFailed'))
    router.push('/brands')
  } finally {
    loading.value = false
  }
}

const loadMarketVisibility = async () => {
  try {
    const res = await getBrandMarketVisibility(brandId)
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
  if (!brand.value) return
  Object.assign(brandForm, {
    id: brand.value.id,
    name: brand.value.name,
    logo: brand.value.logo || '',
    description: brand.value.description || '',
    website: brand.value.website || '',
    trademark_number: brand.value.trademark_number || '',
    trademark_country: brand.value.trademark_country || '',
    enable_page: brand.value.enable_page === 1,
    sort: brand.value.sort || 0
  })
  editDialogVisible.value = true
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
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
        ElMessage.success(t('brands.updateSuccess'))
        editDialogVisible.value = false
        loadBrand()
      } catch (error) {
        handleError(error, t('brands.updateFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

const handleTogglePage = async () => {
  if (!brand.value) return
  try {
    await toggleBrandPage(brand.value.id, brand.value.enable_page === 1)
    ElMessage.success(brand.value.enable_page === 1 ? t('brands.enabledPageSuccess') : t('brands.disabledPageSuccess'))
  } catch (error) {
    handleError(error, t('brands.operationFailed'))
    // Revert the change
    brand.value.enable_page = brand.value.enable_page === 1 ? 0 : 1
  }
}

const handleSaveMarketVisibility = async () => {
  savingMarket.value = true
  try {
    const selectedIds = availableMarkets.value.filter(m => m.selected).map(m => m.id)
    await setBrandMarketVisibility(brandId, selectedIds, true)
    ElMessage.success(t('brands.updateSuccess'))
    showMarketDialog.value = false
    loadMarketVisibility()
  } catch (error) {
    handleError(error, t('brands.updateFailed'))
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
  loadBrand()
  loadMarketVisibility()
})
</script>

<script lang="ts">
import { watch } from 'vue'
export default {
  name: 'BrandDetail'
}
</script>

<style scoped>
.brand-detail-page {
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

.brand-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.image-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  color: #6366F1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.website-link {
  color: #6366F1;
  text-decoration: none;
}

.website-link:hover {
  text-decoration: underline;
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
