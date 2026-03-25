<template>
  <div class="shipping-calculator-page">
    <!-- Page Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <div class="header-left">
          <h2>运费计算器</h2>
          <p class="header-desc">测试不同地址的运费计算结果</p>
        </div>
      </div>
    </el-card>

    <!-- Calculator -->
    <el-card class="calculator-card" shadow="never">
      <el-row :gutter="24">
        <!-- Address Selection -->
        <el-col :span="10" :xs="24">
          <div class="address-section">
            <h3 class="section-title">
              <el-icon><Location /></el-icon>
              收货地址
            </h3>

            <el-form label-width="80px">
              <el-form-item label="省份">
                <el-select
                  v-model="address.province_code"
                  placeholder="请选择省份"
                  filterable
                  @change="handleProvinceChange"
                >
                  <el-option
                    v-for="p in provinces"
                    :key="p.code"
                    :label="p.name"
                    :value="p.code"
                  />
                </el-select>
              </el-form-item>

              <el-form-item label="城市">
                <el-select
                  v-model="address.city_code"
                  placeholder="请选择城市"
                  filterable
                  :disabled="!address.province_code"
                  @change="handleCityChange"
                >
                  <el-option
                    v-for="c in cities"
                    :key="c.code"
                    :label="c.name"
                    :value="c.code"
                  />
                </el-select>
              </el-form-item>

              <el-form-item label="区县">
                <el-select
                  v-model="address.district_code"
                  placeholder="请选择区县"
                  filterable
                  :disabled="!address.city_code"
                >
                  <el-option
                    v-for="d in districts"
                    :key="d.code"
                    :label="d.name"
                    :value="d.code"
                  />
                </el-select>
              </el-form-item>
            </el-form>
          </div>
        </el-col>

        <!-- Test Items -->
        <el-col :span="14" :xs="24">
          <div class="items-section">
            <div class="section-header">
              <h3 class="section-title">
                <el-icon><Goods /></el-icon>
                测试商品
              </h3>
              <el-button type="primary" size="small" @click="addTestItem">
                <el-icon><Plus /></el-icon>
                添加商品
              </el-button>
            </div>

            <div class="items-list">
              <div v-for="(item, index) in testItems" :key="index" class="test-item">
                <div class="item-header">
                  <span class="item-number">商品 {{ index + 1 }}</span>
                  <el-button type="danger" link size="small" @click="removeTestItem(index)">
                    移除
                  </el-button>
                </div>

                <el-row :gutter="16">
                  <el-col :span="12" :xs="24">
                    <el-form-item label="商品">
                      <el-select
                        v-model="item.product_id"
                        placeholder="选择商品"
                        filterable
                        clearable
                        @change="(val) => handleProductSelect(item, val)"
                      >
                        <el-option
                          v-for="p in products"
                          :key="p.id"
                          :label="p.name"
                          :value="p.id"
                        />
                      </el-select>
                    </el-form-item>
                  </el-col>
                  <el-col :span="6" :xs="12">
                    <el-form-item label="数量">
                      <el-input-number v-model="item.quantity" :min="1" :max="99" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="6" :xs="12">
                    <el-form-item label="重量(g)">
                      <el-input-number v-model="item.weight" :min="1" />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="单价">
                  <el-input-number
                    v-model="item.price"
                    :min="0"
                    :precision="2"
                    style="width: 200px"
                  >
                    <template #prefix>¥</template>
                  </el-input-number>
                </el-form-item>
              </div>

              <el-empty v-if="testItems.length === 0" description="请添加测试商品" :image-size="80">
                <el-button type="primary" @click="addTestItem">
                  添加商品
                </el-button>
              </el-empty>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- Calculate Button -->
      <div class="calculate-action">
        <el-button
          type="primary"
          size="large"
          @click="calculateShipping"
          :loading="calculating"
          :disabled="!isFormValid"
        >
          <el-icon><Coin /></el-icon>
          计算运费
        </el-button>
      </div>
    </el-card>

    <!-- Result -->
    <el-card v-if="result" class="result-card" shadow="never">
      <div class="result-header">
        <div class="fee-display">
          <span class="fee-label">运费</span>
          <span class="fee-amount">¥ {{ formatAmount(result.shipping_fee) }}</span>
        </div>
      </div>

      <el-divider />

      <div class="result-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="匹配模板">
            {{ result.template_name }}
          </el-descriptions-item>
          <el-descriptions-item label="匹配区域">
            {{ result.zone_name }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="fee-breakdown">
          <h4 class="breakdown-title">费用明细</h4>
          <div class="breakdown-item">
            <span class="breakdown-label">计费方式：</span>
            <span class="breakdown-value">{{ getFeeTypeLabel(result.fee_detail.fee_type) }}</span>
          </div>
          <div class="breakdown-item" v-if="result.fee_detail.fee_type === 'by_weight'">
            <span class="breakdown-label">计算重量：</span>
            <span class="breakdown-value">{{ result.fee_detail.calculated_weight }}g</span>
          </div>
          <div class="breakdown-item" v-if="result.fee_detail.fee_type === 'by_count'">
            <span class="breakdown-label">计算数量：</span>
            <span class="breakdown-value">{{ result.fee_detail.calculated_units }} 件</span>
          </div>
          <div class="breakdown-item">
            <span class="breakdown-label">基础运费：</span>
            <span class="breakdown-value">¥ {{ formatAmount(result.fee_detail.first_fee) }}</span>
          </div>
          <div class="breakdown-item" v-if="getAdditionalUnits > 0">
            <span class="breakdown-label">续费 ({{ getAdditionalUnits }} 单位)：</span>
            <span class="breakdown-value">¥ {{ formatAdditionalFee }}</span>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Location, Goods, Plus, Coin } from '@element-plus/icons-vue'
import { calculateShippingFee, getRegions, type CalculateResult, type Region } from '@/api/shipping'
import { getProductList } from '@/api/product'
import type { Product } from '@/api/product'

// Types
interface TestItem {
  product_id: number | null
  quantity: number
  weight: number
  price: number
}

// State
const calculating = ref(false)
const address = reactive({
  province_code: '',
  city_code: '',
  district_code: ''
})

const testItems = ref<TestItem[]>([])
const result = ref<CalculateResult | null>(null)

// Region data loaded from API
const allProvinces = ref<Region[]>([])
const allCities = ref<Region[]>([])
const allDistricts = ref<Region[]>([])

// Computed filtered lists
const provinces = computed(() => allProvinces.value)
const cities = computed(() =>
  allCities.value.filter(c => c.parent_code === address.province_code)
)
const districts = computed(() =>
  allDistricts.value.filter(d => d.parent_code === address.city_code)
)

const products = ref<Product[]>([])

// Computed
const isFormValid = computed(() => {
  return address.province_code &&
         address.city_code &&
         address.district_code &&
         testItems.value.length > 0 &&
         testItems.value.every(item => item.weight > 0 && item.price > 0)
})

const getAdditionalUnits = computed(() => {
  if (!result.value) return 0
  return result.value.fee_detail.calculated_units || 0
})

const formatAdditionalFee = computed(() => {
  if (!result.value) return '0'
  const additional = parseFloat(result.value.fee_detail.additional_fee || '0')
  const units = result.value.fee_detail.calculated_units || 0
  return (additional * units).toFixed(2)
})

// Methods
const loadProducts = async () => {
  try {
    const data = await getProductList({ page: 1, page_size: 100 })
    products.value = data.list || []
  } catch (error) {
    console.error('Failed to load products:', error)
  }
}

const loadProvinces = async () => {
  try {
    const data = await getRegions()
    allProvinces.value = data
  } catch (error) {
    console.error('Failed to load provinces:', error)
  }
}

const loadCities = async (provinceCode: string) => {
  try {
    const data = await getRegions(provinceCode)
    allCities.value = data
  } catch (error) {
    console.error('Failed to load cities:', error)
  }
}

const loadDistricts = async (cityCode: string) => {
  try {
    const data = await getRegions(cityCode)
    allDistricts.value = data
  } catch (error) {
    console.error('Failed to load districts:', error)
  }
}

const handleProvinceChange = () => {
  address.city_code = ''
  address.district_code = ''
  allCities.value = []
  allDistricts.value = []
  if (address.province_code) {
    loadCities(address.province_code)
  }
}

const handleCityChange = () => {
  address.district_code = ''
  allDistricts.value = []
  if (address.city_code) {
    loadDistricts(address.city_code)
  }
}

const addTestItem = () => {
  testItems.value.push({
    product_id: null,
    quantity: 1,
    weight: 500,
    price: 0
  })
}

const removeTestItem = (index: number) => {
  testItems.value.splice(index, 1)
}

const handleProductSelect = (item: TestItem, productId: number | null) => {
  if (productId) {
    const product = products.value.find(p => p.id === productId)
    if (product) {
      item.weight = parseInt(product.weight) || 500
      item.price = product.price || 0
    }
  }
}

const calculateShipping = async () => {
  if (!isFormValid.value) {
    ElMessage.warning('请完整填写地址和商品信息')
    return
  }

  calculating.value = true
  try {
    const data = await calculateShippingFee({
      address: {
        province_code: address.province_code,
        city_code: address.city_code,
        district_code: address.district_code
      },
      items: testItems.value.map(item => ({
        product_id: item.product_id || 0,
        quantity: item.quantity,
        weight: item.weight,
        price: String(item.price)
      }))
    })
    result.value = data
    ElMessage.success('计算完成')
  } catch (error) {
    console.error('Failed to calculate:', error)
    ElMessage.error('计算失败，请检查配置')
  } finally {
    calculating.value = false
  }
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toFixed(2)
}

const getFeeTypeLabel = (feeType: string) => {
  const labels: Record<string, string> = {
    fixed: '固定运费',
    by_count: '按件计费',
    by_weight: '按重量计费',
    free: '免运费'
  }
  return labels[feeType] || feeType
}

// Lifecycle
onMounted(() => {
  loadProducts()
  loadProvinces()
  addTestItem() // Add one item by default
})
</script>

<style scoped>
.shipping-calculator-page {
  padding: 0;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
}

.header-desc {
  margin: 4px 0 0;
  font-size: 13px;
  color: #6B7280;
}

.calculator-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.address-section,
.items-section {
  padding: 20px 0;
}

.section-title {
  margin: 0 0 20px;
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.test-item {
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.item-number {
  font-weight: 600;
  color: #6366F1;
}

.calculate-action {
  margin-top: 24px;
  text-align: center;
}

.result-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.result-header {
  text-align: center;
  padding: 20px 0;
}

.fee-display {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
}

.fee-label {
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 8px;
}

.fee-amount {
  font-size: 36px;
  font-weight: 700;
  color: #6366F1;
}

.result-details {
  padding: 16px 0;
}

.fee-breakdown {
  margin-top: 20px;
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
}

.breakdown-title {
  margin: 0 0 16px;
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
}

.breakdown-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #E5E7EB;
}

.breakdown-item:last-child {
  border-bottom: none;
}

.breakdown-label {
  color: #6B7280;
}

.breakdown-value {
  font-weight: 500;
  color: #1E1B4B;
}

/* Descriptions */
:deep(.el-descriptions) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-descriptions__label) {
  background: #F9FAFB;
  font-weight: 500;
}

/* Responsive */
@media (max-width: 768px) {
  .calculator-card :deep(.el-row) {
    flex-direction: column;
  }

  .calculator-card :deep(.el-col) {
    max-width: 100%;
    flex: 0 0 100%;
  }

  .test-item :deep(.el-row) {
    flex-direction: column;
  }

  .test-item :deep(.el-col) {
    max-width: 100%;
    flex: 0 0 100%;
  }
}
</style>