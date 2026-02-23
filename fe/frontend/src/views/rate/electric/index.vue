<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="关键词">
          <el-input v-model="queryParams.keyword" placeholder="编码/名称" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="计费模式">
          <el-select v-model="queryParams.calc_mode" placeholder="请选择" clearable style="width: 120px">
            <el-option label="固定单价" :value="1" />
            <el-option label="分时电价" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 100px">
            <el-option label="启用" :value="1" />
            <el-option label="停用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="table-container">
      <div class="table-toolbar">
        <span class="table-title">电费费率</span>
        <el-button type="primary" :icon="Plus" @click="handleAdd">新增</el-button>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="rate_code" label="费率编码" width="170" show-overflow-tooltip />
        <el-table-column prop="rate_name" label="费率名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="scope" label="适用范围" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.scope === 1" type="primary" size="small">全局</el-tag>
            <el-tag v-else type="warning" size="small">{{ row.merchant_name || '指定商户' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="calc_mode" label="计费模式" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.calc_mode === 1 ? 'info' : 'success'" size="small">
              {{ row.calc_mode === 1 ? '固定单价' : '分时电价' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="unit_price" label="单价(元/kWh)" width="130" align="right">
          <template #default="{ row }">
            <template v-if="row.calc_mode === 1">
              {{ formatPrice(row.unit_price) }}
            </template>
            <el-popover v-else placement="right" :width="320" trigger="hover">
              <template #reference>
                <el-button type="primary" text size="small">查看时段</el-button>
              </template>
              <el-table :data="row.tou_details || []" size="small" border>
                <el-table-column prop="period_name" label="时段" width="70" />
                <el-table-column prop="start_time" label="开始" width="60" />
                <el-table-column prop="end_time" label="结束" width="60" />
                <el-table-column prop="unit_price" label="单价" width="80" align="right">
                  <template #default="{ row: touRow }">
                    {{ formatPrice(touRow.unit_price) }}
                  </template>
                </el-table-column>
              </el-table>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="effective_date" label="生效日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.effective_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" text size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增费率' : '编辑费率'"
      width="750px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-divider content-position="left">基本信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="费率编码" prop="rate_code">
              <el-input v-model="formData.rate_code" placeholder="留空自动生成" :disabled="dialogType === 'edit'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="费率名称" prop="rate_name">
              <el-input v-model="formData.rate_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="适用范围" prop="scope">
              <el-radio-group v-model="formData.scope" @change="handleScopeChange">
                <el-radio :value="1">全局</el-radio>
                <el-radio :value="2">指定商户</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="formData.scope === 2">
            <el-form-item label="所属商户" prop="merchant_id">
              <el-select v-model="formData.merchant_id" placeholder="请选择" filterable style="width: 100%">
                <el-option
                  v-for="item in merchantList"
                  :key="item.id"
                  :label="item.merchant_name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="生效日期" prop="effective_date">
              <el-date-picker
                v-model="formData.effective_date"
                type="date"
                placeholder="请选择"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">启用</el-radio>
                <el-radio :value="0">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">计费配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="计费模式" prop="calc_mode">
              <el-radio-group v-model="formData.calc_mode" @change="handleCalcModeChange">
                <el-radio :value="1">固定单价</el-radio>
                <el-radio :value="2">分时电价</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="formData.calc_mode === 1">
            <el-form-item label="单价" prop="unit_price">
              <el-input-number
                v-model="formData.unit_price"
                :min="0"
                :precision="4"
                :step="0.01"
                placeholder="元/kWh"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 分时电价配置 -->
        <template v-if="formData.calc_mode === 2">
          <el-form-item label="分时时段" :error="touError">
            <div style="width: 100%">
              <el-table :data="formData.tou_details" border size="small" style="margin-bottom: 10px">
                <el-table-column prop="period_name" label="时段名称" width="120">
                  <template #default="{ row }">
                    <el-input v-model="row.period_name" placeholder="如：尖峰" size="small" />
                  </template>
                </el-table-column>
                <el-table-column prop="start_time" label="开始时间" width="130">
                  <template #default="{ row }">
                    <el-time-select
                      v-model="row.start_time"
                      :max-time="row.end_time"
                      placeholder="开始"
                      start="00:00"
                      step="00:30"
                      end="23:30"
                      size="small"
                      style="width: 100%"
                    />
                  </template>
                </el-table-column>
                <el-table-column prop="end_time" label="结束时间" width="130">
                  <template #default="{ row }">
                    <el-time-select
                      v-model="row.end_time"
                      :min-time="row.start_time"
                      placeholder="结束"
                      start="00:30"
                      step="00:30"
                      end="24:00"
                      size="small"
                      style="width: 100%"
                    />
                  </template>
                </el-table-column>
                <el-table-column prop="unit_price" label="单价(元/kWh)" width="140">
                  <template #default="{ row }">
                    <el-input-number
                      v-model="row.unit_price"
                      :min="0"
                      :precision="4"
                      :step="0.01"
                      size="small"
                      :controls="false"
                      style="width: 100%"
                    />
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="70" align="center">
                  <template #default="{ $index }">
                    <el-button type="danger" text size="small" @click="removeTouDetail($index)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-button type="primary" text :icon="Plus" @click="addTouDetail">添加时段</el-button>
              <span style="margin-left: 16px; color: #909399; font-size: 12px">
                当前覆盖: {{ coveredHours }} / 24 小时
              </span>
            </div>
          </el-form-item>
        </template>

        <!-- 服务费配置 -->
        <el-divider content-position="left">服务费配置（可选）</el-divider>
        <el-form-item label="服务费">
          <div style="width: 100%">
            <el-table :data="formData.service_fees" border size="small" style="margin-bottom: 10px">
              <el-table-column prop="fee_name" label="费用名称" min-width="150">
                <template #default="{ row }">
                  <el-input v-model="row.fee_name" placeholder="如：电损费" size="small" />
                </template>
              </el-table-column>
              <el-table-column prop="fee_type" label="收费类型" width="140">
                <template #default="{ row }">
                  <el-select v-model="row.fee_type" size="small" style="width: 100%">
                    <el-option label="百分比" :value="2" />
                    <el-option label="固定金额" :value="1" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column prop="fee_value" label="费用值" width="140">
                <template #default="{ row }">
                  <el-input-number
                    v-model="row.fee_value"
                    :min="0"
                    :precision="row.fee_type === 2 ? 2 : 4"
                    :step="row.fee_type === 2 ? 1 : 0.01"
                    size="small"
                    :controls="false"
                    style="width: 100%"
                  />
                </template>
              </el-table-column>
              <el-table-column label="单位" width="80" align="center">
                <template #default="{ row }">
                  {{ row.fee_type === 2 ? '%' : '元' }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="70" align="center">
                <template #default="{ $index }">
                  <el-button type="danger" text size="small" @click="removeServiceFee($index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-button type="primary" text :icon="Plus" @click="addServiceFee">添加服务费</el-button>
          </div>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="备注" prop="remark">
              <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import * as rateApi from '@/api/modules/rate'
import * as merchantApi from '@/api/modules/merchant'
import type { ElectricRate, ElectricRateForm, TOUDetail, ServiceFee } from '@/types/rate'
import type { Merchant } from '@/types/business'

const loading = ref(false)
const tableData = ref<ElectricRate[]>([])
const total = ref(0)
const merchantList = ref<Merchant[]>([])

// 格式化单价（后端可能返回字符串或数字）
function formatPrice(price: number | string | null | undefined): string {
  if (price === null || price === undefined) return '-'
  const num = typeof price === 'string' ? parseFloat(price) : price
  return isNaN(num) ? '-' : num.toFixed(4)
}

// 格式化日期（后端可能返回完整的ISO格式）
function formatDate(date: string | null | undefined): string {
  if (!date) return '-'
  // 处理 ISO 格式或 YYYY-MM-DD 格式
  return date.substring(0, 10)
}

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  calc_mode: undefined as number | undefined,
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const touError = ref('')

// 表单数据
const getDefaultFormData = (): ElectricRateForm => ({
  rate_code: '',
  rate_name: '',
  scope: 1,
  merchant_id: undefined,
  calc_mode: 1,
  unit_price: 0.85,
  effective_date: new Date().toISOString().split('T')[0],
  expire_date: '',
  status: 1,
  remark: '',
  tou_details: [],
  service_fees: []
})

const formData = ref<ElectricRateForm>(getDefaultFormData())

// 表单验证规则
const formRules = reactive<FormRules>({
  rate_name: [
    { required: true, message: '请输入费率名称', trigger: 'blur' },
    { max: 100, message: '费率名称不能超过100个字符', trigger: 'blur' }
  ],
  scope: [
    { required: true, message: '请选择适用范围', trigger: 'change' }
  ],
  merchant_id: [
    {
      required: true,
      message: '请选择所属商户',
      trigger: 'change',
      validator: (_rule, _value, callback) => {
        if (formData.value.scope === 2 && !formData.value.merchant_id) {
          callback(new Error('请选择所属商户'))
        } else {
          callback()
        }
      }
    }
  ],
  calc_mode: [
    { required: true, message: '请选择计费模式', trigger: 'change' }
  ],
  unit_price: [
    {
      required: true,
      message: '请输入单价',
      trigger: 'blur',
      validator: (_rule, _value, callback) => {
        if (formData.value.calc_mode === 1 && (!formData.value.unit_price || formData.value.unit_price <= 0)) {
          callback(new Error('请输入有效的单价'))
        } else {
          callback()
        }
      }
    }
  ],
  effective_date: [
    { required: true, message: '请选择生效日期', trigger: 'change' }
  ]
})

// 计算分时时段覆盖的小时数
const coveredHours = computed(() => {
  if (!formData.value.tou_details || formData.value.tou_details.length === 0) {
    return 0
  }
  let totalMinutes = 0
  for (const detail of formData.value.tou_details) {
    if (detail.start_time && detail.end_time) {
      const [startH, startM] = detail.start_time.split(':').map(Number)
      const [endH, endM] = detail.end_time.split(':').map(Number)
      let startMinutes = startH * 60 + startM
      let endMinutes = endH * 60 + endM
      // 处理跨天情况
      if (endMinutes <= startMinutes) {
        endMinutes += 24 * 60
      }
      totalMinutes += endMinutes - startMinutes
    }
  }
  return Math.round(totalMinutes / 60 * 10) / 10
})

// 验证分时时段是否覆盖24小时
function validateTouDetails(): boolean {
  if (formData.value.calc_mode !== 2) {
    return true
  }

  const details = formData.value.tou_details || []
  if (details.length === 0) {
    touError.value = '请至少添加一个分时时段'
    return false
  }

  // 检查每个时段是否填写完整
  for (let i = 0; i < details.length; i++) {
    const d = details[i]
    if (!d.period_name || !d.start_time || !d.end_time || d.unit_price === undefined || d.unit_price === null) {
      touError.value = `第${i + 1}行时段信息不完整`
      return false
    }
    if (d.unit_price < 0) {
      touError.value = `第${i + 1}行单价不能为负数`
      return false
    }
  }

  // 验证是否覆盖完整24小时
  if (coveredHours.value < 24) {
    touError.value = `时段配置不完整，当前仅覆盖${coveredHours.value}小时，需覆盖完整24小时`
    return false
  }

  // 检查时段重叠（简单检查）
  const intervals: Array<{ start: number; end: number }> = []
  for (const d of details) {
    const [startH, startM] = d.start_time.split(':').map(Number)
    const [endH, endM] = d.end_time.split(':').map(Number)
    let start = startH * 60 + startM
    let end = endH * 60 + endM
    if (end <= start) {
      end += 24 * 60
    }
    intervals.push({ start, end })
  }

  // 排序并检查重叠
  intervals.sort((a, b) => a.start - b.start)
  for (let i = 1; i < intervals.length; i++) {
    if (intervals[i].start < intervals[i - 1].end) {
      touError.value = '存在时段重叠，请检查配置'
      return false
    }
  }

  touError.value = ''
  return true
}

// 获取费率列表
async function fetchList() {
  loading.value = true
  try {
    const res = await rateApi.getRateList(queryParams)
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

// 获取商户列表
async function fetchMerchantList() {
  try {
    const res = await merchantApi.getAllMerchants()
    merchantList.value = res
  } catch {
    // 如果 getAllMerchants 接口不存在，使用分页接口
    const res = await merchantApi.getMerchantList({ page: 1, pageSize: 1000, status: 1 })
    merchantList.value = res.list
  }
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.keyword = ''
  queryParams.calc_mode = undefined
  queryParams.status = undefined
  handleQuery()
}

// 生成费率编码
function generateRateCode(): string {
  const now = new Date()
  const dateStr = now.toISOString().slice(0, 10).replace(/-/g, '')
  const randomStr = Math.random().toString(36).substring(2, 6).toUpperCase()
  return `RATE${dateStr}${randomStr}`
}

// 适用范围变化
function handleScopeChange(val: number) {
  if (val === 1) {
    formData.value.merchant_id = undefined
  }
}

// 计费模式变化
function handleCalcModeChange(val: number) {
  if (val === 1) {
    formData.value.tou_details = []
    formData.value.unit_price = 0.85
  } else {
    formData.value.unit_price = undefined
    // 添加默认分时时段
    if (!formData.value.tou_details || formData.value.tou_details.length === 0) {
      formData.value.tou_details = [
        { period_name: '尖峰', start_time: '10:00', end_time: '12:00', unit_price: 1.2 },
        { period_name: '高峰', start_time: '08:00', end_time: '10:00', unit_price: 0.95 },
        { period_name: '高峰', start_time: '12:00', end_time: '17:00', unit_price: 0.95 },
        { period_name: '高峰', start_time: '19:00', end_time: '22:00', unit_price: 0.95 },
        { period_name: '平段', start_time: '17:00', end_time: '19:00', unit_price: 0.65 },
        { period_name: '平段', start_time: '07:00', end_time: '08:00', unit_price: 0.65 },
        { period_name: '低谷', start_time: '22:00', end_time: '24:00', unit_price: 0.35 },
        { period_name: '低谷', start_time: '00:00', end_time: '07:00', unit_price: 0.35 }
      ]
    }
  }
  touError.value = ''
}

// 添加分时时段
function addTouDetail() {
  if (!formData.value.tou_details) {
    formData.value.tou_details = []
  }
  formData.value.tou_details.push({
    period_name: '',
    start_time: '',
    end_time: '',
    unit_price: 0
  } as TOUDetail)
}

// 删除分时时段
function removeTouDetail(index: number) {
  formData.value.tou_details?.splice(index, 1)
}

// 添加服务费
function addServiceFee() {
  if (!formData.value.service_fees) {
    formData.value.service_fees = []
  }
  formData.value.service_fees.push({
    fee_name: '',
    fee_type: 2, // 默认百分比
    fee_value: 0
  } as ServiceFee)
}

// 删除服务费
function removeServiceFee(index: number) {
  formData.value.service_fees?.splice(index, 1)
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = getDefaultFormData()
  touError.value = ''
  dialogVisible.value = true
}

// 编辑
async function handleEdit(row: ElectricRate) {
  dialogType.value = 'edit'
  touError.value = ''
  // 获取详情以获取完整的分时配置和服务费
  try {
    const detail = await rateApi.getRateDetail(row.id)
    // 处理单价格式（后端可能返回字符串）
    const unitPrice = typeof detail.unit_price === 'string'
      ? parseFloat(detail.unit_price)
      : detail.unit_price
    // 处理日期格式（后端可能返回完整ISO格式）
    const effectiveDate = detail.effective_date?.substring(0, 10) || ''
    const expireDate = detail.expire_date?.substring(0, 10) || ''
    // 处理分时详情中的单价格式
    const touDetails = (detail.tou_details || []).map(tou => ({
      ...tou,
      unit_price: typeof tou.unit_price === 'string' ? parseFloat(tou.unit_price) : tou.unit_price
    }))
    // 处理服务费中的值格式
    const serviceFees = (detail.service_fees || []).map(fee => ({
      ...fee,
      fee_value: typeof fee.fee_value === 'string' ? parseFloat(fee.fee_value) : fee.fee_value
    }))

    formData.value = {
      id: detail.id,
      rate_code: detail.rate_code,
      rate_name: detail.rate_name,
      scope: detail.scope,
      merchant_id: detail.merchant_id,
      calc_mode: detail.calc_mode,
      unit_price: unitPrice,
      effective_date: effectiveDate,
      expire_date: expireDate,
      status: detail.status,
      remark: detail.remark,
      tou_details: touDetails,
      service_fees: serviceFees
    }
    dialogVisible.value = true
  } catch {
    ElMessage.error('获取费率详情失败')
  }
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return

  // 先验证分时时段
  if (!validateTouDetails()) {
    return
  }

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        // 构建提交数据
        const submitData = { ...formData.value }

        // 如果是新增且没有填写编码，自动生成
        if (dialogType.value === 'add' && !submitData.rate_code) {
          submitData.rate_code = generateRateCode()
        }

        // 过滤掉空的服务费
        submitData.service_fees = (submitData.service_fees || []).filter(
          f => f.fee_name && f.fee_value !== undefined && f.fee_value !== null
        )

        if (dialogType.value === 'add') {
          await rateApi.createRate(submitData)
          ElMessage.success('新增成功')
        } else {
          await rateApi.updateRate(submitData.id!, submitData)
          ElMessage.success('编辑成功')
        }
        dialogVisible.value = false
        fetchList()
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 删除
async function handleDelete(row: ElectricRate) {
  try {
    await ElMessageBox.confirm(`确定要删除费率"${row.rate_name}"吗？`, '提示', { type: 'warning' })
    await rateApi.deleteRate(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error: any) {
    // 如果是用户取消操作，不显示错误
    if (error !== 'cancel' && error?.message !== 'cancel') {
      // 后端可能返回被引用无法删除的错误
      if (error?.message?.includes('使用中') || error?.message?.includes('引用')) {
        ElMessage.warning(error.message)
      }
    }
  }
}

onMounted(() => {
  fetchList()
  fetchMerchantList()
})
</script>
