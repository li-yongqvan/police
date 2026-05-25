<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const users = ref([])
const page = ref(1)
const total = ref(0)
const limit = 20
const logsUserId = ref('')
const logs = ref([])
const levelDraft = ref({})

async function load() {
  const data = await userApi.listUsers(page.value, limit)
  users.value = data.users
  total.value = data.total
}

async function changePage(next) {
  page.value = next
  await load()
}

async function toggleStatus(user) {
  const next = user.status === 'active' ? 'banned' : 'active'
  await adminApi.setUserStatus(user.id, next)
  session.setFlash(next === 'banned' ? '用户已封禁' : '用户已解封', 'success')
  await load()
}

async function saveLevel(user) {
  await adminApi.updateUserLevel(user.id, levelDraft.value[user.id])
  session.setFlash('等级已更新', 'success')
  await load()
}

async function showLogs(user) {
  logsUserId.value = user.id
  logs.value = await adminApi.getUserLogs(user.id)
}

const totalPages = () => Math.max(1, Math.ceil(total.value / limit))

onMounted(load)
</script>

<template>
  <div class="page-stack">
    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">用户管理</p>
          <h3>封禁、解封与等级 · 第 {{ page }} / {{ totalPages() }} 页</h3>
        </div>
        <div class="audit-actions mw-pagination">
          <button class="secondary-button" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <button class="secondary-button" :disabled="page >= totalPages()" @click="changePage(page + 1)">
            下一页
          </button>
        </div>
      </div>

      <div class="user-list">
        <article v-for="user in users" :key="user.id" class="user-row mw-user-row">
          <div>
            <strong>{{ user.name }}</strong>
            <p>{{ user.username }} · {{ user.role }} · Lv.{{ user.level }}</p>
          </div>
          <div class="user-actions stacked-admin-actions">
            <span :class="['badge', user.status === 'active' ? 'good' : 'bad']">{{ user.status }}</span>
            <label class="level-inline">
              <input v-model="levelDraft[user.id]" type="number" min="0" max="5" :placeholder="String(user.level)" />
              <button class="ghost-button" @click="saveLevel(user)">改等级</button>
            </label>
            <button class="secondary-button" @click="toggleStatus(user)">
              {{ user.status === 'active' ? '封禁' : '解封' }}
            </button>
            <button class="ghost-button" @click="showLogs(user)">日志</button>
          </div>
        </article>
      </div>
    </section>

    <section v-if="logsUserId" class="panel content-panel">
      <p class="eyebrow">操作日志 · 用户 {{ logsUserId }}</p>
      <div class="comment-list">
        <article v-for="(log, i) in logs" :key="i" class="comment-card">
          <strong>{{ log.action || log.operation }}</strong>
          <p>{{ log.created_at || log.timestamp }} — {{ log.detail || log.message || '' }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
