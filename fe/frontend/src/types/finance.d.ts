// 账户
// 一个商户可以有多个账户，一个账户可以为多个店铺的表计付费
export interface Account {
  id: number
  account_no: string
  account_name?: string
  merchant_id: number
  balance: number
  total_recharge: number
  total_consumption: number
  status: number
  remark?: string
  created_at: string
  updated_at: string
  // 关联数据
  merchant?: {
    id: number
    merchant_no: string
    merchant_name: string
  }
}

// 账户查询参数
export interface AccountQuery {
  account_no?: string
  account_name?: string
  merchant_id?: number
  status?: number
  page: number
  pageSize: number
}

// 账户表单
export interface AccountForm {
  id?: number
  account_no?: string
  account_name?: string
  merchant_id: number
  status: number
  remark?: string
}

// 充值记录
export interface Recharge {
  id: number
  recharge_no: string
  account_id: number
  amount: number
  balance_before: number
  balance_after: number
  payment_method: number
  payment_no?: string
  operator_id: number
  operator_name?: string
  remark?: string
  created_at: string
  // 关联数据
  account?: Account
}

// 充值查询参数
export interface RechargeQuery {
  recharge_no?: string
  account_id?: number
  payment_method?: number
  start_time?: string
  end_time?: string
  page: number
  pageSize: number
}

// 充值表单
export interface RechargeForm {
  account_id: number
  amount: number
  payment_method: number
  payment_no?: string
  remark?: string
}

// 电费扣费记录
export interface ElectricDeduction {
  id: number
  deduction_no: string
  account_id: number
  meter_id: number
  consumption: number
  unit_price: number
  amount: number
  balance_before: number
  balance_after: number
  deduction_time: string
  status: number
  remark?: string
  created_at: string
  // 关联数据
  account?: Account
  meter?: {
    id: number
    meter_no: string
  }
}

// 电费扣费查询参数
export interface ElectricDeductionQuery {
  deduction_no?: string
  account_id?: number
  meter_id?: number
  status?: number
  start_time?: string
  end_time?: string
  page: number
  pageSize: number
}

// 水费扣费记录
export interface WaterDeduction {
  id: number
  deduction_no: string
  account_id: number
  meter_id: number
  consumption: number
  unit_price: number
  amount: number
  balance_before: number
  balance_after: number
  deduction_time: string
  status: number
  remark?: string
  created_at: string
  // 关联数据
  account?: Account
  meter?: {
    id: number
    meter_no: string
  }
}

// 水费扣费查询参数
export interface WaterDeductionQuery {
  deduction_no?: string
  account_id?: number
  meter_id?: number
  status?: number
  start_time?: string
  end_time?: string
  page: number
  pageSize: number
}

// 工作台统计
export interface DashboardStats {
  merchant_count: number
  merchant_growth: number
  shop_count: number
  shop_growth: number
  total_balance: number
  balance_change: number
  today_recharge: number
  recharge_growth: number
}

// 趋势数据点
export interface TrendPoint {
  date: string
  electric_value: number
  water_value: number
}

// 收入趋势数据点
export interface RevenueTrendPoint {
  date: string
  recharge: number
  consumption: number
}

// 预警项
export interface WarningItem {
  id: number
  type: 'balance' | 'offline'
  title: string
  description: string
  level: 'warning' | 'danger'
  created_at: string
}
