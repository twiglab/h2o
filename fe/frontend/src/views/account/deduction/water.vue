<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="账户编号">
          <el-input v-model="queryParams.accountNo" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="扣费时间">
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
        <span class="table-title">水费扣费记录</span>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="deductionNo" label="扣费流水号" width="180" />
        <el-table-column prop="accountNo" label="账户编号" width="140" />
        <el-table-column prop="meterNo" label="水表号" width="140" />
        <el-table-column prop="consumption" label="用水量(吨)" width="120" align="right" />
        <el-table-column prop="unitPrice" label="单价" width="100" align="right" />
        <el-table-column prop="amount" label="扣费金额" width="120" align="right">
          <template #default="{ row }">
            <span class="text-money negative">-{{ formatMoney(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balanceAfter" label="扣费后余额" width="120" align="right">
          <template #default="{ row }">{{ formatMoney(row.balanceAfter) }}</template>
        </el-table-column>
        <el-table-column prop="deductionTime" label="扣费时间" width="170" />
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleQuery" @current-change="handleQuery" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import { formatMoney } from '@/utils/format'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)

const queryParams = reactive({ page: 1, pageSize: 10, accountNo: '', dateRange: [] as string[] })

async function fetchList() {
  loading.value = true
  try { tableData.value = []; total.value = 0 } finally { loading.value = false }
}

function handleQuery() { queryParams.page = 1; fetchList() }
function handleReset() { queryParams.accountNo = ''; queryParams.dateRange = []; handleQuery() }

onMounted(() => fetchList())
</script>
