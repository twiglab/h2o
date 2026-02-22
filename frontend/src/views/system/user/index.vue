<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="queryParams.realName" placeholder="请输入" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="部门">
          <el-tree-select
            v-model="queryParams.deptId"
            :data="deptTree"
            :props="{ label: 'deptName', value: 'id' }"
            placeholder="请选择"
            clearable
            check-strictly
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择" clearable style="width: 120px">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
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
        <span class="table-title">用户列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'system:user:create'" @click="handleAdd">
            新增
          </el-button>
          <el-button type="danger" :icon="Delete" v-permission="'system:user:delete'" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除
          </el-button>
        </div>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        @selection-change="handleSelectionChange"
        stripe
        border
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="realName" label="姓名" min-width="100" />
        <el-table-column prop="phone" label="手机号" min-width="120" />
        <el-table-column prop="deptName" label="部门" min-width="120" />
        <el-table-column prop="roleNames" label="角色" min-width="150">
          <template #default="{ row }">
            <el-tag v-for="role in row.roleNames" :key="role" size="small" class="role-tag">
              {{ role }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
              v-permission="'system:user:update'"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" v-permission="'system:user:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" text size="small" v-permission="'system:user:reset-password'" @click="handleResetPassword(row)">
              重置密码
            </el-button>
            <el-button type="danger" text size="small" v-permission="'system:user:delete'" @click="handleDelete(row)">
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
        />
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="formData.username" placeholder="请输入" :disabled="dialogType === 'edit'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="realName">
              <el-input v-model="formData.realName" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="dialogType === 'add'">
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="formData.password" type="password" placeholder="请输入" show-password />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="formData.confirmPassword" type="password" placeholder="请确认密码" show-password />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="formData.email" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门" prop="deptId">
              <el-tree-select
                v-model="formData.deptId"
                :data="deptTree"
                :props="{ label: 'deptName', value: 'id' }"
                placeholder="请选择"
                check-strictly
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="roleIds">
              <el-select v-model="formData.roleIds" multiple placeholder="请选择" style="width: 100%">
                <el-option v-for="role in roleList" :key="role.id" :label="role.roleName" :value="role.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">正常</el-radio>
                <el-radio :value="0">禁用</el-radio>
              </el-radio-group>
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
import { Search, Refresh, Plus, Delete } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import { validatePhone, validateEmail } from '@/utils/validate'
import * as userApi from '@/api/modules/user'
import * as deptApi from '@/api/modules/dept'
import * as roleApi from '@/api/modules/role'
import type { User, Dept, Role } from '@/types/system'

const loading = ref(false)
const tableData = ref<User[]>([])
const total = ref(0)
const selectedIds = ref<number[]>([])
const deptTree = ref<Dept[]>([])
const roleList = ref<Role[]>([])

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  username: '',
  realName: '',
  deptId: undefined as number | undefined,
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref({
  id: undefined as number | undefined,
  username: '',
  realName: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: '',
  deptId: undefined as number | undefined,
  roleIds: [] as number[],
  status: 1
})

const formRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度为2-50个字符', trigger: 'blur' }
  ],
  realName: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== formData.value.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  phone: [
    { validator: validatePhone, trigger: 'blur' }
  ],
  email: [
    { validator: validateEmail, trigger: 'blur' }
  ],
  deptId: [
    { required: true, message: '请选择部门', trigger: 'change' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await userApi.getUserList(queryParams)
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

// 获取部门树
async function fetchDeptTree() {
  const res = await deptApi.getDeptTree()
  deptTree.value = res
}

// 获取角色列表
async function fetchRoleList() {
  const res = await roleApi.getRoleList({ page: 1, pageSize: 1000, status: 1 })
  roleList.value = res.list
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.username = ''
  queryParams.realName = ''
  queryParams.deptId = undefined
  queryParams.status = undefined
  handleQuery()
}

// 选择变化
function handleSelectionChange(selection: User[]) {
  selectedIds.value = selection.map(item => item.id)
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    username: '',
    realName: '',
    password: '',
    confirmPassword: '',
    phone: '',
    email: '',
    deptId: undefined,
    roleIds: [],
    status: 1
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: User) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    username: row.username,
    realName: row.realName || '',
    password: '',
    confirmPassword: '',
    phone: row.phone || '',
    email: row.email || '',
    deptId: row.deptId,
    roleIds: row.roleIds || [],
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
          await userApi.createUser(formData.value)
          ElMessage.success('新增成功')
        } else {
          await userApi.updateUser(formData.value.id!, formData.value)
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
async function handleDelete(row: User) {
  await ElMessageBox.confirm(`确定要删除用户"${row.username}"吗？`, '提示', { type: 'warning' })
  await userApi.deleteUser(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

// 批量删除
async function handleBatchDelete() {
  await ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}个用户吗？`, '提示', { type: 'warning' })
  await userApi.batchDeleteUsers(selectedIds.value)
  ElMessage.success('删除成功')
  fetchList()
}

// 状态变更
async function handleStatusChange(row: User) {
  try {
    await userApi.updateUser(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch {
    row.status = row.status === 1 ? 0 : 1
  }
}

// 重置密码
async function handleResetPassword(row: User) {
  await ElMessageBox.confirm(`确定要重置用户"${row.username}"的密码吗？`, '提示', { type: 'warning' })
  await userApi.resetPassword(row.id)
  ElMessage.success('密码已重置为默认密码')
}

// 监听分页变化
watch(
  () => [queryParams.page, queryParams.pageSize],
  ([newPage, newPageSize], [oldPage, oldPageSize]) => {
    // 如果是 pageSize 变化，重置到第一页
    if (newPageSize !== oldPageSize) {
      queryParams.page = 1
    }
    fetchList()
  }
)

onMounted(() => {
  fetchList()
  fetchDeptTree()
  fetchRoleList()
})
</script>

<style lang="scss" scoped>
.role-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}
</style>
