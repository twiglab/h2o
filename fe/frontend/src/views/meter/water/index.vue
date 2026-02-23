<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="表号">
          <el-input v-model="queryParams.meterNo" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="所属商户">
          <el-select v-model="queryParams.merchantId" placeholder="请选择" clearable filterable style="width: 180px">
            <el-option v-for="item in merchantList" :key="item.id" :label="item.merchantName" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="在线状态">
          <el-select v-model="queryParams.onlineStatus" placeholder="请选择" clearable style="width: 100px">
            <el-option label="在线" :value="1" />
            <el-option label="离线" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 100px">
            <el-option label="正常" :value="1" />
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
        <span class="table-title">水表列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'meter:water:create'" @click="handleAdd">
            新增
          </el-button>
          <el-button type="success" :icon="Download" v-permission="'meter:water:export'" @click="handleExport">
            导出
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="meterNo" label="表号" width="140" />
        <el-table-column prop="merchantName" label="所属商户" min-width="120" show-overflow-tooltip />
        <el-table-column prop="shopName" label="关联店铺" min-width="100" show-overflow-tooltip />
        <el-table-column prop="commAddr" label="通信地址" width="120" />
        <el-table-column prop="currentReading" label="当前读数" width="110" align="right">
          <template #default="{ row }">
            <span class="text-money">{{ row.currentReading?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="multiplier" label="倍率" width="80" align="center" />
        <el-table-column prop="onlineStatus" label="在线状态" width="90" align="center">
          <template #default="{ row }">
            <template v-if="row.onlineStatus != null">
              <el-tag :type="row.onlineStatus === 1 ? 'success' : 'danger'" size="small">
                {{ row.onlineStatus === 1 ? '在线' : '离线' }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="lastCollectAt" label="最后采集时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.lastCollectAt) || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <template v-if="row.status != null">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? '正常' : '停用' }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" v-permission="'meter:water:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" text size="small" v-permission="'meter:water:delete'" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleQuery"
          @current-change="handleQuery"
        />
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增水表' : '编辑水表'"
      width="650px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-divider content-position="left">基本信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="表号" prop="meterNo">
              <el-input v-model="formData.meterNo" placeholder="请输入" :disabled="dialogType === 'edit'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属商户" prop="merchantId">
              <el-select v-model="formData.merchantId" placeholder="请选择" filterable style="width: 100%">
                <el-option v-for="item in merchantList" :key="item.id" :label="item.merchantName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="关联店铺" prop="shopId">
              <el-select v-model="formData.shopId" placeholder="请选择" filterable clearable style="width: 100%">
                <el-option v-for="item in shopList" :key="item.id" :label="item.shopName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关联账户" prop="accountId">
              <el-select v-model="formData.accountId" placeholder="请选择" filterable clearable style="width: 100%">
                <el-option v-for="item in accountList" :key="item.id" :label="`${item.accountNo} - ${item.accountName || ''}`" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="费率" prop="rateId">
              <el-select v-model="formData.rateId" placeholder="请选择" filterable style="width: 100%">
                <el-option v-for="item in rateList" :key="item.id" :label="item.rateName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">正常</el-radio>
                <el-radio :value="0">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">通信配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="MQTT主题" prop="mqttTopic">
              <el-input v-model="formData.mqttTopic" placeholder="请输入MQTT主题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通信地址" prop="commAddr">
              <el-input v-model="formData.commAddr" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="通信协议" prop="protocol">
              <el-select v-model="formData.protocol" placeholder="请选择" style="width: 100%">
                <el-option label="Modbus RTU" value="modbus" />
                <el-option label="M-Bus" value="mbus" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">表计参数</el-divider>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="倍率" prop="multiplier">
              <el-input-number v-model="formData.multiplier" :min="1" :precision="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="初始读数" prop="initReading">
              <el-input-number v-model="formData.initReading" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="当前读数" prop="currentReading">
              <el-input-number v-model="formData.currentReading" :min="0" :precision="2" :disabled="dialogType === 'edit'" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">设备信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="品牌" prop="brand">
              <el-input v-model="formData.brand" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="型号" prop="model">
              <el-input v-model="formData.model" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="安装日期" prop="installDate">
              <el-date-picker v-model="formData.installDate" type="date" placeholder="请选择" style="width: 100%" value-format="YYYY-MM-DD" />
            </el-form-item>
          </el-col>
        </el-row>

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
import { ref, reactive, onMounted, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Download } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import * as merchantApi from '@/api/modules/merchant'
import type { Merchant, Shop, WaterMeter } from '@/types/business'
import type { Account } from '@/types/finance'

const loading = ref(false)
const tableData = ref<WaterMeter[]>([])
const total = ref(0)
const merchantList = ref<Merchant[]>([])
const shopList = ref<Shop[]>([])
const accountList = ref<Account[]>([])
const rateList = ref<any[]>([])

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  meterNo: '',
  merchantId: undefined as number | undefined,
  onlineStatus: undefined as number | undefined,
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref({
  id: undefined as number | undefined,
  meterNo: '',
  merchantId: undefined as number | undefined,
  shopId: undefined as number | undefined,
  accountId: undefined as number | undefined,
  rateId: undefined as number | undefined,
  mqttTopic: '',
  commAddr: '',
  protocol: 'modbus',
  multiplier: 1,
  initReading: 0,
  currentReading: 0,
  brand: '',
  model: '',
  installDate: '',
  status: 1,
  remark: ''
})

const formRules = reactive<FormRules>({
  meterNo: [
    { required: true, message: '请输入表号', trigger: 'blur' },
    { max: 32, message: '表号不能超过32个字符', trigger: 'blur' }
  ],
  merchantId: [
    { required: true, message: '请选择所属商户', trigger: 'change' }
  ]
})

// 监听商户变化
watch(() => formData.value.merchantId, async (merchantId) => {
  if (merchantId) {
    // TODO: 加载该商户下的店铺列表和账户列表
  } else {
    shopList.value = []
    accountList.value = []
  }
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    // TODO: 调用API
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 获取商户列表
async function fetchMerchantList() {
  const res = await merchantApi.getMerchantList({ page: 1, pageSize: 1000, status: 1 })
  merchantList.value = res.list
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.meterNo = ''
  queryParams.merchantId = undefined
  queryParams.onlineStatus = undefined
  queryParams.status = undefined
  handleQuery()
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    meterNo: '',
    merchantId: undefined,
    shopId: undefined,
    accountId: undefined,
    rateId: undefined,
    mqttTopic: '',
    commAddr: '',
    protocol: 'modbus',
    multiplier: 1,
    initReading: 0,
    currentReading: 0,
    brand: '',
    model: '',
    installDate: '',
    status: 1,
    remark: ''
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: WaterMeter) {
  dialogType.value = 'edit'
  formData.value = { ...row }
  dialogVisible.value = true
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        // TODO: 调用API
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '编辑成功')
        dialogVisible.value = false
        fetchList()
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 删除
async function handleDelete(row: WaterMeter) {
  await ElMessageBox.confirm(`确定要删除水表"${row.meterNo}"吗？`, '提示', { type: 'warning' })
  // TODO: 调用API
  ElMessage.success('删除成功')
  fetchList()
}

// 导出
async function handleExport() {
  try {
    // TODO: 调用API
    ElMessage.success('导出成功')
  } catch {
    ElMessage.error('导出失败')
  }
}

onMounted(() => {
  fetchList()
  fetchMerchantList()
})
</script>
