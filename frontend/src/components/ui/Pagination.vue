<script setup>
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import Button from './Button.vue'

const props = defineProps({
  currentPage: { type: Number, required: true },
  totalPages: { type: Number, required: true },
  class: { type: null, default: undefined },
})

const emit = defineEmits(['update:currentPage'])

const pages = computed(() => {
  const p = []
  const total = props.totalPages
  const current = props.currentPage
  if (total <= 7) {
    for (let i = 1; i <= total; i++) p.push(i)
  } else {
    p.push(1)
    if (current > 3) p.push('...')
    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)
    for (let i = start; i <= end; i++) p.push(i)
    if (current < total - 2) p.push('...')
    p.push(total)
  }
  return p
})

function go(page) {
  if (page === '...' || page < 1 || page > props.totalPages) return
  emit('update:currentPage', page)
}
</script>

<template>
  <nav :class="cn('flex items-center gap-1', props.class)" role="navigation" aria-label="分页">
    <Button variant="ghost" size="icon" :disabled="currentPage <= 1" @click="go(currentPage - 1)">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
    </Button>
    <Button
      v-for="page in pages"
      :key="page"
      :variant="page === currentPage ? 'default' : 'ghost'"
      size="icon"
      :disabled="page === '...'"
      @click="go(page)"
    >
      {{ page }}
    </Button>
    <Button variant="ghost" size="icon" :disabled="currentPage >= totalPages" @click="go(currentPage + 1)">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
    </Button>
  </nav>
</template>
