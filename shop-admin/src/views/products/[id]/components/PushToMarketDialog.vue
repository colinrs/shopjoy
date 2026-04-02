<template>
  <el-dialog
    :model-value="visible"
    :title="$t('products.pushToMarket')"
    width="500px"
    destroy-on-close
    @update:model-value="emit('update:visible', $event)"
  >
    <el-form
      ref="pushToMarketFormRef"
      :model="pushToMarketForm"
      label-width="120px"
    >
      <el-form-item
        :label="$t('products.market')"
        prop="markets"
        required
      >
        <el-checkbox-group v-model="pushToMarketForm.selectedMarkets">
          <el-checkbox
            v-for="market in availableMarketsForPush"
            :key="market.id"
            :value="market.id"
            :label="market.id"
          >
            {{ market.code }} - {{ market.name }}
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item
        :label="$t('products.priceUSD')"
        prop="price"
        required
      >
        <el-input-number
          v-model="pushToMarketForm.price"
          :min="0"
          :precision="2"
          style="width: 100%"
        />
        <div class="price-note">
          {{ $t('products.priceNote') }}
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleConfirm"
      >
        {{ $t('products.pushToMarket') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { pushToMarket } from '@/api/product'
import { t } from '@/plugins/i18n'
import type { PushToMarketDialogProps, PushToMarketDialogEmits, PushToMarketFormData } from '../types'

const props = defineProps<PushToMarketDialogProps>()
const emit = defineEmits<PushToMarketDialogEmits>()

const pushToMarketFormRef = ref()
const pushToMarketForm = reactive<PushToMarketFormData>({
  selectedMarkets: [],
  price: 0
})

watch(() => props.visible, (newVal) => {
  if (newVal) {
    pushToMarketForm.selectedMarkets = []
    pushToMarketForm.price = parseFloat(props.productPrice) || 0
  }
})

const availableMarketsForPush = computed(() => {
  const existingMarketIds = props.productMarkets.map(pm => pm.market_id)
  return props.markets.filter(m => m.is_active && !existingMarketIds.includes(m.id))
})

const handleClose = () => {
  emit('update:visible', false)
}

const handleConfirm = async () => {
  if (pushToMarketForm.selectedMarkets.length === 0) {
    ElMessage.warning(t('products.selectAtLeastOneMarket'))
    return
  }
  if (pushToMarketForm.price <= 0) {
    ElMessage.warning(t('products.enterValidPrice'))
    return
  }

  try {
    const prices = pushToMarketForm.selectedMarkets.map(() =>
      pushToMarketForm.price.toFixed(2)
    )

    const result = await pushToMarket(props.productId, {
      market_ids: pushToMarketForm.selectedMarkets,
      prices
    })

    ElMessage.success(
      t('products.pushToMarketSuccess', { success: result.success?.length || 0, failed: result.failed?.length || 0 })
    )
    emit('success')
    handleClose()
  } catch (error) {
    console.error('Failed to push to market:', error)
    ElMessage.error(t('products.pushToMarketFailed'))
  }
}
</script>

<style scoped>
.price-note {
  font-size: 12px;
  color: #F59E0B;
  margin-top: 8px;
  line-height: 1.4;
}
</style>
