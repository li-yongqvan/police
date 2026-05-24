<template>
  <div class="admin-board" v-loading="loading">
    <h2>板块管理</h2>
    <el-button type="primary" @click="openCreateDialog">新增板块</el-button>
    <el-table :data="boards" stripe style="margin-top: 16px">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="slug" label="标识" width="150" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="enabled" label="状态" width="100">
        <template #default="{ row }">
          <el-switch v-model="row.enabled" @change="toggleEnabled(row)" />
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="formVisible" :title="editingBoard ? '编辑板块' : '新增板块'" width="500px">
      <el-form :model="boardForm" label-width="80px">
        <el-form-item label="名称"><el-input v-model="boardForm.name" /></el-form-item>
        <el-form-item label="标识"><el-input v-model="boardForm.slug" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="boardForm.description" type="textarea" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="boardForm.sort_order" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listBoards, createBoard, updateBoard, deleteBoard } from '@/api/admin'

const boards = ref([])
const loading = ref(false)
const formVisible = ref(false)
const editingBoard = ref(null)
const boardForm = ref({ name: '', slug: '', description: '', sort_order: 0 })

async function fetchBoards() {
  loading.value = true
  try {
    const res = await listBoards()
    boards.value = res.data.boards || []
  } catch (e) {
    ElMessage.error('获取板块列表失败')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  editingBoard.value = null
  boardForm.value = { name: '', slug: '', description: '', sort_order: 0 }
  formVisible.value = true
}

function openEditDialog(row) {
  editingBoard.value = row
  boardForm.value = { name: row.name, slug: row.slug, description: row.description, sort_order: row.sort_order }
  formVisible.value = true
}

async function handleSave() {
  try {
    if (editingBoard.value) {
      await updateBoard(editingBoard.value.id, boardForm.value)
      ElMessage.success('板块已更新')
    } else {
      await createBoard(boardForm.value)
      ElMessage.success('板块已创建')
    }
    formVisible.value = false
    fetchBoards()
  } catch (e) {
    ElMessage.error('保存失败')
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm('确认删除此板块？', '确认')
    await deleteBoard(row.id)
    ElMessage.success('板块已删除')
    fetchBoards()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

async function toggleEnabled(row) {
  try {
    await updateBoard(row.id, { enabled: row.enabled })
    ElMessage.success('状态已更新')
  } catch (e) {
    ElMessage.error('更新失败')
  }
}

onMounted(fetchBoards)
</script>

<style scoped>
.admin-board h2 { margin-bottom: 16px; }
</style>
