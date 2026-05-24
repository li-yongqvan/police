<template>
  <div class="admin-sensitive-words" v-loading="loading">
    <h2>敏感词管理</h2>
    <div style="margin-bottom: 16px; display: flex; gap: 12px;">
      <el-input v-model="newWord" placeholder="敏感词" style="width: 150px" />
      <el-select v-model="newCategory" style="width: 120px">
        <el-option label="general" value="general" />
        <el-option label="spam" value="spam" />
        <el-option label="abuse" value="abuse" />
        <el-option label="ads" value="ads" />
      </el-select>
      <el-button type="primary" @click="handleAdd">添加敏感词</el-button>
    </div>
    <el-table :data="words" stripe>
      <el-table-column prop="word" label="敏感词" width="150" />
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column prop="created_at" label="添加时间" width="180">
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
import { listSensitiveWords, addSensitiveWord, deleteSensitiveWord } from '@/api/admin'

const words = ref([])
const loading = ref(false)
const newWord = ref('')
const newCategory = ref('general')

async function fetchWords() {
  loading.value = true
  try {
    // Use direct axios since the sensitive words API is on /api/v1/admin/sensitive-words
    const axios = (await import('axios')).default
    const res = await axios.get('/api/v1/admin/sensitive-words', {
      headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` }
    })
    words.value = res.data
  } catch (e) {
    ElMessage.error('获取敏感词失败')
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  if (!newWord.value) {
    ElMessage.warning('请输入敏感词')
    return
  }
  try {
    await addSensitiveWord(newWord.value, newCategory.value)
    ElMessage.success('添加成功')
    newWord.value = ''
    fetchWords()
  } catch (e) {
    ElMessage.error('添加失败')
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm('确认删除此敏感词？', '确认')
    await deleteSensitiveWord(row.id)
    ElMessage.success('已删除')
    fetchWords()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(fetchWords)
</script>

<style scoped>
.admin-sensitive-words h2 { margin-bottom: 16px; }
</style>
