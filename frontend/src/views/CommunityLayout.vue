<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import GxAdminFab from '../components/gx/GxAdminFab.vue'
import GxMobileTabBar from '../components/gx/GxMobileTabBar.vue'
import GxSiteFooter from '../components/gx/GxSiteFooter.vue'
import GxSiteHeader from '../components/gx/GxSiteHeader.vue'
import GxSiteSidebar from '../components/gx/GxSiteSidebar.vue'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'
import { forumApi } from '../api'

const route = useRoute()
const session = useSessionStore()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

const isFeedShell = computed(() =>
  ['community-home', 'board', 'rank', 'campus-circle', 'my-posts', 'my-favorites', 'my-history'].includes(
    route.name,
  ),
)
const boards = ref([])
const unreadCount = ref(0)

onMounted(async () => {
  try {
    boards.value = await forumApi.getBoards()
  } catch {
    boards.value = []
  }
  if (session.currentUser) {
    try {
      const { items } = await forumApi.listNotifications(1, 50)
      unreadCount.value = items.filter((i) => !i.is_read).length
    } catch {
      unreadCount.value = 0
    }
  }
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
