<template>
  <div class="user">
    <div class="user-header">
      <div class="avatar"></div>
      <div class="info">
        <h2>用户名</h2>
        <p>13800138000</p>
      </div>
    </div>
    
    <div class="menu">
      <div class="menu-item" @click="goToOrders">
        <span>我的订单</span>
        <span>&gt;</span>
      </div>
      <div class="menu-item" @click="goToCart">
        <span>购物车</span>
        <span>&gt;</span>
      </div>
      <div class="menu-item">
        <span>收货地址</span>
        <span>&gt;</span>
      </div>
      <div class="menu-item">
        <span>设置</span>
        <span>&gt;</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToOrders = () => {
  router.push('/orders')
}

const goToCart = () => {
  router.push('/cart')
}
</script>

<style scoped>
.user {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
  margin-bottom: 20px;
}

.avatar {
  width: 60px;
  height: 60px;
  background: white;
  border-radius: 50%;
}

.menu {
  border: 1px solid #eee;
  border-radius: 8px;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:hover {
  background: #f9f9f9;
}
</style>
