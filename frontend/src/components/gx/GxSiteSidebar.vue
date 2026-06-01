<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import GxIcon from './GxIcon.vue'
import GxSidebarUserCard from './GxSidebarUserCard.vue'
import {
  GX_NAV_ITEMS,
  GX_SIDEBAR_PERSONAL,
  isNavActive,
  resolveBoardByKey,
} from '../../composables/useGxNav'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  open: { type: Boolean, default: false },
  boards: { type: Array, default: () => [] },
})
const emit = defineEmits(['navigate'])

const route = useRoute()
const router = useRouter()
const session = useSessionStore()

const boardLinks = computed(() =>
  GX_NAV_ITEMS.map((nav) => {
    const board = resolveBoardByKey(props.boards, nav.key)
    return {
      key: nav.key,
      label: board?.name || nav.label,
      to: `/community/boards/${nav.key}`,
      postCount: board?.postCount ?? 0,
    }
  }),
)

function onNav() {
  emit('navigate')
}

function isPersonalActive(item) {
  if (item.name) return route.name === item.name
  return route.path === item.to
}

function isBoardActive(key) {
  return route.name === 'board' && route.params.slug === key
}

function logout() {
  emit('navigate')
  session.logout()
  router.push('/')
}
</script>

<template>
  <aside
    class="gx-community-sidebar gx-sidebar-mockup flex flex-col"
    :class="{ 'is-open': open }"
    aria-label="社区导航"
  >
    <p class="gx-sidebar-nav__section">导航</p>
    <nav class="gx-sidebar-nav">
      <RouterLink
        v-for="item in GX_SIDEBAR_PERSONAL"
        :key="item.id"
        :to="item.to"
        class="gx-sidebar-nav__link"
        :class="{ 'is-active': isPersonalActive(item) }"
        @click="onNav"
      >
        <GxIcon :name="item.icon" :size="20" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <div class="gx-sidebar-nav__sep" />

    <p class="gx-sidebar-nav__section">板块列表</p>
    <nav class="gx-sidebar-boards">
      <RouterLink
        v-for="item in boardLinks"
        :key="item.key"
        :to="item.to"
        class="gx-sidebar-boards__link"
        :class="{ 'is-active': isBoardActive(item.key) }"
        @click="onNav"
      >
        <span>{{ item.label }}</span>
        <GxIcon v-if="isBoardActive(item.key)" name="star" :size="16" class="gx-sidebar-boards__star" />
      </RouterLink>
    </nav>

    <GxSidebarUserCard class="gx-sidebar-user-wrap mt-auto" />

    <nav v-if="session.canAccessAdmin" class="gx-sidebar-nav">
      <RouterLink to="/admin" class="gx-sidebar-nav__link" @click="onNav">
        <GxIcon name="shield" :size="20" />
        <span>管理端</span>
      </RouterLink>
      <button type="button" class="gx-sidebar-nav__link w-full border-0 bg-transparent text-left" @click="logout">
        <GxIcon name="logout" :size="20" />
        <span>退出登录</span>
      </button>
    </nav>
  </aside>
</template>
