<script setup>
import { onMounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import GxAdminSidebar from '../components/gx/GxAdminSidebar.vue'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const router = useRouter()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

onMounted(() => session.refreshMe())

function logout() {
  closeDrawer()
  session.logout()
  router.push('/')
}
</script>

<template>
  <div class="gx-app gx-admin-shell">
    <header class="gx-admin-topbar">
      <button
        type="button"
        class="gx-admin-topbar__menu"
        aria-label="打开菜单"
        @click="toggleDrawer"
      >
        <span class="gx-admin-topbar__menu-bar" />
        <span class="gx-admin-topbar__menu-bar" />
        <span class="gx-admin-topbar__menu-bar" />
      </button>
      <span class="gx-admin-topbar__brand">AI 智联 · 管理端</span>
      <button type="button" class="gx-admin-topbar__logout" @click="logout">
        退出
      </button>
    </header>
    <div class="gx-drawer-backdrop" :class="{ 'is-open': drawerOpen }" @click="closeDrawer" />
    <GxAdminSidebar :open="drawerOpen" @navigate="closeDrawer" @back="logout" />
    <main class="gx-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.gx-admin-topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
  height: 56px;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
  background: var(--color-surface, #fff);
  position: sticky;
  top: 0;
  z-index: 50;
}
.gx-admin-topbar__menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
}
.gx-admin-topbar__menu-bar {
  width: 20px;
  height: 2px;
  background: var(--color-primary, #1a2332);
  border-radius: 1px;
}
.gx-admin-topbar__brand {
  font-weight: 700;
  font-size: 16px;
  color: var(--color-primary, #1a2332);
}
.gx-admin-topbar__logout {
  margin-left: auto;
  padding: 6px 14px;
  border: 1px solid var(--color-border, #e5e7eb);
  border-radius: 6px;
  background: transparent;
  color: var(--color-muted, #6b7280);
  cursor: pointer;
  font-size: 13px;
}
.gx-admin-topbar__logout:hover {
  color: var(--color-danger, #ef4444);
  border-color: var(--color-danger, #ef4444);
}
</style>