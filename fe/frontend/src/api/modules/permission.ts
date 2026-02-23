import { get, post, put, del } from '../request'
import type { Permission, PermissionForm } from '@/types/system'

// 权限查询参数
export interface PermissionQuery {
  permName?: string
  status?: number
}

// 获取权限树
export function getPermissionTree(params?: PermissionQuery): Promise<Permission[]> {
  return get<Permission[]>('/api/permissions', { params })
}

// 获取权限详情
export function getPermissionDetail(id: number): Promise<Permission> {
  return get<Permission>(`/api/permissions/${id}`)
}

// 创建权限
export function createPermission(data: PermissionForm): Promise<Permission> {
  return post<Permission>('/api/permissions', data)
}

// 更新权限
export function updatePermission(id: number, data: PermissionForm): Promise<Permission> {
  return put<Permission>(`/api/permissions/${id}`, data)
}

// 删除权限
export function deletePermission(id: number): Promise<void> {
  return del<void>(`/api/permissions/${id}`)
}
