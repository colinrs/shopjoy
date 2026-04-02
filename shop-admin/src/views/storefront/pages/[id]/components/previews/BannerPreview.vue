<template>
  <div class="banner-preview">
    <div v-if="!config?.images?.length" class="preview-placeholder">
      <el-icon size="32"><Picture /></el-icon>
      <span>{{ $t('storefront.bannerPreview') }} ({{ config?.images?.length || 0 }} {{ $t('points.images') }})</span>
    </div>
    <div v-else class="preview-carousel">
      <el-carousel
        :autoplay="config?.autoplay !== false"
        :interval="config?.interval || 5000"
        height="200px"
        indicator-position="outside"
      >
        <el-carousel-item v-for="(img, idx) in config.images.slice(0, 3)" :key="idx">
          <div class="carousel-image" :style="{ backgroundImage: `url(${img})` }">
            <span v-if="config.images.length > 3 && idx === 2" class="more-overlay">
              +{{ config.images.length - 3 }} {{ $t('points.images') }}
            </span>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Picture } from '@element-plus/icons-vue'

defineProps<{
  config: Record<string, any>
}>()
</script>

<style scoped>
.banner-preview {
  min-height: 120px;
}

.preview-placeholder {
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 8px;
  color: #6366F1;
}

.preview-carousel {
  border-radius: 8px;
  overflow: hidden;
}

.carousel-image {
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-color: #F3F4F6;
  display: flex;
  align-items: center;
  justify-content: center;
}

.more-overlay {
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
}
</style>