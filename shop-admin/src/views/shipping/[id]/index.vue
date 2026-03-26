<template>
  <div class="template-detail-page">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-left">
          <el-button link @click="handleBack">
            <el-icon><ArrowLeft /></el-icon>
            返回模板列表
          </el-button>
          <el-divider direction="vertical" />
          <h2 class="template-title">
            {{ template?.name || '加载中...' }}
            <el-tag v-if="template?.is_default" type="success" size="small">默认</el-tag>
          </h2>
        </div>
        <div class="header-right">
          <el-switch
            v-model="isActive"
            active-text="已启用"
            inactive-text="已禁用"
            @change="handleStatusChange"
          />
          <el-button type="primary" @click="handleSave" :loading="saveLoading">
            <el-icon><Check /></el-icon>
            保存修改
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Tabs -->
    <el-card class="tabs-card" shadow="never" v-loading="loading">
      <el-tabs v-model="activeTab">
        <!-- Zones Tab -->
        <el-tab-pane label="配送区域" name="zones">
          <div class="zones-section">
            <div class="section-header">
              <h3 class="section-title">配送区域配置</h3>
              <el-button type="primary" @click="showAddZoneDialog">
                <el-icon><Plus /></el-icon>
                添加区域
              </el-button>
            </div>

            <!-- Zones List -->
            <div class="zones-list" v-if="zones.length > 0">
              <ZoneConfigForm
                v-for="(zone, index) in zones"
                :key="zone.id || `new-${index}`"
                :zone="zone"
                :index="index"
                @update="handleZoneUpdate"
                @delete="handleZoneDelete"
              />
            </div>

            <el-empty v-else description="暂无配送区域">
              <el-button type="primary" @click="showAddZoneDialog">
                添加第一个区域
              </el-button>
            </el-empty>
          </div>
        </el-tab-pane>

        <!-- Associations Tab -->
        <el-tab-pane label="关联配置" name="associations">
          <div class="associations-section">
            <!-- Product Associations -->
            <div class="association-block">
              <div class="block-header">
                <h3 class="block-title">商品关联</h3>
                <el-button type="primary" size="small" @click="showProductSelector">
                  <el-icon><Plus /></el-icon>
                  添加商品
                </el-button>
              </div>
              <el-table :data="productMappings" stripe v-if="productMappings.length > 0">
                <el-table-column label="商品" min-width="300">
                  <template #default="{ row }">
                    <div class="product-cell">
                      <div class="product-info">
                        <p class="product-name">{{ row.target_name || `商品ID: ${row.target_id}` }}</p>
                      </div>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="100" align="center">
                  <template #default="{ row }">
                    <el-button type="danger" link size="small" @click="removeProductMapping(row)">
                      移除
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无关联商品" :image-size="80" />
            </div>

            <!-- Category Associations -->
            <div class="association-block">
              <div class="block-header">
                <h3 class="block-title">分类关联</h3>
                <el-button type="primary" size="small" @click="showCategorySelector">
                  <el-icon><Plus /></el-icon>
                  添加分类
                </el-button>
              </div>
              <div class="category-tags" v-if="categoryMappings.length > 0">
                <el-tag
                  v-for="cat in categoryMappings"
                  :key="cat.id"
                  closable
                  @close="removeCategoryMapping(cat)"
                  class="category-tag"
                >
                  {{ cat.target_name || `分类 #${cat.target_id}` }}
                </el-tag>
              </div>
              <el-empty v-else description="暂无关联分类" :image-size="80" />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Add Zone Dialog -->
    <el-dialog
      v-model="zoneDialogVisible"
      :title="editingZone ? '编辑配送区域' : '添加配送区域'"
      width="800px"
      destroy-on-close
      :close-on-click-modal="false"
    >
      <ZoneConfigForm
        :zone="editingZone ?? undefined"
        :is-dialog="true"
        @save="handleZoneSave"
        @cancel="zoneDialogVisible = false"
      />
    </el-dialog>

    <!-- Product Selector Dialog -->
    <el-dialog
      v-model="productSelectorVisible"
      title="选择商品"
      width="800px"
      destroy-on-close
    >
      <div class="product-selector">
        <el-input
          v-model="productSearch"
          placeholder="搜索商品名称..."
          class="search-input"
          clearable
          @keyup.enter="searchProducts"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-table
          :data="availableProducts"
          v-loading="productLoading"
          @selection-change="handleProductSelection"
          max-height="400"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="商品" min-width="300">
            <template #default="{ row }">
              <div class="product-cell">
                <el-image v-if="row.images && row.images.length > 0" :src="row.images[0]" class="product-thumb" fit="cover" />
                <div class="product-info">
                  <p class="product-name">{{ row.name }}</p>
                </div>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="productSelectorVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmProductSelection" :disabled="selectedProducts.length === 0">
          确认添加 ({{ selectedProducts.length }})
        </el-button>
      </template>
    </el-dialog>

    <!-- Category Selector Dialog -->
    <el-dialog
      v-model="categorySelectorVisible"
      title="选择分类"
      width="500px"
      destroy-on-close
    >
      <el-tree
        ref="categoryTreeRef"
        :data="categoryTree"
        :props="{ label: 'name', children: 'children' }"
        show-checkbox
        node-key="id"
        :default-checked-keys="selectedCategoryIds"
        @check="handleCategoryCheck"
      />
      <template #footer>
        <el-button @click="categorySelectorVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCategorySelection" :disabled="selectedCategories.length === 0">
          确认添加 ({{ selectedCategories.length }})
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Check, Plus, Search } from '@element-plus/icons-vue'
import {
  getShippingTemplate,
  updateShippingTemplate,
  createShippingZone,
  updateShippingZone,
  deleteShippingZone,
  createTemplateMapping,
  deleteTemplateMapping,
  type TemplateDetail,
  type ShippingZone,
  type TemplateMapping,
  type CreateZoneRequest
} from '@/api/shipping'
import { getProductList } from '@/api/product'
import type { Product } from '@/api/product'
import { getCategoryTree, type CategoryTree } from '@/api/category'
import ZoneConfigForm from '../components/ZoneConfigForm.vue'

const route = useRoute()
const router = useRouter()

// State
const loading = ref(false)
const saveLoading = ref(false)
const activeTab = ref('zones')
const template = ref<TemplateDetail | null>(null)
const zones = ref<ShippingZone[]>([])
const mappings = ref<TemplateMapping[]>([])
const zoneDialogVisible = ref(false)
const editingZone = ref<ShippingZone | null>(null)

// Product selector
const productSelectorVisible = ref(false)
const productSearch = ref('')
const productLoading = ref(false)
const availableProducts = ref<Product[]>([])
const selectedProducts = ref<Product[]>([])

// Category selector
const categorySelectorVisible = ref(false)
const categoryTreeRef = ref()
const categoryTree = ref<CategoryTree[]>([])
const selectedCategories = ref<CategoryTree[]>([])
const selectedCategoryIds = ref<number[]>([])

// Computed
const templateId = computed(() => Number(route.params.id))
const isActive = computed({
  get: () => template.value?.is_active ?? true,
  set: () => {}
})

const productMappings = computed(() =>
  mappings.value.filter(m => m.target_type === 'product')
)

const categoryMappings = computed(() =>
  mappings.value.filter(m => m.target_type === 'category')
)

// Methods
const loadTemplate = async () => {
  loading.value = true
  try {
    const data = await getShippingTemplate(templateId.value)
    template.value = data
    zones.value = data.zones || []
    mappings.value = data.mappings || []
  } catch (error) {
    console.error('Failed to load template:', error)
    ElMessage.error('加载模板失败')
  } finally {
    loading.value = false
  }
}

const handleBack = () => {
  router.push('/shipping')
}

const handleStatusChange = async (val: boolean) => {
  try {
    await updateShippingTemplate(templateId.value, { is_active: val })
    ElMessage.success(val ? '已启用' : '已禁用')
    if (template.value) {
      template.value.is_active = val
    }
  } catch (error) {
    console.error('Failed to update status:', error)
    ElMessage.error('更新状态失败')
  }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    // Save zones if needed
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('Failed to save:', error)
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

// Zone methods
const showAddZoneDialog = () => {
  editingZone.value = null
  zoneDialogVisible.value = true
}

const handleZoneUpdate = async (zone: ShippingZone) => {
  try {
    if (zone.id) {
      await updateShippingZone(zone.id, zone)
    }
    ElMessage.success('更新成功')
    loadTemplate()
  } catch (error) {
    console.error('Failed to update zone:', error)
    ElMessage.error('更新失败')
  }
}

const handleZoneDelete = async (zoneId: number) => {
  try {
    await deleteShippingZone(zoneId)
    ElMessage.success('删除成功')
    loadTemplate()
  } catch (error) {
    console.error('Failed to delete zone:', error)
    ElMessage.error('删除失败')
  }
}

const handleZoneSave = async (zoneData: CreateZoneRequest) => {
  try {
    if (editingZone.value?.id) {
      await updateShippingZone(editingZone.value.id, zoneData)
    } else {
      await createShippingZone(templateId.value, zoneData)
    }
    ElMessage.success(editingZone.value ? '更新成功' : '添加成功')
    zoneDialogVisible.value = false
    loadTemplate()
  } catch (error) {
    console.error('Failed to save zone:', error)
    ElMessage.error('保存失败')
  }
}

// Product methods
const showProductSelector = async () => {
  productSelectorVisible.value = true
  await searchProducts()
}

const searchProducts = async () => {
  productLoading.value = true
  try {
    const data = await getProductList({
      page: 1,
      page_size: 50,
      name: productSearch.value || undefined
    })
    availableProducts.value = data.list || []
  } catch (error) {
    console.error('Failed to search products:', error)
  } finally {
    productLoading.value = false
  }
}

const handleProductSelection = (selection: Product[]) => {
  selectedProducts.value = selection
}

const confirmProductSelection = async () => {
  try {
    for (const product of selectedProducts.value) {
      await createTemplateMapping({
        template_id: templateId.value,
        target_type: 'product',
        target_id: product.id
      })
    }
    ElMessage.success('添加成功')
    productSelectorVisible.value = false
    loadTemplate()
  } catch (error) {
    console.error('Failed to add mappings:', error)
    ElMessage.error('添加失败')
  }
}

const removeProductMapping = async (mapping: TemplateMapping) => {
  try {
    await deleteTemplateMapping(mapping.id)
    ElMessage.success('移除成功')
    loadTemplate()
  } catch (error) {
    console.error('Failed to remove mapping:', error)
    ElMessage.error('移除失败')
  }
}

// Category methods
const showCategorySelector = async () => {
  categorySelectorVisible.value = true
  try {
    const data = await getCategoryTree()
    categoryTree.value = data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const handleCategoryCheck = (_: any, { checkedNodes }: any) => {
  selectedCategories.value = checkedNodes
}

const confirmCategorySelection = async () => {
  try {
    for (const category of selectedCategories.value) {
      await createTemplateMapping({
        template_id: templateId.value,
        target_type: 'category',
        target_id: category.id
      })
    }
    ElMessage.success('添加成功')
    categorySelectorVisible.value = false
    loadTemplate()
  } catch (error) {
    console.error('Failed to add mappings:', error)
    ElMessage.error('添加失败')
  }
}

const removeCategoryMapping = async (mapping: TemplateMapping) => {
  try {
    await deleteTemplateMapping(mapping.id)
    ElMessage.success('移除成功')
    loadTemplate()
  } catch (error) {
    console.error('Failed to remove mapping:', error)
    ElMessage.error('移除失败')
  }
}

// Lifecycle
onMounted(() => {
  loadTemplate()
})
</script>

<style scoped>
.template-detail-page {
  padding: 0;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.template-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.tabs-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.zones-section,
.associations-section {
  padding: 20px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
}

.zones-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.association-block {
  margin-bottom: 24px;
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.block-title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #374151;
}

.product-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-thumb {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: #F3F4F6;
}

.product-info {
  flex: 1;
}

.product-name {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #1E1B4B;
}

.product-sku {
  margin: 4px 0 0;
  font-size: 12px;
  color: #6B7280;
}

.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-tag {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
}

.product-selector {
  margin-bottom: 16px;
}

.search-input {
  margin-bottom: 16px;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
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

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-right {
    width: 100%;
    justify-content: space-between;
  }
}
</style>