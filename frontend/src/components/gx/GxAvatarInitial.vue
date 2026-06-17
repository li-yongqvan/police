<script setup>
import { computed } from 'vue'
import { maskStudentLabel } from '../../utils/displayName'

const props = defineProps({
  name: { type: String, default: '' },
  size: { type: Number, default: 40 },
})

const label = computed(() => maskStudentLabel(props.name))
const initial = computed(() => {
  const t = label.value.replace(/[·*]/g, '').trim()
  return t ? t.charAt(0).toUpperCase() : '?'
})
</script>

<template>
  <span
    class="gx-avatar-initial inline-flex shrink-0 items-center justify-center rounded-full bg-gx-primary/10 font-semibold text-gx-primary"
    :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${Math.round(size * 0.4)}px` }"
    :title="label"
    aria-hidden="true"
  >
    {{ initial }}
  </span>
</template>
