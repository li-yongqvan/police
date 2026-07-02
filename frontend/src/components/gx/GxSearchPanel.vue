<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { onKeyStroke } from '@vueuse/core'
import { useSearchPanel } from '../../composables/useSearchPanel'
import { useSessionStore } from '../../stores/session'
import GxIcon from './GxIcon.vue'

const router = useRouter()
const session = useSessionStore()
const { isOpen, close } = useSearchPanel()

const keyword = ref('')
const inputRef = ref(null)

const role = computed(() => session.currentUser?.role)
const roleLabel = computed(() => {
  if (role.value === 'platform_admin') return '平台管理员'
  if (role.value === 'admin') return '管理员'
  return '学生'
})

const studentCommands = [
  { id: 'community-home', label: '社区首页', icon: 'home', to: '/community' },
  { id: 'board', label: '浏览板块', icon: 'book', to: '/community/boards/study' },
  { id: 'new-post', label: '发新帖', icon: 'edit', to: '/community/posts/new' },
  { id: 'my-posts', label: '我的帖子', icon: 'file', to: '/community/my/posts' },
  { id: 'my-favorites', label: '我的收藏', icon: 'star', to: '/community/my/favorites' },
  { id: 'my-history', label: '浏览历史', icon: 'clock', to: '/community/my/history' },
  { id: 'rank', label: '排行榜', icon: 'award', to: '/community/rank' },
  { id: 'campus-circle', label: '校园圈', icon: 'users', to: '/community/circle' },
  { id: 'messages', label: '消息中心', icon: 'message', to: '/community/messages' },
  { id: 'profile', label: '个人资料', icon: 'user', to: '/community/profile' },
  { id: 'about', label: '关于本站', icon: 'info', to: '/community/about' },
]

const adminCommands = [
  { id: 'admin-overview', label: '数据概览', icon: 'home', to: '/admin' },
  { id: 'admin-stats', label: '趋势统计', icon: 'bar-chart', to: '/admin/stats' },
  { id: 'admin-audit', label: '内容审核', icon: 'shield', to: '/admin/audit' },
  { id: 'admin-reports', label: '举报处理', icon: 'flag', to: '/admin/reports' },
  { id: 'admin-posts', label: '帖子管理', icon: 'edit', to: '/admin/posts' },
  { id: 'admin-users', label: '用户管理', icon: 'user', to: '/admin/users' },
  { id: 'admin-boards', label: '板块管理', icon: 'book', to: '/admin/boards' },
  { id: 'admin-config', label: '运营配置', icon: 'settings', to: '/admin/config' },
  { id: 'community-back', label: '返回社区', icon: 'arrow-left', to: '/community' },
]

const platformCommands = [
  { id: 'admin-invites', label: '邀请码管理', icon: 'key', to: '/admin/invites' },
  { id: 'admin-sensitive', label: '敏感词管理', icon: 'shield', to: '/admin/sensitive' },
  { id: 'admin-roles', label: '角色权限', icon: 'lock', to: '/admin/roles' },
]

const isAdmin = computed(() => role.value === 'admin' || role.value === 'platform_admin')
const isPlatform = computed(() => role.value === 'platform_admin')

const allCommands = computed(() => {
  if (isAdmin.value) {
    const list = [...adminCommands]
    if (isPlatform.value) list.push(...platformCommands)
    return list
  }
  return studentCommands
})

function matchCommand(cmd) {
  if (!keyword.value) return true
  const kw = keyword.value.toLowerCase()
  return cmd.label.toLowerCase().includes(kw) || cmd.id.toLowerCase().includes(kw)
}

const filteredCommands = computed(() => {
  return allCommands.value.filter(matchCommand)
})

const hasKeyword = computed(() => keyword.value.trim().length > 0)

function navigate(cmd) {
  close()
  router.push(cmd.to)
}

function searchPosts() {
  const q = keyword.value.trim()
  if (!q) return
  close()
  router.push({ path: '/community', query: { q } })
}

watch(isOpen, async (val) => {
  if (val) {
    keyword.value = ''
    await nextTick()
    inputRef.value?.focus()
    document.body.classList.add('mw-drawer-open')
  } else {
    document.body.classList.remove('mw-drawer-open')
  }
})

onKeyStroke('Escape', close)

onKeyStroke(['Meta', 'k'], (e) => {
  e.preventDefault()
  isOpen.value ? close() : (isOpen.value = true)
})
onKeyStroke(['Control', 'k'], (e) => {
  e.preventDefault()
  isOpen.value ? close() : (isOpen.value = true)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="search-panel">
      <div v-if="isOpen" class="gx-search-overlay" @click.self="close">
        <!-- Desktop: centered dialog -->
        <div class="gx-search-panel">
          <div class="gx-search-panel__header">
            <GxIcon name="search" :size="18" class="text-gx-muted" />
            <input
              ref="inputRef"
              v-model="keyword"
              type="search"
              placeholder="搜索页面、功能或帖子..."
              class="gx-search-panel__input"
              @keydown.enter="searchPosts"
            />
            <span class="gx-search-panel__hint">ESC 关闭</span>
          </div>
          <div class="gx-search-panel__body">
            <div v-if="filteredCommands.length" class="gx-search-panel__section">
              <p class="gx-search-panel__section-title">
                功能导航
                <span class="gx-search-panel__role-tag">{{ roleLabel }}</span>
              </p>
              <button
                v-for="cmd in filteredCommands"
                :key="cmd.id"
                type="button"
                class="gx-search-panel__item"
                @click="navigate(cmd)"
              >
                <GxIcon :name="cmd.icon" :size="18" />
                <span>{{ cmd.label }}</span>
              </button>
            </div>
            <div v-if="hasKeyword" class="gx-search-panel__section">
              <p class="gx-search-panel__section-title">快速搜索</p>
              <button type="button" class="gx-search-panel__item gx-search-panel__item--search" @click="searchPosts">
                <GxIcon name="search" :size="18" />
                <span>搜索帖子"{{ keyword.trim() }}"</span>
              </button>
            </div>
            <div v-if="!filteredCommands.length && !hasKeyword" class="gx-search-panel__empty">
              没有匹配的功能
            </div>
          </div>
        </div>

        <!-- Mobile: bottom sheet -->
        <div class="gx-search-sheet">
          <div class="gx-search-sheet__handle" />
          <div class="gx-search-sheet__header">
            <GxIcon name="search" :size="18" class="text-gx-muted shrink-0" />
            <input
              ref="inputRef"
              v-model="keyword"
              type="search"
              placeholder="搜索页面、功能或帖子..."
              class="gx-search-sheet__input"
              enterkeyhint="search"
              @keydown.enter="searchPosts"
            />
            <button type="button" class="gx-search-sheet__cancel" @click="close">取消</button>
          </div>
          <div class="gx-search-sheet__body">
            <div v-if="filteredCommands.length" class="gx-search-panel__section">
              <p class="gx-search-panel__section-title">
                功能导航
                <span class="gx-search-panel__role-tag">{{ roleLabel }}</span>
              </p>
              <button
                v-for="cmd in filteredCommands"
                :key="cmd.id"
                type="button"
                class="gx-search-panel__item"
                @click="navigate(cmd)"
              >
                <GxIcon :name="cmd.icon" :size="20" />
                <span>{{ cmd.label }}</span>
              </button>
            </div>
            <div v-if="hasKeyword" class="gx-search-panel__section">
              <p class="gx-search-panel__section-title">快速搜索</p>
              <button type="button" class="gx-search-panel__item gx-search-panel__item--search" @click="searchPosts">
                <GxIcon name="search" :size="20" />
                <span>搜索帖子"{{ keyword.trim() }}"</span>
              </button>
            </div>
            <div v-if="!filteredCommands.length && !hasKeyword" class="gx-search-panel__empty">
              没有匹配的功能
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* ====== Shared ====== */
.gx-search-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.gx-search-panel__section {
  padding: 8px 0;
}
.gx-search-panel__section + .gx-search-panel__section {
  border-top: 1px solid var(--color-border, #e5e7eb);
}

.gx-search-panel__section-title {
  padding: 6px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-muted, #6b7280);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  display: flex;
  align-items: center;
  gap: 8px;
}

.gx-search-panel__role-tag {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 500;
  background: var(--color-primary-bg, rgba(61, 124, 115, 0.1));
  color: var(--color-primary, #3d7c73);
}

.gx-search-panel__item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 16px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--color-primary, #17324f);
  font-size: 14px;
  cursor: pointer;
  transition: background 0.1s ease;
  text-align: left;
}
.gx-search-panel__item:hover,
.gx-search-panel__item:focus-visible {
  background: var(--color-bg, #f5f5f5);
  outline: none;
}
.gx-search-panel__item:active {
  background: var(--color-border, #e5e7eb);
}

.gx-search-panel__item--search {
  color: var(--color-primary, #3d7c73);
  font-weight: 500;
}

.gx-search-panel__empty {
  padding: 24px 16px;
  text-align: center;
  color: var(--color-muted, #9ca3af);
  font-size: 14px;
}

/* ====== Desktop: centered dialog ====== */
.gx-search-panel {
  display: none;
  flex-direction: column;
  width: 480px;
  max-height: 60vh;
  margin: auto;
  background: var(--color-surface, #fff);
  border: 1px solid var(--color-border, #e5e7eb);
  border-radius: 12px;
  box-shadow: 0 16px 48px rgba(15, 43, 91, 0.22);
  overflow: hidden;
}

.gx-search-panel__header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
}

.gx-search-panel__input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 15px;
  color: var(--color-primary, #17324f);
  background: transparent;
}
.gx-search-panel__input::placeholder {
  color: var(--color-muted, #9ca3af);
}

.gx-search-panel__hint {
  font-size: 11px;
  color: var(--color-muted, #9ca3af);
  white-space: nowrap;
  padding: 2px 8px;
  border: 1px solid var(--color-border, #e5e7eb);
  border-radius: 4px;
}

.gx-search-panel__body {
  overflow-y: auto;
  padding: 4px;
}

/* ====== Mobile: bottom sheet ====== */
.gx-search-sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  background: var(--color-surface, #fff);
  border-radius: 20px 20px 0 0;
  box-shadow: 0 -4px 24px rgba(15, 43, 91, 0.16);
  padding-bottom: env(safe-area-inset-bottom, 0px);
}

@media (min-width: 768px) {
  .gx-search-overlay {
    align-items: center;
  }

  .gx-search-panel {
    display: flex;
  }

  .gx-search-sheet {
    display: none;
  }
}

.gx-search-sheet__handle {
  width: 36px;
  height: 5px;
  margin: 10px auto 4px;
  border-radius: 999px;
  background: var(--color-border, #d1d5db);
  flex-shrink: 0;
}

.gx-search-sheet__header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px 12px;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
  flex-shrink: 0;
}

.gx-search-sheet__input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  color: var(--color-primary, #17324f);
  background: transparent;
}
.gx-search-sheet__input::placeholder {
  color: var(--color-muted, #9ca3af);
}

.gx-search-sheet__cancel {
  flex-shrink: 0;
  border: none;
  background: transparent;
  color: var(--color-primary, #3d7c73);
  font-size: 15px;
  font-weight: 500;
  padding: 8px 4px;
  cursor: pointer;
  min-height: var(--mw-tap-min, 44px);
  display: flex;
  align-items: center;
}
.gx-search-sheet__cancel:active {
  opacity: 0.7;
}

.gx-search-sheet__body {
  overflow-y: auto;
  padding: 4px;
  flex: 1;
}

/* ====== Touch ====== */
@media (max-width: 767.98px) {
  .gx-search-overlay {
    background: rgba(0, 0, 0, 0.35);
    -webkit-backdrop-filter: blur(2px);
    backdrop-filter: blur(2px);
  }

  .gx-search-panel__item {
    min-height: var(--mw-tap-min, 44px);
    font-size: 15px;
    padding: 12px 16px;
  }
  .gx-search-panel__item:active {
    transform: scale(0.98);
  }
}

/* ====== Transitions ====== */
.search-panel-enter-active {
  transition: opacity 0.2s ease;
}
.search-panel-enter-active .gx-search-panel {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.search-panel-enter-active .gx-search-sheet {
  transition: transform 0.25s cubic-bezier(0.32, 0.72, 0, 1);
}

.search-panel-leave-active {
  transition: opacity 0.15s ease;
}
.search-panel-leave-active .gx-search-panel {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.search-panel-leave-active .gx-search-sheet {
  transition: transform 0.2s cubic-bezier(0.32, 0.72, 0, 1);
}

.search-panel-enter-from {
  opacity: 0;
}
.search-panel-enter-from .gx-search-panel {
  opacity: 0;
  transform: scale(0.96);
}
.search-panel-enter-from .gx-search-sheet {
  transform: translateY(100%);
}

.search-panel-leave-to {
  opacity: 0;
}
.search-panel-leave-to .gx-search-panel {
  opacity: 0;
  transform: scale(0.96);
}
.search-panel-leave-to .gx-search-sheet {
  transform: translateY(100%);
}
</style>
