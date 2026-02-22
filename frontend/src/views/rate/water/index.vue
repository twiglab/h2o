<template>
  <div class="page-container">
    <div class="search-bar">
      <el-form :model="queryParams" inline>
        <el-form-item label="费率名称">
          <el-input v-model="queryParams.rateName" placeholder="请输入" clearable style="width: 150px" />
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

    <div class="table-container">
      <div class="table-toolbar">
        <span class="table-title">水费费率</span>
        <el-button type="primary" :icon="Plus" @click="handleAdd">新增</el-button>
      </div>

      <el-table v-loading="loading" :data="tableData" stripe border>
        <el-table-column prop="rateCode" label="费率编码" width="140" />
        <el-table-column prop="rateName" label="费率名称" min-width="150" />
        <el-table-column prop="calcMode" label="计费模式" width="120" align="center">
          <template #default="{ row }">
            <template v-if="row.calcMode != null">
              <el-tag :type="row.calcMode === 1 ? 'info' : 'warning'" size="small">{{ row.calcMode === 1 ? '固定单价' : '阶梯水价' }}</el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="unitPrice" label="单价(元/吨)" width="130" align="right" />
        <el-table-column prop="effectiveDate" label="生效日期" width="120" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <template v-if="row.status != null">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" text size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50]" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleQuery" @current-change="handleQuery" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)

const queryParams = reactive({ page: 1, pageSize: 10, rateName: '', status: undefined as number | undefined })

async function fetchList() {
  loading.value = true
  try { tableData.value = []; total.value = 0 } finally { loading.value = false }
}

function handleQuery() { queryParams.page = 1; fetchList() }
function handleReset() { queryParams.rateName = ''; queryParams.status = undefined; handleQuery() }
function handleAdd() { ElMessage.info('功能开发中') }
function handleEdit(_row: any) { ElMessage.info('功能开发中') }
function handleDelete(_row: any) { ElMessage.info('功能开发中') }

onMounted(() => fetchList())
</script>
