<script setup>
import GxFeedSortTabs from './GxFeedSortTabs.vue'
import { FEED_SORT_OPTIONS } from '../../composables/useFeedSort'

const sort = defineModel({ type: String, default: 'hot' })
</script>

<template>
  <div class="gx-feed-sort-bar">
    <GxFeedSortTabs v-model="sort" />
    <label class="gx-feed-sort-bar__extra">
      <span class="gx-feed-sort-bar__extra-label">排序方式：</span>
      <select v-model="sort" class="gx-feed-sort-bar__select" aria-label="排序方式">
        <option v-for="opt in FEED_SORT_OPTIONS" :key="opt.id" :value="opt.id">
          {{ opt.label }}
        </option>
      </select>
    </label>
  </div>
</template>
