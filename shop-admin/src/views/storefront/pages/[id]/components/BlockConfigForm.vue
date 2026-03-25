<template>
  <div class="block-config-form">
    <!-- Banner Config -->
    <template v-if="blockType === 'banner'">
      <el-form-item label="自动播放">
        <el-switch v-model="localConfig.autoplay" />
      </el-form-item>
      <el-form-item label="播放间隔(ms)">
        <el-input-number v-model="localConfig.interval" :min="1000" :max="10000" :step="500" />
      </el-form-item>
      <el-form-item label="图片列表">
        <div class="image-list">
          <div v-for="(img, idx) in localConfig.images" :key="idx" class="image-item">
            <el-input v-model="localConfig.images[idx]" placeholder="图片URL" />
            <el-button text type="danger" @click="removeImage(idx)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <el-button type="primary" text @click="addImage">
            <el-icon><Plus /></el-icon> 添加图片
          </el-button>
        </div>
      </el-form-item>
    </template>

    <!-- Product Grid Config -->
    <template v-else-if="blockType === 'product_grid' || blockType === 'featured_products'">
      <el-form-item label="标题">
        <el-input v-model="localConfig.title" placeholder="区块标题" />
      </el-form-item>
      <el-form-item label="列数">
        <el-select v-model="localConfig.columns">
          <el-option :value="2" label="2列" />
          <el-option :value="3" label="3列" />
          <el-option :value="4" label="4列" />
          <el-option :value="5" label="5列" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="blockType === 'featured_products'" label="显示数量">
        <el-input-number v-model="localConfig.count" :min="4" :max="20" />
      </el-form-item>
      <el-form-item v-else label="商品ID">
        <el-select
          v-model="localConfig.product_ids"
          multiple
          filterable
          placeholder="选择商品"
        >
          <!-- TODO: Load product options -->
        </el-select>
      </el-form-item>
    </template>

    <!-- Rich Text Config -->
    <template v-else-if="blockType === 'rich_text'">
      <el-form-item label="内容">
        <el-input
          v-model="localConfig.content"
          type="textarea"
          :rows="10"
          placeholder="支持HTML内容"
        />
      </el-form-item>
    </template>

    <!-- Divider Config -->
    <template v-else-if="blockType === 'divider'">
      <el-form-item label="样式">
        <el-select v-model="localConfig.style">
          <el-option value="solid" label="实线" />
          <el-option value="dashed" label="虚线" />
          <el-option value="dotted" label="点线" />
        </el-select>
      </el-form-item>
      <el-form-item label="颜色">
        <el-color-picker v-model="localConfig.color" />
      </el-form-item>
    </template>

    <!-- Spacer Config -->
    <template v-else-if="blockType === 'spacer'">
      <el-form-item label="高度(px)">
        <el-slider v-model="localConfig.height" :min="10" :max="200" show-input />
      </el-form-item>
    </template>

    <!-- Video Config -->
    <template v-else-if="blockType === 'video'">
      <el-form-item label="视频URL">
        <el-input v-model="localConfig.url" placeholder="视频链接" />
      </el-form-item>
      <el-form-item label="自动播放">
        <el-switch v-model="localConfig.autoplay" />
      </el-form-item>
    </template>

    <!-- Categories Config -->
    <template v-else-if="blockType === 'categories'">
      <el-form-item label="显示全部">
        <el-switch v-model="localConfig.show_all" />
      </el-form-item>
      <el-form-item label="列数">
        <el-select v-model="localConfig.columns">
          <el-option :value="2" label="2列" />
          <el-option :value="3" label="3列" />
          <el-option :value="4" label="4列" />
        </el-select>
      </el-form-item>
    </template>

    <!-- Custom HTML Config -->
    <template v-else-if="blockType === 'custom_html'">
      <el-form-item label="HTML代码">
        <el-input
          v-model="localConfig.html"
          type="textarea"
          :rows="15"
          placeholder="自定义HTML代码"
        />
      </el-form-item>
    </template>

    <!-- Default Config -->
    <template v-else>
      <el-alert type="info" :closable="false">
        此区块类型暂无可配置项
      </el-alert>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { Delete, Plus } from '@element-plus/icons-vue'

const props = defineProps<{
  blockType: string
  config: Record<string, any>
}>()

const emit = defineEmits<{
  update: [config: Record<string, any>]
}>()

const localConfig = reactive<Record<string, any>>({ ...props.config })

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