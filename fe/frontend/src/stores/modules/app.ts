import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const device = ref<'desktop' | 'mobile'>('desktop')

  // 切换侧边栏
  function toggleSidebar(): void {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 设置侧边栏状态
  function setSidebarCollapsed(collapsed: boolean): void {
    sidebarCollapsed.value = collapsed
  }

  // 设置设备类型
  function setDevice(value: 'desktop' | 'mobile'): void {
    device.value = value
  }

  return {
    sidebarCollapsed,
    device,
    toggleSidebar,
    setSidebarCollapsed,
    setDevice
  }
}, {
  persist: {
    key: 'app',
    paths: ['sidebarCollapsed']
  }
})
