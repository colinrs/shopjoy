<template>
  <div class="orders-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.pending }}</p>
          <p class="stat-label">待处理</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.paid }}</p>
          <p class="stat-label">已支付</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.shipped }}</p>
          <p class="stat-label">已发货</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item">
          <p class="stat-number">{{ orderStats.completed }}</p>
          <p class="stat-label">已完成</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="搜索订单号/买家"
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="订单状态" clearable class="filter-select">
            <el-option label="全部" value="" />
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="待发货" :value="2" />
            <el-option label="已发货" :value="3" />
            <el-option label="已完成" :value="4" />
            <el-option label="已取消" :value="5" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            class="date-picker"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>导出
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Orders Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="orderList" v-loading="loading" stripe>
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="order-detail">
              <el-row :gutter="20">
                <el-col :span="12">
                  <h4>商品信息</h4>
                  <div v-for="item in row.items" :key="item.id" class="order-item">
                    <el-image :src="item.image" class="item-image" fit="cover">
                      <template #error>
                        <div class="image-placeholder">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <div class="item-info">
                      <p class="item-name">{{ item.name }}</p>
                      <p class="item-price">¥{{ (item.price / 100).toFixed(2) }} x {{ item.quantity }}</p>
                    </div>
                  </div>
                </el-col>
                <el-col :span="12">
                  <h4>收货信息</h4>
                  <p><strong>收货人:</strong> {{ row.receiver_name }}</p>
                  <p><strong>电话:</strong> {{ row.receiver_phone }}</p>
                  <p><strong>地址:</strong> {{ row.receiver_address }}</p>
                  <h4 style="margin-top: 20px">支付信息</h4>
                  <p><strong>支付方式:</strong> {{ row.payment_method }}</p>
                  <p><strong>支付时间:</strong> {{ row.paid_at || '-' }}</p>
                </el-col>
              </el-row>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="订单号" min-width="160">
          <template #default="{ row }">
            <div class="order-no-cell">
              <span class="order-no">{{ row.order_no }}</span>
              <el-tag v-if="row.is_urgent" size="small" type="danger" effect="dark">急</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="商品" min-width="200">
          <template #default="{ row }">
            <div class="goods-preview">
              <el-image
                v-for="(item, idx) in row.items.slice(0, 3)"
                :key="idx"
                :src="item.image"
                class="goods-thumb"
                fit="cover"
              >
                <template #error>
                  <div class="thumb-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <span v-if="row.items.length > 3" class="more-goods">+{{ row.items.length - 3 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="买家信息" min-width="150">
          <template #default="{ row }">
            <div class="buyer-info">
              <p class="buyer-name">{{ row.buyer_name }}</p>
              <p class="buyer-phone">{{ row.buyer_phone }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="金额" width="140" align="right">
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="total-amount">¥{{ (row.pay_amount / 100).toFixed(2) }}</p>
              <p class="item-count">共 {{ row.item_count }} 件</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button 
              v-if="row.status === 1" 
              type="primary" 
              size="small"
              @click="handleShip(row)"
            >
              发货
            </el-button>
            <el-button 
              v-if="row.status === 0" 
              type="warning" 
              size="small"
              @click="handleRemind(row)"
            >
              催付
            </el-button>
            <el-button type="primary" link size="small" @click="handleDetail(row)">
              详情
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

    <!-- Ship Dialog -->
    <el-dialog v-model="shipDialogVisible" title="订单发货" width="500px">
      <el-form :model="shipForm" label-width="100px">
        <el-form-item label="物流公司">
          <el-select v-model="shipForm.company" placeholder="请选择物流公司" style="width: 100%">
            <el-option label="顺丰速运" value="sf" />
            <el-option label="中通快递" value="zt" />
            <el-option label="圆通速递" value="yt" />
            <el-option label="韵达快递" value="yd" />
            <el-option label="申通快递" value="st" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号">
          <el-input v-model="shipForm.tracking_no" placeholder="请输入物流单号" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="shipForm.remark" type="textarea" rows="3" placeholder="选填" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Download, Refresh, Picture } from '@element-plus/icons-vue'

const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const dateRange = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(100)
const shipDialogVisible = ref(false)
const currentOrder = ref<any>(null)

const shipForm = reactive({
  company: '',
  tracking_no: '',
  remark: ''
})

const orderStats = ref({
  pending: 12,
  paid: 8,
  shipped: 23,
  completed: 156
})

const orderList = ref([
  {
    id: 1,
    order_no: 'ORD2024031800100',
    buyer_name: '张三',
    buyer_phone: '138****8001',
    pay_amount: 29900,
    item_count: 2,
    status: 1,
    created_at: '2024-03-18 14:30:25',
    is_urgent: true,
    payment_method: '微信支付',
    paid_at: '2024-03-18 14:32:18',
    receiver_name: '张三',
    receiver_phone: '13800138001',
    receiver_address: '广东省深圳市南山区科技园南区A栋1201',
    items: [
      { id: 1, name: '无线蓝牙耳机 Pro', price: 29900, quantity: 1, image: '' }
    ]
  },
  {
    id: 2,
    order_no: 'ORD2024031800099',
    buyer_name: '李四',
    buyer_phone: '139****9002',
    pay_amount: 45600,
    item_count: 3,
    status: 2,
    created_at: '2024-03-18 13:15:36',
    is_urgent: false,
    payment_method: '支付宝',
    paid_at: '2024-03-18 13:16:45',
    receiver_name: '李四',
    receiver_phone: '13900139002',
    receiver_address: '北京市朝阳区建国路88号',
    items: [
      { id: 2, name: '智能手表 Series 7', price: 199900, quantity: 1, image: '' },
      { id: 3, name: '便携充电宝', price: 12900, quantity: 2, image: '' }
    ]
  },
  {
    id: 3,
    order_no: 'ORD2024031800098',
    buyer_name: '王五',
    buyer_phone: '137****7003',
    pay_amount: 12900,
    item_count: 1,
    status: 0,
    created_at: '2024-03-18 12:00:00',
    is_urgent: false,
    items: [
      { id: 4, name: '便携充电宝 20000mAh', price: 12900, quantity: 1, image: '' }
    ]
  },
  {
    id: 4,
    order_no: 'ORD2024031800097',
    buyer_name: '赵六',
    buyer_phone: '136****6004',
    pay_amount: 59900,
    item_count: 1,
    status: 3,
    created_at: '2024-03-18 10:30:15',
    is_urgent: false,
    payment_method: '微信支付',
    paid_at: '2024-03-18 10:32:00',
    receiver_name: '赵六',
    receiver_phone: '13600136004',
    receiver_address: '上海市浦东新区陆家嘴环路1000号',
    items: [
      { id: 5, name: '机械键盘 RGB', price: 45900, quantity: 1, image: '' }
    ]
  },
  {
    id: 5,
    order_no: 'ORD2024031800096',
    buyer_name: '钱七',
    buyer_phone: '135****5005',
    pay_amount: 159900,
    item_count: 1,
    status: 4,
    created_at: '2024-03-17 16:45:30',
    is_urgent: false,
    payment_method: '支付宝',
    paid_at: '2024-03-17 16:47:12',
    receiver_name: '钱七',
    receiver_phone: '13500135005',
    receiver_address: '浙江省杭州市西湖区文三路200号',
    items: [
      { id: 6, name: '4K 高清显示器', price: 159900, quantity: 1, image: '' }
    ]
  }
])

const getStatusType = (status: number) => {
  const types: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: number) => {
  const texts: Record<number, string> = {
    0: '待支付',
    1: '已支付',
    2: '待发货',
    3: '已发货',
    4: '已完成',
    5: '已取消'
  }
  return texts[status] || '未知'
}

const handleExport = () => {
  ElMessage.success('订单导出成功')
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
}

const handleShip = (row: any) => {
  currentOrder.value = row
  shipDialogVisible.value = true
}

const confirmShip = () => {
  shipDialogVisible.value = false
  ElMessage.success('发货成功')
}

const handleRemind = (row: any) => {
  ElMessage.success(`已向 ${row.buyer_name} 发送催付提醒`)
}

const handleDetail = (row: any) => {
  ElMessage.info('查看订单详情: ' + row.order_no)
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

onMounted(() => {
  // Load orders
})
</script>

<style scoped>
.orders-page {
  padding: 0;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-item {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #059669;
  margin: 0 0 4px 0;
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
  width: 220px;
}

.filter-select {
  width: 140px;
}

.date-picker {
  width: 260px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Table */
.table-card {
  margin-bottom: 20px;
}

/* Order Detail Expand */
.order-detail {
  padding: 20px;
  background: #F9FAFB;
  border-radius: 8px;
}

.order-detail h4 {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 12px 0;
}

.order-detail p {
  font-size: 13px;
  color: #4B5563;
  margin: 0 0 8px 0;
}

.order-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #E5E7EB;
}

.order-item:last-child {
  border-bottom: none;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #E5E7EB;
  color: #9CA3AF;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  color: #111827;
  margin: 0 0 4px 0;
}

.item-price {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

/* Order No Cell */
.order-no-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.order-no {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #059669;
  font-weight: 500;
}

/* Goods Preview */
.goods-preview {
  display: flex;
  align-items: center;
  gap: 4px;
}

.goods-thumb {
  width: 50px;
  height: 50px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #E5E7EB;
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #9CA3AF;
}

.more-goods {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  border-radius: 6px;
  font-size: 12px;
  color: #6B7280;
  font-weight: 500;
}

/* Buyer Info */
.buyer-info {
  line-height: 1.5;
}

.buyer-name {
  font-weight: 500;
  color: #111827;
  margin: 0;
}

.buyer-phone {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

/* Amount Cell */
.amount-cell {
  text-align: right;
}

.total-amount {
  font-size: 16px;
  font-weight: 600;
  color: #EF4444;
  margin: 0;
}

.item-count {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

/* Time Text */
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
  border-top: 1px solid #E5E7EB;
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
  
  .goods-preview {
    flex-wrap: wrap;
  }
}
</style>
