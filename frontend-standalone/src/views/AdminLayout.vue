<script setup>
import { onMounted } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const router = useRouter()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

onMounted(async () => {
  await session.refreshMe()
})

function backToCommunity() {
  closeDrawer()
  router.push('/community')
}
</script>

<template>
  <div class="layout-app admin-shell">
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
        <strong>论坛中台</strong>
        <span>{{ session.currentUser?.name }}</span>
      </div>
    </header>

    <div
      class="layout-backdrop"
      :class="{ 'is-visible': drawerOpen }"
      aria-hidden="true"
      @click="closeDrawer"
    />

    <aside class="layout-drawer sidebar admin-sidebar" :class="{ 'is-open': drawerOpen }">
      <div class="brand-card dark">
        <p class="eyebrow">Admin Console</p>
        <h1>论坛中台</h1>
        <p>控制、监管、统计三件事压缩在一套最小管理界面里。</p>
      </div>

      <section class="panel user-panel">
        <div class="user-header">
          <div class="avatar">{{ session.currentUser?.avatar }}</div>
          <div>
            <h2>{{ session.currentUser?.name }}</h2>
            <p>{{ session.currentUser?.department }}</p>
          </div>
        </div>
      </section>

      <nav class="panel nav-panel">
        <RouterLink to="/admin" class="nav-link">数据概览</RouterLink>
        <RouterLink to="/admin/audit" class="nav-link">内容审核</RouterLink>
        <RouterLink to="/admin/config" class="nav-link">系统配置</RouterLink>
        <RouterLink to="/admin/users" class="nav-link">用户管理</RouterLink>
      </nav>

      <button class="secondary-button full-width" @click="backToCommunity">返回前台社区</button>
    </aside>

    <main class="layout-main main-area">
      <RouterView />
    </main>
  </div>
</template>
