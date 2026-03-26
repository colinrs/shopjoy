<template>
  <div class="page-editor" v-loading="pageLoading">
    <!-- Editor Header -->
    <div class="editor-header">
      <div class="header-left">
        <el-button text @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <el-divider direction="vertical" />
        <span class="page-title">{{ pageData?.page?.name || '页面编辑' }}</span>
        <el-tag v-if="pageData?.page?.is_published" type="success" size="small">已发布</el-tag>
        <el-tag v-else type="info" size="small">草稿</el-tag>
      </div>
      <div class="header-center">
        <el-radio-group v-model="viewMode" size="small">
          <el-radio-button label="desktop">
            <el-icon><Monitor /></el-icon>
          </el-radio-button>
          <el-radio-button label="tablet">
            <el-icon><Grid /></el-icon>
          </el-radio-button>
          <el-radio-button label="mobile">
            <el-icon><Iphone /></el-icon>
          </el-radio-button>
        </el-radio-group>
      </div>
      <div class="header-right">
        <el-button @click="showVersionHistory">
          <el-icon><Clock /></el-icon>
          版本历史
        </el-button>
        <el-button @click="handleSaveDraft" :loading="saving">
          <el-icon><Document /></el-icon>
          保存草稿
        </el-button>
        <el-button type="primary" @click="handlePublish" :loading="publishing">
          <el-icon><Promotion /></el-icon>
          发布
        </el-button>
      </div>
    </div>

    <!-- Editor Main -->
    <div class="editor-main">
      <!-- Block Palette -->
      <div class="block-palette">
        <h4>区块组件</h4>
        <div class="palette-grid">
          <div
            v-for="block in BLOCK_TYPES"
            :key="block.type"
            class="palette-item"
            draggable="true"
            @dragstart="onDragStart($event, block.type)"
          >
            <el-icon class="block-icon"><component :is="getBlockIcon(block.icon)" /></el-icon>
            <span class="block-name">{{ block.name }}</span>
          </div>
        </div>
      </div>

      <!-- Canvas Area -->
      <div class="canvas-wrapper" :class="viewMode">
        <div
          class="canvas"
          @dragover.prevent
          @drop="onDrop"
          :class="{ 'is-empty': blocks.length === 0 }"
        >
          <div v-if="blocks.length === 0" class="empty-placeholder">
            <el-icon size="48"><Plus /></el-icon>
            <p>拖拽左侧区块到此处开始编辑</p>
          </div>

          <div v-else class="blocks-container">
            <div
              v-for="(block, index) in blocks"
              :key="block.id || `new-${index}`"
              class="block-wrapper"
              :class="{ 'is-selected': selectedBlockIndex === index }"
              draggable="true"
              @click="selectBlock(index)"
              @dragstart="onBlockDragStart($event, index)"
              @dragover.prevent
              @drop="onBlockReorder($event, index)"
            >
              <div class="block-toolbar">
                <span class="block-type">{{ getBlockName(block.block_type) }}</span>
                <div class="block-actions">
                  <el-button text size="small" @click.stop="moveBlock(index, -1)" :disabled="index === 0">
                    <el-icon><ArrowUp /></el-icon>
                  </el-button>
                  <el-button text size="small" @click.stop="moveBlock(index, 1)" :disabled="index === blocks.length - 1">
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <el-button text size="small" type="danger" @click.stop="deleteBlock(index)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
              <div class="block-content">
                <component
                  :is="getBlockPreviewComponent(block.block_type)"
                  :config="block.block_config"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Config Panel -->
      <div class="config-panel" v-if="selectedBlockIndex !== null">
        <div class="panel-header">
          <h4>{{ getBlockName(blocks[selectedBlockIndex]?.block_type) }} 配置</h4>
          <el-button text size="small" @click="selectedBlockIndex = null">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
        <div class="panel-content">
          <block-config-form
            :block-type="blocks[selectedBlockIndex]?.block_type"
            :config="blocks[selectedBlockIndex]?.block_config"
            @update="updateBlockConfig"
          />
        </div>
      </div>
    </div>

    <!-- Version History Drawer -->
    <el-drawer
      v-model="versionDrawerVisible"
      title="版本历史"
      size="450px"
    >
      <div class="version-drawer-content" v-loading="versionsLoading">
        <div
          v-for="ver in versions"
          :key="ver.id"
          class="version-item"
        >
          <div class="version-info">
            <span class="version-number">v{{ ver.version }}</span>
            <span class="version-time">{{ formatTime(ver.created_at) }}</span>
          </div>
          <div class="version-actions">
            <el-button size="small" text type="primary" @click="viewVersionDetail(ver)">
              查看
            </el-button>
            <el-button
              size="small"
              text
              type="warning"
              @click="handleRestoreVersion(ver)"
              :loading="ver.restoring"
            >
              恢复
            </el-button>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- Version Detail Dialog -->
    <el-dialog v-model="versionDetailVisible" :title="`版本详情 v${selectedVersion?.version}`" width="700px">
      <div class="version-detail" v-if="versionDetail">
        <div
          v-for="(block, index) in versionDetail.blocks"
          :key="index"
          class="version-block"
        >
          <el-tag size="small">{{ block.block_type }}</el-tag>
          <pre>{{ JSON.stringify(block.block_config, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Monitor, Iphone, Clock, Document, Promotion,
  Plus, ArrowUp, ArrowDown, Delete, Close, Picture, Goods, Document as DocIcon,
  VideoPlay, Menu, Minus, Rank
} from '@element-plus/icons-vue'
import {
  getPage,
  saveDraft,
  publishPage,
  listVersions,
  getVersion,
  restoreVersion,
  deleteDecoration,
  BLOCK_TYPES,
  type PageDetailResponse,
  type DecorationDTO,
  type VersionItem,
  type VersionDetailResponse as VersionDetailResp
} from '@/api/storefront'
import BlockConfigForm from './components/BlockConfigForm.vue'
import BannerPreview from './components/previews/BannerPreview.vue'
import ProductGridPreview from './components/previews/ProductGridPreview.vue'
import RichTextPreview from './components/previews/RichTextPreview.vue'
import DividerPreview from './components/previews/DividerPreview.vue'
import SpacerPreview from './components/previews/SpacerPreview.vue'
import DefaultPreview from './components/previews/DefaultPreview.vue'

const route = useRoute()
const router = useRouter()

const pageId = Number(route.params.id)
const pageLoading = ref(true)
const saving = ref(false)
const publishing = ref(false)
const pageData = ref<PageDetailResponse | null>(null)
const blocks = ref<DecorationDTO[]>([])
const selectedBlockIndex = ref<number | null>(null)
const viewMode = ref<'desktop' | 'tablet' | 'mobile'>('desktop')

// Version history
const versionDrawerVisible = ref(false)
const versionsLoading = ref(false)
const versions = ref<(VersionItem & { restoring?: boolean })[]>([])
const versionDetailVisible = ref(false)
const selectedVersion = ref<VersionItem | null>(null)
const versionDetail = ref<VersionDetailResp | null>(null)

// Block preview components
const previewComponents: Record<string, any> = {
  banner: markRaw(BannerPreview),
  product_grid: markRaw(ProductGridPreview),
  rich_text: markRaw(RichTextPreview),
  divider: markRaw(DividerPreview),
  spacer: markRaw(SpacerPreview),
  image_carousel: markRaw(BannerPreview),
  featured_products: markRaw(ProductGridPreview),
  categories: markRaw(DefaultPreview),
  video: markRaw(DefaultPreview),
  custom_html: markRaw(RichTextPreview)
}

const getBlockPreviewComponent = (type: string) => {
  return previewComponents[type] || markRaw(DefaultPreview)
}

const getBlockIcon = (icon: string) => {
  const icons: Record<string, any> = {
    Picture, Goods, Document: DocIcon, VideoPlay, Menu, Minus, Rank
  }
  return icons[icon] || Document
}

const getBlockName = (type: string) => {
  return BLOCK_TYPES.find(b => b.type === type)?.name || type
}

const fetchPage = async () => {
  pageLoading.value = true
  try {
    const res = await getPage(pageId)
    pageData.value = res
    blocks.value = res.decorations || []
  } catch (error) {
    ElMessage.error('获取页面数据失败')
  } finally {
    pageLoading.value = false
  }
}

const goBack = () => {
  router.push('/storefront/pages')
}

// Drag and Drop
let draggedBlockType = ''
let draggedBlockIndex: number | null = null

const onDragStart = (event: DragEvent, blockType: string) => {
  draggedBlockType = blockType
  draggedBlockIndex = null
  event.dataTransfer?.setData('text/plain', blockType)
}

const onBlockDragStart = (event: DragEvent, index: number) => {
  draggedBlockIndex = index
  draggedBlockType = ''
  event.dataTransfer?.setData('text/plain', `block-${index}`)
}

const onDrop = async (_event: DragEvent) => {
  if (draggedBlockType) {
    // Add new block
    const newBlock: DecorationDTO = {
      id: 0,
      block_type: draggedBlockType,
      block_config: getDefaultConfig(draggedBlockType),
      sort_order: blocks.value.length
    }
    blocks.value.push(newBlock)
    selectedBlockIndex.value = blocks.value.length - 1
  }
}

const onBlockReorder = (_event: DragEvent, targetIndex: number) => {
  if (draggedBlockIndex !== null && draggedBlockIndex !== targetIndex) {
    const block = blocks.value.splice(draggedBlockIndex, 1)[0]
    blocks.value.splice(targetIndex, 0, block)
    selectedBlockIndex.value = targetIndex
  }
}

const getDefaultConfig = (type: string): Record<string, any> => {
  const configs: Record<string, Record<string, any>> = {
    banner: { images: [], autoplay: true, interval: 5000 },
    product_grid: { title: '商品推荐', product_ids: [], columns: 4 },
    rich_text: { content: '<p>点击编辑内容...</p>' },
    image_carousel: { images: [] },
    featured_products: { title: '热门商品', count: 8 },
    categories: { show_all: true, columns: 4 },
    divider: { style: 'solid', color: '#E5E7EB' },
    video: { url: '', autoplay: false },
    spacer: { height: 20 },
    custom_html: { html: '' }
  }
  return configs[type] || {}
}

// Block operations
const selectBlock = (index: number) => {
  selectedBlockIndex.value = index
}

const moveBlock = async (index: number, direction: number) => {
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= blocks.value.length) return
  const block = blocks.value.splice(index, 1)[0]
  blocks.value.splice(newIndex, 0, block)
  selectedBlockIndex.value = newIndex
}

const deleteBlock = async (index: number) => {
  try {
    await ElMessageBox.confirm('确定要删除此区块吗？', '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const block = blocks.value[index]
    if (block.id) {
      await deleteDecoration(block.id)
    }
    blocks.value.splice(index, 1)
    if (selectedBlockIndex.value === index) {
      selectedBlockIndex.value = null
    } else if (selectedBlockIndex.value !== null && selectedBlockIndex.value > index) {
      selectedBlockIndex.value--
    }
    ElMessage.success('区块已删除')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const updateBlockConfig = (config: Record<string, any>) => {
  if (selectedBlockIndex.value !== null) {
    blocks.value[selectedBlockIndex.value].block_config = config
  }
}

// Save and Publish
const handleSaveDraft = async () => {
  saving.value = true
  try {
    await saveDraft(pageId, blocks.value)
    ElMessage.success('草稿已保存')
    await fetchPage()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handlePublish = async () => {
  try {
    await ElMessageBox.confirm(
      '发布后访客将能看到最新内容，确定发布吗？',
      '发布确认',
      { confirmButtonText: '发布', cancelButtonText: '取消', type: 'info' }
    )
    publishing.value = true
    // First save draft, then publish
    await saveDraft(pageId, blocks.value)
    await publishPage(pageId)
    ElMessage.success('页面已发布')
    await fetchPage()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '发布失败')
    }
  } finally {
    publishing.value = false
  }
}

// Version History
const showVersionHistory = async () => {
  versionDrawerVisible.value = true
  versionsLoading.value = true
  try {
    const res = await listVersions(pageId)
    versions.value = res.versions || []
  } catch (error) {
    ElMessage.error('获取版本历史失败')
  } finally {
    versionsLoading.value = false
  }
}

const viewVersionDetail = async (ver: VersionItem) => {
  try {
    const res = await getVersion(pageId, ver.version)
    versionDetail.value = res
    selectedVersion.value = ver
    versionDetailVisible.value = true
  } catch (error) {
    ElMessage.error('获取版本详情失败')
  }
}

const handleRestoreVersion = async (ver: VersionItem & { restoring?: boolean }) => {
  try {
    await ElMessageBox.confirm(
      `恢复到版本 v${ver.version} 将覆盖当前内容，确定恢复吗？`,
      '恢复版本',
      { confirmButtonText: '恢复', cancelButtonText: '取消', type: 'warning' }
    )
    ver.restoring = true
    await restoreVersion(pageId, { version: ver.version })
    ElMessage.success('版本已恢复')
    versionDrawerVisible.value = false
    await fetchPage()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '恢复失败')
    }
  } finally {
    ver.restoring = false
  }
}

const formatTime = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchPage()
})
</script>

<style scoped>
.page-editor {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #F5F3FF;
}

.editor-header {
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #E5E7EB;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
}

.header-center {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  gap: 12px;
}

.editor-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.block-palette {
  width: 200px;
  background: #fff;
  border-right: 1px solid #E5E7EB;
  padding: 16px;
  overflow-y: auto;
}

.block-palette h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
}

.palette-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.palette-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 8px;
  cursor: grab;
  transition: all 0.2s ease;
}

.palette-item:hover {
  background: #EEF2FF;
  transform: translateX(4px);
}

.palette-item:active {
  cursor: grabbing;
}

.block-icon {
  font-size: 18px;
  color: var(--color-primary);
}

.block-name {
  font-size: 13px;
  color: #374151;
}

.canvas-wrapper {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 24px;
  overflow-y: auto;
}

.canvas-wrapper.desktop .canvas {
  width: 100%;
  max-width: 1200px;
}

.canvas-wrapper.tablet .canvas {
  width: 768px;
}

.canvas-wrapper.mobile .canvas {
  width: 375px;
}

.canvas {
  background: #fff;
  border-radius: 12px;
  min-height: 600px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.canvas.is-empty {
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-placeholder {
  text-align: center;
  color: #9CA3AF;
}

.empty-placeholder p {
  margin: 12px 0 0 0;
}

.blocks-container {
  padding: 16px;
}

.block-wrapper {
  border: 2px solid transparent;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.2s ease;
}

.block-wrapper:hover {
  border-color: #E5E7EB;
}

.block-wrapper.is-selected {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.block-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #F9FAFB;
  border-radius: 6px 6px 0 0;
  opacity: 0;
  transition: opacity 0.2s;
}

.block-wrapper:hover .block-toolbar,
.block-wrapper.is-selected .block-toolbar {
  opacity: 1;
}

.block-type {
  font-size: 12px;
  font-weight: 500;
  color: #6B7280;
}

.block-actions {
  display: flex;
  gap: 4px;
}

.block-content {
  padding: 16px;
  min-height: 80px;
}

.config-panel {
  width: 360px;
  background: #fff;
  border-left: 1px solid #E5E7EB;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #E5E7EB;
}

.panel-header h4 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

/* Version Drawer */
.version-drawer-content {
  padding: 0 8px;
}

.version-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 8px;
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
  gap: 8px;
}

.version-detail {
  max-height: 500px;
  overflow-y: auto;
}

.version-block {
  background: #F9FAFB;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
}

.version-block pre {
  margin: 8px 0 0 0;
  padding: 12px;
  background: #1E1B4B;
  border-radius: 6px;
  color: #E5E7EB;
  font-size: 12px;
  overflow-x: auto;
}
</style>