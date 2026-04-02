<template>
  <div class="rich-text-preview">
    <div v-if="config?.content" class="content" v-html="config.content"></div>
    <div v-else class="placeholder">
      <el-icon size="24"><Document /></el-icon>
      <span>{{ $t('storefront.richTextBlock') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'

defineProps<{
  config: Record<string, any>
}>()
</script>

<style scoped>
.rich-text-preview {
  min-height: 80px;
}

.content {
  padding: 8px;
  line-height: 1.6;
  color: #374151;
}

.content :deep(p) {
  margin: 0 0 8px 0;
}

.placeholder {
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #F9FAFB;
  border-radius: 8px;
  color: #9CA3AF;
}
</style>