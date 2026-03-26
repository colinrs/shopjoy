<template>
  <el-card class="template-card" shadow="hover" @click="handleClick">
    <!-- Header -->
    <div class="card-header">
      <div class="header-left">
        <el-tag v-if="template.is_default" type="success" effect="dark" size="small" class="default-badge">
          默认
        </el-tag>
        <h3 class="template-name">{{ template.name }}</h3>
      </div>
      <el-tag :type="template.is_active ? 'success' : 'info'" size="small">
        {{ template.is_active ? '已启用' : '已禁用' }}
      </el-tag>
    </div>

    <!-- Stats -->
    <div class="card-stats">
      <div class="stat-item">
        <el-icon><Location /></el-icon>
        <span class="stat-value">{{ template.zone_count }}</span>
        <span class="stat-label">配送区域</span>
      </div>
      <div class="stat-item">
        <el-icon><Goods /></el-icon>
        <span class="stat-value">{{ template.product_count }}</span>
        <span class="stat-label">商品</span>
      </div>
      <div class="stat-item">
        <el-icon><Menu /></el-icon>
        <span class="stat-value">{{ template.category_count }}</span>
        <span class="stat-label">分类</span>
      </div>
    </div>

    <!-- Actions -->
    <div class="card-actions" @click.stop>
      <el-button type="primary" link size="small" @click="$emit('edit', template.id)">
        <el-icon><Edit /></el-icon>
        编辑
      </el-button>
      <el-button
        v-if="!template.is_default"
        type="primary"
        link
        size="small"
        @click="$emit('set-default', template)"
      >
        <el-icon><Star /></el-icon>
        设为默认
      </el-button>
      <el-popconfirm
        v-if="!template.is_default"
        title="确认删除该模板？"
        @confirm="$emit('delete', template)"
      >
        <template #reference>
          <el-button type="danger" link size="small">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </template>
      </el-popconfirm>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Location, Goods, Menu, Edit, Delete, Star } from '@element-plus/icons-vue'
import type { ShippingTemplate } from '@/api/shipping'

const props = defineProps<{
  template: ShippingTemplate
}>()

const emit = defineEmits<{
  edit: [id: number]
  delete: [template: ShippingTemplate]
  'set-default': [template: ShippingTemplate]
}>()

const handleClick = () => {
  emit('edit', props.template.id)
}
</script>

<style scoped>
.template-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
}

.template-card:hover {
  border-color: rgba(99, 102, 241, 0.2);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.default-badge {
  width: fit-content;
}

.template-name {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1E1B4B;
  line-height: 1.4;
}

.card-stats {
  display: flex;
  justify-content: space-between;
  padding: 16px 0;
  border-top: 1px solid #F3F4F6;
  border-bottom: 1px solid #F3F4F6;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.stat-item .el-icon {
  font-size: 18px;
  color: #6366F1;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
}

.stat-label {
  font-size: 12px;
  color: #6B7280;
}

.card-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.card-actions .el-button {
  padding: 4px 8px;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}
</style>