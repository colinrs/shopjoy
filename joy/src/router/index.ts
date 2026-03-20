import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/home/index.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/index.vue')
    },
    {
      path: '/products',
      name: 'products',
      component: () => import('@/views/products/index.vue')
    },
    {
      path: '/products/:id',
      name: 'product-detail',
      component: () => import('@/views/products/detail.vue')
    },
    {
      path: '/cart',
      name: 'cart',
      component: () => import('@/views/cart/index.vue')
    },
    {
      path: '/checkout',
      name: 'checkout',
      component: () => import('@/views/checkout/index.vue')
    },
    {
      path: '/orders',
      name: 'orders',
      component: () => import('@/views/orders/index.vue')
    },
    {
      path: '/user',
      name: 'user',
      component: () => import('@/views/user/index.vue')
    }
  ]
})

export default router
