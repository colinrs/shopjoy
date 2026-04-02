<template>
  <div
    v-loading="loading"
    class="pricing-section"
  >
    <div class="section-header">
      <h3 class="section-title">
        {{ $t('products.marketPrice') }}
      </h3>
      <el-button
        type="primary"
        size="small"
        :loading="pricingSaveLoading"
        @click="handleSavePricing"
      >
        <el-icon><Check /></el-icon>
        {{ $t('products.savePricing') }}
      </el-button>
    </div>
    <el-table
      :data="pricingData"
      stripe
    >
      <el-table-column
        :label="$t('products.market')"
        min-width="150"
      >
        <template #default="{ row }">
          <div class="market-cell">
            <span class="market-code">{{ row.market_code }}</span>
            <span class="market-name">{{ row.market_name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('products.currency')"
        width="80"
        align="center"
      >
        <template #default="{ row }">
          <el-tag size="small">
            {{ row.currency }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('products.basePrice')"
        width="150"
        align="right"
      >
        <template #default="{ row }">
          <el-input-number
            v-model="row.price_value"
            :min="0"
            :precision="2"
            :controls="false"
            size="small"
            style="width: 100px"
          />
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('products.compareAtPrice')"
        width="150"
        align="right"
      >
        <template #default="{ row }">
          <el-input-number
            v-model="row.compare_at_price_value"
            :min="0"
            :precision="2"
            :controls="false"
            size="small"
            style="width: 100px"
          />
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('products.discount')"
        width="100"
        align="center"
      >
        <template #default="{ row }">
          <span
            v-if="row.compare_at_price_value > row.price_value && row.compare_at_price_value > 0"
            class="discount-badge"
          >
            -{{ Math.round((1 - row.price_value / row.compare_at_price_value) * 100) }}%
          </span>
          <span
            v-else
            class="text-muted"
          >-</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.status')"
        width="100"
        align="center"
      >
        <template #default="{ row }">
          <el-switch
            v-model="row.is_enabled"
            size="small"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-empty
      v-if="pricingData.length === 0 && !loading"
      :description="$t('products.notOnAnyMarket')"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { updateProductMarket, type ProductMarket } from '@/api/product'
import { t } from '@/plugins/i18n'
import type { ProductPricingTabProps, ProductPricingTabEmits, PricingRowData } from '../types'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<ProductPricingTabProps>()
const { handleError } = useErrorHandler()
const emit = defineEmits<ProductPricingTabEmits>()

const pricingData = ref<PricingRowData[]>([])
const pricingSaveLoading = ref(false)
const loading = ref(false)

watch(() => props.productMarkets, (newVal) => {
  preparePricingData(newVal)
}, { immediate: true })

const preparePricingData = (markets: ProductMarket[]) => {
  pricingData.value = markets.map(pm => ({
    ...pm,
    price_value: parseFloat(pm.price) || 0,
    compare_at_price_value: parseFloat(pm.compare_at_price) || 0
  }))
}

const handleSavePricing = async () => {
  pricingSaveLoading.value = true
  try {
    // Update each market price
    for (const item of pricingData.value) {
      const original = props.productMarkets.find(pm => pm.market_id === item.market_id)
      if (
        item.price_value !== parseFloat(item.price) ||
        item.compare_at_price_value !== parseFloat(item.compare_at_price || '0') ||
        item.is_enabled !== original?.is_enabled
      ) {
        await updateProductMarket(props.productId, item.market_id, {
          price: item.price_value.toString(),
          compare_at_price: item.compare_at_price_value.toString(),
          is_enabled: item.is_enabled
        })
      }
    }
    ElMessage.success(t('products.savePricingSuccess'))
    emit('refresh')
  } catch (error) {
    handleError(error, t('products.savePricingFailed'))
  } finally {
    pricingSaveLoading.value = false
  }
}

defineExpose({
  preparePricingData
})
</script>

<style scoped>
.pricing-section {
  padding: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.market-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.market-code {
  font-weight: 600;
  font-size: 14px;
}

.market-name {
  font-size: 12px;
  color: #6B7280;
}

.text-muted {
  color: #9CA3AF;
  font-size: 12px;
}

.discount-badge {
  background: #ECFDF5;
  color: #059669;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}
</style>
