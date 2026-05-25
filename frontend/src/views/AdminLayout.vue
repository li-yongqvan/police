<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const router = useRouter()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

const roleLabel = computed(() =>
  session.currentUser?.role === 'platform_admin' ? '中台管理员' : '协会管理员',
)

const isPlatform = computed(() => session.currentUser?.role === 'platform_admin')

const adminNav = computed(() => {
  const items = [
    { to: '/admin', label: '数据概览' },
    { to: '/admin/stats', label: '趋势统计' },
    { to: '/admin/audit', label: '内容审核' },
    { to: '/admin/posts', label: '帖子管理' },
    { to: '/admin/users', label: '用户管理' },
    { to: '/admin/boards', label: '板块管理' },
    { to: '/admin/config', label: isPlatform.value ? '系统配置' : '运营配置' },
  ]
  if (isPlatform.value) {
    items.push(
      { to: '/admin/invites', label: '邀请码' },
      { to: '/admin/sensitive', label: '敏感词' },
      { to: '/admin/roles', label: '角色权限' },
    )
  }
  return items
})

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
        <span>{{ session.currentUser?.name }} · {{ roleLabel }}</span>
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
            <p>{{ roleLabel }}</p>
          </div>
        </div>
      </section>

      <nav class="panel nav-panel">
        <RouterLink
          v-for="item in adminNav"
          :key="item.to"
          :to="item.to"
          class="nav-link"
          @click="closeDrawer"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <button class="secondary-button full-width" @click="backToCommunity">返回前台社区</button>
    </aside>

    <main class="layout-main main-area">
      <RouterView />
    </main>
  </div>
</template>
