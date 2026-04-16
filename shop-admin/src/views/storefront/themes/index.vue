<template>
  <div class="themes-page">
    <page-header
      :title="$t('storefront.themeManagement')"
      :subtitle="$t('storefront.selectAndCustomizeTheme')"
    >
      <template #extra>
        <el-button
          type="primary"
          @click="showCustomThemeDialog"
        >
          <el-icon><Plus /></el-icon>
          {{ $t('storefront.uploadCustomTheme') }}
        </el-button>
      </template>
    </page-header>

    <!-- Current Theme Section -->
    <el-card
      v-if="currentTheme"
      class="current-theme-card"
      shadow="hover"
    >
      <template #header>
        <div class="card-header">
          <span>{{ $t('storefront.currentTheme') }}</span>
          <el-tag
            type="success"
            effect="dark"
            size="small"
          >
            {{ $t('storefront.inUse') }}
          </el-tag>
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
                <el-icon size="48">
                  <Picture />
                </el-icon>
              </div>
            </template>
          </el-image>
        </div>
        <div class="theme-info">
          <h3>{{ currentTheme.theme.name }}</h3>
          <p class="theme-desc">
            {{ currentTheme.theme.description }}
          </p>

          <!-- Theme Config -->
          <div class="theme-config">
            <h4>{{ $t('storefront.themeConfig') }}</h4>
            <el-form
              :model="configForm"
              label-width="100px"
              label-position="left"
            >
              <el-form-item :label="$t('storefront.primaryColor')">
                <div class="color-picker-wrapper">
                  <el-color-picker v-model="configForm.primary_color" />
                  <span class="color-value">{{ configForm.primary_color }}</span>
                </div>
              </el-form-item>
              <el-form-item :label="$t('storefront.secondaryColor')">
                <div class="color-picker-wrapper">
                  <el-color-picker v-model="configForm.secondary_color" />
                  <span class="color-value">{{ configForm.secondary_color }}</span>
                </div>
              </el-form-item>
              <el-form-item :label="$t('storefront.font')">
                <el-select
                  v-model="configForm.font_family"
                  :placeholder="$t('storefront.selectFont')"
                >
                  <el-option
                    label="Inter"
                    value="inter"
                  />
                  <el-option
                    label="Roboto"
                    value="roboto"
                  />
                  <el-option
                    label="Open Sans"
                    value="opensans"
                  />
                  <el-option
                    label="Poppins"
                    value="poppins"
                  />
                  <el-option
                    label="Noto Sans"
                    value="notosans"
                  />
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('storefront.buttonStyle')">
                <el-select
                  v-model="configForm.button_style"
                  :placeholder="$t('storefront.selectButtonStyle')"
                >
                  <el-option
                    :label="$t('storefront.roundedButton')"
                    value="rounded"
                  />
                  <el-option
                    :label="$t('storefront.pillButton')"
                    value="pill"
                  />
                  <el-option
                    :label="$t('storefront.squareButton')"
                    value="square"
                  />
                  <el-option
                    :label="$t('storefront.underlineButton')"
                    value="underline"
                  />
                </el-select>
              </el-form-item>
            </el-form>
            <div class="config-actions">
              <el-button
                type="primary"
                :loading="configLoading"
                @click="saveConfig"
              >
                {{ $t('storefront.saveConfig') }}
              </el-button>
              <el-button @click="resetConfig">
                {{ $t('storefront.resetToDefault') }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- Tabs Section -->
    <el-card
      class="tabs-card"
      shadow="never"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane
          :label="$t('storefront.availableThemes')"
          name="themes"
        >
          <div
            v-loading="loading"
            class="themes-grid"
          >
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
                      <el-icon size="32">
                        <Picture />
                      </el-icon>
                    </div>
                  </template>
                </el-image>
                <div class="theme-overlay">
                  <el-button
                    type="primary"
                    size="small"
                    @click.stop="previewTheme(theme)"
                  >
                    {{ $t('storefront.preview') }}
                  </el-button>
                </div>
                <div
                  v-if="theme.is_current"
                  class="current-badge"
                >
                  <el-icon><Check /></el-icon>
                  {{ $t('storefront.inUse') }}
                </div>
                <div
                  v-if="theme.is_preset"
                  class="preset-badge"
                >
                  {{ $t('storefront.officialTheme') }}
                </div>
              </div>
              <div class="theme-meta">
                <h4>{{ theme.name }}</h4>
                <p>{{ theme.description }}</p>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane
          :label="$t('storefront.auditLogs')"
          name="audit"
        >
          <el-table
            v-loading="auditLoading"
            :data="auditLogs"
            stripe
          >
            <el-table-column
              :label="$t('storefront.action')"
              width="150"
            >
              <template #default="{ row }">
                {{ getActionText(row.action) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('storefront.theme')"
              width="200"
            >
              <template #default="{ row }">
                {{ row.theme_name }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('storefront.user')"
              width="150"
            >
              <template #default="{ row }">
                {{ row.user_name }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.createdAt')">
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-if="!auditLoading && auditLogs.length === 0"
            :description="$t('storefront.noAuditLogs')"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Theme Preview Dialog -->
    <el-dialog
      v-model="previewDialogVisible"
      :title="(previewThemeData?.name || '') + ' - ' + $t('storefront.themePreview')"
      width="900px"
      destroy-on-close
    >
      <div
        v-if="previewThemeData"
        class="preview-dialog-content"
      >
        <el-image
          :src="previewThemeData.preview_image || previewThemeData.thumbnail"
          fit="contain"
          class="preview-large-image"
        >
          <template #error>
            <div class="image-placeholder large">
              <el-icon size="64">
                <Picture />
              </el-icon>
              <span>{{ $t('storefront.noPreviewImage') }}</span>
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
              :loading="switchLoading"
              @click="applyTheme(previewThemeData.id)"
            >
              {{ $t('storefront.applyTheme') }}
            </el-button>
            <el-tag
              v-else
              type="success"
              size="large"
            >
              {{ $t('storefront.currentlyInUse') }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Picture, Check } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  listThemes,
  getCurrentTheme,
  switchTheme,
  updateThemeConfig,
  getThemeAuditLogs,
  type ThemeItem,
  type CurrentThemeResponse,
  type ThemeAuditLog
} from '@/api/storefront'

const { t } = useI18n()

const activeTab = ref('themes')
const loading = ref(false)
const configLoading = ref(false)
const switchLoading = ref(false)
const themes = ref<ThemeItem[]>([])
const currentTheme = ref<CurrentThemeResponse | null>(null)
const previewDialogVisible = ref(false)
const previewThemeData = ref<ThemeItem | null>(null)
const auditLogs = ref<ThemeAuditLog[]>([])
const auditLoading = ref(false)
const isMounted = ref(false)

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
    if (isMounted.value) {
      themes.value = res.themes || []
    }
  } catch (error) {
    if (isMounted.value) {
      ElMessage.error(t('storefront.loadThemesFailed'))
    }
  } finally {
    if (isMounted.value) {
      loading.value = false
    }
  }
}

const fetchCurrentTheme = async () => {
  try {
    const res = await getCurrentTheme()
    if (isMounted.value) {
      currentTheme.value = res
      if (res.config) {
        configForm.primary_color = res.config.primary_color || '#6366F1'
        configForm.secondary_color = res.config.secondary_color || '#818CF8'
        configForm.font_family = res.config.font_family || 'inter'
        configForm.button_style = res.config.button_style || 'rounded'
      }
    }
  } catch (error) {
    if (isMounted.value) {
      ElMessage.error(t('storefront.loadCurrentThemeFailed'))
    }
  }
}

const previewTheme = (theme: ThemeItem) => {
  previewThemeData.value = theme
  previewDialogVisible.value = true
}

const applyTheme = async (themeId: number) => {
  try {
    await ElMessageBox.confirm(
      t('storefront.confirmSwitchTheme'),
      t('storefront.switchTheme'),
      {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
    switchLoading.value = true
    await switchTheme({ theme_id: themeId })
    if (isMounted.value) {
      ElMessage.success(t('storefront.themeSwitchSuccess'))
      previewDialogVisible.value = false
      await fetchThemes()
      await fetchCurrentTheme()
    }
  } catch (error: unknown) {
    if (error !== 'cancel' && isMounted.value) {
      ElMessage.error((error as Error).message || t('storefront.themeSwitchFailed'))
    }
  } finally {
    if (isMounted.value) {
      switchLoading.value = false
    }
  }
}

const saveConfig = async () => {
  configLoading.value = true
  try {
    await updateThemeConfig({ config: { ...configForm } })
    ElMessage.success(t('storefront.themeConfigSaved'))
  } catch (error: unknown) {
    ElMessage.error((error as Error).message || t('storefront.saveConfigFailed'))
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
  ElMessage.info(t('storefront.customThemeComingSoon'))
}

const fetchAuditLogs = async () => {
  auditLoading.value = true
  try {
    const res = await getThemeAuditLogs()
    if (isMounted.value) {
      auditLogs.value = res.logs || []
    }
  } catch (error) {
    if (isMounted.value) {
      ElMessage.error(t('storefront.loadAuditLogsFailed'))
    }
  } finally {
    if (isMounted.value) {
      auditLoading.value = false
    }
  }
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getActionText = (action: string) => {
  const texts: Record<string, string> = {
    'switch': t('storefront.actionSwitch'),
    'update_config': t('storefront.actionUpdateConfig')
  }
  return texts[action] || action
}

// Watch for tab changes to load audit logs
watch(activeTab, (newTab) => {
  if (newTab === 'audit' && auditLogs.value.length === 0) {
    fetchAuditLogs()
  }
})

onMounted(() => {
  isMounted.value = true
  fetchThemes()
  fetchCurrentTheme()
})

onUnmounted(() => {
  isMounted.value = false
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

.tabs-card {
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
