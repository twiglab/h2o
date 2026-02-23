import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi, getUserInfo as getUserInfoApi } from '@/api/modules/auth'
import { setToken, setRefreshToken, clearAuth, getToken } from '@/utils/auth'
import type { LoginParams, LoginResult, UserInfo } from '@/types/system'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>(getToken())
  const userInfo = ref<UserInfo | null>(null)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const permissions = computed(() => userInfo.value?.permissions || [])
  const roles = computed(() => userInfo.value?.roles || [])
  const isSuperAdmin = computed(() => userInfo.value?.userType === 1)

  // 登录
  async function login(params: LoginParams): Promise<LoginResult> {
    const result = await loginApi(params)
    token.value = result.accessToken
    setToken(result.accessToken, result.expiresIn)
    setRefreshToken(result.refreshToken)
    return result
  }

  // 获取用户信息
  async function fetchUserInfo(): Promise<UserInfo> {
    const info = await getUserInfoApi()
    userInfo.value = info
    return info
  }

  // 登出
  async function logout(): Promise<void> {
    try {
      await logoutApi()
    } catch (e) {
      // 忽略登出错误
    } finally {
      resetState()
    }
  }

  // 重置状态
  function resetState(): void {
    token.value = ''
    userInfo.value = null
    clearAuth()
  }

  // 检查权限
  function hasPermission(perm: string | string[]): boolean {
    // 超级管理员拥有所有权限
    if (isSuperAdmin.value) return true

    const perms = permissions.value
    if (Array.isArray(perm)) {
      return perm.some(p => perms.includes(p))
    }
    return perms.includes(perm)
  }

  // 检查所有权限
  function hasAllPermissions(perms: string[]): boolean {
    if (isSuperAdmin.value) return true
    return perms.every(p => permissions.value.includes(p))
  }

  // 修改密码
  async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
    // TODO: 调用API
    console.log('changePassword', oldPassword, newPassword)
    // await changePasswordApi({ oldPassword, newPassword })
  }

  return {
    // 状态
    token,
    userInfo,
    // 计算属性
    isLoggedIn,
    permissions,
    roles,
    isSuperAdmin,
    // 方法
    login,
    fetchUserInfo,
    logout,
    resetState,
    hasPermission,
    hasAllPermissions,
    changePassword
  }
}, {
  persist: {
    key: 'user',
    paths: ['token']
  }
})
