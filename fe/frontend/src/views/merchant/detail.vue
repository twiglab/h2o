<template>
  <div class="page-container">
    <el-page-header @back="goBack" title="返回" :content="'商户详情 - ' + detail.merchant_name" />

    <el-card shadow="never" style="margin-top: 20px">
      <el-descriptions title="基本信息" :column="3" border>
        <el-descriptions-item label="商户编号">{{ detail.merchant_no }}</el-descriptions-item>
        <el-descriptions-item label="商户名称">{{ detail.merchant_name }}</el-descriptions-item>
        <el-descriptions-item label="商户类型">{{ detail.merchant_type === 1 ? '企业' : '个人' }}</el-descriptions-item>
        <el-descriptions-item label="联系人">{{ detail.contact_name }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ detail.contact_phone }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <template v-if="detail.status != null">
            <el-tag :type="detail.status === 1 ? 'success' : 'danger'">{{ detail.status === 1 ? '正常' : '停用' }}</el-tag>
          </template>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="地址" :span="3">{{ detail.address }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as merchantApi from '@/api/modules/merchant'

const route = useRoute()
const router = useRouter()
const detail = ref<any>({})

function goBack() {
  router.push('/customer/merchant')
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (id) {
    detail.value = await merchantApi.getMerchantDetail(id)
  }
})
</script>
