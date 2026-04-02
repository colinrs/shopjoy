<template>
  <div class="reviews-section" v-loading="reviewsLoading">
    <!-- Review Stats Overview -->
    <div class="review-stats-overview">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic :title="$t('products.totalReviews')" :value="productStats?.total_reviews || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.averageRating')" :value="productStats?.average_rating || '0'" :suffix="$t('products.stars')" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.reviewsWithImages')" :value="productStats?.with_image_count || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('products.replyRate')" :value="productStats?.reply_rate || 0" :suffix="$t('products.replyRateUnit')">
            <template #suffix>
              <span style="font-size: 14px">{{ $t('products.replyRateUnit') }}</span>
            </template>
          </el-statistic>
        </el-col>
      </el-row>
    </div>

    <!-- Rating Distribution -->
    <div class="rating-distribution">
      <h3 class="section-title">{{ $t('products.ratingDistribution') }}</h3>
      <div class="rating-bars">
        <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="rating-bar-item">
          <span class="star-label">{{ star }}{{ $t('products.stars') }}</span>
          <div class="bar-container">
            <div
              class="bar-fill"
              :style="{ width: getRatingPercentage(star) + '%' }"
            ></div>
          </div>
          <span class="star-count">{{ getRatingCount(star) }}</span>
        </div>
      </div>
    </div>

    <!-- Rating Details -->
    <div class="rating-details">
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="rating-card">
            <div class="rating-card-title">{{ $t('products.qualityRating') }}</div>
            <div class="rating-card-value">{{ productStats?.quality_avg_rating || '0' }}</div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="rating-card">
            <div class="rating-card-title">{{ $t('products.valueRating') }}</div>
            <div class="rating-card-value">{{ productStats?.value_avg_rating || '0' }}</div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="rating-card">
            <div class="rating-card-title">{{ $t('products.replyCount') }}</div>
            <div class="rating-card-value">{{ productStats?.reply_count || 0 }}</div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getProductStats, type ProductStats } from '@/api/review'
import { t } from '@/plugins/i18n'
import type { ProductReviewsTabProps } from '../types'

const props = defineProps<ProductReviewsTabProps>()

const reviewsLoading = ref(false)
const productStats = ref<ProductStats | null>(null)

const loadProductStats = async () => {
  reviewsLoading.value = true
  try {
    const stats = await getProductStats(props.productId)
    productStats.value = stats
  } catch (error) {
    console.error('Failed to load product stats:', error)
    ElMessage.error(t('products.loadReviewStatsFailed'))
  } finally {
    reviewsLoading.value = false
  }
}

// Get rating count for specific star
const getRatingCount = (star: number): number => {
  if (!productStats.value?.rating_distribution) return 0
  const key = String(star) as '1' | '2' | '3' | '4' | '5'
  return productStats.value.rating_distribution[key] || 0
}

// Get rating percentage for specific star
const getRatingPercentage = (star: number): number => {
  if (!productStats.value?.total_reviews || productStats.value.total_reviews === 0) return 0
  const count = getRatingCount(star)
  return Math.round((count / productStats.value.total_reviews) * 100)
}

onMounted(() => {
  loadProductStats()
})

defineExpose({
  loadProductStats
})
</script>

<style scoped>
.reviews-section {
  padding: 0;
}

.review-stats-overview {
  padding: 20px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 24px;
}

.rating-distribution {
  margin-bottom: 24px;
}

.rating-bars {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-width: 500px;
}

.rating-bar-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.star-label {
  width: 40px;
  font-size: 14px;
  color: #6B7280;
}

.bar-container {
  flex: 1;
  height: 20px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.star-count {
  width: 50px;
  text-align: right;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.rating-details {
  margin-top: 24px;
}

.rating-card {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
}

.rating-card-title {
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 8px;
}

.rating-card-value {
  font-size: 28px;
  font-weight: 700;
  color: #6366F1;
}
</style>
