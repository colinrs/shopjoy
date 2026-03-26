<template>
  <div class="seo-page">
    <page-header title="SEO设置" subtitle="配置店铺的全局和页面级SEO信息">
    </page-header>

    <el-tabs v-model="activeTab" class="seo-tabs">
      <!-- Global SEO -->
      <el-tab-pane label="全局SEO" name="global">
        <el-card shadow="hover" class="seo-card">
          <template #header>
            <div class="card-header">
              <span>全局SEO配置</span>
              <el-tag type="info" size="small">适用于全站</el-tag>
            </div>
          </template>
          <el-form
            ref="globalFormRef"
            :model="globalSEO"
            label-width="100px"
            label-position="left"
            v-loading="globalLoading"
          >
            <el-form-item label="网站标题" prop="title">
              <el-input
                v-model="globalSEO.title"
                placeholder="网站标题，建议不超过60个字符"
                maxlength="60"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="网站描述" prop="description">
              <el-input
                v-model="globalSEO.description"
                type="textarea"
                :rows="4"
                placeholder="网站描述，建议不超过160个字符"
                maxlength="160"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="关键词" prop="keywords">
              <el-input
                v-model="globalSEO.keywords"
                placeholder="关键词，用英文逗号分隔"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveGlobalSEO" :loading="globalSaving">
                保存全局配置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- Page SEO -->
      <el-tab-pane label="页面SEO" name="pages">
        <el-card shadow="hover" class="seo-card" v-loading="pageLoading">
          <div class="page-seo-list">
            <div
              v-for="item in pageSEOList"
              :key="`${item.page_type}-${item.page_id || 'default'}`"
              class="page-seo-item"
            >
              <div class="item-header" @click="toggleExpand(item)">
                <div class="item-info">
                  <span class="page-name">{{ getPageTypeName(item.page_type) }}</span>
                  <el-tag size="small" :type="getPageTypeTag(item.page_type)">
                    {{ item.page_type }}
                  </el-tag>
                </div>
                <el-icon class="expand-icon" :class="{ 'is-expanded': expandedItems.has(`${item.page_type}-${item.page_id || 'default'}`) }">
                  <ArrowDown />
                </el-icon>
              </div>
              <el-collapse-transition>
                <div
                  v-show="expandedItems.has(`${item.page_type}-${item.page_id || 'default'}`)"
                  class="item-content"
                >
                  <el-form
                    :model="item.config"
                    label-width="80px"
                    label-position="left"
                    size="small"
                  >
                    <el-form-item label="页面标题">
                      <el-input v-model="item.config.title" placeholder="页面标题" />
                    </el-form-item>
                    <el-form-item label="页面描述">
                      <el-input
                        v-model="item.config.description"
                        type="textarea"
                        :rows="3"
                        placeholder="页面描述"
                      />
                    </el-form-item>
                    <el-form-item label="关键词">
                      <el-input v-model="item.config.keywords" placeholder="关键词" />
                    </el-form-item>
                    <el-form-item>
                      <el-button
                        type="primary"
                        size="small"
                        @click="savePageSEO(item)"
                        :loading="item.saving"
                      >
                        保存
                      </el-button>
                    </el-form-item>
                  </el-form>
                </div>
              </el-collapse-transition>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- SEO Preview -->
      <el-tab-pane label="预览效果" name="preview">
        <el-card shadow="hover" class="seo-card">
          <template #header>
            <span>搜索引擎预览</span>
          </template>
          <div class="seo-preview">
            <div class="preview-section">
              <h4>Google 搜索结果预览</h4>
              <div class="google-preview">
                <div class="google-title">{{ globalSEO.title || '您的网站标题' }}</div>
                <div class="google-url">https://yourstore.com</div>
                <div class="google-desc">{{ globalSEO.description || '您的网站描述将显示在这里...' }}</div>
              </div>
            </div>
            <el-divider />
            <div class="preview-section">
              <h4>SEO 评分建议</h4>
              <div class="seo-score">
                <el-progress
                  :percentage="seoScore"
                  :color="scoreColor"
                  :stroke-width="12"
                />
                <div class="score-tips">
                  <div v-for="tip in seoTips" :key="tip.text" class="tip-item" :class="tip.type">
                    <el-icon>
                      <component :is="tip.type === 'success' ? 'CircleCheck' : 'Warning'" />
                    </el-icon>
                    <span>{{ tip.text }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  getGlobalSEO,
  updateGlobalSEO,
  listPageSEO,
  updatePageSEO,
  type SEOConfigDTO,
  type PageSEOConfigDTO
} from '@/api/storefront'

const activeTab = ref('global')
const globalLoading = ref(false)
const globalSaving = ref(false)
const pageLoading = ref(false)

const globalSEO = reactive<SEOConfigDTO>({
  title: '',
  description: '',
  keywords: ''
})

const pageSEOList = ref<(PageSEOConfigDTO & { saving?: boolean })[]>([])
const expandedItems = ref(new Set<string>())

const fetchGlobalSEO = async () => {
  globalLoading.value = true
  try {
    const res = await getGlobalSEO()
    globalSEO.title = res.title || ''
    globalSEO.description = res.description || ''
    globalSEO.keywords = res.keywords || ''
  } catch (error) {
    ElMessage.error('获取全局SEO配置失败')
  } finally {
    globalLoading.value = false
  }
}

const fetchPageSEO = async () => {
  pageLoading.value = true
  try {
    const res = await listPageSEO()
    pageSEOList.value = (res.configs || []).map(item => ({
      ...item,
      config: {
        title: item.config?.title || '',
        description: item.config?.description || '',
        keywords: item.config?.keywords || ''
      }
    }))
  } catch (error) {
    ElMessage.error('获取页面SEO配置失败')
  } finally {
    pageLoading.value = false
  }
}

const saveGlobalSEO = async () => {
  globalSaving.value = true
  try {
    await updateGlobalSEO({ ...globalSEO })
    ElMessage.success('全局SEO配置已保存')
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    globalSaving.value = false
  }
}

const savePageSEO = async (item: PageSEOConfigDTO & { saving?: boolean }) => {
  item.saving = true
  try {
    await updatePageSEO(item.page_type, {
      title: item.config.title,
      description: item.config.description,
      keywords: item.config.keywords
    }, item.page_id)
    ElMessage.success('页面SEO配置已保存')
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    item.saving = false
  }
}

const toggleExpand = (item: PageSEOConfigDTO) => {
  const key = `${item.page_type}-${item.page_id || 'default'}`
  if (expandedItems.value.has(key)) {
    expandedItems.value.delete(key)
  } else {
    expandedItems.value.add(key)
  }
}

const getPageTypeName = (type: string) => {
  const names: Record<string, string> = {
    global: '全局设置',
    home: '首页',
    product: '商品详情页',
    category: '分类页面',
    collection: '集合页面',
    custom: '自定义页面'
  }
  return names[type] || type
}

const getPageTypeTag = (type: string) => {
  const tags: Record<string, string> = {
    global: 'danger',
    home: 'warning',
    product: 'success',
    category: 'info'
  }
  return tags[type] || 'info'
}

// SEO Score calculation
const seoScore = computed(() => {
  let score = 0
  if (globalSEO.title.length > 0 && globalSEO.title.length <= 60) score += 30
  if (globalSEO.description.length > 0 && globalSEO.description.length <= 160) score += 30
  if (globalSEO.keywords.length > 0) score += 20
  if (globalSEO.title.length > 30) score += 10
  if (globalSEO.description.length > 120) score += 10
  return score
})

const scoreColor = computed(() => {
  if (seoScore.value >= 80) return '#10B981'
  if (seoScore.value >= 50) return '#F59E0B'
  return '#EF4444'
})

const seoTips = computed(() => {
  const tips: { text: string; type: 'success' | 'warning' }[] = []
  if (globalSEO.title.length === 0) {
    tips.push({ text: '请添加网站标题', type: 'warning' })
  } else if (globalSEO.title.length > 60) {
    tips.push({ text: '标题过长，可能被截断', type: 'warning' })
  } else {
    tips.push({ text: '标题长度合适', type: 'success' })
  }

  if (globalSEO.description.length === 0) {
    tips.push({ text: '请添加网站描述', type: 'warning' })
  } else if (globalSEO.description.length > 160) {
    tips.push({ text: '描述过长，可能被截断', type: 'warning' })
  } else {
    tips.push({ text: '描述长度合适', type: 'success' })
  }

  if (globalSEO.keywords.length === 0) {
    tips.push({ text: '建议添加关键词', type: 'warning' })
  } else {
    tips.push({ text: '关键词已设置', type: 'success' })
  }

  return tips
})

onMounted(() => {
  fetchGlobalSEO()
  fetchPageSEO()
})
</script>

<style scoped>
.seo-page {
  padding: 0;
}

.seo-tabs {
  margin-top: 16px;
}

.seo-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-seo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.page-seo-item {
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  overflow: hidden;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #F9FAFB;
  cursor: pointer;
  transition: background 0.2s;
}

.item-header:hover {
  background: #F3F4F6;
}

.item-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-name {
  font-weight: 500;
  color: #1E1B4B;
}

.expand-icon {
  transition: transform 0.3s;
  color: #9CA3AF;
}

.expand-icon.is-expanded {
  transform: rotate(180deg);
}

.item-content {
  padding: 16px;
  border-top: 1px solid #E5E7EB;
}

.seo-preview {
  padding: 8px;
}

.preview-section {
  margin-bottom: 24px;
}

.preview-section h4 {
  margin: 0 0 16px 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.google-preview {
  padding: 16px;
  background: #fff;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
}

.google-title {
  font-size: 18px;
  color: #1a0dab;
  margin-bottom: 4px;
}

.google-url {
  font-size: 14px;
  color: #006621;
  margin-bottom: 8px;
}

.google-desc {
  font-size: 14px;
  color: #545454;
  line-height: 1.5;
}

.seo-score {
  max-width: 400px;
}

.score-tips {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.tip-item.success {
  color: #10B981;
}

.tip-item.warning {
  color: #F59E0B;
}
</style>