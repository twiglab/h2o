import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    // 页面标题
    title?: string
    // 图标
    icon?: string
    // 权限码
    permission?: string | string[]
    // 是否缓存
    keepAlive?: boolean
    // 是否隐藏菜单
    hidden?: boolean
    // 是否固定在标签页
    affix?: boolean
    // 是否显示面包屑
    breadcrumb?: boolean
    // 激活的菜单路径
    activeMenu?: string
  }
}

export interface MenuItem {
  path: string
  title: string
  icon?: string
  children?: MenuItem[]
}
