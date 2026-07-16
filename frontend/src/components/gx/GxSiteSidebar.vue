<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import GxIcon from './GxIcon.vue'
import { GX_NAV_ITEMS } from '../../composables/useGxNav'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  open: { type: Boolean, default: false },
  boards: { type: Array, default: () => [] },
})
const emit = defineEmits(['navigate'])

const route = useRoute()
const router = useRouter()
const session = useSessionStore()

const HOME_LINK = { id: 'home', name: 'community-home', label: '首页', icon: 'home', to: '/community' }

const boardLinks = computed(() => {
  const enabledBoards = props.boards
    .filter((board) => board?.enabled !== false && board?.slug)
    .slice()
    .sort((a, b) => {
      const order = (a.sortOrder ?? 0) - (b.sortOrder ?? 0)
      if (order !== 0) return order
      return String(a.name || '').localeCompare(String(b.name || ''), 'zh-Hans-CN')
    })

  if (enabledBoards.length) {
    return enabledBoards.map((board) => ({
      key: board.slug,
      label: board.name,
      to: `/community/boards/${board.slug}`,
      postCount: board.postCount ?? 0,
      icon: iconForBoard(board),
      description: board.description || '',
    }))
  }

  return GX_NAV_ITEMS.map((nav) => ({
    key: nav.key,
    label: nav.label,
    to: `/community/boards/${nav.key}`,
    postCount: 0,
    icon: nav.icon,
    description: '',
  }))
})

function iconForBoard(board) {
  const name = board?.name || ''
  const navMatch = GX_NAV_ITEMS.find(
    (nav) => nav.key === board?.slug || nav.keywords?.some((keyword) => name.includes(keyword)),
  )
  if (navMatch?.icon) return navMatch.icon
  if (/公告|通知|活动/.test(name)) return 'bell'
  if (/实训|警务|治理|安全/.test(name)) return 'shield'
  if (/社团|风采|校园|圈/.test(name)) return 'flag'
  if (/排行|榜单|数据/.test(name)) return 'bar-chart'
  if (/问答|交流|讨论|消息/.test(name)) return 'message'
  if (/用户|同学|成员/.test(name)) return 'users'
  return 'book'
}

function onNav() {
  emit('navigate')
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
    id="gx-community-sidebar"
    class="gx-community-sidebar gx-sidebar-mockup flex flex-col"
    :class="{ 'is-open': open }"
    aria-label="社区导航"
  >
    <p class="gx-sidebar-nav__section">导航</p>
    <nav class="gx-sidebar-nav">
      <RouterLink
        :to="HOME_LINK.to"
        class="gx-sidebar-nav__link"
        :class="{ 'is-active': route.name === HOME_LINK.name }"
        @click="onNav"
      >
        <GxIcon :name="HOME_LINK.icon" :size="20" />
        <span>{{ HOME_LINK.label }}</span>
      </RouterLink>
    </nav>

    <div class="gx-sidebar-nav__sep" />

    <div class="gx-sidebar-board-head">
      <p class="gx-sidebar-nav__section">板块列表</p>
      <span>{{ boardLinks.length }} 个</span>
    </div>
    <nav class="gx-sidebar-boards" aria-label="板块列表">
      <RouterLink
        v-for="item in boardLinks"
        :key="item.key"
        :to="item.to"
        class="gx-sidebar-boards__link"
        :class="{ 'is-active': isBoardActive(item.key) }"
        :title="item.description || item.label"
        @click="onNav"
      >
        <span class="gx-sidebar-boards__mark">
          <GxIcon :name="item.icon" :size="18" />
        </span>
        <span class="gx-sidebar-boards__body">
          <span class="gx-sidebar-boards__name">{{ item.label }}</span>
          <span v-if="item.description" class="gx-sidebar-boards__desc">{{ item.description }}</span>
        </span>
        <span class="gx-sidebar-boards__count">{{ item.postCount }}</span>
      </RouterLink>
    </nav>

    <nav v-if="session.canAccessAdmin" class="gx-sidebar-nav gx-sidebar-nav--admin mt-auto">
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
