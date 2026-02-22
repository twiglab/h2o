<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="充值流水号">
          <el-input v-model="queryParams.keyword" placeholder="请输入" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="所属商户">
          <el-select v-model="queryParams.merchant_id" placeholder="请选择" clearable filterable style="width: 180px">
            <el-option v-for="item in merchantList" :key="item.id" :label="item.merchant_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="充值时间">
          <el-date-picker v-model="queryParams.dateRange" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-container">
      <div class="table-toolbar">
        <span class="table-title">充值记录</span>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="recharge_no" label="充值流水号" width="180" />
        <el-table-column prop="account.account_no" label="账户编号" width="160">
          <template #default="{ row }">
            {{ row.account?.account_no || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="account.account_name" label="账户名称" min-width="120">
          <template #default="{ row }">
            {{ row.account?.account_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="充值金额" width="120" align="right">
          <template #default="{ row }">
            <span class="text-money positive">+{{ formatMoney(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_before" label="充值前余额" width="120" align="right">
          <template #default="{ row }">{{ formatMoney(row.balance_before) }}</template>
        </el-table-column>
        <el-table-column prop="balance_after" label="充值后余额" width="120" align="right">
          <template #default="{ row }">{{ formatMoney(row.balance_after) }}</template>
        </el-table-column>
        <el-table-column prop="payment_method" label="支付方式" width="100" align="center">
          <template #default="{ row }">
            {{ ['', '现金', '转账', '微信', '支付宝'][row.payment_method] || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="operator_name" label="操作员" width="100" />
        <el-table-column prop="created_at" label="充值时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="fetchList" @current-change="fetchList" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import { formatMoney, formatDateTime } from '@/utils/format'
import * as accountApi from '@/api/modules/account'
import * as merchantApi from '@/api/modules/merchant'
import type { Recharge } from '@/types/finance'
import type { Merchant } from '@/types/business'

const loading = ref(false)
const tableData = ref<Recharge[]>([])
const total = ref(0)
const merchantList = ref<Merchant[]>([])

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  merchant_id: undefined as number | undefined,
  dateRange: [] as string[]
})

async function fetchList() {
  loading.value = true
  try {
    const params: any = {
      page: queryParams.page,
      pageSize: queryParams.pageSize
    }
    if (queryParams.keyword) {
      params.keyword = queryParams.keyword
    }
    if (queryParams.merchant_id) {
      params.merchant_id = queryParams.merchant_id
    }
    if (queryParams.dateRange && queryParams.dateRange.length === 2) {
      params.start_date = queryParams.dateRange[0]
      params.end_date = queryParams.dateRange[1]
    }
    const res = await accountApi.getRechargeList(params)
    tableData.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

async function fetchMerchantList() {
  merchantList.value = await merchantApi.getAllMerchants()
}

function handleQuery() {
  queryParams.page = 1
  fetchList()
}

function handleReset() {
  queryParams.keyword = ''
  queryParams.merchant_id = undefined
  queryParams.dateRange = []
  handleQuery()
}

onMounted(() => {
  fetchList()
  fetchMerchantList()
})
</script>
