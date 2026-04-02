<template>
  <div class="product-detail-page">
    <!-- Page Header -->
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-left">
          <el-button link @click="handleBack">
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('products.backToList') }}
          </el-button>
          <el-divider direction="vertical" />
          <h2 class="product-title">
            {{ product?.name || $t('common.loading') }}
            <el-tag v-if="product" :type="getStatusType(product.status)" size="small">
              {{ getStatusText(product.status) }}
            </el-tag>
          </h2>
        </div>
        <div class="header-right">
          <el-button
            v-if="product?.status === 'on_sale'"
            @click="handleTakeOffSale"
            :loading="statusLoading"
          >
            <el-icon><Hide /></el-icon>
            {{ $t('products.offSale') }}
          </el-button>
          <el-button
            v-else-if="product?.status === 'off_sale' || product?.status === 'draft'"
            type="success"
            @click="handlePutOnSale"
            :loading="statusLoading"
          >
            <el-icon><View /></el-icon>
            {{ $t('products.onSale') }}
          </el-button>
          <el-button @click="handleSave" type="primary" :loading="saveLoading">
            <el-icon><Check /></el-icon>
            {{ $t('products.saveChanges') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Loading State -->
    <el-skeleton v-if="loading" :rows="10" animated />

    <!-- Tab Layout -->
    <el-card v-else class="tabs-card" shadow="never">
      <el-tabs v-model="activeTab" class="product-tabs">
        <!-- Basic Info Tab -->
        <el-tab-pane :label="$t('products.basicInfo')" name="basic">
          <el-form :model="productForm" label-width="140px" :rules="formRules" ref="formRef">
            <!-- Basic Information Section -->
            <div class="form-section">
              <h3 class="section-title">{{ $t('products.basicInfo') }}</h3>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item :label="$t('products.productName')" prop="name">
                    <el-input v-model="productForm.name" :placeholder="$t('products.enterProductName')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="SKU" prop="sku">
                    <el-input v-model="productForm.sku" :placeholder="$t('products.enterSKU')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.brand')">
                    <el-input v-model="productForm.brand" :placeholder="$t('products.enterBrand')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.categoryId')">
                    <el-input-number v-model="productForm.category_id" :min="0" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item :label="$t('products.productDescription')">
                    <el-input v-model="productForm.description" type="textarea" :rows="4" :placeholder="$t('products.enterProductDescription')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.price')" prop="price">
                    <el-input-number v-model="productForm.price" :min="0" :precision="2" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.currency')">
                    <el-select v-model="productForm.currency" style="width: 100%">
                      <el-option label="USD" value="USD" />
                      <el-option label="EUR" value="EUR" />
                      <el-option label="GBP" value="GBP" />
                      <el-option label="CNY" value="CNY" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.costPrice')">
                    <el-input-number v-model="productForm.cost_price" :min="0" :precision="2" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.stock')">
                    <el-input-number v-model="productForm.stock" :min="0" style="width: 100%" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('common.status')">
                    <el-select v-model="productForm.status" style="width: 100%">
                      <el-option :label="$t('products.draft')" value="draft" />
                      <el-option :label="$t('products.onSale')" value="on_sale" />
                      <el-option :label="$t('products.offSale')" value="off_sale" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.isMatrixProduct')">
                    <el-switch v-model="productForm.is_matrix_product" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- Compliance Section -->
            <div class="form-section">
              <h3 class="section-title">{{ $t('products.complianceInfo') }}</h3>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item :label="$t('products.hsCode')">
                    <el-input v-model="productForm.hs_code" :placeholder="$t('products.enterHsCode')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('products.coo')">
                    <el-input v-model="productForm.coo" :placeholder="$t('products.enterCoo')" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item :label="$t('products.weight')">
                    <el-input v-model="productForm.weight" :placeholder="$t('products.weightPlaceholder')">
                      <template #append>
                        <el-select v-model="productForm.weight_unit" style="width: 80px">
                          <el-option label="kg" value="kg" />
                          <el-option label="g" value="g" />
                          <el-option label="lb" value="lb" />
                        </el-select>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="16">
                  <el-form-item :label="$t('products.dimensions')">
                    <el-row :gutter="8">
                      <el-col :span="8">
                        <el-input v-model="productForm.length" :placeholder="$t('products.length')">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                      <el-col :span="8">
                        <el-input v-model="productForm.width" :placeholder="$t('products.width')">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                      <el-col :span="8">
                        <el-input v-model="productForm.height" :placeholder="$t('products.height')">
                          <template #append>cm</template>
                        </el-input>
                      </el-col>
                    </el-row>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item :label="$t('products.dangerousGoods')">
                    <el-checkbox-group v-model="productForm.dangerous_goods">
                      <el-checkbox label="battery">{{ $t('products.battery') }}</el-checkbox>
                      <el-checkbox label="liquid">{{ $t('products.liquid') }}</el-checkbox>
                      <el-checkbox label="flammable">{{ $t('products.flammable') }}</el-checkbox>
                      <el-checkbox label="magnetic">{{ $t('products.magnetic') }}</el-checkbox>
                      <el-checkbox label="fragile">{{ $t('products.fragile') }}</el-checkbox>
                    </el-checkbox-group>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- Images Section -->
            <div class="form-section">
              <h3 class="section-title">{{ $t('products.productImage') }}</h3>
              <el-form-item :label="$t('products.imageUrl')">
                <div class="image-list">
                  <div v-for="(img, index) in productForm.images" :key="index" class="image-item">
                    <el-image :src="img" fit="cover" class="product-image">
                      <template #error>
                        <div class="image-placeholder">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <el-button
                      type="danger"
                      size="small"
                      circle
                      class="remove-btn"
                      @click="removeImage(index)"
                    >
                      <el-icon><Close /></el-icon>
                    </el-button>
                  </div>
                  <div class="add-image" @click="addImage">
                    <el-icon><Plus /></el-icon>
                    <span>{{ $t('products.addImage') }}</span>
                  </div>
                </div>
              </el-form-item>
            </div>
          </el-form>
        </el-tab-pane>

        <!-- Markets Tab -->
        <el-tab-pane :label="$t('products.markets')" name="markets">
          <div class="markets-section">
            <div class="section-header">
              <h3 class="section-title">{{ $t('products.marketAvailability') }}</h3>
              <el-button type="primary" @click="showPushToMarketDialog">
                <el-icon><Plus /></el-icon>
                {{ $t('products.pushToMarket') }}
              </el-button>
            </div>

            <el-table :data="productMarkets" v-loading="marketsLoading" stripe>
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

            <div v-if="productMarkets.length === 0 && !marketsLoading" class="empty-markets">
              <el-empty :description="$t('products.notOnAnyMarket')">
                <el-button type="primary" @click="showPushToMarketDialog">{{ $t('products.pushToMarket') }}</el-button>
              </el-empty>
            </div>
          </div>
        </el-tab-pane>

        <!-- Variants Tab -->
        <el-tab-pane :label="$t('products.variants')" name="variants">
          <div class="variants-section" v-loading="variantsLoading">
            <div class="section-header">
              <h3 class="section-title">{{ $t('products.productVariants') }}</h3>
              <el-button type="primary" size="small" @click="showAddVariantDialog">
                <el-icon><Plus /></el-icon>
                {{ $t('products.addVariant') }}
              </el-button>
            </div>
            <el-table :data="variants" stripe>
              <el-table-column :label="$t('products.skuCode')" prop="code" min-width="150" />
              <el-table-column :label="$t('products.attributes')" min-width="200">
                <template #default="{ row }">
                  <div class="attribute-tags">
                    <el-tag v-for="(value, key) in row.attributes" :key="key" size="small" class="attribute-tag">
                      {{ key }}: {{ value }}
                    </el-tag>
                    <span v-if="Object.keys(row.attributes || {}).length === 0" class="text-muted">-</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column :label="$t('products.price')" width="120" align="right">
                <template #default="{ row }">
                  {{ row.currency }} {{ row.price }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('products.stock')" prop="stock" width="100" align="center" />
              <el-table-column :label="$t('products.availableStock')" prop="available_stock" width="100" align="center" />
              <el-table-column :label="$t('products.safetyStock')" prop="safety_stock" width="100" align="center" />
              <el-table-column :label="$t('common.status')" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'enabled' ? 'success' : 'info'" size="small">
                    {{ row.status === 'enabled' ? $t('common.enabled') : $t('common.disabled') }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="$t('products.stockAlert')" width="100" align="center">
                <template #default="{ row }">
                  <el-tag v-if="row.is_low_stock" type="danger" size="small">{{ $t('products.lowStock') }}</el-tag>
                  <span v-else class="text-muted">-</span>
                </template>
              </el-table-column>
              <el-table-column :label="$t('common.actions')" width="120" align="center">
                <template #default="{ row }">
                  <el-button type="primary" link size="small" @click="handleEditVariant(row)">{{ $t('common.edit') }}</el-button>
                  <el-button type="danger" link size="small" @click="handleDeleteVariant(row)">{{ $t('common.delete') }}</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="variants.length === 0 && !variantsLoading" :description="$t('products.noVariants')">
              <el-button type="primary" @click="showAddVariantDialog">{{ $t('products.addVariant') }}</el-button>
            </el-empty>
          </div>
        </el-tab-pane>

        <!-- Pricing Tab -->
        <el-tab-pane :label="$t('products.pricing')" name="pricing">
          <div class="pricing-section" v-loading="marketsLoading">
            <div class="section-header">
              <h3 class="section-title">{{ $t('products.marketPrice') }}</h3>
              <el-button type="primary" size="small" @click="handleSavePricing" :loading="pricingSaveLoading">
                <el-icon><Check /></el-icon>
                {{ $t('products.savePricing') }}
              </el-button>
            </div>
            <el-table :data="pricingData" stripe>
              <el-table-column :label="$t('products.market')" min-width="150">
                <template #default="{ row }">
                  <div class="market-cell">
                    <span class="market-code">{{ row.market_code }}</span>
                    <span class="market-name">{{ row.market_name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column :label="$t('products.currency')" width="80" align="center">
                <template #default="{ row }">
                  <el-tag size="small">{{ row.currency }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="$t('products.basePrice')" width="150" align="right">
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
              <el-table-column :label="$t('products.compareAtPrice')" width="150" align="right">
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
              <el-table-column :label="$t('products.discount')" width="100" align="center">
                <template #default="{ row }">
                  <span v-if="row.compare_at_price_value > row.price_value && row.compare_at_price_value > 0" class="discount-badge">
                    -{{ Math.round((1 - row.price_value / row.compare_at_price_value) * 100) }}%
                  </span>
                  <span v-else class="text-muted">-</span>
                </template>
              </el-table-column>
              <el-table-column :label="$t('common.status')" width="100" align="center">
                <template #default="{ row }">
                  <el-switch v-model="row.is_enabled" size="small" />
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="pricingData.length === 0 && !marketsLoading" :description="$t('products.notOnAnyMarket')" />
          </div>
        </el-tab-pane>

        <!-- Localization Tab -->
        <el-tab-pane :label="$t('products.localization')" name="localization">
          <div class="localization-section" v-loading="localizationsLoading">
            <div class="section-header">
              <h3 class="section-title">{{ $t('products.multilingualContent') }}</h3>
              <el-button type="primary" size="small" @click="handleAddLocalization">
                <el-icon><Plus /></el-icon>
                {{ $t('products.addLanguage') }}
              </el-button>
            </div>

            <!-- Language Tabs -->
            <el-tabs v-model="activeLanguage" type="card" class="language-tabs">
              <el-tab-pane
                v-for="lang in supportedLanguages"
                :key="lang.code"
                :label="lang.name"
                :name="lang.code"
              >
                <div v-if="getLocalizationByLang(lang.code)" class="localization-content">
                  <el-form label-width="100px">
                    <el-form-item :label="$t('products.productName')">
                      <el-input
                        v-model="getLocalizationByLang(lang.code)!.name"
                        :placeholder="$t('products.enterLocalizedName')"
                        @change="handleUpdateLocalization(getLocalizationByLang(lang.code)!)"
                      />
                    </el-form-item>
                    <el-form-item :label="$t('products.productDescription')">
                      <el-input
                        v-model="getLocalizationByLang(lang.code)!.description"
                        type="textarea"
                        :rows="4"
                        :placeholder="$t('products.enterLocalizedDescription')"
                        @change="handleUpdateLocalization(getLocalizationByLang(lang.code)!)"
                      />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="danger" size="small" @click="handleDeleteLocalization(getLocalizationByLang(lang.code)!)">
                        {{ $t('products.deleteThisLanguage') }}
                      </el-button>
                    </el-form-item>
                  </el-form>
                </div>
                <div v-else class="no-localization">
                  <el-empty :description="$t('products.noContentForLang', { name: lang.name })">
                    <el-button type="primary" size="small" @click="handleCreateLocalization(lang.code)">
                      {{ $t('products.createContentForLang', { name: lang.name }) }}
                    </el-button>
                  </el-empty>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-tab-pane>

        <!-- Inventory Tab -->
        <el-tab-pane :label="$t('products.stock')" name="inventory">
          <div class="inventory-section" v-loading="inventoryLoading">
            <!-- Inventory Overview -->
            <div class="inventory-overview">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-statistic :title="$t('products.totalStock')" :value="product?.stock || 0" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.availableStock')" :value="skuInventory?.available_stock || 0" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.lockedStock')" :value="skuInventory?.locked_stock || 0" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.safetyStock')" :value="skuInventory?.safety_stock || 0" />
                </el-col>
              </el-row>
            </div>

            <!-- Warehouse Inventory -->
            <div class="warehouse-inventory">
              <div class="section-header">
                <h3 class="section-title">{{ $t('products.warehouseInventory') }}</h3>
                <el-button type="primary" size="small" @click="showAdjustStockDialog">
                  <el-icon><Edit /></el-icon>
                  {{ $t('products.adjustStock') }}
                </el-button>
              </div>
              <el-table :data="skuInventory?.warehouses || []" stripe>
                <el-table-column :label="$t('products.warehouse')" min-width="150">
                  <template #default="{ row }">
                    <span>{{ row.warehouse_name || `${$t('products.warehouse')} ${row.warehouse_id}` }}</span>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('products.availableStock')" prop="available_stock" width="120" align="center" />
                <el-table-column :label="$t('products.lockedStock')" prop="locked_stock" width="120" align="center" />
                <el-table-column :label="$t('products.totalStock')" width="120" align="center">
                  <template #default="{ row }">
                    {{ row.available_stock + row.locked_stock }}
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-if="!skuInventory?.warehouses?.length" :description="$t('products.noWarehouseInventory')" />
            </div>

            <!-- Inventory Logs -->
            <div class="inventory-logs">
              <h3 class="section-title">{{ $t('products.inventoryChangeLog') }}</h3>
              <el-table :data="inventoryLogs" stripe>
                <el-table-column :label="$t('products.time')" width="180">
                  <template #default="{ row }">
                    {{ row.created_at }}
                  </template>
                </el-table-column>
                <el-table-column :label="$t('products.type')" width="100">
                  <template #default="{ row }">
                    <el-tag :type="getLogTypeStyle(row.change_type)" size="small">
                      {{ getLogTypeText(row.change_type) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('products.changeQuantity')" width="120" align="right">
                  <template #default="{ row }">
                    <span :class="row.change_quantity >= 0 ? 'text-success' : 'text-danger'">
                      {{ row.change_quantity >= 0 ? '+' : '' }}{{ row.change_quantity }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('products.beforeChange')" prop="before_stock" width="100" align="center" />
                <el-table-column :label="$t('products.afterChange')" prop="after_stock" width="100" align="center" />
                <el-table-column :label="$t('products.remark')" prop="remark" min-width="150" />
              </el-table>
              <el-pagination
                v-if="inventoryLogsTotal > 0"
                class="pagination"
                background
                layout="total, prev, pager, next"
                :total="inventoryLogsTotal"
                :page-size="inventoryLogsPageSize"
                :current-page="inventoryLogsPage"
                @current-change="handleInventoryLogsPageChange"
              />
            </div>
          </div>
        </el-tab-pane>

        <!-- Reviews Tab -->
        <el-tab-pane :label="$t('products.reviewStats')" name="reviews">
          <div class="reviews-section" v-loading="reviewsLoading">
            <!-- Review Stats Overview -->
            <div class="review-stats-overview">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-statistic :title="$t('products.totalReviews')" :value="productStats?.total_reviews || 0" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.averageRating')" :value="productStats?.average_rating || '0'" :suffix="$t('products.stars')" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.reviewsWithImages')" :value="productStats?.with_image_count || 0" />
                </el-col>
                <el-col :span="6">
                  <el-statistic :title="$t('products.replyRate')" :value="productStats?.reply_rate || 0" :suffix="$t('products.replyRateUnit')">
                    <template #suffix>
                      <span style="font-size: 14px">{{ $t('products.replyRateUnit') }}</span>
                    </template>
                  </el-statistic>
                </el-col>
              </el-row>
            </div>

            <!-- Rating Distribution -->
            <div class="rating-distribution">
              <h3 class="section-title">{{ $t('products.ratingDistribution') }}</h3>
              <div class="rating-bars">
                <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="rating-bar-item">
                  <span class="star-label">{{ star }}{{ $t('products.stars') }}</span>
                  <div class="bar-container">
                    <div
                      class="bar-fill"
                      :style="{ width: getRatingPercentage(star) + '%' }"
                    ></div>
                  </div>
                  <span class="star-count">{{ getRatingCount(star) }}</span>
                </div>
              </div>
            </div>

            <!-- Rating Details -->
            <div class="rating-details">
              <el-row :gutter="20">
                <el-col :span="8">
                  <div class="rating-card">
                    <div class="rating-card-title">{{ $t('products.qualityRating') }}</div>
                    <div class="rating-card-value">{{ productStats?.quality_avg_rating || '0' }}</div>
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="rating-card">
                    <div class="rating-card-title">{{ $t('products.valueRating') }}</div>
                    <div class="rating-card-value">{{ productStats?.value_avg_rating || '0' }}</div>
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="rating-card">
                    <div class="rating-card-title">{{ $t('products.replyCount') }}</div>
                    <div class="rating-card-value">{{ productStats?.reply_count || 0 }}</div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Push to Market Dialog -->
    <el-dialog
      v-model="pushToMarketDialogVisible"
      :title="$t('products.pushToMarket')"
      width="500px"
      destroy-on-close
    >
      <el-form :model="pushToMarketForm" label-width="120px" ref="pushToMarketFormRef">
        <el-form-item :label="$t('products.market')" prop="markets" required>
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
        <el-form-item :label="$t('products.priceUSD')" prop="price" required>
          <el-input-number
            v-model="pushToMarketForm.price"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
          <div class="price-note">{{ $t('products.priceNote') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushToMarketDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleConfirmPushToMarket" :loading="pushToMarketLoading">
          {{ $t('products.pushToMarket') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Add Image Dialog -->
    <el-dialog v-model="addImageDialogVisible" :title="$t('products.addImageUrl')" width="400px">
      <el-input v-model="newImageUrl" :placeholder="$t('products.enterImageUrl')" />
      <template #footer>
        <el-button @click="addImageDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmAddImage">{{ $t('common.add') }}</el-button>
      </template>
    </el-dialog>

    <!-- Adjust Stock Dialog -->
    <el-dialog v-model="adjustStockDialogVisible" :title="$t('products.adjustStock')" width="450px">
      <el-form :model="adjustStockForm" label-width="100px">
        <el-form-item :label="$t('products.warehouse')">
          <el-select v-model="adjustStockForm.warehouse_id" :placeholder="$t('products.selectWarehouse')" style="width: 100%">
            <el-option
              v-for="warehouse in warehouses"
              :key="warehouse.id"
              :label="warehouse.name"
              :value="warehouse.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('products.adjustQuantity')">
          <el-input-number
            v-model="adjustStockForm.quantity"
            :step="1"
            style="width: 100%"
          />
          <div class="adjust-tip">{{ $t('products.stockInOutTip') }}</div>
        </el-form-item>
        <el-form-item :label="$t('products.remark')">
          <el-input v-model="adjustStockForm.remark" type="textarea" :rows="2" :placeholder="$t('products.enterRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustStockDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleAdjustStock">{{ $t('products.confirmAdjustment') }}</el-button>
      </template>
    </el-dialog>

    <!-- Variant Dialog -->
    <el-dialog v-model="variantDialogVisible" :title="isEditVariant ? $t('products.editVariant') : $t('products.addVariant')" width="550px">
      <el-form :model="variantForm" label-width="100px" ref="variantFormRef">
        <el-form-item :label="$t('products.skuCode')" required>
          <el-input v-model="variantForm.code" :placeholder="$t('products.enterSkuCode')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('products.price')">
              <el-input-number v-model="variantForm.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.currency')">
              <el-select v-model="variantForm.currency" style="width: 100%">
                <el-option label="USD" value="USD" />
                <el-option label="EUR" value="EUR" />
                <el-option label="GBP" value="GBP" />
                <el-option label="CNY" value="CNY" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('products.stock')">
              <el-input-number v-model="variantForm.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('products.safetyStock')">
              <el-input-number v-model="variantForm.safety_stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('products.enablePreSale')">
          <el-switch v-model="variantForm.pre_sale_enabled" />
        </el-form-item>
        <el-form-item :label="$t('products.attributes')">
          <div class="attributes-section">
            <div v-for="(value, key) in variantForm.attributes" :key="key" class="attribute-item">
              <span class="attribute-text">{{ key }}: {{ value }}</span>
              <el-button type="danger" link size="small" @click="handleRemoveAttribute(key)">
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <div class="add-attribute-row">
              <el-input v-model="newAttributeKey" :placeholder="$t('products.attributeName')" style="width: 120px" />
              <el-input v-model="newAttributeValue" :placeholder="$t('products.attributeValue')" style="width: 120px" />
              <el-button type="primary" size="small" @click="handleAddAttribute">{{ $t('common.add') }}</el-button>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="variantDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveVariant">{{ $t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Check, Plus, Picture, Close, Hide, View, Edit } from '@element-plus/icons-vue'
import {
  getProduct,
  updateProduct,
  getProductMarkets,
  updateProductMarket,
  pushToMarket,
  removeFromMarket,
  putOnSale,
  takeOffSale,
  getSKUsByProduct,
  createSKU,
  updateSKU,
  deleteSKU,
  getProductLocalizations,
  createProductLocalization,
  updateProductLocalization,
  deleteProductLocalization,
  type Product,
  type ProductMarket,
  type SKU,
  type CreateSKURequest,
  type ProductLocalization
} from '@/api/product'
import { getMarkets, type Market } from '@/api/market'
import {
  getSKUInventory,
  getInventoryLogs,
  adjustStock,
  getWarehouses,
  type SKUInventory,
  type InventoryLog,
  type Warehouse
} from '@/api/inventory'
import { getProductStats, type ProductStats } from '@/api/review'
import { t } from '@/plugins/i18n'

const route = useRoute()
const router = useRouter()

const productId = computed(() => Number(route.params.id))

// State
const loading = ref(false)
const saveLoading = ref(false)
const statusLoading = ref(false)
const marketsLoading = ref(false)
const pushToMarketLoading = ref(false)
const activeTab = ref('basic')
const product = ref<Product | null>(null)
const productMarkets = ref<ProductMarket[]>([])
const markets = ref<Market[]>([])
const pushToMarketDialogVisible = ref(false)
const addImageDialogVisible = ref(false)
const adjustStockDialogVisible = ref(false)
const newImageUrl = ref('')
const formRef = ref()
const pushToMarketFormRef = ref()

// Inventory state
const inventoryLoading = ref(false)
const skuInventory = ref<SKUInventory | null>(null)
const inventoryLogs = ref<InventoryLog[]>([])
const inventoryLogsTotal = ref(0)
const inventoryLogsPage = ref(1)
const inventoryLogsPageSize = ref(10)
const warehouses = ref<Warehouse[]>([])

// Reviews state
const reviewsLoading = ref(false)
const productStats = ref<ProductStats | null>(null)

// Adjust stock form
const adjustStockForm = reactive({
  warehouse_id: 0,
  quantity: 0,
  remark: ''
})

// Pricing state
const pricingSaveLoading = ref(false)
const pricingData = ref<Array<ProductMarket & { price_value: number; compare_at_price_value: number }>>([])

// Variants state
const variantsLoading = ref(false)
const variants = ref<SKU[]>([])
const variantDialogVisible = ref(false)
const variantFormRef = ref()
const variantForm = reactive({
  id: 0,
  code: '',
  price: 0,
  currency: 'USD',
  stock: 0,
  safety_stock: 0,
  pre_sale_enabled: false,
  attributes: {} as Record<string, string>
})
const isEditVariant = ref(false)
const newAttributeKey = ref('')
const newAttributeValue = ref('')

// Localization state
const localizationsLoading = ref(false)
const localizations = ref<ProductLocalization[]>([])
const activeLanguage = ref('en')
// Language code to translation key mapping
const languageNameKeys: Record<string, string> = {
  'en': 'products.langEn',
  'zh-CN': 'products.langZhCN',
  'ja': 'products.langJa',
  'ko': 'products.langKo',
  'es': 'products.langEs'
}

const supportedLanguages = computed(() => [
  { code: 'en', name: t(languageNameKeys['en']) },
  { code: 'zh-CN', name: t(languageNameKeys['zh-CN']) },
  { code: 'ja', name: t(languageNameKeys['ja']) },
  { code: 'ko', name: t(languageNameKeys['ko']) },
  { code: 'es', name: t(languageNameKeys['es']) }
])

// Form
const productForm = reactive({
  name: '',
  description: '',
  price: '0',
  currency: 'USD',
  cost_price: '0',
  stock: 0,
  status: 'draft',
  category_id: 0,
  sku: '',
  brand: '',
  tags: [] as string[],
  images: [] as string[],
  is_matrix_product: false,
  hs_code: '',
  coo: '',
  weight: '',
  weight_unit: 'kg',
  length: '',
  width: '',
  height: '',
  dangerous_goods: [] as string[]
})

// Push to market form
const pushToMarketForm = reactive({
  selectedMarkets: [] as number[],
  price: 0
})

// Available markets for push (excluding existing ones)
const availableMarketsForPush = computed(() => {
  const existingMarketIds = productMarkets.value.map(pm => pm.market_id)
  return markets.value.filter(m => m.is_active && !existingMarketIds.includes(m.id))
})

// Form rules
const formRules = {
  name: [{ required: true, message: () => t('products.enterProductName'), trigger: 'blur' }],
  price: [{ required: true, message: () => t('products.enterPrice'), trigger: 'blur' }]
}

// Helper functions
const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    on_sale: 'success',
    off_sale: 'warning',
    draft: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusKeyMap: Record<string, string> = {
    on_sale: 'products.onSale',
    off_sale: 'products.offSale',
    draft: 'products.draft'
  }
  const key = statusKeyMap[status]
  if (key) {
    const translated = t(key)
    return translated !== key ? translated : status
  }
  return status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

// Navigation
const handleBack = () => {
  router.push('/products')
}

// Load product data
const loadProduct = async () => {
  loading.value = true
  try {
    const data = await getProduct(productId.value)
    product.value = data
    // Populate form - price is already in yuan as string from API
    Object.assign(productForm, {
      name: data.name || '',
      description: data.description || '',
      price: data.price || '0',
      currency: data.currency || 'USD',
      cost_price: data.cost_price || '0',
      stock: data.stock || 0,
      status: data.status || 'draft',
      category_id: data.category_id || 0,
      sku: data.sku || '',
      brand: data.brand || '',
      tags: data.tags || [],
      images: data.images || [],
      is_matrix_product: data.is_matrix_product || false,
      hs_code: data.hs_code || '',
      coo: data.coo || '',
      weight: data.weight || '',
      weight_unit: data.weight_unit || 'kg',
      length: data.length || '',
      width: data.width || '',
      height: data.height || '',
      dangerous_goods: data.dangerous_goods || []
    })
  } catch (error) {
    console.error('Failed to load product:', error)
    ElMessage.error(t('products.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Load product markets
const loadProductMarkets = async () => {
  marketsLoading.value = true
  try {
    const response = await getProductMarkets(productId.value)
    productMarkets.value = response.list || []
  } catch (error) {
    console.error('Failed to load product markets:', error)
    ElMessage.error(t('products.loadMarketsFailed'))
  } finally {
    marketsLoading.value = false
  }
}

// Load all markets
const loadMarkets = async () => {
  try {
    const response = await getMarkets()
    markets.value = response.list || []
  } catch (error) {
    console.error('Failed to load markets:', error)
    ElMessage.error(t('products.loadMarketsFailed'))
  }
}

// Put on sale
const handlePutOnSale = async () => {
  statusLoading.value = true
  try {
    await putOnSale(productId.value)
    ElMessage.success(t('products.onSaleSuccess'))
    loadProduct()
  } catch (error) {
    console.error('Failed to put on sale:', error)
    ElMessage.error(t('products.onSaleFailed'))
  } finally {
    statusLoading.value = false
  }
}

// Take off sale
const handleTakeOffSale = async () => {
  statusLoading.value = true
  try {
    await takeOffSale(productId.value)
    ElMessage.success(t('products.offSaleSuccess'))
    loadProduct()
  } catch (error) {
    console.error('Failed to take off sale:', error)
    ElMessage.error(t('products.offSaleFailed'))
  } finally {
    statusLoading.value = false
  }
}

// Save product
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        await updateProduct({
          id: productId.value,
          name: productForm.name,
          description: productForm.description,
          price: productForm.price, // Already in yuan as string
          currency: productForm.currency,
          category_id: productForm.category_id,
          sku: productForm.sku,
          brand: productForm.brand,
          tags: productForm.tags,
          images: productForm.images,
          is_matrix_product: productForm.is_matrix_product,
          hs_code: productForm.hs_code,
          coo: productForm.coo,
          weight: productForm.weight,
          weight_unit: productForm.weight_unit,
          length: productForm.length,
          width: productForm.width,
          height: productForm.height,
          dangerous_goods: productForm.dangerous_goods
        })
        ElMessage.success(t('products.updateSuccess'))
        loadProduct()
      } catch (error) {
        console.error('Failed to update product:', error)
        ElMessage.error(t('products.updateFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

// Market actions
const handleMarketEnableChange = async (row: ProductMarket, enabled: boolean) => {
  try {
    await updateProductMarket(productId.value, row.market_id, { is_enabled: enabled })
    ElMessage.success(enabled ? t('products.marketEnabled') : t('products.marketDisabled'))
    loadProductMarkets()
  } catch (error) {
    console.error('Failed to update market:', error)
    ElMessage.error(t('products.updateMarketStatusFailed'))
    row.is_enabled = !enabled // Revert
  }
}

const handleMarketPriceChange = async (row: ProductMarket) => {
  try {
    await updateProductMarket(productId.value, row.market_id, {
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
    await removeFromMarket(productId.value, row.market_id)
    ElMessage.success(t('products.removedFromMarket'))
    loadProductMarkets()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove from market:', error)
      ElMessage.error(t('products.removeFromMarketFailed'))
    }
  }
}

// Push to market
const showPushToMarketDialog = () => {
  pushToMarketForm.selectedMarkets = []
  pushToMarketForm.price = parseFloat(productForm.price) || 0
  pushToMarketDialogVisible.value = true
}

const handleConfirmPushToMarket = async () => {
  if (pushToMarketForm.selectedMarkets.length === 0) {
    ElMessage.warning(t('products.selectAtLeastOneMarket'))
    return
  }
  if (pushToMarketForm.price <= 0) {
    ElMessage.warning(t('products.enterValidPrice'))
    return
  }

  pushToMarketLoading.value = true
  try {
    const prices = pushToMarketForm.selectedMarkets.map(() =>
      pushToMarketForm.price.toFixed(2)
    )

    const result = await pushToMarket(productId.value, {
      market_ids: pushToMarketForm.selectedMarkets,
      prices
    })

    pushToMarketDialogVisible.value = false
    ElMessage.success(
      t('products.pushToMarketSuccess', { success: result.success?.length || 0, failed: result.failed?.length || 0 })
    )
    loadProductMarkets()
  } catch (error) {
    console.error('Failed to push to market:', error)
    ElMessage.error(t('products.pushToMarketFailed'))
  } finally {
    pushToMarketLoading.value = false
  }
}

// Image management
const addImage = () => {
  newImageUrl.value = ''
  addImageDialogVisible.value = true
}

const confirmAddImage = () => {
  if (newImageUrl.value.trim()) {
    productForm.images.push(newImageUrl.value.trim())
    addImageDialogVisible.value = false
  }
}

const removeImage = (index: number) => {
  productForm.images.splice(index, 1)
}

// Inventory functions
const loadInventoryData = async () => {
  if (!productForm.sku) return

  inventoryLoading.value = true
  try {
    // Load SKU inventory
    const inventory = await getSKUInventory(productForm.sku)
    skuInventory.value = inventory

    // Load inventory logs
    const logs = await getInventoryLogs({
      page: inventoryLogsPage.value,
      page_size: inventoryLogsPageSize.value,
      sku_code: productForm.sku
    })
    inventoryLogs.value = logs.list || []
    inventoryLogsTotal.value = logs.total || 0
  } catch (error) {
    console.error('Failed to load inventory:', error)
  } finally {
    inventoryLoading.value = false
  }
}

const loadWarehouses = async () => {
  try {
    const response = await getWarehouses()
    warehouses.value = response || []
  } catch (error) {
    console.error('Failed to load warehouses:', error)
  }
}

const handleInventoryLogsPageChange = (page: number) => {
  inventoryLogsPage.value = page
  loadInventoryData()
}

const showAdjustStockDialog = () => {
  adjustStockForm.warehouse_id = warehouses.value.find(w => w.is_default)?.id || warehouses.value[0]?.id || 0
  adjustStockForm.quantity = 0
  adjustStockForm.remark = ''
  adjustStockDialogVisible.value = true
}

const handleAdjustStock = async () => {
  if (!productForm.sku) {
    ElMessage.warning(t('products.skuNotExist'))
    return
  }
  if (adjustStockForm.quantity === 0) {
    ElMessage.warning(t('products.enterAdjustQuantity'))
    return
  }

  try {
    await adjustStock({
      sku_code: productForm.sku,
      warehouse_id: adjustStockForm.warehouse_id,
      quantity: adjustStockForm.quantity,
      remark: adjustStockForm.remark
    })
    ElMessage.success(t('products.adjustSuccess'))
    adjustStockDialogVisible.value = false
    loadInventoryData()
    loadProduct() // Refresh product stock
  } catch (error) {
    console.error('Failed to adjust stock:', error)
    ElMessage.error(t('products.adjustFailed'))
  }
}

const getLogTypeStyle = (type: string) => {
  const styles: Record<string, string> = {
    manual: 'primary',
    order: 'warning',
    return: 'success',
    adjustment: 'info'
  }
  return styles[type] || 'info'
}

const getLogTypeText = (type: string) => {
  const texts: Record<string, string> = {
    manual: t('products.manual'),
    order: t('products.order'),
    return: t('products.return'),
    adjustment: t('products.adjustment')
  }
  return texts[type] || type
}

// Load product review stats
const loadProductStats = async () => {
  reviewsLoading.value = true
  try {
    const stats = await getProductStats(productId.value)
    productStats.value = stats
  } catch (error) {
    console.error('Failed to load product stats:', error)
    ElMessage.error(t('products.loadReviewStatsFailed'))
  } finally {
    reviewsLoading.value = false
  }
}

// Get rating count for specific star
const getRatingCount = (star: number): number => {
  if (!productStats.value?.rating_distribution) return 0
  const key = String(star) as '1' | '2' | '3' | '4' | '5'
  return productStats.value.rating_distribution[key] || 0
}

// Get rating percentage for specific star
const getRatingPercentage = (star: number): number => {
  if (!productStats.value?.total_reviews || productStats.value.total_reviews === 0) return 0
  const count = getRatingCount(star)
  return Math.round((count / productStats.value.total_reviews) * 100)
}

// Watch for tab change to load inventory data
watch(activeTab, (newTab) => {
  if (newTab === 'inventory') {
    loadWarehouses()
    loadInventoryData()
  } else if (newTab === 'pricing') {
    preparePricingData()
  } else if (newTab === 'variants') {
    loadVariants()
  } else if (newTab === 'localization') {
    loadLocalizations()
  } else if (newTab === 'reviews') {
    loadProductStats()
  }
})

// Localization functions
const loadLocalizations = async () => {
  localizationsLoading.value = true
  try {
    const response = await getProductLocalizations(productId.value)
    localizations.value = response.list || []
    // Set first available language as active
    if (localizations.value.length > 0) {
      activeLanguage.value = localizations.value[0].language_code
    }
  } catch (error) {
    console.error('Failed to load localizations:', error)
    ElMessage.error(t('products.loadLocalizationsFailed'))
  } finally {
    localizationsLoading.value = false
  }
}

const getLocalizationByLang = (langCode: string) => {
  return localizations.value.find(loc => loc.language_code === langCode)
}

const handleAddLocalization = () => {
  // Find first language without localization
  const existingLangs = localizations.value.map(loc => loc.language_code)
  const missingLang = supportedLanguages.value.find(lang => !existingLangs.includes(lang.code))
  if (missingLang) {
    activeLanguage.value = missingLang.code
  }
}

const handleCreateLocalization = async (langCode: string) => {
  try {
    await createProductLocalization({
      product_id: productId.value,
      language_code: langCode,
      name: productForm.name,
      description: productForm.description
    })
    ElMessage.success(t('products.createSuccess'))
    loadLocalizations()
  } catch (error) {
    console.error('Failed to create localization:', error)
    ElMessage.error(t('products.createFailed'))
  }
}

const handleUpdateLocalization = async (loc: ProductLocalization) => {
  try {
    await updateProductLocalization({
      id: loc.id,
      name: loc.name,
      description: loc.description
    })
    ElMessage.success(t('products.updateSuccess'))
  } catch (error) {
    console.error('Failed to update localization:', error)
    ElMessage.error(t('products.updateFailed'))
    loadLocalizations() // Reload to reset
  }
}

const handleDeleteLocalization = async (loc: ProductLocalization) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDeleteLocalization', { name: supportedLanguages.value.find(l => l.code === loc.language_code)?.name || loc.language_code }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteProductLocalization(loc.id)
    ElMessage.success(t('products.deleteSuccess'))
    loadLocalizations()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete localization:', error)
      ElMessage.error(t('products.deleteFailed'))
    }
  }
}

// Pricing functions
const preparePricingData = () => {
  pricingData.value = productMarkets.value.map(pm => ({
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
      if (item.price_value !== parseFloat(item.price) ||
          item.compare_at_price_value !== parseFloat(item.compare_at_price || '0') ||
          item.is_enabled !== productMarkets.value.find(pm => pm.market_id === item.market_id)?.is_enabled) {
        await updateProductMarket(productId.value, item.market_id, {
          price: item.price_value.toString(),
          compare_at_price: item.compare_at_price_value.toString(),
          is_enabled: item.is_enabled
        })
      }
    }
    ElMessage.success(t('products.savePricingSuccess'))
    loadProductMarkets()
  } catch (error) {
    console.error('Failed to save pricing:', error)
    ElMessage.error(t('products.savePricingFailed'))
  } finally {
    pricingSaveLoading.value = false
  }
}

// Variant functions
const loadVariants = async () => {
  variantsLoading.value = true
  try {
    const response = await getSKUsByProduct(productId.value)
    variants.value = response.list || []
  } catch (error) {
    console.error('Failed to load variants:', error)
    ElMessage.error(t('products.loadSKUsFailed'))
  } finally {
    variantsLoading.value = false
  }
}

const showAddVariantDialog = () => {
  isEditVariant.value = false
  variantForm.id = 0
  variantForm.code = ''
  variantForm.price = parseFloat(productForm.price) || 0
  variantForm.currency = productForm.currency
  variantForm.stock = 0
  variantForm.safety_stock = 0
  variantForm.pre_sale_enabled = false
  variantForm.attributes = {}
  newAttributeKey.value = ''
  newAttributeValue.value = ''
  variantDialogVisible.value = true
}

const handleEditVariant = (row: SKU) => {
  isEditVariant.value = true
  variantForm.id = row.id
  variantForm.code = row.code
  variantForm.price = parseFloat(row.price) || 0
  variantForm.currency = row.currency
  variantForm.stock = row.stock
  variantForm.safety_stock = row.safety_stock
  variantForm.pre_sale_enabled = row.pre_sale_enabled
  variantForm.attributes = { ...row.attributes }
  newAttributeKey.value = ''
  newAttributeValue.value = ''
  variantDialogVisible.value = true
}

const handleAddAttribute = () => {
  if (newAttributeKey.value && newAttributeValue.value) {
    variantForm.attributes[newAttributeKey.value] = newAttributeValue.value
    newAttributeKey.value = ''
    newAttributeValue.value = ''
  }
}

const handleRemoveAttribute = (key: string) => {
  delete variantForm.attributes[key]
}

const handleSaveVariant = async () => {
  if (!variantForm.code) {
    ElMessage.warning(t('products.enterSkuCode'))
    return
  }
  if (variantForm.price <= 0) {
    ElMessage.warning(t('products.enterValidPrice'))
    return
  }

  try {
    const data: CreateSKURequest = {
      product_id: productId.value,
      code: variantForm.code,
      price: variantForm.price.toFixed(2), // Already in yuan, convert to string
      currency: variantForm.currency,
      stock: variantForm.stock,
      safety_stock: variantForm.safety_stock,
      pre_sale_enabled: variantForm.pre_sale_enabled,
      attributes: variantForm.attributes
    }

    if (isEditVariant.value) {
      await updateSKU({ ...data, id: variantForm.id })
      ElMessage.success(t('products.variantUpdateSuccess'))
    } else {
      await createSKU(data)
      ElMessage.success(t('products.variantCreateSuccess'))
    }

    variantDialogVisible.value = false
    loadVariants()
  } catch (error) {
    console.error('Failed to save variant:', error)
    ElMessage.error(isEditVariant.value ? t('products.variantUpdateFailed') : t('products.variantCreateFailed'))
  }
}

const handleDeleteVariant = async (row: SKU) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDeleteVariant', { code: row.code }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteSKU(row.id)
    ElMessage.success(t('products.variantDeleteSuccess'))
    loadVariants()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete variant:', error)
      ElMessage.error(t('products.variantDeleteFailed'))
    }
  }
}

// Initialize
onMounted(() => {
  loadProduct()
  loadProductMarkets()
  loadMarkets()
  loadWarehouses()
})
</script>

<style scoped>
.product-detail-page {
  padding: 0;
}

/* Header */
.header-card {
  margin-bottom: 20px;
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

.product-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* Tabs Card */
.tabs-card {
  min-height: 500px;
}

/* Form Sections */
.form-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #E5E7EB;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

/* Markets Section */
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

.price-note {
  font-size: 12px;
  color: #F59E0B;
  margin-top: 8px;
  line-height: 1.4;
}

/* Image Management */
.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.image-item {
  position: relative;
  width: 120px;
  height: 120px;
}

.product-image {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #9CA3AF;
}

.remove-btn {
  position: absolute;
  top: -8px;
  right: -8px;
}

.add-image {
  width: 120px;
  height: 120px;
  border: 2px dashed #D1D5DB;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  color: #9CA3AF;
  transition: all 0.2s;
}

.add-image:hover {
  border-color: #409EFF;
  color: #409EFF;
}

/* Inventory Section */
.inventory-section {
  padding: 0;
}

.inventory-overview {
  padding: 20px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 24px;
}

.warehouse-inventory {
  margin-bottom: 24px;
}

.inventory-logs {
  margin-top: 24px;
}

.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}

.text-success {
  color: #10B981;
  font-weight: 500;
}

.text-danger {
  color: #EF4444;
  font-weight: 500;
}

.adjust-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

/* Pricing Section */
.pricing-section {
  padding: 0;
}

.discount-badge {
  background: #ECFDF5;
  color: #059669;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

/* Variants Section */
.variants-section {
  padding: 0;
}

.attribute-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.attribute-tag {
  margin: 2px;
}

.attributes-section {
  width: 100%;
}

.attribute-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.attribute-text {
  font-size: 14px;
}

.add-attribute-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-top: 8px;
}

/* Localization Section */
.localization-section {
  padding: 0;
}

.language-tabs {
  margin-top: 16px;
}

.localization-content {
  padding: 16px 0;
}

.no-localization {
  padding: 20px 0;
}

/* Reviews Section */
.reviews-section {
  padding: 0;
}

.review-stats-overview {
  padding: 20px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 24px;
}

.rating-distribution {
  margin-bottom: 24px;
}

.rating-bars {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-width: 500px;
}

.rating-bar-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.star-label {
  width: 40px;
  font-size: 14px;
  color: #6B7280;
}

.bar-container {
  flex: 1;
  height: 20px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.star-count {
  width: 50px;
  text-align: right;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.rating-details {
  margin-top: 24px;
}

.rating-card {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
}

.rating-card-title {
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 8px;
}

.rating-card-value {
  font-size: 28px;
  font-weight: 700;
  color: #6366F1;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-left {
    flex-direction: column;
    align-items: flex-start;
  }

  .product-title {
    font-size: 18px;
  }
}
</style>