<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import GxIcon from './GxIcon.vue'
import { formatUserChipLabel } from '../../utils/displayName'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  stats: { type: Object, default: null },
})

const session = useSessionStore()

const greeting = computed(() => {
  const who = session.currentUser ? formatUserChipLabel(session.currentUser) : '同学'
  const hour = new Date().getHours()
  if (hour < 12) return `早上好，${who}`
  if (hour < 18) return `下午好，${who}`
  return `晚上好，${who}`
})

const statLine = computed(() => {
  const s = props.stats
  if (!s) return '浏览全校最新讨论'
  return `全校 ${s.total_posts ?? 0} 帖 · 今日 ${s.posts_today ?? 0} 帖 · ${s.online_users ?? 0} 人在线`
})
</script>

<template>
  <header class="gx-home-hero">
    <div class="gx-home-hero__body">
      <p class="gx-home-hero__eyebrow">全校社区</p>
      <h1 class="gx-home-hero__title">{{ greeting }}</h1>
      <p class="gx-home-hero__desc">{{ statLine }}</p>
    </div>
    <RouterLink to="/community/posts/new" class="gx-home-hero__cta">
      <GxIcon name="edit" :size="18" />
      发布帖子
    </RouterLink>
  </header>
</template>
