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
          path: 'products/:id',
          name: 'product-detail',
          component: () => import('@/views/products/[id]/index.vue'),
          meta: { title: '商品详情' }
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
          path: 'users/:id',
          name: 'user-detail',
          component: () => import('@/views/users/[id].vue'),
          meta: { title: '用户详情' }
        },
        {
          path: 'admin-users',
          name: 'admin-users',
          component: () => import('@/views/admin-users/index.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'admin-users/:id',
          name: 'admin-user-detail',
          component: () => import('@/views/admin-users/[id].vue'),
          meta: { title: '管理员详情' }
        },
        {
          path: 'promotions',
          name: 'promotions',
          component: () => import('@/views/promotions/index.vue'),
          meta: { title: '促销管理' }
        },
        {
          path: 'promotions/:id',
          name: 'promotion-detail',
          component: () => import('@/views/promotions/[id]/index.vue'),
          meta: { title: '促销详情' }
        },
        {
          path: 'shop',
          name: 'shop',
          component: () => import('@/views/shop/index.vue'),
          meta: { title: '店铺设置' }
        },
        {
          path: 'settings/markets',
          name: 'settings-markets',
          component: () => import('@/views/settings/markets/index.vue'),
          meta: { title: '市场管理' }
        },
        {
          path: 'categories',
          name: 'categories',
          component: () => import('@/views/categories/index.vue'),
          meta: { title: '分类管理' }
        },
        {
          path: 'brands',
          name: 'brands',
          component: () => import('@/views/brands/index.vue'),
          meta: { title: '品牌管理' }
        },
        {
          path: 'inventory',
          name: 'inventory',
          component: () => import('@/views/inventory/index.vue'),
          meta: { title: '库存管理' }
        },
        // Fulfillment Module Routes
        {
          path: 'fulfillment/shipments',
          name: 'shipments',
          component: () => import('@/views/fulfillment/shipments/index.vue'),
          meta: { title: '发货管理' }
        },
        {
          path: 'fulfillment/shipments/:id',
          name: 'shipment-detail',
          component: () => import('@/views/fulfillment/shipments/[id]/index.vue'),
          meta: { title: '发货详情' }
        },
        {
          path: 'fulfillment/refunds',
          name: 'refunds',
          component: () => import('@/views/fulfillment/refunds/index.vue'),
          meta: { title: '退款管理' }
        },
        {
          path: 'fulfillment/refunds/:id',
          name: 'refund-detail',
          component: () => import('@/views/fulfillment/refunds/[id]/index.vue'),
          meta: { title: '退款详情' }
        },
        {
          path: 'fulfillment/statistics',
          name: 'fulfillment-statistics',
          component: () => import('@/views/fulfillment/statistics/index.vue'),
          meta: { title: '售后统计' }
        },
        {
          path: 'reviews',
          name: 'reviews',
          component: () => import('@/views/reviews/index.vue'),
          meta: { title: '评价管理' }
        },
        // Payment Module Routes
        {
          path: 'payments',
          name: 'payments',
          component: () => import('@/views/payments/index.vue'),
          meta: { title: '支付管理' }
        },
        // Points Module Routes
        {
          path: 'points/dashboard',
          name: 'points-dashboard',
          component: () => import('@/views/points/dashboard/index.vue'),
          meta: { title: '积分概览' }
        },
        {
          path: 'points/earn-rules',
          name: 'points-earn-rules',
          component: () => import('@/views/points/earn-rules/index.vue'),
          meta: { title: '获取规则' }
        },
        {
          path: 'points/redeem-rules',
          name: 'points-redeem-rules',
          component: () => import('@/views/points/redeem-rules/index.vue'),
          meta: { title: '兑换规则' }
        },
        {
          path: 'points/accounts',
          name: 'points-accounts',
          component: () => import('@/views/points/accounts/index.vue'),
          meta: { title: '积分账户' }
        },
        {
          path: 'points/accounts/:id',
          name: 'points-account-detail',
          component: () => import('@/views/points/accounts/[id].vue'),
          meta: { title: '账户详情' }
        },
        {
          path: 'points/transactions',
          name: 'points-transactions',
          component: () => import('@/views/points/transactions/index.vue'),
          meta: { title: '交易记录' }
        },
        {
          path: 'points/redemptions',
          name: 'points-redemptions',
          component: () => import('@/views/points/redemptions/index.vue'),
          meta: { title: '兑换记录' }
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
