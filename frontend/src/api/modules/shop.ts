import { get, post, put, del } from '../request'
import type { Shop, ShopQuery, ShopForm } from '@/types/business'
import type { PageResult } from '@/types/common'

// 获取店铺列表
export function getShopList(params: ShopQuery): Promise<PageResult<Shop>> {
  return get<PageResult<Shop>>('/api/shops', params)
}

// 获取所有店铺（下拉选择用）
export function getAllShops(merchant_id?: number): Promise<Shop[]> {
  return get<Shop[]>('/api/shops/all', { merchant_id })
}

// 获取店铺详情
export function getShopDetail(id: number): Promise<Shop> {
  return get<Shop>(`/api/shops/${id}`)
}

// 创建店铺
export function createShop(data: ShopForm): Promise<Shop> {
  return post<Shop>('/api/shops', data)
}

// 更新店铺
export function updateShop(id: number, data: ShopForm): Promise<Shop> {
  return put<Shop>(`/api/shops/${id}`, data)
}

// 删除店铺
export function deleteShop(id: number): Promise<void> {
  return del<void>(`/api/shops/${id}`)
}

// 获取店铺统计
export function getShopStats(id: number): Promise<{
  electric_meter_count: number
  water_meter_count: number
  balance: number
  total_recharge: number
  total_consumption: number
}> {
  return get(`/api/shops/${id}/stats`)
}
