<script setup>
import { onMounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import GxAdminSidebar from '../components/gx/GxAdminSidebar.vue'
import GxSiteHeader from '../components/gx/GxSiteHeader.vue'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const router = useRouter()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

onMounted(() => session.refreshMe())

function backToCommunity() {
  closeDrawer()
  router.push('/community')
}
</script>

<template>
  <div class="gx-app gx-admin-shell">
    <GxSiteHeader :drawer-open="drawerOpen" @toggle-drawer="toggleDrawer" />
    <div class="gx-drawer-backdrop" :class="{ 'is-open': drawerOpen }" @click="closeDrawer" />
    <GxAdminSidebar :open="drawerOpen" @navigate="closeDrawer" @back="backToCommunity" />
    <main class="gx-main">
      <RouterView />
    </main>
  </div>
</template>
