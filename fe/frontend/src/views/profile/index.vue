<template>
  <div class="page-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="never">
          <div class="profile-header">
            <el-avatar :size="80" :src="userStore.userInfo?.avatar">
              {{ userStore.userInfo?.realName?.charAt(0) || 'U' }}
            </el-avatar>
            <h3>{{ userStore.userInfo?.realName || userStore.userInfo?.username }}</h3>
            <p class="text-muted">{{ userStore.userInfo?.deptName || '未分配部门' }}</p>
          </div>
          <el-divider />
          <div class="profile-info">
            <div class="info-item">
              <el-icon><User /></el-icon>
              <span>用户名：{{ userStore.userInfo?.username }}</span>
            </div>
            <div class="info-item">
              <el-icon><Phone /></el-icon>
              <span>手机号：{{ userStore.userInfo?.phone || '未设置' }}</span>
            </div>
            <div class="info-item">
              <el-icon><Message /></el-icon>
              <span>邮箱：{{ userStore.userInfo?.email || '未设置' }}</span>
            </div>
            <div class="info-item">
              <el-icon><Calendar /></el-icon>
              <span>注册时间：{{ userStore.userInfo?.createdAt }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="never">
          <template #header><span>基本资料</span></template>
          <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" style="max-width: 500px">
            <el-form-item label="姓名" prop="realName">
              <el-input v-model="formData.realName" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="formData.phone" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="formData.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSave">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { User, Phone, Message, Calendar } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'

const userStore = useUserStore()
const formRef = ref<FormInstance>()

const formData = ref({
  realName: '',
  phone: '',
  email: ''
})

const formRules = reactive<FormRules>({
  realName: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
})

function handleSave() {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      // TODO: call API
      ElMessage.success('保存成功')
    }
  })
}

onMounted(() => {
  if (userStore.userInfo) {
    formData.value.realName = userStore.userInfo.realName || ''
    formData.value.phone = userStore.userInfo.phone || ''
    formData.value.email = userStore.userInfo.email || ''
  }
})
</script>

<style scoped>
.profile-header { text-align: center; padding: 20px 0; }
.profile-header h3 { margin: 16px 0 8px; font-size: 18px; }
.profile-info { padding: 0 20px; }
.info-item { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; color: #606266; }
</style>
