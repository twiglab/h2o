// 验证规则

// 手机号正则
export const PHONE_REGEX = /^1[3-9]\d{9}$/

// 邮箱正则
export const EMAIL_REGEX = /^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$/

// 身份证正则
export const ID_CARD_REGEX = /^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/

// 统一社会信用代码正则
export const CREDIT_CODE_REGEX = /^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$/

// 验证手机号
export function isPhone(value: string): boolean {
  return PHONE_REGEX.test(value)
}

// 验证邮箱
export function isEmail(value: string): boolean {
  return EMAIL_REGEX.test(value)
}

// 验证身份证
export function isIdCard(value: string): boolean {
  return ID_CARD_REGEX.test(value)
}

// 验证统一社会信用代码
export function isCreditCode(value: string): boolean {
  return CREDIT_CODE_REGEX.test(value)
}

// 验证URL
export function isUrl(value: string): boolean {
  try {
    new URL(value)
    return true
  } catch {
    return false
  }
}

// 验证是否为空
export function isEmpty(value: any): boolean {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim() === ''
  if (Array.isArray(value)) return value.length === 0
  if (typeof value === 'object') return Object.keys(value).length === 0
  return false
}

// Element Plus 表单验证规则生成器
export function requiredRule(message: string) {
  return { required: true, message, trigger: 'blur' }
}

export function phoneRule(message = '请输入正确的手机号') {
  return {
    pattern: PHONE_REGEX,
    message,
    trigger: 'blur'
  }
}

export function emailRule(message = '请输入正确的邮箱') {
  return {
    pattern: EMAIL_REGEX,
    message,
    trigger: 'blur'
  }
}

export function idCardRule(message = '请输入正确的身份证号') {
  return {
    pattern: ID_CARD_REGEX,
    message,
    trigger: 'blur'
  }
}

export function lengthRule(min: number, max: number, message?: string) {
  return {
    min,
    max,
    message: message || `长度应在 ${min} 到 ${max} 个字符`,
    trigger: 'blur'
  }
}

// Element Plus 自定义验证器（用于 validator 属性）
export function validatePhone(rule: any, value: string, callback: (error?: Error) => void) {
  if (!value) {
    callback()
  } else if (!PHONE_REGEX.test(value)) {
    callback(new Error('请输入正确的手机号'))
  } else {
    callback()
  }
}

export function validateEmail(rule: any, value: string, callback: (error?: Error) => void) {
  if (!value) {
    callback()
  } else if (!EMAIL_REGEX.test(value)) {
    callback(new Error('请输入正确的邮箱'))
  } else {
    callback()
  }
}

export function validateIdCard(rule: any, value: string, callback: (error?: Error) => void) {
  if (!value) {
    callback()
  } else if (!ID_CARD_REGEX.test(value)) {
    callback(new Error('请输入正确的身份证号'))
  } else {
    callback()
  }
}

export function validateCreditCode(rule: any, value: string, callback: (error?: Error) => void) {
  if (!value) {
    callback()
  } else if (!CREDIT_CODE_REGEX.test(value)) {
    callback(new Error('请输入正确的统一社会信用代码'))
  } else {
    callback()
  }
}
