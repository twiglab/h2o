import type { Router } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useUserStore } from '@/stores/modules/user'
import { usePermissionStore } from '@/stores/modules/permission'
import { getToken } from '@/utils/auth'

// 配置 NProgress
NProgress.configure({ showSpinner: false })

// 白名单路由
const whiteList = ['/login', '/403', '/404']

export function setupRouterGuards(router: Router): void {
  router.beforeEach(async (to, from, next) => {
    // 开始进度条
    NProgress.start()

    // 设置页面标题
    const title = import.meta.env.VITE_APP_TITLE
    document.title = to.meta?.title ? `${to.meta.title} - ${title}` : title

    const token = getToken()

    if (token) {
      if (to.path === '/login') {
        // 已登录，跳转到首页
        next({ path: '/' })
        NProgress.done()
      } else {
        const userStore = useUserStore()
        const permissionStore = usePermissionStore()

        // 检查是否已获取用户信息
        if (!userStore.userInfo) {
          try {
            // 获取用户信息
            await userStore.fetchUserInfo()

            // 生成动态路由
            const accessRoutes = permissionStore.generateRoutes(
              userStore.permissions,
              userStore.isSuperAdmin
            )

            // 移除临时占位符路由（在 constantRoutes 中定义）
            if (router.hasRoute('AsyncRoutePlaceholder')) {
              router.removeRoute('AsyncRoutePlaceholder')
            }

            // 添加动态路由
            accessRoutes.forEach(route => {
              router.addRoute(route)
            })

            // 检查是否有登录后的重定向地址
            const loginRedirect = sessionStorage.getItem('login_redirect')
            if (loginRedirect) {
              sessionStorage.removeItem('login_redirect')
              next({ path: loginRedirect, replace: true })
            } else {
              // 使用 fullPath 重新导航，让 Vue Router 用新添加的路由重新解析
              // 不要使用 { ...to } 因为它会保留原来的 matched: [] 状态
              next({ path: to.fullPath, replace: true })
            }
          } catch (error) {
            // 获取用户信息失败，清除 token 并跳转登录页
            userStore.resetState()
            next(`/login?redirect=${to.path}`)
            NProgress.done()
          }
        } else {
          // 检查路由权限
          if (to.meta?.permission) {
            const hasPerm = userStore.hasPermission(to.meta.permission as string | string[])
            if (hasPerm) {
              next()
            } else {
              next('/403')
              NProgress.done()
            }
          } else {
            next()
          }
        }
      }
    } else {
      // 未登录
      if (whiteList.includes(to.path)) {
        next()
      } else {
        next(`/login?redirect=${to.path}`)
        NProgress.done()
      }
    }
  })

  router.afterEach(() => {
    NProgress.done()
  })
}
