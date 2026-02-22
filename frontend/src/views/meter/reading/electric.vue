<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="表号">
          <el-input v-model="queryParams.meterNo" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="抄表时间">
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
        <span class="table-title">电表抄表记录</span>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="meterNo" label="表号" width="140" />
        <el-table-column prop="shopName" label="店铺" min-width="120" />
        <el-table-column prop="readingValue" label="读数(kWh)" width="120" align="right" />
        <el-table-column prop="readingTime" label="抄表时间" width="170" />
        <el-table-column prop="collectType" label="采集方式" width="100" align="center">
          <template #default="{ row }">
            <template v-if="row.collectType != null">
              <el-tag :type="row.collectType === 1 ? 'success' : 'warning'" size="small">{{ row.collectType === 1 ? '自动' : '人工' }}</el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <template v-if="row.status != null">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">{{ row.status === 1 ? '正常' : '异常' }}</el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
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

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)

const queryParams = reactive({ page: 1, pageSize: 10, meterNo: '', dateRange: [] as string[] })

async function fetchList() {
  loading.value = true
  try { tableData.value = []; total.value = 0 } finally { loading.value = false }
}

function handleQuery() { queryParams.page = 1; fetchList() }
function handleReset() { queryParams.meterNo = ''; queryParams.dateRange = []; handleQuery() }

onMounted(() => fetchList())
</script>
