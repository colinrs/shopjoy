<template>
  <div class="default-preview">
    <div class="preview-content">
      <el-icon size="24">
        <component :is="icon" />
      </el-icon>
      <span class="block-name">{{ blockName }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { VideoPlay, Menu, Document } from '@element-plus/icons-vue'
import { BLOCK_TYPES } from '@/api/storefront'

const { t } = useI18n()

const props = defineProps<{
  config: Record<string, any>
  blockType?: string
}>()

const icon = computed(() => {
  const icons: Record<string, any> = {
    video: VideoPlay,
    categories: Menu,
    custom_html: Document,
    image_carousel: Document
  }
  return icons[props.blockType || ''] || Document
})

const blockName = computed(() => {
  return BLOCK_TYPES.find(b => b.type === props.blockType)?.name || props.blockType || t('storefront.block')
})
</script>

<style scoped>
.default-preview {
  min-height: 80px;
}

.preview-content {
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: linear-gradient(135deg, #F9FAFB 0%, #F3F4F6 100%);
  border-radius: 8px;
  color: #9CA3AF;
}

.block-name {
  font-size: 13px;
}
</style>