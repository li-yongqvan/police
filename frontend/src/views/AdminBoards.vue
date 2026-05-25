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
  <div class="page-stack">
    <section class="panel form-panel">
      <p class="eyebrow">新建板块</p>
      <h2>扩展社区结构</h2>
      <div class="form-grid">
        <label>
          <span>名称</span>
          <input v-model="form.name" type="text" />
        </label>
        <label>
          <span>Slug</span>
          <input v-model="form.slug" type="text" placeholder="例如 ai-lab" />
        </label>
        <label class="full-span">
          <span>描述</span>
          <input v-model="form.description" type="text" />
        </label>
      </div>
      <button class="primary-button" @click="create">创建板块</button>
    </section>

    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">板块列表</p>
          <h3>{{ boards.length }} 个板块</h3>
        </div>
      </div>
      <div class="user-list">
        <article v-for="board in boards" :key="board.id" class="user-row">
          <div v-if="editingId === board.id" class="full-span form-grid">
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
            <label>
              <span>启用</span>
              <input v-model="editForm.enabled" type="checkbox" />
            </label>
            <label class="full-span">
              <span>描述</span>
              <input v-model="editForm.description" type="text" />
            </label>
            <div class="audit-actions">
              <button class="primary-button" @click="saveEdit">保存</button>
              <button class="secondary-button" @click="cancelEdit">取消</button>
            </div>
          </div>
          <template v-else>
            <div>
              <strong>{{ board.name }}</strong>
              <p>{{ board.slug }} · {{ board.postCount }} 帖 · {{ board.enabled ? '启用' : '停用' }}</p>
            </div>
            <div class="audit-actions">
              <button class="secondary-button" @click="startEdit(board)">编辑</button>
              <button class="danger-button" @click="remove(board)">删除</button>
            </div>
          </template>
        </article>
      </div>
    </section>
  </div>
</template>
