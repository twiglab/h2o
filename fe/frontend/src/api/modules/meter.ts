import { get, post, put, del } from '../request'
import type { PageResult } from '@/types/common'

// 电表列表查询参数
export interface ElectricMeterQuery {
  page: number
  pageSize: number
  meterNo?: string
  merchantId?: number
  onlineStatus?: number
  status?: number
}

// 水表列表查询参数
export interface WaterMeterQuery {
  page: number
  pageSize: number
  meterNo?: string
  merchantId?: number
  onlineStatus?: number
  status?: number
}

// ==================== 电表 API ====================

// 获取电表列表
export function getElectricMeterList(params: ElectricMeterQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>('/api/electric-meters', params)
}

// 获取电表详情
export function getElectricMeter(id: number): Promise<any> {
  return get<any>(`/api/electric-meters/${id}`)
}

// 新增电表
export function createElectricMeter(data: any): Promise<any> {
  return post('/api/electric-meters', data)
}

// 更新电表
export function updateElectricMeter(id: number, data: any): Promise<any> {
  return put(`/api/electric-meters/${id}`, data)
}

// 删除电表
export function deleteElectricMeter(id: number): Promise<any> {
  return del(`/api/electric-meters/${id}`)
}

// 手工抄表
export function manualReadingElectric(id: number, data: { reading: number; readingTime: string }): Promise<any> {
  return post(`/api/electric-meters/${id}/reading`, data)
}

// 获取电表读数记录
export function getElectricReadings(id: number, params: { page: number; pageSize: number }): Promise<PageResult<any>> {
  return get<PageResult<any>>(`/api/electric-meters/${id}/readings`, params)
}

// ==================== 水表 API ====================

// 获取水表列表
export function getWaterMeterList(params: WaterMeterQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>('/api/water-meters', params)
}

// 获取水表详情
export function getWaterMeter(id: number): Promise<any> {
  return get<any>(`/api/water-meters/${id}`)
}

// 新增水表
export function createWaterMeter(data: any): Promise<any> {
  return post('/api/water-meters', data)
}

// 更新水表
export function updateWaterMeter(id: number, data: any): Promise<any> {
  return put(`/api/water-meters/${id}`, data)
}

// 删除水表
export function deleteWaterMeter(id: number): Promise<any> {
  return del(`/api/water-meters/${id}`)
}

// 手工抄表
export function manualReadingWater(id: number, data: { reading: number; readingTime: string }): Promise<any> {
  return post(`/api/water-meters/${id}/reading`, data)
}

// 获取水表读数记录
export function getWaterReadings(id: number, params: { page: number; pageSize: number }): Promise<PageResult<any>> {
  return get<PageResult<any>>(`/api/water-meters/${id}/readings`, params)
}
