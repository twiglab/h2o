<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="店铺编号">
          <el-input v-model="queryParams.shop_no" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="店铺名称">
          <el-input v-model="queryParams.shop_name" placeholder="请输入" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="所属商户">
          <el-select v-model="queryParams.merchant_id" placeholder="请选择" clearable filterable style="width: 180px">
            <el-option v-for="item in merchantList" :key="item.id" :label="item.merchant_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="楼栋">
          <el-input v-model="queryParams.building" placeholder="请输入" clearable style="width: 100px" />
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
        <span class="table-title">店铺列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'merchant:shop:create'" @click="handleAdd">
            新增
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="shop_no" label="店铺编号" width="140" />
        <el-table-column prop="shop_name" label="店铺名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="merchant_name" label="所属商户" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.merchant?.merchant_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="location" label="位置" min-width="150">
          <template #default="{ row }">
            {{ [row.building, row.floor, row.room_no].filter(Boolean).join('-') || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="contact_name" label="联系人" width="100" />
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
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
            <el-button type="primary" text size="small" v-permission="'merchant:shop:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" text size="small" v-permission="'merchant:shop:delete'" @click="handleDelete(row)">
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
      :title="dialogType === 'add' ? '新增店铺' : dialogType === 'edit' ? '编辑店铺' : '店铺详情'"
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
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="店铺编号" prop="shop_no">
              <el-input v-model="formData.shop_no" placeholder="留空自动生成" :disabled="dialogType !== 'add'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="店铺名称" prop="shop_name">
              <el-input v-model="formData.shop_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属商户" prop="merchant_id">
              <el-select v-model="formData.merchant_id" placeholder="请选择" filterable style="width: 100%">
                <el-option v-for="item in merchantList" :key="item.id" :label="item.merchant_name" :value="item.id" />
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

        <el-divider content-position="left">位置信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="楼栋" prop="building">
              <el-input v-model="formData.building" placeholder="如: A栋" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="楼层" prop="floor">
              <el-input v-model="formData.floor" placeholder="如: 3F" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房号" prop="room_no">
              <el-input v-model="formData.room_no" placeholder="如: 301" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">联系信息</el-divider>
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
import * as shopApi from '@/api/modules/shop'
import type { Merchant, Shop, ShopForm } from '@/types/business'

const loading = ref(false)
const tableData = ref<Shop[]>([])
const total = ref(0)
const merchantList = ref<Merchant[]>([])

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  shop_no: '',
  shop_name: '',
  merchant_id: undefined as number | undefined,
  building: '',
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit' | 'view'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref<ShopForm>({
  shop_no: '',
  shop_name: '',
  merchant_id: undefined as unknown as number,
  building: '',
  floor: '',
  room_no: '',
  contact_name: '',
  contact_phone: '',
  status: 1,
  remark: ''
})

const formRules = reactive<FormRules>({
  shop_name: [
    { required: true, message: '请输入店铺名称', trigger: 'blur' },
    { max: 100, message: '店铺名称不能超过100个字符', trigger: 'blur' }
  ],
  merchant_id: [
    { required: true, message: '请选择所属商户', trigger: 'change' }
  ],
  contact_phone: [
    { validator: validatePhone, trigger: 'blur' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await shopApi.getShopList(queryParams)
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
  queryParams.shop_no = ''
  queryParams.shop_name = ''
  queryParams.merchant_id = undefined
  queryParams.building = ''
  queryParams.status = undefined
  handleQuery()
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    shop_no: '',
    shop_name: '',
    merchant_id: undefined as unknown as number,
    building: '',
    floor: '',
    room_no: '',
    contact_name: '',
    contact_phone: '',
    status: 1,
    remark: ''
  }
  dialogVisible.value = true
}

// 查看
function handleView(row: Shop) {
  dialogType.value = 'view'
  formData.value = {
    id: row.id,
    shop_no: row.shop_no,
    shop_name: row.shop_name,
    merchant_id: row.merchant_id,
    building: row.building,
    floor: row.floor,
    room_no: row.room_no,
    contact_name: row.contact_name,
    contact_phone: row.contact_phone,
    status: row.status,
    remark: row.remark
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Shop) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    shop_no: row.shop_no,
    shop_name: row.shop_name,
    merchant_id: row.merchant_id,
    building: row.building,
    floor: row.floor,
    room_no: row.room_no,
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
          await shopApi.createShop(formData.value)
          ElMessage.success('新增成功')
        } else {
          await shopApi.updateShop(formData.value.id!, formData.value)
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
async function handleDelete(row: Shop) {
  await ElMessageBox.confirm(`确定要删除店铺"${row.shop_name}"吗？`, '提示', { type: 'warning' })
  await shopApi.deleteShop(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
  fetchMerchantList()
})
</script>
