<template>
  <div class="cart">
    <h1>购物车</h1>
    <div v-if="cartItems.length === 0" class="empty">
      <p>购物车是空的</p>
      <button @click="goToShop">去购物</button>
    </div>
    <div v-else class="cart-content">
      <div v-for="item in cartItems" :key="item.id" class="cart-item">
        <div class="item-image"></div>
        <div class="item-info">
          <h3>{{ item.name }}</h3>
          <p class="price">¥{{ item.price.toFixed(2) }}</p>
        </div>
        <div class="item-actions">
          <input type="number" v-model="item.quantity" min="1" />
          <button @click="removeItem(item.id)">删除</button>
        </div>
      </div>
      <div class="cart-footer">
        <p class="total">合计: ¥{{ total.toFixed(2) }}</p>
        <button class="checkout" @click="checkout">去结算</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const cartItems = ref([
  { id: 1, name: '商品1', price: 99, quantity: 2 },
  { id: 2, name: '商品2', price: 199, quantity: 1 }
])

const total = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const removeItem = (id: number) => {
  cartItems.value = cartItems.value.filter(item => item.id !== id)
}

const checkout = () => {
  router.push('/checkout')
}

const goToShop = () => {
  router.push('/products')
}
</script>

<style scoped>
.cart {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.empty {
  text-align: center;
  padding: 60px;
}

.cart-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #eee;
  gap: 16px;
}

.item-image {
  width: 80px;
  height: 80px;
  background: #f5f5f5;
  border-radius: 4px;
}

.item-info {
  flex: 1;
}

.price {
  color: #e53935;
  font-weight: bold;
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

input {
  width: 60px;
  padding: 4px;
}

.cart-footer {
  margin-top: 20px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total {
  font-size: 20px;
  font-weight: bold;
}

.checkout {
  background: #e53935;
  color: white;
  border: none;
  padding: 12px 32px;
  border-radius: 4px;
  cursor: pointer;
}
</style>
