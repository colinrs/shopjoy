<template>
  <div class="reviews-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <StatsCard color="primary">
          <template #icon><ChatDotRound /></template>
          <template #value>{{ stats.total_reviews }}</template>
          <template #label>{{ $t('reviews.totalReviews') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="warning">
          <template #icon><Clock /></template>
          <template #value>{{ stats.pending_reviews }}</template>
          <template #label>{{ $t('reviews.pendingApproval') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="success">
          <template #icon><Star /></template>
          <template #value>{{ parseFloat(stats.average_rating).toFixed(1) }}</template>
          <template #label>{{ $t('reviews.averageRating') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="info">
          <template #icon><Picture /></template>
          <template #value>{{ stats.with_image_count }}</template>
          <template #label>{{ $t('reviews.withImages') }}</template>
        </StatsCard>
      </el-col>
    </el-row>

    <!-- Rating Distribution Card -->
    <el-card class="distribution-card" shadow="never">
      <div class="rating-distribution">
        <div class="distribution-item" v-for="star in 5" :key="star">
          <span class="star-label">{{ 6 - star }} {{ $t('reviews.stars') }}</span>
          <el-progress
            :percentage="getRatingPercentage(6 - star)"
            :stroke-width="10"
            :show-text="false"
            class="rating-progress"
          />
          <span class="count-label">{{ getRatingCount(6 - star) }}</span>
        </div>
      </div>
      <div class="reply-rate">
        <span class="rate-label">{{ $t('reviews.replyRate') }}:</span>
        <span class="rate-value">{{ (stats.reply_rate * 100).toFixed(1) }}%</span>
      </div>
    </el-card>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('reviews.searchPlaceholder')"
            class="search-input"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="filters.status" :placeholder="$t('reviews.status')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('reviews.allStatus')" value="" />
            <el-option :label="$t('reviews.pending')" value="pending" />
            <el-option :label="$t('reviews.approved')" value="approved" />
            <el-option :label="$t('reviews.hidden')" value="hidden" />
          </el-select>
          <el-select v-model="filters.rating_min" :placeholder="$t('reviews.minRating')" clearable class="filter-select" @change="handleSearch">
            <el-option :label="$t('reviews.any')" :value="0" />
            <el-option :label="'1+'" :value="1" />
            <el-option :label="'2+'" :value="2" />
            <el-option :label="'3+'" :value="3" />
            <el-option :label="'4+'" :value="4" />
            <el-option :label="'5'" :value="5" />
          </el-select>
          <el-checkbox v-model="filters.has_image" @change="handleSearch">{{ $t('reviews.hasImages') }}</el-checkbox>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            :range-separator="$t('reviews.startDate')"
            :start-placeholder="$t('reviews.startDate')"
            :end-placeholder="$t('reviews.endDate')"
            class="date-picker"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            @change="handleSearch"
          />
        </div>
        <div class="filter-right">
          <el-button v-if="selectedRows.length > 0" type="success" @click="handleBatchApprove">
            <el-icon><Check /></el-icon>{{ $t('reviews.batchApprove') }} ({{ selectedRows.length }})
          </el-button>
          <el-button v-if="selectedRows.length > 0" type="warning" @click="handleBatchHide">
            <el-icon><Hide /></el-icon>{{ $t('reviews.batchHide') }}
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>{{ $t('reviews.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Reviews Table -->
    <el-card class="table-card" shadow="never">
      <el-table
        ref="tableRef"
        :data="reviewList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column :label="$t('reviews.productInfo')" min-width="180">
          <template #default="{ row }">
            <div class="product-cell">
              <div class="product-info">
                <p class="product-name">{{ row.product_name }}</p>
                <p class="product-sku">{{ $t('reviews.sku') }}: {{ row.sku_code }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.reviewer')" min-width="120">
          <template #default="{ row }">
            <div class="reviewer-cell">
              <p class="reviewer-name">
                {{ row.is_anonymous ? $t('reviews.anonymous') : row.user_name }}
              </p>
              <div class="reviewer-badges">
                <el-tag v-if="row.is_verified" size="small" type="success" effect="plain" class="verified-tag">
                  <el-icon><CircleCheck /></el-icon>{{ $t('reviews.verified') }}
                </el-tag>
                <el-tag v-if="row.is_featured" size="small" type="warning" effect="plain">
                  <el-icon><Star /></el-icon>{{ $t('reviews.featured') }}
                </el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.ratingLabel')" width="160" align="center">
          <template #default="{ row }">
            <div class="rating-cell">
              <div class="rating-stars">
                <el-rate :model-value="parseFloat(row.overall_rating)" disabled :max="5" size="small" />
              </div>
              <div class="rating-detail">
                <span class="rating-item">
                  <span class="rating-label">Q:</span>
                  <span class="rating-value">{{ row.quality_rating }}</span>
                </span>
                <span class="rating-item">
                  <span class="rating-label">V:</span>
                  <span class="rating-value">{{ row.value_rating }}</span>
                </span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.contentLabel')" min-width="200">
          <template #default="{ row }">
            <div class="content-cell">
              <p class="content-text" :class="{ 'is-expanded': expandedRows.has(row.id) }">
                {{ row.content || $t('reviews.noContent') }}
              </p>
              <el-button
                v-if="row.content && row.content.length > 80"
                type="primary"
                link
                size="small"
                @click="toggleContent(row.id)"
              >
                {{ expandedRows.has(row.id) ? $t('reviews.collapse') : $t('reviews.expand') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.images')" width="100" align="center">
          <template #default="{ row }">
            <div v-if="row.images && row.images.length > 0" class="images-cell">
              <el-badge :value="row.images.length" :max="99" class="image-badge">
                <el-image
                  :src="row.images[0]"
                  :preview-src-list="row.images"
                  class="preview-image"
                  fit="cover"
                >
                  <template #error>
                    <div class="thumb-placeholder">
                      <el-icon><Picture /></el-icon>
                    </div>
                  </template>
                </el-image>
              </el-badge>
            </div>
            <span v-else class="no-images">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.status')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.helpful')" width="90" align="center">
          <template #default="{ row }">
            <div class="helpful-cell">
              <el-icon><Star /></el-icon>
              <span>{{ row.helpful_count }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.reply')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.has_reply" type="success" effect="plain" size="small">{{ $t('reviews.replied') }}</el-tag>
            <el-tag v-else type="info" effect="plain" size="small">{{ $t('reviews.noReply') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('reviews.createdAt')" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleViewDetail(row)">
              {{ $t('reviews.detail') }}
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              link
              size="small"
              @click="handleApprove(row)"
            >
              {{ $t('reviews.approve') }}
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              type="warning"
              link
              size="small"
              @click="handleHide(row)"
            >
              {{ $t('reviews.hide') }}
            </el-button>
            <el-button
              v-if="row.status === 'hidden'"
              type="success"
              link
              size="small"
              @click="handleShow(row)"
            >
              {{ $t('reviews.show') }}
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              :type="row.is_featured ? 'warning' : 'primary'"
              link
              size="small"
              @click="handleToggleFeatured(row)"
            >
              {{ row.is_featured ? $t('reviews.unfeature') : $t('reviews.feature') }}
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="handleReply(row)"
            >
              {{ $t('reviews.reply') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Review Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="$t('reviews.reviewDetails')"
      width="700px"
      destroy-on-close
      class="detail-dialog"
    >
      <div v-if="currentReview" class="detail-content">
        <!-- Order Info -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('reviews.orderInformation') }}</h4>
          <div class="order-info">
            <div class="info-item">
              <span class="info-label">{{ $t('reviews.orderId') }}:</span>
              <span class="info-value">{{ currentReview.order_id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('reviews.productLabel') }}:</span>
              <span class="info-value">{{ currentReview.product_name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('reviews.skuCode') }}:</span>
              <span class="info-value">{{ currentReview.sku_code }}</span>
            </div>
          </div>
        </div>

        <!-- Reviewer Info -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('reviews.reviewerInfo') }}</h4>
          <div class="reviewer-info">
            <div class="reviewer-meta">
              <p class="reviewer-name">
                {{ currentReview.is_anonymous ? $t('reviews.anonymousUser') : currentReview.user_name }}
                <el-tag v-if="currentReview.is_anonymous" size="small" type="info" effect="plain">
                  {{ $t('reviews.anonymousUser') }}
                </el-tag>
                <el-tag v-if="currentReview.is_verified" size="small" type="success" effect="plain">
                  <el-icon><CircleCheck /></el-icon>{{ $t('reviews.verifiedPurchase') }}
                </el-tag>
              </p>
              <p class="reviewer-date">{{ $t('reviews.reviewedOn') }} {{ formatDateTime(currentReview.created_at) }}</p>
            </div>
          </div>
        </div>

        <!-- Rating Display -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('reviews.ratings') }}</h4>
          <div class="rating-display">
            <div class="overall-rating">
              <span class="rating-score">{{ parseFloat(currentReview.overall_rating).toFixed(1) }}</span>
              <el-rate :model-value="parseFloat(currentReview.overall_rating)" disabled :max="5" size="large" />
            </div>
            <div class="rating-breakdown">
              <div class="rating-bar-item">
                <span class="bar-label">{{ $t('reviews.quality') }}</span>
                <el-progress
                  :percentage="currentReview.quality_rating * 20"
                  :stroke-width="10"
                  :show-text="false"
                  class="rating-progress"
                />
                <span class="bar-value">{{ currentReview.quality_rating }}/5</span>
              </div>
              <div class="rating-bar-item">
                <span class="bar-label">{{ $t('reviews.value') }}</span>
                <el-progress
                  :percentage="currentReview.value_rating * 20"
                  :stroke-width="10"
                  :show-text="false"
                  class="rating-progress"
                />
                <span class="bar-value">{{ currentReview.value_rating }}/5</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Review Content -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('reviews.reviewContent') }}</h4>
          <p class="review-content">{{ currentReview.content || $t('reviews.noContent') }}</p>
          <div v-if="currentReview.images && currentReview.images.length > 0" class="review-images">
            <el-image
              v-for="(img, idx) in currentReview.images"
              :key="idx"
              :src="img"
              :preview-src-list="currentReview.images"
              :initial-index="idx"
              class="review-image"
              fit="cover"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
        </div>

        <!-- Statistics -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('reviews.statistics') }}</h4>
          <div class="stats-info">
            <div class="stat-item">
              <span class="stat-label">{{ $t('reviews.helpfulCount') }}:</span>
              <span class="stat-value">{{ currentReview.helpful_count }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">{{ $t('reviews.status') }}:</span>
              <el-tag :type="getStatusType(currentReview.status)" effect="light" size="small">
                {{ getStatusText(currentReview.status) }}
              </el-tag>
            </div>
            <div class="stat-item" v-if="currentReview.is_featured">
              <span class="stat-label">{{ $t('reviews.featuredLabel') }}:</span>
              <el-tag type="warning" effect="light" size="small">{{ $t('reviews.yes') }}</el-tag>
            </div>
          </div>
        </div>

        <!-- Merchant Reply -->
        <div class="detail-section" v-if="currentReview.reply">
          <h4 class="section-title">{{ $t('reviews.merchantReply') }}</h4>
          <div class="merchant-reply">
            <div class="reply-header">
              <el-icon class="reply-icon"><ChatLineRound /></el-icon>
              <span class="reply-label">{{ $t('reviews.shopResponse') }}</span>
              <span class="reply-time">{{ formatDateTime(currentReview.reply.created_at) }}</span>
              <el-button type="primary" link size="small" @click="handleEditReply">{{ $t('reviews.editReply') }}</el-button>
              <el-button type="danger" link size="small" @click="handleDeleteReply">{{ $t('reviews.deleteReply') }}</el-button>
            </div>
            <p class="reply-content">{{ currentReview.reply.content }}</p>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button
            v-if="currentReview?.status === 'pending'"
            type="success"
            @click="handleApproveFromDetail"
          >
            {{ $t('reviews.approve') }}
          </el-button>
          <el-button
            v-if="currentReview?.status === 'approved'"
            type="warning"
            @click="handleHideFromDetail"
          >
            {{ $t('reviews.hide') }}
          </el-button>
          <el-button
            v-if="currentReview?.status === 'hidden'"
            type="success"
            @click="handleShowFromDetail"
          >
            {{ $t('reviews.show') }}
          </el-button>
          <el-button
            v-if="currentReview?.status === 'approved'"
            :type="currentReview?.is_featured ? 'warning' : 'primary'"
            @click="handleToggleFeaturedFromDetail"
          >
            {{ currentReview?.is_featured ? $t('reviews.unfeature') : $t('reviews.feature') }}
          </el-button>
          <el-button
            type="primary"
            @click="handleReplyFromDetail"
          >
            {{ currentReview?.reply ? $t('reviews.editReply') : $t('reviews.reply') }}
          </el-button>
          <el-button @click="detailDialogVisible = false">{{ $t('reviews.close') }}</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- Reply Dialog -->
    <el-dialog
      v-model="replyDialogVisible"
      :title="isEditReply ? $t('reviews.editReply') : $t('reviews.replyToReview')"
      width="500px"
      destroy-on-close
    >
      <div v-if="replyReview" class="reply-dialog-content">
        <!-- Original Review Preview -->
        <div class="original-review">
          <div class="review-header">
            <span class="reviewer">{{ replyReview.is_anonymous ? $t('reviews.anonymous') : replyReview.user_name }}</span>
            <el-rate :model-value="parseFloat(replyReview.overall_rating)" disabled :max="5" size="small" />
          </div>
          <p class="review-text">{{ replyReview.content || $t('reviews.noContent') }}</p>
        </div>

        <el-divider />

        <!-- Reply Form -->
        <el-form :model="replyForm" :rules="replyRules" ref="replyFormRef" label-position="top">
          <el-form-item :label="$t('reviews.yourReply')" prop="content">
            <el-input
              v-model="replyForm.content"
              type="textarea"
              :rows="5"
              :placeholder="$t('reviews.enterReply')"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <el-button @click="replyDialogVisible = false">{{ $t('reviews.cancel') }}</el-button>
        <el-button type="primary" @click="submitReply" :loading="replyLoading">
          {{ isEditReply ? $t('reviews.updateReply') : $t('reviews.submitReply') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Hide Reason Dialog -->
    <el-dialog
      v-model="hideDialogVisible"
      :title="$t('reviews.hideReview')"
      width="400px"
    >
      <el-form :model="hideForm" label-position="top">
        <el-form-item :label="$t('reviews.reasonOptional')">
          <el-input
            v-model="hideForm.reason"
            type="textarea"
            :rows="3"
            :placeholder="$t('reviews.enterReason')"
            maxlength="200"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="hideDialogVisible = false">{{ $t('reviews.cancel') }}</el-button>
        <el-button type="warning" @click="confirmHide" :loading="hideLoading">{{ $t('reviews.hide') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import StatsCard from '@/components/common/StatsCard.vue'
import {
  ChatDotRound,
  Clock,
  Star,
  Picture,
  Search,
  Refresh,
  Check,
  CircleCheck,
  ChatLineRound,
  Hide
} from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getReviewList,
  getReviewDetail,
  approveReview,
  hideReview,
  showReview,
  toggleFeatured,
  createReply,
  updateReply,
  deleteReply,
  batchApprove,
  batchHide,
  getReviewStats,
  type ReviewListItem,
  type ReviewDetail,
  type ReviewStats,
  type ListReviewsParams
} from '@/api/review'

const { t } = useI18n()

// State
const loading = ref(false)
const searchQuery = ref('')
const dateRange = ref<[string, string] | null>(null)
const tableRef = ref()
const selectedRows = ref<ReviewListItem[]>([])
const expandedRows = ref<Set<number>>(new Set())

// Statistics
const stats = ref<ReviewStats>({
  total_reviews: 0,
  pending_reviews: 0,
  approved_reviews: 0,
  hidden_reviews: 0,
  average_rating: '0',
  quality_avg_rating: '0',
  value_avg_rating: '0',
  five_star_count: 0,
  four_star_count: 0,
  three_star_count: 0,
  two_star_count: 0,
  one_star_count: 0,
  with_image_count: 0,
  reply_rate: 0,
  featured_count: 0
})

// Filters
const filters = reactive({
  status: '',
  rating_min: 0,
  has_image: false
})

// Pagination
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// Review list
const reviewList = ref<ReviewListItem[]>([])

// Detail Dialog
const detailDialogVisible = ref(false)
const currentReview = ref<ReviewDetail | null>(null)

// Reply Dialog
const replyDialogVisible = ref(false)
const replyReview = ref<ReviewListItem | ReviewDetail | null>(null)
const replyLoading = ref(false)
const replyFormRef = ref()
const replyForm = reactive({
  content: ''
})
const replyRules = {
  content: [
    { required: true, message: t('reviews.enterReply'), trigger: 'blur' },
    { min: 1, message: t('reviews.noContent'), trigger: 'blur' },
    { max: 500, message: t('reviews.submitReplyFailed'), trigger: 'blur' }
  ]
}
const isEditReply = ref(false)

// Hide Dialog
const hideDialogVisible = ref(false)
const hideLoading = ref(false)
const hideForm = reactive({
  reason: ''
})
const hideReviewId = ref<number | null>(null)
const isBatchHide = ref(false)

// Methods
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    'pending': 'warning',
    'approved': 'success',
    'hidden': 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'pending': t('reviews.pending'),
    'approved': t('reviews.approved'),
    'hidden': t('reviews.hidden')
  }
  return texts[status] || status
}

const getRatingPercentage = (star: number): number => {
  const total = stats.value.total_reviews
  if (total === 0) return 0
  return Math.round((getRatingCount(star) / total) * 100)
}

const getRatingCount = (star: number): number => {
  const counts: Record<number, number> = {
    5: stats.value.five_star_count,
    4: stats.value.four_star_count,
    3: stats.value.three_star_count,
    2: stats.value.two_star_count,
    1: stats.value.one_star_count
  }
  return counts[star] || 0
}

const toggleContent = (id: number) => {
  if (expandedRows.value.has(id)) {
    expandedRows.value.delete(id)
  } else {
    expandedRows.value.add(id)
  }
}

const handleSelectionChange = (rows: ReviewListItem[]) => {
  selectedRows.value = rows
}

const handleSearch = () => {
  pagination.page = 1
  loadReviews()
}

const handleRefresh = () => {
  loadReviews()
  loadStats()
  ElMessage.success(t('reviews.refresh') + ' ' + t('common.success'))
}

const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  loadReviews()
}

const handleCurrentChange = (val: number) => {
  pagination.page = val
  loadReviews()
}

const handleViewDetail = async (row: ReviewListItem) => {
  try {
    const detail = await getReviewDetail(row.id)
    currentReview.value = detail
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error(t('reviews.loadDetailFailed'))
  }
}

const handleApprove = async (row: ReviewListItem) => {
  try {
    await ElMessageBox.confirm(t('reviews.confirmApprove'), t('reviews.confirmApproval'), {
      confirmButtonText: t('reviews.approveButton'),
      cancelButtonText: t('reviews.cancelButton'),
      type: 'success'
    })
    await approveReview(row.id)
    ElMessage.success(t('reviews.approveSuccess'))
    loadReviews()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('reviews.approveFailed'))
    }
  }
}

const handleHide = (row: ReviewListItem) => {
  hideReviewId.value = row.id
  hideForm.reason = ''
  isBatchHide.value = false
  hideDialogVisible.value = true
}

const handleShow = async (row: ReviewListItem) => {
  try {
    await ElMessageBox.confirm(t('reviews.confirmShow'), t('reviews.confirmTitle'), {
      confirmButtonText: t('reviews.showButton'),
      cancelButtonText: t('reviews.cancelButton'),
      type: 'info'
    })
    await showReview(row.id)
    ElMessage.success(t('reviews.showSuccess'))
    loadReviews()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('reviews.showFailed'))
    }
  }
}

const confirmHide = async () => {
  if (hideReviewId.value === null) return

  hideLoading.value = true
  try {
    if (isBatchHide.value) {
      const ids = selectedRows.value.map(r => r.id)
      const result = await batchHide({ ids, reason: hideForm.reason || undefined })
      ElMessage.success(`${result.success_count} ${t('reviews.batchHideSuccess')}`)
      tableRef.value?.clearSelection()
    } else {
      await hideReview(hideReviewId.value, { reason: hideForm.reason || undefined })
      ElMessage.success(t('reviews.hideSuccess'))
    }
    hideDialogVisible.value = false
    loadReviews()
    loadStats()
  } catch (error) {
    ElMessage.error(t('reviews.hideFailed'))
  } finally {
    hideLoading.value = false
  }
}

const handleToggleFeatured = async (row: ReviewListItem) => {
  try {
    const newFeatured = !row.is_featured
    await toggleFeatured(row.id, { is_featured: newFeatured })
    ElMessage.success(newFeatured ? t('reviews.featureSuccess') : t('reviews.unfeatureSuccess'))
    loadReviews()
  } catch (error) {
    ElMessage.error(t('reviews.featureFailed'))
  }
}

const handleReply = (row: ReviewListItem | ReviewDetail) => {
  replyReview.value = row
  isEditReply.value = false
  replyForm.content = ''
  replyDialogVisible.value = true
}

const handleEditReply = () => {
  if (currentReview.value?.reply) {
    replyReview.value = currentReview.value
    replyForm.content = currentReview.value.reply.content
    isEditReply.value = true
    replyDialogVisible.value = true
  }
}

const handleDeleteReply = async () => {
  if (!currentReview.value) return

  try {
    await ElMessageBox.confirm(t('reviews.confirmDeleteReply'), t('reviews.confirmDeletion'), {
      confirmButtonText: t('reviews.deleteButton'),
      cancelButtonText: t('reviews.cancelButton'),
      type: 'warning'
    })
    await deleteReply(currentReview.value.id)
    ElMessage.success(t('reviews.deleteReplySuccess'))
    detailDialogVisible.value = false
    loadReviews()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('reviews.deleteReplyFailed'))
    }
  }
}

const submitReply = async () => {
  if (!replyFormRef.value) return
  const review = replyReview.value
  if (!review) return

  await replyFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    replyLoading.value = true
    try {
      if (isEditReply.value) {
        await updateReply(review.id, { content: replyForm.content })
        ElMessage.success(t('reviews.updateReplySuccess'))
      } else {
        await createReply(review.id, { content: replyForm.content })
        ElMessage.success(t('reviews.submitReplySuccess'))
      }
      replyDialogVisible.value = false
      loadReviews()
      loadStats()
      // Refresh detail if open
      if (currentReview.value && currentReview.value.id === review.id) {
        const detail = await getReviewDetail(review.id)
        currentReview.value = detail
      }
    } catch (error) {
      ElMessage.error(isEditReply.value ? t('reviews.updateReplyFailed') : t('reviews.submitReplyFailed'))
    } finally {
      replyLoading.value = false
    }
  })
}

const handleBatchApprove = async () => {
  try {
    await ElMessageBox.confirm(`${t('reviews.batchApproveConfirm')} (${selectedRows.value.length})`, t('reviews.batchApproveTitle'), {
      confirmButtonText: t('reviews.approveButton'),
      cancelButtonText: t('reviews.cancelButton'),
      type: 'success'
    })
    const ids = selectedRows.value.map(r => r.id)
    const result = await batchApprove({ ids })
    ElMessage.success(`${result.success_count} ${t('reviews.batchApproveSuccess')}`)
    tableRef.value?.clearSelection()
    loadReviews()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('reviews.batchApproveFailed'))
    }
  }
}

const handleBatchHide = () => {
  hideForm.reason = ''
  isBatchHide.value = true
  hideDialogVisible.value = true
}

// Detail dialog action handlers
const handleApproveFromDetail = async () => {
  if (currentReview.value) {
    await handleApprove({ id: currentReview.value.id } as ReviewListItem)
    detailDialogVisible.value = false
  }
}

const handleHideFromDetail = () => {
  if (currentReview.value) {
    hideReviewId.value = currentReview.value.id
    hideForm.reason = ''
    isBatchHide.value = false
    detailDialogVisible.value = false
    hideDialogVisible.value = true
  }
}

const handleShowFromDetail = async () => {
  if (currentReview.value) {
    await handleShow({ id: currentReview.value.id } as ReviewListItem)
    detailDialogVisible.value = false
  }
}

const handleToggleFeaturedFromDetail = async () => {
  if (currentReview.value) {
    await handleToggleFeatured({ id: currentReview.value.id, is_featured: currentReview.value.is_featured } as ReviewListItem)
    // Refresh detail
    const detail = await getReviewDetail(currentReview.value.id)
    currentReview.value = detail
  }
}

const handleReplyFromDetail = () => {
  if (currentReview.value) {
    if (currentReview.value.reply) {
      handleEditReply()
    } else {
      handleReply(currentReview.value)
    }
  }
}

const loadReviews = async () => {
  loading.value = true
  try {
    const params: ListReviewsParams = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchQuery.value) params.keyword = searchQuery.value
    if (filters.status) params.status = filters.status
    if (filters.rating_min > 0) {
      params.rating_min = filters.rating_min
      params.rating_max = 5
    }
    if (filters.has_image) params.has_image = true
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }

    const res = await getReviewList(params)
    reviewList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to load reviews:', error)
    ElMessage.error(t('reviews.loadFailed'))
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getReviewStats()
    stats.value = res
  } catch (error) {
    console.error('Failed to load stats:', error)
    ElMessage.error(t('reviews.loadStatsFailed'))
  }
}

onMounted(() => {
  loadReviews()
  loadStats()
})
</script>

<style scoped>
.reviews-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #6B7280;
}

.stat-icon.pending {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.stat-icon.rating {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.stat-icon.images {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Distribution Card */
.distribution-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.rating-distribution {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.distribution-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.star-label {
  width: 60px;
  font-size: 13px;
  color: #6B7280;
}

.rating-progress {
  flex: 1;
}

:deep(.rating-progress .el-progress-bar__inner) {
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
}

.count-label {
  width: 40px;
  font-size: 13px;
  color: #6366F1;
  font-weight: 600;
  text-align: right;
}

.reply-rate {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #F3F4F6;
}

.rate-label {
  font-size: 13px;
  color: #6B7280;
}

.rate-value {
  font-size: 14px;
  font-weight: 600;
  color: #6366F1;
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.search-input {
  width: 220px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 120px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.date-picker {
  width: 260px;
}

.date-picker :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Table */
.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Product Cell */
.product-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.product-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
  font-family: 'Fira Code', monospace;
}

/* Reviewer Cell */
.reviewer-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.reviewer-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.reviewer-badges {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.verified-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

/* Rating Cell */
.rating-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.rating-stars {
  display: flex;
}

.rating-detail {
  display: flex;
  gap: 12px;
  font-size: 12px;
}

.rating-item {
  display: flex;
  gap: 2px;
}

.rating-label {
  color: #9CA3AF;
}

.rating-value {
  color: #6366F1;
  font-weight: 600;
}

/* Content Cell */
.content-cell {
  max-width: 100%;
}

.content-text {
  font-size: 13px;
  color: #4B5563;
  margin: 0 0 4px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.content-text.is-expanded {
  -webkit-line-clamp: unset;
}

/* Images Cell */
.images-cell {
  display: flex;
  justify-content: center;
}

.image-badge {
  cursor: pointer;
}

.preview-image {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #E5E7EB;
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.no-images {
  color: #9CA3AF;
}

/* Helpful Cell */
.helpful-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #6B7280;
}

/* Time Text */
.time-text {
  font-size: 13px;
  color: #6B7280;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
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

/* Detail Dialog Content */
.detail-content {
  max-height: 60vh;
  overflow-y: auto;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #F3F4F6;
}

.order-info {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.info-item {
  display: flex;
  gap: 8px;
}

.info-label {
  color: #6B7280;
  font-size: 13px;
}

.info-value {
  color: #1E1B4B;
  font-size: 13px;
  font-weight: 500;
}

.reviewer-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reviewer-meta {
  flex: 1;
}

.reviewer-meta .reviewer-name {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.reviewer-date {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

.rating-display {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.overall-rating {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rating-score {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.rating-breakdown {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rating-bar-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bar-label {
  width: 60px;
  font-size: 13px;
  color: #6B7280;
}

.bar-value {
  width: 40px;
  font-size: 13px;
  color: #6366F1;
  font-weight: 600;
}

.review-content {
  font-size: 14px;
  color: #4B5563;
  line-height: 1.6;
  margin: 0;
}

.review-images {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.review-image {
  width: 80px;
  height: 80px;
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s;
}

.review-image:hover {
  transform: scale(1.05);
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.stats-info {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-item .stat-label {
  color: #6B7280;
  font-size: 13px;
}

.stat-item .stat-value {
  font-size: 14px;
  font-weight: 500;
  color: #1E1B4B;
}

.merchant-reply {
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  padding: 16px;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.reply-icon {
  color: #6366F1;
}

.reply-label {
  font-weight: 500;
  color: #6366F1;
}

.reply-time {
  margin-left: auto;
  font-size: 12px;
  color: #6B7280;
}

.reply-content {
  font-size: 14px;
  color: #4B5563;
  line-height: 1.6;
  margin: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: wrap;
}

/* Reply Dialog */
.original-review {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.reviewer {
  font-weight: 500;
  color: #1E1B4B;
}

.review-text {
  font-size: 14px;
  color: #4B5563;
  line-height: 1.5;
  margin: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select,
  .date-picker {
    width: 100%;
  }

  .filter-right {
    justify-content: flex-end;
  }

  .stat-card {
    border-radius: 14px;
    padding: 16px;
  }

  .stat-value {
    font-size: 20px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }
}
</style>
