# 前端开发规范

## 1. API 字段命名规范

### 1.1 后端返回字段格式
后端 API 统一使用 **snake_case** (下划线命名) 格式返回 JSON 数据。

示例：
```json
{
  "id": 1,
  "merchant_no": "M20240101001",
  "merchant_name": "测试商户",
  "merchant_type": 1,
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "created_at": "2024-01-01T00:00:00+08:00"
}
```

### 1.2 前端字段使用规范
前端在以下场景中 **必须** 使用与后端一致的 snake_case 字段名：

1. **表格列绑定** (`el-table-column` 的 `prop` 属性)
```vue
<!-- 正确 -->
<el-table-column prop="merchant_name" label="商户名称" />

<!-- 错误 -->
<el-table-column prop="merchantName" label="商户名称" />
```

2. **模板中访问数据**
```vue
<!-- 正确 -->
<span>{{ row.merchant_name }}</span>
<el-tag v-if="row.merchant_type === 1">企业</el-tag>

<!-- 错误 -->
<span>{{ row.merchantName }}</span>
```

3. **API 请求参数**
```typescript
// 正确
const params = {
  merchant_id: 1,
  perm_code: 'system:user:list'
}

// 错误
const params = {
  merchantId: 1,
  permCode: 'system:user:list'
}
```

## 2. 组件使用规范

### 2.1 el-tag 组件 type 属性
当使用动态 `type` 属性时，必须添加 null 检查，防止 Vue 警告：

```vue
<!-- 正确写法 -->
<el-table-column prop="status" label="状态">
  <template #default="{ row }">
    <template v-if="row.status != null">
      <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
        {{ row.status === 1 ? '正常' : '停用' }}
      </el-tag>
    </template>
    <span v-else>-</span>
  </template>
</el-table-column>

<!-- 错误写法 - 可能导致 type 为 null 的警告 -->
<el-table-column prop="status" label="状态">
  <template #default="{ row }">
    <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
      {{ row.status === 1 ? '正常' : '停用' }}
    </el-tag>
  </template>
</el-table-column>
```

### 2.2 el-tree / el-tree-select 组件
配置 props 时使用 snake_case：

```vue
<!-- 正确 -->
<el-tree-select
  :data="permTreeOptions"
  :props="{ label: 'perm_name', value: 'id', children: 'children' }"
/>

<!-- 错误 -->
<el-tree-select
  :data="permTreeOptions"
  :props="{ label: 'permName', value: 'id', children: 'children' }"
/>
```

## 3. 常用字段对照表

| 中文名称 | 后端字段 (snake_case) |
|---------|---------------------|
| 商户编号 | merchant_no |
| 商户名称 | merchant_name |
| 商户类型 | merchant_type |
| 联系人 | contact_name |
| 联系电话 | contact_phone |
| 联系邮箱 | contact_email |
| 账户编号 | account_no |
| 账户名称 | account_name |
| 权限编码 | perm_code |
| 权限名称 | perm_name |
| 权限类型 | perm_type |
| 角色编码 | role_code |
| 角色名称 | role_name |
| 角色类型 | role_type |
| 数据范围 | data_scope |
| 排序序号 | sort_order |
| 在线状态 | online_status |
| 创建时间 | created_at |
| 更新时间 | updated_at |

## 4. TypeScript 类型定义

虽然后端返回 snake_case，但 TypeScript 类型定义可以使用 camelCase（保持代码一致性），前提是需要在 API 层进行转换，或者直接使用 snake_case 类型定义。

**推荐方案**：类型定义直接使用 snake_case，与后端保持一致：

```typescript
// types/business.d.ts
export interface Merchant {
  id: number
  merchant_no: string
  merchant_name: string
  merchant_type: number
  contact_name: string | null
  contact_phone: string | null
  status: number
  created_at: string
  updated_at: string
}
```

## 5. 日期时间处理

后端返回的日期时间格式为 ISO 8601 格式（带时区），前端使用 `formatDateTime` 工具函数进行格式化：

```vue
<el-table-column prop="created_at" label="创建时间">
  <template #default="{ row }">
    {{ formatDateTime(row.created_at) }}
  </template>
</el-table-column>
```
