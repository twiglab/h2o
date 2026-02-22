<template>
  <div class="sidebar-container" :class="{ 'is-collapsed': collapsed }">
    <!-- Logo -->
    <div class="sidebar-logo">
      <router-link to="/">
        <el-icon size="24" color="#fff"><Odometer /></el-icon>
        <span v-show="!collapsed" class="logo-title">水电预付费</span>
      </router-link>
    </div>

    <!-- 菜单 -->
    <el-scrollbar>
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :collapse-transition="false"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        unique-opened
        router
      >
        <template v-for="route in routes" :key="route.path">
          <MenuItem v-if="!route.meta?.hidden" :item="route" :base-path="route.path" />
        </template>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore, usePermissionStore } from '@/stores'
import MenuItem from './MenuItem.vue'

const route = useRoute()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const collapsed = computed(() => appStore.sidebarCollapsed)
const routes = computed(() => permissionStore.routes)
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})
</script>

<style lang="scss" scoped>
.sidebar-container {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: $sidebar-width;
  background-color: #304156;
  z-index: 1001;
  overflow: hidden;
  transition: width $transition-duration $transition-timing;

  &.is-collapsed {
    width: $sidebar-collapsed-width;
  }
}

.sidebar-logo {
  height: $navbar-height;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3649;

  a {
    display: flex;
    align-items: center;
    height: 100%;
    width: 100%;
    padding: 0 16px;
  }

  .logo-title {
    margin-left: 12px;
    font-size: 16px;
    font-weight: 600;
    color: #fff;
    white-space: nowrap;
  }
}

:deep(.el-menu) {
  border-right: none;
  width: 100% !important;
}

:deep(.el-scrollbar) {
  height: calc(100% - #{$navbar-height});
}
</style>
