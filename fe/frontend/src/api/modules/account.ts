import { get, post, put } from '../request'
import type { Account, AccountQuery, AccountForm, RechargeForm, Recharge, RechargeQuery } from '@/types/finance'
import type { PageResult } from '@/types/common'

// 获取账户列表
export function getAccountList(params: AccountQuery): Promise<PageResult<Account>> {
  return get<PageResult<Account>>('/api/accounts', params)
}

// 获取所有账户（下拉选择用）
export function getAllAccounts(merchant_id?: number): Promise<Account[]> {
  return get<Account[]>('/api/accounts/all', { merchant_id })
}

// 获取账户详情
export function getAccountDetail(id: number): Promise<Account> {
  return get<Account>(`/api/accounts/${id}`)
}

// 创建账户
export function createAccount(data: AccountForm): Promise<Account> {
  return post<Account>('/api/accounts', data)
}

// 更新账户
export function updateAccount(id: number, data: AccountForm): Promise<Account> {
  return put<Account>(`/api/accounts/${id}`, data)
}

// 账户充值
export function rechargeAccount(id: number, data: RechargeForm): Promise<Recharge> {
  return post<Recharge>(`/api/accounts/${id}/recharge`, data)
}

// 获取充值记录列表
export function getRechargeList(params: RechargeQuery): Promise<PageResult<Recharge>> {
  return get<PageResult<Recharge>>('/api/accounts/recharges', params)
}

// 获取账户充值记录
export function getAccountRecharges(id: number, params?: { start_date?: string; end_date?: string; page?: number; pageSize?: number }): Promise<PageResult<Recharge>> {
  return get<PageResult<Recharge>>(`/api/accounts/${id}/recharges`, params)
}

// 获取欠费账户列表
export function getArrearsAccounts(params?: { page?: number; pageSize?: number }): Promise<PageResult<Account>> {
  return get<PageResult<Account>>('/api/accounts/arrears', params)
}

// 扣费记录查询参数
export interface DeductionQuery {
  page: number
  pageSize: number
  accountNo?: string
  meterNo?: string
  startDate?: string
  endDate?: string
}

// 获取电费扣费记录列表
export function getElectricDeductionList(params: DeductionQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>('/api/accounts/electric-deductions', params)
}

// 获取水费扣费记录列表
export function getWaterDeductionList(params: DeductionQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>('/api/accounts/water-deductions', params)
}

// 获取账户电费扣费记录
export function getAccountElectricDeductions(id: number, params?: DeductionQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>(`/api/accounts/${id}/electric-deductions`, params)
}

// 获取账户水费扣费记录
export function getAccountWaterDeductions(id: number, params?: DeductionQuery): Promise<PageResult<any>> {
  return get<PageResult<any>>(`/api/accounts/${id}/water-deductions`, params)
}
