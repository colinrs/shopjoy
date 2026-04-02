<template>
  <el-row
    :gutter="16"
    class="stats-row"
  >
    <el-col
      :xs="12"
      :sm="6"
    >
      <StatsCard color="warning">
        <template #icon>
          <PriceTag />
        </template>
        <template #value>
          {{ stats.pending_payment }}
        </template>
        <template #label>
          {{ $t('orders.pendingPayment') }}
        </template>
      </StatsCard>
    </el-col>
    <el-col
      :xs="12"
      :sm="6"
    >
      <StatsCard color="info">
        <template #icon>
          <Edit />
        </template>
        <template #value>
          {{ stats.partial_shipped }}
        </template>
        <template #label>
          {{ $t('orders.partialShipped') }}
        </template>
      </StatsCard>
    </el-col>
    <el-col
      :xs="12"
      :sm="6"
    >
      <StatsCard color="primary">
        <template #icon>
          <Van />
        </template>
        <template #value>
          {{ stats.shipped }}
        </template>
        <template #label>
          {{ $t('orders.shipped') }}
        </template>
      </StatsCard>
    </el-col>
    <el-col
      :xs="12"
      :sm="6"
    >
      <StatsCard color="success">
        <template #icon>
          <Check />
        </template>
        <template #value>
          {{ stats.delivered }}
        </template>
        <template #label>
          {{ $t('orders.delivered') }}
        </template>
      </StatsCard>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { PriceTag, Edit, Van, Check } from '@element-plus/icons-vue'
import StatsCard from '@/components/common/StatsCard.vue'

interface OrderStats {
  pending_payment: number
  partial_shipped: number
  shipped: number
  delivered: number
}

defineProps<{
  stats: OrderStats
}>()
</script>

<style scoped>
.stats-row {
  margin-bottom: 20px;
}
</style>
