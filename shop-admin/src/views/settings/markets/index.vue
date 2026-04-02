<template>
  <div class="markets-page">
    <!-- Header Bar -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2 class="page-title">{{ t('settings.markets.title') }}</h2>
          <p class="page-desc">{{ t('settings.markets.pageDesc') }}</p>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>{{ t('settings.markets.addMarket') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Markets Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="marketList" v-loading="loading" stripe>
        <el-table-column :label="t('settings.markets.market')" min-width="200">
          <template #default="{ row }">
            <div class="market-cell">
              <span class="market-flag">{{ row.flag || getFlagEmoji(row.code) }}</span>
              <div class="market-details">
                <p class="market-name">{{ row.name }}</p>
                <p class="market-code">{{ row.code }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('settings.markets.currency')" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ row.currency }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('settings.markets.defaultLanguage')" width="120" align="center">
          <template #default="{ row }">
            {{ getLanguageLabel(row.default_language) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('settings.markets.taxConfig')" min-width="180">
          <template #default="{ row }">
            <div class="tax-config">
              <span v-if="row.tax_rules?.vat_rate">VAT: {{ row.tax_rules.vat_rate }}%</span>
              <span v-if="row.tax_rules?.gst_rate">GST: {{ row.tax_rules.gst_rate }}%</span>
              <el-tag v-if="row.tax_rules?.ioss_enabled" size="small" type="success">IOSS</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('settings.markets.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_active"
              @change="(val: boolean) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="t('settings.markets.defaultMarket')" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="primary" size="small">{{ t('settings.markets.defaultMarket') }}</el-tag>
            <el-button v-else type="primary" link size="small" @click="handleSetDefault(row)">
              {{ t('settings.markets.setAsDefault') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
              :disabled="row.is_default"
            >
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @size-change="loadMarkets"
          @current-change="loadMarkets"
        />
      </div>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('settings.markets.editMarket') : t('settings.markets.addMarket')"
      width="650px"
      destroy-on-close
    >
      <el-form :model="marketForm" label-width="120px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.marketCode')" prop="code">
              <el-select
                v-model="marketForm.code"
                :placeholder="t('settings.markets.selectMarket')"
                style="width: 100%"
                :disabled="isEdit"
              >
                <el-option
                  v-for="market in availableMarkets"
                  :key="market.code"
                  :label="`${market.flag} ${market.name}`"
                  :value="market.code"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.marketName')" prop="name">
              <el-input v-model="marketForm.name" :placeholder="t('settings.markets.enterMarketName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.currency')" prop="currency">
              <el-select v-model="marketForm.currency" :placeholder="t('settings.markets.selectCurrency')" style="width: 100%">
                <el-option label="CNY - RMB" value="CNY" />
                <el-option label="USD - US Dollar" value="USD" />
                <el-option label="EUR - Euro" value="EUR" />
                <el-option label="GBP - British Pound" value="GBP" />
                <el-option label="JPY - Japanese Yen" value="JPY" />
                <el-option label="AUD - Australian Dollar" value="AUD" />
                <el-option label="CAD - Canadian Dollar" value="CAD" />
                <el-option label="SGD - Singapore Dollar" value="SGD" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.defaultLanguage')" prop="default_language">
              <el-select v-model="marketForm.default_language" :placeholder="t('settings.markets.selectLanguage')" style="width: 100%">
                <el-option :label="t('settings.markets.langZhCN')" value="zh-CN" />
                <el-option :label="t('settings.markets.langEnUS')" value="en-US" />
                <el-option :label="t('settings.markets.langJaJP')" value="ja-JP" />
                <el-option :label="t('settings.markets.langDeDE')" value="de-DE" />
                <el-option :label="t('settings.markets.langFrFR')" value="fr-FR" />
                <el-option :label="t('settings.markets.langEsES')" value="es-ES" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="t('settings.markets.status')">
              <el-switch v-model="marketForm.is_active" />
              <span class="form-hint">{{ t('settings.markets.disableMarketHint') }}</span>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Tax Configuration Section -->
        <el-divider content-position="left">{{ t('settings.markets.taxConfig') }}</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.vatRate')">
              <el-input-number
                v-model="marketForm.tax_rules.vat_rate"
                :min="0"
                :max="100"
                :precision="2"
                style="width: 100%"
              />
              <span class="form-hint">{{ t('settings.markets.taxRateDescription') }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.gstRate')">
              <el-input-number
                v-model="marketForm.tax_rules.gst_rate"
                :min="0"
                :max="100"
                :precision="2"
                style="width: 100%"
              />
              <span class="form-hint">{{ t('settings.markets.gstRateDescription') }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.iossEnabled')">
              <el-switch v-model="marketForm.tax_rules.ioss_enabled" />
              <span class="form-hint">{{ t('settings.markets.iossDescription') }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('settings.markets.priceIncludesTax')">
              <el-switch v-model="marketForm.tax_rules.include_tax" />
              <span class="form-hint">{{ t('settings.markets.priceTaxDescription') }}</span>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getMarkets,
  getMarket,
  createMarket,
  updateMarket,
  deleteMarket,
  type Market,
  type CreateMarketRequest,
  type UpdateMarketRequest
} from '@/api/market'
import { t } from '@/plugins/i18n'

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const formRef = ref()
const marketList = ref<Market[]>([])

const marketForm = reactive({
  id: undefined as number | undefined,
  code: '',
  name: '',
  currency: 'USD',
  default_language: 'en-US',
  is_active: true,
  tax_rules: {
    vat_rate: 0 as number | undefined,
    gst_rate: 0 as number | undefined,
    ioss_enabled: false,
    include_tax: false
  }
})

const formRules = {
  code: [{ required: true, message: t('settings.markets.selectMarket'), trigger: 'change' }],
  name: [{ required: true, message: t('settings.markets.enterMarketName'), trigger: 'blur' }],
  currency: [{ required: true, message: t('settings.markets.selectCurrency'), trigger: 'change' }],
  default_language: [{ required: true, message: t('settings.markets.selectLanguage'), trigger: 'change' }]
}

// Available markets for selection (code and flag only, name is translated)
const availableMarketsRaw = [
  { code: 'CN', flag: '🇨🇳' },
  { code: 'US', flag: '🇺🇸' },
  { code: 'GB', flag: '🇬🇧' },
  { code: 'DE', flag: '🇩🇪' },
  { code: 'FR', flag: '🇫🇷' },
  { code: 'JP', flag: '🇯🇵' },
  { code: 'AU', flag: '🇦🇺' },
  { code: 'CA', flag: '🇨🇦' },
  { code: 'SG', flag: '🇸🇬' },
  { code: 'ES', flag: '🇪🇸' },
  { code: 'IT', flag: '🇮🇹' },
  { code: 'NL', flag: '🇳🇱' }
]

// Language code to translation key mapping
const languageKeys: Record<string, string> = {
  'zh-CN': 'settings.markets.langZhCN',
  'en-US': 'settings.markets.langEnUS',
  'ja-JP': 'settings.markets.langJaJP',
  'de-DE': 'settings.markets.langDeDE',
  'fr-FR': 'settings.markets.langFrFR',
  'es-ES': 'settings.markets.langEsES'
}

const getLanguageLabel = (lang: string) => {
  const key = languageKeys[lang]
  return key ? t(key) : lang
}

// Market code to translation key mapping
const marketNameKeys: Record<string, string> = {
  'CN': 'settings.markets.marketCN',
  'US': 'settings.markets.marketUS',
  'GB': 'settings.markets.marketGB',
  'DE': 'settings.markets.marketDE',
  'FR': 'settings.markets.marketFR',
  'JP': 'settings.markets.marketJP',
  'AU': 'settings.markets.marketAU',
  'CA': 'settings.markets.marketCA',
  'SG': 'settings.markets.marketSG',
  'ES': 'settings.markets.marketES',
  'IT': 'settings.markets.marketIT',
  'NL': 'settings.markets.marketNL'
}

// Computed available markets with translated names
const availableMarkets = computed(() => {
  return availableMarketsRaw.map(m => ({
    ...m,
    name: t(marketNameKeys[m.code] || 'settings.markets.market', { code: m.code })
  }))
})

const getFlagEmoji = (code: string) => {
  const market = availableMarketsRaw.find(m => m.code === code)
  return market?.flag || '🌐'
}

const resetForm = () => {
  Object.assign(marketForm, {
    id: undefined,
    code: '',
    name: '',
    currency: 'USD',
    default_language: 'en-US',
    is_active: true,
    tax_rules: {
      vat_rate: 0,
      gst_rate: 0,
      ioss_enabled: false,
      include_tax: false
    }
  })
}

const loadMarkets = async () => {
  loading.value = true
  try {
    const response = await getMarkets()
    marketList.value = response.list || []
    total.value = response.total || 0
  } catch (error: unknown) {
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.loadingFailed'))
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const handleEdit = async (row: Market) => {
  isEdit.value = true
  try {
    const market = await getMarket(row.id)
    Object.assign(marketForm, {
      id: market.id,
      code: market.code,
      name: market.name,
      currency: market.currency,
      default_language: market.default_language,
      is_active: market.is_active,
      tax_rules: {
        vat_rate: parseFloat(market.tax_rules?.vat_rate || '0'),
        gst_rate: parseFloat(market.tax_rules?.gst_rate || '0'),
        ioss_enabled: market.tax_rules?.ioss_enabled ?? false,
        include_tax: market.tax_rules?.include_tax ?? false
      }
    })
    dialogVisible.value = true
  } catch (error: unknown) {
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.getMarketFailed'))
  }
}

const handleDelete = async (row: Market) => {
  if (row.is_default) {
    ElMessage.warning(t('settings.markets.cannotDeleteDefault'))
    return
  }

  try {
    await ElMessageBox.confirm(
      t('settings.markets.confirmDelete', { name: row.name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    await deleteMarket(row.id)
    ElMessage.success(t('settings.markets.deleteSuccess'))
    loadMarkets()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    if (error !== 'cancel') {
      ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.deleteFailed'))
    }
  }
}

const handleStatusChange = async (row: Market, val: boolean) => {
  try {
    await updateMarket({
      id: row.id,
      is_active: val
    })
    ElMessage.success(val ? t('settings.markets.enabled') : t('settings.markets.disabled'))
  } catch (error: unknown) {
    row.is_active = !val
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.statusUpdateFailed'))
  }
}

const handleSetDefault = async (row: Market) => {
  try {
    await ElMessageBox.confirm(
      t('settings.markets.confirmSetDefault', { name: row.name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'info'
      }
    )

    // Update the market to be default
    await updateMarket({
      id: row.id,
      is_default: true
    })

    ElMessage.success(t('settings.markets.defaultSetSuccess'))
    loadMarkets()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    if (error !== 'cancel') {
      ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.setDefaultFailed'))
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    saveLoading.value = true

    if (isEdit.value && marketForm.id) {
      // Update existing market
      const updateData: UpdateMarketRequest = {
        id: marketForm.id,
        name: marketForm.name,
        is_active: marketForm.is_active,
        tax_rules: {
          vat_rate: String(marketForm.tax_rules?.vat_rate ?? 0),
          gst_rate: String(marketForm.tax_rules?.gst_rate ?? 0),
          ioss_enabled: marketForm.tax_rules?.ioss_enabled,
          include_tax: marketForm.tax_rules?.include_tax
        }
      }
      await updateMarket(updateData)
      ElMessage.success(t('settings.markets.updateSuccess'))
    } else {
      // Create new market
      const createData: CreateMarketRequest = {
        code: marketForm.code,
        name: marketForm.name,
        currency: marketForm.currency,
        default_language: marketForm.default_language,
        tax_rules: {
          vat_rate: String(marketForm.tax_rules?.vat_rate ?? 0),
          gst_rate: String(marketForm.tax_rules?.gst_rate ?? 0),
          ioss_enabled: marketForm.tax_rules?.ioss_enabled,
          include_tax: marketForm.tax_rules?.include_tax
        }
      }
      await createMarket(createData)
      ElMessage.success(t('settings.markets.createSuccess'))
    }

    dialogVisible.value = false
    loadMarkets()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { msg?: string } }; message?: string }
    if (error !== false) {
      ElMessage.error(err.response?.data?.msg || err.message || t('settings.markets.saveFailed'))
    }
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => {
  loadMarkets()
})
</script>

<style scoped>
.markets-page {
  padding: 0;
}

/* Header Card */
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
  flex: 1;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: #6B7280;
}

.header-right {
  display: flex;
  gap: 12px;
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

/* Market Cell */
.market-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.market-flag {
  font-size: 28px;
  line-height: 1;
}

.market-details {
  flex: 1;
  min-width: 0;
}

.market-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-size: 14px;
}

.market-code {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
  font-family: 'Fira Code', monospace;
}

/* Tax Config */
.tax-config {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 12px;
  color: #6B7280;
}

.tax-config span {
  padding: 2px 8px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 6px;
  color: #6366F1;
}

/* Form Hint */
.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #9CA3AF;
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

:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
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

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Divider */
:deep(.el-divider__text) {
  font-weight: 500;
  color: #374151;
}

/* Responsive */
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
    width: 100%;
  }
}
</style>