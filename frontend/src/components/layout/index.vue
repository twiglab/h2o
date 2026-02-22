<template>
  <div class="layout-container" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <!-- 侧边栏 -->
    <Sidebar />

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <Navbar />

      <!-- 标签页 -->
      <TagsView />

      <!-- 页面内容 -->
      <div class="app-main">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <keep-alive :include="cachedViews">
              <component :is="Component" :key="route.path" />
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount } from 'vue'
import { useAppStore, useTagsViewStore } from '@/stores'
import Sidebar from './Sidebar/index.vue'
import Navbar from './Navbar/index.vue'
import TagsView from './TagsView/index.vue'

const appStore = useAppStore()
const tagsViewStore = useTagsViewStore()

const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const cachedViews = computed(() => tagsViewStore.cachedViews)

// 响应式处理
function handleResize() {
  const width = document.body.clientWidth
  if (width < 992) {
    appStore.setDevice('mobile')
    appStore.setSidebarCollapsed(true)
  } else {
    appStore.setDevice('desktop')
  }
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.layout-container {
  display: flex;
  height: 100vh;
  width: 100%;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-left: $sidebar-width;
  transition: margin-left $transition-duration $transition-timing;
}

.sidebar-collapsed .main-container {
  margin-left: $sidebar-collapsed-width;
}

.app-main {
  flex: 1;
  overflow: auto;
  background-color: $bg-page;
}

// 页面过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
