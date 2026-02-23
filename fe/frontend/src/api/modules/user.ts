import { get, post, put, del } from '../request'
import type { User, UserQuery, UserForm } from '@/types/system'
import type { PageResult } from '@/types/common'

// 获取用户列表
export function getUserList(params: UserQuery): Promise<PageResult<User>> {
  return get<PageResult<User>>('/api/users', params)
}

// 获取用户详情
export function getUserDetail(id: number): Promise<User> {
  return get<User>(`/api/users/${id}`)
}

// 创建用户
export function createUser(data: UserForm): Promise<User> {
  return post<User>('/api/users', data)
}

// 更新用户
export function updateUser(id: number, data: UserForm): Promise<User> {
  return put<User>(`/api/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: number): Promise<void> {
  return del<void>(`/api/users/${id}`)
}

// 重置密码
export function resetPassword(id: number, password: string): Promise<void> {
  return put<void>(`/api/users/${id}/reset-password`, { password })
}

// 分配角色
export function assignRoles(id: number, roleIds: number[]): Promise<void> {
  return put<void>(`/api/users/${id}/roles`, { roleIds })
}

// 获取用户角色
export function getUserRoles(id: number): Promise<number[]> {
  return get<number[]>(`/api/users/${id}/roles`)
}
