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
          meta: { title: 'dashboard.title' }
        },
        {
          path: 'products',
          name: 'products',
          component: () => import('@/views/products/index.vue'),
          meta: { title: 'products.title' }
        },
        {
          path: 'products/:id',
          name: 'product-detail',
          component: () => import('@/views/products/[id]/index.vue'),
          meta: { title: 'products.productDetail' }
        },
        {
          path: 'orders',
          name: 'orders',
          component: () => import('@/views/orders/index.vue'),
          meta: { title: 'orders.title' }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/users/index.vue'),
          meta: { title: 'users.title' }
        },
        {
          path: 'users/:id',
          name: 'user-detail',
          component: () => import('@/views/users/[id].vue'),
          meta: { title: 'users.userDetail' }
        },
        {
          path: 'admin-users',
          name: 'admin-users',
          component: () => import('@/views/admin-users/index.vue'),
          meta: { title: 'adminUsers.title' }
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('@/views/roles/index.vue'),
          meta: { title: 'roles.title' }
        },
        {
          path: 'admin-users/:id',
          name: 'admin-user-detail',
          component: () => import('@/views/admin-users/[id].vue'),
          meta: { title: 'adminUsers.adminDetail' }
        },
        {
          path: 'promotions',
          name: 'promotions',
          component: () => import('@/views/promotions/index.vue'),
          meta: { title: 'promotions.title' }
        },
        {
          path: 'promotions/:id',
          name: 'promotion-detail',
          component: () => import('@/views/promotions/[id]/index.vue'),
          meta: { title: 'promotions.promotionDetail' }
        },
        {
          path: 'shop',
          name: 'shop',
          component: () => import('@/views/shop/index.vue'),
          meta: { title: 'shop.title' }
        },
        {
          path: 'settings/markets',
          name: 'settings-markets',
          component: () => import('@/views/settings/markets/index.vue'),
          meta: { title: 'settings.markets.title' }
        },
        {
          path: 'categories',
          name: 'categories',
          component: () => import('@/views/categories/index.vue'),
          meta: { title: 'categories.title' }
        },
        {
          path: 'categories/:id',
          name: 'category-detail',
          component: () => import('@/views/categories/[id].vue'),
          meta: { title: 'categories.categoryDetail' }
        },
        {
          path: 'brands',
          name: 'brands',
          component: () => import('@/views/brands/index.vue'),
          meta: { title: 'brands.title' }
        },
        {
          path: 'brands/:id',
          name: 'brand-detail',
          component: () => import('@/views/brands/[id].vue'),
          meta: { title: 'brands.brandDetail' }
        },
        {
          path: 'inventory',
          name: 'inventory',
          component: () => import('@/views/inventory/index.vue'),
          meta: { title: 'inventory.title' }
        },
        // Fulfillment Module Routes
        {
          path: 'fulfillment/shipments',
          name: 'shipments',
          component: () => import('@/views/fulfillment/shipments/index.vue'),
          meta: { title: 'fulfillment.shipments' }
        },
        {
          path: 'fulfillment/shipments/:id',
          name: 'shipment-detail',
          component: () => import('@/views/fulfillment/shipments/[id]/index.vue'),
          meta: { title: 'fulfillment.shipmentDetail' }
        },
        {
          path: 'fulfillment/refunds',
          name: 'refunds',
          component: () => import('@/views/fulfillment/refunds/index.vue'),
          meta: { title: 'fulfillment.refunds' }
        },
        {
          path: 'fulfillment/refunds/:id',
          name: 'refund-detail',
          component: () => import('@/views/fulfillment/refunds/[id]/index.vue'),
          meta: { title: 'fulfillment.refundDetail' }
        },
        {
          path: 'fulfillment/statistics',
          name: 'fulfillment-statistics',
          component: () => import('@/views/fulfillment/statistics/index.vue'),
          meta: { title: 'fulfillment.statistics' }
        },
        {
          path: 'fulfillment/refund-reasons',
          name: 'refund-reasons',
          component: () => import('@/views/fulfillment/refund-reasons/index.vue'),
          meta: { title: 'refundReasons.title' }
        },
        {
          path: 'reviews',
          name: 'reviews',
          component: () => import('@/views/reviews/index.vue'),
          meta: { title: 'reviews.title' }
        },
        // Payment Module Routes
        {
          path: 'payments',
          name: 'payments',
          component: () => import('@/views/payments/index.vue'),
          meta: { title: 'payments.title' }
        },
        {
          path: 'payments/transactions/:id',
          name: 'payment-transaction-detail',
          component: () => import('@/views/payments/transactions/[id]/index.vue'),
          meta: { title: 'payments.transactionDetail' }
        },
        // Points Module Routes
        {
          path: 'points/dashboard',
          name: 'points-dashboard',
          component: () => import('@/views/points/dashboard/index.vue'),
          meta: { title: 'points.dashboard' }
        },
        {
          path: 'points/earn-rules',
          name: 'points-earn-rules',
          component: () => import('@/views/points/earn-rules/index.vue'),
          meta: { title: 'points.earnRules' }
        },
        {
          path: 'points/redeem-rules',
          name: 'points-redeem-rules',
          component: () => import('@/views/points/redeem-rules/index.vue'),
          meta: { title: 'points.redeemRules' }
        },
        {
          path: 'points/accounts',
          name: 'points-accounts',
          component: () => import('@/views/points/accounts/index.vue'),
          meta: { title: 'points.accounts' }
        },
        {
          path: 'points/accounts/:id',
          name: 'points-account-detail',
          component: () => import('@/views/points/accounts/[id].vue'),
          meta: { title: 'points.accountDetail' }
        },
        {
          path: 'points/transactions',
          name: 'points-transactions',
          component: () => import('@/views/points/transactions/index.vue'),
          meta: { title: 'points.transactions' }
        },
        {
          path: 'points/redemptions',
          name: 'points-redemptions',
          component: () => import('@/views/points/redemptions/index.vue'),
          meta: { title: 'points.redemptions' }
        },
        // Storefront Module Routes
        {
          path: 'storefront/themes',
          name: 'storefront-themes',
          component: () => import('@/views/storefront/themes/index.vue'),
          meta: { title: 'storefront.themes' }
        },
        {
          path: 'storefront/pages',
          name: 'storefront-pages',
          component: () => import('@/views/storefront/pages/index.vue'),
          meta: { title: 'storefront.pages' }
        },
        {
          path: 'storefront/pages/:id/edit',
          name: 'storefront-page-edit',
          component: () => import('@/views/storefront/pages/[id]/edit.vue'),
          meta: { title: 'storefront.pageEdit' }
        },
        {
          path: 'storefront/seo',
          name: 'storefront-seo',
          component: () => import('@/views/storefront/seo/index.vue'),
          meta: { title: 'storefront.seo' }
        },
        // Shipping Module Routes
        {
          path: 'shipping',
          name: 'shipping',
          component: () => import('@/views/shipping/index.vue'),
          meta: { title: 'shipping.title' }
        },
        {
          path: 'shipping/calculator',
          name: 'shipping-calculator',
          component: () => import('@/views/shipping/calculator/index.vue'),
          meta: { title: 'shipping.calculator' }
        },
        {
          path: 'shipping/:id',
          name: 'shipping-detail',
          component: () => import('@/views/shipping/[id]/index.vue'),
          meta: { title: 'shipping.templateDetail' }
        },
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  if (to.path !== '/login' && !userStore.token) {
    next('/login')
  } else {
    next()
  }
})

export default router
