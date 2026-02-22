<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="账户编号">
          <el-input v-model="queryParams.account_no" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="账户名称">
          <el-input v-model="queryParams.account_name" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="所属商户">
          <el-select v-model="queryParams.merchant_id" placeholder="请选择" clearable filterable style="width: 180px">
            <el-option v-for="item in merchantList" :key="item.id" :label="item.merchant_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 100px">
            <el-option label="正常" :value="1" />
            <el-option label="欠费" :value="2" />
            <el-option label="冻结" :value="0" />
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
        <span class="table-title">账户列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'account:create'" @click="handleAdd">
            新增
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="account_no" label="账户编号" width="160" />
        <el-table-column prop="account_name" label="账户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="merchant_name" label="所属商户" min-width="150" show-overflow-tooltip />
        <el-table-column prop="balance" label="余额" width="120" align="right">
          <template #default="{ row }">
            <span :class="['text-money', row.balance < 0 ? 'negative' : 'positive']">
              {{ formatMoney(row.balance) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="total_recharge" label="累计充值" width="120" align="right">
          <template #default="{ row }">
            {{ formatMoney(row.total_recharge) }}
          </template>
        </el-table-column>
        <el-table-column prop="total_consumption" label="累计消费" width="120" align="right">
          <template #default="{ row }">
            {{ formatMoney(row.total_consumption) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <template v-if="row.status != null">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusLabel(row.status) }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="success" text size="small" v-permission="'account:recharge'" @click="handleRecharge(row)">
              充值
            </el-button>
            <el-button type="primary" text size="small" @click="handleView(row)">
              详情
            </el-button>
            <el-button type="primary" text size="small" v-permission="'account:update'" @click="handleEdit(row)">
              编辑
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
      :title="dialogType === 'add' ? '新增账户' : dialogType === 'edit' ? '编辑账户' : '账户详情'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        :disabled="dialogType === 'view'"
      >
        <el-row :gutter="20" v-if="dialogType !== 'add'">
          <el-col :span="12">
            <el-form-item label="账户编号">
              <el-input v-model="formData.account_no" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="账户名称" prop="account_name">
              <el-input v-model="formData.account_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属商户" prop="merchant_id">
              <el-select v-model="formData.merchant_id" placeholder="请选择" filterable style="width: 100%" :disabled="dialogType === 'edit'">
                <el-option v-for="item in merchantList" :key="item.id" :label="item.merchant_name" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="formData.status" placeholder="请选择" style="width: 100%">
                <el-option label="正常" :value="1" />
                <el-option label="冻结" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="备注" prop="remark">
              <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 详情模式显示更多信息 -->
        <template v-if="dialogType === 'view'">
          <el-divider content-position="left">账户统计</el-divider>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="可用余额">
              <span class="text-money" :class="formData.balance >= 0 ? 'positive' : 'negative'">
                {{ formatMoney(formData.balance) }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="累计充值">
              {{ formatMoney(formData.total_recharge) }}
            </el-descriptions-item>
            <el-descriptions-item label="累计消费">
              {{ formatMoney(formData.total_consumption) }}
            </el-descriptions-item>
          </el-descriptions>
        </template>
      </el-form>
      <template #footer v-if="dialogType !== 'view'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 充值弹窗 -->
    <el-dialog v-model="rechargeDialogVisible" title="账户充值" width="500px" destroy-on-close>
      <el-form ref="rechargeFormRef" :model="rechargeForm" :rules="rechargeRules" label-width="100px">
        <el-form-item label="账户编号">
          <el-input v-model="rechargeForm.account_no" disabled />
        </el-form-item>
        <el-form-item label="账户名称">
          <el-input v-model="rechargeForm.account_name" disabled />
        </el-form-item>
        <el-form-item label="当前余额">
          <span class="text-money" :class="rechargeForm.balance >= 0 ? 'positive' : 'negative'">
            {{ formatMoney(rechargeForm.balance) }}
          </span>
        </el-form-item>
        <el-form-item label="充值金额" prop="amount">
          <el-input-number v-model="rechargeForm.amount" :min="0.01" :max="999999" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="支付方式" prop="payment_method">
          <el-radio-group v-model="rechargeForm.payment_method">
            <el-radio :value="1">现金</el-radio>
            <el-radio :value="2">转账</el-radio>
            <el-radio :value="3">微信</el-radio>
            <el-radio :value="4">支付宝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="rechargeForm.remark" type="textarea" :rows="2" placeholder="请输入" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="rechargeLoading" @click="handleRechargeSubmit">确认充值</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { formatDateTime, formatMoney } from '@/utils/format'
import * as merchantApi from '@/api/modules/merchant'
import * as accountApi from '@/api/modules/account'
import type { Merchant } from '@/types/business'
import type { Account } from '@/types/finance'

const loading = ref(false)
const tableData = ref<Account[]>([])
const total = ref(0)
const merchantList = ref<Merchant[]>([])

// 状态
function getStatusLabel(status: number): string {
  const map: Record<number, string> = { 1: '正常', 2: '欠费', 0: '冻结' }
  return map[status] || '-'
}

function getStatusType(status: number | null | undefined): '' | 'success' | 'warning' | 'danger' {
  if (status == null) return ''
  const map: Record<number, '' | 'success' | 'warning' | 'danger'> = { 1: 'success', 2: 'warning', 0: 'danger' }
  return map[status] ?? ''
}

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  account_no: '',
  account_name: '',
  merchant_id: undefined as number | undefined,
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit' | 'view'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref({
  id: undefined as number | undefined,
  account_no: '',
  account_name: '',
  merchant_id: undefined as number | undefined,
  balance: 0,
  total_recharge: 0,
  total_consumption: 0,
  status: 1,
  remark: ''
})

const formRules = reactive<FormRules>({
  account_name: [
    { max: 100, message: '账户名称不能超过100个字符', trigger: 'blur' }
  ],
  merchant_id: [
    { required: true, message: '请选择所属商户', trigger: 'change' }
  ]
})

// 充值弹窗
const rechargeDialogVisible = ref(false)
const rechargeLoading = ref(false)
const rechargeFormRef = ref<FormInstance>()

const rechargeForm = ref({
  account_id: undefined as number | undefined,
  account_no: '',
  account_name: '',
  balance: 0,
  amount: 100,
  payment_method: 1,
  remark: ''
})

const rechargeRules = reactive<FormRules>({
  amount: [
    { required: true, message: '请输入充值金额', trigger: 'blur' }
  ],
  payment_method: [
    { required: true, message: '请选择支付方式', trigger: 'change' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await accountApi.getAccountList(queryParams)
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

// 获取商户列表
async function fetchMerchantList() {
  merchantList.value = await merchantApi.getAllMerchants()
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.account_no = ''
  queryParams.account_name = ''
  queryParams.merchant_id = undefined
  queryParams.status = undefined
  handleQuery()
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    account_no: '',
    account_name: '',
    merchant_id: undefined,
    balance: 0,
    total_recharge: 0,
    total_consumption: 0,
    status: 1,
    remark: ''
  }
  dialogVisible.value = true
}

// 查看
function handleView(row: Account) {
  dialogType.value = 'view'
  formData.value = {
    id: row.id,
    account_no: row.account_no,
    account_name: row.account_name || '',
    merchant_id: row.merchant_id,
    balance: row.balance,
    total_recharge: row.total_recharge,
    total_consumption: row.total_consumption,
    status: row.status,
    remark: row.remark || ''
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Account) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    account_no: row.account_no,
    account_name: row.account_name || '',
    merchant_id: row.merchant_id,
    balance: row.balance,
    total_recharge: row.total_recharge,
    total_consumption: row.total_consumption,
    status: row.status,
    remark: row.remark || ''
  }
  dialogVisible.value = true
}

// 提交
async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (dialogType.value === 'add') {
          await accountApi.createAccount(formData.value)
          ElMessage.success('新增成功')
        } else {
          await accountApi.updateAccount(formData.value.id!, formData.value)
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

// 充值
function handleRecharge(row: Account) {
  rechargeForm.value = {
    account_id: row.id,
    account_no: row.account_no,
    account_name: row.account_name || '',
    balance: row.balance,
    amount: 100,
    payment_method: 1,
    remark: ''
  }
  rechargeDialogVisible.value = true
}

// 提交充值
async function handleRechargeSubmit() {
  if (!rechargeFormRef.value) return
  await rechargeFormRef.value.validate(async (valid) => {
    if (valid) {
      rechargeLoading.value = true
      try {
        await accountApi.rechargeAccount(rechargeForm.value.account_id!, {
          amount: rechargeForm.value.amount,
          payment_method: rechargeForm.value.payment_method,
          remark: rechargeForm.value.remark
        })
        ElMessage.success('充值成功')
        rechargeDialogVisible.value = false
        fetchList()
      } finally {
        rechargeLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchList()
  fetchMerchantList()
})
</script>
