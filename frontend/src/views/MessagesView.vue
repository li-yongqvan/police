<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxIcon from '../components/gx/GxIcon.vue'
import { forumApi } from '../api'

const tab = ref('all')
const items = ref([])
const activeId = ref('')

const tabDefs = [
  { id: 'all', label: '全部消息', icon: 'message', desc: '所有通知' },
  { id: 'reply', label: '回复我的', icon: 'message', desc: '帖子回复提醒' },
  { id: 'system', label: '系统公告', icon: 'bell', desc: '平台与审核通知' },
]

const filteredItems = computed(() => {
  if (tab.value === 'reply') return items.value.filter((i) => i.type === 'reply')
  if (tab.value === 'system') return items.value.filter((i) => i.type === 'system')
  return items.value
})

const activeTab = computed(() => tabDefs.find((t) => t.id === tab.value) || tabDefs[0])

const unreadTotal = computed(() => items.value.filter((i) => !i.is_read).length)

function tabUnread(id) {
  const pool =
    id === 'reply'
      ? items.value.filter((i) => i.type === 'reply')
      : id === 'system'
        ? items.value.filter((i) => i.type === 'system')
        : items.value
  return pool.filter((i) => !i.is_read).length
}

function tabCount(id) {
  if (id === 'reply') return items.value.filter((i) => i.type === 'reply').length
  if (id === 'system') return items.value.filter((i) => i.type === 'system').length
  return items.value.length
}

function typeIcon(type) {
  return type === 'system' ? 'bell' : 'message'
}

function typeLabel(type) {
  return type === 'system' ? '系统' : '回复'
}

async function load() {
  const data = await forumApi.listNotifications()
  items.value = data.items
  if (!activeId.value && data.items.length) {
    activeId.value = data.items[0].id
  }
}

async function markRead(item) {
  if (!item.is_read) {
    await forumApi.markNotificationRead(item.id)
    item.is_read = true
  }
}

function selectItem(item) {
  activeId.value = item.id
}

function switchTab(id) {
  tab.value = id
  const first = filteredItems.value[0]
  activeId.value = first?.id || ''
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-messages-page">
    <div class="gx-messages-shell">
      <aside class="gx-messages-rail" aria-label="消息分类">
        <header class="gx-messages-rail__head">
          <div class="gx-messages-rail__brand">
            <span class="gx-messages-rail__logo" aria-hidden="true">
              <GxIcon name="message" :size="20" />
            </span>
            <div>
              <h1 class="gx-messages-rail__title">消息中心</h1>
              <p v-if="unreadTotal" class="gx-messages-rail__meta">{{ unreadTotal }} 条未读</p>
              <p v-else class="gx-messages-rail__meta">已全部读完</p>
            </div>
          </div>
        </header>

        <nav class="gx-messages-rail__nav">
          <p class="gx-messages-rail__section">收件箱</p>
          <button
            v-for="t in tabDefs"
            :key="t.id"
            type="button"
            class="gx-messages-channel"
            :class="{ 'is-active': tab === t.id }"
            @click="switchTab(t.id)"
          >
            <span class="gx-messages-channel__icon" aria-hidden="true">
              <GxIcon :name="t.icon" :size="18" />
            </span>
            <span class="gx-messages-channel__text">
              <span class="gx-messages-channel__label">{{ t.label }}</span>
              <span class="gx-messages-channel__desc">{{ t.desc }}</span>
            </span>
            <span v-if="tabUnread(t.id)" class="gx-messages-channel__unread">{{ tabUnread(t.id) }}</span>
            <span v-else-if="tabCount(t.id)" class="gx-messages-channel__count">{{ tabCount(t.id) }}</span>
          </button>
        </nav>
      </aside>

      <section class="gx-messages-panel">
        <header class="gx-messages-panel__head">
          <div class="gx-messages-panel__head-main">
            <span class="gx-messages-panel__head-icon" aria-hidden="true">
              <GxIcon :name="activeTab.icon" :size="20" />
            </span>
            <div>
              <h2 class="gx-messages-panel__title">{{ activeTab.label }}</h2>
              <p class="gx-messages-panel__subtitle">
                {{ filteredItems.length }} 条消息
                <template v-if="tabUnread(tab)"> · {{ tabUnread(tab) }} 未读</template>
              </p>
            </div>
          </div>
        </header>

        <div class="gx-messages-panel__body">
          <GxEmptyState
            v-if="!filteredItems.length"
            title="暂无消息"
            description="回复与系统通知会显示在这里"
          />

          <ul v-else class="gx-message-list" role="list">
            <li
              v-for="item in filteredItems"
              :key="item.id"
              class="gx-message-row"
              :class="{
                'is-unread': !item.is_read,
                'is-active': activeId === item.id,
              }"
            >
              <button type="button" class="gx-message-row__main" @click="selectItem(item)">
                <span
                  class="gx-message-row__avatar"
                  :class="item.type === 'system' ? 'gx-message-row__avatar--system' : 'gx-message-row__avatar--reply'"
                  aria-hidden="true"
                >
                  <GxIcon :name="typeIcon(item.type)" :size="18" />
                </span>
                <span class="gx-message-row__content">
                  <span class="gx-message-row__top">
                    <span class="gx-message-row__title-wrap">
                      <strong class="gx-message-row__title">{{ item.title }}</strong>
                      <span class="gx-message-row__tag">{{ typeLabel(item.type) }}</span>
                    </span>
                    <time class="gx-message-row__time" :datetime="item.createdAtIso">{{ item.created_at }}</time>
                  </span>
                  <span class="gx-message-row__preview">{{ item.body }}</span>
                </span>
                <span v-if="!item.is_read" class="gx-message-row__dot" aria-label="未读" />
              </button>

              <div v-if="activeId === item.id" class="gx-message-row__actions">
                <button v-if="!item.is_read" type="button" class="gx-message-row__action" @click="markRead(item)">
                  标为已读
                </button>
                <RouterLink
                  v-if="item.related_post_id"
                  :to="`/community/posts/${item.related_post_id}`"
                  class="gx-message-row__action gx-message-row__action--primary"
                >
                  查看帖子
                </RouterLink>
              </div>
            </li>
          </ul>
        </div>
      </section>
    </div>
  </div>
</template>
