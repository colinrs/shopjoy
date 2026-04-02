<template>
  <div
    class="spacer-preview"
    :style="{ height: `${config?.height || 20}px` }"
  >
    <div class="spacer-indicator">
      <span>{{ config?.height || 20 }}px</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  config: Record<string, any>
}>()
</script>

<style scoped>
.spacer-preview {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spacer-indicator {
  position: absolute;
  padding: 2px 8px;
  background: #EEF2FF;
  border: 1px dashed #C7D2FE;
  border-radius: 4px;
  font-size: 11px;
  color: #6366F1;
  opacity: 0;
  transition: opacity 0.2s;
}

.spacer-preview:hover .spacer-indicator {
  opacity: 1;
}
</style>