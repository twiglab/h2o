<template>
  <div class="page-container">
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #67C23A 0%, #91d5a2 100%)">
            <el-icon size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value success">{{ formatMoney(stats.totalBalance) }}</div>
            <div class="stat-label">账户总余额</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #409EFF 0%, #36D1DC 100%)">
            <el-icon size="24"><Plus /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value primary">{{ formatMoney(stats.todayRecharge) }}</div>
            <div class="stat-label">今日充值</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #E6A23C 0%, #f7d794 100%)">
            <el-icon size="24"><Lightning /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value warning">{{ formatMoney(stats.todayDeduction) }}</div>
            <div class="stat-label">今日消费</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #909399 0%, #b4b4b4 100%)">
            <el-icon size="24"><Connection /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.onlineMeterCount }} / {{ stats.electricMeterCount }}</div>
            <div class="stat-label">在线电表</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 快捷操作 + 今日统计 -->
    <el-row :gutter="16" class="quick-row">
      <el-col :xs="24" :lg="16">
        <el-card shadow="never" class="quick-card">
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <div class="action-item" @click="$router.push('/customer/recharge')">
              <div class="action-icon" style="background-color: #e6f7ff">
                <el-icon color="#1890ff"><CreditCard /></el-icon>
              </div>
              <span>账户充值</span>
            </div>
            <div class="action-item" @click="$router.push('/customer/merchant')">
              <div class="action-icon" style="background-color: #f6ffed">
                <el-icon color="#52c41a"><Plus /></el-icon>
              </div>
              <span>新增商户</span>
            </div>
            <div class="action-item" @click="$router.push('/electric/meter')">
              <div class="action-icon" style="background-color: #fffbe6">
                <el-icon color="#faad14"><Lightning /></el-icon>
              </div>
              <span>电表管理</span>
            </div>
            <div class="action-item" @click="$router.push('/electric/deduction')">
              <div class="action-icon" style="background-color: #fff1f0">
                <el-icon color="#f5222d"><Document /></el-icon>
              </div>
              <span>扣费记录</span>
            </div>
            <div class="action-item" @click="$router.push('/electric/rate')">
              <div class="action-icon" style="background-color: #f0f5ff">
                <el-icon color="#2f54eb"><Setting /></el-icon>
              </div>
              <span>费率配置</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card shadow="never" class="today-card">
          <template #header>
            <span>今日采集</span>
          </template>
          <div class="today-stats">
            <div class="today-item">
              <div class="today-value">{{ formatNumber(stats.todayReadingCount) }}</div>
              <div class="today-label">采集次数</div>
            </div>
            <div class="today-item">
              <div class="today-value">{{ formatNumber(stats.merchantCount) }}</div>
              <div class="today-label">账户数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>用电量趋势（近30天）</span>
              <el-radio-group v-model="consumptionType" size="small">
                <el-radio-button value="electric">电量</el-radio-button>
                <el-radio-button value="water">水量</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="consumptionChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>充值/消费统计（近30天）</span>
            </div>
          </template>
          <div ref="financeChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 预警信息 -->
    <el-row :gutter="16" class="detail-row">
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>欠费账户</span>
              <el-button type="primary" text size="small" @click="$router.push('/customer/account')">
                查看更多
              </el-button>
            </div>
          </template>
          <el-table :data="lowBalanceAccounts" stripe size="small" max-height="200">
            <el-table-column prop="account_no" label="账户编号" min-width="130" />
            <el-table-column prop="account_name" label="账户名称" min-width="140" show-overflow-tooltip />
            <el-table-column prop="balance" label="当前余额" min-width="100" align="right">
              <template #default="{ row }">
                <span class="text-money negative">{{ formatMoney(row.balance) }}</span>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="lowBalanceAccounts.length === 0" description="暂无欠费账户" :image-size="60" />
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>离线电表</span>
              <el-button type="primary" text size="small" @click="$router.push('/electric/meter')">
                查看更多
              </el-button>
            </div>
          </template>
          <el-table :data="offlineMeters" stripe size="small" max-height="200">
            <el-table-column prop="meter_no" label="表号" min-width="130" />
            <el-table-column prop="location" label="商户" min-width="140" show-overflow-tooltip />
            <el-table-column prop="last_collect_at" label="最后采集" min-width="150">
              <template #default="{ row }">
                {{ formatDateTime(row.last_collect_at) || '从未采集' }}
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="offlineMeters.length === 0" description="暂无离线电表" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import { Wallet, Plus, CreditCard, Document, Setting, Connection, Lightning } from '@element-plus/icons-vue'
import { formatNumber, formatMoney, formatDateTime } from '@/utils/format'
import * as dashboardApi from '@/api/modules/dashboard'
import * as accountApi from '@/api/modules/account'
import * as meterApi from '@/api/modules/meter'

// 统计数据
const stats = ref({
  merchantCount: 0,
  totalBalance: 0,
  electricMeterCount: 0,
  waterMeterCount: 0,
  onlineMeterCount: 0,
  todayRecharge: 0,
  todayDeduction: 0,
  todayReadingCount: 0
})

// 用量图表
const consumptionType = ref<'electric' | 'water'>('electric')
const consumptionChartRef = ref<HTMLElement>()
let consumptionChart: echarts.ECharts | null = null
const consumptionData = ref<dashboardApi.ConsumptionStat[]>([])

// 财务图表
const financeMonth = ref('')
const financeChartRef = ref<HTMLElement>()
let financeChart: echarts.ECharts | null = null
const revenueData = ref<dashboardApi.RevenueStat[]>([])

// 余额预警账户
const lowBalanceAccounts = ref<any[]>([])

// 离线设备
const offlineMeters = ref<any[]>([])

// 获取统计数据
async function fetchStats() {
  try {
    const res = await dashboardApi.getDashboard()
    stats.value = {
      merchantCount: res.account?.total || 0,
      totalBalance: parseFloat(res.account?.total_balance) || 0,
      electricMeterCount: res.meter?.total || 0,
      waterMeterCount: 0,
      onlineMeterCount: res.meter?.online || 0,
      todayRecharge: parseFloat(res.today?.recharge) || 0,
      todayDeduction: parseFloat(res.today?.deduction) || 0,
      todayReadingCount: res.today?.reading_count || 0
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取余额预警数据（欠费账户）
async function fetchLowBalanceAccounts() {
  try {
    const res = await accountApi.getArrearsAccounts({ page: 1, pageSize: 5 })
    lowBalanceAccounts.value = res.list || []
  } catch (error) {
    console.error('获取欠费账户失败:', error)
    lowBalanceAccounts.value = []
  }
}

// 获取离线设备数据
async function fetchOfflineMeters() {
  try {
    const res = await meterApi.getElectricMeterList({ page: 1, pageSize: 5, onlineStatus: 0 })
    offlineMeters.value = (res.list || []).map((item: any) => ({
      meter_no: item.meter_no,
      meter_type: 'electric',
      location: item.merchant_name || '-',
      last_collect_at: item.last_collect_at
    }))
  } catch (error) {
    console.error('获取离线设备失败:', error)
    offlineMeters.value = []
  }
}

// 获取用量数据
async function fetchConsumptionData() {
  try {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30)

    const res = await dashboardApi.getConsumptionReport({
      start_date: startDate.toISOString().split('T')[0],
      end_date: endDate.toISOString().split('T')[0],
      group_by: 'day'
    })
    consumptionData.value = res || []
    updateConsumptionChart()
  } catch (error) {
    console.error('获取用量数据失败:', error)
  }
}

// 获取收入数据
async function fetchRevenueData() {
  try {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 30)

    const res = await dashboardApi.getRevenueReport({
      start_date: startDate.toISOString().split('T')[0],
      end_date: endDate.toISOString().split('T')[0],
      group_by: 'day'
    })
    revenueData.value = res || []
    updateFinanceChart()
  } catch (error) {
    console.error('获取收入数据失败:', error)
  }
}

// 初始化用量图表
function initConsumptionChart() {
  if (!consumptionChartRef.value) return
  consumptionChart = echarts.init(consumptionChartRef.value)
  updateConsumptionChart()
}

// 更新用量图表
function updateConsumptionChart() {
  if (!consumptionChart) return

  const dates = consumptionData.value.map(d => d.date)
  const consumption = consumptionData.value.map(d =>
    consumptionType.value === 'electric'
      ? parseFloat(d.electric_consumption) || 0
      : parseFloat(d.water_consumption) || 0
  )

  const option: EChartsOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    legend: {
      data: [consumptionType.value === 'electric' ? '电量(kWh)' : '水量(吨)'],
      bottom: 0
    },
    grid: {
      top: 20,
      left: 50,
      right: 20,
      bottom: 40
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        rotate: 45,
        formatter: (value: string) => value.slice(5)
      }
    },
    yAxis: {
      type: 'value',
      name: consumptionType.value === 'electric' ? 'kWh' : '吨'
    },
    series: [
      {
        name: consumptionType.value === 'electric' ? '电量(kWh)' : '水量(吨)',
        type: 'bar',
        data: consumption,
        itemStyle: {
          color: consumptionType.value === 'electric' ? '#E6A23C' : '#409EFF'
        }
      }
    ]
  }

  consumptionChart.setOption(option, true)
}

// 初始化财务图表
function initFinanceChart() {
  if (!financeChartRef.value) return
  financeChart = echarts.init(financeChartRef.value)
  updateFinanceChart()
}

// 更新财务图表
function updateFinanceChart() {
  if (!financeChart) return

  const dates = revenueData.value.map(d => d.date)
  const recharges = revenueData.value.map(d => parseFloat(d.recharge_amount) || 0)
  const deductions = revenueData.value.map(d => parseFloat(d.deduction_amount) || 0)

  const option: EChartsOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    legend: {
      data: ['充值', '消费'],
      bottom: 0
    },
    grid: {
      top: 20,
      left: 50,
      right: 20,
      bottom: 40
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        rotate: 45,
        formatter: (value: string) => value.slice(5)
      }
    },
    yAxis: {
      type: 'value',
      name: '元'
    },
    series: [
      {
        name: '充值',
        type: 'bar',
        stack: 'total',
        data: recharges,
        itemStyle: { color: '#67C23A' }
      },
      {
        name: '消费',
        type: 'bar',
        stack: 'total',
        data: deductions,
        itemStyle: { color: '#E6A23C' }
      }
    ]
  }

  financeChart.setOption(option, true)
}

// 监听用量类型变化
watch(consumptionType, () => {
  updateConsumptionChart()
})

// 窗口大小变化时重绘图表
function handleResize() {
  consumptionChart?.resize()
  financeChart?.resize()
}

onMounted(() => {
  fetchStats()
  fetchLowBalanceAccounts()
  fetchOfflineMeters()
  fetchConsumptionData()
  fetchRevenueData()

  initConsumptionChart()
  initFinanceChart()

  window.addEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.stat-row {
  margin-bottom: $spacing-lg;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  background: $bg-white;
  padding: $spacing-lg;
  border-radius: $border-radius-lg;
  transition: $transition-base;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-light;
  }

  .stat-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 56px;
    height: 56px;
    border-radius: $border-radius-lg;
    color: #fff;
  }

  .stat-content {
    flex: 1;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    font-family: $font-family-mono;
    line-height: 1.2;

    &.primary { color: $primary-color; }
    &.success { color: $success-color; }
    &.warning { color: $warning-color; }
  }

  .stat-label {
    font-size: $font-size-base;
    color: $text-secondary;
    margin-top: 4px;
  }
}

.quick-row {
  margin-bottom: $spacing-lg;
}

.quick-card {
  height: 100%;

  .quick-actions {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-lg;
  }

  .action-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: $spacing-sm;
    padding: $spacing-md $spacing-lg;
    cursor: pointer;
    border-radius: $border-radius-md;
    transition: $transition-base;

    &:hover {
      background-color: $bg-hover;
    }

    .action-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 48px;
      height: 48px;
      border-radius: $border-radius-md;
      font-size: 24px;
    }

    span {
      font-size: $font-size-base;
      color: $text-regular;
    }
  }
}

.today-card {
  height: 100%;

  .today-stats {
    display: flex;
    justify-content: space-around;
    padding: $spacing-md 0;
  }

  .today-item {
    text-align: center;
  }

  .today-value {
    font-size: 32px;
    font-weight: 600;
    font-family: $font-family-mono;
    color: $primary-color;
    line-height: 1.2;
  }

  .today-label {
    font-size: $font-size-base;
    color: $text-secondary;
    margin-top: 8px;
  }
}

.chart-row {
  margin-bottom: $spacing-lg;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-container {
  height: 300px;
}

.detail-row {
  margin-bottom: $spacing-lg;
}

.text-money {
  font-family: $font-family-mono;

  &.negative {
    color: var(--el-color-danger);
  }
}

@media (max-width: 768px) {
  .stat-card {
    padding: $spacing-md;

    .stat-icon {
      width: 44px;
      height: 44px;
    }

    .stat-value {
      font-size: 22px;
    }
  }

  .chart-container {
    height: 250px;
  }

  .quick-row {
    .el-col {
      margin-bottom: $spacing-md;
    }
  }

  .today-card .today-value {
    font-size: 24px;
  }
}
</style>
