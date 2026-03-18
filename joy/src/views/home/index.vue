<template>
  <div class="home">
    <div class="banner">
      <h1>ShopJoy 商城</h1>
      <p>发现精彩商品</p>
    </div>
    
    <div class="products">
      <h2>热门商品</h2>
      <div class="product-grid">
        <div v-for="i in 8" :key="i" class="product-card">
          <div class="product-image"></div>
          <h3>商品 {{ i }}</h3>
          <p class="price">¥99.00</p>
          <button @click="addToCart(i)">加入购物车</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const addToCart = (id: number) => {
  alert('已加入购物车: 商品 ' + id)
}
</script>

<style scoped>
.home {
  padding: 20px;
}

.banner {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  margin-bottom: 40px;
}

.banner h1 {
  margin: 0 0 10px;
  font-size: 48px;
}

.products h2 {
  margin-bottom: 20px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.product-card {
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}

.product-image {
  width: 100%;
  height: 200px;
  background: #f5f5f5;
  border-radius: 4px;
  margin-bottom: 12px;
}

.price {
  color: #e53935;
  font-size: 20px;
  font-weight: bold;
}

button {
  background: #667eea;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background: #5568d3;
}
</style>
