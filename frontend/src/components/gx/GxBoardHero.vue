<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import GxIcon from './GxIcon.vue'
import { GX_NAV_ITEMS } from '../../composables/useGxNav'

const props = defineProps({
  board: { type: Object, default: null },
  navLabel: { type: String, default: '' },
})

const iconName = computed(() => {
  const key = props.board?.slug || ''
  const nav = GX_NAV_ITEMS.find((n) => n.key === key)
  return nav?.icon || 'book'
})

const followLabel = computed(() => {
  const n = props.board?.postCount ?? 0
  if (n >= 10000) return `${(n / 10000).toFixed(1)}w`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return String(n)
})
</script>

<template>
  <header v-if="board" class="gx-board-hero">
    <div class="gx-board-hero__icon" aria-hidden="true">
      <GxIcon :name="iconName" :size="28" />
    </div>
    <div class="gx-board-hero__body">
      <p class="gx-board-hero__eyebrow">{{ navLabel || '板块' }}</p>
      <h1 class="gx-board-hero__title">{{ board.name }}</h1>
      <p class="gx-board-hero__desc">{{ board.description }}</p>
      <p class="gx-board-hero__stats text-meta">
        <span>关注 {{ followLabel }}</span>
        <span aria-hidden="true">·</span>
        <span>帖子 {{ board.postCount ?? 0 }}</span>
      </p>
    </div>
    <RouterLink
      :to="{ path: '/community/posts/new', query: { boardId: board.id } }"
      class="gx-board-hero__cta gx-board-hero__btn"
    >
      <GxIcon name="edit" :size="18" />
      发布到本板块
    </RouterLink>
    <div class="gx-board-hero__watermark" aria-hidden="true" />
  </header>
</template>
