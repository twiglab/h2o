import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { constantRoutes, asyncRoutes } from '@/router/routes'
import type { MenuItem } from '@/types/router'

export const usePermissionStore = defineStore('permission', () => {
  // 状态
  const routes = ref<RouteRecordRaw[]>([])
  const dynamicRoutes = ref<RouteRecordRaw[]>([])
  const sidebarMenus = ref<MenuItem[]>([])

  // 过滤路由
  function filterAsyncRoutes(routes: RouteRecordRaw[], permissions: string[]): RouteRecordRaw[] {
    const result: RouteRecordRaw[] = []

    routes.forEach(route => {
      const tmp = { ...route }

      // 检查权限
      if (hasPermission(tmp, permissions)) {
        // 递归过滤子路由
        if (tmp.children) {
          tmp.children = filterAsyncRoutes(tmp.children, permissions)
        }
        result.push(tmp)
      }
    })

    return result
  }

  // 检查路由权限
  function hasPermission(route: RouteRecordRaw, permissions: string[]): boolean {
    if (route.meta?.permission) {
      const perm = route.meta.permission
      if (Array.isArray(perm)) {
        return perm.some(p => permissions.includes(p))
      }
      return permissions.includes(perm)
    }
    return true
  }

  // 生成路由
  function generateRoutes(permissions: string[], isSuperAdmin: boolean): RouteRecordRaw[] {
    let accessedRoutes: RouteRecordRaw[]

    if (isSuperAdmin) {
      // 超级管理员拥有所有路由
      accessedRoutes = asyncRoutes
    } else {
      // 根据权限过滤路由
      accessedRoutes = filterAsyncRoutes(asyncRoutes, permissions)
    }

    routes.value = constantRoutes.concat(accessedRoutes)
    dynamicRoutes.value = accessedRoutes
    sidebarMenus.value = generateMenus(routes.value)

    return accessedRoutes
  }

  // 生成菜单
  function generateMenus(routes: RouteRecordRaw[]): MenuItem[] {
    const menus: MenuItem[] = []

    routes.forEach(route => {
      // 跳过隐藏的路由
      if (route.meta?.hidden) return

      const menu: MenuItem = {
        path: route.path,
        title: route.meta?.title || '',
        icon: route.meta?.icon
      }

      // 处理子路由
      if (route.children && route.children.length > 0) {
        const children = generateMenus(route.children)
        if (children.length > 0) {
          menu.children = children
        }
      }

      menus.push(menu)
    })

    return menus
  }

  // 重置
  function resetRoutes(): void {
    routes.value = []
    dynamicRoutes.value = []
    sidebarMenus.value = []
  }

  return {
    routes,
    dynamicRoutes,
    sidebarMenus,
    generateRoutes,
    resetRoutes
  }
})
