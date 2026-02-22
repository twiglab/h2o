// 本地存储封装

const PREFIX = 'prepaid_'

// 设置
export function setStorage(key: string, value: any): void {
  const data = JSON.stringify(value)
  localStorage.setItem(PREFIX + key, data)
}

// 获取
export function getStorage<T = any>(key: string, defaultValue?: T): T | undefined {
  const data = localStorage.getItem(PREFIX + key)
  if (data) {
    try {
      return JSON.parse(data) as T
    } catch {
      return defaultValue
    }
  }
  return defaultValue
}

// 移除
export function removeStorage(key: string): void {
  localStorage.removeItem(PREFIX + key)
}

// 清空
export function clearStorage(): void {
  const keys = Object.keys(localStorage)
  keys.forEach(key => {
    if (key.startsWith(PREFIX)) {
      localStorage.removeItem(key)
    }
  })
}

// Session Storage
export function setSession(key: string, value: any): void {
  const data = JSON.stringify(value)
  sessionStorage.setItem(PREFIX + key, data)
}

export function getSession<T = any>(key: string, defaultValue?: T): T | undefined {
  const data = sessionStorage.getItem(PREFIX + key)
  if (data) {
    try {
      return JSON.parse(data) as T
    } catch {
      return defaultValue
    }
  }
  return defaultValue
}

export function removeSession(key: string): void {
  sessionStorage.removeItem(PREFIX + key)
}

// 记住用户名
const REMEMBER_USER_KEY = 'remembered_user'

export function rememberUser(username: string): void {
  setStorage(REMEMBER_USER_KEY, username)
}

export function getRememberedUser(): string | undefined {
  return getStorage<string>(REMEMBER_USER_KEY)
}

export function clearRememberedUser(): void {
  removeStorage(REMEMBER_USER_KEY)
}
