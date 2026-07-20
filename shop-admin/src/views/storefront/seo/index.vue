<template>
  <div class="seo-page">
    <page-header
      :title="$t('storefront.seoSettings')"
      :subtitle="$t('storefront.configureStoreSEO')"
    />

    <el-card
      shadow="hover"
      class="seo-card"
    >
      <template #header>
        <div class="card-header">
          <span>{{ $t('storefront.globalSEOConfig') }}</span>
          <el-tag
            type="info"
            size="small"
          >
            {{ $t('storefront.appliesToAll') }}
          </el-tag>
        </div>
      </template>
      <el-form
        ref="globalFormRef"
        v-loading="globalLoading"
        :model="globalSEO"
        label-width="100px"
        label-position="left"
      >
        <el-form-item
          :label="$t('storefront.websiteTitle')"
          prop="title"
        >
          <el-input
            v-model="globalSEO.title"
            :placeholder="$t('storefront.websiteTitlePlaceholder')"
            maxlength="60"
            show-word-limit
          />
        </el-form-item>
        <el-form-item
          :label="$t('storefront.websiteDescription')"
          prop="description"
        >
          <el-input
            v-model="globalSEO.description"
            type="textarea"
            :rows="4"
            :placeholder="$t('storefront.websiteDescriptionPlaceholder')"
            maxlength="160"
            show-word-limit
          />
        </el-form-item>
        <el-form-item
          :label="$t('storefront.keywords')"
          prop="keywords"
        >
          <el-input
            v-model="globalSEO.keywords"
            :placeholder="$t('storefront.keywordsPlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="globalSaving"
            @click="saveGlobalSEO"
          >
            {{ $t('storefront.saveGlobalConfig') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getGlobalSEO,
  updateGlobalSEO,
  type SEOConfigDTO
} from '@/api/storefront'

const { t } = useI18n()

const globalLoading = ref(false)
const globalSaving = ref(false)

const globalSEO = reactive<SEOConfigDTO>({
  title: '',
  description: '',
  keywords: ''
})

const fetchGlobalSEO = async () => {
  globalLoading.value = true
  try {
    const res = await getGlobalSEO()
    globalSEO.title = res.title || ''
    globalSEO.description = res.description || ''
    globalSEO.keywords = res.keywords || ''
  } catch (error) {
    ElMessage.error(t('storefront.loadGlobalSEOFailed'))
  } finally {
    globalLoading.value = false
  }
}

const saveGlobalSEO = async () => {
  globalSaving.value = true
  try {
    await updateGlobalSEO({ ...globalSEO })
    ElMessage.success(t('storefront.globalSEOSaved'))
  } catch (error: unknown) {
    ElMessage.error((error as Error).message || t('storefront.saveFailed'))
  } finally {
    globalSaving.value = false
  }
}

onMounted(() => {
  fetchGlobalSEO()
})
</script>

<style scoped>
.seo-page {
  padding: 0;
}

.seo-card {
  border-radius: 12px;
  margin-top: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
