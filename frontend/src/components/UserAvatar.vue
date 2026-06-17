<template>
  <div class="user-avatar" :style="avatarStyle">
    <img v-if="avatarUrl" :src="avatarUrl" :alt="displayName" />
    <span v-else class="avatar-initials">{{ initials }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  user: { type: Object, default: () => ({}) },
  size: { type: String, default: 'md' },
})

const sizeMap = { sm: 32, md: 48, lg: 96 }
const sizePx = sizeMap[props.size] || 48

const avatarUrl = computed(() => props.user?.avatar || null)
const displayName = computed(() => props.user?.nickname || props.user?.username || '?')
const initials = computed(() => {
  const name = displayName.value
  return name ? name.substring(0, 2).toUpperCase() : '?'
})

const colors = ['#4a90d9', '#2ecc71', '#e74c3c', '#9b59b6', '#f39c12', '#1abc9c']
const colorIndex = computed(() => {
  let hash = 0
  for (let i = 0; i < displayName.value.length; i++) {
    hash = displayName.value.charCodeAt(i) + ((hash << 5) - hash)
  }
  return Math.abs(hash) % colors.length
})

const avatarStyle = computed(() => ({
  width: sizePx + 'px',
  height: sizePx + 'px',
  borderRadius: '50%',
  overflow: 'hidden',
  backgroundColor: colors[colorIndex.value],
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}))
</script>

<style scoped>
.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-initials {
  color: white;
  font-weight: bold;
  font-size: inherit;
}
</style>
