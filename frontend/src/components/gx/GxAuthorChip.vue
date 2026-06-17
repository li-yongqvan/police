<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { formatAuthorLabel, formatDisplayTime } from '../../utils/displayName'

const props = defineProps({
  authorId: { type: [String, Number], default: '' },
  authorName: { type: String, default: '' },
  user: { type: Object, default: null },
  post: { type: Object, default: null },
  createdAt: { type: String, default: '' },
  size: { type: String, default: 'md' },
  linkable: { type: Boolean, default: true },
})

const label = computed(() => formatAuthorLabel(props.user, props.post || { authorName: props.authorName }))
const initial = computed(() => {
  const raw = props.user?.name || props.user?.username || props.authorName || props.post?.authorName || '?'
  return String(raw).trim()[0]?.toUpperCase() || '?'
})
const profileTo = computed(() => (props.authorId ? `/community/users/${props.authorId}` : ''))
const timeLabel = computed(() => {
  if (!props.createdAt) return ''
  const raw = String(props.createdAt)
  if (raw.includes('T')) return formatDisplayTime(raw)
  return raw
})
</script>

<template>
  <div class="gx-author-chip" :class="`gx-author-chip--${size}`">
    <component
      :is="linkable && profileTo ? RouterLink : 'div'"
      :to="linkable && profileTo ? profileTo : undefined"
      class="gx-author-chip__avatar"
      :aria-label="linkable && profileTo ? `${label} 的主页` : undefined"
    >
      {{ initial }}
    </component>
    <div class="gx-author-chip__meta">
      <component
        :is="linkable && profileTo ? RouterLink : 'span'"
        :to="linkable && profileTo ? profileTo : undefined"
        class="gx-author-chip__name"
      >
        {{ label }}
      </component>
      <time v-if="timeLabel" class="gx-author-chip__time">{{ timeLabel }}</time>
    </div>
  </div>
</template>
