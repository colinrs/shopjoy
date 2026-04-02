<template>
  <div class="divider-preview">
    <hr :style="{ borderStyle: config?.style || 'solid', borderColor: config?.color || '#E5E7EB' }">
  </div>
</template>

<script setup lang="ts">
defineProps<{
  config: Record<string, any>
}>()
</script>

<style scoped>
.divider-preview {
  padding: 16px 0;
}

hr {
  border: none;
  border-top-width: 1px;
  margin: 0;
}
</style>