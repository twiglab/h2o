import { get, post, put, del } from '../request'
import type { Role, RoleForm } from '@/types/system'
import type { PageResult } from '@/types/common'

// 获取角色列表
export function getRoleList(params?: { page?: number; pageSize?: number }): Promise<PageResult<Role>> {
  return get<PageResult<Role>>('/api/roles', params)
}

// 获取所有角色（下拉选择用）
export function getAllRoles(): Promise<Role[]> {
  return get<Role[]>('/api/roles/all')
}

// 获取角色详情
export function getRoleDetail(id: number): Promise<Role> {
  return get<Role>(`/api/roles/${id}`)
}

// 创建角色
export function createRole(data: RoleForm): Promise<Role> {
  return post<Role>('/api/roles', data)
}

// 更新角色
export function updateRole(id: number, data: RoleForm): Promise<Role> {
  return put<Role>(`/api/roles/${id}`, data)
}

// 删除角色
export function deleteRole(id: number): Promise<void> {
  return del<void>(`/api/roles/${id}`)
}

// 分配权限
export function assignPermissions(id: number, permissionIds: number[]): Promise<void> {
  return put<void>(`/api/roles/${id}/permissions`, { permissionIds })
}

// 获取角色权限
export function getRolePermissions(id: number): Promise<number[]> {
  return get<number[]>(`/api/roles/${id}/permissions`)
}
