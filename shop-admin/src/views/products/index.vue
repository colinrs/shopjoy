<template>
  <div class="products-page">
    <!-- Market Filter Bar -->
    <el-card
      class="market-filter-card"
      shadow="never"
    >
      <div class="market-filter-bar">
        <el-radio-group
          v-model="selectedMarket"
          @change="handleMarketChange"
        >
          <el-radio-button value="">
            {{ $t('common.all') }}
          </el-radio-button>
          <el-radio-button
            v-for="market in markets"
            :key="market.id"
            :value="market.id"
          >
            {{ market.code }}
          </el-radio-button>
        </el-radio-group>
      </div>
    </el-card>

    <!-- Search & Filter Bar -->
    <el-card
      class="filter-card"
      shadow="never"
    >
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('products.searchPlaceholder')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="filterStatus"
            :placeholder="$t('products.status')"
            clearable
            class="filter-select"
            @change="handleSearch"
          >
            <el-option
              :label="$t('common.all')"
              value=""
            />
            <el-option
              :label="$t('products.onSale')"
              value="on_sale"
            />
            <el-option
              :label="$t('products.offSale')"
              value="off_sale"
            />
            <el-option
              :label="$t('products.draft')"
              value="draft"
            />
          </el-select>
          <el-cascader
            v-model="filterCategory"
            :options="categories"
            :props="{ checkStrictly: true, emitPath: false, value: 'id', label: 'name' }"
            :placeholder="$t('products.category')"
            clearable
            class="filter-select"
            @change="handleSearch"
          />
          <el-button
            type="primary"
            @click="handleSearch"
          >
            <el-icon><Search /></el-icon>{{ $t('common.search') }}
          </el-button>
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>{{ $t('common.export') }}
          </el-button>
          <el-button
            type="primary"
            @click="handleAdd"
          >
            <el-icon><Plus /></el-icon>{{ $t('products.addProduct') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Bulk Actions -->
    <div
      v-if="selectedProducts.length > 0"
      class="bulk-actions"
    >
      <span class="selected-count">{{ $t('products.selectedCount', { count: selectedProducts.length }) }}</span>
      <el-button
        size="small"
        @click="handleBatchOnSale"
      >
        {{ $t('products.batchOnSale') }}
      </el-button>
      <el-button
        size="small"
        @click="handleBatchOffSale"
      >
        {{ $t('products.batchOffSale') }}
      </el-button>
      <el-button
        size="small"
        type="success"
        @click="handleBatchPushToMarket"
      >
        {{ $t('products.pushToMarket') }}
      </el-button>
      <el-button
        size="small"
        type="danger"
        @click="handleBatchDelete"
      >
        {{ $t('products.batchDelete') }}
      </el-button>
      <el-button
        size="small"
        @click="handleClearSelection"
      >
        {{ $t('common.clearSelection') }}
      </el-button>
    </div>

    <!-- Products Table -->
    <el-card
      class="table-card"
      shadow="never"
    >
      <!-- Skeleton loading -->
      <TableSkeleton
        v-if="loading && productList.length === 0"
        :rows="10"
        :columns="7"
      />

      <!-- Empty state -->
      <EmptyState
        v-else-if="productList.length === 0"
        :title="$t('products.noProducts')"
        :description="$t('products.noProductsDesc')"
      />

      <!-- Actual table -->
      <Table
        v-if="!loading && productList.length > 0"
        ref="tableRef"
        :data="productList"
        :loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          :label="$t('products.productInfo')"
          min-width="300"
        >
          <template #default="{ row }">
            <div class="product-cell">
              <el-image
                :src="row.images?.[0] || defaultImage"
                fit="cover"
                class="product-thumb"
                :preview-src-list="row.images"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="product-details">
                <p class="product-name">
                  {{ row.name }}
                </p>
                <p class="product-sku">
                  SKU: {{ row.sku || $t('common.noData') }}
                </p>
                <div class="product-tags">
                  <el-tag
                    v-if="row.tags?.includes('hot')"
                    size="small"
                    type="danger"
                    effect="plain"
                  >
                    {{ $t('products.hot') }}
                  </el-tag>
                  <el-tag
                    v-if="row.tags?.includes('new')"
                    size="small"
                    type="success"
                    effect="plain"
                  >
                    {{ $t('products.new') }}
                  </el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.price')"
          width="150"
          align="right"
        >
          <template #default="{ row }">
            <div class="price-cell">
              <p class="sale-price">
                {{ row.currency || 'USD' }}{{ formatPrice(row.price) }}
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="stock"
          :label="$t('products.stock')"
          width="120"
          align="center"
        >
          <template #default="{ row }">
            <el-tag
              :type="getStockType(row.stock)"
              size="small"
            >
              {{ row.stock }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.markets')"
          min-width="150"
          align="center"
        >
          <template #default="{ row }">
            <div class="market-tags">
              <el-tag
                v-for="market in row.markets"
                :key="market.market_id"
                :type="market.is_enabled ? 'success' : 'info'"
                size="small"
                class="market-tag"
              >
                {{ market.market_code }}
              </el-tag>
              <span
                v-if="!row.markets || row.markets.length === 0"
                class="no-markets"
              >
                {{ $t('products.noMarkets') }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('products.categoryId')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            {{ row.category_id || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          :label="$t('products.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'on_sale'"
              :inactive-value="'off_sale'"
              @change="(val: string) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="180"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="handlePreview(row)"
            >
              {{ $t('products.preview') }}
            </el-button>
            <el-dropdown
              trigger="click"
              @command="(cmd: string) => handleCommand(cmd, row)"
            >
              <el-button
                type="primary"
                link
                size="small"
              >
                {{ $t('common.more') }}<el-icon class="el-icon--right">
                  <ArrowDown />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="copy">
                    {{ $t('products.copy') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="top"
                    divided
                  >
                    {{ $t('products.setTop') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    command="delete"
                    type="danger"
                  >
                    {{ $t('common.delete') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </Table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? $t('products.editProduct') : $t('products.addProduct')"
      width="900px"
      destroy-on-close
    >
      <el-tabs
        v-model="activeProductTab"
        class="product-form-tabs"
      >
        <el-tab-pane
          :label="$t('products.basicInfo')"
          name="basic"
        >
          <el-form
            ref="formRef"
            :model="productForm"
            label-width="100px"
            :rules="formRules"
          >
            <el-row :gutter="20">
              <el-col :span="24">
                <el-form-item
                  :label="$t('products.productName')"
                  prop="name"
                >
                  <el-input
                    v-model="productForm.name"
                    :placeholder="$t('products.enterProductName')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item
                  :label="$t('products.productCategory')"
                  prop="category_id"
                >
                  <el-cascader
                    v-model="productForm.category_id"
                    :options="categories"
                    :props="{ checkStrictly: true, emitPath: false, value: 'id', label: 'name' }"
                    :placeholder="$t('products.selectCategory')"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.brand')">
                  <el-input
                    v-model="productForm.brand"
                    :placeholder="$t('products.enterBrand')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item :label="$t('products.productImage')">
                  <el-upload
                    class="avatar-uploader"
                    action="#"
                    :show-file-list="false"
                    :auto-upload="false"
                    :on-change="handleImageChange"
                  >
                    <img
                      v-if="productForm.image"
                      :src="productForm.image"
                      class="avatar"
                    >
                    <el-icon
                      v-else
                      class="avatar-uploader-icon"
                    >
                      <Plus />
                    </el-icon>
                  </el-upload>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item :label="$t('products.productDescription')">
                  <el-input
                    v-model="productForm.description"
                    type="textarea"
                    rows="4"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane
          :label="$t('products.priceStock')"
          name="price"
        >
          <el-form
            :model="productForm"
            label-width="100px"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item
                  :label="$t('products.productPrice')"
                  prop="price"
                >
                  <el-input-number
                    v-model="productForm.price"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.originalPrice')">
                  <el-input-number
                    v-model="productForm.original_price"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item
                  :label="$t('products.stockQuantity')"
                  prop="stock"
                >
                  <el-input-number
                    v-model="productForm.stock"
                    :min="0"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.costPrice')">
                  <el-input-number
                    v-model="productForm.cost_price"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.currency')">
                  <el-select
                    v-model="productForm.currency"
                    style="width: 100%"
                  >
                    <el-option
                      label="USD"
                      value="USD"
                    />
                    <el-option
                      label="EUR"
                      value="EUR"
                    />
                    <el-option
                      label="GBP"
                      value="GBP"
                    />
                    <el-option
                      label="CNY"
                      value="CNY"
                    />
                    <el-option
                      label="JPY"
                      value="JPY"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane
          :label="$t('products.logisticsInfo')"
          name="logistics"
        >
          <el-form
            :model="productForm"
            label-width="100px"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('products.weight')">
                  <el-input-number
                    v-model="productForm.weight"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.weightUnit')">
                  <el-select
                    v-model="productForm.weight_unit"
                    style="width: 100%"
                  >
                    <el-option
                      label="kg"
                      value="kg"
                    />
                    <el-option
                      label="g"
                      value="g"
                    />
                    <el-option
                      label="lb"
                      value="lb"
                    />
                    <el-option
                      label="oz"
                      value="oz"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('products.length')">
                  <el-input-number
                    v-model="productForm.length"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('products.width')">
                  <el-input-number
                    v-model="productForm.width"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('products.height')">
                  <el-input-number
                    v-model="productForm.height"
                    :min="0"
                    :precision="2"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane
          :label="$t('products.complianceInfo')"
          name="compliance"
        >
          <el-form
            :model="productForm"
            label-width="100px"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('products.hsCode')">
                  <el-input
                    v-model="productForm.hs_code"
                    :placeholder="$t('products.enterHsCode')"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('products.coo')">
                  <el-select
                    v-model="productForm.coo"
                    :placeholder="$t('products.selectCoo')"
                    style="width: 100%"
                  >
                    <el-option
                      label="China"
                      value="CN"
                    />
                    <el-option
                      label="United States"
                      value="US"
                    />
                    <el-option
                      label="European Union"
                      value="EU"
                    />
                    <el-option
                      label="Japan"
                      value="JP"
                    />
                    <el-option
                      label="South Korea"
                      value="KR"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item :label="$t('products.dangerousGoods')">
                  <el-checkbox-group v-model="productForm.dangerous_goods">
                    <el-checkbox label="flammable">
                      {{ $t('products.flammable') }}
                    </el-checkbox>
                    <el-checkbox label="explosive">
                      {{ $t('products.explosive') }}
                    </el-checkbox>
                    <el-checkbox label="corrosive">
                      {{ $t('products.corrosive') }}
                    </el-checkbox>
                    <el-checkbox label="radioactive">
                      {{ $t('products.radioactive') }}
                    </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="dialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="saveLoading"
          @click="handleSave"
        >
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Push to Market Dialog -->
    <el-dialog
      v-model="pushToMarketDialogVisible"
      :title="$t('products.pushToMarket')"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="pushToMarketFormRef"
        :model="pushToMarketForm"
        label-width="120px"
      >
        <el-form-item
          :label="$t('products.markets')"
          prop="markets"
          required
        >
          <el-checkbox-group v-model="pushToMarketForm.selectedMarkets">
            <el-checkbox
              v-for="market in availableMarkets"
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
        <el-button @click="pushToMarketDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="pushToMarketLoading"
          @click="handleConfirmPushToMarket"
        >
          {{ $t('products.pushToMarket') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Preview Dialog -->
    <el-dialog
      v-model="previewDialogVisible"
      :title="$t('products.preview')"
      width="700px"
      destroy-on-close
    >
      <div
        v-if="previewProduct"
        class="preview-content"
      >
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="preview-image">
              <el-image
                v-if="previewProduct.images && previewProduct.images.length > 0"
                :src="previewProduct.images[0]"
                fit="cover"
                class="preview-main-image"
                :preview-src-list="previewProduct.images"
              />
              <div
                v-else
                class="preview-image-placeholder"
              >
                <el-icon><Picture /></el-icon>
              </div>
            </div>
          </el-col>
          <el-col :span="16">
            <div class="preview-info">
              <h3 class="preview-title">
                {{ previewProduct.name }}
              </h3>
              <p class="preview-sku">
                SKU: {{ previewProduct.sku || $t('common.noData') }}
              </p>
              <div class="preview-tags">
                <el-tag
                  v-if="previewProduct.tags?.includes('hot')"
                  size="small"
                  type="danger"
                >
                  {{ $t('products.hot') }}
                </el-tag>
                <el-tag
                  v-if="previewProduct.tags?.includes('new')"
                  size="small"
                  type="success"
                >
                  {{ $t('products.new') }}
                </el-tag>
              </div>
              <div class="preview-price">
                <span class="price-label">{{ $t('products.price') }}:</span>
                <span class="price-value">{{ previewProduct.currency || 'USD' }} {{ previewProduct.price }}</span>
              </div>
              <div class="preview-detail-row">
                <span class="detail-label">{{ $t('products.stock') }}:</span>
                <el-tag
                  :type="getStockType(previewProduct.stock)"
                  size="small"
                >
                  {{ previewProduct.stock }}
                </el-tag>
              </div>
              <div class="preview-detail-row">
                <span class="detail-label">{{ $t('products.status') }}:</span>
                <el-tag
                  :type="previewProduct.status === 'on_sale' ? 'success' : 'warning'"
                  size="small"
                >
                  {{ previewProduct.status === 'on_sale' ? $t('products.onSale') : previewProduct.status === 'off_sale' ? $t('products.offSale') : $t('products.draft') }}
                </el-tag>
              </div>
              <div class="preview-detail-row">
                <span class="detail-label">{{ $t('products.brand') }}:</span>
                <span>{{ previewProduct.brand || '-' }}</span>
              </div>
              <div class="preview-detail-row">
                <span class="detail-label">{{ $t('products.categoryId') }}:</span>
                <span>{{ previewProduct.category_id || '-' }}</span>
              </div>
            </div>
          </el-col>
        </el-row>
        <el-row
          :gutter="20"
          class="preview-description-row"
        >
          <el-col :span="24">
            <div class="preview-description">
              <span class="detail-label">{{ $t('products.productDescription') }}:</span>
              <p class="description-content">
                {{ previewProduct.description || $t('common.noData') }}
              </p>
            </div>
          </el-col>
        </el-row>
        <el-row
          v-if="previewProduct.markets && previewProduct.markets.length > 0"
          :gutter="20"
        >
          <el-col :span="24">
            <div class="preview-markets">
              <span class="detail-label">{{ $t('products.markets') }}:</span>
              <div class="market-tags">
                <el-tag
                  v-for="market in previewProduct.markets"
                  :key="market.market_id"
                  :type="market.is_enabled ? 'success' : 'info'"
                  size="small"
                  class="market-tag"
                >
                  {{ market.market_code }}: {{ market.price }} {{ market.currency }}
                </el-tag>
              </div>
            </div>
          </el-col>
        </el-row>
        <el-row
          :gutter="20"
          class="preview-meta-row"
        >
          <el-col :span="12">
            <div class="preview-meta">
              <span class="detail-label">{{ $t('common.createdAt') }}:</span>
              <span>{{ previewProduct.created_at || '-' }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="preview-meta">
              <span class="detail-label">{{ $t('common.updatedAt') }}:</span>
              <span>{{ previewProduct.updated_at || '-' }}</span>
            </div>
          </el-col>
        </el-row>
      </div>
      <template #footer>
        <el-button @click="previewDialogVisible = false">
          {{ $t('common.close') }}
        </el-button>
        <el-button
          type="primary"
          @click="handlePreviewEdit"
        >
          {{ $t('common.edit') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type UploadFile } from 'element-plus'
import { Search, Plus, Download, Picture, ArrowDown } from '@element-plus/icons-vue'
import { getProductList, getProduct, pushToMarket, putOnSale, takeOffSale, createProduct, deleteProduct, exportProductsUrl, batchUpdateProducts, type Product, type ListProductsParams, type BatchUpdateProductRequest } from '@/api/product'
import { getMarkets, type Market } from '@/api/market'
import { getCategoryTree, type CategoryTree } from '@/api/category'
import { uploadImage } from '@/api/upload'
import { TableSkeleton } from '@/components/skeleton'
import Table from '@/components/common/Table.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { t } from '@/plugins/i18n'
import { downloadFile } from '@/utils/download'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { handleError } = useErrorHandler()

const loading = ref(false)
const saveLoading = ref(false)
const pushToMarketLoading = ref(false)
const dialogVisible = ref(false)
const pushToMarketDialogVisible = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const filterStatus = ref('')
const filterCategory = ref<number | ''>('')
const selectedMarket = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedProducts = ref<Product[]>([])
const tableRef = ref()
const formRef = ref()
const pushToMarketFormRef = ref()
const defaultImage = 'https://via.placeholder.com/80'

// Market data
const markets = ref<Market[]>([])
const productList = ref<Product[]>([])

// Category data
const categories = ref<CategoryTree[]>([])

// Push to market form
const pushToMarketForm = reactive({
  selectedMarkets: [] as number[],
  price: 0
})

// Available markets for push to market dialog
const availableMarkets = computed(() => markets.value.filter(m => m.is_active))

const productForm = reactive({
  id: null as number | null,
  name: '',
  price: 0,
  original_price: 0,
  stock: 0,
  category: '',
  category_id: 0,
  image: '',
  images: [] as string[],
  description: '',
  sku: '',
  brand: '',
  tags: [] as string[],
  is_matrix_product: false,
  hs_code: '',
  coo: '',
  weight: '',
  weight_unit: 'kg',
  length: '',
  width: '',
  height: '',
  dangerous_goods: [] as string[],
  currency: 'USD',
  cost_price: 0,
  status: 'draft'
})

const formRules = {
  name: [{ required: true, message: '', trigger: 'blur' }],
  price: [{ required: true, message: '', trigger: 'blur' }],
  stock: [{ required: true, message: '', trigger: 'blur' }],
  category_id: [{ required: true, message: '', trigger: 'change' }]
}

const formatPrice = (price: string) => {
  return price
}

const getStockType = (stock: number) => {
  if (stock === 0) return 'danger'
  if (stock < 20) return 'warning'
  return 'success'
}

const handleSearch = () => {
  currentPage.value = 1
  loadProducts()
}

const handleExport = async () => {
  try {
    loading.value = true
    // Build export params from current filters
    const { url, params } = exportProductsUrl({
      name: searchQuery.value || undefined,
      status: filterStatus.value || undefined,
      category_id: filterCategory.value || undefined,
      market_id: selectedMarket.value || undefined
    })

    await downloadFile(url, params)
  } catch (error) {
    handleError(error)
    // Error message is handled by downloadFile utility
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(productForm, {
    id: null,
    name: '',
    price: 0,
    original_price: 0,
    stock: 0,
    category: '',
    image: '',
    description: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: Product) => {
  router.push(`/products/${row.id}`)
}

const previewDialogVisible = ref(false)
const previewProduct = ref<Product | null>(null)
const previewLoading = ref(false)
const activeProductTab = ref('basic')

const handlePreview = async (row: Product) => {
  previewLoading.value = true
  previewDialogVisible.value = true
  try {
    // Fetch fresh product data from API
    const freshProduct = await getProduct(row.id)
    previewProduct.value = freshProduct
  } catch (error) {
    handleError(error, t('products.previewLoadFailed'))
    previewDialogVisible.value = false
  } finally {
    previewLoading.value = false
  }
}

const handlePreviewEdit = () => {
  if (previewProduct.value) {
    previewDialogVisible.value = false
    router.push(`/products/${previewProduct.value.id}`)
  }
}

const handleCommand = async (cmd: string, row: Product) => {
  switch (cmd) {
    case 'copy':
      // Call API to create a duplicate product
      try {
        const response = await createProduct({
          name: `${row.name} (Copy)`,
          description: row.description || '',
          price: row.price || '0',
          currency: row.currency || 'USD',
          cost_price: row.cost_price || '',
          category_id: row.category_id || 0,
          sku: row.sku ? `${row.sku}-copy` : '',
          brand: row.brand || '',
          tags: row.tags || [],
          images: row.images || [],
          is_matrix_product: row.is_matrix_product || false,
          hs_code: row.hs_code || '',
          coo: row.coo || '',
          weight: row.weight || '',
          weight_unit: row.weight_unit || 'kg',
          length: row.length || '',
          width: row.width || '',
          height: row.height || '',
          dangerous_goods: row.dangerous_goods || []
        })
        ElMessage.success(t('products.copySuccess'))
        // Navigate to edit the new product
        router.push(`/products/${response.id}`)
      } catch (error) {
        handleError(error, t('products.copyFailed'))
      }
      break
    case 'top':
      ElMessage.success(t('products.setTopSuccess'))
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = async (row: Product) => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmDelete', { name: row.name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    await deleteProduct(row.id)
    ElMessage.success(t('products.deleteSuccess'))
    loadProducts()
  } catch (error) {
    if ((error as Error).message !== 'cancel') {
      handleError(error, t('products.deleteFailed'))
    }
  }
}

const handleStatusChange = async (row: Product, val: string) => {
  const statusKey = val === 'on_sale' ? 'products.onSaleSuccess' : 'products.offSaleSuccess'
  try {
    if (val === 'on_sale') {
      await putOnSale(row.id)
    } else {
      await takeOffSale(row.id)
    }
    ElMessage.success(t(statusKey))
    loadProducts()
  } catch (error) {
    const errorKey = val === 'on_sale' ? 'products.onSaleFailed' : 'products.offSaleFailed'
    handleError(error, t(errorKey))
    // Revert the switch
    row.status = val === 'on_sale' ? 'off_sale' : 'on_sale'
  }
}

const handleSelectionChange = (selection: Product[]) => {
  selectedProducts.value = selection
}

const handleClearSelection = () => {
  tableRef.value?.clearSelection()
  selectedProducts.value = []
}

const handleBatchOnSale = async () => {
  try {
    const data: BatchUpdateProductRequest = {
      product_ids: selectedProducts.value.map(p => p.id),
      update_fields: {
        status: 'on_sale'
      }
    }

    const res = await batchUpdateProducts(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('products.batchOnSalePartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('products.batchOnSaleSuccess', { count: successCount }))
    }

    selectedProducts.value = []
    loadProducts()
  } catch (error) {
    handleError(error, t('products.batchOnSaleFailed'))
  }
}

const handleBatchOffSale = async () => {
  try {
    const data: BatchUpdateProductRequest = {
      product_ids: selectedProducts.value.map(p => p.id),
      update_fields: {
        status: 'off_sale'
      }
    }

    const res = await batchUpdateProducts(data)

    const successCount = res.success?.length || 0
    const failedCount = res.failed?.length || 0

    if (failedCount > 0) {
      ElMessage.warning(t('products.batchOffSalePartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('products.batchOffSaleSuccess', { count: successCount }))
    }

    selectedProducts.value = []
    loadProducts()
  } catch (error) {
    handleError(error, t('products.batchOffSaleFailed'))
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      t('products.confirmBatchDeleteWarning', { count: selectedProducts.value.length }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const deletePromises = selectedProducts.value.map(p => deleteProduct(p.id))
    const results = await Promise.allSettled(deletePromises)

    const failedCount = results.filter(r => r.status === 'rejected').length
    const successCount = results.filter(r => r.status === 'fulfilled').length

    if (failedCount > 0) {
      ElMessage.error(t('products.batchDeletePartialSuccess', { success: successCount, failed: failedCount }))
    } else {
      ElMessage.success(t('products.batchDeleteSuccess', { count: successCount }))
    }

    selectedProducts.value = []
    loadProducts()
  } catch (error) {
    if ((error as Error)?.message !== 'cancel') {
      handleError(error, t('products.batchDeleteFailed'))
    }
  }
}

const handleBatchPushToMarket = () => {
  if (selectedProducts.value.length === 0) {
    ElMessage.warning(t('products.selectAtLeastOneProduct'))
    return
  }

  // Reset form
  pushToMarketForm.selectedMarkets = []
  pushToMarketForm.price = 0
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

    let successCount = 0
    let failCount = 0

    for (const product of selectedProducts.value) {
      try {
        const result = await pushToMarket(product.id, {
          market_ids: pushToMarketForm.selectedMarkets,
          prices
        })
        successCount += result.success?.length || 0
        failCount += result.failed?.length || 0
      } catch (error) {
        failCount += pushToMarketForm.selectedMarkets.length
      }
    }

    pushToMarketDialogVisible.value = false
    ElMessage.success(t('products.pushToMarketSuccess', { success: successCount, failed: failCount }))
    loadProducts()
  } catch (error) {
    handleError(error, t('products.pushToMarketFailed'))
  } finally {
    pushToMarketLoading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        await createProduct({
          name: productForm.name,
          description: productForm.description,
          price: String(productForm.price),
          currency: productForm.currency || 'USD',
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
        dialogVisible.value = false
        ElMessage.success(t('products.addSuccess'))
        loadProducts()
      } catch (error) {
        handleError(error, t('products.addFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

const handleImageChange = async (file: UploadFile) => {
  try {
    if (!file.raw) {
      handleError(new Error('No file to upload'), t('products.imageUploadFailed'))
      return
    }
    const response = await uploadImage(file.raw, 'product')
    productForm.image = response.url
  } catch (error) {
    handleError(error, t('products.imageUploadFailed'))
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadProducts()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadProducts()
}

const handleMarketChange = () => {
  currentPage.value = 1
  loadProducts()
}

const loadMarkets = async () => {
  try {
    const response = await getMarkets()
    markets.value = response.list || []
  } catch (error) {
    handleError(error, t('products.loadMarketsFailed'))
  }
}

const fetchCategories = async () => {
  try {
    const response = await getCategoryTree()
    categories.value = response || []
  } catch (error) {
    handleError(error, t('products.loadCategoriesFailed'))
  }
}

const loadProducts = async () => {
  loading.value = true
  try {
    const params: ListProductsParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (searchQuery.value) {
      params.name = searchQuery.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterCategory.value) {
      params.category_id = filterCategory.value
    }
    if (selectedMarket.value) {
      params.market_id = selectedMarket.value
    }

    const response = await getProductList(params)
    productList.value = response.list || []
    total.value = response.total || 0
  } catch (error) {
    handleError(error, t('products.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMarkets()
  fetchCategories()
  loadProducts()
})
</script>

<style scoped>
.products-page {
  padding: 0;
}

/* Market Filter Bar */
.market-filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.market-filter-bar {
  display: flex;
  justify-content: flex-start;
}

.market-filter-bar :deep(.el-radio-button__inner) {
  border-radius: 10px !important;
  border: 1px solid #E5E7EB;
  padding: 8px 18px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.market-filter-bar :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  border-color: #6366F1;
  box-shadow: 0 4px 10px -2px rgba(99, 102, 241, 0.3);
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Bulk Actions */
.bulk-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 12px;
  margin-bottom: 16px;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.selected-count {
  font-size: 14px;
  color: #6366F1;
  font-weight: 600;
}

/* Table Card */
.table-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Product Cell */
.product-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.product-thumb {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid #E5E7EB;
  transition: transform 0.2s ease;
}

.product-thumb:hover {
  transform: scale(1.05);
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.product-details {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-size: 14px;
  line-height: 1.4;
}

.product-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0 0 8px 0;
  font-family: 'Fira Code', monospace;
}

.product-tags {
  display: flex;
  gap: 8px;
}

/* Price Cell */
.price-cell {
  text-align: right;
}

.sale-price {
  font-size: 16px;
  font-weight: 700;
  color: #EF4444;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.original-price {
  font-size: 12px;
  color: #9CA3AF;
  text-decoration: line-through;
  margin: 4px 0 0 0;
}

/* Market Tags */
.market-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  justify-content: center;
}

.market-tag {
  font-size: 11px;
  border-radius: 6px;
}

.no-markets {
  color: #9CA3AF;
  font-size: 12px;
}

.price-note {
  font-size: 12px;
  color: #F59E0B;
  margin-top: 8px;
  line-height: 1.4;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Dialog Styling */
:deep(.el-dialog) {
  border-radius: 16px;
}

/* Product Form Tabs */
.product-form-tabs {
  margin-bottom: 20px;
}

.product-form-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.product-form-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
}

.product-form-tabs :deep(.el-tabs__item) {
  font-weight: 500;
  color: #6B7280;
  transition: all 0.2s ease;
}

.product-form-tabs :deep(.el-tabs__item:hover) {
  color: #6366F1;
}

.product-form-tabs :deep(.el-tabs__item.is-active) {
  color: #6366F1;
  font-weight: 600;
}

.product-form-tabs :deep(.el-tabs__active-bar) {
  height: 3px;
  border-radius: 3px;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Upload */
.avatar-uploader {
  border: 2px dashed #E5E7EB;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.2s ease;
  width: 178px;
  height: 178px;
  background: #F9FAFB;
}

.avatar-uploader:hover {
  border-color: #6366F1;
  background: #F5F3FF;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #9CA3AF;
  width: 178px;
  height: 178px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: cover;
}

/* Switch */
:deep(.el-switch.is-checked .el-switch__core) {
  background-color: #10B981;
}

/* Tags */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: #10B981;
}

:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: #F59E0B;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .product-cell {
    flex-direction: column;
    align-items: flex-start;
  }

  .product-thumb {
    width: 60px;
    height: 60px;
  }

  .market-filter-card,
  .filter-card,
  .table-card {
    border-radius: 14px;
  }
}

/* Preview Dialog */
.preview-content {
  padding: 10px 0;
}

.preview-image {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 12px;
  overflow: hidden;
  background: #F9FAFB;
}

.preview-main-image {
  width: 100%;
  height: 100%;
}

.preview-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
  font-size: 48px;
}

.preview-info {
  padding: 10px 0;
}

.preview-title {
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.preview-sku {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 12px 0;
  font-family: 'Fira Code', monospace;
}

.preview-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.preview-price {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.price-label {
  font-size: 14px;
  color: #6B7280;
}

.price-value {
  font-size: 24px;
  font-weight: 700;
  color: #EF4444;
}

.preview-detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.detail-label {
  font-size: 14px;
  color: #6B7280;
  min-width: 80px;
}

.preview-description-row {
  margin-top: 20px;
}

.preview-description {
  background: #F9FAFB;
  border-radius: 8px;
  padding: 12px;
}

.description-content {
  margin: 8px 0 0 0;
  font-size: 14px;
  color: #374151;
  line-height: 1.6;
  white-space: pre-wrap;
}

.preview-markets {
  margin-top: 16px;
}

.preview-markets .market-tags {
  justify-content: flex-start;
  margin-top: 8px;
}

.preview-meta-row {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #E5E7EB;
}

.preview-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6B7280;
}
</style>
