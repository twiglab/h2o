<template>
  <!-- 只有一个子路由时，直接显示子路由 -->
  <template v-if="hasOneShowingChild(item.children, item)">
    <el-menu-item
      v-if="onlyOneChild && !onlyOneChild.children"
      :index="resolvePath(onlyOneChild.path)"
    >
      <el-icon v-if="onlyOneChild.meta?.icon || item.meta?.icon">
        <component :is="onlyOneChild.meta?.icon || item.meta?.icon" />
      </el-icon>
      <template #title>
        <span>{{ onlyOneChild.meta?.title || item.meta?.title }}</span>
      </template>
    </el-menu-item>
  </template>

  <!-- 有多个子路由时，显示子菜单 -->
  <el-sub-menu v-else :index="resolvePath(item.path)">
    <template #title>
      <el-icon v-if="item.meta?.icon">
        <component :is="item.meta.icon" />
      </el-icon>
      <span>{{ item.meta?.title }}</span>
    </template>

    <template v-for="child in item.children" :key="child.path">
      <MenuItem
        v-if="!child.meta?.hidden"
        :item="child"
        :base-path="resolvePath(child.path)"
      />
    </template>
  </el-sub-menu>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import path from 'path-browserify'

const props = defineProps<{
  item: RouteRecordRaw
  basePath?: string
}>()

const onlyOneChild = ref<RouteRecordRaw | null>(null)

// 检查是否只有一个显示的子路由
function hasOneShowingChild(children: RouteRecordRaw[] = [], parent: RouteRecordRaw): boolean {
  const showingChildren = children.filter(item => !item.meta?.hidden)

  // 只有一个子路由时，显示子路由
  if (showingChildren.length === 1) {
    onlyOneChild.value = showingChildren[0]
    return true
  }

  // 没有子路由时，显示父路由
  if (showingChildren.length === 0) {
    onlyOneChild.value = { ...parent, path: '' }
    return true
  }

  return false
}

// 解析路径
function resolvePath(routePath: string): string {
  if (isExternal(routePath)) {
    return routePath
  }
  if (isExternal(props.basePath || '')) {
    return props.basePath || ''
  }
  return path.resolve(props.basePath || '', routePath)
}

// 判断是否为外部链接
function isExternal(path: string): boolean {
  return /^(https?:|mailto:|tel:)/.test(path)
}
</script>
