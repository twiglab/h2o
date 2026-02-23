import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/stores/modules/user'

// 权限指令
export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    const { value } = binding

    if (value) {
      const userStore = useUserStore()
      const hasPermission = userStore.hasPermission(value)

      if (!hasPermission) {
        // 移除元素
        el.parentNode?.removeChild(el)
      }
    }
  }
}
