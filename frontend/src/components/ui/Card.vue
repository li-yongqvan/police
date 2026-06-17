<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <div :class="cn('rounded-gx-md border border-gx-border bg-gx-surface text-[#222] shadow-[0_2px_8px_rgba(15,43,91,0.08)]', $attrs.class)">
    <slot />
  </div>
</template>
