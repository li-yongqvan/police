<template>
  <div class="admin-user-management" v-loading="loading">
    <h2>用户管理</h2>
    <div style="margin-bottom: 16px; display: flex; gap: 12px;">
      <el-input v-model="searchQuery" placeholder="搜索用户名" style="width: 200px" />
      <el-select v-model="statusFilter" style="width: 120px">
        <el-option label="全部" value="all" />
        <el-option label="正常" value="active" />
        <el-option label="封禁" value="banned" />
      </el-select>
      <el-button @click="fetchUsers">搜索</el-button>
    </div>
    <el-table :data="users" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="nickname" label="昵称" width="120" />
      <el-table-column prop="level" label="等级" width="80" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '正常' : '封禁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="300">
        <template #default="{ row }">
          <el-button
            :type="row.status === 'banned' ? 'success' : 'danger'"
            size="small"
            @click="handleBan(row)"
          >
            {{ row.status === 'banned' ? '解封' : '封禁' }}
          </el-button>
          <el-button size="small" @click="openLevelDialog(row)">等级调整</el-button>
          <el-button size="small" @click="openLogsDrawer(row)">查看日志</el-button>
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

    <el-dialog v-model="banVisible" :title="banData?.status === 'banned' ? '确认解封' : '确认封禁'" width="400px">
      <el-input v-model="banReason" type="textarea" placeholder="请输入原因" />
      <template #footer>
        <el-button @click="banVisible = false">取消</el-button>
        <el-button :type="banData?.status === 'banned' ? 'success' : 'danger'" @click="handleBanSubmit">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="levelVisible" title="调整用户等级" width="400px">
      <p>当前等级: {{ levelData?.level }}</p>
      <el-select v-model="newLevel">
        <el-option v-for="n in 6" :key="n-1" :label="`Level ${n-1}`" :value="n-1" />
      </el-select>
      <template #footer>
        <el-button @click="levelVisible = false">取消</el-button>
        <el-button type="primary" @click="handleLevelSubmit">确认调整</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="logsVisible" :title="logsUser?.username + ' 的操作日志'" size="40%">
      <el-table :data="userLogs">
        <el-table-column prop="action" label="操作" width="120" />
        <el-table-column prop="target_type" label="目标类型" width="100" />
        <el-table-column prop="created_at" label="时间">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listUsers, banUser, updateUserLevel, getUserLogs } from '@/api/admin'

const users = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(20)
const total = ref(0)
const searchQuery = ref('')
const statusFilter = ref('all')

const banVisible = ref(false)
const banData = ref(null)
const banReason = ref('')

const levelVisible = ref(false)
const levelData = ref(null)
const newLevel = ref(0)

const logsVisible = ref(false)
const logsUser = ref(null)
const userLogs = ref([])

async function fetchUsers() {
  loading.value = true
  try {
    const res = await listUsers(page.value, limit.value, statusFilter.value)
    users.value = res.data.users || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

function onPageChange(p) { page.value = p; fetchUsers() }

function handleBan(row) {
  banData.value = row
  banReason.value = ''
  banVisible.value = true
}

async function handleBanSubmit() {
  try {
    await banUser(banData.value.id, banReason.value)
    ElMessage.success(banData.value.status === 'banned' ? '已解封' : '已封禁')
    banVisible.value = false
    fetchUsers()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

function openLevelDialog(row) {
  levelData.value = row
  newLevel.value = row.level
  levelVisible.value = true
}

async function handleLevelSubmit() {
  try {
    await updateUserLevel(levelData.value.id, newLevel.value)
    ElMessage.success('等级已调整')
    levelVisible.value = false
    fetchUsers()
  } catch (e) {
    ElMessage.error('调整失败')
  }
}

async function openLogsDrawer(row) {
  logsUser.value = row
  userLogs.value = []
  logsVisible.value = true
  try {
    const res = await getUserLogs(row.id, 1, 50)
    userLogs.value = res.data.logs || []
  } catch (e) {
    ElMessage.error('获取日志失败')
  }
}

onMounted(fetchUsers)
</script>

<style scoped>
.admin-user-management h2 { margin-bottom: 16px; }
</style>
