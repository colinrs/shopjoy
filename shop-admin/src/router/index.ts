import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/index.vue')
    },
    {
      path: '/',
      name: 'layout',
      component: () => import('@/layouts/MainLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'products',
          name: 'products',
          component: () => import('@/views/products/index.vue'),
          meta: { title: '商品管理' }
        },
        {
          path: 'orders',
          name: 'orders',
          component: () => import('@/views/orders/index.vue'),
          meta: { title: '订单管理' }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/users/index.vue'),
          meta: { title: '顾客管理' }
        },
        {
          path: 'admin-users',
          name: 'admin-users',
          component: () => import('@/views/admin-users/index.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'promotions',
          name: 'promotions',
          component: () => import('@/views/promotions/index.vue'),
          meta: { title: '促销管理' }
        },
        {
          path: 'shop',
          name: 'shop',
          component: () => import('@/views/shop/index.vue'),
          meta: { title: '店铺设置' }
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.path !== '/login' && !userStore.token) {
    next('/login')
  } else {
    next()
  }
})

export default router
