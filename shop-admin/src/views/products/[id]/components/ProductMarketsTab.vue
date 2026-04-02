<template>
  <div class="markets-section">
    <div class="section-header">
      <h3 class="section-title">{{ $t('products.marketAvailability') }}</h3>
      <el-button type="primary" @click="handleShowPushToMarketDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('products.pushToMarket') }}
      </el-button>
    </div>

    <el-table :data="localProductMarkets" v-loading="loading" stripe>
      <el-table-column :label="$t('products.market')" min-width="150">
        <template #default="{ row }">
          <div class="market-cell">
            <span class="market-code">{{ row.market_code }}</span>
            <span class="market-name">{{ row.market_name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.status')" width="120" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.is_enabled"
            @change="(val: boolean) => handleMarketEnableChange(row, val)"
          />
        </template>
      </el-table-column>
      <el-table-column :label="$t('products.price')" width="180" align="right">
        <template #default="{ row }">
          <div class="price-cell">
            <el-input-number
              v-model="row.price"
              :min="0"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 100px"
              @change="() => handleMarketPriceChange(row)"
            />
            <span class="currency">{{ row.currency }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('products.compareAtPrice')" width="180" align="right">
        <template #default="{ row }">
          <div class="price-cell">
            <el-input-number
              v-model="row.compare_at_price"
              :min="0"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 100px"
              @change="() => handleMarketPriceChange(row)"
            />
            <span class="currency">{{ row.currency }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('products.stockAlert')" width="120" align="center">
        <template #default="{ row }">
          <el-input-number
            v-model="row.stock_alert_threshold"
            :min="0"
            :controls="false"
            size="small"
            style="width: 80px"
            @change="() => handleMarketPriceChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column :label="$t('products.publishedAt')" width="120" align="center">
        <template #default="{ row }">
          <span v-if="row.published_at">{{ formatDate(row.published_at) }}</span>
          <span v-else class="text-muted">{{ $t('products.notPublished') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="100" align="center">
        <template #default="{ row }">
          <el-button
            type="danger"
            link
            size="small"
            @click="handleRemoveFromMarket(row)"
          >
            {{ $t('common.remove') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="localProductMarkets.length === 0 && !loading" class="empty-markets">
      <el-empty :description="$t('products.notOnAnyMarket')">
        <el-button type="primary" @click="handleShowPushToMarketDialog">{{ $t('products.pushToMarket') }}</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { updateProductMarket, removeFromMarket, type ProductMarket } from '@/api/product'
import { t } from '@/plugins/i18n'
import type { ProductMarketsTabProps, ProductMarketsTabEmits } from '../types'

const props = defineProps<ProductMarketsTabProps>()
const emit = defineEmits<ProductMarketsTabEmits>()

const localProductMarkets = ref<ProductMarket[]>([...props.productMarkets])

watch(() => props.productMarkets, (newVal) => {
  localProductMarkets.value = [...newVal]
}, { deep: true })

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

const handleMarketEnableChange = async (row: ProductMarket, enabled: boolean) => {
  try {
    await updateProductMarket(props.productId, row.market_id, { is_enabled: enabled })
    ElMessage.success(enabled ? t('products.marketEnabled') : t('products.marketDisabled'))
    emit('refresh')
  } catch (error) {
    console.error('Failed to update market:', error)
    ElMessage.error(t('products.updateMarketStatusFailed'))
    row.is_enabled = !enabled // Revert
  }
}

const handleMarketPriceChange = async (row: ProductMarket) => {
  try {
    await updateProductMarket(props.productId, row.market_id, {
      price: row.price,
      compare_at_price: row.compare_at_price,
      stock_alert_threshold: row.stock_alert_threshold
    })
    ElMessage.success(t('products.marketPriceUpdated'))
  } catch (error) {
    console.error('Failed to update market price:', error)
    ElMessage.error(t('products.updateMarketPriceFailed'))
  }
}

const handleRemoveFromMarket = async (row: ProductMarket) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmRemoveFromMarket', { name: row.market_name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await removeFromMarket(props.productId, row.market_id)
    ElMessage.success(t('products.removedFromMarket'))
    emit('refresh')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove from market:', error)
      ElMessage.error(t('products.removeFromMarketFailed'))
    }
  }
}

const handleShowPushToMarketDialog = () => {
  emit('show-push-to-market')
}

defineExpose({
  getLocalMarkets: () => localProductMarkets.value
})
</script>

<style scoped>
.markets-section {
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

.price-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

.currency {
  font-size: 12px;
  color: #6B7280;
  min-width: 36px;
}

.text-muted {
  color: #9CA3AF;
  font-size: 12px;
}

.empty-markets {
  padding: 40px 0;
}
</style>
