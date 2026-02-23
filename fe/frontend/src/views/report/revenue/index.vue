<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="统计周期">
          <el-date-picker v-model="queryParams.dateRange" type="monthrange" start-placeholder="开始月份" end-placeholder="结束月份" value-format="YYYY-MM" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">查询</el-button>
          <el-button type="success" :icon="Download" @click="handleExport">导出</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="8">
        <el-card shadow="never">
          <div class="stat-value success">{{ formatMoney(stats.totalRecharge) }}</div>
          <div class="stat-label">总充值金额</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <div class="stat-value warning">{{ formatMoney(stats.totalElectric) }}</div>
          <div class="stat-label">电费收入</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <div class="stat-value primary">{{ formatMoney(stats.totalWater) }}</div>
          <div class="stat-label">水费收入</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>收入趋势</span></template>
          <div ref="trendChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>收入构成</span></template>
          <div ref="pieChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import * as echarts from 'echarts'
import { Search, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatMoney } from '@/utils/format'

const trendChartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

const stats = ref({ totalRecharge: 58600, totalElectric: 32800, totalWater: 12500 })
const queryParams = reactive({ dateRange: [] as string[] })

function initCharts() {
  if (trendChartRef.value) {
    trendChart = echarts.init(trendChartRef.value)
    trendChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['充值', '电费', '水费'], bottom: 0 },
      xAxis: { type: 'category', data: ['1月', '2月', '3月', '4月', '5月', '6月'] },
      yAxis: { type: 'value', name: '元' },
      series: [
        { name: '充值', type: 'line', data: [8000, 9200, 10100, 9800, 11200, 10300], smooth: true },
        { name: '电费', type: 'line', data: [5200, 5800, 5400, 5600, 5200, 5600], smooth: true },
        { name: '水费', type: 'line', data: [2100, 2000, 2200, 1900, 2100, 2200], smooth: true }
      ]
    })
  }

  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    pieChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { bottom: 0 },
      series: [{
        type: 'pie', radius: ['40%', '70%'],
        data: [
          { value: 32800, name: '电费收入', itemStyle: { color: '#E6A23C' } },
          { value: 12500, name: '水费收入', itemStyle: { color: '#409EFF' } }
        ]
      }]
    })
  }
}

function handleQuery() { initCharts() }
function handleExport() { ElMessage.info('功能开发中') }

onMounted(() => {
  initCharts()
  window.addEventListener('resize', () => { trendChart?.resize(); pieChart?.resize() })
})
</script>

<style scoped>
.stat-value { font-size: 28px; font-weight: 600; text-align: center; }
.stat-value.success { color: #67C23A; }
.stat-value.warning { color: #E6A23C; }
.stat-value.primary { color: #409EFF; }
.stat-label { font-size: 14px; color: #909399; text-align: center; margin-top: 8px; }
</style>
