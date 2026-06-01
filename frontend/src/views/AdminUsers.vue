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
  <div class="gx-page gx-admin-page">
    <header class="gx-page-head">
      <p class="gx-eyebrow">用户管理</p>
      <h1>封禁、解封与等级</h1>
      <p class="gx-muted">第 {{ page }} / {{ totalPages() }} 页</p>
    </header>

    <section class="gx-card">
      <div class="gx-section-head">
        <span class="gx-muted">共 {{ total }} 人</span>
        <div class="gx-admin-actions">
          <button type="button" class="gx-btn gx-btn--secondary" :disabled="page <= 1" @click="changePage(page - 1)">
            上一页
          </button>
          <button
            type="button"
            class="gx-btn gx-btn--secondary"
            :disabled="page >= totalPages()"
            @click="changePage(page + 1)"
          >
            下一页
          </button>
        </div>
      </div>

      <div class="gx-admin-list">
        <article v-for="user in users" :key="user.id" class="gx-admin-row">
          <div>
            <strong>{{ user.name }}</strong>
            <p class="gx-muted">{{ user.username }} · {{ user.role }} · Lv.{{ user.level }}</p>
          </div>
          <div class="gx-admin-actions">
            <span :class="['gx-badge', user.status === 'active' ? 'gx-badge--ok' : 'gx-badge--bad']">
              {{ user.status }}
            </span>
            <label class="gx-muted">
              <input v-model="levelDraft[user.id]" type="number" min="0" max="5" :placeholder="String(user.level)" />
            </label>
            <button type="button" class="gx-btn gx-btn--ghost" @click="saveLevel(user)">改等级</button>
            <button type="button" class="gx-btn gx-btn--secondary" @click="toggleStatus(user)">
              {{ user.status === 'active' ? '封禁' : '解封' }}
            </button>
            <button type="button" class="gx-btn gx-btn--ghost" @click="showLogs(user)">日志</button>
          </div>
        </article>
      </div>
    </section>

    <section v-if="logsUserId" class="gx-card">
      <h3 class="gx-panel__title">操作日志 · 用户 {{ logsUserId }}</h3>
      <div class="gx-admin-list">
        <article v-for="(log, i) in logs" :key="i" class="gx-comment">
          <strong>{{ log.action || log.operation }}</strong>
          <p class="gx-muted">{{ log.created_at || log.timestamp }} — {{ log.detail || log.message || '' }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
