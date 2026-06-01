<script setup>
import { computed, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import GxAvatarInitial from './GxAvatarInitial.vue'
import GxIcon from './GxIcon.vue'
import { GX_HEADER_NAV } from '../../composables/useGxNav'
import { useSessionStore } from '../../stores/session'

defineProps({
  drawerOpen: { type: Boolean, default: false },
  unreadCount: { type: Number, default: 0 },
})
const emit = defineEmits(['toggle-drawer'])

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const showDropdown = ref(false)
const searchQ = ref(route.query.q || '')

const displayName = computed(() => session.currentUser?.name || session.currentUser?.username || '同学')

function isHeaderActive(item) {
  return item.match?.(route) ?? route.path.startsWith(item.to)
}

function submitSearch() {
  const q = searchQ.value.trim()
  router.push(q ? { path: '/community', query: { q } } : { path: '/community' })
}

function toggleDropdown() {
  showDropdown.value = !showDropdown.value
}

function closeDropdown() {
  showDropdown.value = false
}

function goProfile() {
  closeDropdown()
  router.push('/community/profile')
}

function doLogout() {
  closeDropdown()
  session.logout()
  router.push('/')
}
</script>

<template>
  <header class="gx-header gx-header--mockup">
    <div class="gx-header__inner">
      <button
        type="button"
        class="gx-header__menu gx-tap"
        aria-label="打开菜单"
        @click="emit('toggle-drawer')"
      >
        <span class="gx-header__menu-bar" />
        <span class="gx-header__menu-bar" />
        <span class="gx-header__menu-bar" />
      </button>

      <RouterLink to="/community" class="gx-header__brand" @click="closeDropdown">
        <span class="gx-header__logo" aria-hidden="true">
          <GxIcon name="shield" :size="22" />
        </span>
        <span class="gx-header__site">
          <strong>警院论坛</strong>
          <em>忠诚 · 勤学 · 严谨 · 创新</em>
        </span>
      </RouterLink>

      <nav class="gx-header__tabs hidden lg:flex" aria-label="站点导航">
        <RouterLink
          v-for="item in GX_HEADER_NAV"
          :key="item.id"
          :to="item.to"
          class="gx-header__tab"
          :class="{ 'is-active': isHeaderActive(item) }"
          @click="closeDropdown"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="gx-header__actions">
        <form class="gx-header__search hidden md:flex" @submit.prevent="submitSearch">
          <GxIcon name="search" :size="18" />
          <input v-model="searchQ" type="search" placeholder="搜索帖子、板块或用户" aria-label="搜索" />
        </form>

        <RouterLink
          to="/community/messages"
          class="gx-header__icon-btn gx-header__icon-btn--badge"
          aria-label="通知"
          @click="closeDropdown"
        >
          <GxIcon name="bell" :size="20" />
          <span v-if="unreadCount > 0" class="gx-header__badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
        </RouterLink>

        <RouterLink
          to="/community/messages"
          class="gx-header__icon-btn hidden sm:grid"
          aria-label="消息"
          @click="closeDropdown"
        >
          <GxIcon name="message" :size="20" />
        </RouterLink>

        <div class="gx-avatar-dropdown">
          <button
            type="button"
            class="gx-header__user-chip"
            aria-label="用户菜单"
            :aria-expanded="showDropdown"
            @click="toggleDropdown"
          >
            <GxAvatarInitial :name="displayName" :size="32" />
            <span class="gx-header__user-name hidden md:inline">{{ displayName }}</span>
            <span class="gx-header__caret" aria-hidden="true">▾</span>
          </button>
          <Transition name="gx-dropdown">
            <div v-if="showDropdown" class="gx-avatar-dropdown__menu">
              <button type="button" class="gx-avatar-dropdown__item" @click="goProfile">
                <GxIcon name="user" :size="16" />
                <span>个人中心</span>
              </button>
              <button
                type="button"
                class="gx-avatar-dropdown__item gx-avatar-dropdown__item--danger"
                @click="doLogout"
              >
                <GxIcon name="logout" :size="16" />
                <span>退出登录</span>
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </div>
    <div v-if="showDropdown" class="gx-header__backdrop" aria-hidden="true" @click="closeDropdown" />
  </header>
</template>

<style scoped>
.gx-avatar-dropdown {
  position: relative;
}

.gx-avatar-dropdown__menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 148px;
  padding: 4px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: 0 8px 24px rgba(15, 43, 91, 0.18);
  z-index: 70;
}

.gx-avatar-dropdown__item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-primary);
  font-size: 14px;
  cursor: pointer;
}

.gx-avatar-dropdown__item:hover {
  background: var(--color-bg);
}

.gx-avatar-dropdown__item--danger {
  color: var(--color-danger);
}

.gx-header__backdrop {
  position: fixed;
  inset: 0;
  z-index: 55;
  background: transparent;
}

.gx-dropdown-enter-active,
.gx-dropdown-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.gx-dropdown-enter-from,
.gx-dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
