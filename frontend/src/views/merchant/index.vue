<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="商户编号">
          <el-input v-model="queryParams.merchant_no" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="商户名称">
          <el-input v-model="queryParams.merchant_name" placeholder="请输入" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="商户类型">
          <el-select v-model="queryParams.merchant_type" placeholder="请选择" clearable style="width: 120px">
            <el-option label="企业" :value="1" />
            <el-option label="个人" :value="2" />
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
        <span class="table-title">商户列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'merchant:info:create'" @click="handleAdd">
            新增
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="merchant_no" label="商户编号" width="140" />
        <el-table-column prop="merchant_name" label="商户名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="merchant_type" label="商户类型" width="90" align="center">
          <template #default="{ row }">
            <template v-if="row.merchant_type != null">
              <el-tag :type="row.merchant_type === 1 ? 'primary' : 'success'" size="small">
                {{ row.merchant_type === 1 ? '企业' : '个人' }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="contact_name" label="联系人" width="100" />
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
        <el-table-column prop="shop_count" label="店铺数" width="80" align="center" />
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
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleView(row)">
              详情
            </el-button>
            <el-button type="primary" text size="small" v-permission="'merchant:info:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" text size="small" v-permission="'merchant:info:delete'" @click="handleDelete(row)">
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
      :title="dialogType === 'add' ? '新增商户' : dialogType === 'edit' ? '编辑商户' : '商户详情'"
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
            <el-form-item label="商户编号" prop="merchant_no">
              <el-input v-model="formData.merchant_no" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商户名称" prop="merchant_name">
              <el-input v-model="formData.merchant_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商户类型" prop="merchant_type">
              <el-radio-group v-model="formData.merchant_type">
                <el-radio :value="1">企业</el-radio>
                <el-radio :value="2">个人</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系人" prop="contact_name">
              <el-input v-model="formData.contact_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="contact_phone">
              <el-input v-model="formData.contact_phone" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">正常</el-radio>
                <el-radio :value="0">停用</el-radio>
              </el-radio-group>
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
      </el-form>
      <template #footer v-if="dialogType !== 'view'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import { validatePhone } from '@/utils/validate'
import * as merchantApi from '@/api/modules/merchant'
import type { Merchant, MerchantForm } from '@/types/business'

const loading = ref(false)
const tableData = ref<Merchant[]>([])
const total = ref(0)

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  merchant_no: '',
  merchant_name: '',
  merchant_type: undefined as number | undefined,
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit' | 'view'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref<MerchantForm>({
  merchant_name: '',
  merchant_type: 1,
  contact_name: '',
  contact_phone: '',
  status: 1,
  remark: ''
})

const formRules = reactive<FormRules>({
  merchant_name: [
    { required: true, message: '请输入商户名称', trigger: 'blur' },
    { max: 100, message: '商户名称不能超过100个字符', trigger: 'blur' }
  ],
  merchant_type: [
    { required: true, message: '请选择商户类型', trigger: 'change' }
  ],
  contact_phone: [
    { validator: validatePhone, trigger: 'blur' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await merchantApi.getMerchantList(queryParams)
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.merchant_no = ''
  queryParams.merchant_name = ''
  queryParams.merchant_type = undefined
  queryParams.status = undefined
  handleQuery()
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    merchant_name: '',
    merchant_type: 1,
    contact_name: '',
    contact_phone: '',
    status: 1,
    remark: ''
  }
  dialogVisible.value = true
}

// 查看
function handleView(row: Merchant) {
  dialogType.value = 'view'
  formData.value = {
    id: row.id,
    merchant_no: row.merchant_no,
    merchant_name: row.merchant_name,
    merchant_type: row.merchant_type,
    contact_name: row.contact_name,
    contact_phone: row.contact_phone,
    status: row.status,
    remark: row.remark
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Merchant) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    merchant_no: row.merchant_no,
    merchant_name: row.merchant_name,
    merchant_type: row.merchant_type,
    contact_name: row.contact_name,
    contact_phone: row.contact_phone,
    status: row.status,
    remark: row.remark
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
          await merchantApi.createMerchant(formData.value)
          ElMessage.success('新增成功')
        } else {
          await merchantApi.updateMerchant(formData.value.id!, formData.value)
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
async function handleDelete(row: Merchant) {
  await ElMessageBox.confirm(`确定要删除商户"${row.merchant_name}"吗？`, '提示', { type: 'warning' })
  await merchantApi.deleteMerchant(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>
