import { get, post, put, del } from '../request'
import type { Merchant, MerchantQuery, MerchantForm } from '@/types/business'
import type { PageResult } from '@/types/common'

// 获取商户列表
export function getMerchantList(params: MerchantQuery): Promise<PageResult<Merchant>> {
  return get<PageResult<Merchant>>('/api/merchants', params)
}

// 获取所有商户（下拉选择用）
export function getAllMerchants(): Promise<Merchant[]> {
  return get<Merchant[]>('/api/merchants/all')
}

// 获取商户详情
export function getMerchantDetail(id: number): Promise<Merchant> {
  return get<Merchant>(`/api/merchants/${id}`)
}

// 创建商户
export function createMerchant(data: MerchantForm): Promise<Merchant> {
  return post<Merchant>('/api/merchants', data)
}

// 更新商户
export function updateMerchant(id: number, data: MerchantForm): Promise<Merchant> {
  return put<Merchant>(`/api/merchants/${id}`, data)
}

// 删除商户
export function deleteMerchant(id: number): Promise<void> {
  return del<void>(`/api/merchants/${id}`)
}

// 获取商户统计
export function getMerchantStats(id: number): Promise<{
  shop_count: number
  account_count: number
  electric_meter_count: number
  water_meter_count: number
  total_balance: number
}> {
  return get(`/api/merchants/${id}/stats`)
}
