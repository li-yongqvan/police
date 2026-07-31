<script setup>
import { computed } from "vue"
import { RouterLink, useRoute } from "vue-router"
import GxIcon from "./GxIcon.vue"
import { useSessionStore } from "../../stores/session"

defineProps({
  open: { type: Boolean, default: false },
})
const emit = defineEmits(["navigate", "back"])

const route = useRoute()
const session = useSessionStore()

const isPlatform = computed(() => session.currentUser?.role === "platform_admin")

const DISCOURSE_ADMIN_URL = "http://122.51.233.225:8080/admin"

const adminNav = computed(() => {
  const items = [
    { to: "/admin", label: "管理首页", icon: "home" },
    { to: "/admin/users", label: "用户管理", icon: "user" },
  ]
  if (isPlatform.value) {
    items.push(
      { to: "/admin/invites", label: "邀请码", icon: "flag" },
    )
  }
  return items
})

function onNav() {
  emit("navigate")
}

function back() {
  emit("navigate")
  emit("back")
}
</script>

<template>
  <aside class="gx-admin-sidebar" :class="{ 'is-open': open }" aria-label="管理端导航">
    <p class="gx-sidebar-nav__section">管理端</p>
    <nav class="gx-sidebar-nav">
      <RouterLink
        v-for="item in adminNav"
        :key="item.to"
        :to="item.to"
        class="gx-sidebar-nav__link"
        :class="{ 'is-active': route.path === item.to || (item.to !== '/admin' && route.path.startsWith(item.to)) }"
        @click="onNav"
      >
        <GxIcon :name="item.icon" :size="20" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>
    <div class="gx-sidebar-nav__section">外部管理</div>
    <nav class="gx-sidebar-nav">
      <a
        :href="DISCOURSE_ADMIN_URL"
        class="gx-sidebar-nav__link"
        target="_blank"
        rel="noopener noreferrer"
        @click="onNav"
      >
        <GxIcon name="shield" :size="20" />
        <span>前往 Discourse 管理</span>
      </a>
    </nav>
  </aside>
</template>
