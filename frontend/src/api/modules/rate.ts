import { get, post, put, del } from '../request'
import type { ElectricRate, ElectricRateQuery, ElectricRateForm } from '@/types/rate'
import type { PageResult } from '@/types/common'

// 获取费率列表
export function getRateList(params: ElectricRateQuery): Promise<PageResult<ElectricRate>> {
  return get<PageResult<ElectricRate>>('/api/rates', params)
}

// 获取所有费率（下拉选择用）
export function getAllRates(): Promise<ElectricRate[]> {
  return get<ElectricRate[]>('/api/rates/all')
}

// 获取费率详情
export function getRateDetail(id: number): Promise<ElectricRate> {
  return get<ElectricRate>(`/api/rates/${id}`)
}

// 创建费率
export function createRate(data: ElectricRateForm): Promise<ElectricRate> {
  return post<ElectricRate>('/api/rates', data)
}

// 更新费率
export function updateRate(id: number, data: ElectricRateForm): Promise<ElectricRate> {
  return put<ElectricRate>(`/api/rates/${id}`, data)
}

// 删除费率
export function deleteRate(id: number): Promise<void> {
  return del<void>(`/api/rates/${id}`)
}
