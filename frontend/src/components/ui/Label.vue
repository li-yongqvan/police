<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
  for: { type: String, default: undefined },
})
</script>

<template>
  <label :for="for" :class="cn('text-sm font-medium leading-none text-[#222]', $attrs.class)">
    <slot />
  </label>
</template>
