<template>
  <div class="promotions-page">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <StatsCard color="primary">
          <template #icon><Ticket /></template>
          <template #value>{{ stats.totalCoupons }}</template>
          <template #label>{{ $t('promotions.totalCoupons') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="success">
          <template #icon><CircleCheck /></template>
          <template #value>{{ stats.activeCoupons }}</template>
          <template #label>{{ $t('promotions.activeCoupons') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="info">
          <template #icon><User /></template>
          <template #value>{{ stats.totalUsed }}</template>
          <template #label>{{ $t('promotions.totalUsed') }}</template>
        </StatsCard>
      </el-col>
      <el-col :xs="12" :sm="6">
        <StatsCard color="warning">
          <template #icon><Present /></template>
          <template #value>{{ stats.totalPromotions }}</template>
          <template #label>{{ $t('promotions.totalPromotions') }}</template>
        </StatsCard>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-card class="tabs-card" shadow="never">
      <el-tabs v-model="activeTab" class="promotion-tabs" @tab-change="handleTabChange">
        <!-- Coupons Tab -->
        <el-tab-pane :label="$t('promotions.couponManagement')" name="coupon">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="couponParams.name"
                :placeholder="$t('promotions.searchCouponName')"
                clearable
                class="search-input"
                @clear="loadCoupons"
                @keyup.enter="loadCoupons"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="couponParams.status" :placeholder="$t('promotions.status')" clearable class="filter-select" @change="loadCoupons">
                <el-option :label="$t('promotions.allStatus')" value="" />
                <el-option :label="$t('promotions.unactivated')" value="inactive" />
                <el-option :label="$t('promotions.activatedStatus')" value="active" />
              </el-select>
              <el-select v-model="couponParams.type" :placeholder="$t('promotions.type')" clearable class="filter-select" @change="loadCoupons">
                <el-option :label="$t('promotions.allType')" value="" />
                <el-option :label="$t('promotions.fixedAmount')" value="fixed_amount" />
                <el-option :label="$t('promotions.percentage')" value="percentage" />
              </el-select>
            </div>
            <el-button type="primary" @click="handleAddCoupon">
              <el-icon><Plus /></el-icon>{{ $t('promotions.createCouponButton') }}
            </el-button>
            <el-button @click="showBatchGenerateDialog">
              <el-icon><DocumentCopy /></el-icon>{{ $t('promotions.batchGenerate') }}
            </el-button>
          </div>

          <el-table :data="couponList" v-loading="couponLoading" stripe>
            <el-table-column :label="$t('promotions.couponInfo')" min-width="250">
              <template #default="{ row }">
                <div class="coupon-cell">
                  <div class="coupon-icon" :class="row.type">
                    <el-icon size="24"><Ticket /></el-icon>
                  </div>
                  <div class="coupon-details">
                    <p class="coupon-name">{{ row.name }}</p>
                    <p class="coupon-code">{{ $t('promotions.couponCode') }} {{ row.code }}</p>
                    <div class="coupon-tags">
                      <el-tag size="small" :type="row.type === 'fixed_amount' ? 'success' : 'warning'">
                        {{ row.type === 'fixed_amount' ? $t('promotions.fixedAmountLabel') : $t('promotions.percentageLabel') }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.discountStrength')" width="150" align="center">
              <template #default="{ row }">
                <div class="discount-value">
                  <span v-if="row.type === 'fixed_amount'">¥{{ row.discount_value }}</span>
                  <span v-else>{{ row.discount_value }}%</span>
                </div>
                <div class="min-order" v-if="row.min_order_amount && parseFloat(row.min_order_amount) > 0">
                  {{ $t('promotions.minOrderValue', { min: row.min_order_amount }) }}
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.usageStats')" width="180" align="center">
              <template #default="{ row }">
                <div class="usage-stats">
                  <el-progress
                    :percentage="getUsagePercentage(row)"
                    :status="getProgressStatus(row)"
                    :stroke-width="8"
                  />
                  <p class="usage-text">{{ $t('promotions.usedOfTotal', { used: row.used_count, total: row.usage_limit > 0 ? row.usage_limit : '∞' }) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.validity')" width="200">
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">{{ $t('promotions.toText') }}</p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.statusColumn')" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getCouponStatusType(row.status)" effect="light" size="small">
                  {{ getCouponStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.actionsColumn')" width="240" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditCoupon(row)">
                  {{ $t('promotions.edit') }}
                </el-button>
                <el-button type="primary" link size="small" @click="handleCouponUsage(row)">
                  {{ $t('promotions.data') }}
                </el-button>
                <el-button type="success" link size="small" @click="handleIssueToUser(row)">
                  {{ $t('promotions.issueToUser') }}
                </el-button>
                <el-button type="danger" link size="small" @click="handleDeleteCoupon(row)">
                  {{ $t('promotions.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <TablePagination
            v-model:current-page="couponParams.page"
            v-model:page-size="couponParams.page_size"
            :total="couponTotal"
            @change="handleCouponPageChange"
          />
        </el-tab-pane>

        <!-- Promotions Tab -->
        <el-tab-pane :label="$t('promotions.promotionActivities')" name="promotion">
          <div class="tab-header">
            <div class="tab-filters">
              <el-input
                v-model="promotionParams.name"
                :placeholder="$t('promotions.searchPromotionPlaceholder')"
                clearable
                class="search-input"
                @clear="loadPromotions"
                @keyup.enter="loadPromotions"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="promotionParams.status" :placeholder="$t('promotions.promotionStatus')" clearable class="filter-select" @change="loadPromotions">
                <el-option :label="$t('promotions.allStatus')" value="" />
                <el-option :label="$t('promotions.pending')" value="pending" />
                <el-option :label="$t('promotions.activeStatus')" value="active" />
                <el-option :label="$t('promotions.paused')" value="paused" />
                <el-option :label="$t('promotions.ended')" value="ended" />
              </el-select>
              <el-select v-model="promotionParams.type" :placeholder="$t('promotions.promotionType')" clearable class="filter-select" @change="loadPromotions">
                <el-option :label="$t('promotions.allType')" value="" />
                <el-option :label="$t('promotions.discount')" value="discount" />
                <el-option :label="$t('promotions.flashSale')" value="flash_sale" />
                <el-option :label="$t('promotions.bundle')" value="bundle" />
                <el-option :label="$t('promotions.buyXGetY')" value="buy_x_get_y" />
              </el-select>
            </div>
            <el-button type="primary" @click="handleAddPromotion">
              <el-icon><Plus /></el-icon>{{ $t('promotions.createPromotion') }}
            </el-button>
          </div>

          <el-table :data="promotionList" v-loading="promotionLoading" stripe>
            <el-table-column :label="$t('promotions.promotionInfo')" min-width="250">
              <template #default="{ row }">
                <div class="promo-cell">
                  <div class="promo-icon" :class="row.type">
                    <el-icon size="24"><Present /></el-icon>
                  </div>
                  <div class="promo-details">
                    <p class="promo-name">{{ row.name }}</p>
                    <p class="promo-desc">{{ row.description || $t('promotions.noDescription') }}</p>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.promotionType')" width="120" align="center">
              <template #default="{ row }">
                <el-tag size="small" :type="row.type === 'full_reduce' ? 'danger' : 'primary'">
                  {{ getPromotionTypeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.discountContent')" width="150" align="center">
              <template #default="{ row }">
                <div class="discount-value">
                  <span v-if="row.discount_type === 'fixed_amount'">¥{{ row.discount_value }}</span>
                  <span v-else-if="row.discount_type === 'percentage'">{{ row.discount_value }}%</span>
                  <span v-else>-</span>
                </div>
                <div class="min-order" v-if="row.min_order_amount && parseFloat(row.min_order_amount) > 0">
                  {{ $t('promotions.minOrderValue', { min: row.min_order_amount }) }}
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.validity')" width="200">
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">{{ $t('promotions.toText') }}</p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.statusColumn')" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getPromoStatusType(row.status)" effect="light" size="small">
                  {{ getPromoStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('promotions.actionsColumn')" width="200" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditPromotion(row)">
                  {{ $t('promotions.edit') }}
                </el-button>
                <el-button
                  v-if="row.status === 'pending' || row.status === 'paused'"
                  type="success"
                  link
                  size="small"
                  @click="handleActivatePromotion(row)"
                >
                  {{ $t('promotions.activate') }}
                </el-button>
                <el-button
                  v-if="row.status === 'active'"
                  type="warning"
                  link
                  size="small"
                  @click="handleDeactivatePromotion(row)"
                >
                  {{ $t('promotions.deactivate') }}
                </el-button>
                <el-button
                  v-if="row.status !== 'active'"
                  type="danger"
                  link
                  size="small"
                  @click="handleDeletePromotion(row)"
                >
                  {{ $t('promotions.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <TablePagination
            v-model:current-page="promotionParams.page"
            v-model:page-size="promotionParams.page_size"
            :total="promotionTotal"
            @change="handlePromotionPageChange"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Coupon Dialog -->
    <el-dialog v-model="couponDialogVisible" :title="isEditCoupon ? $t('promotions.editCoupon') : $t('promotions.createCoupon')" width="600px" destroy-on-close>
      <el-form :model="couponForm" label-width="100px" :rules="couponRules" ref="couponFormRef">
        <el-form-item :label="$t('promotions.couponName')" prop="name">
          <el-input v-model="couponForm.name" :placeholder="$t('promotions.enterCouponName')" maxlength="100" />
        </el-form-item>
        <el-form-item :label="$t('promotions.couponCode')" prop="code">
          <el-input v-model="couponForm.code" :placeholder="$t('promotions.enterCouponCode')" maxlength="50">
            <template #append>
              <el-button @click="generateCouponCode">{{ $t('promotions.generate') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('promotions.promotionDescription')">
          <el-input v-model="couponForm.description" type="textarea" :rows="2" :placeholder="$t('promotions.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('promotions.couponType')" prop="type">
          <el-radio-group v-model="couponForm.type">
            <el-radio label="fixed_amount">{{ $t('promotions.fixedAmountType') }}</el-radio>
            <el-radio label="percentage">{{ $t('promotions.percentageType') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="couponForm.type === 'fixed_amount' ? $t('promotions.discountAmount') : $t('promotions.discountRatio')" prop="discount_value">
          <el-input-number
            v-model="couponForm.discount_value_num"
            :min="0"
            :max="couponForm.type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.minOrderAmount')">
          <el-input-number v-model="couponForm.min_order_amount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="couponForm.type === 'percentage'" :label="$t('promotions.maxDiscount')">
          <el-input-number v-model="couponForm.max_discount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.usageLimit')">
          <el-input-number v-model="couponForm.usage_limit" :min="0" style="width: 100%" />
          <div class="form-tip">{{ $t('promotions.zeroUnlimited') }}</div>
        </el-form-item>
        <el-form-item :label="$t('promotions.perUserLimit')">
          <el-input-number v-model="couponForm.per_user_limit" :min="0" style="width: 100%" />
          <div class="form-tip">{{ $t('promotions.zeroUnlimited') }}</div>
        </el-form-item>
        <el-form-item :label="$t('promotions.validityPeriod')" prop="dateRange">
          <el-date-picker
            v-model="couponForm.dateRange"
            type="datetimerange"
            :start-placeholder="$t('promotions.startPlaceholder')"
            :end-placeholder="$t('promotions.endPlaceholder')"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="couponDialogVisible = false">{{ $t('promotions.cancel') }}</el-button>
        <el-button type="primary" @click="saveCoupon" :loading="saveLoading">{{ $t('promotions.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Promotion Dialog -->
    <el-dialog v-model="promotionDialogVisible" :title="isEditPromotion ? $t('promotions.editPromotion') : $t('promotions.createPromotion')" width="700px" destroy-on-close>
      <el-form :model="promotionForm" label-width="100px" :rules="promotionRules" ref="promotionFormRef">
        <el-form-item :label="$t('promotions.promotionName')" prop="name">
          <el-input v-model="promotionForm.name" :placeholder="$t('promotions.enterPromotionName')" maxlength="100" />
        </el-form-item>
        <el-form-item :label="$t('promotions.promotionDescription')">
          <el-input v-model="promotionForm.description" type="textarea" :rows="2" :placeholder="$t('promotions.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('promotions.promotionTypeSelect')" prop="type">
          <el-radio-group v-model="promotionForm.type">
            <el-radio label="discount">{{ $t('promotions.discount') }}</el-radio>
            <el-radio label="flash_sale">{{ $t('promotions.flashSale') }}</el-radio>
            <el-radio label="bundle">{{ $t('promotions.bundle') }}</el-radio>
            <el-radio label="buy_x_get_y">{{ $t('promotions.buyXGetY') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('promotions.discountType')" prop="discount_type">
          <el-radio-group v-model="promotionForm.discount_type">
            <el-radio label="fixed_amount">{{ $t('promotions.fixedAmount') }}</el-radio>
            <el-radio label="percentage">{{ $t('promotions.percentage') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="promotionForm.discount_type === 'fixed_amount' ? $t('promotions.discountAmountLabel') : $t('promotions.discountRatioLabel')" prop="discount_value">
          <el-input-number
            v-model="promotionForm.discount_value_num"
            :min="0"
            :max="promotionForm.discount_type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.lowestConsume')">
          <el-input-number v-model="promotionForm.min_order_amount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="promotionForm.discount_type === 'percentage'" :label="$t('promotions.maxDiscountAmount')">
          <el-input-number v-model="promotionForm.max_discount_num" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.validityPeriod')" prop="dateRange">
          <el-date-picker
            v-model="promotionForm.dateRange"
            type="datetimerange"
            :start-placeholder="$t('promotions.startPlaceholder')"
            :end-placeholder="$t('promotions.endPlaceholder')"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="promotionDialogVisible = false">{{ $t('promotions.cancel') }}</el-button>
        <el-button type="primary" @click="savePromotion" :loading="saveLoading">{{ $t('promotions.save') }}</el-button>
      </template>
    </el-dialog>

    <!-- Coupon Usage Dialog -->
    <el-dialog v-model="usageDialogVisible" :title="$t('promotions.couponUsageRecords')" width="800px">
      <el-table :data="usageList" v-loading="usageLoading" stripe>
        <el-table-column prop="user_id" :label="$t('promotions.userID')" width="100" />
        <el-table-column prop="order_id" :label="$t('promotions.orderID')" width="120" />
        <el-table-column prop="discount_amount" :label="$t('promotions.discountAmountColumn')" width="120">
          <template #default="{ row }">
            ¥{{ row.discount_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="used_at" :label="$t('promotions.usedTime')" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.used_at) }}
          </template>
        </el-table-column>
      </el-table>
      <TablePagination
        v-model:current-page="usageParams.page"
        v-model:page-size="usageParams.page_size"
        :total="usageTotal"
        @change="handleUsagePageChange"
      />
    </el-dialog>

    <!-- Batch Generate Coupons Dialog -->
    <el-dialog v-model="batchGenerateDialogVisible" :title="$t('promotions.batchGenerateCoupons')" width="600px" destroy-on-close>
      <el-form :model="batchGenerateForm" label-width="140px">
        <el-form-item :label="$t('promotions.generateCouponPrefix')">
          <el-input v-model="batchGenerateForm.prefix" :placeholder="$t('promotions.generateCouponPrefixPlaceholder')" maxlength="10" />
        </el-form-item>
        <el-form-item :label="$t('promotions.generateCouponQuantity')">
          <el-input-number v-model="batchGenerateForm.quantity" :min="1" :max="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.generateCouponLength')">
          <el-input-number v-model="batchGenerateForm.length" :min="4" :max="20" style="width: 100%" />
        </el-form-item>
        <el-divider>{{ $t('promotions.generateCouponConfig') }}</el-divider>
        <el-form-item :label="$t('promotions.couponName')">
          <el-input v-model="batchGenerateForm.couponName" :placeholder="$t('promotions.enterCouponName')" />
        </el-form-item>
        <el-form-item :label="$t('promotions.couponType')">
          <el-radio-group v-model="batchGenerateForm.couponType">
            <el-radio label="fixed_amount">{{ $t('promotions.fixedAmountType') }}</el-radio>
            <el-radio label="percentage">{{ $t('promotions.percentageType') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="batchGenerateForm.couponType === 'fixed_amount' ? $t('promotions.discountAmount') : $t('promotions.discountRatio')">
          <el-input-number v-model="batchGenerateForm.discountValue" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.minOrderAmount')">
          <el-input-number v-model="batchGenerateForm.minOrderAmount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.usageLimit')">
          <el-input-number v-model="batchGenerateForm.usageLimit" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('promotions.validityPeriod')">
          <el-date-picker
            v-model="batchGenerateForm.dateRange"
            type="datetimerange"
            :start-placeholder="$t('promotions.startPlaceholder')"
            :end-placeholder="$t('promotions.endPlaceholder')"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchGenerateDialogVisible = false">{{ $t('promotions.cancel') }}</el-button>
        <el-button type="primary" @click="handleBatchGenerate" :loading="batchGenerating">
          {{ $t('promotions.batchGenerate') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Generated Codes Dialog -->
    <el-dialog v-model="generatedCodesDialogVisible" :title="$t('promotions.generatedCodes')" width="700px">
      <div class="generated-codes-info">
        {{ $t('promotions.generatedCodesCount', { count: generatedCodes.length }) }}
        <el-button type="primary" link @click="handleCopyAllCodes">{{ $t('promotions.copyAllCodes') }}</el-button>
      </div>
      <div class="generated-codes-list">
        <el-tag v-for="code in generatedCodes" :key="code" class="code-tag" type="info">
          {{ code }}
        </el-tag>
      </div>
      <template #footer>
        <el-button type="primary" @click="generatedCodesDialogVisible = false">{{ $t('common.close') }}</el-button>
      </template>
    </el-dialog>

    <!-- Issue Coupon to User Dialog -->
    <el-dialog v-model="issueDialogVisible" :title="$t('promotions.issueCoupon')" width="500px" destroy-on-close>
      <el-form :model="issueForm" label-width="100px">
        <el-form-item :label="$t('promotions.currentCoupon')">
          <el-input :value="issueForm.coupon_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('promotions.userId')" required>
          <el-input-number
            v-model="issueForm.user_id"
            :min="1"
            :placeholder="$t('promotions.enterUserId')"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="issueDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleConfirmIssue" :loading="issueLoading">
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Ticket, CircleCheck, User, Search, Plus, Present, DocumentCopy } from '@element-plus/icons-vue'
import {
  getCouponList, createCoupon, updateCoupon, deleteCoupon,
  getPromotionList, createPromotion, updatePromotion, deletePromotion,
  activatePromotion, deactivatePromotion, getCouponUsage,
  generateCouponCodes, issueUserCoupon,
  type Coupon, type Promotion, type CouponUsage,
  type CouponStatus, type CouponType, type PromotionStatus, type PromotionType,
  type GenerateCouponCodesRequest
} from '@/api/promotion'
import TablePagination from '@/components/common/TablePagination.vue'
import StatsCard from '@/components/common/StatsCard.vue'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { handleError } = useErrorHandler()

// Tab state
const activeTab = ref('coupon')

// Stats
const stats = ref({
  totalCoupons: 0,
  activeCoupons: 0,
  totalUsed: 0,
  totalPromotions: 0
})

// Coupon list state
const couponList = ref<Coupon[]>([])
const couponLoading = ref(false)
const couponTotal = ref(0)
const couponParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as CouponStatus | undefined,
  type: undefined as CouponType | undefined
})

// Promotion list state
const promotionList = ref<Promotion[]>([])
const promotionLoading = ref(false)
const promotionTotal = ref(0)
const promotionParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as PromotionStatus | undefined,
  type: undefined as PromotionType | undefined
})

// Usage dialog state
const usageDialogVisible = ref(false)
const usageList = ref<CouponUsage[]>([])
const usageLoading = ref(false)
const usageTotal = ref(0)
const usageParams = reactive({
  id: 0,
  page: 1,
  page_size: 10
})

// Batch generate dialog state
const batchGenerateDialogVisible = ref(false)
const batchGenerating = ref(false)
const batchGenerateForm = reactive({
  prefix: '',
  quantity: 10,
  length: 8,
  couponName: '',
  couponType: 'fixed_amount' as 'fixed_amount' | 'percentage',
  discountValue: 0,
  minOrderAmount: 0,
  usageLimit: 100,
  dateRange: [] as string[]
})

// Generated codes dialog state
const generatedCodesDialogVisible = ref(false)
const generatedCodes = ref<string[]>([])

// Issue coupon dialog state
const issueDialogVisible = ref(false)
const issueLoading = ref(false)
const issueForm = reactive({
  user_id: 0 as number,
  coupon_id: 0 as number,
  coupon_name: ''
})

// Coupon form state
const couponDialogVisible = ref(false)
const isEditCoupon = ref(false)
const saveLoading = ref(false)
const couponFormRef = ref()
const couponForm = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  type: 'fixed_amount' as 'fixed_amount' | 'percentage',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  usage_limit: 0,
  per_user_limit: 0,
  dateRange: [] as string[]
})

const couponRules = {
  name: [{ required: true, message: t('promotions.enterCouponName'), trigger: 'blur' }],
  code: [{ required: true, message: t('promotions.enterCouponCode'), trigger: 'blur' }],
  type: [{ required: true, message: t('promotions.selectCouponType'), trigger: 'change' }],
  discount_value: [{ required: true, message: t('promotions.enterDiscountValue'), trigger: 'blur' }],
  dateRange: [{ required: true, message: t('promotions.selectValidityPeriod'), trigger: 'change' }]
}

// Promotion form state
const promotionDialogVisible = ref(false)
const isEditPromotion = ref(false)
const promotionFormRef = ref()
const promotionForm = reactive({
  id: 0,
  name: '',
  description: '',
  type: 'discount' as 'discount' | 'flash_sale' | 'bundle' | 'buy_x_get_y',
  discount_type: 'fixed_amount' as 'fixed_amount' | 'percentage' | 'buy_x_get_y',
  discount_value_num: 0,
  min_order_amount_num: 0,
  max_discount_num: 0,
  dateRange: [] as string[]
})

const promotionRules = {
  name: [{ required: true, message: t('promotions.enterPromotionName'), trigger: 'blur' }],
  type: [{ required: true, message: t('promotions.selectPromotionType'), trigger: 'change' }],
  discount_type: [{ required: true, message: t('promotions.selectDiscountType'), trigger: 'change' }],
  dateRange: [{ required: true, message: t('promotions.selectValidityPeriod'), trigger: 'change' }]
}

// Load functions
const loadCoupons = async () => {
  couponLoading.value = true
  try {
    const res = await getCouponList(couponParams)
    couponList.value = res.list || []
    couponTotal.value = res.total || 0
    // Update stats
    stats.value.totalCoupons = couponTotal.value
    stats.value.activeCoupons = couponList.value.filter(c => c.status === 'active').length
    stats.value.totalUsed = couponList.value.reduce((sum, c) => sum + c.used_count, 0)
  } catch (error) {
    handleError(error, t('promotions.loadCouponsFailed'))
  } finally {
    couponLoading.value = false
  }
}

const loadPromotions = async () => {
  promotionLoading.value = true
  try {
    const res = await getPromotionList(promotionParams)
    promotionList.value = res.list || []
    promotionTotal.value = res.total || 0
    stats.value.totalPromotions = promotionTotal.value
  } catch (error) {
    handleError(error, t('promotions.loadPromotionsFailed'))
  } finally {
    promotionLoading.value = false
  }
}

const loadCouponUsage = async () => {
  usageLoading.value = true
  try {
    const res = await getCouponUsage(usageParams.id, { page: usageParams.page, page_size: usageParams.page_size })
    usageList.value = res.list || []
    usageTotal.value = res.total || 0
  } catch (error) {
    handleError(error, t('promotions.loadCouponUsageFailed'))
  } finally {
    usageLoading.value = false
  }
}

// Helper functions
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getUsagePercentage = (row: Coupon) => {
  if (row.usage_limit === 0) return 0
  return Math.round((row.used_count / row.usage_limit) * 100)
}

const getProgressStatus = (row: Coupon) => {
  const percentage = getUsagePercentage(row)
  if (percentage >= 90) return 'exception'
  if (percentage >= 70) return 'warning'
  return ''
}

const getCouponStatusType = (status: string) => {
  const types: Record<string, string> = {
    'inactive': 'info',
    'active': 'success',
    'expired': 'warning',
    'depleted': 'danger'
  }
  return types[status] || 'info'
}

const getCouponStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'inactive': t('promotions.unactivated'),
    'active': t('promotions.activatedStatus'),
    'expired': t('promotions.expiredStatus'),
    'depleted': t('promotions.depletedStatus')
  }
  return texts[status] || status
}

const getPromoStatusType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'paused': 'warning',
    'pending': 'info',
    'ended': 'info'
  }
  return types[status] || 'info'
}

const getPromoStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': t('promotions.activeStatus'),
    'paused': t('promotions.paused'),
    'pending': t('promotions.pending'),
    'ended': t('promotions.ended')
  }
  return texts[status] || status
}

const getPromotionTypeText = (type: string) => {
  const texts: Record<string, string> = {
    'discount': t('promotions.discount'),
    'flash_sale': t('promotions.flashSale'),
    'bundle': t('promotions.bundle'),
    'buy_x_get_y': t('promotions.buyXGetY')
  }
  return texts[type] || type
}

const generateCouponCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 10; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  couponForm.code = code
}

// Page change handlers
const handleCouponPageChange = () => {
  loadCoupons()
}

const handlePromotionPageChange = () => {
  loadPromotions()
}

const handleUsagePageChange = () => {
  loadCouponUsage()
}

const handleTabChange = (tab: string) => {
  if (tab === 'coupon') {
    loadCoupons()
  } else if (tab === 'promotion') {
    loadPromotions()
  }
}

// Coupon actions
const handleAddCoupon = () => {
  isEditCoupon.value = false
  Object.assign(couponForm, {
    id: 0,
    name: '',
    code: '',
    description: '',
    type: 'fixed_amount',
    discount_value_num: 0,
    min_order_amount_num: 0,
    max_discount_num: 0,
    usage_limit: 100,
    per_user_limit: 0,
    dateRange: []
  })
  couponDialogVisible.value = true
}

const handleEditCoupon = (row: Coupon) => {
  isEditCoupon.value = true
  Object.assign(couponForm, {
    id: row.id,
    name: row.name,
    code: row.code,
    description: row.description,
    type: row.type,
    discount_value_num: parseFloat(row.discount_value) || 0,
    min_order_amount_num: parseFloat(row.min_order_amount) || 0,
    max_discount_num: parseFloat(row.max_discount) || 0,
    usage_limit: row.usage_limit,
    per_user_limit: row.per_user_limit,
    dateRange: [row.start_time, row.end_time]
  })
  couponDialogVisible.value = true
}

const handleCouponUsage = (row: Coupon) => {
  usageParams.id = row.id
  usageParams.page = 1
  usageDialogVisible.value = true
  loadCouponUsage()
}

const handleDeleteCoupon = async (row: Coupon) => {
  try {
    await ElMessageBox.confirm(t('promotions.deleteCouponConfirm', { name: row.name }), t('common.warning'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await deleteCoupon(row.id)
    ElMessage.success(t('promotions.deleteSuccess'))
    loadCoupons()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('promotions.deleteFailed'))
    }
  }
}

const handleIssueToUser = (row: Coupon) => {
  issueForm.user_id = 0
  issueForm.coupon_id = row.id
  issueForm.coupon_name = row.name
  issueDialogVisible.value = true
}

const handleConfirmIssue = async () => {
  if (!issueForm.user_id) {
    ElMessage.warning(t('promotions.enterUserId'))
    return
  }

  issueLoading.value = true
  try {
    await issueUserCoupon({
      user_id: issueForm.user_id,
      coupon_id: issueForm.coupon_id
    })
    ElMessage.success(t('promotions.issueSuccess'))
    issueDialogVisible.value = false
  } catch (error) {
    handleError(error, t('promotions.issueFailed'))
  } finally {
    issueLoading.value = false
  }
}

const saveCoupon = async () => {
  if (!couponFormRef.value) return

  await couponFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saveLoading.value = true
    try {
      const data = {
        name: couponForm.name,
        code: couponForm.code,
        description: couponForm.description,
        type: couponForm.type,
        discount_value: String(couponForm.discount_value_num),
        min_order_amount: String(couponForm.min_order_amount_num),
        max_discount: String(couponForm.max_discount_num),
        usage_limit: couponForm.usage_limit,
        per_user_limit: couponForm.per_user_limit,
        start_time: couponForm.dateRange[0],
        end_time: couponForm.dateRange[1]
      }

      if (isEditCoupon.value) {
        await updateCoupon({ id: couponForm.id, ...data })
      } else {
        await createCoupon(data)
      }

      ElMessage.success(isEditCoupon.value ? t('promotions.updateSuccess') : t('promotions.createSuccess'))
      couponDialogVisible.value = false
      loadCoupons()
    } catch (error) {
      handleError(error, isEditCoupon.value ? t('promotions.updateCouponFailed') : t('promotions.createCouponFailed'))
    } finally {
      saveLoading.value = false
    }
  })
}

// Batch generate coupon codes
const showBatchGenerateDialog = () => {
  batchGenerateForm.prefix = ''
  batchGenerateForm.quantity = 10
  batchGenerateForm.length = 8
  batchGenerateForm.couponName = ''
  batchGenerateForm.couponType = 'fixed_amount'
  batchGenerateForm.discountValue = 0
  batchGenerateForm.minOrderAmount = 0
  batchGenerateForm.usageLimit = 100
  batchGenerateForm.dateRange = []
  batchGenerateDialogVisible.value = true
}

const handleBatchGenerate = async () => {
  if (!batchGenerateForm.couponName) {
    ElMessage.warning(t('promotions.enterCouponName'))
    return
  }
  if (!batchGenerateForm.dateRange || batchGenerateForm.dateRange.length !== 2) {
    ElMessage.warning(t('promotions.selectValidityPeriod'))
    return
  }

  batchGenerating.value = true
  try {
    const couponConfig = {
      name: batchGenerateForm.couponName,
      type: batchGenerateForm.couponType,
      discount_value: String(batchGenerateForm.discountValue),
      min_order_amount: String(batchGenerateForm.minOrderAmount),
      usage_limit: batchGenerateForm.usageLimit,
      start_time: batchGenerateForm.dateRange[0],
      end_time: batchGenerateForm.dateRange[1]
    }

    const data: GenerateCouponCodesRequest = {
      prefix: batchGenerateForm.prefix || undefined,
      quantity: batchGenerateForm.quantity,
      length: batchGenerateForm.length,
      coupon_config: JSON.stringify(couponConfig)
    }

    const res = await generateCouponCodes(data)
    generatedCodes.value = res.codes || []
    batchGenerateDialogVisible.value = false
    generatedCodesDialogVisible.value = true
    ElMessage.success(t('promotions.generateSuccess'))
  } catch (error) {
    handleError(error, t('promotions.generateFailed'))
  } finally {
    batchGenerating.value = false
  }
}

const handleCopyAllCodes = async () => {
  const codesText = generatedCodes.value.join('\n')
  try {
    await navigator.clipboard.writeText(codesText)
    ElMessage.success(t('promotions.codesCopied'))
  } catch (error) {
    handleError(error)
  }
}

// Promotion actions
const handleAddPromotion = () => {
  isEditPromotion.value = false
  Object.assign(promotionForm, {
    id: 0,
    name: '',
    description: '',
    type: 'discount',
    discount_type: 'fixed_amount',
    discount_value_num: 0,
    min_order_amount_num: 0,
    max_discount_num: 0,
    dateRange: []
  })
  promotionDialogVisible.value = true
}

const handleEditPromotion = (row: Promotion) => {
  router.push(`/promotions/${row.id}`)
}

const handleActivatePromotion = async (row: Promotion) => {
  try {
    await activatePromotion(row.id)
    ElMessage.success(t('promotions.activateSuccess'))
    loadPromotions()
  } catch (error) {
    handleError(error, t('promotions.activatePromotionFailed'))
  }
}

const handleDeactivatePromotion = async (row: Promotion) => {
  try {
    await deactivatePromotion(row.id)
    ElMessage.success(t('promotions.deactivateSuccess'))
    loadPromotions()
  } catch (error) {
    handleError(error, t('promotions.deactivatePromotionFailed'))
  }
}

const handleDeletePromotion = async (row: Promotion) => {
  try {
    await ElMessageBox.confirm(t('promotions.deletePromotionConfirm', { name: row.name }), t('common.warning'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await deletePromotion(row.id)
    ElMessage.success(t('promotions.deleteSuccess'))
    loadPromotions()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('promotions.deleteFailed'))
    }
  }
}

const savePromotion = async () => {
  if (!promotionFormRef.value) return

  await promotionFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saveLoading.value = true
    try {
      const data = {
        name: promotionForm.name,
        description: promotionForm.description,
        type: promotionForm.type,
        discount_type: promotionForm.discount_type,
        discount_value: String(promotionForm.discount_value_num),
        min_order_amount: String(promotionForm.min_order_amount_num),
        max_discount: String(promotionForm.max_discount_num),
        start_time: promotionForm.dateRange[0],
        end_time: promotionForm.dateRange[1]
      }

      if (isEditPromotion.value) {
        await updatePromotion({ id: promotionForm.id, ...data })
      } else {
        await createPromotion(data)
      }

      ElMessage.success(isEditPromotion.value ? t('promotions.updateSuccess') : t('promotions.createSuccess'))
      promotionDialogVisible.value = false
      loadPromotions()
    } catch (error) {
      handleError(error, t('promotions.savePromotionFailed'))
    } finally {
      saveLoading.value = false
    }
  })
}

// Initialize
onMounted(() => {
  loadCoupons()
  loadPromotions()
})
</script>

<style scoped>
.promotions-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #6B7280;
}

.stat-icon.active {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.stat-icon.used {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  color: white;
}

.stat-icon.promotions {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Tabs Card */
.tabs-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Tabs */
:deep(.el-tabs__item.is-active) {
  color: #6366F1;
  font-weight: 600;
}

:deep(.el-tabs__active-bar) {
  background-color: #6366F1;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.tab-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 220px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

/* Coupon Cell */
.coupon-cell, .promo-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.coupon-icon, .promo-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F5F3FF;
  color: #6366F1;
}

.coupon-icon.fixed_amount, .promo-icon.discount {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.coupon-icon.percentage, .promo-icon.full_reduce {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.coupon-details, .promo-details {
  flex: 1;
  min-width: 0;
}

.coupon-name, .promo-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.coupon-code {
  font-size: 12px;
  color: #6B7280;
  margin: 0 0 8px 0;
  font-family: 'Fira Code', monospace;
}

.coupon-tags {
  display: flex;
  gap: 8px;
}

.promo-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Discount Value */
.discount-value {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.min-order {
  font-size: 12px;
  color: #6B7280;
  margin-top: 4px;
}

/* Usage Stats */
.usage-stats {
  text-align: center;
}

.usage-text {
  font-size: 12px;
  color: #6B7280;
  margin: 8px 0 0 0;
}

/* Validity Period */
.validity-period {
  text-align: center;
  font-size: 13px;
  color: #4B5563;
}

.to-text {
  margin: 2px 0;
  color: #9CA3AF;
  font-size: 12px;
}

/* Form */
.form-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

/* Generated Codes */
.generated-codes-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 16px;
  background: #F3F4F6;
  border-radius: 8px;
  font-size: 14px;
  color: #6B7280;
}

.generated-codes-list {
  max-height: 400px;
  overflow-y: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.code-tag {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  padding: 8px 12px;
}

/* Dialog */
:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Progress */
:deep(.el-progress-bar__inner) {
  background: linear-gradient(90deg, #6366F1 0%, #818CF8 100%);
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

:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: #6B7280;
}

:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: #EF4444;
}

:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: #6366F1;
}

/* Responsive */
@media (max-width: 768px) {
  .tab-header {
    flex-direction: column;
    align-items: stretch;
  }

  .tab-filters {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>