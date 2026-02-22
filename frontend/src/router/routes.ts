import type { RouteRecordRaw } from 'vue-router'

// 布局组件
const Layout = () => import('@/components/layout/index.vue')

// 静态路由（不需要权限）
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { hidden: true }
  },
  {
    path: '/redirect/:path(.*)',
    component: () => import('@/views/redirect/index.vue'),
    meta: { hidden: true }
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { hidden: true }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { hidden: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'AsyncRoutePlaceholder',
    component: () => import('@/views/redirect/index.vue'),
    meta: { hidden: true }
  }
]

// 动态路由（需要权限）
export const asyncRoutes: RouteRecordRaw[] = [
  // 1. 工作台
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台', icon: 'HomeFilled', affix: true }
      }
    ]
  },

  // 2. 客户管理
  {
    path: '/customer',
    component: Layout,
    redirect: '/customer/merchant',
    meta: { title: '客户管理', icon: 'User' },
    children: [
      {
        path: 'merchant',
        name: 'MerchantList',
        component: () => import('@/views/merchant/index.vue'),
        meta: { title: '商户列表', icon: 'Briefcase', permission: 'merchant:info:list' }
      },
      {
        path: 'merchant/:id',
        name: 'MerchantDetail',
        component: () => import('@/views/merchant/detail.vue'),
        meta: { title: '商户详情', hidden: true, activeMenu: '/customer/merchant' }
      },
      {
        path: 'shop',
        name: 'ShopList',
        component: () => import('@/views/shop/index.vue'),
        meta: { title: '店铺列表', icon: 'Shop', permission: 'merchant:shop:list' }
      },
      {
        path: 'shop/:id',
        name: 'ShopDetail',
        component: () => import('@/views/shop/detail.vue'),
        meta: { title: '店铺详情', hidden: true, activeMenu: '/customer/shop' }
      },
      {
        path: 'account',
        name: 'AccountList',
        component: () => import('@/views/account/index.vue'),
        meta: { title: '账户列表', icon: 'Wallet', permission: 'finance:account:list' }
      },
      {
        path: 'account/:id',
        name: 'AccountDetail',
        component: () => import('@/views/account/detail.vue'),
        meta: { title: '账户详情', hidden: true, activeMenu: '/customer/account' }
      },
      {
        path: 'recharge',
        name: 'Recharge',
        component: () => import('@/views/account/recharge/index.vue'),
        meta: { title: '充值记录', icon: 'Money', permission: 'finance:recharge:list' }
      }
    ]
  },

  // 3. 电表管理
  {
    path: '/electric',
    component: Layout,
    redirect: '/electric/meter',
    meta: { title: '电表管理', icon: 'Lightning' },
    children: [
      {
        path: 'meter',
        name: 'ElectricMeter',
        component: () => import('@/views/meter/electric/index.vue'),
        meta: { title: '电表列表', icon: 'Odometer', permission: 'meter:electric:list' }
      },
      {
        path: 'reading',
        name: 'ElectricReading',
        component: () => import('@/views/meter/reading/electric.vue'),
        meta: { title: '电表抄表', icon: 'Document', permission: 'meter:reading:list' }
      },
      {
        path: 'rate',
        name: 'ElectricRate',
        component: () => import('@/views/rate/electric/index.vue'),
        meta: { title: '电费费率', icon: 'PriceTag', permission: 'rate:electric:list' }
      },
      {
        path: 'deduction',
        name: 'ElectricDeduction',
        component: () => import('@/views/account/deduction/electric.vue'),
        meta: { title: '电费扣费', icon: 'Tickets', permission: 'finance:deduction:list' }
      }
    ]
  },

  // 4. 水表管理
  {
    path: '/water',
    component: Layout,
    redirect: '/water/meter',
    meta: { title: '水表管理', icon: 'Bowl' },
    children: [
      {
        path: 'meter',
        name: 'WaterMeter',
        component: () => import('@/views/meter/water/index.vue'),
        meta: { title: '水表列表', icon: 'Odometer', permission: 'meter:water:list' }
      },
      {
        path: 'reading',
        name: 'WaterReading',
        component: () => import('@/views/meter/reading/water.vue'),
        meta: { title: '水表抄表', icon: 'Document', permission: 'meter:reading:list' }
      },
      {
        path: 'rate',
        name: 'WaterRate',
        component: () => import('@/views/rate/water/index.vue'),
        meta: { title: '水费费率', icon: 'PriceTag', permission: 'rate:water:list' }
      },
      {
        path: 'deduction',
        name: 'WaterDeduction',
        component: () => import('@/views/account/deduction/water.vue'),
        meta: { title: '水费扣费', icon: 'Tickets', permission: 'finance:deduction:list' }
      }
    ]
  },

  // 5. 统计报表
  {
    path: '/report',
    component: Layout,
    redirect: '/report/consumption',
    meta: { title: '统计报表', icon: 'DataAnalysis' },
    children: [
      {
        path: 'consumption',
        name: 'ConsumptionReport',
        component: () => import('@/views/report/consumption/index.vue'),
        meta: { title: '用量统计', icon: 'TrendCharts', permission: 'report:consumption:view' }
      },
      {
        path: 'revenue',
        name: 'RevenueReport',
        component: () => import('@/views/report/revenue/index.vue'),
        meta: { title: '收入统计', icon: 'Coin', permission: 'report:billing:view' }
      }
    ]
  },

  // 6. 系统管理
  {
    path: '/system',
    component: Layout,
    redirect: '/system/user',
    meta: { title: '系统管理', icon: 'Setting' },
    children: [
      {
        path: 'user',
        name: 'SystemUser',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'User', permission: 'system:user:list' }
      },
      {
        path: 'dept',
        name: 'SystemDept',
        component: () => import('@/views/system/dept/index.vue'),
        meta: { title: '部门管理', icon: 'OfficeBuilding', permission: 'system:dept:list' }
      },
      {
        path: 'role',
        name: 'SystemRole',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', icon: 'UserFilled', permission: 'system:role:list' }
      },
      {
        path: 'permission',
        name: 'SystemPermission',
        component: () => import('@/views/system/permission/index.vue'),
        meta: { title: '权限管理', icon: 'Key', permission: 'system:permission:list' }
      },
      {
        path: 'log',
        name: 'SystemLog',
        component: () => import('@/views/system/log/index.vue'),
        meta: { title: '操作日志', icon: 'Document', permission: 'system:log:list' }
      }
    ]
  },

  // 个人中心
  {
    path: '/profile',
    component: Layout,
    meta: { hidden: true },
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心' }
      }
    ]
  },

  // 404 重定向必须放最后
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    meta: { hidden: true }
  }
]
