import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export interface TagView {
  path: string
  fullPath: string
  name: string
  title: string
  affix?: boolean
}

export const useTagsViewStore = defineStore('tagsView', () => {
  // 状态
  const visitedViews = ref<TagView[]>([])
  const cachedViews = ref<string[]>([])

  // 添加视图
  function addView(view: RouteLocationNormalized): void {
    addVisitedView(view)
    addCachedView(view)
  }

  // 添加已访问视图
  function addVisitedView(view: RouteLocationNormalized): void {
    if (visitedViews.value.some(v => v.path === view.path)) return

    visitedViews.value.push({
      path: view.path,
      fullPath: view.fullPath,
      name: view.name as string,
      title: view.meta?.title as string || 'no-title',
      affix: view.meta?.affix
    })
  }

  // 直接添加 TagView（用于初始化固定标签）
  function addTagView(tag: TagView): void {
    if (visitedViews.value.some(v => v.path === tag.path)) return
    visitedViews.value.push(tag)
  }

  // 添加缓存视图
  function addCachedView(view: RouteLocationNormalized): void {
    const name = view.name as string
    if (!name) return
    if (cachedViews.value.includes(name)) return
    if (view.meta?.keepAlive !== false) {
      cachedViews.value.push(name)
    }
  }

  // 删除视图
  function removeView(view: TagView): Promise<{ visitedViews: TagView[]; cachedViews: string[] }> {
    return new Promise(resolve => {
      removeVisitedView(view)
      removeCachedView(view)
      resolve({
        visitedViews: [...visitedViews.value],
        cachedViews: [...cachedViews.value]
      })
    })
  }

  // 删除已访问视图
  function removeVisitedView(view: TagView): void {
    const index = visitedViews.value.findIndex(v => v.path === view.path)
    if (index > -1) {
      visitedViews.value.splice(index, 1)
    }
  }

  // 删除缓存视图
  function removeCachedView(view: TagView): void {
    const name = view.name
    if (!name) return
    const index = cachedViews.value.indexOf(name)
    if (index > -1) {
      cachedViews.value.splice(index, 1)
    }
  }

  // 删除其他视图
  function removeOtherViews(view: TagView): void {
    visitedViews.value = visitedViews.value.filter(v => v.affix || v.path === view.path)
    cachedViews.value = cachedViews.value.filter(name => name === view.name)
  }

  // 删除所有视图
  function removeAllViews(): void {
    // 保留固定的标签
    visitedViews.value = visitedViews.value.filter(v => v.affix)
    cachedViews.value = []
  }

  // 删除右侧视图
  function removeRightViews(view: TagView): void {
    const index = visitedViews.value.findIndex(v => v.path === view.path)
    if (index === -1) return

    visitedViews.value = visitedViews.value.filter((v, i) => {
      if (v.affix || i <= index) return true
      const name = v.name
      if (name) {
        const cacheIndex = cachedViews.value.indexOf(name)
        if (cacheIndex > -1) {
          cachedViews.value.splice(cacheIndex, 1)
        }
      }
      return false
    })
  }

  return {
    visitedViews,
    cachedViews,
    addView,
    addVisitedView,
    addTagView,
    addCachedView,
    removeView,
    removeVisitedView,
    removeCachedView,
    removeOtherViews,
    removeAllViews,
    removeRightViews
  }
})
