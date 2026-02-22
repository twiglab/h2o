import { createRouter, createWebHistory } from 'vue-router'
import { constantRoutes } from './routes'
import { setupRouterGuards } from './guards'

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 })
})

// 设置路由守卫
setupRouterGuards(router)

export default router

// 重置路由（移除动态添加的路由）
export function resetRouter(): void {
  // 获取所有路由名称
  const routeNames = router.getRoutes().map(route => route.name)

  // 获取静态路由名称
  const staticRouteNames = new Set(
    constantRoutes.filter(r => r.name).map(r => r.name)
  )

  // 移除所有非静态路由
  routeNames.forEach(name => {
    if (name && !staticRouteNames.has(name as string)) {
      router.removeRoute(name)
    }
  })

  // 如果占位符路由不存在，重新添加（用于下次登录）
  if (!router.hasRoute('AsyncRoutePlaceholder')) {
    router.addRoute({
      path: '/:pathMatch(.*)*',
      name: 'AsyncRoutePlaceholder',
      component: () => import('@/views/redirect/index.vue'),
      meta: { hidden: true }
    })
  }
}
