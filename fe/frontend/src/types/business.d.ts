// 商户
export interface Merchant {
  id: number
  merchant_no: string
  merchant_name: string
  merchant_type: number
  contact_name?: string
  contact_phone?: string
  status: number
  remark?: string
  created_at: string
  updated_at: string
  // 关联数据
  shop_count?: number
}

// 商户查询参数
export interface MerchantQuery {
  merchant_no?: string
  merchant_name?: string
  merchant_type?: number
  status?: number
  page: number
  pageSize: number
}

// 商户表单
export interface MerchantForm {
  id?: number
  merchant_no?: string
  merchant_name: string
  merchant_type: number
  contact_name?: string
  contact_phone?: string
  status: number
  remark?: string
}

// 店铺
export interface Shop {
  id: number
  shop_no: string
  shop_name: string
  merchant_id: number
  merchant_name?: string
  building?: string
  floor?: string
  room_no?: string
  contact_name?: string
  contact_phone?: string
  status: number
  remark?: string
  created_at: string
  updated_at: string
  // 关联数据
  merchant?: Merchant
}

// 店铺查询参数
export interface ShopQuery {
  shop_no?: string
  shop_name?: string
  merchant_id?: number
  building?: string
  status?: number
  page: number
  pageSize: number
}

// 店铺表单
export interface ShopForm {
  id?: number
  shop_no?: string
  shop_name: string
  merchant_id: number
  building?: string
  floor?: string
  room_no?: string
  contact_name?: string
  contact_phone?: string
  status: number
  remark?: string
}

// 电表
export interface ElectricMeter {
  id: number
  meter_no: string
  merchant_id: number
  merchant_name?: string
  shop_id?: number
  shop_name?: string
  account_id?: number
  account_no?: string
  rate_id?: number
  rate_name?: string
  mqtt_topic?: string
  comm_addr?: string
  protocol: string
  multiplier: number
  init_reading: number
  current_reading: number
  last_collect_at?: string
  online_status: number
  status: number
  remark?: string
  created_at: string
  updated_at: string
}

// 电表查询参数
export interface ElectricMeterQuery {
  meter_no?: string
  merchant_id?: number
  shop_id?: number
  online_status?: number
  status?: number
  page: number
  pageSize: number
}

// 电表表单
export interface ElectricMeterForm {
  id?: number
  meter_no: string
  merchant_id: number
  shop_id?: number
  account_id?: number
  rate_id?: number
  mqtt_topic?: string
  comm_addr?: string
  protocol: string
  multiplier: number
  init_reading: number
  status: number
  remark?: string
}

// 水表
export interface WaterMeter {
  id: number
  meter_no: string
  merchant_id: number
  merchant_name?: string
  shop_id?: number
  shop_name?: string
  account_id?: number
  account_no?: string
  rate_id?: number
  rate_name?: string
  mqtt_topic?: string
  comm_addr?: string
  protocol: string
  multiplier: number
  init_reading: number
  current_reading: number
  last_collect_at?: string
  online_status: number
  status: number
  remark?: string
  created_at: string
  updated_at: string
}

// 水表查询参数
export interface WaterMeterQuery {
  meter_no?: string
  merchant_id?: number
  shop_id?: number
  online_status?: number
  status?: number
  page: number
  pageSize: number
}

// 水表表单
export interface WaterMeterForm {
  id?: number
  meter_no: string
  merchant_id: number
  shop_id?: number
  account_id?: number
  rate_id?: number
  mqtt_topic?: string
  comm_addr?: string
  protocol: string
  multiplier: number
  init_reading: number
  status: number
  remark?: string
}

// 抄表记录
export interface Reading {
  id: number
  merchant_id: number
  meter_id: number
  meter_no: string
  reading_value: number
  reading_time: string
  collect_type: number
  operator_id: number
  operator_name?: string
  status: number
  remark?: string
  created_at: string
}

// 抄表查询参数
export interface ReadingQuery {
  meter_no?: string
  merchant_id?: number
  meter_id?: number
  collect_type?: number
  status?: number
  start_time?: string
  end_time?: string
  page: number
  pageSize: number
}

// 电费费率
export interface ElectricRate {
  id: number
  rate_code: string
  rate_name: string
  unit_price: number
  effective_date: string
  status: number
  remark?: string
  created_at: string
  updated_at: string
}

// 电费费率查询参数
export interface ElectricRateQuery {
  rate_code?: string
  rate_name?: string
  status?: number
  page: number
  pageSize: number
}

// 电费费率表单
export interface ElectricRateForm {
  id?: number
  rate_code: string
  rate_name: string
  unit_price: number
  effective_date: string
  status: number
  remark?: string
}

// 水费费率
export interface WaterRate {
  id: number
  rate_code: string
  rate_name: string
  unit_price: number
  effective_date: string
  status: number
  remark?: string
  created_at: string
  updated_at: string
}

// 水费费率查询参数
export interface WaterRateQuery {
  rate_code?: string
  rate_name?: string
  status?: number
  page: number
  pageSize: number
}

// 水费费率表单
export interface WaterRateForm {
  id?: number
  rate_code: string
  rate_name: string
  unit_price: number
  effective_date: string
  status: number
  remark?: string
}
