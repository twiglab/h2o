<template>
  <div class="login-container">
    <div class="login-card">
      <!-- Logo -->
      <div class="login-header">
        <el-icon size="40" color="#409EFF"><Odometer /></el-icon>
        <h1 class="login-title">水电预付费管理系统</h1>
        <p class="login-subtitle">Prepaid Utility Management System</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item prop="captcha" v-if="captchaEnabled">
          <div class="captcha-wrapper">
            <el-input
              v-model="loginForm.captcha"
              placeholder="请输入验证码"
              size="large"
              :prefix-icon="Picture"
            />
            <img
              v-if="captchaImage"
              :src="captchaImage"
              class="captcha-image"
              alt="验证码"
              @click="refreshCaptcha"
            />
          </div>
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="javascript:;" class="forgot-password">忘记密码？</a>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 底部信息 -->
      <div class="login-footer">
        <p>Copyright &copy; {{ new Date().getFullYear() }} 水电预付费系统</p>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="login-background">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { User, Lock, Picture } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'
import { getRememberedUser, rememberUser, clearRememberedUser } from '@/utils/storage'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)
const captchaEnabled = ref(false) // 是否启用验证码
const captchaImage = ref('')
const captchaKey = ref('')

const loginForm = ref({
  username: '',
  password: '',
  captcha: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度为2-50个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' }
  ]
}

// 刷新验证码
async function refreshCaptcha() {
  // TODO: 调用验证码API
  // const res = await getCaptcha()
  // captchaImage.value = res.image
  // captchaKey.value = res.key
}

// 登录
async function handleLogin() {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login({
          username: loginForm.value.username,
          password: loginForm.value.password,
          captcha: loginForm.value.captcha,
          captchaKey: captchaKey.value
        })

        // 记住我
        if (rememberMe.value) {
          rememberUser(loginForm.value.username)
        } else {
          clearRememberedUser()
        }

        ElMessage.success('登录成功')

        // 获取重定向地址并存储，让路由守卫在加载路由后处理重定向
        const redirect = route.query.redirect as string
        if (redirect && redirect !== '/' && redirect !== '/login') {
          sessionStorage.setItem('login_redirect', redirect)
        }
        // 先跳转到首页，路由守卫会加载动态路由并处理重定向
        router.push('/')
      } catch (error: any) {
        ElMessage.error(error.message || '登录失败')
        if (captchaEnabled.value) {
          refreshCaptcha()
        }
      } finally {
        loading.value = false
      }
    }
  })
}

// 初始化
onMounted(() => {
  // 读取记住的用户名
  const remembered = getRememberedUser()
  if (remembered) {
    loginForm.value.username = remembered
    rememberMe.value = true
  }

  // 加载验证码
  if (captchaEnabled.value) {
    refreshCaptcha()
  }
})
</script>

<style lang="scss" scoped>
.login-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

.login-card {
  position: relative;
  z-index: 10;
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;

  .login-title {
    margin: 16px 0 8px;
    font-size: 24px;
    font-weight: 600;
    color: $text-primary;
  }

  .login-subtitle {
    font-size: 12px;
    color: $text-secondary;
    text-transform: uppercase;
    letter-spacing: 1px;
  }
}

.login-form {
  .el-input {
    :deep(.el-input__wrapper) {
      padding: 0 15px;
      box-shadow: 0 0 0 1px $border-color inset;

      &:hover {
        box-shadow: 0 0 0 1px $primary-color inset;
      }

      &.is-focus {
        box-shadow: 0 0 0 1px $primary-color inset;
      }
    }

    :deep(.el-input__inner) {
      height: 44px;
    }
  }

  .captcha-wrapper {
    display: flex;
    gap: 12px;
    width: 100%;

    .el-input {
      flex: 1;
    }

    .captcha-image {
      width: 120px;
      height: 44px;
      border-radius: $border-radius-md;
      cursor: pointer;
      object-fit: cover;
    }
  }

  .login-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;

    .forgot-password {
      font-size: 13px;
      color: $primary-color;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .login-button {
    width: 100%;
    height: 44px;
    font-size: 16px;
    font-weight: 500;
    border-radius: $border-radius-md;
  }
}

.login-footer {
  margin-top: 32px;
  text-align: center;

  p {
    font-size: 12px;
    color: $text-secondary;
  }
}

// 背景装饰
.login-background {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;

  .bg-shape {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    animation: float 6s ease-in-out infinite;
  }

  .shape-1 {
    width: 300px;
    height: 300px;
    top: -100px;
    right: -100px;
    animation-delay: 0s;
  }

  .shape-2 {
    width: 200px;
    height: 200px;
    bottom: -50px;
    left: -50px;
    animation-delay: 2s;
  }

  .shape-3 {
    width: 150px;
    height: 150px;
    top: 50%;
    right: 10%;
    animation-delay: 4s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-20px) scale(1.05);
  }
}

// 响应式
@media (max-width: 480px) {
  .login-card {
    width: 90%;
    padding: 30px 20px;
  }

  .login-header {
    .login-title {
      font-size: 20px;
    }
  }
}
</style>
