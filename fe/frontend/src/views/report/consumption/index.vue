<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="统计类型">
          <el-radio-group v-model="queryParams.type">
            <el-radio-button value="electric">用电</el-radio-button>
            <el-radio-button value="water">用水</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="统计周期">
          <el-date-picker v-model="queryParams.dateRange" type="monthrange" start-placeholder="开始月份" end-placeholder="结束月份" value-format="YYYY-MM" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">查询</el-button>
          <el-button type="success" :icon="Download" @click="handleExport">导出</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <span>{{ queryParams.type === 'electric' ? '用电量' : '用水量' }}趋势</span>
          </template>
          <div ref="chartRef" style="height: 400px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <span>用量明细</span>
      </template>
      <el-table :data="tableData" stripe border>
        <el-table-column prop="period" label="统计周期" width="120" />
        <el-table-column prop="merchantName" label="商户" min-width="150" />
        <el-table-column prop="consumption" label="用量" width="120" align="right">
          <template #default="{ row }">
            {{ row.consumption }} {{ queryParams.type === 'electric' ? 'kWh' : '吨' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="费用(元)" width="120" align="right" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import * as echarts from 'echarts'
import { Search, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null
const tableData = ref<any[]>([])

const queryParams = reactive({
  type: 'electric' as 'electric' | 'water',
  dateRange: [] as string[]
})

function initChart() {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'] },
    yAxis: { type: 'value', name: queryParams.type === 'electric' ? 'kWh' : '吨' },
    series: [{ name: queryParams.type === 'electric' ? '用电量' : '用水量', type: 'bar', data: [820, 932, 901, 934, 1290, 1330, 1320, 1200, 1100, 980, 920, 850], itemStyle: { color: queryParams.type === 'electric' ? '#E6A23C' : '#409EFF' } }]
  })
}

function handleQuery() { initChart() }
function handleExport() { ElMessage.info('功能开发中') }

watch(() => queryParams.type, () => initChart())

onMounted(() => {
  initChart()
  window.addEventListener('resize', () => chart?.resize())
})
</script>
