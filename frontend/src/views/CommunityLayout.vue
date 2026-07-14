<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import GxAdminFab from '../components/gx/GxAdminFab.vue'
import GxMobileTabBar from '../components/gx/GxMobileTabBar.vue'
import GxSiteFooter from '../components/gx/GxSiteFooter.vue'
import GxSiteHeader from '../components/gx/GxSiteHeader.vue'
import GxSiteSidebar from '../components/gx/GxSiteSidebar.vue'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useUnreadNotifications } from '../composables/useUnreadNotifications'
import { useSessionStore } from '../stores/session'
import { forumApi } from '../api'

const route = useRoute()
const session = useSessionStore()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()
const { unreadCount, refreshUnreadCount } = useUnreadNotifications()

const isFeedShell = computed(() =>
  ['community-home', 'board', 'rank', 'campus-circle', 'my-posts', 'my-favorites', 'my-history'].includes(
    route.name,
  ),
)
const boards = ref([])

async function loadBoards() {
  try {
    boards.value = await forumApi.getBoards()
  } catch {
    boards.value = []
  }
}

function refreshWhenVisible() {
  if (!document.hidden && session.token) {
    refreshUnreadCount()
  }
}

watch(
  () => session.token,
  () => refreshUnreadCount(),
  { immediate: true },
)

watch(
  () => route.fullPath,
  () => refreshUnreadCount(),
)

onMounted(async () => {
  await loadBoards()
  window.addEventListener('focus', refreshUnreadCount)
  document.addEventListener('visibilitychange', refreshWhenVisible)
})

onUnmounted(() => {
  window.removeEventListener('focus', refreshUnreadCount)
  document.removeEventListener('visibilitychange', refreshWhenVisible)
})
</script>

<template>
  <div class="gx-app gx-community-shell gx-app--with-tabbar">
    <GxSiteHeader
      :drawer-open="drawerOpen"
      :unread-count="unreadCount"
      @toggle-drawer="toggleDrawer"
    />
    <div class="gx-drawer-backdrop" :class="{ 'is-open': drawerOpen }" @click="closeDrawer" />
    <GxSiteSidebar :open="drawerOpen" :boards="boards" @navigate="closeDrawer" />

    <main class="gx-main" :class="{ 'gx-main--feed': isFeedShell }">
      <div class="gx-page-frame" :class="{ 'gx-page-frame--feed': isFeedShell }">
        <RouterView />
      </div>
      <GxSiteFooter v-if="!isFeedShell" />
    </main>
    <GxMobileTabBar />
    <GxAdminFab />
  </div>
</template>
