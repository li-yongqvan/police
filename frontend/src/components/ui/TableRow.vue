<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <tr :class="cn('border-b border-gx-border transition-colors hover:bg-gx-bg/50 data-[state=selected]:bg-gx-bg', $attrs.class)">
    <slot />
  </tr>
</template>
