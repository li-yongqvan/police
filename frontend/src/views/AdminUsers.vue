<script setup>
import { computed, onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi, userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '用户管理' },
]

const users = ref([])
const page = ref(1)
const total = ref(0)
const limit = 20
const logsUserId = ref('')
const logs = ref([])
const levelDraft = ref({})
const activeCount = computed(() => users.value.filter((user) => user.status === 'active').length)
const bannedCount = computed(() => users.value.filter((user) => user.status && user.status !== 'active').length)

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
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="用户管理" title="封禁、解封与等级" :description="`第 ${page} / ${totalPages()} 页`" />

    <section class="gx-admin-summary gx-admin-summary--four">
      <article class="gx-admin-summary__item">
        <span>用户总数</span>
        <strong>{{ total }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>当前页用户</span>
        <strong>{{ users.length }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>本页正常</span>
        <strong>{{ activeCount }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>本页异常</span>
        <strong>{{ bannedCount }}</strong>
      </article>
    </section>

    <section class="gx-card">
      <div class="gx-section-head">
        <strong>用户列表</strong>
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

      <div v-if="!users.length" class="gx-empty">当前没有用户数据。</div>
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
      <div class="gx-section-head">
        <strong>操作日志</strong>
        <span class="gx-muted">用户 {{ logsUserId }}</span>
      </div>
      <div v-if="!logs.length" class="gx-empty">当前用户暂无操作日志。</div>
      <div class="gx-admin-list">
        <article v-for="(log, i) in logs" :key="i" class="gx-admin-row gx-admin-row--stack">
          <strong>{{ log.action || log.operation }}</strong>
          <p class="gx-muted">{{ log.created_at || log.timestamp }} — {{ log.detail || log.message || '' }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
