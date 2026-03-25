<template>
  <div class="pages-page">
    <page-header title="页面管理" subtitle="管理店铺的所有页面">
      <template #extra>
        <el-button type="primary" @click="createPage">
          <el-icon><Plus /></el-icon>
          新建页面
        </el-button>
      </template>
    </page-header>

    <!-- Pages Table -->
    <el-card shadow="hover" class="pages-card">
      <el-table :data="pages" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="页面名称" min-width="150">
          <template #default="{ row }">
            <div class="page-name">
              <el-icon class="page-icon" :class="row.page_type">
                <component :is="getPageIcon(row.page_type)" />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="page_type" label="页面类型" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getPageTypeTag(row.page_type)">
              {{ getPageTypeName(row.page_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="slug" label="URL路径" width="150">
          <template #default="{ row }">
            <code class="slug-code">/{{ row.slug }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="is_published" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'" effect="light" size="small">
              {{ row.is_published ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80">
          <template #default="{ row }">
            <el-tooltip content="点击查看版本历史">
              <el-button text type="primary" size="small" @click="showVersions(row)">
                v{{ row.version }}
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click="editPage(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button
                v-if="!row.is_published"
                size="small"
                type="success"
                @click="handlePublish(row)"
                :loading="row.publishing"
              >
                发布
              </el-button>
              <el-button
                v-else
                size="small"
                type="warning"
                @click="handleUnpublish(row)"
                :loading="row.unpublishing"
              >
                下线
              </el-button>
              <el-button size="small" @click="previewPage(row)">
                <el-icon><View /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Version History Drawer -->
    <el-drawer
      v-model="versionDrawerVisible"
      title="版本历史"
      size="450px"
      :with-header="true"
    >
      <div class="version-drawer-content" v-if="currentPage">
        <div class="version-header">
          <span class="page-title">{{ currentPage.name }}</span>
          <el-tag size="small">当前版本: v{{ currentPage.version }}</el-tag>
        </div>

        <div class="versions-list" v-loading="versionsLoading">
          <div
            v-for="ver in versions"
            :key="ver.id"
            class="version-item"
            :class="{ 'is-current': ver.version === currentPage.version }"
          >
            <div class="version-info">
              <span class="version-number">v{{ ver.version }}</span>
              <span class="version-time">{{ formatTime(ver.created_at) }}</span>
            </div>
            <div class="version-actions">
              <el-button
                size="small"
                text
                type="primary"
                @click="viewVersion(ver)"
              >
                查看
              </el-button>
              <el-button
                v-if="ver.version !== currentPage.version"
                size="small"
                text
                type="warning"
                @click="restoreVersion(ver)"
              >
                恢复
              </el-button>
              <el-tag v-else size="small" type="success">当前</el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- Version Detail Dialog -->
    <el-dialog
      v-model="versionDetailVisible"
      :title="`版本详情 - v${selectedVersion?.version}`"
      width="800px"
    >
      <div class="version-detail" v-if="versionDetail">
        <div class="detail-header">
          <span>创建时间: {{ formatTime(versionDetail.version.created_at) }}</span>
        </div>
        <div class="blocks-preview">
          <div
            v-for="block in versionDetail.blocks"
            :key="block.sort_order"
            class="block-item"
          >
            <div class="block-header">
              <el-tag size="small">{{ block.block_type }}</el-tag>
            </div>
            <pre class="block-config">{{ JSON.stringify(block.block_config, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, View, HomeFilled, Goods, Collection, Document } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import {
  listPages,
  publishPage,
  unpublishPage,
  listVersions,
  getVersion,
  restoreVersion as restoreVersionApi,
  type PageItem,
  type VersionItem,
  type VersionDetailResponse
} from '@/api/storefront'

const router = useRouter()
const loading = ref(false)
const versionsLoading = ref(false)
const pages = ref<(PageItem & { publishing?: boolean; unpublishing?: boolean })[]>([])
const versionDrawerVisible = ref(false)
const versionDetailVisible = ref(false)
const currentPage = ref<PageItem | null>(null)
const versions = ref<VersionItem[]>([])
const selectedVersion = ref<VersionItem | null>(null)
const versionDetail = ref<VersionDetailResponse | null>(null)

const fetchPages = async () => {
  loading.value = true
  try {
    const res = await listPages()
    pages.value = res.pages || []
  } catch (error) {
    ElMessage.error('获取页面列表失败')
  } finally {
    loading.value = false
  }
}

const getPageIcon = (type: string) => {
  const icons: Record<string, any> = {
    home: HomeFilled,
    product: Goods,
    collection: Collection,
    custom: Document
  }
  return icons[type] || Document
}

const getPageTypeName = (type: string) => {
  const names: Record<string, string> = {
    home: '首页',
    product: '商品页',
    collection: '集合页',
    custom: '自定义页'
  }
  return names[type] || '未知'
}

const getPageTypeTag = (type: string) => {
  const tags: Record<string, string> = {
    home: 'danger',
    product: 'warning',
    collection: 'success',
    custom: 'info'
  }
  return tags[type] || 'info'
}

const editPage = (page: PageItem) => {
  router.push(`/storefront/pages/${page.id}/edit`)
}

const previewPage = (page: PageItem) => {
  window.open(`/${page.slug}`, '_blank')
}

const createPage = () => {
  ElMessage.info('新建页面功能即将上线')
}

const handlePublish = async (page: PageItem & { publishing?: boolean }) => {
  try {
    await ElMessageBox.confirm(
      '确定要发布此页面吗？发布后访客将能看到最新内容。',
      '发布页面',
      { confirmButtonText: '确定发布', cancelButtonText: '取消', type: 'info' }
    )
    page.publishing = true
    await publishPage(page.id)
    ElMessage.success('页面已发布')
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '发布失败')
    }
  } finally {
    page.publishing = false
  }
}

const handleUnpublish = async (page: PageItem & { unpublishing?: boolean }) => {
  try {
    await ElMessageBox.confirm(
      '确定要下线此页面吗？下线后访客将无法访问。',
      '下线页面',
      { confirmButtonText: '确定下线', cancelButtonText: '取消', type: 'warning' }
    )
    page.unpublishing = true
    await unpublishPage(page.id)
    ElMessage.success('页面已下线')
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '下线失败')
    }
  } finally {
    page.unpublishing = false
  }
}

const showVersions = async (page: PageItem) => {
  currentPage.value = page
  versionDrawerVisible.value = true
  versionsLoading.value = true
  try {
    const res = await listVersions(page.id)
    versions.value = res.versions || []
  } catch (error) {
    ElMessage.error('获取版本历史失败')
  } finally {
    versionsLoading.value = false
  }
}

const viewVersion = async (ver: VersionItem) => {
  if (!currentPage.value) return
  try {
    const res = await getVersion(currentPage.value.id, ver.version)
    versionDetail.value = res
    selectedVersion.value = ver
    versionDetailVisible.value = true
  } catch (error) {
    ElMessage.error('获取版本详情失败')
  }
}

const restoreVersion = async (ver: VersionItem) => {
  if (!currentPage.value) return
  try {
    await ElMessageBox.confirm(
      `确定要恢复到版本 v${ver.version} 吗？当前内容将被替换。`,
      '恢复版本',
      { confirmButtonText: '确定恢复', cancelButtonText: '取消', type: 'warning' }
    )
    await restoreVersionApi(currentPage.value.id, { version: ver.version })
    ElMessage.success('版本已恢复')
    versionDrawerVisible.value = false
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '恢复失败')
    }
  }
}

const formatTime = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchPages()
})
</script>

<style scoped>
.pages-page {
  padding: 0;
}

.pages-card {
  border-radius: 12px;
}

.page-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-icon {
  font-size: 18px;
  color: var(--color-primary);
}

.page-icon.home { color: #EF4444; }
.page-icon.product { color: #F59E0B; }
.page-icon.collection { color: #10B981; }

.slug-code {
  background: #F3F4F6;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #6B7280;
}

.version-drawer-content {
  padding: 0 8px;
}

.version-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #E5E7EB;
}

.page-title {
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.versions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.version-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #F9FAFB;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.version-item:hover {
  background: #F3F4F6;
}

.version-item.is-current {
  background: #EEF2FF;
  border: 1px solid #C7D2FE;
}

.version-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.version-number {
  font-weight: 600;
  color: #1E1B4B;
}

.version-time {
  font-size: 12px;
  color: #9CA3AF;
}

.version-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.version-detail {
  max-height: 500px;
  overflow-y: auto;
}

.detail-header {
  margin-bottom: 16px;
  color: #6B7280;
  font-size: 13px;
}

.blocks-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.block-item {
  background: #F9FAFB;
  border-radius: 8px;
  padding: 12px;
}

.block-header {
  margin-bottom: 8px;
}

.block-config {
  margin: 0;
  padding: 12px;
  background: #1E1B4B;
  border-radius: 6px;
  color: #E5E7EB;
  font-size: 12px;
  font-family: 'Fira Code', monospace;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>