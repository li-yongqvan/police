<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const boards = ref([])
const form = ref({ name: '', slug: '', description: '', sortOrder: 0 })
const editingId = ref('')
const editForm = ref({ name: '', slug: '', description: '', sortOrder: 0, enabled: true })

async function load() {
  boards.value = await adminApi.listBoards()
}

async function create() {
  await adminApi.createBoard(form.value)
  form.value = { name: '', slug: '', description: '', sortOrder: 0 }
  session.setFlash('板块已创建。', 'success')
  await load()
}

function startEdit(board) {
  editingId.value = board.id
  editForm.value = {
    name: board.name,
    slug: board.slug,
    description: board.description || '',
    sortOrder: board.sortOrder ?? 0,
    enabled: board.enabled !== false,
  }
}

function cancelEdit() {
  editingId.value = ''
}

async function saveEdit() {
  await adminApi.updateBoard(editingId.value, editForm.value)
  session.setFlash('板块已更新。', 'success')
  editingId.value = ''
  await load()
}

async function remove(board) {
  if (!window.confirm(`确定删除板块「${board.name}」？`)) return
  await adminApi.deleteBoard(board.id)
  session.setFlash('板块已删除。', 'success')
  await load()
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page">
    <header class="gx-page-head">
      <p class="gx-eyebrow">板块管理</p>
      <h1>扩展社区结构</h1>
    </header>

    <section class="gx-card gx-form">
      <h3 class="gx-panel__title">新建板块</h3>
      <label>
        <span>名称</span>
        <input v-model="form.name" type="text" />
      </label>
      <label>
        <span>Slug</span>
        <input v-model="form.slug" type="text" placeholder="例如 ai-lab" />
      </label>
      <label>
        <span>描述</span>
        <input v-model="form.description" type="text" />
      </label>
      <button type="button" class="gx-btn gx-btn--primary" @click="create">创建板块</button>
    </section>

    <section class="gx-card">
      <h3 class="gx-panel__title">{{ boards.length }} 个板块</h3>
      <div class="gx-admin-list">
        <article v-for="board in boards" :key="board.id" class="gx-admin-row">
          <div v-if="editingId === board.id" class="gx-form" style="width: 100%">
            <label>
              <span>名称</span>
              <input v-model="editForm.name" type="text" />
            </label>
            <label>
              <span>Slug</span>
              <input v-model="editForm.slug" type="text" />
            </label>
            <label>
              <span>排序</span>
              <input v-model.number="editForm.sortOrder" type="number" />
            </label>
            <label class="gx-admin-row">
              <span>启用</span>
              <input v-model="editForm.enabled" type="checkbox" />
            </label>
            <label>
              <span>描述</span>
              <input v-model="editForm.description" type="text" />
            </label>
            <div class="gx-admin-actions">
              <button type="button" class="gx-btn gx-btn--primary" @click="saveEdit">保存</button>
              <button type="button" class="gx-btn gx-btn--secondary" @click="cancelEdit">取消</button>
            </div>
          </div>
          <template v-else>
            <div>
              <strong>{{ board.name }}</strong>
              <p class="gx-muted">{{ board.slug }} · {{ board.postCount }} 帖 · {{ board.enabled ? '启用' : '停用' }}</p>
            </div>
            <div class="gx-admin-actions">
              <button type="button" class="gx-btn gx-btn--secondary" @click="startEdit(board)">编辑</button>
              <button type="button" class="gx-btn gx-btn--danger" @click="remove(board)">删除</button>
            </div>
          </template>
        </article>
      </div>
    </section>
  </div>
</template>
