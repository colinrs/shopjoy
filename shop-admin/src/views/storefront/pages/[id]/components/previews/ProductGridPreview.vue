<template>
  <div class="product-grid-preview">
    <h4 v-if="config?.title" class="grid-title">{{ config.title }}</h4>
    <div class="grid-container" :style="{ gridTemplateColumns: `repeat(${config?.columns || 4}, 1fr)` }">
      <div v-for="i in (config?.count || 4)" :key="i" class="product-card">
        <div class="product-image">
          <el-icon size="24"><Goods /></el-icon>
        </div>
        <div class="product-info">
          <span class="product-name">商品 {{ i }}</span>
          <span class="product-price">¥99.00</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Goods } from '@element-plus/icons-vue'

defineProps<{
  config: Record<string, any>
}>()
</script>

<style scoped>
.product-grid-preview {
  padding: 8px;
}

.grid-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  text-align: center;
}

.grid-container {
  display: grid;
  gap: 16px;
}

.product-card {
  background: #F9FAFB;
  border-radius: 8px;
  overflow: hidden;
}

.product-image {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  color: #A5B4FC;
}

.product-info {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.product-name {
  font-size: 12px;
  color: #374151;
}

.product-price {
  font-size: 13px;
  font-weight: 600;
  color: #EF4444;
}
</style>