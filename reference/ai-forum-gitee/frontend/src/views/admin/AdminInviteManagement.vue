<template>
  <div class="admin-invite" v-loading="loading">
    <h2>邀请码管理</h2>
    <div style="margin-bottom: 16px; display: flex; gap: 12px;">
      <el-input v-model="searchQuery" placeholder="搜索邀请码" style="width: 200px" />
      <el-select v-model="statusFilter" style="width: 120px">
        <el-option label="全部" value="all" />
        <el-option label="未使用" value="unused" />
        <el-option label="已使用" value="used" />
        <el-option label="已作废" value="voided" />
      </el-select>
      <el-button @click="fetchCodes">搜索</el-button>
    </div>
    <el-table :data="codes" stripe>
      <el-table-column prop="code" label="邀请码" width="200" />
      <el-table-column prop="created_by_name" label="创建者" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusColor(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="used_by_name" label="使用者" width="120" />
      <el-table-column prop="used_at" label="使用时间" width="180">
        <template #default="{ row }">
          {{ row.used_at ? new Date(row.used_at).toLocaleString() : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'unused'"
            type="danger"
            size="small"
            @click="handleVoid(row)"
          >作废</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="page"
      :page-size="limit"
      :total="total"
      layout="prev, pager, next"
      @current-change="onPageChange"
      style="margin-top: 16px"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listInviteCodes, voidInviteCode } from '@/api/admin'

const codes = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(20)
const total = ref(0)
const searchQuery = ref('')
const statusFilter = ref('all')

async function fetchCodes() {
  loading.value = true
  try {
    const res = await listInviteCodes(page.value, limit.value)
    codes.value = res.data.codes || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error('获取邀请码失败')
  } finally {
    loading.value = false
  }
}

function onPageChange(p) { page.value = p; fetchCodes() }

function statusColor(s) {
  return s === 'unused' ? 'success' : s === 'used' ? 'info' : 'danger'
}

async function handleVoid(row) {
  try {
    await ElMessageBox.confirm(`确认作废邀请码 ${row.code}？`, '确认')
    await voidInviteCode(row.code)
    ElMessage.success('邀请码已作废')
    fetchCodes()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

onMounted(fetchCodes)
</script>

<style scoped>
.admin-invite h2 { margin-bottom: 16px; }
</style>
