<script setup>
import { RouterLink } from 'vue-router'
import Badge from '../ui/Badge.vue'
import Card from '../ui/Card.vue'
import GxAvatarInitial from './GxAvatarInitial.vue'
import { boardTagClass } from '../../composables/useGxNav'
import { formatAuthorLabel } from '../../utils/displayName'

defineProps({
  post: { type: Object, required: true },
  pinned: { type: Boolean, default: false },
  announce: { type: Boolean, default: false },
  variant: { type: String, default: 'list' },
})

function badgeVariant(boardName = '') {
  const cls = boardTagClass(boardName)
  if (cls.includes('club') || cls.includes('notice')) return 'gold'
  return 'secondary'
}

function dateIso(post) {
  return post.createdAtIso || post.createdAt || ''
}
</script>

<template>
  <RouterLink
    v-if="variant === 'list'"
    :to="`/community/posts/${post.id}`"
    class="gx-post-row group block rounded-gx-md border border-gx-border bg-gx-surface px-3 py-3 transition-colors hover:border-gx-primary/25 hover:bg-gx-bg/80"
    :class="{ 'gx-post-row--announce': announce, 'gx-post-row--pinned': pinned && !announce }"
  >
    <div class="gx-post-row__grid">
      <GxAvatarInitial :name="formatAuthorLabel(null, post)" :size="40" class="gx-post-row__avatar" />
      <div class="gx-post-row__main min-w-0">
        <div class="flex flex-wrap items-center gap-2">
          <Badge :variant="badgeVariant(post.boardName)">{{ post.boardName || '讨论' }}</Badge>
          <Badge v-if="pinned" variant="gold">置顶</Badge>
        </div>
        <h3 class="mt-1 text-base font-semibold leading-snug text-gx-primary group-hover:text-gx-accent">
          {{ post.title }}
        </h3>
        <p class="mt-0.5 text-meta text-gx-muted">{{ formatAuthorLabel(null, post) }}</p>
      </div>
      <div class="gx-post-row__meta text-right text-meta text-gx-muted">
        <time class="block shrink-0 whitespace-nowrap tabular-nums" :datetime="dateIso(post)">
          {{ post.createdAt }}
        </time>
        <span class="mt-1 block whitespace-nowrap">{{ post.commentCount }} 回复 · {{ post.likeCount }} 赞</span>
      </div>
    </div>
  </RouterLink>

  <Card
    v-else
    class="gx-post-card group overflow-hidden transition-shadow hover:shadow-[0_4px_16px_rgba(15,43,91,0.12)]"
    :class="{ 'gx-post-card--pinned': pinned, 'gx-post-card--announce': announce }"
  >
    <Badge v-if="pinned" variant="gold" class="absolute right-3 top-3 z-10">置顶</Badge>
    <RouterLink :to="`/community/posts/${post.id}`" class="block p-4 pt-5">
      <header class="mb-2 flex items-center justify-between gap-2 text-meta text-gx-muted">
        <Badge :variant="badgeVariant(post.boardName)">{{ post.boardName || '讨论' }}</Badge>
        <time class="shrink-0 whitespace-nowrap tabular-nums" :datetime="dateIso(post)">
          {{ post.createdAt }}
        </time>
      </header>
      <h3 class="mb-2 text-base font-semibold leading-snug text-gx-primary group-hover:text-gx-accent">
        {{ post.title }}
      </h3>
      <p class="mb-3 line-clamp-2 text-body text-gx-muted">{{ post.content }}</p>
      <footer class="flex flex-wrap items-center justify-between gap-2 border-t border-gx-border pt-3 text-meta text-gx-muted">
        <span class="font-medium text-gx-primary">{{ formatAuthorLabel(null, post) }}</span>
        <span>{{ post.commentCount }} 回复 · {{ post.likeCount }} 赞</span>
      </footer>
    </RouterLink>
  </Card>
</template>

<style scoped>
.gx-post-row__grid {
  display: grid;
  grid-template-columns: 40px 1fr auto;
  gap: 0.75rem 1rem;
  align-items: start;
}
.gx-post-row--pinned {
  border-left: 3px solid var(--color-gold);
}
.gx-post-row--announce {
  border-left: 3px solid var(--color-gold);
}
@media (max-width: 480px) {
  .gx-post-row__grid {
    grid-template-columns: 36px 1fr;
  }
  .gx-post-row__meta {
    grid-column: 2;
    text-align: left;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem 1rem;
  }
}
.gx-post-card {
  position: relative;
}
.gx-post-card--pinned {
  border-left: 3px solid var(--color-gold);
}
.gx-post-card--announce {
  border-left: 3px solid var(--color-gold);
}
</style>
