# Payment Module UI Design Specification

## Design Overview

This document specifies the UI components for the ShopJoy admin panel payment module, following the existing design patterns from orders and fulfillment pages.

---

## Design System Reference

### Color Palette

| Token | Light Mode | Usage |
|-------|------------|-------|
| Primary | `#6366F1` (Indigo) | Buttons, links, highlights |
| Primary Light | `#818CF8` | Gradients, hover states |
| Success | `#10B981` | Success status, positive amounts |
| Warning | `#F59E0B` | Pending status, warnings |
| Danger | `#EF4444` | Refunds, errors, negative amounts |
| Text Primary | `#1E1B4B` | Headings, important text |
| Text Secondary | `#6B7280` | Labels, descriptions |
| Background | `#F5F3FF` | Hover states, highlighted rows |
| Card Border | `rgba(99, 102, 241, 0.06)` | Card borders |
| **Stripe Brand** | `#635BFF` | Stripe channel tags, icons |

### Channel Brand Colors

| Channel | Color | Usage |
|---------|-------|-------|
| Stripe | `#635BFF` | Stripe tags, badges |
| Alipay | `#1677FF` | Alipay tags (Phase 2) |
| WeChat | `#07C160` | WeChat tags (Phase 2) |

### Supported Currencies (Phase 1)

| Currency | Code | Stripe Support |
|----------|------|----------------|
| US Dollar | USD | ✅ |
| Euro | EUR | ✅ |
| British Pound | GBP | ✅ |
| Japanese Yen | JPY | ✅ |
| Singapore Dollar | SGD | ✅ |
| Hong Kong Dollar | HKD | ✅ |
| Chinese Yuan | CNY | ❌ (Phase 2 Alipay/WeChat) |

### Typography

| Element | Font | Size | Weight |
|---------|------|------|--------|
| Page Title | System | 24px | 700 |
| Card Title | System | 16px | 600 |
| Stat Number | Fira Sans | 32px | 700 |
| Amount Value | Fira Sans | 16px | 700 |
| Transaction ID | Fira Code | 13px | 500 |
| Body Text | System | 14px | 400 |
| Label | System | 13px | 500 |

### Spacing & Radius

| Element | Value |
|---------|-------|
| Card Radius | 16px |
| Input Radius | 12px |
| Image Radius | 10px |
| Card Padding | 24px |
| Grid Gutter | 16px |
| Section Margin | 20px |

---

## 1. Payment Statistics Card Component

### Component: `PaymentStatsCard.vue`

#### Location
- Dashboard page
- Payment transactions list page header

#### UI Mockup

```
+------------------------------------------------------------------+
|                        Payment Overview                           |
+------------------------------------------------------------------+
|  +----------------+  +----------------+  +----------------+      |
|  | Today Received |  | 7-Day Received |  | Refund Amount  |      |
|  |                |  |                |  |                |      |
|  |    CNY 12,580  |  |   CNY 89,420   |  |    CNY 1,250   |      |
|  |    +15.2%      |  |    +8.3%       |  |    Rate: 1.4%  |      |
|  +----------------+  +----------------+  +----------------+      |
|                                                                   |
|  +----------------------------------------------------------+    |
|  |              Payment Channel Distribution                 |    |
|  |                                                          |    |
|  |  [=======Stripe 45%=======]                              |    |
|  |  [=====Alipay 30%=====]                                   |    |
|  |  [===WeChat 25%===]                                       |    |
|  |                                                          |    |
|  +----------------------------------------------------------+    |
+------------------------------------------------------------------+
```

#### Component Structure

```vue
<template>
  <div class="payment-stats">
    <!-- Amount Stats Row -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="24" :sm="8">
        <div class="stat-item today">
          <div class="stat-header">
            <el-icon><Calendar /></el-icon>
            <span class="stat-label">Today Received</span>
          </div>
          <p class="stat-amount">{{ currency }} {{ formatAmount(stats.todayReceived) }}</p>
          <div class="stat-trend up">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ stats.todayGrowth }}%</span>
            <span class="trend-label">vs yesterday</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-item week">
          <div class="stat-header">
            <el-icon><Timer /></el-icon>
            <span class="stat-label">7-Day Received</span>
          </div>
          <p class="stat-amount">{{ currency }} {{ formatAmount(stats.weekReceived) }}</p>
          <div class="stat-trend up">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ stats.weekGrowth }}%</span>
            <span class="trend-label">vs last week</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-item refund">
          <div class="stat-header">
            <el-icon><RefreshLeft /></el-icon>
            <span class="stat-label">Refund Amount</span>
          </div>
          <p class="stat-amount refund-amount">{{ currency }} {{ formatAmount(stats.refundAmount) }}</p>
          <div class="stat-trend">
            <span class="refund-rate">Rate: {{ stats.refundRate }}%</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Channel Distribution Card -->
    <el-card class="channel-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><PieChart /></el-icon>
            Payment Channel Distribution
          </span>
        </div>
      </template>
      <el-row :gutter="24">
        <el-col :xs="24" :md="12">
          <!-- Progress Bars -->
          <div class="channel-list">
            <div v-for="channel in channelDistribution" :key="channel.name" class="channel-item">
              <div class="channel-info">
                <span class="channel-name">{{ channel.name }}</span>
                <span class="channel-percent">{{ channel.percent }}%</span>
              </div>
              <el-progress
                :percentage="channel.percent"
                :color="channel.color"
                :stroke-width="8"
                :show-text="false"
              />
              <div class="channel-amount">
                {{ currency }} {{ formatAmount(channel.amount) }}
                <span class="channel-count">{{ channel.count }} transactions</span>
              </div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :md="12">
          <!-- Mini Pie Chart Placeholder -->
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon :size="48"><PieChart /></el-icon>
              <p>Channel Distribution Chart</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>
```

#### Props Definition

```typescript
interface PaymentStats {
  todayReceived: number       // Amount in cents
  todayGrowth: number         // Percentage
  weekReceived: number        // Amount in cents
  weekGrowth: number          // Percentage
  refundAmount: number        // Amount in cents
  refundRate: number          // Percentage
}

interface ChannelDistribution {
  name: string               // 'Stripe' | 'Alipay' | 'WeChat'
  percent: number            // 0-100
  amount: number             // Amount in cents
  count: number              // Transaction count
  color: string              // Chart color
}

defineProps<{
  stats: PaymentStats
  channelDistribution: ChannelDistribution[]
  currency?: string          // Default: 'CNY'
}>()
```

#### Styling

```css
.payment-stats {
  margin-bottom: 20px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-item {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-item.today {
  border-left: 4px solid #6366F1;
}

.stat-item.week {
  border-left: 4px solid #10B981;
}

.stat-item.refund {
  border-left: 4px solid #EF4444;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.stat-header .el-icon {
  color: #6366F1;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
}

.stat-amount {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 8px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-amount.refund-amount {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.stat-trend.up {
  color: #10B981;
}

.stat-trend .el-icon {
  font-size: 14px;
}

.trend-label {
  color: #9CA3AF;
  margin-left: 4px;
}

.refund-rate {
  color: #6B7280;
  font-weight: 500;
}

/* Channel Card */
.channel-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.channel-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.channel-item {
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.channel-item:hover {
  background: #F5F3FF;
}

.channel-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.channel-name {
  font-weight: 500;
  color: #1E1B4B;
}

.channel-percent {
  font-weight: 600;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

.channel-amount {
  margin-top: 8px;
  font-size: 13px;
  color: #6B7280;
  display: flex;
  justify-content: space-between;
}

.channel-count {
  color: #9CA3AF;
}

.chart-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  background: #F9FAFB;
  border-radius: 12px;
}

.chart-placeholder {
  text-align: center;
  color: #9CA3AF;
}

/* Responsive */
@media (max-width: 768px) {
  .stat-amount {
    font-size: 24px;
  }

  .stat-item {
    margin-bottom: 12px;
  }
}
```

---

## 2. Order Detail Page - Payment Info Section

### Component: Add to existing order detail page

#### UI Mockup

```
+------------------------------------------------------------------+
|                     Payment Information                           |
+------------------------------------------------------------------+
|                                                                   |
|  Payment Channel:     [Stripe]                                    |
|                                                                   |
|  PaymentIntent ID:    pi_3Oabc...xyz        [Copy]               |
|  Charge ID:           ch_3Oabc...xyz        [Copy]               |
|                                                                   |
|  Paid Amount:         USD 299.00   Transaction Fee: USD 2.99     |
|                                                                   |
|  Payment Time:        2024-03-18T14:32:18Z                        |
|                                                                   |
|  Payment Status:      [Success]                                   |
|                                                                   |
+------------------------------------------------------------------+
```

#### Component Structure

```vue
<template>
  <el-card class="payment-info-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="card-title">
          <el-icon><CreditCard /></el-icon>
          Payment Information
        </span>
        <el-button
          v-if="canRefund"
          type="danger"
          size="small"
          @click="openRefundDialog"
        >
          <el-icon><RefreshLeft /></el-icon>
          Refund
        </el-button>
      </div>
    </template>

    <el-descriptions :column="2" border>
      <el-descriptions-item label="Payment Channel">
        <div class="channel-cell">
          <el-tag :type="getChannelTagType(payment.channel)" effect="plain">
            {{ payment.channel }}
          </el-tag>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="">
        <!-- Empty for alignment -->
      </el-descriptions-item>
      <el-descriptions-item label="PaymentIntent ID">
        <div class="transaction-id-cell">
          <span class="transaction-id" :title="payment.channelIntentId">
            {{ truncateId(payment.channelIntentId) }}
          </span>
          <el-button link size="small" @click="copyToClipboard(payment.channelIntentId, 'PaymentIntent ID')">
            <el-icon><CopyDocument /></el-icon>
          </el-button>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="Charge ID">
        <div class="transaction-id-cell">
          <span class="transaction-id" :title="payment.channelPaymentId">
            {{ truncateId(payment.channelPaymentId) }}
          </span>
          <el-button link size="small" @click="copyToClipboard(payment.channelPaymentId, 'Charge ID')">
            <el-icon><CopyDocument /></el-icon>
          </el-button>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="Paid Amount">
        <span class="paid-amount">
          {{ payment.currency }} {{ formatAmount(payment.paidAmount) }}
        </span>
      </el-descriptions-item>
      <el-descriptions-item label="Transaction Fee">
        <span class="fee-amount">
          {{ payment.feeCurrency || payment.currency }} {{ formatAmount(payment.transactionFee) }}
        </span>
      </el-descriptions-item>
      <el-descriptions-item label="Payment Time">
        <span class="time-text">{{ payment.paidAt }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="Payment Status">
        <status-tag :status="payment.status" :type-map="paymentStatusMap" />
      </el-descriptions-item>
    </el-descriptions>

    <!-- Refund History -->
    <div v-if="payment.refunds && payment.refunds.length > 0" class="refund-history">
      <p class="history-title">
        <el-icon><List /></el-icon>
        Refund History
      </p>
      <el-table :data="payment.refunds" size="small">
        <el-table-column label="Refund No." min-width="140">
          <template #default="{ row }">
            <span class="refund-id">{{ row.refundNo }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Channel Refund ID" min-width="160">
          <template #default="{ row }">
            <div v-if="row.channelRefundId" class="channel-refund-cell">
              <span class="channel-refund-id" :title="row.channelRefundId">
                {{ truncateId(row.channelRefundId) }}
              </span>
              <el-button link size="small" @click="copyToClipboard(row.channelRefundId, 'Channel Refund ID')">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column label="Amount" width="120" align="right">
          <template #default="{ row }">
            <span class="refund-amount">{{ payment.currency }} {{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Status" width="100" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="refundStatusMap" />
          </template>
        </el-table-column>
        <el-table-column label="Time" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ row.refundedAt }}</span>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </el-card>
</template>
```

#### Data Types

```typescript
interface PaymentInfo {
  channel: 'Stripe' | 'Alipay' | 'WeChat'
  channelIntentId: string       // PaymentIntent ID (pi_xxx for Stripe)
  channelPaymentId: string      // Charge ID (ch_xxx for Stripe)
  paidAmount: number            // Amount in cents
  transactionFee: number        // Amount in cents
  currency: string              // 'USD' | 'EUR' | 'GBP' | 'JPY' | 'SGD' | 'HKD'
  feeCurrency: string           // Fee currency
  paidAt: string                // RFC3339 datetime
  status: 'pending' | 'processing' | 'success' | 'failed' | 'cancelled' | 'refunded' | 'partially_refunded' | 'requires_action'
  refunds: RefundRecord[]
}

interface RefundRecord {
  refundNo: string
  channelRefundId: string       // re_xxx for Stripe
  amount: number
  status: 'pending' | 'succeeded' | 'failed'  // Aligned with PRD PaymentRefund status
  refundedAt: string
}

// Payment status mapping (aligned with PRD)
const paymentStatusMap = {
  pending: { type: 'warning', text: 'Pending' },
  processing: { type: 'primary', text: 'Processing' },
  success: { type: 'success', text: 'Success' },
  failed: { type: 'danger', text: 'Failed' },
  cancelled: { type: 'info', text: 'Cancelled' },
  refunded: { type: 'info', text: 'Refunded' },
  partially_refunded: { type: 'warning', text: 'Partially Refunded' },
  requires_action: { type: 'warning', text: 'Requires Action' }  // 3D Secure
}

// Refund status mapping (PaymentRefund entity)
const refundStatusMap = {
  pending: { type: 'warning', text: 'Pending' },
  succeeded: { type: 'success', text: 'Succeeded' },
  failed: { type: 'danger', text: 'Failed' }
}

// Helper: Truncate long Stripe IDs
const truncateId = (id: string) => {
  if (!id || id.length <= 16) return id
  return `${id.slice(0, 8)}...${id.slice(-4)}`
}

// Helper: Copy to clipboard
const copyToClipboard = (text: string, label: string) => {
  navigator.clipboard.writeText(text)
  ElMessage.success(`${label} copied`)
}
```

#### Refund Reason Types (Aligned with fulfillment module)

```typescript
// Use uppercase codes, aligned with existing fulfillment.RefundReason
const refundReasonOptions = [
  { code: 'DEFECTIVE', name: 'Product Defective' },
  { code: 'WRONG_ITEM', name: 'Wrong Item Received' },
  { code: 'NOT_AS_DESCRIBED', name: 'Not As Described' },
  { code: 'DAMAGED', name: 'Damaged in Transit' },
  { code: 'NO_LONGER_NEEDED', name: 'No Longer Needed' },
  { code: 'LATE_DELIVERY', name: 'Late Delivery' },
  { code: 'OTHER', name: 'Other' }
]
```

#### Styling

```css
.payment-info-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title .el-icon {
  color: #6366F1;
}

.channel-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.channel-cell .el-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

.transaction-id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.transaction-id {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.paid-amount {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}

.fee-amount {
  color: #6B7280;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.refund-history {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
}

.history-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 12px 0;
}

.refund-id {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  color: #6366F1;
}

.refund-amount {
  color: #EF4444;
  font-weight: 500;
}
```

---

## 3. Order Detail Page - Refund Dialog

### Component: `RefundDialog.vue`

#### UI Mockup - Step 1: Refund Type Selection

```
+------------------------------------------------------------------+
|                     Initiate Refund                          [X] |
+------------------------------------------------------------------+
|                                                                   |
|  Order: ORD2024031800100                                          |
|  Paid Amount: CNY 299.00                                          |
|  Max Refundable: CNY 299.00                                       |
|                                                                   |
|  Refund Type:                                                     |
|  +--------------------------+  +--------------------------+       |
|  |    [O] Full Refund      |  |    [ ] Partial Refund    |       |
|  |    Refund full amount   |  |    Refund part of amount |       |
|  +--------------------------+  +--------------------------+       |
|                                                                   |
|  Reason:                                                          |
|  +----------------------------------------------------------+    |
|  | Select refund reason...                                v  |    |
|  +----------------------------------------------------------+    |
|                                                                   |
|  Additional Notes (Optional):                                     |
|  +----------------------------------------------------------+    |
|  |                                                          |    |
|  |                                                          |    |
|  +----------------------------------------------------------+    |
|                                                                   |
|                                          [Cancel]  [Next Step >]  |
+------------------------------------------------------------------+
```

#### UI Mockup - Step 2: Partial Refund Amount

```
+------------------------------------------------------------------+
|                     Initiate Refund                          [X] |
+------------------------------------------------------------------+
|                                                                   |
|  Step 2 of 2: Enter Refund Amount                                 |
|                                                                   |
|  +----------------------------------------------------------+    |
|  | Refund Amount                                            |    |
|  |                                                          |    |
|  |  CNY [______150.00______]                                |    |
|  |                                                          |    |
|  |  Max: CNY 299.00 | Original: CNY 299.00                  |    |
|  +----------------------------------------------------------+    |
|                                                                   |
|  Selected Reason: Product Defective                               |
|                                                                   |
|  Summary:                                                         |
|  +----------------------------------------------------------+    |
|  | Order Total:                             CNY 299.00       |    |
|  | Refund Amount:                           CNY 150.00       |    |
|  | Remaining Amount:                        CNY 149.00       |    |
|  +----------------------------------------------------------+    |
|                                                                   |
|                                   [< Back]  [Cancel]  [Confirm]  |
+------------------------------------------------------------------+
```

#### Component Structure

```vue
<template>
  <el-dialog
    v-model="visible"
    title="Initiate Refund"
    width="560px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-steps :active="currentStep" simple class="refund-steps">
      <el-step title="Select Type" />
      <el-step title="Confirm" />
    </el-steps>

    <!-- Step 1: Type & Reason -->
    <div v-show="currentStep === 0" class="step-content">
      <!-- Order Info -->
      <div class="order-info-banner">
        <el-icon><Document /></el-icon>
        <div class="order-info-content">
          <p class="order-no">Order: {{ order.orderNo }}</p>
          <p class="order-amount">Paid: {{ order.currency }} {{ formatAmount(order.paidAmount) }}</p>
        </div>
      </div>

      <!-- Refund Type -->
      <div class="refund-type-section">
        <p class="section-label">Refund Type</p>
        <el-radio-group v-model="form.refundType" class="refund-type-group">
          <el-radio-button value="full">
            <div class="type-option">
              <el-icon><FullScreen /></el-icon>
              <span class="type-title">Full Refund</span>
              <span class="type-desc">Refund full amount</span>
            </div>
          </el-radio-button>
          <el-radio-button value="partial">
            <div class="type-option">
              <el-icon><Compass /></el-icon>
              <span class="type-title">Partial Refund</span>
              <span class="type-desc">Refund part of amount</span>
            </div>
          </el-radio-button>
        </el-radio-group>
      </div>

      <!-- Partial Amount Input -->
      <transition name="slide-down">
        <div v-if="form.refundType === 'partial'" class="amount-input-section">
          <p class="section-label">Refund Amount</p>
          <div class="amount-input-wrapper">
            <span class="currency-label">{{ order.currency }}</span>
            <el-input-number
              v-model="form.refundAmount"
              :min="1"
              :max="maxRefundable / 100"
              :precision="2"
              :controls="false"
              class="amount-input"
            />
          </div>
          <div class="amount-hint">
            <span>Max refundable: {{ order.currency }} {{ formatAmount(maxRefundable) }}</span>
            <el-progress
              :percentage="refundPercentage"
              :stroke-width="4"
              :show-text="false"
              class="amount-progress"
            />
          </div>
        </div>
      </transition>

      <!-- Reason -->
      <div class="reason-section">
        <p class="section-label">Refund Reason</p>
        <el-select v-model="form.reasonType" placeholder="Select refund reason" style="width: 100%">
          <el-option
            v-for="reason in refundReasons"
            :key="reason.code"
            :label="reason.name"
            :value="reason.code"
          />
        </el-select>
      </div>

      <!-- Notes -->
      <div class="notes-section">
        <p class="section-label">Additional Notes (Optional)</p>
        <el-input
          v-model="form.notes"
          type="textarea"
          :rows="3"
          placeholder="Enter any additional information..."
          maxlength="500"
          show-word-limit
        />
      </div>
    </div>

    <!-- Step 2: Confirmation -->
    <div v-show="currentStep === 1" class="step-content">
      <el-alert type="warning" :closable="false" class="confirm-alert">
        <template #title>
          Please confirm the refund details below
        </template>
      </el-alert>

      <div class="confirm-summary">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="Order No.">
            <span class="value-text">{{ order.orderNo }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="Refund Type">
            <el-tag :type="form.refundType === 'full' ? 'primary' : 'warning'" size="small">
              {{ form.refundType === 'full' ? 'Full Refund' : 'Partial Refund' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Refund Amount">
            <span class="confirm-amount">
              {{ order.currency }} {{ formatAmount(actualRefundAmount) }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="Reason">
            {{ getReasonName(form.reasonType) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="form.notes" label="Notes">
            {{ form.notes }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="form.refundType === 'partial'" class="remaining-info">
          <el-icon><InfoFilled /></el-icon>
          <span>Remaining order amount: {{ order.currency }} {{ formatAmount(remainingAmount) }}</span>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button v-if="currentStep > 0" @click="prevStep">
          <el-icon><ArrowLeft /></el-icon>
          Back
        </el-button>
        <el-button @click="handleCancel">Cancel</el-button>
        <el-button
          v-if="currentStep === 0"
          type="primary"
          :disabled="!canProceed"
          @click="nextStep"
        >
          Next Step
          <el-icon><ArrowRight /></el-icon>
        </el-button>
        <el-button
          v-else
          type="danger"
          :loading="submitting"
          @click="confirmRefund"
        >
          <el-icon><Check /></el-icon>
          Confirm Refund
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>
```

#### Props & State

```typescript
import { v4 as uuidv4 } from 'uuid'

interface RefundForm {
  idempotencyKey: string         // Generated by frontend for idempotency
  refundType: 'full' | 'partial'
  refundAmount: number           // In currency units (not cents)
  reasonType: string             // Uppercase: DEFECTIVE, WRONG_ITEM, etc.
  notes: string
}

interface Order {
  orderId: string
  orderNo: string
  paidAmount: number             // In cents
  currency: string               // USD, EUR, GBP, JPY, SGD, HKD
}

const props = defineProps<{
  modelValue: boolean
  order: Order
  maxRefundable: number          // In cents
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

// State
const currentStep = ref(0)
const submitting = ref(false)
const form = reactive<RefundForm>({
  idempotencyKey: uuidv4(),       // Generate once on dialog open
  refundType: 'full',
  refundAmount: 0,
  reasonType: '',
  notes: ''
})

// Regenerate idempotency key on dialog open
watch(() => props.modelValue, (val) => {
  if (val) {
    form.idempotencyKey = uuidv4()
    form.refundType = 'full'
    form.refundAmount = 0
    form.reasonType = ''
    form.notes = ''
    currentStep.value = 0
  }
})

// Computed
const actualRefundAmount = computed(() => {
  return form.refundType === 'full'
    ? props.maxRefundable
    : Math.round(form.refundAmount * 100)
})

const remainingAmount = computed(() => {
  return props.order.paidAmount - actualRefundAmount.value
})

const refundPercentage = computed(() => {
  if (props.maxRefundable === 0) return 0
  return Math.round((actualRefundAmount.value / props.maxRefundable) * 100)
})

const canProceed = computed(() => {
  if (!form.reasonType) return false
  if (form.refundType === 'partial' && form.refundAmount <= 0) return false
  return true
})
```

#### Styling

```css
.refund-steps {
  margin-bottom: 24px;
}

.step-content {
  min-height: 300px;
}

.order-info-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  margin-bottom: 20px;
}

.order-info-banner .el-icon {
  font-size: 24px;
  color: #6366F1;
}

.order-info-content {
  flex: 1;
}

.order-no {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
  font-family: 'Fira Code', monospace;
}

.order-amount {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

.section-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin: 0 0 8px 0;
}

.refund-type-section {
  margin-bottom: 20px;
}

.refund-type-group {
  display: flex;
  gap: 12px;
  width: 100%;
}

.refund-type-group :deep(.el-radio-button) {
  flex: 1;
}

.refund-type-group :deep(.el-radio-button__inner) {
  width: 100%;
  border-radius: 12px !important;
  border: 1px solid #E5E7EB;
  padding: 16px;
}

.refund-type-group :deep(.el-radio-button.is-active .el-radio-button__inner) {
  border-color: #6366F1;
  background: #F5F3FF;
}

.type-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.type-option .el-icon {
  font-size: 24px;
  color: #6366F1;
}

.type-title {
  font-weight: 600;
  color: #1E1B4B;
}

.type-desc {
  font-size: 12px;
  color: #6B7280;
}

.amount-input-section {
  margin-bottom: 20px;
}

.amount-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 12px 16px;
}

.currency-label {
  font-weight: 600;
  color: #6366F1;
  font-size: 16px;
}

.amount-input {
  flex: 1;
}

.amount-input :deep(.el-input__inner) {
  font-size: 24px;
  font-weight: 700;
  text-align: right;
  background: transparent;
  border: none;
  box-shadow: none;
}

.amount-hint {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #6B7280;
}

.amount-progress {
  flex: 1;
}

.reason-section,
.notes-section {
  margin-bottom: 20px;
}

/* Step 2 */
.confirm-alert {
  margin-bottom: 16px;
}

.confirm-summary {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.confirm-amount {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}

.remaining-info {
  margin-top: 12px;
  padding: 12px;
  background: #FEF3C7;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #92400E;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Transitions */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Responsive */
@media (max-width: 576px) {
  .refund-type-group {
    flex-direction: column;
  }
}
```

---

## 4. Refund List Page Enhancement

### Component: Enhance existing `/views/fulfillment/refunds/index.vue`

#### Additional Columns to Add

```vue
<!-- Add after "Amount" column -->
<el-table-column label="Payment Channel" width="120" align="center">
  <template #default="{ row }">
    <el-tag :type="getChannelTagType(row.payment_channel)" effect="plain" size="small">
      {{ row.payment_channel }}
    </el-tag>
  </template>
</el-table-column>

<el-table-column label="Channel Refund ID" min-width="160">
  <template #default="{ row }">
    <div v-if="row.channel_refund_id" class="channel-refund-cell">
      <span class="channel-refund-id">{{ row.channel_refund_id }}</span>
      <el-button link size="small" @click="copyChannelRefundId(row.channel_refund_id)">
        <el-icon><CopyDocument /></el-icon>
      </el-button>
    </div>
    <span v-else class="no-data">-</span>
  </template>
</el-table-column>

<el-table-column label="Refund Fee" width="100" align="right">
  <template #default="{ row }">
    <span v-if="row.refund_fee" class="refund-fee">
      {{ row.currency }} {{ formatAmount(row.refund_fee) }}
    </span>
    <span v-else class="no-data">-</span>
  </template>
</el-table-column>
```

#### Add Retry Button for Failed Refunds

```vue
<!-- Update Actions column -->
<el-table-column label="Actions" width="220" fixed="right">
  <template #default="{ row }">
    <el-button type="primary" link size="small" @click="viewDetail(row)">
      Details
    </el-button>
    <el-button
      v-if="row.status === 0"
      type="success"
      link
      size="small"
      @click="quickApprove(row)"
    >
      Approve
    </el-button>
    <el-button
      v-if="row.status === 0"
      type="danger"
      link
      size="small"
      @click="openRejectDialog(row)"
    >
      Reject
    </el-button>
    <!-- NEW: Retry button for failed refunds -->
    <el-button
      v-if="row.status === 5"
      type="warning"
      link
      size="small"
      @click="retryRefund(row)"
    >
      <el-icon><RefreshRight /></el-icon>
      Retry
    </el-button>
  </template>
</el-table-column>
```

#### Additional State & Methods

```typescript
// Add to status type map
const statusTypeMap = {
  0: { type: 'warning' as const, text: 'Pending' },
  1: { type: 'success' as const, text: 'Approved' },
  2: { type: 'danger' as const, text: 'Rejected' },
  3: { type: 'primary' as const, text: 'Completed' },
  4: { type: 'info' as const, text: 'Cancelled' },
  5: { type: 'danger' as const, text: 'Failed' }  // NEW
}

// New method
const retryRefund = async (row: Refund) => {
  try {
    await ElMessageBox.confirm(
      `Retry refund of ${row.currency} ${(row.amount / 100).toFixed(2)}?`,
      'Retry Refund',
      { type: 'warning' }
    )
    // Call retry API
    await retryRefundApi({ refund_id: row.id })
    ElMessage.success('Refund retry initiated')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to retry refund')
    }
  }
}

const getChannelTagType = (channel: string) => {
  const types: Record<string, string> = {
    'Stripe': 'primary',
    'Alipay': 'success',
    'WeChat': 'warning'
  }
  return types[channel] || 'info'
}

const copyChannelRefundId = (id: string) => {
  navigator.clipboard.writeText(id)
  ElMessage.success('Channel refund ID copied')
}
```

#### Additional Styling

```css
.channel-refund-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.channel-refund-id {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  color: #6366F1;
}

.refund-fee {
  color: #6B7280;
  font-size: 13px;
}

.no-data {
  color: #9CA3AF;
}
```

---

## 5. Payment Transactions List Page (New)

### Component: `/views/payments/transactions/index.vue`

#### UI Mockup

```
+------------------------------------------------------------------+
|                    Payment Transactions                       [+] |
+------------------------------------------------------------------+
|  +----------------+  +----------------+  +----------------+      |
|  |    Success     |  |    Pending     |  |    Failed      |      |
|  |      156       |  |       3        |  |       2        |      |
|  +----------------+  +----------------+  +----------------+      |
+------------------------------------------------------------------+
|                                                                   |
|  [Search transaction ID/order no...] [Status v] [Channel v]       |
|  [Start Date] to [End Date]                    [Export] [Refresh] |
|                                                                   |
+------------------------------------------------------------------+
| Transaction ID    | Order No.      | Channel | Amount  | Fee     |
|-------------------|----------------|---------|---------|---------|
| txn_abc123...     | ORD20240318... | Stripe  | CNY 299 | CNY 2.99|
| txn_def456...     | ORD20240317... | Alipay  | CNY 456 | CNY 4.56|
| txn_ghi789...     | ORD20240317... | WeChat  | CNY 129 | CNY 1.29|
+------------------------------------------------------------------+
| Status      | Created At              | Actions                  |
|-------------|-------------------------|--------------------------|
| [Success]   | 2024-03-18 14:32:18     | [Details]                |
| [Success]   | 2024-03-17 16:45:30     | [Details]                |
| [Pending]   | 2024-03-17 12:00:00     | [Details] [Query Status] |
+------------------------------------------------------------------+
|                                                                   |
|                Total 161 records   [<] 1 2 3 4 5 [>]             |
+------------------------------------------------------------------+
```

#### Component Structure

```vue
<template>
  <div class="transactions-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="8">
        <div class="stat-item success" @click="handleStatusFilter('success')">
          <p class="stat-number">{{ stats.success }}</p>
          <p class="stat-label">Success</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item pending" @click="handleStatusFilter('pending')">
          <p class="stat-number">{{ stats.pending }}</p>
          <p class="stat-label">Pending</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item failed" @click="handleStatusFilter('failed')">
          <p class="stat-number">{{ stats.failed }}</p>
          <p class="stat-label">Failed</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="Search transaction ID / order no."
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="Status" clearable class="filter-select">
            <el-option label="All" value="" />
            <el-option label="Success" value="success" />
            <el-option label="Pending" value="pending" />
            <el-option label="Failed" value="failed" />
          </el-select>
          <el-select v-model="channelFilter" placeholder="Channel" clearable class="filter-select">
            <el-option label="All" value="" />
            <el-option label="Stripe" value="Stripe" />
            <el-option label="Alipay" value="Alipay" />
            <el-option label="WeChat" value="WeChat" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="to"
            start-placeholder="Start Date"
            end-placeholder="End Date"
            class="date-picker"
            value-format="YYYY-MM-DD"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            Export
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            Refresh
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Transactions Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="transactionList" v-loading="loading" stripe>
        <el-table-column prop="transaction_id" label="Transaction ID" min-width="180">
          <template #default="{ row }">
            <div class="transaction-id-cell">
              <span class="transaction-id">{{ row.transaction_id }}</span>
              <el-button link size="small" @click="copyTransactionId(row.transaction_id)">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="Order No." min-width="160">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewOrder(row.order_id)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="Channel" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getChannelTagType(row.channel)" effect="plain" size="small">
              {{ row.channel }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Amount" width="140" align="right">
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="transaction-amount">{{ row.currency }} {{ formatAmount(row.amount) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Fee" width="100" align="right">
          <template #default="{ row }">
            <span class="fee-amount">{{ row.currency }} {{ formatAmount(row.fee) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Status" width="100" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="statusTypeMap" />
          </template>
        </el-table-column>
        <el-table-column label="Created At" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              Details
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="warning"
              link
              size="small"
              @click="queryStatus(row)"
            >
              Query Status
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Download, Refresh, CopyDocument } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getTransactionList, queryTransactionStatus } from '@/api/payment'

const router = useRouter()

// State
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const channelFilter = ref('')
const dateRange = ref<[string, string] | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const stats = ref({
  success: 156,
  pending: 3,
  failed: 2
})

const statusTypeMap = {
  success: { type: 'success' as const, text: 'Success' },
  pending: { type: 'warning' as const, text: 'Pending' },
  failed: { type: 'danger' as const, text: 'Failed' }
}

// Mock data
const transactionList = ref([
  {
    id: 1,
    transaction_id: 'txn_abc123def456',
    order_id: 'ORD001',
    order_no: 'ORD2024031800100',
    channel: 'Stripe',
    amount: 29900,
    fee: 299,
    currency: 'CNY',
    status: 'success',
    created_at: '2024-03-18 14:32:18'
  },
  {
    id: 2,
    transaction_id: 'txn_def456ghi789',
    order_id: 'ORD002',
    order_no: 'ORD2024031700099',
    channel: 'Alipay',
    amount: 45600,
    fee: 456,
    currency: 'CNY',
    status: 'success',
    created_at: '2024-03-17 16:45:30'
  },
  {
    id: 3,
    transaction_id: 'txn_ghi789jkl012',
    order_id: 'ORD003',
    order_no: 'ORD2024031700098',
    channel: 'WeChat',
    amount: 12900,
    fee: 129,
    currency: 'CNY',
    status: 'pending',
    created_at: '2024-03-17 12:00:00'
  }
])

// Methods
const formatAmount = (amount: number) => {
  return (amount / 100).toFixed(2)
}

const getChannelTagType = (channel: string) => {
  const types: Record<string, string> = {
    'Stripe': 'primary',
    'Alipay': 'success',
    'WeChat': 'warning'
  }
  return types[channel] || 'info'
}

const copyTransactionId = (id: string) => {
  navigator.clipboard.writeText(id)
  ElMessage.success('Transaction ID copied')
}

const handleStatusFilter = (status: string) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = () => {
  ElMessage.success('Export successful')
}

const handleRefresh = () => {
  loadData()
}

const viewOrder = (orderId: string) => {
  router.push(`/orders?id=${orderId}`)
}

const viewDetail = (row: any) => {
  router.push(`/payments/transactions/${row.id}`)
}

const queryStatus = async (row: any) => {
  try {
    await queryTransactionStatus(row.id)
    ElMessage.success('Status query initiated')
    loadData()
  } catch (error) {
    ElMessage.error('Failed to query status')
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value,
      channel: channelFilter.value,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    }
    const res = await getTransactionList(params)
    transactionList.value = res.data.list
    total.value = res.data.total
    stats.value = res.data.stats
  } catch (error) {
    // Use mock data
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.transactions-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-item {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-item.success {
  border-left: 4px solid #10B981;
}

.stat-item.pending {
  border-left: 4px solid #F59E0B;
}

.stat-item.failed {
  border-left: 4px solid #EF4444;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 6px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
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
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  width: 240px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 120px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.date-picker {
  width: 260px;
}

.date-picker :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Table */
.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.transaction-id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.transaction-id {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.transaction-amount {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.fee-amount {
  color: #6B7280;
  font-size: 13px;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
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
  .filter-select,
  .date-picker {
    width: 100%;
  }

  .stat-item {
    border-radius: 14px;
    padding: 20px;
  }

  .stat-number {
    font-size: 28px;
  }
}
</style>
```

#### API Types

```typescript
// /api/payment.ts

interface Transaction {
  id: number
  transaction_id: string
  order_id: string
  order_no: string
  channel: 'Stripe' | 'Alipay' | 'WeChat'
  amount: number              // In cents
  fee: number                 // In cents
  currency: string
  status: 'success' | 'pending' | 'failed'
  created_at: string
}

interface TransactionListParams {
  page: number
  page_size: number
  status?: string
  channel?: string
  start_time?: string
  end_time?: string
  search?: string
}

interface TransactionListResponse {
  list: Transaction[]
  total: number
  stats: {
    success: number
    pending: number
    failed: number
  }
}

// API functions
export function getTransactionList(params: TransactionListParams): Promise<{ data: TransactionListResponse }>
export function getTransactionDetail(id: number): Promise<{ data: Transaction }>
export function queryTransactionStatus(id: number): Promise<void>
```

---

## Summary

### Files to Create

| File Path | Description |
|-----------|-------------|
| `src/components/payment/PaymentStatsCard.vue` | Payment statistics dashboard component |
| `src/components/payment/RefundDialog.vue` | Refund initiation dialog with idempotency key |
| `src/views/payments/transactions/index.vue` | Payment transactions list page |

### Files to Modify

| File Path | Changes |
|-----------|---------|
| `src/views/fulfillment/refunds/index.vue` | Add payment channel, channel refund ID, refund fee columns; add retry button |
| `src/views/orders/[id]/index.vue` | Add payment information section (if order detail page exists) |

### Design Tokens Summary

```
Colors:
  Primary:      #6366F1 (Indigo)
  Stripe Brand: #635BFF (Stripe official purple)
  Alipay Brand: #1677FF (Phase 2)
  WeChat Brand: #07C160 (Phase 2)
  Success:      #10B981 (Green)
  Warning:      #F59E0B (Amber)
  Danger:       #EF4444 (Red)
  Text Primary: #1E1B4B
  Text Muted:   #6B7280
  Background:   #F5F3FF

Border Radius:
  Cards:        16px
  Inputs:       12px
  Images:       10px

Typography:
  Numbers:      Fira Sans
  IDs/Codes:    Fira Code
  Body:         System font
```

### Payment Status Colors

| Status | Color | Element Plus Type |
|--------|-------|-------------------|
| Pending | Amber | warning |
| Processing | Indigo | primary |
| Success | Green | success |
| Failed | Red | danger |
| Cancelled | Gray | info |
| Refunded | Gray | info |
| Partially Refunded | Amber | warning |
| Requires Action | Amber | warning |

### Supported Currencies (Phase 1 - Stripe)

| Currency | Code | Note |
|----------|------|------|
| US Dollar | USD | Default |
| Euro | EUR | |
| British Pound | GBP | |
| Japanese Yen | JPY | |
| Singapore Dollar | SGD | |
| Hong Kong Dollar | HKD | |

**Note:** CNY not supported by Stripe. Will be supported in Phase 2 via Alipay/WeChat.

### Refund Reason Codes (Uppercase, aligned with fulfillment module)

| Code | Display Name |
|------|--------------|
| DEFECTIVE | Product Defective |
| WRONG_ITEM | Wrong Item Received |
| NOT_AS_DESCRIBED | Not As Described |
| DAMAGED | Damaged in Transit |
| NO_LONGER_NEEDED | No Longer Needed |
| LATE_DELIVERY | Late Delivery |
| OTHER | Other |

### Component Dependencies

- Element Plus: `el-card`, `el-table`, `el-dialog`, `el-form`, `el-tag`, `el-progress`, `el-descriptions`, `el-timeline`, `el-pagination`
- Custom: `StatusTag` component from `@/components/common/StatusTag.vue`
- Utils: `uuid` for generating idempotency keys