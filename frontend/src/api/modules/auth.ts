import { post, get, put } from '../request'
import type { LoginParams, LoginResult, UserInfo, ChangePasswordParams } from '@/types/system'

// 后端响应格式 (snake_case)
interface BackendLoginResult {
  access_token: string
  refresh_token: string
  user: {
    id: number
    username: string
    real_name: string
  }
}

interface BackendUserInfo {
  id: number
  username: string
  real_name: string
  phone: string
  email: string
  status: number
}

// 登录
export async function login(data: LoginParams): Promise<LoginResult> {
  const res = await post<BackendLoginResult>('/api/auth/login', data)
  // 转换为前端格式
  return {
    accessToken: res.access_token,
    refreshToken: res.refresh_token,
    expiresIn: 7200 // 默认2小时
  }
}

// 登出
export function logout(): Promise<void> {
  return post<void>('/api/auth/logout')
}

// 刷新 Token
export async function refreshToken(token: string): Promise<LoginResult> {
  const res = await post<{ access_token: string; refresh_token: string }>('/api/auth/refresh', { refresh_token: token })
  return {
    accessToken: res.access_token,
    refreshToken: res.refresh_token,
    expiresIn: 7200
  }
}

// 获取当前用户信息
export async function getUserInfo(): Promise<UserInfo> {
  const res = await get<BackendUserInfo>('/api/auth/profile')
  // 转换为前端格式
  return {
    id: res.id,
    username: res.username,
    realName: res.real_name || '',
    phone: res.phone || '',
    email: res.email || '',
    avatar: '',
    deptId: 0,
    deptName: '',
    userType: 1,
    status: res.status,
    roles: [],
    permissions: ['*'], // 暂时给所有权限
    dataScope: 1,
    lastLoginAt: '',
    lastLoginIp: '',
    createdAt: ''
  }
}

// 修改密码
export function changePassword(data: ChangePasswordParams): Promise<void> {
  return put<void>('/api/auth/password', {
    old_password: data.oldPassword,
    new_password: data.newPassword
  })
}
