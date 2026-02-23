import { get } from '../request'

// 仪表盘统计数据
export interface DashboardStats {
  meter: {
    total: number
    online: number
  }
  account: {
    total: number
    arrears_count: number
    total_balance: string
    total_arrears: string
  }
  today: {
    recharge: string
    deduction: string
    reading_count: number
  }
}

// 用量统计
export interface ConsumptionStat {
  date: string
  electric_consumption: string
  water_consumption: string
  electric_count: number
  water_count: number
}

// 收入统计
export interface RevenueStat {
  date: string
  recharge_amount: string
  deduction_amount: string
}

// 获取仪表盘统计
export function getDashboard(): Promise<DashboardStats> {
  return get<DashboardStats>('/api/dashboard')
}

// 获取用量报表
export function getConsumptionReport(params: {
  start_date: string
  end_date: string
  group_by?: 'day' | 'month'
}): Promise<ConsumptionStat[]> {
  return get<ConsumptionStat[]>('/api/reports/consumption', params)
}

// 获取收入报表
export function getRevenueReport(params: {
  start_date: string
  end_date: string
  group_by?: 'day' | 'month'
}): Promise<RevenueStat[]> {
  return get<RevenueStat[]>('/api/reports/revenue', params)
}
