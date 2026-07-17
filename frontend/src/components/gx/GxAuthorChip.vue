<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import GxAvatarInitial from './GxAvatarInitial.vue'
import { formatAuthorLabel, formatDisplayTime } from '../../utils/displayName'

const props = defineProps({
  authorId: { type: [String, Number], default: '' },
  authorName: { type: String, default: '' },
  authorAvatar: { type: String, default: '' },
  user: { type: Object, default: null },
  post: { type: Object, default: null },
  createdAt: { type: String, default: '' },
  size: { type: String, default: 'md' },
  linkable: { type: Boolean, default: true },
})

const label = computed(() => formatAuthorLabel(props.user, props.post || { authorName: props.authorName }))
const avatarSrc = computed(() => props.authorAvatar || props.user?.avatar || props.post?.authorAvatar || props.post?.avatar || '')
const avatarSize = computed(() => {
  if (props.size === 'lg') return 64
  if (props.size === 'sm') return 32
  return 40
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
      <GxAvatarInitial :name="label" :src="avatarSrc" :size="avatarSize" />
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
