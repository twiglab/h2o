<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="账户编号">
          <el-input v-model="queryParams.accountNo" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="电表号">
          <el-input v-model="queryParams.meterNo" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="扣费时间">
          <el-date-picker v-model="dateRange" type="daterange" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-container">
      <div class="table-toolbar">
        <span class="table-title">电费扣费记录</span>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="deduction_no" label="扣费流水号" width="180" />
        <el-table-column prop="account_no" label="账户编号" width="140" />
        <el-table-column prop="account_name" label="账户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="meter_no" label="电表号" width="120" />
        <el-table-column prop="consumption" label="用电量(kWh)" width="110" align="right">
          <template #default="{ row }">
            {{ parseFloat(row.consumption)?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="扣费金额" width="100" align="right">
          <template #default="{ row }">
            <span class="text-money negative">-{{ formatMoney(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="扣费后余额" width="110" align="right">
          <template #default="{ row }">{{ formatMoney(row.balance_after) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : row.status === 2 ? '部分' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deduction_time" label="扣费时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.deduction_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleViewDetail(row)">明细</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleQuery" @current-change="handleQuery" />
      </div>
    </div>

    <!-- 扣费明细弹窗 -->
    <el-dialog v-model="detailVisible" title="扣费明细" width="700px" destroy-on-close>
      <template v-if="currentDetail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="扣费流水号" :span="2">{{ currentDetail.deduction_no }}</el-descriptions-item>

          <el-descriptions-item label="商户名称">{{ currentDetail.merchant_name }}</el-descriptions-item>
          <el-descriptions-item label="店铺名称">{{ currentDetail.shop_name || '-' }}</el-descriptions-item>

          <el-descriptions-item label="账户编号">{{ currentDetail.account_no }}</el-descriptions-item>
          <el-descriptions-item label="账户名称">{{ currentDetail.account_name || '-' }}</el-descriptions-item>

          <el-descriptions-item label="电表号">{{ currentDetail.meter_no }}</el-descriptions-item>
          <el-descriptions-item label="倍率">{{ currentDetail.multiplier }}</el-descriptions-item>

          <el-descriptions-item label="起始读数">{{ parseFloat(currentDetail.start_reading)?.toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="终止读数">{{ parseFloat(currentDetail.end_reading)?.toFixed(2) }}</el-descriptions-item>

          <el-descriptions-item label="用电量(kWh)">
            <span class="text-primary">{{ parseFloat(currentDetail.consumption)?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="计费周期">
            {{ formatDateTime(currentDetail.period_start) }} ~ {{ formatDateTime(currentDetail.period_end) }}
          </el-descriptions-item>

          <el-descriptions-item label="费率名称">{{ currentDetail.rate_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="计费模式">
            <el-tag :type="currentDetail.calc_mode === 1 ? 'info' : 'warning'" size="small">
              {{ currentDetail.calc_mode === 1 ? '固定单价' : '分时电价' }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item v-if="currentDetail.calc_mode === 1" label="单价(元/kWh)">
            {{ parseFloat(currentDetail.unit_price)?.toFixed(4) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="currentDetail.calc_mode === 1" label="基础电费">
            ¥{{ formatMoney(currentDetail.base_amount) }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 分时电价明细 -->
        <template v-if="currentDetail.calc_mode === 2 && touDetails.length > 0">
          <el-divider content-position="left">分时计费明细</el-divider>
          <el-table :data="touDetails" border size="small">
            <el-table-column prop="period_name" label="时段" width="80" align="center" />
            <el-table-column label="时间范围" width="140" align="center">
              <template #default="{ row }">{{ row.start_time }} - {{ row.end_time }}</template>
            </el-table-column>
            <el-table-column prop="consumption" label="用量(kWh)" width="100" align="right">
              <template #default="{ row }">{{ row.consumption?.toFixed(4) }}</template>
            </el-table-column>
            <el-table-column prop="unit_price" label="单价(元)" width="100" align="right">
              <template #default="{ row }">{{ row.unit_price?.toFixed(4) }}</template>
            </el-table-column>
            <el-table-column prop="amount" label="金额(元)" width="100" align="right">
              <template #default="{ row }">{{ row.amount?.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </template>

        <el-descriptions :column="2" border style="margin-top: 16px">
          <el-descriptions-item label="基础电费">¥{{ formatMoney(currentDetail.base_amount) }}</el-descriptions-item>
          <el-descriptions-item label="服务费">¥{{ formatMoney(currentDetail.service_amount) }}</el-descriptions-item>

          <el-descriptions-item label="扣费金额">
            <span class="text-money negative" style="font-size: 16px; font-weight: bold;">-¥{{ formatMoney(currentDetail.amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="扣费状态">
            <el-tag :type="currentDetail.status === 1 ? 'success' : currentDetail.status === 2 ? 'warning' : 'danger'">
              {{ currentDetail.status === 1 ? '扣费成功' : currentDetail.status === 2 ? '部分扣费' : '扣费失败' }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="扣费前余额">¥{{ formatMoney(currentDetail.balance_before) }}</el-descriptions-item>
          <el-descriptions-item label="扣费后余额">¥{{ formatMoney(currentDetail.balance_after) }}</el-descriptions-item>

          <el-descriptions-item label="扣费时间" :span="2">{{ formatDateTime(currentDetail.deduction_time) }}</el-descriptions-item>

          <el-descriptions-item v-if="currentDetail.remark" label="备注" :span="2">{{ currentDetail.remark }}</el-descriptions-item>
        </el-descriptions>
      </template>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import { formatMoney, formatDateTime } from '@/utils/format'
import * as accountApi from '@/api/modules/account'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const dateRange = ref<string[]>([])

// 明细弹窗
const detailVisible = ref(false)
const currentDetail = ref<any>(null)

// 解析分时电价明细
const touDetails = computed(() => {
  if (!currentDetail.value?.tou_detail) return []
  try {
    const details = typeof currentDetail.value.tou_detail === 'string'
      ? JSON.parse(currentDetail.value.tou_detail)
      : currentDetail.value.tou_detail
    return details || []
  } catch {
    return []
  }
})

const queryParams = reactive({
  page: 1,
  pageSize: 10,
  accountNo: '',
  meterNo: '',
  startDate: '',
  endDate: ''
})

async function fetchList() {
  loading.value = true
  try {
    // 处理日期范围
    if (dateRange.value && dateRange.value.length === 2) {
      queryParams.startDate = dateRange.value[0]
      queryParams.endDate = dateRange.value[1]
    } else {
      queryParams.startDate = ''
      queryParams.endDate = ''
    }
    const res = await accountApi.getElectricDeductionList(queryParams)
    tableData.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取扣费记录失败:', error)
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function handleQuery() {
  queryParams.page = 1
  fetchList()
}

function handleReset() {
  queryParams.accountNo = ''
  queryParams.meterNo = ''
  dateRange.value = []
  handleQuery()
}

function handleViewDetail(row: any) {
  currentDetail.value = row
  detailVisible.value = true
}

onMounted(() => fetchList())
</script>

<style lang="scss" scoped>
.text-primary {
  color: var(--el-color-primary);
  font-weight: 500;
}

.text-money {
  font-family: 'Monaco', 'Menlo', monospace;

  &.negative {
    color: var(--el-color-danger);
  }
}

:deep(.el-descriptions__label) {
  width: 100px;
  font-weight: 500;
}

:deep(.el-divider__text) {
  font-weight: 500;
  color: var(--el-text-color-primary);
}
</style>
