<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="权限名称">
          <el-input v-model="queryParams.permName" placeholder="请输入" clearable style="width: 200px" />
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
        <span class="table-title">权限列表</span>
        <div>
          <el-button type="primary" :icon="Plus" v-permission="'system:permission:create'" @click="handleAdd()">
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
        <el-table-column prop="perm_name" label="权限名称" min-width="200" />
        <el-table-column prop="perm_code" label="权限编码" min-width="180" />
        <el-table-column prop="perm_type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <template v-if="row.perm_type != null">
              <el-tag :type="getPermTypeTag(row.perm_type)" size="small">
                {{ getPermTypeLabel(row.perm_type) }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="icon" label="图标" width="60" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.icon"><component :is="row.icon" /></el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" min-width="150" show-overflow-tooltip />
        <el-table-column prop="component" label="组件路径" min-width="180" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="70" align="center" />
        <el-table-column prop="visible" label="可见" width="70" align="center">
          <template #default="{ row }">
            <template v-if="row.visible != null">
              <el-tag :type="row.visible === 1 ? 'success' : 'info'" size="small">
                {{ row.visible === 1 ? '是' : '否' }}
              </el-tag>
            </template>
            <span v-else>-</span>
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" v-permission="'system:permission:create'" v-if="row.perm_type !== 3" @click="handleAdd(row)">
              新增
            </el-button>
            <el-button type="primary" text size="small" v-permission="'system:permission:update'" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" text size="small" v-permission="'system:permission:delete'" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增权限' : '编辑权限'"
      width="650px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="上级权限" prop="parent_id">
              <el-tree-select
                v-model="formData.parent_id"
                :data="permTreeOptions"
                :props="{ label: 'perm_name', value: 'id' }"
                placeholder="请选择上级权限"
                check-strictly
                clearable
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="权限类型" prop="perm_type">
              <el-radio-group v-model="formData.perm_type">
                <el-radio :value="1">目录</el-radio>
                <el-radio :value="2">菜单</el-radio>
                <el-radio :value="3">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="图标" prop="icon" v-if="formData.perm_type !== 3">
              <el-input v-model="formData.icon" placeholder="图标名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="权限名称" prop="perm_name">
              <el-input v-model="formData.perm_name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限编码" prop="perm_code">
              <el-input v-model="formData.perm_code" placeholder="如: system:user:list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="formData.perm_type !== 3">
          <el-col :span="12">
            <el-form-item label="路由路径" prop="path">
              <el-input v-model="formData.path" placeholder="如: /system/user" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件路径" prop="component" v-if="formData.perm_type === 2">
              <el-input v-model="formData.component" placeholder="如: system/user/index" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="formData.perm_type === 3">
          <el-col :span="12">
            <el-form-item label="API路径" prop="api_path">
              <el-input v-model="formData.api_path" placeholder="如: /api/v1/users" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="HTTP方法" prop="api_method">
              <el-select v-model="formData.api_method" placeholder="请选择" style="width: 100%">
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="排序" prop="sort_order">
              <el-input-number v-model="formData.sort_order" :min="0" :max="999" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8" v-if="formData.perm_type !== 3">
            <el-form-item label="是否可见" prop="visible">
              <el-radio-group v-model="formData.visible">
                <el-radio :value="1">是</el-radio>
                <el-radio :value="0">否</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">正常</el-radio>
                <el-radio :value="0">停用</el-radio>
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
import { ref, reactive, onMounted, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Sort } from '@element-plus/icons-vue'
import * as permissionApi from '@/api/modules/permission'
import type { Permission } from '@/types/system'

const loading = ref(false)
const tableData = ref<Permission[]>([])
const expandAll = ref(true)
const refreshTable = ref(true)
const permTreeOptions = ref<Permission[]>([])

// 权限类型
function getPermTypeLabel(type: number): string {
  const map: Record<number, string> = { 1: '目录', 2: '菜单', 3: '按钮' }
  return map[type] || '-'
}

function getPermTypeTag(type: number | null | undefined): '' | 'success' | 'warning' {
  if (type == null) return ''
  const map: Record<number, '' | 'success' | 'warning'> = { 1: '', 2: 'success', 3: 'warning' }
  return map[type] ?? ''
}

// 查询参数
const queryParams = reactive({
  permName: '',
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
  perm_code: '',
  perm_name: '',
  perm_type: 2 as 1 | 2 | 3,
  path: '',
  component: '',
  icon: '',
  api_path: '',
  api_method: '',
  visible: 1,
  sort_order: 0,
  status: 1
})

const formRules = reactive<FormRules>({
  perm_name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { max: 100, message: '权限名称不能超过100个字符', trigger: 'blur' }
  ],
  perm_code: [
    { required: true, message: '请输入权限编码', trigger: 'blur' },
    { max: 100, message: '权限编码不能超过100个字符', trigger: 'blur' }
  ],
  perm_type: [
    { required: true, message: '请选择权限类型', trigger: 'change' }
  ]
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const res = await permissionApi.getPermissionTree(queryParams)
    tableData.value = res
    permTreeOptions.value = [{ id: 0, perm_name: '顶级权限', children: res } as any]
  } finally {
    loading.value = false
  }
}

// 搜索
function handleQuery() {
  fetchList()
}

// 重置
function handleReset() {
  queryParams.permName = ''
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
function handleAdd(row?: any) {
  dialogType.value = 'add'
  formData.value = {
    id: undefined,
    parent_id: row?.id || 0,
    perm_code: '',
    perm_name: '',
    perm_type: row ? (row.perm_type === 1 ? 2 : 3) : 1,
    path: '',
    component: '',
    icon: '',
    api_path: '',
    api_method: '',
    visible: 1,
    sort_order: 0,
    status: 1
  }
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: any) {
  dialogType.value = 'edit'
  formData.value = {
    id: row.id,
    parent_id: row.parent_id || 0,
    perm_code: row.perm_code,
    perm_name: row.perm_name,
    perm_type: row.perm_type,
    path: row.path || '',
    component: row.component || '',
    icon: row.icon || '',
    api_path: row.api_path || '',
    api_method: row.api_method || '',
    visible: row.visible ?? 1,
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
          await permissionApi.createPermission(formData.value)
          ElMessage.success('新增成功')
        } else {
          await permissionApi.updatePermission(formData.value.id!, formData.value)
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
async function handleDelete(row: any) {
  if (row.children && row.children.length > 0) {
    ElMessage.warning('该权限下有子权限，不能删除')
    return
  }
  await ElMessageBox.confirm(`确定要删除权限"${row.perm_name}"吗？`, '提示', { type: 'warning' })
  await permissionApi.deletePermission(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>
