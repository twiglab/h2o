<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="角色名称">
          <el-input v-model="queryParams.roleName" placeholder="请输入" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="角色编码">
          <el-input v-model="queryParams.roleCode" placeholder="请输入" clearable style="width: 160px" />
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
        <span class="table-title">角色列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'system:role:create'" @click="handleAdd">
            新增
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="role_code" label="角色编码" width="150" />
        <el-table-column prop="role_name" label="角色名称" min-width="120" />
        <el-table-column prop="role_type" label="角色类型" width="100" align="center">
          <template #default="{ row }">
            <template v-if="row.role_type != null">
              <el-tag :type="row.role_type === 1 ? 'warning' : 'info'" size="small">
                {{ row.role_type === 1 ? '内置' : '自定义' }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="data_scope" label="数据权限" width="140">
          <template #default="{ row }">
            {{ getDataScopeLabel(row.data_scope) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              :disabled="row.role_type === 1"
              @change="handleStatusChange(row)"
              v-permission="'system:role:update'"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" v-permission="'system:role:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" text size="small" v-permission="'system:role:permission'" @click="handlePermission(row)">
              权限
            </el-button>
            <el-button
              type="danger"
              text
              size="small"
              :disabled="row.roleType === 1"
              v-permission="'system:role:delete'"
              @click="handleDelete(row)"
            >
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
      :title="dialogType === 'add' ? '新增角色' : '编辑角色'"
      width="500px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="角色编码" prop="roleCode">
          <el-input v-model="formData.roleCode" placeholder="请输入" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="formData.roleName" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="数据权限" prop="dataScope">
          <el-select v-model="formData.dataScope" placeholder="请选择" style="width: 100%">
            <el-option v-for="item in dataScopeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="formData.sortOrder" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限分配弹窗 -->
    <el-dialog v-model="permDialogVisible" title="分配权限" width="500px" destroy-on-close>
      <div class="perm-tree-header">
        <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAllChange">
          全选/全不选
        </el-checkbox>
        <el-checkbox v-model="checkStrictly">父子不关联</el-checkbox>
      </div>
      <el-scrollbar height="400px">
        <el-tree
          ref="permTreeRef"
          :data="permissionTree"
          :props="{ label: 'perm_name', children: 'children' }"
          node-key="id"
          show-checkbox
          :check-strictly="checkStrictly"
          default-expand-all
        />
      </el-scrollbar>
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="permSubmitLoading" @click="handlePermSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/format'
import * as roleApi from '@/api/modules/role'
import * as permissionApi from '@/api/modules/permission'
import type { Role, Permission } from '@/types/system'

const loading = ref(false)
const tableData = ref<Role[]>([])
const total = ref(0)
const permissionTree = ref<Permission[]>([])

// 数据权限选项
const dataScopeOptions = [
  { label: '全部数据', value: 1 },
  { label: '本部门及下级部门', value: 2 },
  { label: '仅本部门', value: 3 },
  { label: '仅本人', value: 4 }
]

function getDataScopeLabel(value: number): string {
  return dataScopeOptions.find(item => item.value === value)?.label || '-'
}

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  roleName: '',
  roleCode: '',
  status: undefined as number | undefined
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = ref({
  id: undefined as number | undefined,
  roleCode: '',
  roleName: '',
  dataScope: 4,
  sortOrder: 0,
  status: 1,
  description: ''
})

const formRules = reactive<FormRules>({
  roleCode: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { max: 50, message: '角色编码不能超过50个字符', trigger: 'blur' }
  ],
  roleName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { max: 50, message: '角色名称不能超过50个字符', trigger: 'blur' }
  ],
  dataScope: [
    { required: true, message: '请选择数据权限', trigger: 'change' }
  ]
})

// 权限弹窗
const permDialogVisible = ref(false)
const permSubmitLoading = ref(false)
const permTreeRef = ref<InstanceType<typeof ElTree>>()
const currentRoleId = ref<number>()
const checkAll = ref(false)
const isIndeterminate = ref(false)
const checkStrictly = ref(false)

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await roleApi.getRoleList(queryParams)
    tableData.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

// 获取权限树
async function fetchPermissionTree() {
  const res = await permissionApi.getPermissionTree()
  permissionTree.value = res
}

// 搜索
function handleQuery() {
  queryParams.page = 1
  fetchList()
}

// 重置
function handleReset() {
  queryParams.roleName = ''
  queryParams.roleCode = ''
  queryParams.status = undefined
  handleQuery()
}

// 新增
function handleAdd() {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    roleCode: '',
    roleName: '',
    dataScope: 4,
    sortOrder: 0,
    status: 1,
    description: ''
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Role) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    roleCode: row.roleCode,
    roleName: row.roleName,
    dataScope: row.dataScope,
    sortOrder: row.sortOrder || 0,
    status: row.status,
    description: row.description || ''
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
          await roleApi.createRole(formData.value)
          ElMessage.success('新增成功')
        } else {
          await roleApi.updateRole(formData.value.id!, formData.value)
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
async function handleDelete(row: Role) {
  await ElMessageBox.confirm(`确定要删除角色"${row.roleName}"吗？`, '提示', { type: 'warning' })
  await roleApi.deleteRole(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

// 状态变更
async function handleStatusChange(row: Role) {
  try {
    await roleApi.updateRole(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch {
    row.status = row.status === 1 ? 0 : 1
  }
}

// 分配权限
async function handlePermission(row: Role) {
  currentRoleId.value = row.id
  checkAll.value = false
  isIndeterminate.value = false
  permDialogVisible.value = true

  await nextTick()
  // 获取角色已有权限
  const permIds = await roleApi.getRolePermissions(row.id)
  permTreeRef.value?.setCheckedKeys(permIds)
}

// 全选/全不选
function handleCheckAllChange(val: boolean) {
  if (val) {
    const allKeys = getAllTreeKeys(permissionTree.value)
    permTreeRef.value?.setCheckedKeys(allKeys)
  } else {
    permTreeRef.value?.setCheckedKeys([])
  }
  isIndeterminate.value = false
}

// 获取所有树节点的key
function getAllTreeKeys(tree: Permission[]): number[] {
  let keys: number[] = []
  tree.forEach(node => {
    keys.push(node.id)
    if (node.children && node.children.length > 0) {
      keys = keys.concat(getAllTreeKeys(node.children))
    }
  })
  return keys
}

// 提交权限
async function handlePermSubmit() {
  if (!currentRoleId.value) return
  permSubmitLoading.value = true
  try {
    const checkedKeys = permTreeRef.value?.getCheckedKeys() as number[]
    const halfCheckedKeys = permTreeRef.value?.getHalfCheckedKeys() as number[]
    const allKeys = [...checkedKeys, ...halfCheckedKeys]
    await roleApi.updateRolePermissions(currentRoleId.value, allKeys)
    ElMessage.success('权限分配成功')
    permDialogVisible.value = false
  } finally {
    permSubmitLoading.value = false
  }
}

onMounted(() => {
  fetchList()
  fetchPermissionTree()
})
</script>

<style lang="scss" scoped>
.perm-tree-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: $spacing-md;
  padding-bottom: $spacing-sm;
  border-bottom: 1px solid $border-color-light;
}
</style>
