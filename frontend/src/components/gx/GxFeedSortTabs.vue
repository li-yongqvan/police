<script setup>
import { FEED_SORT_OPTIONS } from '../../composables/useFeedSort'

const sort = defineModel({ type: String, default: 'hot' })
</script>

<template>
  <div class="gx-feed-sort" role="tablist" aria-label="帖子排序">
    <button
      v-for="opt in FEED_SORT_OPTIONS"
      :key="opt.id"
      type="button"
      role="tab"
      class="gx-feed-sort__tab"
      :class="{ 'is-active': sort === opt.id }"
      :aria-selected="sort === opt.id"
      @click="sort = opt.id"
    >
      {{ opt.label }}
    </button>
  </div>
</template>
