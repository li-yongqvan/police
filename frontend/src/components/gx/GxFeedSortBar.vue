<script setup>
import GxFeedSortTabs from './GxFeedSortTabs.vue'
import { FEED_SORT_OPTIONS } from '../../composables/useFeedSort'

const sort = defineModel({ type: String, default: 'hot' })

const sortLabel = (id) => FEED_SORT_OPTIONS.find((o) => o.id === id)?.label || '热门'
</script>

<template>
  <div class="gx-feed-sort-bar">
    <GxFeedSortTabs v-model="sort" />
    <div class="gx-feed-sort-bar__extra">
      <span class="text-meta text-gx-muted">排序方式</span>
      <span class="gx-feed-sort-bar__current">{{ sortLabel(sort) }}</span>
      <span class="gx-feed-sort-bar__caret" aria-hidden="true">▾</span>
    </div>
  </div>
</template>
