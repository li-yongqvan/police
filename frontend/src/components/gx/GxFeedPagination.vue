<script setup>
import { computed } from 'vue'

const props = defineProps({
  page: { type: Number, default: 1 },
  total: { type: Number, default: 0 },
  limit: { type: Number, default: 20 },
})

const emit = defineEmits(['update:page'])

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit) || 1))

const pages = computed(() => {
  const tp = totalPages.value
  const p = props.page
  if (tp <= 7) return Array.from({ length: tp }, (_, i) => i + 1)
  const items = [1]
  if (p > 3) items.push('…')
  const mid = [p - 1, p, p + 1].filter((n) => n > 1 && n < tp)
  items.push(...mid)
  if (p < tp - 2) items.push('…')
  if (tp > 1) items.push(tp)
  return [...new Set(items)]
})

function go(next) {
  const n = Math.min(totalPages.value, Math.max(1, next))
  if (n !== props.page) emit('update:page', n)
}
</script>

<template>
  <nav v-if="totalPages > 1" class="gx-feed-pagination" aria-label="分页">
    <button type="button" class="gx-feed-pagination__btn" :disabled="page <= 1" @click="go(page - 1)">
      ‹
    </button>
    <template v-for="(item, idx) in pages" :key="`${item}-${idx}`">
      <span v-if="item === '…'" class="gx-feed-pagination__ellipsis">…</span>
      <button
        v-else
        type="button"
        class="gx-feed-pagination__btn"
        :class="{ 'is-active': item === page }"
        @click="go(item)"
      >
        {{ item }}
      </button>
    </template>
    <button
      type="button"
      class="gx-feed-pagination__btn"
      :disabled="page >= totalPages"
      @click="go(page + 1)"
    >
      ›
    </button>
  </nav>
</template>
