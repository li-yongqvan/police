<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import Card from '../ui/Card.vue'
import GxIcon from './GxIcon.vue'
import { GX_NAV_ITEMS, resolveBoardByKey } from '../../composables/useGxNav'

const props = defineProps({
  boards: { type: Array, default: () => [] },
})

const items = computed(() =>
  GX_NAV_ITEMS.map((nav) => ({
    ...nav,
    board: resolveBoardByKey(props.boards, nav.key),
    to: `/community/boards/${nav.key}`,
  })),
)
</script>

<template>
  <section class="grid grid-cols-2 gap-3 md:grid-cols-4" aria-label="板块快捷入口">
    <RouterLink v-for="item in items" :key="item.key" :to="item.to" class="group">
      <Card
        class="flex h-full flex-col items-center gap-2 p-4 text-center transition-all hover:-translate-y-0.5 hover:border-gx-primary/30 hover:shadow-md"
      >
        <span class="flex h-12 w-12 items-center justify-center rounded-full bg-gx-primary/8 text-gx-primary transition-colors group-hover:bg-gx-primary group-hover:text-white">
          <GxIcon :name="item.icon" :size="28" />
        </span>
        <strong class="text-sm text-gx-primary">{{ item.label }}</strong>
        <span class="text-xs text-gx-muted">{{ item.board?.name || '进入板块' }}</span>
      </Card>
    </RouterLink>
  </section>
</template>
