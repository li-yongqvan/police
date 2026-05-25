<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { adminApi, forumApi } from '../api'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()
const boards = ref([])
const config = ref(null)

const roleLabel = computed(() => {
  const role = session.currentUser?.role
  if (role === 'platform_admin') return '中台管理员'
  if (role === 'admin') return '协会管理员'
  return '学生用户'
})

const pendingCount = computed(() =>
  session.canAccessAdmin ? '可进入中台' : '前台演示',
)

async function loadMeta() {
  boards.value = await forumApi.getBoards()
  if (session.canAccessAdmin) {
    config.value = await adminApi.getConfig()
  }
}

function logout() {
  closeDrawer()
  session.logout()
  router.push('/')
}

onMounted(() => {
  loadMeta()
  window.addEventListener('forum-config-updated', loadMeta)
})
</script>

<template>
  <div class="layout-app">
    <header class="layout-topbar">
      <button
        type="button"
        class="layout-menu-button"
        aria-label="打开导航菜单"
        @click="toggleDrawer"
      >
        菜单
      </button>
      <div class="layout-topbar-title">
        <strong>AI 智联论坛</strong>
        <span>{{ session.currentUser?.name }} · {{ roleLabel }}</span>
      </div>
    </header>

    <div
      class="layout-backdrop"
      :class="{ 'is-visible': drawerOpen }"
      aria-hidden="true"
      @click="closeDrawer"
    />

    <aside class="layout-drawer sidebar" :class="{ 'is-open': drawerOpen }">
      <div class="brand-card">
        <p class="eyebrow">Campus AI Community</p>
        <h1>AI 智联论坛</h1>
        <p>面向学院内部交流、问答求助、活动运营与最小中台管理。</p>
      </div>

      <section class="panel user-panel">
        <div class="user-header">
          <div class="avatar">{{ session.currentUser?.avatar }}</div>
          <div>
            <h2>{{ session.currentUser?.name }}</h2>
            <p>{{ session.currentUser?.department }}</p>
          </div>
        </div>
        <p class="user-bio">{{ session.currentUser?.bio }}</p>
        <div class="status-row">
          <span class="status-badge">{{ roleLabel }}</span>
          <span class="status-light">{{ pendingCount }}</span>
        </div>
      </section>

      <nav class="panel nav-panel">
        <RouterLink to="/community" class="nav-link">社区总览</RouterLink>
        <RouterLink
          v-for="board in boards"
          :key="board.id"
          :to="`/community/boards/${board.slug}`"
          class="nav-link"
        >
          {{ board.name }}
        </RouterLink>
        <RouterLink to="/community/posts/new" class="nav-link">发布新帖</RouterLink>
        <RouterLink to="/community/profile" class="nav-link">个人主页</RouterLink>
        <RouterLink v-if="session.canAccessAdmin" to="/admin" class="nav-link accent-link">
          进入中台
        </RouterLink>
      </nav>

      <section v-if="config" class="panel info-panel">
        <p class="eyebrow">当前控制状态</p>
        <h3>{{ config.postingEnabled ? '发帖开放中' : '发帖已关闭' }}</h3>
        <p>监管模式：{{ config.moderationMode === 'manual' ? '人工审核' : '自动审核' }}</p>
      </section>

      <button class="ghost-button full-width" @click="logout">退出演示</button>
    </aside>

    <main class="layout-main main-area">
      <RouterView />
    </main>
  </div>
</template>
