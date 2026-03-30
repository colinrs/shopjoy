<template>
  <div class="operation-log">
    <div class="coming-soon">
      <el-empty :description="$t('users.comingSoon')" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  userId?: number
}>()
</script>

<style scoped>
.operation-log {
  padding: 0;
}

.coming-soon {
  padding: 40px 0;
}
</style>
