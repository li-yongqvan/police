<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
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
  <nav class="gx-feed-chips" aria-label="板块快捷入口">
    <RouterLink v-for="item in items" :key="item.key" :to="item.to" class="gx-feed-chips__item">
      {{ item.board?.name || item.label }}
    </RouterLink>
  </nav>
</template>
