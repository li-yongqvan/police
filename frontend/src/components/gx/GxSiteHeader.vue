<script setup>
import { computed, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import GxAvatarInitial from './GxAvatarInitial.vue'
import GxIcon from './GxIcon.vue'
import { GX_HEADER_NAV } from '../../composables/useGxNav'
import { useSearchPanel } from '../../composables/useSearchPanel'
import { formatUserChipLabel } from '../../utils/displayName'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  drawerOpen: { type: Boolean, default: false },
  unreadCount: { type: Number, default: 0 },
})
const emit = defineEmits(['toggle-drawer'])

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const showDropdown = ref(false)
const { open: openSearch } = useSearchPanel()

const displayName = computed(() =>
  session.currentUser ? formatUserChipLabel(session.currentUser) : '同学',
)
const avatarSrc = computed(() => session.currentUser?.avatar || '')
const notificationLabel = computed(() =>
  props.unreadCount > 0 ? `通知，${props.unreadCount > 99 ? '99+' : props.unreadCount} 条未读` : '通知',
)

function isHeaderActive(item) {
  return item.match?.(route) ?? route.path.startsWith(item.to)
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
        aria-controls="gx-community-sidebar"
        :aria-expanded="drawerOpen"
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
        <button
          type="button"
          class="gx-header__search-btn hidden md:flex"
          data-search-trigger
          aria-label="搜索帖子"
          @click="openSearch"
        >
          <GxIcon name="search" :size="18" />
          <span class="gx-header__search-placeholder">搜索</span>
        </button>
        <button
          type="button"
          class="gx-header__icon-btn md:hidden"
          data-search-trigger
          aria-label="搜索帖子"
          @click="openSearch"
        >
          <GxIcon name="search" :size="20" />
        </button>

        <RouterLink
          to="/community/messages"
          class="gx-header__icon-btn gx-header__icon-btn--badge"
          :aria-label="notificationLabel"
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
            <GxAvatarInitial :name="displayName" :src="avatarSrc" :size="32" />
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
  z-index: 70;
  min-width: 148px;
  padding: 4px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  box-shadow: 0 8px 24px rgba(15, 43, 91, 0.18);
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
  cursor: pointer;
  font-size: 14px;
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

.gx-header__search-btn {
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md, 8px);
  background: var(--color-bg, #f5f5f5);
  color: var(--color-muted, #9ca3af);
  cursor: pointer;
  transition: border-color 0.15s ease;
}

.gx-header__search-btn:hover {
  border-color: var(--color-primary, #3d7c73);
}

.gx-header__search-placeholder {
  padding: 1px 6px;
  border: 1px solid var(--color-border, #e5e7eb);
  border-radius: 4px;
  color: var(--color-muted, #9ca3af);
  font-family: monospace;
  font-size: 13px;
}

@media (max-width: 767.98px) {
  .gx-header__icon-btn {
    min-height: var(--mw-tap-min, 44px);
  }
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
