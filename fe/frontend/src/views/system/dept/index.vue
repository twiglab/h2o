<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="部门名称">
          <el-input v-model="queryParams.dept_name" placeholder="请输入" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 120px">
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
        <span class="table-title">部门列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'system:dept:create'" @click="handleAdd()">
            新增
          </el-button>
          <el-button :icon="Sort" @click="toggleExpandAll">
            {{ expandAll ? '收起' : '展开' }}全部
          </el-button>
        </div>
      </div>

      <el-table
        v-if="refreshTable"
        v-loading="loading"
        :data="tableData"
        row-key="id"
        :default-expand-all="expandAll"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        border
      >
        <el-table-column prop="dept_name" label="部门名称" min-width="200" />
        <el-table-column prop="dept_code" label="部门编码" width="150" />
        <el-table-column prop="leader_name" label="负责人" width="120" />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
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
            <el-button type="primary" text size="small" v-permission="'system:dept:create'" @click="handleAdd(row)">
              新增
            </el-button>
            <el-button type="primary" text size="small" v-permission="'system:dept:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" text size="small" v-permission="'system:dept:delete'" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增部门' : '编辑部门'"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="上级部门" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="deptTreeOptions"
            :props="{ label: 'dept_name', value: 'id' }"
            placeholder="请选择上级部门"
            check-strictly
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="部门名称" prop="dept_name">
          <el-input v-model="formData.dept_name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="部门编码" prop="dept_code">
          <el-input v-model="formData.dept_code" placeholder="请输入部门编码" />
        </el-form-item>
        <el-form-item label="负责人" prop="leader_id">
          <el-select v-model="formData.leader_id" placeholder="请选择负责人" clearable style="width: 100%">
            <el-option v-for="user in userList" :key="user.id" :label="user.real_name || user.username" :value="user.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Sort } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import * as deptApi from '@/api/modules/dept'
import * as userApi from '@/api/modules/user'
import type { Dept, User } from '@/types/system'

const loading = ref(false)
const tableData = ref<Dept[]>([])
const expandAll = ref(true)
const refreshTable = ref(true)
const userList = ref<User[]>([])
const deptTreeOptions = ref<Dept[]>([])

// 查询参数
const queryParams = reactive({
  dept_name: '',
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref({
  id: undefined as number | undefined,
  parent_id: 0,
  dept_code: '',
  dept_name: '',
  leader_id: undefined as number | undefined,
  sort_order: 0,
  status: 1
})

const formRules = reactive<FormRules>({
  dept_name: [
    { required: true, message: '请输入部门名称', trigger: 'blur' },
    { max: 100, message: '部门名称不能超过100个字符', trigger: 'blur' }
  ],
  dept_code: [
    { required: true, message: '请输入部门编码', trigger: 'blur' },
    { max: 32, message: '部门编码不能超过32个字符', trigger: 'blur' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await deptApi.getDeptTree(queryParams)
    tableData.value = res
    deptTreeOptions.value = [{ id: 0, dept_name: '顶级部门', children: res } as Dept]
  } finally {
    loading.value = false
  }
}

// 获取用户列表（用于选择负责人）
async function fetchUserList() {
  const res = await userApi.getUserList({ page: 1, pageSize: 1000, status: 1 })
  userList.value = res.list
}

// 搜索
function handleQuery() {
  fetchList()
}

// 重置
function handleReset() {
  queryParams.dept_name = ''
  queryParams.status = undefined
  handleQuery()
}

// 展开/收起
function toggleExpandAll() {
  refreshTable.value = false
  expandAll.value = !expandAll.value
  nextTick(() => {
    refreshTable.value = true
  })
}

// 新增
function handleAdd(row?: Dept) {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    parent_id: row?.id || 0,
    dept_code: '',
    dept_name: '',
    leader_id: undefined,
    sort_order: 0,
    status: 1
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Dept) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    parent_id: row.parent_id || 0,
    dept_code: row.dept_code,
    dept_name: row.dept_name,
    leader_id: row.leader_id,
    sort_order: row.sort_order || 0,
    status: row.status
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
          await deptApi.createDept(formData.value)
          ElMessage.success('新增成功')
        } else {
          await deptApi.updateDept(formData.value.id!, formData.value)
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
async function handleDelete(row: Dept) {
  if (row.children && row.children.length > 0) {
    ElMessage.warning('该部门下有子部门，不能删除')
    return
  }
  await ElMessageBox.confirm(`确定要删除部门"${row.dept_name}"吗？`, '提示', { type: 'warning' })
  await deptApi.deleteDept(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
  fetchUserList()
})
</script>
