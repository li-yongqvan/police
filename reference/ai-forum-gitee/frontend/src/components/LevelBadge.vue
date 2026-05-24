<template>
  <span class="level-badge" :class="levelClass">{{ levelText }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  level: { type: Number, default: 0 },
  size: { type: String, default: 'sm' },
})

const levelText = computed(() => {
  switch (props.level) {
    case 0: return '游客'
    case 1: return '初级'
    case 2: return '中级'
    case 3: return '高级'
    default: return '资深'
  }
})

const levelClass = computed(() => {
  const classes = ['level-badge']
  if (props.size) classes.push(`size-${props.size}`)
  if (props.level === 0) classes.push('level-gray')
  else if (props.level === 1) classes.push('level-blue')
  else if (props.level === 2) classes.push('level-green')
  else if (props.level === 3) classes.push('level-purple')
  else classes.push('level-gold')
  return classes.join(' ')
})
</script>

<style scoped>
.level-badge {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: white;
}

.level-gray { background: #95a5a6; }
.level-blue { background: #3498db; }
.level-green { background: #2ecc71; }
.level-purple { background: #9b59b6; }
.level-gold { background: #f39c12; }

.size-sm { font-size: 0.7rem; padding: 0.1rem 0.4rem; }
.size-lg { font-size: 0.875rem; padding: 0.25rem 0.75rem; }
</style>
