<template>
  <div class="page-container">
    <el-page-header @back="goBack" title="返回" :content="'店铺详情 - ' + detail.shop_name" />

    <el-card shadow="never" style="margin-top: 20px">
      <el-descriptions title="基本信息" :column="3" border>
        <el-descriptions-item label="店铺编号">{{ detail.shop_no }}</el-descriptions-item>
        <el-descriptions-item label="店铺名称">{{ detail.shop_name }}</el-descriptions-item>
        <el-descriptions-item label="所属商户">{{ detail.merchant_name }}</el-descriptions-item>
        <el-descriptions-item label="位置">{{ detail.building }}-{{ detail.floor }}-{{ detail.room_no }}</el-descriptions-item>
        <el-descriptions-item label="联系人">{{ detail.contact_name }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ detail.contact_phone }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <template v-if="detail.status != null">
            <el-tag :type="detail.status === 1 ? 'success' : 'danger'">{{ detail.status === 1 ? '正常' : '停用' }}</el-tag>
          </template>
          <span v-else>-</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as shopApi from '@/api/modules/shop'

const route = useRoute()
const router = useRouter()
const detail = ref<any>({})

function goBack() {
  router.push('/customer/shop')
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (id) {
    detail.value = await shopApi.getShopDetail(id)
  }
})
</script>
