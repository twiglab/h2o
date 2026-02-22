import { get, post, put, del } from '../request'
import type { Dept, DeptForm } from '@/types/system'

// 获取部门树
export function getDeptTree(): Promise<Dept[]> {
  return get<Dept[]>('/api/depts')
}

// 获取部门详情
export function getDeptDetail(id: number): Promise<Dept> {
  return get<Dept>(`/api/depts/${id}`)
}

// 创建部门
export function createDept(data: DeptForm): Promise<Dept> {
  return post<Dept>('/api/depts', data)
}

// 更新部门
export function updateDept(id: number, data: DeptForm): Promise<Dept> {
  return put<Dept>(`/api/depts/${id}`, data)
}

// 删除部门
export function deleteDept(id: number): Promise<void> {
  return del<void>(`/api/depts/${id}`)
}
