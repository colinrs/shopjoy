<template>
  <div class="block-config-form">
    <!-- Banner Config -->
    <template v-if="blockType === 'banner'">
      <el-form-item :label="$t('storefront.autoplay')">
        <el-switch v-model="localConfig.autoplay" />
      </el-form-item>
      <el-form-item :label="$t('storefront.playInterval')">
        <el-input-number v-model="localConfig.interval" :min="1000" :max="10000" :step="500" />
      </el-form-item>
      <el-form-item :label="$t('storefront.imageList')">
        <div class="image-list">
          <div v-for="(_img, idx) in localConfig.images" :key="idx" class="image-item">
            <el-input v-model="localConfig.images[idx]" :placeholder="$t('storefront.imageUrl')" />
            <el-button text type="danger" @click="removeImage(idx)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <el-button type="primary" text @click="addImage">
            <el-icon><Plus /></el-icon> {{ $t('storefront.addImage') }}
          </el-button>
        </div>
      </el-form-item>
    </template>

    <!-- Product Grid Config -->
    <template v-else-if="blockType === 'product_grid' || blockType === 'featured_products'">
      <el-form-item :label="$t('storefront.title')">
        <el-input v-model="localConfig.title" :placeholder="$t('storefront.blockTitle')" />
      </el-form-item>
      <el-form-item :label="$t('storefront.columns')">
        <el-select v-model="localConfig.columns">
          <el-option :value="2" :label="$t('storefront.column2')" />
          <el-option :value="3" :label="$t('storefront.column3')" />
          <el-option :value="4" :label="$t('storefront.column4')" />
          <el-option :value="5" :label="$t('storefront.column5')" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="blockType === 'featured_products'" :label="$t('storefront.displayCount')">
        <el-input-number v-model="localConfig.count" :min="4" :max="20" />
      </el-form-item>
      <el-form-item v-else :label="$t('storefront.productId')">
        <el-select
          v-model="localConfig.product_ids"
          multiple
          filterable
          remote
          :remote-method="loadProducts"
          :loading="loadingProducts"
          :placeholder="$t('storefront.searchAndSelectProduct')"
          style="width: 100%"
        >
          <el-option
            v-for="product in productOptions"
            :key="product.id"
            :label="product.name"
            :value="product.id"
          >
            <span>{{ product.name }}</span>
            <span style="color: #909399; font-size: 12px; margin-left: 8px;">
              {{ product.sku_code }}
            </span>
          </el-option>
        </el-select>
        <div style="font-size: 12px; color: #909399; margin-top: 4px;">
          {{ $t('storefront.selectedProducts', { count: localConfig.product_ids?.length || 0 }) }}
        </div>
      </el-form-item>
    </template>

    <!-- Rich Text Config -->
    <template v-else-if="blockType === 'rich_text'">
      <el-form-item :label="$t('storefront.content')">
        <el-input
          v-model="localConfig.content"
          type="textarea"
          :rows="10"
          :placeholder="$t('storefront.supportHtmlContent')"
        />
      </el-form-item>
    </template>

    <!-- Divider Config -->
    <template v-else-if="blockType === 'divider'">
      <el-form-item :label="$t('storefront.style')">
        <el-select v-model="localConfig.style">
          <el-option value="solid" :label="$t('storefront.solidLine')" />
          <el-option value="dashed" :label="$t('storefront.dashedLine')" />
          <el-option value="dotted" :label="$t('storefront.dottedLine')" />
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('storefront.color')">
        <el-color-picker v-model="localConfig.color" />
      </el-form-item>
    </template>

    <!-- Spacer Config -->
    <template v-else-if="blockType === 'spacer'">
      <el-form-item :label="$t('storefront.heightPx')">
        <el-slider v-model="localConfig.height" :min="10" :max="200" show-input />
      </el-form-item>
    </template>

    <!-- Video Config -->
    <template v-else-if="blockType === 'video'">
      <el-form-item :label="$t('storefront.videoUrl')">
        <el-input v-model="localConfig.url" :placeholder="$t('storefront.videoUrl')" />
      </el-form-item>
      <el-form-item :label="$t('storefront.autoplay')">
        <el-switch v-model="localConfig.autoplay" />
      </el-form-item>
    </template>

    <!-- Categories Config -->
    <template v-else-if="blockType === 'categories'">
      <el-form-item :label="$t('storefront.showAll')">
        <el-switch v-model="localConfig.show_all" />
      </el-form-item>
      <el-form-item :label="$t('storefront.columns')">
        <el-select v-model="localConfig.columns">
          <el-option :value="2" :label="$t('storefront.column2')" />
          <el-option :value="3" :label="$t('storefront.column3')" />
          <el-option :value="4" :label="$t('storefront.column4')" />
        </el-select>
      </el-form-item>
    </template>

    <!-- Custom HTML Config -->
    <template v-else-if="blockType === 'custom_html'">
      <el-form-item :label="$t('storefront.htmlCode')">
        <el-input
          v-model="localConfig.html"
          type="textarea"
          :rows="15"
          :placeholder="$t('storefront.customHtmlCode')"
        />
      </el-form-item>
    </template>

    <!-- Default Config -->
    <template v-else>
      <el-alert type="info" :closable="false">
        {{ $t('storefront.noConfigurableOptions') }}
      </el-alert>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive, onMounted } from 'vue'
import { Delete, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getProductList, type Product } from '@/api/product'

const props = defineProps<{
  blockType: string
  config: Record<string, any>
}>()

const emit = defineEmits<{
  update: [config: Record<string, any>]
}>()

const localConfig = reactive<Record<string, any>>({ ...props.config })

// Product options for product_grid selector
const productOptions = ref<Array<{ id: number; name: string; sku_code: string }>>([])
const loadingProducts = ref(false)

// Load product options
const loadProducts = async () => {
  loadingProducts.value = true
  try {
    const response = await getProductList({ page: 1, page_size: 100 })
    productOptions.value = response.list.map((p: Product) => ({
      id: p.id,
      name: p.name,
      sku_code: p.sku || ''
    }))
  } catch (error) {
    console.error('Failed to load products:', error)
    ElMessage.error('storefront.loadProductsFailed')
  } finally {
    loadingProducts.value = false
  }
}

onMounted(() => {
  if (props.blockType === 'product_grid') {
    loadProducts()
  }
})

watch(() => props.config, (newConfig) => {
  Object.assign(localConfig, newConfig)
}, { deep: true })

watch(localConfig, (newConfig) => {
  emit('update', { ...newConfig })
}, { deep: true })

const addImage = () => {
  if (!localConfig.images) {
    localConfig.images = []
  }
  localConfig.images.push('')
}

const removeImage = (index: number) => {
  localConfig.images.splice(index, 1)
}
</script>

<style scoped>
.block-config-form {
  padding: 0;
}

.image-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.image-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.image-item .el-input {
  flex: 1;
}
</style>
