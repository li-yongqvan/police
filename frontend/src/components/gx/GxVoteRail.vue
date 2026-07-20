<script setup>
import { computed } from 'vue'
import GxIcon from './GxIcon.vue'

const props = defineProps({
  score: { type: Number, default: 0 },
  liked: { type: Boolean, default: false },
  disliked: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  compact: { type: Boolean, default: false },
  vertical: { type: Boolean, default: true },
  trendIcon: { type: String, default: '' },
  trendTone: { type: String, default: 'stable' },
  trendLabel: { type: String, default: '' },
  trendOnly: { type: Boolean, default: false },
})

const emit = defineEmits(['vote', 'dislike'])

const label = computed(() => {
  const n = props.score
  if (n >= 10000) return `${(n / 10000).toFixed(1)}万`
  return String(n)
})

function onVote() {
  if (!props.loading) emit('vote')
}

function onDislike() {
  if (!props.loading) emit('dislike')
}
</script>

<template>
  <div
    class="gx-vote-rail"
    :class="{
      'gx-vote-rail--compact': compact,
      'gx-vote-rail--trend-only': trendOnly,
      'gx-vote-rail--liked': liked,
      'gx-vote-rail--disliked': disliked,
    }"
    @click.stop
  >
    <button
      v-if="!trendOnly"
      type="button"
      class="gx-vote-rail__btn"
      :class="{ 'is-active': liked }"
      :disabled="loading"
      aria-label="点赞"
      @click="onVote"
    >
      <svg viewBox="0 0 20 20" width="18" height="18" aria-hidden="true">
        <path
          fill="currentColor"
          d="M10 3.5 14.5 9H12v7H8V9H5.5L10 3.5z"
        />
      </svg>
    </button>
    <span
      v-if="trendIcon"
      class="gx-vote-rail__trend"
      :class="`is-${trendTone}`"
      :title="trendLabel"
      :aria-label="trendLabel"
    >
      <GxIcon :name="trendIcon" :size="16" />
    </span>
    <span v-else class="gx-vote-rail__score tabular-nums" :class="{ 'is-active': liked }">
      {{ label }}
    </span>
    <button
      v-if="!trendOnly"
      type="button"
      class="gx-vote-rail__btn gx-vote-rail__btn--down"
      :class="{ 'is-active': disliked }"
      :disabled="loading"
      aria-label="点踩"
      @click="onDislike"
    >
      <svg viewBox="0 0 20 20" width="18" height="18" aria-hidden="true">
        <path
          fill="currentColor"
          d="M10 16.5 5.5 11H8V4h4v7h2.5L10 16.5z"
        />
      </svg>
    </button>
  </div>
</template>
