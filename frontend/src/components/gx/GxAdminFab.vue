<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useSessionStore } from '../../stores/session'

const session = useSessionStore()
const route = useRoute()

const visible = computed(
  () =>
    session.currentUser?.role === 'platform_admin' &&
    route.path.startsWith('/community') &&
    !route.path.startsWith('/admin'),
)
</script>

<template>
  <RouterLink v-if="visible" to="/admin" class="gx-admin-fab gx-tap" title="返回管理端">
    管理端
  </RouterLink>
</template>
