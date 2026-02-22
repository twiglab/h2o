// 电费费率
export interface ElectricRate {
  id: number
  rate_code: string
  rate_name: string
  scope: number              // 1=全局, 2=商户
  merchant_id?: number
  merchant_name?: string
  calc_mode: number          // 1=固定单价, 2=分时
  unit_price?: number
  effective_date: string
  expire_date?: string
  status: number
  remark?: string
  tou_details?: TOUDetail[]
  service_fees?: ServiceFee[]
  created_at: string
  updated_at: string
}

// 分时电价详情
export interface TOUDetail {
  id?: number
  rate_id?: number
  period_name: string
  start_time: string
  end_time: string
  unit_price: number
}

// 服务费
export interface ServiceFee {
  id?: number
  rate_id?: number
  fee_name: string
  fee_type: number           // 1=固定金额, 2=百分比
  fee_value: number
}

// 电费费率查询参数
export interface ElectricRateQuery {
  keyword?: string
  scope?: number
  calc_mode?: number
  status?: number
  page: number
  pageSize: number
}

// 电费费率表单
export interface ElectricRateForm {
  id?: number
  rate_code?: string
  rate_name: string
  scope: number
  merchant_id?: number
  calc_mode: number
  unit_price?: number
  effective_date: string
  expire_date?: string
  status: number
  remark?: string
  tou_details?: TOUDetail[]
  service_fees?: ServiceFee[]
}
