<template>
  <div class="user-review-list">
    <el-table v-loading="loading" :data="reviews" stripe>
      <el-table-column :label="$t('users.reviewsColumns.product')" min-width="180">
        <template #default="{ row }">
          <span class="product-name">{{ row.product_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用户" width="100" align="center">
        <template #default="{ row }">
          <span v-if="row.is_anonymous" class="anonymous">{{ $t('users.reviewsColumns.anonymous') }}</span>
          <span v-else>{{ row.user_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviewsColumns.rating')" width="110" align="center">
        <template #default="{ row }">
          <span class="rating">
            <el-icon v-for="i in 5" :key="i" :class="i <= Number(row.overall_rating) ? 'star-filled' : 'star-empty'">
              <Star />
            </el-icon>
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviewsColumns.content')" min-width="220">
        <template #default="{ row }">
          <el-tooltip :content="row.content" placement="top" :disabled="!row.content || row.content.length <= 80">
            <span class="content-text">{{ truncate(row.content, 80) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviewsColumns.status')" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviewsColumns.createdAt')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="100" align="center" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="goReviews(row)">
            {{ $t('common.view') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty v-if="!loading && reviews.length === 0" :description="$t('users.noReviews')" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Star } from '@element-plus/icons-vue'
import { getUserReviews, type UserReviewListItem } from '@/api/user'
import { ElMessage } from 'element-plus'
import { t } from '@/plugins/i18n'

const props = defineProps<{ userId?: string }>()
const router = useRouter()

const loading = ref(false)
const reviews = ref<UserReviewListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserReviews(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    reviews.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    console.error('Failed to load reviews:', err)
    ElMessage.error(t('users.loadReviewsFailed'))
    reviews.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const truncate = (s: string, n: number) => (s && s.length > n ? s.slice(0, n) + '…' : s || '-')

const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'approved' || status === 'visible') return 'success'
  if (status === 'pending') return 'warning'
  if (status === 'hidden') return 'info'
  return 'danger'
}

const goReviews = (row: UserReviewListItem) => {
  router.push({ path: '/reviews', query: { product_id: row.product_id } })
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-review-list { padding: 0; }
.product-name { font-weight: 500; color: #1E1B4B; }
.anonymous { color: #9CA3AF; font-style: italic; }
.rating { display: inline-flex; gap: 2px; }
.star-filled { color: #F59E0B; font-size: 14px; }
.star-empty { color: #E5E7EB; font-size: 14px; }
.content-text { font-size: 13px; color: #374151; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid #F3F4F6;
}
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>