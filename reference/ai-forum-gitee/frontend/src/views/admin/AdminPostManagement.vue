<template>
  <div class="admin-post-management" v-loading="loading">
    <h2>帖子管理</h2>
    <el-table :data="posts" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="标题" width="250" />
      <el-table-column prop="author_name" label="作者" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'warning'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_featured" label="精华" width="80">
        <template #default="{ row }">
          <el-switch v-model="row.is_featured" @change="toggleFeatured(row)" />
        </template>
      </el-table-column>
      <el-table-column prop="is_pinned" label="置顶" width="80">
        <template #default="{ row }">
          <el-switch v-model="row.is_pinned" @change="togglePinned(row)" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="80">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPosts, deletePost, setPostFeatured, setPostPinned } from '@/api/admin'

const posts = ref([])
const loading = ref(false)

async function fetchPosts() {
  loading.value = true
  try {
    const { data } = await listPosts(1, 100)
    posts.value = data.posts || []
  } catch (e) {
    ElMessage.error('获取帖子列表失败')
  } finally {
    loading.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm('确认删除此帖子？', '删除确认')
    await deletePost(row.id)
    ElMessage.success('帖子已删除')
    fetchPosts()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

async function toggleFeatured(row) {
  try {
    await setPostFeatured(row.id, row.is_featured)
    ElMessage.success(row.is_featured ? '已设为精华' : '已取消精华')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function togglePinned(row) {
  try {
    await setPostPinned(row.id, row.is_pinned)
    ElMessage.success(row.is_pinned ? '已设为置顶' : '已取消置顶')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

onMounted(fetchPosts)
</script>

<style scoped>
.admin-post-management h2 { margin-bottom: 16px; }
</style>
