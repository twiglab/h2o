import Cookies from 'js-cookie'

const TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

// 获取 Token
export function getToken(): string {
  return Cookies.get(TOKEN_KEY) || ''
}

// 设置 Token
export function setToken(token: string, expires?: number): void {
  if (expires) {
    Cookies.set(TOKEN_KEY, token, { expires: expires / 86400 }) // 转换为天
  } else {
    Cookies.set(TOKEN_KEY, token)
  }
}

// 移除 Token
export function removeToken(): void {
  Cookies.remove(TOKEN_KEY)
}

// 获取 Refresh Token
export function getRefreshToken(): string {
  return Cookies.get(REFRESH_TOKEN_KEY) || ''
}

// 设置 Refresh Token
export function setRefreshToken(token: string, expires?: number): void {
  if (expires) {
    Cookies.set(REFRESH_TOKEN_KEY, token, { expires: expires / 86400 })
  } else {
    Cookies.set(REFRESH_TOKEN_KEY, token, { expires: 7 }) // 默认7天
  }
}

// 移除 Refresh Token
export function removeRefreshToken(): void {
  Cookies.remove(REFRESH_TOKEN_KEY)
}

// 清除所有认证信息
export function clearAuth(): void {
  removeToken()
  removeRefreshToken()
}
