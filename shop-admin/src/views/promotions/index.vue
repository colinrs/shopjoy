<template>
  <div class="promotions-page">
    <!-- Stats Cards -->
    <el-row
      :gutter="16"
      class="stats-row"
    >
      <el-col
        :xs="12"
        :sm="6"
      >
        <StatsCard color="primary">
          <template #icon>
            <Ticket />
          </template>
          <template #value>
            {{ stats.totalCoupons }}
          </template>
          <template #label>
            {{ $t('promotions.totalCoupons') }}
          </template>
        </StatsCard>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <StatsCard color="success">
          <template #icon>
            <CircleCheck />
          </template>
          <template #value>
            {{ stats.activeCoupons }}
          </template>
          <template #label>
            {{ $t('promotions.activeCoupons') }}
          </template>
        </StatsCard>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <StatsCard color="info">
          <template #icon>
            <User />
          </template>
          <template #value>
            {{ stats.totalUsed }}
          </template>
          <template #label>
            {{ $t('promotions.totalUsed') }}
          </template>
        </StatsCard>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <StatsCard color="warning">
          <template #icon>
            <Present />
          </template>
          <template #value>
            {{ stats.totalPromotions }}
          </template>
          <template #label>
            {{ $t('promotions.totalPromotions') }}
          </template>
        </StatsCard>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-card
      class="tabs-card"
      shadow="never"
    >
      <el-tabs
        v-model="activeTab"
        class="promotion-tabs"
        @tab-change="handleTabChange"
      >
        <!-- Coupons Tab -->
        <el-tab-pane
          :label="$t('promotions.couponManagement')"
          name="coupon"
        >
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
              <el-select
                v-model="couponParams.status"
                :placeholder="$t('promotions.status')"
                clearable
                class="filter-select"
                @change="loadCoupons"
              >
                <el-option
                  :label="$t('promotions.allStatus')"
                  value=""
                />
                <el-option
                  :label="$t('promotions.pending')"
                  value="pending"
                />
                <el-option
                  :label="$t('promotions.activeStatus')"
                  value="active"
                />
                <el-option
                  :label="$t('promotions.paused')"
                  value="paused"
                />
                <el-option
                  :label="$t('promotions.ended')"
                  value="ended"
                />
                <el-option
                  :label="$t('promotions.expiredStatus')"
                  value="expired"
                />
              </el-select>
              <el-select
                v-model="couponParams.market_id"
                :placeholder="$t('promotions.market')"
                clearable
                class="filter-select"
                @change="loadCoupons"
              >
                <el-option
                  v-for="m in marketOptions"
                  :key="m.id"
                  :label="m.name"
                  :value="m.id"
                />
              </el-select>
            </div>
            <el-button
              type="primary"
              @click="handleAddCoupon"
            >
              <el-icon><Plus /></el-icon>{{ $t('promotions.createCouponButton') }}
            </el-button>
            <el-button @click="showBatchGenerateDialog">
              <el-icon><DocumentCopy /></el-icon>{{ $t('promotions.batchGenerate') }}
            </el-button>
          </div>

          <el-table
            v-loading="couponLoading"
            :data="couponList"
            stripe
          >
            <el-table-column
              :label="$t('promotions.couponInfo')"
              min-width="250"
            >
              <template #default="{ row }">
                <div class="coupon-cell">
                  <div
                    class="coupon-icon"
                    :class="row.type"
                  >
                    <el-icon size="24">
                      <Ticket />
                    </el-icon>
                  </div>
                  <div class="coupon-details">
                    <p class="coupon-name">
                      {{ row.name }}
                    </p>
                    <p
                      v-if="row.kind === 'coupon' && row.code"
                      class="coupon-code"
                    >
                      {{ $t('promotions.couponCode') }} {{ row.code }}
                    </p>
                    <div class="coupon-tags">
                      <el-tag
                        size="small"
                        :type="row.type === 'fixed_amount' ? 'success' : 'warning'"
                      >
                        {{ row.type === 'fixed_amount' ? $t('promotions.fixedAmountLabel') : $t('promotions.percentageLabel') }}
                      </el-tag>
                      <el-tag
                        v-if="(row.rules?.length || 0) > 0"
                        size="small"
                        type="info"
                      >
                        {{ $t('promotions.ruleCount', { count: row.rules?.length }) }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.discountStrength')"
              width="150"
              align="center"
            >
              <template #default="{ row }">
                <div class="discount-value">
                  <template v-if="row.rules && row.rules.length > 0">
                    <span
                      v-for="(rule, idx) in row.rules.slice(0, 1)"
                      :key="idx"
                    >
                      <span v-if="rule.action_type === 'fixed_amount'">¥{{ rule.action_value }}</span>
                      <span v-else-if="rule.action_type === 'percentage'">{{ rule.action_value }}%</span>
                      <span v-else>{{ rule.action_type }}</span>
                    </span>
                  </template>
                  <span v-else>-</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              v-if="false"
              :label="$t('promotions.usageStats')"
              width="180"
              align="center"
            >
              <template #default="{ row }">
                <div class="usage-stats">
                  <el-progress
                    :percentage="getUsagePercentage(row)"
                    :status="getProgressStatus(row)"
                    :stroke-width="8"
                  />
                  <p class="usage-text">
                    {{ $t('promotions.usedOfTotal', { used: row.used_count, total: row.usage_limit > 0 ? row.usage_limit : '∞' }) }}
                  </p>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.usageStats')"
              width="180"
              align="center"
            >
              <template #default="{ row }">
                <div class="usage-stats">
                  <el-progress
                    :percentage="getUsagePercentage(row)"
                    :status="getProgressStatus(row)"
                    :stroke-width="8"
                  />
                  <p class="usage-text">
                    {{ $t('promotions.usedOfTotal', {
                      used: row.used_count,
                      total: row.usage_limit > 0 ? row.usage_limit : (row.total_count && row.total_count > 0 ? row.total_count : '∞')
                    }) }}
                  </p>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.validity')"
              width="200"
            >
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">
                    {{ $t('promotions.toText') }}
                  </p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.statusColumn')"
              width="120"
              align="center"
            >
              <template #default="{ row }">
                <el-switch
                  v-if="row.status === 'active' || row.status === 'pending' || row.status === 'paused'"
                  :model-value="row.status === 'active'"
                  :loading="couponToggleLoading[row.id] === true"
                  :active-text="$t('promotions.activatedStatus')"
                  :inactive-text="$t('promotions.unactivated')"
                  inline-prompt
                  @change="(val: boolean) => handleToggleCouponStatus(row, val)"
                />
                <el-tag
                  v-else
                  :type="getPromoStatusType(row.status)"
                  effect="light"
                  size="small"
                >
                  {{ getPromoStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.actionsColumn')"
              width="240"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="handleEditCoupon(row)"
                >
                  {{ $t('promotions.edit') }}
                </el-button>
                <el-button
                  v-if="row.kind === 'coupon'"
                  type="primary"
                  link
                  size="small"
                  @click="handleCouponUsage(row)"
                >
                  {{ $t('promotions.data') }}
                </el-button>
                <el-button
                  v-if="row.kind === 'coupon'"
                  type="success"
                  link
                  size="small"
                  @click="handleIssueToUser(row)"
                >
                  {{ $t('promotions.issueToUser') }}
                </el-button>
                <el-button
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
            v-model:current-page="couponParams.page"
            v-model:page-size="couponParams.page_size"
            :total="couponTotal"
            @change="loadCoupons"
          />
        </el-tab-pane>

        <!-- Promotions Tab -->
        <el-tab-pane
          :label="$t('promotions.promotionActivities')"
          name="promotion"
        >
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
              <el-select
                v-model="promotionParams.status"
                :placeholder="$t('promotions.promotionStatus')"
                clearable
                class="filter-select"
                @change="loadPromotions"
              >
                <el-option
                  :label="$t('promotions.allStatus')"
                  value=""
                />
                <el-option
                  :label="$t('promotions.pending')"
                  value="pending"
                />
                <el-option
                  :label="$t('promotions.activeStatus')"
                  value="active"
                />
                <el-option
                  :label="$t('promotions.paused')"
                  value="paused"
                />
                <el-option
                  :label="$t('promotions.ended')"
                  value="ended"
                />
                <el-option
                  :label="$t('promotions.expiredStatus')"
                  value="expired"
                />
              </el-select>
              <el-select
                v-model="promotionParams.market_id"
                :placeholder="$t('promotions.market')"
                clearable
                class="filter-select"
                @change="loadPromotions"
              >
                <el-option
                  v-for="m in marketOptions"
                  :key="m.id"
                  :label="m.name"
                  :value="m.id"
                />
              </el-select>
              <el-select
                v-model="promotionParams.type"
                :placeholder="$t('promotions.promotionType')"
                clearable
                class="filter-select"
                @change="loadPromotions"
              >
                <el-option
                  :label="$t('promotions.allType')"
                  value=""
                />
                <el-option
                  :label="$t('promotions.discount')"
                  value="discount"
                />
                <el-option
                  :label="$t('promotions.flashSale')"
                  value="flash_sale"
                />
                <el-option
                  :label="$t('promotions.bundle')"
                  value="bundle"
                />
                <el-option
                  :label="$t('promotions.buyXGetY')"
                  value="buy_x_get_y"
                />
              </el-select>
            </div>
            <el-button
              type="primary"
              @click="handleAddPromotion"
            >
              <el-icon><Plus /></el-icon>{{ $t('promotions.createPromotion') }}
            </el-button>
          </div>

          <el-table
            v-loading="promotionLoading"
            :data="promotionList"
            stripe
          >
            <el-table-column
              :label="$t('promotions.promotionInfo')"
              min-width="250"
            >
              <template #default="{ row }">
                <div class="promo-cell">
                  <div
                    class="promo-icon"
                    :class="row.type"
                  >
                    <el-icon size="24">
                      <Present />
                    </el-icon>
                  </div>
                  <div class="promo-details">
                    <p class="promo-name">
                      {{ row.name }}
                    </p>
                    <p class="promo-desc">
                      {{ row.description || $t('promotions.noDescription') }}
                    </p>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.promotionType')"
              width="120"
              align="center"
            >
              <template #default="{ row }">
                <el-tag
                  size="small"
                  :type="row.type === 'full_reduce' ? 'danger' : 'primary'"
                >
                  {{ getPromotionTypeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.discountContent')"
              width="180"
              align="center"
            >
              <template #default="{ row }">
                <div class="discount-value">
                  <template v-if="row.rules && row.rules.length > 0">
                    <span
                      v-for="(rule, idx) in row.rules.slice(0, 1)"
                      :key="idx"
                    >
                      <span v-if="rule.action_type === 'fixed_amount'">¥{{ rule.action_value }}</span>
                      <span v-else-if="rule.action_type === 'percentage'">{{ rule.action_value }}%</span>
                      <span v-else-if="rule.action_type === 'free_shipping'">{{ $t('promotions.freeShipping') }}</span>
                      <span v-else>-</span>
                    </span>
                  </template>
                  <span v-else>-</span>
                </div>
                <div
                  v-if="row.rules && row.rules[0]?.max_discount && parseFloat(row.rules[0].max_discount) > 0"
                  class="min-order"
                >
                  {{ $t('promotions.maxDiscountUpTo', { value: row.rules[0].max_discount }) }}
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.validity')"
              width="200"
            >
              <template #default="{ row }">
                <div class="validity-period">
                  <p>{{ formatDateTime(row.start_time) }}</p>
                  <p class="to-text">
                    {{ $t('promotions.toText') }}
                  </p>
                  <p>{{ formatDateTime(row.end_time) }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.statusColumn')"
              width="120"
              align="center"
            >
              <template #default="{ row }">
                <el-switch
                  v-if="row.status === 'active' || row.status === 'pending' || row.status === 'paused'"
                  :model-value="row.status === 'active'"
                  :loading="promotionToggleLoading[row.id] === true"
                  :active-text="$t('promotions.activatedStatus')"
                  :inactive-text="$t('promotions.unactivated')"
                  inline-prompt
                  @change="(val: boolean) => handleTogglePromotionStatus(row, val)"
                />
                <el-tag
                  v-else
                  :type="getPromoStatusType(row.status)"
                  effect="light"
                  size="small"
                >
                  {{ getPromoStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('promotions.actionsColumn')"
              width="180"
              fixed="right"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="handleEditPromotion(row)"
                >
                  {{ $t('promotions.edit') }}
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
            @change="loadPromotions"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Unified Promotion/Coupon Dialog -->
    <el-dialog
      v-model="formDialogVisible"
      :title="formDialogTitle"
      width="700px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        label-width="120px"
        :rules="formRules"
      >
        <el-form-item
          :label="$t('promotions.promotionName')"
          prop="name"
        >
          <el-input
            v-model="form.name"
            :placeholder="form.kind === 'coupon' ? $t('promotions.enterCouponName') : $t('promotions.enterPromotionName')"
            maxlength="100"
          />
        </el-form-item>

        <!-- Coupon-only: code -->
        <el-form-item
          v-if="form.kind === 'coupon'"
          :label="$t('promotions.couponCode')"
          prop="code"
        >
          <el-input
            v-model="form.code"
            :placeholder="$t('promotions.enterCouponCode')"
            maxlength="50"
          >
            <template #append>
              <el-button @click="generateCouponCode">
                {{ $t('promotions.generate') }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item :label="$t('promotions.promotionDescription')">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            :placeholder="$t('promotions.enterDescription')"
          />
        </el-form-item>

        <!-- Promotion-only: type select -->
        <el-form-item
          v-if="form.kind === 'promotion'"
          :label="$t('promotions.promotionTypeSelect')"
          prop="type"
        >
          <el-radio-group v-model="form.type">
            <el-radio label="discount">
              {{ $t('promotions.discount') }}
            </el-radio>
            <el-radio label="flash_sale">
              {{ $t('promotions.flashSale') }}
            </el-radio>
            <el-radio label="bundle">
              {{ $t('promotions.bundle') }}
            </el-radio>
            <el-radio label="buy_x_get_y">
              {{ $t('promotions.buyXGetY') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- Coupon-only: type radios (fixed_amount/percentage) -->
        <el-form-item
          v-if="form.kind === 'coupon'"
          :label="$t('promotions.couponType')"
          prop="type"
        >
          <el-radio-group v-model="form.type">
            <el-radio label="fixed_amount">
              {{ $t('promotions.fixedAmountType') }}
            </el-radio>
            <el-radio label="percentage">
              {{ $t('promotions.percentageType') }}
            </el-radio>
            <el-radio label="free_shipping">
              {{ $t('promotions.freeShipping') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- Currency (auto-bound from selected market) -->
        <el-form-item :label="$t('promotions.currency')">
          <el-input
            :model-value="form.currency || (selectedMarketCurrency || '—')"
            readonly
            disabled
            :placeholder="$t('promotions.currencyAutoBound')"
          />
        </el-form-item>

        <!-- Market -->
        <el-form-item :label="$t('promotions.market')">
          <el-select
            v-model="form.market_id"
            :placeholder="$t('promotions.selectMarket')"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="m in marketOptions"
              :key="m.id"
              :label="m.name"
              :value="m.id"
            />
          </el-select>
        </el-form-item>

        <!-- Rules -->
        <el-form-item :label="$t('promotions.rules')">
          <div class="rules-editor">
            <el-button
              v-if="!(form.kind === 'coupon' && (form.rules || []).length >= 1)"
              size="small"
              type="primary"
              @click="showAddRuleDialogInForm"
            >
              <el-icon><Plus /></el-icon>{{ $t('promotions.addRule') }}
            </el-button>
            <div
              v-if="form.rules && form.rules.length > 0"
              class="rules-list"
            >
              <div
                v-for="(rule, idx) in form.rules"
                :key="idx"
                class="rule-item"
              >
                <span class="rule-text">{{ formatRuleSummary(rule) }}</span>
                <el-button
                  type="primary"
                  link
                  size="small"
                  @click="editRuleFromForm(idx)"
                >
                  {{ $t('common.edit') }}
                </el-button>
                <el-button
                  type="danger"
                  link
                  size="small"
                  @click="removeRuleFromForm(idx)"
                >
                  {{ $t('common.delete') }}
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>

        <!-- Coupon-only: total_count -->
        <el-form-item
          v-if="form.kind === 'coupon'"
          :label="$t('promotions.totalCount')"
        >
          <el-input-number
            v-model="form.total_count"
            :min="0"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('promotions.totalCountTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="$t('promotions.usageLimit')">
          <el-input-number
            v-model="form.usage_limit"
            :min="0"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('promotions.zeroUnlimited') }}
          </div>
        </el-form-item>

        <el-form-item :label="$t('promotions.perUserLimit')">
          <el-input-number
            v-model="form.per_user_limit"
            :min="0"
            style="width: 100%"
          />
          <div class="form-tip">
            {{ $t('promotions.zeroUnlimited') }}
          </div>
        </el-form-item>

        <el-form-item :label="$t('promotions.scopeType')">
          <el-radio-group v-model="form.scope_type">
            <el-radio label="storewide">
              {{ $t('promotions.storewide') }}
            </el-radio>
            <el-radio label="products">
              {{ $t('promotions.specificProductsScope') }}
            </el-radio>
            <el-radio label="categories">
              {{ $t('promotions.specificCategories') }}
            </el-radio>
            <el-radio label="brands">
              {{ $t('promotions.specificBrands') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="$t('promotions.tags')">
          <el-input
            v-model="form.tagsText"
            :placeholder="$t('promotions.tagsPlaceholder')"
          />
        </el-form-item>

        <el-form-item
          :label="$t('promotions.validityPeriod')"
          prop="dateRange"
        >
          <el-date-picker
            v-model="form.dateRange"
            type="datetimerange"
            :start-placeholder="$t('promotions.startPlaceholder')"
            :end-placeholder="$t('promotions.endPlaceholder')"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">
          {{ $t('promotions.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="saveLoading"
          @click="saveForm"
        >
          {{ $t('promotions.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Add/Edit Rule Dialog (used by inline form) -->
    <el-dialog
      v-model="ruleDialogVisible"
      :title="editingRuleIndex === -1 ? $t('promotions.addRule') : $t('promotions.editRule')"
      width="500px"
      destroy-on-close
    >
      <el-form
        :model="ruleForm"
        label-width="120px"
      >
        <el-form-item :label="$t('promotions.conditionType')">
          <el-select
            v-model="ruleForm.condition_type"
            style="width: 100%"
          >
            <el-option
              :label="$t('promotions.minAmount')"
              value="min_amount"
            />
            <el-option
              :label="$t('promotions.minQuantity')"
              value="min_quantity"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('promotions.conditionValue')">
          <el-input-number
            v-model="ruleForm.condition_value_num"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.actionType')">
          <el-select
            v-model="ruleForm.action_type"
            style="width: 100%"
          >
            <el-option
              :label="$t('promotions.fixedAmountType')"
              value="fixed_amount"
            />
            <el-option
              :label="$t('promotions.percentageType')"
              value="percentage"
            />
            <el-option
              :label="$t('promotions.freeShipping')"
              value="free_shipping"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="ruleForm.action_type !== 'free_shipping'"
          :label="$t('promotions.actionValue')"
        >
          <el-input-number
            v-model="ruleForm.action_value_num"
            :min="0"
            :max="ruleForm.action_type === 'percentage' ? 100 : 99999"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item
          v-if="ruleForm.action_type === 'percentage'"
          :label="$t('promotions.maxDiscount')"
        >
          <el-input-number
            v-model="ruleForm.max_discount_num"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.sortOrder')">
          <el-input-number
            v-model="ruleForm.sort_order"
            :min="0"
            :max="100"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="ruleSaving"
          @click="saveRuleInForm"
        >
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Coupon Usage Dialog -->
    <el-dialog
      v-model="usageDialogVisible"
      :title="$t('promotions.couponUsageRecords')"
      width="800px"
    >
      <el-table
        v-loading="usageLoading"
        :data="usageList"
        stripe
      >
        <el-table-column
          prop="user_id"
          :label="$t('promotions.userID')"
          width="100"
        />
        <el-table-column
          prop="order_id"
          :label="$t('promotions.orderID')"
          width="120"
        />
        <el-table-column
          prop="discount_amount"
          :label="$t('promotions.discountAmountColumn')"
          width="120"
        >
          <template #default="{ row }">
            ¥{{ row.discount_amount }}
          </template>
        </el-table-column>
        <el-table-column
          prop="used_at"
          :label="$t('promotions.usedTime')"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDateTime(row.used_at) }}
          </template>
        </el-table-column>
      </el-table>
      <TablePagination
        v-model:current-page="usageParams.page"
        v-model:page-size="usageParams.page_size"
        :total="usageTotal"
        @change="loadCouponUsage"
      />
    </el-dialog>

    <!-- Batch Generate Coupons Dialog -->
    <el-dialog
      v-model="batchGenerateDialogVisible"
      :title="$t('promotions.batchGenerateCoupons')"
      width="600px"
      destroy-on-close
    >
      <el-form
        :model="batchGenerateForm"
        label-width="140px"
      >
        <el-form-item :label="$t('promotions.generateCouponPrefix')">
          <el-input
            v-model="batchGenerateForm.prefix"
            :placeholder="$t('promotions.generateCouponPrefixPlaceholder')"
            maxlength="10"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.generateCouponQuantity')">
          <el-input-number
            v-model="batchGenerateForm.quantity"
            :min="1"
            :max="1000"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.generateCouponLength')">
          <el-input-number
            v-model="batchGenerateForm.length"
            :min="4"
            :max="20"
            style="width: 100%"
          />
        </el-form-item>
        <el-divider>{{ $t('promotions.generateCouponConfig') }}</el-divider>
        <el-form-item :label="$t('promotions.couponName')">
          <el-input
            v-model="batchGenerateForm.couponName"
            :placeholder="$t('promotions.enterCouponName')"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.couponType')">
          <el-radio-group v-model="batchGenerateForm.couponType">
            <el-radio label="fixed_amount">
              {{ $t('promotions.fixedAmountType') }}
            </el-radio>
            <el-radio label="percentage">
              {{ $t('promotions.percentageType') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="batchGenerateForm.couponType === 'fixed_amount' ? $t('promotions.discountAmount') : $t('promotions.discountRatio')">
          <el-input-number
            v-model="batchGenerateForm.discountValue"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.minOrderAmount')">
          <el-input-number
            v-model="batchGenerateForm.minOrderAmount"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="$t('promotions.usageLimit')">
          <el-input-number
            v-model="batchGenerateForm.usageLimit"
            :min="0"
            style="width: 100%"
          />
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
        <el-button @click="batchGenerateDialogVisible = false">
          {{ $t('promotions.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="batchGenerating"
          @click="handleBatchGenerate"
        >
          {{ $t('promotions.batchGenerate') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Generated Codes Dialog -->
    <el-dialog
      v-model="generatedCodesDialogVisible"
      :title="$t('promotions.generatedCodes')"
      width="700px"
    >
      <div class="generated-codes-info">
        {{ $t('promotions.generatedCodesCount', { count: generatedCodes.length }) }}
        <el-button
          type="primary"
          link
          @click="handleCopyAllCodes"
        >
          {{ $t('promotions.copyAllCodes') }}
        </el-button>
      </div>
      <div class="generated-codes-list">
        <el-tag
          v-for="code in generatedCodes"
          :key="code"
          class="code-tag"
          type="info"
        >
          {{ code }}
        </el-tag>
      </div>
      <template #footer>
        <el-button
          type="primary"
          @click="generatedCodesDialogVisible = false"
        >
          {{ $t('common.close') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Issue Coupon to User Dialog -->
    <el-dialog
      v-model="issueDialogVisible"
      :title="$t('promotions.issueCoupon')"
      width="500px"
      destroy-on-close
    >
      <el-form
        :model="issueForm"
        label-width="100px"
      >
        <el-form-item :label="$t('promotions.currentCoupon')">
          <el-input
            :value="issueForm.coupon_name"
            disabled
          />
        </el-form-item>
        <el-form-item
          :label="$t('promotions.userId')"
          required
        >
          <el-select
            v-model="issueForm.selectedUsers"
            multiple
            filterable
            remote
            :remote-method="searchIssueUsers"
            :loading="issueSearchLoading"
            :placeholder="$t('promotions.searchUserPlaceholder')"
            style="width: 100%"
            value-key="id"
            @change="onIssueUserSelectionChange"
          >
            <el-option
              v-for="u in issueUserOptions"
              :key="u.id"
              :label="formatUserOptionLabel(u)"
              :value="u"
            >
              <div class="user-option">
                <div class="user-option-main">
                  <span class="user-email">{{ u.email }}</span>
                  <span class="user-name">{{ u.name }}</span>
                </div>
                <div
                  v-if="u.phone"
                  class="user-phone"
                >
                  {{ u.phone }}
                </div>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="issueDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="issueLoading"
          @click="handleConfirmIssue"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Ticket, CircleCheck, User, Search, Plus, Present, DocumentCopy } from '@element-plus/icons-vue'
import {
  getPromotionList, createPromotion, updatePromotion, deletePromotion,
  activatePromotion, deactivatePromotion,
  generateCouponCodes, issueUserCoupon, batchIssueUserCoupon,
  type Promotion, type PromotionKind, type PromotionRule, type PromotionRuleRequest,
  type CouponUsage,
  type GenerateCouponCodesRequest
} from '@/api/promotion'
import { getMarkets, type Market } from '@/api/market'
import { getUserList, type User as UserAccount } from '@/api/user'
import TablePagination from '@/components/common/TablePagination.vue'
import StatsCard from '@/components/common/StatsCard.vue'
import { t } from '@/plugins/i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { handleError } = useErrorHandler()

// Tab state
const activeTab = ref<'coupon' | 'promotion'>('coupon')

// Stats
const stats = ref({
  totalCoupons: 0,
  activeCoupons: 0,
  totalUsed: 0,
  totalPromotions: 0
})

// Markets (shared by both tabs and form)
const marketOptions = ref<Market[]>([])

// Coupon list state
const couponList = ref<Promotion[]>([])
const couponLoading = ref(false)
const couponTotal = ref(0)
const couponParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as string | undefined,
  market_id: undefined as string | undefined
})

// Promotion list state
const promotionList = ref<Promotion[]>([])
const promotionLoading = ref(false)
const promotionTotal = ref(0)
const promotionParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  status: undefined as string | undefined,
  market_id: undefined as string | undefined,
  type: undefined as string | undefined
})

// Usage dialog state
const usageDialogVisible = ref(false)
const usageList = ref<CouponUsage[]>([])
const usageLoading = ref(false)
const usageTotal = ref(0)
const usageParams = reactive({
  id: '',
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
const issueUserOptions = ref<UserAccount[]>([])
const issueSearchLoading = ref(false)
let issueSearchSeq = 0
const issueForm = reactive({
  selectedUsers: [] as UserAccount[],
  coupon_id: '' as string,
  coupon_name: ''
})

// =================== Unified form (merged from couponForm + promotionForm) ===================
const formDialogVisible = ref(false)
const isFormEdit = ref(false)
const saveLoading = ref(false)
const formRef = ref()
const form = reactive({
  id: '' as string,
  kind: 'promotion' as PromotionKind,
  name: '',
  description: '',
  code: '',
  type: 'discount' as string,
  currency: 'CNY',
  market_id: '' as string,
  scope_type: 'storewide' as 'storewide' | 'products' | 'categories' | 'brands',
  scope_ids: [] as string[],
  exclude_ids: [] as string[],
  tagsText: '',
  usage_limit: 0,
  per_user_limit: 0,
  total_count: 0,
  status: 'pending' as string,
  dateRange: [] as string[],
  rules: [] as PromotionRuleRequest[]
})

const formRules = {
  name: [{ required: true, message: t('promotions.enterPromotionName'), trigger: 'blur' }],
  type: [{ required: true, message: t('promotions.selectPromotionType'), trigger: 'change' }],
  code: [{
    validator: (_: unknown, value: string, cb: (err?: Error) => void) => {
      if (form.kind === 'coupon' && !value) {
        return cb(new Error(t('promotions.enterCouponCode')))
      }
      cb()
    },
    trigger: 'blur'
  }],
  dateRange: [{ required: true, message: t('promotions.selectValidityPeriod'), trigger: 'change' }]
}

const formDialogTitle = computed(() => {
  if (isFormEdit.value) {
    return form.kind === 'coupon' ? t('promotions.editCoupon') : t('promotions.editPromotion')
  }
  return form.kind === 'coupon' ? t('promotions.createCoupon') : t('promotions.createPromotion')
})

// Currency is auto-bound from the selected market — never user-editable.
// Looked up via marketOptions so changing the market updates form.currency.
const selectedMarketCurrency = computed(() => {
  if (!form.market_id) return ''
  const m = marketOptions.value.find(x => String(x.id) === String(form.market_id))
  return m?.currency || ''
})

watch(() => form.market_id, (newId) => {
  if (!newId) {
    form.currency = ''
    return
  }
  const m = marketOptions.value.find(x => String(x.id) === String(newId))
  if (m?.currency) form.currency = m.currency
})

// Rule editor (inline in form)
const ruleDialogVisible = ref(false)
const ruleSaving = ref(false)
const editingRuleIndex = ref(-1)
const ruleForm = reactive({
  condition_type: 'min_amount' as 'min_amount' | 'min_quantity',
  condition_value_num: 0,
  action_type: 'fixed_amount' as 'fixed_amount' | 'percentage' | 'free_shipping',
  action_value_num: 0,
  max_discount_num: 0,
  sort_order: 0
})

// =================== Loading functions ===================
const loadMarkets = async () => {
  try {
    const res = await getMarkets()
    marketOptions.value = res.list || []
  } catch (error) {
    handleError(error, t('promotions.loadMarketsFailed'))
  }
}

const loadCoupons = async () => {
  couponLoading.value = true
  try {
    const res = await getPromotionList({
      page: couponParams.page,
      page_size: couponParams.page_size,
      kind: 'coupon',
      name: couponParams.name || undefined,
      status: couponParams.status || undefined,
      market_id: couponParams.market_id || undefined
    })
    couponList.value = res.list || []
    couponTotal.value = res.total || 0
    // Update stats
    stats.value.totalCoupons = couponTotal.value
    stats.value.activeCoupons = couponList.value.filter(c => c.status === 'active').length
    stats.value.totalUsed = couponList.value.reduce((sum, c) => sum + (c.used_count || 0), 0)
  } catch (error) {
    handleError(error, t('promotions.loadCouponsFailed'))
  } finally {
    couponLoading.value = false
  }
}

const loadPromotions = async () => {
  promotionLoading.value = true
  try {
    const res = await getPromotionList({
      page: promotionParams.page,
      page_size: promotionParams.page_size,
      kind: 'promotion',
      name: promotionParams.name || undefined,
      status: promotionParams.status || undefined,
      market_id: promotionParams.market_id || undefined,
      type: promotionParams.type || undefined
    })
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
    const res = await (await import('@/api/promotion')).getCouponUsage(usageParams.id, { page: usageParams.page, page_size: usageParams.page_size })
    usageList.value = res.list || []
    usageTotal.value = res.total || 0
  } catch (error) {
    handleError(error, t('promotions.loadCouponUsageFailed'))
  } finally {
    usageLoading.value = false
  }
}

// =================== Helpers ===================
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

const getUsagePercentage = (row: Promotion) => {
  const denom = row.usage_limit > 0 ? row.usage_limit : (row.total_count && row.total_count > 0 ? row.total_count : 0)
  if (!denom) return 0
  return Math.round((row.used_count / denom) * 100)
}

const getProgressStatus = (row: Promotion) => {
  const percentage = getUsagePercentage(row)
  if (percentage >= 90) return 'exception'
  if (percentage >= 70) return 'warning'
  return ''
}

const getPromoStatusType = (status: string) => {
  const types: Record<string, string> = {
    'active': 'success',
    'paused': 'warning',
    'pending': 'info',
    'ended': 'info',
    'expired': 'danger',
    'depleted': 'danger'
  }
  return types[status] || 'info'
}

const getPromoStatusText = (status: string) => {
  const texts: Record<string, string> = {
    'active': t('promotions.activeStatus'),
    'paused': t('promotions.paused'),
    'pending': t('promotions.pending'),
    'ended': t('promotions.ended'),
    'expired': t('promotions.expiredStatus'),
    'depleted': t('promotions.depletedStatus')
  }
  return texts[status] || status
}

const getPromotionTypeText = (type: string) => {
  const texts: Record<string, string> = {
    'discount': t('promotions.discount'),
    'flash_sale': t('promotions.flashSale'),
    'bundle': t('promotions.bundle'),
    'buy_x_get_y': t('promotions.buyXGetY'),
    'fixed_amount': t('promotions.fixedAmount'),
    'percentage': t('promotions.percentage'),
    'free_shipping': t('promotions.freeShipping')
  }
  return texts[type] || type
}

const formatRuleSummary = (rule: PromotionRuleRequest | PromotionRule) => {
  const cond = rule.condition_type === 'min_amount'
    ? t('promotions.minAmountEquals', { value: rule.condition_value })
    : t('promotions.minQuantityEquals', { value: rule.condition_value })
  let action = ''
  if (rule.action_type === 'fixed_amount') action = `减 ¥${rule.action_value}`
  else if (rule.action_type === 'percentage') action = `${rule.action_value}% off`
  else if (rule.action_type === 'free_shipping') action = t('promotions.freeShipping')
  return `${cond} → ${action}`
}

const generateCouponCode = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let code = ''
  for (let i = 0; i < 10; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.code = code
}

const handleTabChange = (tab: string | number) => {
  if (tab === 'coupon') loadCoupons()
  else if (tab === 'promotion') loadPromotions()
}

// =================== Form actions ===================
const resetForm = () => {
  Object.assign(form, {
    id: '',
    kind: 'promotion',
    name: '',
    description: '',
    code: '',
    type: 'discount',
    currency: 'CNY',
    market_id: '',
    scope_type: 'storewide',
    scope_ids: [],
    exclude_ids: [],
    tagsText: '',
    usage_limit: 0,
    per_user_limit: 0,
    total_count: 0,
    status: 'pending',
    dateRange: [],
    rules: []
  })
}

const fillFormFromRow = (row: Promotion) => {
  Object.assign(form, {
    id: row.id,
    kind: row.kind,
    name: row.name,
    description: row.description || '',
    code: row.code || '',
    type: row.type || (row.kind === 'coupon' ? 'fixed_amount' : 'discount'),
    currency: row.currency || 'CNY',
    market_id: row.market_id || '',
    scope_type: row.scope_type || 'storewide',
    scope_ids: row.scope_ids || [],
    exclude_ids: row.exclude_ids || [],
    tagsText: (row.tags || []).join(', '),
    usage_limit: row.usage_limit || 0,
    per_user_limit: row.per_user_limit || 0,
    total_count: row.total_count || 0,
    status: row.status,
    dateRange: [row.start_time, row.end_time],
    rules: (row.rules || []).map(r => ({
      condition_type: r.condition_type,
      condition_value: r.condition_value,
      action_type: r.action_type,
      action_value: r.action_value,
      max_discount: r.max_discount,
      sort_order: r.sort_order
    }))
  })
}

const handleAddCoupon = () => {
  resetForm()
  form.kind = 'coupon'
  form.type = 'fixed_amount'
  isFormEdit.value = false
  formDialogVisible.value = true
}

const handleEditCoupon = async (row: Promotion) => {
  resetForm()
  try {
    const detail = await (await import('@/api/promotion')).getPromotion(row.id)
    fillFormFromRow(detail)
  } catch {
    fillFormFromRow(row)
  }
  // Defensive: handlers know which kind they're editing. Don't rely on the
  // row's kind field being populated if the API omitted it.
  form.kind = 'coupon'
  isFormEdit.value = true
  formDialogVisible.value = true
}

const handleAddPromotion = () => {
  resetForm()
  form.kind = 'promotion'
  form.type = 'discount'
  isFormEdit.value = false
  formDialogVisible.value = true
}

const handleEditPromotion = (row: Promotion) => {
  router.push(`/promotions/${row.id}`)
}

// =================== Rule editor in form ===================
const showAddRuleDialogInForm = () => {
  editingRuleIndex.value = -1
  ruleForm.condition_type = 'min_amount'
  ruleForm.condition_value_num = 0
  ruleForm.action_type = 'fixed_amount'
  ruleForm.action_value_num = 0
  ruleForm.max_discount_num = 0
  ruleForm.sort_order = 0
  ruleDialogVisible.value = true
}

const editRuleFromForm = (idx: number) => {
  const rule = form.rules[idx]
  if (!rule) return
  editingRuleIndex.value = idx
  ruleForm.condition_type = rule.condition_type
  ruleForm.condition_value_num = parseFloat(rule.condition_value) || 0
  ruleForm.action_type = rule.action_type
  ruleForm.action_value_num = parseFloat(rule.action_value) || 0
  ruleForm.max_discount_num = parseFloat(rule.max_discount || '0') || 0
  ruleForm.sort_order = rule.sort_order || 0
  ruleDialogVisible.value = true
}

const removeRuleFromForm = (idx: number) => {
  form.rules.splice(idx, 1)
}

const saveRuleInForm = () => {
  const data: PromotionRuleRequest = {
    condition_type: ruleForm.condition_type,
    condition_value: String(ruleForm.condition_value_num),
    action_type: ruleForm.action_type,
    action_value: ruleForm.action_type === 'free_shipping' ? '0' : String(ruleForm.action_value_num),
    max_discount: ruleForm.action_type === 'percentage' ? String(ruleForm.max_discount_num) : undefined,
    sort_order: ruleForm.sort_order
  }

  if (editingRuleIndex.value === -1) {
    form.rules.push(data)
  } else {
    form.rules[editingRuleIndex.value] = data
  }
  ruleDialogVisible.value = false
}

// =================== Save form ===================
const saveForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    saveLoading.value = true
    try {
      const tags = form.tagsText
        ? form.tagsText.split(',').map(s => s.trim()).filter(Boolean)
        : []

      const baseData = {
        kind: form.kind,
        name: form.name,
        description: form.description,
        type: form.type,
        currency: form.currency,
        market_id: form.market_id || undefined,
        usage_limit: form.usage_limit,
        per_user_limit: form.per_user_limit,
        start_time: form.dateRange[0],
        end_time: form.dateRange[1],
        scope_type: form.scope_type,
        scope_ids: form.scope_ids,
        exclude_ids: form.exclude_ids,
        tags,
        rules: form.rules.length > 0 ? form.rules : undefined
      }

      const data = form.kind === 'coupon'
        ? { ...baseData, code: form.code, total_count: form.total_count }
        : baseData

      if (isFormEdit.value) {
        await updatePromotion({ id: form.id, ...data })
      } else {
        await createPromotion(data as Parameters<typeof createPromotion>[0])
      }

      ElMessage.success(isFormEdit.value ? t('promotions.updateSuccess') : t('promotions.createSuccess'))
      formDialogVisible.value = false
      if (form.kind === 'coupon') loadCoupons()
      else loadPromotions()
    } catch (error) {
      handleError(error, t('promotions.savePromotionFailed'))
    } finally {
      saveLoading.value = false
    }
  })
}

// =================== Toggle / delete / usage ===================
const couponToggleLoading = reactive<Record<string, boolean>>({})
const handleToggleCouponStatus = async (row: Promotion, nextActive: boolean) => {
  const wasActive = row.status === 'active'
  row.status = nextActive ? 'active' : 'pending'
  couponToggleLoading[row.id] = true
  try {
    if (nextActive) {
      await activatePromotion(row.id)
      ElMessage.success(t('promotions.activateSuccess'))
    } else {
      await deactivatePromotion(row.id)
      ElMessage.success(t('promotions.deactivateSuccess'))
    }
    loadCoupons()
  } catch (error) {
    row.status = wasActive ? 'active' : 'pending'
    handleError(
      error,
      nextActive ? t('promotions.activateCouponFailed') : t('promotions.deactivateCouponFailed')
    )
  } finally {
    couponToggleLoading[row.id] = false
  }
}

const promotionToggleLoading = reactive<Record<string, boolean>>({})
const handleTogglePromotionStatus = async (row: Promotion, nextActive: boolean) => {
  const wasActive = row.status === 'active'
  row.status = nextActive ? 'active' : 'paused'
  promotionToggleLoading[row.id] = true
  try {
    if (nextActive) {
      await activatePromotion(row.id)
      ElMessage.success(t('promotions.activateSuccess'))
    } else {
      await deactivatePromotion(row.id)
      ElMessage.success(t('promotions.deactivateSuccess'))
    }
    loadPromotions()
  } catch (error) {
    row.status = wasActive ? 'active' : 'paused'
    handleError(
      error,
      nextActive ? t('promotions.activatePromotionFailed') : t('promotions.deactivatePromotionFailed')
    )
  } finally {
    promotionToggleLoading[row.id] = false
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
    if (row.kind === 'coupon') loadCoupons()
    else loadPromotions()
  } catch (error) {
    if (error !== 'cancel') {
      handleError(error, t('promotions.deleteFailed'))
    }
  }
}

const handleCouponUsage = (row: Promotion) => {
  usageParams.id = row.id
  usageParams.page = 1
  usageDialogVisible.value = true
  loadCouponUsage()
}

const handleIssueToUser = (row: Promotion) => {
  issueForm.selectedUsers = []
  issueForm.coupon_id = row.id
  issueForm.coupon_name = row.name
  issueUserOptions.value = []
  issueDialogVisible.value = true
}

const searchIssueUsers = async (query: string) => {
  const mySeq = ++issueSearchSeq
  issueSearchLoading.value = true
  try {
    const res = await getUserList({
      page: 1,
      page_size: 20,
      keyword: query || undefined,
      phone: /^\d+$/.test(query) ? query : undefined
    })
    if (mySeq !== issueSearchSeq) return
    issueUserOptions.value = res.list || []
  } catch (error) {
    handleError(error, t('promotions.searchUserFailed'))
  } finally {
    if (mySeq === issueSearchSeq) issueSearchLoading.value = false
  }
}

const onIssueUserSelectionChange = (selected: UserAccount[]) => {
  const known = new Map(issueUserOptions.value.map(u => [u.id, u]))
  for (const u of selected) {
    if (!known.has(u.id)) known.set(u.id, u)
  }
  issueUserOptions.value = Array.from(known.values())
}

const formatUserOptionLabel = (u: UserAccount) => {
  return u.name ? `${u.email} (${u.name})` : u.email
}

const handleConfirmIssue = async () => {
  if (issueForm.selectedUsers.length === 0) {
    ElMessage.warning(t('promotions.selectAtLeastOneUser'))
    return
  }

  issueLoading.value = true
  try {
    const userIds = issueForm.selectedUsers.map(u => u.id)
    if (userIds.length === 1) {
      await issueUserCoupon({
        user_id: userIds[0],
        coupon_id: issueForm.coupon_id
      })
      ElMessage.success(t('promotions.issueSuccess'))
    } else {
      const res = await batchIssueUserCoupon({
        coupon_id: issueForm.coupon_id,
        user_ids: userIds
      })
      ElMessage.success(t('promotions.batchIssueSuccess', { count: res.issued }))
    }
    issueDialogVisible.value = false
  } catch (error) {
    handleError(error, t('promotions.issueFailed'))
  } finally {
    issueLoading.value = false
  }
}

// =================== Batch generate ===================
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

onMounted(() => {
  loadMarkets()
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

.rules-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 8px;
}

.rule-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: #F9FAFB;
  border-radius: 8px;
  font-size: 13px;
}

.rule-text {
  flex: 1;
  color: #4B5563;
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
