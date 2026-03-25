<template>
  <div class="themes-page">
    <page-header title="主题管理" subtitle="选择并自定义您的店铺主题风格">
      <template #extra>
        <el-button type="primary" @click="showCustomThemeDialog">
          <el-icon><Plus /></el-icon>
          上传自定义主题
        </el-button>
      </template>
    </page-header>

    <!-- Current Theme Section -->
    <el-card v-if="currentTheme" class="current-theme-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>当前主题</span>
          <el-tag type="success" effect="dark" size="small">使用中</el-tag>
        </div>
      </template>
      <div class="current-theme-content">
        <div class="theme-preview">
          <el-image
            :src="currentTheme.theme.preview_image || currentTheme.theme.thumbnail"
            fit="cover"
            class="preview-image"
          >
            <template #error>
              <div class="image-placeholder">
                <el-icon size="48"><Picture /></el-icon>
              </div>
            </template>
          </el-image>
        </div>
        <div class="theme-info">
          <h3>{{ currentTheme.theme.name }}</h3>
          <p class="theme-desc">{{ currentTheme.theme.description }}</p>

          <!-- Theme Config -->
          <div class="theme-config">
            <h4>主题配置</h4>
            <el-form :model="configForm" label-width="100px" label-position="left">
              <el-form-item label="主色调">
                <div class="color-picker-wrapper">
                  <el-color-picker v-model="configForm.primary_color" />
                  <span class="color-value">{{ configForm.primary_color }}</span>
                </div>
              </el-form-item>
              <el-form-item label="辅助色">
                <div class="color-picker-wrapper">
                  <el-color-picker v-model="configForm.secondary_color" />
                  <span class="color-value">{{ configForm.secondary_color }}</span>
                </div>
              </el-form-item>
              <el-form-item label="字体">
                <el-select v-model="configForm.font_family" placeholder="选择字体">
                  <el-option label="Inter" value="inter" />
                  <el-option label="Roboto" value="roboto" />
                  <el-option label="Open Sans" value="opensans" />
                  <el-option label="Poppins" value="poppins" />
                  <el-option label="Noto Sans" value="notosans" />
                </el-select>
              </el-form-item>
              <el-form-item label="按钮样式">
                <el-select v-model="configForm.button_style" placeholder="选择按钮样式">
                  <el-option label="圆角按钮" value="rounded" />
                  <el-option label="胶囊按钮" value="pill" />
                  <el-option label="方角按钮" value="square" />
                  <el-option label="下划线样式" value="underline" />
                </el-select>
              </el-form-item>
            </el-form>
            <div class="config-actions">
              <el-button type="primary" @click="saveConfig" :loading="configLoading">
                保存配置
              </el-button>
              <el-button @click="resetConfig">重置为默认</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- Available Themes -->
    <div class="themes-section">
      <h3 class="section-title">可用主题</h3>
      <div class="themes-grid" v-loading="loading">
        <div
          v-for="theme in themes"
          :key="theme.id"
          class="theme-card"
          :class="{ 'is-current': theme.is_current }"
          @click="previewTheme(theme)"
        >
          <div class="theme-thumbnail">
            <el-image
              :src="theme.thumbnail"
              fit="cover"
              class="thumbnail-image"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon size="32"><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="theme-overlay">
              <el-button type="primary" size="small" @click.stop="previewTheme(theme)">
                预览
              </el-button>
            </div>
            <div v-if="theme.is_current" class="current-badge">
              <el-icon><Check /></el-icon>
              使用中
            </div>
            <div v-if="theme.is_preset" class="preset-badge">官方主题</div>
          </div>
          <div class="theme-meta">
            <h4>{{ theme.name }}</h4>
            <p>{{ theme.description }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Theme Preview Dialog -->
    <el-dialog
      v-model="previewDialogVisible"
      :title="previewThemeData?.name + ' - 主题预览'"
      width="900px"
      destroy-on-close
    >
      <div class="preview-dialog-content" v-if="previewThemeData">
        <el-image
          :src="previewThemeData.preview_image || previewThemeData.thumbnail"
          fit="contain"
          class="preview-large-image"
        >
          <template #error>
            <div class="image-placeholder large">
              <el-icon size="64"><Picture /></el-icon>
              <span>暂无预览图</span>
            </div>
          </template>
        </el-image>
        <div class="preview-info">
          <h3>{{ previewThemeData.name }}</h3>
          <p>{{ previewThemeData.description }}</p>
          <div class="preview-actions">
            <el-button
              v-if="!previewThemeData.is_current"
              type="primary"
              size="large"
              @click="applyTheme(previewThemeData.id)"
              :loading="switchLoading"
            >
              应用此主题
            </el-button>
            <el-tag v-else type="success" size="large">当前使用</el-tag>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Picture, Check } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  listThemes,
  getCurrentTheme,
  switchTheme,
  updateThemeConfig,
  type ThemeItem,
  type CurrentThemeResponse
} from '@/api/storefront'

const loading = ref(false)
const configLoading = ref(false)
const switchLoading = ref(false)
const themes = ref<ThemeItem[]>([])
const currentTheme = ref<CurrentThemeResponse | null>(null)
const previewDialogVisible = ref(false)
const previewThemeData = ref<ThemeItem | null>(null)

const configForm = reactive({
  primary_color: '#6366F1',
  secondary_color: '#818CF8',
  font_family: 'inter',
  button_style: 'rounded'
})

const fetchThemes = async () => {
  loading.value = true
  try {
    const res = await listThemes()
    themes.value = res.themes || []
  } catch (error) {
    ElMessage.error('获取主题列表失败')
  } finally {
    loading.value = false
  }
}

const fetchCurrentTheme = async () => {
  try {
    const res = await getCurrentTheme()
    currentTheme.value = res
    if (res.config) {
      configForm.primary_color = res.config.primary_color || '#6366F1'
      configForm.secondary_color = res.config.secondary_color || '#818CF8'
      configForm.font_family = res.config.font_family || 'inter'
      configForm.button_style = res.config.button_style || 'rounded'
    }
  } catch (error) {
    ElMessage.error('获取当前主题失败')
  }
}

const previewTheme = (theme: ThemeItem) => {
  previewThemeData.value = theme
  previewDialogVisible.value = true
}

const applyTheme = async (themeId: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要切换到此主题吗？切换后您的自定义配置将被重置。',
      '切换主题',
      {
        confirmButtonText: '确定切换',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    switchLoading.value = true
    await switchTheme({ theme_id: themeId })
    ElMessage.success('主题切换成功')
    previewDialogVisible.value = false
    await fetchThemes()
    await fetchCurrentTheme()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '主题切换失败')
    }
  } finally {
    switchLoading.value = false
  }
}

const saveConfig = async () => {
  configLoading.value = true
  try {
    await updateThemeConfig({ config: { ...configForm } })
    ElMessage.success('主题配置已保存')
  } catch (error: any) {
    ElMessage.error(error.message || '保存配置失败')
  } finally {
    configLoading.value = false
  }
}

const resetConfig = () => {
  if (currentTheme.value?.theme.code === 'classic') {
    configForm.primary_color = '#3B82F6'
    configForm.secondary_color = '#1E40AF'
  } else if (currentTheme.value?.theme.code === 'modern') {
    configForm.primary_color = '#10B981'
    configForm.secondary_color = '#059669'
  } else if (currentTheme.value?.theme.code === 'minimal') {
    configForm.primary_color = '#000000'
    configForm.secondary_color = '#6B7280'
  } else if (currentTheme.value?.theme.code === 'bold') {
    configForm.primary_color = '#8B5CF6'
    configForm.secondary_color = '#6D28D9'
  } else if (currentTheme.value?.theme.code === 'nature') {
    configForm.primary_color = '#059669'
    configForm.secondary_color = '#047857'
  }
  configForm.font_family = 'inter'
  configForm.button_style = 'rounded'
}

const showCustomThemeDialog = () => {
  ElMessage.info('自定义主题功能即将上线')
}

onMounted(() => {
  fetchThemes()
  fetchCurrentTheme()
})
</script>

<style scoped>
.themes-page {
  padding: 0;
}

.current-theme-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.current-theme-content {
  display: flex;
  gap: 32px;
}

.theme-preview {
  flex-shrink: 0;
}

.preview-image {
  width: 400px;
  height: 250px;
  border-radius: 12px;
  overflow: hidden;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #E0E7FF 100%);
  color: #A5B4FC;
}

.theme-info {
  flex: 1;
}

.theme-info h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
}

.theme-desc {
  color: #6B7280;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 24px 0;
}

.theme-config {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 20px;
}

.theme-config h4 {
  margin: 0 0 16px 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.color-picker-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.color-value {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6B7280;
}

.config-actions {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 16px 0;
}

.themes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.theme-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.theme-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.theme-card.is-current {
  border-color: var(--color-primary);
}

.theme-thumbnail {
  position: relative;
  height: 160px;
  overflow: hidden;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
}

.theme-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.theme-card:hover .theme-overlay {
  opacity: 1;
}

.current-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: var(--color-primary);
  color: #fff;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 4px;
}

.preset-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 11px;
}

.theme-meta {
  padding: 16px;
}

.theme-meta h4 {
  margin: 0 0 8px 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.theme-meta p {
  margin: 0;
  font-size: 13px;
  color: #6B7280;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.preview-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.preview-large-image {
  width: 100%;
  max-height: 400px;
  border-radius: 12px;
  overflow: hidden;
}

.image-placeholder.large {
  height: 300px;
  flex-direction: column;
  gap: 12px;
}

.preview-info {
  text-align: center;
}

.preview-info h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
}

.preview-info p {
  margin: 0 0 20px 0;
  color: #6B7280;
  font-size: 14px;
}

.preview-actions {
  display: flex;
  justify-content: center;
}
</style>