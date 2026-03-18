<template>
  <div class="product-detail">
    <div class="product-info">
      <div class="product-image"></div>
      <div class="product-meta">
        <h1>商品详情</h1>
        <p class="desc">这是商品的详细描述...</p>
        <p class="price">¥199.00</p>
        <div class="actions">
          <button class="add-cart" @click="addToCart">加入购物车</button>
          <button class="buy-now" @click="buyNow">立即购买</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const addToCart = () => {
  alert('已加入购物车')
}

const buyNow = () => {
  router.push('/checkout')
}
</script>

<style scoped>
.product-detail {
  padding: 20px;
}

.product-info {
  display: flex;
  gap: 40px;
  max-width: 1000px;
  margin: 0 auto;
}

.product-image {
  width: 400px;
  height: 400px;
  background: #f5f5f5;
  border-radius: 8px;
}

.product-meta {
  flex: 1;
}

.desc {
  color: #666;
  margin: 20px 0;
}

.price {
  font-size: 32px;
  color: #e53935;
  font-weight: bold;
}

.actions {
  margin-top: 30px;
  display: flex;
  gap: 16px;
}

button {
  padding: 12px 32px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.add-cart {
  background: #ff9800;
  color: white;
}

.buy-now {
  background: #e53935;
  color: white;
}
</style>
