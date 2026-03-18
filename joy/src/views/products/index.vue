<template>
  <div class="products">
    <h1>商品列表</h1>
    <div class="product-grid">
      <div v-for="i in 12" :key="i" class="product-card" @click="goToDetail(i)">
        <div class="product-image"></div>
        <h3>商品 {{ i }}</h3>
        <p class="desc">商品描述信息</p>
        <p class="price">¥{{ 50 + i * 10 }}.00</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToDetail = (id: number) => {
  router.push(`/products/${id}`)
}
</script>

<style scoped>
.products {
  padding: 20px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.product-card {
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: box-shadow 0.3s;
}

.product-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.product-image {
  width: 100%;
  height: 150px;
  background: #f5f5f5;
  border-radius: 4px;
  margin-bottom: 12px;
}

.desc {
  color: #666;
  font-size: 14px;
}

.price {
  color: #e53935;
  font-size: 18px;
  font-weight: bold;
  margin-top: 8px;
}
</style>
