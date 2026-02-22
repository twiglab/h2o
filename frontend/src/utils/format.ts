import dayjs from 'dayjs'

// 格式化日期时间
export function formatDateTime(value: string | Date | null | undefined, format = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!value) return '-'
  return dayjs(value).format(format)
}

// 格式化日期
export function formatDate(value: string | Date | null | undefined, format = 'YYYY-MM-DD'): string {
  if (!value) return '-'
  return dayjs(value).format(format)
}

// 格式化数字（千分位）
export function formatNumber(value: number | null | undefined, decimals = 0): string {
  if (value === null || value === undefined) return '0'
  return value.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals
  })
}

// 格式化金额
export function formatMoney(value: number | null | undefined, decimals = 2): string {
  if (value === null || value === undefined) return '-'
  return value.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals
  })
}

// 格式化金额（带符号）
export function formatCurrency(value: number | null | undefined, decimals = 2): string {
  if (value === null || value === undefined) return '-'
  return '¥' + formatMoney(value, decimals)
}

// 格式化百分比
export function formatPercent(value: number | null | undefined, decimals = 2): string {
  if (value === null || value === undefined) return '-'
  return (value * 100).toFixed(decimals) + '%'
}

// 格式化用量
export function formatConsumption(value: number | null | undefined, unit: string, decimals = 2): string {
  if (value === null || value === undefined) return '-'
  return value.toFixed(decimals) + ' ' + unit
}

// 格式化电量
export function formatElectricity(value: number | null | undefined): string {
  return formatConsumption(value, 'kWh')
}

// 格式化水量
export function formatWater(value: number | null | undefined): string {
  return formatConsumption(value, '吨')
}

// 格式化状态
export function formatStatus(value: number, options: { value: number; label: string }[]): string {
  const option = options.find(item => item.value === value)
  return option?.label || '-'
}

// 省略文本
export function ellipsis(text: string | null | undefined, length = 20): string {
  if (!text) return '-'
  if (text.length <= length) return text
  return text.slice(0, length) + '...'
}
