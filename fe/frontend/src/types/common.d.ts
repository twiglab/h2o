// 分页参数
export interface PageParams {
  page: number
  pageSize: number
}

// 分页结果
export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// API 响应
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
export interface PageResponse<T = any> {
  code: number
  message: string
  data: PageResult<T>
}

// 状态选项
export interface StatusOption {
  label: string
  value: number | string
  type?: 'success' | 'warning' | 'danger' | 'info' | ''
}

// 树形节点
export interface TreeNode {
  id: number
  label: string
  children?: TreeNode[]
  [key: string]: any
}

// 选择器选项
export interface SelectOption {
  label: string
  value: number | string
  disabled?: boolean
}

// 表格列配置
export interface TableColumn {
  prop: string
  label: string
  width?: number | string
  minWidth?: number | string
  fixed?: 'left' | 'right' | boolean
  align?: 'left' | 'center' | 'right'
  sortable?: boolean | 'custom'
  formatter?: (row: any, column: any, cellValue: any, index: number) => any
  slot?: string
}
