<template>
  <div class="pages-page">
    <page-header :title="$t('storefront.pageManagement')" :subtitle="$t('storefront.manageAllPages')">
      <template #extra>
        <el-button type="primary" @click="createPage">
          <el-icon><Plus /></el-icon>
          {{ $t('storefront.createPage') }}
        </el-button>
      </template>
    </page-header>

    <!-- Pages Table -->
    <el-card shadow="hover" class="pages-card">
      <el-table :data="pages" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" :label="$t('storefront.pageName')" min-width="150">
          <template #default="{ row }">
            <div class="page-name">
              <el-icon class="page-icon" :class="row.page_type">
                <component :is="getPageIcon(row.page_type)" />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="page_type" :label="$t('storefront.pageType')" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getPageTypeTag(row.page_type)">
              {{ getPageTypeName(row.page_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="slug" :label="$t('storefront.urlPath')" width="150">
          <template #default="{ row }">
            <code class="slug-code">/{{ row.slug }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="is_published" :label="$t('storefront.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_published ? 'success' : 'info'" effect="light" size="small">
              {{ row.is_published ? $t('storefront.published') : $t('storefront.draft') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" :label="$t('storefront.version')" width="80">
          <template #default="{ row }">
            <el-tooltip :content="$t('storefront.viewVersionHistory')">
              <el-button text type="primary" size="small" @click="showVersions(row)">
                v{{ row.version }}
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column :label="$t('storefront.actions')" width="260" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click="editPage(row)">
                <el-icon><Edit /></el-icon>
                {{ $t('storefront.edit') }}
              </el-button>
              <el-button
                v-if="!row.is_published"
                size="small"
                type="success"
                @click="handlePublish(row)"
                :loading="row.publishing"
              >
                {{ $t('storefront.publish') }}
              </el-button>
              <el-button
                v-else
                size="small"
                type="warning"
                @click="handleUnpublish(row)"
                :loading="row.unpublishing"
              >
                {{ $t('storefront.unpublish') }}
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
      :title="$t('storefront.versionHistory')"
      size="450px"
      :with-header="true"
    >
      <div class="version-drawer-content" v-if="currentPage">
        <div class="version-header">
          <span class="page-title">{{ currentPage.name }}</span>
          <el-tag size="small">{{ $t('storefront.currentVersion') }}: v{{ currentPage.version }}</el-tag>
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
                {{ $t('storefront.view') }}
              </el-button>
              <el-button
                v-if="ver.version !== currentPage.version"
                size="small"
                text
                type="warning"
                @click="restoreVersion(ver)"
              >
                {{ $t('storefront.restore') }}
              </el-button>
              <el-tag v-else size="small" type="success">{{ $t('storefront.current') }}</el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- Version Detail Dialog -->
    <el-dialog
      v-model="versionDetailVisible"
      :title="$t('storefront.versionDetail') + ' v' + selectedVersion?.version"
      width="800px"
    >
      <div class="version-detail" v-if="versionDetail">
        <div class="detail-header">
          <span>{{ $t('storefront.createdTime') }}: {{ formatTime(versionDetail.version.created_at) }}</span>
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
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n()

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
    ElMessage.error(t('storefront.loadPagesFailed'))
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
    home: 'storefront.home',
    product: 'storefront.productPage',
    collection: 'storefront.collectionPage',
    custom: 'storefront.customPage'
  }
  const key = names[type]
  if (key === 'storefront.home') return '首页'
  if (key === 'storefront.productPage') return '商品页'
  if (key === 'storefront.collectionPage') return '集合页'
  if (key === 'storefront.customPage') return '自定义页'
  return 'storefront.unknown'
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
  ElMessage.info(t('storefront.createPageComingSoon'))
}

const handlePublish = async (page: PageItem & { publishing?: boolean }) => {
  try {
    await ElMessageBox.confirm(
      t('storefront.confirmPublishMessage'),
      t('storefront.confirmPublish'),
      { confirmButtonText: 'Publish', cancelButtonText: 'Cancel', type: 'info' }
    )
    page.publishing = true
    await publishPage(page.id)
    ElMessage.success(t('storefront.pagePublished'))
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || t('storefront.publishFailed'))
    }
  } finally {
    page.publishing = false
  }
}

const handleUnpublish = async (page: PageItem & { unpublishing?: boolean }) => {
  try {
    await ElMessageBox.confirm(
      t('storefront.confirmUnpublishMessage'),
      t('storefront.confirmUnpublish'),
      { confirmButtonText: 'Unpublish', cancelButtonText: 'Cancel', type: 'warning' }
    )
    page.unpublishing = true
    await unpublishPage(page.id)
    ElMessage.success(t('storefront.pageUnpublished'))
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || t('storefront.unpublishFailed'))
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
    ElMessage.error(t('storefront.loadVersionHistoryFailed'))
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
    ElMessage.error(t('storefront.loadVersionDetailFailed'))
  }
}

const restoreVersion = async (ver: VersionItem) => {
  if (!currentPage.value) return
  try {
    await ElMessageBox.confirm(
      t('storefront.confirmRestoreMessage', { version: ver.version }),
      t('storefront.confirmRestore'),
      { confirmButtonText: 'Restore', cancelButtonText: 'Cancel', type: 'warning' }
    )
    await restoreVersionApi(currentPage.value.id, { version: ver.version })
    ElMessage.success(t('storefront.versionRestored'))
    versionDrawerVisible.value = false
    await fetchPages()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || t('storefront.restoreFailed'))
    }
  }
}

const formatTime = (timestampStr: string) => {
  return new Date(timestampStr).toLocaleString('zh-CN', {
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
