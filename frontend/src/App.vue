<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import GxSearchPanel from './components/gx/GxSearchPanel.vue'
import { useSessionStore } from './stores/session'

const session = useSessionStore()

onMounted(() => {
  if (!session.token) return
  session.refreshMe().catch(() => {})
})
</script>

<template>
  <div class="app-root">
    <transition name="fade">
      <div v-if="session.flashMessage" :class="['toast', session.flashType]">
        {{ session.flashMessage }}
      </div>
    </transition>
    <RouterView />
    <GxSearchPanel />
  </div>
</template>
