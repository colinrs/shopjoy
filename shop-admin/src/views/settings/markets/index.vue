<template>
  <div class="markets-page">
    <!-- Header Bar -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2 class="page-title">市场管理</h2>
          <p class="page-desc">管理多市场配置，包括货币、语言和税务规则</p>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增市场
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Markets Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="marketList" v-loading="loading" stripe>
        <el-table-column label="市场" min-width="200">
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
        <el-table-column label="货币" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ row.currency }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="默认语言" width="120" align="center">
          <template #default="{ row }">
            {{ getLanguageLabel(row.default_language) }}
          </template>
        </el-table-column>
        <el-table-column label="税务配置" min-width="180">
          <template #default="{ row }">
            <div class="tax-config">
              <span v-if="row.tax_rules?.vat_rate">VAT: {{ row.tax_rules.vat_rate }}%</span>
              <span v-if="row.tax_rules?.gst_rate">GST: {{ row.tax_rules.gst_rate }}%</span>
              <el-tag v-if="row.tax_rules?.ioss_enabled" size="small" type="success">IOSS</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_active"
              @change="(val) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="主市场" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="primary" size="small">主市场</el-tag>
            <el-button v-else type="primary" link size="small" @click="handleSetDefault(row)">
              设为主市场
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
              :disabled="row.is_default"
            >
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
      :title="isEdit ? '编辑市场' : '新增市场'"
      width="650px"
      destroy-on-close
    >
      <el-form :model="marketForm" label-width="120px" :rules="formRules" ref="formRef">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="市场代码" prop="code">
              <el-select
                v-model="marketForm.code"
                placeholder="请选择市场"
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
            <el-form-item label="市场名称" prop="name">
              <el-input v-model="marketForm.name" placeholder="请输入市场名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="货币" prop="currency">
              <el-select v-model="marketForm.currency" placeholder="请选择货币" style="width: 100%">
                <el-option label="CNY - 人民币" value="CNY" />
                <el-option label="USD - 美元" value="USD" />
                <el-option label="EUR - 欧元" value="EUR" />
                <el-option label="GBP - 英镑" value="GBP" />
                <el-option label="JPY - 日元" value="JPY" />
                <el-option label="AUD - 澳元" value="AUD" />
                <el-option label="CAD - 加元" value="CAD" />
                <el-option label="SGD - 新加坡元" value="SGD" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认语言" prop="default_language">
              <el-select v-model="marketForm.default_language" placeholder="请选择语言" style="width: 100%">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
                <el-option label="日本語" value="ja-JP" />
                <el-option label="Deutsch" value="de-DE" />
                <el-option label="Français" value="fr-FR" />
                <el-option label="Español" value="es-ES" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="启用状态">
              <el-switch v-model="marketForm.is_active" />
              <span class="form-hint">关闭后该市场将不对外展示</span>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Tax Configuration Section -->
        <el-divider content-position="left">税务配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="VAT 税率">
              <el-input-number
                v-model="marketForm.tax_rules.vat_rate"
                :min="0"
                :max="100"
                :precision="2"
                style="width: 100%"
              />
              <span class="form-hint">增值税税率 (欧盟市场)</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="GST 税率">
              <el-input-number
                v-model="marketForm.tax_rules.gst_rate"
                :min="0"
                :max="100"
                :precision="2"
                style="width: 100%"
              />
              <span class="form-hint">商品及服务税税率 (澳洲市场)</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用 IOSS">
              <el-switch v-model="marketForm.tax_rules.ioss_enabled" />
              <span class="form-hint">欧盟进口一站式申报服务</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="含税价格">
              <el-switch v-model="marketForm.tax_rules.include_tax" />
              <span class="form-hint">商品价格是否含税</span>
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

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const formRef = ref()
const marketList = ref<Market[]>([])

const marketForm = reactive<CreateMarketRequest & { id?: number; is_active?: boolean }>({
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

const formRules = {
  code: [{ required: true, message: '请选择市场', trigger: 'change' }],
  name: [{ required: true, message: '请输入市场名称', trigger: 'blur' }],
  currency: [{ required: true, message: '请选择货币', trigger: 'change' }],
  default_language: [{ required: true, message: '请选择默认语言', trigger: 'change' }]
}

// Available markets for selection
const availableMarkets = [
  { code: 'CN', name: '中国', flag: '🇨🇳' },
  { code: 'US', name: '美国', flag: '🇺🇸' },
  { code: 'GB', name: '英国', flag: '🇬🇧' },
  { code: 'DE', name: '德国', flag: '🇩🇪' },
  { code: 'FR', name: '法国', flag: '🇫🇷' },
  { code: 'JP', name: '日本', flag: '🇯🇵' },
  { code: 'AU', name: '澳大利亚', flag: '🇦🇺' },
  { code: 'CA', name: '加拿大', flag: '🇨🇦' },
  { code: 'SG', name: '新加坡', flag: '🇸🇬' },
  { code: 'ES', name: '西班牙', flag: '🇪🇸' },
  { code: 'IT', name: '意大利', flag: '🇮🇹' },
  { code: 'NL', name: '荷兰', flag: '🇳🇱' }
]

const languageLabels: Record<string, string> = {
  'zh-CN': '简体中文',
  'en-US': 'English',
  'ja-JP': '日本語',
  'de-DE': 'Deutsch',
  'fr-FR': 'Français',
  'es-ES': 'Español'
}

const getLanguageLabel = (lang: string) => {
  return languageLabels[lang] || lang
}

const getFlagEmoji = (code: string) => {
  const market = availableMarkets.find(m => m.code === code)
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
  } catch (error: any) {
    ElMessage.error(error.message || '加载市场列表失败')
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
        ioss_enabled: market.tax_rules?.ioss_enabled || false,
        include_tax: market.tax_rules?.include_tax || false
      }
    })
    dialogVisible.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '获取市场信息失败')
  }
}

const handleDelete = async (row: Market) => {
  if (row.is_default) {
    ElMessage.warning('主市场不能删除')
    return
  }

  try {
    await ElMessageBox.confirm(`确认删除市场 "${row.name}"?`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteMarket(row.id)
    ElMessage.success('删除成功')
    loadMarkets()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleStatusChange = async (row: Market, val: boolean) => {
  try {
    await updateMarket({
      id: row.id,
      is_active: val
    })
    ElMessage.success(val ? '市场已启用' : '市场已禁用')
  } catch (error: any) {
    row.is_active = !val
    ElMessage.error(error.message || '状态更新失败')
  }
}

const handleSetDefault = async (row: Market) => {
  try {
    await ElMessageBox.confirm(`确认将 "${row.name}" 设为主市场?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })

    // Update the market to be default
    await updateMarket({
      id: row.id,
      is_default: true
    })

    ElMessage.success('主市场设置成功')
    loadMarkets()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '设置主市场失败')
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
          vat_rate: marketForm.tax_rules?.vat_rate?.toString(),
          gst_rate: marketForm.tax_rules?.gst_rate?.toString(),
          ioss_enabled: marketForm.tax_rules?.ioss_enabled,
          include_tax: marketForm.tax_rules?.include_tax
        }
      }
      await updateMarket(updateData)
      ElMessage.success('更新成功')
    } else {
      // Create new market
      const createData: CreateMarketRequest = {
        code: marketForm.code,
        name: marketForm.name,
        currency: marketForm.currency,
        default_language: marketForm.default_language,
        tax_rules: {
          vat_rate: marketForm.tax_rules?.vat_rate?.toString(),
          gst_rate: marketForm.tax_rules?.gst_rate?.toString(),
          ioss_enabled: marketForm.tax_rules?.ioss_enabled,
          include_tax: marketForm.tax_rules?.include_tax
        }
      }
      await createMarket(createData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadMarkets()
  } catch (error: any) {
    if (error !== false) {
      ElMessage.error(error.message || '保存失败')
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