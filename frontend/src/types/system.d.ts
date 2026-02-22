// 用户信息
export interface UserInfo {
  id: number
  username: string
  realName: string
  phone: string
  email: string
  avatar: string
  deptId: number
  deptName: string
  userType: number
  status: number
  roles: RoleInfo[]
  permissions: string[]
  dataScope: number
  lastLoginAt: string
  lastLoginIp: string
  createdAt: string
}

// 用户列表项
export interface User {
  id: number
  username: string
  realName: string
  phone: string
  email: string
  avatar: string
  deptId: number
  deptName?: string
  userType: number
  status: number
  lastLoginAt: string
  lastLoginIp: string
  createdAt: string
  updatedAt: string
}

// 用户查询参数
export interface UserQuery {
  username?: string
  realName?: string
  phone?: string
  deptId?: number
  status?: number
  page: number
  pageSize: number
}

// 用户表单
export interface UserForm {
  id?: number
  username: string
  password?: string
  realName: string
  phone: string
  email: string
  deptId: number
  status: number
  roleIds: number[]
}

// 部门
export interface Dept {
  id: number
  parentId: number
  deptCode: string
  deptName: string
  leaderId: number
  leaderName?: string
  sortOrder: number
  status: number
  children?: Dept[]
  createdAt: string
}

// 部门表单
export interface DeptForm {
  id?: number
  parentId: number
  deptCode: string
  deptName: string
  leaderId?: number
  sortOrder: number
  status: number
}

// 角色
export interface Role {
  id: number
  roleCode: string
  roleName: string
  roleType: number
  dataScope: number
  description: string
  status: number
  sortOrder: number
  createdAt: string
}

// 角色信息（简化）
export interface RoleInfo {
  id: number
  roleCode: string
  roleName: string
}

// 角色表单
export interface RoleForm {
  id?: number
  roleCode: string
  roleName: string
  dataScope: number
  description: string
  status: number
  sortOrder: number
  permissionIds: number[]
}

// 权限
export interface Permission {
  id: number
  parentId: number
  permCode: string
  permName: string
  permType: number
  path: string
  component: string
  icon: string
  apiPath: string
  apiMethod: string
  visible: number
  sortOrder: number
  status: number
  children?: Permission[]
}

// 权限表单
export interface PermissionForm {
  id?: number
  parentId: number
  permCode: string
  permName: string
  permType: number
  path: string
  component: string
  icon: string
  apiPath: string
  apiMethod: string
  visible: number
  sortOrder: number
  status: number
}

// 操作日志
export interface OperationLog {
  id: number
  userId: number
  username: string
  module: string
  action: string
  description: string
  targetType: string
  targetId: string
  requestMethod: string
  requestUrl: string
  requestParams: string
  responseCode: number
  oldValue: string
  newValue: string
  ipAddress: string
  userAgent: string
  durationMs: number
  status: number
  errorMsg: string
  createdAt: string
}

// 日志查询参数
export interface LogQuery {
  username?: string
  module?: string
  action?: string
  status?: number
  startTime?: string
  endTime?: string
  page: number
  pageSize: number
}

// 登录参数
export interface LoginParams {
  username: string
  password: string
  rememberMe?: boolean
}

// 登录响应
export interface LoginResult {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

// 修改密码参数
export interface ChangePasswordParams {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}
