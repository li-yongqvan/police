<template>
  <div class="admin-audit" v-loading="loading">
    <h2>内容审核</h2>
    <el-button type="danger" :disabled="selectedPosts.length === 0" @click="showBatchDelete = true">
      批量删除 ({{ selectedPosts.length }})
    </el-button>
    <el-table :data="posts" stripe style="margin-top: 16px">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column prop="author_name" label="作者" width="120" />
      <el-table-column prop="board_name" label="板块" width="150" />
      <el-table-column label="内容预览" width="250">
        <template #default="{ row }">
          {{ (row.content || '').substring(0, 80) }}...
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button type="success" size="small" @click="handleApprove(row)">通过</el-button>
          <el-button type="danger" size="small" @click="openRejectDialog(row)">驳回</el-button>
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

    <el-dialog v-model="rejectVisible" title="驳回原因" width="400px">
      <el-input v-model="rejectReason" type="textarea" placeholder="请输入驳回原因" />
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="danger" @click="handleReject">确认驳回</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showBatchDelete" title="批量删除确认" width="400px">
      <el-input v-model="batchReason" type="textarea" placeholder="请输入删除原因" />
      <template #footer>
        <el-button @click="showBatchDelete = false">取消</el-button>
        <el-button type="danger" @click="handleBatchDelete">确认删除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPendingAudits, approvePost, rejectPost, batchDeletePosts } from '@/api/admin'

const posts = ref([])
const loading = ref(false)
const page = ref(1)
const limit = ref(20)
const total = ref(0)
const rejectVisible = ref(false)
const rejectReason = ref('')
const rejectPostData = ref(null)
const selectedPosts = ref([])
const showBatchDelete = ref(false)
const batchReason = ref('')

async function fetchPosts() {
  loading.value = true
  try {
    const res = await getPendingAudits(page.value, limit.value)
    posts.value = res.data.posts || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error('获取待审核列表失败')
  } finally {
    loading.value = false
  }
}

function onPageChange(p) {
  page.value = p
  fetchPosts()
}

async function handleApprove(row) {
  try {
    await approvePost(row.id)
    ElMessage.success('已审核通过')
    fetchPosts()
  } catch (e) {
    ElMessage.error('审核失败')
  }
}

function openRejectDialog(row) {
  rejectPostData.value = row
  rejectReason.value = ''
  rejectVisible.value = true
}

async function handleReject() {
  try {
    await rejectPost(rejectPostData.value.id, rejectReason.value)
    ElMessage.success('已驳回')
    rejectVisible.value = false
    fetchPosts()
  } catch (e) {
    ElMessage.error('驳回失败')
  }
}

async function handleBatchDelete() {
  try {
    const ids = posts.value.filter(p => selectedPosts.value.includes(p.id)).map(p => p.id)
    await batchDeletePosts(ids, batchReason.value)
    ElMessage.success('批量删除成功')
    showBatchDelete.value = false
    fetchPosts()
  } catch (e) {
    ElMessage.error('批量删除失败')
  }
}

onMounted(fetchPosts)
// Auto-refresh every 30 seconds
setInterval(fetchPosts, 30000)
</script>

<style scoped>
.admin-audit h2 { margin-bottom: 16px; }
</style>
