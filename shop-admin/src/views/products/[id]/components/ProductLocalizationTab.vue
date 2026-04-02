<template>
  <div class="localization-section" v-loading="localizationsLoading">
    <div class="section-header">
      <h3 class="section-title">{{ $t('products.multilingualContent') }}</h3>
      <el-button type="primary" size="small" @click="handleAddLocalization">
        <el-icon><Plus /></el-icon>
        {{ $t('products.addLanguage') }}
      </el-button>
    </div>

    <!-- Language Tabs -->
    <el-tabs v-model="activeLanguage" type="card" class="language-tabs">
      <el-tab-pane
        v-for="lang in supportedLanguages"
        :key="lang.code"
        :label="lang.name"
        :name="lang.code"
      >
        <div v-if="getLocalizationByLang(lang.code)" class="localization-content">
          <el-form label-width="100px">
            <el-form-item :label="$t('products.productName')">
              <el-input
                v-model="getLocalizationByLang(lang.code)!.name"
                :placeholder="$t('products.enterLocalizedName')"
                @change="handleUpdateLocalization(getLocalizationByLang(lang.code)!)"
              />
            </el-form-item>
            <el-form-item :label="$t('products.productDescription')">
              <el-input
                v-model="getLocalizationByLang(lang.code)!.description"
                type="textarea"
                :rows="4"
                :placeholder="$t('products.enterLocalizedDescription')"
                @change="handleUpdateLocalization(getLocalizationByLang(lang.code)!)"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" size="small" @click="handleDeleteLocalization(getLocalizationByLang(lang.code)!)">
                {{ $t('products.deleteThisLanguage') }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <div v-else class="no-localization">
          <el-empty :description="$t('products.noContentForLang', { name: lang.name })">
            <el-button type="primary" size="small" @click="handleCreateLocalization(lang.code)">
              {{ $t('products.createContentForLang', { name: lang.name }) }}
            </el-button>
          </el-empty>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getProductLocalizations,
  createProductLocalization,
  updateProductLocalization,
  deleteProductLocalization,
  type ProductLocalization
} from '@/api/product'
import { t } from '@/plugins/i18n'
import type { ProductLocalizationTabProps, ProductLocalizationTabEmits, SupportedLanguage } from '../types'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<ProductLocalizationTabProps>()
const { handleError } = useErrorHandler()
const emit = defineEmits<ProductLocalizationTabEmits>()

const localizationsLoading = ref(false)
const localizations = ref<ProductLocalization[]>([])
const activeLanguage = ref('en')

const languageNameKeys: Record<string, string> = {
  'en': 'products.langEn',
  'zh-CN': 'products.langZhCN',
  'ja': 'products.langJa',
  'ko': 'products.langKo',
  'es': 'products.langEs'
}

const supportedLanguages = computed<SupportedLanguage[]>(() => [
  { code: 'en', name: t(languageNameKeys['en']) },
  { code: 'zh-CN', name: t(languageNameKeys['zh-CN']) },
  { code: 'ja', name: t(languageNameKeys['ja']) },
  { code: 'ko', name: t(languageNameKeys['ko']) },
  { code: 'es', name: t(languageNameKeys['es']) }
])

const loadLocalizations = async () => {
  localizationsLoading.value = true
  try {
    const response = await getProductLocalizations(props.productId)
    localizations.value = response.list || []
    // Set first available language as active
    if (localizations.value.length > 0) {
      activeLanguage.value = localizations.value[0].language_code
    }
  } catch (error) {
    handleError(error, t('products.loadLocalizationsFailed'))
  } finally {
    localizationsLoading.value = false
  }
}

const getLocalizationByLang = (langCode: string) => {
  return localizations.value.find(loc => loc.language_code === langCode)
}

const handleAddLocalization = () => {
  // Find first language without localization
  const existingLangs = localizations.value.map(loc => loc.language_code)
  const missingLang = supportedLanguages.value.find(lang => !existingLangs.includes(lang.code))
  if (missingLang) {
    activeLanguage.value = missingLang.code
  }
}

const handleCreateLocalization = async (langCode: string) => {
  try {
    await createProductLocalization({
      product_id: props.productId,
      language_code: langCode,
      name: props.productName,
      description: props.productDescription
    })
    ElMessage.success(t('products.createSuccess'))
    loadLocalizations()
    emit('localizations-change')
  } catch (error) {
    handleError(error, t('products.createFailed'))
  }
}

const handleUpdateLocalization = async (loc: ProductLocalization) => {
  try {
    await updateProductLocalization({
      id: loc.id,
      name: loc.name,
      description: loc.description
    })
    ElMessage.success(t('products.updateSuccess'))
  } catch (error) {
    handleError(error, t('products.updateFailed'))
    loadLocalizations() // Reload to reset
  }
}

const handleDeleteLocalization = async (loc: ProductLocalization) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDeleteLocalization', { name: supportedLanguages.value.find(l => l.code === loc.language_code)?.name || loc.language_code }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteProductLocalization(loc.id)
    ElMessage.success(t('products.deleteSuccess'))
    loadLocalizations()
    emit('localizations-change')
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('products.deleteFailed'))
    }
  }
}

onMounted(() => {
  loadLocalizations()
})

defineExpose({
  loadLocalizations
})
</script>

<style scoped>
.localization-section {
  padding: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.language-tabs {
  margin-top: 16px;
}

.localization-content {
  padding: 16px 0;
}

.no-localization {
  padding: 20px 0;
}
</style>
