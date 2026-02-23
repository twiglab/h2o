<template>
  <div class="page-container">
    <el-page-header @back="goBack" title="返回" :content="'账户详情 - ' + detail.account_no" />

    <el-card shadow="never" style="margin-top: 20px">
      <el-descriptions title="账户信息" :column="3" border>
        <el-descriptions-item label="账户编号">{{ detail.account_no }}</el-descriptions-item>
        <el-descriptions-item label="账户名称">{{ detail.account_name }}</el-descriptions-item>
        <el-descriptions-item label="所属商户">{{ detail.merchant_name }}</el-descriptions-item>
        <el-descriptions-item label="可用余额">
          <span class="text-money" :class="detail.balance >= 0 ? 'positive' : 'negative'">{{ formatMoney(detail.balance) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="冻结金额">{{ formatMoney(detail.frozen_amount) }}</el-descriptions-item>
        <el-descriptions-item label="欠费金额">
          <span class="text-money negative">{{ formatMoney(detail.arrears) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="累计充值">{{ formatMoney(detail.total_recharge) }}</el-descriptions-item>
        <el-descriptions-item label="累计电费">{{ formatMoney(detail.total_electric_consumption) }}</el-descriptions-item>
        <el-descriptions-item label="累计水费">{{ formatMoney(detail.total_water_consumption) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <template v-if="detail.status != null">
            <el-tag :type="detail.status === 1 ? 'success' : detail.status === 2 ? 'warning' : 'danger'">
              {{ detail.status === 1 ? '正常' : detail.status === 2 ? '欠费' : '冻结' }}
            </el-tag>
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
import { formatMoney } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const detail = ref<any>({})

function goBack() {
  router.push('/customer/account')
}

onMounted(() => {
  const id = route.params.id
  // TODO: fetch detail by id
  detail.value = { accountNo: 'ACC001', accountName: '测试账户', merchantName: '测试商户', balance: 1500.50, frozenAmount: 0, arrears: 0, totalRecharge: 5000, totalElectricConsumption: 2500, totalWaterConsumption: 1000, status: 1 }
})
</script>
