<script setup>
import { cn } from '@/lib/utils'

defineProps({
  as: { type: String, default: 'div' },
  center: { type: Boolean, default: false },
  class: { type: null, default: undefined },
})
</script>

<template>
  <component
    :is="as"
    :class="cn('gx-reading-column', center && 'gx-reading-column--center', $props.class)"
  >
    <slot />
  </component>
</template>
