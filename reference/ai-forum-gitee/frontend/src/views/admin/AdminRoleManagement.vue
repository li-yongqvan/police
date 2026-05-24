<template>
  <div class="admin-roles" v-loading="loading">
    <h2>角色权限管理</h2>
    <h3>可用角色</h3>
    <el-table :data="roles" stripe style="margin-bottom: 24px">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="角色名" width="150" />
      <el-table-column prop="description" label="描述" width="150" />
      <el-table-column label="权限">
        <template #default="{ row }">
          <el-tag v-for="p in row.permissions" :key="p" size="small" style="margin: 2px">{{ p }}</el-tag>
        </template>
      </el-table-column>
    </el-table>

    <h3>分配角色</h3>
    <div style="margin-bottom: 16px; display: flex; gap: 12px;">
      <el-input v-model="assignUserId" placeholder="用户ID" style="width: 120px" />
      <el-select v-model="assignRoleId" placeholder="选择角色">
        <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
      </el-select>
      <el-button type="primary" @click="handleAssign">分配</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listRoles, assignRole } from '@/api/admin'

const roles = ref([])
const loading = ref(false)
const assignUserId = ref('')
const assignRoleId = ref(null)

async function fetchRoles() {
  loading.value = true
  try {
    const res = await listRoles()
    roles.value = res.data.roles || []
  } catch (e) {
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

async function handleAssign() {
  if (!assignUserId.value || !assignRoleId.value) {
    ElMessage.warning('请输入用户ID并选择角色')
    return
  }
  try {
    await assignRole(parseInt(assignUserId.value), assignRoleId.value)
    ElMessage.success('角色已分配')
    assignUserId.value = ''
    assignRoleId.value = null
  } catch (e) {
    ElMessage.error('分配失败')
  }
}

onMounted(fetchRoles)
</script>

<style scoped>
.admin-roles h2 { margin-bottom: 16px; }
.admin-roles h3 { margin: 16px 0 8px 0; }
</style>
