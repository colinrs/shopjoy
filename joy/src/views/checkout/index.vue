<template>
  <div class="checkout">
    <h1>确认订单</h1>
    
    <div class="section">
      <h2>收货地址</h2>
      <div class="address">
        <p>张三 13800138000</p>
        <p>广东省深圳市南山区科技园</p>
      </div>
    </div>
    
    <div class="section">
      <h2>商品清单</h2>
      <div v-for="item in orderItems" :key="item.id" class="order-item">
        <div class="item-image"></div>
        <div class="item-info">
          <h4>{{ item.name }}</h4>
          <p>¥{{ item.price.toFixed(2) }} x {{ item.quantity }}</p>
        </div>
        <p class="item-total">¥{{ (item.price * item.quantity).toFixed(2) }}</p>
      </div>
    </div>
    
    <div class="section">
      <h2>订单总计</h2>
      <div class="summary">
        <p>商品金额: ¥{{ subtotal.toFixed(2) }}</p>
        <p>运费: ¥{{ shipping.toFixed(2) }}</p>
        <p class="total">实付金额: ¥{{ total.toFixed(2) }}</p>
      </div>
    </div>
    
    <button class="submit" @click="submitOrder">提交订单</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const orderItems = ref([
  { id: 1, name: '商品1', price: 99, quantity: 2 },
  { id: 2, name: '商品2', price: 199, quantity: 1 }
])

const shipping = ref(10)

const subtotal = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const total = computed(() => subtotal.value + shipping.value)

const submitOrder = () => {
  alert('订单提交成功！')
  router.push('/orders')
}
</script>

<style scoped>
.checkout {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
}

.section h2 {
  margin-bottom: 16px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
}

.item-image {
  width: 60px;
  height: 60px;
  background: #f5f5f5;
  border-radius: 4px;
}

.item-info {
  flex: 1;
}

.item-total {
  font-weight: bold;
}

.summary {
  text-align: right;
}

.total {
  font-size: 24px;
  color: #e53935;
  font-weight: bold;
  margin-top: 16px;
}

.submit {
  width: 100%;
  background: #e53935;
  color: white;
  border: none;
  padding: 16px;
  border-radius: 4px;
  font-size: 18px;
  cursor: pointer;
}
</style>
