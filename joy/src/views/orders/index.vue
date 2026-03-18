<template>
  <div class="orders">
    <h1>我的订单</h1>
    <div v-for="order in orders" :key="order.id" class="order-card">
      <div class="order-header">
        <span>订单号: {{ order.orderNo }}</span>
        <span :class="['status', order.status]">{{ order.statusText }}</span>
      </div>
      <div class="order-items">
        <div v-for="item in order.items" :key="item.id" class="item">
          <div class="item-image"></div>
          <div class="item-info">
            <h4>{{ item.name }}</h4>
            <p>¥{{ item.price.toFixed(2) }} x {{ item.quantity }}</p>
          </div>
        </div>
      </div>
      <div class="order-footer">
        <p>合计: ¥{{ order.total.toFixed(2) }}</p>
        <button v-if="order.status === 'pending'" @click="payOrder(order.id)">立即支付</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const orders = ref([
  {
    id: 1,
    orderNo: '202403180001',
    status: 'pending',
    statusText: '待支付',
    total: 397,
    items: [
      { id: 1, name: '商品1', price: 99, quantity: 2 },
      { id: 2, name: '商品2', price: 199, quantity: 1 }
    ]
  },
  {
    id: 2,
    orderNo: '202403170001',
    status: 'completed',
    statusText: '已完成',
    total: 199,
    items: [
      { id: 3, name: '商品3', price: 199, quantity: 1 }
    ]
  }
])

const payOrder = (id: number) => {
  alert('支付订单: ' + id)
}
</script>

<style scoped>
.orders {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.order-card {
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 20px;
  padding: 16px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid #f5f5f5;
}

.status {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 14px;
}

.status.pending {
  background: #fff3e0;
  color: #ff9800;
}

.status.completed {
  background: #e8f5e9;
  color: #4caf50;
}

.item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
}

.item-image {
  width: 60px;
  height: 60px;
  background: #f5f5f5;
  border-radius: 4px;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f5f5f5;
}

button {
  background: #e53935;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
}
</style>
