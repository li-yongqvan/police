<script setup>
import { onMounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { LogOut, Menu } from 'lucide-vue-next'
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
        title="打开菜单"
        @click="toggleDrawer"
      >
        <Menu :size="22" stroke-width="2" aria-hidden="true" />
      </button>
      <span class="gx-admin-topbar__brand" title="AI 智联 · 管理端">
        <span class="gx-admin-topbar__brand-full">AI 智联 · 管理端</span>
        <span class="gx-admin-topbar__brand-short" aria-hidden="true">管理端</span>
      </span>
      <button type="button" class="gx-admin-topbar__logout" aria-label="退出登录" title="退出登录" @click="logout">
        <LogOut :size="18" stroke-width="2" aria-hidden="true" />
        <span class="gx-admin-topbar__logout-text">退出</span>
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
  min-height: 56px;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
  background: var(--color-surface, #fff);
  position: sticky;
  top: 0;
  z-index: 50;
  box-sizing: border-box;
}
.gx-admin-topbar__menu {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 44px;
  width: 44px;
  min-width: 44px;
  height: 44px;
  min-height: 44px;
  padding: 0;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--color-primary, #1a2332);
  cursor: pointer;
}
.gx-admin-topbar__brand {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 700;
  font-size: 16px;
  color: var(--color-primary, #1a2332);
}
.gx-admin-topbar__brand-short {
  display: none;
}
.gx-admin-topbar__logout {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-width: 44px;
  min-height: 44px;
  padding: 0 12px;
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

@media (max-height: 480px) and (orientation: landscape) {
  .gx-admin-topbar {
    gap: 8px;
    height: 52px;
    min-height: 52px;
    padding: 0 12px;
  }
  .gx-admin-topbar__brand-full,
  .gx-admin-topbar__logout-text {
    display: none;
  }
  .gx-admin-topbar__brand-short {
    display: inline;
  }
  .gx-admin-topbar__logout {
    width: 44px;
    padding: 0;
  }
}
</style>
