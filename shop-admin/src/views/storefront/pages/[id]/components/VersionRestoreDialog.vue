<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="700px"
    destroy-on-close
  >
    <div class="version-restore-content" v-loading="loading">
      <!-- Version Info -->
      <div class="version-info" v-if="version">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('storefront.version')">
            <el-tag size="small">v{{ version.version }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('storefront.createdTime')">
            {{ formatTime(version.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Block Preview -->
      <div class="block-preview" v-if="versionDetail">
        <h4>{{ $t('storefront.blockPreview') }}</h4>
        <div class="preview-blocks">
          <div
            v-for="(block, index) in versionDetail.blocks"
            :key="index"
            class="preview-block"
          >
            <div class="block-header">
              <el-tag size="small" type="info">{{ block.block_type }}</el-tag>
              <span class="block-order">#{{ index + 1 }}</span>
            </div>
            <pre class="block-config">{{ JSON.stringify(block.block_config, null, 2) }}</pre>
          </div>
          <el-empty v-if="versionDetail.blocks.length === 0" :description="$t('storefront.noBlocks')" />
        </div>
      </div>

      <!-- JSON Raw View -->
      <div class="json-view" v-if="versionDetail">
        <h4>{{ $t('storefront.rawJson') }}</h4>
        <pre class="json-content">{{ JSON.stringify(versionDetail, null, 2) }}</pre>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRestore" :loading="restoring">
          {{ $t('storefront.restoreThisVersion') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getVersion, restoreVersion, type VersionItem, type VersionDetailResponse } from '@/api/storefront'

const { t } = useI18n()

const props = defineProps<{
  modelValue: boolean
  pageId: number
  version: VersionItem | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'restored'): void
}>()

const loading = ref(false)
const restoring = ref(false)
const versionDetail = ref<VersionDetailResponse | null>(null)

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const dialogTitle = computed(() => {
  if (props.version) {
    return `${t('storefront.versionDetail')} - v${props.version.version}`
  }
  return t('storefront.versionDetail')
})

const formatTime = (timestampStr: string) => {
  return new Date(timestampStr).toLocaleString()
}

const loadVersionDetail = async () => {
  if (!props.version || !props.pageId) return

  loading.value = true
  try {
    versionDetail.value = await getVersion(props.pageId, props.version.version)
  } catch (error) {
    console.error('Failed to load version detail:', error)
    ElMessage.error(t('storefront.loadVersionDetailFailed'))
  } finally {
    loading.value = false
  }
}

const handleRestore = async () => {
  if (!props.version || !props.pageId) return

  try {
    restoring.value = true
    await restoreVersion(props.pageId, { version: props.version.version })
    ElMessage.success(t('storefront.versionRestored'))
    visible.value = false
    emit('restored')
  } catch (error: unknown) {
    console.error('Failed to restore version:', error)
    ElMessage.error((error as Error).message || t('storefront.restoreFailed'))
  } finally {
    restoring.value = false
  }
}

watch(() => props.modelValue, (val) => {
  if (val && props.version) {
    loadVersionDetail()
  }
})

watch(() => props.version, (val) => {
  if (val && visible.value) {
    loadVersionDetail()
  }
})
</script>

<style scoped>
.version-restore-content {
  max-height: 600px;
  overflow-y: auto;
}

.version-info {
  margin-bottom: 20px;
}

.block-preview {
  margin-bottom: 20px;
}

.block-preview h4,
.json-view h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
}

.preview-blocks {
  max-height: 300px;
  overflow-y: auto;
}

.preview-block {
  background: #F9FAFB;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.block-order {
  font-size: 12px;
  color: #9CA3AF;
}

.block-config {
  margin: 0;
  padding: 12px;
  background: #1E1B4B;
  border-radius: 6px;
  color: #E5E7EB;
  font-size: 12px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.json-view {
  margin-top: 20px;
}

.json-content {
  margin: 0;
  padding: 16px;
  background: #1E1B4B;
  border-radius: 8px;
  color: #E5E7EB;
  font-size: 12px;
  overflow-x: auto;
  max-height: 400px;
  white-space: pre-wrap;
  word-break: break-all;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
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
